# poprako-s — Agent Context

`poprako-s` is an **event-driven backend service** written in Go for managing manga (comic) translation projects. It handles teams, worksets, comics, chapters, pages, translation units, assignments, and user/member management.

---

## Code Style Constitution

Before writing or reviewing **any** Go code in this repository, agents **must**
load and follow the code style constitution:

**Skill:** `common-code-style-constitution`  
**Location:** `.agents/skills/common-code-style-constitution/SKILL.md`

The rules defined there are non-negotiable and apply uniformly across every
layer of the codebase. They cover comment language and coverage, identifier
quoting, package naming (`*_iface` / `*_impl` / `*_infra`), constant
documentation, camelCase/PascalCase conventions (no Go-style acronym
uppercasing), context parameter naming, error format, import organization,
and constructor conventions.

Other layer-specific skills (listed at the bottom of the constitution) extend
these base rules for their respective packages.

---

## Module & Runtime

- **Go module**: `poprako-s` (`go 1.25.5`)
- **Entry point**: `main.go` — boots an event bus and registers all event handlers; no HTTP server in this binary (it is the **async side-car** service, not the API server)
- **Build / task runner**: `justfile`

---

## Architecture

```
main.go
└── internal/
    ├── app/            # Application layer: use-case orchestration, VO validation
    │   ├── event_handler/  # Domain event handlers
    │   └── val/            # Input value objects (args/VO types per use-case)
    ├── cfg/            # Config structs (auth, etc.)
    ├── domain/
    │   ├── event/      # Event interface + EventSource interface
    │   ├── ext/        # External service abstractions (oss/)
    │   ├── model/      # Domain models (plain Go structs, also contain business logic)
    │   ├── repo/       # Repository interfaces (one file per aggregate)
    │   └── service/    # Domain services (stateless, business rule validation)
    └── infra/
        ├── event/      # EventBus concrete impl (+ mock/)
        ├── ext/        # External service concrete impls (oss/)
        └── repo/       # GORM repository concrete impls
            ├── entity/ # SQL row structs + ToXxx() domain model converters
            └── mock/   # Mock repos for unit tests
```

---

## Domain Aggregates & Key Models

| Model            | Notes                                                                                              |
| ---------------- | -------------------------------------------------------------------------------------------------- |
| `UserInfo`       | System user; has avatar (OSS two-step), stats (`UserStats`), JWT via `GenToken`                    |
| `TeamInfo`       | Translation group/team; has avatar                                                                 |
| `MemberInfo`     | Team membership with 7 role timestamps (`assigned_raw_provider_at` … `assigned_admin_at`)          |
| `InvitationInfo` | Invite to join a team; `pending` bool + 7 `to_be_*` role booleans                                  |
| `WorksetInfo`    | A batch of comics owned by a team; tracks `comic_count`                                            |
| `ComicInfo`      | A single manga title inside a workset; has `chapter_count`, `pinned_*` replica columns (see below) |
| `ChapterInfo`    | A chapter; has a 9-step workflow timeline + `IsPinned` flag                                        |
| `PageInfo`       | A page inside a chapter; `is_uploaded`, `oss_key`                                                  |
| `UnitInfo`       | A translation unit (text box) on a page; `x_coord`/`y_coord` stored as `REAL`→`int`                |
| `AssignmentInfo` | Role assignment linking a user to a chapter with 7 `*time.Time` role timestamps                    |

### Workflow State Machine (ChapterInfo)

Chapters progress through a non-linear workflow tracked by nullable timestamps:

Every workflow can progress independently(like, finish typeset while proofread has not finished)

Each phase has `<phase>_at` (completed) and/or `<phase>ing_at` (started) columns. `WorkflowPhase` enum: `Pending / Ongoing / Completed`.

---

## Critical Business Invariants

### 1. Unique Pinned Chapter

- Every `chapter_table` row has a `pinned` (bool) column.
- **Every `ChapterRepo.Create` must atomically**:
  1. `UPDATE chapter_table SET pinned=false WHERE comic_id=? AND deleted_at IS NULL AND pinned=TRUE`
  2. `INSERT` new chapter with `pinned=true`
  - Done inside a `gorm.DB.Transaction`.
- `ChapterRepo.FindPinnedByComicID` returns the single pinned chapter for a comic.

