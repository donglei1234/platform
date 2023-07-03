package migrator

import "errors"

var (
	ErrUnexpectedVersionValue       = errors.New("ErrUnexpectedVersionValue")
	ErrUnexpectedFixupVersionValue  = errors.New("ErrUnexpectedFixupVersionValue")
	ErrUnexpectedNilUnderlyingValue = errors.New("ErrUnexpectedNilUnderlyingValue")
	ErrUnexpectedNonLatestVersion   = errors.New("ErrUnexpectedNonLatestVersion")
	ErrUnableToMigrateFromLatest    = errors.New("ErrUnableToMigrateFromLatest")
	ErrVersionMigrationFailure      = errors.New("ErrVersionMigrationFailure")
	ErrFixupFailure                 = errors.New("ErrFixupFailure")
	ErrNoFixupToApply               = errors.New("ErrNoFixupToApply")
)
