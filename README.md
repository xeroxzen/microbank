# Microbank - High Level Design (HLD)

## 1. System Overview

Microbank is a full-stack microservice banking platform consisting of two core backend services built in Go, with a React/Next.js frontend providing both client and admin interfaces.

### 1.1 Core Services

- **Client Service**: Handles user registration, authentication, profile management, and blacklisting
- **Banking Service**: Manages deposits, withdrawals, balance tracking, and transaction history
- **Frontend Applications**: Client dashboard and admin panel

## 2. System Architecture

### 2.1 Architecture Diagram

```
┌─────────────────┐    ┌─────────────────┐
│   Client App    │    │   Admin Panel   │
│   (Next.js)     │    │   (Next.js)     │
└─────────┬───────┘    └─────────┬───────┘
          │                      │
          └──────────┬───────────┘
                     │
          ┌─────────────────┐
          │   API Gateway   │
          │   (Optional)    │
          └─────────┬───────┘
                    │
        ┌───────────┴───────────┐
        │                       │
┌───────▼───────┐       ┌───────▼───────┐
│ Client Service │       │Banking Service│
│    (Go/Gin)    │◄──────┤   (Go/Gin)    │
└───────┬───────┘       └───────┬───────┘
        │                       │
┌───────▼───────┐       ┌───────▼───────┐
│  PostgreSQL    │       │  PostgreSQL   │
│ (Client DB)    │       │ (Banking DB)  │
└────────────────┘       └───────────────┘
```

### 2.2 Technology Stack

| Component         | Technology                      |
| ----------------- | ------------------------------- |
| Frontend          | Next.js 14, React, Tailwind CSS |
| Backend Services  | Go 1.21+, Gin Framework         |
| Databases         | PostgreSQL                      |
| Authentication    | JWT (JSON Web Tokens)           |
| Containerization  | Docker & Docker Compose         |
| API Documentation | Swagger/OpenAPI                 |

## 3. Service Design

### 3.1 Client Service

#### 3.1.1 Responsibilities

- User registration and authentication
- JWT token generation and validation
- User profile management
- Blacklist management (admin only)
- User status verification

#### 3.1.2 API Endpoints

```go
// Public endpoints
POST   /api/v1/auth/register
POST   /api/v1/auth/login
POST   /api/v1/auth/refresh

// Protected endpoints
GET    /api/v1/profile
PUT    /api/v1/profile
GET    /api/v1/auth/validate

// Admin endpoints
GET    /api/v1/admin/clients
PUT    /api/v1/admin/clients/{id}/blacklist
DELETE /api/v1/admin/clients/{id}/blacklist
```

#### 3.1.3 Database Schema

```sql
-- Users table
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email VARCHAR(255) UNIQUE NOT NULL,
    name VARCHAR(255) NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    is_blacklisted BOOLEAN DEFAULT FALSE,
    is_admin BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Refresh tokens table
CREATE TABLE refresh_tokens (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    token_hash VARCHAR(255) NOT NULL,
    expires_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

### 3.2 Banking Service

#### 3.2.1 Responsibilities

- Account balance management
- Deposit and withdrawal processing
- Transaction history tracking
- Overdraft prevention
- Blacklist enforcement

#### 3.2.2 API Endpoints

```go
// Account endpoints
GET    /api/v1/account/balance
GET    /api/v1/account/transactions

// Transaction endpoints
POST   /api/v1/transactions/deposit
POST   /api/v1/transactions/withdraw
GET    /api/v1/transactions/{id}
```

#### 3.2.3 Database Schema

```sql
-- Accounts table
CREATE TABLE accounts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID UNIQUE NOT NULL,
    balance DECIMAL(15,2) DEFAULT 0.00,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Transactions table
CREATE TABLE transactions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    account_id UUID REFERENCES accounts(id) ON DELETE CASCADE,
    user_id UUID NOT NULL,
    type VARCHAR(20) NOT NULL CHECK (type IN ('deposit', 'withdrawal')),
    amount DECIMAL(15,2) NOT NULL CHECK (amount > 0),
    balance_before DECIMAL(15,2) NOT NULL,
    balance_after DECIMAL(15,2) NOT NULL,
    description TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

## 4. Security Design

### 4.1 Authentication Flow

1. User registers/logs in via Client Service
2. Client Service validates credentials and generates JWT
3. JWT contains user ID, email, blacklist status, and expiry
4. All subsequent requests include JWT in Authorization header
5. Services validate JWT before processing requests

### 4.2 JWT Token Structure

```json
{
  "sub": "user-uuid",
  "email": "user@example.com",
  "name": "User Name",
  "is_admin": false,
  "is_blacklisted": false,
  "exp": 1625097600,
  "iat": 1625011200
}
```

### 4.3 Security Measures

- Password hashing using bcrypt
- JWT tokens with short expiry (15 minutes)
- Refresh token mechanism
- Rate limiting on sensitive endpoints
- Input validation and sanitization
- CORS configuration
- HTTPS enforcement in production

## 5. Inter-Service Communication

### 5.1 Service-to-Service Authentication

Banking Service validates user status by:

1. Decoding and validating JWT signature
2. Checking blacklist status in token
3. Optional: Real-time validation with Client Service for critical operations

### 5.2 Communication Pattern

```go
// Banking Service validates user before processing
func (h *TransactionHandler) validateUser(userID string) error {
    // Decode JWT and check blacklist status
    // For critical operations, optionally call Client Service
    return nil
}
```

## 6. Data Flow

