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

func (dr *DummyReader) Read(p []byte) (int, error) {
	n := rand.Intn(len(p))

	for i := 0; i < n; i++ {
		if dr.counter >= len(dr.b) {
			return i, errors.New("End of array")
		}

		p[i] = dr.b[dr.counter]
		dr.counter++
	}

	return n, nil
}

func TestKeepHeaders(t *testing.T) {
	rand.Seed(time.Now().UTC().UnixNano())
	in := randomArray(headersSize + rand.Intn(10*headersSize))
	var out []byte
	dr := &DummyReader{b: in}
	kh := &KeepHeaders{r: dr}

	var err error
	for err == nil {
		p := make([]byte, headersSize)
		var n int
		n, err = kh.Read(p)
		out = append(out, p[:n]...)
	}

	assert.Equal(t, headersSize, kh.read)
	assert.Equal(t, in, out)
	assert.Equal(t, in[:headersSize], kh.headers[:])
}

func randomArray(n int) []byte {
	ret := make([]byte, n)
	rand.Read(ret)
	return ret
}
