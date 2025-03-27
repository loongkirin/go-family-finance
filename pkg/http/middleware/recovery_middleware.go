package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/loongkirin/go-family-finance/pkg/http/response"
	pkglogger "github.com/loongkirin/go-family-finance/pkg/logger"
)

func Recovery(logger pkglogger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				logger.Error("Panic recovered", pkglogger.Fields{
					"error":     err,
					"traceId":   GetTraceID(c),
					"requestId": GetRequestId(c),
					"method":    c.Request.Method,
					"path":      c.Request.URL.Path,
					"message":   err,
				})

				c.AbortWithStatusJSON(500, response.NewResponse(response.ERROR, "Internal server error"))
			}
		}()
		c.Next()
	}
}
