package client

import (
	"context"
	"time"

	"github.com/donglei1234/platform/services/common/utils"
	"go.uber.org/zap"
	"google.golang.org/grpc"

	pb "github.com/donglei1234/platform/services/leaderboard/generated/grpc/go/leaderboard/api"
)

const (
	SUM        = "SUM"
	BETTER     = "BETTER"
	LAST       = "LAST"
	DESCENDING = "DESCENDING"
	ASCENDING  = "ASCENDING"
)

type LeaderboardClient interface {
	GetTopK(ctx context.Context, appId, listId string, k int32) (list *pb.GetLeaderboardResponse, err error)
	GetMToN(ctx context.Context, appId, listId string, m, n int32) (list *pb.GetLeaderboardResponse, err error)
	GetIdRank(ctx context.Context, appId, listId, id string) (rank *pb.GetIdRankResponse, err error)
	UpdateScore(ctx context.Context, appId, listId, id string, score int32) (session *pb.NothingResponse, err error)
	DeleteMember(ctx context.Context, appId, listId, id string) (session *pb.NothingResponse, err error)
	NewLeaderboard(ctx context.Context, appId, listId, method, order string, resetTime,
		updateTime int32) (session *pb.NothingResponse, err error)
	GetLeaderBoardSize(ctx context.Context, appId, listId string) (size *pb.GetLeaderBoardSizeResponse, err error)
	ResetLeaderboard(ctx context.Context, appId, listId string) (session *pb.NothingResponse, err error)
	Close() error
}

type publicClient struct {
	l *zap.Logger
	*client
}

func (p *publicClient) GetTopK(ctx context.Context, appId, listId string, k int32) (list *pb.GetLeaderboardResponse,
	err error) {
	cli := pb.NewLeaderboardServiceClient(p.conn)
	list, err = cli.GetTopK(ctx, &pb.GetTopKRequest{
		AppId:           appId,
		LeaderboardName: listId,
		K:               k,
	})
	return
}

func (p *publicClient) GetMToN(ctx context.Context, appId, listId string, m,
	n int32) (list *pb.GetLeaderboardResponse, err error) {
	cli := pb.NewLeaderboardServiceClient(p.conn)
	list, err = cli.GetMToN(ctx, &pb.GetMToNRequest{
		AppId:           appId,
		LeaderboardName: listId,
		M:               m,
		N:               n,
	})
	return
}

func (p *publicClient) GetIdRank(ctx context.Context, appId, listId, id string) (rank *pb.GetIdRankResponse,
	err error) {
	cli := pb.NewLeaderboardServiceClient(p.conn)
	rank, err = cli.GetIdRank(ctx, &pb.GetIdRankRequest{
		AppId:           appId,
		LeaderboardName: listId,
		Id:              id,
	})
	return
}

func (p *publicClient) UpdateScore(ctx context.Context, appId, listId, id string,
	score int32) (session *pb.NothingResponse, err error) {
	cli := pb.NewLeaderboardServiceClient(p.conn)
	session, err = cli.UpdateScore(ctx, &pb.UpdateScoreRequest{
		AppId:           appId,
		LeaderboardName: listId,
		Id:              id,
		Score:           score,
	})
	return
}

func (p *publicClient) DeleteMember(ctx context.Context, appId, listId, id string) (session *pb.NothingResponse,
	err error) {
	cli := pb.NewLeaderboardServiceClient(p.conn)
	session, err = cli.DeleteMember(ctx, &pb.GetIdRankRequest{
		AppId:           appId,
		LeaderboardName: listId,
		Id:              id,
	})
	return
}

func (p *publicClient) NewLeaderboard(ctx context.Context, appId, listId, method, order string, resetTime,
	updateTime int32) (session *pb.NothingResponse, err error) {
	cli := pb.NewLeaderboardServiceClient(p.conn)
	var m pb.NewLeaderboardRequest_MethodType
	var o pb.NewLeaderboardRequest_OrderType
	switch method {
	case SUM:
		m = pb.NewLeaderboardRequest_SUM
	case BETTER:
		m = pb.NewLeaderboardRequest_BETTER
	case LAST:
		m = pb.NewLeaderboardRequest_LAST
	}
	switch order {
	case DESCENDING:
		o = pb.NewLeaderboardRequest_DESCENDING
	case ASCENDING:
		o = pb.NewLeaderboardRequest_ASCENDING
	}
	session, err = cli.NewLeaderboard(ctx, &pb.NewLeaderboardRequest{
		AppId:           appId,
		LeaderboardName: listId,
		Method:          m,
		Order:           o,
		ResetTime:       resetTime,
		UpdateTime:      updateTime,
	})
	return
}

func (p *publicClient) GetLeaderBoardSize(ctx context.Context, appId,
	listId string) (size *pb.GetLeaderBoardSizeResponse, err error) {
	cli := pb.NewLeaderboardServiceClient(p.conn)
	size, err = cli.GetLeaderBoardSize(ctx, &pb.LeaderboardRequest{
		AppId:           appId,
		LeaderboardName: listId,
	})
	return
}

func (p *publicClient) ResetLeaderboard(ctx context.Context, appId, listId string) (session *pb.NothingResponse,
	err error) {
	cli := pb.NewLeaderboardServiceClient(p.conn)
	session, err = cli.ResetLeaderboard(ctx, &pb.LeaderboardRequest{
		AppId:           appId,
		LeaderboardName: listId,
	})
	return
}

func NewLeaderboardClient(l *zap.Logger, target string, secure bool) (client LeaderboardClient, err error) {
	if c, e := newClient(target, secure); e != nil {
		err = e
	} else {
		client = &publicClient{client: c, l: l}
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
