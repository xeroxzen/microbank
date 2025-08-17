# Microbank Setup Guide

## 🚀 Quick Start

Get Microbank running locally in minutes with Docker Compose!

### Prerequisites

- Docker & Docker Compose
- Go 1.21+ (for local development)
- Node.js 18+ (for local development)
- Git

### 1. Clone and Setup

```bash
git clone <your-repo-url>
cd microbank
```

### 2. Environment Configuration

Create `.env` files for each service:

**Backend Services:**

```bash
# backend/services/client-service/.env
DB_HOST=localhost
DB_PORT=5433
DB_USER=postgres
DB_PASSWORD=password
DB_NAME=client_service
JWT_SECRET=your-super-secret-jwt-key-change-in-production
GIN_MODE=debug

# backend/services/banking-service/.env
DB_HOST=localhost
DB_PORT=5434
DB_USER=postgres
DB_PASSWORD=password
DB_NAME=banking_service
JWT_SECRET=your-super-secret-jwt-key-change-in-production
GIN_MODE=debug
```

**Frontend Apps:**

```bash
# frontend/apps/client-dashboard/.env.local
NEXT_PUBLIC_CLIENT_SERVICE_URL=http://localhost:8081
NEXT_PUBLIC_BANKING_SERVICE_URL=http://localhost:8082

# frontend/apps/admin-panel/.env.local
NEXT_PUBLIC_CLIENT_SERVICE_URL=http://localhost:8081
NEXT_PUBLIC_BANKING_SERVICE_URL=http://localhost:8082
```

### 3. Start All Services

```bash
# Start everything with Docker Compose
docker-compose up -d

# View logs
docker-compose logs -f

# Stop all services
docker-compose down
```

### 4. Access Your Applications

- **Client Dashboard**: http://localhost:3000
- **Admin Panel**: http://localhost:3001
- **Client Service API**: http://localhost:8081
- **Banking Service API**: http://localhost:8082
- **pgAdmin**: http://localhost:5050 (admin@microbank.com / admin123)

## 🛠️ Development Setup

### Backend Development (Go)

#### Install Dependencies

```bash
cd backend/services/client-service
go mod tidy

cd ../banking-service
go mod tidy
```

#### Run Locally

```bash
# Client Service
cd backend/services/client-service
go run cmd/main.go

# Banking Service (in another terminal)
cd backend/services/banking-service
go run cmd/main.go
```

#### Database Migrations

```bash
# Connect to client database
psql -h localhost -p 5433 -U postgres -d client_service

# Connect to banking database
psql -h localhost -p 5434 -U postgres -d banking_service
```

### Frontend Development (Next.js)

#### Install Dependencies

```bash
# Install root dependencies
cd frontend
npm install

# Install app dependencies
cd apps/client-dashboard && npm install
cd ../admin-panel && npm install

# Install shared package dependencies
cd ../packages/ui && npm install
```

#### Run Locally

```bash
# From frontend root directory
npm run dev

# Or run individual apps
cd apps/client-dashboard && npm run dev
cd ../admin-panel && npm run dev
```

## 🧪 Testing

### Backend Tests

```bash
# Run all tests
cd backend/services/client-service
go test ./...

# Run with coverage
go test -cover ./...

# Run specific test
go test -v ./internal/handlers
```

### Frontend Tests

```bash
# Run all tests
cd frontend
npm test

# Run specific app tests
cd apps/client-dashboard && npm test
```

## 📊 Monitoring & Debugging

### Health Checks

- **Client Service**: http://localhost:8081/health
- **Banking Service**: http://localhost:8082/health
- **Client Dashboard**: http://localhost:3000/api/health
- **Admin Panel**: http://localhost:3001/api/health

### Database Management

1. Open pgAdmin at http://localhost:5050
2. Login with admin@microbank.com / admin123
3. Add servers:
   - **Client DB**: localhost:5433, postgres/password
   - **Banking DB**: localhost:5434, postgres/password

### Logs

```bash
# View all service logs
docker-compose logs -f

# View specific service logs
docker-compose logs -f client-service
docker-compose logs -f banking-service
docker-compose logs -f client-dashboard
```

## 🔧 Troubleshooting

### Common Issues

#### Port Already in Use

```bash
# Find process using port
lsof -i :3000
lsof -i :8081

# Kill process
kill -9 <PID>
```

#### Database Connection Issues

```bash
# Check if databases are running
docker-compose ps

# Restart database services
docker-compose restart client-db banking-db

# Check database logs
docker-compose logs client-db
```

#### Go Module Issues

```bash
# Clean module cache
go clean -modcache

# Download dependencies again
go mod download
```

#### Node.js Issues

```bash
# Clear npm cache
npm cache clean --force

# Remove node_modules and reinstall
rm -rf node_modules package-lock.json
npm install
```

### Performance Tuning

#### Database Optimization

```sql
-- Add indexes for better performance
CREATE INDEX CONCURRENTLY idx_users_email ON users(email);
CREATE INDEX CONCURRENTLY idx_transactions_user_id ON transactions(user_id);
CREATE INDEX CONCURRENTLY idx_transactions_created_at ON transactions(created_at);
```

#### Go Service Optimization

```bash
# Build with optimizations
go build -ldflags="-s -w" -o main ./cmd

# Run with profiling
go run -cpuprofile=cpu.prof cmd/main.go
```

## 🚀 Production Deployment

### Environment Variables

```bash
# Production environment variables
GIN_MODE=release
JWT_SECRET=<strong-random-secret>
DB_SSLMODE=require
DB_HOST=<production-db-host>
DB_PASSWORD=<production-db-password>
```

### Security Checklist

- [ ] Change default JWT secret
- [ ] Enable HTTPS
- [ ] Set up proper CORS
- [ ] Configure rate limiting
- [ ] Set up monitoring and alerting
- [ ] Enable database SSL
- [ ] Set up backup strategy

### Scaling Considerations

- Use load balancer for multiple service instances
- Implement Redis for session management
- Set up database read replicas
- Use CDN for static assets
- Implement proper logging aggregation

## 📚 API Documentation

### Swagger/OpenAPI

Once services are running, access API documentation:

- **Client Service**: http://localhost:8081/swagger/index.html
- **Banking Service**: http://localhost:8082/swagger/index.html

### Postman Collection

Import the provided Postman collection for testing all endpoints.

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests
5. Submit a pull request

## 📞 Support

- **Issues**: GitHub Issues
- **Documentation**: README.md
- **Architecture**: HLD in README.md

---

**Happy Coding! 🎉**
