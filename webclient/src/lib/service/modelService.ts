import type { AxiosInstance } from "axios";
import { createAxiosInstance, getAuthHeaders } from "./utils";
import type { Response } from "$lib/types/response";
import { modelsStore } from "$lib/store";
import { getApiUrl } from "$lib/utils";

type Model = {
  displayName: string;
  icon: string;
  identifier: string;
}

class ModelService {
  private axios: AxiosInstance;

  constructor() {
    this.axios = createAxiosInstance();
  }

  async getModels() {
    if (modelsStore.current.models.length > 0 && modelsStore.current.lastFetch && new Date().getTime() - modelsStore.current.lastFetch < 15000) {
      return [modelsStore.current.models, null];
    }

    try {
      const response = await this.axios.get<Response<Model[]>>(`${getApiUrl()}/models`, {
        headers: { ...getAuthHeaders() }
      });
      
      modelsStore.current.models = response.data.data;
      modelsStore.current.lastFetch = new Date().getTime();
      return [response.data.data, null];
    } catch (error) {
      return [null, error];
    }
  }

  getCachedModels() {
    return modelsStore.current.models;
  }
}

export default new ModelService();