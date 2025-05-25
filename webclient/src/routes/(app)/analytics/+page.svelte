<script lang="ts">
  import BaseView from "$lib/widgets/BaseView.svelte";
  import { onMount } from 'svelte';
  import { Chart, registerables } from 'chart.js/auto';
  import type { AnalyticsSummary } from '$lib/service/analyticsService';
  import analyticsService from "$lib/service/analyticsService";
  import LL from "../../../i18n/i18n-svelte";
  import MetricsCard from "./MetricsCard.svelte";
  import ChartWithList from "./ChartWithList.svelte";
  import TopPagesList from "./TopPagesList.svelte";
  import { formatMilliseconds, formatBytes } from "./utils";

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

  let isLoading = $state(true);

  onMount(async () => {
    Chart.register(...registerables);
    
    const [analytics, error] = await analyticsService.getAnalytics();
    if (error) {
      console.error(error);
    } else {
      metrics = analytics;
    }
    isLoading = false;
  });
</script>

<BaseView>
  <div class="page">
    <h1 class="page-title">Analytics Dashboard</h1>

    <section>
      <h2 class="section-title">Metrics</h2>
      <div class="metrics-grid">
        <MetricsCard title={$LL.analytics.metrics.total_visits()} value={metrics.total_visits?.toString() ?? 'N/A'} loading={isLoading} />
        <MetricsCard title={$LL.analytics.metrics.unique_visitors()} value={metrics.unique_visitors?.toString() ?? 'N/A'} loading={isLoading} />
        <MetricsCard title={$LL.analytics.metrics.active_visitors()} value={metrics.active_visitors?.toString() ?? 'N/A'} loading={isLoading} />
        <MetricsCard title={$LL.analytics.metrics.error_rate()} value={((metrics.error_rate ?? 0) * 100).toFixed(2).toString() ?? 'N/A'} loading={isLoading} />
        <MetricsCard title={$LL.analytics.metrics.average_latency()} value={formatMilliseconds(metrics.average_latency_ms ?? 0)} loading={isLoading} />
        <MetricsCard title={$LL.analytics.metrics.latency_95th_percentile()} value={formatMilliseconds(metrics.p95_latency_ms ?? 0)} loading={isLoading} />
        <MetricsCard title={$LL.analytics.metrics.latency_99th_percentile()} value={formatMilliseconds(metrics.p99_latency_ms ?? 0)} loading={isLoading} />
        <MetricsCard title={$LL.analytics.metrics.outgoing_traffic()} value={formatBytes(metrics.total_bytes_sent ?? 0)} loading={isLoading} />
        <MetricsCard title={$LL.analytics.metrics.incoming_traffic()} value={formatBytes(metrics.total_bytes_recv ?? 0)} loading={isLoading} />
      </div>
    </section>

    <TopPagesList pages={metrics.top_pages} loading={isLoading} />


    <section class="grid-container">
      <ChartWithList 
        title={$LL.analytics.browser_distribution.title()} 
        listHeader={$LL.analytics.browser_distribution.header()} 
        data={metrics.browser_usage} 
        canvasId="browser-chart" 
        loading={isLoading}
      />

      <ChartWithList 
        title={$LL.analytics.device_distribution.title()} 
        listHeader={$LL.analytics.device_distribution.header()} 
        data={metrics.os_usage} 
        canvasId="device-chart" 
        loading={isLoading}
      />

      <ChartWithList 
        title={$LL.analytics.referer_distribution.title()} 
        listHeader={$LL.analytics.referer_distribution.header()} 
        data={metrics.referer_usage} 
        canvasId="referer-chart" 
        loading={isLoading}
      />

      <ChartWithList 
        title={$LL.analytics.instance_distribution.title()} 
        listHeader={$LL.analytics.instance_distribution.header()} 
        data={metrics.instance_stats} 
        canvasId="instance-chart" 
        loading={isLoading}
      />
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

  .page-title {
    font-size: 1.8rem;
    font-weight: 700;
    margin-bottom: 0.5rem;
    color: var(--text-color);
  }

  .section-title {
    font-size: 1.25rem;
    font-weight: 600;
    margin-bottom: 1rem;
    color: var(--text-color);
  }
  
  .metrics-grid {
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(200px, 1fr));
    gap: 1rem;
  }
  
  .grid-container {
    display: flex;
    width: 100%;
    flex-wrap: wrap;
    gap: 1.5rem;
  }

</style>