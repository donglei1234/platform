package fxsvcapp

import (
	"go.uber.org/fx"
	"net/http"
	"time"

	"github.com/googollee/go-socket.io/engineio"
	"github.com/googollee/go-socket.io/engineio/transport"
	"github.com/googollee/go-socket.io/engineio/transport/polling"
	"github.com/googollee/go-socket.io/engineio/transport/websocket"

	"github.com/donglei1234/platform/services/common/config"
)

type SocketIOSettings struct {
	fx.In

	PingTimeout  int `name:"PingTimeout"`
	PingInterval int `name:"PingInterval"`
}

func (s *SocketIOSettings) GetSocketIOEngineOption() *engineio.Options {
	pt := polling.Default

	wt := websocket.Default
	wt.CheckOrigin = func(req *http.Request) bool {
		return true
	}
	return &engineio.Options{
		Transports: []transport.Transport{
			pt,
			wt,
		},
		PingTimeout:  time.Duration(s.PingTimeout) * time.Second,
		PingInterval: time.Duration(s.PingInterval) * time.Second,
	}
}

type SocketIOSettingsLoader struct {
	fx.Out
	config.EnvironmentBlock

	PingTimeout  int `name:"PingTimeout" default:"300" envconfig:"PINGTIMEOUT"`
	PingInterval int `name:"PingInterval" default:"2" envconfig:"PINGINTERVAL"`
}

func (g *SocketIOSettingsLoader) LoadFromEnv() (err error) {
	err = config.Load(g)
	return
}

var SocketIOSettingsModule = fx.Provide(
	func() (out SocketIOSettingsLoader, err error) {
		err = out.LoadFromEnv()
		return
	},
)
