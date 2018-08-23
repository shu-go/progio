package progio

import "time"

type Throttler interface {
	CallHandler(func(int64), int64)
}

type PercentThrottling struct {
	Max   int64
	Scale int

	last int
}

func (t *PercentThrottling) CallHandler(handler func(int64), n int64) {
	test := int(float64(n) / float64(t.Max) * 100)
	if t.last+t.Scale <= test {
		handler(n)
		t.last = test
	}
}

type TimeThrottling struct {
	Duration time.Duration

	last time.Time
}

func (t *TimeThrottling) CallHandler(handler func(int64), n int64) {
	test := t.last.Add(t.Duration)
	if t.last.Before(test) {
		handler(n)
		t.last = test
	}
}
