package middleware

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/loongkirin/go-family-finance/pkg/http/response"
	"golang.org/x/time/rate"
)

// RateLimiter 基于 source 的限流器
type RRateLimiter struct {
	sources map[string]*rate.Limiter
	mu      *sync.RWMutex
	r       rate.Limit
	b       int
}

// NewRateLimiter 创建一个新的限流器
func NewRRateLimiter(r rate.Limit, b int) *RRateLimiter {
	return &RRateLimiter{
		sources: make(map[string]*rate.Limiter),
		mu:      &sync.RWMutex{},
		r:       r,
		b:       b,
	}
}

// 添加 Source 到限流器
func (i *RRateLimiter) AddSource(source string) *rate.Limiter {
	i.mu.Lock()
	defer i.mu.Unlock()

	limiter := rate.NewLimiter(i.r, i.b)
	i.sources[source] = limiter
	return limiter
}

// GetLimiter 获取 IP 对应的限流器
func (i *RRateLimiter) GetLimiter(source string) *rate.Limiter {
	i.mu.Lock()
	limiter, exists := i.sources[source]

	if !exists {
		i.mu.Unlock()
		return i.AddSource(source)
	}

	i.mu.Unlock()
	return limiter
}

// RateLimitMiddleware 限流中间件
func RequestRateLimiter(limiter *RRateLimiter) gin.HandlerFunc {
	return func(c *gin.Context) {
		source := fmt.Sprintf("%s:%s:%s", c.ClientIP(), c.Request.Method, c.Request.URL.Path)
		limiter := limiter.GetLimiter(source)
		if !limiter.Allow() {
			c.JSON(http.StatusTooManyRequests, response.NewResponse(response.ERROR, "Too many requests"))
			c.Abort()
			return
		}
		c.Next()
	}
}
