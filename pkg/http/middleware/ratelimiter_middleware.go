package middleware

import (
	"fmt"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/ratelimit"
)

// RateLimiter 基于 source 的限流器
type SourceRateLimiter struct {
	sources map[string]ratelimit.Limiter
	mu      *sync.RWMutex
}

// NewRateLimiter 创建一个新的限流器
func NewSourceRateLimiter() *SourceRateLimiter {
	return &SourceRateLimiter{
		sources: make(map[string]ratelimit.Limiter),
		mu:      &sync.RWMutex{},
	}
}

// GetLimiter 获取 source 对应的限流器
func (i *SourceRateLimiter) GetLimiter(source string, limit int) ratelimit.Limiter {
	i.mu.Lock()
	defer i.mu.Unlock()

	limiter, exists := i.sources[source]
	if !exists {
		limiter = ratelimit.New(limit, ratelimit.Per(time.Second*10)) // 每个 source 的限流速率 limit 次/秒
		i.sources[source] = limiter
	}

	return limiter
}

// RateLimitMiddleware 限流中间件
func RateLimiter(limiter *SourceRateLimiter, limit int) gin.HandlerFunc {
	return func(c *gin.Context) {
		source := fmt.Sprintf("%s:%s:%s", c.ClientIP(), c.Request.Method, c.Request.URL.Path)
		limiter := limiter.GetLimiter(source, limit)
		limiter.Take() // 阻塞直到允许请求
		c.Next()
	}
}
