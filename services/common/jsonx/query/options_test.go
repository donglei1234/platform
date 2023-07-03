package query

import (
	"fmt"
	"strings"
	"testing"
)

func TestWithArrayInsert(t *testing.T) {
	type testCase struct {
		data        interface{}
		value       []byte
		path        []byte
		expOutput   interface{}
		expectedErr error
	}

	testCases := []testCase{
		{
			path:  []byte(`a[1]`),
			value: []byte(`{"g":"h"}`),
			expOutput: map[string]interface{}{
				"a": []interface{}{map[string]interface{}{"b": "c"}, map[string]interface{}{"g": "h"}},
				"d": []interface{}{map[string]interface{}{"e": "f"}},
			},
			data: map[string]interface{}{
				"a": []interface{}{map[string]interface{}{"b": "c"}},
				"d": []interface{}{map[string]interface{}{"e": "f"}},
			},
			expectedErr: nil,
		},
		{
			path:  []byte(`d[1]`),
			value: []byte(`{"g":"h"}`),
			expOutput: map[string]interface{}{
				"a": []interface{}{map[string]interface{}{"b": "c"}},
				"d": []interface{}{map[string]interface{}{"e": "f"}, map[string]interface{}{"g": "h"}},
			},
			data: map[string]interface{}{
				"a": []interface{}{map[string]interface{}{"b": "c"}},
				"d": []interface{}{map[string]interface{}{"e": "f"}},
			},
			expectedErr: nil,
		},
		{
			path:  []byte(`.`),
			value: []byte(`{"g":"h"}`),
			data: map[string]interface{}{
				"a": []interface{}{map[string]interface{}{"b": "c"}},
				"d": []interface{}{map[string]interface{}{"e": "f"}},
			},
			expectedErr: ErrInvalidDestination,
		},
		{
			path:  nil,
			value: []byte(`{"g":"h"}`),
			data: map[string]interface{}{
				"a": []interface{}{map[string]interface{}{"b": "c"}},
				"d": []interface{}{map[string]interface{}{"e": "f"}},
			},
			expectedErr: ErrInvalidPath,
		},
	}
	for i, tc := range testCases {
		var err error
		var opts options
		WithArrayInsert(tc.path, tc.value)(&opts)
		for _, a := range opts.actions {
			if err = a(tc.data); err != nil {
				if tc.expectedErr != err {
					t.Error("Unexpected error executing the provided option action in test case #", i+1, ":", err)
					break
				}
			} else {
				if s1, s2 := fmt.Sprintf("%+v", tc.data), fmt.Sprintf("%+v", tc.expOutput); strings.Compare(s1, s2) != 0 {
					t.Errorf("Test #%d: got %+v; expected %+v", i+1, s1, s2)
				}
			}
			if opts.input != nil {
				t.Error("Input option is not nil when it was not provided in test case #", i+1)
			}
		}
	}
}

func TestWithArrayPushFront(t *testing.T) {
	type testCase struct {
		path        []byte
		value       []byte
		data        interface{}
		expOutput   interface{}
		expectedErr error
	}

	testCases := []testCase{
		{
			path:  []byte(`a`),
			value: []byte(`{"c":"d"}`),
			data: map[string]interface{}{
				"a": []interface{}{map[string]interface{}{"b": "c"}},
			},
			expOutput: map[string]interface{}{
				"a": []interface{}{map[string]interface{}{"c": "d"}, map[string]interface{}{"b": "c"}}},
			expectedErr: nil,
		},
		{
			path:  []byte(`.`),
			value: []byte(`{"g":"h"}`),
			data: map[string]interface{}{
				"a": []interface{}{map[string]interface{}{"b": "c"}},
				"d": []interface{}{map[string]interface{}{"e": "f"}},
			},
			expectedErr: ErrInvalidDestination,
		},
		{
			path:  nil,
			value: []byte(`{"g":"h"}`),
			data: map[string]interface{}{
				"a": []interface{}{map[string]interface{}{"b": "c"}},
				"d": []interface{}{map[string]interface{}{"e": "f"}},
			},
			expectedErr: ErrInvalidPath,
		},
	}

	for i, tc := range testCases {
		var err error
		var opts options
		WithArrayPushFront(tc.path, tc.value)(&opts)
		for _, a := range opts.actions {
			if err = a(tc.data); err != nil {
				if tc.expectedErr != err {
					t.Error("Unexpected error executing the provided option action in test case #", i+1, ":", err)
					break
				}
			} else {
				if s1, s2 := fmt.Sprintf("%+v", tc.data), fmt.Sprintf("%+v", tc.expOutput); strings.Compare(s1, s2) != 0 {
					t.Errorf("Test #%d: got %+v; expected %+v", i+1, s1, s2)
				}
			}
		}

		if opts.input != nil {
			t.Error("Input option is not nil when it was not provided in test case #", i+1)
		}
	}
}

