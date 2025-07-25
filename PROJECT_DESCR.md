# Kyooar - Organization Feedback Platform

## Project Overview

Kyooar is a SaaS platform that enables organizations to collect detailed, product-specific feedback from customers through QR code-driven questionnaires. The platform helps organization owners gain actionable insights to improve their service, menu offerings, and overall customer experience.

## Core Concept

1. **Organization owners** subscribe to the service (monthly fee)
2. **Generate QR codes** for tables, locations, or takeaway orders
3. **Customers scan** the QR code after their meal
4. **Select the product** they ordered from the menu
5. **Answer targeted questions** specific to that product
6. **Organization receives** real-time feedback and analytics

## Key Features

### For Organization Owners

- **Multi-organization support** - Manage multiple organization locations under one account
- **Subscription-based pricing** - Tiered plans based on organization count and features
- **QR code management** - Generate and manage QR codes for different purposes (tables, takeaway, delivery)
- **Product catalog** - Maintain menu items with categories, prices, and availability
- **Custom questionnaires** - Create product-specific questions or use templates
- **Real-time analytics** - View feedback trends, ratings, and insights
- **Team collaboration** - Invite team members with different permission levels
- **Low rating alerts** - Get notified when feedback falls below threshold

### For Customers

- **Frictionless experience** - No app download or registration required
- **Quick feedback** - Complete questionnaire in under 2 minutes
- **Product-specific questions** - Relevant questions based on what they ordered
- **Multiple question types** - Ratings, scales, multiple choice, text feedback
- **Optional contact info** - Provide email/phone for follow-up (optional)
- **Mobile optimized** - Works seamlessly on any device

## Technical Architecture

### Tech Stack

- **Web Front End**: Sveltekit with Svelte 5
- **Backend**: Go 1.21+ with Echo framework
- **Database**: PostgreSQL with GORM ORM
- **Cache**: Redis for session management
- **Authentication**: JWT-based auth
- **API Design**: RESTful API with JSON responses
- **Architecture**: Monolith designed for easy microservices migration

### Key Design Decisions

1. **Monolith First, Microservices Ready**
   - All services in one codebase but organized by domain
   - Clear separation between domains for future extraction
   - Shared database with logical separation

2. **Multi-tenant Architecture**
   - Account-based isolation
   - Subscription limits enforced at service layer
   - Organization-level data separation

3. **UUID Primary Keys**
   - Better for distributed systems
   - No sequential ID exposure
   - Prevents ID enumeration attacks

4. **Repository Pattern**
   - Clean separation of data access logic
   - Easy to mock for testing
   - Database-agnostic interface

5. **Service Layer**
   - Business logic separated from HTTP handlers
   - Subscription limit enforcement
   - Cross-domain orchestration

## Domain Model

### Core Entities

1. **Account**
   - Represents a organization owner/company
   - Has subscription, organizations, team members
   - Email-based authentication

2. **Organization**
   - Individual organization location
   - Belongs to an account
   - Has settings, locations, products, QR codes

3. **Product**
   - Menu items with pricing and availability
   - Can have custom questionnaire
   - Tracks feedback per product

4. **QRCode**
   - Unique codes for customer entry points
   - Types: table, location, takeaway, delivery
   - Tracks scan count and last scan time

5. **Feedback**
   - Customer responses to questionnaires
   - Links to product, organization, and QR code
   - Stores ratings and text responses

6. **Subscription**
   - Controls account limits and features
   - Plans: Starter ($29), Professional ($79), Enterprise ($199)
   - Enforces organization count, feedback limits, etc.

## API Structure

### Public Endpoints (Customer-facing)

- QR code validation
- Organization menu retrieval
- Questionnaire fetching
- Feedback submission

### Protected Endpoints (Owner-facing)

- Authentication (register, login, refresh)
- Organization management (CRUD)
- Product management (CRUD)
- QR code generation
- Feedback analytics
- Team management
- Subscription management

## Subscription Tiers

### Starter - $29/month

