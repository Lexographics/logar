<script>
  import { getModels } from "$lib/service/model";
  import Sidebar from "$lib/Sidebar.svelte";
  import { onMount } from "svelte";
  import Topbar from "./Topbar.svelte";

  let { children } = $props();

  let models = $state([]);

  onMount(async () => {
    const [modelsData, error] = await getModels();
    if (error) {
      console.error(error);
    }

    models = modelsData;
  });
</script>


<div class="base-container">
  <Sidebar {models}></Sidebar>
  
  <div class="main-content">
    <Topbar></Topbar>
    <div class="content-area scrollbar">
      {@render children?.()}
    </div>
  </div>
</div>

<style>
  .base-container {
    display: flex;
    max-width: 100vw;
    max-height: 100vh;
    overflow: hidden;
    background-color: var(--background-color);
    color: var(--text-color);
  }

  .main-content {
    height: 100vh;
    flex-grow: 1;
    min-width: 0;
  }

  .content-area {
    height: calc(100vh - 60px);
    overflow-y: auto;
    padding: 0;
  }
</style>