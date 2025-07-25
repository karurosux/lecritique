# Kyooar API Endpoints Reference

## Base Configuration
- **Base URL**: `http://localhost:8080/api/v1`
- **Authentication**: JWT Bearer Token (for protected endpoints)
- **Rate Limiting**: 100 requests/minute per IP

## Public Endpoints (No Auth Required)

### System Health
| Method | Endpoint | Purpose |
|--------|----------|---------|
| GET | `/api/health` | Service health check |

### Customer-Facing Endpoints
| Method | Endpoint | Purpose |
|--------|----------|---------|
| GET | `/public/qr/:code` | Validate QR code and get organization info |
| GET | `/public/organization/:id/menu` | Get available products for a organization |
| GET | `/public/questionnaire/:organizationId/:productId` | Get feedback questions for a product |
| POST | `/public/feedback` | Submit customer feedback |

#### Request/Response Examples:

**GET /public/qr/:code**
```json
// Response
{
  "success": true,
  "data": {
    "organization": {
      "id": "uuid",
      "name": "Organization Name",
      "description": "...",
      "settings": {...}
    },
    "qr_code": {
      "id": "uuid",
      "label": "Table 1",
      "type": "table"
    }
  }
}
```

**POST /public/feedback**
```json
// Request
{
  "qr_code_id": "uuid",
  "product_id": "uuid",
  "customer_name": "John Doe",
  "customer_email": "john@example.com",
  "customer_phone": "+1234567890",
  "overall_rating": 5,
  "responses": [
    {
      "question_id": "uuid",
      "answer": "Great taste!"
    }
  ]
}
```

## Authentication Endpoints

| Method | Endpoint | Purpose |
|--------|----------|---------|
| POST | `/auth/register` | Create new account |
| POST | `/auth/login` | Login and get JWT token |
| POST | `/auth/refresh` | Refresh JWT token |
| POST | `/auth/send-verification` | Send email verification (requires auth) |
| GET | `/auth/verify-email` | Verify email address with token |
| POST | `/auth/forgot-password` | Send password reset email |
| POST | `/auth/reset-password` | Reset password with token |

### Request/Response Examples:

**POST /auth/register**
```json
// Request
{
  "email": "owner@organization.com",
  "password": "securepassword123",
  "company_name": "My Organization Group"
}

// Response
{
  "success": true,
  "data": {
    "account": {
      "id": "uuid",
      "email": "owner@organization.com",
      "company_name": "My Organization Group"
    },
    "message": "Registration successful. Please check your email to verify your account."
  }
}
```

**POST /auth/login**
```json
// Request
{
  "email": "owner@organization.com",
  "password": "securepassword123"
}

// Response
{
  "success": true,
  "data": {
    "token": "jwt.token.here",
    "account": {
      "id": "uuid",
      "email": "owner@organization.com",
      "company_name": "My Organization Group"
    }
  }
}
```

**POST /auth/send-verification**
```json
// Request (Requires Authorization header)
{}

// Response
{
  "success": true,
  "data": {
    "message": "Verification email sent successfully"
  }
}
```

**GET /auth/verify-email?token=verification_token_here**
```json
// Response (Success)
{
  "success": true,
  "data": {
    "message": "Email verified successfully"
  }
}

// Response (Invalid/Expired Token)
{
  "success": false,
  "error": {
    "code": "INVALID_TOKEN",
    "message": "Invalid or expired verification token"
  }
}
```

**POST /auth/forgot-password**
```json
// Request
{
  "email": "owner@organization.com"
}

// Response
{
  "success": true,
  "data": {
    "message": "If an account with this email exists, a password reset link has been sent"
  }
}
```

**POST /auth/reset-password**
```json
// Request
{
  "token": "reset_token_here",
  "new_password": "newsecurepassword123"
}

// Response (Success)
{
  "success": true,
  "data": {
    "message": "Password reset successfully"
  }
}

// Response (Invalid/Expired Token)
{
  "success": false,
  "error": {
    "code": "INVALID_TOKEN",
    "message": "Invalid or expired reset token"
  }
}
```

