<script>
  import { goto } from '$app/navigation';
  import { base } from '$app/paths';
  import { userStore } from '$lib/store';
  import ThemeProvider from '$lib/widgets/ThemeProvider.svelte';
  import { onMount } from 'svelte';

  let { children } = $props();
  let loaded = $state(false);
  
  onMount(() => {
    if(!userStore.current.token) {
      goto(`${base}/login`);
      return;
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
