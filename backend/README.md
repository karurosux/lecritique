# LeCritique API

A subscription-based feedback platform that enables restaurants to collect dish-specific insights through QR code-driven questionnaires.

## Features

- Multi-restaurant account management
- Subscription-based access control
- QR code generation for tables/locations
- Dish-specific questionnaires
- Real-time analytics and insights
- Team collaboration

## Tech Stack

- **Language**: Go 1.21+
- **Framework**: Echo v4
- **ORM**: GORM
- **Database**: PostgreSQL
- **Cache**: Redis
- **Migrations**: golang-migrate

## API Documentation

See [API_DOCS.md](./API_DOCS.md) for detailed endpoint documentation.

## Getting Started

### Prerequisites

- Go 1.21+
- PostgreSQL 14+
- Redis 7+
- Make
- Docker & Docker Compose (optional)

### Installation

1. Clone the repository
2. Copy environment variables: `cp .env.example .env`
3. Install dependencies: `make deps`
4. Run migrations: `make migrate-up`
5. Run the application: `make run`

### Development

For hot reload during development:
```bash
# Install air
go install github.com/cosmtrek/air@latest

# Run with hot reload
make dev
```

### Docker

```bash
make docker-up
```

### Testing

```bash
make test
```
