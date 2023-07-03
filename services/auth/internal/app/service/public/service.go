package public

import (
	pb "github.com/donglei1234/platform/services/proto/gen/auth/api"
	"go.uber.org/fx"
	"go.uber.org/zap"

	db "github.com/donglei1234/platform/services/auth/internal/db/document"
	bfx "github.com/donglei1234/platform/services/auth/pkg/fx"
	"github.com/donglei1234/platform/services/auth/pkg/metadata"
	"github.com/donglei1234/platform/services/common/access"
	"github.com/donglei1234/platform/services/common/fx/fxsvcapp"
	"github.com/donglei1234/platform/services/common/nosql/document"
	"github.com/donglei1234/platform/services/common/service"
)

type Service struct {
	service.TcpTransport
	appId          string
	logger         *zap.Logger
	deployment     string
	doc            *db.Document
	facebookAppId  string
	facebookSecret string
	googleSecret   string
	jwtSecret      string
	url            string
}

func (s *Service) RegisterWithGrpcServer(server service.HasGrpcServer) error {
	pb.RegisterAuthPublicServer(server.GrpcServer(), s)
	pb.RegisterAuthPrivateServer(server.GrpcServer(), s)
	return nil
}

func (s *Service) RegisterWithGatewayServer(server service.HasGatewayServer) error {
	//if err := pb.RegisterAuthPublicHandlerFromEndpoint(
	//	context.Background(), server.GatewayRuntimeMux(), s.url, server.GatewayOption()); err != nil {
	//	return err
	//}
	//if err := pb.RegisterAuthPrivateHandlerFromEndpoint(
	//	context.Background(), server.GatewayRuntimeMux(), s.url, server.GatewayOption()); err != nil {
	//	return err
	//}
	return nil
}

func NewService(
	l *zap.Logger,
	deployment string,
	as document.DocumentStore,
	facebookAppId string,
	facebookSecret string,
	googleSecret string,
	jwtSecret string,
	url string,
) (result *Service, err error) {
	result = &Service{
		appId:          metadata.AppId,
		logger:         l,
		deployment:     deployment,
		doc:            db.OpenDatabase(l, as),
		facebookAppId:  facebookAppId,
		facebookSecret: facebookSecret,
		googleSecret:   googleSecret,
		jwtSecret:      jwtSecret,
		url:            url,
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
		fs bfx.AuthSettings,
		as fxsvcapp.GlobalAuthStore,
	) (out service.GrpcServiceFactory, gatewayOut service.GatewayServiceFactory, err error) {
		if svc, e := NewService(
			l,
			s.Deployment,
			as.AuthStore,
			s.FacebookAppId,
			s.FacebookAppSecret,
			s.GoogleApiKeys,
			s.JwtVerificationKey,
			fs.AuthUrl,
		); e != nil {
			err = e
		} else {
			err = out.Execute(svc)
			err = gatewayOut.Execute(svc)
		}
		return
	},
)
