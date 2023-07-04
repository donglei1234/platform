package module

import (
	"github.com/donglei1234/platform/services/auth/internal/app/service/public"
	"github.com/donglei1234/platform/services/auth/pkg/metadata"
	"go.uber.org/fx"
)

var (
	Modules = fx.Module(metadata.AppId,
		public.ServiceModule,
	)
)
