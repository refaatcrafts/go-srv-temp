APP_NAME = go-srv-temp
DB_DSN ?= postgres://postgres:postgres@localhost:5432/ecommerce?sslmode=disable
GOOSE = go run github.com/pressly/goose/v3/cmd/goose@v3.27.1

.PHONY: run build test lint migrate-up migrate-down migrate-create

run:
	go run ./cmd/app

build:
	go build -o bin/$(APP_NAME) ./cmd/app

test:
	go test ./... -v -race -count=1

lint:
	golangci-lint run ./...

migrate-up:
	@echo "Running migrations up..."
	$(GOOSE) -dir migrations postgres "$(DB_DSN)" up

migrate-down:
	@echo "Running migrations down..."
	$(GOOSE) -dir migrations postgres "$(DB_DSN)" down

migrate-create:
	@read -p "Migration name: " name; \
	$(GOOSE) -dir migrations create "$$name" sql
