<script lang="ts">
  import { onMount } from 'svelte';
  import { Card, Button } from '$lib/components/ui';
  import { getApiClient, handleApiError } from '$lib/api/client';
  import { auth } from '$lib/stores/auth';
  import { goto } from '$app/navigation';

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

  let loading = true;
  let error = '';
  let restaurants: Restaurant[] = [];

  $: authState = $auth;

  onMount(async () => {
    if (!authState.isAuthenticated) {
      goto('/login');
      return;
    }

    await loadRestaurants();
  });

  async function loadRestaurants() {
    loading = true;
    error = '';

    try {
      const api = getApiClient();
      
      // Use actual API client to get restaurants
      const response = await api.api.v1RestaurantsList();
      
      if (response.data.success && response.data.data) {
        restaurants = response.data.data.map(restaurant => ({
          id: restaurant.id || '',
          name: restaurant.name || '',
          description: restaurant.description || '',
          address: '', // Note: address would come from locations array
          phone: restaurant.phone || '',
          email: restaurant.email || '',
          website: restaurant.website || '',
          cuisine_type: '', // Note: cuisine_type not in API model
          status: restaurant.is_active ? 'active' : 'inactive',
          created_at: restaurant.created_at || '',
          updated_at: restaurant.updated_at || ''
        }));
      } else {
        restaurants = [];
      }

    } catch (err) {
      error = handleApiError(err);
    } finally {
      loading = false;
    }
  }

  function formatDate(dateString: string): string {
    return new Date(dateString).toLocaleDateString();
  }

  function getStatusColor(status: string): string {
    return status === 'active' ? 'bg-green-100 text-green-800' : 'bg-gray-100 text-gray-800';
  }

  async function toggleRestaurantStatus(restaurant: Restaurant) {
    try {
      const api = getApiClient();
      const newStatus = restaurant.status === 'active' ? 'inactive' : 'active';
      
      // Update restaurant status via API
      await api.api.v1RestaurantsUpdate(restaurant.id, {
        is_active: newStatus === 'active'
      });
      
      // Update local state
      restaurants = restaurants.map(r => 
        r.id === restaurant.id 
          ? { ...r, status: newStatus, updated_at: new Date().toISOString() }
          : r
      );
      
    } catch (err) {
      error = handleApiError(err);
    }
  }

  async function deleteRestaurant(restaurant: Restaurant) {
    if (!confirm(`Are you sure you want to delete "${restaurant.name}"? This action cannot be undone.`)) {
      return;
    }

    try {
      const api = getApiClient();
      
      // Delete restaurant via API
      await api.api.v1RestaurantsDelete(restaurant.id);
      
      // Update local state
      restaurants = restaurants.filter(r => r.id !== restaurant.id);
      
    } catch (err) {
      error = handleApiError(err);
    }
  }
</script>

<svelte:head>
  <title>Restaurants - LeCritique</title>
  <meta name="description" content="Manage your restaurants" />
</svelte:head>

