package logger

import (
	"context"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/rs/zerolog"
	"gopkg.in/natefinch/lumberjack.v2"
)

type zerologLogger struct {
	logger *zerolog.Logger
}

func NewZerologLogger(config *LogConfig) (Logger, error) {
	// 设置日志级别
	level, err := zerolog.ParseLevel(config.Level)
	if err != nil {
		level = zerolog.InfoLevel
	}
	zerolog.SetGlobalLevel(level)

	// 配置输出
	var writers []io.Writer

	// 添加控制台输出
	if config.EnableConsole {
		consoleWriter := zerolog.ConsoleWriter{
			Out:        os.Stdout,
			TimeFormat: time.RFC3339,
		}
		writers = append(writers, consoleWriter)
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

	// 创建多输出writer
	var writer io.Writer = io.MultiWriter(writers...)

	// 创建logger
	logger := zerolog.New(writer).With().Timestamp()
	if config.EnableCaller {
		logger = logger.Caller()
	}
	l := logger.Logger()

	return &zerologLogger{
		logger: &l,
	}, nil
}

func (l *zerologLogger) Debug(msg string, fields ...Fields) {
	event := l.logger.Debug()
	l.addFields(event, fields...)
	event.Msg(msg)
}

func (l *zerologLogger) Info(msg string, fields ...Fields) {
	event := l.logger.Info()
	l.addFields(event, fields...)
	event.Msg(msg)
}

func (l *zerologLogger) Warn(msg string, fields ...Fields) {
	event := l.logger.Warn()
	l.addFields(event, fields...)
	event.Msg(msg)
}

func (l *zerologLogger) Error(msg string, fields ...Fields) {
	event := l.logger.Error()
	l.addFields(event, fields...)
	event.Msg(msg)
}

func (l *zerologLogger) Fatal(msg string, fields ...Fields) {
	event := l.logger.Fatal()
	l.addFields(event, fields...)
	event.Msg(msg)
}

func (l *zerologLogger) WithContext(ctx context.Context) Logger {
	newLogger := l.logger.With().Str("trace_id", getTraceID(ctx)).Logger()
	return &zerologLogger{logger: &newLogger}
}

func (l *zerologLogger) WithFields(fields Fields) Logger {
	newLogger := l.logger.With().Fields(fields).Logger()
	return &zerologLogger{logger: &newLogger}
}

func (l *zerologLogger) GetLogger() interface{} {
	return l.logger
}

// func (l *zerologLogger) addFields(event *zerolog.Event, fields ...Fields) {
// 	for _, f := range fields {
// 		event.Fields(f)
// 	}
// }

func (l *zerologLogger) addFields(event *zerolog.Event, fields ...Fields) {
	for _, field := range fields {
		for key, value := range field {
			event.Interface(key, value)
		}
	}
}