## Protected Endpoints (Requires Authentication)

### Organization Management

| Method | Endpoint | Purpose |
|--------|----------|---------|
| POST | `/organizations` | Create new organization |
| GET | `/organizations` | List all organizations for account |
| GET | `/organizations/:id` | Get specific organization details |
| PUT | `/organizations/:id` | Update organization information |
| DELETE | `/organizations/:id` | Delete organization |

### Product Management

| Method | Endpoint | Purpose |
|--------|----------|---------|
| POST | `/products` | Create new product |
| GET | `/organizations/:organizationId/products` | List products for a organization |
| GET | `/products/:id` | Get specific product details |
| PUT | `/products/:id` | Update product information |
| DELETE | `/products/:id` | Delete product |

### Request/Response Examples:

**POST /organizations**
```json
// Request
{
  "name": "Downtown Bistro",
  "description": "Casual dining organization",
  "phone": "+1234567890",
  "email": "info@bistro.com",
  "website": "https://bistro.com"
}

// Response
{
  "success": true,
  "data": {
    "id": "uuid",
    "account_id": "uuid",
    "name": "Downtown Bistro",
    "description": "Casual dining organization",
    "is_active": true,
    "created_at": "2025-06-21T10:00:00Z"
  }
}
```

**POST /products**
```json
// Request
{
  "organization_id": "uuid",
  "name": "Grilled Salmon",
  "description": "Fresh Atlantic salmon with herbs",
  "category": "Main Course",
  "price": 24.99,
  "currency": "USD"
}

// Response
{
  "success": true,
  "data": {
    "id": "uuid",
    "organization_id": "uuid",
    "name": "Grilled Salmon",
    "price": 24.99,
    "is_available": true,
    "is_active": true
  }
}
```

## QR Code Management (Protected)

| Method | Endpoint | Purpose |
|--------|----------|---------|
| POST | `/organizations/:organizationId/qr-codes` | Generate QR codes |
| GET | `/organizations/:organizationId/qr-codes` | List QR codes for organization |
| DELETE | `/qr-codes/:id` | Delete QR code |

### Request/Response Examples:

**POST /organizations/:organizationId/qr-codes**
```json
// Request
{
  "type": "table",
  "label": "Table 5"
}

// Response
{
  "success": true,
  "data": {
    "id": "uuid",
    "organization_id": "uuid",
    "code": "LCQ-abc123-timestamp",
    "label": "Table 5",
    "type": "table",
    "is_active": true,
    "scans_count": 0,
    "created_at": "2025-06-21T10:00:00Z"
  }
}
```

## Feedback Management (Protected)

| Method | Endpoint | Purpose |
|--------|----------|---------|
| GET | `/organizations/:organizationId/feedback` | List feedback with pagination |
| GET | `/organizations/:organizationId/analytics` | Get basic feedback statistics |

### Request/Response Examples:

**GET /organizations/:organizationId/feedback?page=1&limit=20**
```json
// Response
{
  "success": true,
  "data": [
    {
      "id": "uuid",
      "customer_name": "John Doe",
      "overall_rating": 5,
      "created_at": "2025-06-21T10:00:00Z",
      "product": {...},
      "qr_code": {...}
    }
  ],
  "meta": {
    "total": 150,
    "page": 1,
    "limit": 20,
    "total_pages": 8
  }
}
```

## Analytics (Protected)

| Method | Endpoint | Purpose |
|--------|----------|---------|
| GET | `/analytics/organizations/:organizationId` | Get comprehensive organization analytics |
| GET | `/analytics/products/:productId` | Get product-specific analytics |

### Request/Response Examples:

**GET /analytics/organizations/:organizationId**
```json
// Response
{
  "success": true,
  "data": {
    "organization_id": "uuid",
    "organization_name": "Downtown Bistro",
    "total_feedback": 500,
    "average_rating": 4.2,
    "feedback_today": 12,
    "feedback_this_week": 85,
    "feedback_this_month": 320,
    "top_rated_products": [...],
    "lowest_rated_products": [...]
  }
}
```

