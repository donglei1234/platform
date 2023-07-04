package tools

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	ErrUrlError = status.Error(codes.Internal, "ErrUrlError")
)
