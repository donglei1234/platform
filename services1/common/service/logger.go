package service

import (
	"context"

	"github.com/grpc-ecosystem/go-grpc-middleware/logging/zap/ctxzap"

	"go.uber.org/zap"
)

// Logger provides contextual logging for services.
type Logger struct{}

// LazyLoggerForContext returns a new LazyContextLogger for the given context.
func (l *Logger) LazyLoggerForContext(ctx context.Context) LazyContextLogger {
	return LazyContextLogger{
		ctx: ctx,
	}
}

// LoggerForContext creates and returns a new contextual logger for the given context.
func (l *Logger) LoggerForContext(ctx context.Context) *zap.Logger {
	return ctxzap.Extract(ctx)
}

// LazyContextLogger provides a lazy-initialized wrapper around LoggerForContext.  Use this when you may or may not need
// a logger based on service logic.
type LazyContextLogger struct {
	ctx    context.Context
	logger *zap.Logger
}

// Get returns the underlying zap.Logger for this context, creating it if it doesn't exist.
func (l *LazyContextLogger) Get() *zap.Logger {
	if l.logger == nil {
		l.logger = ctxzap.Extract(l.ctx)
	}
	return l.logger
}

// Debug writes a message and optional fields to the logger at Debug level.
func (l *LazyContextLogger) Debug(msg string, fields ...zap.Field) {
	l.Get().Debug(msg, fields...)
}

// Info writes a message and optional fields to the logger at Info level.
func (l *LazyContextLogger) Info(msg string, fields ...zap.Field) {
	l.Get().Info(msg, fields...)
}

// Warn writes a message and optional fields to the logger at Warn level.
func (l *LazyContextLogger) Warn(msg string, fields ...zap.Field) {
	l.Get().Warn(msg, fields...)
}

// Error writes a message and optional fields to the logger at Error level.
func (l *LazyContextLogger) Error(msg string, fields ...zap.Field) {
	l.Get().Error(msg, fields...)
}

// Panic writes a message and optional fields to the logger at Panic level.
func (l *LazyContextLogger) Panic(msg string, fields ...zap.Field) {
	l.Get().Panic(msg, fields...)
}

// Fatal writes a message and optional fields to the logger at Fatal level.
func (l *LazyContextLogger) Fatal(msg string, fields ...zap.Field) {
	l.Get().Fatal(msg, fields...)
}
