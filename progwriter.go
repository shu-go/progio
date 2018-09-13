package progio

import (
	"io"
	"time"
)

// Writer a progress watcher and a wrapper of io.Writer.
type Writer struct {
	io.Writer

	listenerCaller func(int64, time.Duration)

	throttler Throttler

	progress int64
	start    time.Time
}

// NewWriter creates a Writer with listener to handle progress.
// listener is a function func(progress int64) or func(progress int64, duration time.Duration).
// optThrottler is used to throttle listener call.
func NewWriter(dst io.Writer, listener interface{}, optThrottler ...Throttler) *Writer {
	w := &Writer{
		Writer:    dst,
		throttler: &nullThrottling{},
	}

	if h, ok := listener.(func(int64)); ok {
		w.listenerCaller = func(p int64, _ time.Duration) {
			h(p)
		}
	} else if h, ok := listener.(func(int64, time.Duration)); ok {
		w.listenerCaller = func(p int64, d time.Duration) {
			h(p, d)
		}
	} else {
		w.listenerCaller = func(p int64, d time.Duration) {
		}
	}

	if len(optThrottler) > 0 && optThrottler[0] != nil {
		w.throttler = optThrottler[0]
	}

	return w
}

// Write implements (io.Writer).Write.
// It calls progress listener.
func (w *Writer) Write(p []byte) (n int, err error) {
	var zero time.Time
	if w.start == zero {
		w.start = time.Now()
	}
	nn, ee := w.Writer.Write(p)
	w.progress += int64(nn)
	w.throttler.CallListener(w.listenerCaller, w.progress, time.Now().Sub(w.start))
	return nn, ee
}
