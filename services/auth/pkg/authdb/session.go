package authdb

import (
	"context"
	"time"

	"github.com/donglei1234/platform/services/common/access"
	"github.com/donglei1234/platform/services/common/nosql/document"
)

const (
	SessionContextKey = "session"
	SessionExpiry     = 1 * time.Hour
	// maximum number of old sessions to delete at once - "should" never be an issue, but we need some cap
	MaxSessionsToInvalidate = 5
	UserIdKey               = "UserId"
)

type Session struct {
	Key            access.SessionKey
	Token          string
	AppId          string
	Authority      string
	AuthorityToken string
	UserId         string
	Begin          time.Time

	// NB: While Session.Metadata may be unused by any core platform services, it is used by
	//     downstream API's that consume platform services to track arbitrary data between calls
	//     over the session's lifetime, e.g., by GP's session version checks.
	Metadata map[string]string

	Version document.Version `json:"-"`
}

func SessionFromContext(ctx context.Context) (*Session, bool) {
	s, ok := ctx.Value(SessionContextKey).(*Session)
	return s, ok
}

func NewSession(token, appId, authority, authorityToken, userId string) *Session {
	session := &Session{
		Key:            access.ParseSessionKeyFromToken(token),
		Token:          token,
		AppId:          appId,
		Authority:      authority,
		AuthorityToken: authorityToken,
		UserId:         userId,
		Begin:          time.Now(),
		Metadata:       make(map[string]string),
	}
	return session
}

func CopySession(key access.SessionKey, token, appId, authority, authorityToken, userId string, begin time.Time, meta map[string]string) (sess *Session) {
	sess = NewSession(token, appId, authority, authorityToken, userId)
	sess.Begin = begin
	sess.CopyMetadata(meta)
	sess.Key = key
	return
}

// Clobber any metadata we might have in favor of a copy of what is being passed in
func (s *Session) CopyMetadata(meta map[string]string) {
	if len(s.Metadata) > 0 {
		s.Metadata = make(map[string]string)
	}
	if meta != nil && len(meta) > 0 {
		for k, v := range meta {
			s.SetMetadata(k, v)
		}
	}
}

func (s *Session) SetMetadata(key string, val string) {
	s.Metadata[key] = val
}

func (s *Session) GetMetadata(key string) (val string, ok bool) {
	val, ok = s.Metadata[key]
	return
}
