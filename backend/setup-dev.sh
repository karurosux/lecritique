#!/bin/bash

set -e

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m'

echo -e "${GREEN}🚀 Setting up Kyooar Backend Development Environment${NC}"
echo ""

check_command() {
  if ! command -v $1 &>/dev/null; then
    echo -e "${RED}❌ $1 is not installed${NC}"
    return 1
  else
    echo -e "${GREEN}✅ $1 is installed${NC}"
    return 0
  fi
}

echo -e "${YELLOW}Checking prerequisites...${NC}"
all_installed=true

if ! check_command go; then
  all_installed=false
  echo "  Please install Go 1.21+ from https://golang.org/dl/"
fi

if ! check_command docker; then
  all_installed=false
  echo "  Please install Docker from https://docs.docker.com/get-docker/"
fi

# Docker compose now included in docker command, ergo, docker compose is the new command combination.
# if ! check_command docker-compose; then
#     all_installed=false
#     echo "  Please install Docker Compose from https://docs.docker.com/compose/install/"
# fi

if ! check_command make; then
  all_installed=false
  echo "  Please install Make"
fi

if [ "$all_installed" = false ]; then
  echo -e "${RED}Please install missing prerequisites and run this script again.${NC}"
  exit 1
fi

echo ""
echo -e "${YELLOW}Setting up environment...${NC}"

if [ ! -f .env ]; then
  echo "Creating .env file..."
  cp .env.example .env
  echo -e "${GREEN}✅ .env file created${NC}"
else
  echo -e "${YELLOW}⚠️  .env file already exists, skipping...${NC}"
fi

echo ""
echo -e "${YELLOW}Installing Go dependencies...${NC}"
go mod download
go mod tidy
echo -e "${GREEN}✅ Go dependencies installed${NC}"

echo ""
echo -e "${YELLOW}Installing development tools...${NC}"

if ! command -v air &>/dev/null; then
  echo "Installing air for hot reload..."
  go install github.com/cosmtrek/air@latest
  echo -e "${GREEN}✅ air installed${NC}"
else
  echo -e "${GREEN}✅ air already installed${NC}"
fi

if ! command -v migrate &>/dev/null; then
  echo "Installing golang-migrate..."
  go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
  echo -e "${GREEN}✅ migrate installed${NC}"
else
  echo -e "${GREEN}✅ migrate already installed${NC}"
fi

if ! command -v swag &>/dev/null; then
  echo "Installing swag for Swagger docs..."
  go install github.com/swaggo/swag/cmd/swag@latest
  echo -e "${GREEN}✅ swag installed${NC}"
else
  echo -e "${GREEN}✅ swag already installed${NC}"
fi

echo ""
echo -e "${YELLOW}Starting Docker services...${NC}"
docker compose up -d postgres redis
echo "Waiting for PostgreSQL to be ready..."
sleep 5

max_retries=30
retry_count=0
while ! docker compose exec -T postgres pg_isready -U postgres >/dev/null 2>&1; do
  retry_count=$((retry_count + 1))
  if [ $retry_count -gt $max_retries ]; then
    echo -e "${RED}❌ PostgreSQL failed to start${NC}"
    exit 1
  fi
  echo -n "."
  sleep 1
done
echo ""
echo -e "${GREEN}✅ PostgreSQL is ready${NC}"

echo ""
echo -e "${YELLOW}Creating Atlas development database...${NC}"
docker compose exec -T postgres psql -U postgres -c "SELECT 1 FROM pg_database WHERE datname = 'kyooar_atlas_dev'" | grep -q 1 || docker compose exec -T postgres psql -U postgres -c "CREATE DATABASE kyooar_atlas_dev;"
echo -e "${GREEN}✅ Atlas dev database created${NC}"

echo ""
echo -e "${YELLOW}Running database migrations...${NC}"
make migrate-up
echo -e "${GREEN}✅ Migrations completed${NC}"

echo ""
echo -e "${YELLOW}Seeding database with test data...${NC}"
make seed
echo -e "${GREEN}✅ Database seeded${NC}"

echo ""
echo -e "${YELLOW}Generating Swagger documentation...${NC}"
make swagger
echo -e "${GREEN}✅ Swagger docs generated${NC}"

echo ""
echo -e "${GREEN}🎉 Setup complete!${NC}"
echo ""
echo "Your development environment is ready. Here's what was set up:"
echo "  - PostgreSQL database running on port 5432"
echo "  - Redis cache running on port 6379"
echo "  - Go dependencies installed"
echo "  - Database migrations applied"
echo "  - Test user created (check cmd/seed/main.go for credentials)"
echo "  - Development tools installed (air, migrate, swag)"
echo ""
echo "To start the development server with hot reload:"
echo -e "  ${YELLOW}make dev${NC}"
echo ""
echo "To start the production server:"
echo -e "  ${YELLOW}make run${NC}"
echo ""
echo "To stop Docker services when done:"
echo -e "  ${YELLOW}docker compose down${NC}"
echo ""
echo "API will be available at: http://localhost:8080"
echo "API Documentation: http://localhost:8080/swagger/index.html"
