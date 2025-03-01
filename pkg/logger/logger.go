package logger

import (
	"context"
	"fmt"

	"go.opentelemetry.io/otel/trace"
)

// Level 定义日志级别
type Level int

const (
	DebugLevel Level = iota
	InfoLevel
	WarnLevel
	ErrorLevel
	FatalLevel
)

// Fields 定义日志字段类型
type Fields map[string]interface{}

// Logger 定义通用日志接口
type Logger interface {
	// 基础日志方法
	Debug(msg string, fields ...Fields)
	Info(msg string, fields ...Fields)
	Warn(msg string, fields ...Fields)
	Error(msg string, fields ...Fields)
	Fatal(msg string, fields ...Fields)

	// 带上下文的日志方法
	WithContext(ctx context.Context) Logger
	WithFields(fields Fields) Logger

	// 获取原始日志实例
	GetLogger() interface{}
}

// TraceIDKey 是上下文中存储 TraceID 的键
type TraceIDKey struct{}

// getTraceID 从上下文中获取追踪ID
func getTraceID(ctx context.Context) string {
	if ctx == nil {
		return ""
	}

	// 首先尝试从上下文中获取自定义的 TraceID
	if traceID, ok := ctx.Value(TraceIDKey{}).(string); ok && traceID != "" {
		return traceID
	}

	// 尝试获取 OpenTelemetry TraceID
	if spanCtx := trace.SpanContextFromContext(ctx); spanCtx.IsValid() {
		return spanCtx.TraceID().String()
	}

	return ""
}

// WithTraceID 向上下文中添加 TraceID
func WithTraceID(ctx context.Context, traceID string) context.Context {
	return context.WithValue(ctx, TraceIDKey{}, traceID)
}

type LoggerType string

const (
	ZerologType LoggerType = "zerolog"
	SlogType    LoggerType = "slog"
	ZapType     LoggerType = "zap"
)

// LoggerConfig 定义Logger配置
type LoggerConfig struct {
	// 日志类型
	LoggerType LoggerType `mapstructure:"logger_type" json:"logger_type" yaml:"logger_type"`
	// 日志配置
	LogConfig LogConfig `mapstructure:"log_config" json:"log_config" yaml:"log_config"`
}

// LogConfig 定义日志配置
type LogConfig struct {
	// 日志级别
	Level string `mapstructure:"level" json:"level" yaml:"level"`
	// 日志格式 (json/console)
	Format string `mapstructure:"format" json:"format" yaml:"format"`
	// 日志文件路径
	FilePath string `mapstructure:"file_path" json:"file_path" yaml:"file_path"`
	// 是否输出到控制台
	EnableConsole bool `mapstructure:"enable_console" json:"enable_console" yaml:"enable_console"`
	// 是否启用调用者信息
	EnableCaller bool `mapstructure:"enable_caller" json:"enable_caller" yaml:"enable_caller"`
	// 日志轮转配置
	MaxSize    int  `mapstructure:"max_size" json:"max_size" yaml:"max_size"`          // MB
	MaxBackups int  `mapstructure:"max_backups" json:"max_backups" yaml:"max_backups"` // 文件个数
	MaxAge     int  `mapstructure:"max_age" json:"max_age" yaml:"max_age"`             // 天数
	Compress   bool `mapstructure:"compress" json:"compress" yaml:"compress"`          // 是否压缩
}

// NewLogger 创建指定类型的日志实例
func NewLogger(config *LoggerConfig) (Logger, error) {
	switch config.LoggerType {
	case ZerologType:
		return NewZerologLogger(&config.LogConfig)
	case SlogType:
		return NewSlogLogger(&config.LogConfig)
	case ZapType:
		return NewZapLogger(&config.LogConfig)
	default:
		return nil, fmt.Errorf("unsupported logger type: %s", config.LoggerType)
	}
}
