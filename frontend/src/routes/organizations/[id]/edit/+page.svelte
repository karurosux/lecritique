<script lang="ts">
  import { Card, Button, Input } from '$lib/components/ui';
  import { getApiClient, handleApiError } from '$lib/api/client';
  import { auth } from '$lib/stores/auth';
  import { goto } from '$app/navigation';
  import { onMount } from 'svelte';
  import { page } from '$app/stores';

  interface OrganizationForm {
    name: string;
    description: string;
    address: string;
    phone: string;
    email: string;
    website: string;
    cuisine_type: string;
  }

  let loading = true;
  let saving = false;
  let error = '';
  let organizationId = '';
  let formData: OrganizationForm = {
    name: '',
    description: '',
    address: '',
    phone: '',
    email: '',
    website: '',
    cuisine_type: ''
  };

  $: authState = $auth;
  $: organizationId = $page.params.id;

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

  onMount(async () => {
    if (!authState.isAuthenticated) {
      goto('/login');
      return;
    }

    await loadOrganization();
  });

  async function loadOrganization() {
    loading = true;
    error = '';

    try {
      const api = getApiClient();
      
      // Get organization details via API
      const response = await api.api.v1OrganizationsDetail(organizationId);
      
      if (response.data.success && response.data.data) {
        const organization = response.data.data;
        formData = {
          name: organization.name || '',
          description: organization.description || '',
          address: '', // Note: address would come from locations array
          phone: organization.phone || '',
          email: organization.email || '',
          website: organization.website || '',
          cuisine_type: '' // Note: cuisine_type not in API model
        };
      } else {
        error = 'Organization not found';
      }

    } catch (err) {
      error = handleApiError(err);
    } finally {
      loading = false;
    }
  }

  function validateForm(): boolean {
    error = '';

    if (!formData.name.trim()) {
      error = 'Organization name is required';
      return false;
    }

    if (!formData.address.trim()) {
      error = 'Address is required';
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
      
      // Update organization via API
      const response = await api.api.v1OrganizationsUpdate(organizationId, {
        name: formData.name,
        description: formData.description || undefined,
        email: formData.email || undefined,
        phone: formData.phone || undefined,
        website: formData.website || undefined
      });
      
      if (response.data.success) {
        // Redirect to organizations list on success
        goto('/organizations');
      } else {
        error = 'Failed to update organization';
      }
      
    } catch (err) {
      error = handleApiError(err);
    } finally {
      saving = false;
    }
  }

  function handleCancel() {
    goto('/organizations');
  }
</script>

<svelte:head>
  <title>Edit Organization - LeCritique</title>
  <meta name="description" content="Edit organization information" />
</svelte:head>

<div class="min-h-screen bg-gray-50">
  <!-- Header -->
  <div class="bg-white shadow-sm border-b">
    <div class="max-w-3xl mx-auto px-4 sm:px-6 lg:px-8">
      <div class="flex justify-between items-center py-6">
        <div>
          <h1 class="text-2xl font-bold text-gray-900">Edit Organization</h1>
          <p class="text-gray-600">Update organization information</p>
        </div>
        <Button variant="outline" on:click={handleCancel}>
          <svg class="h-4 w-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
          </svg>
          Cancel
        </Button>
      </div>
    </div>
  </div>

  <div class="max-w-3xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
    {#if loading}
      <!-- Loading State -->
      <Card>
        <div class="animate-pulse space-y-6">
          <div class="h-6 bg-gray-200 rounded w-1/4"></div>
          <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
            <div class="md:col-span-2">
              <div class="h-4 bg-gray-200 rounded w-1/6 mb-2"></div>
              <div class="h-10 bg-gray-200 rounded"></div>
            </div>
            <div>
              <div class="h-4 bg-gray-200 rounded w-1/4 mb-2"></div>
              <div class="h-10 bg-gray-200 rounded"></div>
            </div>
            <div>
              <div class="h-4 bg-gray-200 rounded w-1/4 mb-2"></div>
              <div class="h-10 bg-gray-200 rounded"></div>
            </div>
          </div>
          <div class="h-10 bg-gray-200 rounded w-1/4 ml-auto"></div>
        </div>
      </Card>

    {:else if error && !formData.name}
      <!-- Error State -->
      <Card>
        <div class="text-center py-12">
          <svg class="h-12 w-12 text-red-500 mx-auto mb-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-.833-1.964-.833-2.732 0L3.732 16.5c-.77.833.192 3 1.732 3z" />
          </svg>
          <h3 class="text-lg font-medium text-gray-900 mb-2">Failed to load organization</h3>
          <p class="text-gray-600 mb-4">{error}</p>
          <div class="flex justify-center space-x-3">
            <Button on:click={loadOrganization}>Try Again</Button>
            <Button variant="outline" on:click={handleCancel}>Go Back</Button>
          </div>
        </div>
      </Card>

    {:else}
      <!-- Edit Form -->
      <form on:submit|preventDefault={handleSubmit}>
        <Card>
          <div class="space-y-6">
            <!-- Basic Information -->
            <div>
              <h3 class="text-lg font-medium text-gray-900 mb-4">Basic Information</h3>
              
              <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
                <!-- Organization Name -->
                <div class="md:col-span-2">
                  <label for="name" class="block text-sm font-medium text-gray-700 mb-2">
                    Organization Name *
                  </label>
                  <Input
                    id="name"
                    type="text"
                    placeholder="Enter organization name"
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
                    placeholder="Brief description of your organization"
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
                <!-- Address -->
                <div class="md:col-span-2">
                  <label for="address" class="block text-sm font-medium text-gray-700 mb-2">
                    Address *
                  </label>
                  <Input
                    id="address"
                    type="text"
                    placeholder="Enter full address"
                    bind:value={formData.address}
                    required
                  />
                </div>

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
                    placeholder="organization@example.com"
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
                    placeholder="https://yourorganization.com"
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
            <div class="flex justify-end space-x-3 pt-6 border-t border-gray-200">
              <Button variant="outline" type="button" on:click={handleCancel}>
                Cancel
              </Button>
              <Button type="submit" disabled={saving}>
                {#if saving}
                  <svg class="animate-spin -ml-1 mr-3 h-5 w-5 text-white inline" fill="none" viewBox="0 0 24 24">
                    <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
                    <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
                  </svg>
                  Saving...
                {:else}
                  Save Changes
                {/if}
              </Button>
            </div>
          </div>
        </Card>
      </form>
    {/if}
  </div>
</div>
