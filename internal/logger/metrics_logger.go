package logger

import (
	"context"

	"go.uber.org/zap"

	"github.com/AydarZaynutdinov/pow/internal/metrics"
)

type MetricsLogger struct {
	logger  Logger
	metrics *metrics.LoggerMetrics
}

func NewMetricsLogger(logger Logger, metrics *metrics.LoggerMetrics) *MetricsLogger {
	return &MetricsLogger{
		logger:  logger,
		metrics: metrics,
	}
}

func (l *MetricsLogger) Debug(msg string, fields ...zap.Field) {
	l.logger.Debug(msg, fields...)
}

func (l *MetricsLogger) Info(msg string, fields ...zap.Field) {
	l.logger.Info(msg, fields...)
}

func (l *MetricsLogger) Warn(msg string, fields ...zap.Field) {
	l.logger.Warn(msg, fields...)

	l.metrics.LogRequests.WithLabelValues(metrics.LoggerMetricsValueWarn).Inc()
}

func (l *MetricsLogger) Error(msg string, fields ...zap.Field) {
	l.logger.Error(msg, fields...)

	l.metrics.LogRequests.WithLabelValues(metrics.LoggerMetricsValueError).Inc()
}

func (l *MetricsLogger) Fatal(msg string, fields ...zap.Field) {
	l.logger.Fatal(msg, fields...)

	l.metrics.LogRequests.WithLabelValues(metrics.LoggerMetricsValueFatal).Inc()
}

func (l MetricsLogger) With(fields ...zap.Field) Logger {
	return NewMetricsLogger(
		l.logger.With(fields...),
		l.metrics,
	)
}

func (l MetricsLogger) WithContext(ctx context.Context) Logger {
	return NewMetricsLogger(
		l.logger.WithContext(ctx),
		l.metrics,
	)
}
