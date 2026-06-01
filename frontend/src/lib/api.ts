const API_URL = process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080";

async function request<T>(path: string, options?: RequestInit): Promise<T> {
  const res = await fetch(`${API_URL}${path}`, {
    ...options,
    headers: {
      "Content-Type": "application/json",
      ...(options?.headers || {}),
    },
    credentials: options?.credentials ?? "include",
    cache: options?.cache ?? "no-store",
  });
  if (!res.ok) {
    const err = await res.json().catch(() => ({ error: res.statusText }));
    throw new Error(err.error || "Request failed");
  }
  return res.json();
}

export function getToken(): string | null {
  if (typeof window === "undefined") return null;
  return localStorage.getItem("surya_token");
}

export function setToken(token: string) {
  localStorage.setItem("surya_token", token);
}

export function clearToken() {
  localStorage.removeItem("surya_token");
}

/** API list endpoints may return null when empty; normalize to []. */
export function asArray<T>(value: T[] | null | undefined): T[] {
  return Array.isArray(value) ? value : [];
}

function authHeaders(): HeadersInit {
  const token = getToken();
  return token ? { Authorization: `Bearer ${token}` } : {};
}

export const api = {
  get: <T>(path: string) => request<T>(path),
  getAuth: <T>(path: string) =>
    request<T>(path, { headers: { ...authHeaders() } }),
  post: <T>(path: string, body: unknown) =>
    request<T>(path, { method: "POST", body: JSON.stringify(body) }),
  postAuth: <T>(path: string, body?: unknown) =>
    request<T>(path, {
      method: "POST",
      headers: { ...authHeaders() },
      body: body ? JSON.stringify(body) : undefined,
    }),
  putAuth: <T>(path: string, body: unknown) =>
    request<T>(path, {
      method: "PUT",
      headers: { ...authHeaders() },
      body: JSON.stringify(body),
    }),
  deleteAuth: <T>(path: string) =>
    request<T>(path, { method: "DELETE", headers: { ...authHeaders() } }),
  uploadAuth: async <T>(path: string, formData: FormData) => {
    const res = await fetch(`${API_URL}${path}`, {
      method: "POST",
      headers: { ...authHeaders() },
      body: formData,
      credentials: "include",
    });
    if (!res.ok) {
      const err = await res.json().catch(() => ({ error: res.statusText }));
      throw new Error(err.error || "Upload failed");
    }
    return res.json() as Promise<T>;
  },
};

export { API_URL };
