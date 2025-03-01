package telemetry

import (
	"context"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.27.0"
)

func InitTracer(cfg TelemetryConfig) (*sdktrace.TracerProvider, error) {
	// 创建 OTLP exporter
	exporter, err := otlptracegrpc.New(context.Background(),
		otlptracegrpc.WithEndpoint(cfg.CollectorURL),
		otlptracegrpc.WithInsecure(),
	)
	if err != nil {
		return nil, err
	}

	// 创建资源属性
	resource := resource.NewWithAttributes(
		semconv.SchemaURL,
		semconv.ServiceNameKey.String(cfg.ServiceName),
	)

	// 创建 TracerProvider
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(resource),
		sdktrace.WithSampler(sdktrace.TraceIDRatioBased(cfg.TraceSample)),
	)

	// 设置全局 TracerProvider
	otel.SetTracerProvider(tp)

	// 设置全局 propagator
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	))

	return tp, nil
}
