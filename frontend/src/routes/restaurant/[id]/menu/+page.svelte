<script lang="ts">
  import { page } from '$app/stores';
  import { onMount } from 'svelte';
  import { goto } from '$app/navigation';
  import { Card, Button } from '$lib/components/ui';
  import { getApiClient, handleApiError } from '$lib/api/client';

  interface Dish {
    id: string;
    name: string;
    description?: string;
    price: number;
    currency: string;
    category: string;
    is_active: boolean;
    tags?: string[];
  }

  interface Restaurant {
    id: string;
    name: string;
    description?: string;
    email?: string;
    phone?: string;
    website?: string;
    logo?: string;
  }

  interface MenuData {
    restaurant: Restaurant;
    dishes: Dish[];
  }

  let loading = true;
  let error = '';
  let menuData: MenuData | null = null;
  let groupedDishes: Record<string, Dish[]> = {};
  
  $: restaurantId = $page.params.id;
  $: qrCode = $page.url.searchParams.get('qr');
  $: pageTitle = menuData?.restaurant?.name 
    ? `${menuData.restaurant.name} - Menu`
    : 'Restaurant Menu';

  onMount(async () => {
    await loadMenu();
  });

  async function loadMenu() {
    if (!restaurantId) {
      error = 'Restaurant not found';
      loading = false;
      return;
    }

    try {
      loading = true;
      error = '';
      
      const api = getApiClient();
      const response = await api.api.v1PublicRestaurantMenuList(restaurantId);
      
      if (response.data.success && response.data.data) {
        menuData = response.data.data as MenuData;
        
        // Group dishes by category
        groupedDishes = {};
        if (menuData.dishes) {
          menuData.dishes.forEach(dish => {
            if (dish.is_active) {
              const category = dish.category || 'Other';
              if (!groupedDishes[category]) {
                groupedDishes[category] = [];
              }
              groupedDishes[category].push(dish);
            }
          });
        }
      } else {
        error = 'Menu not available';
      }
    } catch (err) {
      error = handleApiError(err);
    } finally {
      loading = false;
    }
  }

  function formatPrice(price: number, currency: string = 'USD') {
    return new Intl.NumberFormat('en-US', {
      style: 'currency',
      currency: currency
    }).format(price);
  }

  function handleGiveFeedback() {
    const params = new URLSearchParams();
    params.set('restaurant', restaurantId);
    if (qrCode) params.set('qr', qrCode);
    goto(`/feedback?${params.toString()}`);
  }

  function handleFeedbackForDish(dishId: string) {
    const params = new URLSearchParams();
    params.set('restaurant', restaurantId);
    params.set('dish', dishId);
    if (qrCode) params.set('qr', qrCode);
    goto(`/feedback?${params.toString()}`);
  }
</script>

<svelte:head>
  <title>{pageTitle}</title>
  <meta name="description" content="View our restaurant menu and give feedback on your favorite dishes" />
</svelte:head>

