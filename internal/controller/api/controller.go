package api

import (
	"context"

	"github.com/go-chi/chi"

	"github.com/AydarZaynutdinov/pow/internal/logger"
	"github.com/AydarZaynutdinov/pow/internal/metrics"
)

type Service interface {
	GenerateChallenge(ctx context.Context) (string, error)
	CheckSolution(ctx context.Context, challenge, solution string) (string, error)
	IsSolutionIncorrect(err error) bool
	GetDifficulty() int
}

type Controller struct {
	handler *Handler
	service Service
	metrics *metrics.APIControllerMetrics
	logger  logger.Logger
}

func NewAPIController(
	handler *Handler,
	service Service,
	metrics *metrics.APIControllerMetrics,
	logger logger.Logger,
) *Controller {
	return &Controller{
		handler: handler,
		service: service,
		metrics: metrics,
		logger:  logger,
	}
}

func (c *Controller) AddRoutes(router chi.Router) {
	router.Route("/", func(r chi.Router) {
		r.Use(c.contextMiddleware)
		r.Use(c.panicRecoveryMiddleware)
		r.Use(c.timeoutMiddleware)
		r.Use(c.metricsMiddleware)

		r.Route("/api", func(r chi.Router) {
			r.Route("/pow", func(r chi.Router) {
				r.Get("/challenge", c.handler.Handle(c.GetChallenge))
				r.Post("/solve", c.handler.Handle(c.Solve))
			})
		})
	})
}
