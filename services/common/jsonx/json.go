// Package jsonx contains packages for querying, diffing, and manipulating JSON data.
package jsonx

import (
	"github.com/json-iterator/go"
)

// Dictionary is a generic JSON dictionary.
type Dictionary = map[string]interface{}

// Array is a generic JSON array.
type Array = []interface{}

// Value is a generic JSON value.
type Value = interface{}

// Stringify marshals a generic interface into a string and returns it.
func Stringify(value interface{}) (json string, err error) {
	return jsoniter.ConfigFastest.MarshalToString(value)
}

// Parse parses the given JSON data into a generic interface.
func Parse(data []byte, value interface{}) error {
	return jsoniter.ConfigFastest.Unmarshal(data, value)
}

// ParseAndReturn parses the given JSON data into a generic interface and returns it.
func ParseAndReturn(data []byte) (result interface{}, err error) {
	err = Parse(data, &result)
	return
}

// Parse parses the given JSON string into a generic interface.
func ParseString(s string, value interface{}) error {
	return Parse([]byte(s), value)
}

// ParseStringAndReturn parses the given JSON string into a generic interface and returns it.
func ParseStringAndReturn(s string) (result interface{}, err error) {
	err = ParseString(s, &result)
	return
}
