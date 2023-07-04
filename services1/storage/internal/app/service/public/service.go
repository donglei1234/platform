package public

import (
	"context"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"go.uber.org/fx"
	"go.uber.org/zap"

	"github.com/donglei1234/platform/services/common/access"
	"github.com/donglei1234/platform/services/common/fx/fxsvcapp"
	"github.com/donglei1234/platform/services/common/service"
	pb "github.com/donglei1234/platform/services/storage/generated/grpc/go/storage/api"
	"github.com/donglei1234/platform/services/storage/internal/app/db"
	nfx "github.com/donglei1234/platform/services/storage/pkg/fx"
	"github.com/donglei1234/platform/services/storage/pkg/metadata"
)

type Service struct {
	service.TcpTransport
	appId      string
	logger     *zap.Logger
	deployment string
	db         *db.Database
	sess       *session.Session
	bucketName string
	url        string
}

func (s *Service) RegisterWithGrpcServer(server service.HasGrpcServer) error {
	pb.RegisterStorageServer(server.GrpcServer(), s)
	return nil
}

func (s *Service) RegisterWithGatewayServer(server service.HasGatewayServer) error {
	if err := pb.RegisterStorageHandlerFromEndpoint(
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
	awsAccessKey string,
	awsAccessSecret string,
	awsRegion string,
	bucketName string,
	url string,
) (result *Service, err error) {
	result = &Service{
		appId:      metadata.AppId,
		logger:     l,
		deployment: deployment,
		db:         db.OpenDatabase(l, redisUrl, redisPwd),
		sess:       NewS3Session(awsAccessKey, awsAccessSecret, awsRegion),
		bucketName: bucketName,
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
		fs nfx.StorageSettings,
	) (out service.GrpcServiceFactory, gatewayOut service.GatewayServiceFactory, err error) {
		if svc, e := NewService(
			l,
			s.Deployment,
			s.RedisUrl,
			s.RedisPwd,
			s.AwsAccessKey,
			s.AwsAccessSecret,
			s.AwsRegion,
			s.AwsBucketName,
			fs.StorageUrl,
		); e != nil {
			err = e
		} else {
			err = out.Execute(svc)
			err = gatewayOut.Execute(svc)
		}
		return
	},
)

func NewS3Session(awsAccessKey, awsAccessSecret, awsRegion string) *session.Session {
	return session.Must(session.NewSession(&aws.Config{
		Credentials: credentials.NewStaticCredentials(awsAccessKey, awsAccessSecret, ""),
		Region:      aws.String(awsRegion),
	}))
}
