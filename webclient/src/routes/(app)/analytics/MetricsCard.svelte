<script lang="ts">
  let {title, value, loading = false} : {title: string, value: string, loading?: boolean} = $props();
</script>

<div class="card">
  {#if loading}
    <div class="skeleton-title"></div>
    <div class="skeleton-value"></div>
  {:else}
    <h3 class="title">{title}</h3>
    <p class="value">{value?.toString() ?? 'N/A'}</p>
  {/if}
</div>

<style>
  .card {
    background-color: var(--card-background);
    border-radius: 8px;
    padding: 1.5rem;
    box-shadow: 0 1px 3px var(--shadow-color);
  }

  .title {
    font-size: 0.9rem;
    font-weight: 600;
    margin-bottom: 0.5rem;
    color: var(--text-secondary-color);
  }

  .value {
    font-size: 1.5rem;
    font-weight: 700;
    color: var(--text-color);
  }

  .skeleton-title,
  .skeleton-value {
    background: linear-gradient(90deg, var(--skeleton-base) 25%, var(--skeleton-highlight) 50%, var(--skeleton-base) 75%);
    background-size: 200% 100%;
    animation: shimmer 2.2s infinite;
    border-radius: 4px;
  }

  .skeleton-title {
    height: 1rem;
    width: 70%;
    margin-bottom: 0.5rem;
  }

  .skeleton-value {
    height: 1.5rem;
    width: 50%;
  }

  @keyframes shimmer {
    0% {
      background-position: -200% 0;
    }
    100% {
      background-position: 200% 0;
    }
  }
</style>