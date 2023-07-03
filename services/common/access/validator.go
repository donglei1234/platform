package access

import (
	"context"
	"fmt"
)

const (
	GrpcServerContextKey = "grpcServer"
)

func ProtectedServiceFromContext(ctx context.Context) (svc ProtectedService, err error) {
	if v := ctx.Value(GrpcServerContextKey); v == nil {
		err = fmt.Errorf("unable to get service from context: %w", ErrInvalidContext)
	} else if s, ok := v.(ProtectedService); !ok {
		err = fmt.Errorf("context provided invalid service: %v: %w", svc, ErrInvalidContext)
	} else {
		svc = s
	}
	return
}

func ValidateSessionToken(svc ProtectedService, token string) (err error) {
	if key := ParseSessionKeyFromToken(token); key.String() != token {
		return fmt.Errorf("unable to parse provided session token: %v: %w", token, ErrMalformedToken)
	} else if key.AccessLevel < svc.AccessLevel() {
		return fmt.Errorf("session lacks required access (%v < %v): %w", key.AccessLevel, svc.AccessLevel(), ErrAccessDenied)
	}
	return
}
