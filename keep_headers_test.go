package main

import (
	"errors"
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type DummyReader struct {
	b       []byte
	counter int
}

func (dr *DummyReader) Read(p []byte) (n int, err error) {
	n = rand.Intn(len(p))

	for i := 0; i < n; i++ {
		if dr.counter >= len(dr.b) {
			return 0, errors.New("End of array")
		}

		p[i] = dr.b[dr.counter]
		dr.counter++
	}

	return
}

func TestKeepHeaders(t *testing.T) {
	dr := &DummyReader{b: randomArray(2 * headersSize)}
	kh := &KeepHeaders{r: dr}

	var err error
	for err == nil {
		p := make([]byte, headersSize)
		_, err = kh.Read(p)
	}

	assert.Equal(t, headersSize, kh.read)
	assert.Equal(t, dr.b[:headersSize], kh.headers[:])
}

func randomArray(n int) []byte {
	ret := make([]byte, n)
	rand.Seed(time.Now().UTC().UnixNano())
	rand.Read(ret)
	return ret
}
