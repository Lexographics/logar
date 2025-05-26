<script lang="ts">
  import LL from "../../../i18n/i18n-svelte";

  interface TopPage {
    path: string;
    visits: number;
    percentage: number;
  }

  interface Props {
    pages: TopPage[];
    loading?: boolean;
  }

  let { pages, loading = false }: Props = $props();

  let displayPages = $state<TopPage[]>([]);
  $effect(() => {
    displayPages = (pages?.length ?? 0) > 0 ? pages : Array(5).fill({path: 'N/A', visits: 0, percentage: 0});
  });
</script>

<section class="grid-container">
  <div style="width: 100%;">
    <h2 class="section-title">{$LL.analytics.top_pages.title()}</h2>
    <div class="card">
      {#if loading}
        {#each Array(5) as _}
          <div class="list-item">
            <div class="skeleton-page-path"></div>
            <div class="skeleton-page-value"></div>
          </div>
        {/each}
      {:else}
        {#each displayPages as page}
          <div class="list-item">
            <span class="link-style">{page.path}</span>
            <span class="list-value">{$LL.analytics.top_pages.page_text(page.visits, page.percentage.toFixed(0))}</span>
          </div>
        {/each}
      {/if}
    </div>
  </div>
</section>

<style>
  .section-title {
    font-size: 1.25rem;
    font-weight: 600;
    margin-bottom: 1rem;
    color: var(--text-color);
  }

  .card {
    background-color: var(--card-background);
    border-radius: 8px;
    padding: 1.5rem;
    box-shadow: 0 1px 3px var(--shadow-color);
  }

  .grid-container {
    display: flex;
    width: 100%;
    flex-wrap: wrap;
    gap: 1.5rem;
  }

  .list-item {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 0.5rem 0;
    border-bottom: 1px solid var(--border-color);
    color: var(--text-color);
  }

  .list-item:last-child {
    border-bottom: none;
  }

  .list-value {
    font-weight: 500;
  }

  .link-style {
    color: #2563eb;
    text-decoration: none;
    cursor: pointer;
  }

  .link-style:hover {
    text-decoration: underline;
  }

  .skeleton-page-path,
  .skeleton-page-value {
    background: linear-gradient(90deg, var(--skeleton-base) 25%, var(--skeleton-highlight) 50%, var(--skeleton-base) 75%);
    background-size: 200% 100%;
    animation: shimmer 2.2s infinite;
    border-radius: 4px;
  }

  .skeleton-page-path {
    height: 1rem;
    width: 65%;
  }

  .skeleton-page-value {
    height: 1rem;
    width: 35%;
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