func TestWithArrayPushBack(t *testing.T) {
	type testCase struct {
		path        []byte
		value       []byte
		data        interface{}
		expOutput   interface{}
		expectedErr error
	}

	testCases := []testCase{
		{
			path:  []byte(`a`),
			value: []byte(`{"c":"d"}`),
			data: map[string]interface{}{
				"a": []interface{}{map[string]interface{}{"b": "c"}},
			},
			expOutput: map[string]interface{}{
				"a": []interface{}{map[string]interface{}{"b": "c"}, map[string]interface{}{"c": "d"}},
			},
			expectedErr: nil,
		},
		{
			path:  []byte(`.`),
			value: []byte(`{"g":"h"}`),
			data: map[string]interface{}{
				"a": []interface{}{map[string]interface{}{"b": "c"}},
				"d": []interface{}{map[string]interface{}{"e": "f"}},
			},
			expectedErr: ErrInvalidDestination,
		},
		{
			path:  nil,
			value: []byte(`{"g":"h"}`),
			data: map[string]interface{}{
				"a": []interface{}{map[string]interface{}{"b": "c"}},
				"d": []interface{}{map[string]interface{}{"e": "f"}},
			},
			expectedErr: ErrInvalidPath,
		},
	}
	for i, tc := range testCases {
		var err error
		var opts options
		WithArrayPushBack(tc.path, tc.value)(&opts)
		for _, a := range opts.actions {
			if err = a(tc.data); err != nil {
				if tc.expectedErr != err {
					t.Error("Unexpected error executing the provided option action in test case #", i+1, ":", err)
					break
				}
			} else {
				if s1, s2 := fmt.Sprintf("%+v", tc.data), fmt.Sprintf("%+v", tc.expOutput); strings.Compare(s1, s2) != 0 {
					t.Errorf("Test #%d: got %+v; expected %+v", i+1, s1, s2)
				}
			}
			if opts.input != nil {
				t.Error("Input option is not nil when it was not provided in test case #", i+1)
			}
		}
	}
}

