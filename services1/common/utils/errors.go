package utils

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	ErrTokenInvalid     = status.Error(codes.Internal, "ErrTokenInvalid")
	ErrDecodeToken      = status.Error(codes.Internal, "ErrDecodeToken")
	ErrUnmarshalToken   = status.Error(codes.Internal, "ErrUnmarshalToken")
	ErrNull             = status.Error(codes.Internal, "ErrNull")
	ErrMetadataDataLost = status.Error(codes.Internal, "ErrMetadataDataLost")
)
