<script lang="ts">
  import { Card, NoDataAvailable } from '$lib/components/ui';
  import {
    BarChart3,
    Package,
    ChevronDown,
    ChevronUp,
    Search,
    Layers,
    Grid,
    Loader2,
  } from 'lucide-svelte';
  import ChartDataWidget from './ChartDataWidget.svelte';
  import ChartContent from './ChartContent.svelte';
  import { getApiClient, handleApiError } from '$lib/api/client';

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
    title = 'Analytics Dashboard',
    groupByProduct = true,
    initiallyExpanded = false,
    searchQuery = $bindable(''),
    showOnlyWithData = $bindable(false),
    viewMode = $bindable<'grouped' | 'all'>(groupByProduct ? 'grouped' : 'all'),
    hideControls = false,
    lazyLoading = false,
    organizationId = '',
    filters = {},
    availableProducts = [],
  }: {
    chartData: OrganizationChartData | null;
    title?: string;
    groupByProduct?: boolean;
    initiallyExpanded?: boolean;
    searchQuery?: string;
    showOnlyWithData?: boolean;
    viewMode?: 'grouped' | 'all';
    hideControls?: boolean;
    lazyLoading?: boolean;
    organizationId?: string;
    filters?: Record<string, any>;
    availableProducts?: any[];
  } = $props();

  let expandedGroups = $state(new Set<string>());
  let loadingGroups = $state(new Set<string>());
  let groupChartData = $state(new Map<string, BackendChartData[]>());
  let lastFilters = $state<string>('');

  // Group charts by product
  const productGroups = $derived(() => {
    const groups = new Map<string, ProductGroup>();

    // For lazy loading, prioritize availableProducts when we have them
    if (lazyLoading && availableProducts.length > 0) {
      availableProducts.forEach(product => {
        const productId = product.id;
        const productName = product.name;
        const loadedCharts = groupChartData.get(productId) || [];

        // Calculate totals from loaded chart data
        let totalResponses = 0;
        let averageScore: number | undefined = undefined;

        if (loadedCharts.length > 0) {
          // Sum up responses from all charts for this product
          totalResponses = loadedCharts.reduce(
            (sum, chart) => sum + (chart.data.total || 0),
            0
          );

          // Calculate average score from loaded charts
          const scoredCharts = loadedCharts.filter(
            chart =>
              (chart.chart_type === 'rating' || chart.chart_type === 'scale') &&
              chart.data.average !== undefined
          );

          if (scoredCharts.length > 0) {
            const weightedSum = scoredCharts.reduce(
              (sum, chart) =>
                sum + chart.data.average! * (chart.data.total || 1),
              0
            );
            const totalWeight = scoredCharts.reduce(
              (sum, chart) => sum + (chart.data.total || 1),
              0
            );
            averageScore =
              totalWeight > 0 ? weightedSum / totalWeight : undefined;
          }
        }

        groups.set(productId, {
          product_id: productId,
          product_name: productName,
          charts: loadedCharts,
          totalResponses: loadedCharts.length > 0 ? totalResponses : -1, // -1 indicates unloaded
          averageScore,
        });
      });

      return Array.from(groups.values());
    } else if (chartData?.charts) {
      // Regular mode or lazy loading with existing chart data
      chartData.charts.forEach(chart => {
        const productId = chart.product_id || 'no-product';
        const productName = chart.product_name || 'General Questions';

        if (!groups.has(productId)) {
          groups.set(productId, {
            product_id: productId,
            product_name: productName,
            charts: [],
            totalResponses: 0,
            averageScore: undefined,
          });
        }

        const group = groups.get(productId)!;

        // For lazy loading, only store summary info initially
        if (lazyLoading) {
          // For summary data, we just want to count unique responses per product
          // The chart.data.total in summary mode represents responses to that specific question
          group.totalResponses += chart.data.total || 0;

          // Don't populate the charts array initially - this will be populated lazily
          if (groupChartData.has(productId)) {
            // If we have loaded chart data, use it
            group.charts = groupChartData.get(productId) || [];

            // Calculate average score from loaded data
            group.charts.forEach(loadedChart => {
              if (
                (loadedChart.chart_type === 'rating' ||
                  loadedChart.chart_type === 'scale') &&
                loadedChart.data.average
              ) {
                if (!group.averageScore) {
                  group.averageScore = loadedChart.data.average;
                } else {
                  group.averageScore =
                    (group.averageScore + loadedChart.data.average) / 2;
                }
              }
            });
          }
          // In lazy mode, keep charts empty until loaded
          if (!groupChartData.has(productId)) {
            group.charts = [];
          }
        } else {
          // Non-lazy loading: use full chart data as before
          group.charts.push(chart);
          group.totalResponses += chart.data.total || 0;

          // Calculate average score for rating/scale questions
          if (
            (chart.chart_type === 'rating' || chart.chart_type === 'scale') &&
            chart.data.average
          ) {
            if (!group.averageScore) {
              group.averageScore = chart.data.average;
            } else {
              // Simple average (could be weighted by responses if needed)
              group.averageScore =
                (group.averageScore + chart.data.average) / 2;
            }
          }
        }
      });
    }

    return Array.from(groups.values()).sort(
      (a, b) => b.totalResponses - a.totalResponses
    );
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
        const matchesQuestion = chart.question_text
          .toLowerCase()
          .includes(query);
        const matchesProduct =
          chart.product_name?.toLowerCase().includes(query) || false;
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

    const totalProducts = new Set(
      chartData.charts.map(c => c.product_id || 'no-product')
    ).size;
    const totalQuestions = chartData.charts.length;
    const totalResponses = chartData.summary.total_responses;
    const avgResponsesPerProduct =
      totalProducts > 0 ? Math.round(totalResponses / totalProducts) : 0;

    // Calculate overall average score
    let scoreSum = 0;
    let scoreCount = 0;
    chartData.charts.forEach(chart => {
      if (
        (chart.chart_type === 'rating' || chart.chart_type === 'scale') &&
        chart.data.average
      ) {
        scoreSum += chart.data.average * (chart.data.total || 1);
        scoreCount += chart.data.total || 1;
      }
    });
    const overallAvgScore = scoreCount > 0 ? scoreSum / scoreCount : 0;

    return {
      totalProducts,
      totalQuestions,
      totalResponses,
      avgResponsesPerProduct,
      overallAvgScore,
    };
  });

  async function toggleGroup(productId: string) {
    if (expandedGroups.has(productId)) {
      expandedGroups.delete(productId);
    } else {
      expandedGroups.add(productId);

      // If lazy loading is enabled and we don't have data for this group, load it
      if (lazyLoading && !groupChartData.has(productId) && organizationId) {
        await loadGroupData(productId);
      }
    }
    expandedGroups = new Set(expandedGroups);
  }

  async function loadGroupData(productId: string) {
    if (loadingGroups.has(productId)) return; // Already loading

    loadingGroups.add(productId);
    loadingGroups = new Set(loadingGroups);

    try {
      const api = getApiClient();
      const chartParams = {
        ...filters,
        product_id: productId === 'no-product' ? undefined : productId,
      };

      const response = await api.api.v1AnalyticsOrganizationsChartsList(
        organizationId,
        chartParams
      );

      if (response.data?.data?.charts) {
        // Filter charts for this specific product
        const productCharts = response.data.data.charts.filter(
          (chart: BackendChartData) => {
            const chartProductId = chart.product_id || 'no-product';
            return chartProductId === productId;
          }
        );

        groupChartData.set(productId, productCharts);
        groupChartData = new Map(groupChartData);
      }
    } catch (err) {
      console.error(`Error loading charts for product ${productId}:`, err);
    } finally {
      loadingGroups.delete(productId);
      loadingGroups = new Set(loadingGroups);
    }
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
    if (
      initiallyExpanded &&
      productGroups().length > 0 &&
      expandedGroups.size === 0
    ) {
      expandAll();
    }
  });

  // Watch for filter changes and reload expanded sections
  $effect(() => {
    if (!lazyLoading || !organizationId) return;

    const currentFilters = JSON.stringify(filters);

    // If filters changed and we have cached data, reload expanded sections
    if (
      lastFilters &&
      lastFilters !== currentFilters &&
      expandedGroups.size > 0
    ) {
      // Clear cached data for all groups
      groupChartData.clear();
      groupChartData = new Map(groupChartData);

      // Reload data for all currently expanded groups
      const expandedArray = Array.from(expandedGroups);
      expandedArray.forEach(productId => {
        loadGroupData(productId);
      });
    }

    lastFilters = currentFilters;
  });
