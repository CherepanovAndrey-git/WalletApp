#!/bin/sh

set -e

echo "Running migrations..."
goose -dir /app/sql postgres "$DB_URL" up

echo "Starting application..."
exec ./main