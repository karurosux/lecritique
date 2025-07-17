# Technical Reference - LeCritique

## API Endpoints

### Public (No Auth)
- `GET /api/health` - Health check
- `POST /api/v1/auth/register` - Register new account
- `POST /api/v1/auth/login` - Login
- `GET /api/v1/auth/verify-email` - Verify email
- `POST /api/v1/auth/resend-verification` - Resend verification
- `POST /api/v1/auth/forgot-password` - Password reset
- `POST /api/v1/team/accept-invite` - Accept team invitation
- `GET /api/v1/plans` - Get subscription plans

### Public Feedback
- `GET /api/v1/public/qr/:code` - Validate QR code
- `GET /api/v1/public/restaurant/:id/menu` - Get restaurant menu
- `GET /api/v1/public/questionnaire/:restaurantId/:dishId` - Get questionnaire
- `GET /api/v1/public/restaurant/:restaurantId/dishes/:dishId/questions` - Get dish questions
- `POST /api/v1/public/feedback` - Submit feedback

### Protected (Auth Required)

#### Auth & Profile
- `POST /api/v1/auth/refresh` - Refresh token
- `PUT /api/v1/auth/profile` - Update profile
- `POST /api/v1/auth/change-email` - Request email change

#### Team Management
- `GET /api/v1/team/members` - List team members
- `POST /api/v1/team/members/invite` - Invite member
- `PUT /api/v1/team/members/:id/role` - Update role
- `DELETE /api/v1/team/members/:id` - Remove member

#### Restaurants
- `GET /api/v1/restaurants` - List restaurants
- `POST /api/v1/restaurants` - Create restaurant
- `GET /api/v1/restaurants/:id` - Get restaurant
- `PUT /api/v1/restaurants/:id` - Update restaurant
- `DELETE /api/v1/restaurants/:id` - Delete restaurant

#### Dishes
- `POST /api/v1/dishes` - Create dish
- `GET /api/v1/dishes/:id` - Get dish
- `PUT /api/v1/dishes/:id` - Update dish
- `DELETE /api/v1/dishes/:id` - Delete dish
- `GET /api/v1/restaurants/:restaurantId/dishes` - Get restaurant dishes

#### Questions
- `POST /api/v1/restaurants/:restaurantId/dishes/:dishId/questions` - Create question
- `GET /api/v1/restaurants/:restaurantId/dishes/:dishId/questions` - Get questions
- `PUT /api/v1/restaurants/:restaurantId/dishes/:dishId/questions/:questionId` - Update
- `DELETE /api/v1/restaurants/:restaurantId/dishes/:dishId/questions/:questionId` - Delete
- `POST /api/v1/restaurants/:restaurantId/dishes/:dishId/questions/reorder` - Reorder

#### QR Codes
- `POST /api/v1/restaurants/:restaurantId/qr-codes` - Generate QR code
- `GET /api/v1/restaurants/:restaurantId/qr-codes` - List QR codes
- `PATCH /api/v1/qr-codes/:id` - Update QR code
- `DELETE /api/v1/qr-codes/:id` - Delete QR code

#### Analytics
- `GET /api/v1/analytics/restaurants/:restaurantId` - Restaurant analytics
- `GET /api/v1/analytics/restaurants/:restaurantId/charts` - Chart data
- `GET /api/v1/analytics/dishes/:dishId` - Dish analytics
- `GET /api/v1/analytics/dashboard/:restaurantId` - Dashboard metrics

#### AI
- `POST /api/v1/ai/generate-questions/:dishId` - Generate questions
- `POST /api/v1/ai/generate-questionnaire/:dishId` - Generate questionnaire

## Database Schema

### Key Tables
```sql
accounts          # Restaurant owners
users             # Individual users
team_members      # Links users to accounts
restaurants       # Restaurant entities
dishes            # Menu items
questionnaires    # Feedback forms (dish_id optional)
questions         # Individual questions
feedbacks         # Customer responses
subscriptions     # Active plans
qr_codes          # QR tracking
```

### Important Relations
- Questionnaire → Dish (optional, for dish-specific)
- Feedback → Dish (required)
- Everything scoped by account_id

## Subscription Plans

### Starter ($29/month)
- 1 restaurant
- 500 feedbacks/month
- 10 QR codes
- 2 team members

### Professional ($79/month)
- 3 restaurants
- 2000 feedbacks/month
- 50 QR codes/location
- 5 team members

### Enterprise ($199/month)
- Unlimited everything
- API access
- Custom branding

## Code Patterns

### Backend Service
```go
type Service interface {
    Create(ctx context.Context, accountID uuid.UUID, model *Model) error
    // Always scope by accountID!
}
```

### Frontend Component (Svelte 5)
```svelte
<script lang="ts">
  let { data } = $props();  // Not export let
  let count = $state(0);    // Not just let
  let double = $derived(count * 2);  // Not $:
</script>
```

### API Client Usage
```ts
import { getApiClient } from '$lib/api';
const api = getApiClient();
const response = await api.api.v1RestaurantsDishesList(restaurantId);
```