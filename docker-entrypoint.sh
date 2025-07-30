#!/bin/sh
set -e

echo "Running database migrations..."
migrate -path=/app/migrations -database="postgres://$DB_USER:$DB_PASSWORD@$DB_HOST:$DB_PORT/$DB_NAME?sslmode=$DB_SSLMODE" up

if [ "$RUN_SEEDS" = "true" ]; then
    echo "Running database seeds..."
    ./seed-plans
    ./seed
fi

echo "Starting backend server..."
exec ./main serve --http=0.0.0.0:8080