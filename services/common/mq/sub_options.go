package mq

// SubOptions contains all the various options that the provided WithXyz functions construct.
type SubOptions struct {
	Decoder           Decoder
	VPtrFactory       ValuePtrFactory
	DeliverySemantics DeliverySemantics
	GroupId           string
}

// SubOption is a closure that updates SubOptions.
type SubOption func(o *SubOptions) error

// NewSubOptions constructs an SubOptions struct from the provided SubOption closures and returns it.
func NewSubOptions(opts ...SubOption) (options SubOptions, err error) {
	o := &options

	for _, opt := range opts {
		if err = opt(o); err != nil {
			break
		}
	}
	return
}

// Use WithDecoder to decode a message's data payload into a value pointer, which is generated by passed-in ValuePtrFactory.
// Note: When using the WithDecoder option, Message data is only available via VPtr(). No data will be found in Data().
func WithDecoder(decoder Decoder, vPtrFactory ValuePtrFactory) SubOption {
	return func(o *SubOptions) error {
		if o.Decoder != nil {
			return ErrDecoderAlreadySet
		} else if o.VPtrFactory != nil {
			return ErrValuePtrAlreadySet
		} else {
			o.Decoder = decoder
			o.VPtrFactory = vPtrFactory
			return nil
		}
	}
}

// Configures delivery semantics of subscription to be at-least-once delivery.
// If a semantics preference is not set, mq implementation will use its default mode.
// Mutually exclusive with WithAtMostOnceDelivery()
func WithAtLeastOnceDelivery() SubOption {
	return func(o *SubOptions) error {
		if o.DeliverySemantics != Unset {
			return ErrSemanticsAlreadySet
		} else {
			o.DeliverySemantics = AtLeastOnce
			return nil
		}
	}
}

// Configures delivery semantics of subscription to be at-most-once delivery.
// If a semantics preference is not set, mq implementation will use its default mode.
// groupId is also optional. Pass in mq.DefaultId to have the mq implementation set a default groupId.
// Mutually exclusive with WithAtLeastOnceDelivery()
func WithAtMostOnceDelivery(groupId GroupId) SubOption {
	return func(o *SubOptions) error {
		if o.DeliverySemantics != Unset {
			return ErrSemanticsAlreadySet
		} else {
			o.DeliverySemantics = AtMostOnce
			o.GroupId = string(groupId)
			return nil
		}
	}
}