## Future Endpoints (Planned)

### Subscription Management
| Method | Endpoint | Purpose |
|--------|----------|---------|
| GET | `/subscriptions/current` | Get current subscription |
| GET | `/subscriptions/plans` | List available plans |
| POST | `/subscriptions/subscribe` | Subscribe to a plan |

### Team Management
| Method | Endpoint | Purpose |
|--------|----------|---------|
| GET | `/team/members` | List team members |
| POST | `/team/invite` | Invite team member |
| DELETE | `/team/members/:id` | Remove team member |

## Error Response Format
```json
{
  "success": false,
  "error": {
    "code": "ERROR_CODE",
    "message": "Human readable error message",
    "details": {} // Optional additional details
  },
  "meta": {
    "timestamp": "2025-06-21T10:00:00Z",
    "version": "1.0"
  }
}
```

## Common Error Codes
- `BAD_REQUEST` - Invalid request data
- `UNAUTHORIZED` - Missing or invalid authentication
- `FORBIDDEN` - Authenticated but not authorized
- `NOT_FOUND` - Resource not found
- `CONFLICT` - Resource already exists
- `VALIDATION_ERROR` - Input validation failed
- `INVALID_TOKEN` - Invalid or expired verification/reset token
- `EMAIL_ALREADY_VERIFIED` - Email address already verified
- `SUBSCRIPTION_LIMIT` - Subscription limit reached
- `RATE_LIMIT` - Too many requests

## Headers
### Required for Protected Endpoints
```
Authorization: Bearer <jwt_token>
Content-Type: application/json
```

### Response Headers
```
X-Request-ID: <unique_request_id>
Content-Type: application/json
```

## Implementation Status

### ✅ Implemented
- Health check endpoint
- Authentication (register, login, refresh)
- Email verification endpoints
- Password reset endpoints
- Organization CRUD operations
- Product CRUD operations
- Public feedback submission
- QR code validation
- Menu retrieval
- QR code generation and management
- Feedback retrieval for organization owners
- Basic analytics and statistics

### 📋 TODO
- Subscription management endpoints
- Team management endpoints
- Advanced analytics endpoints
- Questionnaire customization endpoints
- Location management endpoints
- Notification settings endpoints

## Database Models Available

1. **Account** - Organization owner accounts
2. **User** - Team members
3. **Organization** - Organization entities
4. **Location** - Organization locations
5. **Product** - Menu items
6. **QRCode** - QR codes for tables/locations
7. **Questionnaire** - Custom feedback questionnaires
8. **Question** - Individual questions in questionnaires
9. **Feedback** - Customer feedback submissions
10. **VerificationToken** - Email verification and password reset tokens
11. **Subscription** - Account subscriptions
12. **SubscriptionPlan** - Available subscription plans

## Service Layer Architecture

### Available Services
- **AuthService** - Authentication and JWT management
- **OrganizationService** - Organization business logic
- **ProductService** - Product management
- **QRCodeService** - QR code generation and validation
- **FeedbackService** - Feedback submission and analytics

### Available Repositories
- **AccountRepository**
- **TokenRepository**
- **OrganizationRepository**
- **ProductRepository**
- **QRCodeRepository**
- **FeedbackRepository**
- **SubscriptionRepository**
- **QuestionnaireRepository**

## Middleware Stack
1. Request logging
2. Error recovery
3. CORS handling
4. Request ID generation
5. Security headers
6. Body size limiting (2MB)
7. Gzip compression
8. Rate limiting
9. JWT authentication (for protected routes)

## Development Notes

### Adding New Endpoints
1. Define handler method in appropriate handler file
2. Add route in `internal/server/routes/routes.go`
3. Implement service method if needed
4. Add repository method if needed
5. Update this documentation

### Testing Endpoints
```bash
# Health check
curl http://localhost:8080/api/health

# Register
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"password123","company_name":"Test Organization"}'

# Login
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"password123"}'

# Protected endpoint example
curl http://localhost:8080/api/v1/organizations \
  -H "Authorization: Bearer <token>"
```
