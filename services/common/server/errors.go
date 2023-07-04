package server

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	ErrInitAuthClient            = status.Error(codes.Internal, "ErrInitAuthClient")
	ErrMiddlewareParseJwtFailure = status.Error(codes.Internal, "ErrMiddlewareParseJwtFailure")
)
