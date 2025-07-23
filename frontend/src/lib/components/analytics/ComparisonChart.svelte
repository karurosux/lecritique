<script lang="ts">
  import * as Plot from '@observablehq/plot';
  import { onMount } from 'svelte';
  import { TrendingUp, TrendingDown, Minus, AlertTriangle, Info, CheckCircle } from 'lucide-svelte';

  interface Props {
    data: any;
  }

  let { data }: Props = $props();
  let chartContainer: HTMLDivElement;
  let mounted = $state(false);

  function formatValue(value: number, metricType: string): string {
    if (metricType.includes('rate') || metricType.includes('completion')) {
      return value.toFixed(1) + '%';
    } else if (metricType.includes('time')) {
      return value.toFixed(1) + 'm';
    } else if (metricType.includes('rating')) {
      return value.toFixed(2);
    } else {
      return Math.round(value).toLocaleString();
    }
  }

  function formatChange(change: number, metricType: string): string {
    const absChange = Math.abs(change);
    if (metricType.includes('rate') || metricType.includes('completion')) {
      return absChange.toFixed(1) + '%';
    } else if (metricType.includes('time')) {
      return absChange.toFixed(1) + 'm';
    } else if (metricType.includes('rating')) {
      return absChange.toFixed(2);
    } else {
      return Math.round(absChange).toLocaleString();
    }
  }

  function getTrendIcon(trend: string) {
    switch (trend) {
      case 'improving':
        return TrendingUp;
      case 'declining':
        return TrendingDown;
      default:
        return Minus;
    }
  }

  function getTrendColor(trend: string, changePercent: number): string {
    if (Math.abs(changePercent) < 5) {
      return 'text-gray-600';
    }
    
    switch (trend) {
      case 'improving':
        return 'text-green-600';
      case 'declining':
        return 'text-red-600';
      default:
        return 'text-gray-600';
    }
  }

  function getBackgroundColor(trend: string, changePercent: number): string {
    if (Math.abs(changePercent) < 5) {
      return 'bg-gray-50 border-gray-200';
    }
    
    switch (trend) {
      case 'improving':
        return 'bg-green-50 border-green-200';
      case 'declining':
        return 'bg-red-50 border-red-200';
      default:
        return 'bg-gray-50 border-gray-200';
    }
  }

  function getInsightIcon(severity: string) {
    switch (severity) {
      case 'warning':
        return AlertTriangle;
      case 'info':
        return Info;
      default:
        return CheckCircle;
    }
  }

  function getInsightColor(severity: string): string {
    switch (severity) {
      case 'warning':
        return 'text-orange-600';
      case 'info':
        return 'text-blue-600';
      default:
        return 'text-green-600';
    }
  }

  function formatPeriodDate(dateString: string): string {
    return new Date(dateString).toLocaleDateString('en-US', { 
      month: 'short', 
      day: 'numeric',
      year: 'numeric'
    });
  }

  let comparisons = $derived(data?.comparisons || []);
  let insights = $derived(data?.insights || []);

  // Process data for Observable Plot bar chart
  let plotData = $derived(() => {
    if (!comparisons.length) return [];
    
    return comparisons.flatMap((comp: any) => [
      {
        metric: comp.metric_name,
        period: 'Period 1',
        value: comp.period1.value,
        metric_type: comp.metric_type
      },
      {
        metric: comp.metric_name,
        period: 'Period 2', 
        value: comp.period2.value,
        metric_type: comp.metric_type
      }
    ]);
  });

  function renderChart() {
    if (!mounted || !chartContainer) return;
    
    console.log('ComparisonChart rendering with data:', { plotData, comparisons, data });
    
    // Clear previous chart
    chartContainer.innerHTML = '';
    
    if (plotData.length === 0) {
      chartContainer.innerHTML = '<div class="text-center py-8 text-gray-500">No comparison data for chart</div>';
      return;
    }
    
    try {
      // Create the plot
    const plot = Plot.plot({
      title: "Period Comparison",
      width: 600,
      height: 300,
      marginTop: 40,
      marginRight: 40,
      marginBottom: 80,
      marginLeft: 100,
      x: {
        label: "Metrics"
      },
      y: {
        label: "Value",
        grid: true
      },
      color: {
        legend: true,
        range: ["#3B82F6", "#10B981"]
      },
      marks: [
        Plot.barY(plotData, {
          x: "metric",
          y: "value",
          fill: "period",
          tip: true
        }),
        Plot.ruleY([0])
      ]
    });
    
    chartContainer.appendChild(plot);
    } catch (error) {
      console.error('Error rendering comparison chart:', error);
      chartContainer.innerHTML = '<div class="text-center py-8 text-red-500">Error rendering comparison chart</div>';
    }
  }

  // Re-render when data changes
  $effect(() => {
    console.log('ComparisonChart effect triggered:', { mounted, plotData, data });
    if (mounted) {
      renderChart();
    }
  });

  onMount(() => {
    mounted = true;
    renderChart();
  });
