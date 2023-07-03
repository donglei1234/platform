package options

import (
	"fmt"
	"github.com/pkg/errors"

	errors2 "github.com/donglei1234/platform/services/common/nosql/errors"
)

type ZRangeOptions struct {
	Start int64
	End   int64
}

type ZRangeOption func(o *ZRangeOptions) error

func NewZRangeOption(opts ...ZRangeOption) (options ZRangeOptions, err error) {
	o := &options
	for _, opt := range opts {
		if err = opt(o); err != nil {
			break
		}
	}
	return
}

func WithInterval(start, end int64) ZRangeOption {
	return func(o *ZRangeOptions) error {
		if end > 0 && end < start {
			return errors.Wrap(errors2.ErrInternal, fmt.Sprintf("start:%d<end:%d", start, end))
		} else {
			o.Start = start
			o.End = end
		}
		return nil
	}
}
