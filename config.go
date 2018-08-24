package progio

import "time"

type RWConfig func(rw interface{})

// Percent is a config func for Reader and Writer.
// It throttles progress handling by the scale(%).
//
// This config CHANGES a handler parameter BYTES->PERCENT.
func Percent(max int64, scale int) RWConfig {
	return func(rw interface{}) {
		if r, ok := rw.(*Reader); ok {
			r.Throttler = NewPercentThrottler(max, scale)
		} else if w, ok := rw.(*Writer); ok {
			w.Throttler = NewPercentThrottler(max, scale)
		}
	}
}

// Time is a config func for Reader and Writer.
// It throttles progress handling by duration d.
func Time(d time.Duration) RWConfig {
	t := &timeThrottling{
		Duration: d,
	}

	return func(rw interface{}) {
		if r, ok := rw.(*Reader); ok {
			r.Throttler = NewTimeThrottler(d)
		} else if w, ok := rw.(*Writer); ok {
			w.Throttler = NewTimeThrottler(d)
		}
	}
}
