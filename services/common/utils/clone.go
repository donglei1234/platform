package utils

import "reflect"

// ShallowClone returns a shallow copy of a given interface.
func ShallowClone(i interface{}) interface{} {
	return reflect.Indirect(reflect.ValueOf(i)).Interface()
}
