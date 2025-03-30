package main

import (
	"LoongKirin/go-family-finance/pkg/response"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	// 添加追踪ID中间件
	r.Use(func(c *gin.Context) {
		c.Set("trace_id", "example-trace-id")
		c.Next()
	})

	// 成功响应示例
	r.GET("/success", func(c *gin.Context) {
		data := map[string]interface{}{
			"name": "John",
			"age":  30,
		}
		response.Success(c, data)
	})

	// 错误响应示例
	r.GET("/error", func(c *gin.Context) {
		response.Error(c, 1001, "业务错误")
	})

	// 参数错误示例
	r.GET("/bad-request", func(c *gin.Context) {
		response.BadRequest(c, "请求参数错误")
	})

	// 未授权示例
	r.GET("/unauthorized", func(c *gin.Context) {
		response.Unauthorized(c, "未授权访问")
	})

	// 禁止访问示例
	r.GET("/forbidden", func(c *gin.Context) {
		response.Forbidden(c, "禁止访问")
	})

	// 资源不存在示例
	r.GET("/not-found", func(c *gin.Context) {
		response.NotFound(c, "资源不存在")
	})

	// 服务器错误示例
	r.GET("/internal-error", func(c *gin.Context) {
		response.InternalServerError(c, "服务器内部错误")
	})

	// 创建成功示例
	r.POST("/create", func(c *gin.Context) {
		data := map[string]interface{}{
			"id":   1,
			"name": "New Item",
		}
		response.Created(c, data)
	})

	// 无内容示例
	r.DELETE("/delete", func(c *gin.Context) {
		response.NoContent(c)
	})

	// 已接受示例
	r.POST("/async", func(c *gin.Context) {
		data := map[string]interface{}{
			"task_id": "123",
			"status":  "processing",
		}
		response.Accepted(c, data)
	})

	// 部分内容示例
	r.GET("/partial", func(c *gin.Context) {
		data := map[string]interface{}{
			"items": []string{"item1", "item2"},
			"total": 5,
		}
		response.PartialContent(c, data)
	})

	// 重置内容示例
	r.POST("/reset", func(c *gin.Context) {
		response.ResetContent(c)
	})

	// 多状态示例
	r.GET("/multi", func(c *gin.Context) {
		data := map[string]interface{}{
			"results": []map[string]interface{}{
				{"id": 1, "status": "success"},
				{"id": 2, "status": "failed"},
			},
		}
		response.MultiStatus(c, data)
	})

	// 已报告示例
	r.GET("/reported", func(c *gin.Context) {
		data := map[string]interface{}{
			"report_id": "123",
			"status":    "reported",
		}
		response.AlreadyReported(c, data)
	})

	// IM已使用示例
	r.GET("/im-used", func(c *gin.Context) {
		data := map[string]interface{}{
			"message": "IM已使用",
		}
		response.IMUsed(c, data)
	})

	// 启动服务器
	r.Run(":8080")
}
