package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/loongkirin/go-family-finance/internal/api/router"
	"github.com/loongkirin/go-family-finance/internal/app"
	"github.com/loongkirin/go-family-finance/pkg/logger"
)

func main() {
	app.InitAppContext()

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
