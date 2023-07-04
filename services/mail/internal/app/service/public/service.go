package public

import (
	"context"
	"github.com/donglei1234/platform/services/common/mq"
	"github.com/donglei1234/platform/services/mail/internal/app/db"
	pb2 "github.com/donglei1234/platform/services/proto/gen/mail/api"
	"github.com/go-redis/redis"
	"go.uber.org/fx"
	"go.uber.org/zap"

	"github.com/donglei1234/platform/services/common/access"
	"github.com/donglei1234/platform/services/common/fx/fxsvcapp"
	"github.com/donglei1234/platform/services/common/service"
	bfx "github.com/donglei1234/platform/services/mail/pkg/fx"
	"github.com/donglei1234/platform/services/mail/pkg/metadata"
)

type Service struct {
	service.TcpTransport
	appId      string
	logger     *zap.Logger
	deployment string
	db         *db.Database
	url        string
	mq         mq.MessageQueue
}

func (s *Service) RegisterWithGrpcServer(server service.HasGrpcServer) error {
	pb2.RegisterMailServiceServer(
		server.GrpcServer(),
		s,
	)
	return nil
}

func (s *Service) RegisterWithGatewayServer(server service.HasGatewayServer) error {
	return pb2.RegisterMailServiceHandlerFromEndpoint(
		context.Background(), server.GatewayRuntimeMux(), s.url, server.GatewayOption(),
	)
}

func NewService(
	l *zap.Logger,
	deployment string,
	redis *redis.Client,
	url string,
	mq mq.MessageQueue,
) (result *Service, err error) {
	result = &Service{
		appId:      metadata.AppId,
		logger:     l,
		deployment: deployment,
		db:         db.OpenDatabase(l, redis),
		url:        url,
		mq:         mq,
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
		fs bfx.MailSettings,
		memory fxsvcapp.GlobalRedisParams,
		mq fxsvcapp.GlobalMQ,
	) (out service.GrpcServiceFactory, gatewayOut service.GatewayServiceFactory, err error) {
		if svc, e := NewService(
			l,
			s.Deployment,
			memory.Redis,
			fs.MailUrl,
			mq.MessageQueue,
		); e != nil {
			err = e
		} else {
			err = out.Execute(svc)
			err = gatewayOut.Execute(svc)
		}
		return
	},
)
