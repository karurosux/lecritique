#!/bin/bash

set -e

echo "🐳 Building deployment images for Kyooar..."

# Check if .env exists
if [ ! -f .env ]; then
    echo "⚠️  .env file not found. Copying from .env.prod.example..."
    cp .env.prod.example .env
    echo "📝 Please edit .env with your actual values before running again"
    exit 1
fi

# Build images
echo "🔨 Building Docker images..."
docker-compose -f docker-compose.prod.yml build

echo "✅ Images built successfully!"

# Ask if user wants to start services
read -p "🚀 Start services now? (y/n): " -n 1 -r
echo
if [[ $REPLY =~ ^[Yy]$ ]]; then
    echo "🚀 Starting services..."
    docker-compose -f docker-compose.prod.yml up -d
    echo "✅ Services started!"
    echo "🌐 Frontend: http://localhost"
    echo "🔗 API: http://localhost:8080"
    echo ""
    echo "📊 Check status: docker-compose -f docker-compose.prod.yml ps"
    echo "📋 View logs: docker-compose -f docker-compose.prod.yml logs -f"
else
    echo "💡 To start services later: docker-compose -f docker-compose.prod.yml up -d"
fi