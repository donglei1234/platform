package public

import (
	"context"
	"encoding/json"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sns"
	"go.uber.org/fx"
	"go.uber.org/zap"

	"github.com/donglei1234/platform/services/common/access"
	"github.com/donglei1234/platform/services/common/fx/fxsvcapp"
	"github.com/donglei1234/platform/services/common/service"
	pb "github.com/donglei1234/platform/services/notification/generated/grpc/go/notification/api"
	"github.com/donglei1234/platform/services/notification/internal/app/db"
	nfx "github.com/donglei1234/platform/services/notification/pkg/fx"
	"github.com/donglei1234/platform/services/notification/pkg/metadata"
)

type Service struct {
	service.TcpTransport
	appId                string
	logger               *zap.Logger
	deployment           string
	awsSNSApplicationARN map[string]map[string]string
	topicArn             map[string]string
	db                   *db.Database
	sns                  *sns.SNS
	url                  string
}

func (s *Service) RegisterWithGrpcServer(server service.HasGrpcServer) error {
	pb.RegisterNotificationServer(server.GrpcServer(), s)
	return nil
}

func (s *Service) RegisterWithGatewayServer(server service.HasGatewayServer) error {
	if err := pb.RegisterNotificationHandlerFromEndpoint(
		context.Background(), server.GatewayRuntimeMux(), s.url, server.GatewayOption()); err != nil {
		return err
	}
	return nil
}

func NewService(
	l *zap.Logger,
	deployment string,
	awsAccessKey string,
	awsAccessSecret string,
	awsRegion string,
	awsSNSApplicationARN string,
	topicArn string,
	redisUrl string,
	redisPwd string,
	listeningUrl string,
) (result *Service, err error) {
	tempTopicArn, err := ParseTopicArn(topicArn)
	if err != nil {
		return nil, err
	}
	tempApplicationARN, err := ParseApplicationArn(awsSNSApplicationARN)
	if err != nil {
		return nil, err
	}
	result = &Service{
		appId:                metadata.AppId,
		logger:               l,
		deployment:           deployment,
		awsSNSApplicationARN: tempApplicationARN,
		topicArn:             tempTopicArn,
		db:                   db.OpenDatabase(l, redisUrl, redisPwd),
		sns:                  NewSNS(awsAccessKey, awsAccessSecret, awsRegion),
		url:                  listeningUrl,
	}
	return
}

func (s *Service) AccessLevel() access.AccessLevel {
	return access.AccessUndefined
}

func NewSNS(awsAccessKey, awsAccessSecret, awsRegion string) *sns.SNS {
	return sns.New(session.Must(session.NewSession(&aws.Config{
		Credentials: credentials.NewStaticCredentials(awsAccessKey, awsAccessSecret, ""),
		Region:      aws.String(awsRegion),
	})))
}

var ServiceModule = fx.Provide(
	func(
		l *zap.Logger,
		s fxsvcapp.GlobalSettings,
		fs nfx.NotificationSettings,
	) (out service.GrpcServiceFactory, gatewayOut service.GatewayServiceFactory, err error) {
		if svc, e := NewService(
			l,
			s.Deployment,
			s.AwsAccessKey,
			s.AwsAccessSecret,
			s.AwsRegion,
			s.AwsSNSApplicationARN,
			s.TopicArn,
			s.RedisUrl,
			s.RedisPwd,
			fs.NotificationUrl,
		); e != nil {
			err = e
		} else {
			err = out.Execute(svc)
			err = gatewayOut.Execute(svc)
		}
		return
	},
)

func ParseApplicationArn(awsApplicationARN string) (res map[string]map[string]string, err error) {
	err = json.Unmarshal([]byte(awsApplicationARN), &res)
	return
}

func ParseTopicArn(topicArn string) (res map[string]string, err error) {
	err = json.Unmarshal([]byte(topicArn), &res)
	return
}
