package jsonx

import (
	"errors"
	"fmt"
	"reflect"
	"unsafe"
)

var (
	ErrInvalidDestinationForJson = errors.New("ErrInvalidDestinationForJson")
)

type ParseError struct {
	message  string
	location []byte

	line   int
	column int
}

func (e ParseError) Error() string {
	return fmt.Sprintf("%d:%d: %s", e.line, e.column, e.message)
}

func (e *ParseError) locate(data []byte) {
	dataHeader := (*reflect.SliceHeader)(unsafe.Pointer(&data))
	locationHeader := (*reflect.SliceHeader)(unsafe.Pointer(&e.location))

	if locationHeader.Data >= dataHeader.Data && (locationHeader.Data < dataHeader.Data+uintptr(dataHeader.Len)) {
		offset := int(locationHeader.Data - dataHeader.Data)
		line, col := 1, 0

		for index, ch := range data {
			if ch == '\n' {
				line++
				col = 0
			}

			col++

			if index == offset {
				break
			}
		}

		e.line = line
		e.column = col
	}
}
