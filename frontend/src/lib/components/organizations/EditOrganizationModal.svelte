<script lang="ts">
  import { Modal, Button, Input, Card } from '$lib/components/ui';
  import { AlertTriangle } from 'lucide-svelte';
  import { getApiClient, handleApiError } from '$lib/api/client';
  import { toast } from 'svelte-sonner';

  let {
    isOpen = $bindable(true),
    organization,
    clickOrigin = null,
    onclose = () => {},
    onupdated = () => {},
  }: {
    isOpen?: boolean;
    organization: any;
    clickOrigin?: { x: number; y: number } | null;
    onclose?: () => void;
    onupdated?: () => void;
  } = $props();

  let formData = $state({
    name: organization?.name || '',
    description: organization?.description || '',
    address: organization?.address || '',
    phone: organization?.phone || '',
    email: organization?.email || '',
    website: organization?.website || '',
  });

  let loading = $state(false);
  let error = $state('');

  function handleClose() {
    if (!loading) {
      isOpen = false;
      onclose();
    }
  }

  async function handleSubmit(event: Event) {
    event.preventDefault();
    if (!formData.name.trim() || !formData.address.trim()) {
      error = 'Organization name and address are required.';
      return;
    }

    loading = true;
    error = '';

    try {
      const api = getApiClient();

      const organizationData = {
        name: formData.name.trim(),
        description: formData.description.trim() || undefined,
        address: formData.address.trim(),
        phone: formData.phone.trim() || undefined,
        email: formData.email.trim() || undefined,
        website: formData.website.trim() || undefined,
      };

      const response = await api.api.v1OrganizationsUpdate(
        organization.id,
        organizationData
      );

      toast.success('Organization updated successfully');
      onupdated();
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
      handleSubmit(event);
    }
  }

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

<Modal
  bind:isOpen
  title="Edit Organization"
  {clickOrigin}
  size="lg"
  onclose={handleClose}>
  <div class="space-y-6">
    <!-- Error Message -->
    {#if error}
      <Card variant="minimal" class="border-red-200 bg-red-50">
        <div class="flex items-center space-x-2">
          <AlertTriangle class="h-5 w-5 text-red-500 flex-shrink-0" />
          <p class="text-red-700 text-sm">{error}</p>
        </div>
      </Card>
    {/if}

    <form onsubmit={handleSubmit}>
      <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
        <!-- Organization Name -->
        <div class="md:col-span-2">
          <label
            for="name"
            class="block text-sm font-medium text-gray-700 mb-2">
            Organization Name <span class="text-red-500">*</span>
          </label>
          <Input
            id="name"
            bind:value={formData.name}
            placeholder="Enter organization name"
            disabled={loading}
            required
            class="w-full" />
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
            bind:value={formData.description}
            placeholder="Brief description of the organization"
            rows="3"
            disabled={loading}
            class="w-full px-4 py-3 border border-gray-200 rounded-xl bg-white focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500 transition-all duration-200 resize-none scrollbar-thin scrollbar-thumb-gray-300 scrollbar-track-gray-100"
          ></textarea>
        </div>

        <!-- Address -->
        <div class="md:col-span-2">
          <label
            for="address"
            class="block text-sm font-medium text-gray-700 mb-2">
            Address <span class="text-red-500">*</span>
          </label>
          <Input
            id="address"
            bind:value={formData.address}
            placeholder="Organization address"
            disabled={loading}
            required
            class="w-full" />
        </div>

        <!-- Phone -->
        <div>
          <label
            for="phone"
            class="block text-sm font-medium text-gray-700 mb-2">
            Phone
          </label>
          <Input
            id="phone"
            bind:value={formData.phone}
            placeholder="Phone number"
            disabled={loading}
            class="w-full" />
        </div>

        <!-- Email -->
        <div>
          <label
            for="email"
            class="block text-sm font-medium text-gray-700 mb-2">
            Email
          </label>
          <Input
            id="email"
            type="email"
            bind:value={formData.email}
            placeholder="contact@organization.com"
            disabled={loading}
            class="w-full" />
        </div>

        <!-- Website -->
        <div>
          <label
            for="website"
            class="block text-sm font-medium text-gray-700 mb-2">
            Website
          </label>
          <Input
            id="website"
            type="url"
            bind:value={formData.website}
            placeholder="https://organization.com"
            disabled={loading}
            class="w-full" />
        </div>
      </div>

      <!-- Form Actions -->
      <div
        class="flex items-center justify-end space-x-3 mt-8 pt-6 border-t border-gray-200">
        <Button
          type="button"
          variant="outline"
          onclick={handleClose}
          disabled={loading}>
          Cancel
        </Button>
        <Button
          type="submit"
          variant="gradient"
          {loading}
          disabled={loading ||
            !formData.name.trim() ||
            !formData.address.trim()}
          class="min-w-24">
          {loading ? 'Updating...' : 'Update Organization'}
        </Button>
      </div>
    </form>
  </div>
</Modal>
