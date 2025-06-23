<script lang="ts">
  import { onMount } from 'svelte';
  import { Card, Button, UserMenu } from '$lib/components/ui';
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

<div class="min-h-screen bg-gradient-to-br from-slate-50 via-blue-50/50 to-indigo-50">
  <!-- Header -->
  <div class="relative bg-gradient-to-r from-white/95 to-gray-50/95 backdrop-blur-xl border-b border-white/20 shadow-lg shadow-gray-900/5 z-50">
    <div class="absolute inset-0 bg-gradient-to-r from-blue-500/5 to-purple-500/5"></div>
    <div class="relative max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
      <div class="flex justify-between items-center py-8">
        <div class="space-y-2">
          <h1 class="text-3xl font-bold bg-gradient-to-r from-gray-900 to-gray-700 bg-clip-text text-transparent">
            Dashboard
          </h1>
          <p class="text-gray-600 font-medium">Welcome back! Here's an overview of your restaurant's feedback.</p>
        </div>
        <div class="flex items-center space-x-4">
          <div class="flex space-x-3">
            <Button variant="glass" size="lg" on:click={() => goto('/restaurants')}>
              <svg class="h-5 w-5 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 21V5a2 2 0 00-2-2H7a2 2 0 00-2 2v16m14 0h2m-2 0h-4m-5 0H9m11 0a2 2 0 01-2 2H7a2 2 0 01-2-2m2-4h2.01M7 16h6M7 8h6v4H7V8z" />
              </svg>
              Manage Restaurants
            </Button>
            <Button variant="gradient" size="lg" on:click={() => goto('/analytics')}>
              <svg class="h-5 w-5 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v14a2 2 0 01-2 2h-2a2 2 0 01-2-2z" />
              </svg>
              View Analytics
            </Button>
          </div>
          <UserMenu />
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
      <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6 mb-12">
        <!-- Total Feedback -->
        <Card variant="gradient" hover interactive>
          <div class="flex items-center justify-between">
            <div class="space-y-2">
              <p class="text-sm font-semibold text-gray-600 uppercase tracking-wide">Total Feedback</p>
              <p class="text-3xl font-bold bg-gradient-to-r from-blue-600 to-purple-600 bg-clip-text text-transparent">
                {stats.totalFeedback}
              </p>
              <div class="flex items-center space-x-1 text-green-600">
                <svg class="h-4 w-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 7h8m0 0v8m0-8l-8 8-4-4-6 6" />
                </svg>
                <span class="text-sm font-medium">All time</span>
              </div>
            </div>
            <div class="h-16 w-16 bg-gradient-to-br from-blue-500 to-purple-600 rounded-2xl flex items-center justify-center shadow-lg shadow-blue-500/25">
              <svg class="h-8 w-8 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 12h.01M12 12h.01M16 12h.01M21 12c0 4.418-4.03 8-9 8a9.863 9.863 0 01-4.255-.949L3 20l1.395-3.72C3.512 15.042 3 13.574 3 12c0-4.418 4.03-8 9-8s9 3.582 9 8z" />
              </svg>
            </div>
          </div>
        </Card>

        <!-- Average Rating -->
        <Card variant="gradient" hover interactive>
          <div class="flex items-center justify-between">
            <div class="space-y-2">
              <p class="text-sm font-semibold text-gray-600 uppercase tracking-wide">Average Rating</p>
              <p class="text-3xl font-bold bg-gradient-to-r from-yellow-600 to-orange-600 bg-clip-text text-transparent">
                {stats.averageRating.toFixed(1)}
              </p>
              <div class="flex items-center space-x-1">
                <div class="flex text-yellow-400">
                  {#each Array(5) as _, i}
                    <svg class="h-4 w-4 {i < Math.round(stats.averageRating) ? 'text-yellow-400' : 'text-gray-300'}" fill="currentColor" viewBox="0 0 24 24">
                      <path d="M11.049 2.927c.3-.921 1.603-.921 1.902 0l1.519 4.674a1 1 0 00.95.69h4.915c.969 0 1.371 1.24.588 1.81l-3.976 2.888a1 1 0 00-.363 1.118l1.518 4.674c.3.922-.755 1.688-1.538 1.118l-3.976-2.888a1 1 0 00-1.176 0l-3.976 2.888c-.783.57-1.838-.197-1.538-1.118l1.518-4.674a1 1 0 00-.363-1.118l-3.976-2.888c-.784-.57-.38-1.81.588-1.81h4.914a1 1 0 00.951-.69l1.519-4.674z" />
                    </svg>
                  {/each}
                </div>
              </div>
            </div>
            <div class="h-16 w-16 bg-gradient-to-br from-yellow-500 to-orange-500 rounded-2xl flex items-center justify-center shadow-lg shadow-yellow-500/25">
              <svg class="h-8 w-8 text-white" fill="currentColor" viewBox="0 0 24 24">
                <path d="M11.049 2.927c.3-.921 1.603-.921 1.902 0l1.519 4.674a1 1 0 00.95.69h4.915c.969 0 1.371 1.24.588 1.81l-3.976 2.888a1 1 0 00-.363 1.118l1.518 4.674c.3.922-.755 1.688-1.538 1.118l-3.976-2.888a1 1 0 00-1.176 0l-3.976 2.888c-.783.57-1.838-.197-1.538-1.118l1.518-4.674a1 1 0 00-.363-1.118l-3.976-2.888c-.784-.57-.38-1.81.588-1.81h4.914a1 1 0 00.951-.69l1.519-4.674z" />
              </svg>
            </div>
          </div>
        </Card>

        <!-- Today's Feedback -->
        <Card variant="gradient" hover interactive>
          <div class="flex items-center justify-between">
            <div class="space-y-2">
              <p class="text-sm font-semibold text-gray-600 uppercase tracking-wide">Today's Feedback</p>
              <p class="text-3xl font-bold bg-gradient-to-r from-green-600 to-emerald-600 bg-clip-text text-transparent">
                {stats.feedbackToday}
              </p>
              <div class="flex items-center space-x-1 text-green-600">
                <svg class="h-4 w-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 3v1m0 16v1m9-9h-1M4 12H3m15.364 6.364l-.707-.707M6.343 6.343l-.707-.707m12.728 0l-.707.707M6.343 17.657l-.707.707" />
                </svg>
                <span class="text-sm font-medium">Today</span>
              </div>
            </div>
            <div class="h-16 w-16 bg-gradient-to-br from-green-500 to-emerald-500 rounded-2xl flex items-center justify-center shadow-lg shadow-green-500/25">
              <svg class="h-8 w-8 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 7h8m0 0v8m0-8l-8 8-4-4-6 6" />
              </svg>
            </div>
          </div>
        </Card>

        <!-- Active QR Codes -->
        <Card variant="gradient" hover interactive>
          <div class="flex items-center justify-between">
            <div class="space-y-2">
              <p class="text-sm font-semibold text-gray-600 uppercase tracking-wide">Active QR Codes</p>
              <p class="text-3xl font-bold bg-gradient-to-r from-purple-600 to-indigo-600 bg-clip-text text-transparent">
                {stats.activeQRCodes}
              </p>
              <div class="flex items-center space-x-1 text-purple-600">
                <svg class="h-4 w-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
                </svg>
                <span class="text-sm font-medium">Active</span>
              </div>
            </div>
            <div class="h-16 w-16 bg-gradient-to-br from-purple-500 to-indigo-500 rounded-2xl flex items-center justify-center shadow-lg shadow-purple-500/25">
              <svg class="h-8 w-8 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v1m6 11h2m-6 0h-2v4m0-11v3m0 0h.01M12 12h4.01M16 20h4M4 12h4m12 0h.01M5 8h2a1 1 0 001-1V5a1 1 0 00-1-1H5a1 1 0 00-1 1v2a1 1 0 001 1zm12 0h2a1 1 0 001-1V5a1 1 0 00-1-1h-2a1 1 0 00-1 1v2a1 1 0 001 1zM5 20h2a1 1 0 001-1v-2a1 1 0 00-1-1H5a1 1 0 00-1 1v2a1 1 0 001 1z" />
              </svg>
            </div>
          </div>
        </Card>
      </div>

      <!-- Quick Actions & Recent Feedback -->
      <div class="grid grid-cols-1 lg:grid-cols-3 gap-8">
        <!-- Quick Actions -->
        <div class="lg:col-span-1">
          <Card variant="glass" padding={false}>
            <div class="p-6 lg:p-8 space-y-6">
              <div class="space-y-2">
                <h3 class="text-xl font-bold bg-gradient-to-r from-gray-900 to-gray-700 bg-clip-text text-transparent">
                  Quick Actions
                </h3>
                <p class="text-gray-600 font-medium">Common tasks and shortcuts</p>
              </div>
              
              <div class="space-y-4">
                <Button variant="ghost" class="w-full justify-start group hover:bg-gradient-to-r hover:from-blue-50 hover:to-purple-50" on:click={() => goto('/restaurants')}>
                  <div class="h-10 w-10 bg-gradient-to-br from-blue-500 to-purple-600 rounded-xl flex items-center justify-center mr-4 group-hover:scale-110 transition-transform">
                    <svg class="h-5 w-5 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 21V5a2 2 0 00-2-2H7a2 2 0 00-2 2v16m14 0h2m-2 0h-4m-5 0H9m11 0a2 2 0 01-2 2H7a2 2 0 01-2-2m2-4h2.01M7 16h6M7 8h6v4H7V8z" />
                    </svg>
                  </div>
                  <div class="text-left">
                    <p class="font-semibold text-gray-900">Manage Restaurants</p>
                    <p class="text-sm text-gray-600">Add, edit, and configure restaurants</p>
                  </div>
                </Button>
                
                <Button variant="ghost" class="w-full justify-start group hover:bg-gradient-to-r hover:from-green-50 hover:to-emerald-50" on:click={() => goto('/feedback/manage')}>
                  <div class="h-10 w-10 bg-gradient-to-br from-green-500 to-emerald-600 rounded-xl flex items-center justify-center mr-4 group-hover:scale-110 transition-transform">
                    <svg class="h-5 w-5 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 12h.01M12 12h.01M16 12h.01M21 12c0 4.418-4.03 8-9 8a9.863 9.863 0 01-4.255-.949L3 20l1.395-3.72C3.512 15.042 3 13.574 3 12c0-4.418 4.03-8 9-8s9 3.582 9 8z" />
                    </svg>
                  </div>
                  <div class="text-left">
                    <p class="font-semibold text-gray-900">Manage Feedback</p>
                    <p class="text-sm text-gray-600">Review and respond to feedback</p>
                  </div>
                </Button>
                
                <Button variant="ghost" class="w-full justify-start group hover:bg-gradient-to-r hover:from-purple-50 hover:to-pink-50" on:click={() => goto('/analytics')}>
                  <div class="h-10 w-10 bg-gradient-to-br from-purple-500 to-pink-600 rounded-xl flex items-center justify-center mr-4 group-hover:scale-110 transition-transform">
                    <svg class="h-5 w-5 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v14a2 2 0 01-2 2h-2a2 2 0 01-2-2z" />
                    </svg>
                  </div>
                  <div class="text-left">
                    <p class="font-semibold text-gray-900">View Analytics</p>
                    <p class="text-sm text-gray-600">Detailed insights and reports</p>
                  </div>
                </Button>
              </div>

              {#if stats.topRatedDish}
                <div class="relative p-6 bg-gradient-to-br from-green-50 to-emerald-50 rounded-2xl border border-green-200/50 shadow-lg shadow-green-500/10">
                  <div class="absolute inset-0 bg-gradient-to-br from-green-500/5 to-emerald-500/5 rounded-2xl"></div>
                  <div class="relative flex items-center space-x-4">
                    <div class="h-12 w-12 bg-gradient-to-br from-green-500 to-emerald-600 rounded-xl flex items-center justify-center shadow-lg shadow-green-500/25">
                      <svg class="h-6 w-6 text-white" fill="currentColor" viewBox="0 0 24 24">
                        <path d="M11.049 2.927c.3-.921 1.603-.921 1.902 0l1.519 4.674a1 1 0 00.95.69h4.915c.969 0 1.371 1.24.588 1.81l-3.976 2.888a1 1 0 00-.363 1.118l1.518 4.674c.3.922-.755 1.688-1.538 1.118l-3.976-2.888a1 1 0 00-1.176 0l-3.976 2.888c-.783.57-1.838-.197-1.538-1.118l1.518-4.674a1 1 0 00-.363-1.118l-3.976-2.888c-.784-.57-.38-1.81.588-1.81h4.914a1 1 0 00.951-.69l1.519-4.674z" />
                      </svg>
                    </div>
                    <div>
                      <p class="text-sm font-semibold text-green-800 uppercase tracking-wide">Top Rated Dish</p>
                      <p class="font-bold text-green-900">{stats.topRatedDish.name}</p>
                      <div class="flex items-center space-x-1 mt-1">
                        <div class="flex text-yellow-400">
                          {#each Array(5) as _, i}
                            <svg class="h-3 w-3 {i < Math.round(stats.topRatedDish.rating) ? 'text-yellow-400' : 'text-gray-300'}" fill="currentColor" viewBox="0 0 24 24">
                              <path d="M11.049 2.927c.3-.921 1.603-.921 1.902 0l1.519 4.674a1 1 0 00.95.69h4.915c.969 0 1.371 1.24.588 1.81l-3.976 2.888a1 1 0 00-.363 1.118l1.518 4.674c.3.922-.755 1.688-1.538 1.118l-3.976-2.888a1 1 0 00-1.176 0l-3.976 2.888c-.783.57-1.838-.197-1.538-1.118l1.518-4.674a1 1 0 00-.363-1.118l-3.976-2.888c-.784-.57-.38-1.81.588-1.81h4.914a1 1 0 00.951-.69l1.519-4.674z" />
                            </svg>
                          {/each}
                        </div>
                        <span class="text-sm font-medium text-green-700">{stats.topRatedDish.rating.toFixed(1)}</span>
                      </div>
                    </div>
                  </div>
                </div>
              {/if}
            </div>
          </Card>
        </div>

        <!-- Recent Feedback -->
        <div class="lg:col-span-2">
          <Card variant="elevated" padding={false}>
            <div class="p-6 lg:p-8">
              <div class="mb-6 flex justify-between items-center">
                <div class="space-y-2">
                  <h3 class="text-xl font-bold bg-gradient-to-r from-gray-900 to-gray-700 bg-clip-text text-transparent">
                    Recent Feedback
                  </h3>
                  <p class="text-gray-600 font-medium">Latest customer reviews and ratings</p>
                </div>
                <Button variant="gradient" size="sm" on:click={() => goto('/feedback/manage')}>
                  <svg class="h-4 w-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5l7 7-7 7" />
                  </svg>
                  View All
                </Button>
              </div>

              <div class="space-y-4 max-h-96 overflow-y-auto">
                {#each recentFeedback as feedback}
                  <div class="relative group">
                    <div class="absolute inset-0 bg-gradient-to-r {feedback.rating >= 4 ? 'from-green-500/10 to-emerald-500/10' : feedback.rating >= 3 ? 'from-yellow-500/10 to-orange-500/10' : 'from-red-500/10 to-pink-500/10'} rounded-xl opacity-0 group-hover:opacity-100 transition-opacity duration-300"></div>
                    <div class="relative bg-white/50 backdrop-blur-sm rounded-xl border border-gray-100 p-5 hover:shadow-lg hover:shadow-gray-900/5 transition-all duration-300 group-hover:border-gray-200">
                      <div class="flex items-start justify-between mb-3">
                        <div class="flex items-center space-x-3">
                          <div class="h-10 w-10 bg-gradient-to-br {feedback.rating >= 4 ? 'from-green-500 to-emerald-600' : feedback.rating >= 3 ? 'from-yellow-500 to-orange-600' : 'from-red-500 to-pink-600'} rounded-xl flex items-center justify-center shadow-lg">
                            <svg class="h-5 w-5 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 12h.01M12 12h.01M16 12h.01M21 12c0 4.418-4.03 8-9 8a9.863 9.863 0 01-4.255-.949L3 20l1.395-3.72C3.512 15.042 3 13.574 3 12c0-4.418 4.03-8 9-8s9 3.582 9 8z" />
                            </svg>
                          </div>
                          <div>
                            <div class="flex items-center space-x-2 mb-1">
                              <div class="flex text-yellow-400">
                                {#each Array(5) as _, i}
                                  <svg class="h-4 w-4 {i < feedback.rating ? 'text-yellow-400' : 'text-gray-300'}" fill="currentColor" viewBox="0 0 24 24">
                                    <path d="M11.049 2.927c.3-.921 1.603-.921 1.902 0l1.519 4.674a1 1 0 00.95.69h4.915c.969 0 1.371 1.24.588 1.81l-3.976 2.888a1 1 0 00-.363 1.118l1.518 4.674c.3.922-.755 1.688-1.538 1.118l-3.976-2.888a1 1 0 00-1.176 0l-3.976 2.888c-.783.57-1.838-.197-1.538-1.118l1.518-4.674a1 1 0 00-.363-1.118l-3.976-2.888c-.784-.57-.38-1.81.588-1.81h4.914a1 1 0 00.951-.69l1.519-4.674z" />
                                  </svg>
                                {/each}
                              </div>
                              <span class="text-sm font-semibold text-gray-600">{feedback.rating}/5</span>
                            </div>
                            {#if feedback.dish_name}
                              <div class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-blue-100 text-blue-800">
                                {feedback.dish_name}
                              </div>
                            {/if}
                          </div>
                        </div>
                        <div class="text-xs text-gray-500 font-medium">
                          {formatDate(feedback.created_at)}
                        </div>
                      </div>
                      
                      {#if feedback.comment}
                        <div class="bg-gray-50/50 rounded-lg p-4 mb-3 border border-gray-100">
                          <p class="text-gray-700 text-sm leading-relaxed">"{feedback.comment}"</p>
                        </div>
                      {/if}
                      
                      <div class="flex items-center justify-between">
                        <div class="flex items-center space-x-2 text-xs text-gray-500">
                          <svg class="h-4 w-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z" />
                          </svg>
                          <span class="font-medium">
                            {#if feedback.customer_email}
                              {feedback.customer_email}
                            {:else}
                              Anonymous Customer
                            {/if}
                          </span>
                        </div>
                        <div class="flex items-center space-x-1 text-xs text-gray-400">
                          <svg class="h-3 w-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z" />
                          </svg>
                          <span>Just now</span>
                        </div>
                      </div>
                    </div>
                  </div>
                {/each}

                {#if recentFeedback.length === 0}
                  <div class="text-center py-12">
                    <div class="h-20 w-20 bg-gradient-to-br from-gray-100 to-gray-200 rounded-2xl flex items-center justify-center mx-auto mb-4">
                      <svg class="h-10 w-10 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 12h.01M12 12h.01M16 12h.01M21 12c0 4.418-4.03 8-9 8a9.863 9.863 0 01-4.255-.949L3 20l1.395-3.72C3.512 15.042 3 13.574 3 12c0-4.418 4.03-8 9-8s9 3.582 9 8z" />
                      </svg>
                    </div>
                    <h4 class="text-lg font-semibold text-gray-900 mb-2">No feedback yet</h4>
                    <p class="text-gray-600 mb-4">Start collecting feedback by sharing your QR codes!</p>
                    <Button variant="gradient" on:click={() => goto('/restaurants')}>
                      <svg class="h-4 w-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4" />
                      </svg>
                      Get Started
                    </Button>
                  </div>
                {/if}
              </div>
            </div>
          </Card>
        </div>
      </div>
    {/if}
  </div>
</div>