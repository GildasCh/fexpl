package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDuplicates(t *testing.T) {
	f1 := &File{
		Name: "dummy name",
		Hash: "123",
	}
	f2 := &File{
		Name: "dummy name2",
	}

	data = &Data{
		Files: map[string][]*File{
			"123": []*File{f1, f2},
		}}

	assert.Equal(t, []*File{f1, f2}, f1.Duplicates())
}

func TestFormatSize(t *testing.T) {
	in := []int64{
		434243234, 443, 453592249, 53842893237624722,
		11223, 0, 3432, 1, 2, 3, 50}
	expected := []string{
		"414.13M", "443.00", "432.58M", "48969.83T",
		"10.96K", "0.00", "3.35K", "1.00", "2.00", "3.00", "50.00"}

	actual := []string{}
	for i := range in {
		actual = append(actual, (&File{Size: in[i]}).FormatSize())
	}

	assert.Equal(t, expected, actual)

}
