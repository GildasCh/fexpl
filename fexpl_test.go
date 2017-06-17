package main

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

var fc *Collection

func TestExplore(t *testing.T) {
	if fc == nil {
		fc = Explore("dummy name", "data")
	}

	assert.Equal(t, "dummy name", fc.Name)
	assert.Equal(t, 18, len(fc.Files))

	// Check that hashes of a copy is the same
	var h1, h2 string
	for _, f := range fc.Files {
		if f.Name == "Double-compound-pendulum.gif" {
			h1 = f.Hash
		}
		if f.Name == "Copy_of-the-pendulum.gif" {
			h2 = f.Hash
		}
	}

	assert.Equal(t, h1, h2)
}

func TestExploreDataExists(t *testing.T) {
	if fc == nil {
		fc = Explore("dummy name", "data")
	}

	fc.ExportToJSON("data/fexpl.json")
	defer func() {
		err := os.Remove("data/fexpl.json")
		if err != nil {
			fmt.Println(err)
		}
	}()

	fc2 := Explore("", "data")
	for i := range fc.Files {
		fc.Files[i].Modified = fc.Files[i].Modified.UTC()
	}
	for i := range fc2.Files {
		fc2.Files[i].Modified = fc2.Files[i].Modified.UTC()
	}
	assert.Equal(t, fc, fc2)
}

func TestExploreErrorOnNoNameAndNoData(t *testing.T) {
	fc2 := Explore("", "data")
	assert.EqualValues(t, (*Collection)(nil), fc2)
}

func TestExportAndImport(t *testing.T) {
	if fc == nil {
		fc = Explore("dummy name", "data")
	}

	err := fc.ExportToJSON("dummy.json")
	assert.NoError(t, err)
	fc2, err := ImportFromJSON("dummy.json")
	assert.NoError(t, err)

	// fix bug on time.Time.loc
	for i := range fc.Files {
		fc.Files[i].Modified = fc.Files[i].Modified.UTC()
	}
	for i := range fc2.Files {
		fc2.Files[i].Modified = fc2.Files[i].Modified.UTC()
	}

	assert.Equal(t, fc, fc2)
}
