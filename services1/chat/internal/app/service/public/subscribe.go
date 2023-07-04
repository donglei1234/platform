package public

import (
	"github.com/donglei1234/platform/services/common/mq"
	pb "github.com/donglei1234/platform/services/proto/gen/chat/api"
	"github.com/golang/protobuf/proto"
	"go.uber.org/zap"
)

func (s *Service) createSubCallback(stream ChatServer) mq.SubResponseHandler {
	return func(msg mq.Message, err error) mq.ConsumptionCode {
		message := &pb.ChatMessage{}
		if err := proto.Unmarshal(msg.Data(), message); err != nil {
			return mq.ConsumeNackPersistentFailure
		}
		if message.Destination.Channel == pb.ChatChannel_Private {
			if msg, err := s.msgCache.GetAndDeleteChatMessage(message.Destination.Id); err != nil {
				return mq.ConsumeNackPersistentFailure
			} else {
				message.Message = msg
			}
		}
		if err := stream.Send(&pb.ChatResponse{Kind: &pb.ChatResponse_Message{Message: message}}); err != nil {
			return mq.ConsumeNackPersistentFailure
		}
		return mq.ConsumeAck
	}
}

func (s *Service) getOffLineMsg(uid string, stream ChatServer) error {
	if res, err := s.msgCache.GetAndDeleteChatMessage(uid); err != nil {
		return err
	} else if len(res) != 0 {
		if err := stream.Send(&pb.ChatResponse{
			Kind: &pb.ChatResponse_Message{
				Message: &pb.ChatMessage{
					Destination: &pb.Destination{
						Channel: pb.ChatChannel_Private,
						Id:      uid,
					},
					Message: res,
				},
			},
		}); err != nil {
			return err
		}
	}
	return nil
}

func (s *Service) subscribe(
	logger *zap.Logger,
	stream ChatServer,
	client *Client,
	channel pb.ChatChannel,
	id string,
) error {
	topic := s.chatRoomTopic(channel, id)
	if sub, err := s.mq.Subscribe(mq.NatsProtocol+topic, s.createSubCallback(stream)); err != nil {
		logger.Warn("Fast MQ Subscribe failed",
			zap.String("Topic", topic),
			zap.Error(err),
		)
		return ErrSubscribeChatFailed
	} else {
		client.AddSubscription(topic, sub)
		logger.Info("MQ subscribe ",
			zap.String("Topic", topic),
		)
		return nil
	}

}

func (s *Service) unSubscribe(
	client *Client,
	in *pb.ChatRequest_UnSubscribe,
) error {
	dest := in.Destination
	topic := s.chatRoomTopic(dest.Channel, dest.GetId())
	client.Unsub(topic)
	return nil
}
