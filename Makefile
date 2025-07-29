.PHONY: help docker-build docker-build-backend docker-build-frontend docker-push

help:
	@echo "Available commands:"
	@echo "  make docker-build     - Build Docker images for deployment"
	@echo "  make docker-build-backend - Build backend Docker image"
	@echo "  make docker-build-frontend - Build frontend Docker image"
	@echo "  make docker-push      - Push Docker images to registry"

docker-build-backend:
	@echo "Building backend Docker image..."
	@docker build -t kyooar-backend:latest ./backend

docker-build-frontend:
	@echo "Building frontend Docker image..."
	@docker build -t kyooar-frontend:latest ./frontend

docker-build: docker-build-backend docker-build-frontend
	@echo "Docker images built successfully"

docker-push:
	@echo "Pushing Docker images to registry..."
	@docker push kyooar-backend:latest
	@docker push kyooar-frontend:latest
