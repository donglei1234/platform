package query

import (
	"bytes"
	"fmt"
	"strconv"
)

var pathSeparator = []byte{'.'}

type PathNode struct {
	Index  int    // -1 indicates no Index
	Key    string // "" indicates no Key
	Parent interface{}
}

func parsePathElement(element []byte) (name string, index int, err error) {
	if start := bytes.IndexByte(element, '['); start != -1 {
		if end := bytes.IndexByte(element, ']'); end > start {
			index, err = strconv.Atoi(string(element[start+1 : end]))
			name = string(element[:start])
		} else {
			err = ErrInvalidPath
		}
	} else {
		index = -1
		name = string(element)
	}
	return
}

func pathToProperty(root interface{}, path []byte) (nodes []PathNode, err error) {
	// Special case - path is '.'
	if len(path) == 1 && path[0] == '.' {
		newNode := PathNode{
			Parent: root,
			Key:    "",
			Index:  -1,
		}
		nodes = append([]PathNode{newNode}, nodes...)
		return
	}

	var parent interface{}
	switch root.(type) {
	case *interface{}:
		parent = *root.(*interface{})
	default:
		parent = root
	}
	elements := bytes.Split(path, pathSeparator)

	for i := 0; i < len(elements); i++ {
		// Parse path element for later use
		if name, index, err := parsePathElement(elements[i]); err != nil {
			return nil, err
		} else if name == "" && index == -1 {
			return nil, ErrInvalidPath
		} else {
			// Handle name.
			if name != "" {
				switch v := parent.(type) {
				case map[string]interface{}:
					newNode := PathNode{
						Parent: v,
						Key:    string(name),
						Index:  -1,
					}
					nodes = append([]PathNode{newNode}, nodes...)
					parent = v[string(name)]
				default:
					return nil, ErrInvalidPath
				}
			}

			// Handle index.
			if index != -1 {
				switch v := parent.(type) {
				case []interface{}:
					newNode := PathNode{
						Parent: v,
						Index:  index,
						Key:    "",
					}
					nodes = append([]PathNode{newNode}, nodes...)
					if index < len(v) {
						parent = v[index]
					} else {
						parent = nil
					}
				default:
					if v != nil {
						fmt.Println(v)
					}
				}
			}
		}
	}

	return
}
