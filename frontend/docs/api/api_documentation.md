# API Documentation

## Base URL
```
Production: https://api.bishalpujasewa.com/api/v1
Development: http://localhost:8080/api/v1
```

## Authentication

All authenticated endpoints require a Bearer token in the Authorization header:
```
Authorization: Bearer <access_token>
```

### Token Types
- **Access Token**: Expires in 1 hour
- **Refresh Token**: Expires in 7 days

## Endpoints

### Authentication

#### Register
```http
POST /auth/register
Content-Type: application/json

{
  "email": "user@example.com",
  "password": "SecurePass123!",
  "full_name": "Ram Sharma",
  "phone": "9841234567"
}
```

#### Login
```http
POST /auth/login
Content-Type: application/json

{
  "email": "user@example.com",
  "password": "SecurePass123!"
}
```

#### Refresh Token
```http
POST /auth/refresh
Content-Type: application/json

{
  "refresh_token": "eyJhbGciOiJIUzI1NiIs..."
}
```

### User Profile

#### Get Profile
```http
GET /profile
Authorization: Bearer <token>
```

#### Update Profile
```http
PUT /profile
Authorization: Bearer <token>
Content-Type: application/json

{
  "full_name": "Updated Name",
  "phone": "9847654321"
}
```

### Pandits

#### Register as Pandit
```http
POST /pandit/register
Authorization: Bearer <token>
Content-Type: application/json

{
  "bio": "Experienced Vedic pandit",
  "experience_years": 15,
  "specialization": "Vivah, Puja, Griha Pravesh",
  "languages": ["Nepali", "Sanskrit", "Hindi"],
  "base_price": 1500,
  "service_area": "Kathmandu Valley"
}
```

#### List Pandits
```http
GET /pandits?page=1&limit=20
```

### Bookings

#### Create Booking
```http
POST /bookings
Authorization: Bearer <token>
Content-Type: application/json

{
  "pandit_id": "uuid",
  "ritual_id": "uuid",
  "scheduled_date": "2025-01-15",
  "start_time": "08:00",
  "end_time": "10:00",
  "address": "Kathmandu, Nepal",
  "special_notes": "Please bring necessary items"
}
```

### Payments

#### Initiate Payment
```http
POST /payments/initiate
Authorization: Bearer <token>
Content-Type: application/json

{
  "booking_id": "uuid",
  "gateway": "esewa|khalti|mock"
}
```

#### Verify Payment
```http
POST /payments/verify
Authorization: Bearer <token>
Content-Type: application/json

{
  "transaction_id": "uuid",
  "status": "completed"
}
```

## Response Format

### Success
```json
{
  "success": true,
  "message": "Operation successful",
  "data": {}
}
```

### Error
```json
{
  "success": false,
  "message": "Error description",
  "error": "Detailed error message"
}
```

### Paginated
```json
{
  "success": true,
  "message": "List retrieved",
  "data": [],
  "total": 100,
  "page": 1,
  "limit": 20
}
```

## Rate Limiting

- 30 requests per second per IP
- Burst up to 50 requests
- Returns 429 Too Many Requests when exceeded

## Security Headers

All responses include:
- X-Frame-Options: DENY
- X-Content-Type-Options: nosniff
- X-XSS-Protection: 1; mode=block
- Strict-Transport-Security: max-age=31536000
- Content-Security-Policy: restricted
