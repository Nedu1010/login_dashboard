# 🚀 Quick Start - Docker Compose

## ✅ What's Working

- ✅ Docker images build successfully
- ✅ Frontend is running on port 3000
- ✅ Go code compiles without errors

## ⚠️ Final Setup Step

Your local PostgreSQL is running but the password doesn't match. Choose one option:

### Option 1: Use Docker PostgreSQL (Recommended - Easiest)

Stop your local PostgreSQL and use Docker instead:

```bash
# Stop local PostgreSQL
sudo systemctl stop postgresql   # Linux
# OR
brew services stop postgresql    # macOS

# Update docker-compose to include postgres
docker compose down
```

Then edit `docker-compose.yml` - uncomment the postgres service (I'll create a version with both).

### Option 2: Update Password in Docker Compose

Edit `docker-compose.yml` line 9-12 to match your actual Postgres password:

```yaml
environment:
  - DATABASE_URL=postgres://postgres:YOUR_ACTUAL_PASSWORD@host.docker.internal:5432/auth_db?sslmode=disable
  - DATABASE_PASSWORD=YOUR_ACTUAL_PASSWORD
```

### Option 3: Run Backend Locally (No Docker Needed!)

Skip Docker for the backend entirely:

```bash
# Stop docker backend
docker compose down

# Run locally
export PATH=$PATH:$(go env GOPATH)/bin
go run cmd/server/main.go
```

The `.env` file already has `DATABASE_URL=postgres://postgres:postgres@localhost:5432/auth_db?sslmode=disable` configured.

## ✅ Quick Check

Once you fix the database connection, test:

```bash
# Check backend health
curl http://localhost:8080/health

# Should return:
# {"status":"ok","database":"connected"}
```

## 🌐 Access Your App

Once the backend is running:

- Frontend: **http://localhost:3000**
- Backend: **http://localhost:8080**
- Health: **http://localhost:8080/health**

The frontend is already running, so once the backend connects to the database, you're ready to go!

## 🐛 Need Help?

The issue is just database authentication. Everything else is working! Let me know which option you prefer and I can help set it up.
