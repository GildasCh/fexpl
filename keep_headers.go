package main

import "io"

type KeepHeaders struct {
	r       io.Reader
	headers [headersSize]byte
	read    int
}

const headersSize = 261

func (kh *KeepHeaders) Read(p []byte) (n int, err error) {
	n, err = kh.r.Read(p)

	if kh.read < headersSize {
		if n >= headersSize-kh.read {
			copy(kh.headers[kh.read:], p[:headersSize-kh.read])
			kh.read += headersSize - kh.read
		} else {
			copy(kh.headers[kh.read:], p)
			kh.read += n
		}
	}

	return
}
