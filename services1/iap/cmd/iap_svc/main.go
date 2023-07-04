package main

import (
	"math/rand"
	"os"
	"path"
	"time"

	afx "github.com/donglei1234/platform/services/auth/pkg/fx"
	auth "github.com/donglei1234/platform/services/auth/pkg/module"
	"github.com/donglei1234/platform/services/common/fx/fxsvcapp"
	"github.com/donglei1234/platform/services/iap/internal/app/service/public"
	ffx "github.com/donglei1234/platform/services/iap/pkg/fx"
)

func main() {
	rand.Seed(time.Now().Unix())
	home := os.Getenv("HOME")
	setEnvVariable("DOCUMENT_STORE_URL", "badger://badger/"+path.Join(home, "iap"))

	fxsvcapp.StandardMain(
		ffx.IAPSettingsModule,
		public.ServiceModule,
		auth.PublicModule,
		afx.AuthSettingsModule,
		fxsvcapp.AuthStoreModule,
		fxsvcapp.AuthClientModule,
	)
}

func setEnvVariable(key string, value string) {
	if os.Getenv(key) == "" {
		_ = os.Setenv(key, value)
	}
}
