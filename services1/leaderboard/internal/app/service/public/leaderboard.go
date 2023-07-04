package public

import (
	"context"

	"go.uber.org/zap"

	pb "github.com/donglei1234/platform/services/leaderboard/generated/grpc/go/leaderboard/api"
)

type ListItem struct {
	Id    string
	Score float64
}

type ListOption struct {
	Id            string
	MethodType    int32
	OrderType     int32
	LatestVersion int32
}

func (s *Service) GetTopK(ctx context.Context, request *pb.GetTopKRequest) (*pb.GetLeaderboardResponse, error) {
	s.logger.Info("get top k")
	res, err := s.db.GetRankFromMToN(request.GetAppId(), request.GetLeaderboardName(), 0, int64(request.GetK()-1))
	if err != nil {
		s.logger.Error("get top k err", zap.Error(err))
		return &pb.GetLeaderboardResponse{}, err
	}
	var result pb.GetLeaderboardResponse
	for _, val := range res {
		result.Leaderboard = append(result.Leaderboard, &pb.GetLeaderboardResponse_Item{Id: val.Member.(string),
			Score: int32(val.Score)})
	}
	s.logger.Info("get top k success ",
		zap.String("app id ", request.GetAppId()),
		zap.String("leaderboard name ", request.GetLeaderboardName()),
		zap.String("k ", string(request.GetK())))
	return &result, nil
}

func (s *Service) GetMToN(ctx context.Context, request *pb.GetMToNRequest) (*pb.GetLeaderboardResponse, error) {
	s.logger.Info("get top k")
	res, err := s.db.GetRankFromMToN(request.GetAppId(), request.GetLeaderboardName(), int64(request.GetM()-1),
		int64(request.GetN()-1))
	if err != nil {
		s.logger.Error("get m to n err", zap.Error(err))
		return &pb.GetLeaderboardResponse{}, err
	}
	var result pb.GetLeaderboardResponse
	for _, val := range res {
		result.Leaderboard = append(result.Leaderboard, &pb.GetLeaderboardResponse_Item{Id: val.Member.(string),
			Score: int32(val.Score)})
	}
	s.logger.Info("get m to n success ",
		zap.String("app id ", request.GetAppId()),
		zap.String("leaderboard name ", request.GetLeaderboardName()),
		zap.String("m ", string(request.GetM())),
		zap.String("n ", string(request.GetN())))
	return &result, nil
}

func (s *Service) GetIdRank(ctx context.Context, request *pb.GetIdRankRequest) (*pb.GetIdRankResponse, error) {
	s.logger.Info("get id rank")
	res, err := s.db.GetRankById(request.GetAppId(), request.GetLeaderboardName(), request.GetId())
	if err != nil {
		s.logger.Error("get id rank err", zap.Error(err))
		return &pb.GetIdRankResponse{}, err
	}
	s.logger.Info("get id rank success ",
		zap.String("app id ", request.GetAppId()),
		zap.String("leaderboard name ", request.GetLeaderboardName()),
		zap.String("id  ", request.GetId()))
	return &pb.GetIdRankResponse{Rank: res + 1}, nil
}

func (s *Service) UpdateScore(ctx context.Context, request *pb.UpdateScoreRequest) (*pb.NothingResponse, error) {
	s.logger.Info("update score")
	err := s.db.UpdateScore(request.GetAppId(), request.GetLeaderboardName(), request.GetId(), float64(request.GetScore()))
	if err != nil {
		s.logger.Error("update score err", zap.Error(err))
		return &pb.NothingResponse{}, err
	}
	s.logger.Info("update score success ",
		zap.String("app id ", request.GetAppId()),
		zap.String("leaderboard name ", request.GetLeaderboardName()),
		zap.String("id  ", request.GetId()),
		zap.String("score ", string(request.GetScore())))
	return &pb.NothingResponse{}, nil
}

func (s *Service) DeleteMember(ctx context.Context, request *pb.GetIdRankRequest) (*pb.NothingResponse, error) {
	s.logger.Info("delete member")
	err := s.db.DeleteMember(request.GetAppId(), request.GetLeaderboardName(), request.GetId())
	if err != nil {
		s.logger.Error("delete member err", zap.Error(err))
		return &pb.NothingResponse{}, err
	}
	s.logger.Info("delete member success ",
		zap.String("app id ", request.GetAppId()),
		zap.String("leaderboard name ", request.GetLeaderboardName()),
		zap.String("id  ", request.GetId()))
	return &pb.NothingResponse{}, nil
}

func (s *Service) NewLeaderboard(ctx context.Context, request *pb.NewLeaderboardRequest) (*pb.NothingResponse, error) {
	s.logger.Info("new leaderboard")
	err := s.db.NewLeaderboard(request.GetAppId(), request.GetLeaderboardName(), request.GetMethod().String(),
		request.GetOrder().String(), request.GetResetTime(), request.GetUpdateTime())
	if err != nil {
		s.logger.Error("new leaderboard err", zap.Error(err))
		return &pb.NothingResponse{}, err
	}
	s.logger.Info("new leaderboard success ",
		zap.String("app id ", request.GetAppId()),
		zap.String("leaderboard name ", request.GetLeaderboardName()))
	return &pb.NothingResponse{}, nil
}

func (s *Service) GetLeaderBoardSize(ctx context.Context, request *pb.LeaderboardRequest) (*pb.GetLeaderBoardSizeResponse, error) {
	s.logger.Info("get leaderboard size")
	res, err := s.db.GetLeaderboardSize(request.GetAppId(), request.GetLeaderboardName())
	if err != nil {
		s.logger.Error("get leaderboard size err", zap.Error(err))
		return &pb.GetLeaderBoardSizeResponse{}, err
	}
	s.logger.Info("get leaderboard size success ",
		zap.String("app id ", request.GetAppId()),
		zap.String("leaderboard name ", request.GetLeaderboardName()))
	return &pb.GetLeaderBoardSizeResponse{Size: res}, nil
}

func (s *Service) ResetLeaderboard(ctx context.Context, request *pb.LeaderboardRequest) (*pb.NothingResponse, error) {
	s.logger.Info("reset leaderboard")
	err := s.db.ResetLeaderboard(request.GetAppId(), request.GetLeaderboardName())
	if err != nil {
		s.logger.Error("reset leaderboard err", zap.Error(err))
		return &pb.NothingResponse{}, err
	}
	s.logger.Info("reset leaderboard success ",
		zap.String("app id ", request.GetAppId()),
		zap.String("leaderboard name ", request.GetLeaderboardName()))
	return &pb.NothingResponse{}, nil
}
