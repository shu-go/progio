package progio

import "io"

// Reader a progress watcher and a wrapper for io.Reader.
type Reader struct {
	io.Reader
	Handler func(progress int64)

	Throttler Throttler

	progress int64
}

// NewReader creates a Reader with handler to handle progress.
// configs is currently used to set Throttler property.
func NewReader(src io.Reader, handler func(progress int64), configs ...RWConfig) *Reader {
	r := &Reader{
		Reader:    src,
		Handler:   handler,
		Throttler: &nullThrottling{},
	}

	for _, c := range configs {
		c(r)
	}

	return r
}

// Read implements (io.Reader).Read.
// It calls progress handler.
func (r *Reader) Read(p []byte) (n int, err error) {
	nn, ee := r.Reader.Read(p)
	r.progress += int64(nn)
	r.Throttler.CallHandler(r.Handler, r.progress)
	return nn, ee
}
