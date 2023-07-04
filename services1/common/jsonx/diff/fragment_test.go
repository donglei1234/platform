package diff

import (
	"testing"

	"github.com/donglei1234/platform/services/common/jsonx"
)

func TestFragment(t *testing.T) {
	var testCases = []struct {
		id             int
		depth          int
		path           string
		value          []byte
		valueType      jsonx.ValueType
		expectedString string
	}{
		// test case 1
		{
			id:             1,
			depth:          1,
			path:           "[1]",
			value:          []byte{123, 123, 123, 123, 123, 123, 123, 123},
			valueType:      jsonx.ValueTypeObject,
			expectedString: "1: object cf43908fc92aedd3: [1]",
		},

		// test case 2
		{
			id:             2,
			depth:          2,
			path:           "[2].id",
			value:          []byte{50},
			valueType:      jsonx.ValueTypeNumber,
			expectedString: "2: number 6021b5621680598b: [2].id = 2",
		},

		// test case 3
		{
			id:    3,
			depth: 3,
			path:  "object1.id",
			// Below is equivalent to: []byte{byte('"'), byte('o'), byte('b'), byte('j'), byte('1'), byte('"')},
			value:          []byte{111, 98, 106, 49},
			valueType:      jsonx.ValueTypeString,
			expectedString: "3: string e82b7ef1c9b541df: object1.id = obj1",
		},

		// test case 4
		{
			id:    4,
			depth: 4,
			path:  "car",
			// Below is equivalent to: []byte{byte('n'), byte('u'), byte('l'), byte('l')},
			value:          []byte{110, 117, 108, 108},
			valueType:      jsonx.ValueTypeNull,
			expectedString: "4: null 3ec9e10063179f3a: car",
		},

		// test case 5
		{
			id:    5,
			depth: 5,
			path:  "value",
			// Below is equivalent to: []byte{byte('t'), byte('r'), byte('u'), byte('e')},
			value:          []byte{116, 114, 117, 101},
			valueType:      jsonx.ValueTypeTrue,
			expectedString: "5: true d7c9b97948142e4a: value",
		},

		// test case 6
		{
			id:    6,
			depth: 6,
			path:  "value",
			// Below is equivalent to: []byte{byte('f'), byte('a'), byte('l'), byte('s'), byte('e')},
			value:          []byte{102, 97, 108, 115, 101},
			valueType:      jsonx.ValueTypeFalse,
			expectedString: "6: false 6d3f99ccc0c03a7a: value",
		},

		// test case 7
		{
			id:    7,
			depth: 7,
			path:  "",
			value: []byte{91, 123, 34, 105, 100, 34,
				58, 49, 125, 44, 32, 123, 34,
				105, 100, 34, 58, 50, 125, 93},
			valueType:      jsonx.ValueTypeArray,
			expectedString: "7: array ac97e06cd6378d14: ",
		},
	}
	for i, tc := range testCases {
		// Create a new fragment
		fragment := NewFragment(tc.id, tc.depth, tc.path, tc.value, tc.valueType)
		// Check if the fragment is equal to itself
		if fragment.IsEqual(&fragment) == false {
			t.Error("Error: A Fragment is not IsEqual() to itself. Issue encountered in test case #", i+1)
		}
		// Convert the fragment into a string
		returnedString := fragment.String()
		if tc.expectedString != returnedString {
			t.Error("Error: The returned string is not identical to the expected string. Issue encountered in test case #", i+1)
		}
	}
}
