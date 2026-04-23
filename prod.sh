#!/bin/sh
set -eu

# Production orchestration: build and start prod-postgres, wait for health,
# then build and start prod-main-server.

echo "==> Building and starting prod-postgres..."
docker compose --profile prod up -d --build --wait prod-postgres

echo "==> Building and starting prod-main-server..."
docker compose --profile prod up -d --build prod-main-server

echo "==> Production stack is up."
docker compose --profile prod ps
