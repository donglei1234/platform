package tools

import (
	"context"
	"strings"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func ParseProfileId(ctx context.Context) (string, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", status.Errorf(codes.DataLoss, "ParseProfileId: failed to get metadata")
	}

	var profileId []string
	if profileId, ok = md["profileid"]; !ok {
		return "", status.Errorf(codes.DataLoss, "ParseProfileId: failed to get profileId")
	}

	return strings.Join(profileId, ""), nil
}

func ParseAppId(ctx context.Context) (string, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", status.Errorf(codes.DataLoss, "ParseAppId: failed to get metadata")
	}

	var appId []string
	if appId, ok = md["appid"]; !ok {
		return "", status.Errorf(codes.DataLoss, "ParseAppId: failed to get appId")
	}

	return strings.Join(appId, ""), nil
}
