.PHONY: help docker-build docker-build-backend docker-build-frontend docker-push

help:
	@echo "Available commands:"
	@echo "  make docker-build     - Build Docker images for deployment"
	@echo "  make docker-build-backend - Build backend Docker image"
	@echo "  make docker-build-frontend - Build frontend Docker image"
	@echo "  make docker-push      - Push Docker images to registry"

build:
	@docker compose build
