<script lang="ts">
  import ThemeProvider from '$lib/widgets/ThemeProvider.svelte';
  import { onMount } from 'svelte';
  import { loadAllLocalesAsync } from '../i18n/i18n-util.async';
  import { settingsStore } from '$lib/store';
  import { locales } from '../i18n/i18n-util';
  import { setLocale } from '../i18n/i18n-svelte';
  import { setMomentLocale } from '$lib/moment';
  import type { Response } from '$lib/types/response';
  import { PUBLIC_API_URL } from '$env/static/public';
  import type { Locales } from '../i18n/i18n-types';
  import axios from 'axios';

  let { children } = $props();
  let loaded = $state(false);
  
  onMount(async () => {
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

    loaded = true;
  });
</script>

<ThemeProvider>
{#if loaded}
    {@render children?.()}
  {:else}
    <div class="loading-container">
      <div class="spinner"></div>
    </div>
  {/if}
</ThemeProvider>

<style>
  .loading-container {
    position: fixed;
    top: 0;
    left: 0;
    width: 100%;
    height: 100%;
    display: flex;
    justify-content: center;
    align-items: center;
    background-color: var(--background-color);
  }
  
  .spinner {
    width: 50px;
    height: 50px;
    border: 5px solid var(--border-color);
    border-top: 5px solid var(--primary-color);
    border-radius: 50%;
    animation: spin 1s linear infinite;
  }

  @keyframes spin {
    0% { transform: rotate(0deg); }
    100% { transform: rotate(360deg); }
  }
</style>
