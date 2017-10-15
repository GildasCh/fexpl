package files

type Reader func(b []byte) (n int, err error)

func (r Reader) Read(b []byte) (n int, err error) {
	return r(b)
}
