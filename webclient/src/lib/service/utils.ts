import { goto } from "$app/navigation";
import { userStore } from "$lib/store";
import { showToast } from "$lib/toast";
import type { AxiosError, AxiosInstance, AxiosResponse } from "axios";
import axios from "axios";
import { getBasePath } from "$lib/utils";
import { StatusCode, type Response } from "$lib/types/response";

export async function checkSession(response : AxiosResponse) {
  
  const responseData = response.data as Response<any>;
  if (responseData.status_code === StatusCode.SessionExpired) {
    userStore.current = {
      token: null,
    };

    showToast("Session expired, please login again");

    goto(`${getBasePath()}/login`);
  }
}

export function createAxiosInstance() : AxiosInstance {
  const axiosInstance = axios.create({});

  axiosInstance.interceptors.response.use(
    (response: AxiosResponse) => {
      checkSession(response);
      return response;
    },
    (error: AxiosError) => {
      checkSession(error.response);
      return Promise.reject(error);
    }
  );

  return axiosInstance;
}

export function getAuthHeaders(): Record<string, string> {
  return {
    Authorization: `Bearer ${userStore.current.token}`,
  };
}