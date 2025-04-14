import { PUBLIC_API_URL } from "$env/static/public";
import { showToast } from "$lib/toast";
import axios from "axios";

export async function login(username, password) {
  try {
    const form = new FormData();
    form.append('username', username);
    form.append('password', password);

    const response = await axios.post(`${PUBLIC_API_URL}/auth/login`, form);

    return [response.data, null];
  } catch (error) {
    showToast(error.response.data);
    return [null, error];
  }
}
