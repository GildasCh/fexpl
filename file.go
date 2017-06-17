package main

import "time"

type File struct {
	Name       string
	Hash       string
	Size       int64
	Modified   time.Time
	MIME       string
	collection *Collection
}