<div class="min-h-screen bg-gray-50">
  {#if loading}
    <!-- Loading State -->
    <div class="flex items-center justify-center min-h-screen">
      <div class="text-center">
        <svg class="animate-spin h-8 w-8 text-blue-600 mx-auto mb-4" fill="none" viewBox="0 0 24 24">
          <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
          <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 714 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
        </svg>
        <p class="text-gray-600">Loading menu...</p>
      </div>
    </div>
  
  {:else if error}
    <!-- Error State -->
    <div class="flex items-center justify-center min-h-screen px-4">
      <Card>
        <div class="text-center py-12">
          <svg class="h-12 w-12 text-red-500 mx-auto mb-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-2.5L13.732 4c-.77-.833-1.964-.833-2.732 0L3.732 16.5c-.77.833.192 2.5 1.732 2.5z" />
          </svg>
          <h2 class="text-xl font-semibold text-gray-900 mb-2">Menu Not Available</h2>
          <p class="text-gray-600">{error}</p>
        </div>
      </Card>
    </div>
  
  {:else if menuData}
    <!-- Restaurant Header -->
    <div class="bg-white shadow-sm">
      <div class="max-w-4xl mx-auto px-4 py-6">
        <div class="flex items-center space-x-4">
          {#if menuData.restaurant.logo}
            <img 
              src={menuData.restaurant.logo} 
              alt="{menuData.restaurant.name} logo"
              class="h-16 w-16 rounded-full object-cover"
            />
          {:else}
            <div class="h-16 w-16 rounded-full bg-blue-100 flex items-center justify-center">
              <svg class="h-8 w-8 text-blue-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 21V5a2 2 0 00-2-2H7a2 2 0 00-2 2v16m14 0h2m-2 0h-4m-5 0H9m0 0H5m5 0v-4a1 1 0 011-1h2a1 1 0 011 1v4M7 7h3m3 0h3m-6 4h3m3 0h3" />
              </svg>
            </div>
          {/if}
          
          <div class="flex-1">
            <h1 class="text-2xl font-bold text-gray-900">{menuData.restaurant.name}</h1>
            {#if menuData.restaurant.description}
              <p class="text-gray-600 mt-1">{menuData.restaurant.description}</p>
            {/if}
            
            <div class="flex items-center space-x-4 mt-2 text-sm text-gray-500">
              {#if menuData.restaurant.phone}
                <span class="flex items-center">
                  <svg class="h-4 w-4 mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 5a2 2 0 012-2h3.28a1 1 0 01.948.684l1.498 4.493a1 1 0 01-.502 1.21l-2.257 1.13a11.042 11.042 0 005.516 5.516l1.13-2.257a1 1 0 011.21-.502l4.493 1.498a1 1 0 01.684.949V19a2 2 0 01-2 2h-1C9.716 21 3 14.284 3 6V5z" />
                  </svg>
                  {menuData.restaurant.phone}
                </span>
              {/if}
              {#if menuData.restaurant.website}
                <a href={menuData.restaurant.website} target="_blank" rel="noopener noreferrer" class="flex items-center hover:text-blue-600">
                  <svg class="h-4 w-4 mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 12a9 9 0 01-9 9m9-9a9 9 0 00-9-9m9 9H3m9 9v-9m0-9v9m0 9c1.657 0 3-4.03 3-9s-1.343-9-3-9m0 18c-1.657 0-3-4.03-3-9s1.343-9 3-9m-9 9a9 9 0 019-9" />
                  </svg>
                  Website
                </a>
              {/if}
            </div>
          </div>
          
          <Button variant="primary" on:click={handleGiveFeedback}>
            Give Feedback
          </Button>
        </div>
      </div>
    </div>

    <!-- Menu Content -->
    <div class="max-w-4xl mx-auto px-4 py-8">
      {#if Object.keys(groupedDishes).length === 0}
        <Card>
          <div class="text-center py-12">
            <svg class="h-12 w-12 text-gray-400 mx-auto mb-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
            </svg>
            <h3 class="text-lg font-medium text-gray-900 mb-2">No Menu Items Available</h3>
            <p class="text-gray-600">The menu is currently being updated. Please check back later.</p>
          </div>
        </Card>
      {:else}
        <div class="space-y-8">
          {#each Object.entries(groupedDishes) as [category, dishes]}
            <div>
              <h2 class="text-xl font-semibold text-gray-900 mb-4 pb-2 border-b border-gray-200">
                {category}
              </h2>
              
              <div class="grid gap-4">
                {#each dishes as dish}
                  <Card hover>
                    <div class="flex items-center justify-between p-4">
                      <div class="flex-1">
                        <div class="flex items-start justify-between">
                          <div class="flex-1 pr-4">
                            <h3 class="text-lg font-medium text-gray-900">{dish.name}</h3>
                            {#if dish.description}
                              <p class="text-gray-600 mt-1 text-sm">{dish.description}</p>
                            {/if}
                            {#if dish.tags && dish.tags.length > 0}
                              <div class="flex flex-wrap gap-1 mt-2">
                                {#each dish.tags as tag}
                                  <span class="inline-flex items-center px-2 py-1 rounded-full text-xs font-medium bg-gray-100 text-gray-800">
                                    {tag}
                                  </span>
                                {/each}
                              </div>
                            {/if}
                          </div>
                          
                          <div class="text-right">
                            <div class="text-lg font-semibold text-gray-900">
                              {formatPrice(dish.price, dish.currency)}
                            </div>
                            <Button 
                              variant="outline" 
                              size="sm" 
                              class="mt-2"
                              on:click={() => handleFeedbackForDish(dish.id)}
                            >
                              Feedback
                            </Button>
                          </div>
                        </div>
                      </div>
                    </div>
                  </Card>
                {/each}
              </div>
            </div>
          {/each}
        </div>
      {/if}
    </div>
    
    <!-- Floating Feedback Button for Mobile -->
    <div class="fixed bottom-6 right-6 md:hidden">
      <Button 
        variant="primary" 
        size="lg"
        class="rounded-full shadow-lg"
        on:click={handleGiveFeedback}
      >
        <svg class="h-5 w-5 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 12h.01M12 12h.01M16 12h.01M21 12c0 4.418-4.03 8-9 8a9.863 9.863 0 01-4.255-.949L3 20l1.395-3.72C3.512 15.042 3 13.574 3 12c0-4.418 4.03-8 9-8s9 3.582 9 8z" />
        </svg>
        Feedback
      </Button>
    </div>
  {/if}
</div>