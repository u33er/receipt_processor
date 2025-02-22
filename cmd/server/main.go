package main

import (
	"context"
	"errors"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	golog "log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"ticket-processor/internal/api/handlers"

	"ticket-processor/internal/api"
	"ticket-processor/internal/config"
	"ticket-processor/internal/services"
	"ticket-processor/internal/storage"
	"ticket-processor/pkg/logger"
)

func main() {
	cfg, err := config.MustLoad()
	if err != nil {
		golog.Fatalf("Failed to load configuration: %v", err)
	}

	log := logger.SetupLogger(cfg.Env)
	defer log.Sync()

	store := storage.NewInMemoryStore()
	cache := storage.NewInMemoryCache(log)
	receiptProcessor := services.NewReceiptProcessor(log, store, cache)
	receiptHandler := handlers.NewReceiptHandler(log, receiptProcessor)

	e := api.SetupRouter(log, cfg, receiptHandler)

	gracefulShutdown(e, log, cfg)
}

func gracefulShutdown(e *echo.Echo, log *zap.Logger, cfg *config.Config) {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	go func() {
		log.Info("Starting server", zap.String("address", cfg.HTTPServer.Address))
		if err := e.Start(cfg.HTTPServer.Address); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Error("Server error", zap.Error(err))
		}
	}()

	<-ctx.Done()

	log.Info("Shutting down server...")

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), cfg.HTTPServer.ShutdownTimeout)
	defer shutdownCancel()

	if err := e.Shutdown(shutdownCtx); err != nil {
		log.Error("Server shutdown error", zap.Error(err))
	}

	log.Info("Server stopped")
}
