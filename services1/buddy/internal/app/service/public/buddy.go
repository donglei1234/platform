package public

import (
	"context"
	"github.com/donglei1234/platform/services/buddy/internal/db"
	"github.com/donglei1234/platform/services/common/jsonx"
	"github.com/donglei1234/platform/services/common/mq"
	util "github.com/donglei1234/platform/services/common/utils"
	"time"

	"go.uber.org/zap"

	//"github.com/donglei1234/platform/services/buddy/generated/grpc/go/buddy/api"
	pb "github.com/donglei1234/platform/services/proto/gen/buddy/api"

	"github.com/donglei1234/platform/services/buddy/internal/db/schemas/latest"
)

var BuddyChangesTopic = "BuddyChanges_"

var BuddyRewardTopic = "BuddyReward_"

// AddBuddy adds a buddy to the current user's queue.
// The requested buddy must accept via ReplyAddBuddy to fulfill the request.
func (s *Service) AddBuddy(ctx context.Context, req *pb.AddBuddyRequest) (*pb.AddBuddyResponse, error) {
	// 1、获取token
	token, err := util.ParseJwtToken(ctx)
	if err != nil {
		s.logger.Error("parse jwt token err", zap.Error(err))
		return &pb.AddBuddyResponse{}, ErrGeneralFailure
	}
	// 2、解析token，获取uid
	profileId, err := util.DecodeToken(token)
	if err != nil {
		s.logger.Error("jwt token decode to uid err", zap.Error(err))
		return &pb.AddBuddyResponse{}, ErrGeneralFailure
	}
	logger := s.logger.With(
		zap.String("function_name", "AddBuddy"),
		zap.String("profile_id", profileId),
		zap.String("profile_app_id", s.appId),
		zap.Strings("target_name", req.Uid),
	)
	if len(req.Uid) <= 0 {
		err = ErrInvalidParameter
		logger.Error("add failed because target is empty")
		return &pb.AddBuddyResponse{}, err
	}
	bqSelf, e := s.db.LoadBuddyQueue(s.appId, profileId)
	if e != nil {
		logger.Error(
			"error encountered while LoadOrCreate self BuddyQueue",
			zap.Error(e),
		)
		err = ErrGeneralFailure
	}
	reqInfo := latest.NewInviter(
		bqSelf.Uid,
		req.ReqInfo,
		time.Now().Unix(),
	)
	adds := make([]string, 0)
	failed := make([]string, 0)
	for _, v := range req.Uid {
		if profileId == v {
			continue
		}
		if err := s.addBuddy(bqSelf, v, reqInfo); err != nil {
			logger.Error("add buddy failed", zap.Error(err))
			failed = append(failed, v)
		} else {
			adds = append(adds, v)
		}
	}

	if err := bqSelf.Update(func() bool {
		for _, v := range adds {
			bqSelf.AddInviteSend(v, reqInfo)
		}
		return true
	}); err != nil {
		logger.Error("update self BuddyQueue failed", zap.Error(err))
		return &pb.AddBuddyResponse{}, ErrGeneralFailure
	}

	return &pb.AddBuddyResponse{
		Failed: failed,
	}, nil
}

func (s *Service) addBuddy(bqSelf *db.BuddyQueue, target string, reqInfo *latest.Inviter) error {
	if ids := bqSelf.FilterBlocked(target); len(ids) <= 0 {
		return ErrInSelfBlockedList
	} else if bqSelf.GetInviterNum() >= s.maxInviter {
		return ErrSelfInviterTopLimit
	} else if bqSelf.GetMemberCounts() >= s.maxBuddies {
		return ErrSelfBuddiesTopLimit
	} else if found := bqSelf.IsContains(target); found {
		return ErrBuddyAlreadyAdded
	} else if found := bqSelf.IsContainsInviter(target); found {
		return ErrBuddyAlreadyInYourRequestList
	} else if bqTarget, e := s.db.LoadBuddyQueue(
		s.appId,
		target,
	); e != nil {
		s.logger.Error("error encountered while LoadOrCreate target BuddyQueue", zap.Error(e))
		return ErrGeneralFailure
	} else if ids := bqTarget.FilterBlocked(bqSelf.Uid); len(ids) <= 0 {
		return ErrInTargetBlockedList
	} else if isFound := bqTarget.IsContains(bqSelf.Uid); isFound {
		return ErrBuddyAlreadyRequested
	} else if isFound := bqTarget.IsContainsInviter(bqSelf.Uid); isFound {
		return ErrBuddyAlreadyRequested
	} else if !bqTarget.Settings.AllowToBeAdded {
		return ErrNotAllowed
	} else if bqTarget.GetMemberCounts() >= s.maxBuddies {
		return ErrTargetBuddiesTopLimit
	} else {
		if err := bqTarget.Update(func() bool {
			bqTarget.AddInviter(reqInfo)
			return true
		}); err != nil {
			return ErrGeneralFailure
		} else {
			if err := s.mq.Publish(mq.NatsProtocol+BuddyChangesTopic+target, mq.WithJson(&pb.BuddyChanges{
				InviterAdded: []*pb.Inviter{reqInfo.ToProto()},
			})); err != nil {
				s.logger.Error("publish buddy changes failed", zap.Error(err))
			}
			return nil
		}
	}
}

