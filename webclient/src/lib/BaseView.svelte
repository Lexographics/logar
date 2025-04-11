<script>
  import { getModels } from "$lib/service/model";
  import Sidebar from "$lib/Sidebar.svelte";
  import Topbar from "./Topbar.svelte";

  let models = getModels();

  let { children } = $props();
</script>


<div style="display: flex; max-height: 100vh; overflow: hidden;">
  <Sidebar {models}></Sidebar>
  
  <div style="width: 100%; height: 100vh;">
    <Topbar></Topbar>
    {@render children?.()}
  </div>
</div>