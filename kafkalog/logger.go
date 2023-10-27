package kafkalog

import (
	"context"
	"log/slog"

	"github.com/twmb/franz-go/pkg/kgo"
)

// Logger is a wrapper around slog.Logger which implements the kgo logger interface.
// The log level is determined by the slog.Leveler used to create the slog.Logger.
type Logger struct {
	logger *slog.Logger
}

// NewLogger creates a new logger that satisfies the kgo logger interface
// and uses the provided slog.Logger to log messages.
func NewLogger(logger *slog.Logger) *Logger {
	return &Logger{
		logger: logger,
	}
}

func (l Logger) Level() kgo.LogLevel {
	switch {
	case l.logger.Enabled(context.Background(), slog.LevelDebug):
		return kgo.LogLevelDebug
	case l.logger.Enabled(context.Background(), slog.LevelInfo):
		return kgo.LogLevelInfo
	case l.logger.Enabled(context.Background(), slog.LevelWarn):
		return kgo.LogLevelWarn
	case l.logger.Enabled(context.Background(), slog.LevelError):
		return kgo.LogLevelError
	default:
		return kgo.LogLevelNone
	}
}

func (l Logger) Log(level kgo.LogLevel, msg string, keyvals ...any) {
	l.logger.Log(context.Background(), kgoLogLevelToSlogLogLevel(level), msg, keyvals...)
}

func kgoLogLevelToSlogLogLevel(lvl kgo.LogLevel) slog.Level {
	switch lvl {
	case kgo.LogLevelDebug:
		return slog.LevelDebug
	case kgo.LogLevelInfo:
		return slog.LevelInfo
	case kgo.LogLevelWarn:
		return slog.LevelWarn
	case kgo.LogLevelError:
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}
