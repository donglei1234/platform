package utils

import (
	"math/rand"
	"time"
)

// ExponentialSleep encapsulates a simple backoff algorithm.
type ExponentialSleep struct {
	delay, min, max time.Duration
	jitter          float64
}

func NewExponentialSleep(min, max time.Duration, jitter float64) ExponentialSleep {
	if min < 0 || max < 0 {
		min = 0
		max = 0
	} else if min > max {
		min = max
	}

	return ExponentialSleep{
		delay:  min,
		min:    min,
		max:    max,
		jitter: jitter,
	}
}

func (e *ExponentialSleep) Reset() {
	e.delay = e.min
}

func (e *ExponentialSleep) Sleep() {

	if e.jitter != 0 {
		j := e.delay.Nanoseconds()
		j = j + int64(float64(j)*((rand.Float64()*e.jitter)-(e.jitter/2)))

		time.Sleep(time.Duration(j) * time.Nanosecond)
	} else {
		time.Sleep(e.delay)
	}

	e.delay = e.delay + e.delay

	if e.delay >= e.max {
		e.delay = e.max
	}
}
