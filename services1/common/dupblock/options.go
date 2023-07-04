package dupblock

import "bytes"

type textOptions struct {
	buffer *bytes.Buffer
	data   []byte
}

func (o *textOptions) build(opts ...TextOption) error {
	for _, f := range opts {
		if err := f(o); err != nil {
			return err
		}
	}
	return nil
}

type TextOption func(o *textOptions) error

func WithBuffer(buffer *bytes.Buffer) TextOption {
	return func(o *textOptions) error {
		o.buffer = buffer
		return nil
	}
}

func WithBytes(data []byte) TextOption {
	return func(o *textOptions) error {
		o.data = data
		return nil
	}
}
