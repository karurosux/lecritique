# Backend Cleanup Checklist

## 🔴 Critical Issues - High Priority

### 1. Product Module - Package Naming Issues
- [ ] Change `package menu` to `package product` in `internal/product/module.go`
- [ ] Rename all "dish" files to "product":
  - [ ] `handlers/dish.go` → `handlers/product.go`
  - [ ] `models/dish.go` → `models/product.go`
  - [ ] `repositories/dish.go` → `repositories/product.go`
  - [ ] `services/dish.go` → `services/product.go`
- [ ] Rename `MenuPublicHandler` to `ProductPublicHandler` in `handlers/public.go`
- [ ] Update all references from "menu" and "dish" terminology to "product"

### 2. Database Cleanup
- [ ] Create migration to drop unused `locations` table
- [ ] Review if all 16 tables in initial migration are actually used
- [ ] Document which tables are active vs legacy

### 3. Missing Product Routes
- [ ] Add POST route for creating products at `/api/v1/organizations/:organizationId/products`
- [ ] Add GET route for fetching products by organization at `/api/v1/organizations/:organizationId/products`
- [ ] Fix public menu endpoint path from `/organization/:id/menu` to `/public/organization/:id/menu`
- [ ] Ensure all handlers have corresponding routes registered

## 🟡 Medium Priority - Code Organization

### 4. Module Pattern Standardization
- [ ] Convert product module to use `new_module.go` pattern like other modules
- [ ] Remove manual product registration from `internal/providers/providers.go`
- [ ] Make product module self-register like organization, qrcode, subscription modules
- [ ] Standardize all modules to use consistent registration pattern

### 5. AI Module Decision
- [ ] Decide if AI should be a first-class module or just utilities
- [ ] If keeping as module:
  - [ ] Create proper `internal/ai/module.go`
  - [ ] Register routes properly in server.go
  - [ ] Move AI endpoints from organization controller to AI controller
- [ ] If not keeping as module:
  - [ ] Remove unused providers (keep only question_generator.go)
  - [ ] Remove `anthropic_provider.go` if not used
  - [ ] Remove `openai_provider.go` if not used
  - [ ] Remove `gemini_provider.go` if not used
  - [ ] Clean up AI configuration from config if not needed

### 6. Swagger Documentation
- [ ] Remove Swagger annotations for non-existent endpoints
- [ ] Update Swagger docs for actual routes only
- [ ] Regenerate docs with `swag init` after cleanup
- [ ] Verify all documented endpoints actually exist

## 🟢 Low Priority - General Cleanup

### 7. Code Consistency
- [ ] Review and remove any remaining restaurant/menu specific terminology
- [ ] Ensure all error messages are generic (not restaurant-specific)
- [ ] Check for any hardcoded strings that reference old domain

### 8. Provider Registration Cleanup
- [ ] Review `internal/providers/providers.go` for any other manual registrations
- [ ] Remove outdated comments like "providers are now in X module"
- [ ] Ensure consistent dependency injection patterns

### 9. File Organization
- [ ] Check for duplicate or similar functionality across modules
- [ ] Remove any empty or unnecessary files
- [ ] Ensure consistent file naming conventions across all modules

## 📊 Progress Tracking

### Modules Status
- ✅ **Auth** - Keep as is
- ✅ **Organization** - Keep as is
- ⚠️ **Product** - Needs major refactoring
- ✅ **Feedback** - Keep as is
- ✅ **QRCode** - Keep as is
- ✅ **Analytics** - Keep as is
- ✅ **Subscription** - Keep as is
- ⚠️ **AI** - Needs decision on structure
- ✅ **Shared** - Keep as is

## 🚀 Execution Order

1. **Phase 1 - Critical Fixes**
   - Fix product module naming and routes
   - Clean up database tables

2. **Phase 2 - Standardization**
   - Standardize module patterns
   - Organize AI functionality

3. **Phase 3 - Documentation**
   - Update Swagger docs
   - Clean up code comments

4. **Phase 4 - Final Cleanup**
   - Remove all leftovers
   - Ensure consistency

## Notes

- Each checkbox should be checked off when completed
- Add any new findings to appropriate sections
- Update status as work progresses
- Consider creating git branches for each phase