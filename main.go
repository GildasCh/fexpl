package main

import (
	"fmt"
	"html/template"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"text/tabwriter"

	"github.com/gorilla/mux"
)

var fileCollections map[string]*FileCollection

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

	fc := Explore(os.Args[2], os.Args[3])

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

	fc, err := ImportFromJSON(os.Args[2])
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
	fileCollections = make(map[string]*FileCollection)

	for _, root := range os.Args[2:] {
		fc := Explore("", root)
		if _, ok := fileCollections[fc.Name]; ok {
			fmt.Printf("Warning: duplicate collection name %q\n", fc.Name)
			continue
		}
		if fc != nil {
			fileCollections[fc.Name] = fc
		}
	}

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
	"size": func(size int64) string {
		units := []string{"", "K", "M", "G", "T"}
		unit := 0
		fsize := float64(size)
		for fsize > 1024 {
			fsize /= 1024
			unit++
			if unit == len(units)-1 {
				break
			}
		}
		return fmt.Sprintf("%.2f%s", fsize, units[unit])
	}}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("indexHandler")

	t, err := template.New("index.html").Funcs(funcMap).ParseFiles("html/index.html")
	if err != nil {
		fmt.Println(err)
	}
	err = t.Execute(w, fileCollections)
	if err != nil {
		fmt.Println(err)
	}
}

func pathHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("pathHandler", mux.Vars(r))

	collection := mux.Vars(r)["collection"]
	fc, ok := fileCollections[collection]
	if !ok {
		fmt.Printf("%q not found\n", collection)
		w.WriteHeader(http.StatusNotFound)
	}

	path := mux.Vars(r)["path"]
	if path != "" {
		fc, ok = extractPath(fc, collection, path)
		if !ok {
			w.WriteHeader(http.StatusNotFound)
		}
	}

	t, err := template.New("index.html").Funcs(funcMap).ParseFiles("html/index.html")
	if err != nil {
		fmt.Println(err)
	}
	err = t.Execute(w, map[string]*FileCollection{collection: fc})
	if err != nil {
		fmt.Println(err)
	}
}

func extractPath(in *FileCollection, collection, path string) (ret *FileCollection, ok bool) {
	ret = &FileCollection{}
	ret.Name = in.Name
	for _, inf := range in.Files {
		if !strings.HasPrefix(inf.Name, filepath.Join(collection, path)) {
			continue
		}
		ok = true
		ret.Files = append(ret.Files, inf)
	}

	ok = false
	return
}
