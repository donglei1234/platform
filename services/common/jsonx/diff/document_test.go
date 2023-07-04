package diff

import (
	"fmt"
	"testing"

	"github.com/donglei1234/platform/services/common/jsonx"
)

func TestDocument(t *testing.T) {
	var testCases = []struct {
		completeJsonString string
		isValidJson        bool
		isEmpty            bool
	}{
		{
			completeJsonString: `{"name":"Tester", "age":999}}`,
			isValidJson:        true,
			isEmpty:            false,
		},
		{
			completeJsonString: `{ "name":"John", "age":30, "car":null }`,
			isValidJson:        true,
			isEmpty:            false,
		},
		{
			completeJsonString: `{}`,
			isValidJson:        true,
			isEmpty:            false,
		},
		{
			completeJsonString: `{ }`,
			isValidJson:        true,
			isEmpty:            false,
		},
	}

	for i, tc := range testCases {
		if document, err := NewDocument([]byte(tc.completeJsonString)); err != nil {
			t.Fatal("Error encountered creating a document in TestDocument() test case # ", i, ":", err)
		} else {
			if tc.isValidJson && tc.isEmpty == false {
				if document.Empty() {
					t.Error("Error: Document.Empty() returns unexpected value false for test case #", i+1)
				}

				if len(document.Data()) == 0 {
					t.Error("Error: Document.Data() returns unexpected len of 0 for test case #", i+1)
				}

				if document.CompareFragments(document) == false {
					t.Error("Document.CompareFragments() returns unexpected value of false for same document in test case #", i+1)
				}
			}
		}
	}
}

func TestDocumentDebugDump(t *testing.T) {
	var testCases = []struct {
		completeJsonString string
		isValidJson        bool
		isEmpty            bool
	}{
		{
			completeJsonString: `{"name":"Tester", "age":999}}`,
			isValidJson:        true,
			isEmpty:            false,
		},
		{
			completeJsonString: `{ "name":"John", "age":30, "car":null }`,
			isValidJson:        true,
			isEmpty:            false,
		},
	}
	for i, tc := range testCases {
		if testDoc, err := NewDocument([]byte(tc.completeJsonString)); err != nil {
			t.Fatal("Error encountered creating a document in TestDocumentDebugDump() test case # ", i, ":", err)
		} else {
			if tc.isValidJson && tc.isEmpty == false {
				if testDoc.Empty() {
					t.Error("Error: Document.Empty() returns unexpected value false for test case #", i+1)
				}

				if len(testDoc.Data()) == 0 {
					t.Error("Error: Document.Data() returns unexpected len of 0 for test case #", i+1)
				}

				if testDoc.CompareFragments(testDoc) == false {
					t.Error("Error: Document.CompareFragments() returns unexpected value of false for same document.",
						"Issue encountered in test case #", i+1)
				}

				testDoc.DebugDump()
			}
		}
	}
}

func TestContainsPath(t *testing.T) {
	var testCases = []struct {
		completeJsonString string
		isValidJson        bool
		isEmpty            bool
	}{
		{
			completeJsonString: `{"name":"Tester", "age":999}}`,
			isValidJson:        true,
			isEmpty:            false,
		},
		{
			completeJsonString: `{ "name":"John", "age":30, "car":null }`,
			isValidJson:        true,
			isEmpty:            false,
		},
	}
	for i, tc := range testCases {
		if testDoc, err := NewDocument([]byte(tc.completeJsonString)); err != nil {
			t.Fatal("Error encountered creating a document in TestContainsPath() test case # ", i, ":", err)
		} else {
			if tc.isValidJson && tc.isEmpty == false {
				if testDoc.Empty() {
					t.Error("Error: Document.Empty() returns unexpected value false for test case #", i+1)
				}

				if len(testDoc.Data()) == 0 {
					t.Error("Error: Document.Data() returns unexpected len of 0 for test case #", i+1)
				}

				if testDoc.CompareFragments(testDoc) == false {
					t.Error("Error: Document.CompareFragments() returns unexpected value of false for same document.",
						"Issue encountered in test case #", i+1)
				}
			}
			// The path can be any information contained within the JSON data; e.g. a field name
			testDoc.ContainsPath("age")
		}
	}
}

