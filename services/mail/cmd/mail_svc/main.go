package main

import (
	"github.com/donglei1234/platform/services/common/fx/fxsvcapp"
	"github.com/donglei1234/platform/services/common/mq/mock"
	"github.com/donglei1234/platform/services/mail/internal/app/service/external"
	"github.com/donglei1234/platform/services/mail/internal/app/service/public"
	ffx "github.com/donglei1234/platform/services/mail/pkg/fx"
	"os"
	"path"
)

func main() {
	home := os.Getenv("HOME")
	setEnvVariable("DOCUMENT_STORE_URL", "badger://badger/"+path.Join(home, "mail"))
	setEnvVariable("MEMORY_STORE_URL", "redis://127.0.0.1:6379")

	fxsvcapp.StandardMain(
		ffx.MailSettingsModule,
		public.ServiceModule,
		external.ServicesModules,
		fxsvcapp.MessageQueueModule,
		mock.MQModule,
	)
}

func setEnvVariable(key string, value string) {
	if os.Getenv(key) == "" {
		_ = os.Setenv(key, value)
	}
}
