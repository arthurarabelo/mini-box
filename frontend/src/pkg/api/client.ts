import { useAuthStore } from "../stores/auth";

const BASE_URL = "http://localhost:8080/api";

// Wrapper de fetch que injeta o token automaticamente
// No CERNBox isso é o axios interceptor em client.ts
async function request<T>(path: string, options: RequestInit = {}): Promise<T> {
  const authStore = useAuthStore();

  const headers: Record<string, string> = {
    "Content-Type": "application/json",
    ...(options.headers as Record<string, string>),
  };

  if (authStore.accessToken) {
    headers["Authorization"] = `Bearer ${authStore.accessToken}`;
  }

  const res = await fetch(`${BASE_URL}${path}`, { ...options, headers });

  if (res.status === 401) {
    // Token expirado ou inválido — mesmo comportamento do auth guard do CERNBox
    authStore.setToken(null);
    window.location.href = "/login";
    throw new Error("Não autenticado");
  }

  if (!res.ok) {
    const text = await res.text();
    let message = `HTTP ${res.status}`;
    try {
      const json = JSON.parse(text);
      message = json.error || json.message || message;
    } catch {
      if (text) message = text;
    }
    throw new Error(message);
  }

  const text = await res.text();
  return text ? JSON.parse(text) : null;
}

// API de autenticação
export const authApi = {
  login: (username: string, password: string) =>
    request<{ token: string; name: string; email: string }>("/auth/login", {
      method: "POST",
      body: JSON.stringify({ username, password }),
    }),

  register: (name: string, username: string, email: string, password: string) =>
    request<{
      token: string;
      name: string;
      email: string;
      username: string;
      password: string;
    }>("/auth/register", {
      method: "POST",
      body: JSON.stringify({ name, username, email, password }),
    }),
};

export const graphApi = {
  me: () =>
    request<{
      id: string;
      displayName: string;
      mail: string;
      onPremisesSamAccountName: string;
    }>("/graph/v1.0/me"),
};

// API de arquivos
export interface FileItem {
  name: string;
  type: "file" | "folder";
  size?: number;
  modified: string;
  mimeType?: string;
}

export const filesApi = {
  list: (path: string) =>
    request<FileItem[]>(`/files?path=${encodeURIComponent(path)}`),

  upload: (path: string, file: File) => {
    const form = new FormData();
    form.append("file", file);
    // Para upload não usamos o wrapper (precisa remover Content-Type para multipart)
    const authStore = useAuthStore();
    return fetch(`${BASE_URL}/files/upload?path=${encodeURIComponent(path)}`, {
      method: "POST",
      headers: { Authorization: `Bearer ${authStore.accessToken}` },
      body: form,
    });
  },

  downloadUrl: (path: string) =>
    `${BASE_URL}/files/download?path=${encodeURIComponent(path)}`,
};
