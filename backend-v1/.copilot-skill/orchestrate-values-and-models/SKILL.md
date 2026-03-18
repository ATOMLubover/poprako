# orchestrate-values-and-models

Purpose

- Provide a clear, repeatable workflow for designing `value` and `model` objects consistent with this repository's conventions.
- Reduce ambiguity when converting between layers and keep constructors/validation consistent.

Scope

- Workspace-scoped skill for Go services in this repo. Applies to `internal/value/*` and `domain/model/*`.

Principles

- Single responsibility: `value` objects are DTOs (API layer), `model` objects are domain models (business logic + repository).
- Construction via `New` functions only: always expose `NewXxx(...)` constructors; avoid direct struct literals outside constructors.
- Minimal pointers: prefer value fields unless mutability/optional semantics require `*T`.
- Args exception: `Args` types (request payloads) may be deserialized directly and do not require `New` constructors.
- Explicit validation: provide a `Validate() error` method for all `Args` and for any object requiring invariant checks.
- Clear JSON tagging: use snake_case JSON tags, apply `omitempty` according to the PUT/POST semantics used in the project.

Conversion patterns

- Provide deterministic conversion functions between layers, e.g. `NewWorksetInfoFromModel(model.WorksetInfo) value.WorksetInfo`.
- Keep conversions near the `value` package (factory functions that take `model` types) to decouple API and domain packages.
- Avoid cyclic imports by using minimal types in function signatures and pointer fields where necessary.

Constructor pattern

- Always return fully-initialized, valid objects:
  - `func NewWorksetInfoFromModel(m model.WorksetInfo) WorksetInfo { ... }`
  - Do not allow callers to construct core objects via struct literals.

Args & Update semantics

- `Args` types represent raw input; `Validate()` returns human-friendly errors.
- For update structs use PUT semantics: `nil` = set NULL, non-nil = update value, empty string = set to empty string. Document this in the type godoc.

Documentations

- Add godoc to handlers and to the `New` and `Validate` functions.

Checklist for creating a new pair (`value` + `model`)

- Create `domain/model/<name>.go` with domain invariants and repository interfaces.
- Create `internal/value/<name>.go` with DTOs and `New...FromModel` conversion functions.
- Add `New` constructors for model/value where appropriate.
- Implement `Validate()` on `Args` structs and test them.
- Update related handlers in `internal/api/http` with comments explaining mapping.

Suggested prompts to run this skill

- "Create a new `value` + `model` pair for `Project` following the skill conventions."
- "Review `domain/model/task.go` and produce missing `New` constructors and `value` conversions."

Ambiguities to confirm with owner

- Pointer vs value preference for large structs (confirm size threshold).
- Preferred error message language (English vs Chinese) for `Validate()` results.

Outcome

- Produces a short checklist, concrete code patterns, and a small template for `New`/`Validate`/conversion functions to keep `value` and `model` layering consistent.
