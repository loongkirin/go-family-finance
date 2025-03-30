package telemetry

import (
	"context"
	"fmt"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp"
	"go.opentelemetry.io/otel/metric"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.27.0"
)

var (
	meter = otel.GetMeterProvider().Meter("family-finance")

	// HTTP 指标
	httpRequestDuration, _ = meter.Float64Histogram(
		"http.server.duration",
		metric.WithDescription("HTTP request duration in seconds"),
		metric.WithUnit("s"),
	)

	httpRequestSize, _ = meter.Int64Histogram(
		"http.server.request.size",
		metric.WithDescription("HTTP request size in bytes"),
		metric.WithUnit("By"),
	)

	httpResponseSize, _ = meter.Int64Histogram(
		"http.server.response.size",
		metric.WithDescription("HTTP response size in bytes"),
		metric.WithUnit("By"),
	)

	httpRequestsInFlight, _ = meter.Int64UpDownCounter(
		"http.server.requests_in_flight",
		metric.WithDescription("Current number of HTTP requests being served"),
	)

	httpRequestsTotal, _ = meter.Int64Counter(
		"http.server.requests_total",
		metric.WithDescription("Total number of HTTP requests"),
	)

	// 数据库指标
	dbQueryDuration, _ = meter.Float64Histogram(
		"db.query.duration",
		metric.WithDescription("Database query duration in seconds"),
		metric.WithUnit("s"),
	)

	dbConnectionsInUse, _ = meter.Int64UpDownCounter(
		"db.connections.in_use",
		metric.WithDescription("Number of database connections currently in use"),
	)

	dbConnectionsTotal, _ = meter.Int64Counter(
		"db.connections.total",
		metric.WithDescription("Total number of database connections"),
	)

	// 业务指标
	userLoginAttempts, _ = meter.Int64Counter(
		"user.login.attempts",
		metric.WithDescription("Total number of user login attempts"),
	)

	userLoginSuccess, _ = meter.Int64Counter(
		"user.login.success",
		metric.WithDescription("Total number of successful user logins"),
	)

	userLoginFailures, _ = meter.Int64Counter(
		"user.login.failures",
		metric.WithDescription("Total number of failed user logins"),
	)

	// 系统指标
	systemMemoryUsage, _ = meter.Int64UpDownCounter(
		"system.memory.usage",
		metric.WithDescription("Current system memory usage in bytes"),
		metric.WithUnit("By"),
	)

	systemCPUUsage, _ = meter.Float64UpDownCounter(
		"system.cpu.usage",
		metric.WithDescription("Current system CPU usage percentage"),
		metric.WithUnit("%"),
	)

	systemGoroutines, _ = meter.Int64UpDownCounter(
		"system.goroutines",
		metric.WithDescription("Current number of goroutines"),
	)
)

// RecordHTTPRequest 记录 HTTP 请求指标
func RecordHTTPRequest(ctx context.Context, method, path string, status int, duration time.Duration, requestSize, responseSize int64) {
	attrs := attribute.NewSet(
		attribute.String("http.method", method),
		attribute.String("http.route", path),
		attribute.Int("http.status_code", status),
	)

	httpRequestDuration.Record(ctx, duration.Seconds(), metric.WithAttributes(attrs.ToSlice()...))
	httpRequestSize.Record(ctx, requestSize, metric.WithAttributes(attrs.ToSlice()...))
	httpResponseSize.Record(ctx, responseSize, metric.WithAttributes(attrs.ToSlice()...))
	httpRequestsTotal.Add(ctx, 1, metric.WithAttributes(attrs.ToSlice()...))
}

// RecordDBQuery 记录数据库查询指标
func RecordDBQuery(ctx context.Context, operation string, duration time.Duration) {
	attrs := attribute.NewSet(
		attribute.String("db.operation", operation),
	)

	dbQueryDuration.Record(ctx, duration.Seconds(), metric.WithAttributes(attrs.ToSlice()...))
}

// RecordUserLogin 记录用户登录指标
func RecordUserLogin(ctx context.Context, success bool) {
	attrs := attribute.NewSet(
		attribute.Bool("success", success),
	)

	userLoginAttempts.Add(ctx, 1, metric.WithAttributes(attrs.ToSlice()...))
	if success {
		userLoginSuccess.Add(ctx, 1, metric.WithAttributes(attrs.ToSlice()...))
	} else {
		userLoginFailures.Add(ctx, 1, metric.WithAttributes(attrs.ToSlice()...))
	}
}

// RecordSystemMetrics 记录系统指标
func RecordSystemMetrics(ctx context.Context, memoryUsage int64, cpuUsage float64, goroutines int64) {
	systemMemoryUsage.Add(ctx, memoryUsage)
	systemCPUUsage.Add(ctx, cpuUsage)
	systemGoroutines.Add(ctx, goroutines)
}

// UpdateDBConnections 更新数据库连接指标
func UpdateDBConnections(ctx context.Context, inUse, total int64) {
	dbConnectionsInUse.Add(ctx, inUse)
	dbConnectionsTotal.Add(ctx, total)
}

func InitMetrics(ctx context.Context, config TelemetryConfig) (metric.Meter, error) {
	// 创建资源
	res, err := resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceNameKey.String(config.ServiceName),
			semconv.ServiceVersionKey.String(config.ServiceVersion),
			semconv.DeploymentEnvironmentNameKey.String(config.ServiceEnvironment),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create resource: %w", err)
	}

	// 创建指标导出器
	var exporter sdkmetric.Exporter
	switch config.CollectorType {
	case "grpc":
		exporter, err = otlpmetricgrpc.New(ctx,
			otlpmetricgrpc.WithEndpoint(config.CollectorURL),
			otlpmetricgrpc.WithTimeout(time.Second*10),
		)
	case "http":
		exporter, err = otlpmetrichttp.New(ctx,
			otlpmetrichttp.WithEndpoint(config.CollectorURL),
			otlpmetrichttp.WithTimeout(time.Second*10),
		)
	default:
		return nil, fmt.Errorf("unsupported exporter type: %s", config.CollectorType)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to create exporter: %w", err)
	}

	// 创建指标处理器
	processor := sdkmetric.NewPeriodicReader(exporter,
		sdkmetric.WithInterval(time.Second*10),
	)

	// 创建指标提供者
	provider := sdkmetric.NewMeterProvider(
		sdkmetric.WithResource(res),
		sdkmetric.WithReader(processor),
	)

	// 设置全局指标提供者
	otel.SetMeterProvider(provider)

	return provider.Meter(config.ServiceName), nil
}

// ShutdownMetrics 关闭指标收集
func ShutdownMetrics(ctx context.Context) error {
	if provider, ok := otel.GetMeterProvider().(*sdkmetric.MeterProvider); ok {
		return provider.Shutdown(ctx)
	}
	return nil
}