func TestWithCopy(t *testing.T) {
	type testCase struct {
		pathFrom    []byte
		pathTo      []byte
		data        interface{}
		expOutput   interface{}
		expectedErr error
	}

	testCases := []testCase{
		{
			pathFrom: []byte(`a`),
			pathTo:   []byte(`b`),
			data: map[string]interface{}{
				"a": map[string]interface{}{"c": "d"},
				"b": map[string]interface{}{},
			},
			expOutput: map[string]interface{}{
				"a": map[string]interface{}{"c": "d"},
				"b": map[string]interface{}{"c": "d"},
			},
			expectedErr: nil,
		},
		{
			pathFrom: []byte(`b`),
			pathTo:   []byte(`a`),
			data: map[string]interface{}{
				"a": map[string]interface{}{},
				"b": map[string]interface{}{"c": "d"},
			},
			expOutput: map[string]interface{}{
				"a": map[string]interface{}{"c": "d"},
				"b": map[string]interface{}{"c": "d"},
			},
			expectedErr: nil,
		},
		{
			pathFrom: []byte(`a`),
			pathTo:   []byte(`b`),
			data: map[string]interface{}{
				"a": []interface{}{map[string]interface{}{"e": "f"}},
				"b": []interface{}{map[string]interface{}{}},
			},
			expOutput: map[string]interface{}{
				"a": []interface{}{map[string]interface{}{"e": "f"}},
				"b": []interface{}{map[string]interface{}{"e": "f"}},
			},
			expectedErr: nil,
		},
		{
			pathFrom: []byte(`b`),
			pathTo:   []byte(`a`),
			data: map[string]interface{}{
				"a": []interface{}{map[string]interface{}{}},
				"b": []interface{}{map[string]interface{}{"e": "f"}},
			},
			expOutput: map[string]interface{}{
				"a": []interface{}{map[string]interface{}{"e": "f"}},
				"b": []interface{}{map[string]interface{}{"e": "f"}},
			},
			expectedErr: nil,
		},
		{
			// Copying from a nonexistent path wipes the interface - do we want to handle it this way???
			pathFrom: []byte(`c`),
			pathTo:   []byte(`a`),
			data: map[string]interface{}{
				"a": []interface{}{map[string]interface{}{}},
				"b": []interface{}{map[string]interface{}{"e": "f"}},
			},
			expOutput: map[string]interface{}{
				"a": nil,
				"b": []interface{}{map[string]interface{}{"e": "f"}},
			},
			expectedErr: nil,
		},
		{
			pathFrom: []byte(`a`),
			pathTo:   []byte(`.`),
			data: map[string]interface{}{
				"a": []interface{}{map[string]interface{}{}},
				"b": []interface{}{map[string]interface{}{"e": "f"}},
			},
			expOutput: map[string]interface{}{
				"a": []interface{}{map[string]interface{}{}},
				"b": []interface{}{map[string]interface{}{"e": "f"}},
			},
			expectedErr: ErrInvalidDestination,
		},
		{
			pathFrom: nil,
			pathTo:   []byte(`.`),
			data: map[string]interface{}{
				"a": []interface{}{map[string]interface{}{}},
				"b": []interface{}{map[string]interface{}{"e": "f"}},
			},
			expOutput: map[string]interface{}{
				"a": []interface{}{map[string]interface{}{}},
				"b": []interface{}{map[string]interface{}{"e": "f"}},
			},
			expectedErr: ErrInvalidPath,
		},
		{
			pathFrom: []byte(`a`),
			pathTo:   nil,
			data: map[string]interface{}{
				"a": []interface{}{map[string]interface{}{}},
				"b": []interface{}{map[string]interface{}{"e": "f"}},
			},
			expOutput: map[string]interface{}{
				"a": []interface{}{map[string]interface{}{}},
				"b": []interface{}{map[string]interface{}{"e": "f"}},
			},
			expectedErr: ErrInvalidPath,
		},
		{
			pathFrom: []byte(``),
			pathTo:   []byte(`.`),
			data: map[string]interface{}{
				"a": []interface{}{map[string]interface{}{}},
				"b": []interface{}{map[string]interface{}{"e": "f"}},
			},
			expOutput: map[string]interface{}{
				"a": []interface{}{map[string]interface{}{}},
				"b": []interface{}{map[string]interface{}{"e": "f"}},
			},
			expectedErr: ErrInvalidPath,
		},
		{
			pathFrom: []byte(`a`),
			pathTo:   []byte(``),
			data: map[string]interface{}{
				"a": []interface{}{map[string]interface{}{}},
				"b": []interface{}{map[string]interface{}{"e": "f"}},
			},
			expOutput: map[string]interface{}{
				"a": []interface{}{map[string]interface{}{}},
				"b": []interface{}{map[string]interface{}{"e": "f"}},
			},
			expectedErr: ErrInvalidPath,
		},
	}
	for i, tc := range testCases {
		var err error
		var opts options
		WithCopy(tc.pathFrom, tc.pathTo)(&opts)
		for _, a := range opts.actions {
			if err = a(tc.data); err != nil {
				if tc.expectedErr != err {
					t.Error("Unexpected error executing the provided option action in test case #", i+1, ":", err)
					break
				}
			} else {
				if s1, s2 := fmt.Sprintf("%+v", tc.data), fmt.Sprintf("%+v", tc.expOutput); strings.Compare(s1, s2) != 0 {
					t.Errorf("Test #%d: got %+v; expected %+v", i+1, s1, s2)
				}
			}
			if opts.input != nil {
				t.Error("Input option is not nil when it was not provided in test case #", i+1)
			}
		}
	}
}

