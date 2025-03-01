package middleware

import (
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

func Tracing(tracer trace.Tracer) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		spanCtx, span := tracer.Start(ctx, c.Request.URL.Path,
			trace.WithAttributes(
				attribute.String("http.method", c.Request.Method),
				attribute.String("http.url", c.Request.URL.String()),
				attribute.String("traceId", GetTraceID(c)),
				attribute.String("requestId", GetRequestId(c)),
			),
		)
		defer span.End()

		// 将 span 上下文传递给后续处理器
		c.Request = c.Request.WithContext(spanCtx)
		c.Next()

		// 记录响应状态
		span.SetAttributes(
			attribute.Int("http.status_code", c.Writer.Status()),
		)
	}
}
