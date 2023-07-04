package public

import (
	"context"

	"go.uber.org/fx"
	"go.uber.org/zap"

	"github.com/donglei1234/platform/services/common/access"
	"github.com/donglei1234/platform/services/common/fx/fxsvcapp"
	"github.com/donglei1234/platform/services/common/service"
	pb "github.com/donglei1234/platform/services/iap/generated/grpc/go/iap/api"
	"github.com/donglei1234/platform/services/iap/internal/app/db"
	afx "github.com/donglei1234/platform/services/iap/pkg/fx"
	"github.com/donglei1234/platform/services/iap/pkg/metadata"
)

type Service struct {
	service.TcpTransport
	appId         string
	logger        *zap.Logger
	deployment    string
	db            *db.Database
	iapCredential string
	statHost      string
	url           string
}

func (s *Service) RegisterWithGrpcServer(server service.HasGrpcServer) error {
	pb.RegisterIAPPublicServer(server.GrpcServer(), s)
	return nil
}

func (s *Service) RegisterWithGatewayServer(server service.HasGatewayServer) error {
	if err := pb.RegisterIAPPublicHandlerFromEndpoint(
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
	iapCredential string,
	statHost string,
	url string,
) (result *Service, err error) {
	result = &Service{
		appId:         metadata.AppId,
		logger:        l,
		deployment:    deployment,
		db:            db.OpenDatabase(l, redisUrl, redisPwd),
		iapCredential: iapCredential,
		statHost:      statHost,
		url:           url,
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
		fs afx.IAPSettings,
	) (out service.GrpcServiceFactory, gatewayOut service.GatewayServiceFactory, err error) {
		if svc, e := NewService(
			l,
			s.Deployment,
			s.RedisUrl,
			s.RedisPwd,
			s.IAPCredential,
			s.StatHost,
			fs.IAPUrl,
		); e != nil {
			err = e
		} else {
			err = out.Execute(svc)
			err = gatewayOut.Execute(svc)
		}
		return
	},
)
