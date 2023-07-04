package external

import (
	cfx "github.com/donglei1234/platform/services/chat/pkg/fx"
	"github.com/donglei1234/platform/services/common/fx/fxsvcapp"
	"github.com/donglei1234/platform/services/common/mq/nats"
	"go.uber.org/fx"
)

var ServiceModules = fx.Options(
	fxsvcapp.MessageQueueModule,
	nats.MQModule,
	cfx.ChatSettingsModule,
)
