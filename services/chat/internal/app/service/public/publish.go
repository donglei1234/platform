package public

import (
	"github.com/golang/protobuf/proto"
	"go.uber.org/zap"

	"github.com/donglei1234/platform/services/common/mq"
	pb "github.com/donglei1234/platform/services/proto/gen/chat/api"
)

func (s *Service) publish(logger *zap.Logger, in *pb.ChatRequest_Chat) error {
	msg := in.Message
	dest := msg.Destination
	if dest.Channel == pb.ChatChannel_Private || dest.Channel == pb.ChatChannel_Room {
		if dest.Id == "" {
			logger.Error("Publish message data error, Destination ID empty")
			return ErrChatMessageFailure
		}
	}
	if _, ok := pb.ChatChannel_value[dest.Channel.String()]; !ok {
		logger.Error("Publish message data error, Channel not supported", zap.String("Channel", dest.Channel.String()))
		return ErrDestinationFailure
	}
	if dest.Channel == pb.ChatChannel_Private {
		if err := s.msgCache.PushChatMessage(dest.GetId(), msg.GetMessage()...); err != nil {
			return err
		}
	}
	topic := s.chatRoomTopic(dest.Channel, dest.Id)
	if data, err := proto.Marshal(msg); err != nil {
		return err
	} else if err := s.mq.Publish(
		mq.NatsProtocol+topic,
		mq.WithBytes(data),
	); err != nil {
		return ErrPublishChatFailed
	} else {
		logger.Info("MQ Publish ",
			zap.String("Topic", topic),
		)
		return nil
	}
}
