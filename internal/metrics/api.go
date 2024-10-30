package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

const (
	apiControllerLabelPath = "path"
	apiControllerLabelCode = "code"
)

type APIControllerMetrics struct {
	Requests         *prometheus.CounterVec
	RequestsDuration *prometheus.HistogramVec
}

func NewAPIControllerMetrics() *APIControllerMetrics {
	metrics := &APIControllerMetrics{
		Requests: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "http_requests_total",
				Help: "Counter of HTTP requests for any HTTP-based requests",
			},
			[]string{
				apiControllerLabelPath,
				apiControllerLabelCode,
			},
		),
		RequestsDuration: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Name: "http_response_times",
				Help: "Response times in ms",
			},
			[]string{
				apiControllerLabelPath,
				apiControllerLabelCode,
			},
		),
	}

	metrics.register()

	return metrics
}

func (m *APIControllerMetrics) register() {
	prometheus.MustRegister(
		m.Requests,
		m.RequestsDuration,
	)
}
