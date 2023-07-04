package client

import (
	"context"
	"time"

	"go.uber.org/zap"
	"google.golang.org/grpc"

	"github.com/donglei1234/platform/services/common/utils"
	pb "github.com/donglei1234/platform/services/notification/generated/grpc/go/notification/api"
)

type NotificationClient interface {
	RegisterArn(ctx context.Context, profileId, deviceToken, deviceId, region, appId, deviceType string) (token *pb.NothingResponse, err error)
	DeleteArn(ctx context.Context, delType int32, publishId, deviceId, appId string) (session *pb.NothingResponse, err error)
	PublishMessage(ctx context.Context, pubType int32, message []byte, publishId, appId string) (session *pb.NothingResponse, err error)
	SubscribeTopic(ctx context.Context, subType int32, topicName, profileId, appId string) (session *pb.NothingResponse, err error)
	Close() error
}

type publicClient struct {
	l *zap.Logger
	*client
}

func NewNotificationClient(l *zap.Logger, target string, secure bool) (client NotificationClient, err error) {
	if c, e := newClient(target, secure); e != nil {
		err = e
	} else {
		client = &publicClient{client: c, l: l}
	}

	return
}

func (c *publicClient) RegisterArn(ctx context.Context, profileId, deviceToken, deviceId, region, appId, deviceType string) (token *pb.NothingResponse, err error) {
	cli := pb.NewNotificationClient(c.conn)
	token, err = cli.RegisterArn(ctx, &pb.RegisterArnRequest{
		ProfileId:   &pb.ProfileId{ProfileId: profileId},
		DeviceToken: deviceToken,
		DeviceId:    deviceId,
		Region:      region,
		AppId:       appId,
		DeviceType:  deviceType,
	})

	return
}

func (c *publicClient) PublishMessage(ctx context.Context, pubType int32, message []byte, publishId, appId string) (session *pb.NothingResponse, err error) {
	cli := pb.NewNotificationClient(c.conn)
	session, err = cli.PublishMessage(ctx, &pb.PublishMessageRequest{
		PubType:   pb.PublishMessageRequest_PubType(pubType),
		PublishId: publishId,
		Message:   message,
		AppId:     appId,
	})

	return
}

func (c *publicClient) DeleteArn(ctx context.Context, delType int32, publishId, deviceId, appId string) (session *pb.NothingResponse, err error) {
	cli := pb.NewNotificationClient(c.conn)
	session, err = cli.DeleteArn(ctx, &pb.DeleteArnRequest{
		DelType:   pb.DeleteArnRequest_DelType(delType),
		PublishId: publishId,
		DeviceID:  deviceId,
		AppId:     appId,
	})

	return
}

func (c *publicClient) SubscribeTopic(ctx context.Context, subType int32, topicName, profileId, appId string) (session *pb.NothingResponse, err error) {
	cli := pb.NewNotificationClient(c.conn)
	session, err = cli.SubscribeTopic(ctx, &pb.SubscribeTopicRequest{
		SubType:   pb.SubscribeTopicRequest_SubType(subType),
		TopicName: topicName,
		ProfileId: &pb.ProfileId{ProfileId: profileId},
		AppId:     appId,
	})
	return
}

type client struct {
	conn *grpc.ClientConn
}

func newClient(target string, secure bool) (cli *client, err error) {
	if conn, e := utils.Dial(
		target,
		utils.TransportSecurity(secure),
		grpc.WithBackoffMaxDelay(5*time.Second),
	); e != nil {
		err = e
	} else {
		cli = &client{conn: conn}
	}

	return
}

func (c *client) Close() error {
	return c.conn.Close()
}
