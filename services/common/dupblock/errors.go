package dupblock

import "errors"

var (
	ErrDataRequired     = errors.New("ErrDataRequired")
	ErrMalformedCommand = errors.New("ErrMalformedCommand")
	ErrUnhandledAction  = errors.New("ErrUnhandledAction")
	ErrUnknownDUPAction = errors.New("ErrUnknownDUPAction")
)
