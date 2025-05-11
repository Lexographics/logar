import type { AxiosInstance } from "axios";
import { createAxiosInstance, getAuthHeaders } from "./utils";
import { getApiUrl } from "$lib/utils";
import type { Response } from "$lib/types/response";

export type PageStats = {
  path: string;
  visits: number;
  percentage: number;
};

export type AnalyticsSummary = {
  total_visits: number;
  unique_visitors: number;
  active_visitors: number;
  error_rate: number;
  average_latency_ms: number;
  p95_latency_ms: number;
  p99_latency_ms: number;
  total_bytes_sent: number;
  total_bytes_recv: number;
  top_pages: PageStats[];
  os_usage: Record<string, number>;
  browser_usage: Record<string, number>;
  referer_usage: Record<string, number>;
  instance_stats: Record<string, number>;
};

class AnalyticsService {
  private axios: AxiosInstance;

  constructor() {
    this.axios = createAxiosInstance();
  }

  async getAnalytics(): Promise<[AnalyticsSummary, Error]> {
    try {
      const response = await this.axios.get<Response<AnalyticsSummary>>(`${getApiUrl()}/analytics`, {
        headers: { ...getAuthHeaders() }
      });
      return [response.data.data, null];
    } catch (error: any) {
      return [null, error];
    }
  }
}

export default new AnalyticsService();
