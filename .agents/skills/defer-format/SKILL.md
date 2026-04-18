---
name: defer-format
description: |
  Enforce and document the repository convention: every `defer` statement
  must have at least one blank line both above and below the `defer` line.
  Use this skill to generate guidance, a quick-check script, and example edits.
---

# Defer Statement Formatting Skill

## Purpose

Enforce a simple formatting rule for `defer` statements in Go files:
there must be at least one empty line immediately before and after the `defer` line.

## Why

This improves visual separation of cleanup logic and reduces accidental
line-joining by automated edits.

## Rule (precise)

- For any line matching `^\s*defer\b.*`, the previous non-eof line must be blank,
  and the following non-eof line must be blank.

## How to apply

1. Run the included quick-fix PowerShell script (or similar) to insert blank
   lines where missing.
2. Review git changes and run `gofmt`/linters if desired.

## Quick-check script (PowerShell)

```powershell
$root = 'D:\\my_projects\\poprako\\poprako-s\\internal'
Get-ChildItem -Path $root -Recurse -Include *.go | ForEach-Object {
  $lines = Get-Content -Path $_.FullName -Encoding UTF8
  for ($i=0; $i -lt $lines.Count; $i++) {
    if ($lines[$i] -match '^[ \t]*defer\b') {
      $prev = if ($i-1 -ge 0) { $lines[$i-1].Trim() } else { '' }
      $next = if ($i+1 -lt $lines.Count) { $lines[$i+1].Trim() } else { '' }
      if ($prev -ne '' -or $next -ne '') { Write-Output "Violation: $($_.FullName):$($i+1)" }
    }
  }
}
```

## Example prompt

- "Run the `defer-format` skill and fix all `defer` spacing issues under `internal/`."

---

Edit this SKILL.md to include additional automation or stricter checks.
