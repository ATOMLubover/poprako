---
name: "repository-tests"
description: "Generate or refactor repository tests that follow real application usage and shared test infrastructure"
argument-hint: "Target repositories or files, for example: team member invitation"
agent: agent
model: "GPT-5 (copilot)"
---

Write or refactor repository-layer tests for the requested targets.

Task focus:

- Work on repository tests only.
- Prefer large linear scenarios over many tiny isolated tests.
- Reuse existing shared test infrastructure when possible.
- Run the relevant tests after editing.

Required workflow:

1. Read the repository interface and implementation for each target.
2. Read the related entity mappings, query options, domain models, and SQL migrations.
3. Read the application-layer call sites that use the target repositories.
4. Infer the real query and mutation patterns from application usage before writing tests.
5. If tests expose a schema or mapping mismatch, fix the root cause with minimal changes.
6. Run the affected tests and report the outcome.

Hard constraints:

- Do not invent repository call combinations that the application layer does not use unless the repository API clearly guarantees that combination.
- Do not optimize for method-count coverage. Optimize for realistic behavior coverage.
- Do not create one Docker container per test case if a shared package-level setup already exists.
- Keep tests linearly readable: seed a batch of related data first, then verify a sequence of behaviors.
- Prefer validating externally visible behavior over asserting GORM implementation details.
- If an include or join option is only used in one direction in the application layer, test that real usage path instead of forcing unsupported combinations.
- If migrations and repository mappings disagree on column names, stop treating that as a test problem and fix the mismatch.

Testing heuristics:

- For create/list/get/update/delete flows, build one or two scenario tests that cover the full lifecycle.
- For join/include flows, test only the include combinations the application actually requests.
- For existence checks, use the same query option combinations the application uses.
- For soft-delete repositories, verify both direct lookup behavior and list exclusion behavior.
- For invitation or relation tables, seed the minimum required parent rows first so foreign keys remain realistic.

Expected output:

- Implement the test changes directly.
- Keep the final summary short.
- Report any repository bug or schema inconsistency that had to be fixed.
- Include the exact test command that was run and whether it passed.

Preferred summary format:

- Changed files
- Behavioral coverage added
- Root-cause fixes made
- Test command and result
