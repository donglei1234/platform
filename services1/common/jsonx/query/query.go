package query

import (
	"github.com/donglei1234/platform/services/common/jsonx"
)

func Execute(opts ...Option) (output []byte, err error) {
	options := newOptions(opts...)
	var obj jsonx.Value

	if e := jsonx.Parse(options.input, &obj); e != nil {
		err = e
	} else {
		for _, a := range options.actions {
			if e := a(&obj); e != nil {
				err = e
				break
			}
		}

		if err == nil {
			output, err = jsonx.Marshal(obj)
		}
	}

	return
}
