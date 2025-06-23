<script lang="ts">
  import { getApiClient, handleApiError } from '$lib/api/client';
  import { auth } from '$lib/stores/auth';
  import { goto } from '$app/navigation';
  import { page } from '$app/stores';
  import DishesHeader from '$lib/components/dishes/DishesHeader.svelte';
  import DishesSearchAndFilters from '$lib/components/dishes/DishesSearchAndFilters.svelte';
  import DishesList from '$lib/components/dishes/DishesList.svelte';
  import AddDishModal from '$lib/components/dishes/AddDishModal.svelte';

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

  let loading = $state(true);
  let error = $state('');
  let restaurantName = $state('');
  let dishes = $state<Dish[]>([]);
  let categories = $state<Category[]>([]);
  let showAddDishModal = $state(false);
  let editingDish = $state<Dish | null>(null);
  let modalLoading = $state(false);
  let modalError = $state('');
  let hasInitialized = $state(false);

  // Search and filter state
  let searchQuery = $state('');
  let categoryFilter = $state('all');
  let availabilityFilter = $state('all');
  let sortBy = $state('name');

  let authState = $derived($auth);
  let restaurantId = $derived($page.params.id);
  
  // Filter and sort dishes using $state and $effect instead of $derived
  let filteredAndSortedDishes = $state<Dish[]>([]);
  let filteredCategories = $state<Category[]>([]);

  // Update filtered dishes when dependencies change
  $effect(() => {
    console.log('🔥 FILTERING EFFECT RUNNING!');
    console.log('Input dishes:', $state.snapshot(dishes));
    console.log('Filter criteria - search:', searchQuery, 'category:', categoryFilter, 'availability:', availabilityFilter, 'sort:', sortBy);
    
    let filtered = dishes.filter(dish => {
      const matchesSearch = dish.name.toLowerCase().includes(searchQuery.toLowerCase()) ||
                           dish.description?.toLowerCase().includes(searchQuery.toLowerCase());
      const matchesCategory = categoryFilter === 'all' || dish.category === categoryFilter;
      const matchesAvailability = availabilityFilter === 'all' || 
                                 (availabilityFilter === 'available' && dish.is_available) ||
                                 (availabilityFilter === 'unavailable' && !dish.is_available);
      
      console.log(`Dish ${dish.name}: search=${matchesSearch}, category=${matchesCategory}, availability=${matchesAvailability}`);
      return matchesSearch && matchesCategory && matchesAvailability;
    });

    console.log('Filtered dishes:', $state.snapshot(filtered));

    // Sort dishes
    filtered.sort((a, b) => {
      switch (sortBy) {
        case 'price':
          return a.price - b.price;
        case 'category':
          return a.category.localeCompare(b.category);
        case 'created_at':
          return new Date(b.created_at).getTime() - new Date(a.created_at).getTime();
        default:
          return a.name.localeCompare(b.name);
      }
    });

    filteredAndSortedDishes = filtered;
    
    // Update categories
    console.log('🔥 UPDATING CATEGORIES!');
    const categoryMap = new Map<string, Dish[]>();
    
    filtered.forEach(dish => {
      if (!categoryMap.has(dish.category)) {
        categoryMap.set(dish.category, []);
      }
      categoryMap.get(dish.category)!.push(dish);
    });

    const categoryResult = Array.from(categoryMap.entries()).map(([name, dishes]) => ({
      id: name.toLowerCase().replace(/\s+/g, '-'),
      name,
      dishes
    })).sort((a, b) => a.name.localeCompare(b.name));
    
    console.log('filteredCategories result:', $state.snapshot(categoryResult));
    filteredCategories = categoryResult;
  });

  // Handle authentication and initial load
  $effect(() => {
    console.log('Effect running - authState:', authState, 'hasInitialized:', hasInitialized);
    
    if (!authState.isAuthenticated) {
      console.log('Not authenticated, redirecting to login');
      goto('/login');
      return;
    }
    
    if (authState.isAuthenticated && !hasInitialized) {
      console.log('Starting to load dishes...');
      hasInitialized = true;
      loadDishes();
    }
  });

  async function loadDishes() {
    loading = true;
    error = '';

    try {
      const api = getApiClient();
      
      console.log('Loading dishes for restaurant:', restaurantId);
      
      const [restaurantResponse, dishesResponse] = await Promise.all([
        api.api.v1RestaurantsDetail(restaurantId),
        api.api.v1RestaurantsDishesList(restaurantId)
      ]);
      
      console.log('Restaurant response:', restaurantResponse);
      console.log('Dishes response:', dishesResponse);
      
      // Set restaurant name
      if (restaurantResponse.data.success && restaurantResponse.data.data) {
        restaurantName = restaurantResponse.data.data.name || 'Unknown Restaurant';
      } else {
        restaurantName = 'Unknown Restaurant';
        console.warn('Failed to load restaurant data:', restaurantResponse);
      }
      
      // Process dishes
      if (dishesResponse.data.success && dishesResponse.data.data) {
        console.log('Raw dishes data:', dishesResponse.data.data);
        console.log('Sample dish structure:', JSON.stringify(dishesResponse.data.data[0], null, 2));
        dishes = dishesResponse.data.data.map((dish: any) => ({
          id: dish.id || '',
          name: dish.name || '',
          description: dish.description || '',
          price: dish.price || 0,
          category: dish.category || 'Uncategorized',
          is_available: dish.is_available !== false,
          allergens: [],
          preparation_time: 0,
          created_at: dish.created_at || '',
          updated_at: dish.updated_at || ''
        }));
        console.log('Processed dishes:', dishes);
      } else {
        dishes = [];
        console.warn('No dishes found or API error:', dishesResponse);
        
        // For debugging: Always add test data to verify components work
        console.log('Adding test dishes for debugging...');
        dishes = [
            {
              id: 'test-1',
              name: 'Test Burger',
              description: 'A delicious test burger',
              price: 12.99,
              category: 'Main Courses',
              is_available: true,
              allergens: ['gluten'],
              preparation_time: 15,
              created_at: new Date().toISOString(),
              updated_at: new Date().toISOString()
            },
            {
              id: 'test-2',
              name: 'Test Salad',
              description: 'A fresh test salad',
              price: 8.99,
              category: 'Salads',
              is_available: true,
              allergens: [],
              preparation_time: 5,
              created_at: new Date().toISOString(),
              updated_at: new Date().toISOString()
            }
          ];
      }

      // Group dishes by category
      updateCategories();
      console.log('Categories after update:', $state.snapshot(categories));
      console.log('Final dishes array:', $state.snapshot(dishes));
      console.log('Final filteredCategories - will show after derived runs');
      
      console.log('Dishes loaded successfully!');

    } catch (err) {
      console.error('Error loading dishes:', err);
      error = handleApiError(err);
    } finally {
      loading = false;
      console.log('Loading finished. Final state - loading:', loading, 'dishes:', dishes.length, 'categories:', categories.length);
    }
  }

  function updateCategories() {
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
    })).sort((a, b) => a.name.localeCompare(b.name));
  }

  function handleAddDish() {
    editingDish = null;
    modalError = '';
    showAddDishModal = true;
  }

  function handleEditDish(dish: Dish) {
    editingDish = dish;
    modalError = '';
    showAddDishModal = true;
  }

  function handleCloseModal() {
    showAddDishModal = false;
    editingDish = null;
    modalError = '';
  }

  async function handleSaveDish(dishData: any) {
    modalLoading = true;
    modalError = '';

    try {
      const api = getApiClient();
      
      if (editingDish) {
        // Update existing dish
        await api.api.v1DishesUpdate(editingDish.id, {
          name: dishData.name,
          description: dishData.description || undefined,
          price: dishData.price,
          category: dishData.category
        });

        // Update local state
        dishes = dishes.map(dish => 
          dish.id === editingDish!.id 
            ? { ...dish, ...dishData, updated_at: new Date().toISOString() }
            : dish
        );
      } else {
        // Create new dish
        const response = await api.api.v1DishesCreate({
          name: dishData.name,
          description: dishData.description || undefined,
          price: dishData.price,
          category: dishData.category,
          restaurant_id: restaurantId
        });

        // Add to local state
        const newDish: Dish = {
          id: response.data.id || Date.now().toString(),
          name: dishData.name,
          description: dishData.description,
          price: dishData.price,
          category: dishData.category,
          is_available: dishData.is_available,
          allergens: dishData.allergens || [],
          preparation_time: dishData.preparation_time,
          created_at: new Date().toISOString(),
          updated_at: new Date().toISOString()
        };

        dishes = [newDish, ...dishes];
      }
      
      updateCategories();
      handleCloseModal();
      
    } catch (err) {
      modalError = handleApiError(err);
    } finally {
      modalLoading = false;
    }
  }

  async function handleToggleAvailability(dish: Dish) {
    try {
      const api = getApiClient();
      
      // Note: API might not support is_available field, 
      // so this call might fail - we'll update locally regardless
      try {
        await api.api.v1DishesUpdate(dish.id, {
          is_available: !dish.is_available
        });
      } catch (updateError) {
        console.warn('API does not support is_available field update:', updateError);
      }
      
      // Update local state
      dishes = dishes.map(d => 
        d.id === dish.id 
          ? { ...d, is_available: !d.is_available, updated_at: new Date().toISOString() }
          : d
      );
      
      updateCategories();
      
    } catch (err) {
      error = handleApiError(err);
    }
  }

  async function handleDeleteDish(dish: Dish) {
    if (!confirm(`Are you sure you want to delete "${dish.name}"? This action cannot be undone.`)) {
      return;
    }

    try {
      const api = getApiClient();
      
      await api.api.v1DishesDelete(dish.id);
      
      // Update local state
      dishes = dishes.filter(d => d.id !== dish.id);
      updateCategories();
      
    } catch (err) {
      error = handleApiError(err);
    }
  }

  function handleRetry() {
    loadDishes();
  }
</script>

<svelte:head>
  <title>Menu Management - {restaurantName} - LeCritique</title>
  <meta name="description" content="Manage restaurant dishes and menu items for {restaurantName}" />
</svelte:head>

<div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
  <!-- Dishes Header -->
  <DishesHeader
    {restaurantName}
    {restaurantId}
    {dishes}
    {loading}
    onadddish={handleAddDish}
  />

  <!-- Search and Filters -->
  <DishesSearchAndFilters
    bind:searchQuery
    bind:categoryFilter
    bind:availabilityFilter
    bind:sortBy
    {categories}
    totalDishes={dishes.length}
    filteredCount={filteredAndSortedDishes.length}
  />

  <!-- Dishes List -->
  <DishesList
    dishes={filteredAndSortedDishes}
    {loading}
    {error}
    onadddish={handleAddDish}
    oneditdish={handleEditDish}
    ontoggleavailability={handleToggleAvailability}
    ondeletedish={handleDeleteDish}
    onretry={handleRetry}
  />
</div>

<!-- Add/Edit Dish Modal -->
<AddDishModal
  bind:isOpen={showAddDishModal}
  {editingDish}
  loading={modalLoading}
  error={modalError}
  onclose={handleCloseModal}
  onsave={handleSaveDish}
/>