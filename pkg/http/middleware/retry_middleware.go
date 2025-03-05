package middleware

import (
	"fmt"
	"net/http"
	"time"

	"github.com/avast/retry-go"
	"github.com/gin-gonic/gin"
	pkglogger "github.com/loongkirin/go-family-finance/pkg/logger"
)

func Retry(logger pkglogger.Logger, maxRetries uint, retryDelay time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		err := retry.Do(
			func() error {
				// 执行下一个处理程序
				c.Next()

				// 如果没有错误，直接返回
				if c.Writer.Status() == http.StatusOK || c.Writer.Status() == http.StatusNotFound {
					return nil
				}

				// 返回错误
				return fmt.Errorf("request failed with status %d", c.Writer.Status())
			},
			retry.Attempts(maxRetries), // 重试次数
			retry.Delay(retryDelay),    // 重试间隔
			retry.OnRetry(func(n uint, err error) {
				errMsg := fmt.Sprintf("Retry #%d: %s\n", n, err)
				logger.Error("request failed with retry", pkglogger.Fields{
					"error":     err,
					"traceId":   GetTraceID(c),
					"requestId": GetRequestId(c),
					"method":    c.Request.Method,
					"path":      c.Request.URL.Path,
					"message":   errMsg,
				})
			}),
		)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Request failed after retries",
			})
		}
	}
}
