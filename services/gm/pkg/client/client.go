package client

import (
	"context"
	pb "github.com/donglei1234/platform/services/proto/gen/gm/api"
	"time"

	"github.com/donglei1234/platform/services/common/utils"
	"google.golang.org/grpc"
)

type GmClient interface {
	// SetProfilesBanStatus 封停、解封、解封时间、封停原因
	SetProfilesBanStatus(ctx context.Context, in *pb.SetProfilesBanStatusRequest) (*pb.Nothing, error)
	// GetProfilesBanStatus 获取当前封停的账号信息
	GetProfilesBanStatus(ctx context.Context, in *pb.GetProfilesBanStatusRequest) (*pb.GetProfilesBanStatusResponse, error)
	// SendBulletin 发送游戏内公告
	SendBulletin(ctx context.Context, in *pb.SendBulletinRequest) (*pb.Nothing, error)
	// GetBulletin 获取游戏当前公告信息
	GetBulletin(ctx context.Context, in *pb.GetBulletinRequest) (*pb.GetBulletinResponse, error)
	// GetPurchaseRecord 根据账号或订单号查询充值金额、充值时间、是否到账
	GetPurchaseRecord(ctx context.Context, in *pb.GetPurchaseRecordRequest) (*pb.GetPurchaseRecordResponse, error)
	// ReissuePurchase 根据掉单订单号进行补单(每个订单只可进行一次成功补单）
	ReissuePurchase(ctx context.Context, in *pb.ReissuePurchaseRequest) (*pb.Nothing, error)
	// WatchProfilesBan watch 黑名单信息变化
	WatchProfilesBan(ctx context.Context, in *pb.Nothing) (pb.GmService_WatchProfilesBanClient, error)
}

type publicClient struct {
	*client
}

func (c *publicClient) WatchProfilesBan(ctx context.Context, in *pb.Nothing) (pb.GmService_WatchProfilesBanClient, error) {
	//TODO implement me
	panic("implement me")
}

func (c *publicClient) SetProfilesBanStatus(ctx context.Context, in *pb.SetProfilesBanStatusRequest) (*pb.Nothing, error) {
	//TODO implement me
	panic("implement me")
}

func (c *publicClient) GetProfilesBanStatus(ctx context.Context, in *pb.GetProfilesBanStatusRequest) (*pb.GetProfilesBanStatusResponse, error) {
	cli := pb.NewGmServiceClient(c.conn)
	return cli.GetProfilesBanStatus(ctx, in)
}

func (c *publicClient) SendBulletin(ctx context.Context, in *pb.SendBulletinRequest) (*pb.Nothing, error) {
	//TODO implement me
	panic("implement me")
}

func (c *publicClient) GetBulletin(ctx context.Context, in *pb.GetBulletinRequest) (*pb.GetBulletinResponse, error) {
	cli := pb.NewGmServiceClient(c.conn)
	return cli.GetBulletin(ctx, in)
}

func (c *publicClient) GetPurchaseRecord(ctx context.Context, in *pb.GetPurchaseRecordRequest) (*pb.GetPurchaseRecordResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (c *publicClient) ReissuePurchase(ctx context.Context, in *pb.ReissuePurchaseRequest) (*pb.Nothing, error) {
	//TODO implement me
	panic("implement me")
}

func NewGMClient(target string, secure bool) (client GmClient, err error) {
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
