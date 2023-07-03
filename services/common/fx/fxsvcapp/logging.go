package fxsvcapp

import (
	"reflect"

	"github.com/donglei1234/platform/services/common/server"
	"github.com/donglei1234/platform/services/common/service"
	"go.uber.org/zap"
)

// LogCreateResource generates a standard log message for creating a local resource.
func LogCreateResource(l *zap.Logger, tag string, ctx string) {
	l.Info("create",
		zap.String("kind", "resource"),
		zap.String("tag", tag),
		zap.String("context", ctx),
	)
}

// LogConnectService generates a standard log message for connecting to a remote service.
func LogConnectService(l *zap.Logger, service string, url string) {
	l.Info("connect",
		zap.String("kind", "service"),
		zap.String("tag", service),
		zap.String("url", url),
	)
}

// LogOpenDocumentStore generates a standard log message for opening a document store.
func LogOpenDocumentStore(l *zap.Logger, name string, tag string) {
	l.Info("open",
		zap.String("kind", "nosql.DocumentStore"),
		zap.String("name", name),
		zap.String("tag", tag),
	)
}

// LogOpenMemoryStore generates a standard log message for opening a document store.
func LogOpenMemoryStore(l *zap.Logger, name string, tag string) {
	l.Info("open",
		zap.String("kind", "nosql.MemoryStore"),
		zap.String("name", name),
		zap.String("tag", tag),
	)
}

// LogBindService generates a standard log message when binding a service to a server.
func LogBindService(l *zap.Logger, svc service.Service, svr server.Server) {
	svcType := reflect.TypeOf(svc).Elem()
	svrType := reflect.TypeOf(svr).Elem()
	l.Info("bind",
		zap.String("service", svcType.String()),
		zap.String("server", svrType.Name()),
		zap.Stringer("transport", svc.ServiceTransport()),
	)
}

func LogLoadVault(l *zap.Logger) {

}
