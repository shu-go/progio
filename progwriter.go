package progio

import (
	"io"
	"time"
)

// Writer a progress watcher and a wrapper of io.Writer.
type Writer struct {
	io.Writer

	handler   func(progress int64)
	throttler Throttler

	progress int64
	start    time.Time
}

// NewWriter creates a Writer with handler to handle progress.
// optThrottler is used to throttle handler call.
func NewWriter(dst io.Writer, handler func(progress int64), optThrottler ...Throttler) *Writer {
	w := &Writer{
		Writer:    dst,
		handler:   handler,
		throttler: &nullThrottling{},
	}

	if len(optThrottler) > 0 && optThrottler[0] != nil {
		w.throttler = optThrottler[0]
	}

	return w
}

// Write implements (io.Writer).Write.
// It calls progress handler.
func (r *Writer) Write(p []byte) (n int, err error) {
	var zero time.Time
	if r.start == zero {
		r.start = time.Now()
	}
	nn, ee := r.Writer.Write(p)
	r.progress += int64(nn)
	r.throttler.CallHandler(r.handler, r.progress, time.Now().Sub(r.start))
	return nn, ee
}
