package jsonx

import (
	"fmt"
)

// Return signal from an InterateCallback handler.
type IterateCallbackResult int

const (
	// Continue iterating recursively into the current value.
	IterateCallbackContinue IterateCallbackResult = iota

	// Skip iterating recursively into the current value.
	IterateCallbackSkip IterateCallbackResult = iota

	// Stop iterating.
	IterateCallbackStop IterateCallbackResult = iota
)

const (
	pathLimit = 1024
)

func appendKey(path []byte, key []byte) (newPath []byte) {
	if len(path) > 0 {
		return append(append(path, '.'), key...)
	} else {
		return append(path, key...)
	}
}

// Process one step of the iteration process.
type IterateCallback func(depth int, path []byte, key []byte, value []byte, valueType ValueType) IterateCallbackResult

// Iterate over JSON data and invoke the provided callback function for each value encountered.  Iteration
// is top-down and recursive.  The callback can choose how to proceed with the current value by returning
// the appropriate IterateCallbackResult.
func Iterate(data []byte, callback IterateCallback) error {
	var result error

	func() {
		defer func() {
			if r := recover(); r != nil {
				switch v := r.(type) {
				case *ParseError:
					v.locate(data)
					result = v
				case error:
					result = v
				}
			}
		}()

		path := [pathLimit]byte{}
		iterateValue(0, path[:0], nil, data, callback)
	}()

	return result
}

// Iterate over the keys / value pairs in the object that's contained in the provided slice.  The expectation is that
// the slice starts just after the opening '{' and ends just before the closing '}'.
//
// {"property":"value"}
//
//	^                ^
func iterateObject(depth int, path []byte, data []byte, callback IterateCallback) IterateCallbackResult {
	result := IterateCallbackContinue

	for len(data) > 0 && result == IterateCallbackContinue {
		data = skipWhitespace(data)

		// handle empty object case
		if len(data) == 0 {
			break
		}

		if key, err := scanString(data); err != nil {
			panic(err)
		} else {
			// consume the keys and any whitespace
			data = skipWhitespace(data[len(key):])

			// remove quotes from keys
			key = key[1 : len(key)-1]

			if data[0] == ':' {
				// consume the colon and any whitespace
				data = skipWhitespace(data[1:])

				// iterate the child value
				var value []byte
				result, value, _ = iterateValue(depth, appendKey(path, key), key, data, callback)

				if result == IterateCallbackContinue {
					// eat the value and whitespace
					data = skipWhitespace(data[len(value):])

					// next has to be a comma or end of object
					if len(data) == 0 {
						break
					} else if data[0] == ',' {
						data = data[1:]
					} else {
						panic(&ParseError{
							message:  `',' or '}' expected`,
							location: data,
						})
					}
				}
			} else {
				panic(&ParseError{
					message:  `':' expected`,
					location: data,
				})
			}
		}
	}

	return result
}

// Iterate over the array values in the object that's contained in the provided slice.  The expectation is that
// the slice starts just after the opening '[' and ends just before the closing ']'.
//
// [value,value,value]
//
//	^               ^
func iterateArray(depth int, path []byte, data []byte, callback IterateCallback) IterateCallbackResult {
	result := IterateCallbackContinue

	index := 0
	buffer := [16]byte{}

	for len(data) > 0 && result == IterateCallbackContinue {
		data = skipWhitespace(data)

		key := indexPath(index, buffer[:])

		var value []byte
		result, value, _ = iterateValue(depth, append(path, key...), key, data, callback)

		if result == IterateCallbackContinue {
			// eat the value and whitespace
			data = skipWhitespace(data[len(value):])

			// next has to be a comma or end of array
			if len(data) == 0 {
				break
			} else if data[0] == ',' {
				data = data[1:]
				index++
			} else {
				panic(&ParseError{
					message:  `',' or ']' expected`,
					location: data,
				})
			}
		}
	}

	return result
}

