#!/bin/sh
set -e  # Exit immediately if any command fails

echo "Running database migrations..."
/usr/local/bin/dbmate up --migrations-dir="/app/db/migrations"

echo "Starting application..."
exec /app/main  # Use exec to replace shell with your Go app
