package progio

import (
	"io"
	"time"
)

// Reader a progress watcher and a wrapper of io.Reader.
type Reader struct {
	io.Reader

	handlerCaller func(int64, time.Duration)

	throttler Throttler

	progress int64
	start    time.Time
}

// NewReader creates a Reader with handler to handle progress.
// optThrottler is used to throttle handler call.
func NewReader(src io.Reader, handler interface{}, optThrottler ...Throttler) *Reader {
	r := &Reader{
		Reader:    src,
		throttler: &nullThrottling{},
	}

	if h, ok := handler.(func(int64)); ok {
		r.handlerCaller = func(p int64, _ time.Duration) {
			h(p)
		}
	} else if h, ok := handler.(func(int64, time.Duration)); ok {
		r.handlerCaller = func(p int64, d time.Duration) {
			h(p, d)
		}
	} else {
		r.handlerCaller = func(p int64, d time.Duration) {
		}
	}

	if len(optThrottler) > 0 && optThrottler[0] != nil {
		r.throttler = optThrottler[0]
	}

	return r
}

// Read implements (io.Reader).Read.
// It calls progress handler.
func (r *Reader) Read(p []byte) (n int, err error) {
	var zero time.Time
	if r.start == zero {
		r.start = time.Now()
	}
	nn, ee := r.Reader.Read(p)
	r.progress += int64(nn)
	r.throttler.CallHandler(r.handlerCaller, r.progress, time.Now().Sub(r.start))
	return nn, ee
}
