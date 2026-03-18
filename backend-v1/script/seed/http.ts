import { CONFIG } from "./config";

let token: string | null = null;

export function setToken(t: string) {
  token = t;
}

export async function request(method: string, path: string, body?: any) {
  const res = await fetch(`${CONFIG.apiBase}${path}`, {
    method,
    headers: {
      "Content-Type": "application/json",
      ...(token ? { Authorization: `Bearer ${token}` } : {}),
    },
    body: body ? JSON.stringify(body) : undefined,
  });

  if (!res.ok) {
    const text = await res.text();
    throw new Error(`${method} ${path} failed: ${text}`);
  }

  if (res.status === 204) return null;

  // Try to parse JSON. The server uses a wrapper: { code, message, data }
  // Unwrap `data` on success, and surface server error messages when `code` != 0.
  const text = await res.text();
  if (!text) return null;

  let parsed: any;
  try {
    parsed = JSON.parse(text);
  } catch (err) {
    // Not JSON — return raw text
    return text;
  }

  if (
    parsed &&
    typeof parsed === "object" &&
    "code" in parsed &&
    "message" in parsed
  ) {
    // Server uses HTTP-like codes (2xx) for success, but some endpoints
    // may use 0. Accept 0 or any 2xx code as success.
    const success =
      parsed.code === 0 ||
      (typeof parsed.code === "number" &&
        parsed.code >= 200 &&
        parsed.code < 300);
    if (!success) {
      throw new Error(`${method} ${path} failed: ${JSON.stringify(parsed)}`);
    }

    return parsed.data ?? null;
  }

  return parsed;
}
