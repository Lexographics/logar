<script>
  import { browser } from "$app/environment";
  import { page } from "$app/stores";
  import BaseView from "$lib/BaseView.svelte";
  import { getLogs, connectToLogStream } from "$lib/service/logs";
  import { getModels } from "$lib/service/model";
  import Sidebar from "$lib/Sidebar.svelte";
  import { SessionStorage } from "$lib/storage.svelte";
  import moment from "moment";
  import { mount, onMount } from "svelte";
  import { fade, fly } from 'svelte/transition';

  let models = getModels();

  let model = $state("");
  let filters = $state([]);
  
  $effect(() => {
    model = $page.url.searchParams.get("model") || "";
    logs = [];
    loading = false;
    lastCursor = 0;
    hasMore = true;
  });
  
  function addMessageFilter(filter) {
    filters = [...filters, {
      type: "message",
      value: filter,
    }];
  }

  function getSeverityClass(severity) {
    switch (severity) {
      case 1:
        return "log";
      case 2:
        return "info";
      case 3:
        return "warning";
      case 4:
        return "error";
      case 5:
        return "fatal";
      case 6:
        return "trace";
      default:
        return "";
    }
  }

  let logs = $state([]);
  let loading = $state(false);
  let lastCursor = $state(0);
  let hasMore = $state(true);
  let isStreamConnected = $state(false);
  let streamCleanup = null;

  let currentFilterInput = $state("");

  function removeFilter(index) {
    filters = filters.filter((_, i) => i !== index);
    logs = [];
    lastCursor = 0;
    hasMore = true;
    loadMore();
    if (isStreamConnected) {
      disconnectLogStream();
      setupLogStream();
    }
  }

  function handleAddFilter() {
    if (currentFilterInput.trim()) {
      addMessageFilter(currentFilterInput);
      currentFilterInput = "";
      logs = [];
      lastCursor = 0;
      hasMore = true;
      loadMore();
      if (isStreamConnected) {
        disconnectLogStream();
        setupLogStream();
      }
    }
  }

  async function loadMore() {
    if (loading || !hasMore) return;
    
    const modelName = model;

    loading = true;
    const [data, err] = await getLogs(modelName, lastCursor, filters);

    if (err) {
      console.error(err);
      loading = false;
      return;
    }

    if (data?.Logs?.length < 20) {
      hasMore = false;
    }

    if (data?.Logs?.length > 0) {
      logs = [...logs, ...data.Logs];
      lastCursor = data.Logs[data.Logs.length - 1].ID;
      
      if (!isStreamConnected) {
        setupLogStream();
      }
    }

    loading = false;
  }

  function handleNewLogs(data) {
    if (data.Logs && data.Logs.length > 0) {
      logs = [...data.Logs, ...logs];
    }
  }

  function setupLogStream() {
    if (streamCleanup) {
      streamCleanup();
    }
    
    streamCleanup = connectToLogStream(
      model, 
      filters, 
      handleNewLogs,
      (error) => console.error("Stream error:", error)
    );
    
    isStreamConnected = true;
  }

  function disconnectLogStream() {
    if (streamCleanup) {
      streamCleanup();
      streamCleanup = null;
    }
    isStreamConnected = false;
  }

  $effect(() => {
    model = $page.url.searchParams.get("model") || "";
    logs = [];
    loading = false;
    lastCursor = 0;
    hasMore = true;
    
    disconnectLogStream();
  });

  let observerTarget;
  onMount(() => {
    const observer = new IntersectionObserver(
      (entries) => {
        if (entries[0].isIntersecting) {
          loadMore();
        }
      },
      { threshold: 0.1 },
    );

    if (observerTarget) {
      observer.observe(observerTarget);
    }

    return () => {
      observer.disconnect();
      disconnectLogStream();
    };
  });

</script>

