package public

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	ErrGeneralFailure = status.Error(codes.Internal, "ErrGeneralFailure")
	ErrStreamClosed   = status.Error(codes.Internal, "ErrStreamClosed")
	InvalidReq        = status.Error(codes.InvalidArgument, "InvalidReq")
	InvalidProfileId  = status.Error(codes.InvalidArgument, "InvalidProfileId")
)
