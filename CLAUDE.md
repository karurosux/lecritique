# Kyooar - AI Assistant Guidelines

## 🚨 CRITICAL RULES - NEVER IGNORE THESE

1. **NEVER EXECUTE APPS OR BUILD COMMANDS UNLESS INDICATED BY USER** - User handles all execution
2. **Use Lucide icons exclusively** - No SVG icons, only import from lucide-svelte
3. **Always run `make generate-client-api` after backend API changes**
4. **Use TodoWrite tool for multi-step tasks** - Track progress proactively

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

1. Always check subscription limits before operations, using middlware for backend and gates and guards for frontend.
2. Document all endpoints for Swagger
3. Add identifier classes to root elements (e.g., "user-list")
4. Follow existing code patterns and conventions
5. Check neighboring files before adding libraries

## 🔍 Quick References

- [TASKS.md](./TASKS.md) - Current work and todo items
- [TECHNICAL_REFERENCE.md](./TECHNICAL_REFERENCE.md) - API endpoints, database schema
- [PROJECT_DESCR.md](./PROJECT_DESCR.md) - Full project description
- [QUICK_START.md](./QUICK_START.md) - Essential info and commands

## 💻 Development Principles

- Use good coding principles like SOLID, YAGNI, KISS, and DRY, don't over-engineer solutions
- Questionnaires MUST be product-specific (core feature!)
- Never commit secrets or expose sensitive data

## 🎨 Code Style Rules

- **No AI comments** - Don't add explanatory comments to code
- **Remove existing comments** - Unless they provide important context or are TODO items
- **No emojis as icons** - Use Lucide icons only
- **Add identifier classes** - For root elements (e.g., "user-list")

## 📝 Task Management

- **Use TodoWrite for complex tasks** - 3+ steps or multi-step processes
- **Track progress in real-time** - Mark tasks completed immediately
- **One task in_progress at a time** - Focus on current work
