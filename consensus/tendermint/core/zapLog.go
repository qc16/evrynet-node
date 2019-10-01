package core

import (
	"context"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type loggerKeyType int

const (
	loggerKey loggerKeyType = iota
)

func initLogger() *zap.Logger {
	// First, define our level-handling logic.
	highPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl >= zapcore.ErrorLevel
	})
	lowPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return lvl < zapcore.ErrorLevel
	})

	encoder := zapcore.NewConsoleEncoder(zap.NewProductionEncoderConfig())
	consoleDebugging := zapcore.Lock(os.Stdout)
	consoleErrors := zapcore.Lock(os.Stderr)

	core := zapcore.NewTee(
		zapcore.NewCore(encoder, consoleDebugging, lowPriority),
		zapcore.NewCore(encoder, consoleErrors, highPriority),
	)
	logger := zap.New(core)
	defer logger.Sync()

	return logger
}

// NewLoggerContext returns a context has a zap logger with the extra fields added
func NewLoggerContext(ctx context.Context, fields ...zap.Field) context.Context {
	return context.WithValue(ctx, loggerKey, NewLogger(ctx).With(fields...))
}

// NewLogger returns a zap logger with as much context as possible
func NewLogger(c context.Context) *zap.Logger {
	if c != nil {
		if cLogger, ok := c.Value(loggerKey).(*zap.Logger); ok {
			return cLogger
		}
	}
	return initLogger()
}

// NewSugaredLogger creates a new sugared logger
func NewSugaredLogger(c context.Context) *zap.SugaredLogger {
	log := NewLogger(c)
	sugar := log.Sugar()

	return sugar
}
