import { getStoredToken } from "./session";

export type ErrorResponse = {
  code?: string;
  message?: string;
};

export async function apiRequest<T>(
  path: string,
  init?: RequestInit,
): Promise<T> {
  const token = getStoredToken();
  const headers: Record<string, string> = {};
  if (token) {
    headers["Authorization"] = `Bearer ${token}`;
  }
  if (!(init?.body instanceof FormData)) {
    headers["Content-Type"] = "application/json";
  }

  const res = await fetch(path, {
    ...init,
    headers: {
      ...headers,
      ...(init?.headers as Record<string, string> ?? {}),
    },
  });

  const text = await res.text();
  const json = text ? safeJsonParse(text) : null;

  if (!res.ok) {
    const msg =
      (json && typeof json === "object" && "message" in json && json.message
        ? String(json.message)
        : null) ?? `Request failed (${res.status})`;
    throw new Error(msg);
  }

  return json as T;
}

function safeJsonParse(input: string) {
  try {
    return JSON.parse(input);
  } catch {
    return null;
  }
}
