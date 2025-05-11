<script lang="ts">
  import BaseView from "$lib/widgets/BaseView.svelte";
  import { onMount } from 'svelte';
  import { Chart, registerables } from 'chart.js/auto';
  import type { AnalyticsSummary } from '$lib/service/analyticsService';
  import analyticsService from "$lib/service/analyticsService";

  let metrics : AnalyticsSummary = $state({
    total_visits: 0,
    unique_visitors: 0,
    active_visitors: 0,
    error_rate: 0,
    average_latency_ms: 0,
    p95_latency_ms: 0,
    p99_latency_ms: 0,
    total_bytes_sent: 0,
    total_bytes_recv: 0,
    top_pages: [],
    os_usage: {},
    browser_usage: {},
    referer_usage: {},
    instance_stats: {},
  }); 

  onMount(async () => {
    const [analytics, error] = await analyticsService.getAnalytics();
    if (error) {
      console.error(error);
    } else {
      metrics = analytics;
      createCharts();
    }

  });

  function formatMilliseconds(milliseconds: number) : string {
    const seconds = milliseconds / 1000;
    const minutes = seconds / 60;
    const hours = minutes / 60;
    const days = hours / 24;

    if (days >= 1) {
      return `${days.toFixed(1)}d`;
    } else if (hours >= 1) {
      return `${hours.toFixed(1)}h`;
    } else if (minutes >= 1) {
      return `${minutes.toFixed(1)}m`;
    } else if (seconds >= 1) {
      return `${seconds.toFixed(1)}s`;
    } else {
      return `${milliseconds.toFixed(0)}ms`;
    }
  }

  function formatBytes(bytes: number) : string {
    const units = ['B', 'KB', 'MB', 'GB', 'TB'];
    let index = 0;
    let value = bytes;

    while (value >= 1024 && index < units.length - 1) {
      value /= 1024;
      index++;
    }

    return `${value.toFixed(2)} ${units[index]}`;
  }
  let browserChartCanvas = $state<HTMLCanvasElement | null>(null);
  let browserChart = $state<Chart | null>(null);
    
  let deviceChartCanvas = $state<HTMLCanvasElement | null>(null);
  let deviceChart = $state<Chart | null>(null);

  let refererChartCanvas = $state<HTMLCanvasElement | null>(null);
  let refererChart = $state<Chart | null>(null);

  let instanceChartCanvas = $state<HTMLCanvasElement | null>(null);
  let instanceChart = $state<Chart | null>(null);

  function createCharts() {
    const textColor = window.getComputedStyle(document.body).getPropertyValue('--text-color');

    const browserCtx = browserChartCanvas.getContext('2d');
    browserChart = new Chart(browserCtx, {
      type: 'doughnut',
      data: {
        labels: Object.keys(metrics.browser_usage),
        datasets: [{
          label: ' Percentage',
          data: Object.values(metrics.browser_usage),
          backgroundColor: [
            ...Object.keys(metrics.browser_usage).map(label => {
              const hash = label.split('').reduce((acc, char) => acc + Math.pow(2, char.charCodeAt(0)), 0);
              const hue = (hash * 137.5) % 360;
              const saturation = 50 + (hash % 30);
              const lightness = 40 + (hash % 20);
              return `hsl(${hue}, ${saturation}%, ${lightness}%)`;
            })
          ],
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

    const deviceCtx = deviceChartCanvas.getContext('2d');
    deviceChart = new Chart(deviceCtx, {
      type: 'doughnut',
      data: {
        labels: Object.keys(metrics.os_usage),
        datasets: [{
          label: ' Percentage',
          data: Object.values(metrics.os_usage),
          backgroundColor: [
            ...Object.keys(metrics.os_usage).map(label => {
              const hash = label.split('').reduce((acc, char) => acc + Math.pow(2, char.charCodeAt(0)), 0);
              const hue = (hash * 137.5) % 360;
              const saturation = 50 + (hash % 30);
              const lightness = 40 + (hash % 20);
              return `hsl(${hue}, ${saturation}%, ${lightness}%)`;
            })
          ],
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

    const refererCtx = refererChartCanvas.getContext('2d');
    refererChart = new Chart(refererCtx, {
      type: 'doughnut',
      data: {
        labels: Object.keys(metrics.referer_usage),
        datasets: [{
          label: ' Percentage',
          data: Object.values(metrics.referer_usage),
          backgroundColor: [
            ...Object.keys(metrics.referer_usage).map(label => {
              const hash = label.split('').reduce((acc, char) => acc + Math.pow(2, char.charCodeAt(0)), 0);
              const hue = (hash * 137.5) % 360;
              const saturation = 50 + (hash % 30);
              const lightness = 40 + (hash % 20);
              return `hsl(${hue}, ${saturation}%, ${lightness}%)`;
            })
          ],
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

    const instanceCtx = instanceChartCanvas.getContext('2d');
    instanceChart = new Chart(instanceCtx, {
      type: 'doughnut',
      data: {
        labels: Object.keys(metrics.instance_stats),
        datasets: [{
          label: ' Percentage',
          data: Object.values(metrics.instance_stats),
          backgroundColor: [
            ...Object.keys(metrics.instance_stats).map(label => {
              const hash = label.split('').reduce((acc, char) => acc + Math.pow(2, char.charCodeAt(0)), 0);
              const hue = (hash * 137.5) % 360;
              const saturation = 50 + (hash % 30);
              const lightness = 40 + (hash % 20);
              return `hsl(${hue}, ${saturation}%, ${lightness}%)`;
            })
          ],
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

  onMount(() => {
    Chart.register(...registerables);

    // deviceChart.update();
    return () => {
      deviceChart?.destroy();
      browserChart?.destroy();
      refererChart?.destroy();
      instanceChart?.destroy();
    };
  });
</script>

<BaseView>
  <div class="page">
    <h1 class="page-title">Analytics Dashboard</h1>

    <section>
      <h2 class="section-title">Metrics</h2>
      <div class="metrics-grid">
        <div class="card">
          <h3 class="metric-title">Total Visits</h3>
          <p class="metric-value">{metrics.total_visits?.toString() ?? 'N/A'}</p>
        </div>
        <div class="card">
          <h3 class="metric-title">Unique Visitors</h3>
          <p class="metric-value">{metrics.unique_visitors?.toString() ?? 'N/A'}</p>
        </div>
        <div class="card">
          <h3 class="metric-title">Active Visitors</h3>
          <p class="metric-value">{metrics.active_visitors?.toString() ?? 'N/A'}</p>
        </div>
        <div class="card">
          <h3 class="metric-title">Error Rate</h3>
          <p class="metric-value">%{((metrics.error_rate ?? 0) * 100).toFixed(2).toString() ?? 'N/A'}</p>
        </div>
        <div class="card">
          <h3 class="metric-title">Average Latency</h3>
          <p class="metric-value">{formatMilliseconds(metrics.average_latency_ms ?? 0)}</p>
        </div>
        <div class="card">
          <h3 class="metric-title">95th Percentile Latency</h3>
          <p class="metric-value">{formatMilliseconds(metrics.p95_latency_ms ?? 0)}</p>
        </div>
        <div class="card">
          <h3 class="metric-title">99th Percentile Latency</h3>
          <p class="metric-value">{formatMilliseconds(metrics.p99_latency_ms ?? 0)}</p>
        </div>
        <div class="card">
          <h3 class="metric-title">Outgoing Traffic</h3>
          <p class="metric-value">{formatBytes(metrics.total_bytes_sent ?? 0)}</p>
        </div>
        <div class="card">
          <h3 class="metric-title">Incoming Traffic</h3>
          <p class="metric-value">{formatBytes(metrics.total_bytes_recv ?? 0)}</p>
        </div>
      </div>
    </section>

    <section class="grid-container">
      <div style="width: 100%;">
          <h2 class="section-title">Top Pages</h2>
          <div class="card">
            {#each metrics.top_pages as page}
              <div class="list-item">
                <span class="link-style">{page.path}</span>
                <span class="list-value">{page.visits} views ({page.percentage.toFixed(0)}%)</span>
              </div>
            {/each}
          </div>
      </div>
    </section>


    <section class="grid-container">
      <div>
        <h2 class="section-title">Browser Distribution</h2>
          <div class="chart-card-split">
            <div class="list-section">
              <h3 class="card-title">Browser</h3>
              {#each Object.entries(metrics.browser_usage).sort((a, b) => b[1] - a[1]).slice(0, 4) as [browser, percentage]}
                <div class="list-item">
                  <span style="margin-right: 10px;">{browser}</span>
                  <span class="list-value">{percentage.toFixed(1)}%</span>
                </div>
              {/each}
            </div>
             <div class="chart-section">
               <canvas bind:this={browserChartCanvas}></canvas>
             </div>
          </div>
      </div>

      <div>
        <h2 class="section-title">Device Distribution</h2>
          <div class="chart-card-split">
            <div class="list-section">
              <h3 class="card-title">Device</h3>
              {#each Object.entries(metrics.os_usage).sort((a, b) => b[1] - a[1]).slice(0, 4) as [os, percentage]}
                <div class="list-item">
                  <span style="margin-right: 10px;">{os}</span>
                  <span class="list-value">{percentage.toFixed(1)}%</span>
                </div>
              {/each}
            </div>
             <div class="chart-section">
               <canvas bind:this={deviceChartCanvas}></canvas>
             </div>
          </div>
      </div>

      <div>
        <h2 class="section-title">Referer Distribution</h2>
          <div class="chart-card-split">
            <div class="list-section">
              <h3 class="card-title">Referer</h3>
              {#each Object.entries(metrics.referer_usage).sort((a, b) => b[1] - a[1]).slice(0, 4) as [referer, percentage]}
                <div class="list-item">
                  <span style="margin-right: 10px;">{referer}</span>
                  <span class="list-value">{percentage.toFixed(1)}%</span>
                </div>
              {/each}
            </div>
             <div class="chart-section">
               <canvas bind:this={refererChartCanvas}></canvas>
             </div>
          </div>
      </div>

      <div>
        <h2 class="section-title">Instance Distribution</h2>
          <div class="chart-card-split">
            <div class="list-section">
              <h3 class="card-title">Instance</h3>
              {#each Object.entries(metrics.instance_stats).sort((a, b) => b[1] - a[1]).slice(0, 4) as [instance, percentage]}
                <div class="list-item">
                  <span style="margin-right: 10px;">{instance}</span>
                  <span class="list-value">{percentage.toFixed(1)}%</span>
                </div>
              {/each}
            </div>
             <div class="chart-section">
               <canvas bind:this={instanceChartCanvas}></canvas>
             </div>
          </div>
      </div>
    </section>

  </div>
</BaseView>

<style>
  .page {
    padding: 1.5rem;
    display: flex;
    flex-direction: column;
    gap: 2rem;
    color: var(--text-color);
  }

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

  .metric-title {
    font-size: 0.9rem;
    font-weight: 600;
    margin-bottom: 0.5rem;
    color: var(--text-secondary-color);
  }

  .metric-value {
    font-size: 1.5rem;
    font-weight: 700;
    color: var(--text-color);
  }
  
  .page-title {
    font-size: 1.8rem;
    font-weight: 700;
    margin-bottom: 0.5rem;
    color: var(--text-color);
  }
  
  .metrics-grid {
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(200px, 1fr));
    gap: 1rem;
  }
  
  .chart-card-split {
    background-color: var(--card-background);
    border-radius: 8px;
    padding: 1.5rem;
    box-shadow: 0 1px 3px var(--shadow-color);
    margin-bottom: 1.5rem;
  }
  
  .chart-section {
    height: 200px;
  }
  
  .grid-container {
    display: flex;
    width: 100%;
    flex-wrap: wrap;
    gap: 1.5rem;
  }
  
  .chart-card-split {
    display: flex;
    flex-direction: column;
    flex-wrap: wrap;
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

  .link-style {
    color: #2563eb;
    text-decoration: none;
    cursor: pointer;
  }
  .link-style:hover {
    text-decoration: underline;
  }

  .card-title {
    font-size: 1.125rem;
    font-weight: 500;
    margin-bottom: 0.75rem;
  }

</style>