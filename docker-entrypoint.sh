#!/bin/sh
set -e

echo "🗄️  Running database migrations..."
migrate -path ./migrations -database "$DATABASE_URL" up || echo "Migrations already applied"

echo "🚀 Starting Go server..."
exec ./server
