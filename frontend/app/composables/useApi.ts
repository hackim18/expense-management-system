type RequestOptions = {
  method?: string;
  body?: unknown;
  headers?: Record<string, string>;
  auth?: boolean;
};

export const useApi = () => {
  const config = useRuntimeConfig();
  const auth = useAuth();
  const apiBase = config.public?.apiBase || "http://localhost:8080";

  const request = async <T>(path: string, options: RequestOptions = {}) => {
    if (import.meta.client) {
      auth.init();
    }

    const headers: Record<string, string> = {
      "Content-Type": "application/json",
      ...options.headers,
    };

    if (options.auth !== false && auth.token.value) {
      headers.Authorization = `Bearer ${auth.token.value}`;
    }

    const response = await fetch(`${apiBase}${path}`, {
      method: options.method || "GET",
      headers,
      body: options.body ? JSON.stringify(options.body) : undefined,
    });

    const isJson = response.headers.get("content-type")?.includes("application/json") ?? false;
    const payload = isJson ? await response.json() : null;

    if (!response.ok) {
      const message = payload?.errors || payload?.message || response.statusText;
      throw new Error(message);
    }

    return (payload?.data ?? payload) as T;
  };

  return { request };
};
