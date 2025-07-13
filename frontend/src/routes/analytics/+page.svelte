<script lang="ts">
  import { getApiClient, handleApiError } from '$lib/api/client';
  import { auth } from '$lib/stores/auth';
  import { goto } from '$app/navigation';
  import { Card, Button, Select } from '$lib/components/ui';
  import FeedbackChartWidget from '$lib/components/analytics/FeedbackChartWidget.svelte';
  import { BarChart3, Activity, Download, QrCode, Utensils, Calendar } from 'lucide-svelte';

  let loading = $state(true);
  let error = $state('');
  let restaurants = $state<any[]>([]);
  let selectedRestaurant = $state('');
  let feedbacks = $state<any[]>([]);
  let analyticsData = $state<any>(null);
  let hasInitialized = $state(false);
  
  // Filter states
  let filters = $state({
    days: 'all', // 'all', '7', '30', '90'
    dishId: ''
  });
  
  let authState = $derived($auth);
  let availableDishes = $state<any[]>([]);

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
          await loadDishes();
          loadAnalytics();
        }
      }
    } catch (err) {
      error = handleApiError(err);
    }
  }

  async function loadDishes() {
    if (!selectedRestaurant) return;
    
    try {
      const api = getApiClient();
      const response = await api.api.v1RestaurantsDishesList(selectedRestaurant);
      
      if (response.data.success && response.data.data) {
        availableDishes = response.data.data;
      }
    } catch (err) {
      // Error loading dishes is not critical
      console.error('Error loading dishes:', err);
    }
  }

  async function loadAnalytics() {
    if (!selectedRestaurant) return;

    loading = true;
    error = '';

    try {
      const api = getApiClient();
      
      // Build filter parameters
      const filterParams = {};
      
      // Date filter
      if (filters.days !== 'all') {
        const today = new Date();
        const daysAgo = new Date(today.getTime() - parseInt(filters.days) * 24 * 60 * 60 * 1000);
        filterParams.date_from = daysAgo.toISOString().split('T')[0];
      }
      
      // Dish filter
      if (filters.dishId) {
        filterParams.dish_id = filters.dishId;
      }
      
      // Load analytics data and feedbacks in parallel
      const [analyticsResponse, feedbackResponse] = await Promise.all([
        api.api.v1AnalyticsRestaurantsDetail(selectedRestaurant),
        api.api.v1RestaurantsFeedbackList(selectedRestaurant, { 
          limit: 100,
          ...filterParams
        })
      ]);
      
      // Process analytics data
      if (analyticsResponse.data?.data) {
        analyticsData = analyticsResponse.data.data;
      }
      
      // Process feedback data
      if (feedbackResponse.data?.data) {
        feedbacks = feedbackResponse.data.data;
      } else if (Array.isArray(feedbackResponse.data)) {
        feedbacks = feedbackResponse.data;
      } else {
        feedbacks = [];
      }

    } catch (err) {
      console.error('Error loading analytics:', err);
      error = handleApiError(err);
    } finally {
      loading = false;
    }
  }

  function handleRestaurantChange() {
    resetFilters();
    loadDishes();
    loadAnalytics();
  }

  function handleExportReport() {
    // TODO: Implement export functionality
    console.log('Export report for restaurant:', selectedRestaurant);
  }

  function resetFilters() {
    filters.days = 'all';
    filters.dishId = '';
  }

  function applyFilters() {
    loadAnalytics();
  }

</script>

<svelte:head>
  <title>Dish Analytics - LeCritique</title>
</svelte:head>

