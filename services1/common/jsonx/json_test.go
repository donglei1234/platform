package jsonx

import (
	"testing"
)

func TestParse(t *testing.T) {
	type TestCase struct {
		json        string
		errExpected bool
	}

	cases := []TestCase{
		{
			json:        `{"a":"b"}`,
			errExpected: false,
		},
		{
			json:        `{"a":{"b":{"c":true}}}`,
			errExpected: false,
		},
		{
			json:        `a`,
			errExpected: true,
		},
		{
			json:        ``,
			errExpected: true,
		},
	}
	var result interface{}
	for i, tc := range cases {
		if err := Parse([]byte(tc.json), &result); tc.errExpected && err != nil {
			continue
		} else if err != nil {
			t.Fatal("Error encountered in Parse() call in test case #", i+1, ":", err)
		} else if stringResult, err := Stringify(result); err != nil {
			t.Fatal("Error encountered in Stringify() call in test case #", i+1, ":", err)
		} else if stringResult == "" {
			t.Fatal("Error encountered in test case #", i+1, ": Result string cannot be nil / empty.")
		} else if stringResult != tc.json {
			t.Fatal("Result string does not match the provided test case string in test case #", i+1)
		}
	}
}

func TestParseString(t *testing.T) {
	type TestCase struct {
		json        string
		errExpected bool
	}

	cases := []TestCase{
		{
			json:        `{"a":"b"}`,
			errExpected: false,
		},
		{
			json:        `{"a":{"b":{"c":true}}}`,
			errExpected: false,
		},
		{
			json:        `a`,
			errExpected: true,
		},
		{
			json:        ``,
			errExpected: true,
		},
	}

	for i, tc := range cases {
		var v Value
		if err := ParseString(tc.json, &v); tc.errExpected == true && err != nil {
			continue
		} else if tc.errExpected == true && err == nil {
			t.Fatal("Error was not encountered where an error was expected in test case", i+1)
		} else if err != nil {
			t.Fatal("Error encountered in test case", i+1, ":", err)
		}
	}
}

func TestStringify(t *testing.T) {
	type TestCase struct {
		input, expectedOutput string
		errExpected           bool
	}

	cases := []TestCase{
		{
			input:          `{"a":"b"}`,
			expectedOutput: `{"a":"b"}`,
			errExpected:    false,
		},
		{
			input:       ``,
			errExpected: true,
			// empty input string
		},
	}

	for i, tc := range cases {
		var v Value
		if err := ParseString(tc.input, &v); tc.errExpected && err != nil {
			continue
		} else if err != nil {
			t.Fatal("Error encountered in test case", i+1, ":", err)
		} else {
			if out, err := Stringify(v); err != nil {
				t.Fatal("Error encountered in test case", i+1, ":", err)
			} else if out != tc.expectedOutput && tc.errExpected == true {
				continue
			} else if out != tc.expectedOutput && tc.errExpected == false {
				t.Fatal("Error encountered in test case", i+1,
					":", "Stringify() output", out, "differs from the expected output", tc.expectedOutput)
			}
		}
	}
}

func TestParseAndReturn(t *testing.T) {
	type testCase struct {
		json        []byte
		errExpected bool
	}
	cases := []testCase{
		{
			json:        []byte{50},
			errExpected: false,
		},
		{
			json:        []byte{116, 114, 117, 101},
			errExpected: false,
		},
		{
			json:        []byte{111, 98, 106, 49},
			errExpected: true,
		},
		{
			json:        []byte{},
			errExpected: true,
		},
	}

	for i, tc := range cases {
		if result, err := ParseAndReturn(tc.json); tc.errExpected == true && err != nil {
			continue
		} else if err == nil && tc.errExpected == true {
			t.Fatal("Error was not encountered where error was expected in test case", i+1)
		} else if result == nil {
			t.Fatal("JSON parse result cannot be nil.")
		} else if stringResult, err := Stringify(result); err != nil {
			t.Fatal("Error encountered in Stringify() call in test case", i+1, ":", err)
		} else if stringResult == "" {
			t.Fatal("Error encountered in test case #", i+1, ": Result string cannot be nil / empty.")
		} else if stringResult != string(tc.json) {
			t.Fatal("Result string does not match the provided test case string in test case #", i+1)
		}
	}
}

func TestParseStringAndReturn(t *testing.T) {
	type testCase struct {
		json        string
		errExpected bool
	}
	cases := []testCase{
		{
			json:        `{"a":"b"}`,
			errExpected: false,
		},
		{
			json:        `{"a":{"b":{"c":true}}}`,
			errExpected: false,
		},
		{
			json:        `{"a":{"b":{"c":{"d":{"e":true}}}}}`,
			errExpected: false,
		},
		{
			json:        `a`,
			errExpected: true,
		},
		{
			json:        ``,
			errExpected: true,
		},
		{
			json:        `{"a":"b"} a`,
			errExpected: true,
		},
		{
			json:        `The entire contents of the Encyclopedia Brittanica`,
			errExpected: true,
		},
	}

	for i, tc := range cases {
		if result, err := ParseStringAndReturn(tc.json); err != nil && tc.errExpected == true {
			continue
		} else if err == nil && tc.errExpected == true {
			t.Fatal("Error was not encountered where error was expected in test case", i+1)
		} else if result == nil {
			t.Fatal("JSON parse result cannot be nil.")
		} else if stringResult, err := Stringify(result); err != nil {
			t.Fatal("Error encountered in Stringify() call in test case", i+1, ":", err)
		} else if stringResult == "" {
			t.Fatal("Error encountered in test case #", i+1, ": Result string cannot be nil / empty.")
		} else if stringResult != string(tc.json) {
			t.Fatal("Result string does not match the provided test case string in test case #", i+1)
		}
	}
}
