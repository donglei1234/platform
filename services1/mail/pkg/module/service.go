package module

import (
	"github.com/donglei1234/platform/services/mail/internal/app/service/public"
	fx3 "github.com/donglei1234/platform/services/mail/pkg/fx"
	"go.uber.org/fx"
)

var (
	PublicModule = fx.Options(
		public.ServiceModule,
		fx3.MailClientModule,
		fx3.MailSettingsModule,
	)
)
