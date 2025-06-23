<script lang="ts">
  import { Modal, Button, Input, Card } from '$lib/components/ui';

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

  let { 
    isOpen = $bindable(false),
    editingDish = null,
    loading = false,
    error = '',
    onclose = () => {},
    onsave = (dishData: any) => {}
  }: {
    isOpen?: boolean;
    editingDish?: Dish | null;
    loading?: boolean;
    error?: string;
    onclose?: () => void;
    onsave?: (dishData: any) => void;
  } = $props();

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

  let formData = $state({
    name: '',
    description: '',
    price: 0,
    category: '',
    is_available: true,
    allergens: '',
    preparation_time: 0
  });

  // Reset form when modal opens/closes or when editing dish changes
  $effect(() => {
    if (isOpen) {
      if (editingDish) {
        formData.name = editingDish.name;
        formData.description = editingDish.description || '';
        formData.price = editingDish.price;
        formData.category = editingDish.category;
        formData.is_available = editingDish.is_available;
        formData.allergens = editingDish.allergens?.join(', ') || '';
        formData.preparation_time = editingDish.preparation_time || 0;
      } else {
        formData.name = '';
        formData.description = '';
        formData.price = 0;
        formData.category = '';
        formData.is_available = true;
        formData.allergens = '';
        formData.preparation_time = 0;
      }
    }
  });

  function handleClose() {
    if (!loading) {
      onclose();
    }
  }

  function handleSubmit(event: Event) {
    event.preventDefault();
    
    if (!formData.name.trim() || !formData.category || formData.price <= 0) {
      return;
    }

    const dishData = {
      name: formData.name.trim(),
      description: formData.description.trim() || undefined,
      price: formData.price,
      category: formData.category,
      is_available: formData.is_available,
      allergens: formData.allergens ? formData.allergens.split(',').map(a => a.trim()).filter(Boolean) : [],
      preparation_time: formData.preparation_time || undefined
    };

    onsave(dishData);
  }

  let isFormValid = $derived(
    formData.name.trim() !== '' && 
    formData.category !== '' && 
    formData.price > 0
  );
</script>

