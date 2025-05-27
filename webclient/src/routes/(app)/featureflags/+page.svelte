<script lang="ts">
  import BaseView from "$lib/widgets/BaseView.svelte";
  import Modal from "$lib/widgets/Modal.svelte";
  import featureFlagService from "$lib/service/featureFlagsService";
  import type { FeatureFlag } from "$lib/service/featureFlagsService";
  import { onMount } from "svelte";

  let featureFlags = $state<FeatureFlag[]>([]);
  let flagModal = $state<any>(null);
  let currentFlag = $state<FeatureFlag | null>(null);
  let isNewFlag = $state(false);
  let isLoading = $state(false);
  let error = $state<string | null>(null);

  onMount(async () => {
    await loadFeatureFlags();
  });

  async function loadFeatureFlags() {
    isLoading = true;
    error = null;

    const [featureFlagsData, err] = await featureFlagService.getFeatureFlags();
    if (err) {
      error = err.message;
    } else {
      featureFlags = featureFlagsData;
    }
    isLoading = false;
  }

  function addNewFlag() {
    currentFlag = {
      id: 0,
      name: "",
      enabled: true,
      condition: ""
    };
    isNewFlag = true;
    flagModal?.openModal();
  }

  function startEdit(flag: FeatureFlag) {
    currentFlag = { ...flag };
    isNewFlag = false;
    flagModal?.openModal();
  }

  function onModalClose() {
    error = null;
    currentFlag = null;
    isNewFlag = false;
  }

  async function saveFlag() {
    if (!currentFlag) return;

    if (isNewFlag) {
      const [newFlag, err] = await featureFlagService.createFeatureFlag({
        id: 0,
        name: currentFlag.name,
        enabled: currentFlag.enabled,
        condition: currentFlag.condition
      });
      if (err) {
        error = err.message;
      } else {
        featureFlags = [...featureFlags, newFlag];
        flagModal?.closeModal();
      }
    } else {
      const [updatedFlag, err] = await featureFlagService.updateFeatureFlag({
        id: currentFlag.id,
        name: currentFlag.name,
        enabled: currentFlag.enabled,
        condition: currentFlag.condition
      });
      if (err) {
        error = err.message;
      } else {
        featureFlags = featureFlags.map(f => 
          f.id === updatedFlag.id ? updatedFlag : f
        );
        flagModal?.closeModal();
      }
    }
  }

  async function deleteFlag(flagIdToDelete: number) {
    const flagToDelete = featureFlags.find(f => f.id === flagIdToDelete);
    if (flagToDelete && confirm(`Are you sure you want to delete the flag "${flagToDelete.name}"?`)) {
      const err = await featureFlagService.deleteFeatureFlag(flagIdToDelete);
      if (err) {
        error = err.message;
        } else {
        featureFlags = featureFlags.filter(flag => flag.id !== flagIdToDelete);
      }
    }
  }
</script>

