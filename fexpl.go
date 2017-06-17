package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	filetype "gopkg.in/h2non/filetype.v1"

	"github.com/GildasCh/fcomp"
)

type File struct {
	Name       string
	Hash       string
	Size       int64
	Modified   time.Time
	MIME       string
	collection *Collection
}

type Collection struct {
	Files []*File
	Name  string
}

type Data struct {
	Collections map[string]*Collection // collection name to collection
	Files       map[string][]*File     // hash to file
}

func Explore(name, root string) *Collection {
	if _, err := os.Stat(filepath.Join(root, "fexpl.json")); err == nil {
		fmt.Printf("File exists on root %q\n", root)
		ret, err := ImportFromJSON(filepath.Join(root, "fexpl.json"))
		if err == nil {
			return ret
		}
		fmt.Println(err)
	}

	if name == "" {
		return nil
	}

	ret := &Collection{}
	ret.Name = name
	filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		f := &File{
			Name:       path,
			Size:       info.Size(),
			Modified:   info.ModTime(),
			MIME:       "dir",
			collection: ret,
		}
		if !info.IsDir() {
			file, _ := os.Open(path)
			kh := &KeepHeaders{r: file}
			hash := fcomp.Hash(kh)
			ft, _ := filetype.Get(kh.headers[:])
			f.Hash = hash
			f.MIME = ft.MIME.Value
		}
		ret.Files = append(ret.Files, f)
		return nil
	})

	return ret
}

func (fc *Collection) ExportToJSON(output string) error {
	f, err := os.Create(output)
	if err != nil {
		return err
	}

	return json.NewEncoder(f).Encode(fc)
}

func ImportFromJSON(output string) (*Collection, error) {
	f, err := os.Open(output)
	if err != nil {
		return nil, err
	}

	var fc Collection
	err = json.NewDecoder(f).Decode(&fc)
	if err != nil {
		return nil, err
	}

	// Create links to parent Collection
	for _, file := range fc.Files {
		file.collection = &fc
	}

	return &fc, err
}

func DataFromJSON(paths []string) *Data {
	ret := &Data{}
	ret.Collections = make(map[string]*Collection)
	ret.Files = make(map[string][]*File)
	for _, p := range paths {
		c, err := ImportFromJSON(p)
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
