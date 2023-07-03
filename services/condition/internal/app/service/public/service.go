package public

import (
	"github.com/donglei1234/platform/services/common/mq"
	pb "github.com/donglei1234/platform/services/condition/gen/condition/api"
	"github.com/donglei1234/platform/services/condition/internal/app/db"
	fx2 "github.com/donglei1234/platform/services/condition/pkg/fx"
	"github.com/donglei1234/platform/services/condition/pkg/metadata"
	"go.uber.org/fx"
	"go.uber.org/zap"

	"github.com/donglei1234/platform/services/common/access"
	"github.com/donglei1234/platform/services/common/fx/fxsvcapp"
	"github.com/donglei1234/platform/services/common/nosql/memory"
	"github.com/donglei1234/platform/services/common/service"
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
	pb.RegisterConditionServiceServer(server.GrpcServer(), s)
	return nil
}

func (s *Service) RegisterWithGatewayServer(server service.HasGatewayServer) error {
	return nil
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
		fs fx2.ConditionSettings,
		memory fxsvcapp.GlobalConditionServerStore,
		mq fxsvcapp.GlobalMQ,
	) (out service.GrpcServiceFactory, err error) {
		if svc, e := NewService(
			l,
			s.Deployment,
			memory.ConditionServerStore,
			fs.ConditionUrl,
			mq.MessageQueue,
		); e != nil {
			err = e
		} else {
			err = out.Execute(svc)
		}
		return
	},
)
