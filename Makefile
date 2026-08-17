.PHONY: all build-backend build-frontend run-backend run-frontend dev clean test docker-up docker-down docker-reset seed setup

all: build-backend build-frontend

build-backend:
	cd backend && go build -o bin/server ./cmd/main.go

build-frontend:
	cd frontend && npm run build

run-backend:
	cd backend && go run ./cmd/main.go

run-frontend:
	cd frontend && npm run dev

dev:
	@echo "Starting both backend and frontend..."
	@cd backend && go run ./cmd/main.go &
	@cd frontend && npm run dev &
	@wait

clean:
	rm -rf backend/bin frontend/dist

test:
	cd backend && go test -v -race -count=1 ./...

docker-up:
	docker compose up -d --build postgres redis backend frontend
	docker compose run --rm seed

docker-down:
	docker compose down

docker-reset:
	docker compose down -v
	docker compose up -d --build postgres redis backend frontend
	docker compose run --rm seed

seed:
	docker compose run --rm seed

setup:
	./setup.sh
