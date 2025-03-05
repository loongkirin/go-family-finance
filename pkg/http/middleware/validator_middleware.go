package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	pkglogger "github.com/loongkirin/go-family-finance/pkg/logger"
)

// 自定义验证错误消息
func customErrorMessage(e validator.FieldError) string {
	switch e.Tag() {
	case "required":
		return "字段 " + e.Field() + " 是必填的"
	case "min":
		return "字段 " + e.Field() + " 的最小值是 " + e.Param()
	case "max":
		return "字段 " + e.Field() + " 的最大值是 " + e.Param()
	case "email":
		return "字段 " + e.Field() + " 必须是有效的邮箱地址"
	}
	return e.Error()
}

// 验证中间件
func ValidateRequest[T any](logger pkglogger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req T

		// 绑定并验证请求参数
		if err := c.ShouldBindJSON(&req); err != nil {
			// 处理验证错误
			if validationErrors, ok := err.(validator.ValidationErrors); ok {
				for _, e := range validationErrors {
					errMsg := customErrorMessage(e)
					logger.Error("无效的请求参数", pkglogger.Fields{
						"error":     e,
						"traceId":   GetTraceID(c),
						"requestId": GetRequestId(c),
						"method":    c.Request.Method,
						"path":      c.Request.URL.Path,
						"message":   errMsg,
					})
					c.JSON(http.StatusBadRequest, gin.H{
						"error": errMsg,
					})
					c.Abort()
					return
				}
			} else {
				logger.Error("无效的请求参数", pkglogger.Fields{
					"error":     err,
					"traceId":   GetTraceID(c),
					"requestId": GetRequestId(c),
					"method":    c.Request.Method,
					"path":      c.Request.URL.Path,
					"message":   err.Error(),
				})
			}
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "无效的请求参数",
			})
			c.Abort()
			return
		}

		// 将验证后的请求参数存储到上下文中
		c.Set("validatedRequest", req)

		// 继续处理请求
		c.Next()
	}
}
