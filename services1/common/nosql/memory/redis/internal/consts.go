package internal

import "time"

const (
	LockTTL        = time.Millisecond * 200
	LockMaxRetries = 5
)
