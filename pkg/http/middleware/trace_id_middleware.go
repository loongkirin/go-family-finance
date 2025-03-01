package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/loongkirin/go-family-finance/pkg/util"
)

func TraceIdMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		traceId := GetTraceID(ctx)
		if len(traceId) == 0 {
			traceId = util.GenerateId()
		}
		SetTraceID(ctx, traceId)
		ctx.Next()
	}
}
