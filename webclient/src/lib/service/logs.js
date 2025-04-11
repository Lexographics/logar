import axios from "axios";


export async function getLogs(model, cursor = 0, filters = []) {
  try {
    let url = `http://localhost:3000/logger/${model}/json?cursor=${cursor}`;
    if (filters.length > 0) {
      const filterParams = filters.map(filter => `${filter.type}=${filter.value}`).join('&');
      url += `&${filterParams}`;
    }

    const response = await axios.get(url, );

    return [response.data, null];
  } catch (error) {
    console.error("Error fetching logs:", error);
    return [null, error];
  }
}

/**
 * Connect to the SSE logs endpoint and receive real-time log updates
 * @param {string} model
 * @param {Array} filters
 * @param {Function} onLogsReceived Callback function when new logs are received
 * @param {Function} onError Callback function when an error occurs
 * @returns {function} Function to close the SSE connection
 */
export function connectToLogStream(model, filters = [], onLogsReceived, onError) {
  let url = `http://localhost:3000/logger/${model}/sse`;
  if (filters.length > 0) {
    const filterParams = filters.map(filter => `${filter.type}=${filter.value}`).join('&');
    url += `?${filterParams}`;
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
  };

  return () => {
    if (eventSource) {
      eventSource.close();
    }
  };
}