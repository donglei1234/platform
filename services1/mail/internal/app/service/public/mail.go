package public

import (
	"context"
	"github.com/donglei1234/platform/services/common/jsonx"
	"github.com/donglei1234/platform/services/common/mq"
	"github.com/donglei1234/platform/services/common/nosql/errors"
	pb "github.com/donglei1234/platform/services/proto/gen/mail/api"
	"github.com/samber/lo"
	"sort"
	"time"

	"go.uber.org/zap"

	util "github.com/donglei1234/platform/services/common/utils"
)

func (s *Service) Watch(nothing *pb.Nothing, server pb.MailService_WatchServer) error {
	// 1、获取token
	token, err := util.ParseJwtToken(server.Context())
	if err != nil {
		s.logger.Error("join mail parse jwt token err", zap.Error(err))
		return InvalidReq
	}
	// 2、解析token，获取uid
	profileId, err := util.DecodeToken(token)
	if err != nil {
		s.logger.Error("join mail jwt token decode to uid err", zap.Error(err))
		return InvalidProfileId
	}
	// 3、订阅私人邮件
	if subscribe, err := s.mq.Subscribe(mq.NatsProtocol+MailPrivateTopic.String()+profileId,
		s.handleMails(server),
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

	// 4、订阅公共邮件
	if subscribe, err := s.mq.Subscribe(mq.NatsProtocol+MailPublicTopic.String(),
		s.handleMails(server),
		mq.WithAtMostOnceDelivery(mq.GroupId(profileId)),
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

	if err := s.db.MergePublicAndPrivateMail(profileId); err != nil {
		return err
	}

	mails := s.db.GetSelfMails(profileId)
	s.checkMailLimit(profileId, mails)

	mailLst := make([]*pb.Mail, 0)
	for _, v := range mails {
		mailLst = append(mailLst, v)
	}

	err = server.Send(&pb.WatchMailResponse{Mails: mailLst})
	if err != nil {
		return err
	}
	<-server.Context().Done()
	return nil
}

var mailNumLimit = 100

func (s *Service) checkMailLimit(profileId string, mails map[int64]*pb.Mail) {
	if len(mails) <= mailNumLimit {
		return
	}
	list := lo.MapToSlice(mails, func(key int64, value *pb.Mail) *pb.Mail {
		return value
	})
	// sort mails by date
	sort.Slice(list, func(i, j int) bool {
		return list[i].Date < list[j].Date
	})

	deleteNum := len(list) - mailNumLimit

	// delete status rewarded mails
	removedNum := 0
	for i := 0; i < deleteNum; i++ {
		if list[i].Status == pb.MailStatus_REWARDED {
			delete(mails, list[i].Id)
			// remove from list
			list = append(list[:i], list[i+1:]...)
			removedNum++
		}
	}
	deleteNum -= removedNum
	if deleteNum <= 0 {
		return
	}

	for i := 0; i < deleteNum; i++ {
		delete(mails, list[i].Id)
		list = append(list[:i], list[i+1:]...)
	}

	err := s.db.AddMails(profileId, list...)
	if err != nil {
		s.logger.Error("add mails err", zap.Error(err))
	}
}

func (s *Service) handleMails(server pb.MailService_WatchServer) mq.SubResponseHandler {
	return func(msg mq.Message, err error) mq.ConsumptionCode {
		mails := make([]*pb.Mail, 0)
		if err := jsonx.Unmarshal(msg.Data(), &mails); err != nil {
			return mq.ConsumeNackPersistentFailure
		}
		err = server.Send(&pb.WatchMailResponse{Mails: mails})
		if err != nil {
			s.logger.Error("send change err", zap.Error(err))
			return mq.ConsumeNackPersistentFailure
		}
		return mq.ConsumeAck
	}
}

func (s *Service) SendMail(ctx context.Context, request *pb.SendMailRequest) (*pb.Nothing, error) {
	// 1、获取token
	_, err := util.ParseJwtToken(ctx)
	if err != nil {
		s.logger.Error("send mail parse jwt token err", zap.Error(err))
		return nil, InvalidReq
	}

	mail := s.initDefault(request.Mail)
	if len(request.Targets) == 0 {
		err := s.savePublicMails(mail)
		if err != nil {
			return nil, err
		}
		if err = s.mq.Publish(
			mq.NatsProtocol+MailPublicTopic.String(),
			mq.WithJson([]*pb.Mail{mail}),
		); err != nil {
			return nil, err
		}
	} else {
		for _, v := range request.Targets {
			err := s.savePrivateMail(v, mail)
			if err != nil {
				return nil, err
			}
			if err = s.mq.Publish(
				mq.NatsProtocol+MailPrivateTopic.String()+v,
				mq.WithJson([]*pb.Mail{mail}),
			); err != nil {
				return nil, err
			}
		}
	}

	return &pb.Nothing{}, nil
}

func (s *Service) initDefault(mail *pb.Mail) *pb.Mail {
	mail.Id = time.Now().Unix()
	if mail.Date == 0 {
		mail.Date = time.Now().Unix()
	}
	if mail.Expire <= 0 {
		DefaultExpire := time.Hour * 24 * 100
		mail.Expire = mail.Date + DefaultExpire.Milliseconds()
	} else {
		duration := time.Hour * time.Duration(mail.Expire)
		mail.Expire = mail.Date + duration.Milliseconds()
	}
	return mail
}

func (s *Service) savePublicMails(mail *pb.Mail) error {
	return s.db.PushMailToPublic(mail)
}

func (s *Service) savePrivateMail(target string, mail *pb.Mail) error {
	if target == "" {
		return errors.ErrInternal
	}

	if err := s.db.AddMails(target, mail); err != nil {
		return err
	}
	if err := s.mq.Publish(
		mq.NatsProtocol+MailPrivateTopic.String()+target,
		mq.WithJson([]*pb.Mail{mail}),
	); err != nil {
		return err
	}
	return nil
}

func (s *Service) UpdateMail(ctx context.Context, request *pb.UpdateMailRequest) (*pb.UpdateMailResponse, error) {
	// 1、获取token
	token, err := util.ParseJwtToken(ctx)
	if err != nil {
		s.logger.Error("update mail parse jwt token err", zap.Error(err))
		return nil, InvalidReq
	}
	// 2、解析token，获取uid
	profileId, err := util.DecodeToken(token)
	if err != nil {
		s.logger.Error("update mail jwt token decode to uid err", zap.Error(err))
		return nil, InvalidProfileId
	}
	if err := s.db.MergePublicAndPrivateMail(profileId); err != nil {
		return nil, err
	}

	updates := make([]*pb.Mail, 0)
	if request.Id == 0 {
		mails, err := s.db.UpdateAllStatus(profileId, request.Status)
		if err != nil {
			return nil, err
		}
		updates = mails
	} else {
		mail, err := s.db.UpdateOneStatus(profileId, request.Id, request.Status)
		if err != nil {
			return nil, err
		}
		updates = append(updates, mail)
	}

	if err = s.mq.Publish(
		mq.NatsProtocol+MailPrivateTopic.String()+profileId,
		mq.WithJson(updates),
	); err != nil {
		return nil, err
	}
	rewards := make([]*pb.MailReward, 0)
	if request.Status == pb.MailStatus_REWARDED {
		for _, v := range updates {
			if v.Rewards == nil {
				continue
			}
			rewards = append(rewards, v.Rewards...)
		}
	}
	return &pb.UpdateMailResponse{Rewards: rewards}, nil
}
