# Kyooar Quick Start

## 🎯 Project Overview
SaaS platform for organizations to collect product-specific feedback via QR codes.

**Core Flow**: Customer scans QR → Selects product → Answers questions → Organization gets insights

## 🚀 Getting Started

### Backend (Go)
```bash
cd backend
make dev                    # Start server (localhost:8080)
make seed                   # Create test user: admin@kyooar.com / admin123
make generate-frontend-api  # Update Swagger docs
```

### Frontend (SvelteKit)
```bash
cd frontend
npm run dev          # Start dev server (localhost:5173)
npm run generate-api # Update API types after backend changes
```

## 🏗️ Architecture
- **Backend**: Go + Echo + PostgreSQL + GORM
- **Frontend**: SvelteKit + Svelte 5 (use runes!)
- **Auth**: JWT with refresh tokens
- **IDs**: All UUIDs, not integers

## 💰 Subscription Tiers
- **Starter** ($29): 1 location, 3 products, basic analytics
- **Professional** ($79): 3 locations, unlimited products, advanced analytics
- **Enterprise** ($199): Unlimited everything, AI insights, API access

## 📋 Key Commands
```bash
# Lint & Check
make lint          # Backend
npm run check      # Frontend

# Database
make migrate       # Run migrations
make seed-force    # Recreate test data
```

## 🔗 Important URLs
- Backend API: http://localhost:8080
- Frontend: http://localhost:5173
- Swagger: http://localhost:8080/swagger/index.html

## ⚡ Remember
- Questionnaires are **product-specific** (core feature!)
- Always check subscription limits
- Use existing UI components from lib/components
- Multi-tenant: scope queries to account_id
