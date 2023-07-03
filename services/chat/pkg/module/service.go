package module

import (
	"github.com/donglei1234/platform/services/chat/internal/app/service/public"
	fx2 "github.com/donglei1234/platform/services/chat/pkg/fx"
	"go.uber.org/fx"
)

var Module = fx.Module("chat", fx.Options(
	public.ServiceModule,
	fx2.ChatSettingsModule,
	fx2.StoreProviderModule,
	fx2.ChatMemoryModule,
))
