package public

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// error code
var (
	ErrRequestSysType = status.Error(codes.NotFound, "ErrRequestSysType")
	ErrRequestToken   = status.Error(codes.InvalidArgument, "ErrRequestToken")
	ErrApproveRequest = status.Error(codes.NotFound, "ErrApproveRequest")
	ErrUrlError       = status.Error(codes.Internal, "ErrUrlError")
	ErrJsonMarshal    = status.Error(codes.Internal, "ErrJsonMarshal")
	ErrGetUserId      = status.Error(codes.Internal, "ErrGetUserToken")
)

const (
	OKCODE int32  = 0
	OKMSG  string = "ok"
)
