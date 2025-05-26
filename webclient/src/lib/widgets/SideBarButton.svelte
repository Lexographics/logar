<script>
  import { page } from "$app/stores";
  import { getBasePath } from "$lib/utils";

  let { href, icon, text, children, onclick, tooltip = '', end, submenu = false, active = null } = $props();
</script>

<li>
  <a class="link {submenu ? 'submenu-item' : ''}" href={`${getBasePath()}${href}`} class:active={active === null ? $page.url.pathname.startsWith(`${getBasePath()}${href}`) : active} onclick={onclick}>
    <i class={icon}></i>
    <span class="text">{text}</span>
    {#if tooltip}
      <span class="tooltip">{tooltip}</span>
    {/if}
    {#if end}
      {@render end()}
    {/if}
  </a>

  {@render children?.()}
</li>

<style>
  .tooltip {
    background-color: #c45d3e;
    color: #ffffff;
    border-radius: 12px;
    font-size: 0.7rem;
    margin-left: 5px;
    padding: 5px;
    min-width: 1.1rem;
    height: 1.1rem;
    display: inline-flex;
    align-items: center;
    justify-content: center;
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

  :global(.sidebar) li {
    margin: 10px 0;
  }

  :global(.sidebar) .text, :global(.sidebar) .tooltip {
    opacity: 0;
    visibility: hidden;
    transition: opacity 0.2s ease, visibility 0.2s ease;
  }

  :global(.sidebar:hover) .text, :global(.sidebar.locked) .text,
  :global(.sidebar:hover) .tooltip, :global(.sidebar.locked) .tooltip {
    opacity: 1;
    visibility: visible;
  }

  :global(.submenu) li {
    margin: 5px 0;
    padding-left: 0;
    transition: padding-left 0.3s ease;
  }

  :global(.sidebar.locked) :global(.submenu) li,
  :global(.sidebar:hover) :global(.submenu) li {
    padding-left: 1rem;
  }

  :global(.submenu) a {
    font-size: 0.9em;
    padding: 5px 10px;
  }
</style>