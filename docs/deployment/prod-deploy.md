# Production Deployment (Scheme B)

This repository uses image bundle delivery (`tar.gz`) and remote `docker load` for production.

## Goals

- No registry push/pull in CI/CD.
- No build/compile on target server.
- Separate database runtime from schema migration execution.
- Support incremental schema migration during deployment.

## Components

- `docker/compose.prod.yml`
  - `prod-postgres`: long-running DB container.
  - `prod-db-migrate`: one-shot migration runner.
  - `prod-main-server`: long-running app container.
- `docker/prod-database-migrate.sh`
  - tracks applied versions in `public.schema_migration_table`.
  - supports baseline for pre-existing schema.
  - validates checksum consistency for applied versions.

## Server directory layout

Suggested target path: `/opt/poprako-s`

- `releases/<sha>/` release bundles
- `shared/.env` runtime env file
- `shared/docker/compose.prod.yml` compose file
- `shared/migrations/` SQL migrations

## GitHub Actions workflow

Workflow file: `.github/workflows/deploy-prod.yml`

For manual deployment outside CI, use:

- `scripts/package-release.sh` to build, export, package, and upload the release bundle.
- `scripts/remote-switch-release.sh` on the server to load images, run migrations, and switch the app container.

Build and package on runner:

- build app/db images
- `docker save | gzip` image bundles
- package migrations + migration runner + prod compose as `tar.gz`

Deploy to server:

- upload bundles via SSH key auth (`SSH_PRIVATE_KEY`)
- `docker load` both images
- extract migration bundle into `shared/`
- write `shared/.env`
- run `docker compose up -d --wait prod-postgres`
- run `docker compose run --rm prod-db-migrate`
- run `docker compose up -d --force-recreate prod-main-server`

## Required GitHub Secrets

- `SSH_PRIVATE_KEY`
- `SERVER_HOST`
- `SERVER_USER`
- `POSTGRES_PASSWORD`
- `DATABASE_PASSWORD`
- `JWT_SECRET`
- `JWT_EXPIRATION_HOURS`
- `OSS_PLATFORM`
- optional OSS provider vars (`R2_*`, `ALIYUN_*`)

## Migration behavior

- New database (empty volume):
  - `prod-db-migrate` applies all `migrations/*.up.sql` in filename order.
- Existing database without tracking table:
  - if `public.t_assignment_invitation` exists, treat as initialized schema.
  - mark existing migration files as already applied (baseline), then only apply new ones.
- Existing database with tracking table:
  - apply only unapplied migration files.
  - fail if an already applied version has checksum drift.

## Hot schema update flow

For a release with new migration file(s):

1. DB container stays running.
2. One-shot migration container executes incremental SQL.
3. App container is recreated afterward.

This ensures schema changes are applied without resetting DB volume.