<Modal bind:isOpen title={editingDish ? 'Edit Dish' : 'Add New Dish'} size="lg" onclose={handleClose}>
  <form onsubmit={handleSubmit} class="space-y-6">
    <!-- Error Message -->
    {#if error}
      <Card variant="minimal" class="border-red-200 bg-red-50 p-4">
        <div class="flex items-center space-x-3">
          <svg class="h-5 w-5 text-red-500 flex-shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-.833-1.964-.833-2.732 0L3.732 16.5c-.77.833.192 3 1.732 3z" />
          </svg>
          <p class="text-red-700 text-sm font-medium">{error}</p>
        </div>
      </Card>
    {/if}

    <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
      <!-- Dish Name -->
      <div class="md:col-span-2">
        <label for="dish-name" class="block text-sm font-semibold text-gray-700 mb-2">
          Dish Name <span class="text-red-500">*</span>
        </label>
        <Input
          id="dish-name"
          type="text"
          placeholder="Enter dish name"
          bind:value={formData.name}
          disabled={loading}
          required
          variant="default"
          class="w-full"
        />
      </div>

      <!-- Description -->
      <div class="md:col-span-2">
        <label for="dish-description" class="block text-sm font-semibold text-gray-700 mb-2">
          Description
        </label>
        <textarea
          id="dish-description"
          rows="3"
          class="w-full px-4 py-3 border-2 border-gray-200 rounded-xl bg-white focus:outline-none focus:border-blue-500 focus:ring-4 focus:ring-blue-500/20 transition-all duration-300 resize-none disabled:opacity-50"
          placeholder="Brief description of the dish..."
          bind:value={formData.description}
          disabled={loading}
        ></textarea>
      </div>

      <!-- Price -->
      <div>
        <label for="dish-price" class="block text-sm font-semibold text-gray-700 mb-2">
          Price <span class="text-red-500">*</span>
        </label>
        <div class="relative">
          <span class="absolute left-4 top-1/2 transform -translate-y-1/2 text-gray-500 text-lg font-medium">$</span>
          <input
            id="dish-price"
            type="number"
            step="0.01"
            min="0"
            class="w-full pl-8 pr-4 py-3 border-2 border-gray-200 rounded-xl bg-white focus:outline-none focus:border-blue-500 focus:ring-4 focus:ring-blue-500/20 transition-all duration-300 disabled:opacity-50"
            placeholder="0.00"
            bind:value={formData.price}
            required
            disabled={loading}
          />
        </div>
      </div>

      <!-- Category -->
      <div>
        <label for="dish-category" class="block text-sm font-semibold text-gray-700 mb-2">
          Category <span class="text-red-500">*</span>
        </label>
        <select
          id="dish-category"
          class="w-full px-4 py-3 border-2 border-gray-200 rounded-xl bg-white focus:outline-none focus:border-blue-500 focus:ring-4 focus:ring-blue-500/20 transition-all duration-300 cursor-pointer disabled:opacity-50"
          bind:value={formData.category}
          required
          disabled={loading}
        >
          <option value="">Select category</option>
          {#each dishCategories as category}
            <option value={category}>{category}</option>
          {/each}
        </select>
      </div>

      <!-- Preparation Time -->
      <div>
        <label for="dish-prep-time" class="block text-sm font-semibold text-gray-700 mb-2">
          Preparation Time
        </label>
        <div class="relative">
          <input
            id="dish-prep-time"
            type="number"
            min="0"
            class="w-full px-4 py-3 pr-16 border-2 border-gray-200 rounded-xl bg-white focus:outline-none focus:border-blue-500 focus:ring-4 focus:ring-blue-500/20 transition-all duration-300 disabled:opacity-50"
            placeholder="0"
            bind:value={formData.preparation_time}
            disabled={loading}
          />
          <span class="absolute right-4 top-1/2 transform -translate-y-1/2 text-gray-500 text-sm">minutes</span>
        </div>
      </div>

      <!-- Availability -->
      <div>
        <label class="block text-sm font-semibold text-gray-700 mb-3">
          Availability
        </label>
        <label class="flex items-center space-x-3 cursor-pointer group">
          <input
            type="checkbox"
            class="h-5 w-5 text-blue-600 focus:ring-blue-500 border-2 border-gray-300 rounded cursor-pointer transition-all duration-200 disabled:opacity-50"
            bind:checked={formData.is_available}
            disabled={loading}
          />
          <span class="text-sm text-gray-700 group-hover:text-gray-900 transition-colors">Available for ordering</span>
        </label>
      </div>
    </div>

    <!-- Allergens -->
    <div>
      <label for="dish-allergens" class="block text-sm font-semibold text-gray-700 mb-2">
        Allergens
      </label>
      <Input
        id="dish-allergens"
        type="text"
        placeholder="e.g., gluten, dairy, nuts (comma-separated)"
        bind:value={formData.allergens}
        disabled={loading}
        variant="default"
        class="w-full"
      />
      <p class="text-xs text-gray-500 mt-2 flex items-center">
        <svg class="h-3 w-3 mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 16h-1v-4h-1m1-4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
        </svg>
        Separate multiple allergens with commas
      </p>
    </div>

    <!-- Form Actions -->
    <div class="flex items-center justify-end space-x-4 pt-6 border-t border-gray-200">
      <Button 
        variant="outline" 
        type="button" 
        onclick={handleClose}
        disabled={loading}
      >
        Cancel
      </Button>
      <Button 
        type="submit"
        variant="gradient"
        disabled={loading || !isFormValid}
        class="min-w-32 shadow-lg"
      >
        {#if loading}
          <svg class="animate-spin h-4 w-4 mr-2" fill="none" viewBox="0 0 24 24">
            <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
            <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 714 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
          </svg>
          {editingDish ? 'Saving...' : 'Creating...'}
        {:else}
          {editingDish ? 'Save Changes' : 'Create Dish'}
        {/if}
      </Button>
    </div>
  </form>
</Modal>