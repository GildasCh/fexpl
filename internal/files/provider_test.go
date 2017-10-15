package files

import (
	"io"
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
)

type dummyProvider struct{}

func (dummyProvider) reader(path string) (Reader, error) {
	return Reader(func(b []byte) (n int, err error) {
		if len(b) >= 1 {
			b[0] = byte('a')
		}
		return 1, io.EOF
	}), nil
}

func TestProtocolOK(t *testing.T) {
	protocols["test protocol"] = func(string) (adapter, error) {
		return dummyProvider{}, nil
	}

	p, err := New("test protocol", "")
	assert.NoError(t, err)

	r, err := p.Reader("")
	assert.NoError(t, err)

	b, err := ioutil.ReadAll(r)
	assert.NoError(t, err)
	assert.Equal(t, "a", string(b))
}

func TestProtocolNotFound(t *testing.T) {
	_, err := New("protocol not found", "")
	assert.EqualError(t, err, `unknow protocol "protocol not found"`)
}
