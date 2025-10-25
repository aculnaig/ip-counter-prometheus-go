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

// export returns the spike test scenario
func exportSpikeScenario() scenario.Vus {
	return scenario.Vus{
		{
			Nu:   200,  // 200 virtual users for sudden spike
			Rate: 2000, // spawn all within 2 seconds
			Fu:   spikeScenario,
		},
	}
}

func spikeScenario(ctx context.Context, vui int) {
	client, err := httpClient.NewHttpClient(ctx, "spike-test")
	if err != nil {
		log.Printf("[VU-%d] Failed to create HTTP client: %v", vui, err)
		return
	}

	logURL := "http://localhost:5000/logs"

	// Phase 1: Sudden spike (30 seconds)
	spikePhase := time.After(30 * time.Second)

	ips := generateIPs(vui, 20)

	log.Printf("[VU-%d] Starting spike phase", vui)

	for {
		select {
		case <-spikePhase:
			log.Printf("[VU-%d] Spike test completed", vui)
			return
		default:
			logEntry := map[string]string{
				"timestamp": time.Now().UTC().Format(time.RFC3339),
				"ip":        ips[rand.Intn(len(ips))],
				"url":       "/api/endpoint",
			}

			payload, _ := json.Marshal(logEntry)
			headers := map[string]string{"Content-Type": "application/json"}

			go client.Post(ctx, logURL, payload, headers)

			// 150 requests per second per VU during spike
			dis.SleepRatePoisson(150)
		}
	}
}

func generateIPs(vui, count int) []string {
	ips := make([]string, count)
	for i := 0; i < count; i++ {
		ips[i] = fmt.Sprintf("10.%d.%d.%d", vui/256, vui%256, i)
	}
	return ips
}
