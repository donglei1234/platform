package client

import (
	"context"
	pb "github.com/donglei1234/platform/services/proto/gen/mail/api"
	"time"

	"github.com/donglei1234/platform/services/common/utils"
	"google.golang.org/grpc"
)

type MailClient interface {
	Watch(ctx context.Context) (pb.MailService_WatchClient, error)
	SendMail(ctx context.Context, mail *pb.Mail, targets []string) (*pb.Nothing, error)
	UpdateMail(ctx context.Context, id int64, status pb.MailStatus) ([]*pb.MailReward, error)
	Close() error
}

type publicClient struct {
	*client
}

func (c *publicClient) UpdateMail(ctx context.Context, id int64, status pb.MailStatus) (
	[]*pb.MailReward, error) {
	cli := pb.NewMailServiceClient(c.conn)
	req := &pb.UpdateMailRequest{
		Id:     id,
		Status: status,
	}
	resp, err := cli.UpdateMail(ctx, req)
	if err != nil {
		return nil, err
	}

	return resp.Rewards, nil
}

func (c *publicClient) Watch(ctx context.Context) (pb.MailService_WatchClient, error) {
	cli := pb.NewMailServiceClient(c.conn)
	stream, err := cli.Watch(ctx, &pb.Nothing{})
	if err != nil {
		return nil, err
	}

	return stream, nil
}

func (c *publicClient) SendMail(ctx context.Context, mail *pb.Mail, targets []string) (*pb.Nothing, error) {
	cli := pb.NewMailServiceClient(c.conn)
	req := &pb.SendMailRequest{
		Mail:    mail,
		Targets: targets,
	}
	return cli.SendMail(ctx, req)
}

func NewMailClient(target string, secure bool) (client MailClient, err error) {
	if c, e := newClient(target, secure); e != nil {
		err = e
	} else {
		client = &publicClient{client: c}
	}

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
