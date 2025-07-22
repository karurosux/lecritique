<script lang="ts">
  import { Card, Button } from '$lib/components/ui';
  import { Star, BarChart3, TrendingUp, Utensils, ChevronDown, ChevronUp, Search, Filter, Layers, Grid, Eye, EyeOff } from 'lucide-svelte';
  import ChartDataWidget from './ChartDataWidget.svelte';

  interface BackendChartData {
    question_id: string;
    question_text: string;
    question_type: string;
    chart_type: string;
    product_id?: string;
    product_name?: string;
    data: {
      scale?: number;
      distribution?: Record<string, number>;
      average?: number;
      total?: number;
      percentages?: Record<string, number>;
      options?: Record<string, number>;
      is_multi_choice?: boolean;
      combinations?: Array<{
        options: string[];
        count: number;
        percentage: number;
      }>;
      positive?: number;
      neutral?: number;
      negative?: number;
      samples?: string[];
      keywords?: string[];
    };
  }

  interface OrganizationChartData {
    organization_id: string;
    charts: BackendChartData[];
    summary: {
      total_responses: number;
      date_range: {
        start: string;
        end: string;
      };
      filters_applied: Record<string, any>;
    };
  }

  interface ProductGroup {
    product_id: string;
    product_name: string;
    charts: BackendChartData[];
    totalResponses: number;
    averageScore?: number;
  }

  let {
    chartData = null,
    title = "Analytics Dashboard",
    groupByProduct = true,
    initiallyExpanded = false,
    searchQuery = $bindable(''),
    showOnlyWithData = $bindable(false),
    viewMode = $bindable<'grouped' | 'all'>(groupByProduct ? 'grouped' : 'all'),
    hideControls = false
  }: {
    chartData: OrganizationChartData | null;
    title?: string;
    groupByProduct?: boolean;
    initiallyExpanded?: boolean;
    searchQuery?: string;
    showOnlyWithData?: boolean;
    viewMode?: 'grouped' | 'all';
    hideControls?: boolean;
  } = $props();

  let expandedGroups = $state(new Set<string>());

  // Group charts by product
  const productGroups = $derived(() => {
    if (!chartData?.charts) return [];
    
    const groups = new Map<string, ProductGroup>();
    
    chartData.charts.forEach(chart => {
      const productId = chart.product_id || 'no-product';
      const productName = chart.product_name || 'General Questions';
      
      if (!groups.has(productId)) {
        groups.set(productId, {
          product_id: productId,
          product_name: productName,
          charts: [],
          totalResponses: 0,
          averageScore: undefined
        });
      }
      
      const group = groups.get(productId)!;
      group.charts.push(chart);
      group.totalResponses += chart.data.total || 0;
      
      // Calculate average score for rating/scale questions
      if ((chart.chart_type === 'rating' || chart.chart_type === 'scale') && chart.data.average) {
        if (!group.averageScore) {
          group.averageScore = chart.data.average;
        } else {
          // Simple average (could be weighted by responses if needed)
          group.averageScore = (group.averageScore + chart.data.average) / 2;
        }
      }
    });
    
    return Array.from(groups.values()).sort((a, b) => b.totalResponses - a.totalResponses);
  });

  // Filter groups based on search and data availability
  const filteredGroups = $derived(() => {
    return productGroups().filter(group => {
      // Search filter
      if (searchQuery) {
        const query = searchQuery.toLowerCase();
        const matchesProduct = group.product_name.toLowerCase().includes(query);
        const matchesQuestion = group.charts.some(chart => 
          chart.question_text.toLowerCase().includes(query)
        );
        if (!matchesProduct && !matchesQuestion) return false;
      }
      
      // Data filter
      if (showOnlyWithData && group.totalResponses === 0) {
        return false;
      }
      
      return true;
    });
  });

  // Filtered charts for "all" view
  const filteredCharts = $derived(() => {
    if (!chartData?.charts) return [];
    
    return chartData.charts.filter(chart => {
      // Search filter
      if (searchQuery) {
        const query = searchQuery.toLowerCase();
        const matchesQuestion = chart.question_text.toLowerCase().includes(query);
        const matchesProduct = chart.product_name?.toLowerCase().includes(query) || false;
        if (!matchesQuestion && !matchesProduct) return false;
      }
      
      // Data filter
      if (showOnlyWithData && (!chart.data.total || chart.data.total === 0)) {
        return false;
      }
      
      return true;
    });
  });

  // Summary stats
  const summaryStats = $derived(() => {
    if (!chartData?.charts) return null;
    
    const totalProductes = new Set(chartData.charts.map(c => c.product_id || 'no-product')).size;
    const totalQuestions = chartData.charts.length;
    const totalResponses = chartData.summary.total_responses;
    const avgResponsesPerProduct = totalProductes > 0 ? Math.round(totalResponses / totalProductes) : 0;
    
    // Calculate overall average score
    let scoreSum = 0;
    let scoreCount = 0;
    chartData.charts.forEach(chart => {
      if ((chart.chart_type === 'rating' || chart.chart_type === 'scale') && chart.data.average) {
        scoreSum += chart.data.average * (chart.data.total || 1);
        scoreCount += chart.data.total || 1;
      }
    });
    const overallAvgScore = scoreCount > 0 ? scoreSum / scoreCount : 0;
    
    return {
      totalProductes,
      totalQuestions,
      totalResponses,
      avgResponsesPerProduct,
      overallAvgScore
    };
  });

  function toggleGroup(productId: string) {
    if (expandedGroups.has(productId)) {
      expandedGroups.delete(productId);
    } else {
      expandedGroups.add(productId);
    }
    expandedGroups = new Set(expandedGroups);
  }

  function expandAll() {
    filteredGroups().forEach(group => {
      expandedGroups.add(group.product_id);
    });
    expandedGroups = new Set(expandedGroups);
  }

  function collapseAll() {
    expandedGroups.clear();
    expandedGroups = new Set(expandedGroups);
  }

  // Initialize expanded state
  $effect(() => {
    if (initiallyExpanded && productGroups().length > 0 && expandedGroups.size === 0) {
      expandAll();
    }
  });
