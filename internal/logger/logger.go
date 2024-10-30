package logger

import (
	"context"
	"fmt"

	"github.com/AydarZaynutdinov/pow/internal/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	loggerTimeKey       = "ts"
	loggerLevelKey      = "level"
	loggerNameKey       = "logger"
	loggerCallerKey     = "caller"
	loggerMessageKey    = "message"
	loggerStacktraceKey = "stacktrace"

	loggerErrorLevel = "error"
	loggerDebugLevel = "debug"

	loggerRFC3339Nano = "rfc3339nano"

	loggerFieldAppName = "app"
	appName            = "pow"
)

type Logger interface {
	Debug(msg string, fields ...zap.Field)
	Info(msg string, fields ...zap.Field)
	Warn(msg string, fields ...zap.Field)
	Error(msg string, fields ...zap.Field)
	Fatal(msg string, fields ...zap.Field)
	With(fields ...zap.Field) Logger
	WithContext(ctx context.Context) Logger
}

type ZapLogger struct {
	log *zap.Logger
}

func NewZapLogger(loggerConfig *config.LoggerConfig) (*ZapLogger, error) {
	var (
		level                          = zapcore.DebugLevel
		timeFormat zapcore.TimeEncoder = zapcore.EpochTimeEncoder
	)

	if loggerConfig.Level != "" {
		if err := level.Set(loggerConfig.Level); err != nil {
			return nil, fmt.Errorf("failed to set logger's level (expected: debug, info, warn, error, dpanic, panic, fatal): %w", err)
		}
	}

	if loggerConfig.Time != "" {
		if err := timeFormat.UnmarshalText([]byte(loggerConfig.Time)); err != nil {
			return nil, fmt.Errorf("failed to set logger's time format (expected: rfc3339, rfc3339nano, iso8601, millis, nanos): %w", err)
		}
	}

	cfg := zap.NewProductionConfig()
	cfg.EncoderConfig = zapcore.EncoderConfig{
		MessageKey:     loggerMessageKey,
		LevelKey:       loggerLevelKey,
		TimeKey:        loggerTimeKey,
		NameKey:        loggerNameKey,
		CallerKey:      loggerCallerKey,
		FunctionKey:    zapcore.OmitKey,
		StacktraceKey:  loggerStacktraceKey,
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.RFC3339NanoTimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	// by default zap sends logs to stderr
	cfg.OutputPaths = []string{"stdout"}
	cfg.ErrorOutputPaths = []string{"stdout"}

	cfg.Level = zap.NewAtomicLevelAt(level)
	cfg.EncoderConfig.EncodeTime = timeFormat

	logger, err := cfg.Build()
	if err != nil {
		return nil, err
	}

	logger = logger.With(
		zap.String(loggerFieldAppName, appName),
	)

	return &ZapLogger{
		log: logger,
	}, nil
}

func NewErrorLogger() (*ZapLogger, error) {
	return NewZapLogger(
		&config.LoggerConfig{
			Level: loggerErrorLevel,
			Time:  loggerRFC3339Nano,
		},
	)
}

func NewTestLogger() Logger {
	rawLogger, _ := NewZapLogger(
		&config.LoggerConfig{
			Level: loggerDebugLevel,
			Time:  loggerRFC3339Nano,
		},
	)
	return rawLogger
}

func (l *ZapLogger) Debug(msg string, fields ...zap.Field) {
	l.log.Debug(msg, fields...)
}

func (l *ZapLogger) Info(msg string, fields ...zap.Field) {
	l.log.Info(msg, fields...)
}

func (l *ZapLogger) Warn(msg string, fields ...zap.Field) {
	l.log.Warn(msg, fields...)
}

func (l *ZapLogger) Error(msg string, fields ...zap.Field) {
	l.log.Error(msg, fields...)
}

func (l *ZapLogger) Fatal(msg string, fields ...zap.Field) {
	l.log.Fatal(msg, fields...)
}

func (l *ZapLogger) GetZap() *zap.Logger {
	return l.log
}

func (l ZapLogger) With(fields ...zap.Field) Logger {
	l.log = l.log.With(fields...)

	return &l
}

func (l ZapLogger) WithContext(ctx context.Context) Logger {
	fields := getLogFieldsFromContext(ctx)
	if fields == nil {
		return &l
	}

	l.log = l.log.With(fields...)

	return &l
}
