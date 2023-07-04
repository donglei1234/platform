package public

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	ErrGeneralFailure            = status.Error(codes.Internal, "ErrGeneralFailure")
	ErrSetTokenFailure           = status.Error(codes.Internal, "ErrSetTokenFailure")
	ErrGenerateJwtFailure        = status.Error(codes.Internal, "ErrGenerateJwtFailure")
	ErrSetIDFailure              = status.Error(codes.Internal, "ErrSetIDFailure")
	ErrGetIDFailure              = status.Error(codes.Internal, "ErrGetIDFailure")
	ErrMiddlewareParseJwtFailure = status.Error(codes.Internal, "ErrMiddlewareParseJwtFailure")
	ErrGetFacebookInfoFailure    = status.Error(codes.Internal, "ErrGetFacebookInfoFailure")
	ErrGetSessionFailure         = status.Error(codes.Internal, "ErrGetSessionFailure")
	ErrGenerateIDFailure         = status.Error(codes.Internal, "ErrGenerateIDFailure")
	ErrGenerateTokenFailure      = status.Error(codes.Internal, "ErrGenerateTokenFailure")
	InvalidProfileId             = status.Error(codes.InvalidArgument, "InvalidProfileId")
	InvalidJwtToken              = status.Error(codes.InvalidArgument, "InvalidJwtToken")
	ErrFacebookTokenTimeOut      = status.Error(codes.DeadlineExceeded, "ErrFacebookTokenTimeOut")
	ErrDecodeFacebookInfoFailure = status.Error(codes.Internal, "ErrDecodeFacebookInfoFailure")
	ErrFacebookInvalid           = status.Error(codes.Internal, "ErrFacebookInvalid")
	ErrInitGoogleServiceFailure  = status.Error(codes.Internal, "ErrInitGoogleServiceFailure")
	ErrGetGoogleInfoFailure      = status.Error(codes.Internal, "ErrGetGoogleInfoFailure")
	ErrGoogleAudienceInvalid     = status.Error(codes.Internal, "ErrGoogleAudienceInvalid")
	ErrGoogleTokenExpired        = status.Error(codes.Internal, "ErrGoogleTokenExpired")
	ErrParseJwtTokenFailure      = status.Error(codes.Internal, "ErrParseJwtTokenFailure")
	ErrTokenSigningMethod        = status.Error(codes.Internal, "ErrTokenSigningMethod")
	ErrTokenExpired              = status.Error(codes.Internal, "ErrTokenExpired")
	ErrTokenMalformed            = status.Error(codes.Internal, "ErrTokenMalformed")
	ErrTokenHandle               = status.Error(codes.Internal, "ErrTokenHandle")
)
