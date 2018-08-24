package progio

import "time"

type Throttler interface {
	CallHandler(func(int64), int64)
}

type nullThrottling struct{}

func (t *nullThrottling) CallHandler(handler func(int64), n int64) {
	handler(n)
}

type percentThrottling struct {
	Max   int64
	Scale int

	last int
}

func (t *percentThrottling) percentInScale(n int64) int {
	return int(float64(n)/float64(t.Max)*100/float64(t.Scale)) * t.Scale
}

func (t *percentThrottling) CallHandler(handler func(int64), n int64) {
	test := t.percentInScale(n)
	if t.last+t.Scale <= test {
		handler( /*n*/ int64(test))
		t.last = test
	}
}

type timeThrottling struct {
	Duration time.Duration

	last time.Time
}

func (t *timeThrottling) CallHandler(handler func(int64), n int64) {
	test := time.Now()
	if t.last.Add(t.Duration).Before(test) {
		handler(n)
		t.last = test
	}
}
