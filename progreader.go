package progio

import "io"

type Reader struct {
	io.Reader
	Handler func(progress int64)

	Throttling Throttler

	progress int64
}

func NewReader(src io.Reader, handler func(progress int64)) *Reader {
	return &Reader{
		Reader:  src,
		Handler: handler,
	}
}

func (r *Reader) Read(p []byte) (n int, err error) {
	nn, ee := r.Reader.Read(p)
	r.progress += int64(nn)
	if r.Throttling == nil {
		r.Handler(r.progress)
	} else {
		r.Throttling.CallHandler(r.Handler, r.progress)
	}

	return nn, ee
}
