package maintenance

import (
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Controller struct{}

func NewMaintenanceController() *Controller {
	return &Controller{}
}

func (c *Controller) AddRoutes(router chi.Router) {
	router.Mount("/metrics", promhttp.Handler())
	router.Mount("/debug", middleware.Profiler())

	// IMP probes.
	router.HandleFunc("/health/ready", ReadinessHandler)
	router.HandleFunc("/health/live", LivenessHandler)

	// DevPlatform probes.
	router.HandleFunc("/readyz", ReadinessHandler)
	router.HandleFunc("/healthz", LivenessHandler)
}
