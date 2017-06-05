package main

import (
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/GildasCh/fcomp"
)

type File struct {
	Name     string
	Hash     string
	Size     int64
	Modified time.Time
}

type FileCollection struct {
	Files []File
	Name  string
}

func Explore(name, root string) FileCollection {
	ret := FileCollection{}
	filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		file, _ := os.Open(path)
		hash := fcomp.Hash(file)

		ret.Files = append(ret.Files, File{
			Name:     strings.TrimPrefix(path, root),
			Hash:     hash,
			Size:     info.Size(),
			Modified: info.ModTime(),
		})
		return nil
	})

	return ret
}
