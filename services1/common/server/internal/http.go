package internal

import (
	"context"
	"github.com/donglei1234/platform/services/common/tracing"
	"net/http"
	"net/http/httptest"

	"go.uber.org/zap"
)

type TcpHttpServer struct {
	logger   *zap.Logger
	server   *http.Server
	mux      HttpServeMux
	listener HasHttpListener
	port     Port
}

func (s *TcpHttpServer) StartServing(ctx context.Context) error {
	if listener, err := s.listener.HttpListener(); err != nil {
		return err
	} else {
		s.logger.Info(
			"serving http",
			zap.String("network", listener.Addr().Network()),
			zap.String("address", listener.Addr().String()),
			zap.Int("port", s.port.Value()),
		)

		go func() {
			if err := s.server.Serve(listener); err != http.ErrServerClosed {
				panic(err)
			}
		}()
	}

	return nil
}

func (s *TcpHttpServer) StopServing(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}

func (s *TcpHttpServer) HttpServeMux() HttpServeMux {
	return s.mux
}

func (s *TcpHttpServer) Port() Port {
	return s.port
}

func NewTcpHttpServer(
	logger *zap.Logger,
	tracer tracing.Tracer,
	listener HasHttpListener,
	port Port,
) (result *TcpHttpServer, err error) {
	mux := tracer.ServeMux()
	result = &TcpHttpServer{
		logger:   logger,
		listener: listener,
		port:     port,
		mux:      mux,
		server: &http.Server{
			Handler: mux,
		},
	}
	return
}

type UnixHttpServer struct {
	logger   *zap.Logger
	server   *http.Server
	mux      HttpServeMux
	listener HasHttpListener
	socket   Socket
}

func (s *UnixHttpServer) StartServing(ctx context.Context) error {
	if listener, err := s.listener.HttpListener(); err != nil {
		return err
	} else {
		s.logger.Info(
			"serving http",
			zap.String("network", listener.Addr().Network()),
			zap.String("address", listener.Addr().String()),
		)

		go func() {
			if err := s.server.Serve(listener); err != nil {
				panic(err)
			}
		}()
	}
	return nil
}

func (s *UnixHttpServer) StopServing(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}

func (s *UnixHttpServer) HttpServeMux() HttpServeMux {
	return s.mux
}

func (s *UnixHttpServer) Socket() Socket {
	return s.socket
}

func NewUnixHttpServer(
	logger *zap.Logger,
	tracer tracing.Tracer,
	listener HasHttpListener,
	socket Socket,
) (result *UnixHttpServer, err error) {
	mux := tracer.ServeMux()
	result = &UnixHttpServer{
		logger:   logger,
		listener: listener,
		socket:   socket,
		mux:      mux,
		server: &http.Server{
			Handler: mux,
		},
	}
	return
}

type TestHttpServer struct {
	logger *zap.Logger
	server *httptest.Server
	mux    HttpServeMux
	port   Port
}

func (s *TestHttpServer) StartServing(ctx context.Context) error {
	s.logger.Info(
		"serving test http",
		zap.Int("port", s.port.Value()),
	)
	s.server.Start()
	return nil
}

func (s *TestHttpServer) StopServing(ctx context.Context) error {
	s.server.Close()
	return nil
}

func (s *TestHttpServer) HttpServeMux() HttpServeMux {
	return s.mux
}

func NewTestHttpServer(
	logger *zap.Logger,
	tracer tracing.Tracer,
	port Port,
) (result *TestHttpServer, err error) {
	mux := tracer.ServeMux()
	result = &TestHttpServer{
		logger: logger,
		port:   port,
		mux:    mux,
		server: httptest.NewUnstartedServer(mux),
	}
	return
}
