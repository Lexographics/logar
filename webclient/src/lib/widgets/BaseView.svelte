<script>
  import Sidebar from "$lib/widgets/Sidebar.svelte";
  import { onMount } from "svelte";
  import Topbar from "./Topbar.svelte";
  import modelService from "$lib/service/modelService";

  let { children, loaded = true } = $props();

  let models = $state([]);

  onMount(async () => {
    models = modelService.getCachedModels();
    const [modelsData, error] = await modelService.getModels();
    if (error) {
      console.error(error);
    }

    models = modelsData;
  });
</script>

{#if loaded}
<div class="base-container">
  <Sidebar {models}></Sidebar>
  
  <div class="main-content">
    <Topbar></Topbar>
    <div class="content-area scrollbar">
      {@render children?.()}
    </div>
  </div>
</div>
{/if}

<style>
  .base-container {
    display: flex;
    max-width: 100vw;
    max-height: 100dvh;
    overflow: hidden;
    background: var(--background-color);
    color: var(--text-color);
  }

  .main-content {
    height: 100dvh;
    overflow-y: hidden;
    flex-grow: 1;
    min-width: 0;
  }

  .content-area {
    height: calc(100dvh - 60px);
    overflow-y: auto;
    padding: 0;
  }
</style>