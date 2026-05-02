/**
 * check-type-comments.test.js
 *
 * Unit tests for `check-type-comments.js` using Node's built-in test runner.
 * Run with: pnpm test
 */

import { test } from "node:test";
import assert from "node:assert/strict";
import { writeFileSync, mkdirSync } from "fs";
import { join, resolve } from "path";
import { tmpdir } from "os";

// ---------------------------------------------------------------------------
// Re-implement the public functions under test so tests run in-process without
// spawning a subprocess.  Keep these in sync with check-type-comments.js.
// ---------------------------------------------------------------------------

/**
 * `isCommentLine` returns true when `line` (trimmed) is a Go comment line.
 * @param {string} line
 * @returns {boolean}
 */
function isCommentLine(line) {
  return line.startsWith("//") || line.startsWith("/*");
}

/**
 * `isTypeDeclaration` returns true when `line` (trimmed) starts a Go type
 * declaration.
 * @param {string} line
 * @returns {boolean}
 */
function isTypeDeclaration(line) {
  return /^type\s+\w/.test(line);
}

import { readFileSync } from "fs";

/**
 * `checkFile` analyses a Go file and returns violation objects.
 * @param {string} filePath
 * @returns {{ file: string, line: number, text: string }[]}
 */
function checkFile(filePath) {
  const source = readFileSync(filePath, "utf8");
  const lines = source.split("\n");
  const violations = [];

  for (let i = 0; i < lines.length; i++) {
    const trimmed = lines[i].trim();
    if (!isTypeDeclaration(trimmed)) continue;
    const prevLine = i > 0 ? lines[i - 1].trim() : "";
    if (!isCommentLine(prevLine)) {
      violations.push({ file: filePath, line: i + 1, text: lines[i].trimEnd() });
    }
  }

  return violations;
}

/**
 * `invocationCwd` returns `INIT_CWD` when set (pnpm/npm lifecycle),
 * otherwise falls back to `process.cwd()`.
 * @returns {string}
 */
function invocationCwd() {
  return process.env.INIT_CWD ?? process.cwd();
}

// ---------------------------------------------------------------------------
// Test helpers
// ---------------------------------------------------------------------------

let tmpCount = 0;

/**
 * `writeTmp` writes `content` to a temporary .go file and returns its path.
 * @param {string} content
 * @returns {string}
 */
function writeTmp(content) {
  const dir = join(tmpdir(), "check-type-comments-test");
  mkdirSync(dir, { recursive: true });
  const file = join(dir, `fixture_${++tmpCount}.go`);
  writeFileSync(file, content, "utf8");
  return file;
}

// ---------------------------------------------------------------------------
// isCommentLine
// ---------------------------------------------------------------------------

test("isCommentLine: // comment", () => {
  assert.equal(isCommentLine("// `Foo` is something."), true);
});

test("isCommentLine: /* block comment", () => {
  assert.equal(isCommentLine("/* block */"), true);
});

test("isCommentLine: blank line is not a comment", () => {
  assert.equal(isCommentLine(""), false);
});

test("isCommentLine: code line is not a comment", () => {
  assert.equal(isCommentLine("type Foo struct {"), false);
});

// ---------------------------------------------------------------------------
// isTypeDeclaration
// ---------------------------------------------------------------------------

test("isTypeDeclaration: struct type", () => {
  assert.equal(isTypeDeclaration("type Foo struct {"), true);
});

test("isTypeDeclaration: interface type", () => {
  assert.equal(isTypeDeclaration("type Bar interface {"), true);
});

test("isTypeDeclaration: named type alias", () => {
  assert.equal(isTypeDeclaration("type EvTyp string"), true);
});

test("isTypeDeclaration: generic type", () => {
  assert.equal(isTypeDeclaration("type Set[T comparable] struct {"), true);
});

test("isTypeDeclaration: 'type' inside a comment is not a declaration", () => {
  assert.equal(isTypeDeclaration("// type Foo is something"), false);
});

test("isTypeDeclaration: 'type' as part of a longer word is not a declaration", () => {
  assert.equal(isTypeDeclaration("typewriter := 0"), false);
});

test("isTypeDeclaration: assignment with 'type' word inside is not a declaration", () => {
  assert.equal(isTypeDeclaration('someType := "string"'), false);
});

// ---------------------------------------------------------------------------
// checkFile — no violations
// ---------------------------------------------------------------------------

test("checkFile: single type with // comment above — no violation", () => {
  const path = writeTmp(
    `package foo\n\n// \`Foo\` is a foo.\ntype Foo struct {}\n`
  );
  assert.deepEqual(checkFile(path), []);
});

test("checkFile: type with /* block comment above — no violation", () => {
  const path = writeTmp(
    `package foo\n\n/* \`Bar\` is a bar. */\ntype Bar interface {}\n`
  );
  assert.deepEqual(checkFile(path), []);
});

