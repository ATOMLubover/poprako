import { BASE_URL } from "./config";
import type { ApiResponse } from "./types";

// ─── Colour helpers (ANSI) ───────────────────────────────────────────────────

const RESET = "\x1b[0m";
const CYAN = "\x1b[36m";
const GREEN = "\x1b[32m";
const YELLOW = "\x1b[33m";
const RED = "\x1b[31m";
const GREY = "\x1b[90m";
const BOLD = "\x1b[1m";

function colorize(text: string, color: string): string {
  return `${color}${text}${RESET}`;
}

// ─── Logger ──────────────────────────────────────────────────────────────────

export function logStep(phase: string, step: string): void {
  console.log(`\n${colorize(`[${phase}]`, CYAN)}${BOLD} ${step}${RESET}`);
}

export function logOk(label: string, value?: unknown): void {
  const suffix =
    value !== undefined ? ` ${colorize(JSON.stringify(value), GREY)}` : "";
  console.log(`  ${colorize("✓", GREEN)} ${label}${suffix}`);
}

export function logInfo(label: string, value?: unknown): void {
  const suffix =
    value !== undefined ? ` ${colorize(JSON.stringify(value), GREY)}` : "";
  console.log(`  ${colorize("·", YELLOW)} ${label}${suffix}`);
}

export function logError(label: string, detail?: unknown): void {
  console.error(`  ${colorize("✗", RED)} ${label}`);
  if (detail !== undefined) {
    console.error(`    ${colorize(JSON.stringify(detail, null, 2), RED)}`);
  }
}

// ─── HTTP client ─────────────────────────────────────────────────────────────

type Method = "GET" | "POST" | "PUT" | "PATCH" | "DELETE";

interface RequestOptions {
  token?: string;
  query?: Record<string, string | number | boolean | undefined>;
}

/**
 * Core request helper. Throws on HTTP errors or API-level non-200 codes.
 * Returns the parsed `data` field from the envelope, or null for empty responses.
 */
export async function api<T>(
  method: Method,
  path: string,
  body: unknown,
  opts: RequestOptions = {},
): Promise<T> {
  const url = new URL(BASE_URL + path);

  if (opts.query) {
    for (const [k, v] of Object.entries(opts.query)) {
      if (v !== undefined) {
        url.searchParams.set(k, String(v));
      }
    }
  }

  const headers: Record<string, string> = {
    "Content-Type": "application/json",
  };

  if (opts.token) {
    headers["Authorization"] = `Bearer ${opts.token}`;
  }

  const init: RequestInit = {
    method,
    headers,
    body: body != null ? JSON.stringify(body) : undefined,
  };

  const res = await fetch(url.toString(), init);

  // Parse JSON envelope
  let envelope: ApiResponse<T> | null = null;
  try {
    envelope = (await res.json()) as ApiResponse<T>;
  } catch {
    throw new Error(
      `[${method} ${path}] HTTP ${res.status} — response is not JSON`,
    );
  }

  if (!res.ok || (envelope && envelope.code >= 400)) {
    const errMsg = envelope?.message ?? `HTTP ${res.status}`;
    throw new Error(`[${method} ${path}] FAILED (${res.status}): ${errMsg}`);
  }

  return envelope!.data as T;
}

/**
 * GET convenience wrapper.
 */
export async function apiGet<T>(
  path: string,
  opts: RequestOptions = {},
): Promise<T> {
  return api<T>("GET", path, null, opts);
}