// RemoveBuddy removes a buddy from the current user's queue.
// No corresponding acknowledgement is needed from the removed buddy.
func (s *Service) RemoveBuddy(ctx context.Context, req *pb.RemoveBuddyRequest) (response *pb.Nothing, err error) {
	// 1、获取token
	token, err := util.ParseJwtToken(ctx)
	if err != nil {
		s.logger.Error("parse jwt token err", zap.Error(err))
		return &pb.Nothing{}, ErrGeneralFailure
	}
	// 2、解析token，获取uid
	profileId, err := util.DecodeToken(token)
	if err != nil {
		s.logger.Error("jwt token decode to uid err", zap.Error(err))
		return &pb.Nothing{}, ErrGeneralFailure
	}

	logger := s.logger.With(
		zap.String("function_name", "RemoveBuddy"),
		zap.String("profile_id", profileId),
		zap.String("profile_app_id", s.appId),
		zap.String("target_name", req.Uid),
	)
	if bqSelf, e := s.db.LoadBuddyQueue(s.appId, profileId); e != nil {
		logger.Error(
			"error encountered while LoadOrCreate self BuddyQueue",
			zap.Error(e),
		)
		err = ErrGeneralFailure
	} else if e := bqSelf.Update(func() bool {
		bqSelf.Delete(req.Uid)
		return true
	}); e != nil {
		logger.Error(
			"error encountered while updating self BuddyQueue",
			zap.Error(e),
		)
		err = ErrGeneralFailure
	} else if bqTarget, e := s.db.LoadBuddyQueue(s.appId, req.Uid); e != nil {
		logger.Error(
			"error encountered while LoadOrCreate target BuddyQueue",
			zap.Error(e),
		)
		err = ErrGeneralFailure
	} else if e := bqTarget.Update(func() bool {
		bqTarget.Delete(profileId)
		return true
	}); e != nil {
		logger.Error(
			"error encountered while updating target BuddyQueue delete",
			zap.Error(e),
		)
		err = ErrGeneralFailure
	} else {
		// 发送通知
		if err := s.mq.Publish(mq.NatsProtocol+BuddyChangesTopic+req.Uid, mq.WithJson(&pb.BuddyChanges{
			Removed: []string{profileId},
		})); err != nil {
			logger.Error("error encountered while publishing message", zap.Error(err))
		}

		if err := s.mq.Publish(mq.NatsProtocol+BuddyChangesTopic+profileId, mq.WithJson(&pb.BuddyChanges{
			Removed: []string{req.Uid},
		})); err != nil {
			logger.Error("error encountered while publishing message", zap.Error(err))
		}

		response = NothingResponse
	}

	return
}

// GetBuddies returns the current user's buddies.
func (s *Service) GetBuddies(ctx context.Context, request *pb.Nothing) (response *pb.Buddies, err error) {
	// 1、获取token
	token, err := util.ParseJwtToken(ctx)
	if err != nil {
		s.logger.Error("parse jwt token err", zap.Error(err))
		return &pb.Buddies{}, ErrGeneralFailure
	}
	// 2、解析token，获取uid
	profileId, err := util.DecodeToken(token)
	if err != nil {
		s.logger.Error("jwt token decode to uid err", zap.Error(err))
		return &pb.Buddies{}, ErrGeneralFailure
	}
	logger := s.logger.With(
		zap.String("function_name", "GetBuddies"),
		zap.String("profile_id", profileId),
		zap.String("profile_app_id", s.appId),
	)
	if bq, e := s.db.LoadBuddyQueue(s.appId, profileId); e != nil {
		logger.Error(
			"error encountered while LoadOrCreate self BuddyQueue",
			zap.Error(e),
		)
		err = ErrGeneralFailure
	} else {
		response = &pb.Buddies{}
		for _, buddy := range bq.Buddies {
			response.Buddies = append(response.Buddies, buddy.ToProto())
		}
		for _, invite := range bq.Inviters {
			response.Inviters = append(response.Inviters, invite.ToProto())
		}
		response.InviterSends = make(map[string]*pb.Inviter)
		for k, inviteSend := range bq.InviterSends {
			response.InviterSends[k] = inviteSend.ToProto()
		}
		for _, blocked := range bq.BlockedProfiles {
			response.Blocked = append(response.Blocked, blocked.ToProto())
		}
	}
	return
}