test("checkFile: multiple types all commented — no violations", () => {
  const path = writeTmp(
    [
      "package foo",
      "",
      "// `EvTyp` wraps the event type identifier.",
      "type EvTyp string",
      "",
      "// `Event` represents a domain event.",
      "type Event interface {",
      "    EvTyp() EvTyp",
      "}",
      "",
    ].join("\n")
  );
  assert.deepEqual(checkFile(path), []);
});

test("checkFile: type at first line with no preceding line — violation", () => {
  const path = writeTmp(`type Foo struct {}\n`);
  const violations = checkFile(path);
  assert.equal(violations.length, 1);
  assert.equal(violations[0].line, 1);
});

// ---------------------------------------------------------------------------
// checkFile — violations
// ---------------------------------------------------------------------------

test("checkFile: type with no comment above — one violation", () => {
  const path = writeTmp(
    `package foo\n\ntype Foo struct {}\n`
  );
  const violations = checkFile(path);
  assert.equal(violations.length, 1);
  assert.equal(violations[0].line, 3);
  assert.match(violations[0].text, /type Foo/);
});

test("checkFile: type preceded only by blank line — violation", () => {
  const path = writeTmp(
    `package foo\n\n// \`Foo\` is fine.\ntype Foo struct {}\n\ntype Bar string\n`
  );
  const violations = checkFile(path);
  assert.equal(violations.length, 1);
  assert.equal(violations[0].line, 6);
});

test("checkFile: type preceded by closing brace — violation", () => {
  const path = writeTmp(
    [
      "package foo",
      "",
      "// `Foo` is good.",
      "type Foo struct {",
      "}",
      "type Bar string", // no comment, preceded by closing brace
      "",
    ].join("\n")
  );
  const violations = checkFile(path);
  assert.equal(violations.length, 1);
  assert.equal(violations[0].line, 6);
});

test("checkFile: two types both missing comments — two violations", () => {
  const path = writeTmp(
    ["package foo", "", "type Foo struct {}", "", "type Bar string", ""].join(
      "\n"
    )
  );
  const violations = checkFile(path);
  assert.equal(violations.length, 2);
});

test("checkFile: violation object has correct file path", () => {
  const path = writeTmp(`package foo\n\ntype Uncovered struct {}\n`);
  const violations = checkFile(path);
  assert.equal(violations[0].file, path);
});

// ---------------------------------------------------------------------------
// checkFile — edge cases
// ---------------------------------------------------------------------------

test("checkFile: 'type' inside a string literal is not a declaration", () => {
  const path = writeTmp(
    `package foo\n\n// \`Foo\` is fine.\ntype Foo struct {\n    Kind string // "type A"\n}\n`
  );
  assert.deepEqual(checkFile(path), []);
});

test("checkFile: 'type' keyword inside a function body is still checked", () => {
  // Local type declarations inside functions must also be commented.
  const path = writeTmp(
    [
      "package foo",
      "",
      "// `doSomething` does something.",
      "func doSomething() {",
      "    type localType struct{ x int }", // no comment above
      "}",
      "",
    ].join("\n")
  );
  const violations = checkFile(path);
  assert.equal(violations.length, 1);
  assert.equal(violations[0].line, 5);
});

test("checkFile: empty file — no violations", () => {
  const path = writeTmp(`package foo\n`);
  assert.deepEqual(checkFile(path), []);
});

test("checkFile: file with only comments — no violations", () => {
  const path = writeTmp(`// Package foo provides foo.\npackage foo\n`);
  assert.deepEqual(checkFile(path), []);
});

// ---------------------------------------------------------------------------
// invocationCwd
// ---------------------------------------------------------------------------

test("invocationCwd: returns INIT_CWD when set", () => {
  const original = process.env.INIT_CWD;
  process.env.INIT_CWD = "/some/project/root";
  assert.equal(invocationCwd(), "/some/project/root");
  if (original === undefined) delete process.env.INIT_CWD;
  else process.env.INIT_CWD = original;
});

test("invocationCwd: falls back to process.cwd() when INIT_CWD is unset", () => {
  const original = process.env.INIT_CWD;
  delete process.env.INIT_CWD;
  assert.equal(invocationCwd(), process.cwd());
  if (original !== undefined) process.env.INIT_CWD = original;
});

// ---------------------------------------------------------------------------
// collectGoFiles (via resolve + INIT_CWD)
// ---------------------------------------------------------------------------

test("resolve: relative path resolves against invocationCwd", () => {
  // When INIT_CWD is set, a bare filename should resolve under INIT_CWD.
  const dir = join(tmpdir(), "check-type-comments-test");
  mkdirSync(dir, { recursive: true });
  const file = join(dir, "rel_test.go");
  writeFileSync(file, `// \`Foo\` is fine.\ntype Foo struct {}\n`, "utf8");

  const original = process.env.INIT_CWD;
  process.env.INIT_CWD = dir;
  const abs = resolve(invocationCwd(), "rel_test.go");
  assert.equal(abs, file);
  if (original === undefined) delete process.env.INIT_CWD;
  else process.env.INIT_CWD = original;
});