func TestWithMove(t *testing.T) {
	type testCase struct {
		pathFrom    []byte
		pathTo      []byte
		data        interface{}
		expOutput   interface{}
		expectedErr error
	}

	testCases := []testCase{
		{
			pathFrom: []byte(`a`),
			pathTo:   []byte(`b`),
			data: map[string]interface{}{
				"a": map[string]interface{}{"c": "d"},
				"b": map[string]interface{}{},
			},
			expOutput: map[string]interface{}{
				"b": map[string]interface{}{"c": "d"},
			},
			expectedErr: nil,
		},
		{
			pathFrom: []byte(`b`),
			pathTo:   []byte(`a`),
			data: map[string]interface{}{
				"a": []interface{}{map[string]interface{}{}},
				"b": []interface{}{map[string]interface{}{"e": "f"}},
			},
			expOutput: map[string]interface{}{
				"a": []interface{}{map[string]interface{}{"e": "f"}},
			},
			expectedErr: nil,
		},
		{
			pathFrom: []byte(`b`),
			pathTo:   []byte(`c`),
			data: map[string]interface{}{
				"a": []interface{}{map[string]interface{}{}},
				"b": []interface{}{map[string]interface{}{"e": "f"}},
			},
			expOutput: map[string]interface{}{
				"a": []interface{}{map[string]interface{}{}},
				"c": []interface{}{map[string]interface{}{"e": "f"}},
			},
			expectedErr: nil,
		},
		{
			pathFrom: []byte(`b`),
			pathTo:   []byte(`.`),
			data: map[string]interface{}{
				"a": []interface{}{map[string]interface{}{}},
				"b": []interface{}{map[string]interface{}{"e": "f"}},
			},
			expOutput: map[string]interface{}{
				"a": []interface{}{map[string]interface{}{}},
				"c": []interface{}{map[string]interface{}{"e": "f"}},
			},
			expectedErr: ErrInvalidDestination,
		},
		{
			pathFrom: nil,
			pathTo:   []byte(`.`),
			data: map[string]interface{}{
				"a": []interface{}{map[string]interface{}{}},
				"b": []interface{}{map[string]interface{}{"e": "f"}},
			},
			expOutput: map[string]interface{}{
				"a": []interface{}{map[string]interface{}{}},
				"b": []interface{}{map[string]interface{}{"e": "f"}},
			},
			expectedErr: ErrInvalidPath,
		},
		{
			pathFrom: []byte(`a`),
			pathTo:   nil,
			data: map[string]interface{}{
				"a": []interface{}{map[string]interface{}{}},
				"b": []interface{}{map[string]interface{}{"e": "f"}},
			},
			expOutput: map[string]interface{}{
				"a": []interface{}{map[string]interface{}{}},
				"b": []interface{}{map[string]interface{}{"e": "f"}},
			},
			expectedErr: ErrInvalidPath,
		},
		{
			pathFrom: []byte(``),
			pathTo:   []byte(`.`),
			data: map[string]interface{}{
				"a": []interface{}{map[string]interface{}{}},
				"b": []interface{}{map[string]interface{}{"e": "f"}},
			},
			expOutput: map[string]interface{}{
				"a": []interface{}{map[string]interface{}{}},
				"b": []interface{}{map[string]interface{}{"e": "f"}},
			},
			expectedErr: ErrInvalidPath,
		},
		{
			pathFrom: []byte(`a`),
			pathTo:   []byte(``),
			data: map[string]interface{}{
				"a": []interface{}{map[string]interface{}{}},
				"b": []interface{}{map[string]interface{}{"e": "f"}},
			},
			expOutput: map[string]interface{}{
				"a": []interface{}{map[string]interface{}{}},
				"b": []interface{}{map[string]interface{}{"e": "f"}},
			},
			expectedErr: ErrInvalidPath,
		},
	}
	for i, tc := range testCases {
		var err error
		var opts options
		WithMove(tc.pathFrom, tc.pathTo)(&opts)
		for _, a := range opts.actions {
			if err = a(tc.data); err != nil {
				if tc.expectedErr != err {
					t.Error("Unexpected error executing the provided option action in test case #", i+1, ":", err)
					break
				}
			} else {
				if s1, s2 := fmt.Sprintf("%+v", tc.data), fmt.Sprintf("%+v", tc.expOutput); strings.Compare(s1, s2) != 0 {
					t.Errorf("Test #%d: got %+v; expected %+v", i+1, s1, s2)
				}
			}
			if opts.input != nil {
				t.Error("Input option is not nil when it was not provided in test case #", i+1)
			}
		}
	}
}

