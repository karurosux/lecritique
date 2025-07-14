# LeCritique - AI Context Guide

## 🎯 Quick Overview
LeCritique is a SaaS platform for restaurants to collect dish-specific feedback via QR codes. Customers scan → select dish → answer targeted questions → restaurants get insights.

## 🛠️ Tech Stack
- **Backend**: Go 1.23, Echo, PostgreSQL, GORM, JWT auth
- **Frontend**: SvelteKit, Svelte 5 (use runes: $state, $derived, $effect)
- **Key Libraries**: Stripe (payments), AI providers (Anthropic/OpenAI/Gemini)

## 📁 Key Directories
backend/
├── internal/
│   ├── auth/          # Authentication & teams
│   ├── restaurant/    # Restaurant management  
│   ├── menu/          # Dishes management
│   ├── feedback/      # Questionnaires & feedback
│   ├── qrcode/        # QR code system
│   ├── subscription/  # Billing & plans
│   └── shared/        # Common utilities

frontend/
├── src/
│   ├── lib/
│   │   ├── api/       # Generated API client
│   │   ├── components/# Reusable components
│   │   └── stores/    # Svelte stores
│   └── routes/        # SvelteKit pages

## 🚨 Important Notes

### Always Remember
1. **Multi-tenant**: All queries must be scoped to account_id
2. **Svelte 5**: Use runes ($state, $derived, $effect) not old syntax
3. **API Client**: Run `npm run generate-api` after backend changes
4. **Subscription Limits**: Check limits before allowing actions
5. **UUID Keys**: All IDs are UUIDs, not integers
6. **UI Components**: Always use ui components in lib, if it does not exists and is a dumb component, create it
7. **Endpoint Updates**: Whenever updating backend endpoints, regenerate client and types
8. **Documentation**: Always document endpoints for swagger documentation

### Current Architecture
- Domain-Driven Design with clear module separation
- Repository pattern for data access
- Service layer for business logic
- JWT auth with refresh tokens
- Subscription-based feature flags

## 📊 Active Development
See [CURRENT_SPRINT.md](./CURRENT_SPRINT.md) for what's being worked on.

## 📚 References
- [TECHNICAL_REFERENCE.md](./TECHNICAL_REFERENCE.md) - API endpoints, database schema
- [PROGRESS_LOG.md](./PROGRESS_LOG.md) - Completed features and history
- [PROJECT_DESCR.md](./PROJECT_DESCR.md) - Full project description

## 🔧 Common Commands
```bash
# Backend
cd backend
make dev                    # Start with hot reload
make generate-frontend-api  # Generate Swagger docs

# Frontend  
cd frontend
npm run dev           # Start dev server
npm run generate-api  # Update API client
```

## 💡 Quick Context
When working on LeCritique:
1. Questionnaires MUST be dish-specific (core feature, not optional!)
2. Feedback flow: QR scan → select dish → answer questions
3. AI generates contextual questions for specific dishes
4. Subscription tiers limit features (Starter/Professional/Enterprise)
5. Team members can have different roles

### Design Guidelines
- Lets keep the same kind of design styles and rules all time
- Lets always use lucide icons in UI, no random SVG as possible

### Development Guidelines
- When adding elements, if its a root element, we should add an identifier class that allows you to identify the components, for example if its a user list the class "user-list" would be ideal

### Deployment Guidelines
- Dont execute apps, neither frontend or backend, user always executing it out of AI context

## 🚨 AI Interaction Guidelines
- Stop adding things I did not ask, it just make me spend more money