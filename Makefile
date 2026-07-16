include .env
export

# Goose picks up these environment variables automatically.
# See: https://github.com/pressly/goose#environment-variables
export GOOSE_DBSTRING := host=$(DB_HOST) port=$(DB_PORT) user=$(DB_USER) password=$(DB_PASSWORD) dbname=$(DB_NAME) sslmode=disable
export GOOSE_DRIVER := postgres
export GOOSE_MIGRATION_DIR := database/migrations

## help: Show available make targets
.PHONY: help
help:
	@echo "Usage:"
	@echo "  make migrate-up         - Apply all pending migrations"
	@echo "  make migrate-down       - Roll back the last migration"
	@echo "  make migrate-status     - Show current migration status"
	@echo "  make migrate-reset      - Roll back all migrations"
	@echo "  make migrate-create     - Create a new migration file (usage: make migrate-create name=add_email_to_users)"
	@echo "  make build              - Build the application binary"
	@echo "  make run                - Run the application"
	@echo "  make test               - Run all tests"

## migrate-up: Apply all pending migrations
.PHONY: migrate-up
migrate-up:
	goose up

## migrate-down: Roll back the last migration
.PHONY: migrate-down
migrate-down:
	goose down

## migrate-status: Show current migration status
.PHONY: migrate-status
migrate-status:
	goose status

## migrate-reset: Roll back all migrations
.PHONY: migrate-reset
migrate-reset:
	goose reset

## migrate-create: Create a new migration file (usage: make migrate-create name=add_email_to_users)
.PHONY: migrate-create
migrate-create:
	@if [ -z "$(name)" ]; then echo "Error: provide a migration name, e.g. make migrate-create name=add_email_to_users"; exit 1; fi
	goose create $(name) sql

## build: Build the application binary
.PHONY: build
build:
	go build -o bin/task-manager-go main.go

## run: Run the application
.PHONY: run
run:
	go run main.go

## test: Run all tests
.PHONY: test
test:
	go test ./...
