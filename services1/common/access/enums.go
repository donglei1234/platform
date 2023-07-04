package access

type SessionAuthority int

const (
	Unknown_SessionAuthority = SessionAuthority(iota)
	Token_SessionAuthority
	Username_SessionAuthority
	JWT_SessionAuthority
	Steam_SessionAuthority
)

type AccessLevel int

const (
	AccessUndefined = AccessLevel(iota)
	AccessPublic
	AccessProtected
	AccessPrivate
)