</script>

<div class="comparison-chart">
  <div class="mb-6">
    <h3 class="text-lg font-semibold text-gray-900 mb-4">Period Comparison</h3>
    
    {#if comparisons.length === 0}
      <div class="text-center py-8 text-gray-500">
        <p>No comparison data available</p>
      </div>
    {:else}
      <!-- Period Summary -->
      {#if data?.request}
        <div class="bg-blue-50 border border-blue-200 rounded-lg p-4 mb-6">
          <h4 class="font-medium text-blue-900 mb-2">Comparison Periods</h4>
          <div class="grid grid-cols-1 md:grid-cols-2 gap-4 text-sm">
            <div>
              <span class="font-medium text-blue-700">Period 1:</span>
              <span class="ml-2">
                {formatPeriodDate(data.request.period1_start)} - {formatPeriodDate(data.request.period1_end)}
              </span>
            </div>
            <div>
              <span class="font-medium text-blue-700">Period 2:</span>
              <span class="ml-2">
                {formatPeriodDate(data.request.period2_start)} - {formatPeriodDate(data.request.period2_end)}
              </span>
            </div>
          </div>
        </div>
      {/if}
      
      <!-- Comparison Cards -->
      <div class="grid grid-cols-1 lg:grid-cols-2 gap-6 mb-8">
        {#each comparisons as comparison}
          <div class="border rounded-lg {getBackgroundColor(comparison.trend, comparison.change_percent)}">
            <div class="p-6">
              <div class="flex items-start justify-between mb-4">
                <h4 class="font-semibold text-gray-900">{comparison.metric_name}</h4>
                <div class="flex items-center gap-2">
                  {#snippet trendIcon()}
                    {@const TrendIcon = getTrendIcon(comparison.trend)}
                    <TrendIcon class="w-5 h-5 {getTrendColor(comparison.trend, comparison.change_percent)}" />
                  {/snippet}
                  {@render trendIcon()}
                  <span class="font-semibold {getTrendColor(comparison.trend, comparison.change_percent)}">
                    {comparison.change_percent > 0 ? '+' : ''}{comparison.change_percent.toFixed(1)}%
                  </span>
                </div>
              </div>
              
              <div class="grid grid-cols-2 gap-6">
                <!-- Period 1 -->
                <div>
                  <h5 class="text-sm font-medium text-gray-700 mb-2">Period 1</h5>
                  <div class="space-y-1 text-sm">
                    <div class="flex justify-between">
                      <span class="text-gray-500">Total:</span>
                      <span class="font-semibold">
                        {formatValue(comparison.period1.value, comparison.metric_type)}
                      </span>
                    </div>
                    <div class="flex justify-between">
                      <span class="text-gray-500">Average:</span>
                      <span class="font-semibold">
                        {formatValue(comparison.period1.average, comparison.metric_type)}
                      </span>
                    </div>
                    <div class="flex justify-between">
                      <span class="text-gray-500">Count:</span>
                      <span class="font-semibold">{comparison.period1.count.toLocaleString()}</span>
                    </div>
                  </div>
                </div>
                
                <!-- Period 2 -->
                <div>
                  <h5 class="text-sm font-medium text-gray-700 mb-2">Period 2</h5>
                  <div class="space-y-1 text-sm">
                    <div class="flex justify-between">
                      <span class="text-gray-500">Total:</span>
                      <span class="font-semibold">
                        {formatValue(comparison.period2.value, comparison.metric_type)}
                      </span>
                    </div>
                    <div class="flex justify-between">
                      <span class="text-gray-500">Average:</span>
                      <span class="font-semibold">
                        {formatValue(comparison.period2.average, comparison.metric_type)}
                      </span>
                    </div>
                    <div class="flex justify-between">
                      <span class="text-gray-500">Count:</span>
                      <span class="font-semibold">{comparison.period2.count.toLocaleString()}</span>
                    </div>
                  </div>
                </div>
              </div>
              
              <!-- Change Summary -->
              <div class="mt-4 pt-4 border-t border-gray-200">
                <div class="flex items-center justify-between">
                  <span class="text-sm text-gray-600">Absolute Change:</span>
                  <span class="font-semibold {getTrendColor(comparison.trend, comparison.change_percent)}">
                    {comparison.change > 0 ? '+' : ''}{formatChange(comparison.change, comparison.metric_type)}
                  </span>
                </div>
              </div>
            </div>
          </div>
        {/each}
      </div>
      
      <!-- Insights -->
      {#if insights.length > 0}
        <div class="mb-6">
          <h4 class="text-lg font-semibold text-gray-900 mb-4">Insights & Recommendations</h4>
          <div class="space-y-4">
            {#each insights as insight}
              <div class="border rounded-lg p-4 {insight.severity === 'warning' ? 'bg-orange-50 border-orange-200' : insight.severity === 'info' ? 'bg-blue-50 border-blue-200' : 'bg-green-50 border-green-200'}">
                <div class="flex items-start gap-3">
                  {#snippet insightIcon()}
                    {@const InsightIcon = getInsightIcon(insight.severity)}
                    <InsightIcon class="w-5 h-5 mt-0.5 {getInsightColor(insight.severity)}" />
                  {/snippet}
                  {@render insightIcon()}
                  <div class="flex-1">
                    <div class="font-medium text-gray-900">{insight.message}</div>
                    {#if insight.recommendation}
                      <div class="mt-2 text-sm text-gray-600">
                        <strong>Recommendation:</strong> {insight.recommendation}
                      </div>
                    {/if}
                    <div class="mt-1 text-xs text-gray-500">
                      {insight.metric_type} • {Math.abs(insight.change).toFixed(1)}% change
                    </div>
                  </div>
                </div>
              </div>
            {/each}
          </div>
        </div>
      {/if}
      
      <!-- Visual Comparison Bar Chart -->
      <div class="bg-white border rounded-lg p-6">
        <h4 class="text-lg font-semibold text-gray-900 mb-4">Visual Comparison</h4>
        
        {#if plotData.length > 0}
          <!-- Observable Plot Chart -->
          <div bind:this={chartContainer} class="chart-container w-full mb-6"></div>
        {/if}
        
        <!-- Detailed Bar Charts -->
        <div class="space-y-6">
          {#each comparisons as comparison}
            <div>
              <div class="flex items-center justify-between mb-2">
                <h5 class="font-medium text-gray-700">{comparison.metric_name}</h5>
                <span class="text-sm font-semibold {getTrendColor(comparison.trend, comparison.change_percent)}">
                  {comparison.change_percent > 0 ? '+' : ''}{comparison.change_percent.toFixed(1)}%
                </span>
              </div>
              
              <!-- Bar chart -->
              <div class="flex items-center gap-4">
                <div class="flex-1">
                  <div class="flex items-center gap-2 mb-1">
                    <span class="text-xs text-gray-500 w-16">Period 1</span>
                    <div class="flex-1 bg-gray-200 rounded-full h-2">
                      <div 
                        class="bg-blue-500 h-2 rounded-full" 
                        style="width: {Math.max(5, (comparison.period1.value / Math.max(comparison.period1.value, comparison.period2.value)) * 100)}%"
                      ></div>
                    </div>
                    <span class="text-xs font-medium text-gray-700 w-20 text-right">
                      {formatValue(comparison.period1.value, comparison.metric_type)}
                    </span>
                  </div>
                  
                  <div class="flex items-center gap-2">
                    <span class="text-xs text-gray-500 w-16">Period 2</span>
                    <div class="flex-1 bg-gray-200 rounded-full h-2">
                      <div 
                        class="bg-green-500 h-2 rounded-full" 
                        style="width: {Math.max(5, (comparison.period2.value / Math.max(comparison.period1.value, comparison.period2.value)) * 100)}%"
                      ></div>
                    </div>
                    <span class="text-xs font-medium text-gray-700 w-20 text-right">
                      {formatValue(comparison.period2.value, comparison.metric_type)}
                    </span>
                  </div>
                </div>
              </div>
            </div>
          {/each}
        </div>
      </div>
    {/if}
  </div>
</div>

<style>
  .chart-container :global(svg) {
    width: 100% !important;
    height: auto !important;
  }
  
  .chart-container :global(.plot-title) {
    font-size: 16px;
    font-weight: 600;
    fill: #1f2937;
  }
</style>