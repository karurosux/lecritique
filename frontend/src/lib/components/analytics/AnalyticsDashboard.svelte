<script lang="ts">
  import { Card } from '$lib/components/ui';
  import { getContext } from 'svelte';
  import type { analyticsStore as AnalyticsStoreType } from '$lib/stores/analytics';
  import { ANALYTICS_CONTEXT_KEY } from '$lib/stores/analytics';
  import { TrendingUpIcon, TrendingDownIcon, MinusIcon } from 'lucide-svelte';

  const analytics = getContext<typeof AnalyticsStoreType>(
    ANALYTICS_CONTEXT_KEY
  );
  const { filters, filteredData, computedMetrics, comparisonData } = analytics;

  let data = $derived($filteredData);
  let metrics = $derived($computedMetrics);
  let comparison = $derived($comparisonData);
  let timeframe = $derived($filters.timeframe);

  const keyMetrics = $derived(() => {
    if (!data || !metrics) return [];

    const metricsArray = [
      {
        label: 'Satisfaction Index',
        value: metrics?.satisfactionIndex || 0,
        format: (v: number) => (v != null ? `${v.toFixed(1)}%` : '0.0%'),
        change: comparison?.satisfaction?.change || 0,
        icon: 'satisfaction',
      },
      {
        label: 'Response Rate',
        value: metrics?.responseRate || 0,
        format: (v: number) => (v != null ? v.toString() : '0'),
        change: comparison?.responseRate?.change || 0,
        icon: 'responses',
      },
      {
        label: 'Sentiment Score',
        value: metrics?.sentimentScore || 0,
        format: (v: number) => (v != null ? `${v.toFixed(1)}%` : '0.0%'),
        change: comparison?.sentiment?.change || 0,
        icon: 'sentiment',
      },
      {
        label: 'Improvement Rate',
        value: metrics?.improvementRate || 0,
        format: (v: number) =>
          v != null ? `${v > 0 ? '+' : ''}${v.toFixed(1)}%` : '0.0%',
        change: 0,
        icon: 'improvement',
      },
    ];

    return metricsArray;
  });

  function getChangeIcon(change: number) {
    if (change > 0) return TrendingUpIcon;
    if (change < 0) return TrendingDownIcon;
    return MinusIcon;
  }

  function getChangeColor(change: number) {
    if (change > 0) return 'text-green-600';
    if (change < 0) return 'text-red-600';
    return 'text-gray-500';
  }

  function handleMetricClick(metric: string) {
    analytics.updateSelection({ highlightedMetric: metric });
  }
</script>

<div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4 mb-6">
  {#each keyMetrics() as metric}
    <Card
      variant="default"
      class="cursor-pointer hover:shadow-lg transition-all analytics-metric-card"
      onclick={() => handleMetricClick(metric.icon)}>
      <div class="p-4">
        <div class="flex items-center justify-between mb-2">
          <span class="text-sm font-medium text-gray-600">{metric.label}</span>
          {#if metric.change !== 0}
            {@const Icon = getChangeIcon(metric.change)}
            <Icon class="w-4 h-4 {getChangeColor(metric.change)}" />
          {/if}
        </div>

        <div class="flex items-baseline justify-between">
          <span class="text-2xl font-bold text-gray-900">
            {metric.format(metric.value)}
          </span>

          {#if metric.change !== 0}
            <span class="text-sm {getChangeColor(metric.change)}">
              {metric.change > 0 ? '+' : ''}{metric.change?.toFixed(1) ||
                '0.0'}%
            </span>
          {/if}
        </div>

        {#if comparison}
          <div class="mt-2 text-xs text-gray-500">
            vs. previous {timeframe === '24h'
              ? 'day'
              : timeframe === '7d'
                ? 'week'
                : 'month'}
          </div>
        {/if}
      </div>
    </Card>
  {/each}
</div>

<!-- Time-based Insights -->
{#if data?.trends}
  <Card variant="default" class="analytics-trends-card">
    <div class="p-6">
      <h3 class="text-lg font-semibold text-gray-900 mb-4">
        Performance Trends
      </h3>

      <div class="grid grid-cols-1 md:grid-cols-3 gap-4">
        <!-- Peak Hours -->
        <div class="bg-blue-50 rounded-lg p-4">
          <h4 class="text-sm font-medium text-blue-900 mb-2">Peak Hours</h4>
          <div class="space-y-1 text-sm text-blue-700">
            <div>Lunch: 12-2 PM (45% of feedback)</div>
            <div>Dinner: 7-9 PM (38% of feedback)</div>
          </div>
        </div>

        <!-- Best Days -->
        <div class="bg-green-50 rounded-lg p-4">
          <h4 class="text-sm font-medium text-green-900 mb-2">
            Best Performing Days
          </h4>
          <div class="space-y-1 text-sm text-green-700">
            <div>Friday: 4.6 avg rating</div>
            <div>Saturday: 4.5 avg rating</div>
          </div>
        </div>

        <!-- Areas of Concern -->
        <div class="bg-amber-50 rounded-lg p-4">
          <h4 class="text-sm font-medium text-amber-900 mb-2">
            Attention Needed
          </h4>
          <div class="space-y-1 text-sm text-amber-700">
            <div>Monday mornings: Lower ratings</div>
            <div>Late service: Decreasing trend</div>
          </div>
        </div>
      </div>
    </div>
  </Card>
{/if}
