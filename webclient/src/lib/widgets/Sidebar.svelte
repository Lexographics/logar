<script>
  import { onMount } from 'svelte';
  import { navigationStore } from '../store';
  import { page } from '$app/stores';
  import { getBasePath } from '$lib/utils';
  import LL from '../../i18n/i18n-svelte';
  import SideBarButton from './SideBarButton.svelte';

  let { models = [] } = $props();

  let loaded = $state(false);
  
  function toggleSidebar() {
    navigationStore.current.isSidebarLocked = !navigationStore.current.isSidebarLocked;
  }

  function toggleLogs() {
    navigationStore.current.isLogsExpanded = !navigationStore.current.isLogsExpanded;
  }

  onMount(() => {
    if (navigationStore.current.isSidebarLocked == null) {
      navigationStore.current.isSidebarLocked = window.innerWidth > 768; // Open by default on larger screens
    }

    setTimeout(() => {
      loaded = true;
    }, 0);
  });
</script>

<div class="sidebar-container">
  <div class="sidebar {navigationStore.current.isSidebarLocked || !loaded ? 'locked' : ''}">
    <div class="sidebar-header">
      <h2>
        <a class="linktext" href={`${getBasePath()}/`}>
          <i class="fas fa-truck"></i>
          <span class="text">Logar</span>
        </a>
      </h2>

      <button aria-label="Toggle sidebar" class="toggle-button" onclick={toggleSidebar}>
        <i class="fas {navigationStore.current.isSidebarLocked ? 'fa-lock' : 'fa-unlock'}"></i>
      </button>
    </div>
    
    <nav>
      <ul>
        <SideBarButton href={`/dashboard`} icon="fas fa-home" text={$LL.dashboard.title()}/>
        <SideBarButton href={`javascript:void(0)`} icon="fas fa-list-alt" text={$LL.logs.title()} onclick={toggleLogs} active={$page.url.pathname.startsWith(`${getBasePath()}/logs`)}>
          {#snippet end()}
            <i class="fas {navigationStore.current.isLogsExpanded ? 'fa-chevron-down' : 'fa-chevron-right'} chevron"></i>
          {/snippet}

          <ul class="scrollbar submenu {navigationStore.current.isLogsExpanded ? 'expanded' : ''}">
            {#each models as model}
              <SideBarButton href={`${getBasePath()}/logs?model=${model.identifier}`} icon={model.icon ? model.icon : 'fa-solid fa-cube'} text={model.displayName || model.identifier} active={$page.url.pathname.startsWith(`${getBasePath()}/logs`) && $page.url.searchParams.get('model') === model.identifier} />
            {:else}
              <p style="color: var(--sidebar-text); text-align: center; font-size: 0.8rem; padding: 10px;">No models found</p>
            {/each}
          </ul>
        </SideBarButton>
        <SideBarButton href={`/analytics`} icon="fas fa-chart-bar" text={$LL.analytics.title()}/>
        <SideBarButton href={`/actions`} icon="fa-solid fa-server" text={$LL.remote_actions.title()}/>
        <SideBarButton href={`/featureflags`} icon="fa-solid fa-flag" text={$LL.featureflags.title()}/>
        <SideBarButton href={`/globals`} icon="fa-solid fa-earth-asia" text="Globals"/>
        <SideBarButton href={`/user`} icon="fa-solid fa-users" text={$LL.user_sessions.title()}/>
        <SideBarButton href={`/settings`} icon="fas fa-cog" text={$LL.settings.title()}/>
        <SideBarButton href={`/help`} icon="fas fa-question-circle" text={$LL.help.title()}/>
      </ul>
    </nav>
    
    <div class="sidebar-footer">
      <p><span class="text">© 2025 Logar</span></p>
    </div>
  </div>
</div>

<style>
  .sidebar {
    height: 100vh;
    background-color: var(--sidebar-background);
    padding: 20px 15px;
    box-shadow: 2px 0 5px var(--shadow-color);
    /*
    position: fixed;
    top: 0;
    left: 0;
    */
    z-index: 10;
    overflow-x: hidden;
    transition: all 0.3s ease;
    width: 250px;
  }

  .sidebar {
    width: 70px;
  }

  .sidebar:hover, .sidebar.locked {
    width: 250px;
  }

  .sidebar .text {
    opacity: 0;
    visibility: hidden;
    transition: opacity 0.2s ease, visibility 0.2s ease;
  }

  .sidebar:hover .text, .sidebar.locked .text {
    opacity: 1;
    visibility: visible;
  }

  .toggle-button {
    z-index: 11;
    border: none;
    border-radius: 50%;
    background: var(--sidebar-background);
    color: var(--sidebar-text);
    cursor: pointer;
    width: 30px;
    height: 30px;
    box-shadow: 0 0 5px var(--shadow-color);
    display: flex;
    align-items: center;
    justify-content: center;
  }

  .toggle-button:hover {
    background-color: var(--sidebar-hover);
  }

  .toggle-button i {
    margin: 0;
  }

  .sidebar ul {
    list-style: none;
    padding: 0;
  }
  
  .sidebar-header {
    margin-bottom: 20px;
    padding-bottom: 10px;
    border-bottom: 1px solid var(--border-color);
    padding-left: 5px;
    display: flex;
    justify-content: space-between;
  }
  
  .sidebar-footer {
    margin-top: 20px;
    padding-top: 10px;
    border-top: 1px solid var(--border-color);
    font-size: 0.8em;
    overflow-x: visible;
    white-space: nowrap;
  }

  .linktext {
    text-decoration: none;
    color: var(--sidebar-text);
  }
  
  i {
    min-width: 20px;
    text-align: center;
    margin-right: 10px;
  }
  
  .sidebar-header h2 {
    display: flex;
    align-items: center;
    white-space: nowrap;
  }

  .chevron {
    margin-left: auto;
    font-size: 0.8em;
    transition: transform 0.3s ease;
  }

  /*
  .rotate-90 {
    transform: rotate(90deg);
  }
  */

  .submenu {
    max-height: 0;
    overflow-x: hidden;
    overflow-y: hidden;
    transition: max-height 0.3s ease-out;
    padding-left: 20px;
  }

  .submenu:has(:nth-child(9)) {
    overflow-y: scroll;
  }

  .submenu.expanded {
    /* overflow-y: auto; */
    max-height: 300px;
  }
</style>