package public

import (
	"context"
	"fmt"
	"github.com/donglei1234/platform/services/chat/internal/db"
	"github.com/go-redis/redis"
	"go.uber.org/zap"
	"sync"
	"time"

	"go.uber.org/fx"

	"github.com/donglei1234/platform/services/auth/pkg/authdb"
	cfx "github.com/donglei1234/platform/services/chat/pkg/fx"
	"github.com/donglei1234/platform/services/chat/pkg/metadata"
	"github.com/donglei1234/platform/services/common/access"
	"github.com/donglei1234/platform/services/common/fx/fxsvcapp"
	"github.com/donglei1234/platform/services/common/mq"
	"github.com/donglei1234/platform/services/common/service"
	pb "github.com/donglei1234/platform/services/proto/gen/chat/api"
)

const (
	// session metadata key where desired chat profile id may be specified
	chatProfileMetaKey = "chatId"
)

type Service struct {
	service.TcpTransport
	service.Logger
	appId      string
	deployment string
	mq         mq.MessageQueue
	msgCache   *db.Database

	clients *sync.Map

	chatInterval time.Duration
}

func NewService(
	l *zap.Logger,
	//as document.DocumentStore,
	rClient *redis.Client,
	deployment string,
	mq mq.MessageQueue,
	chatInterval int,
) (result *Service, err error) {
	result = &Service{
		msgCache:     db.OpenDatabase(l, rClient),
		appId:        metadata.AppId,
		deployment:   deployment,
		mq:           mq,
		clients:      &sync.Map{},
		chatInterval: time.Duration(chatInterval) * time.Second,
	}
	return
}

func (s *Service) AccessLevel() access.AccessLevel {
	return access.AccessUndefined
}

func (s *Service) RegisterWithGrpcServer(server service.HasGrpcServer) error {
	pb.RegisterChatServiceServer(server.GrpcServer(), s)
	return nil
}

func getProfileId(session *authdb.Session) (profileId string) {
	// If the session metadata has a chat profile ID set (by a consuming api), use that - otherwise fall back to the
	// session's user/account ID field instead.
	if metaId, ok := session.GetMetadata(chatProfileMetaKey); ok {
		profileId = metaId
	} else {
		profileId = session.UserId
	}
	return profileId
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
		cm cfx.ChatMemoryParams,
		mq fxsvcapp.GlobalMQ,
		s fxsvcapp.GlobalSettings,
		cfx cfx.ChatSettings,

	) (out service.GrpcServiceFactory, err error) {
		if svc, e := NewService(
			l,
			cm.MemoryStore,
			s.Deployment,
			mq.MessageQueue,
			cfx.ChatInterval,
		); e != nil {
			err = e
		} else {
			out.GrpcService = svc
		}
		return
	},
)

func (s *Service) chatRoomTopic(channel pb.ChatChannel, id string) string {
	if channel == pb.ChatChannel_System || channel == pb.ChatChannel_World {
		return fmt.Sprintf("%s.%s.%s.empty", s.appId, s.deployment, channel.String())
	}
	return fmt.Sprintf("%s.%s.%s.%s", s.appId, s.deployment, channel.String(), id)
}
