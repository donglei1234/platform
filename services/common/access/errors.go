package access

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	ErrAccessDenied   = status.Error(codes.PermissionDenied, "ErrAccessDenied")
	ErrInvalidContext = status.Error(codes.InvalidArgument, "ErrInvalidContext")
	ErrMalformedToken = status.Error(codes.InvalidArgument, "ErrMalformedToken")
)
