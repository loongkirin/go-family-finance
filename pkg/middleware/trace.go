package middleware

import (
	"bytes"
	"encoding/json"
	"io"
	"net"
	"net/http"
	"time"

	"LoongKirin/go-family-finance/pkg/logger"
	"bufio"

	"github.com/gin-gonic/gin"
)

// TraceMiddleware 追踪中间件
func TraceMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 开始时间
		start := time.Now()

		// 记录请求体
		var requestBody []byte
		if c.Request.Body != nil {
			requestBody, _ = io.ReadAll(c.Request.Body)
			// 重新设置请求体，因为读取后需要重置
			c.Request.Body = io.NopCloser(bytes.NewBuffer(requestBody))
		}

		// 记录响应体
		writer := &responseWriter{
			ResponseWriter: c.Writer,
			body:           &bytes.Buffer{},
		}
		c.Writer = writer

		// 处理请求
		c.Next()

		// 结束时间
		end := time.Now()
		duration := end.Sub(start)

		// 构建追踪信息
		traceInfo := map[string]interface{}{
			"method":     c.Request.Method,
			"path":       c.Request.URL.Path,
			"query":      c.Request.URL.RawQuery,
			"ip":         c.ClientIP(),
			"user_agent": c.Request.UserAgent(),
			"duration":   duration.String(),
			"status":     c.Writer.Status(),
		}

		// 添加请求头
		headers := make(map[string]string)
		for k, v := range c.Request.Header {
			if len(v) > 0 {
				headers[k] = v[0]
			}
		}
		traceInfo["headers"] = headers

		// 添加请求体
		if len(requestBody) > 0 {
			var prettyJSON bytes.Buffer
			if err := json.Indent(&prettyJSON, requestBody, "", "  "); err == nil {
				traceInfo["request_body"] = prettyJSON.String()
			} else {
				traceInfo["request_body"] = string(requestBody)
			}
		}

		// 添加响应体
		if writer.body.Len() > 0 {
			var prettyJSON bytes.Buffer
			if err := json.Indent(&prettyJSON, writer.body.Bytes(), "", "  "); err == nil {
				traceInfo["response_body"] = prettyJSON.String()
			} else {
				traceInfo["response_body"] = writer.body.String()
			}
		}

		// 添加错误信息
		if len(c.Errors) > 0 {
			errors := make([]string, len(c.Errors))
			for i, err := range c.Errors {
				errors[i] = err.Error()
			}
			traceInfo["errors"] = errors
		}

		// 记录追踪信息
		logger.Info(c.Request.Context(), "HTTP Request Trace", traceInfo)
	}
}

// responseWriter 自定义响应写入器
type responseWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

// Write 重写Write方法以捕获响应内容
func (w *responseWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

// WriteString 重写WriteString方法以捕获响应内容
func (w *responseWriter) WriteString(s string) (int, error) {
	w.body.WriteString(s)
	return w.ResponseWriter.WriteString(s)
}

// WriteHeader 重写WriteHeader方法
func (w *responseWriter) WriteHeader(code int) {
	w.ResponseWriter.WriteHeader(code)
}

// Flush 实现Flush接口
func (w *responseWriter) Flush() {
	w.ResponseWriter.Flush()
}

// Hijack 实现Hijack接口
func (w *responseWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	return w.ResponseWriter.Hijack()
}

// CloseNotify 实现CloseNotify接口
func (w *responseWriter) CloseNotify() <-chan bool {
	return w.ResponseWriter.CloseNotify()
}

// Pusher 实现Pusher接口
func (w *responseWriter) Pusher() http.Pusher {
	return w.ResponseWriter.Pusher()
}
