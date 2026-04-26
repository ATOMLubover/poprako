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

# Start dev-postgres and apply all up-migrations (safe to re-run)
setup:
    docker compose --profile db up -d --wait dev-postgres
    docker compose --profile db --profile setup run --rm dev-database-setup

# Format, regenerate Swagger docs, then build & run dev-main-server
dev: fmt swag
    docker compose --profile app up --build dev-main-server

# Tear down all containers and volumes (full clean slate)
reset:
    docker compose --profile db --profile setup --profile app --profile prod down -v

# Run seed scripts against the running backend
seed:
    cd tests/seed && bun run seed

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
  psql -U devuser -d poprako_s_db
