<script lang="ts">
  import { onMount } from 'svelte';
  import { Card, Button, Input } from '$lib/components/ui';
  import { getApiClient, handleApiError } from '$lib/api/client';
  import { auth } from '$lib/stores/auth';
  import { goto } from '$app/navigation';
  import { page } from '$app/stores';

  interface Dish {
    id: string;
    name: string;
    description?: string;
    price: number;
    category: string;
    is_available: boolean;
    image_url?: string;
    allergens?: string[];
    preparation_time?: number;
    created_at: string;
    updated_at: string;
  }

  interface Category {
    id: string;
    name: string;
    dishes: Dish[];
  }

  let loading = true;
  let error = '';
  let restaurantId = '';
  let restaurantName = '';
  let categories: Category[] = [];
  let showAddDishModal = false;
  let editingDish: Dish | null = null;

  // New dish form
  let newDishForm = {
    name: '',
    description: '',
    price: 0,
    category: '',
    is_available: true,
    allergens: '',
    preparation_time: 0
  };

  $: authState = $auth;
  $: restaurantId = $page.params.id;

  const dishCategories = [
    'Appetizers',
    'Salads',
    'Soups',
    'Main Courses',
    'Seafood',
    'Vegetarian',
    'Pasta',
    'Pizza',
    'Desserts',
    'Beverages',
    'Specials'
  ];

  onMount(async () => {
    if (!authState.isAuthenticated) {
      goto('/login');
      return;
    }

    await loadDishes();
  });

  async function loadDishes() {
    loading = true;
    error = '';

    try {
      const api = getApiClient();
      
      // Get restaurant details and dishes via API
      const [restaurantResponse, dishesResponse] = await Promise.all([
        api.api.v1RestaurantsDetail(restaurantId),
        api.api.v1RestaurantsDishesList(restaurantId)
      ]);
      
      // Set restaurant name
      if (restaurantResponse.data.success && restaurantResponse.data.data) {
        restaurantName = restaurantResponse.data.data.name || 'Unknown Restaurant';
      } else {
        restaurantName = 'Unknown Restaurant';
      }
      
      // Process dishes
      if (dishesResponse.data.success && dishesResponse.data.data) {
        const dishes: Dish[] = dishesResponse.data.data.map((dish: any) => ({
          id: dish.id || '',
          name: dish.name || '',
          description: dish.description || '',
          price: dish.price || 0,
          category: dish.category || 'Uncategorized',
          is_available: dish.is_available !== false, // Default to true if not specified
          allergens: [], // Note: allergens not in API model
          preparation_time: 0, // Note: preparation_time not in API model
          created_at: dish.created_at || '',
          updated_at: dish.updated_at || ''
        }));

        // Group dishes by category
        const categoryMap = new Map<string, Dish[]>();
        
        dishes.forEach(dish => {
          if (!categoryMap.has(dish.category)) {
            categoryMap.set(dish.category, []);
          }
          categoryMap.get(dish.category)!.push(dish);
        });

        categories = Array.from(categoryMap.entries()).map(([name, dishes]) => ({
          id: name.toLowerCase().replace(/\s+/g, '-'),
          name,
          dishes: dishes.sort((a, b) => a.name.localeCompare(b.name))
        }));
      } else {
        categories = [];
      }

    } catch (err) {
      error = handleApiError(err);
    } finally {
      loading = false;
    }
  }

  function openAddDishModal() {
    newDishForm = {
      name: '',
      description: '',
      price: 0,
      category: '',
      is_available: true,
      allergens: '',
      preparation_time: 0
    };
    editingDish = null;
    showAddDishModal = true;
  }

  function openEditDishModal(dish: Dish) {
    newDishForm = {
      name: dish.name,
      description: dish.description || '',
      price: dish.price,
      category: dish.category,
      is_available: dish.is_available,
      allergens: dish.allergens?.join(', ') || '',
      preparation_time: dish.preparation_time || 0
    };
    editingDish = dish;
    showAddDishModal = true;
  }

  function closeModal() {
    showAddDishModal = false;
    editingDish = null;
  }

  async function saveDish() {
    try {
      const api = getApiClient();
      
      if (editingDish) {
        // Update existing dish
        await api.api.v1DishesUpdate(editingDish.id, {
          name: newDishForm.name,
          description: newDishForm.description || undefined,
          price: newDishForm.price,
          category: newDishForm.category
        });
      } else {
        // Create new dish
        await api.api.v1DishesCreate({
          name: newDishForm.name,
          description: newDishForm.description || undefined,
          price: newDishForm.price,
          category: newDishForm.category,
          restaurant_id: restaurantId
        });
      }
      
      closeModal();
      await loadDishes();
      
    } catch (err) {
      error = handleApiError(err);
    }
  }

  async function toggleDishAvailability(dish: Dish) {
    try {
      const api = getApiClient();
      
      // Update dish availability via API
      await api.api.v1DishesUpdate(dish.id, {
        is_available: !dish.is_available
      });
      
      // Update local state
      categories = categories.map(category => ({
        ...category,
        dishes: category.dishes.map(d => 
          d.id === dish.id 
            ? { ...d, is_available: !d.is_available, updated_at: new Date().toISOString() }
            : d
        )
      }));
      
    } catch (err) {
      error = handleApiError(err);
    }
  }

  async function deleteDish(dish: Dish) {
    if (!confirm(`Are you sure you want to delete "${dish.name}"? This action cannot be undone.`)) {
      return;
    }

    try {
      const api = getApiClient();
      
      // Delete dish via API
      await api.api.v1DishesDelete(dish.id);
      
      // Update local state
      categories = categories.map(category => ({
        ...category,
        dishes: category.dishes.filter(d => d.id !== dish.id)
      })).filter(category => category.dishes.length > 0);
      
    } catch (err) {
      error = handleApiError(err);
    }
  }

  function formatPrice(price: number): string {
    return `$${price.toFixed(2)}`;
  }
