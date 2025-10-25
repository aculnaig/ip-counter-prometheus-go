package server

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/aculnaig/log-tracker/internal/config"
	"github.com/aculnaig/log-tracker/internal/tracker"
	"github.com/aculnaig/log-tracker/pkg/middleware"
)

type MetricsServer struct {
	server  *http.Server
	tracker *tracker.IPTracker
	logger  Logger
}

func NewMetricsServer(cfg config.ServerConfig, tracker *tracker.IPTracker, logger Logger) *MetricsServer {
	ms := &MetricsServer{
		tracker: tracker,
		logger:  logger,
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/metrics", ms.handleMetrics)
	mux.HandleFunc("/health", ms.handleHealth)

	ms.server = &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.Port),
		Handler:      middleware.Chain(mux, middleware.Logging(logger), middleware.Recovery(logger)),
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
		IdleTimeout:  cfg.IdleTimeout,
	}

	return ms
}

func (s *MetricsServer) Start() {
	go func() {
		s.logger.Info("metrics server starting", "addr", s.server.Addr)
		if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			s.logger.Fatal("metrics server failed", "error", err)
		}
	}()
}

func (s *MetricsServer) Shutdown(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}

func (s *MetricsServer) handleMetrics(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	count := s.tracker.Count()

	// Prometheus metric format (fixed typo: Promethues -> Prometheus)
	metrics := fmt.Sprintf("# HELP unique_ip_addresses Total number of unique IP addresses seen\n"+
		"# TYPE unique_ip_addresses gauge\n"+
		"unique_ip_addresses %d\n", count)

	w.Header().Set("Content-Type", "text/plain; version=0.0.4; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, metrics)
}

func (s *MetricsServer) handleHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, `{"status":"healthy","timestamp":"%s"}`, time.Now().UTC().Format(time.RFC3339))
}
