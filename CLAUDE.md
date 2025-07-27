# Kyooar - AI Assistant Guidelines

## 🤖 Core Principles
1. **Do exactly what was asked** - No more, no less
2. **Never create files unless necessary** - Always prefer editing existing files
3. **No unsolicited documentation** - Only create docs when explicitly requested
4. **UI Components** - Always use existing components from lib/components

## 🛠️ Key Technical Context
- **Frontend**: SvelteKit with Svelte 5 (use runes: $state, $derived, $effect)
- **Backend**: Go 1.23, Echo, PostgreSQL, GORM
- **IDs**: All IDs are UUIDs, not integers
- **Multi-tenant**: All queries must be scoped to account_id

## 📋 Workflow Rules
1. After backend API changes: Run `npm run generate-api`
2. Always check subscription limits before operations
3. Document all endpoints for Swagger
4. Use lucide icons exclusively in UI
5. Add identifier classes to root elements (e.g., "user-list")
6. Never execute apps - user handles execution

## 🔍 Quick References
- [TASKS.md](./TASKS.md) - Current work and todo items
- [TECHNICAL_REFERENCE.md](./TECHNICAL_REFERENCE.md) - API endpoints, database schema
- [PROJECT_DESCR.md](./PROJECT_DESCR.md) - Full project description
- [QUICK_START.md](./QUICK_START.md) - Essential info and commands

## ⚡ Critical Reminders
- Questionnaires MUST be product-specific (core feature!)
- Follow existing code patterns and conventions
- Check neighboring files before adding libraries
- Never commit secrets or expose sensitive data

## 💻 Development Principles
- Use good coding principles like SOLID, YAGNI, KISS, and DRY, don't over-engineer solutions

## 💬 Code Commentary
- Dont add AI comments to the code

## 📝 Comment Handling
- If found comments, remove them unless it really provides important context or is a TODO item

## 🚫 Design Guidelines
- Dont use emojis as icons