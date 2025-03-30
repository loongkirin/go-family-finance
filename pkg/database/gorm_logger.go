package database

import (
	"context"
	"time"

	"github.com/loongkirin/go-family-finance/pkg/logger"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	gormlogger "gorm.io/gorm/logger"
)

// GormLogger GORM 自定义日志记录器
type GormLogger struct {
	logger        logger.Logger
	LogLevel      gormlogger.LogLevel
	SlowThreshold time.Duration
}

// NewGormLogger 创建新的 GORM 日志记录器
func NewGormLogger(log logger.Logger, level gormlogger.LogLevel, slowThreshold time.Duration) *GormLogger {
	return &GormLogger{
		logger:        log,
		LogLevel:      level,
		SlowThreshold: slowThreshold,
	}
}

// LogMode 设置日志级别
func (l *GormLogger) LogMode(level gormlogger.LogLevel) gormlogger.Interface {
	newLogger := *l
	newLogger.LogLevel = level
	return &newLogger
}

// Info 记录信息日志
func (l *GormLogger) Info(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= gormlogger.Info {
		l.logger.WithContext(ctx).Info(msg, logger.Fields{"data": data})
	}
}

// Warn 记录警告日志
func (l *GormLogger) Warn(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= gormlogger.Warn {
		l.logger.WithContext(ctx).Warn(msg, logger.Fields{"data": data})
	}
}

// Error 记录错误日志
func (l *GormLogger) Error(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= gormlogger.Error {
		l.logger.WithContext(ctx).Error(msg, logger.Fields{"data": data})
	}
}

// Trace 记录 SQL 查询日志
func (l *GormLogger) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	if l.LogLevel <= gormlogger.Silent {
		return
	}

	elapsed := time.Since(begin)
	sql, rows := fc()

	// 创建 OpenTelemetry span
	tr := otel.Tracer("gorm")
	ctx, span := tr.Start(ctx, "gorm.query")
	defer span.End()

	// 设置 span 属性
	span.SetAttributes(
		attribute.String("db.system", "postgresql"),
		attribute.String("db.statement", sql),
		attribute.Int64("db.rows_affected", rows),
		attribute.Int64("db.duration_ms", elapsed.Milliseconds()),
	)

	// 记录日志
	fields := logger.Fields{
		"sql":     sql,
		"rows":    rows,
		"elapsed": elapsed,
	}

	if err != nil {
		fields["error"] = err.Error()
		span.SetStatus(codes.Error, err.Error())
		span.RecordError(err)
		span.SetAttributes(attribute.String("error", err.Error()))
		l.logger.WithContext(ctx).Error("GORM query error", fields)
	} else {
		span.SetStatus(codes.Ok, "success")
		// 检查是否是慢查询
		if l.SlowThreshold > 0 && elapsed > l.SlowThreshold {
			l.logger.WithContext(ctx).Warn("Slow query detected", fields)
		} else {
			l.logger.WithContext(ctx).Debug("GORM query", fields)
		}
	}
}

// ParamsFilter 过滤敏感参数
func (l *GormLogger) ParamsFilter(ctx context.Context, sql string, params ...interface{}) (string, []interface{}) {
	if l.LogLevel == gormlogger.Silent {
		return "", nil
	}
	return sql, params
}
