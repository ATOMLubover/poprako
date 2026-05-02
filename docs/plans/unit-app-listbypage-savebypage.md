# Unit App: ListByPage & SaveByPage Implementation Plan

## Overview

Implement the two core unit app use-cases — `ListByPage` (read) and `SaveByPage` (write) — along with all required infrastructure, repository extensions, and transactional unit-count maintenance across the Page and Chapter aggregates.

---

## Architecture Summary

```
HTTP  GET  /page/{page_id}/units       → ListPageUnits   → UnitApp.ListByPage
HTTP  POST /page/{page_id}/units       → SavePageUnits   → UnitApp.SaveByPage
```

**SaveByPage** is the high-impact operation. It must, within a single transaction:

1. Snapshot the page's current unit counts (pre-mutation baseline)
2. Validate and execute the diff via `UnitSvc.ApplyOps` (creates, saves, deletes, reindex)
3. Recompute the page's unit counts from the DB (post-mutation ground truth)
4. Persist the **page** unit counts (absolute `SET`)
5. Compute the delta: `Δ = new_page_counts − old_page_counts`
6. Persist the **chapter** unit counts via atomic `SET x = x + Δ` (no lost updates)
7. Touch the comic's `last_active`

**ListByPage** is read-only: load permissions, return all units for a page ordered by `index ASC`.

---

## Core Optimization: Delta Propagation to Chapter

Rather than listing every page of the chapter to sum their counts (which requires either a chapter-level lock or risks lost updates under READ COMMITTED), we:

1. **Read** the page's current `TotalUnitCount` / `TranslatedUnitCount` / `ProofreadUnitCount` from the loaded page aggregate at transaction start. These are the committed pre-mutation values.
2. **Execute** `ApplyOps` — all unit mutations and reindexing.
3. **Recompute** the page's true unit counts via `CountByPage` (a SQL aggregation on `t_unit WHERE page_id = ?`).
4. **Write** the page counts as absolute values via `pageRepo.SetUnitCounts`.
5. **Compute delta**: `Δtotal = newTotal − oldTotal`, same for translated/proofread.
6. **Apply delta** to the chapter via `chapterRepo.AdjustUnitCounts(chapterId, Δtotal, Δtranslated, Δproofread)`.

The chapter adjustment SQL is:
```sql
UPDATE t_chapter
SET total_unit_count      = total_unit_count      + ?,
    translated_unit_count = translated_unit_count + ?,
    proofread_unit_count  = proofread_unit_count  + ?,
    updated_at            = ?
WHERE id = ? AND deleted_at IS NULL
```

### Why this avoids lost updates

With `SET x = x + Δ`, PostgreSQL atomically reads the current committed value and adds the delta. Two concurrent transactions modifying different pages of the same chapter each compute their own page's delta and apply it independently. The sum of deltas converges to the correct total.

Example:
- Initial: Page A has 5 units, Page B has 3 units. Chapter total = 8.
- T1 modifies Page A: old=5, new=7, Δ=+2
- T2 modifies Page B: old=3, new=2, Δ=−1
- T1 commits: chapter total = 8 + 2 = 10
- T2 commits: chapter total = 10 + (−1) = 9
- True state: Page A=7, Page B=2, sum=9 ✓

No page listing, no chapter-level lock, no lost updates.

---

## Permission Model

| Operation | Permission | Rationale |
|-----------|-----------|-----------|
| `ListByPage` | Team member (`CanListChapter`) | Same as listing pages |
| `SaveByPage` | Has any assignment on the chapter | Any assigned contributor (translator, proofreader, etc.) can edit units |

The `CanEditPageUnits` check mirrors `CanResvPages` / `CanMarkImageUploaded` in `PageSvc`: use `assignmentRepo.GetByChapterUserId` and verify the assignment exists (non-nil). Any role on the chapter is sufficient.

---

## Files to Create

### 1. `internal/infra/repo/entity/unit.go` — GORM entity rows

```
UnitRow            — read model mapping t_unit columns, ToUnitAggr() converter
UnitCreRow         — write model for INSERT (Id from cre.LocalId, which holds the server-generated id)
UnitUpdRow         — write model for UPSERT (ON CONFLICT DO UPDATE all fields)
UnitReindexUpdRow  — write model for single-index UPDATE
```

