package main

import (
	"context"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/AydarZaynutdinov/pow/internal/app"
	"github.com/AydarZaynutdinov/pow/internal/config"
	"github.com/AydarZaynutdinov/pow/internal/repository"
	"github.com/AydarZaynutdinov/pow/internal/service"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cfg, err := config.ParseServer()
	if err != nil {
		log.Fatal("Failed to parse server config:", err)
	}

	quoteRepo, err := repository.NewQuote(cfg.QuoteList)
	if err != nil {
		log.Fatal("Failed to create quote repository:", err)
	}
	powService := service.NewPoW(cfg.PoW.Complexity)
	server := app.NewServer(ctx, cfg, quoteRepo, powService)

	stopChan := make(chan os.Signal, 1)
	signal.Notify(
		stopChan,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
		syscall.SIGKILL,
	)

	go func() {
		if err := server.Start(); err != nil {
			log.Fatal("Failed to start server:", err)
		}
	}()

	slog.Info("Server is running")

	<-stopChan
	slog.Info("Shutting signal received")

	cancel()
	server.Shutdown()
	slog.Info("Server shutdown complete")
}
