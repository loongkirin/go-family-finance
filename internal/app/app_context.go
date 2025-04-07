package app

import (
	"context"
	"fmt"

	"github.com/fsnotify/fsnotify"
	"github.com/loongkirin/gdk/cache/redis"
	gdkgorm "github.com/loongkirin/gdk/database/gorm"
	"github.com/loongkirin/gdk/logger"
	"github.com/loongkirin/gdk/telemetry"
	"github.com/spf13/viper"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/trace"
	"golang.org/x/sync/singleflight"
)

type appContext struct {
	APP_CONFIG                 AppConfig
	APP_VP                     *viper.Viper
	APP_REDIS                  *redis.RedisClient
	APP_DbContext              gdkgorm.DbContext
	APP_Concurrency_Controller *singleflight.Group
	APP_LOGGER                 logger.Logger
	APP_TRACER                 trace.Tracer
	APP_METRICS                metric.Meter
}

var AppContext appContext

func InitAppContext() {
	AppContext = appContext{
		APP_Concurrency_Controller: &singleflight.Group{},
	}
	AppContext.initViper()
	AppContext.initLogger()
	AppContext.initTracer()
	AppContext.initMetrics()
	AppContext.initRedis()
	AppContext.initDbContext()
}

func (ctx *appContext) initViper() {
	vp := viper.New()
	vp.SetConfigName("config")
	vp.SetConfigType("yaml")
	vp.AddConfigPath("./")
	err := vp.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("fatal error when read config file: %s", err))
	}

	vp.WatchConfig()

	vp.OnConfigChange(func(e fsnotify.Event) {
		fmt.Println("config file changed:", e.Name)
		if err := vp.Unmarshal(&ctx.APP_CONFIG); err != nil {
			fmt.Println(err)
		}
	})

	if err := vp.Unmarshal(&ctx.APP_CONFIG); err != nil {
		fmt.Println(err)
		panic(fmt.Errorf("fatal error when unmarshal config file: %s", err))
	}

	ctx.APP_VP = vp
}

func (ctx *appContext) initLogger() {
	applogger, err := logger.NewLogger(&ctx.APP_CONFIG.LoggerConfig)
	if err != nil {
		fmt.Println(err)
	}
	ctx.APP_LOGGER = applogger
}

func (ctx *appContext) initTracer() {
	tp, err := telemetry.InitTracer(context.Background(), ctx.APP_CONFIG.TelemetryConfig)
	if err != nil {
		fmt.Println(err)
	} else {
		ctx.APP_TRACER = tp.Tracer(ctx.APP_CONFIG.TelemetryConfig.ServiceName)
		// defer telemetry.ShutdownTracer(context.Background())
	}
}

func (ctx *appContext) initMetrics() {
	metrics, err := telemetry.InitMetrics(context.Background(), ctx.APP_CONFIG.TelemetryConfig)
	if err != nil {
		fmt.Println(err)
	} else {
		ctx.APP_METRICS = metrics.Meter(ctx.APP_CONFIG.TelemetryConfig.ServiceName)
		// defer telemetry.ShutdownMetrics(context.Background())
	}
}

func (ctx *appContext) initRedis() {
	redisCfg := ctx.APP_CONFIG.RedisConfig
	client, err := redis.NewRedisClient(&redisCfg)
	if err != nil {
		fmt.Println(err)
	} else {
		ctx.APP_REDIS = client
	}
}

func (ctx *appContext) initDbContext() {
	dbContext := gdkgorm.CreateDbContext(&ctx.APP_CONFIG.DbConfig)
	ctx.APP_DbContext = dbContext
}
