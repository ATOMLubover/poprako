#!/usr/bin/env node
/**
 * check-type-comments.js
 *
 * Checks that every `type` declaration in Go source files has at least one
 * comment line immediately above it.
 *
 * Usage:
 *   node check-type-comments.js <path1> [path2 ...]
 *   pnpm check-types ./internal/event ./internal/domain
 *
 * Paths can be files or directories. Directories are walked recursively and
 * all *.go files within them are checked.
 *
 * Exit code:
 *   0 — all type declarations are properly commented
 *   1 — one or more violations found (details printed to stdout)
 */

import { readFileSync, statSync, readdirSync } from "fs";
import { join, resolve } from "path";

// ---------------------------------------------------------------------------
// File collection
// ---------------------------------------------------------------------------

/**
 * `invocationCwd` returns the directory from which the user originally invoked
 * pnpm/npm. When pnpm runs a lifecycle script it changes CWD to the package
 * directory and sets `INIT_CWD` to the original working directory. Falling
 * back to `process.cwd()` covers direct `node` invocations.
 * @returns {string}
 */
function invocationCwd() {
  return process.env.INIT_CWD ?? process.cwd();
}

/**
 * `collectGoFiles` recursively collects all *.go file paths under `root`.
 * If `root` itself is a .go file it is returned directly.
 * Relative paths are resolved against the directory from which the user ran
 * pnpm (i.e. `INIT_CWD`), not the package directory.
 * @param {string} root - Absolute or relative path to a file or directory.
 * @returns {string[]} Sorted list of absolute .go file paths.
 */
function collectGoFiles(root) {
  const abs = resolve(invocationCwd(), root);
  const stat = statSync(abs);

  if (stat.isFile()) {
    return abs.endsWith(".go") ? [abs] : [];
  }

  if (stat.isDirectory()) {
    const entries = readdirSync(abs, { withFileTypes: true });
    return entries.flatMap((e) => {
      const full = join(abs, e.name);
      if (e.isDirectory()) return collectGoFiles(full);
      if (e.isFile() && e.name.endsWith(".go")) return [full];
      return [];
    });
  }

  return [];
}

// ---------------------------------------------------------------------------
// Per-file analysis
// ---------------------------------------------------------------------------

/**
 * `isCommentLine` returns true when `line` is a Go comment line (// or /*).
 * Blank lines and closing-comment lines (*‌/) are NOT treated as comments.
 * @param {string} line - Raw line string (already trimmed).
 * @returns {boolean}
 */
function isCommentLine(line) {
  return line.startsWith("//") || line.startsWith("/*");
}

/**
 * `isTypeDeclaration` returns true when `line` starts a Go type declaration,
 * e.g. `type Foo struct`, `type Bar interface`, `type Baz string`.
 * Lines inside generic constraints or composite literals that happen to
 * contain the word "type" are excluded by requiring it to be the first token.
 * @param {string} line - Raw line string (already trimmed).
 * @returns {boolean}
 */
function isTypeDeclaration(line) {
  return /^type\s+\w/.test(line);
}

/**
 * `checkFile` analyses a single Go source file and returns all violations.
 * A violation is a `type` declaration that has no comment line immediately
 * above it (ignoring blank lines between the comment and the declaration is
 * intentionally NOT done — the comment must be on the line directly above).
 *
 * @param {string} filePath - Absolute path to the .go file.
 * @returns {{ file: string, line: number, text: string }[]} List of violations.
 */
function checkFile(filePath) {
  const source = readFileSync(filePath, "utf8");
  const lines = source.split("\n");
  const violations = [];

  for (let i = 0; i < lines.length; i++) {
    const trimmed = lines[i].trim();

    // Skip lines that are not type declarations.
    if (!isTypeDeclaration(trimmed)) continue;

    // Check the line directly above (1-indexed: line i is line number i+1).
    const prevLine = i > 0 ? lines[i - 1].trim() : "";
    if (!isCommentLine(prevLine)) {
      violations.push({
        file: filePath,
        line: i + 1, // 1-indexed for human readability
        text: lines[i].trimEnd(),
      });
    }
  }

  return violations;
}

// ---------------------------------------------------------------------------
// Entry point
// ---------------------------------------------------------------------------

/**
 * `main` is the CLI entry point.
 * Reads paths from `process.argv`, collects Go files, checks each one, and
 * prints a report. Exits with code 1 if any violations were found.
 */
function main() {
  const args = process.argv.slice(2);

  if (args.length === 0) {
    console.error("Usage: check-type-comments.js <path1> [path2 ...]");
    process.exit(1);
  }

  // Collect all .go files from all supplied paths.
  const files = args.flatMap((arg) => collectGoFiles(arg));

  if (files.length === 0) {
    console.error("No .go files found in the supplied paths.");
    process.exit(1);
  }

  // Run the check on every file.
  const allViolations = files.flatMap((f) => checkFile(f));

  if (allViolations.length === 0) {
    console.log(`✓ Checked ${files.length} file(s) — no violations found.`);
    process.exit(0);
  }

  // Print violations grouped by file for readability.
  console.error(
    `✗ Found ${allViolations.length} type declaration(s) missing a comment:\n`
  );

  let currentFile = null;
  for (const v of allViolations) {
    if (v.file !== currentFile) {
      currentFile = v.file;
      console.error(`  ${v.file}`);
    }
    console.error(`    line ${v.line}: ${v.text}`);
  }

  process.exit(1);
}

main();