</script>

{#if lazyLoading ? availableProducts.length === 0 : !chartData || !chartData.charts || chartData.charts.length === 0}
  <NoDataAvailable
    title="No Analytics Data Available"
    description="Start collecting customer feedback to unlock powerful insights and beautiful analytics visualizations."
    icon={BarChart3} />
{:else}
  <div class="analytics-charts-grouped space-y-6">
    {#if title}
      <div class="mb-6">
        <div class="flex items-center gap-3 mb-3">
          <div
            class="h-8 w-8 bg-gradient-to-br from-purple-500 to-pink-600 rounded-lg flex items-center justify-center">
            <BarChart3 class="h-4 w-4 text-white" />
          </div>
          <div>
            <h2 class="text-xl font-semibold text-gray-900">{title}</h2>
            <p class="text-sm text-gray-600">
              {chartData.summary.total_responses.toLocaleString()} responses
              {#if chartData.summary.date_range.start}
                • {new Date(
                  chartData.summary.date_range.start
                ).toLocaleDateString()}
                to {new Date(
                  chartData.summary.date_range.end
                ).toLocaleDateString()}
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
          <Search
            class="absolute left-3 top-1/2 transform -translate-y-1/2 h-4 w-4 text-gray-400" />
          <input
            type="text"
            bind:value={searchQuery}
            placeholder="Search..."
            class="w-full pl-9 pr-3 py-2 text-sm border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-purple-500 focus:border-transparent" />
        </div>

        <div class="flex items-center gap-1 bg-gray-100 rounded-lg p-1">
          <button
            class="px-3 py-1.5 text-sm font-medium rounded-md transition-colors {viewMode ===
            'grouped'
              ? 'bg-white text-gray-900 shadow-sm'
              : 'text-gray-600 hover:text-gray-900'}"
            onclick={() => (viewMode = 'grouped')}>
            <Layers class="h-4 w-4 inline mr-1" />
            Grouped
          </button>
          <button
            class="px-3 py-1.5 text-sm font-medium rounded-md transition-colors {viewMode ===
            'all'
              ? 'bg-white text-gray-900 shadow-sm'
              : 'text-gray-600 hover:text-gray-900'}"
            onclick={() => (viewMode = 'all')}>
            <Grid class="h-4 w-4 inline mr-1" />
            All
          </button>
        </div>

        <label class="flex items-center gap-2 cursor-pointer text-sm">
          <input
            type="checkbox"
            bind:checked={showOnlyWithData}
            class="w-4 h-4 text-purple-600 bg-gray-100 border-gray-300 rounded focus:ring-purple-500" />
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
            <Card variant="minimal" padding={false} class="overflow-hidden">
              <!-- Group Header -->
              <button
                onclick={() => toggleGroup(group.product_id)}
                class="w-full px-6 py-4 flex items-center justify-between hover:bg-gray-50 transition-colors cursor-pointer">
                <div class="flex items-center gap-4">
                  <div
                    class="h-10 w-10 bg-gradient-to-br from-purple-500 to-pink-600 rounded-xl flex items-center justify-center shadow-sm">
                    <Package class="h-5 w-5 text-white" />
                  </div>
                  <div class="text-left">
                    <h3 class="text-lg font-semibold text-gray-900">
                      {group.product_name}
                    </h3>
                    <div
                      class="flex items-center gap-4 text-sm text-gray-600 mt-1">
                      <span
                        >{group.totalResponses === -1
                          ? '--'
                          : group.charts.length} questions</span>
                      <span>•</span>
                      <span
                        >{group.totalResponses === -1
                          ? '--'
                          : group.totalResponses.toLocaleString()} responses</span>
                    </div>
                  </div>
                </div>

                <div class="flex items-center gap-3">
                  {#if group.totalResponses === 0 && group.charts.length > 0}
                    <span
                      class="px-2 py-1 bg-gray-100 text-gray-600 text-xs font-medium rounded">
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
                  {#if lazyLoading && loadingGroups.has(group.product_id)}
                    <!-- Loading state for this group -->
                    <div class="flex items-center justify-center py-12">
                      <div class="flex items-center gap-3 text-gray-600">
                        <Loader2 class="h-5 w-5 animate-spin" />
                        <span class="text-sm font-medium"
                          >Loading charts for {group.product_name}...</span>
                      </div>
                    </div>
                  {:else}
                    {@const chartsToRender = lazyLoading
                      ? groupChartData.get(group.product_id) || []
                      : group.charts}
                    {#if chartsToRender.length > 0}
                      <div class="grid grid-cols-1 lg:grid-cols-2 gap-8">
                        {#each chartsToRender as chart (chart.question_id)}
                          <div
                            class="p-6 border border-gray-200 rounded-lg bg-white">
                            <ChartContent {chart} />
                          </div>
                        {/each}
                      </div>
                    {:else}
                      <div class="py-8">
                        <NoDataAvailable
                          title="No Charts Available"
                          description="No analytics data found for {group.product_name}."
                          icon={BarChart3}
                          variant="inline" />
                      </div>
                    {/if}
                  {/if}
                </div>
              {/if}
            </Card>
          {/each}
        </div>
      {:else}
        <!-- No Results -->
        <NoDataAvailable
          title="No Results Found"
          description="No products or questions match your search criteria."
          icon={Search} />
      {/if}
    {:else}
      <!-- All Charts View -->
      {#if filteredCharts().length > 0}
        <ChartDataWidget
          chartData={{
            ...chartData,
            charts: filteredCharts(),
          }}
          title="" />
      {:else}
        <!-- No Results -->
        <NoDataAvailable
          title="No Results Found"
          description="No charts match your search criteria."
          icon={Search} />
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
