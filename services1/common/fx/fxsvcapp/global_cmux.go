package fxsvcapp

import (
	"github.com/donglei1234/platform/services/common/server"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

// GlobalConnectionMux loads all of the global connection mux instances from the fx dependency graph.
type GlobalConnectionMux struct {
	fx.In
	TcpConnectionMux  server.ConnectionMux `name:"TcpConnectionMux"`
	UnixConnectionMux server.ConnectionMux `name:"UnixConnectionMux"`
}

// GlobalConnectionMuxFactory injects connection mux instances into the fx dependency graph based on values in GlobalSettings.
type GlobalConnectionMuxFactory struct {
	fx.Out
	TcpConnectionMux  server.ConnectionMux `name:"TcpConnectionMux"`
	UnixConnectionMux server.ConnectionMux `name:"UnixConnectionMux"`
}

func (f *GlobalConnectionMuxFactory) Execute(l *zap.Logger, g GlobalSettings, s SecuritySettings) (err error) {
	newTcpConnectionMux := func(out *server.ConnectionMux) {
		if err != nil {
			return
		}

		port := server.Port(g.Port)

		if g.AppTestMode {
			*out, err = server.NewTestTcpConnectionMux()
		} else if s.TlsCert != "" && s.TlsKey != "" {
			*out, err = server.NewTlsTcpConnectionMux(l, port, s.TlsCert, s.TlsKey)
		} else {
			*out, err = server.NewTcpConnectionMux(l, port)
		}
	}

	newUnixConnectionMux := func(out *server.ConnectionMux) {
		if err != nil {
			// do nothing
		} else if g.AppTestMode {
			l.Warn("unix connection mux not supported in test mode")
		} else {
			socket := server.Socket(g.Socket)
			*out, err = server.NewUnixConnectionMux(l, socket)
		}
	}

	newTcpConnectionMux(&f.TcpConnectionMux)
	newUnixConnectionMux(&f.UnixConnectionMux)

	return
}

var ConnectionMuxModule = fx.Provide(
	func(l *zap.Logger, g GlobalSettings, s SecuritySettings) (out GlobalConnectionMuxFactory, err error) {
		err = out.Execute(l, g, s)
		return
	},
)
