import { settingsStore } from "$lib/store";
import axios from "axios";
import { setLocale } from "../i18n/i18n-svelte";
import { loadAllLocalesAsync } from "../i18n/i18n-util.async";
import { PUBLIC_API_URL } from "$env/static/public";
import { locales } from "../i18n/i18n-util";
import type { Response } from "$lib/types/response";
import type { Locales } from "../i18n/i18n-types";
import { setMomentLocale } from "$lib/moment";
import { browser } from "$app/environment";

export const prerender = true;
export const trailingSlash = 'always';

if (browser) {

  (async () => {
    await loadAllLocalesAsync();
  
    if (settingsStore.current.selectedLanguage && locales.includes(settingsStore.current.selectedLanguage)) {
      setLocale(settingsStore.current.selectedLanguage);
      await setMomentLocale(settingsStore.current.selectedLanguage);
  
    } else {
      setLocale("en");
      await setMomentLocale("en");
  
      axios.get<Response<string>>(`${PUBLIC_API_URL}/language`).then(async (res) => {
        settingsStore.current.selectedLanguage = res.data.data;
        setLocale(res.data.data as Locales);
        await setMomentLocale(res.data.data);
      });
    }
  })();
}
