package fxsvcapp

import (
	"context"

	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/pkg/errors"
	"go.uber.org/fx"
	"go.uber.org/zap"

	"github.com/donglei1234/platform/services/common/server"
	"github.com/donglei1234/platform/services/common/service"
)

type ServiceBinder struct {
	fx.In

	TcpConnectionMux  server.ConnectionMux `name:"TcpConnectionMux"`
	UnixConnectionMux server.ConnectionMux `name:"UnixConnectionMux"`

	TcpGrpcServer     server.GrpcServer     `name:"TcpGrpcServer"`
	TcpHttpServer     server.HttpServer     `name:"TcpHttpServer"`
	TcpSocketIOServer server.SocketIOServer `name:"TcpSocketIOServer"`
	TcpGatewayServer  server.GatewayServer  `name:"TcpGatewayServer"`
	UnixGrpcServer    server.GrpcServer     `name:"UnixGrpcServer"`
	UnixHttpServer    server.HttpServer     `name:"UnixHttpServer"`
	ZinxTcpServer     server.ZinxTcpServer  `name:"ZinxServer"`

	HttpServices     []service.HttpService     `group:"HttpService"`
	GrpcServices     []service.GrpcService     `group:"GrpcService"`
	SocketIOServices []service.SocketIOService `group:"SocketIOService"`
	GatewayServices  []service.GatewayService  `group:"GatewayService"`
	ZinxServices     []service.ZinxTcpService  `group:"ZinxTcpService"`
}

// Execute registers all configured services with their servers and hooks the app lifecycle for servers that
// have active services.
func (g *ServiceBinder) Execute(l *zap.Logger, lc fx.Lifecycle) error {
	if tcp, unix, hooks, err := appendBinders(
		l,
		lc,
		g.bindGrpcServices,
		g.bindSocketIOServices,
		g.bindHttpServices,
		g.bindGatewayServices,
		g.bindTcpServices,
	); err != nil {
		return err
	} else {
		if tcp > 0 {
			connectionMuxHook(lc, g.TcpConnectionMux)
		}

		if unix > 0 {
			connectionMuxHook(lc, g.UnixConnectionMux)
		}

		for _, h := range hooks {
			h(lc)
		}
	}

	return nil
}

func (g *ServiceBinder) bindGrpcServices(
	l *zap.Logger,
	lc fx.Lifecycle,
) (tcp, unix int, hooks []LifecycleHook, err error) {
	for _, s := range g.GrpcServices {
		switch s.ServiceTransport() {
		case service.Tcp:
			LogBindService(l, s, g.TcpGrpcServer)
			s.RegisterWithGrpcServer(g.TcpGrpcServer)
			tcp++
		case service.Unix:
			LogBindService(l, s, g.UnixGrpcServer)
			s.RegisterWithGrpcServer(g.UnixGrpcServer)
			unix++
		default:
			err = errors.Errorf("unsupported grpc service transport '%s'", s.ServiceTransport())
		}
	}

	if tcp > 0 {
		grpc_prometheus.Register(g.TcpGrpcServer.GrpcServer())
		hooks = append(hooks, makeServerHook(g.TcpGrpcServer))
	}

	if unix > 0 {
		grpc_prometheus.Register(g.UnixGrpcServer.GrpcServer())
		hooks = append(hooks, makeServerHook(g.UnixGrpcServer))
	}

	return
}

func (g *ServiceBinder) bindHttpServices(
	l *zap.Logger,
	lc fx.Lifecycle,
) (tcp, unix int, hooks []LifecycleHook, err error) {
	for _, s := range g.HttpServices {
		switch s.ServiceTransport() {
		case service.Tcp:
			LogBindService(l, s, g.TcpHttpServer)
			s.RegisterWithHttpServer(g.TcpHttpServer)
			tcp++
		case service.Unix:
			LogBindService(l, s, g.UnixHttpServer)
			s.RegisterWithHttpServer(g.UnixHttpServer)
			unix++
		default:
			err = errors.Errorf("unsupported http service transport '%s'", s.ServiceTransport())
		}
	}

	if tcp > 0 {
		hooks = append(hooks, makeServerHook(g.TcpHttpServer))
	}

	if unix > 0 {
		hooks = append(hooks, makeServerHook(g.UnixHttpServer))
	}

	return
}

func (g *ServiceBinder) bindSocketIOServices(
	l *zap.Logger,
	lc fx.Lifecycle,
) (tcp, unix int, hooks []LifecycleHook, err error) {
	for _, s := range g.SocketIOServices {
		switch s.ServiceTransport() {
		case service.Tcp:
			LogBindService(l, s, g.TcpSocketIOServer)
			s.RegisterWithSocketIOServer(g.TcpSocketIOServer)
			tcp++
		default:
			err = errors.Errorf("unsupported socket.io service transport '%s'", s.ServiceTransport())
		}
	}

	if tcp > 0 {
		hooks = append(hooks, makeServerHook(g.TcpSocketIOServer))
	}

	return
}

func (g *ServiceBinder) bindTcpServices(
	l *zap.Logger,
	lc fx.Lifecycle,
) (tcp, unix int, hooks []LifecycleHook, err error) {
	for _, s := range g.ZinxServices {
		switch s.ServiceTransport() {
		case service.Tcp:
			LogBindService(l, s, g.ZinxTcpServer)
			s.RegisterWithTCPServer(g.ZinxTcpServer)
			tcp++
		default:
			err = errors.Errorf("unsupported tcp service transport '%s'", s.ServiceTransport())
		}
	}

	if tcp > 0 {
		hooks = append(hooks, makeServerHook(g.ZinxTcpServer))
	}

	return
}

func (g *ServiceBinder) bindGatewayServices(
	l *zap.Logger,
	lc fx.Lifecycle,
) (tcp, unix int, hooks []LifecycleHook, err error) {
	for _, s := range g.GatewayServices {
		switch s.ServiceTransport() {
		case service.Tcp:
			LogBindService(l, s, g.TcpGatewayServer)
			err := s.RegisterWithGatewayServer(g.TcpGatewayServer)
			if err != nil {
				return 0, 0, nil, err
			}
			tcp++
		default:
			err = errors.Errorf("unsupported gateway service transport '%s'", s.ServiceTransport())
		}
	}

	if tcp > 0 {
		hooks = append(hooks, makeServerHook(g.TcpGatewayServer))
	}

	return
}

type LifecycleHook = func(lc fx.Lifecycle)

func makeServerHook(s server.Server) LifecycleHook {
	return func(lc fx.Lifecycle) {
		lc.Append(fx.Hook{
			OnStart: func(ctx context.Context) error {
				return s.StartServing(ctx)
			},
			OnStop: func(ctx context.Context) error {
				return s.StopServing(ctx)
			},
		})
	}
}

func connectionMuxHook(lc fx.Lifecycle, m server.ConnectionMux) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			return m.StartServing(ctx)
		},
		OnStop: func(ctx context.Context) error {
			return m.StopServing(ctx)
		},
	})
}

func appendBinders(
	l *zap.Logger,
	lc fx.Lifecycle,
	fs ...func(*zap.Logger, fx.Lifecycle) (int, int, []LifecycleHook, error),
) (a int, b int, hs []LifecycleHook, err error) {
	for _, f := range fs {
		if x, y, gs, e := f(l, lc); e != nil {
			err = e
			break
		} else {
			a += x
			b += y
			hs = append(hs, gs...)
		}
	}

	return
}
