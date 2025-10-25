package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"time"

	httpClient "github.com/gobench-io/gobench/clients/http"
	"github.com/gobench-io/gobench/dis"
	"github.com/gobench-io/gobench/executor/scenario"
)

// export returns the stress test scenario
func exportStressScenario() scenario.Vus {
	return scenario.Vus{
		{
			Nu:   100,  // 100 virtual users
			Rate: 1000, // spawn within 1 second
			Fu:   highLoadScenario,
		},
	}
}

func highLoadScenario(ctx context.Context, vui int) {
	client, err := httpClient.NewHttpClient(ctx, fmt.Sprintf("stress-vu-%d", vui))
	if err != nil {
		log.Printf("[VU-%d] Failed to create HTTP client: %v", vui, err)
		return
	}

	logURL := "http://localhost:5000/logs"
	timeout := time.After(5 * time.Minute)

	// Generate unique IPs per VU to test IP tracking scalability
	baseIP := fmt.Sprintf("10.%d.%d", vui/256, vui%256)

	for {
		select {
		case <-timeout:
			return
		default:
			logEntry := map[string]string{
				"timestamp": time.Now().UTC().Format(time.RFC3339),
				"ip":        fmt.Sprintf("%s.%d", baseIP, rand.Intn(255)),
				"url":       "/api/endpoint",
			}

			payload, _ := json.Marshal(logEntry)
			headers := map[string]string{"Content-Type": "application/json"}

			go client.Post(ctx, logURL, payload, headers)

			// 100 requests per second per VU
			dis.SleepRatePoisson(100)
		}
	}
}
