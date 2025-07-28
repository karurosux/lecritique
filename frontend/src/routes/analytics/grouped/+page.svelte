<script lang="ts">
  import { getApiClient, handleApiError } from '$lib/api/client';
  import { auth } from '$lib/stores/auth';
  import { goto } from '$app/navigation';
  import { Card, Button, Select } from '$lib/components/ui';
  import ChartDataWidgetGrouped from '$lib/components/analytics/ChartDataWidgetGrouped.svelte';
  import { BarChart3, Activity, Download, Package, Calendar } from 'lucide-svelte';

  let loading = $state(true);
  let error = $state('');
  let organizations = $state<any[]>([]);
  let selectedOrganization = $state('');
  let chartData = $state<any>(null);
  let hasInitialized = $state(false);
  
  // Filter states
  let filters = $state({
    days: 'all', // 'all', '7', '30', '90'
    productId: ''
  });
  
  let authState = $derived($auth);
  let availableProducts = $state<any[]>([]);

  $effect(() => {
    if (!authState.isAuthenticated) {
      goto('/login');
      return;
    }
    
    if (authState.isAuthenticated && !hasInitialized) {
      hasInitialized = true;
      loadOrganizations();
    }
  });

  async function loadOrganizations() {
    try {
      const api = getApiClient();
      const response = await api.api.v1OrganizationsList();
      
      if (response.data.success && response.data.data) {
        organizations = response.data.data;
        if (organizations.length > 0) {
          selectedOrganization = organizations[0].id;
          await loadProducts();
          loadAnalytics();
        }
      }
    } catch (err) {
      error = handleApiError(err);
    }
  }

  async function loadProducts() {
    if (!selectedOrganization) return;
    
    try {
      const api = getApiClient();
      const response = await api.api.v1OrganizationsProductsList(selectedOrganization);
      
      if (response.data.success && response.data.data) {
        availableProducts = response.data.data;
      }
    } catch (err) {
      console.error('Error loading products:', err);
    }
  }

  async function loadAnalytics() {
    if (!selectedOrganization) return;

    loading = true;
    error = '';

    try {
      const api = getApiClient();
      
      // Build filter parameters
      const chartParams = {};
      
      // Date filter
      if (filters.days !== 'all') {
        const today = new Date();
        const daysAgo = new Date(today.getTime() - parseInt(filters.days) * 24 * 60 * 60 * 1000);
        chartParams.date_from = daysAgo.toISOString().split('T')[0];
      }
      
      // Product filter
      if (filters.productId) {
        chartParams.product_id = filters.productId;
      }
      
      // Load chart data
      const chartResponse = await api.api.v1AnalyticsOrganizationsChartsList(selectedOrganization, chartParams);
      
      // Process chart data
      if (chartResponse.data?.data) {
        chartData = chartResponse.data.data;
      } else {
        chartData = null;
      }

    } catch (err) {
      console.error('Error loading analytics:', err);
      error = handleApiError(err);
    } finally {
      loading = false;
    }
  }

  function handleOrganizationChange() {
    resetFilters();
    loadProducts();
    loadAnalytics();
  }

  function handleExportReport() {
    // TODO: Implement report export functionality
  }

  function resetFilters() {
    filters.days = 'all';
    filters.productId = '';
  }

  function applyFilters() {
    loadAnalytics();
  }
</script>

<svelte:head>
  <title>Grouped Analytics Demo - Kyooar</title>
</svelte:head>

