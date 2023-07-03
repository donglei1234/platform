package utils

import (
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func ShouldRetryGrpcError(e error) bool {
	errStatus, ok := status.FromError(e)
	if ok {
		switch errStatus.Code() {
		case codes.Unavailable:
			return true
		}
	}
	return false
}

type Retry struct {
	delay      ExponentialSleep
	maxRetries int
}

func NewRetry(minDelay, maxDelay time.Duration, maxRetries int, jitter float64) *Retry {
	return &Retry{delay: NewExponentialSleep(minDelay, maxDelay, jitter), maxRetries: maxRetries}
}

func (r *Retry) Retry(f func() (bool, error)) (err error) {
	retries := 0
	for {
		retry, e := f()

		if retry {
			if retries < r.maxRetries {
				r.delay.Sleep()
				retries++
				continue
			}
		}
		err = e
		break
	}
	return
}
