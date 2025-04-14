import { PUBLIC_API_URL } from "$env/static/public";
import { modelsStore, userStore } from "$lib/store";
import axios from "axios";
import { checkSession } from "./service";


export async function getModels() {
  if (modelsStore.current.models.length > 0 && modelsStore.current.lastFetch && new Date() - new Date(modelsStore.current.lastFetch) < 15000) {
    return [modelsStore.current.models, null];
  }

  try {
    const response = await axios.get(`${PUBLIC_API_URL}/models`, {
      headers: {
        Authorization: `Bearer ${userStore.current.token}`,
      },
    });

    modelsStore.current = {
      models: response.data,
      lastFetch: new Date(),
    };

    return [response.data, null];
  } catch (error) {
    checkSession(error.response);
    return [null, error];
  }
}