<div class="analytics-grouped-demo max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
  <!-- Demo Header -->
  <div class="mb-8">
    <div class="flex flex-col lg:flex-row lg:justify-between lg:items-center gap-6">
      <div class="space-y-3">
        <div class="flex items-center space-x-3">
          <div class="h-12 w-12 bg-gradient-to-br from-purple-500 to-pink-600 rounded-2xl flex items-center justify-center shadow-lg shadow-purple-500/25">
            <Package class="h-6 w-6 text-white" />
          </div>
          <div>
            <h1 class="text-3xl font-bold bg-gradient-to-r from-gray-900 to-gray-700 bg-clip-text text-transparent">
              Grouped Analytics Demo
            </h1>
            <p class="text-gray-600 font-medium">Charts grouped by product with collapsible sections</p>
          </div>
        </div>
      </div>
      
      <div class="flex items-center space-x-3">
        <Button 
          variant="outline"
          href="/analytics"
          class="group"
        >
          Back to Original
        </Button>
        <Button 
          variant="gradient" 
          size="lg" 
          class="group relative overflow-hidden shadow-lg hover:shadow-xl transition-all duration-300" 
          onclick={handleExportReport}
          disabled={!chartData}
        >
          <div class="absolute inset-0 bg-gradient-to-r from-purple-600 to-pink-600 opacity-0 group-hover:opacity-100 transition-opacity duration-300"></div>
          <Download class="h-5 w-5 mr-2 relative z-10 group-hover:scale-110 transition-transform duration-200" />
          <span class="relative z-10">Export Report</span>
        </Button>
      </div>
    </div>
  </div>

  <!-- Organization Selection -->
  <div class="mb-8">
    <Card variant="elevated">
      <div class="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4">
        <div class="flex items-center space-x-4">
          <div class="h-10 w-10 bg-gradient-to-br from-indigo-500 to-purple-600 rounded-xl flex items-center justify-center">
            <Activity class="h-5 w-5 text-white" />
          </div>
          <div>
            <p class="text-sm font-medium text-gray-600">Organization</p>
            <p class="text-lg font-semibold text-gray-900">
              {organizations.find(r => r.id === selectedOrganization)?.name || 'Select a organization'}
            </p>
            {#if chartData?.summary?.total_responses}
              <p class="text-sm text-gray-500 mt-1">{chartData.summary.total_responses} total responses</p>
            {/if}
          </div>
        </div>
        <div class="flex items-center gap-4">
          {#if organizations.length > 0}
            <Select
              bind:value={selectedOrganization}
              options={organizations.map(r => ({ value: r.id, label: r.name }))}
              onchange={handleOrganizationChange}
            />
          {/if}
          <Button 
            variant="outline" 
            onclick={loadAnalytics} 
            disabled={loading}
            class="group"
          >
            <svg class="h-4 w-4 mr-2 {loading ? 'animate-spin' : 'group-hover:rotate-180 transition-transform duration-300'}" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15" />
            </svg>
            {loading ? 'Loading...' : 'Refresh'}
          </Button>
        </div>
      </div>
    </Card>
  </div>

  <!-- Filters -->
  <div class="mb-8">
    <Card variant="elevated">
      <div class="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4">
        <!-- Time Period Filter -->
        <div class="flex items-center space-x-4">
          <div class="h-10 w-10 bg-gradient-to-br from-blue-500 to-purple-600 rounded-xl flex items-center justify-center">
            <Calendar class="h-5 w-5 text-white" />
          </div>
          <div>
            <p class="text-sm font-medium text-gray-600">Time Period</p>
            <div class="flex items-center gap-2 mt-1">
              <Button 
                variant={filters.days === 'all' ? 'default' : 'outline'} 
                size="sm" 
                onclick={() => { filters.days = 'all'; applyFilters(); }}
              >
                All Time
              </Button>
              <Button 
                variant={filters.days === '7' ? 'default' : 'outline'} 
                size="sm" 
                onclick={() => { filters.days = '7'; applyFilters(); }}
              >
                7 Days
              </Button>
              <Button 
                variant={filters.days === '30' ? 'default' : 'outline'} 
                size="sm" 
                onclick={() => { filters.days = '30'; applyFilters(); }}
              >
                30 Days
              </Button>
              <Button 
                variant={filters.days === '90' ? 'default' : 'outline'} 
                size="sm" 
                onclick={() => { filters.days = '90'; applyFilters(); }}
              >
                90 Days
              </Button>
            </div>
          </div>
        </div>

        <!-- Product Filter -->
        <div class="flex items-center space-x-4">
          <div class="h-10 w-10 bg-gradient-to-br from-purple-500 to-pink-600 rounded-xl flex items-center justify-center">
            <Package class="h-5 w-5 text-white" />
          </div>
          <Select
            bind:value={filters.productId}
            options={[
              { value: '', label: 'All Products' },
              ...availableProducts.map(d => ({ value: d.id, label: d.name }))
            ]}
            onchange={applyFilters}
            minWidth="min-w-48"
            placeholder="Select a product..."
          />
        </div>
      </div>
    </Card>
  </div>

  <!-- Demo Notice -->
  <div class="mb-8">
    <Card variant="gradient" class="bg-gradient-to-r from-purple-50 to-pink-50 border-purple-200">
      <div class="flex items-start gap-3">
        <div class="h-8 w-8 bg-white rounded-lg flex items-center justify-center shadow-sm">
          <BarChart3 class="h-4 w-4 text-purple-600" />
        </div>
        <div>
          <h3 class="text-lg font-semibold text-purple-900 mb-1">New Grouped View</h3>
          <p class="text-purple-700">
            This demo showcases the new grouped analytics view with:
          </p>
          <ul class="mt-2 space-y-1 text-sm text-purple-700">
            <li>• Charts grouped by product with collapsible sections</li>
            <li>• Summary statistics for each product</li>
            <li>• Search and filter functionality</li>
            <li>• Toggle between grouped and all charts view</li>
            <li>• Expand/collapse all controls</li>
          </ul>
        </div>
      </div>
    </Card>
  </div>

  {#if loading}
    <!-- Loading State -->
    <div class="space-y-8">
      <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
        {#each Array(4) as _}
          <Card variant="gradient">
            <div class="animate-pulse">
              <div class="flex items-center justify-between mb-4">
                <div class="h-4 bg-gray-200 rounded w-1/2"></div>
                <div class="h-10 w-10 bg-gray-200 rounded-xl"></div>
              </div>
              <div class="h-8 bg-gray-200 rounded w-3/4 mb-2"></div>
              <div class="h-4 bg-gray-200 rounded w-1/3"></div>
            </div>
          </Card>
        {/each}
      </div>
    </div>

  {:else if error}
    <!-- Error State -->
    <Card variant="elevated">
      <div class="text-center py-12">
        <div class="h-16 w-16 bg-red-100 rounded-2xl flex items-center justify-center mx-auto mb-4">
          <svg class="h-8 w-8 text-red-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-.833-1.964-.833-2.732 0L3.732 16.5c-.77.833.192 3 1.732 3z" />
          </svg>
        </div>
        <h3 class="text-lg font-semibold text-gray-900 mb-2">Failed to load analytics</h3>
        <p class="text-gray-600 mb-6">{error}</p>
        <Button onclick={loadAnalytics} variant="gradient" size="lg">
          Try Again
        </Button>
      </div>
    </Card>

  {:else}
    <!-- Grouped Chart Analytics -->
    <ChartDataWidgetGrouped 
      chartData={chartData} 
      title="" 
      groupByProduct={true}
      initiallyExpanded={false}
    />
  {/if}
</div>
