package ilog

import (
	"context"

	"go.uber.org/zap"
)

const TypeValue = "service"

type serviceLogContextKeyType string

const contextKey = serviceLogContextKeyType(TypeValue)

// WithLogger returns a copy of the provided context with the provided Logger included as a value.
func WithLogger(ctx context.Context, logger *zap.Logger) context.Context {
	return context.WithValue(ctx, contextKey, logger)
}

func FromContext(ctx context.Context) *zap.Logger {
	if logger, ok := ctx.Value(contextKey).(*zap.Logger); ok {
		return logger
	}
	return DefaultLogger()
}

func DefaultLogger() *zap.Logger {
	logger, _ := zap.NewProduction()
	return logger
}
