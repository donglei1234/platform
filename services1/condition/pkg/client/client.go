package client

import (
	"context"
	"time"

	"github.com/donglei1234/platform/services/common/utils"
	pb2 "github.com/donglei1234/platform/services/condition/gen/condition/api"
	"google.golang.org/grpc"
)

type ConditionClient interface {
	Watch(ctx context.Context) (pb2.ConditionService_WatchClient, error)
	Update(ctx context.Context, update *pb2.UpdateRequest) error
	Register(ctx context.Context, condition *pb2.RegisterRequest) error
	Unregister(ctx context.Context, condition *pb2.UnRegisterRequest) error
	Close() error
}

type publicClient struct {
	*client
}

func NewConditionClient(target string, secure bool) (client ConditionClient, err error) {
	if c, e := newClient(target, secure); e != nil {
		err = e
	} else {
		client = &publicClient{client: c}
	}

	return
}

func (c *publicClient) Watch(ctx context.Context) (pb2.ConditionService_WatchClient, error) {
	cli := pb2.NewConditionServiceClient(c.conn)
	stream, err := cli.Watch(ctx, &pb2.Nothing{})
	if err != nil {
		return nil, err
	}
	return stream, nil
}

func (c *publicClient) Update(ctx context.Context, update *pb2.UpdateRequest) error {
	cli := pb2.NewConditionServiceClient(c.conn)
	_, err := cli.Update(ctx, update)
	if err != nil {
		return err
	}
	return nil
}

func (c *publicClient) Register(ctx context.Context, condition *pb2.RegisterRequest) error {
	cli := pb2.NewConditionServiceClient(c.conn)
	_, err := cli.Register(ctx, condition)
	if err != nil {
		return err
	}
	return nil

}

func (c *publicClient) Unregister(ctx context.Context, req *pb2.UnRegisterRequest) error {
	cli := pb2.NewConditionServiceClient(c.conn)
	_, err := cli.Unregister(ctx, req)
	if err != nil {
		return err
	}
	return nil
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
