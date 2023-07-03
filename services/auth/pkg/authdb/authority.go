package authdb

import (
	"fmt"
	"github.com/donglei1234/platform/services/common/nosql/document"
)

type Authority interface {
	Name() string
	TokenString(token Token) (string, error)

	Enabled() bool
	SetEnabled(enabled bool)
	// AutoInvalidatePreviousSessions should return true if users should be restricted to only one session at a time
	AutoInvalidatePreviousSessions() bool
}

type InternalAuthority interface {
	Authority
	Authenticate(appId string, token Token, authStore document.DocumentStore) (UserId, error)
}

type ExternalAuthority interface {
	Authority

	CredentialsRequired() bool
	EmptyCredentials() Credentials
	CacheCredentials(Credentials) error
	GetCredentials() (Credentials, bool)

	Authenticate(credentials Credentials, token Token) (UserId, error)
}

type Credentials interface{}

type Token interface{}

type UserId interface {
	fmt.Stringer
}

func NewUserId(id string) UserId {
	return stringer(id)
}

type stringer string

func (s stringer) String() string {
	return string(s)
}
