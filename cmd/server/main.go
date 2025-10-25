package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/aculnaig/log-tracker/internal/config"
	"github.com/aculnaig/log-tracker/internal/server"
	"github.com/aculnaig/log-tracker/internal/tracker"
	"github.com/aculnaig/log-tracker/pkg/logger"
)

func main() {
	// Initialize structured logger
	lgr := logger.New(os.Getenv("LOG_LEVEL"))

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		lgr.Fatal("failed to load configuration", "error", err)
	}

	// Initialize IP tracker
	ipTracker := tracker.NewIPTracker(lgr)

	// Create servers
	logServer := server.NewLogServer(cfg.LogServer, ipTracker, lgr)
	metricsServer := server.NewMetricsServer(cfg.MetricsServer, ipTracker, lgr)

	// Start servers
	logServer.Start()
	metricsServer.Start()

	lgr.Info("servers started successfully",
		"log_port", cfg.LogServer.Port,
		"metrics_port", cfg.MetricsServer.Port)

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	lgr.Info("shutting down servers...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := logServer.Shutdown(ctx); err != nil {
		lgr.Error("log server shutdown error", "error", err)
	}

	if err := metricsServer.Shutdown(ctx); err != nil {
		lgr.Error("metrics server shutdown error", "error", err)
	}

	lgr.Info("servers stopped")
}
