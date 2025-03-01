package middleware

import (
	"github.com/gin-gonic/gin"
)

const (
	traceHeaderKey     = "x-trace-id"
	requestIdHeaderKey = "x-request-id"
)

func GetTraceID(c *gin.Context) string {
	traceId := c.GetHeader(traceHeaderKey)
	return traceId
}

func SetTraceID(c *gin.Context, traceId string) {
	c.Header(traceHeaderKey, traceId)
}

func GetRequestId(c *gin.Context) string {
	requestId := c.GetHeader(requestIdHeaderKey)
	return requestId
}

func SetRequestId(c *gin.Context, requestId string) {
	c.Header(requestIdHeaderKey, requestId)
}
