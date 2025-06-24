<script lang="ts">
  import { getApiClient, handleApiError } from '$lib/api/client';
  import { auth } from '$lib/stores/auth';
  import { goto } from '$app/navigation';
  import { Card, Button } from '$lib/components/ui';
  import AnalyticsHeader from '$lib/components/analytics/AnalyticsHeader.svelte';
  import AnalyticsFilters from '$lib/components/analytics/AnalyticsFilters.svelte';
  import AnalyticsSummaryCards from '$lib/components/analytics/AnalyticsSummaryCards.svelte';
  import AnalyticsCharts from '$lib/components/analytics/AnalyticsCharts.svelte';
  import AnalyticsRecentFeedback from '$lib/components/analytics/AnalyticsRecentFeedback.svelte';
  import AnalyticsInsights from '$lib/components/analytics/AnalyticsInsights.svelte';

  interface AnalyticsData {
    total_feedback: number;
    average_rating: number;
    feedback_today: number;
    feedback_this_week: number;
    feedback_this_month: number;
    rating_distribution: Record<string, number>;
    top_dishes: Array<{
      id: string;
      name: string;
      average_rating: number;
      feedback_count: number;
    }>;
    recent_feedback: Array<{
      id: string;
      rating: number;
      comment?: string;
      dish_name?: string;
      created_at: string;
    }>;
    trends: {
      daily_feedback: Array<{ date: string; count: number; average_rating: number }>;
      weekly_feedback: Array<{ week: string; count: number; average_rating: number }>;
      monthly_feedback: Array<{ month: string; count: number; average_rating: number }>;
    };
  }

  interface Restaurant {
    id: string;
    name: string;
  }

  let loading = $state(true);
  let error = $state('');
  let selectedRestaurant = $state('');
  let restaurants = $state<Restaurant[]>([]);
  let analyticsData = $state<AnalyticsData | null>(null);
  let selectedTimeframe = $state('7d');
  let hasInitialized = $state(false);

  let authState = $derived($auth);

  // Handle authentication and initial load
  $effect(() => {
    if (!authState.isAuthenticated) {
      goto('/login');
      return;
    }
    
    // Only load data once when authenticated
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
        restaurants = response.data.data.map((r: any) => ({
          id: r.id || '',
          name: r.name || ''
        }));
        if (restaurants.length > 0 && !selectedRestaurant) {
          selectedRestaurant = restaurants[0].id;
          // Load analytics for the first restaurant
          loadAnalytics();
        }
      }
    } catch (err) {
      console.error('Error loading restaurants:', err);
      error = handleApiError(err);
    }
  }

  async function loadAnalytics() {
    if (!selectedRestaurant) return;

    loading = true;
    error = '';

    try {
      const api = getApiClient();
      const response = await api.api.v1AnalyticsRestaurantsDetail(selectedRestaurant);
      
      if (response.data) {
        // Map the response to our analytics structure
        analyticsData = {
          total_feedback: response.data.total_feedback || 0,
          average_rating: response.data.average_rating || 0,
          feedback_today: response.data.feedback_today || 0,
          feedback_this_week: response.data.feedback_this_week || 0,
          feedback_this_month: response.data.feedback_this_month || 0,
          rating_distribution: response.data.rating_distribution || { '1': 0, '2': 0, '3': 0, '4': 0, '5': 0 },
          top_dishes: response.data.top_dishes || [],
          recent_feedback: response.data.recent_feedback || [],
          trends: response.data.trends || {
            daily_feedback: [],
            weekly_feedback: [],
            monthly_feedback: []
          }
        };
      } else {
        // Fallback data structure
        analyticsData = {
          total_feedback: 0,
          average_rating: 0,
          feedback_today: 0,
          feedback_this_week: 0,
          feedback_this_month: 0,
          rating_distribution: { '1': 0, '2': 0, '3': 0, '4': 0, '5': 0 },
          top_dishes: [],
          recent_feedback: [],
          trends: {
            daily_feedback: [],
            weekly_feedback: [],
            monthly_feedback: []
          }
        };
      }

    } catch (err) {
      error = handleApiError(err);
    } finally {
      loading = false;
    }
  }

  function handleRestaurantChange() {
    loadAnalytics();
  }

  function handleTimeframeChange() {
    // In the future, this could reload analytics with different timeframe
    // For now, just keeping the state updated
  }

  function exportReport() {
    if (!analyticsData) return;

    const reportData = {
      restaurant: restaurants.find(r => r.id === selectedRestaurant)?.name,
      generated_at: new Date().toISOString(),
      timeframe: selectedTimeframe,
      summary: {
        total_feedback: analyticsData.total_feedback,
        average_rating: analyticsData.average_rating,
        feedback_today: analyticsData.feedback_today,
        feedback_this_week: analyticsData.feedback_this_week,
        feedback_this_month: analyticsData.feedback_this_month
      },
      rating_distribution: analyticsData.rating_distribution,
      top_dishes: analyticsData.top_dishes
    };

    const blob = new Blob([JSON.stringify(reportData, null, 2)], { type: 'application/json' });
    const link = document.createElement('a');
    const url = URL.createObjectURL(blob);
    link.setAttribute('href', url);
    link.setAttribute('download', `analytics-report-${new Date().toISOString().split('T')[0]}.json`);
    link.style.visibility = 'hidden';
    document.body.appendChild(link);
    link.click();
    document.body.removeChild(link);
  }
</script>

<svelte:head>
  <title>Analytics & Reports - LeCritique</title>
  <meta name="description" content="Advanced analytics and reporting for restaurant feedback" />
</svelte:head>

<div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
  <!-- Analytics Header -->
  <AnalyticsHeader
    {analyticsData}
    {loading}
    onexportreport={exportReport}
  />

  <!-- Filters -->
  <AnalyticsFilters
    {restaurants}
    bind:selectedRestaurant
    bind:selectedTimeframe
    onrestaurantchange={handleRestaurantChange}
    ontimeframechange={handleTimeframeChange}
  />
  {#if error}
    <!-- Error State -->
    <Card>
      <div class="text-center py-12">
        <svg class="h-12 w-12 text-red-500 mx-auto mb-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-.833-1.964-.833-2.732 0L3.732 16.5c-.77.833.192 3 1.732 3z" />
        </svg>
        <h3 class="text-lg font-medium text-gray-900 mb-2">Failed to load analytics</h3>
        <p class="text-gray-600 mb-4">{error}</p>
        <Button onclick={loadAnalytics}>Try Again</Button>
      </div>
    </Card>
  {:else}
    <!-- Summary Cards -->
    <AnalyticsSummaryCards {analyticsData} {loading} />

    <!-- Charts -->
    <AnalyticsCharts {analyticsData} {loading} />

    <!-- Recent Feedback -->
    <AnalyticsRecentFeedback {analyticsData} {loading} />

    <!-- Insights -->
    <AnalyticsInsights {analyticsData} {loading} />
  {/if}
</div>