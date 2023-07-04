package fxsvcapp

import (
	"runtime"

	"go.uber.org/fx"
	"go.uber.org/zap"
	"google.golang.org/grpc"

	"github.com/donglei1234/platform/services/common/server"
)

// GlobalServers loads all of the global server instances from the fx dependency graph.
type GlobalServers struct {
	fx.In
	TcpGrpcServer     server.GrpcServer     `name:"TcpGrpcServer"`
	TcpHttpServer     server.HttpServer     `name:"TcpHttpServer"`
	UnixGrpcServer    server.GrpcServer     `name:"UnixGrpcServer"`
	UnixHttpServer    server.HttpServer     `name:"UnixHttpServer"`
	TcpSocketIOServer server.SocketIOServer `name:"TcpSocketIOServer"`
	TcpGatewayServer  server.GatewayServer  `name:"TcpGatewayServer"`
	ZinxServer        server.ZinxTcpServer  `name:"ZinxServer"`
}

// GlobalServersFactory injects server instances into the fx dependency graph based on values in GlobalSettings.
type GlobalServersFactory struct {
	fx.Out
	TcpGrpcServer     server.GrpcServer     `name:"TcpGrpcServer"`
	TcpHttpServer     server.HttpServer     `name:"TcpHttpServer"`
	UnixGrpcServer    server.GrpcServer     `name:"UnixGrpcServer"`
	UnixHttpServer    server.HttpServer     `name:"UnixHttpServer"`
	TcpSocketIOServer server.SocketIOServer `name:"TcpSocketIOServer"`
	TcpGatewayServer  server.GatewayServer  `name:"TcpGatewayServer"`
	ZinxServer        server.ZinxTcpServer  `name:"ZinxServer"`
}

func (f *GlobalServersFactory) Execute(
	l *zap.Logger,
	tr GlobalTracer,
	s GlobalSettings,
	a GlobalAuthClient,
	m GlobalConnectionMux,
) (err error) {
	{
		newTcpGrpcServer := func(mux server.ConnectionMux, out *server.GrpcServer) {
			if err != nil {
				return
			}
			switch t := mux.(type) {
			case server.TcpServer:
				if s.AppTestMode {
					*out, err = server.NewTestTcpGrpcServer(
						l,
						tr.Tracer,
						t.Port(),
						s.Version,
						a.AuthClient,
						grpc.KeepaliveEnforcementPolicy(s.GetKeepaliveEnforcementPolicy()),
						grpc.KeepaliveParams(s.GetKeepaliveServerParameters()),
					)
				} else {
					*out, err = server.NewTcpGrpcServer(
						l,
						tr.Tracer,
						mux,
						t.Port(),
						s.Version,
						a.AuthClient,
						grpc.KeepaliveEnforcementPolicy(s.GetKeepaliveEnforcementPolicy()),
						grpc.KeepaliveParams(s.GetKeepaliveServerParameters()),
					)
				}
			default:
				err = ErrTcpGrpcMuxMismatch
			}

		}

		newUnixGrpcServer := func(mux server.ConnectionMux, out *server.GrpcServer) {
			if err != nil {
				// do nothing
			} else if s.AppTestMode {
				l.Warn("unix grpc server not supported in test mode")
			} else if runtime.GOOS == "windows" {
				l.Warn("unix grpc server not supported on windows")
			} else {
				switch t := mux.(type) {
				case server.UnixServer:
					*out, err = server.NewUnixGrpcServer(
						l,
						tr.Tracer,
						mux,
						t.Socket(),
						s.Version,
						a.AuthClient,
						grpc.KeepaliveEnforcementPolicy(s.GetKeepaliveEnforcementPolicy()),
						grpc.KeepaliveParams(s.GetKeepaliveServerParameters()),
					)
				default:
					err = ErrUnixGrpcMuxMismatch
				}
			}
		}

		newTcpGrpcServer(m.TcpConnectionMux, &f.TcpGrpcServer)
		newUnixGrpcServer(m.UnixConnectionMux, &f.UnixGrpcServer)
	}

	{
		newTcpHttpServer := func(mux server.ConnectionMux, out *server.HttpServer) {
			if err != nil {
				return
			}

			switch t := mux.(type) {
			case server.TcpServer:
				if s.AppTestMode {
					*out, err = server.NewTestTcpHttpServer(l, tr.Tracer, t.Port())
				} else {
					*out, err = server.NewTcpHttpServer(l, tr.Tracer, mux, t.Port())
				}
			default:
				err = ErrTcpHttpMuxMismatch
			}
		}

		newUnixHttpServer := func(mux server.ConnectionMux, out *server.HttpServer) {
			if err != nil {
				// do nothing
			} else if s.AppTestMode {
				l.Warn("unix grpc server not supported in test mode")
			} else if runtime.GOOS == "windows" {
				l.Warn("unix http server not supported on windows")
			} else {
				switch t := mux.(type) {
				case server.UnixServer:
					*out, err = server.NewUnixHttpServer(l, tr.Tracer, mux, t.Socket())
				default:
					err = ErrUnixHttpMuxMismatch
				}
			}
		}

		newTcpHttpServer(m.TcpConnectionMux, &f.TcpHttpServer)
		newUnixHttpServer(m.UnixConnectionMux, &f.UnixHttpServer)
	}
	{
		newTcpSocketIOServer := func(mux server.ConnectionMux, out *server.SocketIOServer) {
			if err != nil {
				return
			}

			switch t := mux.(type) {
			case server.TcpServer:
				*out, err = server.NewTcpSocketIOServer(
					l,
					mux,
					t.Port(),
					s.GetSocketIOEngineOption(),
					//nil,
				)
			default:
				err = ErrTcpHttpMuxMismatch
			}
		}

		newTcpSocketIOServer(m.TcpConnectionMux, &f.TcpSocketIOServer)
	}
	{
		newTcpServer := func(mux server.ConnectionMux, out *server.ZinxTcpServer) {
			if err != nil {
				return
			}
			switch t := mux.(type) {
			case server.TcpServer:
				*out, err = server.NewZinxTcpServer(
					l,
					mux,
					t.Port(),
				)
			default:
				err = ErrTcpHttpMuxMismatch
			}
		}

		newTcpServer(m.TcpConnectionMux, &f.ZinxServer)
	}
	{
		newTcpGatewayServer := func(mux server.ConnectionMux, out *server.GatewayServer) {
			if err != nil {
				return
			}
			switch t := mux.(type) {
			case server.TcpServer:
				*out, err = server.NewTcpGatewayServer(
					l,
					mux,
					t.Port(),
				)
			default:
				err = ErrTcpHttpMuxMismatch
			}
		}
		newTcpGatewayServer(m.TcpConnectionMux, &f.TcpGatewayServer)
	}
	return
}

var ServersModule = fx.Provide(
	func(
		l *zap.Logger,
		t GlobalTracer,
		g GlobalSettings,
		a GlobalAuthClient,
		s SecuritySettings,
		m GlobalConnectionMux,
	) (out GlobalServersFactory, err error) {
		err = out.Execute(l, t, g, a, m)
		return
	},
)
