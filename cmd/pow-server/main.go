package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"runtime"
	"sync"
	"syscall"
	"time"

	"github.com/AydarZaynutdinov/pow/internal/cache"
	"github.com/AydarZaynutdinov/pow/internal/config"
	"github.com/AydarZaynutdinov/pow/internal/controller/api"
	"github.com/AydarZaynutdinov/pow/internal/controller/maintenance"
	"github.com/AydarZaynutdinov/pow/internal/infrastructure/server"
	"github.com/AydarZaynutdinov/pow/internal/logger"
	"github.com/AydarZaynutdinov/pow/internal/metrics"
	challengeService "github.com/AydarZaynutdinov/pow/internal/service/challenge"

	"github.com/go-playground/validator/v10"
	"go.uber.org/automaxprocs/maxprocs"
	"go.uber.org/zap"
)

var (
	configFilePath string

	signals = []os.Signal{syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT}
)

func init() {
	// parsing config file path
	flagConfig := flag.String("config", "config/config.dev.yaml", "config file path")
	flag.Parse()

	configFilePath = *flagConfig
}

func main() {
	errLogger, err := logger.NewErrorLogger()
	if err != nil {
		log.Fatalf("failed to create error logger: %s", err)
	}

	validate := validator.New()

	cfg, err := config.New(configFilePath, validate)
	if err != nil {
		errLogger.Fatal("failed to load main config", zap.Error(err))
	}

	rawLogger, err := logger.NewZapLogger(&cfg.Logger)
	if err != nil {
		errLogger.Fatal("failed to create main logger", zap.Error(err))
	}

	// metrics
	loggerMetrics := metrics.NewLoggerMetrics()
	apiControllerMetrics := metrics.NewAPIControllerMetrics()
	cacheMetrics := metrics.NewCacheMetrics()

	mainLogger := logger.NewMetricsLogger(rawLogger, loggerMetrics)

	// set GOMAXPROCS according to CPU
	mainLogger.Info("current GOMAXPROCS", zap.Int("GOMAXPROCS", runtime.GOMAXPROCS(0)))
	if _, err = maxprocs.Set(maxprocs.Logger(func(format string, args ...any) {
		mainLogger.Info(fmt.Sprintf(format, args...))
	})); err != nil {
		mainLogger.Warn("failed to change GOMAXPROCS according to cpu quotas", zap.Error(err))
	}

	// check signals to stop the app
	sigChan := make(chan os.Signal, len(signals))
	signal.Notify(sigChan, signals...)
	signalContext, signalContextCancel := context.WithCancel(context.Background())
	go func() {
		sig := <-sigChan

		mainLogger.Info("system signal was received", zap.String("signal", sig.String()))
		signalContextCancel()
	}()

	wg := &sync.WaitGroup{}

	// cache
	appCache, err := cache.NewCache(signalContext, cfg.Cache, cacheMetrics)
	if err != nil {
		errLogger.Fatal("failed to setup cache", zap.Error(err))
	}

	// services
	challenge := challengeService.NewService(appCache, cfg.QuotesList, cfg.Challenge, mainLogger)

	// handlers
	handler := api.NewHandler(mainLogger)

	// controllers
	apiController := api.NewAPIController(
		handler,
		challenge,
		apiControllerMetrics,
		mainLogger,
	)
	maintenanceController := maintenance.NewMaintenanceController()

	// servers
	apiServer := server.NewServer(
		cfg.AppServer,
		mainLogger,
	)
	err = apiServer.AddControllerRoutes(apiController)
	if err != nil {
		mainLogger.Fatal("failed to add routes to api server", zap.Error(err))
	}
	apiServer.Run(signalContext, "pow", wg)

	maintenanceServer := server.NewServer(
		cfg.MaintenanceServer,
		mainLogger,
	)
	err = maintenanceServer.AddControllerRoutes(maintenanceController)
	if err != nil {
		mainLogger.Fatal("failed to add routes to maintenance server", zap.Error(err))
	}
	maintenanceServer.Run(signalContext, "maintenance", wg)

	<-signalContext.Done()

	if waitTimeout(wg, cfg.AppServer.ShutdownDuration) {
		mainLogger.Info("timed out while waiting for graceful shutdown, resorting to forceful shutdown")
	}

	mainLogger.Info("service shutdown")
}

// waitTimeout waits for the WaitGroup for the specified max timeout.
// Returns true if waiting timed out.
func waitTimeout(wg *sync.WaitGroup, timeout time.Duration) bool {
	c := make(chan struct{})
	go func() {
		defer close(c)
		wg.Wait()
	}()

	select {
	case <-c:
		return false // completed normally
	case <-time.After(timeout):
		return true // timed out
	}
}
