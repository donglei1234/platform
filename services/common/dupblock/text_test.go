package dupblock

import (
	"io"
	"testing"
)

func TestTextReader_Read(t *testing.T) {
	type testCase struct {
		data string
	}

	testCases := []testCase{
		{"key cheese\nset .\n[0,1,2,3,4]\n*\n"},
		{"key apple\n"},
		{"cpy source dest\n"},
		{"mov source dest\n"},
		{"swp source dest\n"},
	}

	for i, tc := range testCases {
		if r, err := NewTextReader(WithBytes([]byte(tc.data))); err != nil {
			t.Fatal(err)
		} else {
			command := Command{}
			j := 0

			for {
				if err := r.Read(&command); err != nil {
					if err == io.EOF {
						break
					} else {
						t.Fatal(i, j, err)
					}
				} else {
					// @TODO: SNICHOLS: validate the command read
					t.Log(i, j, command)
					j++
				}
			}
		}
	}
}
