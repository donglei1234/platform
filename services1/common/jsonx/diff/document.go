package diff

import (
	"fmt"
	"sort"

	"github.com/donglei1234/platform/services/common/jsonx"
)

var jsonPathSeparator = []byte{'.'}

// Document is a read-only parsed representation of a JSON document.
type Document struct {
	data            []byte
	fragments       []Fragment
	fragmentsByPath []Fragment
	fragmentsByHash []Fragment
}

func NewDocument(data []byte) (Document, error) {
	result := Document{}
	return result, result.Parse(data)
}

func (d *Document) Parse(data []byte) error {
	d.data = data

	// make a slice of all fragments
	{
		fragments := make([]Fragment, 0, 64)
		fragmentIndex := 0

		if err := jsonx.Iterate(data, func(depth int, path []byte, key []byte, value []byte, valueType jsonx.ValueType) jsonx.IterateCallbackResult {
			fragments = append(fragments, NewFragment(fragmentIndex, depth, string(path), value, valueType))
			fragmentIndex++

			return jsonx.IterateCallbackContinue
		}); err != nil {
			return err
		}

		d.fragments = fragments
	}

	// make a slice of fragments sorted by path
	{
		fbp := make([]Fragment, len(d.fragments))
		copy(fbp, d.fragments)

		sort.Slice(fbp, func(i, j int) bool {
			return fbp[i].Path < fbp[j].Path
		})

		d.fragmentsByPath = fbp
	}

	// make a slice of fragments sorted by hash
	{
		fbh := make([]Fragment, len(d.fragments))
		copy(fbh, d.fragments)

		sort.Slice(fbh, func(i, j int) bool {
			return fbh[i].Hash < fbh[j].Hash
		})

		d.fragmentsByHash = fbh
	}

	return nil
}

func (d Document) DebugDump() {
	for _, fragment := range d.fragments {
		fmt.Println(fragment)
	}
}

func (d Document) ContainsPath(path string) bool {
	index := sort.Search(len(d.fragmentsByPath), func(i int) bool {
		return d.fragmentsByPath[i].Path >= path
	})

	if index == len(d.fragmentsByPath) || d.fragmentsByPath[index].Path != path {
		return false
	} else {
		return true
	}
}

func (d Document) FindFragmentByPath(path string) (*Fragment, bool) {
	index := sort.Search(len(d.fragmentsByPath), func(i int) bool {
		return d.fragmentsByPath[i].Path >= path
	})

	if index == len(d.fragmentsByPath) || d.fragmentsByPath[index].Path != path {
		return nil, false
	} else {
		return &d.fragmentsByPath[index], true
	}
}

func (d Document) FindFragmentByHash(hash uint64) (Fragment, bool) {
	index := sort.Search(len(d.fragmentsByHash), func(i int) bool {
		return d.fragmentsByHash[i].Hash >= hash
	})

	if index == len(d.fragmentsByHash) || d.fragmentsByHash[index].Hash != hash {
		return Fragment{}, false
	} else {
		return d.fragmentsByHash[index], true
	}
}

func (d Document) ForEachFragmentWithHash(hash uint64, cb func(fragment *Fragment) bool) {
	index := sort.Search(len(d.fragmentsByHash), func(i int) bool {
		return d.fragmentsByHash[i].Hash >= hash
	})

	for index < len(d.fragmentsByHash) && d.fragmentsByHash[index].Hash == hash {
		if !cb(&d.fragmentsByHash[index]) {
			return
		}

		index++
	}
}

func (d Document) ContainsFragment(fragment *Fragment) bool {
	index := sort.Search(len(d.fragmentsByPath), func(i int) bool {
		return d.fragmentsByPath[i].Path >= fragment.Path
	})

	return index < len(d.fragmentsByPath) && d.fragmentsByPath[index].Path == fragment.Path && d.fragmentsByPath[index].IsEqual(fragment)
}

func (d Document) ForEachChildFragment(index int, cb func(index int, fragment *Fragment) bool) {
	rootDepth := d.fragments[index].Depth
	childDepth := rootDepth + 1
	index++

	childIndex := 0

	for index < len(d.fragments) {
		fragment := &d.fragments[index]

		if fragment.Depth == rootDepth {
			break
		}

		if fragment.Depth == childDepth {
			if !cb(childIndex, fragment) {
				break
			}

			childIndex++
		}

		index++
	}
}

func (d Document) CountChildFragments(index int) int {
	rootDepth := d.fragments[index].Depth
	childDepth := rootDepth + 1
	index++

	count := 0

	for index < len(d.fragments) {
		fragment := &d.fragments[index]

		if fragment.Depth == rootDepth {
			break
		}

		if fragment.Depth == childDepth {
			count++
		}

		index++
	}

	return count
}

func (d Document) Data() []byte {
	return d.data
}

func (d Document) Empty() bool {
	return d.data == nil
}

func (d Document) CompareFragments(o Document) bool {
	if len(d.fragments) != len(o.fragments) {
		return false
	}

	for i, left := range d.fragmentsByPath {
		if !left.IsEqual(&o.fragmentsByPath[i]) {
			return false
		}
	}

	return true
}
