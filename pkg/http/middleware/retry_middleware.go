package middleware

import (
	"fmt"
	"net"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"time"

	"math"

	"github.com/avast/retry-go"
	"github.com/gin-gonic/gin"
	"github.com/loongkirin/go-family-finance/pkg/http/response"
	pkglogger "github.com/loongkirin/go-family-finance/pkg/logger"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func Retry(logger pkglogger.Logger, maxRetries uint, retryDelay time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		traceId := GetTraceID(c)
		requestId := GetRequestId(c)
		err := retry.Do(
			func() error {
				ctx := c.Copy()
				// 执行下一个处理程序
				ctx.Next()

				// 如果没有错误，直接返回
				// if ctx.Writer.Status() == http.StatusOK || ctx.Writer.Status() == http.StatusNotFound {
				// 	return nil
				// }

				// // 返回错误
				// lastErr = fmt.Errorf("request failed with status %d", ctx.Writer.Status())
				// return lastErr

				status := ctx.Writer.Status()
				fmt.Println("status", status)
				if status >= http.StatusInternalServerError {
					return fmt.Errorf("http status %d", status)
				}
				return nil
			},
			retry.Attempts(maxRetries), // 重试次数
			retry.Delay(retryDelay),    // 重试间隔
			retry.MaxJitter(retryDelay),
			retry.MaxDelay(15*time.Second),
			retry.DelayType(func(n uint, err error, cof *retry.Config) time.Duration {
				return time.Duration(math.Pow(2, float64(n))) * time.Second
			}), //指数退避策略
			retry.RetryIf(func(err error) bool {
				// 只对网络错误和5xx状态码重试
				if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
					return true
				}
				if statusErr, ok := err.(interface{ StatusCode() int }); ok {
					return statusErr.StatusCode() >= http.StatusInternalServerError
				}
				return false
			}),
			retry.OnRetry(func(n uint, err error) {
				fmt.Println("retry", n, err)
				errMsg := fmt.Sprintf("Retry #%d: %s\n", n, err)
				logger.Error("request failed with retry", pkglogger.Fields{
					"error":     err,
					"traceId":   traceId,
					"requestId": requestId,
					"method":    c.Request.Method,
					"path":      c.Request.URL.Path,
					"message":   errMsg,
				})
				metrics.Increment("http_retries_total",
					"path", c.Request.URL.Path,
					"status", fmt.Sprint(getStatusCode(err)),
				)
				span := trace.SpanFromContext(c)
				span.AddEvent(fmt.Sprintf("Retry #%d", n),
					trace.WithAttributes(
						attribute.String("error", err.Error()),
					))
			}),
			retry.Context(c.Request.Context()),
		)

		if err != nil {
			c.AbortWithStatusJSON(c.Writer.Status(), response.NewResponse(response.ERROR, fmt.Sprintf("Request failed after retries: %s", err.Error())))
			// c.JSON(http.StatusInternalServerError, response.NewResponse(c.Writer.Status(), fmt.Sprintf("Request failed after retries: %s", err.Error())))
		}
	}
}

func getStatusCode(err error) int {
	// 类型断言优先顺序很重要
	switch e := err.(type) {
	case interface{ StatusCode() int }: // 自定义错误类型
		return e.StatusCode()

	case *url.Error: // http.Client 返回的错误
		if e.Err != nil {
			return getStatusCode(e.Err) // 递归处理
		}
		return http.StatusBadGateway

	case net.Error:
		if e.Timeout() {
			return http.StatusGatewayTimeout
		}
		return http.StatusServiceUnavailable

	default:
		// 处理 gRPC 状态码（如果使用 gRPC）
		if s, ok := status.FromError(err); ok {
			return httpStatusCodeFromGRPC(s.Code())
		}

		// 解析错误字符串中的状态码
		if code := parseStatusCodeFromError(err); code != 0 {
			return code
		}

		return http.StatusInternalServerError
	}
}

// 辅助函数：从错误字符串解析状态码
func parseStatusCodeFromError(err error) int {
	str := err.Error()
	re := regexp.MustCompile(`(\d{3})`)
	matches := re.FindStringSubmatch(str)
	if len(matches) > 1 {
		code, _ := strconv.Atoi(matches[1])
		if code >= 100 && code < 600 {
			return code
		}
	}
	return 0
}

// gRPC 状态码转 HTTP 状态码
func httpStatusCodeFromGRPC(grpcCode codes.Code) int {
	switch grpcCode {
	case codes.OK:
		return http.StatusOK
	case codes.InvalidArgument:
		return http.StatusBadRequest
	case codes.DeadlineExceeded:
		return http.StatusGatewayTimeout
	// ... 其他转换
	default:
		return http.StatusInternalServerError
	}
}
