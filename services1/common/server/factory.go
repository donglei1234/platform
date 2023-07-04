package server

import (
	"github.com/googollee/go-socket.io/engineio"
	"go.uber.org/zap"
	"google.golang.org/grpc"

	auth "github.com/donglei1234/platform/services/auth/pkg/client"
	"github.com/donglei1234/platform/services/common/server/internal"
	"github.com/donglei1234/platform/services/common/tracing"
)

func NewTcpGrpcServer(
	logger *zap.Logger,
	tracer tracing.Tracer,
	listener HasGrpcListener,
	port Port,
	version string,
	authClient auth.AuthClient,
	opts ...grpc.ServerOption,
) (result GrpcServer, err error) {
	return internal.NewTcpGrpcServer(
		logger,
		listener,
		port,
		addInterceptorOptions(logger, tracer, version, authClient, opts...)...,
	)
}

func NewTestTcpGrpcServer(
	logger *zap.Logger,
	tracer tracing.Tracer,
	port Port,
	version string,
	authClient auth.AuthClient,
	opts ...grpc.ServerOption,
) (result GrpcServer, err error) {
	result = internal.NewTestGrpcServer(
		logger,
		port,
		addInterceptorOptions(logger, tracer, version, authClient, opts...)...,
	)
	return
}

func NewUnixGrpcServer(
	logger *zap.Logger,
	tracer tracing.Tracer,
	listener HasGrpcListener,
	socket Socket,
	//authClient auth.PrivateClient,
	version string,
	authClient auth.AuthClient,
	opts ...grpc.ServerOption,
) (result GrpcServer, err error) {
	return internal.NewUnixGrpcServer(
		logger,
		listener,
		socket,
		addInterceptorOptions(logger, tracer, version, authClient, opts...)...,
	)
}

func NewTcpHttpServer(
	logger *zap.Logger,
	tracer tracing.Tracer,
	listener HasHttpListener,
	port Port,
) (result HttpServer, err error) {
	return internal.NewTcpHttpServer(logger, tracer, listener, port)
}

func NewTestTcpHttpServer(
	logger *zap.Logger,
	tracer tracing.Tracer,
	port Port,
) (result HttpServer, err error) {
	return internal.NewTestHttpServer(logger, tracer, port)
}

func NewUnixHttpServer(
	logger *zap.Logger,
	tracer tracing.Tracer,
	listener HasHttpListener,
	socket Socket,
) (result HttpServer, err error) {
	return internal.NewUnixHttpServer(logger, tracer, listener, socket)
}

func NewTcpConnectionMux(
	logger *zap.Logger,
	port Port,
) (result ConnectionMux, err error) {
	return internal.NewTcpConnectionMux(logger, port)
}

func NewTestTcpConnectionMux() (result ConnectionMux, err error) {
	return internal.NewTestTcpConnectionMux()
}

func NewTlsTcpConnectionMux(
	logger *zap.Logger,
	port Port,
	tlsCert string,
	tlsKey string,
) (result ConnectionMux, err error) {
	return internal.NewTlsTcpConnectionMux(logger, port, tlsCert, tlsKey)
}

func NewUnixConnectionMux(
	logger *zap.Logger,
	socket Socket,
) (result ConnectionMux, err error) {
	return internal.NewUnixConnectionMux(logger, socket)
}

func NewSocket(name string) Socket {
	return internal.NewSocket(name)
}

func NewTcpSocketIOServer(
	logger *zap.Logger,
	listener HasHttpListener,
	port Port,
	opt *engineio.Options,
) (result SocketIOServer, err error) {
	return internal.NewTcpSocketIOServer(logger, listener, port, opt)
}

func NewZinxTcpServer(
	logger *zap.Logger,
	listener HasHttpListener,
	port Port,
) (result ZinxTcpServer, err error) {
	return internal.NewZinxTcpServer(logger, listener, port)
}

func NewTcpGatewayServer(
	logger *zap.Logger,
	listener HasHttpListener,
	port Port,
) (result GatewayServer, err error) {
	return internal.NewTcpGatewayServer(logger, listener, port)
}
