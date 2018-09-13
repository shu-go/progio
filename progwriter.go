package progio

import (
	"io"
	"time"
)

// Writer a progress watcher and a wrapper of io.Writer.
type Writer struct {
	io.Writer

	handlerCaller func(int64, time.Duration)

	throttler Throttler

	progress int64
	start    time.Time
}

// NewWriter creates a Writer with handler to handle progress.
// optThrottler is used to throttle handler call.
func NewWriter(dst io.Writer, handler interface{}, optThrottler ...Throttler) *Writer {
	w := &Writer{
		Writer:    dst,
		throttler: &nullThrottling{},
	}

	if h, ok := handler.(func(int64)); ok {
		w.handlerCaller = func(p int64, _ time.Duration) {
			h(p)
		}
	} else if h, ok := handler.(func(int64, time.Duration)); ok {
		w.handlerCaller = func(p int64, d time.Duration) {
			h(p, d)
		}
	} else {
		w.handlerCaller = func(p int64, d time.Duration) {
		}
	}

	if len(optThrottler) > 0 && optThrottler[0] != nil {
		w.throttler = optThrottler[0]
	}

	return w
}

// Write implements (io.Writer).Write.
// It calls progress handler.
func (w *Writer) Write(p []byte) (n int, err error) {
	var zero time.Time
	if w.start == zero {
		w.start = time.Now()
	}
	nn, ee := w.Writer.Write(p)
	w.progress += int64(nn)
	w.throttler.CallHandler(w.handlerCaller, w.progress, time.Now().Sub(w.start))
	return nn, ee
}
