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

const DEFAULT_TIMEOUT = 30_000;
const CHAT_TIMEOUT = 180_000;

export function friendlyError(e: unknown): Error {
  if (e instanceof DOMException && e.name === "AbortError") {
    return new Error("Sunucu zaman aşımına uğradı. Lütfen tekrar deneyin.");
  }
  const msg = (e as Error)?.message || "";
  if (msg === "Failed to fetch" || msg.includes("NetworkError") || msg.includes("Network request failed")) {
    return new Error("Sunucuya bağlanılamadı. İnternet bağlantınızı kontrol edin ve tekrar deneyin.");
  }
  return e as Error;
}

async function fetchWithTimeout(url: string, init: RequestInit, timeoutMs: number): Promise<Response> {
  const controller = new AbortController();
  const timer = setTimeout(() => controller.abort(), timeoutMs);
  try {
    return await fetch(url, { ...init, signal: controller.signal });
  } finally {
    clearTimeout(timer);
  }
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

  if (res.status === 204 || text.trim() === "") {
    return undefined;
  }

  if (!res.ok) {
    try {
      const data = JSON.parse(text);
      throw new Error(data.error || `HTTP ${res.status}`);
    } catch (e: any) {
      throw new Error(
        e?.message?.startsWith?.("HTTP") || !e?.message?.includes?.("JSON")
          ? e.message || `HTTP ${res.status}`
          : `HTTP ${res.status}: ${text.slice(0, 200)}`
      );
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
    try {
      const res = await fetchWithTimeout(`${API_BASE}${path}`, { headers: getHeaders() }, DEFAULT_TIMEOUT);
      return await handleResponse(res);
    } catch (e) {
      throw friendlyError(e);
    }
  },

  async post(path: string, body?: unknown, timeoutMs: number = DEFAULT_TIMEOUT) {
    try {
      const res = await fetchWithTimeout(
        `${API_BASE}${path}`,
        {
          method: "POST",
          headers: getHeaders(),
          body: body ? JSON.stringify(body) : undefined,
        },
        timeoutMs
      );
      return await handleResponse(res);
    } catch (e) {
      throw friendlyError(e);
    }
  },

  async put(path: string, body?: unknown) {
    try {
      const res = await fetchWithTimeout(
        `${API_BASE}${path}`,
        {
          method: "PUT",
          headers: getHeaders(),
          body: body ? JSON.stringify(body) : undefined,
        },
        DEFAULT_TIMEOUT
      );
      return await handleResponse(res);
    } catch (e) {
      throw friendlyError(e);
    }
  },

  async delete(path: string) {
    try {
      const res = await fetchWithTimeout(`${API_BASE}${path}`, {
        method: "DELETE",
        headers: getHeaders(),
      }, DEFAULT_TIMEOUT);
      return await handleResponse(res);
    } catch (e) {
      throw friendlyError(e);
    }
  },
};

export const chatApi = {
  async listConversations() {
    return api.get("/api/v1/chat/conversations");
  },
  async createConversation(title?: string, cvId?: string, jdId?: string, matchId?: string) {
    return api.post("/api/v1/chat/conversations", {
      title: title || "New Chat",
      cv_id: cvId,
      jd_id: jdId,
      match_id: matchId,
    });
  },
  async getConversation(id: string) {
    return api.get(`/api/v1/chat/conversations/${id}`);
  },
  async deleteConversation(id: string) {
    return api.delete(`/api/v1/chat/conversations/${id}`);
  },
  async sendMessage(id: string, content: string) {
    return api.post(`/api/v1/chat/conversations/${id}/messages`, { content }, CHAT_TIMEOUT);
  },
};
