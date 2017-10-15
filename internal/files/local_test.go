package files

import (
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLocalReadSuccessful(t *testing.T) {
	l, err := newLocal(".")
	assert.NoError(t, err)

	r, err := l.reader("test_dummy_file.txt")
	assert.NoError(t, err)

	b, err := ioutil.ReadAll(r)
	assert.NoError(t, err)
	assert.Equal(t, "dummy\n", string(b))
}

func TestLocalRootNotFound(t *testing.T) {
	_, err := newLocal("not found")
	assert.EqualError(t, err, "error stating local file root: stat not found: no such file or directory")
}

func TestLocalFileNotFound(t *testing.T) {
	l, err := newLocal(".")
	assert.NoError(t, err)

	_, err = l.reader("file_that_does_not_exist.txt")
	assert.EqualError(t, err, "error opening file file_that_does_not_exist.txt: open file_that_does_not_exist.txt: no such file or directory")
}
