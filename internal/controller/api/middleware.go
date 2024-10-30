package api

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi"
	"github.com/urfave/negroni"
	"go.uber.org/zap"

	"github.com/AydarZaynutdinov/pow/internal/logger"
)

const (
	apiControllerRequestTimeout = 30 * time.Second
)

func (c *Controller) contextMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := logger.AddRequestInfoToContext(r.Context(), r)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (c *Controller) panicRecoveryMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				ctx := r.Context()
				_ = c.handler.InternalServerError(ctx, w)
				c.logger.WithContext(ctx).Error("http handler panic recovery", zap.Any("error", err))
			}
		}()

		next.ServeHTTP(w, r)
	})
}

func (c *Controller) timeoutMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx, cancel := context.WithTimeout(r.Context(), apiControllerRequestTimeout)
		defer cancel()

		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}

func (c *Controller) metricsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		wrappedWriter := negroni.NewResponseWriter(w)

		start := time.Now()
		next.ServeHTTP(wrappedWriter, r)
		duration := time.Since(start).Seconds()

		routeContext := chi.RouteContext(ctx)
		routePattern := routeContext.RoutePattern()

		c.metrics.Requests.WithLabelValues(
			routePattern,
			strconv.Itoa(wrappedWriter.Status()),
		).Inc()
		c.metrics.RequestsDuration.WithLabelValues(
			routePattern,
			strconv.Itoa(wrappedWriter.Status()),
		).Observe(duration)
	})
}
