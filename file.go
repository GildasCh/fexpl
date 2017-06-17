package main

import (
	"fmt"
	"time"
)

type File struct {
	Name       string
	Hash       string
	Size       int64
	Modified   time.Time
	MIME       string
	collection *Collection
}

func (f *File) Duplicates() []*File {
	return data.Files[f.Hash]
}

func (f *File) FormatSize() string {
	units := []string{"", "K", "M", "G", "T"}
	unit := 0
	size := float64(f.Size)
	for size > 1024 {
		size /= 1024
		unit++
		if unit == len(units)-1 {
			break
		}
	}
	return fmt.Sprintf("%.2f%s", size, units[unit])
}
