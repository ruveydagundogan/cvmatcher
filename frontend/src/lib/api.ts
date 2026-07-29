import { API_BASE_URL } from "./config";

export const API_BASE = API_BASE_URL;

function getToken(): string | null {
  if (typeof window === "undefined") return null;
  return localStorage.getItem("token");
}

function clearAuth(): void {
  if (typeof window === "undefined") return;
  localStorage.removeItem("token");
  localStorage.removeItem("userId");
  localStorage.removeItem("userName");
}

function getHeaders(): Record<string, string> {
  const headers: Record<string, string> = { "Content-Type": "application/json" };
  const token = getToken();
  if (token) headers["Authorization"] = `Bearer ${token}`;
  return headers;
}

async function handleResponse(res: Response) {
  const text = await res.text();

  if (res.status === 401) {
    clearAuth();
    if (typeof window !== "undefined") {
      window.location.href = "/";
    }
    throw new Error("Session expired. Please login again.");
  }

  if (!res.ok) {
    try {
      const data = JSON.parse(text);
      throw new Error(data.error || `HTTP ${res.status}`);
    } catch (e: any) {
      throw new Error(e.message || `HTTP ${res.status}: ${text.slice(0, 100)}`);
    }
  }

  let data: any;
  try {
    data = JSON.parse(text);
  } catch {
    throw new Error(`invalid JSON response: ${text.slice(0, 100)}`);
  }

  if (!data.success) {
    throw new Error(data.error || "request failed");
  }

  return data.data;
}

export const api = {
  async get(path: string) {
    const res = await fetch(`${API_BASE}${path}`, { headers: getHeaders() });
    return handleResponse(res);
  },

  async post(path: string, body?: unknown) {
    const res = await fetch(`${API_BASE}${path}`, {
      method: "POST",
      headers: getHeaders(),
      body: body ? JSON.stringify(body) : undefined,
    });
    return handleResponse(res);
  },

  async put(path: string, body?: unknown) {
    const res = await fetch(`${API_BASE}${path}`, {
      method: "PUT",
      headers: getHeaders(),
      body: body ? JSON.stringify(body) : undefined,
    });
    return handleResponse(res);
  },

  async delete(path: string) {
    const res = await fetch(`${API_BASE}${path}`, {
      method: "DELETE",
      headers: getHeaders(),
    });
    return handleResponse(res);
  },
};
