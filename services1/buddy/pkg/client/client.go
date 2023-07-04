package client

import (
	"context"
	"time"

	"github.com/donglei1234/platform/services/common/tracing"
	"github.com/donglei1234/platform/services/common/tracing/noop"

	"github.com/donglei1234/platform/services/common/utils"
	pb "github.com/donglei1234/platform/services/proto/gen/buddy/api"

	"google.golang.org/grpc"
)

type PublicClient interface {
	AddBuddy(ctx context.Context, text string, names ...string) ([]string, error)
	RemoveBuddy(ctx context.Context, name string) error
	GetBuddies(ctx context.Context) (*pb.Buddies, error)
	WatchBuddies(ctx context.Context) (BuddyUpdateStream, error)
	ReplyAddBuddy(ctx context.Context, response bool, names ...string) ([]string, error)
	Remark(ctx context.Context, name string, remark string) error
	UpdateBuddySettings(ctx context.Context, allowToBeAdded bool) error
	GetBlockedProfiles(ctx context.Context) (*pb.ProfileIds, error)
	AddBlockedProfiles(ctx context.Context, ids ...string) error
	RemoveBlockedProfiles(ctx context.Context, ids ...string) error
	AddToRecentMet(ctx context.Context, ids ...string) error
	StarBuddy(ctx context.Context, friendValue int32, ids ...string) error
	FavoriteBuddy(ctx context.Context, isFavor bool, ids ...string) error
	CollectBuddyReward(ctx context.Context, ids ...string) (int32, error)
	IsInvited(ctx context.Context, id string) (bool, error)
	Close() error
}

type PrivateClient interface {
	// Deprecated: Do not use.
	GetBuddies(ctx context.Context, name string) (*pb.Buddies, error)
	// Deprecated: Do not use.
	WatchBuddies(ctx context.Context, name string) (BuddyUpdateStream, error)
	GetProfileBuddies(ctx context.Context, profileID, appID string) (*pb.Buddies, error)
	WatchProfileBuddies(ctx context.Context, profileID, appID string) (BuddyUpdateStream, error)
	GetProfileBlockedList(ctx context.Context, profileID, appID string) (*pb.GetProfileBlockedListResponse, error)
	Close() error
}

type BuddyUpdateStream interface {
	Recv() (*pb.BuddyChanges, error)
	grpc.ClientStream
}

type publicClient struct {
	*client
}

func NewPublicClient(target string, secure bool) (client PublicClient, err error) {
	return NewPublicClientWithTracer(noop.NewTracer(), target, secure)
}

func NewPublicClientWithTracer(tracer tracing.Tracer, target string, secure bool) (client PublicClient, err error) {
	if c, e := newClient(tracer, target, secure); e != nil {
		err = e
	} else {
		client = &publicClient{client: c}
	}

	return
}

func (c *publicClient) AddBuddy(ctx context.Context, text string, names ...string) ([]string, error) {
	cli := pb.NewPublicServiceClient(c.conn)
	if resp, err := cli.AddBuddy(ctx, &pb.AddBuddyRequest{Uid: names, ReqInfo: text}); err != nil {
		return nil, err
	} else {
		return resp.GetFailed(), nil
	}
}

func (c *publicClient) RemoveBuddy(ctx context.Context, name string) (err error) {
	cli := pb.NewPublicServiceClient(c.conn)
	_, err = cli.RemoveBuddy(ctx, &pb.RemoveBuddyRequest{Uid: name})

	return
}

func (c *publicClient) GetBuddies(ctx context.Context) (buddies *pb.Buddies, err error) {
	cli := pb.NewPublicServiceClient(c.conn)
	buddies, err = cli.GetBuddies(ctx, &pb.Nothing{})

	return
}

func (c *publicClient) WatchBuddies(ctx context.Context) (stream BuddyUpdateStream, err error) {
	cli := pb.NewPublicServiceClient(c.conn)
	stream, err = cli.WatchBuddies(ctx, &pb.Nothing{})
	return
}

func (c *publicClient) ReplyAddBuddy(ctx context.Context, response bool, names ...string) ([]string, error) {
	cli := pb.NewPublicServiceClient(c.conn)
	if resp, err := cli.ReplyAddBuddy(ctx, &pb.ReplyAddBuddyRequest{Uid: names, Response: response}); err != nil {
		return nil, err
	} else {
		return resp.GetFailed(), nil
	}
}

func (c *publicClient) Remark(ctx context.Context, name string, remark string) (err error) {
	cli := pb.NewPublicServiceClient(c.conn)
	_, err = cli.Remark(ctx, &pb.RemarkRequest{Uid: name, Remark: remark})

	return
}

func (c *publicClient) UpdateBuddySettings(ctx context.Context, allowToBeAdded bool) (err error) {
	cli := pb.NewPublicServiceClient(c.conn)
	_, err = cli.UpdateBuddySettings(ctx, &pb.UpdateBuddySettingsRequest{AllowToBeAdded: allowToBeAdded})
	return
}

func (c *publicClient) GetBlockedProfiles(ctx context.Context) (profileIds *pb.ProfileIds, err error) {
	cli := pb.NewPublicServiceClient(c.conn)
	profileIds, err = cli.GetBlockedProfiles(ctx, &pb.Nothing{})

	return
}

func (c *publicClient) IsInvited(ctx context.Context, id string) (bool, error) {
	cli := pb.NewPublicServiceClient(c.conn)
	resp, err := cli.IsInvited(ctx, &pb.IsInvitedRequest{Uid: id})
	if err != nil {
		return false, err
	}

	return resp.IsInvited, nil
}

