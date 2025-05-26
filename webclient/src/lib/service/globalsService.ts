import type { AxiosInstance } from "axios";
import { createAxiosInstance, getAuthHeaders } from "./utils";
import { getApiUrl } from "$lib/utils";
import type { Response } from "$lib/types/response";

export type Global = {
  Key: string;
  Value: any;
}

class GlobalsService {
  private axios: AxiosInstance;

  constructor() {
    this.axios = createAxiosInstance();
  }

  async getGlobals(): Promise<[Global[], Error]> {
    try {
      const response = await this.axios.get<Response<Global[]>>(`${getApiUrl()}/globals`, {
        headers: { ...getAuthHeaders() },
      });
      return [response.data.data, null];
    } catch (error: any) {
      return [null, error];
    }
  }

  async updateGlobal(key: string, value: any): Promise<[Global, Error]> {
    try {
      const response = await this.axios.put<Response<Global>>(`${getApiUrl()}/globals?key=${key}`, value, {
        headers: { ...getAuthHeaders() },
      });
      return [response.data.data, null];
    } catch (error: any) {
      return [null, error];
    }
  }
  
  async deleteGlobal(key: string): Promise<Error> {
    try {
      await this.axios.delete(`${getApiUrl()}/globals?key=${key}`, {
        headers: { ...getAuthHeaders() },
      });
      return null;
    } catch (error: any) {
      return error;
    }
  }
}

export const globalsService = new GlobalsService();