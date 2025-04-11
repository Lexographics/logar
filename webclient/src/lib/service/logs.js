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