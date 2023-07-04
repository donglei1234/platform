package client

import (
	"context"
	chat "github.com/donglei1234/platform/services/proto/gen/chat/api"
	"time"

	"github.com/donglei1234/platform/services/common/utils"
	"google.golang.org/grpc"
)

type Client struct {
	conn *grpc.ClientConn
}

type ChatStream interface {
	// who subscribe which channel
	Subscribe(profileID string, target *chat.Destination) error
	// who unSubscribe which channel
	UnSubscribe(profileID string, target *chat.Destination) error

	SendMessage(content *chat.ChatMessage_Message, dest *chat.Destination) error
	SendIn(*chat.ChatRequest) error
	RecvIn() (*chat.ChatResponse, error)
	grpc.ClientStream
}

type chatStream struct {
	chat.ChatService_ChatClient
}

func NewClient(target string, secure bool) (client *Client, err error) {
	if conn, e := utils.Dial(
		target,
		utils.TransportSecurity(secure),
		grpc.WithBackoffMaxDelay(5*time.Second),
	); e != nil {
		err = e
	} else {
		client = &Client{conn: conn}
	}

	return
}

func (c *Client) Close() error {
	return c.conn.Close()
}

func newChatStream(chatProxyClient chat.ChatService_ChatClient) ChatStream {
	return &chatStream{chatProxyClient}
}

func CreateChatStream(client *Client, ctx context.Context) (stream ChatStream, err error) {
	cli := chat.NewChatServiceClient(client.conn)
	if s, e := cli.Chat(ctx); e != nil {
		err = e
	} else {
		stream = newChatStream(s)
	}
	return
}

func (c *chatStream) Subscribe(profileID string, target *chat.Destination) error {
	chatChannel := &chat.ChatRequest_Subscribe_{
		Subscribe: &chat.ChatRequest_Subscribe{
			ProfileID:   profileID,
			Destination: target,
		},
	}

	if err := c.Send(&chat.ChatRequest{Kind: chatChannel}); err != nil {
		return err
	}
	return nil
}

func (c *chatStream) UnSubscribe(profileID string, target *chat.Destination) error {
	chatChannel := &chat.ChatRequest_Unsubscribe{
		Unsubscribe: &chat.ChatRequest_UnSubscribe{
			ProfileID:   profileID,
			Destination: target,
		},
	}

	if err := c.Send(&chat.ChatRequest{Kind: chatChannel}); err != nil {
		return err
	}
	return nil
}

func (c *chatStream) SendMessage(content *chat.ChatMessage_Message, dest *chat.Destination) error {
	msg := &chat.ChatRequest{
		Kind: &chat.ChatRequest_Message{
			Message: &chat.ChatRequest_Chat{
				Message: &chat.ChatMessage{
					Destination: dest,
					Message:     []*chat.ChatMessage_Message{content},
				},
			},
		},
	}
	return c.Send(msg)
}

func (c *chatStream) SendIn(req *chat.ChatRequest) error {
	return c.Send(req)
}

func (c *chatStream) RecvIn() (*chat.ChatResponse, error) {
	return c.Recv()
}
