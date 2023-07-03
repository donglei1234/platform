package public

import (
	"context"
	"github.com/donglei1234/platform/services/common/mq"
	"github.com/donglei1234/platform/services/common/nosql/document"
	pb "github.com/donglei1234/platform/services/proto/gen/buddy/api"

	"github.com/donglei1234/platform/services/auth/pkg/authdb"

	"github.com/donglei1234/platform/services/buddy/internal/db"
	bfx "github.com/donglei1234/platform/services/buddy/pkg/fx"
	"github.com/donglei1234/platform/services/buddy/pkg/metadata"
	"github.com/donglei1234/platform/services/common/access"
	"github.com/donglei1234/platform/services/common/fx/fxsvcapp"
	"github.com/donglei1234/platform/services/common/service"

	"go.uber.org/fx"
	"go.uber.org/zap"
)

var (
	NothingResponse = &pb.Nothing{}
)

type Service struct {
	service.TcpTransport
	appId      string
	logger     *zap.Logger
	db         db.Database
	mq         mq.MessageQueue
	deployment string
	maxInviter int32
	maxBuddies int32
	maxBlocked int32
}

func NewService(
	l *zap.Logger,
	rosStore document.DocumentStore,
	deployment string,
	mq mq.MessageQueue,
	maxInviter int32,
	maxBuddies int32,
	maxBlocked int32,
) (result *Service, err error) {
	result = &Service{
		appId:      metadata.AppId,
		logger:     l,
		db:         db.OpenDatabase(l, rosStore),
		deployment: deployment,
		mq:         mq,
		maxBuddies: maxBuddies,
		maxBlocked: maxBlocked,
		maxInviter: maxInviter,
	}
	return
}

func (s *Service) AccessLevel() access.AccessLevel {
	return access.AccessUndefined
}

func (s *Service) ValidateSettings() error {
	return s.db.ValidateSettings()
}

func (s *Service) RegisterWithGrpcServer(server service.HasGrpcServer) error {
	pb.RegisterPublicServiceServer(server.GrpcServer(), s)
	return nil
}

func (s *Service) sessionFromContext(ctx context.Context) (*authdb.Session, error) {
	if session, ok := authdb.SessionFromContext(ctx); !ok {
		return nil, ErrNoMetaData
	} else {
		return session, nil
	}
}

func (s *Service) loadBuddyQueueFromContext(ctx context.Context) (bq *db.BuddyQueue, err error) {
	if session, e := s.sessionFromContext(ctx); e != nil {
		err = e
	} else {
		bq, err = s.db.LoadOrCreateBuddyQueue(session.AppId, session.UserId)
	}
	return
}

var ServiceModule = fx.Provide(
	func(
		l *zap.Logger,
		rdb bfx.StoreParams,
		mq fxsvcapp.GlobalMQ,
		s fxsvcapp.GlobalSettings,
		fs bfx.BuddySettings,
	) (out service.GrpcServiceFactory, err error) {
		if svc, e := NewService(
			l,
			rdb.ServerStore,
			s.Deployment,
			mq.MessageQueue,
			fs.InviterMaxCount,
			fs.BuddyMaxCount,
			fs.BlockedMaxCount,
		); e != nil {
			err = e
		} else if e = svc.ValidateSettings(); e != nil {
			err = e
		} else {
			out.GrpcService = svc
		}
		return
	},
)