// Iterate over an arbitrary value in the provided slice.
func iterateValue(
	depth int,
	path []byte,
	key []byte,
	data []byte,
	callback IterateCallback,
) (IterateCallbackResult, []byte, ValueType) {
	if value, valueType, err := scanValue(data); err != nil {
		panic(err)
	} else {
		result := callback(depth, path, key, value, valueType)

		if result == IterateCallbackContinue {
			switch valueType {
			case ValueTypeObject:
				return iterateObject(depth+1, path, value[1:len(value)-1], callback), value, valueType
			case ValueTypeArray:
				return iterateArray(depth+1, path, value[1:len(value)-1], callback), value, valueType
			}
		}

		return result, value, valueType
	}
}

// Skip any prefix whitespace in the provided slice.  Return a subslice with the whitepsace removed.
func skipWhitespace(data []byte) []byte {
	if l := len(data); l > 0 {
		switch data[0] {
		case ' ', '\t', '\n', '\r':
			for index := 1; index < l; index++ {
				switch data[index] {
				case ' ', '\t', '\n', '\r':
					continue
				default:
					return data[index:]
				}
			}
		default:
			return data
		}
	}

	return data[len(data):]
}

// Scan for a JSON string escape sequence in the given slice.  Return a subslice that contains the sequence.
func scanEscapeSequence(data []byte) (value []byte, err error) {
	if len(data) > 0 && data[0] == '\\' {
		length := len(data)

		if length > 1 {
			switch data[1] {
			case 'b', 'f', 'n', 't', 'r', '\\', '"', '\'':
				return data[:2], nil
			case 'u':
				if length > 5 {
					return data[:6], nil
				} else {
					return nil, &ParseError{
						message:  `invalid unicode escape sequence`,
						location: data,
					}
				}
			default:
				return nil, &ParseError{
					message:  `invalid escape sequence`,
					location: data,
				}
			}
		} else {
			return nil, &ParseError{
				message:  `invalid escape sequence`,
				location: data,
			}
		}
	} else {
		return nil, &ParseError{
			message:  `'\' expected`,
			location: data,
		}
	}
}

// Scan for a string in the given slice.  Return a subslice that contains the string.
func scanString(data []byte) (value []byte, err error) {
	if l := len(data); l > 0 && data[0] == '"' {
		for index := 1; index < l; index++ {
			c := data[index]

			if c == '"' {
				return data[:index+1], nil
			} else if c == '\\' {
				return scanStringSlow(data, l, index)
			}
		}
	}

	return nil, &ParseError{
		message:  `" expected`,
		location: data,
	}
}

func scanStringSlow(data []byte, l int, index int) (value []byte, err error) {
	for index < l {
		c := data[index]

		if c == '\\' {
			if sequence, err := scanEscapeSequence(data[index:]); err != nil {
				return nil, err
			} else {
				index += len(sequence)
			}
		} else if c == '"' {
			return data[:index+1], nil
		}
	}

	return nil, &ParseError{
		message:  `" expected`,
		location: data,
	}
}

// Scan for an object started at the provided slice.  Returns a subslice that contains the object.
func scanObject(data []byte) (value []byte, err error) {
	if len(data) > 0 && data[0] == '{' {
		depth := 0
		index := 0

		// scan for the matching closing '}'
		for index < len(data) {
			switch data[index] {
			case '"':
				if stringContent, err := scanString(data[index:]); err != nil {
					return nil, err
				} else {
					index += len(stringContent)
				}
			case '[', '{':
				depth++
				index++
			case ']':
				depth--
				index++
			case '}':
				depth--

				if depth == 0 {
					// done parsing, return
					return data[:index+1], nil
				} else {
					index++
				}
			default:
				index++
			}
		}

		return nil, &ParseError{
			message:  `'}' expected`,
			location: data,
		}
	} else {
		return nil, &ParseError{
			message:  `'{' expected`,
			location: data,
		}
	}
}

