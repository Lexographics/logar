import { SessionStorage } from "./storage.svelte";

export const navigationStore = new SessionStorage("navigation", {
  isSidebarLocked: null,
  isLogsExpanded: false,
})

