package utils_test

import (
	"math"
	"runtime/debug"
	"testing"
	"time"

	"google.golang.org/grpc"

	"github.com/donglei1234/platform/services/common/utils"
)

func TestTransportSecurity(t *testing.T) {
	testcases := map[string]struct {
		secure    bool
		expResult grpc.DialOption
		expPanic  bool
	}{"Secure": {
		secure:    true,
		expResult: utils.SecureDialOption,
	}, "Insecure": {
		secure:    false,
		expResult: utils.InsecureDialOption,
	}}
	for i, tc := range testcases {
		func() {
			defer func() {
				if r := recover(); r != nil && !tc.expPanic {
					t.Errorf("Test %s: Unexpected Panic occurred: %s\n%s", i, r, debug.Stack())
				}
			}()
			if opt := utils.TransportSecurity(tc.secure); opt != tc.expResult {
				t.Errorf("Test %s: Received %v; expected %v", i, opt, tc.expResult)
			}
			if tc.expPanic {
				t.Errorf("Test %s: Expected Panic did not occur", i)
			}
		}()
	}
}

func TestSetBlockingDialTimeout(t *testing.T) {
	testcases := map[string]time.Duration{
		"ZeroDuration": time.Duration(0),
		"Nanosecond":   time.Nanosecond,
		"Microsecond":  time.Microsecond,
		"Millisecond":  time.Millisecond,
		"Second":       time.Second,
		"Minute":       time.Minute,
		"Hour":         time.Hour,
		"Minus1":       time.Duration(-1),
		"MaxInt64":     time.Duration(math.MaxInt64),
		"MinInt64":     time.Duration(math.MinInt64),
	}
	for i, tc := range testcases {
		if utils.SetBlockingDialTimeout(tc); *utils.BlockingDialTimeout != tc {
			t.Errorf("Test %s: BlockingDialTimeOut was %v; Expected %v", i, *utils.BlockingDialTimeout, tc)
		}
	}
}
