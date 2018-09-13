package progio

import (
	"time"
)

// Throttler throttles handling of Reader/Writer progress.
type Throttler interface {
	// CallHandler calls a progress handler if needed.
	CallHandler(handler func(int64, time.Duration), progress int64, duration time.Duration)
}

type nullThrottling struct {
	duration time.Duration
}

func (t *nullThrottling) CallHandler(handler func(int64, time.Duration), p int64, d time.Duration) {
	handler(p, d)
}

////////////////////////////////////////////////////////////////////////////////

// Percent throttles progress handling by the scale(%).
// Param max is the size of target Reader to calc percent.
// e.g. response.ContentLength
//
// This throttler CHANGES a handler parameter BYTES->PERCENT.
func Percent(max int64, scale int) Throttler {
	if scale < 0 || 100 < scale {
		panic("0 <= scale <= 100")
	}
	return &percentThrottling{
		max:   max,
		scale: scale,
	}
}

type percentThrottling struct {
	max   int64
	scale int

	last int
}

func (t *percentThrottling) percentInScale(p int64) int {
	return int(float64(p)/float64(t.max)*100/float64(t.scale)) * t.scale
}

func (t *percentThrottling) CallHandler(handler func(int64, time.Duration), p int64, d time.Duration) {
	test := t.percentInScale(p)
	if t.last+t.scale <= test {
		handler( /*p*/ int64(test), d)
	}
	t.last = test
}

////////////////////////////////////////////////////////////////////////////////

// Time throttles progress handling by duration d.
func Time(d time.Duration) Throttler {
	return &timeThrottling{
		d: d,
	}
}

type timeThrottling struct {
	d time.Duration

	last time.Time
}

func (t *timeThrottling) CallHandler(handler func(int64, time.Duration), p int64, d time.Duration) {
	test := time.Now()
	if t.last.Add(t.d).Before(test) {
		handler(p, d)
		t.last = test
	}
}
