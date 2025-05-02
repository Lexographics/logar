import type { AxiosInstance } from "axios";
import { createAxiosInstance, getAuthHeaders } from "./utils";
import type { Response } from "$lib/types/response";
import { getApiUrl } from "$lib/utils";
import { userStore } from "$lib/store";

type Log = {
  id: number;
  message: string;
  timestamp: string;
  level: string;
  model: string;
}

type LogFilter = {
  field: string;
  operator: string;
  value: string;
}

class LogService {
  private axios: AxiosInstance;

  constructor() {
    this.axios = createAxiosInstance();
  }

  async getLogs(model: string, cursor: number, filters: LogFilter[]) {
    try {
      const response = await this.axios.get<Response<Log[]>>(`${getApiUrl()}/logs/${model}`, {
        params: {
          cursor,
          filters: JSON.stringify(filters),
        },
        headers: { ...getAuthHeaders() },
      });

      return [response.data.data, null];
    } catch (error) {
      return [null, error];
    }
  }

  connectToLogStream(model: string, filters: LogFilter[], onLogsReceived: (logs: Log[]) => void, onError: (error: Error) => void, onCloseCallback: () => void) {
    let url = `${getApiUrl()}/logs/${model}/sse?token=${userStore.current.token}`;
    if (filters.length > 0) {
      url += `&filters=${JSON.stringify(filters)}`;
    }

    const eventSource = new EventSource(url);

    eventSource.addEventListener('logs', (event) => {
      try {
        const data = JSON.parse(event.data);
        onLogsReceived(data);

      } catch (error) {
        console.error("Error parsing SSE log data:", error);
        onError(error);
      }
    });

    eventSource.onerror = (error) => {
      console.error("SSE connection error:", error);
      onError(new Error(error.toString()));

      eventSource.close();
      onCloseCallback();
    };

    eventSource.addEventListener('close', () => {
      console.log("SSE connection closed by server");
      onCloseCallback();
      eventSource.close();
    });

    eventSource.onopen = () => {
      
    };

    return () => {
      if (eventSource) {
        eventSource.close();
      }
    };
  }
}

export default new LogService();