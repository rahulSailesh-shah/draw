import { jwtDecode } from "jwt-decode";
import { authClient } from "./auth-client";

export type APIResponse<T> = {
  data: T | null;
  error: string | null;
  status: number;
};

export class ApiError extends Error {
  public status: number;
  constructor(message: string, status: number) {
    super(message);
    this.name = "ApiError";
    this.status = status;
  }
}

class ApiClient {
  private readonly baseURL: string;
  private readonly headers: Record<string, string>;
  private cachedToken: string | null = null;

  constructor(baseURL: string) {
    this.baseURL = baseURL;
    this.cachedToken = null;
    this.headers = {
      "Content-Type": "application/json",
    };
  }

  private isTokenValid() {
    if (!this.cachedToken) {
      return false;
    }

    const decodedToken = jwtDecode(this.cachedToken);
    if (!decodedToken.exp) {
      return false;
    }
    const expirationTime = decodedToken.exp;
    const currentTime = Date.now() / 1000;
    return expirationTime > currentTime;
  }

  private async getToken() {
    if (this.isTokenValid()) {
      return this.cachedToken;
    }

    const { data, error } = await authClient.token();
    if (error || !data?.token) {
      return null;
    }
    this.cachedToken = data.token;
    return data.token;
  }

  public async get<T>(
    url: string,
    params?: Record<string, string | number | boolean | undefined>
  ): Promise<APIResponse<T>> {
    const token = await this.getToken();
    if (!token) {
      return {
        error: "Failed to fetch token",
        status: 401,
        data: null,
      };
    }

    let fullUrl = `${this.baseURL}${url}`;

    if (params && Object.keys(params).length > 0) {
      const searchParams = new URLSearchParams();

      Object.entries(params).forEach(([key, value]) => {
        if (value !== undefined && value !== null && value !== "") {
          searchParams.append(key, String(value)); // Convert to string here
        }
      });

      const queryString = searchParams.toString();
      if (queryString) {
        fullUrl += `?${queryString}`;
      }
    }

    const response = await fetch(fullUrl, {
      credentials: "include",
      method: "GET",
      mode: "cors",
      headers: {
        ...this.headers,
        Authorization: `Bearer ${token}`,
      },
    });

    const data = await response.json();

    if (!response.ok) {
      return {
        error: data.error || data.message || "Failed to fetch data",
        status: response.status,
        data: null,
      };
    }

    return {
      data: data.data,
      error: null,
      status: response.status,
    };
  }

  public async post<T>(url: string, body: unknown): Promise<APIResponse<T>> {
    const token = await this.getToken();
    if (!token) {
      return {
        error: "Failed to fetch token",
        status: 401,
        data: null,
      };
    }
    const response = await fetch(`${this.baseURL}${url}`, {
      method: "POST",
      headers: {
        ...this.headers,
        Authorization: `Bearer ${token}`,
      },
      body: JSON.stringify(body),
    });
    const data = await response.json();
    if (!response.ok) {
      return {
        error: data.error || data.message || "Request failed",
        status: response.status,
        data: null,
      };
    }
    return {
      data: data.data,
      error: null,
      status: response.status,
    };
  }

  public async delete(url: string): Promise<APIResponse<void>> {
    const token = await this.getToken();
    if (!token) {
      return {
        error: "Failed to fetch token",
        status: 401,
        data: null,
      };
    }
    const response = await fetch(`${this.baseURL}${url}`, {
      method: "DELETE",
      headers: {
        ...this.headers,
        Authorization: `Bearer ${token}`,
      },
    });
    const data = await response.json();
    if (!response.ok) {
      return {
        error: data.error || data.message || "Request failed",
        status: response.status,
        data: null,
      };
    }
    return {
      data: null,
      error: null,
      status: response.status,
    };
  }

  public async put<T>(url: string, body: unknown): Promise<APIResponse<T>> {
    const token = await this.getToken();
    if (!token) {
      return {
        error: "Failed to fetch token",
        status: 401,
        data: null,
      };
    }
    const response = await fetch(`${this.baseURL}${url}`, {
      method: "PUT",
      headers: {
        ...this.headers,
        Authorization: `Bearer ${token}`,
      },
      body: JSON.stringify(body),
    });
    const data = await response.json();
    if (!response.ok) {
      return {
        error: data.error || data.message || "Request failed",
        status: response.status,
        data: null,
      };
    }
    return {
      data: data.data,
      error: null,
      status: response.status,
    };
  }
}

const serverUrl = import.meta.env.VITE_SERVER_URL as string;
export const apiClient = new ApiClient(serverUrl);
