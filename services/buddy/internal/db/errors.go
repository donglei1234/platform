package db

import "errors"

var (
	ErrLoadBuddyQueueFailed                  = errors.New("ErrLoadBuddyQueueFailed")
	ErrFixupFailed                           = errors.New("ErrFixupFailed")
	ErrUnhandledBuddyQueueVersion            = errors.New("ErrUnhandledBuddyQueueVersion")
	ErrUnexpectedEmptyBuddyQueue             = errors.New("ErrUnexpectedEmptyBuddyQueue")
	ErrBuddyQueueVersionExceedsLatestVersion = errors.New("ErrBuddyQueueVersionExceedsLatestVersion")
	ErrBuddyQueueUnmarshalFailure            = errors.New("ErrBuddyQueueUnmarshalFailure")
	ErrBuddyQueueNotFound                    = errors.New("ErrBuddyQueueNotFound")
	ErrNoDataChange                          = errors.New("ErrNoDataChange")
)
