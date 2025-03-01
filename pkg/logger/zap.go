package logger

import (
	"context"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

type zapLogger struct {
	logger *zap.Logger
	sugar  *zap.SugaredLogger
}

func NewZapLogger(config *LogConfig) (Logger, error) {
	// 配置编码器
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.TimeKey = "timestamp"
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	// 配置输出
	var encoder zapcore.Encoder
	if config.Format == "json" {
		encoder = zapcore.NewJSONEncoder(encoderConfig)
	} else {
		encoder = zapcore.NewConsoleEncoder(encoderConfig)
	}

	// 配置写入器
	var cores []zapcore.Core

	// 添加控制台输出
	if config.EnableConsole {
		cores = append(cores, zapcore.NewCore(
			encoder,
			zapcore.AddSync(os.Stdout),
			parseZapLevel(config.Level),
		))
	}

	// 添加文件输出
	if config.FilePath != "" {
		fileWriter := &lumberjack.Logger{
			Filename:   config.FilePath,
			MaxSize:    config.MaxSize,
			MaxBackups: config.MaxBackups,
			MaxAge:     config.MaxAge,
			Compress:   config.Compress,
		}

		cores = append(cores, zapcore.NewCore(
			encoder,
			zapcore.AddSync(fileWriter),
			parseZapLevel(config.Level),
		))
	}

	// 创建logger
	core := zapcore.NewTee(cores...)
	logger := zap.New(core)

	if config.EnableCaller {
		logger = logger.WithOptions(zap.AddCaller())
	}

	return &zapLogger{
		logger: logger,
		sugar:  logger.Sugar(),
	}, nil
}

func (l *zapLogger) Debug(msg string, fields ...Fields) {
	l.sugar.Debugw(msg, l.convertFields(fields...)...)
}

func (l *zapLogger) Info(msg string, fields ...Fields) {
	l.sugar.Infow(msg, l.convertFields(fields...)...)
}

func (l *zapLogger) Warn(msg string, fields ...Fields) {
	l.sugar.Warnw(msg, l.convertFields(fields...)...)
}

func (l *zapLogger) Error(msg string, fields ...Fields) {
	l.sugar.Errorw(msg, l.convertFields(fields...)...)
}

func (l *zapLogger) Fatal(msg string, fields ...Fields) {
	l.sugar.Fatalw(msg, l.convertFields(fields...)...)
}

func (l *zapLogger) WithContext(ctx context.Context) Logger {
	return &zapLogger{
		logger: l.logger,
		sugar:  l.sugar.With("trace_id", getTraceID(ctx)),
	}
}

func (l *zapLogger) WithFields(fields Fields) Logger {
	return &zapLogger{
		logger: l.logger,
		sugar:  l.sugar.With(l.convertFields(fields)...),
	}
}

func (l *zapLogger) GetLogger() interface{} {
	return l.logger
}

func (l *zapLogger) convertFields(fields ...Fields) []interface{} {
	if len(fields) == 0 {
		return nil
	}

	args := make([]interface{}, 0)
	for _, f := range fields {
		for k, v := range f {
			args = append(args, k, v)
		}
	}
	return args
}

func parseZapLevel(level string) zapcore.Level {
	switch level {
	case "debug":
		return zapcore.DebugLevel
	case "info":
		return zapcore.InfoLevel
	case "warn":
		return zapcore.WarnLevel
	case "error":
		return zapcore.ErrorLevel
	default:
		return zapcore.InfoLevel
	}
}
