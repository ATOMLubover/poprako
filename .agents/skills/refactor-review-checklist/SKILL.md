---
name: refactor-review-checklist
description: >
  Use when reviewing uncommitted or newly added refactor slices under `refactor/`
  especially when the review must verify migration schema correctness index coverage
  entity-field alignment typed include completeness transaction and domain-event flow
  and API HTTP swagger or `val` comment accuracy against real behavior
---

# Refactor Review Checklist

## Summary

Use this skill to review one refactor slice end to end before merge

Treat `refactor/migrations` as schema truth

Treat old project code under `internal/` as behavior reference only

Do not silently accept schema or behavior drift without calling it out

## Review Order

1. Lock the review scope from git diff and keep the slice narrow
2. Read related migrations first and extract exact table columns nullability defaults foreign keys and indexes
3. Read the comparable implemented slice in `refactor/` to confirm local structure conventions
4. Read the old implementation in `internal/` only for behavior parity and async side effects
5. Review domain app infra http from inner semantics to outer transport
6. Validate findings with the cheapest focused check available such as `go test` or targeted compile errors when needed

## Schema Review

Check these points against migrations first

- Every entity field and GORM column tag must map to a real migrated column name and nullable rule
- Every repo query path introduced by the slice must be covered by a suitable index or partial index
- For composite indexes check left-prefix usefulness against actual `WHERE` plus `ORDER BY`
- For list APIs verify filtering columns and default ordering columns are indexed together when the query shape needs it
- For uniqueness rules verify the migration enforces the invariant rather than only app logic
- If a query relies on a partial predicate such as `pending = TRUE` or `read = FALSE` confirm the index predicate matches the query predicate exactly
- If an update path depends on `PUT` semantics confirm nullable fields can really be written as SQL `NULL`

## Entity And Repo Review

Check infra implementation against schema and domain contract

- `entity` row structs must cover all persisted columns needed by the slice without inventing extra fields
- `ToXxxAggr` and `NewXxxRowFromAggr` mappings must preserve nullability and not coerce `NULL` into zero values
- Repo method signatures should use typed query objects and typed includes instead of raw strings or many booleans
- If an aggregate exposes relation fields then the repo contract and infra implementation must both support typed includes completely
- If no relation field exists the review should explicitly note that include completeness is not applicable rather than skipping the check silently
- Repo list methods should use stable base queries deterministic ordering and precise selected columns
- `Delete` must mean hard delete and `Remove` must mean soft delete
- Transaction-scoped app flows should obtain repos through `prov.XxxRepo()` inside `repo_iface.RunWithTxn[T]`

## Logic Review

Check business flow completeness rather than only happy-path CRUD

- Permission checks must align with old behavior and current domain rules
- Transaction closures should decide business rejection near the failing step and return stable app-facing messages
- Expected side effects that must happen only after commit should stay outside the transaction closure
- Expected side effects that must remain atomic with writes should stay inside the transaction and be justified
- If the old system emitted domain events or async follow-up work for this slice verify the refactor slice still triggers them or explicitly record the gap
- Check join or upsert flows for race windows that rely on DB uniqueness without handling duplicate-key fallback
- Check list and mark-read style operations for ownership filters so callers cannot read or mutate other users' data

## HTTP And Documentation Review

Check transport shape against real app signatures and behavior

- Swagger `@Param` declarations must match the real handler inputs exactly including path query and body fields
- Swagger success and failure response docs must describe the real `res.HttpRes` wrapper rather than naked app payloads
- Auth mode in godoc must match runtime behavior including header and cookie fallback
- `val` struct comments must explain non-obvious fields and encodings such as `role_mask` bit semantics
- For list endpoints confirm pagination params are documented and forwarded explicitly
- For body structs confirm field comments are sufficient for clients who cannot infer semantics from field names alone

## Findings Format

Return findings first ordered by severity

For each finding include

1. The broken contract or risk
2. The concrete file anchor
3. Why it is incorrect relative to migrations old behavior or refactor conventions
4. The likely impact or missing validation

If no issue is found in one checklist area say that explicitly and mention any remaining uncertainty or test gap

## Example Prompts

- Review the uncommitted refactor slice for schema index entity include txn event and swagger completeness
- Compare this refactor feature against old `internal/` behavior and list all regressions without fixing them
- Audit whether the new `val` structs and swagger docs fully explain `role_mask` and wrapped HTTP responses