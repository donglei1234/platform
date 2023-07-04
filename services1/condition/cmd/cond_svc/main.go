package main

import (
	auth "github.com/donglei1234/platform/services/auth/pkg/module"
	"github.com/donglei1234/platform/services/common/fx/fxsvcapp"
	"github.com/donglei1234/platform/services/common/mq/mock"
	"github.com/donglei1234/platform/services/condition/internal/app/service/external"
	"github.com/donglei1234/platform/services/condition/internal/app/service/public"
	"github.com/donglei1234/platform/services/condition/pkg/fx"
	"os"
	"path"
)

func main() {
	home := os.Getenv("HOME")
	setEnvVariable("DOCUMENT_STORE_URL", "badger://badger/"+path.Join(home, "room"))
	setEnvVariable("MEMORY_STORE_URL", "redis://127.0.0.1:6379/0")

	fxsvcapp.StandardMain(
		fx.ConditionSettingsModule,
		public.ServiceModule,
		external.ServicesModules,
		auth.PublicModule,
		fxsvcapp.AuthStoreModule,
		fxsvcapp.MessageQueueModule,
		mock.MQModule,
	)
}

func setEnvVariable(key string, value string) {
	if os.Getenv(key) == "" {
		_ = os.Setenv(key, value)
	}
}
