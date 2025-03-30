package main

import (
	"LoongKirin/go-family-finance/pkg/middleware"
	"LoongKirin/go-family-finance/pkg/response"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	// 添加追踪中间件
	r.Use(middleware.TraceMiddleware())

	// 添加追踪ID中间件
	r.Use(func(c *gin.Context) {
		c.Set("trace_id", "example-trace-id")
		c.Next()
	})

	// GET 请求示例
	r.GET("/users", func(c *gin.Context) {
		response.Success(c, map[string]interface{}{
			"users": []map[string]interface{}{
				{"id": 1, "name": "John"},
				{"id": 2, "name": "Jane"},
			},
		})
	})

	// POST 请求示例
	r.POST("/users", func(c *gin.Context) {
		var user struct {
			Name  string `json:"name"`
			Email string `json:"email"`
		}

		if err := c.ShouldBindJSON(&user); err != nil {
			response.BadRequest(c, "无效的请求参数")
			return
		}

		response.Created(c, map[string]interface{}{
			"id":    3,
			"name":  user.Name,
			"email": user.Email,
		})
	})

	// 错误处理示例
	r.GET("/error", func(c *gin.Context) {
		response.Error(c, 1001, "业务错误")
	})

	// 启动服务器
	r.Run(":8080")
}
