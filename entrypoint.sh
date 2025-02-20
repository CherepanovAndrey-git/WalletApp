#!/bin/sh

set -e

echo "Waiting for the database to be ready..."
/app/wait-for-db.sh db

echo "Running migrations..."
/app/goose -dir /app/sql postgres "$DB_URL" up

echo "Starting application..."
/app/main