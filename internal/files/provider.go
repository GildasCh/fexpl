package files

import (
	"io"

	"github.com/pkg/errors"
)

type Provider struct {
	protocol string
	url      string

	a adapter
}

var protocols = map[string](func(url string) (adapter, error)){
	"local": func(url string) (adapter, error) { return newLocal(url) },
}

type adapter interface {
	reader(path string) (Reader, error)
}

func New(protocol string, url string) (*Provider, error) {
	init, ok := protocols[protocol]

	if !ok {
		return nil, errors.Errorf("unknow protocol %q", protocol)
	}

	a, err := init(url)
	if err != nil {
		return nil, errors.Wrapf(err, "error creating provider %q, %q", protocol, url)
	}
	return &Provider{
		protocol: protocol,
		url:      url,
		a:        a,
	}, nil
}

func (p *Provider) Reader(path string) (io.Reader, error) {
	return p.a.reader(path)
}
