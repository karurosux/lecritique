<script lang="ts">
  import { getApiClient, handleApiError } from '$lib/api/client';
  import { auth } from '$lib/stores/auth';
  import { goto } from '$app/navigation';
  import { Card, Button, Select } from '$lib/components/ui';
  import ChartDataWidget from '$lib/components/analytics/ChartDataWidget.svelte';
  import ChartDataWidgetGrouped from '$lib/components/analytics/ChartDataWidgetGrouped.svelte';
  import { BarChart3, Activity, Layers, Grid3x3 } from 'lucide-svelte';

  let loading = $state(true);
  let error = $state('');
  let restaurants = $state<any[]>([]);
  let selectedRestaurant = $state('');
  let chartData = $state<any>(null);
  let hasInitialized = $state(false);
  let viewMode = $state<'side-by-side' | 'original' | 'grouped'>('side-by-side');
  
  let authState = $derived($auth);

  $effect(() => {
    if (!authState.isAuthenticated) {
      goto('/login');
      return;
    }
    
    if (authState.isAuthenticated && !hasInitialized) {
      hasInitialized = true;
      loadRestaurants();
    }
  });

  async function loadRestaurants() {
    try {
      const api = getApiClient();
      const response = await api.api.v1RestaurantsList();
      
      if (response.data.success && response.data.data) {
        restaurants = response.data.data;
        if (restaurants.length > 0) {
          selectedRestaurant = restaurants[0].id;
          loadAnalytics();
        }
      }
    } catch (err) {
      error = handleApiError(err);
    }
  }

  async function loadAnalytics() {
    if (!selectedRestaurant) return;

    loading = true;
    error = '';

    try {
      const api = getApiClient();
      const chartResponse = await api.api.v1AnalyticsRestaurantsChartsList(selectedRestaurant, {});
      
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

  function handleRestaurantChange() {
    loadAnalytics();
  }
</script>

<svelte:head>
  <title>Analytics Comparison - LeCritique</title>
</svelte:head>

<div class="analytics-comparison max-w-[1600px] mx-auto px-4 sm:px-6 lg:px-8 py-8">
  <!-- Header -->
  <div class="mb-8">
    <div class="flex flex-col lg:flex-row lg:justify-between lg:items-center gap-6">
      <div class="space-y-3">
        <div class="flex items-center space-x-3">
          <div class="h-12 w-12 bg-gradient-to-br from-purple-500 to-pink-600 rounded-2xl flex items-center justify-center shadow-lg shadow-purple-500/25">
            <Layers class="h-6 w-6 text-white" />
          </div>
          <div>
            <h1 class="text-3xl font-bold bg-gradient-to-r from-gray-900 to-gray-700 bg-clip-text text-transparent">
              Analytics View Comparison
            </h1>
            <p class="text-gray-600 font-medium">Compare original and grouped chart views</p>
          </div>
        </div>
      </div>
      
      <div class="flex items-center gap-3">
        <Button
          variant={viewMode === 'side-by-side' ? 'default' : 'outline'}
          onclick={() => viewMode = 'side-by-side'}
          class="group"
        >
          <Grid3x3 class="h-4 w-4 mr-2" />
          Side by Side
        </Button>
        <Button
          variant={viewMode === 'original' ? 'default' : 'outline'}
          onclick={() => viewMode = 'original'}
        >
          Original Only
        </Button>
        <Button
          variant={viewMode === 'grouped' ? 'default' : 'outline'}
          onclick={() => viewMode = 'grouped'}
        >
          Grouped Only
        </Button>
      </div>
    </div>
  </div>

  <!-- Restaurant Selection -->
  <div class="mb-8">
    <Card variant="elevated">
      <div class="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4">
        <div class="flex items-center space-x-4">
          <div class="h-10 w-10 bg-gradient-to-br from-indigo-500 to-purple-600 rounded-xl flex items-center justify-center">
            <Activity class="h-5 w-5 text-white" />
          </div>
          <div>
            <p class="text-sm font-medium text-gray-600">Restaurant</p>
            <p class="text-lg font-semibold text-gray-900">
              {restaurants.find(r => r.id === selectedRestaurant)?.name || 'Select a restaurant'}
            </p>
            {#if chartData?.summary?.total_responses}
              <p class="text-sm text-gray-500 mt-1">{chartData.summary.total_responses} total responses</p>
            {/if}
          </div>
        </div>
        <div class="flex items-center gap-4">
          {#if restaurants.length > 0}
            <Select
              bind:value={selectedRestaurant}
              options={restaurants.map(r => ({ value: r.id, label: r.name }))}
              onchange={handleRestaurantChange}
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

  {#if loading}
    <!-- Loading State -->
    <div class="grid grid-cols-1 lg:grid-cols-2 gap-6">
      {#each Array(2) as _}
        <Card variant="elevated">
          <div class="animate-pulse p-8">
            <div class="h-6 bg-gray-200 rounded w-1/3 mb-4"></div>
            <div class="space-y-4">
              {#each Array(3) as _}
                <div>
                  <div class="h-4 bg-gray-200 rounded w-1/2 mb-2"></div>
                  <div class="h-32 bg-gray-100 rounded"></div>
                </div>
              {/each}
            </div>
          </div>
        </Card>
      {/each}
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
    <!-- Comparison View -->
    {#if viewMode === 'side-by-side'}
      <div class="grid grid-cols-1 xl:grid-cols-2 gap-8">
        <!-- Original View -->
        <div>
          <Card variant="gradient" class="mb-4 bg-gradient-to-r from-blue-50 to-indigo-50 border-blue-200">
            <div class="flex items-center gap-3">
              <div class="h-8 w-8 bg-white rounded-lg flex items-center justify-center shadow-sm">
                <BarChart3 class="h-4 w-4 text-blue-600" />
              </div>
              <div>
                <h3 class="text-lg font-semibold text-blue-900">Original View</h3>
                <p class="text-sm text-blue-700">All charts in a flat grid layout</p>
              </div>
            </div>
          </Card>
          <ChartDataWidget chartData={chartData} title="" />
        </div>

        <!-- Grouped View -->
        <div>
          <Card variant="gradient" class="mb-4 bg-gradient-to-r from-purple-50 to-pink-50 border-purple-200">
            <div class="flex items-center gap-3">
              <div class="h-8 w-8 bg-white rounded-lg flex items-center justify-center shadow-sm">
                <Layers class="h-4 w-4 text-purple-600" />
              </div>
              <div>
                <h3 class="text-lg font-semibold text-purple-900">Grouped View</h3>
                <p class="text-sm text-purple-700">Charts grouped by dish with collapsible sections</p>
              </div>
            </div>
          </Card>
          <ChartDataWidgetGrouped chartData={chartData} title="" groupByDish={true} initiallyExpanded={false} />
        </div>
      </div>
    {:else if viewMode === 'original'}
      <ChartDataWidget chartData={chartData} title="Original Analytics View" />
    {:else}
      <ChartDataWidgetGrouped chartData={chartData} title="Grouped Analytics View" groupByDish={true} initiallyExpanded={false} />
    {/if}
  {/if}
</div>