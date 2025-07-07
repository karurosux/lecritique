# Progress Log - LeCritique

## 2025

### January
- **Jan 8**: Fixed questionnaire builder
  - Added question CRUD API endpoints
  - Integrated with frontend
  - Added dish-questionnaire UI
  - Dishes show questionnaire status
  
- **Jan 7**: Questionnaire section added
  - Basic questionnaire creation
  - QR functionality ready

- **Jan 6**: Team features
  - Invite members working
  - Fixed page reloading issue

### December 2024
- Settings page (partial)
- Error handling improvements
- Background design updates

## Completed Features

### ✅ Core System
- Authentication (JWT, refresh tokens)
- Multi-tenant architecture
- Subscription management (Stripe)
- Team members with roles

### ✅ Restaurant Features  
- Restaurant CRUD
- Multi-location support
- Menu/dish management
- QR code generation & tracking

### ✅ Feedback System
- Questionnaire builder
- AI question generation
- Multiple question types
- Customer feedback submission

### ✅ Analytics
- Basic feedback analytics
- QR scan tracking
- Response metrics

## Lessons Learned
- Svelte 5 runes are cleaner than old syntax
- Always scope queries by account_id
- Swagger annotations save time
- Dish-specific questionnaires are the core value prop