// Scan for an object started at the provided slice.  Returns a subslice that contains the array.
func scanArray(data []byte) (value []byte, err error) {
	if len(data) > 0 && data[0] == '[' {
		depth := 0
		index := 0

		// scan for the matching closing ']'
		for index < len(data) {
			switch data[index] {
			case '"':
				if stringContent, err := scanString(data[index:]); err != nil {
					return nil, err
				} else {
					index += len(stringContent)
				}
			case '[', '{':
				depth++
				index++
			case '}':
				depth--
				index++
			case ']':
				depth--

				if depth == 0 {
					// done parsing, return
					return data[:index+1], nil
				} else {
					index++
				}
			default:
				index++
			}
		}

		return nil, fmt.Errorf("']' expected")
	} else {
		return nil, fmt.Errorf("'[' expected")

	}
}

// Scan for "true" in the given slice.  Returns a subslice containing it.
func scanTrue(data []byte) (value []byte, err error) {
	if len(data) >= 4 && data[0] == 't' && data[1] == 'r' && data[2] == 'u' && data[3] == 'e' {
		return data[:4], nil
	} else {
		return nil, &ParseError{
			message:  `true expected`,
			location: data,
		}
	}
}

// Scan for "false" in the given slice.  Returns a subslice containing it.
func scanFalse(data []byte) (value []byte, err error) {
	if len(data) >= 5 && data[0] == 'f' && data[1] == 'a' && data[2] == 'l' && data[3] == 's' && data[4] == 'e' {
		return data[:5], nil
	} else {
		return nil, &ParseError{
			message:  `false expected`,
			location: data,
		}
	}
}

// Scan for "null" in the given slice.  Returns a subslice containing it.
func scanNull(data []byte) (value []byte, err error) {
	if len(data) >= 4 && data[0] == 'n' && data[1] == 'u' && data[2] == 'l' && data[3] == 'l' {
		return data[:4], nil
	} else {
		return nil, &ParseError{
			message:  `null expected`,
			location: data,
		}
	}
}

func isNumberByte(ch byte) bool {
	switch ch {
	case '-', '+', '0', '1', '2', '3', '4', '5', '6', '7', '8', '9', '.', 'e', 'E':
		return true
	default:
		return false
	}
}

func scanNumber(data []byte) (value []byte, err error) {
	if len(data) > 0 && (data[0] == '-' || (data[0] >= '0' && data[0] <= '9')) {
		index := 0
		var ch byte

		for _, ch = range data {
			if !isNumberByte(ch) {
				break
			}
			index++
		}

		return data[:index], nil
	} else {
		return nil, &ParseError{
			message:  `number expected`,
			location: data,
		}
	}
}

// Scan for a value in the given slice.  Returns a subslice containing the value as well as the value type.
func scanValue(data []byte) (value []byte, valueType ValueType, err error) {
	data = skipWhitespace(data)

	switch data[0] {
	case '{':
		valueType = ValueTypeObject
		value, err = scanObject(data)
	case '[':
		valueType = ValueTypeArray
		value, err = scanArray(data)
	case '"':
		valueType = ValueTypeString
		value, err = scanString(data)
	case 'n':
		valueType = ValueTypeNull
		value, err = scanNull(data)
	case 't':
		valueType = ValueTypeTrue
		value, err = scanTrue(data)
	case 'f':
		valueType = ValueTypeFalse
		value, err = scanFalse(data)
	case '-', '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		valueType = ValueTypeNumber
		value, err = scanNumber(data)
	default:
		valueType = ValueTypeUndefined
		err = &ParseError{
			message:  `value expected`,
			location: data,
		}
	}

	return
}

func indexPath(value int, buffer []byte) []byte {
	buffer[0] = '['
	index := 1

	for {
		digit := byte(value % 10)
		buffer[index] = digit + '0'
		index++
		value /= 10

		if value == 0 {
			break
		}
	}

	buffer[index] = ']'

	if index > 2 {
		// swap bytes
		left, right := 1, index-1

		for left < right {
			tmp := buffer[left]
			buffer[left] = buffer[right]
			buffer[right] = tmp
			left++
			right--
		}
	}

	return buffer[:index+1]
}