- 1 organization
- 1 location
- 10 QR codes
- 500 feedbacks/month
- 2 team members
- Basic analytics

### Professional - $79/month

- 3 organizations
- 3 locations per organization
- 50 QR codes per location
- 2000 feedbacks/month
- 5 team members
- Advanced analytics

### Enterprise - $199/month

- Unlimited organizations
- Unlimited locations
- Unlimited QR codes
- Unlimited feedbacks
- Unlimited team members
- Advanced analytics
- Custom branding
- API access
- Priority support

## Development Workflow

### Local Development

```bash
# Start services
make docker-up

# Run migrations
make migrate-up

# Start server with hot reload
make dev
```

### Project Structure

```
kyooar-api/
├── cmd/api/          # Application entry point
├── internal/         # Private application code
│   ├── config/       # Configuration management
│   ├── models/       # Database models
│   ├── repositories/ # Data access layer
│   ├── services/     # Business logic
│   ├── handlers/     # HTTP handlers
│   └── middleware/   # HTTP middleware
├── pkg/              # Public packages
├── migrations/       # Database migrations
└── tests/           # Test files
```

## Future Enhancements

### Phase 1 (Current)

- ✅ Core authentication system
- ✅ Organization and product management
- ✅ Basic feedback collection
- ✅ QR code generation
- 🚧 Basic analytics

### Phase 2 (Next)

- [ ] Email notifications
- [ ] Advanced analytics dashboard
- [ ] Custom questionnaire builder
- [ ] Multi-language support
- [ ] White-label options

### Phase 3 (Future)

- [ ] Mobile apps for owners
- [ ] AI-powered insights
- [ ] Integration with POS systems
- [ ] Automated response to feedback
- [ ] Loyalty program integration

## Business Model

### Revenue Streams

1. **Subscription fees** - Monthly recurring revenue
2. **Usage-based pricing** - Additional fees for exceeding limits
3. **Enterprise features** - Custom pricing for large chains
4. **Add-on services** - SMS notifications, custom domains

### Target Market

1. **Small organizations** - Single location, basic needs
2. **Organization groups** - Multiple locations, standardization
3. **Organization chains** - Enterprise features, API integration
4. **Food trucks & pop-ups** - Mobile QR codes, simple setup

## Security Considerations

1. **Authentication** - JWT with refresh tokens
2. **Authorization** - Account-based access control
3. **Rate limiting** - 100 requests/minute per IP
4. **Data isolation** - Tenant data separation
5. **Input validation** - Request validation on all endpoints
6. **HTTPS only** - Encrypted communication
7. **CORS policy** - Restricted origins in production

## Monitoring & Analytics

### For Platform

- Request/response logging
- Error tracking
- Performance metrics
- Usage analytics
- Subscription metrics

### For Organizations

- Feedback volume trends
- Average ratings by product
- Customer sentiment analysis
- Response time metrics
- Popular products tracking

## Deployment Strategy

### Infrastructure

- **Hosting**: AWS/GCP with auto-scaling
- **Database**: Managed PostgreSQL
- **Cache**: Managed Redis
- **CDN**: CloudFlare for static assets
- **Monitoring**: Prometheus + Grafana

### CI/CD Pipeline

1. Git push triggers build
2. Run tests and linting
3. Build Docker image
4. Deploy to staging
5. Run integration tests
6. Deploy to production
7. Run smoke tests

## Success Metrics

### Platform KPIs

- Monthly Recurring Revenue (MRR)
- Customer Acquisition Cost (CAC)
- Customer Lifetime Value (CLV)
- Churn rate
- Active organizations
- Feedback volume

### Organization KPIs

- Average rating improvement
- Feedback response rate
- Time to resolution
- Customer satisfaction score
- Repeat feedback rate

## Contact & Support

- **Documentation**: API_ENDPOINTS.md
- **Support Email**: <support@kyooar.com>
- **Developer Portal**: docs.kyooar.com
- **Status Page**: status.kyooar.com