<div class="analytics-page max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
  <!-- Analytics Header -->
  <div class="mb-8">
    <div class="flex flex-col lg:flex-row lg:justify-between lg:items-center gap-6">
      <div class="space-y-3">
        <div class="flex items-center space-x-3">
          <div class="h-12 w-12 bg-gradient-to-br from-purple-500 to-pink-600 rounded-2xl flex items-center justify-center shadow-lg shadow-purple-500/25">
            <Utensils class="h-6 w-6 text-white" />
          </div>
          <div>
            <h1 class="text-3xl font-bold bg-gradient-to-r from-gray-900 to-gray-700 bg-clip-text text-transparent">
              Dish Performance Analytics
            </h1>
            <p class="text-gray-600 font-medium">Detailed insights on how your dishes are performing</p>
          </div>
        </div>
      </div>
      
      <div class="flex items-center space-x-3">
        <Button 
          variant="gradient" 
          size="lg" 
          class="group relative overflow-hidden shadow-lg hover:shadow-xl transition-all duration-300" 
          onclick={handleExportReport}
          disabled={!analyticsData}
        >
          <div class="absolute inset-0 bg-gradient-to-r from-purple-600 to-pink-600 opacity-0 group-hover:opacity-100 transition-opacity duration-300"></div>
          <Download class="h-5 w-5 mr-2 relative z-10 group-hover:scale-110 transition-transform duration-200" />
          <span class="relative z-10">Export Report</span>
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
            {#if feedbacks.length > 0}
              <p class="text-sm text-gray-500 mt-1">{feedbacks.length} total responses across all dishes</p>
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

  <!-- Simple Filters -->
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

        <!-- Dish Filter -->
        <div class="flex items-center space-x-4">
          <div class="h-10 w-10 bg-gradient-to-br from-purple-500 to-pink-600 rounded-xl flex items-center justify-center">
            <Utensils class="h-5 w-5 text-white" />
          </div>
          <Select
            bind:value={filters.dishId}
            options={[
              { value: '', label: 'All Dishes' },
              ...availableDishes.map(d => ({ value: d.id, label: d.name }))
            ]}
            onchange={applyFilters}
            minWidth="min-w-48"
            placeholder="Select a dish..."
          />
        </div>
      </div>
    </Card>
  </div>

  {#if loading}
    <!-- Loading State -->
    <div class="space-y-8">
      <!-- Summary Cards Loading -->
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

      <!-- Charts Loading -->
      <div class="grid grid-cols-1 lg:grid-cols-2 gap-6">
        {#each Array(2) as _}
          <Card variant="elevated">
            <div class="animate-pulse">
              <div class="h-6 bg-gray-200 rounded w-1/3 mb-4"></div>
              <div class="h-64 bg-gray-100 rounded"></div>
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
          <svg class="h-4 w-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15" />
          </svg>
          Try Again
        </Button>
      </div>
    </Card>

  {:else if restaurants.length === 0}
    <!-- No Restaurants State -->
    <Card variant="elevated">
      <div class="text-center py-16">
        <div class="h-20 w-20 bg-gray-100 rounded-3xl flex items-center justify-center mx-auto mb-6">
          <svg class="h-10 w-10 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 21V5a2 2 0 00-2-2H7a2 2 0 00-2 2v16m14 0h2m-2 0h-5m-9 0H3m2 0h5M9 7h1m-1 4h1m4-4h1m-1 4h1m-5 10v-5a1 1 0 011-1h2a1 1 0 011 1v5m-4 0h4" />
          </svg>
        </div>
        <h3 class="text-xl font-semibold text-gray-900 mb-2">No Restaurants Yet</h3>
        <p class="text-gray-600 mb-6 max-w-md mx-auto">
          Create your first restaurant to start collecting feedback and viewing analytics.
        </p>
        <Button href="/restaurants" variant="gradient" size="lg">
          <svg class="h-5 w-5 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 6v6m0 0v6m0-6h6m-6 0H6" />
          </svg>
          Create Restaurant
        </Button>
      </div>
    </Card>

  {:else}
    <!-- Analytics Content -->
    <div class="space-y-8">
      <!-- Dish Analytics Charts -->
      {#if feedbacks.length > 0}
        <FeedbackChartWidget feedbacks={feedbacks} title="" />
      {:else}
        <!-- No Data State -->
        <Card variant="elevated">
          <div class="text-center py-16">
            <div class="h-20 w-20 bg-gray-100 rounded-3xl flex items-center justify-center mx-auto mb-6">
              <ChartBar class="h-10 w-10 text-gray-400" />
            </div>
            <h3 class="text-xl font-semibold text-gray-900 mb-2">No Analytics Data Available</h3>
            <p class="text-gray-600 mb-6 max-w-md mx-auto">
              Start collecting customer feedback to see detailed analytics and insights about your restaurant performance.
            </p>
            <div class="flex items-center justify-center gap-4">
              <Button href={`/restaurants/${selectedRestaurant}/qr-codes`} variant="gradient" size="lg">
                <QrCode class="h-5 w-5 mr-2" />
                View QR Codes
              </Button>
              <Button onclick={loadAnalytics} variant="outline" size="lg">
                <svg class="h-5 w-5 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 4v5h.582m15.356 2A8.001 8.001 0 004.582 9m0 0H9m11 11v-5h-.581m0 0a8.003 8.003 0 01-15.357-2m15.357 2H15" />
                </svg>
                Check for Data
              </Button>
            </div>
          </div>
        </Card>
      {/if}
    </div>
  {/if}
</div>