---
name: refactor-app-txn-pattern
description: >
  Use when writing or refactoring transaction flows in refactor/internal/app/impl.
  Apply this skill whenever a use case needs repo transaction orchestration,
  provider-based repo retrieval, precise app-level error mapping, and commit-after side effects.
---

# Refactor App Transaction Pattern

## Summary

This skill standardizes app-layer transaction orchestration in refactor mode

It mirrors the improved user registration transaction style and should be reused for all similar flows

## Required Pattern

1. Use `repo_iface.RunWithTxn[T]` instead of `txnCtrl.RunWithTxn(func(context.Context) error)`
2. In transaction closure, accept `prov repo_iface.Prov` and get repositories by `prov.XxxRepo()`
3. Return typed `res.AppRes[...]` from the transaction closure so business rejection is decided inside transaction
4. Return actual infra error for rollback and logging, but keep user-facing message stable
5. Use `res.DefErr()` only as a rollback trigger for expected business rejection where internal error details should not leak
6. Keep event publish or external side effects after transaction commit unless business requires in-transaction write
7. Never call `a.txnCtrl.RunWithTxn(...)` directly in app impl methods
8. Never use local `errCode` fallback branching for transaction error mapping

## Error Mapping Rules

1. Expected business rejection:
- Return `res.Reject(..., res.BadRequest, msg)` with a stable message
- Return `res.DefErr()` as transaction error

2. Infra or system failure:
- Return `res.Reject(..., res.ServerError, msg)`
- Return original `err` for diagnostics and rollback

3. Outer function handling:
- If transaction fails, log with the returned `err`
- Return the `AppRes` produced by transaction closure
- Do not reconstruct a second rejection branch outside transaction

4. Check style consistency:
- Use nested `if err != nil { ... }` blocks for all error checks in one flow
- Avoid split style like `if err != nil && condition` followed by another `if err != nil`
- Keep expected-business checks inside the same error block before server-error fallback

## Closure Shape

```go
if re, err := repo_iface.RunWithTxn[res.AppRes[val.SomeRes]](
    a.txnCtrl,
    func(prov repo_iface.Prov) (res.AppRes[val.SomeRes], error) {
        someRepo := prov.SomeRepo()

        if err := someRepo.DoSomething(); err != nil {
            if repo_infra.IsNotFound(err) {
                return res.Reject[val.SomeRes](res.BadRequest, "业务错误消息"), res.DefErr()
            }

            return res.Reject[val.SomeRes](res.ServerError, "系统错误消息"), err
        }

        return res.Accept(&val.SomeRes{}), nil
    },
); err != nil {
    lgr.Error("[someAppImpl.SomeMethod] failed to run transaction", zap.Error(err))

    return re
}
```

## Review Checklist

1. No `TxnXxxRepo(cx)` in new transaction code
2. No repeated repo retrieval error branches in app impl
3. Error decision is near the failing step and mapped by business semantics
4. Outer layer only logs and returns transaction result
5. Commit-after side effects are outside transaction closure
6. Function naming and signature follow refactor app-layer conventions

## Example Prompts

- Refactor this app transaction to provider-based generic transaction pattern
- Align this use case with user registration transaction style
- Improve error mapping in this transaction flow with AppRes return inside transaction
