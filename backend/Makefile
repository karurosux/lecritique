.PHONY: build run test clean migrate-up migrate-down migrate-create docker-build docker-up docker-down swagger generate-frontend-api seed seed-force

# Variables
APP_NAME=lecritique
BUILD_DIR=bin
MAIN_FILE=cmd/api/main.go
DB_URL=postgres://postgres:postgres@localhost:5432/lecritique?sslmode=disable

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

# Seed database with default test data
seed:
	@echo "Creating default test user..."
	@go run cmd/seed/main.go

# Force recreate default test user
seed-force:
	@echo "Force recreating default test user..."
	@go run cmd/seed/main.go --force
