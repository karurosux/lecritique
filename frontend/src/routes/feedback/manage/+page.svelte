<script lang="ts">
  import { onMount } from 'svelte';
  import { page } from '$app/stores';
  import { Card, Button, Input, Select } from '$lib/components/ui';
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
  let selectedRatingMin = '1';
  let selectedRatingMax = '5';
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
        ...(parseInt(selectedRatingMin) > 1 && { rating_min: parseInt(selectedRatingMin) }),
        ...(parseInt(selectedRatingMax) < 5 && { rating_max: parseInt(selectedRatingMax) }),
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

        if (parseInt(selectedRatingMin) > 1) {
          mappedFeedback = mappedFeedback.filter((fb: Feedback) => fb.rating >= parseInt(selectedRatingMin));
        }

        if (parseInt(selectedRatingMax) < 5) {
          mappedFeedback = mappedFeedback.filter((fb: Feedback) => fb.rating <= parseInt(selectedRatingMax));
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
    selectedRatingMin = '1';
    selectedRatingMax = '5';
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

  function getRatingBgColor(rating: number): string {
    if (rating >= 4) return 'bg-green-100 text-green-800';
    if (rating >= 3) return 'bg-yellow-100 text-yellow-800';
    return 'bg-red-100 text-red-800';
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

<div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
    <!-- Page Header -->
    <div class="mb-8">
      <div class="flex flex-col lg:flex-row lg:justify-between lg:items-center gap-6">
        <div class="space-y-3">
          <div class="flex items-center space-x-3">
            <div class="h-12 w-12 bg-gradient-to-br from-blue-500 to-purple-600 rounded-2xl flex items-center justify-center shadow-lg shadow-blue-500/25">
              <svg class="h-6 w-6 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 12h.01M12 12h.01M16 12h.01M21 12c0 4.418-4.03 8-9 8a9.863 9.863 0 01-4.255-.949L3 20l1.395-3.72C3.512 15.042 3 13.574 3 12c0-4.418 4.03-8 9-8s9 3.582 9 8z" />
              </svg>
            </div>
            <div>
              <h1 class="text-3xl font-bold bg-gradient-to-r from-gray-900 to-gray-700 bg-clip-text text-transparent">
                Feedback Management
              </h1>
              <div class="flex items-center space-x-4 mt-1">
                <p class="text-gray-600 font-medium">Review and analyze customer feedback from all your restaurants</p>
                {#if !loading && feedback.length > 0}
                  <div class="flex items-center space-x-3 text-sm">
                    <div class="flex items-center space-x-1">
                      <div class="w-2 h-2 bg-blue-400 rounded-full"></div>
                      <span class="text-gray-600">{totalCount} Total</span>
                    </div>
                    <div class="flex items-center space-x-1">
                      <div class="w-2 h-2 bg-purple-400 rounded-full"></div>
                      <span class="text-gray-600">{feedback.filter(f => f.rating >= 4).length} Positive</span>
                    </div>
                  </div>
                {/if}
              </div>
            </div>
          </div>
        </div>
        
        <div class="flex items-center space-x-3">
          <!-- Export CSV Button -->
          <Button 
            variant="gradient" 
            size="lg" 
            class="group relative overflow-hidden shadow-lg hover:shadow-xl transition-all duration-300" 
            onclick={exportToCSV}
            disabled={feedback.length === 0}
          >
            <div class="absolute inset-0 bg-gradient-to-r from-blue-600 to-purple-600 opacity-0 group-hover:opacity-100 transition-opacity duration-300"></div>
            <svg class="h-5 w-5 mr-2 relative z-10 group-hover:scale-110 transition-transform duration-200" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 10v6m0 0l-3-3m3 3l3-3m2 8H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
            </svg>
            <span class="relative z-10">Export CSV</span>
          </Button>
        </div>
      </div>
    </div>
    <!-- Filters -->
    <Card variant="default" hover interactive class="mb-6 group transform transition-all duration-300 animate-fade-in-up">
      <div class="space-y-4">
        <!-- First Row: Search and Restaurant -->
        <div class="grid grid-cols-1 lg:grid-cols-2 gap-4">
          <!-- Search -->
          <div class="space-y-1">
            <label class="text-xs font-medium text-gray-500 uppercase tracking-wider">Search</label>
            <Input
              type="text"
              placeholder="Search comments, emails, dishes..."
              bind:value={searchQuery}
              onchange={handleFilterChange}
              class="w-full"
            />
          </div>

          <!-- Restaurant Filter -->
          <div class="space-y-1">
            <label class="text-xs font-medium text-gray-500 uppercase tracking-wider">Restaurant</label>
            <Select
              bind:value={selectedRestaurant}
              options={[
                { value: '', label: 'All Restaurants' },
                ...restaurants.map(r => ({ value: r.id, label: r.name }))
              ]}
              onchange={handleFilterChange}
              minWidth="min-w-full"
            />
          </div>
        </div>

        <!-- Second Row: Rating and Date Range -->
        <div class="grid grid-cols-1 lg:grid-cols-2 gap-4">
          <!-- Rating Range -->
          <div class="space-y-1">
            <label class="text-xs font-medium text-gray-500 uppercase tracking-wider">Rating Range</label>
            <div class="flex items-center gap-3">
              <div class="flex-1">
                <Select
                  bind:value={selectedRatingMin}
                  options={[1, 2, 3, 4, 5].map(r => ({ value: r.toString(), label: `${r}★` }))}
                  onchange={handleFilterChange}
                  minWidth="min-w-full"
                />
              </div>
              <span class="text-sm text-gray-500 font-medium">to</span>
              <div class="flex-1">
                <Select
                  bind:value={selectedRatingMax}
                  options={[1, 2, 3, 4, 5].map(r => ({ value: r.toString(), label: `${r}★` }))}
                  onchange={handleFilterChange}
                  minWidth="min-w-full"
                />
              </div>
            </div>
          </div>

          <!-- Date Range -->
          <div class="space-y-1">
            <label class="text-xs font-medium text-gray-500 uppercase tracking-wider">Date Range</label>
            <div class="flex items-center gap-3">
              <div class="flex-1">
                <Input
                  type="date"
                  bind:value={dateFrom}
                  onchange={handleFilterChange}
                  placeholder="From"
                  class="w-full"
                />
              </div>
              <span class="text-sm text-gray-500 font-medium">to</span>
              <div class="flex-1">
                <Input
                  type="date"
                  bind:value={dateTo}
                  onchange={handleFilterChange}
                  placeholder="To"
                  class="w-full"
                />
              </div>
            </div>
          </div>
        </div>

        <!-- Filter Summary and Actions -->
        <div class="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-3 pt-4 border-t border-gray-100">
          <div class="flex items-center gap-4 text-sm">
            <div class="flex items-center gap-2">
              <svg class="h-4 w-4 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5H7a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2" />
              </svg>
              <span class="text-gray-600 font-medium">
                {feedback.length} {feedback.length === 1 ? 'entry' : 'entries'}
              </span>
            </div>
            {#if totalCount > feedback.length}
              <div class="h-4 w-px bg-gray-200"></div>
              <span class="text-gray-500">
                {totalCount} total
              </span>
            {/if}
          </div>
          <Button variant="outline" size="sm" onclick={clearFilters}>
            <svg class="h-4 w-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
            </svg>
            Clear Filters
          </Button>
        </div>
      </div>
    </Card>

    {#if loading}
      <!-- Loading State -->
      <div class="space-y-4">
        {#each Array(5) as _}
          <Card variant="default" class="opacity-50">
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
      <Card variant="default" hover interactive class="group">
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
      <Card variant="default" hover interactive class="group">
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
        {#each feedback as fb, index}
          <Card variant="default" hover interactive class="group transform transition-all duration-300 animate-fade-in-up" style="animation-delay: {index * 50}ms">
            <div class="flex justify-between items-start mb-4">
              <div class="flex-1">
                <div class="flex items-center space-x-4 mb-2">
                  <div class="flex items-center space-x-3">
                    <span class="inline-flex items-center px-3 py-1 rounded-full text-sm font-semibold {getRatingBgColor(fb.rating)}">
                      <span class="text-lg mr-1">{renderStars(fb.rating)}</span>
                      {fb.rating}/5
                    </span>
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
              <div class="bg-gradient-to-r from-gray-50 to-gray-100/50 rounded-xl p-4 mb-4 border border-gray-200/50">
                <p class="text-gray-700 text-sm italic">"{fb.comment}"</p>
              </div>
            {/if}
            
            {#if fb.responses && Object.keys(fb.responses).length > 0}
              <div class="border-t border-gray-100 pt-4">
                <h4 class="text-sm font-semibold text-gray-700 mb-3 flex items-center">
                  <svg class="h-4 w-4 mr-2 text-purple-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5H7a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2" />
                  </svg>
                  Questionnaire Responses
                </h4>
                <div class="grid grid-cols-1 md:grid-cols-2 gap-3">
                  {#each Object.entries(fb.responses) as [key, value]}
                    <div class="bg-gray-50 rounded-lg p-3 group-hover:bg-gray-100 transition-colors">
                      <span class="text-xs text-gray-500 uppercase tracking-wider">{key}</span>
                      <p class="text-sm text-gray-900 font-medium mt-1">{typeof value === 'object' ? JSON.stringify(value) : value}</p>
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

<style>
  @keyframes fade-in-up {
    from {
      opacity: 0;
      transform: translateY(10px);
    }
    to {
      opacity: 1;
      transform: translateY(0);
    }
  }

  .animate-fade-in-up {
    animation: fade-in-up 0.6s ease-out forwards;
    opacity: 0;
  }
</style>