package jsonx

import (
	"fmt"
	"strings"
	"testing"
)

func TestMarshal(t *testing.T) {
	type testCase struct {
		data interface{}
	}

	cases := []testCase{
		{data: "test"},
		{data: ""},
		{data: []string{"test", "data"}},
		{data: []byte("test")},
		{data: 100},
		{data: 3.14},
		//{data: 0.0000000000000123},  // TODO: @PBAUMAN: Issue #993
		//{data: -0.0000000000000123}, // TODO: @PBAUMAN: Issue #993
		{data: true},
		{data: false},
	}

	for i, tc := range cases {

		if result, err := Marshal(tc.data); err != nil {
			t.Fatal("Could not marshal data in test case #", i+1, ":", err)
		} else {
			switch tc.data.(type) {
			case string:
				var unmarshalResultString string
				if err := Unmarshal(result, &unmarshalResultString); err != nil {
					t.Fatal("Couldn't unmarshal data in test case #", i+1, ":", err)
				} else if unmarshalResultString != string(`"`+tc.data.(string)+`"`) {
					t.Fatal("Result does not match provided data in test case", i+1, "; provided data:", tc.data,
						"; result data:", unmarshalResultString)
				}
			case []string:
				var unmarshalResultSliceString []string
				if err := Unmarshal(result, &unmarshalResultSliceString); err != nil {
					t.Fatal("Couldn't unmarshal data in test case #", i+1, ":", err)
				} else {
					for resultIndex, resultPiece := range unmarshalResultSliceString {
						if !strings.Contains(tc.data.([]string)[resultIndex], resultPiece) {
							t.Fatal("Result does not match provided data in test case", i+1)
						}
					}

				}
			case []byte:
				var unmarshalResultByte []byte
				if err := Unmarshal(result, &unmarshalResultByte); err != nil {
					t.Fatal("Couldn't unmarshal data in test case #", i+1, ":", err)
				} else {
					for resultIndex, resultPiece := range unmarshalResultByte {
						if resultPiece != tc.data.([]byte)[resultIndex] {
							t.Fatal("Result does not match provided data in test case", i+1)
						}
					}
				}
			case int:
				var unmarshalResultInt int
				if err := Unmarshal(result, &unmarshalResultInt); err != nil {
					t.Fatal("Couldn't unmarshal data in test case #", i+1, ":", err)
				} else if unmarshalResultInt != tc.data.(int) {
					t.Fatal("Result does not match provided data in test case", i+1)
				}

			case float64:
				var unmarshalResultFloat float64
				if err := Unmarshal(result, &unmarshalResultFloat); err != nil {
					t.Fatal("Couldn't unmarshal data in test case #", i+1, ":", err)
				} else {
					if unmarshalResultFloat != tc.data.(float64) {
						fmt.Println(unmarshalResultFloat)
						fmt.Println(tc.data.(float64))
						t.Fatal("Result does not match provided data in test case", i+1)
					}
				}
			case bool:
				var unmarshalResultBool bool
				if err := Unmarshal(result, &unmarshalResultBool); err != nil {
					t.Fatal("Couldn't unmarshal data in test case #", i+1, ":", err)
				} else if unmarshalResultBool != tc.data {
					t.Fatal("Result does not match provided data in test case", i+1)
				}

			default:
				var unmarshalResult interface{}
				if err := Unmarshal(result, &unmarshalResult); err != nil {
					t.Fatal("Couldn't unmarshal data in test case #", i+1, ":", err)
				}
			}
		}
	}
}

func TestUnmarshal(t *testing.T) {
	type testCase struct {
		data []byte
	}

	cases := []testCase{
		{data: []byte(`[{"Test": "Name", "Data": "TestData"}]`)},
		{data: []byte(`{"More": "Data", "Value": 10}`)},
		{data: []byte(`{"Message": "Hello", "Array": [1, 2, 3], "Null": null, "Number": 1.234}`)},
	}

	for i, tc := range cases {
		var result []byte
		if err := Unmarshal(tc.data, &result); err != nil {
			t.Fatal("Could not unmarshal JSON data in test case #", i+1, ":", err)
		} else if string(result) != string(tc.data) {
			t.Fatal("Result does not match provided data in test case", i+1)
		}
	}
}
