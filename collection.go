package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	filetype "gopkg.in/h2non/filetype.v1"

	"github.com/gildasch/fcomp"
)

type Collection struct {
	Files []*File
	Name  string
}

func Explore(name, root string, hidden bool, maxSize int64) *Collection {
	if _, err := os.Stat(filepath.Join(root, "fexpl.json")); err == nil {
		fmt.Printf("File exists on root %q\n", root)
		ret, err := importr.Import(filepath.Join(root, "fexpl.json"))
		if err == nil {
			return ret
		}
		fmt.Println(err)
	}

	if name == "" {
		return nil
	}

	counterChan := make(chan string, 100)
	go func() {
		counter := 0
		last := "none"
		timeChan := time.After(time.Second)
		fmt.Println()
		for {
			select {
			case <-timeChan:
				fmt.Printf("\rWalking over %d files, last was %s", counter, last)
				timeChan = time.After(100 * time.Millisecond)
			case name := <-counterChan:
				if name == "" {
					fmt.Println()
					return
				}
				counter++
				last = name
			}
		}
	}()

	ret := &Collection{}
	ret.Name = name
	filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if !hidden && strings.Contains(path, "/.") {
			return nil // skip hidden file
		}
		counterChan <- path
		f := &File{
			Name:       path,
			Size:       info.Size(),
			Modified:   info.ModTime(),
			MIME:       "dir",
			collection: ret,
		}
		if !info.IsDir() {
			var kh *KeepHeaders

			// Hash calculation
			if info.Size() <= maxSize {
				file, _ := os.Open(path)
				kh = &KeepHeaders{r: file}
				hash := fcomp.Hash(kh)
				f.Hash = hash
			}

			// Extract file type from header
			if info.Size() > 261 {
				if kh != nil {
					ft, _ := filetype.Get(kh.headers[:])
					f.MIME = ft.MIME.Value
				} else {
					file, _ := os.Open(path)
					buf := make([]byte, 261)
					file.Read(buf)
					ft, _ := filetype.Get(buf[:261])
					f.MIME = ft.MIME.Value
				}
			}
		}
		ret.Files = append(ret.Files, f)
		return nil
	})
	counterChan <- ""

	return ret
}

func (fc *Collection) ExportToJSON(output string) error {
	f, err := os.Create(output)
	if err != nil {
		return err
	}

	return json.NewEncoder(f).Encode(fc)
}
