# Plan: Replace Root Codebase with Refactor

## Summary

Replace the root (original) layered-architecture codebase with the refactor DDD-architecture codebase
from `refactor/`. Some directories are merged rather than replaced.

---

## Phase 1: Direct Replacements (refactor overwrites root)

### 1.1 `internal/` — Entire directory

**Action**: Delete root `internal/`, copy `refactor/internal/` in its place.

**Why**: The refactor is a complete rewrite of the application layer with a different package structure.
No files from the old `internal/` are needed.

**New packages introduced** (not in root):

- `internal/pkg/util/` — ID generation utilities
- `internal/api/http/middleware/` — separated auth + log middleware
- `internal/api/http/res/` — HTTP response helpers
- `internal/api/state/` — AppState for HTTP layer
- `internal/app/impl/` — concrete app implementations + logging decorators
- `internal/app/res/` — result types
- `internal/domain/model/aggr/` — aggregate roots
- `internal/domain/model/enum/` — domain enums
- `internal/domain/model/event/` — domain event types
- `internal/domain/model/query/` — query/filter models
- `internal/domain/ext/token/` — token parser interface
- `internal/domain/svc/res/` — service result types
- `internal/infra/ext/token/` — JWT token parser implementation
- `internal/event/` — event system core (moved from domain/event)

### 1.2 `main.go`

**Action**: Replace root `main.go` with `refactor/main.go`.

**Key differences in wiring**:
| Aspect | Root | Refactor |
|---|---|---|
| Config loading | `cfg.Load()` returns `(*AppCfg, error)` | `cfg.NewAppCfg()` returns `*AppCfg`, panics internally |
| Logger | `lgr.Init(appCfg)` | `lgr.New(appCfg)` + `lgr.SetGlobal(appLgr)` |
| DB connection | `repo_infra.NewGDB(appCfg.DB.DSN)` | `repo_infra.NewPgGdb(appCfg.Db)` + `repo_infra.ApplyCfg(gdb, appCfg.Db)` |
| Token handling | Inline in auth handler | `token_infra.NewJwtParser()` separate infra component |
| App constructors | `app.NewXxxApp(...)` | `app_impl.NewXxxLogApp(app_impl.NewXxxApp(...))` decorator pattern |
| Event system | `event_infra.NewEventBus()`, `eventBus.SubUnsafe(h)` | `event_infra.NewEvBus()`, `evBus.Sub(h)`, `evBus.Run()` |
| Worker | `ossWorker.Start(ctx)` via signal.NotifyContext | `go worker_infra.NewOssWorker(...).Run(runCx)` goroutine before server |
| State | `app_state.NewAppState(...)` | `state.NewAppState(...)` in api/state |

**New apps introduced**: `sysMailApp`, `assignmentInvApp` (split from chapterApp)

**Apps renamed**: `chapterExportApp` + `chapterImportApp` → `chapterPortApp`

### 1.3 `app_config.json`

**Action**: Replace with `refactor/app_config.json`.

**Schema change**:

```jsonc
// OLD (root)
{ "http_address": "127.0.0.1:8080", "auth": { "expiration_hours": 144 } }

// NEW (refactor)
{
  "environment": "development",          // was APP_ENVIRONMENT env var
  "database": {                          // was DATABASE_URL env var
    "max_open": 24, "max_idle": 4,
    "max_life": 84600, "max_idle_time": 60
  },
  "http": { "host": "localhost", "port": 8080 }  // was http_address
}
```

### 1.4 `go.mod` / `go.sum`

