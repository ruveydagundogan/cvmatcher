function getAPIBase(): string {
  if (process.env.NEXT_PUBLIC_API_URL) {
    return process.env.NEXT_PUBLIC_API_URL;
  }
  return "http://localhost:8080";
}

export const API_BASE = getAPIBase();

function getHeaders(): Record<string, string> {
  const headers: Record<string, string> = { "Content-Type": "application/json" };
  const token = localStorage.getItem("token");
  if (token) headers["Authorization"] = `Bearer ${token}`;
  return headers;
}

async function handleResponse(res: Response) {
  const text = await res.text();

  if (res.status === 401) {
    localStorage.removeItem("token");
    localStorage.removeItem("userId");
    localStorage.removeItem("userName");
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
