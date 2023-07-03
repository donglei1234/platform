package module

import (
	"github.com/donglei1234/platform/services/buddy/internal/app/service/private"
	"github.com/donglei1234/platform/services/buddy/internal/app/service/public"
	fx2 "github.com/donglei1234/platform/services/buddy/pkg/fx"
	"github.com/donglei1234/platform/services/buddy/pkg/metadata"
	"go.uber.org/fx"
)

var (
	PublicModule = fx.Module(metadata.AppId,
		public.ServiceModule,
		fx2.BuddyClientModule,
		fx2.BuddySettingsModule,
		fx2.StoreProviderModule,
		fx2.BuddyMemoryModule,
	)
	PrivateModule = private.ServiceModule
)
