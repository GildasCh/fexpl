package main

import "fmt"

type Data struct {
	Collections map[string]*Collection // collection name to collection
	Files       map[string][]*File     // hash to file
}

func DataFromJSON(paths []string) *Data {
	ret := &Data{}
	ret.Collections = make(map[string]*Collection)
	ret.Files = make(map[string][]*File)
	for _, p := range paths {
		c, err := importr.Import(p)
		if err != nil {
			fmt.Printf("Error loading data from %q\n", p)
			continue
		}
		if _, ok := ret.Collections[c.Name]; ok {
			fmt.Printf("Could not import %q, collection %q already loaded\n", p, c.Name)
			continue
		}
		ret.Collections[c.Name] = c

		for _, f := range c.Files {
			ret.Files[f.Hash] = append(ret.Files[f.Hash], f)
		}

		fmt.Printf("Imported %d files from %q\n", len(c.Files), p)
	}

	return ret
}
