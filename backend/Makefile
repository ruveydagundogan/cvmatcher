.PHONY: build run test lint migrate migrate-up migrate-down docker-up docker-down clean mlc-llm-start

build:
	go build -ldflags="-s -w" -trimpath -o bin/server ./cmd/server

run:
	go run ./cmd/server

test:
	go test ./... -v

lint:
	go vet ./...
	staticcheck ./...

migrate-up:
	psql -U llmdecisionscore -d llmdecisionscore -f migrations/001_initial_schema.sql
	psql -U llmdecisionscore -d llmdecisionscore -f migrations/002_allow_anonymous_scoring.sql

migrate-down:
	@echo "Run drop statements manually"

docker-up:
	docker compose -f ../docker-compose.yml up -d postgres backend prometheus grafana

docker-down:
	docker compose -f ../docker-compose.yml down

docker-up-all:
	docker compose -f ../docker-compose.yml up -d

docker-grafana:
	docker compose -f ../docker-compose.yml up -d prometheus grafana

grafana-open:
	open http://localhost:3001

mlc-llm-start:
	pip install --pre -U -f https://mlc.ai/wheels mlc-llm-nightly mlc-ai-nightly
	python -m mlc_llm serve HF://mlc-ai/gemma-2b-it-q4f16_1-MLC --host 0.0.0.0 --port 8000 --mode interactive

clean:
	rm -rf bin/
