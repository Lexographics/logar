import type { AxiosInstance } from "axios";
import { createAxiosInstance, getAuthHeaders } from "./utils";
import type { Response } from "$lib/types/response";
import { getApiUrl } from "$lib/utils";

type ActionArg = {
  kind: string;
  type: string;
}

type Action = {
  path: string;
  description: string;
  args: ActionArg[];
}

class ActionService {
  private axios: AxiosInstance;

  constructor() {
    this.axios = createAxiosInstance();
  }

  async getActions() {
    try {
      const response = await this.axios.get<Response<Action[]>>(`${getApiUrl()}/actions`, {
        headers: { ...getAuthHeaders() },
      });
  
      return [response.data.data, null];
    } catch (error) {
      return [null, error];
    }
  }

  async invokeAction(path: string, args: any[]) {
    try {
      const stringArgs = args.map(arg => {
        if (typeof arg === 'boolean') {
          return arg ? 'true' : 'false';
        }

        if (typeof arg === 'object' && arg !== null) {
          try {
            return JSON.stringify(arg);
          } catch (e) {
            console.warn("Could not JSON stringify argument:", arg, e);
            return String(arg);
          }
        }
        return String(arg);
      });

      const response = await this.axios.post<Response<any>>(`${getApiUrl()}/actions/invoke`, {
        path: path,
        args: stringArgs
      }, {
        headers: { ...getAuthHeaders() },
      });

      return [response.data.data, null];
    } catch (error) {
      return [null, error];
    }
  }
}

export default new ActionService();