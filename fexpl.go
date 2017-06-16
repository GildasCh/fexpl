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
	Name     string
	Hash     string
	Size     int64
	Modified time.Time
	MIME     string
}

type FileCollection struct {
	Files []File
	Name  string
}

func Explore(name, root string) *FileCollection {
	if _, err := os.Stat(filepath.Join(root, "fexpl.json")); err == nil {
		// File exists
		ret, err := ImportFromJSON(filepath.Join(root, "fexpl.json"))
		if err == nil {
			return ret
		}
		fmt.Println(err)
	}

	if name == "" {
		return nil
	}

	ret := &FileCollection{}
	ret.Name = name
	filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() {
			ret.Files = append(ret.Files, File{
				Name:     path,
				Modified: info.ModTime(),
				MIME:     "dir",
			})
			return nil
		}
		file, _ := os.Open(path)
		kh := &KeepHeaders{r: file}
		hash := fcomp.Hash(kh)
		ft, _ := filetype.Get(kh.headers[:])
		ret.Files = append(ret.Files, File{
			Name:     path,
			Hash:     hash,
			Size:     info.Size(),
			Modified: info.ModTime(),
			MIME:     ft.MIME.Value,
		})
		return nil
	})

	return ret
}

func (fc *FileCollection) ExportToJSON(output string) error {
	f, err := os.Create(output)
	if err != nil {
		return err
	}

	return json.NewEncoder(f).Encode(fc)
}

func ImportFromJSON(output string) (*FileCollection, error) {
	f, err := os.Open(output)
	if err != nil {
		return nil, err
	}

	var fc FileCollection
	err = json.NewDecoder(f).Decode(&fc)
	return &fc, err
}
