<script lang="ts">
  import { onMount } from 'svelte';
  import { Card, Button } from '$lib/components/ui';
  import { getApiClient, handleApiError } from '$lib/api/client';
  import { auth } from '$lib/stores/auth';
  import { goto } from '$app/navigation';

  interface DashboardStats {
    totalFeedback: number;
    averageRating: number;
    feedbackToday: number;
    activeQRCodes: number;
    topRatedDish?: {
      name: string;
      rating: number;
    };
    recentFeedbackCount: number;
  }

  interface RecentFeedback {
    id: string;
    customer_email?: string;
    rating: number;
    comment?: string;
    dish_name?: string;
    restaurant_name?: string;
    created_at: string;
  }

  let loading = true;
  let error = '';
  let stats: DashboardStats = {
    totalFeedback: 0,
    averageRating: 0,
    feedbackToday: 0,
    activeQRCodes: 0,
    recentFeedbackCount: 0
  };
  let recentFeedback: RecentFeedback[] = [];

  $: authState = $auth;

  onMount(async () => {
    // Check if user is authenticated
    if (!authState.isAuthenticated) {
      goto('/login');
      return;
    }

    await loadDashboardData();
  });

  async function loadDashboardData() {
    loading = true;
    error = '';

    try {
      const api = getApiClient();
      
      // Get all restaurants for the account
      const restaurantsResponse = await api.api.v1RestaurantsList();
      
      if (restaurantsResponse.data.success && restaurantsResponse.data.data) {
        const restaurants = restaurantsResponse.data.data;
        
        // For now, using the first restaurant for demo purposes
        // In a real app, you'd either aggregate across all restaurants or let user select
        if (restaurants.length > 0) {
          const firstRestaurant = restaurants[0];
          
          // Get QR codes and analytics for the first restaurant
          const [qrCodesResponse, analyticsResponse] = await Promise.all([
            api.api.v1RestaurantsQrCodesList(firstRestaurant.id!),
            api.api.v1AnalyticsRestaurantsDetail(firstRestaurant.id!)
          ]);
          
          // Calculate stats
          const activeQRCodes = qrCodesResponse.data.success && qrCodesResponse.data.data 
            ? qrCodesResponse.data.data.filter(qr => qr.is_active).length 
            : 0;
          
          // Parse analytics data
          const analyticsData = analyticsResponse.data;
          const totalFeedback = analyticsData?.total_feedback || 0;
          const averageRating = analyticsData?.average_rating || 0;
          const feedbackToday = analyticsData?.feedback_today || 0;
          const topDish = analyticsData?.top_dishes?.[0];
          const recentFeedbackData = analyticsData?.recent_feedback || [];
          
          stats = {
            totalFeedback,
            averageRating,
            feedbackToday,
            activeQRCodes,
            topRatedDish: topDish ? {
              name: topDish.name,
              rating: topDish.average_rating
            } : undefined,
            recentFeedbackCount: recentFeedbackData.length
          };
          
          // Map recent feedback
          recentFeedback = recentFeedbackData.slice(0, 5).map((fb: any) => ({
            id: fb.id,
            customer_email: fb.customer_email,
            rating: fb.rating,
            comment: fb.comment,
            dish_name: fb.dish_name,
            restaurant_name: firstRestaurant.name,
            created_at: fb.created_at
          }));
        } else {
          // No restaurants yet
          stats = {
            totalFeedback: 0,
            averageRating: 0,
            feedbackToday: 0,
            activeQRCodes: 0,
            recentFeedbackCount: 0
          };
          recentFeedback = [];
        }
      } else {
        // Fallback to zero stats
        stats = {
          totalFeedback: 0,
          averageRating: 0,
          feedbackToday: 0,
          activeQRCodes: 0,
          recentFeedbackCount: 0
        };
        recentFeedback = [];
      }

    } catch (err) {
      error = handleApiError(err);
    } finally {
      loading = false;
    }
  }

  function handleLogout() {
    auth.logout();
    goto('/login');
  }

  function formatDate(dateString: string): string {
    try {
      const date = new Date(dateString);
      return date.toLocaleDateString() + ' ' + date.toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' });
    } catch {
      return dateString;
    }
  }

  function renderStars(rating: number): string {
    return '★'.repeat(rating) + '☆'.repeat(5 - rating);
  }

  function getRatingColor(rating: number): string {
    if (rating >= 4) return 'text-green-600';
    if (rating >= 3) return 'text-yellow-600';
    return 'text-red-600';
  }
</script>

<svelte:head>
  <title>Dashboard - LeCritique</title>
  <meta name="description" content="LeCritique restaurant management dashboard" />
</svelte:head>

