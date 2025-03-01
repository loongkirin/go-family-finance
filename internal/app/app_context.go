package app

import (
	"fmt"

	"github.com/fsnotify/fsnotify"
	"github.com/loongkirin/go-family-finance/pkg/cache/redis"
	"github.com/loongkirin/go-family-finance/pkg/database"
	"github.com/loongkirin/go-family-finance/pkg/database/postgres"
	"github.com/loongkirin/go-family-finance/pkg/logger"
	"github.com/loongkirin/go-family-finance/pkg/telemetry"
	"github.com/spf13/viper"
	"go.opentelemetry.io/otel/trace"
	"golang.org/x/sync/singleflight"
)

type appContext struct {
	APP_CONFIG                 AppConfig
	APP_VP                     *viper.Viper
	APP_REDIS                  *redis.RedisClient
	APP_DbContext              database.DbContext
	APP_Concurrency_Controller *singleflight.Group
	APP_LOGGER                 logger.Logger
	APP_TRACER                 trace.Tracer
}

var AppContext appContext

func InitAppContext() {
	AppContext = appContext{
		APP_Concurrency_Controller: &singleflight.Group{},
	}
	AppContext.initViper()
	AppContext.initRedis()
	AppContext.initDbContext()
	AppContext.initLogger()
	AppContext.initTracer()
}

func (ctx *appContext) initViper() {
	vp := viper.New()
	vp.SetConfigName("config")
	vp.SetConfigType("yaml")
	vp.AddConfigPath("./")
	err := vp.ReadInConfig()
	if err != nil {
		panic(fmt.Errorf("Fatal error when read config file: %s", err))
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
		panic(fmt.Errorf("Fatal error when unmarshal config file: %s", err))
	}

	ctx.APP_VP = vp
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
	ctx.APP_DbContext = createDbContext(ctx.APP_CONFIG.DbConfig)
}

func createDbContext(cfg database.DbConfig) database.DbContext {
	var dbcontext database.DbContext
	if cfg.DbType == "postgres" {
		postgresDbCtx, err := postgres.NewPostgresDbContext(&cfg)
		if err != nil {
			fmt.Println(err)
		}
		dbcontext = postgresDbCtx
	}

	return dbcontext
}

func (ctx *appContext) initLogger() {
	applogger, err := logger.NewLogger(&ctx.APP_CONFIG.LoggerConfig)
	if err != nil {
		fmt.Println(err)
	}
	ctx.APP_LOGGER = applogger
}

func (ctx *appContext) initTracer() {
	tp, err := telemetry.InitTracer(ctx.APP_CONFIG.TelemetryConfig)
	if err != nil {
		fmt.Println(err)
	} else {
		ctx.APP_TRACER = tp.Tracer("go-family-finance")
	}
}
