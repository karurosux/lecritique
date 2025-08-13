<script lang="ts">
  import OrganizationCard from './OrganizationCard.svelte';
  import { NoDataAvailable } from '$lib/components/ui';
  import { createEventDispatcher } from 'svelte';

  interface Organization {
    id: string;
    name: string;
    description?: string;
    address?: string;
    phone?: string;
    email?: string;
    website?: string;
    cuisine_type?: string;
    status: 'active' | 'inactive';
    created_at: string;
    updated_at: string;
  }

  export let organizations: Organization[] = [];
  export let loading = false;
  export let viewMode: 'grid' | 'list' = 'grid';

  const dispatch = createEventDispatcher();

  function handleOrganizationClick(organization: Organization) {
    dispatch('organizationClick', organization);
  }

  function handleOrganizationEdit(
    organization: Organization,
    event?: MouseEvent
  ) {
    dispatch('organizationEdit', { organization, event });
  }

  function handleOrganizationToggleStatus(organization: Organization) {
    dispatch('organizationToggleStatus', organization);
  }

  function handleOrganizationDelete(
    organization: Organization,
    event?: MouseEvent
  ) {
    dispatch('organizationDelete', { organization, event });
  }
</script>

{#if loading}
  
  <div
    class="grid {viewMode === 'grid'
      ? 'grid-cols-1 md:grid-cols-2 lg:grid-cols-3'
      : 'grid-cols-1'} gap-6">
    {#each Array(6) as _, i}
      <div class="bg-white rounded-xl border border-gray-200 p-6 animate-pulse">
        <div class="flex items-start justify-between mb-4">
          <div class="flex items-center space-x-3">
            <div class="w-12 h-12 bg-gray-200 rounded-xl"></div>
            <div class="space-y-2">
              <div class="h-4 bg-gray-200 rounded w-32"></div>
              <div class="h-3 bg-gray-200 rounded w-20"></div>
            </div>
          </div>
        </div>
        <div class="space-y-2 mb-4">
          <div class="h-3 bg-gray-200 rounded w-full"></div>
          <div class="h-3 bg-gray-200 rounded w-3/4"></div>
        </div>
        <div class="space-y-2">
          <div class="h-3 bg-gray-200 rounded w-2/3"></div>
          <div class="h-3 bg-gray-200 rounded w-1/2"></div>
        </div>
      </div>
    {/each}
  </div>
{:else if organizations.length === 0}
  
  <NoDataAvailable
    title="No organizations found"
    description="You haven't added any organizations yet, or no organizations match your current filters."
    primaryAction={{
      label: 'Add Your First Organization',
      onClick: () => dispatch('addOrganization'),
    }}
    secondaryAction={{
      label: 'Clear Filters',
      onClick: () => dispatch('clearFilters'),
    }} />
{:else}
  
  <div
    class="grid {viewMode === 'grid'
      ? 'grid-cols-1 md:grid-cols-2 lg:grid-cols-3'
      : 'grid-cols-1'} gap-6">
    {#each organizations as organization, index}
      <OrganizationCard
        {organization}
        {viewMode}
        {index}
        onclick={handleOrganizationClick}
        onedit={handleOrganizationEdit}
        ontogglestatus={handleOrganizationToggleStatus}
        ondelete={handleOrganizationDelete} />
    {/each}
  </div>
{/if}
