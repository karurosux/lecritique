<script lang="ts">
  import { Card, Button } from '$lib/components/ui';
  import { RoleGate } from '$lib/components/auth';
  import { goto } from '$app/navigation';

  interface Organization {
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

  let {
    organization,
    viewMode = 'grid',
    index = 0,
    onclick = (organization: Organization) => {},
    onedit = (organization: Organization, event?: MouseEvent) => {},
    ontogglestatus = (organization: Organization) => {},
    ondelete = (organization: Organization, event?: MouseEvent) => {}
  }: {
    organization: Organization;
    viewMode?: 'grid' | 'list';
    index?: number;
    onclick?: (organization: Organization) => void;
    onedit?: (organization: Organization, event?: MouseEvent) => void;
    ontogglestatus?: (organization: Organization) => void;
    ondelete?: (organization: Organization, event?: MouseEvent) => void;
  } = $props();

  function handleClick() {
    onclick(organization);
  }

  function handleEdit(event?: MouseEvent) {
    onedit(organization, event);
  }

  function handleToggleStatus() {
    ontogglestatus(organization);
  }

  function handleDelete(event?: MouseEvent) {
    ondelete(organization, event);
  }

  function handleViewDetails() {
    // Navigate to organization details/products page
    goto(`/organizations/${organization.id}/products`);
  }

  function formatDate(dateString: string): string {
    return new Date(dateString).toLocaleDateString();
  }

  function getStatusColor(status: string): string {
    return status === 'active' ? 'bg-green-100 text-green-800' : 'bg-gray-100 text-gray-800';
  }

  function getInitials(name: string): string {
    return name.split(' ').map(word => word[0]).join('').toUpperCase().slice(0, 2);
  }
</script>

{#if viewMode === 'grid'}
  <!-- Grid View -->
  <Card 
    variant="default" 
    hover 
    interactive 
    class="group transform transition-all duration-300 animate-fade-in-up"
    style="animation-delay: {index * 100}ms"
    onclick={handleClick}
  >
    <div class="flex items-start justify-between mb-4">
      <div class="flex items-center space-x-3">
        <div class="w-12 h-12 bg-gradient-to-br from-blue-500 to-purple-600 rounded-xl flex items-center justify-center text-white font-bold text-sm shadow-lg shadow-blue-500/25 group-hover:scale-110 transition-transform duration-200">
          {getInitials(organization.name)}
        </div>
        <div class="flex-1 min-w-0">
          <h3 class="font-bold text-lg text-gray-900 group-hover:text-blue-600 transition-colors duration-200 truncate">
            {organization.name}
          </h3>
          <div class="flex items-center space-x-2 mt-1">
            <span class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium {getStatusColor(organization.status)}">
              {organization.status}
            </span>
            {#if organization.cuisine_type}
              <span class="text-xs text-gray-500 bg-gray-100 px-2 py-1 rounded-md">
                {organization.cuisine_type}
              </span>
            {/if}
          </div>
        </div>
      </div>
      <RoleGate roles={['OWNER', 'ADMIN', 'MANAGER']}>
        <div class="flex items-center space-x-1 opacity-0 group-hover:opacity-100 transition-opacity duration-200">
          <button
            class="p-2 text-gray-400 hover:text-blue-600 hover:bg-blue-50 rounded-lg transition-all duration-200"
            onclick={(e) => { e.stopPropagation(); handleEdit(e); }}
            aria-label="Edit organization"
          >
            <svg class="h-4 w-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z" />
            </svg>
          </button>
          <button
            class="p-2 text-gray-400 hover:text-green-600 hover:bg-green-50 rounded-lg transition-all duration-200"
            onclick={(e) => { e.stopPropagation(); handleToggleStatus(); }}
            aria-label="Toggle organization status"
          >
            <svg class="h-4 w-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
            </svg>
          </button>
          <button
            class="p-2 text-gray-400 hover:text-red-600 hover:bg-red-50 rounded-lg transition-all duration-200"
            onclick={(e) => { e.stopPropagation(); handleDelete(e); }}
            aria-label="Delete organization"
          >
            <svg class="h-4 w-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
            </svg>
          </button>
        </div>
      </RoleGate>
    </div>

    {#if organization.description}
      <p class="text-gray-600 text-sm mb-4 line-clamp-2 group-hover:text-gray-700 transition-colors duration-200">
        {organization.description}
      </p>
    {/if}

    <div class="grid grid-cols-1 gap-3 text-sm text-gray-600">
      {#if organization.email}
        <div class="flex items-center space-x-2">
          <svg class="h-4 w-4 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 8l7.89 4.26a2 2 0 002.22 0L21 8M5 19h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v10a2 2 0 002 2z" />
          </svg>
          <span class="truncate">{organization.email}</span>
        </div>
      {/if}
      {#if organization.phone}
        <div class="flex items-center space-x-2">
          <svg class="h-4 w-4 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 5a2 2 0 012-2h3.28a1 1 0 01.948.684l1.498 4.493a1 1 0 01-.502 1.21l-2.257 1.13a11.042 11.042 0 005.516 5.516l1.13-2.257a1 1 0 011.21-.502l4.493 1.498a1 1 0 01.684.949V19a2 2 0 01-2 2h-1C9.716 21 3 14.284 3 6V5z" />
          </svg>
          <span>{organization.phone}</span>
        </div>
      {/if}
      {#if organization.website}
        <div class="flex items-center space-x-2">
          <svg class="h-4 w-4 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 12a9 9 0 01-9 9m9-9a9 9 0 00-9-9m9 9H3m9 9v-9m0-9v9m0 9c-5 0-9-4-9-9s4-9 9-9" />
          </svg>
          <a href={organization.website} target="_blank" rel="noopener noreferrer" class="text-blue-600 hover:text-blue-800 truncate">
            {organization.website}
          </a>
        </div>
      {/if}
    </div>

    <div class="flex items-center justify-between mt-4 pt-4 border-t border-gray-100">
      <span class="text-xs text-gray-500">
        Created {formatDate(organization.created_at)}
      </span>
      <div class="opacity-0 group-hover:opacity-100 transition-opacity duration-200">
        <Button size="sm" variant="outline" onclick={(e) => { e.stopPropagation(); handleViewDetails(); }}>
          Open
        </Button>
      </div>
    </div>
  </Card>

{:else}
  <!-- List View -->
  <Card 
    variant="default" 
    hover 
    interactive 
    class="group transition-all duration-300 animate-fade-in-up"
    style="animation-delay: {index * 50}ms"
    onclick={handleClick}
  >
    <div class="flex items-center justify-between">
      <div class="flex items-center space-x-4 flex-1 min-w-0">
        <div class="w-10 h-10 bg-gradient-to-br from-blue-500 to-purple-600 rounded-lg flex items-center justify-center text-white font-bold text-sm shadow-md shadow-blue-500/25 group-hover:scale-110 transition-transform duration-200">
          {getInitials(organization.name)}
        </div>
        
        <div class="flex-1 min-w-0 grid grid-cols-1 md:grid-cols-4 gap-4">
          <div class="min-w-0">
            <h3 class="font-semibold text-gray-900 group-hover:text-blue-600 transition-colors duration-200 truncate">
              {organization.name}
            </h3>
            <div class="flex items-center space-x-2 mt-1">
              <span class="inline-flex items-center px-2 py-0.5 rounded-full text-xs font-medium {getStatusColor(organization.status)}">
                {organization.status}
              </span>
              {#if organization.cuisine_type}
                <span class="text-xs text-gray-500 bg-gray-100 px-2 py-0.5 rounded">
                  {organization.cuisine_type}
                </span>
              {/if}
            </div>
          </div>
          
          <div class="min-w-0">
            {#if organization.email}
              <div class="flex items-center space-x-1 text-sm text-gray-600">
                <svg class="h-3 w-3 text-gray-400 flex-shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 8l7.89 4.26a2 2 0 002.22 0L21 8M5 19h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v10a2 2 0 002 2z" />
                </svg>
                <span class="truncate">{organization.email}</span>
              </div>
            {/if}
            {#if organization.phone}
              <div class="flex items-center space-x-1 text-sm text-gray-600 mt-1">
                <svg class="h-3 w-3 text-gray-400 flex-shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 5a2 2 0 012-2h3.28a1 1 0 01.948.684l1.498 4.493a1 1 0 01-.502 1.21l-2.257 1.13a11.042 11.042 0 005.516 5.516l1.13-2.257a1 1 0 011.21-.502l4.493 1.498a1 1 0 01.684.949V19a2 2 0 01-2 2h-1C9.716 21 3 14.284 3 6V5z" />
                </svg>
                <span>{organization.phone}</span>
              </div>
            {/if}
          </div>
          
          <div class="min-w-0">
            {#if organization.description}
              <p class="text-sm text-gray-600 line-clamp-1">
                {organization.description}
              </p>
            {/if}
            {#if organization.website}
              <a href={organization.website} target="_blank" rel="noopener noreferrer" class="text-sm text-blue-600 hover:text-blue-800 truncate block mt-1">
                {organization.website}
              </a>
            {/if}
          </div>
          
          <div class="text-right">
            <span class="text-xs text-gray-500">
              {formatDate(organization.created_at)}
            </span>
          </div>
        </div>
      </div>
      
      <RoleGate roles={['OWNER', 'ADMIN', 'MANAGER']}>
        <div class="flex items-center space-x-2 ml-4 opacity-0 group-hover:opacity-100 transition-opacity duration-200">
          <button
            class="p-2 text-gray-400 hover:text-blue-600 hover:bg-blue-50 rounded-lg transition-all duration-200"
            onclick={(e) => { e.stopPropagation(); handleEdit(); }}
            aria-label="Edit organization"
          >
            <svg class="h-4 w-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z" />
            </svg>
          </button>
          <button
            class="p-2 text-gray-400 hover:text-green-600 hover:bg-green-50 rounded-lg transition-all duration-200"
            onclick={(e) => { e.stopPropagation(); handleToggleStatus(); }}
            aria-label="Toggle organization status"
          >
            <svg class="h-4 w-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
            </svg>
          </button>
          <button
            class="p-2 text-gray-400 hover:text-red-600 hover:bg-red-50 rounded-lg transition-all duration-200"
            onclick={(e) => { e.stopPropagation(); handleDelete(e); }}
            aria-label="Delete organization"
          >
            <svg class="h-4 w-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
            </svg>
          </button>
          <Button size="sm" variant="outline" onclick={(e) => { e.stopPropagation(); handleViewDetails(); }}>
            Open
          </Button>
        </div>
      </RoleGate>
    </div>
  </Card>
{/if}

<style>
  @keyframes fade-in-up {
    from {
      opacity: 0;
      transform: translateY(10px);
    }
    to {
      opacity: 1;
      transform: translateY(0);
    }
  }

  .animate-fade-in-up {
    animation: fade-in-up 0.6s ease-out forwards;
    opacity: 0;
  }

  .line-clamp-1 {
    display: -webkit-box;
    -webkit-line-clamp: 1;
    -webkit-box-orient: vertical;
    overflow: hidden;
  }

  .line-clamp-2 {
    display: -webkit-box;
    -webkit-line-clamp: 2;
    -webkit-box-orient: vertical;
    overflow: hidden;
  }
</style>
