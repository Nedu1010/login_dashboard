# 🔐 Authentication Service - Production-Grade Go + React

A modern, secure authentication service implementing industry best practices with a beautiful React frontend and robust Go backend.

![Login Page Design](/home/user/.gemini/antigravity/brain/5718bfa2-5950-4634-ae54-dadc3d8c4746/login_page_design_1770210846388.png)

## ✨ Features

### 🛡️ Security

- **JWT-based authentication** with short-lived access tokens (5 min)
- **Refresh token rotation** for enhanced security
- **HTTP-only secure cookies** (JavaScript cannot access access)
- **CSRF protection** using double-submit pattern
- **bcrypt password hashing** (cost 12)
- **Rate limiting** on authentication endpoints
- **SQL injection prevention** via prepared statements

### 🎨 Modern UI

- **Glassmorphism design** with dark theme
- **Smooth animations** and transitions
- **Password strength indicator**
- **Responsive** mobile-first design
- **Real-time validation** and error handling

### 🏗️ Architecture

- **Clean architecture** with separated layers (domain, repository, service, handler)
- **PostgreSQL** for data persistence
- **Graceful shutdown** handling
- **Structured logging** with request tracing -**Health check** endpoint

## 📚 Tech Stack

### Backend

- Go 1.21+
- Gin Web Framework
- PostgreSQL + pgx driver
- JWT (golang-jwt/jwt)
- Viper (configuration)
- golang-migrate (migrations)

### Frontend

- React 18 + TypeScript
- Vite (build tool)
- React Router v6
- Axios (HTTP client)
- Custom CSS (glassmorphism)

## 🚀 Quick Start

### Prerequisites

- Go 1.21 or higher
- Node.js 18+ and npm
- PostgreSQL 15+
- Docker (optional, for local database)

### 1. Database Setup

**Option A: Docker**

```bash
make docker-up
```

**Option B: Local PostgreSQL**

```SQL
CREATE DATABASE auth_db;
```

### 2. Backend Setup

```bash
# Install dependencies
go mod download

# Copy environment file
cp .env.example .env

# Generate a secure JWT secret
echo "JWT_SECRET=$(openssl rand -base64 64)" >> .env

# Run database migrations
make migrate-up

# Start the server
make run
```

The backend will be running at `http://localhost:8080`

### 3. Frontend Setup

```bash
# Navigate to web directory
cd web

# Install dependencies
npm install

# Start dev server
npm run dev
```

The frontend will be running at `http://localhost:3000`

## 📝 API Endpoints

### Public Endpoints

| Method | Endpoint             | Description                   |
| :----- | :------------------- | :---------------------------- |
| POST   | `/api/auth/register` | Register new user             |
| POST   | `/api/auth/login`    | Login with credentials        |
| POST   | `/api/auth/refresh`  | Refresh access token          |
| POST   | `/api/auth/logout`   | Logout (revoke refresh token) |
| GET    | `/health`            | Health check                  |

### Protected Endpoints

| Method | Endpoint       | Description      |
| :----- | :------------- | :--------------- |
| GET    | `/api/user/me` | Get current user |

## 🔑 Authentication Flow

1. **Register**: User creates account → Password hashed → User stored in DB
2. **Login**: Credentials validated → Access + Refresh tokens generated → Tokens set as HTTP-only cookies
3. **Access Protected Route**: Browser sends cookies automatically → Middleware validates access token
4. **Token Expired**: Frontend intercepts 401 → Calls `/auth/refresh` → New tokens issued
5. **Logout**: Refresh token revoked in DB → All cookies cleared

## 🍪 Cookie Strategy

| Cookie          | HttpOnly | Secure | SameSite | Expiry | Purpose            |
| :-------------- | :------- | :----- | :------- | :----- | :----------------- |
| `access_token`  | ✅ Yes   | ✅ Yes | Strict   | 5 min  | Short-lived access |
| `refresh_token` | ✅ Yes   | ✅ Yes | Strict   | 7 days | Token renewal      |
| `csrf_token`    | ❌ No    | ✅ Yes | Strict   | 24 hrs | CSRF protection    |

## 🧪 Testing

### Backend Tests

```bash
# Run all tests
make test

# With coverage
make test-coverage
```

### Manual Testing

**1. Register a User**

```bash
curl -X POST http://localhost:8080/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"SecureP@ss123"}'
```

**2. Login**

```bash
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"SecureP@ss123"}' \
  -c cookies.txt
```

**3. Access Protected Route**

```bash
curl -X GET http://localhost:8080/api/user/me \
  -b cookies.txt
```

## 🗂️ Project Structure

```
/data/go_p/login_flow/
├── cmd/server/          # Application entry point
├── internal/            # Private application code
│   ├── config/         # Configuration
│   ├── domain/         # Business models
│   ├── repository/     # Data access layer
│   ├── service/        # Business logic
│   ├── handler/        # HTTP handlers
│   ├── middleware/     # HTTP middleware
│   └── util/           # Utility functions
├── pkg/                # Public libraries
│   ├── jwt/           # JWT utilities
│   ├── crypto/        # Cryptography
│   └── validator/     # Input validation
├── migrations/         # Database migrations
├── web/                # React frontend
│   └── src/
│       ├── api/       # API client
│       ├── components/# React components
│       ├── pages/     # Page components
│       └── styles/    # CSS stylesheets
└── Makefile           # Common tasks
```

## 🛠️ Development Commands

### Backend

```bash
make run              # Start server
make build            # Build binary
make test             # Run tests
make migrate-up       # Apply migrations
make migrate-down     # Rollback migrations
make docker-up        # Start PostgreSQL
make docker-down      # Stop PostgreSQL
make clean            # Clean artifacts
```

### Frontend

```bash
npm run dev           # Start dev server
npm run build         # Production build
npm run preview       # Preview production build
```

## 🔒 Security Best Practices Implemented

✅ **Passwords**: Hashed with bcrypt (cost 12)  
✅ **Tokens**: Short-lived JWT + rotating refresh tokens  
✅ **Cookies**: HTTP-only, Secure, SameSite=Strict  
✅ **CSRF**: Double-submit cookie pattern  
✅ **SQL**: Prepared statements (no injection)  
✅ **XSS**: Input sanitization and validation  
✅ **Rate Limiting**: Login attempts limited  
✅ **Secrets**: Environment variables (never committed)

## 📦 Production Deployment

1. **Environment Variables**: Update `.env` with production values
2. **Database**: Use managed PostgreSQL (AWS RDS, Google Cloud SQL)
3. **Secrets**: Use proper secret management (AWS Secrets Manager, Vault)
4. **HTTPS**: Enable SSL/TLS (set `COOKIE_SECURE=true`)
5. **Monitoring**: Add application monitoring (Prometheus, Datadog)
6. **Logging**: Configure centralized logging

### Build for Production

```bash
# Backend
make build
./bin/server

# Frontend
cd web
npm run build
# Serve dist/ with nginx or CDN
```

## 🤝 Contributing

This is a reference implementation. Feel free to fork and customize for your needs!

## 📄 License

MIT

## 🙏 Acknowledgments

- Authentication strategy based on AWS/GCP best practices
- UI inspired by modern SaaS applications
- Architecture follows Go community standards
