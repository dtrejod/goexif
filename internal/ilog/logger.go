package ilog

import (
	"context"

	"go.uber.org/zap"
)

const typeValue = "service"

type serviceLogContextKeyType string

const contextKey = serviceLogContextKeyType(typeValue)

// WithLogger returns a copy of the provided context with the provided Logger included as a value.
func WithLogger(ctx context.Context, logger *zap.Logger) context.Context {
	return context.WithValue(ctx, contextKey, logger)
}

// FromContext returns a copy of the Logger from the provided context
func FromContext(ctx context.Context) *zap.Logger {
	if logger, ok := ctx.Value(contextKey).(*zap.Logger); ok {
		return logger
	}
	return defaultLogger()
}

func defaultLogger() *zap.Logger {
	logger, _ := zap.NewProduction()
	return logger
}
