package jsonx

import (
	"github.com/json-iterator/go"
)

var (
	jsonConfig = jsoniter.Config{
		EscapeHTML:                    false,
		MarshalFloatWith6Digits:       true,
		ObjectFieldMustBeSimpleString: true,
		SortMapKeys:                   true,
	}.Froze()
)

// DecodeFunction allows a closure to be used as the value interface in Unmarshal.  Implement your own custom
// decoding with it.
type DecodeFunction func(data []byte) error

// EncodeFunction allows a closure to be used when calling Marshal. Implement your own custom encoding using it.
type EncodeFunction func() (data []byte, err error)

// Unmarshal unmarshals json data into an interface.  How the data is processed depends on the receiving interface
// type.
func Unmarshal(data []byte, value interface{}) error {
	switch v := value.(type) {
	case *string:
		*v = string(data)
		return nil
	case *[]byte:
		*v = data
		return nil
	case *interface{}:
		return jsonConfig.Unmarshal(data, v)
	case interface{}:
		return jsonConfig.Unmarshal(data, v)
	case DecodeFunction:
		return v(data)
	default:
		return ErrInvalidDestinationForJson
	}
}

// Marshal marshals an interface into json data.  How the interface is marshaled depends on the interface type.
// Note: []byte type indicates already in JSON format - no marshalling necessary.
func Marshal(value interface{}) (data []byte, err error) {
	switch v := value.(type) {
	case []byte:
		data = v
	case *[]byte:
		data = *v
	case EncodeFunction:
		data, err = v()
	default:
		data, err = jsonConfig.Marshal(v)
	}
	return
}