func (c *publicClient) AddBlockedProfiles(ctx context.Context, ids ...string) (err error) {
	var profileIds []*pb.ProfileId

	for _, v := range ids {
		pi := &pb.ProfileId{
			ProfileId: v,
		}
		profileIds = append(profileIds, pi)
	}

	cli := pb.NewPublicServiceClient(c.conn)
	_, err = cli.AddBlockedProfiles(ctx, &pb.ProfileIds{ProfileIds: profileIds})

	return
}

func (c *publicClient) RemoveBlockedProfiles(ctx context.Context, ids ...string) (err error) {
	var profileIds []*pb.ProfileId

	for _, v := range ids {
		pi := &pb.ProfileId{
			ProfileId: v,
		}
		profileIds = append(profileIds, pi)
	}

	cli := pb.NewPublicServiceClient(c.conn)
	_, err = cli.RemoveBlockedProfiles(ctx, &pb.ProfileIds{ProfileIds: profileIds})

	return
}

func (c *publicClient) AddToRecentMet(ctx context.Context, ids ...string) (err error) {
	var profileIds []*pb.ProfileId

	for _, v := range ids {
		pi := &pb.ProfileId{
			ProfileId: v,
		}
		profileIds = append(profileIds, pi)
	}

	cli := pb.NewPublicServiceClient(c.conn)
	_, err = cli.AddToRecentMet(ctx, &pb.ProfileIds{ProfileIds: profileIds})

	return
}
func (c *publicClient) StarBuddy(ctx context.Context, friendValue int32, ids ...string) error {
	cli := pb.NewPublicServiceClient(c.conn)
	if _, err := cli.StarBuddy(ctx, &pb.StarBuddyRequest{
		Uid:         ids,
		FriendValue: friendValue,
	}); err != nil {
		return err
	}
	return nil
}

func (c *publicClient) FavoriteBuddy(ctx context.Context, isFavor bool, ids ...string) error {
	cli := pb.NewPublicServiceClient(c.conn)
	if _, err := cli.FavoriteBuddy(ctx, &pb.FavoriteBuddyRequest{
		Uid:        ids,
		IsFavorite: isFavor,
	}); err != nil {
		return err
	}
	return nil
}

func (c *publicClient) CollectBuddyReward(ctx context.Context, ids ...string) (int32, error) {
	cli := pb.NewPublicServiceClient(c.conn)
	profileIds := make([]*pb.ProfileId, len(ids))
	for i, id := range ids {
		profileIds[i] = &pb.ProfileId{ProfileId: id}
	}

	if resp, err := cli.CollectBuddyReward(ctx, &pb.ProfileIds{
		ProfileIds: profileIds,
	}); err != nil {
		return 0, err
	} else {
		return resp.Num, nil
	}
}

type privateClient struct {
	*client
}

func NewPrivateClient(target string, secure bool) (client PrivateClient, err error) {
	return NewPrivateClientWithTracer(noop.NewTracer(), target, secure)
}

func NewPrivateClientWithTracer(tracer tracing.Tracer, target string, secure bool) (client PrivateClient, err error) {
	if c, e := newClient(tracer, target, secure); e != nil {
		err = e
	} else {
		client = &privateClient{c}
	}

	return
}

func (c *privateClient) GetProfileBuddies(ctx context.Context, profileID, appID string) (buddies *pb.Buddies, err error) {
	cli := pb.NewPrivateServiceClient(c.conn)
	buddies, err = cli.GetProfileBuddies(
		ctx,
		&pb.GetProfileBuddiesRequest{
			ProfileId: &pb.ProfileId{ProfileId: profileID},
			AppId:     &pb.AppId{AppId: appID},
		},
	)
	return
}

func (c *privateClient) WatchProfileBuddies(ctx context.Context, profileID, appID string) (stream BuddyUpdateStream, err error) {
	//cli := pb.NewPrivateServiceClient(c.conn)
	//stream, err = cli.WatchProfileBuddies(
	//	ctx,
	//	&pb.WatchProfileBuddiesRequest{
	//		ProfileId: &pb.ProfileId{ProfileId: profileID},
	//		AppId:     &pb.AppId{AppId: appID},
	//	},
	//)
	return
}

func (c *privateClient) GetProfileBlockedList(
	ctx context.Context,
	profileID,
	appID string,
) (buddies *pb.GetProfileBlockedListResponse, err error) {
	cli := pb.NewPrivateServiceClient(c.conn)
	buddies, err = cli.GetProfileBlockedList(
		ctx,
		&pb.GetProfileBlockedListRequest{
			ProfileId: &pb.ProfileId{ProfileId: profileID},
			AppId:     &pb.AppId{AppId: appID},
		},
	)
	return
}

// Deprecated: Do not use.
func (c *privateClient) GetBuddies(ctx context.Context, name string) (buddies *pb.Buddies, err error) {
	cli := pb.NewPrivateServiceClient(c.conn)
	buddies, err = cli.GetBuddies(ctx, &pb.Name{Uid: name})
	return
}

// Deprecated: Do not use.
func (c *privateClient) WatchBuddies(ctx context.Context, name string) (stream BuddyUpdateStream, err error) {
	//cli := pb.NewPrivateServiceClient(c.conn)
	//stream, err = cli.WatchBuddies(ctx, &pb.Uid{Uid: name})
	return
}

type client struct {
	conn *grpc.ClientConn
}

func newClient(tracer tracing.Tracer, target string, secure bool) (cli *client, err error) {
	if conn, e := utils.Dial(
		target,
		utils.TransportSecurity(secure),
		grpc.WithBackoffMaxDelay(5*time.Second),
		grpc.WithStreamInterceptor(tracer.StreamClientInterceptor()),
		grpc.WithUnaryInterceptor(tracer.UnaryClientInterceptor()),
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
