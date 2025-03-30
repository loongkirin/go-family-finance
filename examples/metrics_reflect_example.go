package main

import (
	"context"
	"log"
	"time"

	"LoongKirin/go-family-finance/pkg/telemetry"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

func main() {
	// 初始化 OpenTelemetry
	ctx := context.Background()
	meter := otel.GetMeterProvider().Meter("example")
	dm := telemetry.NewDynamicMeter(meter)

	// 创建自定义计数器
	err := dm.CreateMetric(telemetry.MetricDefinition{
		Name:        "custom.counter",
		Description: "A custom counter metric",
		Unit:        "1",
		Type:        telemetry.MetricTypeCounter,
		Attributes: []attribute.KeyValue{
			attribute.String("service", "example"),
		},
	})
	if err != nil {
		log.Fatal(err)
	}

	// 创建自定义直方图
	err = dm.CreateMetric(telemetry.MetricDefinition{
		Name:        "custom.histogram",
		Description: "A custom histogram metric",
		Unit:        "ms",
		Type:        telemetry.MetricTypeHistogram,
		Attributes: []attribute.KeyValue{
			attribute.String("service", "example"),
		},
	})
	if err != nil {
		log.Fatal(err)
	}

	// 创建自定义仪表盘
	err = dm.CreateMetric(telemetry.MetricDefinition{
		Name:        "custom.gauge",
		Description: "A custom gauge metric",
		Unit:        "%",
		Type:        telemetry.MetricTypeGauge,
		Attributes: []attribute.KeyValue{
			attribute.String("service", "example"),
		},
	})
	if err != nil {
		log.Fatal(err)
	}

	// 记录计数器值
	err = dm.RecordMetric(ctx, "custom.counter", telemetry.MetricValue{
		Value: int64(1),
		Attributes: []attribute.KeyValue{
			attribute.String("label", "value"),
		},
		Timestamp: time.Now(),
	})
	if err != nil {
		log.Fatal(err)
	}

	// 记录直方图值
	err = dm.RecordMetric(ctx, "custom.histogram", telemetry.MetricValue{
		Value: float64(100.5),
		Attributes: []attribute.KeyValue{
			attribute.String("label", "value"),
		},
		Timestamp: time.Now(),
	})
	if err != nil {
		log.Fatal(err)
	}

	// 记录仪表盘值
	err = dm.RecordMetric(ctx, "custom.gauge", telemetry.MetricValue{
		Value: float64(75.5),
		Attributes: []attribute.KeyValue{
			attribute.String("label", "value"),
		},
		Timestamp: time.Now(),
	})
	if err != nil {
		log.Fatal(err)
	}

	// 模拟持续记录
	for i := 0; i < 10; i++ {
		// 记录计数器
		err = dm.RecordMetric(ctx, "custom.counter", telemetry.MetricValue{
			Value: int64(i + 1),
			Attributes: []attribute.KeyValue{
				attribute.String("iteration", "loop"),
			},
			Timestamp: time.Now(),
		})
		if err != nil {
			log.Printf("Error recording counter: %v", err)
		}

		// 记录直方图
		err = dm.RecordMetric(ctx, "custom.histogram", telemetry.MetricValue{
			Value: float64(i) * 10.5,
			Attributes: []attribute.KeyValue{
				attribute.String("iteration", "loop"),
			},
			Timestamp: time.Now(),
		})
		if err != nil {
			log.Printf("Error recording histogram: %v", err)
		}

		// 记录仪表盘
		err = dm.RecordMetric(ctx, "custom.gauge", telemetry.MetricValue{
			Value: float64(i) * 10.0,
			Attributes: []attribute.KeyValue{
				attribute.String("iteration", "loop"),
			},
			Timestamp: time.Now(),
		})
		if err != nil {
			log.Printf("Error recording gauge: %v", err)
		}

		time.Sleep(time.Second)
	}
}
