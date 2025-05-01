<script>
  import { onMount } from 'svelte';
  import { navigationStore } from '../store';
  import { page } from '$app/stores';
  import { base } from '$app/paths';
  import LL from '../../i18n/i18n-svelte';
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
        <a class="linktext" href={`${base}/`}>
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
        <li><a class="link" href={`${base}/`} class:active={$page.url.pathname.startsWith(`${base}/dashboard`)}><i class="fas fa-home"></i>  <span class="text">{$LL.dashboard.title()}</span></a></li>
        <li>
          <a href={"javascript:void(0)"} onclick={toggleLogs} class="menu-item link" class:active={$page.url.pathname.startsWith(`${base}/logs`)}>
            <i class="fas fa-list-alt"></i>
            <span class="text">{$LL.logs.title()}</span>
            <i class="fas {navigationStore.current.isLogsExpanded ? 'fa-chevron-down' : 'fa-chevron-right'} chevron"></i>
          </a>
          <ul class="scrollbar submenu {navigationStore.current.isLogsExpanded ? 'expanded' : ''}">
            {#each models as model}
              <li>
                <a class="link submenu-item" href={`${base}/logs?model=${model.identifier}`} class:active={$page.url.pathname.startsWith(`${base}/logs`) && $page.url.searchParams.get('model') === model.identifier}>
                  <i class="{model.icon ? model.icon : 'fa-solid fa-cube'}"></i>
                  <span class="text">{model.displayName || model.identifier}</span>
                </a>
              </li>
            {/each}
          </ul>
        </li>
        <li><a class="link" href={`${base}/analytics`} class:active={$page.url.pathname.startsWith(`${base}/analytics`)}><i class="fas fa-chart-bar"></i> <span class="text">{$LL.analytics.title()}</span></a></li>
        <li><a class="link" href={`${base}/actions`} class:active={$page.url.pathname.startsWith(`${base}/actions`)}><i class="fa-solid fa-server"></i> <span class="text">{$LL.remote_actions.title()}</span></a></li>
        <li><a class="link" href={`${base}/user`} class:active={$page.url.pathname.startsWith(`${base}/user`)}><i class="fa-solid fa-users"></i> <span class="text">User / Sessions</span></a></li>
        <li><a class="link" href={`${base}/settings`} class:active={$page.url.pathname.startsWith(`${base}/settings`)}><i class="fas fa-cog"></i> <span class="text">{$LL.settings.title()}</span></a></li>
        <li><a class="link" href={`${base}/help`} class:active={$page.url.pathname.startsWith(`${base}/help`)}><i class="fas fa-question-circle"></i> <span class="text">{$LL.help.title()}</span></a></li>
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
  .sidebar li {
    margin: 10px 0;
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
  
  .link {
    text-decoration: none;
    color: var(--sidebar-text);
    display: flex;
    align-items: center;
    padding: 8px 10px;
    border-radius: 4px;
    transition: background-color 0.2s;
    white-space: nowrap;
  }
  
  .link:hover {
    background-color: var(--sidebar-hover);
  }

  .link.active {
    background-color: var(--sidebar-active);
    font-weight: bold;
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

  .menu-item {
    display: flex;
    justify-content: space-between;
    align-items: center;
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

  .submenu li {
    margin: 5px 0;
    padding-left: 0;
    transition: padding-left 0.3s ease;
  }

  .sidebar.locked .submenu li,
  .sidebar:hover .submenu li {
    padding-left: 1rem;
  }

  .submenu a {
    font-size: 0.9em;
    padding: 5px 10px;
  }
</style>