// ReplyAddBuddy is called to accept or reject a buddy add request.
func (s *Service) ReplyAddBuddy(ctx context.Context, req *pb.ReplyAddBuddyRequest) (*pb.ReplyAddBuddyResponse, error) {
	// 1、获取token
	token, err := util.ParseJwtToken(ctx)
	if err != nil {
		s.logger.Error("parse jwt token err", zap.Error(err))
		return &pb.ReplyAddBuddyResponse{}, ErrGeneralFailure
	}
	// 2、解析token，获取uid
	profileId, err := util.DecodeToken(token)
	if err != nil {
		s.logger.Error("jwt token decode to uid err", zap.Error(err))
		return &pb.ReplyAddBuddyResponse{}, ErrGeneralFailure
	}
	logger := s.logger.With(
		zap.String("function_name", "ReplyAddBuddy"),
		zap.String("profile_id", profileId),
		zap.String("profile_app_id", s.appId),
		zap.Strings("target_name", req.Uid),
	)

	selfBq, e := s.db.LoadBuddyQueue(s.appId, profileId)
	if e != nil {
		return &pb.ReplyAddBuddyResponse{}, ErrGeneralFailure

	}
	addUids := make([]string, 0)
	if len(req.GetUid()) == 0 {
		addUids = selfBq.GetSortedInviters()
	} else {
		addUids = append(addUids, req.GetUid()...)
	}
	failed := make([]string, 0)
	addeds := make([]*latest.Buddy, 0)
	for _, v := range addUids {
		if req.Response {
			if add, err := s.acceptBuddy(selfBq, v); err != nil {
				logger.Error(
					"error encountered while acceptBuddy",
					zap.Error(err),
				)
				failed = append(failed, v)
			} else {
				addeds = append(addeds, add)
			}
		} else {
			if err = s.rejectBuddy(selfBq, v); err != nil {
				logger.Error(
					"error encountered while rejectBuddy",
					zap.Error(err),
				)
			}
		}
	}

	notice := make([]*pb.Buddy, 0)
	if err := selfBq.Update(func() bool {
		for _, v := range addeds {
			selfBq.RemoveInviter(v.Uid)
			selfBq.AddBuddy(v)
			notice = append(notice, v.ToProto())
		}
		return true
	}); err != nil {
		return &pb.ReplyAddBuddyResponse{}, ErrGeneralFailure
	}
	if err := s.mq.Publish(mq.NatsProtocol+BuddyChangesTopic+profileId, mq.WithJson(&pb.BuddyChanges{
		Added: notice,
	})); err != nil {
		logger.Error("error encountered while publishing message", zap.Error(err))
	}
	return &pb.ReplyAddBuddyResponse{
		Failed: failed,
	}, nil
}

func (s *Service) rejectBuddy(selfBq *db.BuddyQueue, inviter string) error {
	targetBuddyQueue, e := s.db.LoadBuddyQueue(s.appId, inviter)
	if e != nil {
		s.logger.Error("load buddy queue error", zap.Error(e))
		return ErrGeneralFailure
	}

	if err := targetBuddyQueue.Update(func() bool {
		targetBuddyQueue.RemoveInviteSend(selfBq.Uid)
		return true
	}); err != nil {
		s.logger.Error(
			"error encountered while updating target BuddyQueue",
			zap.Error(err),
		)
	}
	if e := selfBq.Update(func() bool {
		selfBq.RemoveInviter(inviter)
		return true
	}); e != nil {
		s.logger.Error(
			"error encountered while updating self BuddyQueue",
			zap.Error(e),
		)
		return ErrGeneralFailure
	}
	return nil
}

