<script lang="ts">
  import { globalsService, type Global } from '$lib/service/globalsService';
  import BaseView from '$lib/widgets/BaseView.svelte';
  import { onMount } from 'svelte';

  let globals = $state<Global[]>([]);
  let loading = $state(true);
  let error = $state<string | null>(null);
  let newKey = $state('');
  let newValue = $state('');
  let editingKey = $state<string | null>(null);
  let editingValue = $state('');
  let editingExported = $state<boolean>(false);

  async function loadGlobals() {
    loading = true;
    error = null;
    const [data, err] = await globalsService.getGlobals();
    if (err) {
      error = err.message;
    } else {
      globals = data;
    }
    loading = false;
  }

  async function addGlobal() {
    if (!newKey || !newValue) return;
    let jsonData: any;
    try {
      jsonData = JSON.parse(newValue);
    } catch (err) {
      error = 'Invalid JSON';
      return;
    }
    
    error = null;
    const [data, err] = await globalsService.updateGlobal(newKey, jsonData, false);
    if (err) {
      error = err.message;
    } else {
      globals.push({
        Key: newKey,
        Value: newValue,
        Exported: false,
      });
      newKey = '';
      newValue = '';
    }
  }

  async function updateGlobal(key: string) {
    if (!editingValue) return;
    let jsonData: any;
    try {
      jsonData = JSON.parse(editingValue);
    } catch (err) {
      error = 'Invalid JSON';
      return;
    }
    error = null;
    const [data, err] = await globalsService.updateGlobal(key, jsonData, editingExported);
    if (err) {
      error = err.message;
    } else {
      await loadGlobals();
      editingKey = null;
      editingValue = '';
      editingExported = false;
    }
  }

  async function deleteGlobal(key: string) {
    error = null;
    const err = await globalsService.deleteGlobal(key);
    if (err) {
      error = err.message;
    } else {
      await loadGlobals();
    }
  }

  onMount(() => {
    loadGlobals();
  });
</script>

