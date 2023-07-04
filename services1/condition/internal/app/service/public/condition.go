package public

import (
	"context"
	"github.com/donglei1234/platform/services/common/jsonx"
	"github.com/donglei1234/platform/services/common/mq"
	util "github.com/donglei1234/platform/services/common/utils"
	pb "github.com/donglei1234/platform/services/condition/gen/condition/api"
	"go.uber.org/zap"
)

var ConditionTopic = "Condition_"

func (s *Service) Watch(_ *pb.Nothing, server pb.ConditionService_WatchServer) error {
	// 1、获取token
	token, err := util.ParseJwtToken(server.Context())
	if err != nil {
		s.logger.Error("parse jwt token err", zap.Error(err))
		return InvalidReq
	}
	// 2、解析token，获取uid
	profileId, err := util.DecodeToken(token)
	if err != nil {
		s.logger.Error("jwt token decode to uid err", zap.Error(err))
		return InvalidProfileId
	}
	subscribe, err := s.mq.Subscribe(mq.NatsProtocol+ConditionTopic+profileId,
		s.handleConditions(profileId, server),
		mq.WithAtMostOnceDelivery(mq.GroupId("condition"+profileId)),
	)
	if err != nil {
		return err
	}
	<-server.Context().Done()
	err = s.db.Clear(profileId)
	if err != nil {
		return err
	}
	err = subscribe.Unsubscribe()
	if err != nil {
		return err
	}
	return nil
}

func (s *Service) handleConditions(profileId string, server pb.ConditionService_WatchServer) func(msg mq.Message, err error) mq.ConsumptionCode {
	return func(msg mq.Message, err error) mq.ConsumptionCode {
		update := make([]*pb.Condition, 0)
		err = jsonx.Unmarshal(msg.Data(), &update)
		if err != nil {
			s.logger.Error("Unmarshal mq msg err", zap.Error(err))
			return mq.ConsumeNackPersistentFailure
		}
		s.logger.Debug("watch condition", zap.Any("condition", update))
		res := make([]*pb.Condition, 0)
		for _, v := range update {
			cds, err := s.db.GetAndDeleteFinishedConditions(profileId, v)
			if err != nil {
				s.logger.Error("get condition err", zap.Error(err))
				continue
			}
			if len(cds) <= 0 {
				continue
			}
			res = append(res, cds...)
		}
		if len(res) == 0 {
			return mq.ConsumeAck
		}
		err = server.Send(&pb.Changes{Conditions: res})
		if err != nil {
			s.logger.Error("send change err", zap.Error(err))
		}
		return mq.ConsumeAck
	}
}

func (s *Service) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.Nothing, error) {
	// 1、获取token
	token, err := util.ParseJwtToken(ctx)
	if err != nil {
		s.logger.Error("parse jwt token err", zap.Error(err))
		return &pb.Nothing{}, InvalidReq
	}
	// 2、解析token，获取uid
	profileId, err := util.DecodeToken(token)
	if err != nil {
		s.logger.Error("jwt token decode to uid err", zap.Error(err))
		return &pb.Nothing{}, InvalidProfileId
	}
	if len(req.Conditions) <= 0 {
		s.logger.Error("register condition is empty")
		return &pb.Nothing{}, InvalidReq
	}
	for _, v := range req.Conditions {
		if len(v.Params) < 3 {
			s.logger.Error("register condition param is not enough", zap.Any("params", v.Params))
			return &pb.Nothing{}, InvalidReq
		}
	}

	err = s.db.AddCondition(profileId, req.Conditions...)
	if err != nil {
		s.logger.Error("add condition err", zap.Error(err))
		return &pb.Nothing{}, err
	}
	err = s.mq.Publish(mq.NatsProtocol+ConditionTopic+profileId, mq.WithJson(req.Conditions))
	if err != nil {
		return nil, err
	}
	return &pb.Nothing{}, nil
}

func (s *Service) Unregister(ctx context.Context, req *pb.UnRegisterRequest) (*pb.Nothing, error) {
	// 1、获取token
	token, err := util.ParseJwtToken(ctx)
	if err != nil {
		s.logger.Error("parse jwt token err", zap.Error(err))
		return &pb.Nothing{}, InvalidReq
	}
	// 2、解析token，获取uid
	profileId, err := util.DecodeToken(token)
	if err != nil {
		s.logger.Error("jwt token decode to uid err", zap.Error(err))
		return &pb.Nothing{}, InvalidProfileId
	}

	err = s.db.DeleteCondition(profileId, req.Conditions...)
	if err != nil {
		s.logger.Error("delete condition err", zap.Error(err))
		return &pb.Nothing{}, err
	}
	return &pb.Nothing{}, nil
}

func (s *Service) Update(ctx context.Context, update *pb.UpdateRequest) (*pb.Nothing, error) {
	// 1、获取token
	token, err := util.ParseJwtToken(ctx)
	if err != nil {
		s.logger.Error("parse jwt token err", zap.Error(err))
		return &pb.Nothing{}, InvalidReq
	}
	// 2、解析token，获取uid
	profileId, err := util.DecodeToken(token)
	if err != nil {
		s.logger.Error("jwt token decode to uid err", zap.Error(err))
		return &pb.Nothing{}, InvalidProfileId
	}
	if update.Update == nil || len(update.Update) == 0 {
		s.logger.Error("update condition is empty", zap.Any("update", update))
		return &pb.Nothing{}, InvalidReq
	}
	for _, v := range update.Update {
		if len(v.Params) < 2 {
			s.logger.Error("update condition param is not enough", zap.Any("params", v.Params))
			return &pb.Nothing{}, InvalidReq
		}
	}

	err = s.db.UpdateCondition(profileId, update.Update...)
	if err != nil {
		s.logger.Error("update condition err", zap.Error(err))
		return &pb.Nothing{}, err
	}
	//TODO 可以只传数据变更的Type
	err = s.mq.Publish(mq.NatsProtocol+ConditionTopic+profileId, mq.WithJson(update.GetUpdate()))
	if err != nil {
		return nil, err
	}
	return &pb.Nothing{}, nil
}
