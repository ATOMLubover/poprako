set shell := ["sh", "-c"]

default:
    just --list

swag:
    swag init

check paths="./...":
    golangci-lint run {{ paths }}

fmt:
    golangci-lint fmt ./...

cloc:
    cloc internal/

# Start postgres and apply all up-migrations
setup-database:
    docker compose up -d postgres --wait
    docker compose --profile setup run --rm database-setup

# Format, regenerate Swagger docs, then build & run backend container
# (schema setup runs first and is safe to repeat)
dev: fmt swag
    just setup-database
    docker compose --profile app up --build backend

# Run seed scripts against the running backend
seed:
    cd tests/seed && bun run seed

# Bring up the full production stack
product:
    docker compose -f docker/docker-compose.yml up

build:
    GOOS=linux GOARCH=amd64 go build -o output/poprako-s

mgr-add script-name:
    sqlx migrate add -r {{script-name}}

mgr-run:
    sqlx migrate run

mgr-rvt mode="step":
    {{if mode == "all" {
        "sqlx migrate revert --target-version 0"
    } else {
        "sqlx migrate revert"
    }}}

psql:
  psql -U devuser -d poprako-s-db
