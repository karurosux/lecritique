<script lang="ts">
  import { Card, Button, Input } from '$lib/components/ui';
  import { getApiClient, handleApiError } from '$lib/api/client';
  import { auth } from '$lib/stores/auth';
  import { goto } from '$app/navigation';
  import { onMount } from 'svelte';

  interface OrganizationForm {
    name: string;
    description: string;
    address: string;
    phone: string;
    email: string;
    website: string;
    cuisine_type: string;
  }

  let loading = false;
  let error = '';
  let formData: OrganizationForm = {
    name: '',
    description: '',
    address: '',
    phone: '',
    email: '',
    website: '',
    cuisine_type: '',
  };

  $: authState = $auth;

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
    'Other',
  ];

  onMount(() => {
    if (!authState.isAuthenticated) {
      goto('/login');
      return;
    }
  });

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

    loading = true;
    error = '';

    try {
      const api = getApiClient();

      // Create organization via API
      const response = await api.api.v1OrganizationsCreate({
        name: formData.name,
        description: formData.description || undefined,
        email: formData.email || undefined,
        phone: formData.phone || undefined,
        website: formData.website || undefined,
      });

      if (response.data.success) {
        // Redirect to organizations list on success
        goto('/organizations');
      } else {
        error = 'Failed to create organization';
      }
    } catch (err) {
      error = handleApiError(err);
    } finally {
      loading = false;
    }
  }

  function handleCancel() {
    goto('/organizations');
  }
</script>

<svelte:head>
  <title>Add Organization - Kyooar</title>
  <meta name="description" content="Add a new organization" />
</svelte:head>

<div class="min-h-screen bg-gray-50">
  <!-- Header -->
  <div class="bg-white shadow-sm border-b">
    <div class="max-w-3xl mx-auto px-4 sm:px-6 lg:px-8">
      <div class="flex justify-between items-center py-6">
        <div>
          <h1 class="text-2xl font-bold text-gray-900">Add Organization</h1>
          <p class="text-gray-600">Create a new organization profile</p>
        </div>
        <Button variant="outline" on:click={handleCancel}>
          <svg
            class="h-4 w-4 mr-2"
            fill="none"
            stroke="currentColor"
            viewBox="0 0 24 24">
            <path
              stroke-linecap="round"
              stroke-linejoin="round"
              stroke-width="2"
              d="M6 18L18 6M6 6l12 12" />
          </svg>
          Cancel
        </Button>
      </div>
    </div>
  </div>

  <div class="max-w-3xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
    <form on:submit|preventDefault={handleSubmit}>
      <Card>
        <div class="space-y-6">
          <!-- Basic Information -->
          <div>
            <h3 class="text-lg font-medium text-gray-900 mb-4">
              Basic Information
            </h3>

            <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
              <!-- Organization Name -->
              <div class="md:col-span-2">
                <label
                  for="name"
                  class="block text-sm font-medium text-gray-700 mb-2">
                  Organization Name *
                </label>
                <Input
                  id="name"
                  type="text"
                  placeholder="Enter organization name"
                  bind:value={formData.name}
                  required />
              </div>

              <!-- Description -->
              <div class="md:col-span-2">
                <label
                  for="description"
                  class="block text-sm font-medium text-gray-700 mb-2">
                  Description
                </label>
                <textarea
                  id="description"
                  rows="3"
                  class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
                  placeholder="Brief description of your organization"
                  bind:value={formData.description}></textarea>
              </div>

              <!-- Cuisine Type -->
              <div>
                <label
                  for="cuisine_type"
                  class="block text-sm font-medium text-gray-700 mb-2">
                  Cuisine Type
                </label>
                <select
                  id="cuisine_type"
                  class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500 cursor-pointer"
                  bind:value={formData.cuisine_type}>
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
            <h3 class="text-lg font-medium text-gray-900 mb-4">
              Contact Information
            </h3>

            <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
              <!-- Address -->
              <div class="md:col-span-2">
                <label
                  for="address"
                  class="block text-sm font-medium text-gray-700 mb-2">
                  Address *
                </label>
                <Input
                  id="address"
                  type="text"
                  placeholder="Enter full address"
                  bind:value={formData.address}
                  required />
              </div>

              <!-- Phone -->
              <div>
                <label
                  for="phone"
                  class="block text-sm font-medium text-gray-700 mb-2">
                  Phone Number
                </label>
                <Input
                  id="phone"
                  type="tel"
                  placeholder="+1 (555) 123-4567"
                  bind:value={formData.phone} />
              </div>

              <!-- Email -->
              <div>
                <label
                  for="email"
                  class="block text-sm font-medium text-gray-700 mb-2">
                  Email Address
                </label>
                <Input
                  id="email"
                  type="email"
                  placeholder="organization@example.com"
                  bind:value={formData.email} />
              </div>

              <!-- Website -->
              <div class="md:col-span-2">
                <label
                  for="website"
                  class="block text-sm font-medium text-gray-700 mb-2">
                  Website
                </label>
                <Input
                  id="website"
                  type="url"
                  placeholder="https://yourorganization.com"
                  bind:value={formData.website} />
              </div>
            </div>
          </div>

          <!-- Error Display -->
          {#if error}
            <div class="bg-red-50 border border-red-200 rounded-md p-4">
              <div class="flex">
                <div class="flex-shrink-0">
                  <svg
                    class="h-5 w-5 text-red-400"
                    fill="currentColor"
                    viewBox="0 0 20 20">
                    <path
                      fill-rule="evenodd"
                      d="M10 18a8 8 0 100-16 8 8 0 000 16zM8.707 7.293a1 1 0 00-1.414 1.414L8.586 10l-1.293 1.293a1 1 0 101.414 1.414L10 11.414l1.293 1.293a1 1 0 001.414-1.414L11.414 10l1.293-1.293a1 1 0 00-1.414-1.414L10 8.586 8.707 7.293z"
                      clip-rule="evenodd" />
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
            <Button type="submit" disabled={loading}>
              {#if loading}
                <svg
                  class="animate-spin -ml-1 mr-3 h-5 w-5 text-white inline"
                  fill="none"
                  viewBox="0 0 24 24">
                  <circle
                    class="opacity-25"
                    cx="12"
                    cy="12"
                    r="10"
                    stroke="currentColor"
                    stroke-width="4"></circle>
                  <path
                    class="opacity-75"
                    fill="currentColor"
                    d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"
                  ></path>
                </svg>
                Creating...
              {:else}
                Create Organization
              {/if}
            </Button>
          </div>
        </div>
      </Card>
    </form>
  </div>
</div>
