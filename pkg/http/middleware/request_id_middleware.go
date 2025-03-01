package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/loongkirin/go-family-finance/pkg/util"
)

func RequestIdMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		requestId := GetRequestId(ctx)
		if len(requestId) == 0 {
			requestId = util.GenerateId()
		}
		SetRequestId(ctx, requestId)
		ctx.Next()
	}
}
