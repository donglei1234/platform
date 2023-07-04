package access

import (
	"fmt"

	"github.com/donglei1234/platform/services/common/nosql/document"
	"github.com/donglei1234/platform/services/common/utils"
)

type SessionKey struct {
	Key    document.Key
	legacy bool
	token  string

	AccessLevel AccessLevel
	Authority   SessionAuthority
	Metadata    int
}

func ParseSessionKey(key document.Key) SessionKey {
	sk := ParseSessionKeyFromToken(key.Base())
	sk.Key = key
	return sk
}

func ParseSessionKeyFromToken(token string) SessionKey {
	sk := SessionKey{}

	if n, err := fmt.Sscanf(token, "%d$%x$%s", &sk.AccessLevel, &sk.Metadata, &sk.token); err != nil || n != 3 {
		// we did not parse a new style key, flag it as old and move on
		sk.legacy = true
		sk.token = token
		// zero out values we might have accidentally parsed
		sk.AccessLevel = AccessUndefined
		sk.Metadata = 0
	} else {
		// pop the first digit into authority and keep the rest
		sk.Authority = SessionAuthority(sk.Metadata & 0x0F)
		sk.Metadata >>= 4
	}

	return sk
}

func (sk *SessionKey) String() string {
	meta := sk.Metadata<<4 + int(sk.Authority)
	return fmt.Sprintf("%d$%02x$%s", sk.AccessLevel, meta, sk.token)
}

func (sk *SessionKey) IsLegacy() bool {
	return sk.legacy
}

// Generate the random unique portion of a session key
func generateSessionToken() string {
	token, _ := utils.NewUUIDString()
	return token
}

func NewSessionKey(acc AccessLevel, meta int, auth SessionAuthority) SessionKey {
	token := generateSessionToken()
	sk := SessionKey{
		legacy:      false,
		token:       token,
		AccessLevel: acc,
		Authority:   auth,
		Metadata:    meta,
	}

	// NB: if nosql.Namespace is not set, this can fail silently
	//sk.Key, _ = nosql.SessionKey(sk.String())

	return sk
}
