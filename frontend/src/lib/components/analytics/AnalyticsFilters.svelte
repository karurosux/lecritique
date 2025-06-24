<script lang="ts">
  import { Card } from '$lib/components/ui';

  interface Restaurant {
    id: string;
    name: string;
  }

  let {
    restaurants = [],
    selectedRestaurant = $bindable(''),
    selectedTimeframe = $bindable('7d'),
    onrestaurantchange = () => {},
    ontimeframechange = () => {}
  }: {
    restaurants?: Restaurant[];
    selectedRestaurant?: string;
    selectedTimeframe?: string;
    onrestaurantchange?: () => void;
    ontimeframechange?: () => void;
  } = $props();

  let restaurantCount = $derived(restaurants.length);
  let selectedRestaurantName = $derived(
    restaurants.find(r => r.id === selectedRestaurant)?.name || 'Select Restaurant'
  );

  function handleRestaurantChange(e: Event) {
    const target = e.target as HTMLSelectElement;
    selectedRestaurant = target.value;
    onrestaurantchange();
  }

  function handleTimeframeChange(e: Event) {
    const target = e.target as HTMLSelectElement;
    selectedTimeframe = target.value;
    ontimeframechange();
  }
</script>

<Card variant="glass" class="mb-6">
  <div class="flex flex-col sm:flex-row gap-4 items-start sm:items-center justify-between">
    <div class="flex flex-col sm:flex-row gap-4 items-start sm:items-center">
      <!-- Restaurant Selector -->
      <div class="space-y-1">
        <label for="restaurant-select" class="text-xs font-medium text-gray-500 uppercase tracking-wider">
          Restaurant
        </label>
        <select
          id="restaurant-select"
          bind:value={selectedRestaurant}
          onchange={handleRestaurantChange}
          class="px-4 py-2 pr-10 border border-gray-200 rounded-xl bg-white focus:outline-none focus:ring-2 focus:ring-purple-500 focus:border-transparent cursor-pointer shadow-sm hover:shadow-md transition-all duration-200 text-gray-700 font-medium min-w-[200px]"
        >
          {#each restaurants as restaurant}
            <option value={restaurant.id}>{restaurant.name}</option>
          {/each}
        </select>
      </div>

      <!-- Timeframe Selector -->
      <div class="space-y-1">
        <label for="timeframe-select" class="text-xs font-medium text-gray-500 uppercase tracking-wider">
          Time Period
        </label>
        <select
          id="timeframe-select"
          bind:value={selectedTimeframe}
          onchange={handleTimeframeChange}
          class="px-4 py-2 pr-10 border border-gray-200 rounded-xl bg-white focus:outline-none focus:ring-2 focus:ring-purple-500 focus:border-transparent cursor-pointer shadow-sm hover:shadow-md transition-all duration-200 text-gray-700 font-medium"
        >
          <option value="1d">Last 24 hours</option>
          <option value="7d">Last 7 days</option>
          <option value="30d">Last 30 days</option>
          <option value="90d">Last 90 days</option>
          <option value="1y">Last year</option>
        </select>
      </div>
    </div>

    <div class="flex items-center space-x-4 text-sm">
      <div class="flex items-center space-x-2">
        <svg class="h-4 w-4 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 21V5a2 2 0 00-2-2H7a2 2 0 00-2 2v16m14 0h2m-2 0h-4m-5 0H9m11 0a2 2 0 01-2 2H7a2 2 0 01-2-2m2-4h2.01M7 16h6M7 8h6v4H7V8z" />
        </svg>
        <span class="text-gray-600 font-medium">{restaurantCount} {restaurantCount === 1 ? 'Restaurant' : 'Restaurants'}</span>
      </div>
      <div class="h-4 w-px bg-gray-200"></div>
      <div class="flex items-center space-x-2">
        <svg class="h-4 w-4 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z" />
        </svg>
        <span class="text-gray-600 font-medium">
          {#if selectedTimeframe === '1d'}
            Last 24 hours
          {:else if selectedTimeframe === '7d'}
            Last 7 days
          {:else if selectedTimeframe === '30d'}
            Last 30 days
          {:else if selectedTimeframe === '90d'}
            Last 90 days
          {:else if selectedTimeframe === '1y'}
            Last year
          {/if}
        </span>
      </div>
    </div>
  </div>
</Card>