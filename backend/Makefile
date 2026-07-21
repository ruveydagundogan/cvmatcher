.PHONY: build run test lint migrate migrate-down docker-up docker-down clean

build:
	go build -o bin/server ./cmd/server

run:
	go run ./cmd/server

test:
	go test ./... -v

lint:
	go vet ./...

migrate:
	psql -U llmdecisionscore -d llmdecisionscore -f migrations/001_initial_schema.sql

docker-up:
	docker-compose -f deployments/docker-compose.yml up -d

docker-down:
	docker-compose -f deployments/docker-compose.yml down

clean:
	rm -rf bin/
