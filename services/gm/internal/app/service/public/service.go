package public

import (
	"context"
	"github.com/donglei1234/platform/services/common/mq"
	"github.com/donglei1234/platform/services/gm/internal/app/db"
	"github.com/donglei1234/platform/services/gm/internal/app/service/external"
	pb "github.com/donglei1234/platform/services/proto/gen/gm/api"
	"go.uber.org/fx"
	"go.uber.org/zap"

	"github.com/donglei1234/platform/services/common/access"
	"github.com/donglei1234/platform/services/common/fx/fxsvcapp"
	"github.com/donglei1234/platform/services/common/nosql/memory"
	"github.com/donglei1234/platform/services/common/service"
	bfx "github.com/donglei1234/platform/services/gm/pkg/fx"
	"github.com/donglei1234/platform/services/gm/pkg/metadata"
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
	pb.RegisterGmServiceServer(
		server.GrpcServer(),
		s,
	)
	return nil
}

func (s *Service) RegisterWithGatewayServer(server service.HasGatewayServer) error {
	return pb.RegisterGmServiceHandlerFromEndpoint(
		context.Background(), server.GatewayRuntimeMux(), s.url, server.GatewayOption(),
	)
}

func NewService(
	l *zap.Logger,
	deployment string,
	memory memory.MemoryStore,
	url string,
	mq mq.MessageQueue,
) (result *Service, err error) {
	result = &Service{
		appId:      metadata.AppId,
		logger:     l,
		deployment: deployment,
		db:         db.OpenDatabase(l, memory),
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
		fs bfx.GMSettings,
		memory external.GlobalGMServerStoreParams,
		mq fxsvcapp.GlobalMQ,
	) (out service.GrpcServiceFactory, gatewayOut service.GatewayServiceFactory, err error) {
		if svc, e := NewService(
			l,
			s.Deployment,
			memory.GMServerStore,
			fs.GmUrl,
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