</script>

<svelte:head>
  <title>Dishes - {restaurantName} - LeCritique</title>
  <meta name="description" content="Manage restaurant dishes and menu items" />
</svelte:head>

<div class="min-h-screen bg-gray-50">
  <!-- Header -->
  <div class="bg-white shadow-sm border-b">
    <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
      <div class="flex justify-between items-center py-6">
        <div>
          <h1 class="text-2xl font-bold text-gray-900">Dishes</h1>
          <p class="text-gray-600">{restaurantName} - Menu Management</p>
        </div>
        <div class="flex space-x-3">
          <Button variant="outline" on:click={() => goto('/restaurants')}>
            <svg class="h-4 w-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 19l-7-7 7-7" />
            </svg>
            Back to Restaurants
          </Button>
          <Button on:click={openAddDishModal}>
            <svg class="h-4 w-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4" />
            </svg>
            Add Dish
          </Button>
        </div>
      </div>
    </div>
  </div>

  <div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
    {#if loading}
      <!-- Loading State -->
      <div class="space-y-8">
        {#each Array(3) as _}
          <div class="space-y-4">
            <div class="h-6 bg-gray-200 rounded w-1/4 animate-pulse"></div>
            <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
              {#each Array(3) as _}
                <Card>
                  <div class="animate-pulse space-y-3">
                    <div class="h-5 bg-gray-200 rounded w-3/4"></div>
                    <div class="h-4 bg-gray-200 rounded w-full"></div>
                    <div class="h-4 bg-gray-200 rounded w-1/2"></div>
                  </div>
                </Card>
              {/each}
            </div>
          </div>
        {/each}
      </div>

    {:else if error}
      <!-- Error State -->
      <Card>
        <div class="text-center py-12">
          <svg class="h-12 w-12 text-red-500 mx-auto mb-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-.833-1.964-.833-2.732 0L3.732 16.5c-.77.833.192 3 1.732 3z" />
          </svg>
          <h3 class="text-lg font-medium text-gray-900 mb-2">Failed to load dishes</h3>
          <p class="text-gray-600 mb-4">{error}</p>
          <Button on:click={loadDishes}>Try Again</Button>
        </div>
      </Card>

    {:else if categories.length === 0}
      <!-- Empty State -->
      <Card>
        <div class="text-center py-16">
          <svg class="h-16 w-16 text-gray-400 mx-auto mb-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5H7a2 2 0 00-2 2v10a2 2 0 002 2h8a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2" />
          </svg>
          <h3 class="text-xl font-medium text-gray-900 mb-2">No dishes yet</h3>
          <p class="text-gray-600 mb-6">Start building your menu by adding your first dish</p>
          <Button on:click={openAddDishModal}>
            <svg class="h-4 w-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4" />
            </svg>
            Add Your First Dish
          </Button>
        </div>
      </Card>

    {:else}
      <!-- Dishes by Category -->
      <div class="space-y-8">
        {#each categories as category}
          <div>
            <div class="flex items-center justify-between mb-4">
              <h2 class="text-xl font-semibold text-gray-900">{category.name}</h2>
              <span class="text-sm text-gray-500">{category.dishes.length} items</span>
            </div>

            <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
              {#each category.dishes as dish}
                <Card class="relative">
                  <!-- Availability Badge -->
                  <div class="absolute top-4 right-4">
                    <span class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium {dish.is_available ? 'bg-green-100 text-green-800' : 'bg-red-100 text-red-800'}">
                      {dish.is_available ? 'Available' : 'Unavailable'}
                    </span>
                  </div>

                  <div class="pr-20">
                    <h3 class="text-lg font-medium text-gray-900 mb-2">{dish.name}</h3>
                    
                    {#if dish.description}
                      <p class="text-sm text-gray-600 mb-3 line-clamp-2">{dish.description}</p>
                    {/if}

                    <div class="space-y-2 mb-4">
                      <!-- Price -->
                      <div class="flex items-center text-lg font-semibold text-green-600">
                        {formatPrice(dish.price)}
                      </div>

                      <!-- Preparation Time -->
                      {#if dish.preparation_time}
                        <div class="flex items-center text-sm text-gray-600">
                          <svg class="h-4 w-4 mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z" />
                          </svg>
                          {dish.preparation_time} min
                        </div>
                      {/if}

                      <!-- Allergens -->
                      {#if dish.allergens && dish.allergens.length > 0}
                        <div class="flex flex-wrap gap-1">
                          {#each dish.allergens as allergen}
                            <span class="inline-flex items-center px-2 py-0.5 rounded text-xs font-medium bg-yellow-100 text-yellow-800">
                              {allergen}
                            </span>
                          {/each}
                        </div>
                      {/if}
                    </div>

                    <!-- Action Buttons -->
                    <div class="flex flex-wrap gap-2">
                      <Button 
                        variant="outline" 
                        size="sm" 
                        on:click={() => openEditDishModal(dish)}
                      >
                        <svg class="h-3 w-3 mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z" />
                        </svg>
                        Edit
                      </Button>

                      <Button 
                        variant={dish.is_available ? 'secondary' : 'primary'} 
                        size="sm" 
                        on:click={() => toggleDishAvailability(dish)}
                      >
                        {dish.is_available ? 'Hide' : 'Show'}
                      </Button>

                      <Button 
                        variant="destructive" 
                        size="sm" 
                        on:click={() => deleteDish(dish)}
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
          </div>
        {/each}
      </div>
    {/if}
  </div>
</div>

<!-- Add/Edit Dish Modal -->
{#if showAddDishModal}
  <div class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center p-4 z-50">
    <div class="bg-white rounded-lg max-w-2xl w-full max-h-screen overflow-y-auto">
      <div class="p-6">
        <div class="flex justify-between items-center mb-6">
          <h3 class="text-lg font-medium text-gray-900">
            {editingDish ? 'Edit Dish' : 'Add New Dish'}
          </h3>
          <button on:click={closeModal} class="text-gray-400 hover:text-gray-600 cursor-pointer">
            <svg class="h-6 w-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
            </svg>
          </button>
        </div>

        <form on:submit|preventDefault={saveDish} class="space-y-4">
          <!-- Dish Name -->
          <div>
            <label for="dish-name" class="block text-sm font-medium text-gray-700 mb-1">
              Dish Name *
            </label>
            <Input
              id="dish-name"
              type="text"
              placeholder="Enter dish name"
              bind:value={newDishForm.name}
              required
            />
          </div>

          <!-- Description -->
          <div>
            <label for="dish-description" class="block text-sm font-medium text-gray-700 mb-1">
              Description
            </label>
            <textarea
              id="dish-description"
              rows="3"
              class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
              placeholder="Brief description of the dish"
              bind:value={newDishForm.description}
            ></textarea>
          </div>

          <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
            <!-- Price -->
            <div>
              <label for="dish-price" class="block text-sm font-medium text-gray-700 mb-1">
                Price *
              </label>
              <div class="relative">
                <span class="absolute left-3 top-2 text-gray-500">$</span>
                <input
                  id="dish-price"
                  type="number"
                  step="0.01"
                  min="0"
                  class="w-full pl-8 pr-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
                  placeholder="0.00"
                  bind:value={newDishForm.price}
                  required
                />
              </div>
            </div>

            <!-- Category -->
            <div>
              <label for="dish-category" class="block text-sm font-medium text-gray-700 mb-1">
                Category *
              </label>
              <select
                id="dish-category"
                class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500 cursor-pointer"
                bind:value={newDishForm.category}
                required
              >
                <option value="">Select category</option>
                {#each dishCategories as category}
                  <option value={category}>{category}</option>
                {/each}
              </select>
            </div>

            <!-- Preparation Time -->
            <div>
              <label for="dish-prep-time" class="block text-sm font-medium text-gray-700 mb-1">
                Preparation Time (minutes)
              </label>
              <input
                id="dish-prep-time"
                type="number"
                min="0"
                class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500 cursor-pointer"
                placeholder="0"
                bind:value={newDishForm.preparation_time}
              />
            </div>

            <!-- Availability -->
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-3">
                Availability
              </label>
              <label class="flex items-center cursor-pointer">
                <input
                  type="checkbox"
                  class="h-4 w-4 text-blue-600 focus:ring-blue-500 border-gray-300 rounded cursor-pointer"
                  bind:checked={newDishForm.is_available}
                />
                <span class="ml-2 text-sm text-gray-700">Available for ordering</span>
              </label>
            </div>
          </div>

          <!-- Allergens -->
          <div>
            <label for="dish-allergens" class="block text-sm font-medium text-gray-700 mb-1">
              Allergens
            </label>
            <Input
              id="dish-allergens"
              type="text"
              placeholder="e.g., gluten, dairy, nuts (comma-separated)"
              bind:value={newDishForm.allergens}
            />
            <p class="text-xs text-gray-500 mt-1">Separate multiple allergens with commas</p>
          </div>

          <!-- Form Actions -->
          <div class="flex justify-end space-x-3 pt-4 border-t border-gray-200">
            <Button variant="outline" type="button" on:click={closeModal}>
              Cancel
            </Button>
            <Button type="submit">
              {editingDish ? 'Save Changes' : 'Add Dish'}
            </Button>
          </div>
        </form>
      </div>
    </div>
  </div>
{/if}