### 2. Comic Replica Fields (`pinned_*`)

- `comic_table` mirrors the pinned chapter's workflow timestamps as `pinned_uploaded_at`, `pinned_transalating_at`, `pinned_translated_at`, etc. (10 columns total, plus `has_pinned_chapter` bool).
- These are **read-only for query filtering** inside `ComicRepo.List / Count`.
- They are **written only by event handlers** (`ChapterPublishedHandler`, etc.) — **never inside repo Create/Update directly**.

---

## Infrastructure Conventions

### Repo layer (`internal/infra/repo/`)

- Package name: `repo_infra`
- Constructor pattern: `NewXxxRepo(gdb *gorm.DB) iface.XxxRepo`
- All repos support `FromTxnCx(cx context.Context)` extracting a `*gorm.DB` from context (`txnKey`) for cross-repo transactions.
- **Insert pattern**: use `map[string]any` (never a struct pointer) to avoid GORM skipping zero values.
- **Atomic counters**: `gorm.Expr("col + ?", delta)`.
- **Upsert (stats)**: `clause.OnConflict{DoNothing: true}`.
- **Table name constants** live in `entity/` package, e.g. `entity.ComicTable = "comic_table"`.

### Soft-delete tables

`user_table`, `team_table`, `member_table`, `comic_table`, `chapter_table` — always add `AND deleted_at IS NULL` to queries.

### Hard-delete tables

`invitation_table`, `workset_table`, `assignment_table`, `page_table`, `unit_table`.

### Entity layer (`internal/infra/repo/entity/`)

- One file per aggregate: `user.go`, `team.go`, …
- Each file exports: a `const XxxTable`, a `XxxInfoRow` struct with GORM column tags, and a `ToXxxInfo(row XxxInfoRow) model.XxxInfo` converter.

---

## Testing

- Unit tests live in `internal/app/*_test.go` (one per app-layer file).
- Tests use mock repos from `internal/infra/repo/mock/` and mock event bus from `internal/infra/event/mock/`.
- Integration / infra tests are **not yet implemented** (no test DB setup).
- Test helper constructors are in `constructors_test.go` and `test_helpers_test.go`.
- See `/memories/repo/testing.md` for repo-level test notes.

---

## Key Files for Orientation

| File                                | Purpose                                                 |
| ----------------------------------- | ------------------------------------------------------- |
| `internal/domain/model/workflow.go` | Workflow phase constants and transition strings         |
| `internal/domain/model/chapter.go`  | Chapter state machine + event emission                  |
| `internal/infra/repo/chapter.go`    | Pinned transaction logic                                |
| `internal/infra/repo/comic.go`      | `applyComicWorkflowFilter` helper + replica field usage |
| `internal/app/event_handler/`       | Handlers for chapter lifecycle events (side-effects)    |
| `docs/plans/impl-repo-impl.md`      | Detailed implementation plan reference                  |

---

## Migration Files

Located in `migrations/`. Follow the pattern `YYYYMMDDHHMMSS_<description>.{up,down}.sql`. All schema is already applied; no pending migrations as of April 2026.

---

## Refactor Hard Conventions

The following rules are mandatory for the `refactor/` implementation:

- For current-user identifiers in app/http signatures and local variables, use `currUid` naming. Do not use `currUserId`.
- If an aggregate has relation fields (for example `Workset.Team`), the corresponding repo contract and infra implementation must support typed `includes` via `enum.XxxIncl` and preload mapping.
- Nullable DB columns must stay nullable through domain and infra mappings. Do not coerce nullable fields into non-null defaults during assemble/convert.
- `Put`-style update payloads must overwrite target fields fully. For nullable fields, `nil` means write SQL `NULL` (not "skip update").
- In app transaction flows, follow the `user` app pattern: determine reject code via local state and normal error flow. Do not invent sentinel business errors such as `errNotAdmin`.
- Naming semantics are fixed: `Delete` means hard delete, `Remove` means soft delete.
- Common abbreviations are mandatory in Go identifiers. For example: use `Desc`/`desc` instead of `Description`/`description`. SQL column names and SQL literals are excluded.
- In app signatures, when business inputs beyond `cx` and `currUid` are more than one field, they must be wrapped into `val` args structs.
- All list interfaces must carry pagination explicitly.
