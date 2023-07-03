package tools

import (
	"context"
	pb "github.com/donglei1234/platform/services/condition/gen/condition/api"
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

func CreateRegisterCondition(
	owner, id, cType, progress int32,
	updateStrategy pb.Condition_UpdateStrategy,
	theme string, param0, param1, param2, param3 int32,
) *pb.Condition {
	return &pb.Condition{
		OwnerId:        owner,
		Id:             id,
		Params:         []int32{param0, param1, param2, param3},
		Type:           cType,
		Theme:          theme,
		Progress:       progress,
		UpdateStrategy: updateStrategy,
	}
}

func CreateUpdateCondition(cType, progress int32, param0, param1 int32) *pb.Condition {
	return &pb.Condition{
		Params:   []int32{param0, param1},
		Type:     cType,
		Progress: progress,
	}
}
