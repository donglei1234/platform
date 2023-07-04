package internal

import (
	"net"

	"context"

	"crypto/tls"
	"encoding/base64"
	"os"

	"github.com/pkg/errors"
	"github.com/soheilhy/cmux"
	//cert "tata-ol/services/certificates/cmd"
	"go.uber.org/zap"
	"google.golang.org/grpc/test/bufconn"
)

var (
	httpMatcher     cmux.Matcher
	grpcMatchWriter cmux.MatchWriter
)

func init() {
	grpcMatchWriter = cmux.HTTP2MatchHeaderFieldSendSettings("content-type", "application/grpc")
	httpMatcher = cmux.HTTP1Fast()
}

type TcpConnectionMux struct {
	logger    *zap.Logger
	listener  net.Listener
	mux       cmux.CMux
	port      Port
	tlsConfig *tls.Config
}

func NewTcpConnectionMux(logger *zap.Logger, port Port) (result *TcpConnectionMux, err error) {
	result = &TcpConnectionMux{
		logger: logger,
		port:   port,
	}

	return
}

func NewTlsTcpConnectionMux(
	logger *zap.Logger,
	port Port,
	tlsCert string,
	tlsKey string,
) (result *TcpConnectionMux, err error) {
	if config, e := makeTlsConfig(tlsCert, tlsKey); e != nil {
		err = e
	} else {
		result = &TcpConnectionMux{
			logger:    logger,
			port:      port,
			tlsConfig: config,
		}
	}

	return
}

func makeTlsConfig(tlsCert, tlsKey string) (config *tls.Config, err error) {
	servercert, _ := base64.StdEncoding.DecodeString(tlsCert)
	serverkey, _ := base64.StdEncoding.DecodeString(tlsKey)
	//cacert, _ := base64.StdEncoding.DecodeString(cab64)

	if certificate, e := tls.X509KeyPair(servercert, serverkey); e != nil {
		err = e
	} else {
		//certPool := x509.NewCertPool()
		//certPool.AppendCertsFromPEM([]byte(cert.GetRootCerts()))

		config = &tls.Config{
			//ClientAuth:   tls.RequireAndVerifyClientCert,
			Certificates: []tls.Certificate{certificate},
			//RootCAs:    certPool,
		}
	}

	return
}

func (m *TcpConnectionMux) Network() Network {
	return TcpNetwork
}

func (m *TcpConnectionMux) GrpcListener() (listener net.Listener, err error) {
	if m.mux == nil {
		err = errors.New("connection mux is not serving")
	} else {
		listener = m.mux.MatchWithWriters(grpcMatchWriter)
	}

	return
}

func (m *TcpConnectionMux) HttpListener() (listener net.Listener, err error) {
	if m.mux == nil {
		err = errors.New("connection mux is not serving")
	} else {
		listener = m.mux.Match(httpMatcher)
	}

	return
}

func (m *TcpConnectionMux) StartServing(ctx context.Context) error {
	if listener, err := NewTcpListener(m.port); err != nil {
		return err
	} else {
		if m.tlsConfig != nil {
			m.listener = tls.NewListener(listener, m.tlsConfig)
		} else {
			m.listener = listener
		}
	}

	m.logger.Info(
		"multiplexing traffic",
		zap.String("network", m.listener.Addr().Network()),
		zap.String("address", m.listener.Addr().String()),
		zap.Int("port", m.port.Value()),
		zap.Bool("tls", m.tlsConfig != nil),
	)

	m.mux = cmux.New(m.listener)

	go func() {
		m.mux.Serve()
	}()

	return nil
}

func (m *TcpConnectionMux) StopServing(ctx context.Context) error {
	return m.listener.Close()
}

func (m *TcpConnectionMux) Port() Port {
	return m.port
}

type UnixConnectionMux struct {
	logger   *zap.Logger
	listener net.Listener
	mux      cmux.CMux
	socket   Socket
}

func NewUnixConnectionMux(logger *zap.Logger, socket Socket) (result *UnixConnectionMux, err error) {
	result = &UnixConnectionMux{
		logger: logger,
		socket: socket,
	}

	return
}

func (m *UnixConnectionMux) Network() Network {
	return UnixNetwork
}

func (m *UnixConnectionMux) GrpcListener() (listener net.Listener, err error) {
	if m.mux == nil {
		err = errors.New("connection mux is not serving")
	} else {
		listener = m.mux.MatchWithWriters(grpcMatchWriter)
	}

	return
}

func (m *UnixConnectionMux) HttpListener() (listener net.Listener, err error) {
	if m.mux == nil {
		err = errors.New("connection mux is not serving")
	} else {
		listener = m.mux.Match(httpMatcher)
	}

	return
}

func (m *UnixConnectionMux) StartServing(ctx context.Context) error {
	if listener, err := NewUnixListener(m.socket); err != nil {
		return err
	} else {
		m.listener = listener
	}

	m.logger.Info(
		"multiplexing traffic",
		zap.String("network", m.listener.Addr().Network()),
		zap.String("address", m.listener.Addr().String()),
		zap.String("socket", m.socket.ListenAddress()),
	)

	m.mux = cmux.New(m.listener)

	go func() {
		m.mux.Serve()
	}()

	return nil
}

func (m *UnixConnectionMux) StopServing(ctx context.Context) error {
	err := m.listener.Close()
	if e := os.Remove(m.socket.ListenAddress()); e != nil {
		m.logger.Warn(
			"error removing unix socket",
			zap.String("socket", m.socket.ListenAddress()),
			zap.Error(e),
		)
	}

	return err
}

func (m *UnixConnectionMux) Socket() Socket {
	return m.socket
}

type TestTcpConnectionMux struct {
	listener *bufconn.Listener
}

func NewTestTcpConnectionMux() (result *TestTcpConnectionMux, err error) {
	result = &TestTcpConnectionMux{
		listener: bufconn.Listen(256 * 1024),
	}
	return
}

func (m *TestTcpConnectionMux) Network() Network {
	return TcpNetwork
}

func (m *TestTcpConnectionMux) Port() Port {
	return Port(0)
}

func (m *TestTcpConnectionMux) GrpcListener() (net.Listener, error) {
	return m.listener, nil
}

func (m *TestTcpConnectionMux) HttpListener() (net.Listener, error) {
	return m.listener, nil
}

func (m *TestTcpConnectionMux) StartServing(ctx context.Context) error {
	return nil
}

func (m *TestTcpConnectionMux) StopServing(ctx context.Context) error {
	return nil
}

func (m *TestTcpConnectionMux) Dial() (net.Conn, error) {
	return m.listener.Dial()
}
