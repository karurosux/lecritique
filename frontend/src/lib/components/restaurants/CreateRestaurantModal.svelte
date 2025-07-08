<script lang="ts">
  import { Modal, Button, Input, Card } from '$lib/components/ui';
  import { getApiClient, handleApiError } from '$lib/api/client';

  let { 
    isOpen = $bindable(false),
    clickOrigin = null,
    onclose = () => {},
    onrestaurantcreated = (restaurant: any) => {}
  }: {
    isOpen?: boolean;
    clickOrigin?: { x: number; y: number } | null;
    onclose?: () => void;
    onrestaurantcreated?: (restaurant: any) => void;
  } = $props();

  let formData = $state({
    name: '',
    description: '',
    address: '',
    phone: '',
    email: '',
    website: '',
    cuisine_type: ''
  });

  let loading = $state(false);
  let error = $state('');

  const cuisineTypes = [
    'Italian', 'Mexican', 'Chinese', 'Japanese', 'Indian', 'French', 'Thai', 'Greek',
    'American', 'Mediterranean', 'Korean', 'Vietnamese', 'Spanish', 'Turkish', 'Lebanese', 'Other'
  ];

  function resetForm() {
    formData.name = '';
    formData.description = '';
    formData.address = '';
    formData.phone = '';
    formData.email = '';
    formData.website = '';
    formData.cuisine_type = '';
    error = '';
  }

  function handleClose() {
    if (!loading) {
      resetForm();
      isOpen = false;
      onclose();
    }
  }

  async function handleSubmit(event: Event) {
    event.preventDefault();
    if (!formData.name.trim() || !formData.address.trim()) {
      error = 'Restaurant name and address are required.';
      return;
    }

    loading = true;
    error = '';

    try {
      const api = getApiClient();
      
      const restaurantData = {
        name: formData.name.trim(),
        description: formData.description.trim() || undefined,
        address: formData.address.trim(),
        phone: formData.phone.trim() || undefined,
        email: formData.email.trim() || undefined,
        website: formData.website.trim() || undefined,
        cuisine_type: formData.cuisine_type || undefined,
        is_active: true
      };

      const response = await api.api.v1RestaurantsCreate(restaurantData);
      
      // Create restaurant object matching the expected interface
      const newRestaurant = {
        id: response.data.id,
        name: response.data.name,
        description: response.data.description || '',
        address: response.data.address,
        phone: response.data.phone || '',
        email: response.data.email || '',
        website: response.data.website || '',
        cuisine_type: response.data.cuisine_type || '',
        status: response.data.is_active ? 'active' : 'inactive',
        created_at: response.data.created_at || new Date().toISOString(),
        updated_at: response.data.updated_at || new Date().toISOString()
      };

      onrestaurantcreated(newRestaurant);
      resetForm();
      isOpen = false;
      onclose();
    } catch (err) {
      error = handleApiError(err);
    } finally {
      loading = false;
    }
  }

  function handleKeyDown(event: KeyboardEvent) {
    if (event.key === 'Enter' && (event.metaKey || event.ctrlKey)) {
      handleSubmit();
    }
  }

  // Bind keyboard shortcut on mount
  import { onMount, onDestroy } from 'svelte';
  import { browser } from '$app/environment';
  
  onMount(() => {
    if (browser) {
      window.addEventListener('keydown', handleKeyDown);
    }
  });
  
  onDestroy(() => {
    if (browser) {
      window.removeEventListener('keydown', handleKeyDown);
    }
  });
</script>

