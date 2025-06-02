<script>
  import { onMount } from 'svelte';
  import { navigationStore } from '../store';
  import { page } from '$app/stores';
  import { getBasePath } from '$lib/utils';
  import LL from '../../i18n/i18n-svelte';
  import SideBarButton from './SideBarButton.svelte';

  let { models = [] } = $props();

  let loaded = $state(false);
  let isMobile = $state(false);
  let isSidebarOpen = $state(false);
  
  function toggleSidebar() {
    if (isMobile) {
      isSidebarOpen = !isSidebarOpen;
    } else {
      navigationStore.current.isSidebarLocked = !navigationStore.current.isSidebarLocked;
    }
  }

  function toggleLogs() {
    navigationStore.current.isLogsExpanded = !navigationStore.current.isLogsExpanded;
  }

  function handleResize() {
    isMobile = window.innerWidth <= 768;
    if (isMobile) {
      navigationStore.current.isSidebarLocked = false;
    }
    isSidebarOpen = false;
  }

  onMount(() => {
    if (navigationStore.current.isSidebarLocked == null) {
      navigationStore.current.isSidebarLocked = window.innerWidth > 768;
    }

    handleResize();
    window.addEventListener('resize', handleResize);

    setTimeout(() => {
      loaded = true;
    }, 0);

    return () => {
      window.removeEventListener('resize', handleResize);
    };
  });
</script>

{#if isMobile && isSidebarOpen}
  <div class="backdrop" onclick={toggleSidebar} onkeydown={e => e.key === 'Enter' && toggleSidebar()} tabindex="0" role="button" aria-label="Close menu"></div>
{/if}

<div class="sidebar-container">
  {#if isMobile && !isSidebarOpen}
    <button class="mobile-menu-button" onclick={toggleSidebar} onkeydown={e => e.key === 'Enter' && toggleSidebar()} aria-label="Toggle menu">
      <i class="fas fa-bars"></i>
    </button>
  {/if}

  <div class="sidebar {isMobile ? (isSidebarOpen ? 'locked' : '') : (navigationStore.current.isSidebarLocked || !loaded ? 'locked' : '')} {isMobile ? 'mobile' : ''} {isSidebarOpen ? 'open' : ''}">
    <div class="sidebar-header">
      <h2>
        <a class="linktext" href={`${getBasePath()}/`} style="display: flex; align-items: center;">
          <img src={`${getBasePath()}/favicon.png`} alt="Logar" style="width: 2rem; height: 2rem; margin-right: 8px;">
          <span class="text">Logar</span>
        </a>
      </h2>

      {#if !isMobile}
        <button aria-label="Toggle sidebar" class="toggle-button" onclick={toggleSidebar}>
          <i class="fas {navigationStore.current.isSidebarLocked ? 'fa-lock' : 'fa-unlock'}"></i>
        </button>
      {/if}
    </div>
    
    <nav>
      <ul>
        <SideBarButton href={`/dashboard`} icon="fas fa-home" text={$LL.dashboard.title()}/>
        <SideBarButton icon="fas fa-list-alt" text={$LL.logs.title()} onclick={toggleLogs} active={$page.url.pathname.startsWith(`${getBasePath()}/logs`)}>
          {#snippet end()}
            <i class="fas {navigationStore.current.isLogsExpanded ? 'fa-chevron-down' : 'fa-chevron-right'} chevron"></i>
          {/snippet}

          <ul class="scrollbar submenu {navigationStore.current.isLogsExpanded ? 'expanded' : ''}">
            {#each models as model}
              <SideBarButton href={`/logs?model=${model.identifier}`} icon={model.icon ? model.icon : 'fa-solid fa-cube'} text={model.displayName || model.identifier} active={$page.url.pathname.startsWith(`${getBasePath()}/logs`) && $page.url.searchParams.get('model') === model.identifier} />
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
  .sidebar-container {
    position: relative;
  }

  .backdrop {
    position: fixed;
    top: 0;
    left: 0;
    right: 0;
    bottom: 0;
    background-color: rgba(0, 0, 0, 0.5);
    z-index: 9;
  }

  .mobile-menu-button {
    position: fixed;
    top: 0;
    left: 15px;
    height: 60px;
    display: flex;
    align-items: center;
    z-index: 9;
    background: transparent;
    border: none;
    color: var(--sidebar-text);
    width: 40px;
    cursor: pointer;
    transition: transform 0.3s ease;
  }

  .mobile-menu-button i {
    font-size: 1.2rem;
    margin: 0;
  }

  .sidebar {
    height: 100dvh;
    background-color: var(--sidebar-background);
    padding: 20px 15px;
    box-shadow: 2px 0 5px var(--shadow-color);
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

  .sidebar.mobile {
    position: fixed;
    top: 0;
    left: 0;
    transform: translateX(-100%);
    width: 250px;
  }

  .sidebar.mobile.open {
    transform: translateX(0);
  }

  .sidebar .text {
    opacity: 0;
    visibility: hidden;
    transition: opacity 0.2s ease, visibility 0.2s ease;
  }

  .sidebar:hover .text, .sidebar.locked .text, .sidebar.mobile .text {
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
    align-items: center;
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
    margin: 0;
  }

  .chevron {
    margin-left: auto;
    font-size: 0.8em;
    transition: transform 0.3s ease;
  }

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
    max-height: 300px;
  }

  @media (max-width: 768px) {
    .sidebar {
      width: 250px;
    }

    .sidebar:hover {
      width: 250px;
    }
  }
</style>