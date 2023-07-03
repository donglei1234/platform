package client

import (
	"context"
	"time"

	"go.uber.org/zap"
	"google.golang.org/grpc"

	"github.com/donglei1234/platform/services/common/utils"
	pb "github.com/donglei1234/platform/services/iap/generated/grpc/go/iap/api"
)

type IAPClient interface {
	CheckIAPToken(ctx context.Context, request *pb.IAPRequest) (response *pb.IAPResponse, err error)
	ConsumeCallBack(ctx context.Context, req *pb.IAPRequest) (response *pb.IAPResponse, err error)
	Close() error
}

type publicClient struct {
	l *zap.Logger
	*client
}

func NewIAPClient(l *zap.Logger, target string, secure bool) (client IAPClient, err error) {
	if c, e := newClient(target, secure); e != nil {
		err = e
	} else {
		client = &publicClient{client: c, l: l}
	}

	return
}

func (c *publicClient) CheckIAPToken(ctx context.Context, request *pb.IAPRequest) (response *pb.IAPResponse, err error) {
	cli := pb.NewIAPPublicClient(c.conn)
	response, err = cli.CheckIAPToken(ctx, request)

	return
}

func (c *publicClient) ConsumeCallBack(ctx context.Context, request *pb.IAPRequest) (response *pb.IAPResponse, err error) {
	cli := pb.NewIAPPublicClient(c.conn)
	response, err = cli.ConsumeCallBack(ctx, request)

	return
}

func (c *client) Close() error {
	return c.conn.Close()
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
