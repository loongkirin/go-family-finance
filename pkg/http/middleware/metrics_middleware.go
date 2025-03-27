package middleware

import (
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/prometheus/client_golang/prometheus"
)

var (
	httpRequestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests",
		},
		[]string{"method", "path", "status"},
	)

	httpRequestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "HTTP request duration in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method", "path"},
	)

	httpRequestsInFlight = prometheus.NewGauge(
		prometheus.GaugeOpts{
			Name: "http_requests_in_flight",
			Help: "Current number of HTTP requests being served",
		},
	)
)

func init() {
	prometheus.MustRegister(httpRequestsTotal)
	prometheus.MustRegister(httpRequestDuration)
	prometheus.MustRegister(httpRequestsInFlight)
}

// Metrics 中间件用于收集 HTTP 请求的 Prometheus 指标
func Metrics() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		// 增加当前请求计数
		httpRequestsInFlight.Inc()
		defer httpRequestsInFlight.Dec()

		// 处理请求
		c.Next()

		// 记录请求持续时间
		duration := time.Since(start).Seconds()
		httpRequestDuration.WithLabelValues(c.Request.Method, c.FullPath()).Observe(duration)

		// 记录请求总数
		status := strconv.Itoa(c.Writer.Status())
		httpRequestsTotal.WithLabelValues(c.Request.Method, c.FullPath(), status).Inc()
	}
}
