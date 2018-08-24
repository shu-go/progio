package progio

import "time"

type RWConfig func(rw interface{})

func Percent(max int64, scale int) RWConfig {
	return func(rw interface{}) {
		if r, ok := rw.(*Reader); ok {
			r.Throttler = NewPercentThrottler(max, scale)
		} else if w, ok := rw.(*Writer); ok {
			w.Throttler = NewPercentThrottler(max, scale)
		}
	}
}

func Time(d time.Duration) RWConfig {
	return func(rw interface{}) {
		if r, ok := rw.(*Reader); ok {
			r.Throttler = NewTimeThrottler(d)
		} else if w, ok := rw.(*Writer); ok {
			w.Throttler = NewTimeThrottler(d)
		}
	}
}
