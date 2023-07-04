package main

import (
	"github.com/donglei1234/platform/services/chat/internal/app/service/external"
	"github.com/donglei1234/platform/services/chat/internal/app/service/public"
	"github.com/donglei1234/platform/services/chat/pkg/fx"
	"github.com/donglei1234/platform/services/common/fx/fxsvcapp"
)

func main() {
	fxsvcapp.StandardMain(
		external.ServiceModules,
		public.ServiceModule,
		fx.StoreProviderModule,
	)
}
