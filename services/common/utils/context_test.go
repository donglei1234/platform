package utils_test

import (
	"context"
	"errors"
	"runtime/debug"
	"testing"

	"github.com/donglei1234/platform/services/common/utils"
)

func TestUntilContextDone(t *testing.T) {
	testcases := map[string]struct {
		context.Context
		f        func() error
		expErr   error
		expPanic bool
	}{"FunctionReturnsError": {
		f:      func() error { return utils.TestErrFunc },
		expErr: utils.TestErrFunc,
	}, "ChannelChecks": {
		Context: &utils.TickingContext,
		f:       func() error { return nil },
		expErr:  utils.TestErrContext,
	}, "NilFunction": {
		Context:  nil,
		f:        nil,
		expPanic: true,
	}, "NilContext": {
		Context:  nil,
		f:        func() error { return nil },
		expPanic: true,
	}}
	for i, tc := range testcases {
		func() {
			defer func() {
				if r := recover(); r != nil && !tc.expPanic {
					t.Errorf("Test %s: Unexpected Panic occurred: %v\n%s", i, r, debug.Stack())
				}
			}()
			if err := utils.UntilContextDone(tc.Context, tc.f); !errors.Is(err, tc.expErr) {
				t.Errorf("Test %s: Received error %v; expected error %v", i, err, tc.expErr)
			}
		}()
	}
}
