# Kyooar API Documentation

## Base URL
```
http://localhost:8080/api/v1
```

## Authentication
Most endpoints require authentication via JWT token in the Authorization header:
```
Authorization: Bearer <token>
```

## Endpoints

### Public Endpoints

#### Validate QR Code
```
GET /public/qr/:code
```
Response:
```json
{
  "success": true,
  "data": {
    "organization": {...},
    "qr_code": {
      "id": "uuid",
      "label": "Table 1",
      "type": "table"
    }
  }
}
```

#### Get Organization Menu
```
GET /public/organization/:id/menu
```

#### Submit Feedback
```
POST /public/feedback
```
Body:
```json
{
  "qr_code_id": "uuid",
  "product_id": "uuid",
  "customer_name": "John Doe",
  "customer_email": "john@example.com",
  "overall_rating": 5,
  "responses": [
    {
      "question_id": "uuid",
      "answer": "Great taste!"
    }
  ]
}
```

### Authentication Endpoints

#### Register
```
POST /auth/register
```
Body:
```json
{
  "email": "user@example.com",
  "password": "password123",
  "company_name": "My Organization"
}
```

#### Login
```
POST /auth/login
```
Body:
```json
{
  "email": "user@example.com",
  "password": "password123"
}
```

### Protected Endpoints

#### Organizations

##### Create Organization
```
POST /organizations
```

##### List Organizations
```
GET /organizations
```

##### Get Organization
```
GET /organizations/:id
```

##### Update Organization
```
PUT /organizations/:id
```

##### Delete Organization
```
DELETE /organizations/:id
```

#### Products

##### Create Product
```
POST /products
```

##### Get Products by Organization
```
GET /organizations/:organizationId/products
```

##### Update Product
```
PUT /products/:id
```

##### Delete Product
```
DELETE /products/:id
```

## Response Format

Success Response:
```json
{
  "success": true,
  "data": {...},
  "meta": {
    "timestamp": "2025-06-21T10:00:00Z",
    "version": "1.0"
  }
}
```

Error Response:
```json
{
  "success": false,
  "error": {
    "code": "ERROR_CODE",
    "message": "Error message",
    "details": {...}
  },
  "meta": {
    "timestamp": "2025-06-21T10:00:00Z",
    "version": "1.0"
  }
}
```

## Rate Limiting
API is rate limited to 100 requests per minute per IP address.
