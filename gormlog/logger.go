package gormlog

import (
	"context"
	"log/slog"
	"time"

	"gorm.io/gorm/logger"
)

type loggerOptions struct {
	slowThreshold time.Duration
}

var defaultLoggerOptions = loggerOptions{
	slowThreshold: 2 * time.Second,
}

type LoggerOption func(*loggerOptions)

func WithSlowThreshold(threshold time.Duration) LoggerOption {
	return func(o *loggerOptions) {
		o.slowThreshold = threshold
	}
}

type Logger struct {
	logger *slog.Logger
	loggerOptions
}

// NewLogger returns a new logger that satisfies the GORM logger interface
// and uses the provided slog.Logger to log messages.
func NewLogger(logger *slog.Logger, opts ...LoggerOption) *Logger {
	l := &Logger{
		logger:        logger,
		loggerOptions: defaultLoggerOptions,
	}

	for _, opt := range opts {
		opt(&l.loggerOptions)
	}

	return l
}

// LogMode returns a copy of the logger. The provided log level is ignored.
//
// Don't use this method as it is implemented only to satisfy the GORM logger interface.
// Instead the log level is deternmined by the slog logger's internal slog.Leveler.
func (l Logger) LogMode(_ logger.LogLevel) logger.Interface {
	return l
}

func (l Logger) Info(ctx context.Context, msg string, data ...any) {
	l.logger.InfoContext(ctx, msg, data...)
}

func (l Logger) Warn(ctx context.Context, msg string, data ...any) {
	l.logger.WarnContext(ctx, msg, data...)
}

func (l Logger) Error(ctx context.Context, msg string, data ...any) {
	l.logger.ErrorContext(ctx, msg, data...)
}

func (l Logger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	elapsed := time.Since(begin)
	switch {
	case err != nil && l.logger.Enabled(ctx, slog.LevelError):
		attrs := buildTraceAttrs(elapsed, fc)
		l.logger.ErrorContext(ctx, err.Error(), attrs...)
	case elapsed > l.slowThreshold && l.logger.Enabled(ctx, slog.LevelWarn):
		attrs := buildTraceAttrs(elapsed, fc)
		l.logger.WarnContext(ctx, "SLOW QUERY", attrs...)
	case l.logger.Enabled(ctx, slog.LevelInfo):
		attrs := buildTraceAttrs(elapsed, fc)
		l.logger.InfoContext(ctx, "QUERY", attrs...)
	}
}

func buildTraceAttrs(elapsed time.Duration, fc func() (string, int64)) []any {
	sql, rows := fc()
	attrs := []any{
		slog.Duration("elapsed", elapsed),
		slog.String("sql", sql),
	}
	if rows != -1 {
		attrs = append(attrs, slog.Int64("rows", rows))
	}
	return attrs
}
