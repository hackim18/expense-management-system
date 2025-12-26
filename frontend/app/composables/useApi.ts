type RequestOptions = {
  method?: string;
  body?: unknown;
  headers?: Record<string, string>;
  auth?: boolean;
};

type ApiResponse<T, P = unknown> = {
  message?: string;
  data?: T;
  paging?: P;
  errors?: string;
};

export const useApi = () => {
  const config = useRuntimeConfig();
  const auth = useAuth();
  const apiBase = config.public?.apiBase || "http://localhost:8080";

  const requestWithMeta = async <T, P = unknown>(path: string, options: RequestOptions = {}) => {
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

    return (payload ?? {}) as ApiResponse<T, P>;
  };

  const request = async <T>(path: string, options: RequestOptions = {}) => {
    const payload = await requestWithMeta<T>(path, options);
    if (payload && typeof payload === "object" && "data" in payload) {
      return (payload.data ?? null) as T;
    }
    return payload as unknown as T;
  };

  return { request, requestWithMeta };
};
