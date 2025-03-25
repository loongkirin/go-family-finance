package middleware

import (
	"bytes"
	"io"
	"time"

	"github.com/gin-gonic/gin"
	pkglogger "github.com/loongkirin/go-family-finance/pkg/logger"
	"github.com/loongkirin/go-family-finance/pkg/util"
)

func Logger(logger pkglogger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		traceId := GetTraceID(c)
		if len(traceId) == 0 {
			traceId = util.GenerateId()
		}
		SetTraceID(c, traceId)

		requestId := GetRequestId(c)
		if len(requestId) == 0 {
			requestId = util.GenerateId()
		}
		SetRequestId(c, requestId)
		body, err := c.GetRawData()
		if err != nil {
			logger.Error("Failed to get raw data", pkglogger.Fields{"error": err})
			body = []byte{}
		}
		c.Request.Body = io.NopCloser(bytes.NewBuffer(body))
		// 处理请求
		c.Next()

		// 记录请求信息
		duration := time.Since(start).Microseconds()
		// ctxLogger := logger.With().Fields(map[string]interface{}{
		// 	"traceId":   traceId,
		// 	"requestId": requestId,
		// }).Logger()

		logger.Info("HTTP Request", pkglogger.Fields{
			"traceId":     traceId,
			"requestId":   requestId,
			"method":      c.Request.Method,
			"path":        c.Request.URL.Path,
			"status":      c.Writer.Status(),
			"duration":    duration,
			"clientIp":    c.ClientIP(),
			"userAgent":   c.Request.UserAgent(),
			"requestSize": c.Request.ContentLength,
			"requestBody": string(body),
			"headers":     c.Request.Header,
		})
	}
}
