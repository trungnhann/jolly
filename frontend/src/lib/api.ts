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
  const res = await fetch(path, {
    ...init,
    headers: {
      "Content-Type": "application/json",
      ...(token ? { Authorization: `Bearer ${token}` } : {}),
      ...(init?.headers ?? {}),
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
