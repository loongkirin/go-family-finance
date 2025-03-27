package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/metric"
)

var (
	otelMeter = otel.GetMeterProvider().Meter("gin-http-metrics")

	// OpenTelemetry metrics
	otelRequestDuration, _ = otelMeter.Float64Histogram(
		"http.server.duration",
		metric.WithDescription("HTTP request duration in seconds"),
		metric.WithUnit("s"),
	)

	otelRequestSize, _ = otelMeter.Int64Histogram(
		"http.server.request.size",
		metric.WithDescription("HTTP request size in bytes"),
		metric.WithUnit("By"),
	)

	otelResponseSize, _ = otelMeter.Int64Histogram(
		"http.server.response.size",
		metric.WithDescription("HTTP response size in bytes"),
		metric.WithUnit("By"),
	)

	otelRequestsInFlight, _ = otelMeter.Int64UpDownCounter(
		"http.server.requests_in_flight",
		metric.WithDescription("Current number of HTTP requests being served"),
	)

	otelRequestsTotal, _ = otelMeter.Int64Counter(
		"http.server.requests_total",
		metric.WithDescription("Total number of HTTP requests"),
	)
)

// OtelMetrics 中间件用于收集 OpenTelemetry 指标
func OtelMetrics() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		ctx := c.Request.Context()

		// 记录请求大小
		requestSize := c.Request.ContentLength
		if requestSize > 0 {
			otelRequestSize.Record(ctx, requestSize)
		}

		// 增加当前请求计数
		otelRequestsInFlight.Add(ctx, 1)
		defer otelRequestsInFlight.Add(ctx, -1)

		// 处理请求
		c.Next()

		// 记录响应大小
		responseSize := int64(c.Writer.Size())
		if responseSize > 0 {
			otelResponseSize.Record(ctx, responseSize)
		}

		// 记录请求持续时间
		duration := time.Since(start).Seconds()
		otelRequestDuration.Record(ctx, duration)

		// 记录请求总数
		otelRequestsTotal.Add(ctx, 1)
	}
}
