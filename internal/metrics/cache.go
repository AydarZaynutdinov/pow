package metrics

import "github.com/prometheus/client_golang/prometheus"

const (
	cacheLabelMethod = "method"
	cacheLabelStatus = "status"
)

type CacheMetrics struct {
	Requests         *prometheus.CounterVec
	RequestsDuration *prometheus.HistogramVec
}

func NewCacheMetrics() *CacheMetrics {
	metrics := &CacheMetrics{
		Requests: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "cache_requests_total",
				Help: "Counter of cache requests with any method",
			},
			[]string{
				cacheLabelMethod,
				cacheLabelStatus,
			},
		),
		RequestsDuration: prometheus.NewHistogramVec(
			prometheus.HistogramOpts{
				Name: "cache_response_times",
				Help: "Response times in ms",
			},
			[]string{
				cacheLabelMethod,
				cacheLabelStatus,
			},
		),
	}

	metrics.register()

	return metrics
}

func (m *CacheMetrics) register() {
	prometheus.MustRegister(
		m.Requests,
		m.RequestsDuration,
	)
}
