package utils_test

import (
	"reflect"
	"runtime/debug"
	"strings"
	"testing"

	"github.com/donglei1234/platform/services/common/utils"
)

func TestShallowClone(t *testing.T) {
	testcases := map[string]struct {
		in       interface{}
		expPanic bool
	}{"NilTest": {
		in:       nil,
		expPanic: true,
	}, "StringTest": {
		in:       "string",
		expPanic: false,
	}, "StructTest": {
		in: struct {
			int
			string
			float32
		}{
			int:     34,
			string:  "stringTest",
			float32: 1.1,
		},
		expPanic: false,
	}, "IntTest": {
		in:       5,
		expPanic: false,
	}, "TestStruct": {
		in:       utils.TestStruct{Name: "Name", Age: 34},
		expPanic: false,
	}, "TestInterface": {
		in:       utils.TestInterface(utils.TestStruct{Name: "Name", Age: 34}),
		expPanic: false,
	}}
	for i, tc := range testcases {
		func() {
			defer func() {
				if r := recover(); r != nil && !tc.expPanic {
					t.Errorf("Test %s: panic occurred: %v\n%s", i, r, debug.Stack())
				}
			}()
			if res := utils.ShallowClone(tc.in); res == tc.in {
				tc.in = "newValue"
				if res == tc.in {
					t.Errorf("Test %s: Received the same object back, expected a clone of the object", i)
				}
			} else {
				v1, v2 := reflect.ValueOf(res).String(), reflect.ValueOf(tc.in).String()
				if strings.Compare(v1, v2) != 0 {
					t.Errorf("test %s: received %#v; expected %#v", i, res, tc.in)
				}
			}
			if tc.expPanic {
				t.Errorf("Test %s: Expected panic did not occur", i)
			}
		}()
	}
}
