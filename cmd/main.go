package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/loongkirin/gdk/logger"
	"github.com/loongkirin/go-family-finance/internal/api/router"
	"github.com/loongkirin/go-family-finance/internal/app"
	"github.com/loongkirin/go-family-finance/internal/migrations"
)

func main() {
	// quertFilter := request.NewQueryFilter("name", []interface{}{"test"}, request.EQ)
	// queryWhere := request.NewQueryWhere([]*request.QueryFilter{quertFilter}, request.AND)
	// queryOrderBy := request.NewQueryOrderBy("created_at", true)
	// query := request.NewQuery([]*request.QueryWhere{queryWhere}, 100, 1, []*request.QueryOrderBy{queryOrderBy})

	// dbQuery := &database.DbQuery{}
	// // err := copier.Copy(dbQuery, query)
	// err := copier.CopyWithOption(dbQuery, query, copier.Option{DeepCopy: true})
	// if err != nil {
	// 	fmt.Println(err)
	// }

	// fmt.Println(dbQuery)

	app.InitAppContext()

	app.AppContext.APP_LOGGER.Info("migration database...", logger.Fields{})

	migrations.Migrate(app.AppContext.APP_DbContext.GetMasterDb())

	app.AppContext.APP_LOGGER.Info("Init router...", logger.Fields{})
	// 初始化路由
	apiRouter := router.NewRouter()
	apiRouter.InitRouter()

	app.AppContext.APP_LOGGER.Info("Start server...", logger.Fields{})
	// 启动服务器
	go func() {
		if err := apiRouter.Run(fmt.Sprintf(":%d", app.AppContext.APP_CONFIG.ServerConfig.Port)); err != nil {
			app.AppContext.APP_LOGGER.Error("Failed to start server", logger.Fields{"error": err})
		}
	}()

	// 优雅关机
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	app.AppContext.APP_LOGGER.Info("Shutting down server...", logger.Fields{})
}
