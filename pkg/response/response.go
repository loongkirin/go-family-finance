package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Response 统一响应结构
type Response struct {
	Code    int         `json:"code"`     // 业务码
	Message string      `json:"message"`  // 提示信息
	Data    interface{} `json:"data"`     // 数据
	TraceID string      `json:"trace_id"` // 追踪ID
}

// Success 成功响应
func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code:    0,
		Message: "success",
		Data:    data,
		TraceID: c.GetString("trace_id"),
	})
}

// Error 错误响应
func Error(c *gin.Context, code int, message string) {
	c.JSON(http.StatusOK, Response{
		Code:    code,
		Message: message,
		Data:    nil,
		TraceID: c.GetString("trace_id"),
	})
}

// BadRequest 请求参数错误
func BadRequest(c *gin.Context, message string) {
	c.JSON(http.StatusBadRequest, Response{
		Code:    400,
		Message: message,
		Data:    nil,
		TraceID: c.GetString("trace_id"),
	})
}

// Unauthorized 未授权
func Unauthorized(c *gin.Context, message string) {
	c.JSON(http.StatusUnauthorized, Response{
		Code:    401,
		Message: message,
		Data:    nil,
		TraceID: c.GetString("trace_id"),
	})
}

// Forbidden 禁止访问
func Forbidden(c *gin.Context, message string) {
	c.JSON(http.StatusForbidden, Response{
		Code:    403,
		Message: message,
		Data:    nil,
		TraceID: c.GetString("trace_id"),
	})
}

// NotFound 资源不存在
func NotFound(c *gin.Context, message string) {
	c.JSON(http.StatusNotFound, Response{
		Code:    404,
		Message: message,
		Data:    nil,
		TraceID: c.GetString("trace_id"),
	})
}

// InternalServerError 服务器内部错误
func InternalServerError(c *gin.Context, message string) {
	c.JSON(http.StatusInternalServerError, Response{
		Code:    500,
		Message: message,
		Data:    nil,
		TraceID: c.GetString("trace_id"),
	})
}

// ServiceUnavailable 服务不可用
func ServiceUnavailable(c *gin.Context, message string) {
	c.JSON(http.StatusServiceUnavailable, Response{
		Code:    503,
		Message: message,
		Data:    nil,
		TraceID: c.GetString("trace_id"),
	})
}

// GatewayTimeout 网关超时
func GatewayTimeout(c *gin.Context, message string) {
	c.JSON(http.StatusGatewayTimeout, Response{
		Code:    504,
		Message: message,
		Data:    nil,
		TraceID: c.GetString("trace_id"),
	})
}

// TooManyRequests 请求过多
func TooManyRequests(c *gin.Context, message string) {
	c.JSON(http.StatusTooManyRequests, Response{
		Code:    429,
		Message: message,
		Data:    nil,
		TraceID: c.GetString("trace_id"),
	})
}

// Conflict 资源冲突
func Conflict(c *gin.Context, message string) {
	c.JSON(http.StatusConflict, Response{
		Code:    409,
		Message: message,
		Data:    nil,
		TraceID: c.GetString("trace_id"),
	})
}

// Created 创建成功
func Created(c *gin.Context, data interface{}) {
	c.JSON(http.StatusCreated, Response{
		Code:    0,
		Message: "created",
		Data:    data,
		TraceID: c.GetString("trace_id"),
	})
}

// NoContent 无内容
func NoContent(c *gin.Context) {
	c.Status(http.StatusNoContent)
}

// Accepted 已接受
func Accepted(c *gin.Context, data interface{}) {
	c.JSON(http.StatusAccepted, Response{
		Code:    0,
		Message: "accepted",
		Data:    data,
		TraceID: c.GetString("trace_id"),
	})
}

// PartialContent 部分内容
func PartialContent(c *gin.Context, data interface{}) {
	c.JSON(http.StatusPartialContent, Response{
		Code:    0,
		Message: "partial content",
		Data:    data,
		TraceID: c.GetString("trace_id"),
	})
}

// ResetContent 重置内容
func ResetContent(c *gin.Context) {
	c.Status(http.StatusResetContent)
}

// MultiStatus 多状态
func MultiStatus(c *gin.Context, data interface{}) {
	c.JSON(http.StatusMultiStatus, Response{
		Code:    0,
		Message: "multi status",
		Data:    data,
		TraceID: c.GetString("trace_id"),
	})
}

// AlreadyReported 已报告
func AlreadyReported(c *gin.Context, data interface{}) {
	c.JSON(http.StatusAlreadyReported, Response{
		Code:    0,
		Message: "already reported",
		Data:    data,
		TraceID: c.GetString("trace_id"),
	})
}

// IMUsed IM已使用
func IMUsed(c *gin.Context, data interface{}) {
	c.JSON(http.StatusIMUsed, Response{
		Code:    0,
		Message: "im used",
		Data:    data,
		TraceID: c.GetString("trace_id"),
	})
}
