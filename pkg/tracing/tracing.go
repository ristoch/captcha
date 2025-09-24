package tracing

import (
	"context"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	oteltrace "go.opentelemetry.io/otel/trace"
)

var globalTracer oteltrace.Tracer

func Init(serviceName, jaegerEndpoint string) error {
	exporter, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(jaegerEndpoint)))
	if err != nil {
		return err
	}

	tp := trace.NewTracerProvider(
		trace.WithBatcher(exporter),
		trace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(serviceName),
			semconv.ServiceVersionKey.String("1.0.0"),
		)),
	)

	otel.SetTracerProvider(tp)
	globalTracer = tp.Tracer(serviceName)
	return nil
}

func Get() oteltrace.Tracer {
	if globalTracer == nil {
		globalTracer = otel.Tracer("captcha-service")
	}
	return globalTracer
}

func StartSpan(ctx context.Context, name string) (context.Context, oteltrace.Span) {
	return Get().Start(ctx, name)
}
