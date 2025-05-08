package router

import (
	"net/http"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/gzip"

	"github.com/gin-gonic/gin"
	"github.com/loongkirin/gdk/net/http/gin/middleware"
	"github.com/loongkirin/gdk/telemetry"
	"github.com/loongkirin/gdk/util"
	"github.com/loongkirin/go-family-finance/internal/api/controller"
	"github.com/loongkirin/go-family-finance/internal/app"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"golang.org/x/time/rate"
)

type Router struct {
	engine *gin.Engine
}

func NewRouter() *Router {
	r := gin.New()

	// gin.SetMode(gin.ReleaseMode)
	// CORS middleware
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// Gzip compression
	r.Use(gzip.Gzip(gzip.DefaultCompression))

	r.Use(middleware.RequestId())
	r.Use(middleware.TraceId())
	r.Use(middleware.Recovery(app.AppContext.APP_LOGGER))
	r.Use(middleware.Logger(app.AppContext.APP_LOGGER))
	r.Use(middleware.Tracing(app.AppContext.APP_TRACER))
	int64Meter := telemetry.NewDynamicMeter[int64](app.AppContext.APP_METRICS)
	r.Use(middleware.BlockRateLimiter(middleware.NewIPBlockRateLimiter(middleware.RateLimiterConfig{
		Limit:   30,
		Timeout: time.Second * 45,
		Meter:   int64Meter,
		Logger:  app.AppContext.APP_LOGGER,
	})))
	r.Use(middleware.BreakRateLimiter(middleware.NewIPBreakRateLimiter(middleware.BreakRateLimiterConfig{
		Limit:  rate.Limit(30),
		Burst:  45,
		Meter:  int64Meter,
		Logger: app.AppContext.APP_LOGGER,
	})))
	dynamicMeter := telemetry.NewDynamicMeter[float64](app.AppContext.APP_METRICS)
	r.Use(middleware.Retry(app.AppContext.APP_LOGGER, 3, time.Second*3))
	r.Use(middleware.Metrics(dynamicMeter))
	r.Use(middleware.Validator(util.MaxLenValidator, util.MinLenValidator, util.PasswordValidator))
	return &Router{engine: r}
}

func (r *Router) InitRouter() {
	pubGp := r.engine.Group("")
	{
		// 健康监测
		pubGp.GET("/health", func(c *gin.Context) {
			c.JSON(http.StatusOK, "ok")
		})
		pubGp.GET("/metrics", gin.WrapH(promhttp.Handler()))
	}
	v1 := r.engine.Group("/api/v1")
	initAuthorityRouter(v1)
}

func initAuthorityRouter(router *gin.RouterGroup) (R gin.IRoutes) {
	authRouter := router.Group("auth")
	// oauthMaker, err := oauth.NewPasetoMaker(app.AppContext.APP_CONFIG.OAuthConfig)
	// if err != nil {
	// 	panic(err)
	// }
	// authRouter.Use(middleware.OAuth(oauthMaker))
	authApi := controller.NewAuthorityController()
	authRouter.GET("captcha", authApi.Captcha)
	authRouter.POST("login", authApi.Login)
	authRouter.POST("register", authApi.Register)
	return authRouter
}

func (r *Router) Run(addr string) error {
	return r.engine.Run(addr)
}