### 2. `internal/infra/repo/unit.go` — GORM UnitRepo implementation

Implements the expanded `UnitRepo` interface:

- `ListByPage` — `SELECT * FROM t_unit WHERE page_id = ? ORDER BY index ASC`
- `ListIndicesByPage` — `SELECT id, index FROM t_unit WHERE page_id = ?`
- `CountByPage` — aggregation returning (total, translated, proofread)
- `Create` — INSERT from UnitCreRow
- `Save` — UPSERT via GORM `clause.OnConflict{UpdateAll: true}` (reconstruct deleted units)
- `Reindex` — batch UPDATE via raw SQL `CASE WHEN`
- `Delete` — `DELETE FROM t_unit WHERE id = ?`

Constructor: `NewUnitRepo(gdb *gorm.DB) repo_iface.UnitRepo`

### 3. `internal/app/unit.go` — UnitApp interface

```go
type UnitApp interface {
    ListByPage(cx context.Context, currUid string, args *val.ListPageUnitsArgs) app_res.AppRes[val.ListPageUnitsRes]
    SaveByPage(cx context.Context, currUid string, args *val.SavePageUnitsArgs) app_res.AppRes[val.SavePageUnitsRes]
}
```

### 4. `internal/app/impl/unit.go` — UnitApp implementation

`unitAppImpl` struct with dependencies:
- `txnCtrl`, `unitSvc`, `chapterSvc`
- `memberRepo`, `worksetRepo`, `comicRepo`, `chapterRepo`, `pageRepo`, `unitRepo`, `assignmentRepo`
- `errClsf`

**`ListByPage` flow:**
1. Validate args (pageId non-empty)
2. Load page → chapter → comic (with Workset preload, i.e. `ComicInclWorkset`) → workset
3. Permission: `chapterSvc.CanListChapter(currUid, ws.TeamId, memberRepo, errClsf)`
4. `unitRepo.ListByPage(pageId)` → assemble `UnitVal` slice
5. Return `ListPageUnitsRes{Units, TotalUnitCount, TranslatedUnitCount, ProofreadUnitCount}` from the loaded page aggregate

**`SaveByPage` flow (within transaction):**
1. Validate args (pageId non-empty, diff non-nil, diff.PageId == args.PageId)
2. Within `repo_iface.RunWithTxn`:
   a. Load page via `pageRepo.GetById(args.PageId)` → verify exists
   b. **Snapshot** old page counts: `oldTotal := page.TotalUnitCount`, etc.
   c. Load chapter → comic (with Workset preload) → workset (ownership chain)
   d. Permission: `unitSvc.CanEditPageUnits(currUid, page.ChapterId, assignmentRepo, errClsf)`
   e. `unitSvc.ApplyOps(diff, unitRepo)` — validate + execute all ops + reindex
   f. `unitRepo.CountByPage(args.PageId)` → `(newTotal, newTranslated, newProofread)`
   g. `pageRepo.SetUnitCounts(args.PageId, newTotal, newTranslated, newProofread)`
   h. Compute deltas: `Δtotal := newTotal - oldTotal`, etc.
   i. `chapterRepo.AdjustUnitCounts(page.ChapterId, Δtotal, Δtranslated, Δproofread)`
   j. `comicRepo.TouchLastActive(comic.Id)`
   k. Return `SavePageUnitsRes{newTotal, newTranslated, newProofread}`
3. No domain events (counts are sync within the txn, no side effects)

### 5. `internal/app/impl/unit_log.go` — Logging decorator

Follows the standard `*LogApp` pattern: enrich logger with `curr_uid`, `page_id`, forward call.

### 6. `internal/app/impl/unit_util.go` — Assembly helpers

- `asmUnitVal(unit *aggr.Unit) val.UnitVal` — aggregate → value object
- `vfyListPageUnitsArgs(args)` — validation
- `vfySavePageUnitsArgs(args)` — validation + ensure `args.Diff.PageId == args.PageId`

### 7. `internal/app/val/unit.go` — Value objects

