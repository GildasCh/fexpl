package main

import (
	"encoding/json"
	"os"
)

type Importer interface {
	Import(input string) (*Collection, error)
}

type JSONImporter struct{}

func (_ *JSONImporter) Import(input string) (*Collection, error) {
	f, err := os.Open(input)
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
