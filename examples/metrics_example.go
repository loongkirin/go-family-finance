package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	"LoongKirin/go-family-finance/pkg/telemetry"
)

func main() {
	// 创建上下文
	ctx := context.Background()

	// 初始化指标
	config := &telemetry.MetricsConfig{
		ServiceName:      "family-finance",
		ServiceVersion:   "1.0.0",
		Environment:      "development",
		ExporterType:     "http",
		ExporterEndpoint: "http://localhost:4318/v1/metrics",
		ExportInterval:   10 * time.Second,
		ExportTimeout:    5 * time.Second,
	}

	if err := telemetry.InitMetrics(ctx, config); err != nil {
		log.Fatalf("Failed to initialize metrics: %v", err)
	}
	defer telemetry.ShutdownMetrics(ctx)

	// 创建信号通道
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// 创建指标记录通道
	done := make(chan struct{})

	// 启动指标记录
	go func() {
		ticker := time.NewTicker(10 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-done:
				return
			case <-ticker.C:
				// 记录系统指标
				var m runtime.MemStats
				runtime.ReadMemStats(&m)

				telemetry.RecordSystemMetrics(
					ctx,
					int64(m.Alloc),
					getCPUUsage(),
					int64(runtime.NumGoroutine()),
				)

				// 记录数据库连接指标
				telemetry.UpdateDBConnections(ctx, 5, 10)

				// 记录 HTTP 请求指标
				telemetry.RecordHTTPRequest(
					ctx,
					"GET",
					"/api/v1/users",
					200,
					100*time.Millisecond,
					1024,
					2048,
				)

				// 记录数据库查询指标
				telemetry.RecordDBQuery(
					ctx,
					"find_user_by_id",
					50*time.Millisecond,
				)

				// 记录用户登录指标
				telemetry.RecordUserLogin(ctx, true)
			}
		}
	}()

	// 等待中断信号
	<-sigChan

	// 停止指标记录
	close(done)

	// 优雅关闭
	log.Println("Shutting down...")
}

// getCPUUsage 获取 CPU 使用率
func getCPUUsage() float64 {
	// 这里应该实现实际的 CPU 使用率计算
	// 这只是一个示例返回值
	return 25.5
}
