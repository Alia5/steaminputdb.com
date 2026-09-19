package logging

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

type gormLogger struct {
	slowThreshold time.Duration
}

func NewGormLogger(slowThreshold time.Duration) gormlogger.Interface {
	return &gormLogger{
		slowThreshold: slowThreshold,
	}
}

func (l *gormLogger) LogMode(gormlogger.LogLevel) gormlogger.Interface {
	return l
}

func (l *gormLogger) Info(ctx context.Context, msg string, args ...any) {
	slog.InfoContext(ctx, fmt.Sprintf(msg, args...))
}

func (l *gormLogger) Warn(ctx context.Context, msg string, args ...any) {
	slog.WarnContext(ctx, fmt.Sprintf(msg, args...))
}

func (l *gormLogger) Error(ctx context.Context, msg string, args ...any) {
	slog.ErrorContext(ctx, fmt.Sprintf(msg, args...))
}

func (l *gormLogger) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	duration := time.Since(begin)

	var msg string
	level := slog.LevelDebug
	switch {
	case err != nil && !errors.Is(err, gorm.ErrRecordNotFound):
		msg, level = "SQL Query Error", slog.LevelError
	case l.slowThreshold > 0 && duration > l.slowThreshold:
		msg, level = "Slow SQL Query", slog.LevelWarn
	default:
		msg = "SQL Query"
	}

	if !slog.Default().Enabled(ctx, level) {
		return
	}

	query, rows := fc()
	attrs := []any{
		"query", query,
		"rows", rows,
		"duration", duration,
	}
	if level == slog.LevelError {
		attrs = append(attrs, "error", err)
	}
	slog.Log(ctx, level, msg, attrs...)
}
