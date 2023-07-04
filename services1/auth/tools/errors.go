package tools

import (
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	ErrMetadataDataLost   = errors.New("ErrMetadataDataLost")
	ErrTokenSigningMethod = errors.New("ErrTokenSigningMethod")
	ErrTokenExpired       = errors.New("ErrTokenExpired")
	ErrTokenMalformed     = errors.New("ErrTokenMalformed")
	ErrTokenHandle        = errors.New("ErrTokenHandle")
	ErrNull               = status.Error(codes.Internal, "errNull")
	ErrTokenInvalid       = status.Error(codes.Internal, "ErrTokenInvalid")
	ErrDecodeToken        = status.Error(codes.Internal, "ErrDecodeToken")
	ErrUnmarshalToken     = status.Error(codes.Internal, "ErrUnmarshalToken")
	ErrSignedString       = errors.New("ErrSignedString")
)
