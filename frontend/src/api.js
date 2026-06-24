const BASE_URL = import.meta.env.VITE_API_BASE_URL || "http://localhost:8080";
const TOKEN_KEY = "rfid_dor_lock_token";

export function getToken() {
  return localStorage.getItem(TOKEN_KEY);
}

export function setToken(token) {
  localStorage.setItem(TOKEN_KEY, token);
}

export function clearToken() {
  localStorage.removeItem(TOKEN_KEY);
}

async function request(path, options = {}) {
  const token = getToken();
  const res = await fetch(`${BASE_URL}${path}`, {
    headers: {
      "Content-Type": "application/json",
      ...(token ? { Authorization: `Bearer ${token}` } : {}),
    },
    ...options,
  });

  if (res.status === 401) {
    clearToken();
    throw new Error("Sesi berakhir, silakan login kembali");
  }

  if (!res.ok) {
    const body = await res.json().catch(() => ({}));
    throw new Error(body.error || `Request failed: ${res.status}`);
  }

  if (res.status === 204) return null;
  return res.json();
}

export const api = {
  login: (username, password) =>
    request("/api/v1/login", {
      method: "POST",
      body: JSON.stringify({ username, password }),
    }),
  listCards: () => request("/api/v1/cards"),
  createCard: (uid, label) =>
    request("/api/v1/cards", {
      method: "POST",
      body: JSON.stringify({ uid, label }),
    }),
  deleteCard: (id) => request(`/api/v1/cards/${id}`, { method: "DELETE" }),
  listLogs: () => request("/api/v1/logs"),
  unlock: () => request("/api/v1/unlock", { method: "POST" }),
};
