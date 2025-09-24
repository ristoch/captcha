package tracing

import (
	"context"

	oteltrace "go.opentelemetry.io/otel/trace"
)

type contextKey string

const tracerKey contextKey = "tracer"

func WithContext(ctx context.Context) context.Context {
	return context.WithValue(ctx, tracerKey, Get())
}

func FromContext(ctx context.Context) oteltrace.Tracer {
	if tracer, ok := ctx.Value(tracerKey).(oteltrace.Tracer); ok {
		return tracer
	}
	return Get()
}
