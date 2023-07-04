package utils_test

import (
	"github.com/donglei1234/platform/services/common/utils"
	"os"
	"runtime/debug"
	"testing"
)

const tmpFileName = "07ad83a8ba3b476a8a559b90bbbc7f60"

func TestFileExists(t *testing.T) {
	testcases := map[string]struct {
		statFunc  func(string) (os.FileInfo, error)
		file      string
		expResult bool
		expPanic  bool
	}{
		"StatErr": {
			statFunc: utils.FileExists_StatFuncErr,
			file:     "亜美",
		},
		"IsDirectory": {
			statFunc: utils.FileExists_StatFuncIsDir,
			file:     ".",
		},
		"IsFile": {
			statFunc:  utils.FileExists_StatFuncNotDir,
			file:      "",
			expResult: true,
		},
	}
	for i, tc := range testcases {
		//Save and restore the actual Stat function
		tmpStatFunc := *utils.Status
		defer func() { *utils.Status = tmpStatFunc }()
		func() {
			defer func() {
				if r := recover(); r != nil && !tc.expPanic {
					t.Errorf("Test %s: unexpected panic occurred: %s\n%s", i, r, debug.Stack())
				}
			}()
			*utils.Status = tc.statFunc
			if res := utils.FileExists(tc.file); res != tc.expResult {
				t.Errorf("Test %s: Received %#v; expected %#v", i, res, tc.expResult)
			}
			if tc.expPanic {
				t.Errorf("Test %s: Unexpectidly did not panic", i)
			}
		}()
	}
}
