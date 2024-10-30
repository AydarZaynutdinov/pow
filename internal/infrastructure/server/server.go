package server

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"sync"

	"github.com/AydarZaynutdinov/pow/internal/config"

	"github.com/go-chi/chi"
	"go.uber.org/zap"

	"github.com/AydarZaynutdinov/pow/internal/logger"
)

var ErrServerStarted = errors.New("server has already started")

type Controller interface {
	AddRoutes(router chi.Router)
}

type Server struct {
	server *http.Server
	config config.Server
	router chi.Router
	logger logger.Logger
}

func NewServer(config config.Server, logger logger.Logger) *Server {
	return &Server{
		config: config,
		router: chi.NewRouter(),
		logger: logger,
	}
}

func (s *Server) AddControllerRoutes(controller Controller) error {
	if s.server != nil {
		return ErrServerStarted
	}

	controller.AddRoutes(s.router)

	return nil
}

func (s *Server) Run(ctx context.Context, name string, wg *sync.WaitGroup) {
	wg.Add(1)

	go func() {
		defer wg.Done()

		s.run(ctx, name)
	}()
}

func (s *Server) run(ctx context.Context, name string) {
	s.server = &http.Server{
		Addr:         fmt.Sprintf("%s:%d", s.config.Host, s.config.Port),
		ReadTimeout:  s.config.ReadTimeout,
		WriteTimeout: s.config.WriteTimeout,
		IdleTimeout:  s.config.IdleTimeout,
		Handler:      s.router,
	}

	serverLogger := s.logger.WithContext(ctx).With(
		zap.String("name", name),
		zap.String("host", s.config.Host),
		zap.Int("port", s.config.Port),
	)

	go func() {
		<-ctx.Done()

		serverLogger.Info("shutting down http server")

		shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), s.config.ShutdownDuration)
		defer shutdownCancel()

		err := s.server.Shutdown(shutdownCtx)
		if err != nil {
			serverLogger.Error("failed to shutdown http server", zap.Error(err))
		} else {
			serverLogger.Info("completed http server shutdown")
		}
	}()

	serverLogger.Info("starting http server")

	if err := s.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		serverLogger.Error("failed while listening and serving http server", zap.Error(err))
	}
}
