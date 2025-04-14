import { goto } from "$app/navigation";
import { base } from "$app/paths";
import { userStore } from "$lib/store";
import { showToast } from "$lib/toast";

export async function checkSession(response) {
  if (response.status === 401) {
    userStore.current = {
      token: null,
    };

    showToast("Session expired, please login again");

    goto(`${base}/login`);
  }
}