```go
type UnitVal struct {
    Id                 string  `json:"id"`
    PageId             string  `json:"page_id"`
    Index              int     `json:"index"`
    IsBubble           bool    `json:"is_bubble"`
    IsProofread        bool    `json:"is_proofread"`
    XCoord             float64 `json:"x_coord"`
    YCoord             float64 `json:"y_coord"`
    TranslatedText     *string `json:"translated_text,omitempty"`
    TranslatorComment  *string `json:"translator_comment,omitempty"`
    LastTranslatorId   *string `json:"last_translator_id,omitempty"`
    ProofreadText      *string `json:"proofread_text,omitempty"`
    ProofreaderComment *string `json:"proofreader_comment,omitempty"`
    LastProofreaderId  *string `json:"last_proofreader_id,omitempty"`
    CreatedAt          int64   `json:"created_at"`
    UpdatedAt          int64   `json:"updated_at"`
}

type ListPageUnitsArgs struct {
    PageId string `url:"page_id"`
}

type ListPageUnitsRes struct {
    Units               []UnitVal `json:"units"`
    TotalUnitCount      int       `json:"total_unit_count"`
    TranslatedUnitCount int       `json:"translated_unit_count"`
    ProofreadUnitCount  int       `json:"proofread_unit_count"`
}

type SavePageUnitsArgs struct {
    PageId string         `json:"page_id"`
    Diff   *aggr.UnitDiff `json:"diff"`
}

type SavePageUnitsRes struct {
    TotalUnitCount      int `json:"total_unit_count"`
    TranslatedUnitCount int `json:"translated_unit_count"`
    ProofreadUnitCount  int `json:"proofread_unit_count"`
}
```

### 8. `internal/api/http/unit.go` — HTTP handlers

- `ListPageUnits(st) iris.Handler` — `GET /page/{page_id}/units`
- `SavePageUnits(st) iris.Handler` — `POST /page/{page_id}/units`

Standard handler pattern: extract `currUid`, parse path/body params, call app, map `AppRes` → `HttpRes`.

---

## Files to Modify

### Domain Layer

#### `internal/domain/repo/unit.go` — Add 2 methods to UnitRepo

```go
ListByPage(pageId string) ([]*aggr.Unit, RepoErr)
CountByPage(pageId string) (total, translated, proofread int, RepoErr)
```

#### `internal/domain/repo/page.go` — Add 1 method to PageRepo

```go
// SetUnitCounts overwrites the three unit count columns on one page.
SetUnitCounts(id string, total, translated, proofread int) RepoErr
```

#### `internal/domain/repo/chapter.go` — Add 1 method to ChapterRepo

```go
// AdjustUnitCounts atomically increments the three unit count columns by the given deltas.
// Uses SQL `SET x = x + delta` to avoid lost updates under concurrent page edits.
AdjustUnitCounts(id string, deltaTotal, deltaTranslated, deltaProofread int) RepoErr
```

#### `internal/domain/repo/prov.go` — Add UnitRepo to Prov

```go
UnitRepo() UnitRepo
```

#### `internal/domain/svc/unit.go` — Add CanEditPageUnits

```go
func (UnitSvc) CanEditPageUnits(currUid string, chapterId string, assignmentRepo repo_iface.AssignmentRepo, clsf repo_iface.ErrClsf) svc_res.SvcRes
```

Verifies the user has any assignment (non-nil) on the chapter. Uses `assignmentRepo.GetByChapterUserId`. Pattern matches `CanResvPages` in `PageSvc`.

### Mock Layer

#### `internal/domain/repo/mock/unit.go` — Add ListByPage, CountByPage

- `ListByPage` — clone units for page, sort by Index
- `CountByPage` — iterate map, count total / `translated_text != nil` / `proofread_text != nil`

### Infra Layer

#### `internal/infra/repo/prov.go` — Wire UnitRepo

```go
func (p *provImpl) UnitRepo() repo_iface.UnitRepo {
    return NewUnitRepo(p.gdb)
}
```

#### `internal/infra/repo/page.go` — Implement SetUnitCounts

```go
func (r *pageRepoImpl) SetUnitCounts(id string, total, translated, proofread int) repo_iface.RepoErr {
    updRow := &entity.PageUnitCountsUpdRow{
        TotalUnitCount:      total,
        TranslatedUnitCount: translated,
        ProofreadUnitCount:  proofread,
        UpdatedAt:           time.Now(),
    }
    return r.gdb.Table(entity.PAGE_TABLE).
        Where("id = ?", id).
        Select("total_unit_count", "translated_unit_count", "proofread_unit_count", "updated_at").
        Updates(updRow).Error
}
```

