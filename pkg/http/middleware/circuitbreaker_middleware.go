package middleware

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	breaker "github.com/sony/gobreaker/v2"
)

// CircuitBreakerConfig 断路器配置
type CircuitBreakerConfig struct {
	Name         string
	MaxRequests  uint32        // 半开状态下允许的最大请求数
	Interval     time.Duration // 统计时间窗口
	Timeout      time.Duration // 断路器打开后，多久后尝试半开
	FailureRatio float64       // 触发断路器的失败率阈值
	MinRequests  uint32        // 最小请求数阈值
}

// DefaultCircuitBreakerConfig 默认配置
func DefaultCircuitBreakerConfig() *CircuitBreakerConfig {
	return &CircuitBreakerConfig{
		Name:         "default",
		MaxRequests:  100,
		Interval:     10 * time.Second,
		Timeout:      60 * time.Second,
		FailureRatio: 0.6,
		MinRequests:  10,
	}
}

// CircuitBreakerMiddleware 创建断路器中间件
func CircuitBreaker(config *CircuitBreakerConfig) gin.HandlerFunc {
	if config == nil {
		config = DefaultCircuitBreakerConfig()
	}

	cb := breaker.NewCircuitBreaker[interface{}](breaker.Settings{
		Name:        config.Name,
		MaxRequests: config.MaxRequests,
		Interval:    config.Interval,
		Timeout:     config.Timeout,
		ReadyToTrip: func(counts breaker.Counts) bool {
			failureRatio := float64(counts.TotalFailures) / float64(counts.Requests)
			return counts.Requests >= config.MinRequests && failureRatio >= config.FailureRatio
		},
		OnStateChange: func(name string, from breaker.State, to breaker.State) {
			// 更新 Prometheus 指标
			state := float64(0) // Closed
			if to == breaker.StateHalfOpen {
				state = 1
			} else if to == breaker.StateOpen {
				state = 2
			}
			circuitBreakerState.WithLabelValues(name).Set(state)

			fmt.Printf("Circuit Breaker '%s' state changed from %s to %s\n", name, from, to)
		},
	})

	return func(c *gin.Context) {
		result, err := cb.Execute(func() (interface{}, error) {
			ch := make(chan struct {
				err error
			}, 1)

			go func() {
				c.Next()
				var err error
				if len(c.Errors) > 0 {
					err = c.Errors.Last()
					// 记录失败次数
					circuitBreakerFailures.WithLabelValues(config.Name).Inc()
				}
				ch <- struct{ err error }{err: err}
			}()

			select {
			case result := <-ch:
				return nil, result.err
			case <-time.After(30 * time.Second):
				circuitBreakerFailures.WithLabelValues(config.Name).Inc()
				return nil, fmt.Errorf("request timeout")
			}
		})

		if err != nil {
			if err == breaker.ErrOpenState {
				c.AbortWithStatusJSON(http.StatusServiceUnavailable, gin.H{
					"error": "Service is unavailable",
					"state": "circuit breaker open",
				})
				return
			}

			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		if c.Writer.Written() {
			return
		}

		if result != nil {
			c.JSON(http.StatusOK, result)
		}
	}
}

// CircuitBreakerByPath 为不同路径创建不同的断路器
func CircuitBreakerByPath(configs map[string]*CircuitBreakerConfig) gin.HandlerFunc {
	breakers := make(map[string]*breaker.CircuitBreaker[interface{}])
	defaultConfig := DefaultCircuitBreakerConfig()

	return func(c *gin.Context) {
		path := c.FullPath()
		if path == "" {
			path = c.Request.URL.Path
		}

		cb, exists := breakers[path]
		if !exists {
			config := defaultConfig
			if pathConfig, ok := configs[path]; ok {
				config = pathConfig
			}

			cb = breaker.NewCircuitBreaker[interface{}](breaker.Settings{
				Name:        path,
				MaxRequests: config.MaxRequests,
				Interval:    config.Interval,
				Timeout:     config.Timeout,
				ReadyToTrip: func(counts breaker.Counts) bool {
					failureRatio := float64(counts.TotalFailures) / float64(counts.Requests)
					return counts.Requests >= config.MinRequests && failureRatio >= config.FailureRatio
				},
				OnStateChange: func(name string, from breaker.State, to breaker.State) {
					state := float64(0)
					if to == breaker.StateHalfOpen {
						state = 1
					} else if to == breaker.StateOpen {
						state = 2
					}
					circuitBreakerState.WithLabelValues(name).Set(state)

					fmt.Printf("Circuit Breaker '%s' state changed from %s to %s\n", name, from, to)
				},
			})
			breakers[path] = cb
		}

		result, err := cb.Execute(func() (interface{}, error) {
			ch := make(chan struct {
				err error
			}, 1)

			go func() {
				c.Next()
				var err error
				if len(c.Errors) > 0 {
					err = c.Errors.Last()
					circuitBreakerFailures.WithLabelValues(path).Inc()
				}
				ch <- struct{ err error }{err: err}
			}()

			select {
			case result := <-ch:
				return nil, result.err
			case <-time.After(30 * time.Second):
				circuitBreakerFailures.WithLabelValues(path).Inc()
				return nil, fmt.Errorf("request timeout")
			}
		})

		if err != nil {
			if err == breaker.ErrOpenState {
				c.AbortWithStatusJSON(http.StatusServiceUnavailable, gin.H{
					"error": "Service is unavailable",
					"state": "circuit breaker open",
					"path":  path,
				})
				return
			}

			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"error": err.Error(),
			})
			return
		}

		if c.Writer.Written() {
			return
		}

		if result != nil {
			c.JSON(http.StatusOK, result)
		}
	}
}
