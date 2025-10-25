..PHONY: build test run docker-build docker-run clean lint bench-setup bench-run bench-stress bench-normal

build:
	go build -o bin/server ./cmd/server

test:
	go test -v -race -coverprofile=coverage.out ./...

test-coverage:
	go test -v -race -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out

run:
	go run ./cmd/server

docker-build:
	docker build -t log-tracker:latest .

docker-run:
	docker-compose up

clean:
	rm -rf bin/
	rm -f coverage.out

lint:
	golangci-lint run

bench-setup:
	@echo "Installing Gobench..."
	go install github.com/gobench-io/gobench@master

bench-run:
	@echo "Starting Gobench server..."
	@echo "Open http://localhost:8080 to access the dashboard"
	@echo ""
	@echo "Available scenarios:"
	@echo "  - Normal Load: benchmarks/gobench/normal_load.go"
	@echo "  - Stress Test: benchmarks/gobench/stress_test.go"
	@echo "  - Spike Test:  benchmarks/gobench/spike_test.go"
	@echo ""
	@echo "Copy the desired scenario into Gobench UI"
	gobench
