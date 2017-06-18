package main

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDataFromJSON(t *testing.T) {
	expected := map[string]*Collection{
		"A": &Collection{nil, "A"},
		"B": &Collection{nil, "B"},
		"C": &Collection{
			[]*File{&File{Name: "dummy file", Hash: "123"}},
			"C"},
	}

	importFunc = func(input string) (*Collection, error) {
		return expected[input], nil
	}

	d := DataFromJSON([]string{"A", "C", "B"})

	assert.Equal(t, expected, d.Collections)
	assert.Equal(t, map[string][]*File{"123": []*File{&File{Name: "dummy file", Hash: "123"}}}, d.Files)
}

func TestDataFromJSONImportError(t *testing.T) {
	importFunc = func(input string) (*Collection, error) {
		return nil, fmt.Errorf("dummy error")
	}

	d := DataFromJSON([]string{"A", "C", "B"})

	assert.Equal(t, 0, len(d.Collections))
	assert.Equal(t, 0, len(d.Files))
}

func TestDataFromJSONDuplicateCollectionNameKeepFirst(t *testing.T) {
	expected := map[string]*Collection{
		"A": &Collection{[]*File{&File{Name: "dummy file 1", Hash: "123"}}, "A"},
		"B": &Collection{nil, "B"},
	}

	importFunc = func(input string) (*Collection, error) {
		if input == "C" {
			return &Collection{[]*File{&File{Name: "dummy file 2", Hash: "123"}}, "A"}, nil
		}
		return expected[input], nil
	}

	d := DataFromJSON([]string{"A", "C", "B"})

	assert.Equal(t, expected, d.Collections)
	assert.Equal(t, map[string][]*File{"123": []*File{&File{Name: "dummy file 1", Hash: "123"}}}, d.Files)
}
