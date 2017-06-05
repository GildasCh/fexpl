package main

import (
	"os"
	"path/filepath"
	"strings"
	"time"

	filetype "gopkg.in/h2non/filetype.v1"
	filetypes "gopkg.in/h2non/filetype.v1/types"

	"github.com/GildasCh/fcomp"
)

type File struct {
	Name     string
	Hash     string
	Size     int64
	Modified time.Time
	Header   filetypes.Type
}

type FileCollection struct {
	Files []File
	Name  string
}

func Explore(name, root string) FileCollection {
	ret := FileCollection{}
	filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		file, _ := os.Open(path)
		kh := &KeepHeaders{r: file}
		hash := fcomp.Hash(kh)
		ft, _ := filetype.Get(kh.headers[:])
		ret.Files = append(ret.Files, File{
			Name:     strings.TrimPrefix(path, root),
			Hash:     hash,
			Size:     info.Size(),
			Modified: info.ModTime(),
			Header:   ft,
		})
		return nil
	})

	return ret
}
