.PHONY: up down migrate-up migrate-down run-api run-worker build

up:
	docker-compose up -d

down:
	docker-compose down

migrate-up:
	go run cmd/migrate/main.go up

migrate-down:
	go run cmd/migrate/main.go down

run-api:
	go run cmd/api/main.go

run-worker:
	go run cmd/worker/main.go

build:
	go build -o bin/api ./cmd/api
	go build -o bin/worker ./cmd/worker