func (s *Service) acceptBuddy(selfBq *db.BuddyQueue, inviter string) (*latest.Buddy, error) {
	if selfBq.GetMemberCounts() >= s.maxBuddies {
		return nil, ErrSelfBuddiesTopLimit
	}
	if found := selfBq.IsContains(inviter); found {
		return nil, ErrBuddyAlreadyAdded
	}
	if found := selfBq.IsContainsInviter(inviter); !found {
		return nil, ErrInviterNotFound
	}

	targetBuddyQueue, e := s.db.LoadBuddyQueue(s.appId, inviter)
	if e != nil {
		s.logger.Error("load buddy queue error", zap.Error(e))
		return nil, ErrGeneralFailure
	}
	if targetBuddyQueue.GetMemberCounts() >= s.maxBuddies {
		return nil, ErrTargetBuddiesTopLimit
	}

	if targetBuddyQueue.IsBlocked(selfBq.Uid) {
		return nil, ErrInTargetBlockedList
	}
	self := latest.NewBuddy(selfBq.Uid, "", time.Now().Unix())
	if err := targetBuddyQueue.Update(func() bool {
		targetBuddyQueue.AddBuddy(self)
		targetBuddyQueue.RemoveInviteSend(selfBq.Uid)
		return true
	}); err != nil {
		return nil, ErrGeneralFailure
	}
	target := latest.NewBuddy(inviter, "", time.Now().Unix())

	if err := s.mq.Publish(mq.NatsProtocol+BuddyChangesTopic+inviter, mq.WithJson(&pb.BuddyChanges{
		Added: []*pb.Buddy{self.ToProto()},
	})); err != nil {
		s.logger.Error("error encountered while publishing message", zap.Error(err))
	}
	return target, nil
}

// WatchBuddies returns a stream on which changes to the current user's buddies will be sent.
func (s *Service) WatchBuddies(request *pb.Nothing, server pb.PublicService_WatchBuddiesServer) error {
	ctx := server.Context()
	// 1、获取token
	token, err := util.ParseJwtToken(ctx)
	if err != nil {
		s.logger.Error("parse jwt token err", zap.Error(err))
		return ErrNoMetaData
	}
	// 2、解析token，获取uid
	profileId, err := util.DecodeToken(token)
	if err != nil {
		s.logger.Error("jwt token decode to uid err", zap.Error(err))
		return ErrNoMetaData
	}
	logger := s.logger.With(
		zap.String("function_name", "WatchBuddies"),
		zap.String("profile_id", profileId),
		zap.String("profile_app_id", s.appId),
	)

	if buddyData, err := s.db.LoadOrCreateBuddyQueue(s.appId, profileId); err != nil {
		return err
	} else {
		changes := &pb.BuddyChanges{}
		for _, buddy := range buddyData.Buddies {
			changes.Added = append(changes.Added, buddy.ToProto())
		}
		for _, inviter := range buddyData.Inviters {
			changes.InviterAdded = append(changes.InviterAdded, inviter.ToProto())
		}
		if err := server.Send(changes); err != nil {
			logger.Error("send msg err", zap.Error(err))
			return ErrGeneralFailure
		}
	}

	if sChange, err := s.mq.Subscribe(mq.NatsProtocol+BuddyChangesTopic+profileId,
		s.handleBuddyChanges(profileId, server),
		mq.WithAtMostOnceDelivery(mq.DefaultId),
	); err != nil {
		logger.Info("subscribe err", zap.Error(err))
		return err
	} else if sReward, err := s.mq.Subscribe(mq.NatsProtocol+BuddyRewardTopic+profileId,
		s.handleBuddyChanges(profileId, server),
		mq.WithAtMostOnceDelivery(mq.DefaultId),
	); err != nil {
		logger.Info("subscribe err", zap.Error(err))
		return err
	} else {
		<-ctx.Done()
		if err := sChange.Unsubscribe(); err != nil {
			logger.Error("sChange unsubscribe err", zap.Error(err))
		}
		if err := sReward.Unsubscribe(); err != nil {
			logger.Error("sReward unsubscribe err", zap.Error(err))
		}

	}
	return nil
}

func (s *Service) handleBuddyChanges(profileId string, server pb.PublicService_WatchBuddiesServer) func(msg mq.Message, err error) mq.ConsumptionCode {
	return func(msg mq.Message, err error) mq.ConsumptionCode {
		update := &pb.BuddyChanges{}
		err = jsonx.Unmarshal(msg.Data(), &update)
		if err != nil {
			s.logger.Error("Unmarshal mq msg err", zap.Error(err))
			return mq.ConsumeNackPersistentFailure
		}
		s.logger.Debug("watch ", zap.Any("buddy changes", update))
		if update.RewardUpdate != nil {
			if buddyData, err := s.db.LoadBuddyQueue(s.appId, profileId); err != nil {
				return mq.ConsumeNackPersistentFailure
			} else {
				updates := make([]*pb.Buddy, len(update.RewardUpdate))
				for i, buddy := range update.RewardUpdate {
					updates[i] = buddyData.GetBuddy(buddy.Uid).ToProto()
				}
				if updates != nil {
					update.RewardUpdate = updates
				}
			}
		}

		if err := server.Send(update); err != nil {
			s.logger.Error("send msg err", zap.Error(err))
			return mq.ConsumeNackPersistentFailure
		}
		return mq.ConsumeAck
	}
}

