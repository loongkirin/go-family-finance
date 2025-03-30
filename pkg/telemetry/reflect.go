package telemetry

import (
	"context"
	"fmt"
	"reflect"
	"time"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
)

// MetricType 指标类型
type MetricType string

const (
	MetricTypeCounter       MetricType = "counter"
	MetricTypeUpDownCounter MetricType = "updowncounter"
	MetricTypeHistogram     MetricType = "histogram"
	MetricTypeGauge         MetricType = "gauge"
)

// MetricDefinition 指标定义
type MetricDefinition struct {
	Name        string
	Description string
	Unit        string
	Type        MetricType
	Attributes  []attribute.KeyValue
}

// MetricValue 指标值
type MetricValue struct {
	Value      interface{}
	Attributes []attribute.KeyValue
	Timestamp  time.Time
}

// DynamicMeter 动态指标记录器
type DynamicMeter struct {
	meter   metric.Meter
	metrics map[string]interface{}
}

// NewDynamicMeter 创建动态指标记录器
func NewDynamicMeter(meter metric.Meter) *DynamicMeter {
	return &DynamicMeter{
		meter:   meter,
		metrics: make(map[string]interface{}),
	}
}

// CreateMetric 创建指标
func (dm *DynamicMeter) CreateMetric(def MetricDefinition) error {
	var err error
	switch def.Type {
	case MetricTypeCounter:
		dm.metrics[def.Name], err = dm.meter.Int64Counter(
			def.Name,
			metric.WithDescription(def.Description),
			metric.WithUnit(def.Unit),
		)
	case MetricTypeUpDownCounter:
		dm.metrics[def.Name], err = dm.meter.Int64UpDownCounter(
			def.Name,
			metric.WithDescription(def.Description),
			metric.WithUnit(def.Unit),
		)
	case MetricTypeHistogram:
		dm.metrics[def.Name], err = dm.meter.Float64Histogram(
			def.Name,
			metric.WithDescription(def.Description),
			metric.WithUnit(def.Unit),
		)
	case MetricTypeGauge:
		dm.metrics[def.Name], err = dm.meter.Float64Gauge(
			def.Name,
			metric.WithDescription(def.Description),
			metric.WithUnit(def.Unit),
		)
	default:
		return fmt.Errorf("unsupported metric type: %s", def.Type)
	}
	return err
}

// RecordMetric 记录指标值
func (dm *DynamicMeter) RecordMetric(ctx context.Context, name string, value MetricValue) error {
	metric, exists := dm.metrics[name]
	if !exists {
		return fmt.Errorf("metric not found: %s", name)
	}

	attrs := append(value.Attributes, attribute.String("timestamp", value.Timestamp.Format(time.RFC3339)))

	switch m := metric.(type) {
	case metric.Int64Counter:
		if v, ok := value.Value.(int64); ok {
			m.Add(ctx, v, metric.WithAttributes(attrs...))
		} else {
			return fmt.Errorf("invalid value type for counter: %v", reflect.TypeOf(value.Value))
		}
	case metric.Int64UpDownCounter:
		if v, ok := value.Value.(int64); ok {
			m.Add(ctx, v, metric.WithAttributes(attrs...))
		} else {
			return fmt.Errorf("invalid value type for updowncounter: %v", reflect.TypeOf(value.Value))
		}
	case metric.Float64Histogram:
		if v, ok := value.Value.(float64); ok {
			m.Record(ctx, v, metric.WithAttributes(attrs...))
		} else {
			return fmt.Errorf("invalid value type for histogram: %v", reflect.TypeOf(value.Value))
		}
	case metric.Float64Gauge:
		if v, ok := value.Value.(float64); ok {
			m.Set(ctx, v, metric.WithAttributes(attrs...))
		} else {
			return fmt.Errorf("invalid value type for gauge: %v", reflect.TypeOf(value.Value))
		}
	default:
		return fmt.Errorf("unsupported metric type: %v", reflect.TypeOf(metric))
	}
	return nil
}

// Example usage:
/*
func ExampleDynamicMeter() {
	ctx := context.Background()
	meter := otel.GetMeterProvider().Meter("example")
	dm := NewDynamicMeter(meter)

	// 创建自定义指标
	err := dm.CreateMetric(MetricDefinition{
		Name:        "custom.counter",
		Description: "A custom counter metric",
		Unit:        "1",
		Type:        MetricTypeCounter,
		Attributes: []attribute.KeyValue{
			attribute.String("service", "example"),
		},
	})
	if err != nil {
		log.Fatal(err)
	}

	// 记录指标值
	err = dm.RecordMetric(ctx, "custom.counter", MetricValue{
		Value:      int64(1),
		Attributes: []attribute.KeyValue{
			attribute.String("label", "value"),
		},
		Timestamp: time.Now(),
	})
	if err != nil {
		log.Fatal(err)
	}
}
*/
