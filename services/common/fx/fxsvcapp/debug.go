package fxsvcapp

import (
	"github.com/donglei1234/platform/services/common/access"
	"github.com/donglei1234/platform/services/common/service"
	"net/http/pprof"
)

type TcpDebugHttpService struct {
	service.TcpTransport
}

func (s *TcpDebugHttpService) RegisterWithHttpServer(server service.HasHttpServeMux) {
	h := server.HttpServeMux()
	h.HandleFunc("/debug/pprof/", pprof.Index)
	h.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	h.HandleFunc("/debug/pprof/profile", pprof.Profile)
	h.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	h.HandleFunc("/debug/pprof/trace", pprof.Trace)
}

func NewTcpDebugHttpService() (out service.HttpServiceFactory) {
	out.HttpService = &TcpDebugHttpService{}
	return
}

func (s *TcpDebugHttpService) AccessLevel() access.AccessLevel {
	return access.AccessUndefined
}
