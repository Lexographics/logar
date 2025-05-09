<script>
  import { browser } from "$app/environment";
  import { page } from "$app/stores";
  import BaseView from "$lib/widgets/BaseView.svelte";
  import logService from "$lib/service/logService";
  import Sidebar from "$lib/widgets/Sidebar.svelte";
  import { SessionStorage } from "$lib/storage.svelte";
  import moment from "moment";
  import { mount, onMount, untrack } from "svelte";
  import { fade, fly } from 'svelte/transition';
  import LL from "../../../i18n/i18n-svelte";

  let model = $state("");
  let filters = $state([]);
  
  function ModelChanged() {
    model = $page.url.searchParams.get("model") || "";
    logs = [];
    loading = false;
    lastCursor = 0;
    hasMore = true;
    connectionAttempted = false;
    
    disconnectLogStream();
    setupLogStream();
    loadMore();
  }

  $effect(() => {
    model = $page.url.searchParams.get("model") || "";

    untrack(() => {
      ModelChanged();
    });
  });
  
  const FILTER_FIELDS = $state([
    { value: 'id', label: $LL.logs.fields.id() },
    { value: 'created_at', label: $LL.logs.fields.timestamp() },
    { value: 'category', label: $LL.logs.fields.category() },
    { value: 'message', label: $LL.logs.fields.message() },
    { value: 'severity', label: $LL.logs.fields.severity() },
  ]);

  const FILTER_OPERATORS = $state([
    // Single Value
    { value: '=', label: '=', requires: 1 },
    { value: '!=', label: '!=', requires: 1 },
    { value: '>', label: '>', requires: 1 },
    { value: '>=', label: '>=', requires: 1 },
    { value: '<', label: '<', requires: 1 },
    { value: '<=', label: '<=', requires: 1 },
    { value: 'contains', label: $LL.logs.operators.contains(), requires: 1 },
    { value: 'not_contains', label: $LL.logs.operators.not_contains(), requires: 1 },
    { value: 'starts_with', label: $LL.logs.operators.starts_with(), requires: 1 },
    { value: 'ends_with', label: $LL.logs.operators.ends_with(), requires: 1 },
    // Double Value
    { value: 'between', label: $LL.logs.operators.between(), requires: 2 },
    { value: 'not_between', label: $LL.logs.operators.not_between(), requires: 2 },
    // Multi Value
    { value: 'in', label: $LL.logs.operators.in(), requires: 'multi' },
    { value: 'not_in', label: $LL.logs.operators.not_in(), requires: 'multi' },
  ]);

  let selectedField = $state(FILTER_FIELDS[0].value);
  let selectedOperator = $state(FILTER_OPERATORS[0].value);
  let filterValue1 = $state(''); // For single value or first value of range
  let filterValue2 = $state(''); // For second value of range
  let filterValuesMulti = $state(''); // For multi value input (comma separated)

  let currentOperatorInfo = $derived(FILTER_OPERATORS.find(op => op.value === selectedOperator));

  function getCurrentTimestamp() {
    return moment().format("DD-MM-YYYY HH:mm:ss.SSS");
  }

  function getSeverityClass(severity) {
    switch (severity) {
      case 1:
        return "trace";
      case 2:
        return "log";
      case 3:
        return "info";
      case 4:
        return "warning";
      case 5:
        return "error";
      case 6:
        return "fatal";
      default:
        return "what";
    }
  }

  let logs = $state([]);
  let loading = $state(false);
  let lastCursor = $state(0);
  let hasMore = $state(true);
  let isStreamConnected = $state(false);
  let connectionAttempted = $state(false);
  let streamCleanup = null;

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
    const field = selectedField;
    const operatorInfo = currentOperatorInfo;
    if (!operatorInfo) return;

    let value;
    let valid = true;

    switch (operatorInfo.requires) {
      case 1:
        if (filterValue1.trim() === '') {
          valid = false;
          alert(`Please enter a value for the ${operatorInfo.label} operator.`);
        } else {
          value = [filterValue1.trim()];
        }
        break;
      case 2:
        if (filterValue1.trim() === '' || filterValue2.trim() === '') {
          valid = false;
          alert(`Please enter two values for the ${operatorInfo.label} operator.`);
        } else {
          value = [filterValue1.trim(), filterValue2.trim()];
        }
        break;
      case 'multi':
        if (filterValuesMulti.trim() === '') {
          valid = false;
          alert(`Please enter comma-separated values for the ${operatorInfo.label} operator.`);
        } else {
          value = filterValuesMulti.split(',').map(v => v.trim()).filter(v => v !== '');
          if (value.length === 0) {
             valid = false;
             alert(`Please enter valid comma-separated values for the ${operatorInfo.label} operator.`);
          }
        }
        break;
      default:
        valid = false;
        console.error("Invalid operator requirement:", operatorInfo.requires);
    }

    if (valid) {
      const newFilter = {
        field: field,
        operator: operatorInfo.value,
        value: value,
      };
      filters = [...filters, newFilter];

      filterValue1 = '';
      filterValue2 = '';
      filterValuesMulti = '';

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
    
    const [data, err] = await logService.getLogs(modelName, lastCursor, filters);

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
    }

    loading = false;
  }

  function handleNewLogs(data) {
    if (data.Logs && data.Logs.length > 0) {
      logs = [...data.Logs.map(l => ({ ...l, isNew: true })), ...logs.map(l => ({ ...l, isNew: false }))];
    }
  }

  function setupLogStream() {
    if (streamCleanup) {
      streamCleanup();
    }
    
    streamCleanup = logService.connectToLogStream(
      model, 
      filters, 
      handleNewLogs,
      (error) => {
        console.error("Stream error:", error);
        isStreamConnected = false;
        connectionAttempted = true;
      },
      () => {
        console.log("Stream connection closed.");
        isStreamConnected = false; 
        connectionAttempted = true;
      }
    );
    
    isStreamConnected = true;
    connectionAttempted = true;
  }

  function disconnectLogStream() {
    if (streamCleanup) {
      streamCleanup();
      streamCleanup = null;
    }
    isStreamConnected = false;
    logs = [];
    loading = false;
    lastCursor = 0;
    hasMore = true;
    connectionAttempted = false;
  }

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
  <div class="page-content">
     <div style="padding-bottom: 2.5rem;"></div>

    <div class="filter-creator">
      <select bind:value={selectedField} class="filter-select">
        {#each FILTER_FIELDS as field}
          <option value={field.value}>{field.label}</option>
        {/each}
      </select>

      {#if selectedField === 'created_at'}
        <p class="filter-tip">{$LL.logs.time_format_tip()}</p>
      {:else if selectedField === 'severity'}
        <p class="filter-tip">{$LL.logs.severity_tip()}</p>
      {/if}

      <select bind:value={selectedOperator} class="filter-select">
        {#each FILTER_OPERATORS as operator}
          <option value={operator.value}>{operator.label}</option>
        {/each}
      </select>

      {#if currentOperatorInfo?.requires === 1}
        <input
          type="text"
          bind:value={filterValue1}
          placeholder={$LL.logs.value()}
          class="filter-input"
          onkeydown={(e) => e.key === 'Enter' && handleAddFilter()}
        />
        {#if selectedField === 'created_at'}
          <button class="now-button" onclick={() => filterValue1 = getCurrentTimestamp()}>{$LL.logs.now()}</button>
        {/if}
      {:else if currentOperatorInfo?.requires === 2}
        <input
          type="text"
          bind:value={filterValue1}
          placeholder={$LL.logs.from()}
          class="filter-input filter-input-range"
        />
        {#if selectedField === 'created_at'}
          <button class="now-button" onclick={() => filterValue1 = getCurrentTimestamp()}>{$LL.logs.now()}</button>
        {/if}
        <input
          type="text"
          bind:value={filterValue2}
          placeholder={$LL.logs.to()}
          class="filter-input filter-input-range"
          onkeydown={(e) => e.key === 'Enter' && handleAddFilter()}
        />
        {#if selectedField === 'created_at'}
          <button class="now-button" onclick={() => filterValue2 = getCurrentTimestamp()}>{$LL.logs.now()}</button>
        {/if}
      {:else if currentOperatorInfo?.requires === 'multi'}
        <input
          type="text"
          bind:value={filterValuesMulti}
          placeholder={$LL.logs.values_comma_separated()}
          class="filter-input filter-input-multi"
          onkeydown={(e) => e.key === 'Enter' && handleAddFilter()}
        />
      {/if}

      <button class="filter-button" onclick={handleAddFilter}>
        {$LL.logs.add_filter()}
      </button>
    </div>

    {#if filters.length > 0}
      <div class="active-filters">
        {#each filters as filter, i}
          {@const fieldLabel = FILTER_FIELDS.find(f => f.value === filter.field)?.label || filter.field}
          {@const operatorLabel = FILTER_OPERATORS.find(op => op.value === filter.operator)?.label || filter.operator}
          {@const valueDisplay = Array.isArray(filter.value) ? filter.value.join(', ') : filter.value}
          <div class="filter-tag">
            <span>{fieldLabel} {operatorLabel} {valueDisplay}</span>
            <button class="remove-filter" onclick={() => removeFilter(i)}>×</button>
          </div>
        {/each}
      </div>
    {/if}

    <div class="sse-status-container">
      <span class="status-dot {isStreamConnected ? 'connected' : 'disconnected'}"></span>
      <span class="status-text">{isStreamConnected ? $LL.logs.live_stream_connected() : $LL.logs.live_stream_disconnected()}</span>
      {#if connectionAttempted && !isStreamConnected}
        <button class="reconnect-button" onclick={setupLogStream} title={$LL.logs.reconnect_tip()}>
          {$LL.logs.reconnect()}
        </button>
      {/if}
    </div>

    <div style="flex-grow: 1; min-height: 0; max-width: 100%; overflow-y: scroll;" class="scrollbar">
      <table style="max-width: 100%;">
        <thead>
          <tr>
            <th style="text-align: center;">{$LL.logs.fields.id()}</th>
            <th style="text-align: center;"><i class="fa-solid fa-signal"></i></th>
            <th style="text-align: center;">{$LL.logs.fields.timestamp()}</th>
            <th style="text-align: left;">{$LL.logs.fields.message()}</th>
            <th style="text-align: center;">{$LL.logs.fields.category()}</th>
          </tr>
        </thead>
        <tbody>
          {#key [model, JSON.stringify(filters)]}
            {#each logs as log}
              <tr
                class="row {getSeverityClass(log.Severity)}"
                in:fly={{ x: -50, duration: 300 }}
                out:fade={{ duration: 200 }}
              >
                <td style="width: 1%;">{log.ID}</td>
                <td style="width: 1%;">{getSeverityClass(log.Severity).toUpperCase()}</td>
                <td style="width: 24ch;">{moment(log.CreatedAt).format("DD-MM-YYYY HH:mm:ss.SSS")}</td>
                <td style="width: 60%; word-break: break-all; text-align: left;">{log.Message}</td>
                <td style="width: 1%; white-space: nowrap; overflow: hidden; text-overflow: ellipsis;">{log.Category}</td>
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
  </div>
</BaseView>

<style>
  .filter-creator {
    padding: 0 2rem 1rem 2rem;
    display: flex;
    gap: 0.5rem;
    align-items: center;
    flex-wrap: wrap;
  }

  .filter-select,
  .filter-input {
    padding: 0.5rem;
    font-size: 1rem;
    border: 1px solid var(--input-border);
    border-radius: 4px;
    box-sizing: border-box;
    background-color: var(--input-background);
    color: var(--input-text);
  }

  .filter-select {
    min-width: 100px;
  }

  option {
    background-color: var(--input-background-opaque);
  }

  .filter-input {
    flex-grow: 1;
  }

  .filter-input-range {
     flex-grow: 0.5;
  }

  .filter-input-multi {
    min-width: 200px;
  }

  .filter-select:focus,
  .filter-input:focus {
    outline: none;
    border-color: var(--primary-color);
    box-shadow: 0 0 3px var(--shadow-color);
  }

  table {
    width: 100%;
    border-collapse: collapse;
  }

  th,
  td {
    padding: 0.5rem;
    text-align: center;
    vertical-align: text-top;
    border-bottom: 1px solid var(--border-color);
    color: var(--text-color);
  }

  th {
    background-color: var(--card-background);
    color: var(--text-color);
    position: sticky;
    top: 0;
    z-index: 1;
  }

  .row.trace {
    --row-bg: var(--row-trace-bg);
    --row-striped-bg: var(--row-trace-striped-bg);
  }

  .row.log {
    --row-bg: var(--row-log-bg);
    --row-striped-bg: var(--row-log-striped-bg);
  }

  .row.info {
    --row-bg: var(--row-info-bg);
    --row-striped-bg: var(--row-info-striped-bg);
  }

  .row.warning {
    --row-bg: var(--row-warning-bg);
    --row-striped-bg: var(--row-warning-striped-bg);
  }

  .row.error {
    --row-bg: var(--row-error-bg);
    --row-striped-bg: var(--row-error-striped-bg);
  }

  .row.fatal {
    --row-bg: var(--row-fatal-bg);
    --row-striped-bg: var(--row-fatal-striped-bg);
    font-weight: bold;
  }

  .row {
    background-color: var(--row-bg);
  }

  .row:nth-child(odd) {
    background-color: var(--row-striped-bg);
  }

  .loader {
    padding: 2rem;
    text-align: center;
  }

  .loading-spinner {
    width: 40px;
    height: 40px;
    margin: 0 auto;
    border: 3px solid var(--input-background);
    border-top: 3px solid var(--primary-color);
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

  .filter-button {
    padding: 0.5rem 1rem;
    background-color: var(--primary-color);
    color: white;
    border: none;
    border-radius: 4px;
    cursor: pointer;
    white-space: nowrap; /* Prevent button text wrapping */
  }

  .filter-button:hover {
    background-color: var(--primary-hover-color);
  }

  .active-filters {
    padding: 0 2rem 1rem 2rem;
    display: flex;
    flex-wrap: wrap;
    gap: 0.5rem;
  }

  .filter-tag {
    background-color: var(--input-background);
    color: var(--text-color);
    padding: 0.3rem 0.6rem;
    border-radius: 4px;
    display: flex;
    align-items: center;
    gap: 0.5rem;
    border: 1px solid var(--border-color);
  }

  .remove-filter {
    border: none;
    background: none;
    color: var(--text-secondary-color);
    cursor: pointer;
    padding: 0;
    font-size: 1.2rem;
    line-height: 1;
  }

  .remove-filter:hover {
    color: var(--text-color);
  }

  .filter-tip {
    font-size: 0.8rem;
    color: var(--text-secondary-color);
    margin: 0.2rem 0 0 0;
    width: 100%;
    order: 3;
    text-align: left;
    padding-left: 0.5rem;
  }
  
  /* Now button */
  .now-button {
    padding: 0.4rem 0.6rem;
    font-size: 0.9rem;
    background-color: var(--input-background);
    color: var(--text-color);
    border: 1px solid var(--input-border);
    border-radius: 4px;
    cursor: pointer;
    margin-left: -5px;
    z-index: 1;
    white-space: nowrap;
  }

  .now-button:hover {
    background-color: var(--primary-color);
    color: white;
  }

  .sse-status-container {
    display: flex;
    align-items: center;
    gap: 0.75rem;
    padding: 0.5rem 2rem 1rem 2rem;
    border-bottom: 1px solid var(--border-color);
    margin-bottom: 1rem;
  }

  .status-text {
    font-size: 0.9em;
    color: var(--text-secondary-color);
  }

  .reconnect-button {
    padding: 0.3rem 0.8rem;
    font-size: 0.85em;
    background-color: var(--primary-color);
    color: white;
    border: none;
    border-radius: 4px;
    cursor: pointer;
  }

  .reconnect-button:hover {
    background-color: var(--primary-hover-color);
  }

  .status-dot {
    width: 12px;
    height: 12px;
    border-radius: 50%;
    display: inline-block;
    flex-shrink: 0;
  }

  .status-dot.connected {
    background-color: #2ecc71;
    animation: blink 1.5s infinite ease-in-out;
    box-shadow: 0 0 5px rgba(46, 204, 113, 0.7);
  }

  .status-dot.disconnected {
    background-color: #e74c3c;
    box-shadow: 0 0 5px rgba(231, 76, 60, 0.7);
  }

  @keyframes blink {
    0%, 100% { opacity: 1; }
    50% { opacity: 0.3; }
  }

  .page-content {
    display: flex;
    flex-direction: column;
    height: 100%;
  }
</style>
