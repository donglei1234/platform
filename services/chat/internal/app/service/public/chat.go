package public

import (
	"fmt"
	util "github.com/donglei1234/platform/services/common/utils"
	"go.uber.org/zap"

	pb "github.com/donglei1234/platform/services/proto/gen/chat/api"
)

// Chat stream handler
func (s *Service) Chat(server ChatServer) error {
	ctx := server.Context()
	logger := s.LoggerForContext(ctx)
	// 1、获取token
	token, err := util.ParseJwtToken(ctx)
	if err != nil {
		logger.Error("parse jwt token err", zap.Error(err))
		return ErrNoMetaData
	}

	// 2、解析token，获取uid
	profileId, err := util.DecodeToken(token)
	if err != nil {
		logger.Error("jwt token decode to uid err", zap.Error(err))
		return ErrNoMetaData
	}
	cli, ok := s.clients.LoadOrStore(profileId, func() *Client {
		return CreateClient(logger, profileId, server, s.chatInterval)
	}())
	if ok {
		return ErrGeneralFailure
	}
	client := cli.(*Client)
	if err := s.subPublicChannels(logger, server, client); err != nil {
		return ErrGeneralFailure
	}
	go func() {
		for {
			if req, err := server.Recv(); err != nil {
				logger.Warn("Chat receive error",
					zap.Error(err),
				)
				break
			} else {
				if err = s.receive(logger, server, client, req); err != nil {
					logger.Warn("Chat message deal error",
						zap.Error(err),
						zap.String("Message", req.String()),
					)
				} else {
					logger.Info("Chat received",
						zap.String("Message", req.String()),
					)
				}
			}
		}
	}()
	<-ctx.Done()
	s.clients.Delete(profileId)
	client.UnsubAll()
	return nil
}

func (s *Service) subPublicChannels(
	logger *zap.Logger,
	server ChatServer,
	client *Client,
) error {
	channels := []pb.ChatChannel{
		pb.ChatChannel_System,
		pb.ChatChannel_World,
	}
	for _, channel := range channels {
		if err := s.subscribe(
			logger,
			server,
			client,
			channel,
			"",
		); err != nil {
			return err
		}
	}
	if err := s.subscribe(
		logger,
		server,
		client,
		pb.ChatChannel_Private,
		client.profileID,
	); err != nil {
		return err
	}

	return nil
}

func (s *Service) receive(
	logger *zap.Logger,
	server ChatServer,
	client *Client,
	in *pb.ChatRequest,
) error {
	switch kind := in.Kind.(type) {
	case *pb.ChatRequest_Subscribe_:
		if msg := kind.Subscribe; msg == nil {
			return ErrSubscribeParamFailed
		} else if dest := msg.Destination; dest == nil {
			return ErrDestinationFailure
		} else {
			if err := s.subscribe(
				logger,
				server,
				client,
				msg.Destination.GetChannel(),
				msg.Destination.GetId(),
			); err != nil {
				return err
			}
		}

	case *pb.ChatRequest_Unsubscribe:
		if msg := kind.Unsubscribe; msg == nil {
			return ErrUnSubscribeParamFailed
		} else if dest := msg.Destination; dest == nil {
			return ErrDestinationFailure
		} else {
			if err := s.unSubscribe(client, msg); err != nil {
				return fmt.Errorf("unable to unsubscribe %w", err)
			}
		}

	case *pb.ChatRequest_Message:
		if msg := kind.Message; msg == nil {
			return ErrChatMessageFailure
		} else if chatMsg := msg.Message; chatMsg == nil {
			return ErrGeneralFailure
		} else if dest := chatMsg.Destination; dest == nil {
			return ErrDestinationFailure
		} else {
			if ok := client.CheckInterval(msg.Message.Destination.Channel.String()); !ok {
				return ErrChatIntervalFailed
			}
			if err := s.publish(logger, msg); err != nil {
				return err
			}
		}
	}
	return nil
}