func TestWithDelete(t *testing.T) {
	type testCase struct {
		path        []byte
		data        interface{}
		expOutput   interface{}
		expectedErr error
	}

	testCases := []testCase{
		{
			path: []byte(`b`),
			data: map[string]interface{}{
				"a": map[string]interface{}{"c": "d"},
				"b": map[string]interface{}{"e": "f"},
			},
			expOutput: map[string]interface{}{
				"a": map[string]interface{}{"c": "d"},
			},
			expectedErr: nil,
		},
		{
			path: []byte(`a`),
			data: map[string]interface{}{
				"a": []interface{}{map[string]interface{}{"c": "d"}},
				"b": []interface{}{map[string]interface{}{"e": "f"}},
			},
			expOutput: map[string]interface{}{
				"b": []interface{}{map[string]interface{}{"e": "f"}},
			},
			expectedErr: nil,
		},
		{
			path: []byte(`c`),
			data: map[string]interface{}{
				"a": []interface{}{map[string]interface{}{"c": "d"}},
				"b": []interface{}{map[string]interface{}{"e": "f"}},
			},
			expOutput: map[string]interface{}{
				"a": []interface{}{map[string]interface{}{"c": "d"}},
				"b": []interface{}{map[string]interface{}{"e": "f"}},
			},
			expectedErr: nil,
		},
		{
			path: []byte(`.`),
			data: map[string]interface{}{
				"a": []interface{}{map[string]interface{}{"c": "d"}},
				"b": []interface{}{map[string]interface{}{"e": "f"}},
			},
			expOutput:   map[string]interface{}{},
			expectedErr: ErrInvalidDestination,
		},
		{
			path: []byte(``),
			data: map[string]interface{}{
				"a": []interface{}{map[string]interface{}{"c": "d"}},
				"b": []interface{}{map[string]interface{}{"e": "f"}},
			},
			expOutput: map[string]interface{}{
				"a": []interface{}{map[string]interface{}{"c": "d"}},
				"b": []interface{}{map[string]interface{}{"e": "f"}},
			},
			expectedErr: ErrInvalidPath,
		},
		{
			path: nil,
			data: map[string]interface{}{
				"a": []interface{}{map[string]interface{}{"c": "d"}},
				"b": []interface{}{map[string]interface{}{"e": "f"}},
			},
			expOutput: map[string]interface{}{
				"a": []interface{}{map[string]interface{}{"c": "d"}},
				"b": []interface{}{map[string]interface{}{"e": "f"}},
			},
			expectedErr: ErrInvalidPath,
		},
	}
	for i, tc := range testCases {
		var err error
		var opts options
		WithDelete(tc.path)(&opts)
		for _, a := range opts.actions {
			if err = a(tc.data); err != nil {
				if tc.expectedErr != err {
					t.Error("Unexpected error executing the provided option action in test case #", i+1, ":", err)
					break
				}
			} else {
				if s1, s2 := fmt.Sprintf("%+v", tc.data), fmt.Sprintf("%+v", tc.expOutput); strings.Compare(s1, s2) != 0 {
					t.Errorf("Test #%d: got %+v; expected %+v", i+1, s1, s2)
				}
			}
			if opts.input != nil {
				t.Error("Input option is not nil when it was not provided in test case #", i+1)
			}
		}
	}
}

