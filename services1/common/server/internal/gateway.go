package internal

import (
	"context"
	"fmt"
	"golang.org/x/exp/slog"
	"net/http"
	"strings"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type TcpGatewayServer struct {
	logger   *zap.Logger
	server   *http.Server
	listener HasHttpListener
	port     Port
	mux      *runtime.ServeMux
	opts     []grpc.DialOption
}

func (s *TcpGatewayServer) StartServing(ctx context.Context) error {
	if listener, err := s.listener.HttpListener(); err != nil {
		return err
	} else {
		s.logger.Info(
			"serving grpc-gateway",
			zap.String("network", listener.Addr().Network()),
			zap.String("address", listener.Addr().String()),
			zap.Int("port", s.port.Value()),
		)
		go func() {
			if err := s.server.ListenAndServe(); err != nil {
				panic(err)
			}
		}()
	}

	return nil
}

func (s *TcpGatewayServer) StopServing(ctx context.Context) error {
	if err := s.server.Shutdown(ctx); err != nil {
		return err
	}
	return nil
}

func (s *TcpGatewayServer) GatewayServer() *http.Server {
	return s.server
}

func (s *TcpGatewayServer) GatewayRuntimeMux() *runtime.ServeMux {
	return s.mux
}

func (s *TcpGatewayServer) GatewayOption() []grpc.DialOption {
	return s.opts
}

func NewTcpGatewayServer(
	logger *zap.Logger,
	listener HasHttpListener,
	port Port,
) (result *TcpGatewayServer, err error) {
	mux := runtime.NewServeMux(
		runtime.WithIncomingHeaderMatcher(Matcher),
		runtime.WithOutgoingHeaderMatcher(Matcher),
	)
	opts := []grpc.DialOption{grpc.WithInsecure()}
	// 代理默认端口+1
	port = port + 1
	server := &http.Server{
		Addr:    fmt.Sprintf("0.0.0.0:%d", port.Value()),
		Handler: allowCORS(withLogger(mux)),
	}
	result = &TcpGatewayServer{
		logger:   logger,
		server:   server,
		listener: listener,
		port:     port,
		mux:      mux,
		opts:     opts,
	}
	return
}
func allowCORS(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if origin := r.Header.Get("Origin"); origin != "" {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			if r.Method == "OPTIONS" && r.Header.Get("Access-Control-Request-Method") != "" {
				preflightHandler(w, r)
				return
			}
		}
		h.ServeHTTP(w, r)
	})
}

func preflightHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")

	headers := []string{"*"}
	w.Header().Set("Access-Control-Allow-Headers", strings.Join(headers, ","))

	methods := []string{"GET", "HEAD", "POST", "PUT", "DELETE"}
	w.Header().Set("Access-Control-Allow-Methods", strings.Join(methods, ","))

	slog.Info("preflight request", "http_path", r.URL.Path)
}

func withLogger(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		slog.Info("Run request", "http_method", r.Method, "http_url", r.URL)
		h.ServeHTTP(w, r)
	})
}

func Matcher(key string) (string, bool) {
	switch key {
	case "Token":
		return key, true
	default:
		return runtime.DefaultHeaderMatcher(key)
	}
}
