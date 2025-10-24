.PHONY: build test run docker-build docker-run clean lint

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