func TestFindFragmentByPath(t *testing.T) {
	var testCases = []struct {
		completeJsonString string
		isValidJson        bool
		isEmpty            bool
		ID                 int
		Depth              int
		Path               string
		Value              []byte
		ValueType          jsonx.ValueType
	}{
		{
			completeJsonString: `{"name":"Tester", "age":999}}`,
			isValidJson:        true,
			isEmpty:            false,
			ID:                 1,
			Depth:              5,
			Path:               "test path",
			Value:              []byte("value"),
			ValueType:          jsonx.ValueTypeString,
		},
		{
			completeJsonString: `{ "name":"John", "age":30, "car":null },{ "name":"Egbert", "age":22, "car":null }`,
			isValidJson:        true,
			isEmpty:            false,
			ID:                 2,
			Depth:              4,
			Path:               "test path 2",
			Value:              []byte("value"),
			ValueType:          jsonx.ValueTypeString,
		},
	}
	for i, tc := range testCases {
		if testDoc, err := NewDocument([]byte(tc.completeJsonString)); err != nil {
			t.Fatal("Error encountered creating a document in TestFindFragmentByPath() test case # ", i, ":", err)
		} else {
			if tc.isValidJson && tc.isEmpty == false {
				if testDoc.Empty() {
					t.Error("Error: Document.Empty() returns unexpected value false for test case #", i+1)
				}

				if len(testDoc.Data()) == 0 {
					t.Error("Error: Document.Data() returns unexpected len of 0 for test case #", i+1)
				}

				if testDoc.CompareFragments(testDoc) == false {
					t.Error("Error: Document.CompareFragments() returns unexpected value of false for same document.",
						"Issue encountered in test case #", i+1)
				}
			}
			testDoc.FindFragmentByPath("age")
		}
	}
}

func TestFindFragmentByHash(t *testing.T) {
	var testCases = []struct {
		completeJsonString string
		isValidJson        bool
		isEmpty            bool
		ID                 int
		Depth              int
		Path               string
		Value              []byte
		ValueType          jsonx.ValueType
	}{
		{
			completeJsonString: `{"name":"Tester", "age":999}}`,
			isValidJson:        true,
			isEmpty:            false,
			ID:                 1,
			Depth:              5,
			Path:               "test path",
			Value:              []byte("value"),
			ValueType:          jsonx.ValueTypeString,
		},
		{
			completeJsonString: `{ "name":"John", "age":30, "car":null },{ "name":"Egbert", "age":22, "car":null }`,
			isValidJson:        true,
			isEmpty:            false,
			ID:                 2,
			Depth:              4,
			Path:               "test path 2",
			Value:              []byte("value"),
			ValueType:          jsonx.ValueTypeString,
		},
	}
	for i, tc := range testCases {
		if testDoc, err := NewDocument([]byte(tc.completeJsonString)); err != nil {
			t.Fatal("Error encountered creating a document in TestFindFragmentByHash() test case # ", i, ":", err)
		} else {
			if tc.isValidJson && tc.isEmpty == false {
				if testDoc.Empty() {
					t.Error("Error: Document.Empty() returns unexpected value false for test case #", i+1)
				}

				if len(testDoc.Data()) == 0 {
					t.Error("Error: Document.Data() returns unexpected len of 0 for test case #", i+1)
				}

				if testDoc.CompareFragments(testDoc) == false {
					t.Error("Error: Document.CompareFragments() returns unexpected value of false for same document.",
						"Issue encountered in test case #", i+1)
				}
			}
			testFragment := NewFragment(tc.ID, tc.Depth, tc.Path, tc.Value, tc.ValueType)
			testDoc.FindFragmentByHash(testFragment.Hash)
		}
	}
}