#### `internal/infra/repo/chapter.go` — Implement AdjustUnitCounts

```go
func (r *chapterRepoImpl) AdjustUnitCounts(id string, deltaTotal, deltaTranslated, deltaProofread int) repo_iface.RepoErr {
    now := time.Now()
    return r.gdb.Table(entity.CHAPTER_TABLE).
        Where("id = ? AND deleted_at IS NULL", id).
        Updates(map[string]any{
            "total_unit_count":      gorm.Expr("total_unit_count + ?", deltaTotal),
            "translated_unit_count": gorm.Expr("translated_unit_count + ?", deltaTranslated),
            "proofread_unit_count":  gorm.Expr("proofread_unit_count + ?", deltaProofread),
            "updated_at":            now,
        }).Error
}
```

This uses `gorm.Expr` to produce `SET total_unit_count = total_unit_count + ?` — the update is atomic at the database level.

#### `internal/infra/repo/entity/page.go` — Add PageUnitCountsUpdRow

```go
type PageUnitCountsUpdRow struct {
    TotalUnitCount      int       `gorm:"column:total_unit_count"`
    TranslatedUnitCount int       `gorm:"column:translated_unit_count"`
    ProofreadUnitCount  int       `gorm:"column:proofread_unit_count"`
    UpdatedAt           time.Time `gorm:"column:updated_at"`
}
```

No entity row needed for chapter's `AdjustUnitCounts` — it uses a `map[string]any` with `gorm.Expr` directly (no struct can express `col = col + ?` via GORM tags).

### App Layer Wiring

#### `internal/api/http/http.go` — Add routes under the existing `page` party

```go
page.Get("/{page_id}/units", ListPageUnits(st))
page.Post("/{page_id}/units", SavePageUnits(st))
```

#### `internal/api/state/state.go` — Add UnitApp field

```go
UnitApp app_iface.UnitApp
```

Add to `NewAppState` parameters and struct literal.

#### `main.go` — Wire UnitApp

```go
unitSvc := svc.UnitSvc{}
unitRepo := repo_infra.NewUnitRepo(gdb)

unitApp := app_impl.NewUnitLogApp(
    app_impl.NewUnitApp(
        txnCtrl,
        unitSvc,
        chapterSvc,
        memberRepo,
        worksetRepo,
        comicRepo,
        chapterRepo,
        pageRepo,
        unitRepo,
        assignmentRepo,
        errClsf,
    ),
)
```

Add `unitApp` to `state.NewAppState(...)` call.

---

## Key Design Decisions

### 1. Chapter Counts: Delta via Atomic SQL

Chapter counts use `SET x = x + Δ` — never an absolute overwrite. This avoids lost updates without requiring a chapter-level lock or a full page-listing query. Each concurrent page edit computes its own page's delta; PostgreSQL serializes the atomic column additions.

### 2. Page Counts: Absolute Overwrite After Recompute

Page counts use `CountByPage` (SQL aggregation on `t_unit`) followed by an absolute `SetUnitCounts`. This is safe because:
- The page is only modified by us within this transaction (we hold the implicit row lock from modifying its units).
- Recomputing from `t_unit` is the authoritative ground truth — no delta tracking for the page itself.

### 3. UnitCre.LocalId reused as server ID

`UnitSvc.ApplyOps` overwrites `UnitCre.LocalId` with the server-generated ID before calling `repo.Create`. The repo reads `cre.LocalId` as the database `id`. The mock already follows this convention.

### 4. Save = UPSERT (reconstruct-on-delete)

`UnitRepo.Save` implements `INSERT ... ON CONFLICT (id) DO UPDATE`. If a unit was deleted by a concurrent submitter, a subsequent save from a stale client reconstructs it with the client's latest field values. The mock already implements this fallback.

### 5. No domain events for unit changes

Unit count updates and reindexing happen synchronously within the transaction. No notifications, stats, or external side effects are triggered by unit CRUD. Events exist at higher granularity (assignment creation, chapter publish).