// Remark adds a buddy remark.
func (s *Service) Remark(ctx context.Context, req *pb.RemarkRequest) (response *pb.Nothing, err error) {
	// 1、获取token
	token, err := util.ParseJwtToken(ctx)
	if err != nil {
		s.logger.Error("parse jwt token err", zap.Error(err))
		return &pb.Nothing{}, ErrGeneralFailure
	}
	// 2、解析token，获取uid
	profileId, err := util.DecodeToken(token)
	if err != nil {
		s.logger.Error("jwt token decode to uid err", zap.Error(err))
		return &pb.Nothing{}, ErrGeneralFailure
	}
	logger := s.logger.With(
		zap.String("function_name", "Remark"),
		zap.String("profile_id", profileId),
		zap.String("profile_app_id", s.appId),
	)

	if bq, e := s.db.LoadBuddyQueue(s.appId, profileId); e != nil {
		logger.Error(
			"error encountered while LoadOrCreate buddy queue document",
			zap.Error(e),
		)
		err = ErrGeneralFailure
	} else if e := bq.Update(func() bool {
		bq.UpdateRemark(req.Uid, req.Remark)
		return true
	}); e != nil {
		logger.Error(
			"error encountered while updating buddy remark",
			zap.Error(e),
		)
		err = ErrGeneralFailure
	} else {
		response = NothingResponse
	}

	return
}

func (s *Service) UpdateBuddySettings(
	ctx context.Context,
	req *pb.UpdateBuddySettingsRequest) (response *pb.Nothing, err error) {
	// 1、获取token
	token, err := util.ParseJwtToken(ctx)
	if err != nil {
		s.logger.Error("parse jwt token err", zap.Error(err))
		return &pb.Nothing{}, ErrGeneralFailure
	}
	// 2、解析token，获取uid
	profileId, err := util.DecodeToken(token)
	if err != nil {
		s.logger.Error("jwt token decode to uid err", zap.Error(err))
		return &pb.Nothing{}, ErrGeneralFailure
	}
	logger := s.logger.With(
		zap.String("function_name", "UpdateBuddySettings"),
		zap.String("profile_id", profileId),
		zap.String("profile_app_id", s.appId),
	)
	if bq, e := s.db.LoadBuddyQueue(s.appId, profileId); e != nil {
		logger.Error(
			"error encountered while LoadOrCreate buddy queue document",
			zap.Error(e),
		)
		err = ErrGeneralFailure
	} else if e := bq.Update(func() bool {
		bq.UpdateSettings(req.AllowToBeAdded)
		return true
	}); e != nil {
		logger.Error(
			"error encountered while updating buddy settings",
			zap.Error(e),
		)
		err = ErrGeneralFailure
	} else {
		response = NothingResponse
	}
	return
}

// GetBlockedProfiles get blocked users from the current user's blocked queue.
func (s *Service) GetBlockedProfiles(ctx context.Context, req *pb.Nothing) (response *pb.ProfileIds, err error) {
	// 1、获取token
	token, err := util.ParseJwtToken(ctx)
	if err != nil {
		s.logger.Error("parse jwt token err", zap.Error(err))
		return &pb.ProfileIds{}, ErrGeneralFailure
	}
	// 2、解析token，获取uid
	profileId, err := util.DecodeToken(token)
	if err != nil {
		s.logger.Error("jwt token decode to uid err", zap.Error(err))
		return &pb.ProfileIds{}, ErrGeneralFailure
	}
	logger := s.logger.With(
		zap.String("function_name", "GetBlockedProfiles"),
		zap.String("profile_id", profileId),
		zap.String("profile_app_id", s.appId),
	)
	if bq, e := s.db.LoadBuddyQueue(s.appId, profileId); e != nil {
		logger.Error(
			"error encountered while LoadOrCreate buddy queue document",
			zap.Error(e),
		)
		err = ErrGeneralFailure
	} else {
		profileIds := make([]*pb.ProfileId, 0)
		for _, v := range bq.BlockedProfiles {
			profileIds = append(profileIds, &pb.ProfileId{
				ProfileId: v.ID,
			})
		}

		response = &pb.ProfileIds{
			ProfileIds: profileIds,
		}
	}
	return
}