<Modal bind:isOpen title="Add New Restaurant" {clickOrigin} size="lg" onclose={handleClose}>
  <div class="space-y-6">
    <!-- Error Message -->
    {#if error}
      <Card variant="minimal" class="border-red-200 bg-red-50">
        <div class="flex items-center space-x-2">
          <svg class="h-5 w-5 text-red-500 flex-shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-.833-1.964-.833-2.732 0L3.732 16.5c-.77.833.192 3 1.732 3z" />
          </svg>
          <p class="text-red-700 text-sm">{error}</p>
        </div>
      </Card>
    {/if}

    <form onsubmit={handleSubmit}>
      <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
        <!-- Restaurant Name -->
        <div class="md:col-span-2">
          <label for="name" class="block text-sm font-medium text-gray-700 mb-2">
            Restaurant Name <span class="text-red-500">*</span>
          </label>
          <Input
            id="name"
            bind:value={formData.name}
            placeholder="Enter restaurant name"
            disabled={loading}
            required
            class="w-full"
          />
        </div>

        <!-- Description -->
        <div class="md:col-span-2">
          <label for="description" class="block text-sm font-medium text-gray-700 mb-2">
            Description
          </label>
          <textarea
            id="description"
            bind:value={formData.description}
            placeholder="Brief description of the restaurant"
            rows="3"
            disabled={loading}
            class="w-full px-4 py-3 border border-gray-200 rounded-xl bg-white focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500 transition-all duration-200 resize-none"
          ></textarea>
        </div>

        <!-- Address -->
        <div class="md:col-span-2">
          <label for="address" class="block text-sm font-medium text-gray-700 mb-2">
            Address <span class="text-red-500">*</span>
          </label>
          <Input
            id="address"
            bind:value={formData.address}
            placeholder="Restaurant address"
            disabled={loading}
            required
            class="w-full"
          />
        </div>

        <!-- Phone -->
        <div>
          <label for="phone" class="block text-sm font-medium text-gray-700 mb-2">
            Phone
          </label>
          <Input
            id="phone"
            bind:value={formData.phone}
            placeholder="Phone number"
            disabled={loading}
            class="w-full"
          />
        </div>

        <!-- Email -->
        <div>
          <label for="email" class="block text-sm font-medium text-gray-700 mb-2">
            Email
          </label>
          <Input
            id="email"
            type="email"
            bind:value={formData.email}
            placeholder="contact@restaurant.com"
            disabled={loading}
            class="w-full"
          />
        </div>

        <!-- Website -->
        <div>
          <label for="website" class="block text-sm font-medium text-gray-700 mb-2">
            Website
          </label>
          <Input
            id="website"
            type="url"
            bind:value={formData.website}
            placeholder="https://restaurant.com"
            disabled={loading}
            class="w-full"
          />
        </div>

        <!-- Cuisine Type -->
        <div>
          <label for="cuisine_type" class="block text-sm font-medium text-gray-700 mb-2">
            Cuisine Type
          </label>
          <div class="relative">
            <select
              id="cuisine_type"
              bind:value={formData.cuisine_type}
              disabled={loading}
              class="appearance-none w-full px-4 py-3 pr-10 border border-gray-200 rounded-xl bg-white focus:outline-none focus:ring-2 focus:ring-blue-500 cursor-pointer"
            >
              <option value="">Select cuisine type</option>
              {#each cuisineTypes as cuisine}
                <option value={cuisine}>{cuisine}</option>
              {/each}
            </select>
            <div class="absolute inset-y-0 right-0 flex items-center pr-3 pointer-events-none">
              <svg class="h-4 w-4 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 9l-7 7-7-7" />
              </svg>
            </div>
          </div>
        </div>
      </div>

      <!-- Form Actions -->
      <div class="flex items-center justify-end space-x-3 mt-8 pt-6 border-t border-gray-200">
        <Button
          type="button"
          variant="outline"
          onclick={handleClose}
          disabled={loading}
        >
          Cancel
        </Button>
        <Button
          type="submit"
          variant="gradient"
          {loading}
          disabled={loading || !formData.name.trim() || !formData.address.trim()}
          class="min-w-24"
        >
          {#if loading}
            <svg class="animate-spin h-4 w-4 mr-2" fill="none" viewBox="0 0 24 24">
              <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
              <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
            </svg>
            Creating...
          {:else}
            Create Restaurant
          {/if}
        </Button>
      </div>
    </form>
  </div>
</Modal>

<style>
  /* Custom scrollbar for textarea */
  textarea::-webkit-scrollbar {
    width: 8px;
  }
  
  textarea::-webkit-scrollbar-track {
    background: #f1f1f1;
    border-radius: 4px;
  }
  
  textarea::-webkit-scrollbar-thumb {
    background: #c1c1c1;
    border-radius: 4px;
  }
  
  textarea::-webkit-scrollbar-thumb:hover {
    background: #a8a8a8;
  }
</style>