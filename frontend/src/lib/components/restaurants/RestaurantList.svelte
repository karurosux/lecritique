<script lang="ts">
  import RestaurantCard from './RestaurantCard.svelte';
  import { createEventDispatcher } from 'svelte';

  interface Restaurant {
    id: string;
    name: string;
    description?: string;
    address: string;
    phone?: string;
    email?: string;
    website?: string;
    cuisine_type?: string;
    status: 'active' | 'inactive';
    created_at: string;
    updated_at: string;
  }

  export let restaurants: Restaurant[] = [];
  export let loading = false;
  export let viewMode: 'grid' | 'list' = 'grid';

  const dispatch = createEventDispatcher();

  function handleRestaurantClick(restaurant: Restaurant) {
    dispatch('restaurantClick', restaurant);
  }

  function handleRestaurantEdit(restaurant: Restaurant) {
    dispatch('restaurantEdit', restaurant);
  }

  function handleRestaurantToggleStatus(restaurant: Restaurant) {
    dispatch('restaurantToggleStatus', restaurant);
  }

  function handleRestaurantDelete(restaurant: Restaurant) {
    dispatch('restaurantDelete', restaurant);
  }
</script>

{#if loading}
  <!-- Loading State -->
  <div class="grid {viewMode === 'grid' ? 'grid-cols-1 md:grid-cols-2 lg:grid-cols-3' : 'grid-cols-1'} gap-6">
    {#each Array(6) as _, i}
      <div class="bg-white rounded-xl border border-gray-200 p-6 animate-pulse">
        <div class="flex items-start justify-between mb-4">
          <div class="flex items-center space-x-3">
            <div class="w-12 h-12 bg-gray-200 rounded-xl"></div>
            <div class="space-y-2">
              <div class="h-4 bg-gray-200 rounded w-32"></div>
              <div class="h-3 bg-gray-200 rounded w-20"></div>
            </div>
          </div>
        </div>
        <div class="space-y-2 mb-4">
          <div class="h-3 bg-gray-200 rounded w-full"></div>
          <div class="h-3 bg-gray-200 rounded w-3/4"></div>
        </div>
        <div class="space-y-2">
          <div class="h-3 bg-gray-200 rounded w-2/3"></div>
          <div class="h-3 bg-gray-200 rounded w-1/2"></div>
        </div>
      </div>
    {/each}
  </div>
{:else if restaurants.length === 0}
  <!-- Empty State -->
  <div class="text-center py-12">
    <div class="w-24 h-24 mx-auto bg-gradient-to-br from-gray-100 to-gray-200 rounded-2xl flex items-center justify-center mb-6">
      <svg class="h-12 w-12 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 21V5a2 2 0 00-2-2H7a2 2 0 00-2 2v16m14 0h2m-2 0h-4m-5 0H9m11 0a2 2 0 01-2 2H7a2 2 0 01-2-2m2-4h2.01M7 16h6M7 8h6v4H7V8z" />
      </svg>
    </div>
    <h3 class="text-xl font-semibold text-gray-900 mb-2">No restaurants found</h3>
    <p class="text-gray-600 mb-6 max-w-md mx-auto">
      You haven't added any restaurants yet, or no restaurants match your current filters.
    </p>
    <div class="flex flex-col sm:flex-row items-center justify-center space-y-3 sm:space-y-0 sm:space-x-4">
      <button 
        class="px-6 py-3 bg-gradient-to-r from-blue-600 to-purple-600 text-white rounded-xl font-medium hover:from-blue-700 hover:to-purple-700 transition-all duration-200 shadow-lg hover:shadow-xl"
        on:click={() => dispatch('addRestaurant')}
      >
        <svg class="h-5 w-5 mr-2 inline" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4" />
        </svg>
        Add Your First Restaurant
      </button>
      <button 
        class="px-6 py-3 text-gray-600 hover:text-gray-800 border border-gray-300 rounded-xl font-medium hover:border-gray-400 transition-all duration-200"
        on:click={() => dispatch('clearFilters')}
      >
        Clear Filters
      </button>
    </div>
  </div>
{:else}
  <!-- Restaurant Cards -->
  <div class="grid {viewMode === 'grid' ? 'grid-cols-1 md:grid-cols-2 lg:grid-cols-3' : 'grid-cols-1'} gap-6">
    {#each restaurants as restaurant, index}
      <RestaurantCard 
        {restaurant} 
        {viewMode}
        {index}
        onclick={handleRestaurantClick}
        onedit={handleRestaurantEdit}
        ontogglestatus={handleRestaurantToggleStatus}
        ondelete={handleRestaurantDelete}
      />
    {/each}
  </div>
{/if}