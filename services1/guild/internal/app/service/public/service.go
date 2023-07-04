package public

import (
	"context"

	"go.uber.org/fx"
	"go.uber.org/zap"

	"github.com/donglei1234/platform/services/common/access"
	"github.com/donglei1234/platform/services/common/fx/fxsvcapp"
	"github.com/donglei1234/platform/services/common/service"
	pb "github.com/donglei1234/platform/services/guild/generated/grpc/go/guild/api"
	"github.com/donglei1234/platform/services/guild/internal/app/db"
	nfx "github.com/donglei1234/platform/services/guild/pkg/fx"
	"github.com/donglei1234/platform/services/guild/pkg/metadata"
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
	pb.RegisterGuildServer(server.GrpcServer(), s)
	return nil
}

func (s *Service) RegisterWithGatewayServer(server service.HasGatewayServer) error {
	if err := pb.RegisterGuildHandlerFromEndpoint(
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
	listeningUrl string,
) (result *Service, err error) {
	if err != nil {
		return nil, err
	}
	if err != nil {
		return nil, err
	}
	result = &Service{
		appId:      metadata.AppId,
		logger:     l,
		deployment: deployment,
		db:         db.OpenDatabase(l, redisUrl, redisPwd),
		url:        listeningUrl,
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
		fs nfx.GuildSettings,
	) (out service.GrpcServiceFactory, gatewayOut service.GatewayServiceFactory, err error) {
		if svc, e := NewService(
			l,
			s.Deployment,
			s.RedisUrl,
			s.RedisPwd,
			fs.GuildUrl,
		); e != nil {
			err = e
		} else {
			err = out.Execute(svc)
			err = gatewayOut.Execute(svc)
		}
		return
	},
)
