package progio

import "io"

type Writer struct {
	io.Writer
	Handler func(progress int64)

	progress int64
}

func NewWriter(src io.Writer, handler func(progress int64)) *Writer {
	return &Writer{
		Writer:  src,
		Handler: handler,
	}
}

func (r *Writer) Write(p []byte) (n int, err error) {
	nn, ee := r.Writer.Write(p)
	r.progress += int64(nn)
	r.Handler(r.progress)
	return nn, ee
}
