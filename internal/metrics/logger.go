package metrics

import "github.com/prometheus/client_golang/prometheus"

const loggerMetricsLabelLevel = "level"

const (
	LoggerMetricsValueWarn  = "warn"
	LoggerMetricsValueError = "error"
	LoggerMetricsValueFatal = "fatal"
)

type LoggerMetrics struct {
	LogRequests *prometheus.CounterVec
}

func NewLoggerMetrics() *LoggerMetrics {
	metrics := &LoggerMetrics{
		LogRequests: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "errors_total",
				Help: "Total number of errors written by logger",
			},
			[]string{
				loggerMetricsLabelLevel,
			},
		),
	}

	metrics.register()

	return metrics
}

func (m *LoggerMetrics) register() {
	prometheus.MustRegister(
		m.LogRequests,
	)
}
