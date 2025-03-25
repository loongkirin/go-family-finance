package middleware

import (
	"bytes"
	"io"
	"time"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

func Tracing(tracer trace.Tracer) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		ctx := c.Request.Context()
		body, err := c.GetRawData()
		if err != nil {
			body = []byte{}
		}

		spanCtx, span := tracer.Start(ctx, c.Request.URL.Path,
			trace.WithAttributes(
				attribute.String("http.method", c.Request.Method),
				attribute.String("http.url", c.Request.URL.String()),
				attribute.String("traceId", GetTraceID(c)),
				attribute.String("requestId", GetRequestId(c)),
				attribute.Int64("requestSize", c.Request.ContentLength),
				attribute.String("clientIp", c.ClientIP()),
				attribute.String("userAgent", c.Request.UserAgent()),
				attribute.String("requestBody", string(body)),
			),
		)
		defer span.End()

		c.Request.Body = io.NopCloser(bytes.NewBuffer(body))

		// 将 span 上下文传递给后续处理器
		c.Request = c.Request.WithContext(spanCtx)
		c.Next()

		duration := time.Since(start)
		// 记录响应状态
		span.SetAttributes(
			attribute.Int("http.status_code", c.Writer.Status()),
			attribute.Int("http.duration", int(duration.Milliseconds())),
		)
	}
}