<div class="min-h-screen bg-gray-50">
  <!-- Header -->
  <div class="bg-white shadow-sm border-b">
    <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
      <div class="flex justify-between items-center py-6">
        <div>
          <h1 class="text-2xl font-bold text-gray-900">Restaurants</h1>
          <p class="text-gray-600">Manage your restaurant locations and information</p>
        </div>
        <div class="flex space-x-3">
          <Button variant="outline" on:click={() => goto('/dashboard')}>
            <svg class="h-4 w-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 19l-7-7 7-7" />
            </svg>
            Back to Dashboard
          </Button>
          <Button on:click={() => goto('/restaurants/new')}>
            <svg class="h-4 w-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4" />
            </svg>
            Add Restaurant
          </Button>
        </div>
      </div>
    </div>
  </div>

  <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
    {#if loading}
      <!-- Loading State -->
      <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
        {#each Array(3) as _}
          <Card>
            <div class="animate-pulse">
              <div class="h-6 bg-gray-200 rounded w-3/4 mb-3"></div>
              <div class="h-4 bg-gray-200 rounded w-full mb-2"></div>
              <div class="h-4 bg-gray-200 rounded w-2/3 mb-4"></div>
              <div class="h-8 bg-gray-200 rounded w-1/3"></div>
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
          <h3 class="text-lg font-medium text-gray-900 mb-2">Failed to load restaurants</h3>
          <p class="text-gray-600 mb-4">{error}</p>
          <Button on:click={loadRestaurants}>Try Again</Button>
        </div>
      </Card>

    {:else if restaurants.length === 0}
      <!-- Empty State -->
      <Card>
        <div class="text-center py-16">
          <svg class="h-16 w-16 text-gray-400 mx-auto mb-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 21V5a2 2 0 00-2-2H7a2 2 0 00-2 2v16m14 0h2m-2 0h-4m-5 0H9m11 0a2 2 0 01-2 2H7a2 2 0 01-2-2m2-4h2.01M7 16h6M7 8h6v4H7V8z" />
          </svg>
          <h3 class="text-xl font-medium text-gray-900 mb-2">No restaurants yet</h3>
          <p class="text-gray-600 mb-6">Get started by adding your first restaurant location</p>
          <Button on:click={() => goto('/restaurants/new')}>
            <svg class="h-4 w-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4" />
            </svg>
            Add Your First Restaurant
          </Button>
        </div>
      </Card>

    {:else}
      <!-- Restaurant Grid -->
      <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
        {#each restaurants as restaurant}
          <Card class="relative">
            <!-- Status Badge -->
            <div class="absolute top-4 right-4">
              <span class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium {getStatusColor(restaurant.status)}">
                {restaurant.status}
              </span>
            </div>

            <div class="pr-16">
              <h3 class="text-lg font-medium text-gray-900 mb-2">{restaurant.name}</h3>
              
              {#if restaurant.description}
                <p class="text-sm text-gray-600 mb-3 line-clamp-2">{restaurant.description}</p>
              {/if}

              <div class="space-y-2 mb-4">
                <!-- Address -->
                <div class="flex items-start text-sm text-gray-600">
                  <svg class="h-4 w-4 mr-2 mt-0.5 flex-shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17.657 16.657L13.414 20.9a1.998 1.998 0 01-2.827 0l-4.244-4.243a8 8 0 1111.314 0z" />
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 11a3 3 0 11-6 0 3 3 0 016 0z" />
                  </svg>
                  <span class="break-words">{restaurant.address}</span>
                </div>

                <!-- Phone -->
                {#if restaurant.phone}
                  <div class="flex items-center text-sm text-gray-600">
                    <svg class="h-4 w-4 mr-2 flex-shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 5a2 2 0 012-2h3.28a1 1 0 01.948.684l1.498 4.493a1 1 0 01-.502 1.21l-2.257 1.13a11.042 11.042 0 005.516 5.516l1.13-2.257a1 1 0 011.21-.502l4.493 1.498a1 1 0 01.684.949V19a2 2 0 01-2 2h-1C9.716 21 3 14.284 3 6V5z" />
                    </svg>
                    <span>{restaurant.phone}</span>
                  </div>
                {/if}

                <!-- Cuisine Type -->
                {#if restaurant.cuisine_type}
                  <div class="flex items-center text-sm text-gray-600">
                    <svg class="h-4 w-4 mr-2 flex-shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M7 7h.01M7 3h5c.512 0 1.024.195 1.414.586l7 7a2 2 0 010 2.828l-7 7a2 2 0 01-2.828 0l-7-7A1.994 1.994 0 013 12V7a4 4 0 014-4z" />
                    </svg>
                    <span>{restaurant.cuisine_type}</span>
                  </div>
                {/if}
              </div>

              <div class="text-xs text-gray-500 mb-4">
                Created: {formatDate(restaurant.created_at)}
              </div>

              <!-- Action Buttons -->
              <div class="flex flex-wrap gap-2">
                <Button 
                  variant="outline" 
                  size="sm" 
                  on:click={() => goto(`/restaurants/${restaurant.id}/edit`)}
                >
                  <svg class="h-3 w-3 mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z" />
                  </svg>
                  Edit
                </Button>

                <Button 
                  variant="outline" 
                  size="sm" 
                  on:click={() => goto(`/restaurants/${restaurant.id}/dishes`)}
                >
                  <svg class="h-3 w-3 mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5H7a2 2 0 00-2 2v10a2 2 0 002 2h8a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2" />
                  </svg>
                  Dishes
                </Button>

                <Button 
                  variant={restaurant.status === 'active' ? 'secondary' : 'primary'} 
                  size="sm" 
                  on:click={() => toggleRestaurantStatus(restaurant)}
                >
                  {restaurant.status === 'active' ? 'Deactivate' : 'Activate'}
                </Button>

                <Button 
                  variant="destructive" 
                  size="sm" 
                  on:click={() => deleteRestaurant(restaurant)}
                >
                  <svg class="h-3 w-3 mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
                  </svg>
                  Delete
                </Button>
              </div>
            </div>
          </Card>
        {/each}
      </div>
    {/if}
  </div>
</div>