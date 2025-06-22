# LeCritique API Documentation

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
    "restaurant": {...},
    "qr_code": {
      "id": "uuid",
      "label": "Table 1",
      "type": "table"
    }
  }
}
```

#### Get Restaurant Menu
```
GET /public/restaurant/:id/menu
```

#### Submit Feedback
```
POST /public/feedback
```
Body:
```json
{
  "qr_code_id": "uuid",
  "dish_id": "uuid",
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
  "company_name": "My Restaurant"
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

#### Restaurants

##### Create Restaurant
```
POST /restaurants
```

##### List Restaurants
```
GET /restaurants
```

##### Get Restaurant
```
GET /restaurants/:id
```

##### Update Restaurant
```
PUT /restaurants/:id
```

##### Delete Restaurant
```
DELETE /restaurants/:id
```

#### Dishes

##### Create Dish
```
POST /dishes
```

##### Get Dishes by Restaurant
```
GET /restaurants/:restaurantId/dishes
```

##### Update Dish
```
PUT /dishes/:id
```

##### Delete Dish
```
DELETE /dishes/:id
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