<div class="min-h-screen bg-gray-50">
  <!-- Header -->
  <div class="bg-white shadow-sm border-b">
    <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
      <div class="flex justify-between items-center py-6">
        <div>
          <h1 class="text-2xl font-bold text-gray-900">Dashboard</h1>
          <p class="text-gray-600">Welcome back! Here's an overview of your restaurant's feedback.</p>
        </div>
        <div class="flex space-x-3">
          <Button variant="outline" on:click={() => goto('/restaurants')}>
            <svg class="h-4 w-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4" />
            </svg>
            Manage Restaurants
          </Button>
          <Button on:click={() => goto('/analytics')}>
            View Analytics
          </Button>
          <Button variant="outline" on:click={handleLogout}>
            Logout
          </Button>
        </div>
      </div>
    </div>
  </div>

  <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
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
          <h3 class="text-lg font-medium text-gray-900 mb-2">Failed to load dashboard</h3>
          <p class="text-gray-600 mb-4">{error}</p>
          <Button on:click={loadDashboardData}>Try Again</Button>
        </div>
      </Card>

    {:else}
      <!-- Stats Cards -->
      <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6 mb-8">
        <!-- Total Feedback -->
        <Card>
          <div class="flex items-center">
            <div class="flex-shrink-0">
              <div class="h-8 w-8 bg-blue-500 rounded-lg flex items-center justify-center">
                <svg class="h-5 w-5 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 12h.01M12 12h.01M16 12h.01M21 12c0 4.418-4.03 8-9 8a9.863 9.863 0 01-4.255-.949L3 20l1.395-3.72C3.512 15.042 3 13.574 3 12c0-4.418 4.03-8 9-8s9 3.582 9 8z" />
                </svg>
              </div>
            </div>
            <div class="ml-4">
              <p class="text-sm font-medium text-gray-600">Total Feedback</p>
              <p class="text-2xl font-semibold text-gray-900">{stats.totalFeedback}</p>
            </div>
          </div>
        </Card>

        <!-- Average Rating -->
        <Card>
          <div class="flex items-center">
            <div class="flex-shrink-0">
              <div class="h-8 w-8 bg-yellow-500 rounded-lg flex items-center justify-center">
                <svg class="h-5 w-5 text-white" fill="currentColor" viewBox="0 0 24 24">
                  <path d="M11.049 2.927c.3-.921 1.603-.921 1.902 0l1.519 4.674a1 1 0 00.95.69h4.915c.969 0 1.371 1.24.588 1.81l-3.976 2.888a1 1 0 00-.363 1.118l1.518 4.674c.3.922-.755 1.688-1.538 1.118l-3.976-2.888a1 1 0 00-1.176 0l-3.976 2.888c-.783.57-1.838-.197-1.538-1.118l1.518-4.674a1 1 0 00-.363-1.118l-3.976-2.888c-.784-.57-.38-1.81.588-1.81h4.914a1 1 0 00.951-.69l1.519-4.674z" />
                </svg>
              </div>
            </div>
            <div class="ml-4">
              <p class="text-sm font-medium text-gray-600">Average Rating</p>
              <p class="text-2xl font-semibold text-gray-900">{stats.averageRating.toFixed(1)}</p>
            </div>
          </div>
        </Card>

        <!-- Today's Feedback -->
        <Card>
          <div class="flex items-center">
            <div class="flex-shrink-0">
              <div class="h-8 w-8 bg-green-500 rounded-lg flex items-center justify-center">
                <svg class="h-5 w-5 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 7h8m0 0v8m0-8l-8 8-4-4-6 6" />
                </svg>
              </div>
            </div>
            <div class="ml-4">
              <p class="text-sm font-medium text-gray-600">Today's Feedback</p>
              <p class="text-2xl font-semibold text-gray-900">{stats.feedbackToday}</p>
            </div>
          </div>
        </Card>

        <!-- Active QR Codes -->
        <Card>
          <div class="flex items-center">
            <div class="flex-shrink-0">
              <div class="h-8 w-8 bg-purple-500 rounded-lg flex items-center justify-center">
                <svg class="h-5 w-5 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v1m6 11h2m-6 0h-2v4m0-11v3m0 0h.01M12 12h4.01M16 20h4M4 12h4m12 0h.01M5 8h2a1 1 0 001-1V5a1 1 0 00-1-1H5a1 1 0 00-1 1v2a1 1 0 001 1zm12 0h2a1 1 0 001-1V5a1 1 0 00-1-1h-2a1 1 0 00-1 1v2a1 1 0 001 1zM5 20h2a1 1 0 001-1v-2a1 1 0 00-1-1H5a1 1 0 00-1 1v2a1 1 0 001 1z" />
                </svg>
              </div>
            </div>
            <div class="ml-4">
              <p class="text-sm font-medium text-gray-600">Active QR Codes</p>
              <p class="text-2xl font-semibold text-gray-900">{stats.activeQRCodes}</p>
            </div>
          </div>
        </Card>
      </div>

      <!-- Quick Actions & Recent Feedback -->
      <div class="grid grid-cols-1 lg:grid-cols-3 gap-8">
        <!-- Quick Actions -->
        <div class="lg:col-span-1">
          <Card>
            <div class="mb-4">
              <h3 class="text-lg font-medium text-gray-900">Quick Actions</h3>
              <p class="text-sm text-gray-600">Common tasks and shortcuts</p>
            </div>
            
            <div class="space-y-3">
              <Button variant="outline" class="w-full justify-start" on:click={() => goto('/restaurants')}>
                <svg class="h-4 w-4 mr-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 21V5a2 2 0 00-2-2H7a2 2 0 00-2 2v16m14 0h2m-2 0h-4m-5 0H9m11 0a2 2 0 01-2 2H7a2 2 0 01-2-2m2-4h2.01M7 16h6M7 8h6v4H7V8z" />
                </svg>
                Manage Restaurants
              </Button>
              
              <Button variant="outline" class="w-full justify-start" on:click={() => goto('/feedback/manage')}>
                <svg class="h-4 w-4 mr-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 12h.01M12 12h.01M16 12h.01M21 12c0 4.418-4.03 8-9 8a9.863 9.863 0 01-4.255-.949L3 20l1.395-3.72C3.512 15.042 3 13.574 3 12c0-4.418 4.03-8 9-8s9 3.582 9 8z" />
                </svg>
                Manage Feedback
              </Button>
              
              <Button variant="outline" class="w-full justify-start" on:click={() => goto('/analytics')}>
                <svg class="h-4 w-4 mr-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v14a2 2 0 01-2 2h-2a2 2 0 01-2-2z" />
                </svg>
                View Analytics
              </Button>
            </div>

            {#if stats.topRatedDish}
              <div class="mt-6 p-4 bg-green-50 rounded-lg border border-green-200">
                <div class="flex items-center">
                  <svg class="h-5 w-5 text-green-500 mr-2" fill="currentColor" viewBox="0 0 24 24">
                    <path d="M11.049 2.927c.3-.921 1.603-.921 1.902 0l1.519 4.674a1 1 0 00.95.69h4.915c.969 0 1.371 1.24.588 1.81l-3.976 2.888a1 1 0 00-.363 1.118l1.518 4.674c.3.922-.755 1.688-1.538 1.118l-3.976-2.888a1 1 0 00-1.176 0l-3.976 2.888c-.783.57-1.838-.197-1.538-1.118l1.518-4.674a1 1 0 00-.363-1.118l-3.976-2.888c-.784-.57-.38-1.81.588-1.81h4.914a1 1 0 00.951-.69l1.519-4.674z" />
                  </svg>
                  <div>
                    <p class="text-sm font-medium text-green-800">Top Rated Dish</p>
                    <p class="text-sm text-green-700">{stats.topRatedDish.name}</p>
                    <p class="text-xs text-green-600">{renderStars(Math.round(stats.topRatedDish.rating))}</p>
                  </div>
                </div>
              </div>
            {/if}
          </Card>
        </div>

        <!-- Recent Feedback -->
        <div class="lg:col-span-2">
          <Card>
            <div class="mb-4 flex justify-between items-center">
              <div>
                <h3 class="text-lg font-medium text-gray-900">Recent Feedback</h3>
                <p class="text-sm text-gray-600">Latest customer reviews and ratings</p>
              </div>
              <Button variant="outline" size="sm" on:click={() => goto('/feedback/manage')}>
                View All
              </Button>
            </div>

            <div class="space-y-4">
              {#each recentFeedback as feedback}
                <div class="border-l-4 border-l-blue-500 pl-4 py-2">
                  <div class="flex items-start justify-between">
                    <div class="flex-1">
                      <div class="flex items-center space-x-2 mb-1">
                        <span class="text-lg {getRatingColor(feedback.rating)}">{renderStars(feedback.rating)}</span>
                        <span class="text-sm text-gray-500">{feedback.rating}/5</span>
                        {#if feedback.dish_name}
                          <span class="text-sm text-gray-400">•</span>
                          <span class="text-sm font-medium text-gray-700">{feedback.dish_name}</span>
                        {/if}
                      </div>
                      
                      {#if feedback.comment}
                        <p class="text-gray-600 text-sm mb-2">"{feedback.comment}"</p>
                      {/if}
                      
                      <div class="flex items-center text-xs text-gray-500">
                        <span>{formatDate(feedback.created_at)}</span>
                        {#if feedback.customer_email}
                          <span class="ml-2">• {feedback.customer_email}</span>
                        {:else}
                          <span class="ml-2">• Anonymous</span>
                        {/if}
                      </div>
                    </div>
                  </div>
                </div>
              {/each}

              {#if recentFeedback.length === 0}
                <div class="text-center py-8">
                  <svg class="h-12 w-12 text-gray-400 mx-auto mb-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 12h.01M12 12h.01M16 12h.01M21 12c0 4.418-4.03 8-9 8a9.863 9.863 0 01-4.255-.949L3 20l1.395-3.72C3.512 15.042 3 13.574 3 12c0-4.418 4.03-8 9-8s9 3.582 9 8z" />
                  </svg>
                  <p class="text-gray-500">No feedback yet</p>
                  <p class="text-sm text-gray-400">Start collecting feedback by sharing your QR codes!</p>
                </div>
              {/if}
            </div>
          </Card>
        </div>
      </div>
    {/if}
  </div>
</div>