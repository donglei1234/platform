package module

import (
	"github.com/donglei1234/platform/services/common/fx/fxsvcapp"
	"github.com/donglei1234/platform/services/condition/internal/app/service/public"
	fx2 "github.com/donglei1234/platform/services/condition/pkg/fx"
	"github.com/donglei1234/platform/services/condition/pkg/metadata"
	"go.uber.org/fx"
)

var (
	PublicModule = fx.Module(
		metadata.AppId,
		public.ServiceModule,
		fx2.ConditionClientModule,
		fx2.ConditionSettingsModule,
		fxsvcapp.ConditionStoreModule,
	)
)
