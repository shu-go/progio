package progio

import "time"

type Throttler interface {
	CallHandler(func(int64), int64)
}

type nullThrottling struct{}

func (t *nullThrottling) CallHandler(handler func(int64), n int64) {
	handler(n)
}

////////////////////////////////////////////////////////////////////////////////

// Percent throttles progress handling by the scale(%).
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

func (t *percentThrottling) percentInScale(n int64) int {
	return int(float64(n)/float64(t.max)*100/float64(t.scale)) * t.scale
}

func (t *percentThrottling) CallHandler(handler func(int64), n int64) {
	test := t.percentInScale(n)
	if t.last+t.scale <= test {
		handler( /*n*/ int64(test))
		t.last = test
	}
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

func (t *timeThrottling) CallHandler(handler func(int64), n int64) {
	test := time.Now()
	if t.last.Add(t.d).Before(test) {
		handler(n)
		t.last = test
	}
}
