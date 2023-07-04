package utils_test

import (
	"errors"
	"github.com/donglei1234/platform/services/common/utils"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"runtime/debug"
	"testing"
	"time"
)

func TestNewRetry(t *testing.T) {
	testcases := map[string]struct {
		minDelay, maxDelay time.Duration
		maxRetries         int
		jitter             float64
		expErr             error
		expResult          utils.Retry
		expPanic           bool
	}{"Standard": {
		minDelay:   34,
		maxDelay:   99,
		maxRetries: 3,
		jitter:     3900,
		expResult:  utils.RetryStandard,
		expPanic:   false,
	}}
	for i, tc := range testcases {
		func() {
			defer func() {
				if r := recover(); r != nil && !tc.expPanic {
					t.Errorf("Test %s: unexpected panic occurred: %s\n%s", i, r, debug.Stack())
				}
			}()
			res := utils.NewRetry(tc.minDelay, tc.maxDelay, tc.maxRetries, tc.jitter)
			if !res.Equals(&tc.expResult) {
				t.Errorf("Test %s: Received %#v; expected %#v", i, res, tc.expResult)
			}
			if tc.expPanic {
				t.Errorf("Test %s: Unexpectidly did not panic", i)
			}
		}()
	}
}

func TestShouldRetryGrpcError(t *testing.T) {
	testcases := map[string]struct {
		error
		expResult bool
		expPanic  bool
	}{
		"FromErrorOkUnavailable": {
			error:     status.Error(codes.Unavailable, ""),
			expResult: true,
			expPanic:  false,
		},
		"FromErrorNotOk": {
			error: errors.New("TestErr"),
		},
	}
	for i, tc := range testcases {
		func() {
			defer func() {
				if r := recover(); r != nil && !tc.expPanic {
					t.Errorf("Test %s: unexpected panic occurred: %s\n%s", i, r, debug.Stack())
				}
			}()
			if res := utils.ShouldRetryGrpcError(tc.error); res != tc.expResult {
				t.Errorf("Test %s: Received %t; expected %t", i, res, tc.expResult)
			}
			if tc.expPanic {
				t.Errorf("Test %s: Unexpectidly did not panic", i)
			}
		}()
	}
}

func TestRetry_Retry(t *testing.T) {
	testcases := map[string]struct {
		retry    utils.Retry
		f        func() (bool, error)
		expErr   error
		expPanic bool
	}{
		"FTrue": {
			retry: utils.RetryStandard,
			f:     func() (bool, error) { return true, nil },
		},
		"FFalse": {
			retry:  utils.RetryStandard,
			f:      func() (bool, error) { return false, utils.TestErrRetryErr },
			expErr: utils.TestErrRetryErr,
		},
	}
	for i, tc := range testcases {
		func() {
			defer func() {
				if r := recover(); r != nil && !tc.expPanic {
					t.Errorf("Test %s: unexpected panic occurred: %s\n%s", i, r, debug.Stack())
				}
			}()
			err := tc.retry.Retry(tc.f)
			if !errors.Is(err, tc.expErr) {
				t.Errorf("Test %s: Received error %#v; expected error%#v", i, err, tc.expErr)
			}
			if tc.expPanic {
				t.Errorf("Test %s: Unexpectidly did not panic", i)
			}
		}()
	}
}
