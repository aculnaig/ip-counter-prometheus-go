package server

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/aculnaig/log-tracker/internal/config"
	"github.com/aculnaig/log-tracker/internal/models"
	"github.com/aculnaig/log-tracker/internal/tracker"
	"github.com/aculnaig/log-tracker/pkg/middleware"
)

type LogServer struct {
	server  *http.Server
	tracker *tracker.IPTracker
	logger  Logger
}

type Logger interface {
	Debug(msg string, keysAndValues ...interface{})
	Info(msg string, keysAndValues ...interface{})
	Error(msg string, keysAndValues ...interface{})
	Fatal(msg string, keysAndValues ...interface{})
}

func NewLogServer(cfg config.ServerConfig, tracker *tracker.IPTracker, logger Logger) *LogServer {
	ls := &LogServer{
		tracker: tracker,
		logger:  logger,
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/logs", ls.handleLogs)
	mux.HandleFunc("/health", ls.handleHealth)

	ls.server = &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.Port),
		Handler:      middleware.Chain(mux, middleware.Logging(logger), middleware.Recovery(logger)),
		ReadTimeout:  cfg.ReadTimeout,
		WriteTimeout: cfg.WriteTimeout,
		IdleTimeout:  cfg.IdleTimeout,
	}

	return ls
}

func (s *LogServer) Start() {
	go func() {
		s.logger.Info("log server starting", "addr", s.server.Addr)
		if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			s.logger.Fatal("log server failed", "error", err)
		}
	}()
}

func (s *LogServer) Shutdown(ctx context.Context) error {
	return s.server.Shutdown(ctx)
}

func (s *LogServer) handleLogs(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var logEntry models.LogEntry
	if err := json.NewDecoder(r.Body).Decode(&logEntry); err != nil {
		s.logger.Debug("invalid JSON payload", "error", err)
		http.Error(w, "invalid JSON", http.StatusBadRequest)
		return
	}

	s.tracker.Add(logEntry.IP)

	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(map[string]string{"status": "accepted"})
}

func (s *LogServer) handleHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":    "healthy",
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	})
}
