package main

import (
	"fmt"
	"html/template"
	"net/http"
	"os"
	"text/tabwriter"

	"github.com/gorilla/mux"
)

var fileCollections []*FileCollection

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
		kind := f.Header.MIME.Value
		if kind == "" {
			kind = "unknown"
		}
		fmt.Fprintf(w, "%s\t%s\n", f.Name, kind)
	}
	w.Flush()
}

func serve() {
	for _, root := range os.Args[2:] {
		fc := Explore("", root)
		if fc != nil {
			fileCollections = append(fileCollections, fc)
		}
	}

	r := mux.NewRouter()
	r.HandleFunc("/", indexHandler).Methods("GET")
	r.PathPrefix("/").Handler(http.FileServer(http.Dir("./static/")))

	http.Handle("/", r)
	port := "8080"
	fmt.Printf("Listening on %s...\n", port)
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		fmt.Println(err)
	}
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	t, err := template.New("index.html").Funcs(template.FuncMap{
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
		}}).ParseFiles("html/index.html")
	if err != nil {
		fmt.Println(err)
	}
	err = t.Execute(w, fileCollections)
	if err != nil {
		fmt.Println(err)
	}
}