</script>

{#if !chartData || !chartData.charts || chartData.charts.length === 0}
  <Card variant="elevated">
    <div class="text-center py-20">
      <div class="relative mb-8">
        <div class="h-24 w-24 bg-gradient-to-br from-gray-100 to-gray-200 rounded-3xl flex items-center justify-center mx-auto shadow-lg">
          <BarChart3 class="h-12 w-12 text-gray-400" />
        </div>
      </div>
      <h3 class="text-2xl font-bold bg-gradient-to-r from-gray-900 to-gray-700 bg-clip-text text-transparent mb-3">
        No Analytics Data Available
      </h3>
      <p class="text-gray-600 max-w-lg mx-auto text-lg leading-relaxed">
        Start collecting customer feedback to unlock powerful insights and beautiful analytics visualizations.
      </p>
    </div>
  </Card>
{:else}
  <div class="analytics-charts-grouped space-y-6">
    {#if title}
      <div class="mb-6">
        <div class="flex items-center gap-3 mb-3">
          <div class="h-8 w-8 bg-gradient-to-br from-purple-500 to-pink-600 rounded-lg flex items-center justify-center">
            <BarChart3 class="h-4 w-4 text-white" />
          </div>
          <div>
            <h2 class="text-xl font-semibold text-gray-900">{title}</h2>
            <p class="text-sm text-gray-600">
              {chartData.summary.total_responses.toLocaleString()} responses
              {#if chartData.summary.date_range.start}
                • {new Date(chartData.summary.date_range.start).toLocaleDateString()} 
                to {new Date(chartData.summary.date_range.end).toLocaleDateString()}
              {/if}
            </p>
          </div>
        </div>
      </div>
    {/if}


    {#if !hideControls}
      <!-- Compact Search Bar -->
      <div class="flex flex-wrap items-center gap-3 mb-4">
        <div class="relative flex-1 max-w-xs">
          <Search class="absolute left-3 top-1/2 transform -translate-y-1/2 h-4 w-4 text-gray-400" />
          <input
            type="text"
            bind:value={searchQuery}
            placeholder="Search..."
            class="w-full pl-9 pr-3 py-2 text-sm border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-purple-500 focus:border-transparent"
          />
        </div>
        
        <div class="flex items-center gap-1 bg-gray-100 rounded-lg p-1">
          <button
            class="px-3 py-1.5 text-sm font-medium rounded-md transition-colors {viewMode === 'grouped' ? 'bg-white text-gray-900 shadow-sm' : 'text-gray-600 hover:text-gray-900'}"
            onclick={() => viewMode = 'grouped'}
          >
            <Layers class="h-4 w-4 inline mr-1" />
            Grouped
          </button>
          <button
            class="px-3 py-1.5 text-sm font-medium rounded-md transition-colors {viewMode === 'all' ? 'bg-white text-gray-900 shadow-sm' : 'text-gray-600 hover:text-gray-900'}"
            onclick={() => viewMode = 'all'}
          >
            <Grid class="h-4 w-4 inline mr-1" />
            All
          </button>
        </div>
        
        <label class="flex items-center gap-2 cursor-pointer text-sm">
          <input
            type="checkbox"
            bind:checked={showOnlyWithData}
            class="w-4 h-4 text-purple-600 bg-gray-100 border-gray-300 rounded focus:ring-purple-500"
          />
          <span class="text-gray-700">Data only</span>
        </label>
      </div>
    {/if}

    {#if viewMode === 'grouped'}
      <!-- Grouped View -->
      {#if filteredGroups().length > 0}

        <!-- Product Groups -->
        <div class="space-y-4">
          {#each filteredGroups() as group (group.product_id)}
            {@const isExpanded = expandedGroups.has(group.product_id)}
            <Card variant="elevated" padding={false} class="overflow-hidden">
              <!-- Group Header -->
              <button
                onclick={() => toggleGroup(group.product_id)}
                class="w-full px-6 py-4 flex items-center justify-between hover:bg-gray-50 transition-colors"
              >
                <div class="flex items-center gap-4">
                  <div class="h-10 w-10 bg-gradient-to-br from-purple-500 to-pink-600 rounded-xl flex items-center justify-center shadow-sm">
                    <Utensils class="h-5 w-5 text-white" />
                  </div>
                  <div class="text-left">
                    <h3 class="text-lg font-semibold text-gray-900">
                      {group.product_name}
                    </h3>
                    <div class="flex items-center gap-4 text-sm text-gray-600 mt-1">
                      <span>{group.charts.length} questions</span>
                      <span>•</span>
                      <span>{group.totalResponses.toLocaleString()} responses</span>
                    </div>
                  </div>
                </div>
                
                <div class="flex items-center gap-3">
                  {#if group.totalResponses === 0}
                    <span class="px-2 py-1 bg-gray-100 text-gray-600 text-xs font-medium rounded">
                      No data
                    </span>
                  {/if}
                  {#if isExpanded}
                    <ChevronUp class="h-5 w-5 text-gray-400" />
                  {:else}
                    <ChevronDown class="h-5 w-5 text-gray-400" />
                  {/if}
                </div>
              </button>

              <!-- Group Content -->
              {#if isExpanded}
                <div class="border-t border-gray-200 p-6">
                  <ChartDataWidget 
                    chartData={{
                      ...chartData,
                      charts: group.charts
                    }}
                    title=""
                  />
                </div>
              {/if}
            </Card>
          {/each}
        </div>
      {:else}
        <!-- No Results -->
        <Card variant="elevated">
          <div class="text-center py-12">
            <Search class="h-12 w-12 text-gray-400 mx-auto mb-4" />
            <p class="text-gray-600">No products or questions match your search criteria.</p>
            {#if searchQuery || showOnlyWithData}
              <Button
                variant="outline"
                size="sm"
                onclick={() => {
                  searchQuery = '';
                  showOnlyWithData = false;
                }}
                class="mt-4"
              >
                Clear Filters
              </Button>
            {/if}
          </div>
        </Card>
      {/if}
    {:else}
      <!-- All Charts View -->
      {#if filteredCharts().length > 0}
        <ChartDataWidget 
          chartData={{
            ...chartData,
            charts: filteredCharts()
          }}
          title=""
        />
      {:else}
        <!-- No Results -->
        <Card variant="elevated">
          <div class="text-center py-12">
            <Search class="h-12 w-12 text-gray-400 mx-auto mb-4" />
            <p class="text-gray-600">No charts match your search criteria.</p>
            {#if searchQuery || showOnlyWithData}
              <Button
                variant="outline"
                size="sm"
                onclick={() => {
                  searchQuery = '';
                  showOnlyWithData = false;
                }}
                class="mt-4"
              >
                Clear Filters
              </Button>
            {/if}
          </div>
        </Card>
      {/if}
    {/if}
  </div>
{/if}

<style>
  /* Smooth transition for expand/collapse */
  .analytics-charts-grouped :global(.transition-all) {
    transition: all 0.3s ease;
  }
</style>
