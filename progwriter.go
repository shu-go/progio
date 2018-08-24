package progio

import "io"

type Writer struct {
	io.Writer
	Handler func(progress int64)

	Throttler Throttler

	progress int64
}

func NewWriter(src io.Writer, handler func(progress int64), configs ...RWConfig) *Writer {
	w := &Writer{
		Writer:    src,
		Handler:   handler,
		Throttler: &nullThrottling{},
	}

	for _, c := range configs {
		c(w)
	}

	return w
}

func (r *Writer) Write(p []byte) (n int, err error) {
	nn, ee := r.Writer.Write(p)
	r.progress += int64(nn)
	r.Throttler.CallHandler(r.Handler, r.progress)
	return nn, ee
}
