package logger

import (
	"context"
	"io"
	"log/slog"
	"os"
	"path/filepath"

	"gopkg.in/natefinch/lumberjack.v2"
)

type slogLogger struct {
	logger *slog.Logger
}

func NewSlogLogger(config *LogConfig) (Logger, error) {
	var handler slog.Handler

	// 配置输出
	var writers []io.Writer

	// 添加控制台输出
	if config.EnableConsole {
		writers = append(writers, os.Stdout)
	}

	// 添加文件输出
	if config.FilePath != "" {
		if err := os.MkdirAll(filepath.Dir(config.FilePath), 0744); err != nil {
			return nil, err
		}

		fileWriter := &lumberjack.Logger{
			Filename:   config.FilePath,
			MaxSize:    config.MaxSize,
			MaxBackups: config.MaxBackups,
			MaxAge:     config.MaxAge,
			Compress:   config.Compress,
		}
		writers = append(writers, fileWriter)
	}

	writer := io.MultiWriter(writers...)

	// 设置日志格式
	opts := &slog.HandlerOptions{
		Level:     parseLevel(config.Level),
		AddSource: config.EnableCaller,
	}

	if config.Format == "json" {
		handler = slog.NewJSONHandler(writer, opts)
	} else {
		handler = slog.NewTextHandler(writer, opts)
	}

	logger := slog.New(handler)

	return &slogLogger{
		logger: logger,
	}, nil
}

func (l *slogLogger) Debug(msg string, fields ...Fields) {
	l.logger.Debug(msg, l.convertFields(fields...)...)
}

func (l *slogLogger) Info(msg string, fields ...Fields) {
	l.logger.Info(msg, l.convertFields(fields...)...)
}

func (l *slogLogger) Warn(msg string, fields ...Fields) {
	l.logger.Warn(msg, l.convertFields(fields...)...)
}

func (l *slogLogger) Error(msg string, fields ...Fields) {
	l.logger.Error(msg, l.convertFields(fields...)...)
}

func (l *slogLogger) Fatal(msg string, fields ...Fields) {
	l.logger.Error(msg, l.convertFields(fields...)...)
	os.Exit(1)
}

func (l *slogLogger) WithContext(ctx context.Context) Logger {
	return &slogLogger{
		logger: l.logger.With("trace_id", getTraceID(ctx)),
	}
}

func (l *slogLogger) WithFields(fields Fields) Logger {
	return &slogLogger{
		logger: l.logger.With(l.convertFields(fields)...),
	}
}

func (l *slogLogger) GetLogger() interface{} {
	return l.logger
}

func (l *slogLogger) convertFields(fields ...Fields) []any {
	if len(fields) == 0 {
		return nil
	}

	attrs := make([]any, 0)
	for _, f := range fields {
		for k, v := range f {
			attrs = append(attrs, slog.Any(k, v))
		}
	}
	return attrs
}

func parseLevel(level string) slog.Level {
	switch level {
	case "debug":
		return slog.LevelDebug
	case "info":
		return slog.LevelInfo
	case "warn":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}
