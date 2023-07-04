package client

import (
	"context"
	pb "github.com/donglei1234/platform/services/proto/gen/auth/api"
	"time"

	"go.uber.org/zap"
	"google.golang.org/grpc"

	"github.com/donglei1234/platform/services/common/utils"
)

type AuthClient interface {
	Auth(ctx context.Context, name, appID string) (token *pb.AuthenticateResponse, err error)
	Bind(ctx context.Context, bindRequest *pb.BindRequest) (resp *pb.BindResponse, err error)
	AuthenticateRoom(ctx context.Context, in *pb.RoomInfo) (string, error)

	ValidateRoomToken(ctx context.Context, token string) (*pb.RoomInfo, error)
	ValidateToken(ctx context.Context, jwtToken string) (resp *pb.Session, err error)
	Close() error
}

type publicClient struct {
	l *zap.Logger
	*client
}

func (c *publicClient) AuthenticateRoom(ctx context.Context, in *pb.RoomInfo) (string, error) {
	cli := pb.NewAuthPublicClient(c.conn)
	room, err := cli.AuthenticateRoom(ctx, &pb.AuthenticateRoomRequest{RoomInfo: in})
	if err != nil {
		return "", err
	}
	return room.Token, nil
}

func (c *publicClient) ValidateRoomToken(ctx context.Context, token string) (*pb.RoomInfo, error) {
	cli := pb.NewAuthPublicClient(c.conn)
	room, err := cli.ValidateRoomToken(ctx, &pb.ValidateRoomTokenRequest{Token: token})
	if err != nil {
		return nil, err
	}
	return room.GetRoomInfo(), nil
}

func NewAuthClient(l *zap.Logger, target string, secure bool) (client AuthClient, err error) {
	if c, e := newClient(target, secure); e != nil {
		err = e
	} else {
		client = &publicClient{client: c, l: l}
	}

	return
}

func (c *publicClient) Auth(ctx context.Context, name, appID string) (token *pb.AuthenticateResponse, err error) {
	cli := pb.NewAuthPublicClient(c.conn)
	token, err = cli.Authenticate(ctx, &pb.AuthenticateRequest{Username: name, AppId: appID})

	return
}
func (c *publicClient) Bind(ctx context.Context, bindRequest *pb.BindRequest) (resp *pb.BindResponse, err error) {
	cli := pb.NewAuthPrivateClient(c.conn)
	resp, err = cli.Bind(ctx, bindRequest)
	return
}
func (c *publicClient) ValidateToken(ctx context.Context, jwtToken string) (*pb.Session, error) {
	cli := pb.NewAuthPublicClient(c.conn)
	resp, err := cli.ValidateToken(ctx, &pb.ValidateTokenRequest{JwtToken: jwtToken})
	if err != nil {
		return nil, err
	}
	return resp.Session, nil
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
