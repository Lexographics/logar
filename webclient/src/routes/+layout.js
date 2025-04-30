import { settingsStore } from "$lib/store";
import { setLocale } from "../i18n/i18n-svelte";
import { loadAllLocalesAsync } from "../i18n/i18n-util.async";

export const prerender = true;
export const trailingSlash = 'always';

(async () => {
  await loadAllLocalesAsync();
  setLocale(settingsStore.current.selectedLanguage || "en");
})();
