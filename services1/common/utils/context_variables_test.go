package utils

import (
	"errors"
	"time"
)

var (
	TestErrContext = errors.New("TestErrContext")
	TestErrFunc    = errors.New("TestErrFunc")

	TickingContext = TickableTestContext{
		TestContext: TestContext{
			done:  make(chan struct{}, 1),
			error: TestErrContext,
			Time:  time.Time{},
			ok:    false,
			value: nil,
		},
		ticker: 1,
	}
)

type TickableTestContext struct {
	TestContext
	ticker int
}

func (t *TickableTestContext) Done() <-chan struct{} {
	if t.ticker > 0 {
		t.ticker--
		return nil
	} else {
		t.done <- struct{}{}
		return t.done
	}
}

type TestContext struct {
	done chan struct{}
	error
	time.Time
	isDone bool
	ok     bool
	value  interface{}
}

func (t *TestContext) Deadline() (deadline time.Time, ok bool) { return t.Time, t.ok }
func (t *TestContext) Done() <-chan struct{}                   { return t.done }
func (t *TestContext) Err() error                              { return t.error }
func (t *TestContext) Value(key interface{}) interface{}       { return t.value }
