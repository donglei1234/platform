package public

import (
	"context"

	"go.uber.org/fx"
	"go.uber.org/zap"

	"github.com/donglei1234/platform/services/common/access"
	"github.com/donglei1234/platform/services/common/fx/fxsvcapp"
	"github.com/donglei1234/platform/services/common/service"
	pb "github.com/donglei1234/platform/services/leaderboard/generated/grpc/go/leaderboard/api"
	"github.com/donglei1234/platform/services/leaderboard/internal/app/db"
	bfx "github.com/donglei1234/platform/services/leaderboard/pkg/fx"
	"github.com/donglei1234/platform/services/leaderboard/pkg/metadata"
)

type Service struct {
	service.TcpTransport
	appId      string
	logger     *zap.Logger
	deployment string
	db         *db.Database
	url        string
}

func (s *Service) RegisterWithGrpcServer(server service.HasGrpcServer) error {
	pb.RegisterLeaderboardServiceServer(server.GrpcServer(), s)
	return nil
}

func (s *Service) RegisterWithGatewayServer(server service.HasGatewayServer) error {
	if err := pb.RegisterLeaderboardServiceHandlerFromEndpoint(
		context.Background(), server.GatewayRuntimeMux(), s.url, server.GatewayOption()); err != nil {
		return err
	}
	return nil
}

func NewService(
	l *zap.Logger,
	deployment string,
	redisUrl string,
	redisPwd string,
	url string,
) (result *Service, err error) {
	result = &Service{
		appId:      metadata.AppId,
		logger:     l,
		deployment: deployment,
		db:         db.OpenDatabase(l, redisUrl, redisPwd),
		url:        url,
	}
	return
}

func (s *Service) AccessLevel() access.AccessLevel {
	return access.AccessUndefined
}

var ServiceModule = fx.Provide(
	func(
		l *zap.Logger,
		s fxsvcapp.GlobalSettings,
		fs bfx.LeaderboardSettings,
	) (out service.GrpcServiceFactory, gatewayOut service.GatewayServiceFactory, err error) {
		if svc, e := NewService(
			l,
			s.Deployment,
			s.RedisUrl,
			s.RedisPwd,
			fs.LeaderboardUrl,
		); e != nil {
			err = e
		} else {
			err = out.Execute(svc)
			err = gatewayOut.Execute(svc)
		}
		return
	},
)
