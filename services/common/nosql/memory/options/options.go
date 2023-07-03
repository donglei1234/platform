package options

import (
	"reflect"
	"time"

	"github.com/donglei1234/platform/services/common/nosql/errors"
)

// Options contains all the various options that the provided WithXyz functions construct.
type Options struct {
	AnyVersion  bool
	TTL         time.Duration
	Source      interface{}
	SourceList  []interface{}
	Sources     map[string]interface{}
	Destination interface{}
	//DestinationList []interface{}
	DestinationMap map[string]interface{}
}

// Option is a closure that updates Options.
type Option func(o *Options) error

// NewOptions constructs an Options struct from the provided Option closures and returns it.
func NewOptions(opts ...Option) (options Options, err error) {
	o := &options
	for _, opt := range opts {
		if err = opt(o); err != nil {
			break
		}
	}
	return
}

// Use WithTTL to set the expiration time of a value during an operation.
func WithTTL(ttl time.Duration) Option {
	return func(o *Options) error {
		o.TTL = ttl
		return nil
	}
}

// WithSource provides an interface to source data when updating a value.
func WithSource(src interface{}) Option {
	return func(o *Options) (err error) {
		o.Source = src
		return
	}
}

func WithSourceList(src []interface{}) Option {
	return func(o *Options) (err error) {
		o.SourceList = src
		return
	}
}

func WithMultipleSource(src map[string]interface{}) Option {
	return func(o *Options) (err error) {
		o.Sources = src
		return
	}
}

// WithDestination provides an interface for receiving data when getting a value.
func WithDestination(dst interface{}) Option {
	return func(o *Options) error {
		if dst == nil {
			return errors.ErrDestIsNil
		} else if reflect.TypeOf(dst).Kind() != reflect.Ptr {
			return errors.ErrDestMustBePointer
		} else {
			o.Destination = dst
			return nil
		}
	}
}

//// WithDestinationList provides a interface list for receiving data when getting a value.
//func WithDestinationList(dst []interface{}) Option {
//	return func(o *Options) error {
//		if dst == nil {
//			return errors.ErrDestIsNil
//		} else {
//			o.DestinationList = dst
//			return nil
//		}
//	}
//}

// WithDestinationList provides a interface list for receiving data when getting a value.
func WithDestinationMap(dst map[string]interface{}) Option {
	return func(o *Options) error {
		if dst == nil {
			return errors.ErrDestIsNil
		} else {
			o.DestinationMap = dst
			return nil
		}
	}
}
