<script>
  import { getModels } from "$lib/service/model";
  import Sidebar from "$lib/Sidebar.svelte";
  import Topbar from "./Topbar.svelte";

  let models = getModels();

  let { children } = $props();
</script>


<div style="display: flex; max-width: 100vw; max-height: 100vh; overflow: hidden;">
  <Sidebar {models}></Sidebar>
  
  <div style="height: 100vh; flex-grow: 1; min-width: 0;">
    <Topbar></Topbar>
    <div style="height: calc(100vh - 60px); overflow-y: scroll;">
      {@render children?.()}
    </div>
  </div>
</div>