### 6. Reindex uses batch UPDATE with CASE WHEN

Rather than N individual UPDATE statements, the GORM implementation issues a single `UPDATE t_unit SET index = CASE id WHEN 'id1' THEN 0 WHEN 'id2' THEN 1 ... END WHERE id IN (...)`.

---

## Edge Cases & Error Handling

| Scenario | Behavior |
|----------|----------|
| `SaveByPage` with empty diff (no ops) | `ApplyOps` validates CandOrder (empty). No DB writes. old=new, Δ=(0,0,0). Chapter unchanged. |
| `SaveByPage` where page has no units | `CountByPage` returns (0,0,0). old=(0,0,0). Δ=(0,0,0). Page/chapter unchanged. |
| `SaveByPage` where page initially has units, all deleted | old=(n,t,p), new=(0,0,0), Δ=(−n,−t,−p). Chapter correctly reduced. |
| `SaveByPage` on a page that doesn't exist | `pageRepo.GetById` returns NotFound → `BadRequest` |
| Concurrent `SaveByPage` on same page | Serialized by row-level locks on `t_unit` rows (the UPDATE/DELETE/INSERT ops in ApplyOps). `reorderUnits` handles concurrent-edit semantics. |
| Concurrent `SaveByPage` on different pages, same chapter | Each computes its own page delta. Chapter `SET x = x + Δ` atomically sums both. Correct. |
| `SaveByPage` with stale CandOrder (missing concurrent units) | `reorderUnits` places unknown units near their original neighbors by proportional slot mapping. |
| `ListByPage` on page with no units | Returns empty units slice with page counts (all 0). |
| `ListByPage` where user is not team member | `CanListChapter` rejects → `Forbidden` |
| `SaveByPage` where user has no chapter assignment | `CanEditPageUnits` rejects → `Forbidden` |
| `SaveByPage` where diff.PageId != path page_id | Validation rejects → `BadRequest` |

---

## SQL Count Query Detail

`CountByPage` SQL:
```sql
SELECT
    COUNT(*) AS total,
    COUNT(CASE WHEN translated_text IS NOT NULL THEN 1 END) AS translated,
    COUNT(CASE WHEN proofread_text IS NOT NULL THEN 1 END) AS proofread
FROM t_unit
WHERE page_id = ?
```

This maps to GORM via `Raw().Scan()`:
```go
type unitCountResult struct {
    Total      int
    Translated int
    Proofread  int
}
var res unitCountResult
r.gdb.Raw(`
    SELECT COUNT(*) as total,
           COUNT(CASE WHEN translated_text IS NOT NULL THEN 1 END) as translated,
           COUNT(CASE WHEN proofread_text IS NOT NULL THEN 1 END) as proofread
    FROM t_unit WHERE page_id = ?
`, pageId).Scan(&res)
```

---

## Implementation Order

1. **Repository extensions** (domain interfaces + mock)
   - `UnitRepo.ListByPage`, `CountByPage`
   - `PageRepo.SetUnitCounts`
   - `ChapterRepo.AdjustUnitCounts`
   - `Prov.UnitRepo()`
   - Update mock

2. **Infra GORM implementations**
   - `entity/unit.go` — all row types
   - `repo/unit.go` — UnitRepo SQL impl
   - `entity/page.go` — PageUnitCountsUpdRow
   - `repo/page.go` — SetUnitCounts
   - `repo/chapter.go` — AdjustUnitCounts (no entity row, uses `gorm.Expr`)
   - `prov.go` — wire UnitRepo

3. **Domain service extension**
   - `svc/unit.go` — CanEditPageUnits

4. **App layer**
   - `val/unit.go` — value objects
   - `impl/unit_util.go` — validators + assembly
   - `impl/unit.go` — UnitApp implementation
   - `impl/unit_log.go` — logging decorator
   - `app/unit.go` — UnitApp interface

5. **HTTP + wiring**
   - `api/http/unit.go` — handlers
   - `api/http/http.go` — routes
   - `api/state/state.go` — AppState
   - `main.go` — DI wiring

6. **Tests**
   - Unit tests for `CanEditPageUnits`
   - Unit tests for `CountByPage` mock
   - Integration test sketch for `SaveByPage` delta logic
