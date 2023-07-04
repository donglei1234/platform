package public

import (
	"context"
	"github.com/donglei1234/platform/services/common/jsonx"
	"github.com/donglei1234/platform/services/common/mq"
	pb "github.com/donglei1234/platform/services/proto/gen/gm/api"
	"go.uber.org/zap"
)

func (s *Service) SetProfilesBanStatus(ctx context.Context, request *pb.SetProfilesBanStatusRequest) (*pb.Nothing, error) {
	if err := s.db.SetProfileBanStatus(request.Profiles, request.Status); err != nil {
		return nil, err
	}
	return &pb.Nothing{}, nil
}

func (s *Service) GetProfilesBanStatus(ctx context.Context, request *pb.GetProfilesBanStatusRequest) (*pb.GetProfilesBanStatusResponse, error) {
	if infos, err := s.db.GetProfilesBanInfo(request.GetProfileIds()); err != nil {
		return nil, err
	} else {
		return &pb.GetProfilesBanStatusResponse{Profiles: infos}, nil
	}
}

func (s *Service) WatchProfilesBan(nothing *pb.Nothing, server pb.GmService_WatchProfilesBanServer) error {
	if subscribe, err := s.mq.Subscribe(mq.NatsProtocol+GMTopicBan.String(),
		s.handleChanges(server),
		mq.WithAtMostOnceDelivery(mq.DefaultId),
	); err != nil {
		return err
	} else {
		defer func(subscribe mq.Subscription) {
			err := subscribe.Unsubscribe()
			if err != nil {
				s.logger.Error("unsubscribe err", zap.Error(err))
			}
		}(subscribe)
	}

	banData, err := s.db.GetAllProfilesBan()
	if err != nil {
		return err
	}
	adds := make([]string, len(banData))
	for i, v := range banData {
		adds[i] = v.ProfileId
	}
	err = server.Send(&pb.ProfilesBanChanges{})
	if err != nil {
		return err
	}
	<-server.Context().Done()
	return nil
}

func (s *Service) handleChanges(server pb.GmService_WatchProfilesBanServer) mq.SubResponseHandler {
	return func(msg mq.Message, err error) mq.ConsumptionCode {
		changes := &pb.ProfilesBanChanges{}
		if err := jsonx.Unmarshal(msg.Data(), changes); err != nil {
			return mq.ConsumeNackPersistentFailure
		}
		err = server.Send(changes)
		if err != nil {
			s.logger.Error("send change err", zap.Error(err))
			return mq.ConsumeNackPersistentFailure
		}
		return mq.ConsumeAck
	}
}

func (s *Service) SendBulletin(ctx context.Context, request *pb.SendBulletinRequest) (*pb.Nothing, error) {
	info := request.GetBulletin()
	if err := s.db.SetBulletin(info); err != nil {
		return nil, err
	}
	return &pb.Nothing{}, nil
}

func (s *Service) GetBulletin(ctx context.Context, request *pb.GetBulletinRequest) (*pb.GetBulletinResponse, error) {
	if infos, err := s.db.GetBulletins(); err != nil {
		return nil, err
	} else {
		return &pb.GetBulletinResponse{Bulletins: infos}, nil
	}
}

func (s *Service) GetPurchaseRecord(ctx context.Context, request *pb.GetPurchaseRecordRequest) (*pb.GetPurchaseRecordResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (s *Service) ReissuePurchase(ctx context.Context, request *pb.ReissuePurchaseRequest) (*pb.Nothing, error) {
	//TODO implement me
	panic("implement me")
}
