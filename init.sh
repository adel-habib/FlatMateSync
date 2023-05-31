#!/bin/sh 

set -e 

DB_USER=${DATABASE_USER}
DB_PASS=${DATABASE_PASSWORD}

DB_SOURCE="postgres://${DB_USER}:${DB_PASS}@${DATABASE_HOST}:${DATABASE_PORT}/${DATABASE_NAME}?sslmode=disable"
DB_LOG="postgres://${DATABASE_HOST}:${DATABASE_PORT}/${DATABASE_NAME}"
echo "running db migrations. Database Source: ${DB_LOG}"
/app/migrate -path /app/migrations -database "${DB_SOURCE}" --verbose up

exec "$@"