// AddBlockedProfile add  blocked users to the current user's blocked  queue.
func (s *Service) AddBlockedProfiles(ctx context.Context, req *pb.ProfileIds) (response *pb.Nothing, err error) {
	// 1、获取token
	token, err := util.ParseJwtToken(ctx)
	if err != nil {
		s.logger.Error("parse jwt token err", zap.Error(err))
		return &pb.Nothing{}, ErrGeneralFailure
	}
	// 2、解析token，获取uid
	profileId, err := util.DecodeToken(token)
	if err != nil {
		s.logger.Error("jwt token decode to uid err", zap.Error(err))
		return &pb.Nothing{}, ErrGeneralFailure
	}
	logger := s.logger.With(
		zap.String("function_name", "AddBlockedProfiles"),
		zap.String("profile_id", profileId),
		zap.String("profile_app_id", s.appId),
	)

	for _, v := range req.ProfileIds {
		if profileId == v.ProfileId {
			logger.Error(
				"can not add self to blocked list",
				zap.Error(ErrCanNotAddSelf),
			)
			err = ErrCanNotAddSelf
			return
		}
	}

	if bqSelf, e := s.db.LoadBuddyQueue(s.appId, profileId); e != nil {
		logger.Error(
			"error encountered while LoadOrCreate self BuddyQueue",
			zap.Error(e),
		)
		err = ErrGeneralFailure
	} else {
		var blockedProfiles []*latest.BlockedProfile
		var ids []string
		for _, v := range req.ProfileIds {
			bp := &latest.BlockedProfile{
				ID:      v.ProfileId,
				AddTime: time.Now().Unix(),
			}
			blockedProfiles = append(blockedProfiles, bp)
			ids = append(ids, v.ProfileId)
		}

		if ids = bqSelf.FilterBlocked(ids...); len(ids) <= 0 {
			response = NothingResponse
		} else if bqSelf.GetBlockedNum()+int32(len(ids)) > s.maxBlocked {
			err = ErrBlockedNumExceed
		} else {
			for _, v := range ids {
				if blockedProfile, err := s.db.LoadBuddyQueue(s.appId, v); err != nil {
					logger.Error(
						"error encountered while load or create blocked profile buddy queue",
						zap.String("blocked_id", v),
						zap.Error(err),
					)
				} else if err := blockedProfile.Update(func() bool {
					blockedProfile.Delete(profileId)
					return true
				}); err != nil {
					logger.Error(
						"error encountered while update blocked profile buddy queue",
						zap.String("blocked_id", v),
						zap.Error(err),
					)
				}
			}

			if e := bqSelf.Update(func() bool {
				bqSelf.AddBlockedProfiles(ids...)
				bqSelf.DeleteBuddies(ids...)
				bqSelf.DeleteRecentProfiles(ids...)
				return true
			}); e != nil {
				logger.Error(
					"error encountered while update add blocked profiles",
					zap.Error(e),
				)
				err = ErrGeneralFailure
			} else {
				// 发送删除通知
				notice := &pb.BuddyChanges{
					Removed: ids,
				}
				if err := s.mq.Publish(mq.NatsProtocol+BuddyChangesTopic+profileId, mq.WithJson(notice)); err != nil {
					logger.Error("error encountered while publishing message", zap.Error(err))
				}
				response = NothingResponse
			}

		}
	}
	return
}

// RemoveBlockedProfiles remove a blocked profile to the current profile's blocked  queue.
func (s *Service) RemoveBlockedProfiles(ctx context.Context, req *pb.ProfileIds) (response *pb.Nothing, err error) {
	// 1、获取token
	token, err := util.ParseJwtToken(ctx)
	if err != nil {
		s.logger.Error("parse jwt token err", zap.Error(err))
		return &pb.Nothing{}, ErrGeneralFailure
	}
	// 2、解析token，获取uid
	profileId, err := util.DecodeToken(token)
	if err != nil {
		s.logger.Error("jwt token decode to uid err", zap.Error(err))
		return &pb.Nothing{}, ErrGeneralFailure
	}
	logger := s.logger.With(
		zap.String("function_name", "RemoveBlockedProfiles"),
		zap.String("profile_id", profileId),
		zap.String("profile_app_id", s.appId),
	)
	var ids []string
	for _, v := range req.ProfileIds {
		ids = append(ids, v.ProfileId)
	}

	if bqSelf, e := s.db.LoadBuddyQueue(s.appId, profileId); e != nil {
		logger.Error(
			"error encountered while LoadOrCreate self BuddyQueue",
			zap.Error(e),
		)
		err = ErrGeneralFailure
	} else if e := bqSelf.Update(func() bool {
		bqSelf.DeleteBlockedProfiles(ids...)
		return true
	}); e != nil {
		logger.Error(
			"error encountered while update delete blocked profiles",
			zap.Error(e),
		)
		err = ErrGeneralFailure
	} else {
		response = NothingResponse
	}
	return
}