<BaseView>
  <div class="page">
    <h1>Global Variables</h1>

    {#if error}
      <div class="error-message">
        {error}
      </div>
    {/if}

    <div class="add-form">
      <input
        type="text"
        bind:value={newKey}
        placeholder="Key"
        class="input"
      />
      <input
        type="text"
        bind:value={newValue}
        placeholder="Value"
        class="input"
      />
      <button onclick={addGlobal} class="add-button">
        Add Global
      </button>
    </div>

    {#if loading}
      <div class="loading">Loading...</div>
    {:else}
      <div class="globals">
        {#each globals as global (global.Key)}
          <div class="item">
            {#if editingKey === global.Key}
              <input
                type="text"
                bind:value={editingValue}
                class="input"
              />
              <label class="export-label">
                <input
                  type="checkbox"
                  bind:checked={editingExported}
                />
                Exported
              </label>
              <button onclick={() => updateGlobal(global.Key)} class="save-button">
                Save
              </button>
              <button onclick={() => { editingKey = null; editingValue = ''; editingExported = false; }} class="cancel-button">
                Cancel
              </button>
            {:else}
              <div class="content">
                <span class="key">{global.Key}</span>
                <span class="value">{global.Value}</span>
                <div class="export-status">
                  <span class="export-label">
                    <input type="checkbox" checked={global.Exported} disabled />
                    {global.Exported ? 'Exported' : 'Not Exported'}
                  </span>
                </div>
              </div>
              <div class="action-buttons">
                <button onclick={() => { editingKey = global.Key; editingValue = global.Value; editingExported = global.Exported; }} class="edit-button">
                  Edit
                </button>
                <button onclick={() => deleteGlobal(global.Key)} class="delete-button">
                  Delete
                </button>
              </div>
            {/if}
          </div>
        {/each}
      </div>
    {/if}
  </div>
</BaseView>

<style>
  .page {
    max-width: 800px;
    margin: 2rem auto;
    padding: 0 1rem;
  }

  h1 {
    color: var(--text-color);
    margin-bottom: 2rem;
    font-size: 2rem;
  }

  .error-message {
    background-color: var(--error-color);
    color: white;
    padding: 1rem;
    border-radius: 4px;
    margin-bottom: 1rem;
  }

  .add-form {
    display: flex;
    gap: 1rem;
    margin-bottom: 2rem;
  }

  .input {
    flex: 1;
    padding: 0.5rem;
    border: 1px solid var(--border-color);
    border-radius: 4px;
    background-color: var(--input-background);
    color: var(--input-text);
  }

  .add-button {
    background-color: var(--primary-color);
    color: white;
    border: none;
    padding: 0.5rem 1rem;
    border-radius: 4px;
    cursor: pointer;
    transition: background-color 0.2s;
  }

  .add-button:hover {
    background-color: var(--primary-hover-color);
  }

  .loading {
    text-align: center;
    color: var(--text-secondary-color);
    padding: 2rem;
  }

  .globals {
    display: flex;
    flex-direction: column;
    gap: 1rem;
  }

  .item {
    background-color: var(--card-background);
    border: 1px solid var(--border-color);
    border-radius: 4px;
    padding: 1rem;
    display: flex;
    justify-content: space-between;
    align-items: center;
    gap: 1rem;
  }

  .content {
    flex: 1;
    display: flex;
    flex-direction: column;
    gap: 0.5rem;
  }

  .key {
    font-weight: bold;
    color: var(--text-color);
  }

  .value {
    color: var(--text-secondary-color);
    background-color: var(--input-background);
    padding: 0.4rem 0.6rem;
    border-radius: 4px;
    border: 1px solid var(--border-color);
    font-family: monospace;
    white-space: pre-wrap;
    word-break: break-all;
  }

  .action-buttons {
    display: flex;
    gap: 0.5rem;
  }

  .edit-button,
  .save-button {
    background-color: var(--primary-color);
    color: white;
    border: none;
    padding: 0.5rem 1rem;
    border-radius: 4px;
    cursor: pointer;
    transition: background-color 0.2s;
  }

  .edit-button:hover,
  .save-button:hover {
    background-color: var(--primary-hover-color);
  }

  .delete-button {
    background-color: var(--error-color);
    color: white;
    border: none;
    padding: 0.5rem 1rem;
    border-radius: 4px;
    cursor: pointer;
    transition: background-color 0.2s;
  }

  .delete-button:hover {
    background-color: #c0392b;
  }

  .cancel-button {
    background-color: var(--text-secondary-color);
    color: white;
    border: none;
    padding: 0.5rem 1rem;
    border-radius: 4px;
    cursor: pointer;
    transition: background-color 0.2s;
  }

  .cancel-button:hover {
    background-color: #7f8c8d;
  }

  .export-status {
    margin-top: 0.5rem;
  }

  .export-label {
    display: flex;
    align-items: center;
    gap: 0.5rem;
    color: var(--text-secondary-color);
    font-size: 0.9rem;
  }

  .export-label input[type="checkbox"] {
    width: 1rem;
    height: 1rem;
  }

  @media (max-width: 768px) {
    .page {
      margin: 1rem auto;
    }

    h1 {
      font-size: 1.5rem;
      margin-bottom: 1.5rem;
    }

    .add-form {
      flex-direction: column;
      gap: 0.5rem;
    }

    .input {
      width: 100%;
    }

    .add-button {
      width: 100%;
    }

    .item {
      flex-direction: column;
      align-items: stretch;
    }

    .content {
      margin-bottom: 1rem;
    }

    .value {
      width: 100%;
      box-sizing: border-box;
    }

    .action-buttons {
      flex-direction: column;
      width: 100%;
    }

    .edit-button,
    .save-button,
    .delete-button,
    .cancel-button {
      width: 100%;
    }
  }
</style>
