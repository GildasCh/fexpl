package main

import (
	"fmt"
	"os"
	"text/tabwriter"
)

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
