package public

import (
	"context"
	pb "github.com/donglei1234/platform/services/guild/generated/grpc/go/guild/api"
	"github.com/donglei1234/platform/services/guild/internal/app/db"
)

func (s *Service) JoinGuild(ctx context.Context, request *pb.JoinRequest) (response *pb.NothingResponse, err error) {
	err = s.db.JoinGuild(request.AppId, request.GuildId, request.ProfileId, request.UserMessage.UserId, request.GuildAttr)
	return &pb.NothingResponse{}, err
}

func (s *Service) CreateGuild(ctx context.Context, request *pb.ModifyGuildRequest) (response *pb.GuildMessage, err error) {
	var idx string
	idx, err = s.db.Cre_Mod_Guild(request.AppId, request.ProfileId, request.UMessage.UserId, request.Message.GuildId,
		request.Message.Name, request.Message.Attribute, request.Message.Icon, request.Message.Notice, request.Mode.String())
	return &pb.GuildMessage{GuildId: idx}, err
}

func (s *Service) ModifyGuild(ctx context.Context, request *pb.ModifyGuildRequest) (response *pb.NothingResponse, err error) {
	_, err = s.db.Cre_Mod_Guild(request.AppId, request.ProfileId, request.UMessage.UserId, request.Message.GuildId,
		request.Message.Name, request.Message.Attribute, request.Message.Icon, request.Message.Notice, request.Mode.String())
	return &pb.NothingResponse{}, err
}

func (s *Service) SearchGuild(ctx context.Context, request *pb.SearchRequest) (*pb.SearchResponse, error) {
	var guilds []*pb.GuildMessage
	res, err := s.db.SearchGuild(request.AppId, request.SearchInput, request.Number)
	for _, v := range res {
		var tmp = pb.GuildMessage{
			GuildId: v.Uguildid,
			Name:    v.Uguildname,
			Notice:  v.Uguildnotice,
			Icon:    v.Uguildicon,
		}
		if len(guilds) == 0 {
			guilds = []*pb.GuildMessage{&tmp}
		} else {
			guilds = append(guilds, &tmp)
		}
	}
	response := pb.SearchResponse{Guilds: guilds}
	return &response, err
}

func (s *Service) DeleteGuild(ctx context.Context, request *pb.DelRequest) (response *pb.NothingResponse, err error) {
	err = s.db.DelGuild(request.AppId, request.GuildId)
	return &pb.NothingResponse{}, err
}

func (s *Service) ChangeMemberGuild(ctx context.Context, request *pb.UpdateProfileRequest) (response *pb.NothingResponse, err error) {
	err = s.db.ChangeMemberGuild(request.AppId, request.GuildId, request.ProfileId, request.GuildAttr, request.MemberId, request.ExtraUserId)
	return &pb.NothingResponse{}, err
}

func (s *Service) GetMember(ctx context.Context, request *pb.UserListRequest) (response *pb.UserListResponse, err error) {
	var ret []*db.UMember
	ret, err = s.db.GetMemberList(request.AppId, request.Idx)
	var mem_list []*pb.Member
	for _, v := range ret {
		var mem *pb.Member
		mem.UserMessage.UserId = v.UserId
		mem.ProfileId = v.ProfileId
		mem_list = append(mem_list, mem)
	}
	response.Users = mem_list
	return response, err
}

func (s *Service) Apply(ctx context.Context, request *pb.ApplyRequest) (response *pb.NothingResponse, err error) {
	err = s.db.Apply(request.AppId, request.UserId, request.GuildId)
	return &pb.NothingResponse{}, err
}

func (s *Service) Reply(ctx context.Context, request *pb.ReplyRequest) (response *pb.NothingResponse, err error) {
	err = s.db.Reply(request.AppId, request.ApplyId, request.GuildId, request.Mode, request.IfJoin.ProfileId,
		request.IfJoin.UserMessage.UserId, request.IfJoin.GuildAttr)
	return &pb.NothingResponse{}, err
}

func (s *Service) GetApply(ctx context.Context, request *pb.GetApplyRequest) (response *pb.GetApplyResponse, err error) {
	response, err = s.db.GetApply(request.AppId, request.GuildId)
	return response, err
}
