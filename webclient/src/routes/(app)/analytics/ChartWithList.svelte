<script lang="ts">
  import { onDestroy, untrack } from 'svelte';
  import { Chart } from 'chart.js/auto';

  interface Props {
    title: string;
    listHeader: string;
    data: Record<string, number>;
    canvasId: string;
    loading?: boolean;
  }

  let { title, listHeader, data, canvasId, loading = false }: Props = $props();

  let chartCanvas = $state<HTMLCanvasElement | null>(null);
  let chart = $state<Chart | null>(null);

  function generateColor(label: string): string {
    const hash = label.split('').reduce((acc, char) => acc + Math.pow(2, char.charCodeAt(0)), 0);
    const hue = (hash * 137.5) % 360;
    const saturation = 50 + (hash % 30);
    const lightness = 40 + (hash % 20);
    return `hsl(${hue}, ${saturation}%, ${lightness}%)`;
  }

  function createChart() {
    if (!chartCanvas || !data || Object.keys(data).length === 0) return;

    const textColor = window.getComputedStyle(document.body).getPropertyValue('--text-color');
    const ctx = chartCanvas.getContext('2d');
    
    if (chart) {
      chart.destroy();
    }

    chart = new Chart(ctx, {
      type: 'doughnut',
      data: {
        labels: Object.keys(data),
        datasets: [{
          label: ' Percentage',
          data: Object.values(data),
          backgroundColor: Object.keys(data).map(label => generateColor(label)),
          hoverOffset: 4
        }]
      },
      options: {
        responsive: true,
        maintainAspectRatio: true,
        plugins: {
          legend: { labels: { color: textColor } }
        }
      }
    });
  }

  $effect(() => {
    if (chartCanvas && data) {
      untrack(() => {
        createChart();
      });
    }
  });

  onDestroy(() => {
    chart?.destroy();
  });
</script>

<div>
  <h2 class="section-title">{title}</h2>
  <div class="chart-card-split">
    <div class="list-section">
      <h3 class="card-title">{listHeader}</h3>
      {#if loading}
        {#each Array(4) as _}
          <div class="list-item">
            <div class="skeleton-item-name"></div>
            <div class="skeleton-item-value"></div>
          </div>
        {/each}
      {:else}
        {#each Object.entries(data).sort((a, b) => b[1] - a[1]).slice(0, 4) as [item, percentage]}
          <div class="list-item">
            <span style="margin-right: 10px;">{item}</span>
            <span class="list-value">{percentage.toFixed(1)}%</span>
          </div>
        {/each}
      {/if}
    </div>
    <div class="chart-section">
      {#if loading}
        <div class="skeleton-chart"></div>
      {:else}
        <canvas bind:this={chartCanvas} id={canvasId}></canvas>
      {/if}
    </div>
  </div>
</div>

<style>
  .section-title {
    font-size: 1.25rem;
    font-weight: 600;
    margin-bottom: 1rem;
    color: var(--text-color);
  }

  .chart-card-split {
    background-color: var(--card-background);
    border-radius: 8px;
    padding: 1.5rem;
    box-shadow: 0 1px 3px var(--shadow-color);
    margin-bottom: 1.5rem;
    display: flex;
    flex-direction: column;
    flex-wrap: wrap;
  }

  .chart-section {
    height: 200px;
  }

  @media (min-width: 768px) {
    .chart-card-split {
      flex-direction: row;
    }
    
    .list-section {
      flex: 1;
      padding-right: 1rem;
    }
    
    .chart-section {
      flex: 1;
    }
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

  .card-title {
    font-size: 1.125rem;
    font-weight: 500;
    margin-bottom: 0.75rem;
  }

  .skeleton-item-name,
  .skeleton-item-value,
  .skeleton-chart {
    background: linear-gradient(90deg, var(--skeleton-base) 25%, var(--skeleton-highlight) 50%, var(--skeleton-base) 75%);
    background-size: 200% 100%;
    animation: shimmer 1.5s infinite;
    border-radius: 4px;
  }

  .skeleton-item-name {
    height: 1rem;
    width: 60%;
  }

  .skeleton-item-value {
    height: 1rem;
    width: 30%;
  }

  .skeleton-chart {
    width: 100%;
    height: 100%;
    border-radius: 8px;
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