package database

import (
	"context"
	"gorm.io/gorm/logger"
	"log/slog"
	"time"
)

type slogLogger struct {
	slog *slog.Logger
}

func (l slogLogger) LogMode(level logger.LogLevel) logger.Interface {
	return l
}

func (l slogLogger) Info(ctx context.Context, s string, i ...interface{}) {
	l.slog.InfoContext(ctx, s, i...)
}

func (l slogLogger) Warn(ctx context.Context, s string, i ...interface{}) {
	l.slog.WarnContext(ctx, s, i...)
}

func (l slogLogger) Error(ctx context.Context, s string, i ...interface{}) {
	l.slog.ErrorContext(ctx, s, i...)
}

func (l slogLogger) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	elapsed := time.Since(begin)
	sql, rows := fc()
	if err != nil {
		l.slog.ErrorContext(ctx, err.Error(), "query", sql, "elapsed", elapsed, "rows", rows)
	} else {
		l.slog.InfoContext(ctx, "SQL query executed", "query", sql, "elapsed", elapsed, "rows", rows)
	}
}