// AddToRecentMet add a  profile to the current profile's recent met queue.
func (s *Service) AddToRecentMet(ctx context.Context, req *pb.ProfileIds) (response *pb.Nothing, err error) {
	// 1、获取token
	token, err := util.ParseJwtToken(ctx)
	if err != nil {
		s.logger.Error("parse jwt token err", zap.Error(err))
		return &pb.Nothing{}, ErrGeneralFailure
	}
	// 2、解析token，获取uid
	profileId, err := util.DecodeToken(token)
	if err != nil {
		s.logger.Error("jwt token decode to uid err", zap.Error(err))
		return &pb.Nothing{}, ErrGeneralFailure
	}
	logger := s.logger.With(
		zap.String("function_name", "AddToRecentMet"),
		zap.String("profile_id", profileId),
		zap.String("profile_app_id", s.appId),
	)
	var ids []string
	for _, v := range req.ProfileIds {
		ids = append(ids, v.ProfileId)
	}

	if bq, e := s.db.LoadBuddyQueue(s.appId, profileId); e != nil {
		logger.Error(
			"error encountered while LoadOrCreate self BuddyQueue",
			zap.Error(e),
		)
		err = ErrGeneralFailure
	} else if ids = bq.FilterBlocked(ids...); len(ids) <= 0 {
		logger.Info("all profiles in self blocked list")
		response = NothingResponse
	} else if e := bq.Update(func() bool {
		bq.AddRecentProfiles(ids...)
		return true
	}); e != nil {
		logger.Error(
			"error encountered while update add recent profiles",
			zap.Error(e),
		)
		err = ErrGeneralFailure
	} else {
		response = NothingResponse
	}
	return
}

// StarBuddy star a buddy
func (s *Service) StarBuddy(ctx context.Context, req *pb.StarBuddyRequest) (response *pb.Nothing, err error) {
	// 1、获取token
	token, err := util.ParseJwtToken(ctx)
	if err != nil {
		s.logger.Error("parse jwt token err", zap.Error(err))
		return &pb.Nothing{}, ErrGeneralFailure
	}
	// 2、解析token，获取uid
	profileId, err := util.DecodeToken(token)
	if err != nil {
		s.logger.Error("jwt token decode to uid err", zap.Error(err))
		return &pb.Nothing{}, ErrGeneralFailure
	}
	logger := s.logger.With(
		zap.String("function_name", "StarBuddy"),
		zap.String("profile_id", profileId),
		zap.String("profile_app_id", s.appId),
	)
	if bq, e := s.db.LoadBuddyQueue(s.appId, profileId); e != nil {
		logger.Error(
			"error encountered while LoadOrCreate self BuddyQueue",
			zap.Error(e),
		)
		err = ErrGeneralFailure
	} else if isContains := bq.IsContainsProfiles(req.Uid...); !isContains {
		logger.Error("profiles not in self buddy list", zap.Strings("profile_ids", req.Uid))
		err = ErrBuddiesNotFound
	} else if e := bq.Update(func() bool {
		bq.AddFriendValue(req.FriendValue)
		return true
	}); e != nil {
		logger.Error(
			"error encountered while update add starred profiles",
			zap.Error(e),
		)
		err = ErrGeneralFailure
	} else {
		for _, v := range req.Uid {
			err := s.db.IncrBuddyReward(s.appId, v, "buddies."+profileId+".receiveRewardNum", 1)
			if err != nil {
				logger.Error("incr buddy reward err", zap.Error(err))
				continue
			}
			err = s.db.IncrBuddyReward(s.appId, v, "buddies."+profileId+".friendValue", req.FriendValue)
			if err != nil {
				logger.Error("incr buddy reward err", zap.Error(err))
				continue
			}

			if err := s.mq.Publish(mq.NatsProtocol+BuddyRewardTopic+v, mq.WithJson(&pb.BuddyChanges{
				RewardUpdate: []*pb.Buddy{
					{
						Uid: profileId,
					},
				},
			})); err != nil {
				logger.Error("error encountered while publishing message", zap.Error(err))
			}

		}
		response = NothingResponse
	}
	return
}

