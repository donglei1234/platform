package client

import (
	"context"
	"fmt"
	"time"

	"go.uber.org/zap"
	"google.golang.org/grpc"

	"github.com/donglei1234/platform/services/common/utils"
	pb "github.com/donglei1234/platform/services/guild/generated/grpc/go/guild/api"
)

const master = "0"

type GuildClient interface {
	JoinGuild(ctx context.Context, appid, userid, guildid, profileid, guildattr string) error
	CreateGuild(ctx context.Context, name, notice, icon, attribute,
		appid, userid string) (string, error)
	ModifyGuild(ctx context.Context, guildid, name, notice, icon, attribute,
		appid, userid string, profileid string) error
	SearchGuild(ctx context.Context, appid, searchinput string, number int64) (*pb.SearchResponse, error)
	DeleteGuild(ctx context.Context, appid, guildid string) error
	ChangeMemberGuild(ctx context.Context, appid, memberid, profileid, guildattr, extra_appid, guildid string) error
	GetMember(ctx context.Context, appid, idx string) (*pb.UserListResponse, error)
	Apply(ctx context.Context, appid, userid, guildid string) error
	Reply(ctx context.Context, appid, applyid string, mode bool, guildid, profileid, guildattr string) error
	GetApply(ctx context.Context, appid, guildid string) (*pb.GetApplyResponse, error)
	Close() error
}

type publicClient struct {
	l *zap.Logger
	*client
}

func NewGuildClient(l *zap.Logger, target string, secure bool) (client GuildClient, err error) {
	if c, e := newClient(target, secure); e != nil {
		err = e
	} else {
		client = &publicClient{client: c, l: l}
	}

	return
}

func (u *publicClient) JoinGuild(ctx context.Context, appid, userid,
	guildid, profileid, guildattr string) (err error) {
	cli := pb.NewGuildClient(u.conn)
	_, err = cli.JoinGuild(ctx, &pb.JoinRequest{
		AppId:     appid,
		GuildId:   guildid,
		ProfileId: profileid,
		GuildAttr: guildattr,
		UserMessage: &pb.UserMessage{
			UserId: userid,
		},
	})
	return err
}

func (u *publicClient) CreateGuild(ctx context.Context, name, notice, icon, attribute,
	appid, userid string) (idx string, err error) {
	var tmp *pb.GuildMessage
	cli := pb.NewGuildClient(u.conn)
	tmp, err = cli.CreateGuild(ctx, &pb.ModifyGuildRequest{
		AppId:     appid,
		ProfileId: master,
		Mode:      pb.ModifyGuildRequest_create,
		UMessage: &pb.UserMessage{
			UserId: userid,
		},
		Message: &pb.GuildMessage{
			Name:      name,
			Notice:    notice,
			Icon:      icon,
			Attribute: attribute,
		},
	})
	idx = tmp.GuildId
	return
}

func (u *publicClient) ModifyGuild(ctx context.Context, guildid, name, notice, icon, attribute,
	appid, userid, profileid string) (err error) {
	cli := pb.NewGuildClient(u.conn)
	_, err = cli.ModifyGuild(ctx, &pb.ModifyGuildRequest{
		AppId:     appid,
		ProfileId: profileid,
		Mode:      pb.ModifyGuildRequest_modify,
		UMessage: &pb.UserMessage{
			UserId: userid,
		},
		Message: &pb.GuildMessage{
			GuildId:   guildid,
			Name:      name,
			Notice:    notice,
			Icon:      icon,
			Attribute: attribute,
		},
	})
	return
}

func (u *publicClient) SearchGuild(ctx context.Context, appid, searchinput string, number int64) (token *pb.SearchResponse, err error) {
	cli := pb.NewGuildClient(u.conn)
	token, err = cli.SearchGuild(ctx, &pb.SearchRequest{
		AppId:       appid,
		SearchInput: searchinput,
		Number:      number,
	})
	fmt.Println(token.Guilds)
	return
}

func (u *publicClient) DeleteGuild(ctx context.Context, appid, guildid string) (err error) {
	cli := pb.NewGuildClient(u.conn)
	_, err = cli.DeleteGuild(ctx, &pb.DelRequest{
		AppId:   appid,
		GuildId: guildid,
	})
	return
}

func (u *publicClient) ChangeMemberGuild(ctx context.Context, appid, memberid, profileid, guildattr, extra_appid, guildid string) (err error) {
	cli := pb.NewGuildClient(u.conn)
	_, err = cli.ChangeMemberGuild(ctx, &pb.UpdateProfileRequest{
		AppId:       appid,
		ProfileId:   profileid,
		MemberId:    memberid,
		GuildAttr:   guildattr,
		ExtraUserId: extra_appid,
		GuildId:     guildid,
	})
	return
}

func (u *publicClient) GetMember(ctx context.Context, appid, idx string) (token *pb.UserListResponse, err error) {
	cli := pb.NewGuildClient(u.conn)
	token, err = cli.GetMember(ctx, &pb.UserListRequest{
		Idx:   idx,
		AppId: appid,
	})
	return
}

func (u *publicClient) Apply(ctx context.Context, appid, userid, guildid string) (err error) {
	cli := pb.NewGuildClient(u.conn)
	_, err = cli.Apply(ctx, &pb.ApplyRequest{
		UserId:  userid,
		GuildId: guildid,
		AppId:   appid,
	})
	return
}

func (u *publicClient) Reply(ctx context.Context, appid, applyid string, mode bool, guildid, profileid, guildattr string) (err error) {
	cli := pb.NewGuildClient(u.conn)
	_, err = cli.Reply(ctx, &pb.ReplyRequest{
		AppId:   appid,
		ApplyId: applyid,
		GuildId: guildid,
		Mode:    mode,
		IfJoin: &pb.JoinRequest{
			GuildId:   guildid,
			ProfileId: profileid,
			GuildAttr: guildattr,
			AppId:     appid,
			UserMessage: &pb.UserMessage{
				UserId: applyid,
			},
		},
	})
	return
}

func (u *publicClient) GetApply(ctx context.Context, appid, guildid string) (token *pb.GetApplyResponse, err error) {
	cli := pb.NewGuildClient(u.conn)
	token, err = cli.GetApply(ctx, &pb.GetApplyRequest{
		GuildId: guildid,
		AppId:   appid})
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
