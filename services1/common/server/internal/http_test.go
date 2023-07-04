package internal

import (
	"testing"

	"github.com/donglei1234/platform/services/common/tracing"
	"github.com/donglei1234/platform/services/common/tracing/noop"

	"context"

	"net"

	"go.uber.org/zap"
)

type HttpServerFactory func(logger *zap.Logger, tracer tracing.Tracer, listener HasHttpListener, port Port) (HttpServer, error)

type HttpListener struct{}

func (l HttpListener) HttpListener() (net.Listener, error) {
	return newListener(), nil
}

func testNewHttpServer(t *testing.T, h HttpServerFactory) {
	listener := HttpListener{}
	logger := zap.NewNop()
	tracer := noop.NewTracer()

	type TestCase struct {
		listener HasHttpListener
		port     Port
		logger   *zap.Logger
		tracer   tracing.Tracer
		validate func(i int, server HttpServer, err error)
	}

	validate := func(i int, server HttpServer, err error) {
		if server == nil || err != nil {
			t.Fatal(i, err, server)
		}
	}

	cases := []TestCase{
		{
			listener: listener,
			port:     newPort(),
			logger:   logger,
			tracer:   tracer,
			validate: validate,
		},
	}

	for i, tc := range cases {
		server, err := h(tc.logger, tc.tracer, tc.listener, tc.port)
		tc.validate(i, server, err)
	}
}

func createNewTcpHttpServer(t *testing.T) *TcpHttpServer {
	listener := HttpListener{}
	port := newPort()
	logger := zap.NewNop()
	tracer := noop.NewTracer()

	testServer, err := NewTcpHttpServer(logger, tracer, listener, port)
	if err != nil {
		t.Fatal("Couldn't create the grpc server:", err)
	}
	return testServer
}

func createNewTestHttpServer(t *testing.T) *TestHttpServer {
	port := newPort()
	logger := zap.NewNop()
	tracer := noop.NewTracer()

	testServer, err := NewTestHttpServer(logger, tracer, port)
	if err != nil {
		t.Fatal("Couldn't create the grpc server:", err)
	}
	return testServer
}

func TestNewTcpHttpServer(t *testing.T) {
	testNewHttpServer(t, func(logger *zap.Logger, tracer tracing.Tracer, listener HasHttpListener, port Port) (HttpServer, error) {
		return NewTcpHttpServer(logger, tracer, listener, port)
	})
}

func TestNewTestHttpServer(t *testing.T) {
	testNewHttpServer(t, func(logger *zap.Logger, tracer tracing.Tracer, listener HasHttpListener, port Port) (HttpServer, error) {
		return NewTestHttpServer(logger, tracer, port)
	})
}

func TestTcpHttpServerServing(t *testing.T) {
	ctx := context.Background()
	testServer := createNewTcpHttpServer(t)
	testServer.StartServing(ctx)
	testServer.StopServing(ctx)
}

func TestTestHttpServerServing(t *testing.T) {
	ctx := context.Background()
	testServer := createNewTestHttpServer(t)
	testServer.StartServing(ctx)
	testServer.StopServing(ctx)
}
