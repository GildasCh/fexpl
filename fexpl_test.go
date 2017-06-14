package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExplore(t *testing.T) {
	fc := Explore("dummy name", "data")

	assert.Equal(t, "dummy name", fc.Name)
	assert.Equal(t, 15, len(fc.Files))

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
