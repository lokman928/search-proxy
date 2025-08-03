package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type ILogger interface {
	WithField(string, any) *Logger
	WithFields(map[string]any) *Logger
	Debug(msg string)
	Info(msg string)
	Warn(msg string)
	Error(msg string)
}

var _ ILogger = (*Logger)(nil)

type Logger struct {
	logger        *zap.Logger
	sugaredLogger *zap.SugaredLogger
}

// NewLogger creates a new Logger instance with the given name.
func NewLogger(name string) *Logger {
	logger, loggerErr := buildZapLogger()
	if loggerErr != nil {
		panic("failed to create logger: " + loggerErr.Error())
	}

	logger = logger.Named(name)

	return &Logger{
		logger:        logger,
		sugaredLogger: logger.Sugar(),
	}
}

func buildZapLogger() (*zap.Logger, error) {
	cfg := zap.NewProductionConfig()

	// Change time to human-readable format.
	cfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	logger, loggerErr := cfg.Build(
		zap.WithCaller(false),
	)

	return logger, loggerErr
}

// GetZapLogger returns the underlying zap.Logger instance.
func (l *Logger) GetZapLogger() *zap.Logger {
	return l.logger
}

// WithField creates a new Logger with an additional field.
func (l *Logger) WithField(key string, value any) *Logger {
	return l.WithFields(map[string]any{key: value})
}

// WithFields creates a new Logger with additional fields.
func (l *Logger) WithFields(fields map[string]any) *Logger {
	// Prepare zap fields for core logger (typed).
	zapFields := make([]zap.Field, 0, len(fields))
	for k, v := range fields {
		zapFields = append(zapFields, zap.Any(k, v))
	}

	// Prepare key-value pairs for sugared logger.
	kvFields := make([]interface{}, 0, len(fields)*2)
	for k, v := range fields {
		kvFields = append(kvFields, k, v)
	}

	return &Logger{
		logger:        l.logger.With(zapFields...),
		sugaredLogger: l.sugaredLogger.With(kvFields...),
	}
}

func (l *Logger) Debug(msg string) {
	l.sugaredLogger.Debug(msg)
}

func (l *Logger) Info(msg string) {
	l.sugaredLogger.Info(msg)
}

func (l *Logger) Warn(msg string) {
	l.sugaredLogger.Warn(msg)
}

func (l *Logger) Error(msg string) {
	l.sugaredLogger.Error(msg)
}
