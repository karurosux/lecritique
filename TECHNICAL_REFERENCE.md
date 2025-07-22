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
- `GET /api/v1/public/organization/:id/menu` - Get organization menu
- `GET /api/v1/public/questionnaire/:organizationId/:productId` - Get questionnaire
- `GET /api/v1/public/organization/:organizationId/products/:productId/questions` - Get product questions
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

#### Organizations
- `GET /api/v1/organizations` - List organizations
- `POST /api/v1/organizations` - Create organization
- `GET /api/v1/organizations/:id` - Get organization
- `PUT /api/v1/organizations/:id` - Update organization
- `DELETE /api/v1/organizations/:id` - Delete organization

#### Productes
- `POST /api/v1/products` - Create product
- `GET /api/v1/products/:id` - Get product
- `PUT /api/v1/products/:id` - Update product
- `DELETE /api/v1/products/:id` - Delete product
- `GET /api/v1/organizations/:organizationId/products` - Get organization products

#### Questions
- `POST /api/v1/organizations/:organizationId/products/:productId/questions` - Create question
- `GET /api/v1/organizations/:organizationId/products/:productId/questions` - Get questions
- `PUT /api/v1/organizations/:organizationId/products/:productId/questions/:questionId` - Update
- `DELETE /api/v1/organizations/:organizationId/products/:productId/questions/:questionId` - Delete
- `POST /api/v1/organizations/:organizationId/products/:productId/questions/reorder` - Reorder

#### QR Codes
- `POST /api/v1/organizations/:organizationId/qr-codes` - Generate QR code
- `GET /api/v1/organizations/:organizationId/qr-codes` - List QR codes
- `PATCH /api/v1/qr-codes/:id` - Update QR code
- `DELETE /api/v1/qr-codes/:id` - Delete QR code

#### Analytics
- `GET /api/v1/analytics/organizations/:organizationId` - Organization analytics
- `GET /api/v1/analytics/organizations/:organizationId/charts` - Chart data
- `GET /api/v1/analytics/products/:productId` - Product analytics
- `GET /api/v1/analytics/dashboard/:organizationId` - Dashboard metrics

#### AI
- `POST /api/v1/ai/generate-questions/:productId` - Generate questions
- `POST /api/v1/ai/generate-questionnaire/:productId` - Generate questionnaire

## Database Schema

### Key Tables
```sql
accounts          # Organization owners
users             # Individual users
team_members      # Links users to accounts
organizations       # Organization entities
products            # Menu items
questionnaires    # Feedback forms (product_id optional)
questions         # Individual questions
feedbacks         # Customer responses
subscriptions     # Active plans
qr_codes          # QR tracking
```

### Important Relations
- Questionnaire → Product (optional, for product-specific)
- Feedback → Product (required)
- Everything scoped by account_id

## Subscription Plans

### Starter ($29/month)
- 1 organization
- 500 feedbacks/month
- 10 QR codes
- 2 team members

### Professional ($79/month)
- 3 organizations
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
const response = await api.api.v1OrganizationsProductesList(organizationId);
```
