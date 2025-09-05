<script>
  import { browser } from "$app/environment";
  import { page } from "$app/stores";
  import BaseView from "$lib/widgets/BaseView.svelte";
  import logService from "$lib/service/logService";
  import Sidebar from "$lib/widgets/Sidebar.svelte";
  import { SessionStorage } from "$lib/storage.svelte";
  import moment from "moment";
  import { mount, onMount, untrack } from "svelte";
  import { fade, fly, slide } from 'svelte/transition';
  import LL from "../../../i18n/i18n-svelte";
  import MessageView from "./MessageView.svelte";

  let model = $state("");
  let filters = $state([]);
  let showFilters = $state(false);
  
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

  onMount(() => {
    return () => {
      disconnectLogStream();
    };
  });

</script>

<BaseView>
  <div class="page-content">
     <div style="padding-bottom: 0rem;"></div>

    <div class="filters-section">
      <button class="toggle-filters" onclick={() => showFilters = !showFilters}>
        <i class="fa-solid {showFilters ? 'fa-chevron-up' : 'fa-chevron-down'}"></i>
        {$LL.logs.filters()}
      </button>

      {#if showFilters}
        <div class="filters-panel" transition:slide={{ duration: 200 }}>
          <div class="filter-creator">
            <div class="filter-selects-row">
              <select bind:value={selectedField} class="filter-select">
                {#each FILTER_FIELDS as field}
                  <option value={field.value}>{field.label}</option>
                {/each}
              </select>

              <select bind:value={selectedOperator} class="filter-select">
                {#each FILTER_OPERATORS as operator}
                  <option value={operator.value}>{operator.label}</option>
                {/each}
              </select>
              
              {#if currentOperatorInfo?.requires === 1}
              <div class="input-with-button">
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
              </div>
            {:else if currentOperatorInfo?.requires === 2}
              <div class="input-with-button">
                <input
                  type="text"
                  bind:value={filterValue1}
                  placeholder={$LL.logs.from()}
                  class="filter-input filter-input-range"
                />
                {#if selectedField === 'created_at'}
                <button class="now-button" onclick={() => filterValue1 = getCurrentTimestamp()}>{$LL.logs.now()}</button>
                {/if}
              </div>
              <div class="input-with-button">
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
              </div>
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
              {#if selectedField === 'created_at'}
                <p class="filter-tip">{$LL.logs.time_format_tip()}</p>
              {:else if selectedField === 'severity'}
                <p class="filter-tip">{$LL.logs.severity_tip()}</p>
              {/if}
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
        </div>
      {/if}
    </div>

    <div style="flex-grow: 1; min-height: 0; max-width: 100%; overflow-y: scroll;" class="scrollbar">
      <div class="logs-container">
        <div class="logs-header">
          <div class="log-cell id-cell">{$LL.logs.fields.id()}</div>
          <div class="log-cell severity-cell"><i class="fa-solid fa-signal"></i></div>
          <div class="log-cell timestamp-cell">{$LL.logs.fields.timestamp()}</div>
          <div class="log-cell message-cell">{$LL.logs.fields.message()}</div>
          <div class="log-cell category-cell">{$LL.logs.fields.category()}</div>
        </div>

        {#key [model, JSON.stringify(filters)]}
          {#each logs as log}
            <div
              class="log-row {getSeverityClass(log.Severity)}"
              in:fly={{ x: -50, duration: 300 }}
              out:fade={{ duration: 200 }}
            >
              <div class="log-cell id-cell">{log.ID}</div>
              <div class="log-cell severity-cell">{getSeverityClass(log.Severity).toUpperCase()}</div>
              <div class="log-cell timestamp-cell">{moment(log.CreatedAt).format("DD-MM-YYYY HH:mm:ss.SSS")}</div>
              <div class="log-cell message-cell">
                <MessageView message={log.Message} />
              </div>
              <div class="log-cell category-cell">{log.Category}</div>
            </div>
          {/each}
        {/key}
      </div>

      {#if hasMore}
        <div class="load-more-container">
          <button class="load-more-button" onclick={loadMore} disabled={loading}>
            {#if loading}
              <div class="loading-spinner-small"></div>
              {$LL.logs.loading()}
            {:else}
              {$LL.logs.load_more()}
            {/if}
          </button>
        </div>
      {:else}
        <div class="no-more-logs">
          <p>{$LL.logs.no_more_logs()}</p>
        </div>
      {/if}
    </div>
  </div>
</BaseView>

<style>
  .filter-creator {
    padding: 0 2rem 0.5rem 2rem;
    display: flex;
    gap: 0.5rem;
    align-items: center;
    flex-wrap: wrap;
  }

  .filter-select,
  .filter-input {
    padding: 0.35rem 0.5rem;
    font-size: 0.9rem;
    border: 1px solid var(--input-border);
    border-radius: 4px;
    box-sizing: border-box;
    background-color: var(--input-background);
    color: var(--input-text);
    height: 32px;
  }

  .filter-select {
    min-width: 100px;
  }

  .filter-input {
    padding: 0.35rem 0.5rem;
    font-size: 0.9rem;
    border: 1px solid var(--input-border);
    border-radius: 4px;
    box-sizing: border-box;
    background-color: var(--input-background);
    color: var(--input-text);
    height: 32px;
    flex-grow: 1;
  }

  .filter-input-range {
    flex-grow: 0.5;
  }

  .filter-input-multi {
    min-width: 200px;
  }

  .input-with-button {
    display: flex;
    align-items: center;
    gap: 0.25rem;
    flex-grow: 1;
  }

  .now-button {
    padding: 0 0.4rem;
    font-size: 0.8rem;
    background-color: var(--input-background);
    color: var(--text-color);
    border: 1px solid var(--input-border);
    border-radius: 4px;
    cursor: pointer;
    white-space: nowrap;
    height: 32px;
    line-height: 1;
    flex-shrink: 0;
    box-sizing: border-box;
    display: flex;
    align-items: center;
  }

  .filter-tip {
    font-size: 0.75rem;
    color: var(--text-secondary-color);
    margin: 0.1rem 0 0 0;
    width: 100%;
    order: 3;
    text-align: left;
    padding-left: 0.5rem;
  }

  .filter-selects-row {
    display: flex;
    flex-direction: column;
    gap: 0.5rem;
    width: 100%;
  }

  .filter-button {
    padding: 0.35rem 0.8rem;
    background-color: var(--primary-color);
    color: white;
    border: none;
    border-radius: 4px;
    cursor: pointer;
    white-space: nowrap;
    height: 32px;
  }

  .filter-button:hover {
    background-color: var(--primary-hover-color);
  }

  .active-filters {
    padding: 0 2rem 0.5rem 2rem;
    display: flex;
    flex-wrap: wrap;
    gap: 0.4rem;
  }

  .filter-tag {
    background-color: var(--input-background);
    color: var(--text-color);
    padding: 0.2rem 0.4rem;
    border-radius: 4px;
    display: flex;
    align-items: center;
    gap: 0.4rem;
    border: 1px solid var(--border-color);
    font-size: 0.85rem;
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

  .toggle-filters {
    display: flex;
    align-items: center;
    gap: 0.5rem;
    padding: 0.5rem 2rem;
    background: none;
    border: none;
    color: var(--text-secondary-color);
    cursor: pointer;
    font-size: 0.9rem;
    width: 100%;
    text-align: left;
  }

  .toggle-filters:hover {
    color: var(--text-color);
  }

  .toggle-filters i {
    font-size: 0.8rem;
    transition: transform 0.2s ease;
  }

  .filters-panel {
    border-bottom: 1px solid var(--border-color);
  }

  .sse-status-container {
    display: flex;
    align-items: center;
    gap: 0.75rem;
    padding: 0.5rem 2rem;
    border-top: 1px solid var(--border-color);
    margin-top: 0.5rem;
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

  .logs-container {
    width: 100%;
    display: flex;
    flex-direction: column;
  }

  .logs-header {
    display: flex;
    background-color: var(--card-background);
    color: var(--text-color);
    position: sticky;
    top: 0;
    z-index: 1;
    border-bottom: 1px solid var(--border-color);
  }

  .log-row {
    display: flex;
    border-bottom: 1px solid var(--border-color);
  }

  .log-cell {
    padding: 0.5rem;
    overflow: hidden;
    text-overflow: ellipsis;
  }

  .id-cell {
    width: 5%;
    min-width: 60px;
    text-align: center;
  }

  .severity-cell {
    width: 5%;
    min-width: 60px;
    text-align: center;
  }

  .timestamp-cell {
    width: 20%;
    min-width: 200px;
    white-space: nowrap;
    text-align: center;
  }

  .message-cell {
    flex: 1;
    min-width: 200px;
    word-break: break-all;
    text-align: left;
  }

  .category-cell {
    width: 10%;
    min-width: 100px;
    white-space: nowrap;
    text-overflow: ellipsis;
    text-align: center;
  }

  .log-row.trace {
    --row-bg: var(--row-trace-bg);
    --row-striped-bg: var(--row-trace-striped-bg);
  }

  .log-row.log {
    --row-bg: var(--row-log-bg);
    --row-striped-bg: var(--row-log-striped-bg);
  }

  .log-row.info {
    --row-bg: var(--row-info-bg);
    --row-striped-bg: var(--row-info-striped-bg);
  }

  .log-row.warning {
    --row-bg: var(--row-warning-bg);
    --row-striped-bg: var(--row-warning-striped-bg);
  }

  .log-row.error {
    --row-bg: var(--row-error-bg);
    --row-striped-bg: var(--row-error-striped-bg);
  }

  .log-row.fatal {
    --row-bg: var(--row-fatal-bg);
    --row-striped-bg: var(--row-fatal-striped-bg);
    font-weight: bold;
  }

  .log-row {
    background-color: var(--row-bg);
  }

  .log-row:nth-child(odd) {
    background-color: var(--row-striped-bg);
  }

  @media (max-width: 768px) {
    .filter-creator {
      padding: 0 1rem 0.5rem 1rem;
      flex-direction: column;
      align-items: stretch;
      gap: 0.4rem;
    }

    .filter-select,
    .filter-input,
    .filter-input-range,
    .filter-input-multi {
      width: 100%;
      min-width: unset;
    }

    .filter-select {
      width: auto;
      flex: 1;
    }

    .filter-selects-row {
      display: flex;
      gap: 0.4rem;
      width: 100%;
    }

    .input-with-button {
      width: 100%;
    }

    .filter-button {
      width: 100%;
    }

    .active-filters {
      padding: 0 1rem 0.5rem 1rem;
    }

    .sse-status-container {
      padding: 0.5rem 1rem;
      flex-wrap: wrap;
    }

    .reconnect-button {
      width: 100%;
      margin-top: 0.3rem;
    }

    .logs-container {
      min-width: 100%;
    }

    .logs-header {
      display: none;
    }

    .log-row {
      flex-direction: column;
      padding: 0.5rem;
      gap: 0.5rem;
    }

    .log-cell {
      width: 100%;
      min-width: 100%;
      padding: 0.25rem 0;
      text-align: left;
    }

    .log-cell::before {
      content: attr(data-label);
      font-weight: bold;
      margin-right: 0.5rem;
      color: var(--text-secondary-color);
    }

    .id-cell::before { content: "ID: "; }
    .severity-cell::before { content: "Severity: "; }
    .timestamp-cell::before { content: "Time: "; }
    .message-cell::before { content: "Message: "; }
    .category-cell::before { content: "Category: "; }

    .log-row {
      border: 1px solid var(--border-color);
      border-radius: 4px;
      margin-bottom: 0.5rem;
    }

    .toggle-filters {
      padding: 0.5rem 1rem;
    }

    .filter-input {
      flex: 1;
    }

    .filter-input-range {
      flex: 1;
    }

    .filter-input-multi {
      width: 100%;
    }
  }

  .load-more-container {
    padding: 2rem;
    text-align: center;
    display: flex;
    justify-content: center;
  }

  .load-more-button {
    padding: 0.75rem 2rem;
    background-color: var(--primary-color);
    color: white;
    border: none;
    border-radius: 6px;
    cursor: pointer;
    font-size: 0.95rem;
    font-weight: 500;
    transition: background-color 0.2s ease;
    display: flex;
    align-items: center;
    gap: 0.5rem;
    min-width: 120px;
    justify-content: center;
  }

  .load-more-button:hover:not(:disabled) {
    background-color: var(--primary-hover-color);
  }

  .load-more-button:disabled {
    opacity: 0.6;
    cursor: not-allowed;
  }

  .no-more-logs {
    padding: 2rem;
    text-align: center;
    color: var(--text-secondary-color);
    font-style: italic;
  }

  .loading-spinner-small {
    width: 16px;
    height: 16px;
    border: 2px solid rgba(255, 255, 255, 0.3);
    border-top: 2px solid white;
    border-radius: 50%;
    animation: spin 1s linear infinite;
  }

  .loader {
    padding: 1rem;
    text-align: center;
    min-height: 60px;
    display: flex;
    align-items: center;
    justify-content: center;
  }

  .loading-spinner {
    width: 30px;
    height: 30px;
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
</style>
