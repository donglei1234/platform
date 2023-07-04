package utils

import (
	"errors"
	"math"
	"time"
)

var (
	TestErrRetryErr = errors.New("TestErrRetryErr")
	RetryStandard   = Retry{
		delay: ExponentialSleep{
			min:    34,
			max:    99,
			jitter: 3900,
			delay:  34,
		},
		maxRetries: 3,
	}
)

func (r *Retry) GetDelay() ExponentialSleep { return r.delay }
func (r *Retry) GetMaxRetries() int         { return r.maxRetries }

// redundant; remove after exponential_sleep tests are merged in
func (e ExponentialSleep) Equals(sleep ExponentialSleep) bool {
	return e.min == sleep.min &&
		e.max == sleep.max &&
		e.delay == sleep.delay &&
		(e.jitter == sleep.jitter || (math.IsNaN(e.jitter) && math.IsNaN(sleep.jitter)))
}

func (r *Retry) Equals(rr *Retry) bool {
	return r.delay.Equals(rr.delay) && r.maxRetries == rr.maxRetries
}

type TestSleeper struct {
	SleepTime int64
}

func (t *TestSleeper) Sleep(d time.Duration) {
	if n := d.Nanoseconds(); n < 0 {
		t.SleepTime = 0
	} else {
		t.SleepTime = n
	}
}