func TestForEachFragmentWithHash(t *testing.T) {
	var testCases = []struct {
		completeJsonString string
		isValidJson        bool
		isEmpty            bool
		ID                 int
		Depth              int
		Path               string
		Value              []byte
		ValueType          jsonx.ValueType
	}{
		{
			completeJsonString: `{"name":"Tester", "age":999}}`,
			isValidJson:        true,
			isEmpty:            false,
			ID:                 1,
			Depth:              5,
			Path:               "test path",
			Value:              []byte("value"),
			ValueType:          jsonx.ValueTypeString,
		},
		{
			completeJsonString: `{ "name":"John", "age":30, "car":null },{ "name":"Egbert", "age":22, "car":null }`,
			isValidJson:        true,
			isEmpty:            false,
			ID:                 2,
			Depth:              4,
			Path:               "test path 2",
			Value:              []byte("value"),
			ValueType:          jsonx.ValueTypeString,
		},
		{
			completeJsonString: `{"name":"Tester Three", "age":972}}`,
			isValidJson:        true,
			isEmpty:            false,
			ID:                 1,
			Depth:              5,
			Path:               "test path 3",
			Value:              []byte("value"),
			ValueType:          jsonx.ValueTypeString,
		},
	}
	var cb func(fragment *Fragment) bool
	for i, tc := range testCases {
		if testDoc, err := NewDocument([]byte(tc.completeJsonString)); err != nil {
			t.Fatal("Error encountered creating a document in TestForEachFragmentWithHash() test case # ", i, ":", err)
		} else {
			if tc.isValidJson && tc.isEmpty == false {
				if testDoc.Empty() {
					t.Error("Error: Document.Empty() returns unexpected value false for test case #", i+1)
				}

				if len(testDoc.Data()) == 0 {
					t.Error("Error: Document.Data() returns unexpected len of 0 for test case #", i+1)
				}

				if testDoc.CompareFragments(testDoc) == false {
					t.Error("Error: Document.CompareFragments() returns unexpected value of false for same document.",
						"Issue encountered in test case #", i+1)
				}
			}

			testFragment := NewFragment(tc.ID, tc.Depth, tc.Path, tc.Value, tc.ValueType)
			testDoc.ForEachFragmentWithHash(testFragment.Hash, cb)
		}
	}
}

func TestContainsFragment(t *testing.T) {
	var testCases = []struct {
		completeJsonString string
		isValidJson        bool
		isEmpty            bool
		ID                 int
		Depth              int
		Path               string
		Value              []byte
		ValueType          jsonx.ValueType
	}{
		{
			completeJsonString: `{"name":"Tester", "age":999}}`,
			isValidJson:        true,
			isEmpty:            false,
			ID:                 1,
			Depth:              5,
			Path:               "test path",
			Value:              []byte("value"),
			ValueType:          jsonx.ValueTypeString,
		},
		{
			completeJsonString: `{ "name":"John", "age":30, "car":null },{ "name":"Egbert", "age":22, "car":null }`,
			isValidJson:        true,
			isEmpty:            false,
			ID:                 2,
			Depth:              4,
			Path:               "test path 2",
			Value:              []byte("value"),
			ValueType:          jsonx.ValueTypeString,
		},
	}
	for i, tc := range testCases {
		if testDoc, err := NewDocument([]byte(tc.completeJsonString)); err != nil {
			t.Fatal("Error encountered creating a document in TestContainsFragment() test case # ", i, ":", err)
		} else {
			if tc.isValidJson && tc.isEmpty == false {
				if testDoc.Empty() {
					t.Error("Error: Document.Empty() returns unexpected value false for test case #", i+1)
				}

				if len(testDoc.Data()) == 0 {
					t.Error("Error: Document.Data() returns unexpected len of 0 for test case #", i+1)
				}

				if testDoc.CompareFragments(testDoc) == false {
					t.Error("Error: Document.CompareFragments() returns unexpected value of false for same document.",
						"Issue encountered in test case #", i+1)
				}
			}
			testFragment := NewFragment(tc.ID, tc.Depth, tc.Path, tc.Value, tc.ValueType)
			testFragPointer := &testFragment
			testDoc.ContainsFragment(testFragPointer)
		}
	}
}

func TestForEachChildFragment(t *testing.T) {
	var testCases = []struct {
		completeJsonString string
		isValidJson        bool
		isEmpty            bool
		ID                 int
		Depth              int
		Path               string
		Value              []byte
		ValueType          jsonx.ValueType
	}{
		{
			completeJsonString: `{"name":"Tester", "age":999}}`,
			isValidJson:        true,
			isEmpty:            false,
			ID:                 1,
			Depth:              5,
			Path:               "test path",
			Value:              []byte("value"),
			ValueType:          jsonx.ValueTypeString,
		},
		{
			completeJsonString: `{ "name":"John", "age":30, "car":null },{ "name":"Egbert", "age":22, "car":null }`,
			isValidJson:        true,
			isEmpty:            false,
			ID:                 2,
			Depth:              4,
			Path:               "test path 2",
			Value:              []byte("value"),
			ValueType:          jsonx.ValueTypeString,
		},
	}

	var testCB func(index int, fragment *Fragment) bool

	for i, tc := range testCases {
		if testDoc, err := NewDocument([]byte(tc.completeJsonString)); err != nil {
			t.Fatal("Error encountered creating a document in TestForEachChildFragment() test case # ", i, ":", err)
		} else {
			if tc.isValidJson && tc.isEmpty == false {
				if testDoc.Empty() {
					t.Error("Error: Document.Empty() returns unexpected value false for test case #", i+1)
				}

				if len(testDoc.Data()) == 0 {
					t.Error("Error: Document.Data() returns unexpected len of 0 for test case #", i+1)
				}

				if testDoc.CompareFragments(testDoc) == false {
					t.Error("Error: Document.CompareFragments() returns unexpected value of false for same document.",
						"Issue encountered in test case #", i+1)
				}
			}
			testFragment := NewFragment(tc.ID, tc.Depth, tc.Path, tc.Value, tc.ValueType)
			testDoc.ForEachChildFragment(testFragment.ID, testCB)

		}
	}
}

