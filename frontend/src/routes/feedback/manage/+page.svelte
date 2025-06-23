<script lang="ts">
  import { onMount } from 'svelte';
  import { page } from '$app/stores';
  import { Card, Button, Input } from '$lib/components/ui';
  import { getApiClient, handleApiError } from '$lib/api/client';
  import { auth } from '$lib/stores/auth';
  import { goto } from '$app/navigation';

  interface Feedback {
    id: string;
    customer_email?: string;
    rating: number;
    comment?: string;
    dish_name?: string;
    restaurant_name?: string;
    location_name?: string;
    qr_code?: string;
    responses?: Record<string, any>;
    created_at: string;
  }

  interface FeedbackFilters {
    restaurant_id?: string;
    location_id?: string;
    rating_min?: number;
    rating_max?: number;
    date_from?: string;
    date_to?: string;
    search?: string;
  }

  let loading = true;
  let error = '';
  let feedback: Feedback[] = [];
  let restaurants: any[] = [];
  let totalCount = 0;
  let currentPage = 1;
  let itemsPerPage = 20;
  
  // Filters
  let filters: FeedbackFilters = {};
  let searchQuery = '';
  let selectedRestaurant = '';
  let selectedRatingMin = 1;
  let selectedRatingMax = 5;
  let dateFrom = '';
  let dateTo = '';

  $: authState = $auth;

  onMount(async () => {
    if (!authState.isAuthenticated) {
      goto('/login');
      return;
    }

    await loadRestaurants();
    await loadFeedback();
  });

  async function loadRestaurants() {
    try {
      const api = getApiClient();
      const response = await api.api.v1RestaurantsList();
      
      if (response.data.success && response.data.data) {
        restaurants = response.data.data;
      }
    } catch (err) {
      console.error('Error loading restaurants:', err);
    }
  }

  async function loadFeedback() {
    loading = true;
    error = '';

    try {
      const api = getApiClient();
      
      // Build filters
      const currentFilters: FeedbackFilters = {
        ...(selectedRestaurant && { restaurant_id: selectedRestaurant }),
        ...(selectedRatingMin > 1 && { rating_min: selectedRatingMin }),
        ...(selectedRatingMax < 5 && { rating_max: selectedRatingMax }),
        ...(dateFrom && { date_from: dateFrom }),
        ...(dateTo && { date_to: dateTo }),
        ...(searchQuery && { search: searchQuery })
      };

      // For demo purposes, we'll use the analytics endpoint to get recent feedback
      // In a real implementation, there would be a dedicated feedback list endpoint
      if (restaurants.length > 0) {
        const restaurantId = selectedRestaurant || restaurants[0].id;
        const analyticsResponse = await api.api.v1AnalyticsRestaurantsDetail(restaurantId);
        
        const analyticsData = analyticsResponse.data;
        const recentFeedbackData = analyticsData?.recent_feedback || [];
        
        // Map and filter feedback
        let mappedFeedback = recentFeedbackData.map((fb: any) => ({
          id: fb.id,
          customer_email: fb.customer_email,
          rating: fb.rating,
          comment: fb.comment,
          dish_name: fb.dish_name,
          restaurant_name: restaurants.find(r => r.id === restaurantId)?.name,
          location_name: fb.location_name,
          qr_code: fb.qr_code,
          responses: fb.responses,
          created_at: fb.created_at
        }));

        // Apply client-side filters for demo
        if (searchQuery) {
          const query = searchQuery.toLowerCase();
          mappedFeedback = mappedFeedback.filter((fb: Feedback) =>
            fb.comment?.toLowerCase().includes(query) ||
            fb.customer_email?.toLowerCase().includes(query) ||
            fb.dish_name?.toLowerCase().includes(query)
          );
        }

        if (selectedRatingMin > 1) {
          mappedFeedback = mappedFeedback.filter((fb: Feedback) => fb.rating >= selectedRatingMin);
        }

        if (selectedRatingMax < 5) {
          mappedFeedback = mappedFeedback.filter((fb: Feedback) => fb.rating <= selectedRatingMax);
        }

        if (dateFrom) {
          const fromDate = new Date(dateFrom);
          mappedFeedback = mappedFeedback.filter((fb: Feedback) => new Date(fb.created_at) >= fromDate);
        }

        if (dateTo) {
          const toDate = new Date(dateTo);
          mappedFeedback = mappedFeedback.filter((fb: Feedback) => new Date(fb.created_at) <= toDate);
        }

        feedback = mappedFeedback;
        totalCount = mappedFeedback.length;
      } else {
        feedback = [];
        totalCount = 0;
      }

    } catch (err) {
      error = handleApiError(err);
    } finally {
      loading = false;
    }
  }

  function handleFilterChange() {
    loadFeedback();
  }

  function clearFilters() {
    searchQuery = '';
    selectedRestaurant = '';
    selectedRatingMin = 1;
    selectedRatingMax = 5;
    dateFrom = '';
    dateTo = '';
    loadFeedback();
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

  function exportToCSV() {
    const headers = ['ID', 'Date', 'Customer Email', 'Rating', 'Dish', 'Restaurant', 'Comment'];
    const csvContent = [
      headers.join(','),
      ...feedback.map(fb => [
        fb.id,
        fb.created_at,
        fb.customer_email || 'Anonymous',
        fb.rating,
        fb.dish_name || '',
        fb.restaurant_name || '',
        `"${fb.comment?.replace(/"/g, '""') || ''}"`
      ].join(','))
    ].join('\n');

    const blob = new Blob([csvContent], { type: 'text/csv;charset=utf-8;' });
    const link = document.createElement('a');
    const url = URL.createObjectURL(blob);
    link.setAttribute('href', url);
    link.setAttribute('download', `feedback-export-${new Date().toISOString().split('T')[0]}.csv`);
    link.style.visibility = 'hidden';
    document.body.appendChild(link);
    link.click();
    document.body.removeChild(link);
  }
</script>

<svelte:head>
  <title>Feedback Management - LeCritique</title>
  <meta name="description" content="Manage and analyze customer feedback" />
</svelte:head>

<div class="min-h-screen bg-gray-50">
  <!-- Header -->
  <div class="bg-white shadow-sm border-b">
    <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
      <div class="flex justify-between items-center py-6">
        <div>
          <h1 class="text-2xl font-bold text-gray-900">Feedback Management</h1>
          <p class="text-gray-600">Review and analyze customer feedback from all your restaurants.</p>
        </div>
        <div class="flex space-x-3">
          <Button variant="outline" on:click={exportToCSV} disabled={feedback.length === 0}>
            <svg class="h-4 w-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 10v6m0 0l-3-3m3 3l3-3m2 8H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
            </svg>
            Export CSV
          </Button>
          <Button on:click={() => goto('/dashboard')}>
            <svg class="h-4 w-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 12l2-2m0 0l7-7 7 7M5 10v10a1 1 0 001 1h3m10-11l2 2m-2-2v10a1 1 0 01-1 1h-3m-6 0a1 1 0 001-1v-4a1 1 0 011-1h2a1 1 0 011 1v4a1 1 0 001 1m-6 0h6" />
            </svg>
            Dashboard
          </Button>
        </div>
      </div>
    </div>
  </div>

  <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
    <!-- Filters -->
    <Card class="mb-6">
      <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">
        <!-- Search -->
        <div>
          <label class="block text-sm font-medium text-gray-700 mb-1">Search</label>
          <Input
            type="text"
            placeholder="Search comments, emails, dishes..."
            bind:value={searchQuery}
            on:input={handleFilterChange}
          />
        </div>

        <!-- Restaurant Filter -->
        <div>
          <label class="block text-sm font-medium text-gray-700 mb-1">Restaurant</label>
          <select
            bind:value={selectedRestaurant}
            on:change={handleFilterChange}
            class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
          >
            <option value="">All Restaurants</option>
            {#each restaurants as restaurant}
              <option value={restaurant.id}>{restaurant.name}</option>
            {/each}
          </select>
        </div>

        <!-- Rating Range -->
        <div>
          <label class="block text-sm font-medium text-gray-700 mb-1">Rating Range</label>
          <div class="flex space-x-2">
            <select
              bind:value={selectedRatingMin}
              on:change={handleFilterChange}
              class="flex-1 px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
            >
              {#each [1, 2, 3, 4, 5] as rating}
                <option value={rating}>{rating}★</option>
              {/each}
            </select>
            <span class="self-center text-gray-500">to</span>
            <select
              bind:value={selectedRatingMax}
              on:change={handleFilterChange}
              class="flex-1 px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
            >
              {#each [1, 2, 3, 4, 5] as rating}
                <option value={rating}>{rating}★</option>
              {/each}
            </select>
          </div>
        </div>

        <!-- Date Range -->
        <div>
          <label class="block text-sm font-medium text-gray-700 mb-1">Date Range</label>
          <div class="flex space-x-2">
            <input
              type="date"
              bind:value={dateFrom}
              on:change={handleFilterChange}
              class="flex-1 px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
            />
            <input
              type="date"
              bind:value={dateTo}
              on:change={handleFilterChange}
              class="flex-1 px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
            />
          </div>
        </div>
      </div>

      <div class="mt-4 flex justify-between items-center">
        <p class="text-sm text-gray-600">
          Showing {feedback.length} of {totalCount} feedback entries
        </p>
        <Button variant="outline" size="sm" on:click={clearFilters}>
          Clear Filters
        </Button>
      </div>
    </Card>

    {#if loading}
      <!-- Loading State -->
      <div class="space-y-4">
        {#each Array(5) as _}
          <Card>
            <div class="animate-pulse">
              <div class="flex justify-between items-start mb-4">
                <div class="space-y-2">
                  <div class="h-4 bg-gray-200 rounded w-32"></div>
                  <div class="h-4 bg-gray-200 rounded w-48"></div>
                </div>
                <div class="h-6 bg-gray-200 rounded w-20"></div>
              </div>
              <div class="h-16 bg-gray-200 rounded"></div>
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
          <h3 class="text-lg font-medium text-gray-900 mb-2">Failed to load feedback</h3>
          <p class="text-gray-600 mb-4">{error}</p>
          <Button on:click={loadFeedback}>Try Again</Button>
        </div>
      </Card>

    {:else if feedback.length === 0}
      <!-- Empty State -->
      <Card>
        <div class="text-center py-12">
          <svg class="h-12 w-12 text-gray-400 mx-auto mb-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 12h.01M12 12h.01M16 12h.01M21 12c0 4.418-4.03 8-9 8a9.863 9.863 0 01-4.255-.949L3 20l1.395-3.72C3.512 15.042 3 13.574 3 12c0-4.418 4.03-8 9-8s9 3.582 9 8z" />
          </svg>
          <h3 class="text-lg font-medium text-gray-900 mb-2">No feedback found</h3>
          <p class="text-gray-600 mb-4">No feedback matches your current filters.</p>
          <Button variant="outline" on:click={clearFilters}>Clear Filters</Button>
        </div>
      </Card>

    {:else}
      <!-- Feedback List -->
      <div class="space-y-4">
        {#each feedback as fb}
          <Card>
            <div class="flex justify-between items-start mb-4">
              <div class="flex-1">
                <div class="flex items-center space-x-4 mb-2">
                  <div class="flex items-center space-x-2">
                    <span class="text-lg {getRatingColor(fb.rating)}">{renderStars(fb.rating)}</span>
                    <span class="text-sm text-gray-500">{fb.rating}/5</span>
                  </div>
                  
                  {#if fb.dish_name}
                    <span class="text-sm text-gray-400">•</span>
                    <span class="text-sm font-medium text-gray-700">{fb.dish_name}</span>
                  {/if}
                  
                  {#if fb.restaurant_name}
                    <span class="text-sm text-gray-400">•</span>
                    <span class="text-sm text-gray-600">{fb.restaurant_name}</span>
                  {/if}
                </div>
                
                <div class="flex items-center text-xs text-gray-500 space-x-4">
                  <span>{formatDate(fb.created_at)}</span>
                  {#if fb.customer_email}
                    <span>• {fb.customer_email}</span>
                  {:else}
                    <span>• Anonymous</span>
                  {/if}
                  {#if fb.qr_code}
                    <span>• QR: {fb.qr_code}</span>
                  {/if}
                </div>
              </div>
              
              <div class="flex items-center space-x-2">
                <span class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-blue-100 text-blue-800">
                  ID: {fb.id.slice(0, 8)}
                </span>
              </div>
            </div>
            
            {#if fb.comment}
              <div class="bg-gray-50 rounded-lg p-4 mb-4">
                <p class="text-gray-700 text-sm">"{fb.comment}"</p>
              </div>
            {/if}
            
            {#if fb.responses && Object.keys(fb.responses).length > 0}
              <div class="border-t pt-4">
                <h4 class="text-sm font-medium text-gray-700 mb-2">Questionnaire Responses:</h4>
                <div class="grid grid-cols-1 md:grid-cols-2 gap-2 text-sm">
                  {#each Object.entries(fb.responses) as [key, value]}
                    <div class="flex justify-between">
                      <span class="text-gray-600">{key}:</span>
                      <span class="text-gray-900 font-medium">{typeof value === 'object' ? JSON.stringify(value) : value}</span>
                    </div>
                  {/each}
                </div>
              </div>
            {/if}
          </Card>
        {/each}
      </div>
    {/if}
  </div>
</div>