<BaseView>
  <div class="page">
    <div class="header-actions">
      <h2 class="title">Feature Flags</h2>
      <div class="buttons-container">
        <button class="action-button primary" onclick={addNewFlag}>Create New Flag</button>
      </div>
    </div>

    {#if isLoading}
      <div class="loading">Loading feature flags...</div>
    {:else}
      <div class="flags-container">
        {#each featureFlags as flag (flag.id)}
          <div class="flag-card">
            <div class="flag-header">
              <h3 class="flag-name">{flag.name}</h3>
              <label class="switch">
                <input type="checkbox" bind:checked={flag.enabled} disabled={true} />
                <span class="slider round"></span>
              </label>
            </div>
            <div class="condition-preview">
              <label for="flag-condition">Condition</label>
              <pre class="condition-code">{flag.condition || 'No condition set'}</pre>
            </div>
            <div class="flag-actions">
              <button class="action-button" onclick={() => startEdit(flag)}>Edit</button>
              <button class="action-button danger" onclick={() => deleteFlag(flag.id)}>Delete</button>
            </div>
          </div>
        {/each}
      </div>
    {/if}
  </div>
</BaseView>

<Modal 
  bind:this={flagModal} 
  title={isNewFlag ? "Create New Flag" : "Edit Flag"} 
  onClose={onModalClose}
>
  {#if currentFlag}
    <div class="modal-form">
      <div class="form-group">
        <label for="flag-name">Name</label>
        <input 
          type="text" 
          id="flag-name"
          bind:value={currentFlag.name} 
          placeholder="Enter flag name"
        />
      </div>
      
      <div class="form-group">
        <label class="switch-label">
          <span>Enabled</span>
          <label class="switch">
            <input type="checkbox" bind:checked={currentFlag.enabled} />
            <span class="slider round"></span>
          </label>
        </label>
      </div>

      <div class="form-group">
        <label for="flag-condition">Condition</label>
        <textarea 
          id="flag-condition"
          bind:value={currentFlag.condition} 
          rows="6"
          placeholder="Enter condition expression"
        ></textarea>
      </div>


    {#if error}
      <p class="error-message">{error}</p>
    {/if}

      <div class="modal-actions">
        <button class="action-button" onclick={() => flagModal?.closeModal()}>Cancel</button>
        <button class="action-button success" onclick={saveFlag}>Save</button>
      </div>
    </div>

  {/if}
</Modal>

<style>
  .page {
    padding: 24px;
    display: flex;
    flex-direction: column;
    color: var(--text-color);
    gap: 16px;
  }

  .header-actions {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 16px;
  }

  .title {
    font-size: 1.5rem;
    font-weight: 600;
    color: var(--text-color);
    margin-bottom: 0;
  }

  .buttons-container {
    display: flex;
    gap: 12px;
  }

  .action-button {
    padding: 8px 12px;
    border: 1px solid var(--border-color);
    background-color: var(--primary-color);
    color: #fff;
    border-radius: 4px;
    cursor: pointer;
    font-size: 0.9rem;
    transition: background-color 0.2s ease;
  }

  .action-button:hover {
    background-color: var(--primary-hover-color);
  }

  .action-button.primary {
    background-color: var(--primary-color);
    color: #fff;
    border-color: var(--primary-color);
  }

  .action-button.primary:hover {
    background-color: var(--primary-hover-color);
  }

  .action-button.success {
    background-color: var(--success-color);
    color: #fff;
    border-color: var(--success-color);
  }

  .action-button.success:hover {
    background-color: #4CAF50;
  }

  .action-button.danger {
    background-color: var(--error-color);
    color: #fff;
    border-color: var(--error-color);
  }

  .action-button.danger:hover {
    background-color: #D32F2F;
  }

  .flags-container {
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(350px, 1fr));
    gap: 20px;
  }

  .flag-card {
    background-color: var(--card-background);
    border: 1px solid var(--border-color);
    border-radius: 8px;
    padding: 16px;
    display: flex;
    flex-direction: column;
    gap: 12px;
  }

  .flag-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
  }

  .flag-name {
    font-size: 1.1rem;
    font-weight: 500;
    color: var(--text-color);
    margin: 0;
  }

  .condition-preview {
    display: flex;
    flex-direction: column;
    gap: 8px;
  }

  .condition-preview label {
    font-size: 0.9rem;
    color: var(--text-color);
  }

  .condition-code {
    background-color: var(--input-background);
    padding: 8px;
    border-radius: 4px;
    font-family: monospace;
    font-size: 0.9rem;
    margin: 0;
    white-space: pre-wrap;
    word-break: break-all;
    max-height: 100px;
    overflow-y: auto;
  }

  .flag-actions {
    display: flex;
    gap: 8px;
    justify-content: flex-end;
    margin-top: auto;
  }

  /* Switch styles */
  .switch {
    position: relative;
    display: inline-block;
    width: 40px;
    height: 20px;
    flex-shrink: 0;
  }

  .switch input {
    opacity: 0;
    width: 0;
    height: 0;
  }

  .slider {
    position: absolute;
    cursor: pointer;
    top: 0;
    left: 0;
    right: 0;
    bottom: 0;
    background-color: var(--input-background);
    transition: .4s;
  }

  .slider:before {
    position: absolute;
    content: "";
    height: 14px;
    width: 14px;
    left: 3px;
    bottom: 3px;
    background-color: var(--text-color);
    transition: .4s;
  }

  input:checked + .slider {
    background-color: var(--primary-hover-color);
  }

  input:focus + .slider {
    box-shadow: 0 0 1px var(--primary-hover-color);
  }

  input:checked + .slider:before {
    transform: translateX(20px);
  }

  .slider.round {
    border-radius: 20px;
  }

  .slider.round:before {
    border-radius: 50%;
  }

  input:disabled:checked + .slider {
    background-color: var(--primary-color);
    opacity: 0.6;
  }

  input:disabled + .slider:before {
    background-color: var(--text-color);
  }

  input:disabled:checked + .slider:before {
    background-color: var(--input-background);
  }

  .modal-form {
    display: flex;
    flex-direction: column;
    gap: 16px;
    padding: 16px 0;
  }

  .form-group {
    display: flex;
    flex-direction: column;
    gap: 8px;
  }

  .form-group label {
    font-weight: 500;
    color: var(--text-color);
  }

  .form-group input[type="text"],
  .form-group textarea {
    padding: 8px 12px;
    border: 1px solid var(--border-color);
    border-radius: 4px;
    background-color: var(--input-background);
    color: var(--text-color);
    font-size: 0.9rem;
  }

  .form-group textarea {
    font-family: monospace;
    resize: vertical;
  }

  .switch-label {
    display: flex;
    justify-content: space-between;
    align-items: center;
  }

  .modal-actions {
    display: flex;
    justify-content: flex-end;
    gap: 8px;
    margin-top: 8px;
  }

  .error-message {
    background-color: var(--error-color);
    color: var(--text-color);
    padding: 12px;
    border-radius: 4px;
    margin-bottom: 16px;
  }

  .loading {
    text-align: center;
    padding: 24px;
    color: var(--text-color);
    font-style: italic;
  }

  @media (max-width: 768px) {
    .page {
      padding: 16px;
    }

    .header-actions {
      flex-direction: column;
      gap: 16px;
      align-items: stretch;
    }

    .title {
      text-align: center;
    }

    .buttons-container {
      justify-content: center;
    }

    .action-button {
      width: 100%;
    }

    .flags-container {
      grid-template-columns: 1fr;
    }

    .flag-card {
      padding: 12px;
    }

    .flag-header {
      flex-direction: column;
      gap: 12px;
      align-items: flex-start;
    }

    .flag-actions {
      flex-direction: column;
      width: 100%;
    }

    .action-button {
      width: 100%;
    }

    .modal-actions {
      flex-direction: column;
    }

    .modal-actions .action-button {
      width: 100%;
    }

    .form-group input[type="text"],
    .form-group textarea {
      width: 100%;
    }
  }
</style>