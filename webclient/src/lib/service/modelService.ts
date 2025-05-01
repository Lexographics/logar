import type { AxiosInstance } from "axios";
import { createAxiosInstance, getAuthHeaders } from "./utils";
import type { Response } from "$lib/types/response";
import { modelsStore } from "$lib/store";

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
      const response = await this.axios.get<Response<Model[]>>("/models", {
        headers: { ...getAuthHeaders() }
      });
      return [response.data.data, null];
    } catch (error) {
      return [null, error];
    }
  }
}

export default new ModelService();