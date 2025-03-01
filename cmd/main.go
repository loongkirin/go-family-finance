package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/loongkirin/go-family-finance/internal/app"
	"github.com/loongkirin/go-family-finance/pkg/logger"
)

func main() {
	fmt.Println("Hello, World!")
	app.InitAppContext()

	// 优雅关机
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	app.AppContext.APP_LOGGER.Info("Shutting down server...", logger.Fields{})
}
