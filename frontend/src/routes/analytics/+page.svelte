<script lang="ts">
  import { onMount } from 'svelte';
  import { Card, Button } from '$lib/components/ui';
  import { getApiClient, handleApiError } from '$lib/api/client';
  import { auth } from '$lib/stores/auth';
  import { goto } from '$app/navigation';

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

  let loading = true;
  let error = '';
  let selectedRestaurant = '';
  let restaurants: any[] = [];
  let analyticsData: AnalyticsData | null = null;
  let selectedTimeframe = '7d'; // 7d, 30d, 90d, 1y

  $: authState = $auth;

  onMount(async () => {
    if (!authState.isAuthenticated) {
      goto('/login');
      return;
    }

    await loadRestaurants();
    await loadAnalytics();
  });

  async function loadRestaurants() {
    try {
      const api = getApiClient();
      const response = await api.api.v1RestaurantsList();
      
      if (response.data.success && response.data.data) {
        restaurants = response.data.data;
        if (restaurants.length > 0 && !selectedRestaurant) {
          selectedRestaurant = restaurants[0].id;
        }
      }
    } catch (err) {
      console.error('Error loading restaurants:', err);
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

  function getPercentage(value: number, total: number): number {
    return total > 0 ? (value / total) * 100 : 0;
  }

  function formatDate(dateString: string): string {
    try {
      const date = new Date(dateString);
      return date.toLocaleDateString();
    } catch {
      return dateString;
    }
  }

  function renderStars(rating: number): string {
    return '★'.repeat(Math.round(rating)) + '☆'.repeat(5 - Math.round(rating));
  }

  function getRatingColor(rating: number): string {
    if (rating >= 4) return 'text-green-600';
    if (rating >= 3) return 'text-yellow-600';
    return 'text-red-600';
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
    <!-- Page Header -->
    <div class="mb-8">
      <div class="flex justify-between items-center">
        <div>
          <h1 class="text-3xl font-bold bg-gradient-to-r from-gray-900 to-gray-700 bg-clip-text text-transparent">
            Analytics & Reports
          </h1>
          <p class="text-gray-600 font-medium mt-2">Detailed insights and performance metrics for your restaurants</p>
        </div>
        <div class="flex space-x-3">
          <select
            bind:value={selectedRestaurant}
            on:change={handleRestaurantChange}
            class="px-4 py-2 border border-gray-300 rounded-lg bg-white focus:outline-none focus:ring-2 focus:ring-blue-500 cursor-pointer shadow-sm"
          >
            {#each restaurants as restaurant}
              <option value={restaurant.id}>{restaurant.name}</option>
            {/each}
          </select>
          
          <Button variant="glass" size="lg" on:click={exportReport} disabled={!analyticsData}>
            <svg class="h-5 w-5 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 10v6m0 0l-3-3m3 3l3-3m2 8H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
            </svg>
            Export Report
          </Button>
        </div>
      </div>
    </div>
    {#if loading}
      <!-- Loading State -->
      <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6 mb-8">
        {#each Array(4) as _}
          <Card>
            <div class="animate-pulse">
              <div class="h-4 bg-gray-200 rounded w-3/4 mb-2"></div>
              <div class="h-8 bg-gray-200 rounded w-1/2"></div>
            </div>
          </Card>
        {/each}
      </div>

    {:else if error}
      <!-- Error State -->
      <Card>
        <div class="text-center py-12">
          <svg class="h-12 w-12 text-red-500 mx-auto mb-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-.833-1.964-.833-2.732 0L3.732 16.5c-.77.833.192 3 1.732 3z" />
          </svg>
          <h3 class="text-lg font-medium text-gray-900 mb-2">Failed to load analytics</h3>
          <p class="text-gray-600 mb-4">{error}</p>
          <Button on:click={loadAnalytics}>Try Again</Button>
        </div>
      </Card>

    {:else if analyticsData}
      <!-- Analytics Content -->
      
      <!-- Summary Stats -->
      <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-5 gap-6 mb-8">
        <Card>
          <div class="text-center">
            <div class="text-3xl font-bold text-blue-600 mb-1">{analyticsData.total_feedback}</div>
            <div class="text-sm text-gray-600">Total Feedback</div>
          </div>
        </Card>
        
        <Card>
          <div class="text-center">
            <div class="text-3xl font-bold text-yellow-600 mb-1">{analyticsData.average_rating.toFixed(1)}</div>
            <div class="text-sm text-gray-600">Average Rating</div>
            <div class="text-xs text-yellow-600 mt-1">{renderStars(analyticsData.average_rating)}</div>
          </div>
        </Card>
        
        <Card>
          <div class="text-center">
            <div class="text-3xl font-bold text-green-600 mb-1">{analyticsData.feedback_today}</div>
            <div class="text-sm text-gray-600">Today</div>
          </div>
        </Card>
        
        <Card>
          <div class="text-center">
            <div class="text-3xl font-bold text-purple-600 mb-1">{analyticsData.feedback_this_week}</div>
            <div class="text-sm text-gray-600">This Week</div>
          </div>
        </Card>
        
        <Card>
          <div class="text-center">
            <div class="text-3xl font-bold text-indigo-600 mb-1">{analyticsData.feedback_this_month}</div>
            <div class="text-sm text-gray-600">This Month</div>
          </div>
        </Card>
      </div>

      <div class="grid grid-cols-1 lg:grid-cols-2 gap-8 mb-8">
        <!-- Rating Distribution -->
        <Card>
          <div class="mb-4">
            <h3 class="text-lg font-medium text-gray-900">Rating Distribution</h3>
            <p class="text-sm text-gray-600">Breakdown of customer ratings</p>
          </div>
          
          <div class="space-y-3">
            {#each [5, 4, 3, 2, 1] as rating}
              {@const count = analyticsData.rating_distribution[rating.toString()] || 0}
              {@const percentage = getPercentage(count, analyticsData.total_feedback)}
              <div class="flex items-center">
                <div class="w-12 text-sm text-gray-600">{rating} ★</div>
                <div class="flex-1 mx-3">
                  <div class="bg-gray-200 rounded-full h-2">
                    <div
                      class="h-2 rounded-full {rating >= 4 ? 'bg-green-500' : rating >= 3 ? 'bg-yellow-500' : 'bg-red-500'}"
                      style="width: {percentage}%"
                    ></div>
                  </div>
                </div>
                <div class="w-16 text-sm text-gray-600 text-right">
                  {count} ({percentage.toFixed(1)}%)
                </div>
              </div>
            {/each}
          </div>
        </Card>

        <!-- Top Performing Dishes -->
        <Card>
          <div class="mb-4">
            <h3 class="text-lg font-medium text-gray-900">Top Performing Dishes</h3>
            <p class="text-sm text-gray-600">Highest rated dishes with feedback</p>
          </div>
          
          <div class="space-y-4">
            {#each analyticsData.top_dishes.slice(0, 5) as dish, index}
              <div class="flex items-center justify-between border-b border-gray-100 pb-3 last:border-b-0">
                <div class="flex items-center space-x-3">
                  <div class="flex items-center justify-center w-8 h-8 bg-blue-100 rounded-full text-sm font-medium text-blue-600">
                    {index + 1}
                  </div>
                  <div>
                    <div class="font-medium text-gray-900">{dish.name}</div>
                    <div class="text-sm text-gray-500">{dish.feedback_count} reviews</div>
                  </div>
                </div>
                <div class="text-right">
                  <div class="font-medium {getRatingColor(dish.average_rating)}">{dish.average_rating.toFixed(1)}</div>
                  <div class="text-xs text-gray-500">{renderStars(dish.average_rating)}</div>
                </div>
              </div>
            {:else}
              <div class="text-center py-8 text-gray-500">
                <p>No dish data available</p>
              </div>
            {/each}
          </div>
        </Card>
      </div>

      <!-- Recent Feedback Trends -->
      <Card class="mb-8">
        <div class="mb-4">
          <h3 class="text-lg font-medium text-gray-900">Recent Feedback Activity</h3>
          <p class="text-sm text-gray-600">Latest customer feedback and comments</p>
        </div>
        
        <div class="space-y-4 max-h-96 overflow-y-auto">
          {#each analyticsData.recent_feedback.slice(0, 10) as feedback}
            <div class="border-l-4 {feedback.rating >= 4 ? 'border-green-500' : feedback.rating >= 3 ? 'border-yellow-500' : 'border-red-500'} pl-4 py-2">
              <div class="flex items-start justify-between">
                <div class="flex-1">
                  <div class="flex items-center space-x-2 mb-1">
                    <span class="text-sm {getRatingColor(feedback.rating)}">{renderStars(feedback.rating)}</span>
                    <span class="text-xs text-gray-500">{feedback.rating}/5</span>
                    {#if feedback.dish_name}
                      <span class="text-xs text-gray-400">•</span>
                      <span class="text-xs font-medium text-gray-700">{feedback.dish_name}</span>
                    {/if}
                  </div>
                  
                  {#if feedback.comment}
                    <p class="text-gray-600 text-sm mb-1">"{feedback.comment}"</p>
                  {/if}
                  
                  <p class="text-xs text-gray-500">{formatDate(feedback.created_at)}</p>
                </div>
              </div>
            </div>
          {:else}
            <div class="text-center py-8 text-gray-500">
              <p>No recent feedback available</p>
            </div>
          {/each}
        </div>
      </Card>

      <!-- Summary Insights -->
      <Card>
        <div class="mb-4">
          <h3 class="text-lg font-medium text-gray-900">Key Insights</h3>
          <p class="text-sm text-gray-600">Automated analysis and recommendations</p>
        </div>
        
        <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
          {#if analyticsData}
            {@const highRatings = (analyticsData.rating_distribution['4'] || 0) + (analyticsData.rating_distribution['5'] || 0)}
            {@const highRatingPercentage = getPercentage(highRatings, analyticsData.total_feedback)}
            <!-- High Rating Percentage -->
            <div class="bg-green-50 border border-green-200 rounded-lg p-4">
            <div class="flex items-center">
              <svg class="h-5 w-5 text-green-500 mr-2" fill="currentColor" viewBox="0 0 24 24">
                <path d="M11.049 2.927c.3-.921 1.603-.921 1.902 0l1.519 4.674a1 1 0 00.95.69h4.915c.969 0 1.371 1.24.588 1.81l-3.976 2.888a1 1 0 00-.363 1.118l1.518 4.674c.3.922-.755 1.688-1.538 1.118l-3.976-2.888a1 1 0 00-1.176 0l-3.976 2.888c-.783.57-1.838-.197-1.538-1.118l1.518-4.674a1 1 0 00-.363-1.118l-3.976-2.888c-.784-.57-.38-1.81.588-1.81h4.914a1 1 0 00.951-.69l1.519-4.674z" />
              </svg>
              <div>
                <div class="font-medium text-green-800">{highRatingPercentage.toFixed(1)}% High Ratings</div>
                <div class="text-sm text-green-600">4-5 star reviews</div>
              </div>
            </div>
          </div>

          <!-- Customer Satisfaction -->
          <div class="bg-blue-50 border border-blue-200 rounded-lg p-4">
            <div class="flex items-center">
              <svg class="h-5 w-5 text-blue-500 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M14.828 14.828a4 4 0 01-5.656 0M9 10h1.01M15 10h1.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
              </svg>
              <div>
                <div class="font-medium text-blue-800">
                  {analyticsData.average_rating >= 4 ? 'Excellent' : analyticsData.average_rating >= 3 ? 'Good' : 'Needs Improvement'}
                </div>
                <div class="text-sm text-blue-600">Overall satisfaction</div>
              </div>
            </div>
          </div>

          <!-- Feedback Volume -->
          <div class="bg-purple-50 border border-purple-200 rounded-lg p-4">
            <div class="flex items-center">
              <svg class="h-5 w-5 text-purple-500 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 7h8m0 0v8m0-8l-8 8-4-4-6 6" />
              </svg>
              <div>
                <div class="font-medium text-purple-800">
                  {analyticsData.feedback_this_week > analyticsData.feedback_today * 7 / 2 ? 'Growing' : 'Stable'}
                </div>
                <div class="text-sm text-purple-600">Feedback trend</div>
              </div>
            </div>
          </div>
          {/if}
        </div>
      </Card>
    {/if}
</div>