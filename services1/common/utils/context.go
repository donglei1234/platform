package utils

import (
	"context"
	"runtime"
)

func UntilContextDone(ctx context.Context, f func() error) error {
	for {
		if err := f(); err != nil {
			return err
		}

		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			runtime.Gosched()
		}
	}
}
