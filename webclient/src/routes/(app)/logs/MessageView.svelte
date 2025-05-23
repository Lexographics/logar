<script lang="ts">
  let { message }: { message: string } = $props();
  
  let isJson = $state(false);
  let parsedJson: any = $state(null);
  let isExpanded = $state(false);

  $effect(() => {
    try {
      parsedJson = JSON.parse(message);
      isJson = true;
    } catch {
      isJson = false;
      parsedJson = null;
    }
  });

  function toggleExpand() {
    isExpanded = !isExpanded;
  }

  function renderValue(value: any): string {
    if (value === null) return 'null';
    if (value === undefined) return 'undefined';
    if (typeof value === 'string') return `"${value}"`;
    return String(value);
  }

  function renderJsonNode(node: any, depth: number = 0): string {
    if (node === null || node === undefined) return renderValue(node);
    
    if (Array.isArray(node)) {
      if (node.length === 0) return '[]';
      return `[
${node.map(item => ' '.repeat((depth + 1) * 2) + renderJsonNode(item, depth + 1)).join(',\n')}
${' '.repeat(depth * 2)}]`;
    }
    
    if (typeof node === 'object') {
      const entries = Object.entries(node);
      if (entries.length === 0) return '{}';
      return `{
${entries.map(([key, value]) => ' '.repeat((depth + 1) * 2) + `"${key}": ${renderJsonNode(value, depth + 1)}`).join(',\n')}
${' '.repeat(depth * 2)}}`;
    }
    
    return renderValue(node);
  }
</script>

<div class="message-view">
  {#if isJson}
    <div class="json-view" class:expanded={isExpanded}>
      <button class="expand-button" onclick={toggleExpand}>
        <span class="icon">{isExpanded ? '▼' : '▶'}</span>
        <span class="preview">{isExpanded ? 'Collapse' : 'Expand'} JSON</span>
      </button>
      {#if isExpanded}
        <pre class="json-content">{renderJsonNode(parsedJson)}</pre>
      {:else}
        <pre class="json-preview">{message.slice(0, 100)}{message.length > 100 ? '...' : ''}</pre>
      {/if}
    </div>
  {:else}
    {message}
  {/if}
</div>

<style>
  .message-view {
    font-size: 0.9rem;
    line-height: 1.4;
  }

  .json-view {
    background-color: var(--card-background);
    border-radius: 4px;
    overflow: hidden;
  }

  .expand-button {
    display: flex;
    align-items: center;
    gap: 8px;
    width: 100%;
    padding: 4px 8px;
    border: none;
    background: none;
    cursor: pointer;
    color: var(--text-color);
    font-size: 0.9rem;
    text-align: left;
  }

  .expand-button:hover {
    background-color: var(--input-background);
  }

  .icon {
    font-size: 0.8rem;
    color: var(--text-secondary-color);
  }

  .preview {
    color: var(--text-secondary-color);
  }

  .json-content {
    margin: 0;
    padding: 8px;
    white-space: pre;
    overflow-x: auto;
    border-top: 1px solid var(--border-color);
  }

  .json-preview {
    margin: 0;
    padding: 4px 8px;
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
    color: var(--text-secondary-color);
    font-style: italic;
  }
</style>
