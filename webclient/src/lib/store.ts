import { LocalStorage, SessionStorage } from "./storage.svelte";
import type { User } from "./types/user";

type NavigationStore = {
  isSidebarLocked: boolean | null;
  isLogsExpanded: boolean;
}
export const navigationStore = new SessionStorage<NavigationStore>("navigation", {
  isSidebarLocked: null,
  isLogsExpanded: false,
})

type SettingsStore = {
  selectedLanguage: string | null;
  currentTheme: string;
}
export const settingsStore = new LocalStorage<SettingsStore>("settings", {
  selectedLanguage: null,
  currentTheme: "light",
})

export type UserStore = {
  user: User | null;
  token: string | null;
}
export const userStore = new LocalStorage<UserStore>("user", {
  user: null,
  token: null,
})

type Model = {
  displayName: string;
  identifier: string;
  icon: string;
}
type ModelsStore = {
  models: Model[];
  lastFetch: number | null;
}

export const modelsStore = new LocalStorage<ModelsStore>("models", {
  models: [],
  lastFetch: null,
})


