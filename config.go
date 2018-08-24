package progio

import "time"

type RWConfig func(rw interface{})

// Percent is a config func for Reader and Writer.
// It throttles progress handling by the scale(%).
//
// This config CHANGES a handler parameter BYTES->PERCENT.
func Percent(max int64, scale int) RWConfig {
	if scale < 0 || 100 < scale {
		panic("0 <= scale <= 100")
	}
	t := &percentThrottling{
		Max:   max,
		Scale: scale,
	}

	return func(rw interface{}) {
		if r, ok := rw.(*Reader); ok {
			r.Throttler = t
		} else if w, ok := rw.(*Writer); ok {
			w.Throttler = t
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
			r.Throttler = t
		} else if w, ok := rw.(*Writer); ok {
			w.Throttler = t
		}
	}
}