<BaseView>
  <h1 style="padding: 1rem 2rem;">Logs</h1>

  <div class="search-container">
    <input
      type="text"
      bind:value={currentFilterInput}
      placeholder="Add filter..."
      class="search-input"
      onkeydown={(e) => e.key === 'Enter' && handleAddFilter()}
    />
    <button class="filter-button" onclick={handleAddFilter}>
      Add Filter
    </button>
  </div>

  {#if filters.length || 0 > 0}
    <div class="active-filters">
      {#each filters as filter, i}
        <div class="filter-tag">
          <span>{filter.value}</span>
          <button class="remove-filter" onclick={() => removeFilter(i)}>×</button>
        </div>
      {/each}
    </div>
  {/if}

  <div style="max-width: 100%; height: 75%; overflow-y: scroll;">
    <table style="max-width: 100%;">
      <thead>
        <tr>
          <th style="text-align: center;">ID</th>
          <th style="text-align: center;"><i class="fa-solid fa-signal"></i></th>
          <th style="text-align: center;">Timestamp</th>
          <th style="text-align: left;">Message</th>
          <th style="text-align: center;">Category</th>
        </tr>
      </thead>
      <tbody>
        {#key [model, JSON.stringify(filters)]}
          {#each logs as log (log.ID)}
            <tr
              class="row {getSeverityClass(log.Severity)}"
              in:fly={{ x: -50, duration: 300 }}
              out:fade={{ duration: 200 }}
            >
              <td style="width: 1%;">{log.ID}</td>
              <td style="width: 1%;">{getSeverityClass(log.Severity).toUpperCase()}</td>
              <td style="width: 50ch;">{moment(log.CreatedAt).format("DD-MM-YYYY HH:mm:ss.SSS")}</td>
              <td style="width: 70%; word-break: break-all;">{log.Message}</td>
              <td style="width: 1%;">{log.Category}</td>
            </tr>
          {/each}
        {/key}
      </tbody>
    </table>

    <div bind:this={observerTarget} class="loader">
      {#if loading}
        <div class="loading-spinner"></div>
      {:else if !hasMore}
        <p>No more logs</p>
      {/if}
    </div>
  </div>
</BaseView>

<style>
  table {
    width: 100%;
    border-collapse: collapse;
  }

  th,
  td {
    padding: 0.5rem;
    text-align: left;
    border-bottom: 1px solid #ddd;
  }

  th {
    background-color: #f4f4f4;
  }

  .row.log {
    --row-bg: #d1e7dd;
    --row-striped-bg: #c7dbd2;
  }

  .row.info {
    --row-bg: #c5daf2;
    --row-striped-bg: #d8e6f8;
  }

  .row.warn {
    --row-bg: #fff3cd;
    --row-striped-bg: #f7e6b5;
  }

  .row.error {
    --row-bg: #f8d7da;
    --row-striped-bg: #eccccf;
  }

  .row.crit {
    --row-bg: #d98a8a;
    --row-striped-bg: #eba8a8;
  }

  .row.fatal {
    --row-bg: #d68c8c;
    --row-striped-bg: #e6b1b1;

    font-weight: bold;
  }

  .row.debug {
    --row-bg: #c7d9f0;
    --row-striped-bg: #d9e7f7;
  }

  .row.trace {
    --row-bg: #e6e6e6;
    --row-striped-bg: #d9d9d9;
  }

  .row {
    background-color: var(--row-bg);
  }

  .row:nth-child(odd) {
    background-color: var(--row-striped-bg);
  }

  th {
    position: sticky;
    top: 0;
    background: linear-gradient(to bottom, #f4f4f4, #e0e0e0);
  }

  .loader {
    padding: 2rem;
    text-align: center;
  }

  .loading-spinner {
    width: 40px;
    height: 40px;
    margin: 0 auto;
    border: 3px solid #f3f3f3;
    border-top: 3px solid #3498db;
    border-radius: 50%;
    animation: spin 1s linear infinite;
  }

  @keyframes spin {
    0% {
      transform: rotate(0deg);
    }
    100% {
      transform: rotate(360deg);
    }
  }

  .search-container {
    padding: 0 2rem 1rem 2rem;
    display: flex;
    gap: 1rem;
  }

  .search-input {
    width: 100%;
    padding: 0.5rem;
    font-size: 1rem;
    border: 1px solid #ddd;
    border-radius: 4px;
    box-sizing: border-box;
  }

  .search-input:focus {
    outline: none;
    border-color: #3498db;
    box-shadow: 0 0 3px rgba(52, 152, 219, 0.5);
  }

  .filter-button {
    padding: 0.5rem 1rem;
    background-color: #3498db;
    color: white;
    border: none;
    border-radius: 4px;
    cursor: pointer;
  }

  .filter-button:hover {
    background-color: #2980b9;
  }

  .active-filters {
    padding: 0 2rem 1rem 2rem;
    display: flex;
    flex-wrap: wrap;
    gap: 0.5rem;
  }

  .filter-tag {
    background-color: #e1e1e1;
    padding: 0.3rem 0.6rem;
    border-radius: 4px;
    display: flex;
    align-items: center;
    gap: 0.5rem;
  }

  .remove-filter {
    border: none;
    background: none;
    color: #666;
    cursor: pointer;
    padding: 0;
    font-size: 1.2rem;
    line-height: 1;
  }

  .remove-filter:hover {
    color: #333;
  }
</style>
