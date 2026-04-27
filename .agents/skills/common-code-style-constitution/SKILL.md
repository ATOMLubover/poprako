---
name: common-code-style-constitution
description: >
  The authoritative Go code style constitution for this repository.
  Load this skill before writing, editing, or reviewing ANY Go code in this
  workspace — even trivial changes. Covers comment coverage, identifier quoting,
  package naming (*_iface/*_impl/*_infra), constant documentation, naming
  case conventions, error format, imports, and function/method contracts.
  Use it proactively whenever generating or touching Go source files.
---

# Common Code Style Constitution

This is the **non-negotiable** style guide for all Go code in this repository.
Every rule below applies uniformly across all layers — domain, infra, app, and
test code alike. When in doubt, check this file first.

---

## 1. Comment Language

All comments are written in **English**. This includes doc comments, inline
body comments, block comments, and `// NOTE:` / `// FIXME:` annotations.

---

## 2. Identifier Quoting in Comments

Every Go identifier referenced inside a comment — type names, interface names,
struct names, method names, function names, field names, variable names, package
names — **must** be wrapped in backticks.

**Correct:**
```go
// `PushEv` appends one event to the pending list of `EvBase`.
// `EventBus` manages all events and their corresponding handlers.
// The `cx` context is passed through to all sub-operations.
```

**Wrong:**
```go
// PushEv appends one event to the pending list of EvBase.
// EventBus manages all events and their corresponding handlers.
```

This rule applies everywhere: doc comments on types, methods, functions,
fields, constants, inline body comments, and multi-line `NOTE`/`FIXME` blocks.

---

## 3. Comment Coverage — Every Declaration

Every function and method, whether exported or unexported, must have its own
doc comment. No exceptions, including one-liners and trivial accessors.

Every exported type declaration (interface, struct, named type, type alias)
must have a doc comment.

**Correct:**
```go
// `clearEv` clears all events in the pending list of `EvBase`.
// In case of memory leak, never ref a **single** element of `b.ev` directly.
func (b *EvBase) clearEv() {
    clear(b.ev)
}

// `EvTyp` wraps the type identifier of a domain event.
type EvTyp string
```

**Wrong:**
```go
func (b *EvBase) clearEv() {   // no comment — not allowed
    clear(b.ev)
}
```

---

## 4. Step-by-Step Body Comments

Any function or method that contains two or more distinct logical steps must
have an inline `//` comment immediately before each step. The comment describes
**what** the step does and **why** if non-obvious.

**Correct:**
```go
func (c *UserCreds) VerifyPwd(pwd string) error {
    // Compare the stored password hash against the raw input.
    if err := bcrypt.CompareHashAndPassword(
        []byte(c.PwdHash),
        []byte(pwd),
    ); err != nil {
        return fmt.Errorf("[UserCreds.VerifyPwd] unmatched password: %w", err)
    }

    // Emit a login event so downstream handlers can react.
    c.PushEv(&event.UserLoginEv{UserID: c.ID})

    return nil
}
```

Single-step functions only need the doc comment; no extra inline comment is
required inside the body.

---

## 5. Package Naming — `*_iface` / `*_impl` / `*_infra`

When a package has a separate interface variant and a concrete implementation
variant, the **package declaration** (not just the import alias) must follow
this naming scheme:

| Variant | Package declaration | Typical path |
|---|---|---|
| Interface / abstraction | `package xxx_iface` | `internal/xxx/` |
| Domain model implementation | `package xxx_impl` | `internal/domain/model/xxx/` |
| Infrastructure implementation | `package xxx_infra` | `internal/infra/xxx/` |

**Examples from this codebase:**
```go
// internal/event/event.go
package event_iface

// internal/domain/model/event/user.go
package event_impl

// internal/infra/event/user.go
package event_infra
```

Pure domain packages with no split (e.g. aggregate models) use a plain
descriptive name: `package aggr`, `package model`, `package service`.

### Import Aliases

When importing an `*_iface` or `*_impl` package, the import alias **must**
match the package declaration exactly:

```go
import (
    event_iface "poprako-s/internal/event"
    event_impl  "poprako-s/internal/domain/model/event"
)
```

Never import these packages without an alias or with an arbitrary alias.

---

## 6. Import Organization

Imports are arranged in **three groups**, separated by blank lines:

1. Standard library
2. Internal module packages (`poprako-s/internal/...`)
3. Third-party packages

```go
import (
    "context"
    "errors"
    "fmt"

    event_iface "poprako-s/internal/event"
    event_impl  "poprako-s/internal/domain/model/event"

    "golang.org/x/crypto/bcrypt"
    "gorm.io/gorm"
)
```

Omit a group entirely if it is empty. Never mix groups on the same line or
collapse them into one.

---

## 7. Constant Documentation

Every `const` block and every individual constant must be documented.

- A **block** with a shared meaning gets a single leading comment describing
  the group.
- An **individual constant** within a block, if its meaning differs from the
  group theme or needs clarification, gets its own inline `//` comment on the
  line above it.
- A **standalone constant** always gets a one-line doc comment above it.

**Correct:**
```go
// Event type identifiers for user-domain events.
const (
    // `EvUserLogin` is the event type for a successful user login.
    EvUserLogin event_iface.EvTyp = "event:user_login"

    // `EvUserCreated` is the event type emitted when a new user is registered.
    EvUserCreated event_iface.EvTyp = "event:user_created"
)
```

**Wrong:**
```go
const (
    EvUserLogin event_iface.EvTyp = "event:user_login"  // no comment — not allowed
)
```

---

## 8. Naming Case — camelCase / PascalCase (Not Go-Acronym Style)

All identifiers follow standard **camelCase** and **PascalCase** rules.
Acronyms are **not** fully uppercased — only the first letter of an acronym is
capitalised (or lowercased) according to its position.

| Wrong (Go-style) | Correct |
|---|---|
| `HTTPHandler` | `HttpHandler` |
| `URLParser` | `UrlParser` |
| `parseURL` | `parseUrl` |
| `userID` | `userId` |
| `getHTTPClient` | `getHttpClient` |
| `OSS` (as standalone exported name) | `Oss` |
| `JSONData` | `JsonData` |

This applies to struct fields, function names, variable names, type names,
constant names, and interface method names uniformly.

---

## 9. Context Parameter Name

The context parameter is **always** named `cx`, never `ctx` or any other name.

```go
func (b *EventBus) Pub(cx context.Context, ev []Event) { ... }
func NewRepoFromCx(cx context.Context) (Repo, error)   { ... }
```

---

## 10. Receiver Names

Method receivers use a **single lowercase letter** that semantically represents
the type. Consistent conventions used in this codebase:

| Type pattern | Receiver |
|---|---|
| `UserBase`, `UserCreds` | `u`, `c` |
| `EvBase` | `b` |
| `*xxxRepoImpl` | `r` |
| `*xxxAppImpl` | `a` |
| `*xxxServiceImpl` | `s` |
| `*xxxHandler` | `h` |
| event structs (`*UserLoginEv`) | `e` |

Never use `self`, `this`, or the full type name as a receiver.

---

## 11. Constructor Conventions

Constructors are named `NewXxx` and **always return the interface type**, never
the concrete struct pointer:

```go
// Correct — returns the interface:
func NewEventBus() event_iface.EventBus { return &eventBusImpl{} }

// Wrong — returns the concrete type:
func NewEventBus() *eventBusImpl { return &eventBusImpl{} }
```

If construction can fail, return `(InterfaceType, error)`.

No named return values anywhere.

---

## 12. Error Conventions

Errors originating inside a method are prefixed with `[Pkg.Method]` in
brackets to make log triage easy:

```go
// Wrapping an external error:
return fmt.Errorf("[UserCreds.VerifyPwd] unmatched password: %w", err)

// Fresh error with no underlying cause:
return errors.New("[EventBus.Pub] event bus is closed")
```

- Use `fmt.Errorf("...: %w", err)` when wrapping an underlying error.
- Use `errors.New("...")` when creating a fresh error with no underlying cause.
- Never return a raw repo/infra error without a bracket-prefixed wrapper at
  the boundary where the layer transition happens.

---

## 13. Automated Checks

A script is bundled to verify rule §3 (type comment coverage) mechanically.

**Script:** `scripts/check-type-comments.js`  
**Run from the project root:**
```bash
pnpm --dir .agents/skills/common-code-style-constitution check-types <path1> [path2 ...]
```

The script accepts relative paths (resolved from your current working
directory), recurses into directories, and exits with code `1` when any `type`
declaration is found without a comment line directly above it.

---

## 14. Related Skills

These companion skills cover sub-topics in more depth. Load them when the
relevant situation arises:

- **`defer-format`** — Every `defer` statement must have at least one blank
  line both above and below it.
- **`comment-style`** — Detailed grammar and punctuation rules for Go comments.
- **`domain-repo-style`** — Naming and comment conventions specific to
  `domain/repo` interface files.
- **`domain-service-style`** — Code standards for `domain/service` implementations.
- **`app-layer-style`** — Conventions for the `app` layer (use-case orchestration).
