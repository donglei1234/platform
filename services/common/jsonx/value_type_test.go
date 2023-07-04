package jsonx

import (
	"testing"
)

func TestValueType(t *testing.T) {
	var testCases = []struct {
		valueType               ValueType
		expectedString          string
		expectedIsBasic         bool
		expectedCanHaveChildren bool
	}{
		{
			valueType:               ValueTypeUndefined,
			expectedString:          "undefined",
			expectedIsBasic:         true,
			expectedCanHaveChildren: false,
		},
		{
			valueType:               ValueTypeNull,
			expectedString:          "null",
			expectedIsBasic:         true,
			expectedCanHaveChildren: false,
		},
		{
			valueType:               ValueTypeTrue,
			expectedString:          "true",
			expectedIsBasic:         true,
			expectedCanHaveChildren: false,
		},
		{
			valueType:               ValueTypeFalse,
			expectedString:          "false",
			expectedIsBasic:         true,
			expectedCanHaveChildren: false,
		},
		{
			valueType:               ValueTypeString,
			expectedString:          "string",
			expectedIsBasic:         true,
			expectedCanHaveChildren: false,
		},
		{
			valueType:               ValueTypeNumber,
			expectedString:          "number",
			expectedIsBasic:         true,
			expectedCanHaveChildren: false,
		},
		{
			valueType:               ValueTypeObject,
			expectedString:          "object",
			expectedIsBasic:         false,
			expectedCanHaveChildren: true,
		},
		{
			valueType:               ValueTypeArray,
			expectedString:          "array",
			expectedIsBasic:         false,
			expectedCanHaveChildren: true,
		},
	}

	for _, tc := range testCases {
		if tc.valueType.String() != tc.expectedString {
			t.Error("Error: ValueType.toString() response did not match expected string value. Expected:",
				tc.expectedString, "| Actual:", tc.valueType.String())
		}

		if tc.valueType.IsBasic() != tc.expectedIsBasic {
			t.Error("Error: ValueType.IsBasic() response did not match expected bool value. Expected:",
				tc.expectedIsBasic, "| Actual:", tc.valueType.IsBasic())
		}

		if tc.valueType.CanHaveChildren() != tc.expectedCanHaveChildren {
			t.Error("Error: ValueType.CanHaveChildren() response did not match expected bool value. Expected:",
				tc.expectedCanHaveChildren, "| Actual:", tc.valueType.CanHaveChildren())
		}
	}
}
