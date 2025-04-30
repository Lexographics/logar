import { PUBLIC_API_URL } from "$env/static/public";
import { userStore } from "$lib/store";
import { showToast } from "$lib/toast";
import axios from "axios";
import { checkSession } from "./service";

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

export async function logout() {
  try {
    const response = await axios.post(`${PUBLIC_API_URL}/auth/logout`, {}, {
      headers: {
        Authorization: `Bearer ${userStore.current.token}`,
      },
    });
    return null;
  } catch (error) {
    showToast(error?.response?.data || error.message);
    return error;
  }
}

export async function revokeSession(sessionId) {
  try {
    const form = new FormData();
    form.append('session_id', sessionId);
    const response = await axios.post(`${PUBLIC_API_URL}/auth/revoke-session`, form, {
      headers: {
        Authorization: `Bearer ${userStore.current.token}`,
      },
    });
    return null;
  } catch (error) {
    checkSession(error);
    return error;
  }
}

export async function getActiveSessions() {
  try {
    const response = await axios.get(`${PUBLIC_API_URL}/auth/sessions`, {
      headers: {
        Authorization: `Bearer ${userStore.current.token}`,
      },
    });
    return [response.data, null];
  } catch (error) {
    checkSession(error);
    return [null, error];
  }
}

export async function updateUser({displayName}) {
  try {
    const form = new FormData();
    if (displayName) {
      form.append('display_name', displayName);
    }

    const response = await axios.put(`${PUBLIC_API_URL}/auth/user`, form, {
      headers: {
        Authorization: `Bearer ${userStore.current.token}`,
      },
    });

    userStore.current.user = response.data;
    return null;
  } catch (error) {
    checkSession(error);
    return error;
  }
}
