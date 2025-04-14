import { PUBLIC_API_URL } from "$env/static/public";
import { userStore } from "$lib/store";
import axios from "axios";
import { checkSession } from "./service";

export async function getLogs(model, cursor = 0, filters = []) {
  try {
    const response = await axios.get(`${PUBLIC_API_URL}/logs/${model}`, {
      params: {
        cursor: cursor,
        filters: JSON.stringify(filters),
      },
      headers: {
        Authorization: `Bearer ${userStore.current.token}`,
      },
    });
    
    return [response.data, null];
  } catch (error) {
    checkSession(error.response);

    console.error("Error fetching logs:", error);
    return [null, error];
  }
}

/**
 * @param {string} model
 * @param {Array} filters
 * @param {Function} onLogsReceived
 * @param {Function} onError
 * @param {Function} onCloseCallback
 * @returns {function} Function to close the SSE connection
 */
export function connectToLogStream(model, filters = [], onLogsReceived, onError, onCloseCallback) {
  let url = `${PUBLIC_API_URL}/logs/${model}/sse?token=${userStore.current.token}`;
  if (filters.length > 0) {
    url += `&filters=${JSON.stringify(filters)}`;
  }

  const eventSource = new EventSource(url);

  eventSource.addEventListener('logs', (event) => {
    try {
      const data = JSON.parse(event.data);
      if (onLogsReceived && typeof onLogsReceived === 'function') {
        onLogsReceived(data);
      }
    } catch (error) {
      console.error("Error parsing SSE log data:", error);
      if (onError && typeof onError === 'function') {
        onError(error);
      }
    }
  });

  eventSource.onerror = (error) => {
    console.error("SSE connection error:", error);
    if (onError && typeof onError === 'function') {
      onError(error);
    }

    eventSource.close();

    if (onCloseCallback && typeof onCloseCallback === 'function') {
      onCloseCallback();
    }
  };

  eventSource.addEventListener('close', () => {
    console.log("SSE connection closed by server");
    if (onCloseCallback && typeof onCloseCallback === 'function') {
      onCloseCallback();
    }
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