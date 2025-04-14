import { LocalStorage, SessionStorage } from "./storage.svelte";

export const navigationStore = new SessionStorage("navigation", {
  isSidebarLocked: null,
  isLogsExpanded: false,
})

export const settingsStore = new LocalStorage("settings", {
  selectedLanguage: "en",
  currentTheme: "light",
  emailNotifications: true,
  pushNotifications: false,
})

export const userStore = new LocalStorage("user", {
  user: null,
  token: null,
})

export const modelsStore = new LocalStorage("models", {
  models: [],
  lastFetch: null,
})
