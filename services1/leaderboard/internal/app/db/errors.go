package db

import (
	"errors"
)

var (
	ErrDuplicateListName = errors.New("ErrDuplicateListName")
	ErrListNotExists     = errors.New("ErrListNotExists")
	ErrMemberNotExists   = errors.New("ErrMemberNotExists")
)