func TestCountChildFragments(t *testing.T) {
	var testCases = []struct {
		completeJsonString string
		isValidJson        bool
		isEmpty            bool
		ID                 int
		Depth              int
		Path               string
		Value              []byte
		ValueType          jsonx.ValueType
	}{
		{
			completeJsonString: `{"name":"Tester", "age":999}}`,
			isValidJson:        true,
			isEmpty:            false,
			ID:                 1,
			Depth:              5,
			Path:               "test path",
			Value:              []byte("value"),
			ValueType:          jsonx.ValueTypeString,
		},
		{
			completeJsonString: `{ "name":"John", "age":30, "car":null },{ "name":"Egbert", "age":22, "car":null }`,
			isValidJson:        true,
			isEmpty:            false,
			ID:                 2,
			Depth:              4,
			Path:               "test path 2",
			Value:              []byte("value"),
			ValueType:          jsonx.ValueTypeString,
		},
	}

	for i, tc := range testCases {
		if testDoc, err := NewDocument([]byte(tc.completeJsonString)); err != nil {
			t.Fatal("Error encountered creating a document in TestCountChildFragments() test case # ", i, ":", err)
		} else {
			if tc.isValidJson && tc.isEmpty == false {
				if testDoc.Empty() {
					t.Error("Error: Document.Empty() returns unexpected value false for test case #", i+1)
				}

				if len(testDoc.Data()) == 0 {
					t.Error("Error: Document.Data() returns unexpected len of 0 for test case #", i+1)
				}

				if testDoc.CompareFragments(testDoc) == false {
					t.Error("Error: Document.CompareFragments() returns unexpected value of false for same document.",
						"Issue encountered in test case #", i+1)
				}
			}

			testDoc.CountChildFragments(i)
		}
	}
}

func TestCompareFragments(t *testing.T) {
	var testCases = []struct {
		completeJsonString string
		isValidJson        bool
		isEmpty            bool
	}{
		{`{"name":"Tester", "age":999}`, true, false},
		{`{"name":"Tester McBester", "age":999999}`, true, false},
	}
	for i, tc := range testCases {
		if testDoc, err := NewDocument([]byte(tc.completeJsonString)); err != nil {
			t.Fatal("Error encountered creating a document in TestCompareFragments() test case # ", i, ":", err)
		} else {
			if tc.isValidJson && tc.isEmpty == false {
				if testDoc.Empty() {
					t.Error("Error: Document.Empty() returns unexpected value false for test case #", i+1)
				}

				if len(testDoc.Data()) == 0 {
					t.Error("Error: Document.Data() returns unexpected len of 0 for test case #", i+1)
				}

				if testDoc.CompareFragments(testDoc) == false {
					t.Error("Error: Document.CompareFragments() returns unexpected value of false for same document.",
						"Issue encountered in test case #", i+1)
				}
			}
			if otherTestDoc, err := NewDocument([]byte(tc.completeJsonString)); err != nil {
				t.Fatal("Error encountered creating a second document in TestCompareFragments() test case # ", i, ":", err)

			} else {
				if tc.isValidJson && tc.isEmpty == false {
					if otherTestDoc.Empty() {
						t.Error("Error: Document.Empty() returns unexpected value false for test case #", i+1)
					}

					if len(otherTestDoc.Data()) == 0 {
						t.Error("Error: Document.Data() returns unexpected len of 0 for test case #", i+1)
					}

					if otherTestDoc.CompareFragments(otherTestDoc) == false {
						t.Error("Error: Document.CompareFragments() returns unexpected value of false for same document.",
							"Issue encountered in test case #", i+1)
					}

					fmt.Println(testDoc.fragments)
					fmt.Println(otherTestDoc.fragments)
					otherTestDoc.CompareFragments(testDoc)

				}
			}
		}
	}
}
