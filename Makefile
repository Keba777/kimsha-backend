.PHONY: dev build docker-up docker-down tidy

dev:
	go run ./cmd/server

build:
	go build -o bin/kimsha ./cmd/server

docker-up:
	docker-compose up -d

docker-down:
	docker-compose down

tidy:
	go mod tidy

lint:
	golangci-lint run ./...
