#!/bin/sh

set -e

echo "run db migration"
migrate -path /app/sql/migrations -database "$ENV_DB_SOURCE" -verbose up

echo "start the app"
exec "$@"