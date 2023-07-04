package diff

import (
	"fmt"

	"github.com/cespare/xxhash"
	"github.com/donglei1234/platform/services/common/jsonx"
)

// Fragment encapsulates a JSON document fragment.
type Fragment struct {
	ID        int
	Depth     int
	Path      string
	Hash      uint64
	Value     []byte
	ValueType jsonx.ValueType
}

func NewFragment(id int, depth int, path string, value []byte, valueType jsonx.ValueType) Fragment {
	return Fragment{
		ID:        id,
		Depth:     depth,
		Path:      path,
		Hash:      xxhash.Sum64(value),
		Value:     value,
		ValueType: valueType,
	}
}

func (f Fragment) String() string {
	switch f.ValueType {
	case jsonx.ValueTypeString, jsonx.ValueTypeNumber:
		return fmt.Sprintf(`%0x: %s %x: %s = %s`, f.Depth, f.ValueType, f.Hash, f.Path, string(f.Value))
	default:
		return fmt.Sprintf("%0x: %s %x: %s", f.Depth, f.ValueType, f.Hash, f.Path)
	}
}

func (f *Fragment) IsEqual(o *Fragment) bool {
	return f.ValueType == o.ValueType && f.Hash == o.Hash //&& bytes.Equal(f.Value, o.Value)
}
