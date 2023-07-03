package internal

import (
	"context"
	"github.com/aceld/zinx/ziface"
	"net"
	"net/http"

	socketio "github.com/googollee/go-socket.io"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
)

type Server interface {
	StartServing(ctx context.Context) error
	StopServing(ctx context.Context) error
}

type TcpServer interface {
	Server
	Port() Port
}

type UnixServer interface {
	Server
	Socket() Socket
}

type GrpcServer interface {
	Server
	GrpcServer() *grpc.Server
	Dial(target string, opts ...grpc.DialOption) (conn *grpc.ClientConn, err error)
}

type HttpServer interface {
	Server
	HttpServeMux() HttpServeMux
}

type SocketIOServer interface {
	HttpServer
	SocketIOServer() *socketio.Server
}

type ZinxServer interface {
	TcpServer
	ZinxTcpServer() ziface.IServer
}

type GatewayServer interface {
	Server
	GatewayServer() *http.Server
	GatewayRuntimeMux() *runtime.ServeMux
	GatewayOption() []grpc.DialOption
}

type ConnectionMux interface {
	Server
	HasGrpcListener
	HasHttpListener
	Network() Network
}

type HasGrpcListener interface {
	GrpcListener() (net.Listener, error)
}

type HasHttpListener interface {
	HttpListener() (net.Listener, error)
}

type HttpServeMux interface {
	http.Handler

	Handle(string, http.Handler)
	HandleFunc(pattern string, handler func(http.ResponseWriter, *http.Request))
}
