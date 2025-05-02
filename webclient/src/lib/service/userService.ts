import { getApiUrl } from "$lib/utils";
import { type AxiosInstance } from "axios";
import type { User } from "$lib/types/user";
import type { Session } from "$lib/types/session";
import { createAxiosInstance, getAuthHeaders } from "./utils";
import type { Response } from "$lib/types/response";


type LoginResponse = {
  token: string;
  user: User;
}

class UserService {
  private axios: AxiosInstance;

  constructor() {
    this.axios = createAxiosInstance();
  }

  async login(username: string, password: string): Promise<[LoginResponse, Error]> {
    try {
      const form = new FormData();
      form.append("username", username);
      form.append("password", password);

      const response = await this.axios.post<Response<LoginResponse>>(`${getApiUrl()}/auth/login`, form);
      return [response.data.data, null];
    } catch (error: any) {
      return [null, error];
    }
  }

  async logout(): Promise<Error | null> {
    try {
      await this.axios.post(`${getApiUrl()}/auth/logout`, null, {
        headers: { ...getAuthHeaders() },
      });
      return null;
    } catch (error: any) {
      return error;
    }
  }

  async revokeSession(sessionId: number): Promise<Error | null> {
    try {
      const form = new FormData();
      form.append("session_id", sessionId.toString());
      await this.axios.post(`${getApiUrl()}/auth/revoke-session`, form, {
        headers: { ...getAuthHeaders() },
      });
      return null;
    } catch (error: any) {
      return error;
    }
  }

  async getActiveSessions(): Promise<[Session[], Error]> {
    try {
      const response = await this.axios.get<Response<Session[]>>(`${getApiUrl()}/auth/sessions`, {
        headers: { ...getAuthHeaders() },
      });
      return [response.data.data, null];
    } catch (error: any) {
      return [null, error];
    }
  }

  async updateUser(displayName: string): Promise<[User, Error]> {
    try {
      const form = new FormData();
      form.append("display_name", displayName);

      const response = await this.axios.put<Response<User>>(`${getApiUrl()}/user`, form, {
        headers: { ...getAuthHeaders() },
      });
      return [response.data.data, null];
    } catch (error: any) {
      return [null, error];
    }
  }

  async getAllUsers(): Promise<[User[], Error]> {
    try {
      const response = await this.axios.get<Response<User[]>>(`${getApiUrl()}/user`, {
        headers: { ...getAuthHeaders() },
      });
      return [response.data.data, null];
    } catch (error: any) {
      return [null, error];
    }
  }

  async createUser(username: string, password: string, displayName: string, isAdmin: boolean): Promise<[User, Error]> {
    try {
      const form = new FormData();
      form.append("username", username);
      form.append("password", password);
      form.append("display_name", displayName);
      form.append("is_admin", isAdmin.toString());

      const response = await this.axios.post<Response<User>>(`${getApiUrl()}/user`, form, {
        headers: { ...getAuthHeaders() },
      });
      return [response.data.data, null];
    } catch (error: any) {
      return [null, error];
    }
  }
}

export default new UserService();