### 6.1 User Registration Flow

```
Client App → Client Service → Database → JWT Response → Client App
```

### 6.2 Banking Transaction Flow

```
Client App → Banking Service → Validate JWT → Check Blacklist →
Process Transaction → Update Balance → Record Transaction → Response
```

### 6.3 Admin Blacklist Flow

```
Admin Panel → Client Service → Update User Status →
Optional: Notify Banking Service → Response
```

## 7. Error Handling

### 7.1 Error Response Format

```json
{
  "error": {
    "code": "INSUFFICIENT_FUNDS",
    "message": "Account balance insufficient for withdrawal",
    "details": {
      "requested_amount": 1000.0,
      "current_balance": 500.0
    }
  }
}
```

### 7.2 Common Error Scenarios

- Invalid credentials
- Blacklisted user attempts
- Insufficient funds
- Invalid JWT tokens
- Service unavailability

## 8. Deployment Architecture

### 8.1 Docker Compose Structure

```yaml
version: "3.8"
services:
  client-service:
    build: ./services/client-service
    ports:
      - "8081:8080"
    environment:
      - DB_HOST=client-db
      - JWT_SECRET=your-secret-key
    depends_on:
      - client-db

  banking-service:
    build: ./services/banking-service
    ports:
      - "8082:8080"
    environment:
      - DB_HOST=banking-db
      - JWT_SECRET=your-secret-key
    depends_on:
      - banking-db

  client-app:
    build: ./client
    ports:
      - "3000:3000"
    environment:
      - NEXT_PUBLIC_CLIENT_SERVICE_URL=http://localhost:8081
      - NEXT_PUBLIC_BANKING_SERVICE_URL=http://localhost:8082

  client-db:
    image: postgres:15
    environment:
      POSTGRES_DB: client_service
      POSTGRES_USER: postgres
      POSTGRES_PASSWORD: password

  banking-db:
    image: postgres:15
    environment:
      POSTGRES_DB: banking_service
      POSTGRES_USER: postgres
      POSTGRES_PASSWORD: password
```

## 9. Testing Strategy

### 9.1 Unit Tests

- Service layer business logic
- Repository layer database operations
- Handler layer HTTP request/response
- Utility functions

### 9.2 Integration Tests

- API endpoint testing
- Database integration
- JWT authentication flow
- Service-to-service communication

### 9.3 Test Structure (Go)

```go
func TestDepositHandler(t *testing.T) {
    tests := []struct {
        name           string
        userID         string
        amount         float64
        initialBalance float64
        expectedError  string
    }{
        {
            name:           "successful deposit",
            userID:         "user-123",
            amount:         100.00,
            initialBalance: 500.00,
            expectedError:  "",
        },
        // Additional test cases...
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Test implementation
        })
    }
}
```

## 10. Monitoring and Logging

### 10.1 Logging Structure

```go
// Structured logging with context
log.WithFields(log.Fields{
    "user_id":        userID,
    "transaction_id": txnID,
    "amount":         amount,
    "type":           "deposit",
}).Info("Transaction processed successfully")
```

### 10.2 Metrics to Track

- Transaction success/failure rates
- Response times
- Authentication failures
- Blacklist violations
- Database connection health

## 11. Bonus Features Implementation

### 11.1 Message Queue for Blacklisting

- Use Redis or RabbitMQ for asynchronous blacklist updates
- Client Service publishes blacklist events
- Banking Service subscribes and updates local cache

### 11.2 Rate Limiting

```go
// Using golang.org/x/time/rate
var limiter = rate.NewLimiter(rate.Limit(10), 10) // 10 requests per second

func rateLimitMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        if !limiter.Allow() {
            c.JSON(429, gin.H{"error": "Rate limit exceeded"})
            c.Abort()
            return
        }
        c.Next()
    }
}
```

## 12. Development Setup

### 12.1 Project Structure

```
microbank/
├── client/                    # Next.js frontend
│   ├── components/
│   ├── pages/
│   └── package.json
├── services/
│   ├── client-service/        # Go client service
│   │   ├── cmd/
│   │   ├── internal/
│   │   ├── pkg/
│   │   └── go.mod
│   └── banking-service/       # Go banking service
│       ├── cmd/
│       ├── internal/
│       ├── pkg/
│       └── go.mod
├── docker-compose.yml
└── README.md
```

### 12.1 Getting Started

1. Clone repository
2. Run `docker-compose up` to start all services
3. Access client app at http://localhost:3000
4. API documentation available at service endpoints with `/swagger`

## Port Configuration

### Local Development

- **Client Service**: http://localhost:8081 (authentication, user management)
- **Banking Service**: http://localhost:8080 (accounts, transactions)
- **Client Dashboard**: http://localhost:3000
- **Admin Panel**: http://localhost:3001

### Docker Environment

- **Client Service**: http://localhost:8082 (external) → 8080 (internal)
- **Banking Service**: http://localhost:8081 (external) → 8080 (internal)
- **Client Dashboard**: http://localhost:3000
- **Admin Panel**: http://localhost:3001

### Environment Variables for Frontend

```bash
# Local Development
NEXT_PUBLIC_CLIENT_SERVICE_URL=http://localhost:8081
NEXT_PUBLIC_BANKING_SERVICE_URL=http://localhost:8080

# Docker Environment
NEXT_PUBLIC_CLIENT_SERVICE_URL=http://localhost:8082
NEXT_PUBLIC_BANKING_SERVICE_URL=http://localhost:8081
```
