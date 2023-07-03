package public

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	InvalidReq = status.Error(codes.InvalidArgument, "InvalidReq")
)