func TestWithSwap(t *testing.T) {
	type testCase struct {
		pathFrom    []byte
		pathTo      []byte
		data        interface{}
		expOutput   interface{}
		expectedErr error
	}

	testCases := []testCase{
		{
			pathFrom: []byte(`a`),
			pathTo:   []byte(`b`),
			data: map[string]interface{}{
				"a": map[string]interface{}{"c": "d"},
				"b": map[string]interface{}{"e": "f"},
			},
			expOutput: map[string]interface{}{
				"a": map[string]interface{}{"e": "f"},
				"b": map[string]interface{}{"c": "d"},
			},
			expectedErr: nil,
		},
		{
			pathFrom: []byte(`a`),
			pathTo:   []byte(`b`),
			data: map[string]interface{}{
				"a": []interface{}{map[string]interface{}{"c": "d"}},
				"b": []interface{}{map[string]interface{}{"e": "f"}},
			},
			expOutput: map[string]interface{}{
				"a": []interface{}{map[string]interface{}{"e": "f"}},
				"b": []interface{}{map[string]interface{}{"c": "d"}},
			},
			expectedErr: nil,
		},
		{
			pathFrom: []byte(`c`),
			pathTo:   []byte(`a`),
			data: map[string]interface{}{
				"a": []interface{}{map[string]interface{}{"c": "d"}},
				"b": []interface{}{map[string]interface{}{"e": "f"}},
			},
			expOutput: map[string]interface{}{
				"a": nil,
				"b": []interface{}{map[string]interface{}{"e": "f"}},
				"c": []interface{}{map[string]interface{}{"c": "d"}},
			},
			expectedErr: nil,
		},
		{
			pathFrom: []byte(`a`),
			pathTo:   []byte(`.`),
			data: map[string]interface{}{
				"a": []interface{}{map[string]interface{}{"c": "d"}},
				"b": []interface{}{map[string]interface{}{"e": "f"}},
			},
			expectedErr: ErrInvalidDestination,
		},
	}
	for i, tc := range testCases {
		var err error
		var opts options
		WithSwap(tc.pathFrom, tc.pathTo)(&opts)
		for _, a := range opts.actions {
			if err = a(tc.data); err != nil {
				if tc.expectedErr != err {
					t.Error("Unexpected error executing the provided option action in test case #", i+1, ":", err)
					break
				}
			} else {
				if s1, s2 := fmt.Sprintf("%+v", tc.data), fmt.Sprintf("%+v", tc.expOutput); strings.Compare(s1, s2) != 0 {
					t.Errorf("Test #%d: got %+v; expected %+v", i+1, s1, s2)
				}
			}
			if opts.input != nil {
				t.Error("Input option is not nil when it was not provided in test case #", i+1)
			}
		}
	}
}

func TestWithInput(t *testing.T) {
	type testCase struct {
		input     interface{}
		expOutput interface{}
	}
	testCases := []testCase{
		{
			input: map[string]interface{}{
				"a": map[string]interface{}{"b": "c"},
			},
			expOutput: map[string]interface{}{
				"a": map[string]interface{}{"b": "c"},
			},
		},
		{
			input: map[string]interface{}{
				"a": []interface{}{map[string]interface{}{"b": "c"}},
			},
			expOutput: map[string]interface{}{
				"a": []interface{}{map[string]interface{}{"b": "c"}},
			},
		},
		{
			input: map[string]interface{}{
				"": []interface{}{map[string]interface{}{"": ""}},
			},
			expOutput: map[string]interface{}{
				"": []interface{}{map[string]interface{}{"": ""}},
			},
		},
		{
			input:     nil,
			expOutput: nil,
		},
	}
	for i, tc := range testCases {
		var opts options
		WithInput(tc.input)(&opts)
		if s1, s2 := fmt.Sprintf("%+v", tc.input), fmt.Sprintf("%+v", tc.expOutput); strings.Compare(s1, s2) != 0 {
			t.Errorf("Test #%d: got %+v; expected %+v", i+1, s1, s2)
		}

		if len(opts.actions) != 0 {
			t.Errorf("Test #%d: `actions` was not nil", i+1)
		}
	}
}

func TestNewOptions(t *testing.T) {
	type testCase struct {
		path      []byte
		value     []byte
		data      interface{}
		expOutput interface{}
	}
	testCases := []testCase{
		{
			path:  []byte(`a`),
			value: []byte(`b`),
			data: map[string]interface{}{
				"a": map[string]interface{}{"b": "c"},
			},
			expOutput: map[string]interface{}{},
		},
	}

	for i, tc := range testCases {
		opts := newOptions(WithDelete(tc.path))
		var err error
		for _, a := range opts.actions {
			if err = a(tc.data); err != nil {
				t.Error("Error executing the provided option action in test case #", i+1, ":", err)
				break
			} else {
				if s1, s2 := fmt.Sprintf("%+v", tc.data), fmt.Sprintf("%+v", tc.expOutput); strings.Compare(s1, s2) != 0 {
					t.Errorf("Test #%d: got %+v; expected %+v", i+1, s1, s2)
				}
			}
		}
		if opts.input != nil {
			t.Error("Input option is not nil when it was not provided in test case #", i+1)
		}
	}
}
