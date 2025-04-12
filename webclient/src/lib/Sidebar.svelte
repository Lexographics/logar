<script>
  import { onMount } from 'svelte';
  import { navigationStore } from './store';

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
        <a class="linktext" href="/">
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
        <li><a class="link" href="/"><i class="fas fa-home"></i>  <span class="text">Dashboard</span></a></li>
        <li>
          <a href={"javascript:void(0)"} onclick={toggleLogs} class="menu-item link">
            <i class="fas fa-list-alt"></i>
            <span class="text">Logs</span>
            <i class="fas {navigationStore.current.isLogsExpanded ? 'fa-chevron-down' : 'fa-chevron-right'} chevron"></i>
          </a>
          <ul class="submenu {navigationStore.current.isLogsExpanded ? 'expanded' : ''}">
            {#each models as model}
              <li>
                <a class="link submenu-item" href={`/logs?model=${model.identifier}`}>
                  <i class="fas fa-cube"></i>
                  <span class="text">{model.name || model.identifier}</span>
                </a>
              </li>
            {/each}
          </ul>
        </li>
        <li><a class="link" href="/analytics"><i class="fas fa-chart-bar"></i> <span class="text">Analytics</span></a></li>
        <li><a class="link" href="/settings"><i class="fas fa-cog"></i> <span class="text">Settings</span></a></li>
        <li><a class="link" href="/help"><i class="fas fa-question-circle"></i> <span class="text">Help</span></a></li>
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
    background-color: #f4f4f4;
    padding: 20px;
    box-shadow: 2px 0 5px rgba(0, 0, 0, 0.1);
    /*
    position: fixed;
    top: 0;
    left: 0;
    */
    z-index: 1000;
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
    z-index: 1001;
    border: none;
    border-radius: 50%;
    background: #f4f4f4;
    cursor: pointer;
    width: 30px;
    height: 30px;
    box-shadow: 0 0 5px rgba(0, 0, 0, 0.2);
    display: flex;
    align-items: center;
    justify-content: center;
  }

  .toggle-button:hover {
    background-color: #e0e0e0;
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
    border-bottom: 1px solid #ddd;
    display: flex;
    justify-content: space-between;
  }
  
  .sidebar-footer {
    margin-top: 20px;
    padding-top: 10px;
    border-top: 1px solid #ddd;
    font-size: 0.8em;
    overflow-x: visible;
    white-space: nowrap;
  }

  .linktext {
    text-decoration: none;
    color: #333;
  }
  
  .link {
    text-decoration: none;
    color: #333;
    display: flex;
    align-items: center;
    padding: 8px 10px;
    border-radius: 4px;
    transition: background-color 0.2s;
    white-space: nowrap;
  }
  
  .link:hover {
    background-color: #e0e0e0;
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

  .submenu {
    max-height: 0;
    overflow: hidden;
    transition: max-height 0.3s ease-out;
    padding-left: 20px;
  }

  .submenu.expanded {
    max-height: 150px;
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