**Action**: Replace with refactor versions. The dependency graph is identical except root
has `testify` as an indirect dependency (used by root's test files which are removed).
Run `go mod tidy` after replacement.

### 1.5 `.golangci.yml`

**Action**: Copy `refactor/.golangci.yml` to root (root doesn't have one).

---

## Phase 2: Directory Merges

### 2.1 `docs/`

**Action**: Merge — keep all root docs, add refactor docs that don't conflict.

| Root only           | Refactor only          | Both (conflict) |
| ------------------- | ---------------------- | --------------- |
| `deployment/`       | `ddd-architecture.md`  | `swagger.json`  |
| `explain/`          | `perm-list.md`         | `swagger.yaml`  |
| `rfc/`              | `todos.md`             | `docs.go`       |
| `tests/`            | `plan/`, `plans/`      |                 |
| `export-sample.txt` | `progress/`, `review/` |                 |
| `import-sample.txt` |                        |                 |

**Resolution for conflicts**:

- `swagger.json` / `swagger.yaml`: Use refactor versions (they match the new API surface with sys_mail, assignment_invitation, chapter_port endpoints)
- `docs.go`: Use refactor version

### 2.2 `migrations/`

**Action**: Use refactor migrations as the new baseline. This is the highest-risk change.

**Table mapping (root → refactor)**:

| Root Table                 | Refactor Table                | Change                                          |
| -------------------------- | ----------------------------- | ----------------------------------------------- |
| `features_enable`          | `features_enable`             | Same                                            |
| `user_table`               | `user_table`                  | Same                                            |
| `team_table`               | `team_table`                  | Same                                            |
| `member_table`             | `member_table`                | Same                                            |
| `member_invitation_table`  | `member_invitation_table`     | Same                                            |
| `workset_table`            | `workset_table`               | Same                                            |
| `comic_table`              | `comic_table`                 | Same                                            |
| `chapter_table`            | `chapter_table`               | Same                                            |
| `page_table`               | `page_table`                  | Same                                            |
| `assignment_table`         | `assignment_table`            | Same                                            |
| `unit_table`               | `unit_table`                  | Same                                            |
| `user_stats_table`         | `user_stats_table`            | Same                                            |
| `oss_message_table`        | `oss_message_table`           | Same                                            |
| `notification_table`       | —                             | **Removed**                                     |
| `chapter_invitation_table` | —                             | **Removed** (replaced by assignment_invitation) |
| —                          | `assignment_invitation_table` | **New**                                         |
| —                          | `system_mail_table`           | **New**                                         |

**DB name change**: Refactor hardcodes `db_poprako_s` in `repo_infra.NewPgGdb()`.
Root uses `poprako_s_db` (dev) / `poprako` (prod). This must be harmonized —
see Phase 3.

**Migration strategy for existing databases**: Since this is a breaking schema change,
existing databases need either:

- (a) Fresh setup (drop & recreate) — recommended for dev
- (b) Manual migration script that: creates `assignment_invitation_table` and
  `system_mail_table`, and drops `notification_table` + `chapter_invitation_table`
  (if the new code doesn't reference them)

批注：没有需要保护的已有数据，直接完整替换 migrations 文件夹即可。

**For dev**: `dev-database-setup.sh` should be updated to check for the correct final
table (`assignment_invitation_table` instead of `chapter_invitation_table`).

**For prod**: `prod-database-migrate.sh` baseline mechanism needs the baseline table
changed from `chapter_invitation` to `assignment_invitation_table`.

---

## Phase 3: Configuration & Environment Variable Migration

This is the most critical phase. The refactor uses **different env var names** and reads
them from viper (which is configured with `viper.AutomaticEnv()` in `cfg.NewAppCfg()`).

### 3.1 Env Var Name Mapping

| Root Env Var           | Refactor Env Var                                                       | Used In                      |
| ---------------------- | ---------------------------------------------------------------------- | ---------------------------- |
| `APP_ENVIRONMENT`      | _(in app_config.json)_                                                 | `cfg/cfg.go`                 |
| `JWT_SECRET_KEY`       | `JWT_SECRET`                                                           | `infra/ext/token/parser.go`  |
| _(in app_config auth)_ | `JWT_EXPIRATION_HOURS`                                                 | `infra/ext/token/parser.go`  |
| `DATABASE_URL`         | `DATABASE_USER`, `DATABASE_PASSWORD`, `DATABASE_HOST`, `DATABASE_PORT` | `infra/repo/repo.go`         |
| _(N/A — new)_          | `DATABASE_NAME` (see 3.2)                                              | `infra/repo/repo.go`         |
| `OSS_PLATFORM`         | `OSS_PLATFORM` _(unchanged)_                                           | `infra/ext/oss/oss.go`       |
| `R2_*`                 | `R2_*` _(unchanged)_                                                   | `infra/ext/oss/r2_client.go` |
| `ALIYUN_*`             | `ALIYUN_*` _(unchanged)_                                               | `infra/ext/oss/`             |

### 3.2 Database Name Configuration

批注：所有的 .env 以 refactor 的为准，.env.sample 同步更改。

**Problem**: Refactor hardcodes database name as `db_poprako_s` in `repo_infra.NewPgGdb()`.
Root uses `poprako_s_db` (dev) / `poprako` (prod) from `DATABASE_URL`.

**Fix**: Make the database name configurable via `DATABASE_NAME` env var, with default
`db_poprako_s`. Modify `refactor/internal/infra/repo/repo.go`:

```go
dbName := viper.GetString("DATABASE_NAME")
if dbName == "" {
    dbName = "db_poprako_s"
}
```

批注：数据库名就直接硬编码，不允许修改。

### 3.3 Files That Must Be Updated

#### `docker-compose.yml`

**dev-main-server** environment block — replace:

```yaml
environment:
  DATABASE_URL: postgresql://devuser:devpassword@dev-postgres:5432/poprako_s_db?sslmode=disable
  JWT_SECRET_KEY: hardcoded_test_key
```

with:

```yaml
environment:
  DATABASE_USER: devuser
  DATABASE_PASSWORD: devpassword
  DATABASE_HOST: dev-postgres
  DATABASE_PORT: "5432"
  DATABASE_NAME: poprako_s_db
  JWT_SECRET: hardcoded_test_key
  JWT_EXPIRATION_HOURS: "336"
```

**prod-main-server** environment block — replace:

```yaml
environment:
  APP_ENVIRONMENT: production
  DATABASE_URL: postgres://poprako:${POSTGRES_PASSWORD}@prod-postgres:5432/poprako?sslmode=disable
  JWT_SECRET_KEY: ${JWT_SECRET_KEY}
```

with:

```yaml
environment:
  DATABASE_USER: poprako
  DATABASE_PASSWORD: ${POSTGRES_PASSWORD}
  DATABASE_HOST: prod-postgres
  DATABASE_PORT: "5432"
  DATABASE_NAME: poprako
  JWT_SECRET: ${JWT_SECRET}
  JWT_EXPIRATION_HOURS: ${JWT_EXPIRATION_HOURS:-336}
```

#### `docker/poprako-s-main/app_config.json`

Replace old format:

```json
{ "http_address": "0.0.0.0:8080", "auth": { "expiration_hours": 144 } }
```

with production-appropriate new format:

```json
{
  "environment": "production",
  "database": {
    "max_open": 24,
    "max_idle": 4,
    "max_life": 84600,
    "max_idle_time": 60
  },
  "http": {
    "host": "0.0.0.0",
    "port": 8080
  }
}
```

#### `docker/compose.prod.yml`

Same env var changes as `prod-main-server` in `docker-compose.yml`:

- `APP_ENVIRONMENT: production` → remove (now in app_config.json)
- `DATABASE_URL: ...` → `DATABASE_USER`, `DATABASE_PASSWORD`, `DATABASE_HOST`, `DATABASE_PORT`, `DATABASE_NAME`
- `JWT_SECRET_KEY` → `JWT_SECRET`
- Add `JWT_EXPIRATION_HOURS`

#### `.env.sample`

Replace with refactor-convention env vars:

```bash
# Database
DATABASE_USER=devuser
DATABASE_PASSWORD=devpassword
DATABASE_HOST=localhost
DATABASE_PORT=5432
DATABASE_NAME=poprako_s_db

# JWT
JWT_SECRET=change_me_for_prod_and_dev_app
JWT_EXPIRATION_HOURS=336

# OSS backend: r2 | aliyun | noop
OSS_PLATFORM=r2

# Cloudflare R2
R2_ACCOUNT_ID=
R2_ACCESS_KEY_ID=
R2_SECRET_ACCESS_KEY=
R2_BUCKET_NAME=
R2_CUSTOM_DOMAIN=

# Aliyun OSS (optional)
# ALIYUN_OSS_ACCESS_KEY_ID=
# ALIYUN_OSS_ACCESS_KEY_SECRET=
# ALIYUN_OSS_REGION=
# ALIYUN_OSS_BUCKET_NAME=
# ALIYUN_OSS_CUSTOM_DOMAIN=
```

#### `.github/workflows/deploy-prod.yml`

Replace env var references in both the `deploy` step and the heredoc:

- `JWT_SECRET_KEY` → `JWT_SECRET`
- `DATABASE_URL` construction → separate vars
- Add `JWT_EXPIRATION_HOURS` secret reference
- Remove `APP_ENVIRONMENT` from env (now in app_config.json)
- Update heredoc that writes `.env` file on server

#### `prod.sh`

No env var changes needed (it just orchestrates docker compose),
but verify it reads the correct env file.

---

## Phase 4: Docker & Deployment Updates

### 4.1 Root `Dockerfile`

No structural changes needed (it builds `./main.go` and copies `docs/`). However,
verify that the `docker/poprako-s-main/app_config.json` being copied is the updated one.

### 4.2 `docker/dev-database-setup.sh`

Update `final_table` check from `public.chapter_invitation_table` to
`public.assignment_invitation_table` (the last migration in the refactor set).

### 4.3 `docker/prod-database-migrate.sh`

Update `BASELINE_TABLE` default from `public.chapter_invitation` to
`public.assignment_invitation_table`.

### 4.4 `scripts/package-release.sh`

No changes needed (it just packages files). However, make sure the new
`docker/poprako-s-main/app_config.json` is picked up by the Docker build.

### 4.5 `scripts/remote-switch-release.sh`

No env var changes needed (uses the env file on the server). But the `.env` file
written by CI must use the new variable names.

---

## Phase 5: justfile & Development Tooling

### 5.1 `justfile`

Refactor's `justfile` is a subset of root's. It's missing:

- `build-main`, `build-database`, `save-main`, `save-database`
- `package-release`, `switch-release`

**Action**: Use root's `justfile` as base (keep the Docker build/save/release tasks),
update the `dev` command to remove the `DOCKER_BUILDKIT=0` flags from refactor.

The refactor adds `max-line` command — keep it.

### 5.2 Swagger

The refactor `swagger.json` / `swagger.yaml` includes new endpoints (sys_mail,
assignment_invitation, chapter_port). After replacement, regenerate with `swag init`
to ensure consistency, then copy to `docs/`.

---

## Phase 6: Post-Replacement Verification

### 6.1 Go Build

```bash
go mod tidy
go build ./...
```

### 6.2 Lint

```bash
golangci-lint run ./...
```

### 6.3 Dev Environment

```bash
just setup    # start postgres + run migrations
just dev      # build & run server
```

### 6.4 Production Build

```bash
docker build -f docker/poprako-s-main/Dockerfile -t poprako-s-main:test .
```

---

## Risk Assessment

| Risk                                                          | Severity   | Mitigation                                                      |
| ------------------------------------------------------------- | ---------- | --------------------------------------------------------------- |
| Env var name mismatch breaks existing deployments             | **High**   | Phase 3 checklist must be exhaustive; test with `just dev`      |
| Migration incompatibility with existing DBs                   | **High**   | Document that existing DBs need fresh setup or manual migration |
| Hardcoded DB name `db_poprako_s` doesn't match Docker configs | **Medium** | Make DB name configurable via env var (Phase 3.2)               |
| Docker env vars using old names cause runtime failures        | **High**   | Update all 3 compose files + CI + scripts                       |
| Refactor introduces new aggregates not in old DB schema       | **Medium** | Ensure migrations run before app starts                         |
| go.mod/go.sum divergence after `go mod tidy`                  | **Low**    | Just run `go mod tidy` after replacement                        |

---

## Execution Order

1. Make DB name configurable in refactor code
2. Replace `internal/`, `main.go`, `app_config.json`, `go.mod`, `go.sum`, `.golangci.yml`
3. Merge `docs/` (use refactor swagger, keep root deployment docs)
4. Replace `migrations/` with refactor migrations
5. Update all env var references (Phase 3.3 checklist)
6. Update Docker configs (Phase 4 checklist)
7. Update `justfile` (Phase 5)
8. Update `dev-database-setup.sh` and `prod-database-migrate.sh`
9. Run `go mod tidy` + verify build
10. Test dev environment end-to-end
