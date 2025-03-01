package router

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/loongkirin/go-family-finance/internal/api/controller"
	"github.com/loongkirin/go-family-finance/internal/app"
	"github.com/loongkirin/go-family-finance/pkg/http/middleware"
)

type Router struct {
	engine *gin.Engine
}

func NewRouter() *Router {
	engine := gin.Default()
	engine.Use(middleware.RequestIdMiddleware())
	engine.Use(middleware.TraceIdMiddleware())
	engine.Use(middleware.Recovery(app.AppContext.APP_LOGGER))
	engine.Use(middleware.Logger(app.AppContext.APP_LOGGER))
	engine.Use(middleware.Tracing(app.AppContext.APP_TRACER))
	return &Router{engine: engine}
}

func (r *Router) InitRouter() {
	pubGp := r.engine.Group("")
	{
		// 健康监测
		pubGp.GET("/health", func(c *gin.Context) {
			c.JSON(http.StatusOK, "ok")
		})
	}
	v1 := r.engine.Group("/api/v1")
	initAuthorityRouter(v1)
}

func initAuthorityRouter(router *gin.RouterGroup) (R gin.IRoutes) {
	authRouter := router.Group("auth")
	authApi := controller.NewAuthorityController()
	authRouter.GET("captcha", authApi.Captcha)
	return authRouter
}

func (r *Router) Run(addr string) error {
	return r.engine.Run(addr)
}