func (s *Service) FavoriteBuddy(ctx context.Context, req *pb.FavoriteBuddyRequest) (*pb.Nothing, error) {
	// 1、获取token
	token, err := util.ParseJwtToken(ctx)
	if err != nil {
		s.logger.Error("parse jwt token err", zap.Error(err))
		return &pb.Nothing{}, ErrGeneralFailure
	}
	// 2、解析token，获取uid
	profileId, err := util.DecodeToken(token)
	if err != nil {
		s.logger.Error("jwt token decode to uid err", zap.Error(err))
		return &pb.Nothing{}, ErrGeneralFailure
	}
	logger := s.logger.With(
		zap.String("function_name", "FavoriteBuddy"),
		zap.String("profile_id", profileId),
		zap.String("profile_app_id", s.appId),
	)

	if selfData, err := s.db.LoadBuddyQueue(s.appId, profileId); err != nil {
		logger.Error("error encountered while LoadOrCreate self BuddyQueue",
			zap.Error(err),
		)
		err = ErrGeneralFailure
	} else {
		if err := selfData.Update(func() bool {
			selfData.Favorite(req.IsFavorite, req.Uid...)
			return true
		}); err != nil {
			logger.Error("error encountered while update add favorite profiles",
				zap.Error(err),
			)
		}
	}
	return NothingResponse, err
}

func (s *Service) CollectBuddyReward(ctx context.Context, req *pb.ProfileIds) (*pb.CollectBuddyRewardResponse, error) {
	// 1、获取token
	token, err := util.ParseJwtToken(ctx)
	if err != nil {
		s.logger.Error("parse jwt token err", zap.Error(err))
		return &pb.CollectBuddyRewardResponse{}, ErrGeneralFailure
	}
	// 2、解析token，获取uid
	profileId, err := util.DecodeToken(token)
	if err != nil {
		s.logger.Error("jwt token decode to uid err", zap.Error(err))
		return &pb.CollectBuddyRewardResponse{}, ErrGeneralFailure
	}
	logger := s.logger.With(
		zap.String("function_name", "CollectBuddyReward"),
		zap.String("profile_id", profileId),
		zap.String("profile_app_id", s.appId),
	)

	if selfData, err := s.db.LoadBuddyQueue(s.appId, profileId); err != nil {
		logger.Error("error encountered while LoadOrCreate self BuddyQueue",
			zap.Error(err),
		)
		err = ErrGeneralFailure
	} else {
		res := int32(0)
		ids := make([]string, len(req.ProfileIds))
		for i, v := range req.ProfileIds {
			ids[i] = v.ProfileId
		}

		if err := selfData.Update(func() bool {
			if len(ids) <= 0 {
				res, ids = selfData.ClearRewardNum()
			} else {
				res = selfData.CollectRewardNum(ids)
			}
			return true
		}); err != nil {
			logger.Error("error encountered while update add favorite profiles",
				zap.Error(err),
			)
		}
		changes := make([]*pb.Buddy, 0)
		for _, v := range ids {
			changes = append(changes, &pb.Buddy{
				Uid: v,
			})
		}

		if res > 0 {
			if err := s.mq.Publish(mq.NatsProtocol+BuddyRewardTopic+profileId, mq.WithJson(&pb.BuddyChanges{
				RewardUpdate: changes,
			})); err != nil {
				logger.Error("error encountered while publishing message", zap.Error(err))
			}
		}
		return &pb.CollectBuddyRewardResponse{
			Num: res,
		}, err
	}
	return &pb.CollectBuddyRewardResponse{}, err
}

func (s *Service) IsInvited(ctx context.Context, request *pb.IsInvitedRequest) (*pb.IsInvitedResponse, error) {
	// 1、获取token
	token, err := util.ParseJwtToken(ctx)
	if err != nil {
		s.logger.Error("parse jwt token err", zap.Error(err))
		return &pb.IsInvitedResponse{}, ErrGeneralFailure
	}
	// 2、解析token，获取uid
	profileId, err := util.DecodeToken(token)
	if err != nil {
		s.logger.Error("jwt token decode to uid err", zap.Error(err))
		return &pb.IsInvitedResponse{}, ErrGeneralFailure
	}
	logger := s.logger.With(
		zap.String("function_name", "IsInvited"),
		zap.String("profile_id", profileId),
		zap.String("profile_app_id", s.appId),
	)
	// 3、获取好友列表
	if bq, e := s.db.LoadBuddyQueue(s.appId, profileId); e != nil {
		logger.Error(
			"error encountered while LoadOrCreate self BuddyQueue",
			zap.Error(e),
		)
		err = ErrGeneralFailure
	} else {
		if isOk := bq.IsContains(request.Uid); isOk {
			return &pb.IsInvitedResponse{
				IsInvited: true,
			}, nil
		} else if isInvited := bq.IsInvited(request.Uid); isInvited {
			return &pb.IsInvitedResponse{
				IsInvited: isInvited,
			}, nil
		}
	}
	return &pb.IsInvitedResponse{}, err
}
