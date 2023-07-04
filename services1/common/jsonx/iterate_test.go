package jsonx

import (
	"strings"
	"testing"
)

func TestIterate(t *testing.T) {
	type testCase struct {
		data              []byte
		errExpected       bool
		expectedErrString string
	}
	testCases := []testCase{
		// Success cases
		{
			data:        []byte(`{"Test": "Name", "Data": "TestData"}`),
			errExpected: false,
		},
		{
			data:        []byte(`{"More": "Data"}`),
			errExpected: false,
		},
		{
			data:        []byte(`{"Message": "Hello", "Array": [1, 2, 3], "Null": null, "Number": 1.234}`),
			errExpected: false,
		},
		{
			data:        []byte(`{"Message": "Hello", "Array": [1, 2, 3], "ValueReal": true, "Number": 1.234}`),
			errExpected: false,
		},
		{
			data:        []byte(`{"Message": "Hello", "Array": [1, 2, 3], "ValueReal": false, "Number": 1.234}`),
			errExpected: false,
		},
		{
			data:        []byte(`{"Message": "Hello World/\\"}`),
			errExpected: false,
		},
		{
			data:        []byte(`"Test": "Name", "Data": "TestData"`),
			errExpected: false,
		},
		{
			data:        []byte("[]"),
			errExpected: false,
		},
		{
			data:        []byte(`[{}]`),
			errExpected: false,
		},

		// Error cases
		{
			data:              []byte(`[{"Test": "Name", "Data": "TestData"}`),
			errExpected:       true,
			expectedErrString: "']' expected",
		},
		{
			data:              []byte(`{"Test": "Name", "Data": "TestData"`),
			errExpected:       true,
			expectedErrString: "'}' expected",
		},
		{
			data:              []byte(`[{"Test": "Name, "Data": "TestData"}]`),
			errExpected:       true,
			expectedErrString: `" expected`,
		},
		{
			data:              []byte(`[{"Message": "Hello", "Array": [1, 2, 3], "ValueReal": , "Number": 1.234}]`),
			errExpected:       true,
			expectedErrString: "value expected",
		},
		{
			data:              []byte(`[{"Test" "Name"}]`),
			errExpected:       true,
			expectedErrString: "':' expected",
		},
		{
			data:              []byte(`[{"Test": "Name" "Value": "Object"}]`),
			errExpected:       true,
			expectedErrString: "',' or '}' expected",
		},
		{
			data:              []byte(`[{"Test": "Name/\/"}]`),
			errExpected:       true,
			expectedErrString: "invalid escape sequence",
		},
	}

	iteratedArrays := 0
	iteratedStrings := 0
	iteratedObjects := 0

	// Number of expected elements to be iterated over. Change as new test cases are added.
	expectedArrays := 9
	expectedStrings := 10
	expectedObjects := 10

	for i, tc := range testCases {
		if err := Iterate(tc.data, func(
			depth int,
			path []byte,
			key []byte,
			value []byte,
			valueType ValueType,
		) IterateCallbackResult {
			if depth < 0 {
				t.Error("Iterated depth was unexpectedly less than 0 in test case #", i+1)
			}
			if valueType == ValueTypeArray {
				iteratedArrays++
			} else if valueType == ValueTypeString {
				iteratedStrings++
			} else if valueType == ValueTypeObject {
				iteratedObjects++
			} else if valueType == ValueTypeUndefined {
				t.Error("Encountered unknown value type", valueType, "in test case #", i+1)
			}
			return IterateCallbackContinue
		}); err != nil {
			if tc.errExpected {
				if !strings.Contains(err.Error(), tc.expectedErrString) {
					t.Error("Expected error was not encountered in test case #", i+1,
						"- error encountered:", err)
				}
			} else {
				t.Error("Unexpected error encountered iterating over the data in test case #", i+1, ":", err)
			}
		} else if tc.errExpected {
			t.Error("Error was not encountered where one was expected in test case #", i+1)
		}
	}
	if iteratedStrings != expectedStrings {
		t.Error("The number of strings iterated over does not match the expected number of strings. Iterated:",
			iteratedStrings, "Expected:", expectedStrings)
	} else if iteratedArrays != expectedArrays {
		t.Error("The number of arrays iterated over does not match the expected number of arrays. Iterated:",
			iteratedArrays, "Expected:", expectedArrays)
	} else if iteratedObjects != expectedObjects {
		t.Error("The number of objects iterated over does not match the expected number of objects. Iterated:",
			iteratedObjects, "Expected:", expectedObjects)
	}
}
