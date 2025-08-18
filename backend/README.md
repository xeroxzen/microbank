# Microbank Backend Services

This directory contains the Go microservices for the Microbank platform.

## Architecture

The backend consists of two main services:

- **Client Service** (`services/client-service/`): Handles user authentication, registration, and profile management
- **Banking Service** (`services/banking-service/`): Manages bank accounts, transactions, and balances

## Quick Start

### Prerequisites

- Go 1.21+
- PostgreSQL 15+
- Docker & Docker Compose (optional)

### 1. Environment Setup

Copy the example environment files and configure them:

```bash
# Client Service
cd services/client-service
cp env.example .env
# Edit .env with your database credentials

# Banking Service
cd ../banking-service
cp env.example .env
# Edit .env with your database credentials
```

### 2. Install Dependencies

```bash
# Client Service
cd services/client-service
go mod tidy

# Banking Service
cd ../banking-service
go mod tidy
```

### 3. Start Services

```bash
# Client Service (Terminal 1)
cd services/client-service
go run cmd/main.go

# Banking Service (Terminal 2)
cd services/banking-service
go run cmd/main.go
```

## API Documentation

### Client Service API

#### Authentication Endpoints

**POST** `/api/v1/auth/register`

```json
{
  "email": "andile.mbele@example.com",
  "name": "Andile Mbele",
  "password": "securepassword123"
}
```

**POST** `/api/v1/auth/login`

```json
{
  "email": "andile.mbele@example.com",
  "password": "securepassword123"
}
```

**POST** `/api/v1/auth/refresh`

```json
{
  "refresh_token": "your-refresh-token"
}
```

**GET** `/api/v1/auth/validate` _(Protected)_

#### Profile Endpoints

**GET** `/api/v1/profile` _(Protected)_
**PUT** `/api/v1/profile` _(Protected)_

```json
{
  "name": "Andile Mbele"
}
```

#### Admin Endpoints

**GET** `/api/v1/admin/clients` _(Admin)_
**PUT** `/api/v1/admin/clients/{id}/blacklist` _(Admin)_
**DELETE** `/api/v1/admin/clients/{id}/blacklist` _(Admin)_

### Banking Service API

#### Account Endpoints

**GET** `/api/v1/account/balance` _(Protected)_
**GET** `/api/v1/account/transactions` _(Protected)_

#### Transaction Endpoints

**POST** `/api/v1/transactions/deposit` _(Protected)_

```json
{
  "amount": 100.0,
  "description": "Salary deposit"
}
```

**POST** `/api/v1/transactions/withdraw` _(Protected)_

```json
{
  "amount": 50.0,
  "description": "ATM withdrawal"
}
```

**GET** `/api/v1/transactions/{id}` _(Protected)_

## Authentication

All protected endpoints require a valid JWT token in the Authorization header:

```
Authorization: Bearer <your-jwt-token>
```

### JWT Token Structure

```json
{
  "user_id": "uuid",
  "email": "user@example.com",
  "name": "User Name",
  "is_admin": false,
  "is_blacklisted": false,
  "exp": 1625097600,
  "iat": 1625011200,
  "type": "access"
}
```

## 🗄️ Database Schema

### Client Service Database

#### Users Table

```sql
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
```

#### Refresh Tokens Table

```sql
CREATE TABLE refresh_tokens (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    token_hash VARCHAR(255) NOT NULL,
    expires_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

### Banking Service Database

#### Accounts Table

```sql
CREATE TABLE accounts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID UNIQUE NOT NULL,
    balance DECIMAL(15,2) DEFAULT 0.00,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

#### Transactions Table

```sql
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

## Testing

### Run Tests

```bash
# Client Service
cd services/client-service
go test ./...

# Banking Service
cd ../banking-service
go test ./...

# Run with coverage
go test -cover ./...
```

### Test Structure

- **Unit Tests**: Test individual functions and methods
- **Integration Tests**: Test database operations and API endpoints
- **Mock Tests**: Test service layer with mocked repositories

## Development

### Project Structure

```
services/
├── client-service/
│   ├── cmd/           # Application entry point
│   ├── internal/      # Private application code
│   │   ├── handlers/  # HTTP request handlers
│   │   ├── middleware/# HTTP middleware
│   │   ├── models/    # Data models
│   │   ├── repository/# Database operations
│   │   └── services/  # Business logic
│   ├── go.mod         # Go module file
│   └── env.example    # Environment variables template
└── banking-service/
    ├── cmd/           # Application entry point
    ├── internal/      # Private application code
    │   ├── handlers/  # HTTP request handlers
    │   ├── middleware/# HTTP middleware
    │   ├── models/    # Data models
    │   ├── repository/# Database operations
    │   └── services/  # Business logic
    ├── go.mod         # Go module file
    └── env.example    # Environment variables template
```

### Adding New Features

1. **Models**: Define data structures in `internal/models/`
2. **Repository**: Implement database operations in `internal/repository/`
3. **Service**: Add business logic in `internal/services/`
4. **Handler**: Create HTTP endpoints in `internal/handlers/`
5. **Middleware**: Add request processing in `internal/middleware/`

### Code Style

- Use Go modules for dependency management
- Follow Go naming conventions
- Write comprehensive error handling
- Include proper logging and monitoring
- Add unit tests for new functionality

## Deployment

### Docker

```bash
# Build and run with Docker Compose
docker-compose up -d

# Build individual services
docker build -t client-service ./services/client-service
docker build -t banking-service ./services/banking-service
```

### Production Considerations

- Set `GIN_MODE=release`
- Using strong JWT secrets
- Enabling database SSL
- Configure proper CORS policies
- Set up monitoring and logging
- Use environment-specific configurations

## 📊 Monitoring

### Health Checks

- **Client Service**: `GET /health`
- **Banking Service**: `GET /health`

### Logging

- Structured logging with request IDs
- Error tracking and monitoring
- Performance metrics collection

### Metrics

- Request/response times
- Error rates
- Database connection health
- Transaction success rates

## Security

### Implemented Security Features

- JWT token authentication
- Password hashing with bcrypt
- CORS configuration
- Input validation and sanitization
- SQL injection prevention
- Rate limiting (configurable)
