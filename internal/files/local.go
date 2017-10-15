package files

import (
	"os"
	"path"

	"github.com/pkg/errors"
)

type local struct {
	root string
}

func newLocal(url string) (*local, error) {
	info, err := os.Stat(url)
	if err != nil {
		return nil, errors.Wrap(err, "error stating local file root")
	}
	if !info.IsDir() {
		return nil, errors.Wrap(err, "root must be a directory")
	}
	return &local{root: url}, nil
}

func (l *local) reader(relPath string) (Reader, error) {
	fullPath := path.Join(l.root, relPath)
	f, err := os.Open(fullPath)
	if err != nil {
		return nil, errors.Wrapf(err, "error opening file %s", fullPath)
	}
	return Reader(f.Read), nil
}
