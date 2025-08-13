# Backend Cleanup TODO List

## Overview
This file tracks all identified issues and cleanup tasks for the backend codebase. The analysis found significant leftovers from a previous restaurant/menu application that need to be cleaned up.

## 🔴 CRITICAL - Breaking Issues & Major Leftovers

### 1. Product Module - Wrong Naming (Restaurant Leftovers)
- [ ] File `internal/product/module.go` declares `package menu` instead of `package product`
- [ ] Rename files from "dish" to "product":
  - [ ] `internal/product/handlers/dish.go` → `product.go`
  - [ ] `internal/product/models/dish.go` → `product.go`
  - [ ] `internal/product/repositories/dish.go` → `product.go`
  - [ ] `internal/product/services/dish.go` → `product.go`
- [ ] Change `MenuPublicHandler` to `ProductPublicHandler`
- [ ] Update all imports and references

### 2. Missing Product API Routes
- [ ] POST `/api/v1/organizations/:organizationId/products` - Create product (handler exists at line 51 in handlers/dish.go but route not registered)
- [ ] GET `/api/v1/organizations/:organizationId/products` - List products (handler exists at line 99 but route not registered)
- [ ] Fix public endpoint path: `/organization/:id/menu` should be `/public/organization/:id/menu`

### 3. Unused Database Tables
- [ ] Drop `locations` table (exists in migration but no code uses it)
- [ ] Review all 16 tables for usage:
  ```
  accounts, subscription_plans, subscriptions, organizations,
  locations (UNUSED), products, qr_codes, questionnaires,
  questions, question_templates, feedbacks, verification_tokens,
  subscription_usage, usage_events, team_invitations, team_members
  ```

## 🟡 MEDIUM - Code Organization Issues

### 4. Inconsistent Module Patterns
Current state:
- New pattern (with new_module.go): organization, qrcode, subscription, feedback
- Old pattern (with module.go): auth, analytics, product

Tasks:
- [ ] Convert product module to new_module.go pattern
- [ ] Remove product registration from `internal/providers/providers.go` (lines 132-135)
- [ ] Make product module self-register

### 5. AI Module Not Properly Integrated
Current state:
- `internal/ai/` exists with providers but not registered as module
- Only `question_generator.go` is actually used
- Routes exist in organization controller instead of dedicated controller

Tasks:
- [ ] Decision: Keep as module or just utilities?
- [ ] If keeping: Create proper module structure with routes
- [ ] If removing: Delete unused files:
  - [ ] `internal/ai/services/anthropic_provider.go`
  - [ ] `internal/ai/services/openai_provider.go`
  - [ ] `internal/ai/services/gemini_provider.go`

### 6. Invalid Swagger Documentation
- [ ] Remove docs for non-existent routes:
  - [ ] `/api/v1/organizations/{organizationId}/products` POST (line 50 handlers/dish.go)
  - [ ] `/api/v1/organizations/{organizationId}/products` GET (line 98 handlers/dish.go)
- [ ] Regenerate swagger docs after fixes

## 🟢 LOW - General Cleanup

### 7. Restaurant Terminology Cleanup
- [ ] Search and replace all "menu" references with appropriate terms
- [ ] Search and replace all "dish" references with "product"
- [ ] Check for restaurant-specific error messages

### 8. Provider Registration Comments
- [ ] Clean up outdated comments in `providers.go`:
  - Line 137: "Feedback repositories, services and handlers are now provided by the new feedback module"
  - Line 139: "QR code providers are now in the qrcode module"
  - Line 147: "Subscription providers are now in the subscription module"

## 📁 File Structure Issues Found

```
backend/internal/
├── ai/                    # Not properly integrated as module
├── analytics/             # ✅ OK
├── auth/                  # ✅ OK (uses old pattern but works)
├── feedback/              # ✅ OK
├── organization/          # ✅ OK
├── product/               # ❌ Uses "menu" package, missing routes
├── providers/             # ⚠️ Has manual registrations that should be removed
├── qrcode/                # ✅ OK
├── shared/                # ✅ OK
└── subscription/          # ✅ OK
```

## 🔍 Code Locations for Reference

- Server routes setup: `internal/shared/server/server.go:73-115`
- Product module: `internal/product/module.go:18-31`
- Providers registration: `internal/providers/providers.go:32-148`
- Initial migration: `backend/migrations/20250722051522_init_from_current_db.up.sql`

## ✅ Modules to Keep As-Is
- Auth
- Organization
- Feedback
- QRCode
- Analytics
- Subscription
- Shared

## 🚫 Confirmed Leftovers to Remove
- "menu" package naming
- "dish" file naming
- "locations" database table
- Unused AI providers (if not needed)
- MenuPublicHandler naming

## Execution Priority
1. Fix breaking issues (product routes, naming)
2. Clean database
3. Standardize module patterns
4. Clean up documentation
5. Remove all terminology leftovers

---
Last Analysis Date: 2025-01-13
Total Issues Found: ~25 items
Estimated Cleanup Effort: 2-3 hours