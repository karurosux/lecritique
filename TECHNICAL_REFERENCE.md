# Technical Reference - LeCritique

## API Endpoints

### Public (No Auth)
- `GET /public/qr/:code` - Validate QR
- `GET /public/restaurant/:id/menu` - Get menu
- `GET /public/questionnaire/:restaurantId/:dishId` - Get questionnaire
- `POST /public/feedback` - Submit feedback

### Protected (Auth Required)
See Swagger docs at `/swagger` for full list. Key endpoints:
- Auth: `/api/auth/*`
- Restaurants: `/api/restaurants/*`
- Dishes: `/api/restaurants/:id/dishes/*`
- Questionnaires: `/api/restaurants/:id/questionnaires/*`
- QR Codes: `/api/restaurants/:id/qr-codes/*`

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