<script>
  import { onMount } from 'svelte';
  import { settingsStore } from '$lib/store';

  let { children } = $props();

  function setTheme(theme) {
    document.documentElement.setAttribute('data-theme', theme);
  }

  $effect(() => {
    setTheme(settingsStore.current.currentTheme || 'dark');
  });

  onMount(() => {
    if (!settingsStore.current.currentTheme) {
      settingsStore.current.currentTheme = 'dark';
    }
    setTheme(settingsStore.current.currentTheme);
  });
</script>

{@render children()}