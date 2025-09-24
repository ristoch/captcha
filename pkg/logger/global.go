package logger

import (
	"context"

	"go.uber.org/zap"
)

type contextKey string

const loggerKey contextKey = "logger"

func WithContext(ctx context.Context) context.Context {
	return context.WithValue(ctx, loggerKey, Get())
}

func FromContext(ctx context.Context) *zap.Logger {
	if logger, ok := ctx.Value(loggerKey).(*zap.Logger); ok {
		return logger
	}
	return Get()
}
