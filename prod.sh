#!/bin/sh
set -eu

# Production orchestration:
# 1) build images
# 2) start db
# 3) run one-shot incremental migration runner
# 4) recreate app

echo "==> Building production images..."
docker compose --profile prod build prod-postgres prod-main-server

echo "==> Starting prod-postgres..."
docker compose --profile prod up -d --wait prod-postgres

echo "==> Running incremental migrations..."
docker compose --profile prod run --rm prod-db-migrate

echo "==> Recreating prod-main-server..."
docker compose --profile prod up -d --force-recreate prod-main-server

echo "==> Production stack is up."
docker compose --profile prod ps
