<script lang="ts">
  import * as Plot from '@observablehq/plot';
  import { onMount } from 'svelte';
  import { TrendingUp, TrendingDown, Minus, AlertTriangle, Info, CheckCircle, BarChart3, Maximize2, Filter, Eye, EyeOff } from 'lucide-svelte';

  interface Props {
    data: any;
  }

  let { data }: Props = $props();
  let chartContainer: HTMLDivElement;
  let mounted = $state(false);

  // Interactive features
  let hoveredMetric = $state<string | null>(null);
  let showDetailedBreakdown = $state(false);
  let filteredMetricTypes = $state<string[]>([]);
  let sortBy = $state<'name' | 'change' | 'value'>('change');
  let sortOrder = $state<'asc' | 'desc'>('desc');
  let viewMode = $state<'cards' | 'chart' | 'table'>('cards');

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

  // Interactive functions

  function toggleMetricFilter(metricType: string) {
    if (filteredMetricTypes.includes(metricType)) {
      filteredMetricTypes = filteredMetricTypes.filter(t => t !== metricType);
    } else {
      filteredMetricTypes = [...filteredMetricTypes, metricType];
    }
  }

  function sortComparisons(comparisons: any[]) {
    return [...comparisons].sort((a, b) => {
      let aValue: number, bValue: number;
      
      switch (sortBy) {
        case 'name':
          aValue = a.metric_name.localeCompare(b.metric_name);
          bValue = 0;
          break;
        case 'change':
          aValue = Math.abs(a.change_percent);
          bValue = Math.abs(b.change_percent);
          break;
        case 'value':
          aValue = a.period2.value;
          bValue = b.period2.value;
          break;
        default:
          return 0;
      }
      
      const result = sortBy === 'name' ? aValue : (aValue - bValue);
      return sortOrder === 'asc' ? result : -result;
    });
  }

  let allComparisons = $derived(data?.comparisons || []);
  let insights = $derived(data?.insights || []);
  
  // Filter and sort comparisons
  let comparisons = $derived(() => {
    let filtered = allComparisons;
    
    // Apply metric type filters
    if (filteredMetricTypes.length > 0) {
      filtered = filtered.filter((comp: any) => 
        filteredMetricTypes.some(type => comp.metric_type.includes(type))
      );
    }
    
    return sortComparisons(filtered);
  });

  // Get unique metric types for filtering
  let availableMetricTypes = $derived(() => {
    const types = new Set<string>();
    allComparisons.forEach((comp: any) => {
      if (comp.metric_type.startsWith('question_')) {
        types.add('Questions');
      } else if (comp.metric_type.includes('survey')) {
        types.add('Survey');
      } else {
        types.add('General');
      }
    });
    return Array.from(types);
  });

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
      
      <!-- Interactive Controls -->
      <div class="bg-white border border-gray-200 rounded-xl p-4 mb-6 shadow-sm">
        <div class="flex flex-wrap items-center justify-between gap-4">
          <!-- View Mode Controls -->
          <div class="flex items-center gap-2">
            <span class="text-sm font-medium text-gray-700">View:</span>
            <div class="flex bg-gray-100 rounded-lg p-1">
              <button
                class="px-3 py-1.5 text-sm font-medium rounded-md transition-all {viewMode === 'cards' ? 'bg-white text-gray-900 shadow-sm' : 'text-gray-600 hover:text-gray-900'}"
                onclick={() => viewMode = 'cards'}
              >
                Cards
              </button>
              <button
                class="px-3 py-1.5 text-sm font-medium rounded-md transition-all {viewMode === 'chart' ? 'bg-white text-gray-900 shadow-sm' : 'text-gray-600 hover:text-gray-900'}"
                onclick={() => viewMode = 'chart'}
              >
                Chart
              </button>
              <button
                class="px-3 py-1.5 text-sm font-medium rounded-md transition-all {viewMode === 'table' ? 'bg-white text-gray-900 shadow-sm' : 'text-gray-600 hover:text-gray-900'}"
                onclick={() => viewMode = 'table'}
              >
                Table
              </button>
            </div>
          </div>

          <!-- Sort Controls -->
          <div class="flex items-center gap-2">
            <span class="text-sm font-medium text-gray-700">Sort by:</span>
            <select 
              bind:value={sortBy}
              class="text-sm border border-gray-300 rounded-md px-2 py-1 bg-white"
            >
              <option value="change">Change %</option>
              <option value="value">Current Value</option>
              <option value="name">Name</option>
            </select>
            <button
              class="p-1 text-gray-600 hover:text-blue-600 rounded transition-colors"
              onclick={() => sortOrder = sortOrder === 'asc' ? 'desc' : 'asc'}
              title="Toggle sort order"
            >
              <svg class="w-4 h-4 {sortOrder === 'desc' ? 'rotate-180' : ''} transition-transform" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M7 16V4m0 0L3 8m4-4l4 4m6 0v12m0 0l4-4m-4 4l-4-4"></path>
              </svg>
            </button>
          </div>

          <!-- Filter Controls -->
          {#if availableMetricTypes.length > 1}
            <div class="flex items-center gap-2">
              <span class="text-sm font-medium text-gray-700">Filter:</span>
              <div class="flex gap-1">
                {#each availableMetricTypes as metricType}
                  <button
                    class="px-2 py-1 text-xs rounded-md border transition-all {
                      filteredMetricTypes.length === 0 || filteredMetricTypes.includes(metricType.toLowerCase())
                        ? 'bg-blue-100 border-blue-300 text-blue-700'
                        : 'bg-gray-100 border-gray-300 text-gray-600 hover:bg-gray-200'
                    }"
                    onclick={() => toggleMetricFilter(metricType.toLowerCase())}
                  >
                    {metricType}
                  </button>
                {/each}
              </div>
            </div>
          {/if}
        </div>
      </div>
      
      <!-- Dynamic View Based on Mode -->
      {#if viewMode === 'cards'}
        <!-- Enhanced Interactive Comparison Cards -->
        <div class="grid grid-cols-1 lg:grid-cols-2 gap-6 mb-8">
          {#each comparisons as comparison}
            <div 
              class="border rounded-lg transition-all duration-200 hover:shadow-lg {
                getBackgroundColor(comparison.trend, comparison.change_percent)
              }"
              onmouseenter={() => hoveredMetric = comparison.metric_name}
              onmouseleave={() => hoveredMetric = null}
            >
              <div class="p-6">
                <div class="flex items-start justify-between mb-4">
                  <div class="flex-1">
                    <h4 class="font-semibold text-gray-900 mb-1 {hoveredMetric === comparison.metric_name ? 'text-blue-700' : ''} transition-colors">
                      {comparison.metric_name}
                    </h4>
                    <div class="text-xs text-gray-500 flex items-center gap-2">
                      <span class="px-2 py-1 bg-gray-100 rounded-full capitalize">
                        {comparison.metric_type.replace('_', ' ')}
                      </span>
                    </div>
                  </div>
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
                
                <!-- Period 1 -->
                <div class="transition-all duration-200 {hoveredMetric === comparison.metric_name ? 'bg-white bg-opacity-50 rounded-lg p-2' : ''}">
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
                <div class="transition-all duration-200 {hoveredMetric === comparison.metric_name ? 'bg-white bg-opacity-50 rounded-lg p-2' : ''}">
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

      {:else if viewMode === 'table'}
        <!-- Table View -->
        <div class="bg-white border border-gray-200 rounded-xl overflow-hidden shadow-sm mb-8">
          <div class="overflow-x-auto">
            <table class="w-full">
              <thead class="bg-gray-50 border-b border-gray-200">
                <tr>
                  <th class="px-4 py-3 text-left text-xs font-semibold text-gray-600 uppercase tracking-wider">Metric</th>
                  <th class="px-4 py-3 text-left text-xs font-semibold text-gray-600 uppercase tracking-wider">Type</th>
                  <th class="px-4 py-3 text-right text-xs font-semibold text-gray-600 uppercase tracking-wider">Period 1</th>
                  <th class="px-4 py-3 text-right text-xs font-semibold text-gray-600 uppercase tracking-wider">Period 2</th>
                  <th class="px-4 py-3 text-right text-xs font-semibold text-gray-600 uppercase tracking-wider">Change</th>
                  <th class="px-4 py-3 text-center text-xs font-semibold text-gray-600 uppercase tracking-wider">Trend</th>
                </tr>
              </thead>
              <tbody class="divide-y divide-gray-200">
                {#each comparisons as comparison}
                  <tr 
                    class="hover:bg-gray-50 transition-colors"
                  >
                    <td class="px-4 py-3">
                      <div class="font-medium text-gray-900">{comparison.metric_name}</div>
                    </td>
                    <td class="px-4 py-3">
                      <span class="inline-block px-2 py-1 text-xs bg-gray-100 text-gray-700 rounded-full capitalize">
                        {comparison.metric_type.replace('_', ' ')}
                      </span>
                    </td>
                    <td class="px-4 py-3 text-right font-medium">
                      {formatValue(comparison.period1.value, comparison.metric_type)}
                    </td>
                    <td class="px-4 py-3 text-right font-medium">
                      {formatValue(comparison.period2.value, comparison.metric_type)}
                    </td>
                    <td class="px-4 py-3 text-right">
                      <span class="font-semibold {getTrendColor(comparison.trend, comparison.change_percent)}">
                        {comparison.change_percent > 0 ? '+' : ''}{comparison.change_percent.toFixed(1)}%
                      </span>
                    </td>
                    <td class="px-4 py-3 text-center">
                      {#snippet trendIcon()}
                        {@const TrendIcon = getTrendIcon(comparison.trend)}
                        <TrendIcon class="w-4 h-4 mx-auto {getTrendColor(comparison.trend, comparison.change_percent)}" />
                      {/snippet}
                      {@render trendIcon()}
                    </td>
                  </tr>
                {/each}
              </tbody>
            </table>
          </div>
        </div>

      {:else if viewMode === 'chart'}
        <!-- Enhanced Chart View -->
        <div class="bg-white border border-gray-200 rounded-xl p-6 mb-8 shadow-sm">
          <div class="flex items-center justify-between mb-4">
            <h4 class="text-lg font-semibold text-gray-900">Visual Comparison</h4>
            <div class="flex items-center gap-2 text-sm text-gray-600">
              <BarChart3 class="w-4 h-4" />
              <span>{comparisons.length} metrics</span>
            </div>
          </div>
          
          {#if plotData.length > 0}
            <!-- Observable Plot Chart -->
            <div bind:this={chartContainer} class="chart-container w-full mb-6"></div>
          {/if}
        </div>
      {/if}
      
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