import type { AxiosInstance } from "axios";
import { createAxiosInstance, getAuthHeaders } from "./utils";
import { getApiUrl } from "$lib/utils";
import type { Response } from "$lib/types/response";

export type FeatureFlag = {
  id: number;
  name: string;
  enabled: boolean;
  condition: string;
}

class FeatureFlagsService {
  private axios: AxiosInstance;

  constructor() {
    this.axios = createAxiosInstance();
  }

  async getFeatureFlags(): Promise<[FeatureFlag[], Error | null]> {
    try {
      const response = await this.axios.get<Response<FeatureFlag[]>>(`${getApiUrl()}/feature-flags`, {
        headers: { ...getAuthHeaders() }
      });
      return [response.data.data, null];
    } catch (error) {
      return [null, error];
    }
  }

  async createFeatureFlag(featureFlag: FeatureFlag): Promise<[FeatureFlag, Error | null]> {
    try {
      const formData = new FormData();
      formData.append("name", featureFlag.name);
      formData.append("enabled", featureFlag.enabled.toString());
      formData.append("condition", featureFlag.condition);

      const response = await this.axios.post<Response<FeatureFlag>>(`${getApiUrl()}/feature-flags`, formData, {
        headers: { ...getAuthHeaders() }
      });
      return [response.data.data, null];
    } catch (error) {
      return [null, error];
    }
  }

  async updateFeatureFlag(featureFlag: FeatureFlag): Promise<[FeatureFlag, Error | null]> {
    try {
      const formData = new FormData();
      formData.append("id", featureFlag.id.toString());
      formData.append("name", featureFlag.name);
      formData.append("enabled", featureFlag.enabled.toString());
      formData.append("condition", featureFlag.condition);

      const response = await this.axios.put<Response<FeatureFlag>>(`${getApiUrl()}/feature-flags`, formData, {
        headers: { ...getAuthHeaders() }
      });
      return [response.data.data, null];
    } catch (error) {
      return [null, error];
    }
  }

  async deleteFeatureFlag(id: number): Promise<Error | null> {
    try {
      await this.axios.delete<Response<void>>(`${getApiUrl()}/feature-flags?id=${id}`, {
        headers: { ...getAuthHeaders() }
      });
      return null;
    } catch (error) {
      return error;
    }
  }
}

export default new FeatureFlagsService();