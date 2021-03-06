package main

import (
	"fmt"
	"html/template"
	"net/http"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/gorilla/mux"
)

var data *Data
var importr Importer = &JSONImporter{}

func main() {
	if len(os.Args) < 2 {
		usage()
		return
	}

	switch os.Args[1] {
	case "scan":
		scan()
	case "ls":
		ls()
	case "serve":
		serve()
	default:
		usage()
	}

}

func usage() {
	fmt.Println("Bad command")
}

func scan() {
	if len(os.Args) < 5 {
		usage()
		return
	}

	fc := Explore(os.Args[2], os.Args[3], false, 50*1024*1024)

	err := fc.ExportToJSON(os.Args[4])
	if err != nil {
		fmt.Println(err)
	}

	fmt.Printf("Successfully exported collection of %d file to %s\n", len(fc.Files), os.Args[4])
}

func ls() {
	if len(os.Args) < 3 {
		usage()
		return
	}

	fc, err := importr.Import(os.Args[2])
	if err != nil {
		fmt.Println(err)
	}

	w := tabwriter.NewWriter(os.Stdout, 1, 1, 0, ' ', 0)
	for _, f := range fc.Files {
		kind := f.MIME
		if kind == "" {
			kind = "unknown"
		}
		fmt.Fprintf(w, "%s\t%s\n", f.Name, kind)
	}
	w.Flush()
}

func serve() {
	data = DataFromJSON(os.Args[2:])

	fmt.Println(data.Files["3KrlE+DuPPWnw8sDmRHzbP4gvJ+XqDUy3eiT8eKMFLyaclIkgDUYa6kZK+tS9KExzBRZ7eGmOl/dF5uA4ml6lg=="])

	r := mux.NewRouter()
	r.HandleFunc("/", indexHandler).Methods("GET")
	r.HandleFunc("/browse/{collection:[a-zA-Z0-9 _-]+}/{path:.*}", pathHandler).Methods("GET")
	r.HandleFunc("/browse/{collection:[a-zA-Z0-9 _-]+}", pathHandler).Methods("GET")
	r.PathPrefix("/").Handler(http.FileServer(http.Dir("./static/")))

	http.Handle("/", r)
	port := "8080"
	fmt.Printf("Listening on %s...\n", port)
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		fmt.Println(err)
	}
}

var funcMap = template.FuncMap{
	"duplicates": func(f *File) string {
		matching := f.Duplicates()
		if len(matching) <= 0 {
			return ""
		}
		where := ""
		for _, m := range matching {
			where += m.collection.Name + ", "
		}
		where = strings.Trim(where, ", ")
		return fmt.Sprintf("%d (%s)", len(matching), where)
	}}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("indexHandler")

	t, err := template.New("index.html").Funcs(funcMap).ParseFiles("html/index.html")
	if err != nil {
		fmt.Println(err)
	}
	err = t.Execute(w, data)
	if err != nil {
		fmt.Println(err)
	}
}

func pathHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("pathHandler", mux.Vars(r))

	collection := mux.Vars(r)["collection"]
	fc, ok := data.Collections[collection]
	if !ok {
		fmt.Printf("%q not found, has %v\n", collection, data.Collections)
		w.WriteHeader(http.StatusNotFound)
		return
	}

	path := mux.Vars(r)["path"]
	if path != "" {
		fc, ok = extractPath(fc, path)
		if !ok {
			w.WriteHeader(http.StatusNotFound)
			return
		}
	}

	t, err := template.New("index.html").Funcs(funcMap).ParseFiles("html/index.html")
	if err != nil {
		fmt.Println(err)
	}
	err = t.Execute(w, Data{Collections: map[string]*Collection{fc.Name: fc}, Files: data.Files})
	if err != nil {
		fmt.Println(err)
	}
}

func extractPath(in *Collection, path string) (ret *Collection, ok bool) {
	ret = &Collection{}
	ret.Name = in.Name
	ok = false
	for _, inf := range in.Files {
		if !strings.HasPrefix(strings.TrimLeft(inf.Name, "/"), path) {
			fmt.Printf("Filtering out %q (path %q)\n", strings.TrimLeft(inf.Name, "/"), path)
			continue
		}
		ok = true
		ret.Files = append(ret.Files, inf)
	}

	fmt.Println("Output dataset is empty")

	return
}
