<script lang="ts">
  import { Modal, Button, Input } from '$lib/components/ui';
  import { getApiClient, handleApiError } from '$lib/api/client';
  import { Loader2 } from 'lucide-svelte';
  import { toast } from 'svelte-sonner';

  interface RestaurantForm {
    name: string;
    description: string;
    address: string;
    phone: string;
    email: string;
    website: string;
    cuisine_type: string;
  }

  let {
    restaurant,
    onclose,
    onupdated
  }: {
    restaurant: any;
    onclose: () => void;
    onupdated: () => void;
  } = $props();

  let saving = false;
  let error = '';
  let formData: RestaurantForm = {
    name: restaurant?.name || '',
    description: restaurant?.description || '',
    address: '', // Note: address would come from locations array
    phone: restaurant?.phone || '',
    email: restaurant?.email || '',
    website: restaurant?.website || '',
    cuisine_type: '' // Note: cuisine_type not in API model
  };

  const cuisineTypes = [
    'American',
    'Italian',
    'Mexican',
    'Chinese',
    'Japanese',
    'Thai',
    'Indian',
    'French',
    'Mediterranean',
    'Greek',
    'Spanish',
    'Korean',
    'Vietnamese',
    'Middle Eastern',
    'Brazilian',
    'Seafood',
    'Steakhouse',
    'BBQ',
    'Pizza',
    'Fast Food',
    'Cafe',
    'Bakery',
    'Vegetarian',
    'Vegan',
    'Farm to Table',
    'Fusion',
    'Other'
  ];

  function validateForm(): boolean {
    error = '';

    if (!formData.name.trim()) {
      error = 'Restaurant name is required';
      return false;
    }

    if (formData.email && !isValidEmail(formData.email)) {
      error = 'Please enter a valid email address';
      return false;
    }

    if (formData.website && !isValidUrl(formData.website)) {
      error = 'Please enter a valid website URL';
      return false;
    }

    return true;
  }

  function isValidEmail(email: string): boolean {
    const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
    return emailRegex.test(email);
  }

  function isValidUrl(url: string): boolean {
    try {
      new URL(url);
      return true;
    } catch {
      return false;
    }
  }

  async function handleSubmit() {
    if (!validateForm()) {
      return;
    }

    saving = true;
    error = '';

    try {
      const api = getApiClient();
      
      // Update restaurant via API
      const response = await api.api.v1RestaurantsUpdate(restaurant.id, {
        name: formData.name,
        description: formData.description || undefined,
        email: formData.email || undefined,
        phone: formData.phone || undefined,
        website: formData.website || undefined
      });
      
      if (response.data.success) {
        toast.success('Restaurant updated successfully');
        onupdated();
      } else {
        error = 'Failed to update restaurant';
      }
      
    } catch (err) {
      error = handleApiError(err);
    } finally {
      saving = false;
    }
  }

  function handleClose() {
    onclose();
  }
</script>

<Modal 
  isOpen={true}
  title="Edit Restaurant"
  size="lg"
  onclose={handleClose}
>
  <form on:submit|preventDefault={handleSubmit} class="space-y-6">
    <!-- Basic Information -->
    <div>
      <h3 class="text-lg font-medium text-gray-900 mb-4">Basic Information</h3>
      
      <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
        <!-- Restaurant Name -->
        <div class="md:col-span-2">
          <label for="name" class="block text-sm font-medium text-gray-700 mb-2">
            Restaurant Name <span class="text-red-500">*</span>
          </label>
          <Input
            id="name"
            type="text"
            placeholder="Enter restaurant name"
            bind:value={formData.name}
            required
          />
        </div>

        <!-- Description -->
        <div class="md:col-span-2">
          <label for="description" class="block text-sm font-medium text-gray-700 mb-2">
            Description
          </label>
          <textarea
            id="description"
            rows="3"
            class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
            placeholder="Brief description of your restaurant"
            bind:value={formData.description}
          ></textarea>
        </div>

        <!-- Cuisine Type -->
        <div>
          <label for="cuisine_type" class="block text-sm font-medium text-gray-700 mb-2">
            Cuisine Type
          </label>
          <select
            id="cuisine_type"
            class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
            bind:value={formData.cuisine_type}
          >
            <option value="">Select cuisine type</option>
            {#each cuisineTypes as cuisine}
              <option value={cuisine}>{cuisine}</option>
            {/each}
          </select>
        </div>
      </div>
    </div>

    <!-- Contact Information -->
    <div>
      <h3 class="text-lg font-medium text-gray-900 mb-4">Contact Information</h3>
      
      <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
        <!-- Phone -->
        <div>
          <label for="phone" class="block text-sm font-medium text-gray-700 mb-2">
            Phone Number
          </label>
          <Input
            id="phone"
            type="tel"
            placeholder="+1 (555) 123-4567"
            bind:value={formData.phone}
          />
        </div>

        <!-- Email -->
        <div>
          <label for="email" class="block text-sm font-medium text-gray-700 mb-2">
            Email Address
          </label>
          <Input
            id="email"
            type="email"
            placeholder="restaurant@example.com"
            bind:value={formData.email}
          />
        </div>

        <!-- Website -->
        <div class="md:col-span-2">
          <label for="website" class="block text-sm font-medium text-gray-700 mb-2">
            Website
          </label>
          <Input
            id="website"
            type="url"
            placeholder="https://yourrestaurant.com"
            bind:value={formData.website}
          />
        </div>
      </div>
    </div>

    <!-- Error Display -->
    {#if error}
      <div class="bg-red-50 border border-red-200 rounded-md p-4">
        <div class="flex">
          <div class="flex-shrink-0">
            <svg class="h-5 w-5 text-red-400" fill="currentColor" viewBox="0 0 20 20">
              <path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zM8.707 7.293a1 1 0 00-1.414 1.414L8.586 10l-1.293 1.293a1 1 0 101.414 1.414L10 11.414l1.293 1.293a1 1 0 001.414-1.414L11.414 10l1.293-1.293a1 1 0 00-1.414-1.414L10 8.586 8.707 7.293z" clip-rule="evenodd" />
            </svg>
          </div>
          <div class="ml-3">
            <p class="text-sm text-red-800">{error}</p>
          </div>
        </div>
      </div>
    {/if}

    <!-- Form Actions -->
    <div class="mt-6 pt-6 border-t border-gray-200 flex justify-end space-x-3">
      <Button onclick={handleClose} variant="outline">
        Cancel
      </Button>
      <Button type="submit" disabled={saving} variant="gradient">
        {#if saving}
          <Loader2 class="w-4 h-4 mr-2 animate-spin" />
          Saving...
        {:else}
          Save Changes
        {/if}
      </Button>
    </div>
  </form>
</Modal>