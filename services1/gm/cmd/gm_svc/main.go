package main

import (
	auth "github.com/donglei1234/platform/services/auth/pkg/module"
	"github.com/donglei1234/platform/services/common/fx/fxsvcapp"
	"github.com/donglei1234/platform/services/common/mq/mock"
	"github.com/donglei1234/platform/services/gm/internal/app/service/external"
	"github.com/donglei1234/platform/services/gm/internal/app/service/public"
	ffx "github.com/donglei1234/platform/services/gm/pkg/fx"
	"os"
	"path"
)

func main() {
	home := os.Getenv("HOME")
	setEnvVariable("DOCUMENT_STORE_URL", "badger://badger/"+path.Join(home, "mail"))
	//setEnvVariable("MEMORY_STORE_URL", "redis://127.0.0.1:6379")
	setEnvVariable("MEMORY_STORE_URL", "redis://test:wg1q2w3e@192.168.37.135:6379")

	fxsvcapp.StandardMain(
		ffx.GmSettingsModule,
		public.ServiceModule,
		auth.PublicModule,
		external.GMStoreModule,
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
