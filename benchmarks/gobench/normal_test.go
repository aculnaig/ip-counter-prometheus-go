package main

import (
	"context"
	"encoding/json"
	"log"
	"math/rand"
	"time"

	httpClient "github.com/gobench-io/gobench/clients/http"
	"github.com/gobench-io/gobench/dis"
	"github.com/gobench-io/gobench/executor/scenario"
)

// export returns the normal load test scenarios
func exportNormalScenario() scenario.Vus {
	return scenario.Vus{
		{
			Nu:   10,  // 10 virtual users
			Rate: 100, // spawn all users within 100ms
			Fu:   logIngestionScenario,
		},
		{
			Nu:   5,   // 5 virtual users
			Rate: 100, // spawn all users within 100ms
			Fu:   metricsPollingScenario,
		},
	}
}

// logIngestionScenario simulates log ingestion traffic
func logIngestionScenario(ctx context.Context, vui int) {
	client, err := httpClient.NewHttpClient(ctx, "log-ingestion")
	if err != nil {
		log.Printf("[VU-%d] Failed to create HTTP client: %v", vui, err)
		return
	}

	// Target the log server
	logURL := "http://localhost:5000/logs"

	// Run for 2 minutes
	timeout := time.After(2 * time.Minute)

	// Sample IPs for realistic simulation
	ips := []string{
		"192.168.1.1", "192.168.1.2", "192.168.1.3",
		"10.0.0.1", "10.0.0.2", "10.0.0.3",
		"172.16.0.1", "172.16.0.2", "172.16.0.3",
		"203.0.113.1", "203.0.113.2", "203.0.113.3",
	}

	urls := []string{
		"/api/users", "/api/posts", "/api/comments",
		"/api/auth/login", "/api/auth/logout",
		"/health", "/metrics", "/api/data",
	}

	for {
		select {
		case <-timeout:
			log.Printf("[VU-%d] Log ingestion scenario completed", vui)
			return
		default:
			// Create realistic log entry
			logEntry := map[string]string{
				"timestamp": time.Now().UTC().Format(time.RFC3339),
				"ip":        ips[rand.Intn(len(ips))],
				"url":       urls[rand.Intn(len(urls))],
			}

			payload, _ := json.Marshal(logEntry)

			// Send POST request
			go func() {
				headers := map[string]string{
					"Content-Type": "application/json",
				}
				_, err := client.Post(ctx, logURL, payload, headers)
				if err != nil {
					log.Printf("[VU-%d] POST request failed: %v", vui, err)
				}
			}()

			// Send 20 requests per second with Poisson distribution
			dis.SleepRatePoisson(20)
		}
	}
}

// metricsPollingScenario simulates periodic metrics polling
func metricsPollingScenario(ctx context.Context, vui int) {
	client, err := httpClient.NewHttpClient(ctx, "metrics-polling")
	if err != nil {
		log.Printf("[VU-%d] Failed to create HTTP client: %v", vui, err)
		return
	}

	// Target the metrics server
	metricsURL := "http://localhost:9102/metrics"
	healthURL := "http://localhost:9102/health"

	// Run for 2 minutes
	timeout := time.After(2 * time.Minute)
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	log.Printf("[VU-%d] Starting metrics polling scenario", vui)

	for {
		select {
		case <-timeout:
			log.Printf("[VU-%d] Metrics polling scenario completed", vui)
			return
		case <-ticker.C:
			// Poll metrics endpoint
			go func() {
				_, err := client.Get(ctx, metricsURL, nil)
				if err != nil {
					log.Printf("[VU-%d] Metrics GET failed: %v", vui, err)
				}
			}()

			// Poll health endpoint
			go func() {
				_, err := client.Get(ctx, healthURL, nil)
				if err != nil {
					log.Printf("[VU-%d] Health GET failed: %v", vui, err)
				}
			}()
		}
	}
}
