.PHONY: build run test clean migrate-up migrate-down migrate-create docker-build docker-up docker-down swagger generate-frontend-api seed seed-force setup

# Variables
APP_NAME=kyooar
BUILD_DIR=bin
MAIN_FILE=cmd/api/main.go
DB_URL=postgres://postgres:postgres@localhost:5432/kyooar?sslmode=disable

# Build
build:
	@echo "Building..."
	@go build -o $(BUILD_DIR)/$(APP_NAME) $(MAIN_FILE)

# Run
run: build
	@echo "Running..."
	@./$(BUILD_DIR)/$(APP_NAME)

# Development run with hot reload
dev:
	@air

# Test
test:
	@echo "Testing..."
	@go test -v ./...

# Clean
clean:
	@echo "Cleaning..."
	@rm -rf $(BUILD_DIR)

# Database migrations
migrate-up:
	@migrate -path=migrations -database="$(DB_URL)" up

migrate-down:
	@migrate -path=migrations -database="$(DB_URL)" down

migrate-create:
	@migrate create -ext sql -dir migrations $(name)

# Generate migration from current database state using Atlas
migrate-from-db:
	@echo "Generating migration from current database state..."
	@../atlas migrate diff init_from_current_state \
		--dir "file://migrations" \
		--to "$(DB_URL)" \
		--dev-url "postgres://postgres:postgres@localhost:5432/kyooar_atlas_dev?sslmode=disable" \
		--format golang-migrate

# Docker
docker-build:
	@docker-compose build

docker-up:
	@docker-compose up -d

docker-down:
	@docker-compose down

# Install dependencies
deps:
	@go mod download
	@go mod tidy

# Install migrate CLI tool
install-migrate:
	@echo "Installing golang-migrate..."
	@go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
	@echo "migrate installed successfully!"

# Generate Swagger documentation
swagger:
	@echo "Generating Swagger documentation..."
	@swag init -g main.go -o docs

# Generate TypeScript API client for frontend
generate-frontend-api: swagger
	@echo "Generating TypeScript API client..."
	@cd ../frontend && npm run generate-api

# Seed subscription plans
seed-plans:
	@echo "Creating subscription plans..."
	@go run cmd/seed-plans/main.go

# Seed database with default test data
seed:
	@echo "Creating subscription plans..."
	@go run cmd/seed-plans/main.go
	@echo "Creating default test user..."
	@go run cmd/seed/main.go

# Force recreate default test user
seed-force:
	@echo "Creating subscription plans..."
	@go run cmd/seed-plans/main.go
	@echo "Force recreating default test user..."
	@go run cmd/seed/main.go --force

# Setup development environment
setup:
	@./setup-dev.sh
