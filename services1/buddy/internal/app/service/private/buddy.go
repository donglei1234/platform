package private

import (
	"context"
	pb "github.com/donglei1234/platform/services/proto/gen/buddy/api"
	"go.uber.org/zap"

	"github.com/donglei1234/platform/services/buddy/internal/db"
)

// GetProfileBuddies returns the provided profile's buddies.
func (s *Service) GetProfileBuddies(ctx context.Context, request *pb.GetProfileBuddiesRequest) (*pb.Buddies, error) {
	logger := s.LoggerForContext(ctx)
	logger = logger.With(zap.String("profile_id", request.ProfileId.ProfileId))
	if bq, e := s.db.LoadBuddyQueue(request.AppId.AppId, request.ProfileId.ProfileId); e != nil {
		logger.Error("load or create buddy queue failed", zap.Error(e))
		return nil, ErrGeneralFailure
	} else {
		response := &pb.Buddies{}
		for _, n := range bq.Buddies {
			buddy := pb.Buddy{Uid: n.Uid}
			//if n.State == latest.MemberTypeBuddy {
			//	buddy.Kind = &pb.Buddy_IsBuddy_{IsBuddy: &pb.Buddy_IsBuddy{Remark: n.Remark}}
			//} else {
			//	buddy.Kind = &pb.Buddy_IsInvited_{
			//		IsInvited: &pb.Buddy_IsInvited{ReqInfo: n.ReqInfo, Time: n.ReqTime},
			//	}
			//}
			response.Buddies = append(response.Buddies, &buddy)
		}
		return response, nil
	}
}

// WatchProfileBuddies returns a stream on which changes to the provided profile's buddies will be sent.
func (s *Service) WatchProfileBuddies(request *pb.WatchProfileBuddiesRequest, server pb.PrivateService_WatchProfileBuddiesServer) error {
	ctx := server.Context()
	logger := s.LoggerForContext(ctx)
	logger = logger.With(zap.String("profile_id", request.ProfileId.ProfileId))

	if _, err := db.NewRelativeBuddyQueuePath(request.AppId.AppId, request.ProfileId.ProfileId); err != nil {
		logger.Error(
			"new relative buddy queue path failed",
			zap.Error(err),
		)
		return ErrGeneralFailure
	} else {
		//if e := s.rosClient.WatchDocument(
		//	ctx,
		//	path,
		//	sendUpdate,
		//); e != nil {
		//	logger.Error(
		//		"watch document failed",
		//		zap.Error(e),
		//	)
		//	return ErrGeneralFailure
		//}

	}
	return nil
}

// GetProfilesBlockedList returns the provided profile's blocked list.
func (s *Service) GetProfileBlockedList(
	ctx context.Context,
	request *pb.GetProfileBlockedListRequest,
) (*pb.GetProfileBlockedListResponse, error) {
	logger := s.LoggerForContext(ctx)
	logger = logger.With(zap.String("profile_id", request.ProfileId.ProfileId))
	if bq, e := s.db.LoadBuddyQueue(request.AppId.AppId, request.ProfileId.ProfileId); e != nil {
		logger.Error("load or create buddy queue failed", zap.Error(e))
		return nil, ErrGeneralFailure
	} else {
		response := &pb.GetProfileBlockedListResponse{}
		for _, v := range bq.BlockedProfiles {
			bp := &pb.BlockedProfile{
				ProfileId: &pb.ProfileId{
					ProfileId: v.ID,
				},
				AddTime: v.AddTime,
			}
			response.Profiles = append(response.Profiles, bp)
		}
		return response, nil
	}
}

// Deprecated: GetBuddies returns the provided user's buddies.
func (s *Service) GetBuddies(ctx context.Context, request *pb.Name) (*pb.Buddies, error) {
	if session, e := s.sessionFromContext(ctx); e != nil {
		return nil, e
	} else if bq, e := s.db.LoadBuddyQueue(session.AppId, request.Uid); e != nil {
		return nil, e
	} else {
		response := &pb.Buddies{}

		for _, n := range bq.Buddies {
			buddy := pb.Buddy{Uid: n.Uid}
			//if n.State == latest.MemberTypeBuddy {
			//	buddy.Kind = &pb.Buddy_IsBuddy_{IsBuddy: &pb.Buddy_IsBuddy{Remark: n.Remark}}
			//} else {
			//	buddy.Kind = &pb.Buddy_IsInvited_{
			//		IsInvited: &pb.Buddy_IsInvited{ReqInfo: n.ReqInfo, Time: n.ReqTime},
			//	}
			//}
			response.Buddies = append(response.Buddies, &buddy)
		}
		return response, nil
	}
}

// Deprecated:WatchBuddies returns a stream on which changes to the provided user's buddies will be sent.
func (s *Service) WatchBuddies(request *pb.Name, server pb.PrivateService_WatchBuddiesServer) error {
	//ctx := server.Context()
	//logger := s.LoggerForContext(ctx)
	//if session, e := s.sessionFromContext(ctx); e != nil {
	//	return e
	//} else {
	//	update := pb.ROSUpdate{}
	//sendUpdate := func(data []byte) {
	//	update.Data = data
	//	if err := server.SendMsg(&update); err != nil {
	//		logger.Error(
	//			"could not send buddies change",
	//			zap.Error(err),
	//		)
	//	}
	//}

	//if path, err := db.NewRelativeBuddyQueuePath(session.AppId, request.Uid); err != nil {
	//	return err
	//} else {
	//	if e := s.rosClient.WatchDocument(
	//		ctx,
	//		path,
	//		sendUpdate,
	//	); e != nil {
	//		return e
	//	}
	//}
	//}
	return nil
}
