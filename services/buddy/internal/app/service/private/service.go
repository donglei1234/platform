package private

import (
	"context"
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

type Service struct {
	service.Logger
	service.TcpTransport
	appId      string
	db         db.Database
	deployment string
}

func NewService(
	l *zap.Logger,
	rosStore document.DocumentStore,
	deployment string,
) (result *Service, err error) {
	result = &Service{
		appId:      metadata.AppId,
		db:         db.OpenDatabase(l, rosStore),
		deployment: deployment,
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
	pb.RegisterPrivateServiceServer(server.GrpcServer(), s)
	return nil
}

func (s *Service) sessionFromContext(ctx context.Context) (*authdb.Session, error) {
	if session, ok := authdb.SessionFromContext(ctx); !ok {
		return nil, ErrNoMetaData
	} else {
		return session, nil
	}
}

var ServiceModule = fx.Provide(
	func(
		l *zap.Logger,
		ac fxsvcapp.GlobalAuthClient,
		rdb fxsvcapp.GlobalRosDataStore,
		s fxsvcapp.GlobalSettings,
		fs bfx.BuddySettings,
	) (out service.GrpcServiceFactory, err error) {
		if svc, e := NewService(
			l,
			rdb.RosDataStore,
			s.Deployment,
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
