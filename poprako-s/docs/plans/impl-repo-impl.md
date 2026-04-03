# Repo Infra Full Implementation Plan (backend-v1 -> poprako-s)

## 1. Scope and Source of Truth

This implementation follows this priority strictly:

1. First source: SQL schemas in migrations/
2. Second source: current domain repo interfaces and domain models in poprako-s
3. Third source: backend-v1 implementation style as reference only

Implementation target folder:

- internal/infra/repo
- internal/infra/repo/entity

Non-target in this phase:

- app layer
- domain layer
- event handler logic (except infra data shape needed by handlers)

## 2. Confirmed Business Invariants

### 2.1 Unique pinned chapter per comic

- chapter_table has pinned column (boolean)
- At any time, each comic should have at most one non-deleted pinned chapter
- Create chapter should always make the new chapter become the unique pinned one
- This requires transaction semantics at app layer; repo methods must provide deterministic operations to support this

Repo-level support in this phase:

- ChapterRepo.Create inserts chapter with pinned value from creation payload
- ChapterRepo.Update supports explicit pinned field update
- ChapterRepo.FindPinnedByComicID queries by comic_id + pinned=true + deleted_at is null

Note: unique-pinned enforcement across all chapters of one comic is orchestrated in transaction flows above repo layer.

### 2.2 Comic pinned replica fields

- comic_table replica fields mirror pinned chapter workflow state
- Consistency must rely on transaction + sync event handler
- Repo needs to support fast conditional comic query on replica fields

Repo-level support in this phase:

- ComicRepo.List and ComicRepo.Count include workflow filters against comic_table pinned replica columns
- No implicit replica write in repo methods

## 3. Data Access Design

### 3.1 entity package structure

For each table, define:

- table name const
- row struct with gorm column tags
- conversion function row -> domain model

Files to implement:

- entity/user.go
- entity/team.go
- entity/member.go
- entity/invitation.go
- entity/workset.go
- entity/comic.go
- entity/chapter.go
- entity/page.go
- entity/unit.go
- entity/assignment.go

### 3.2 Repo implementation files

Implement full methods in:

- user.go
- team.go (new)
- member.go
- invitation.go
- workset.go
- comic.go
- chapter.go
- page.go
- unit.go
- assignment.go

Transaction context behavior:

- Keep existing NewXxxRepoFromCx and FromTxnCx pattern
- Reuse txnKey from tx_mgr.go

## 4. Method-Level Execution Plan

### 4.1 UserRepo

Implement:

- GetCredsByQQ
- GetByID
- GetByQQ
- List
- Create
- Update
- Remove (soft delete)
- RefreshLastLogin
- PreFillAvatarOSSKey
- ConfirmAvatarUploaded
- GetOrCreateStats (upsert style)
- PatchStats (delta by gorm.Expr)

### 4.2 TeamRepo

Create new infra file and implement:

- GetByID
- List
- Create
- Update
- Delete (hard delete)
- PreFillAvatarOSSKey
- ConfirmAvatarUploaded
- FromTxnCx

### 4.3 MemberRepo

Implement:

- GetByID
- Get
- List
- Exist
- Create
- Update
- Delete (soft delete)

### 4.4 InvitationRepo

Implement:

- GetByInviteeQQ
- List
- Create
- Update
- Invalidate
- Delete (hard delete)

### 4.5 WorksetRepo

Implement:

- GetByID
- List
- Count
- Create
- Update
- UpdateComicCount (delta)
- Delete (hard delete)

### 4.6 ComicRepo

Implement:

- GetByID
- List
- Count
- Create
- Update
- UpdateChapterCount (delta)
- Delete (hard delete)

List/Count filtering includes:

- workset_id required filter
- optional id filter
- optional fuzzy title filter on composed_title
- optional workflow filters mapped to pinned\_\* replica fields
- pagination from model.Pagination

### 4.7 ChapterRepo

Implement:

- GetByID
- FindPinnedByComicID
- List
- Count
- Create
- Update
- Remove (soft delete)
- UpdateStats

### 4.8 PageRepo

Implement:

- GetByID
- List
- GetStatsByID
- CreateBatch
- Update
- UpdateStats
- Delete
- DeleteBatch

### 4.9 UnitRepo

Implement:

- List
- CreateBatch
- PatchBatch (partial update by non-nil fields)
- DeleteBatch

### 4.10 AssignmentRepo

Implement:

- GetByID
- Get
- List
- Exist
- Create
- Update
- Delete

## 5. Query and Error Conventions

- Keep current repo_infra package naming and constructor style
- Keep simple error returns from gorm
- Use gorm.ErrRecordNotFound transparently
- Respect soft delete tables by adding deleted_at filter in queries
- Update updated_at explicitly for Update/Patch style operations

## 6. Validation and Done Criteria

Done criteria:

1. All repo methods compile with no not implemented placeholder
2. go build ./... passes
3. No added errors in modified files
4. Implementations match domain interfaces and migration schema

Execution sequence:

1. entity files
2. team repo file creation
3. existing repo files method completion
4. compile and fix
5. final self-check against this plan
