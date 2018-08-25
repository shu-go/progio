package progio

import "io"

// Reader a progress watcher and a wrapper of io.Reader.
type Reader struct {
	io.Reader

	handler   func(progress int64)
	throttler Throttler

	progress int64
}

// NewReader creates a Reader with handler to handle progress.
// optThrottler is used to throttle handler call.
func NewReader(src io.Reader, handler func(progress int64), optThrottler ...Throttler) *Reader {
	r := &Reader{
		Reader:    src,
		handler:   handler,
		throttler: &nullThrottling{},
	}

	if len(optThrottler) > 0 && optThrottler[0] != nil {
		r.throttler = optThrottler[0]
	}

	return r
}

// Read implements (io.Reader).Read.
// It calls progress handler.
func (r *Reader) Read(p []byte) (n int, err error) {
	nn, ee := r.Reader.Read(p)
	r.progress += int64(nn)
	r.throttler.CallHandler(r.handler, r.progress)
	return nn, ee
}
