<script lang="ts">
  import { Card, Button } from '$lib/components/ui';
  import { RoleGate } from '$lib/components/auth';
  import { goto } from '$app/navigation';
  import { Mail, Phone, Globe, Clock, Edit2, CheckCircle, Trash2 } from 'lucide-svelte';

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
    variant="gradient" 
    hover 
    interactive 
    class="group relative transform transition-all duration-300 animate-fade-in-up !pb-3"
    style="animation-delay: {index * 100}ms"
    onclick={handleClick}
  >
    <!-- Header Section -->
    <div class="flex items-center space-x-4 mb-4">
      <div class="h-16 w-16 bg-gradient-to-br from-blue-500 to-purple-600 rounded-2xl flex items-center justify-center shadow-lg shadow-blue-500/25">
        <span class="text-white font-bold text-lg">
          {getInitials(organization.name)}
        </span>
      </div>
      <div class="space-y-1 flex-1">
        <p class="text-sm font-semibold text-gray-600 uppercase tracking-wide">Organization</p>
        <p class="text-xl font-bold bg-gradient-to-r from-blue-600 to-purple-600 bg-clip-text text-transparent truncate max-w-[200px]" title={organization.name}>
          {organization.name}
        </p>
        <div class="flex items-center space-x-1">
          <span class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium {getStatusColor(organization.status)}">
            {organization.status}
          </span>
        </div>
      </div>
    </div>

    <!-- Description Section -->
    {#if organization.description}
      <div class="mb-4">
        <p class="text-gray-600 text-sm line-clamp-2 leading-relaxed">
          {organization.description}
        </p>
      </div>
    {/if}

    <!-- Contact Information Section -->
    <div class="space-y-2 mb-4">
      {#if organization.email}
        <div class="flex items-center space-x-3">
          <div class="h-8 w-8 bg-blue-100 rounded-lg flex items-center justify-center">
            <Mail class="h-4 w-4 text-blue-600" />
          </div>
          <span class="text-sm text-gray-700 truncate font-medium">{organization.email}</span>
        </div>
      {/if}
      {#if organization.phone}
        <div class="flex items-center space-x-3">
          <div class="h-8 w-8 bg-green-100 rounded-lg flex items-center justify-center">
            <Phone class="h-4 w-4 text-green-600" />
          </div>
          <span class="text-sm text-gray-700 font-medium">{organization.phone}</span>
        </div>
      {/if}
      {#if organization.website}
        <div class="flex items-center space-x-3">
          <div class="h-8 w-8 bg-purple-100 rounded-lg flex items-center justify-center">
            <Globe class="h-4 w-4 text-purple-600" />
          </div>
          <a href={organization.website} target="_blank" rel="noopener noreferrer" class="text-sm text-blue-600 hover:text-blue-800 truncate font-medium transition-colors">
            {organization.website}
          </a>
        </div>
      {/if}
      <!-- Creation Date -->
      <div class="flex items-center space-x-3">
        <div class="h-8 w-8 bg-gray-100 rounded-lg flex items-center justify-center">
          <Clock class="h-4 w-4 text-gray-600" />
        </div>
        <span class="text-sm text-gray-500 font-medium">
          Created {formatDate(organization.created_at)}
        </span>
      </div>
    </div>

    <!-- Footer Section -->
    <div class="flex items-center justify-between pt-4 border-t border-gray-200">
      <!-- Action Buttons -->
      <RoleGate roles={['OWNER', 'ADMIN', 'MANAGER']}>
        <div class="flex items-center space-x-1 opacity-0 group-hover:opacity-100 transition-opacity duration-200">
          <button
            class="p-1.5 text-gray-400 hover:text-blue-600 hover:bg-blue-50 rounded-lg transition-all duration-200"
            onclick={(e) => { e.stopPropagation(); handleEdit(e); }}
            aria-label="Edit organization"
          >
            <Edit2 class="h-3.5 w-3.5" />
          </button>
          <button
            class="p-1.5 text-gray-400 hover:text-green-600 hover:bg-green-50 rounded-lg transition-all duration-200"
            onclick={(e) => { e.stopPropagation(); handleToggleStatus(); }}
            aria-label="Toggle organization status"
          >
            <CheckCircle class="h-3.5 w-3.5" />
          </button>
          <button
            class="p-1.5 text-gray-400 hover:text-red-600 hover:bg-red-50 rounded-lg transition-all duration-200"
            onclick={(e) => { e.stopPropagation(); handleDelete(e); }}
            aria-label="Delete organization"
          >
            <Trash2 class="h-3.5 w-3.5" />
          </button>
        </div>
      </RoleGate>
      <Button size="sm" variant="gradient" onclick={(e) => { e.stopPropagation(); handleViewDetails(); }}>
        View Details
      </Button>
    </div>
  </Card>

{:else}
  <!-- List View -->
  <Card 
    variant="gradient" 
    hover 
    interactive 
    class="group relative transition-all duration-300 animate-fade-in-up !pb-3 !pt-3"
    style="animation-delay: {index * 50}ms"
    onclick={handleClick}
  >
    <div class="flex items-center justify-between">
      <div class="flex items-center space-x-4 flex-1 min-w-0">
        <!-- Organization Icon -->
        <div class="w-12 h-12 bg-gradient-to-br from-blue-500 to-purple-600 rounded-xl flex items-center justify-center text-white font-bold text-sm shadow-lg shadow-blue-500/25">
          {getInitials(organization.name)}
        </div>
        
        <!-- Organization Info -->
        <div class="flex-1 min-w-0 grid grid-cols-1 md:grid-cols-5 gap-4 items-center">
          <!-- Name & Status -->
          <div class="min-w-0">
            <p class="text-xs font-semibold text-gray-600 uppercase tracking-wide mb-1 whitespace-nowrap">Org</p>
            <h3 class="font-bold text-gray-900 truncate max-w-[160px]" title={organization.name}>
              {organization.name}
            </h3>
            <div class="flex items-center mt-1">
              <span class="inline-flex items-center px-2 py-0.5 rounded-full text-xs font-medium {getStatusColor(organization.status)}">
                {organization.status}
              </span>
            </div>
          </div>
          
          <!-- Contact Info -->
          <div class="min-w-0 space-y-1">
            {#if organization.email}
              <div class="flex items-center space-x-2">
                <div class="h-5 w-5 bg-blue-100 rounded flex items-center justify-center">
                  <Mail class="h-3 w-3 text-blue-600" />
                </div>
                <span class="text-xs text-gray-700 truncate font-medium">{organization.email}</span>
              </div>
            {/if}
            {#if organization.phone}
              <div class="flex items-center space-x-2">
                <div class="h-5 w-5 bg-green-100 rounded flex items-center justify-center">
                  <Phone class="h-3 w-3 text-green-600" />
                </div>
                <span class="text-xs text-gray-700 font-medium">{organization.phone}</span>
              </div>
            {/if}
          </div>
          
          <!-- Description & Website -->
          <div class="min-w-0 space-y-1">
            {#if organization.description}
              <p class="text-xs text-gray-600 line-clamp-1 font-medium">
                {organization.description}
              </p>
            {/if}
            {#if organization.website}
              <div class="flex items-center space-x-2">
                <div class="h-5 w-5 bg-purple-100 rounded flex items-center justify-center">
                  <Globe class="h-3 w-3 text-purple-600" />
                </div>
                <a href={organization.website} target="_blank" rel="noopener noreferrer" class="text-xs text-blue-600 hover:text-blue-800 truncate font-medium transition-colors">
                  {organization.website}
                </a>
              </div>
            {/if}
          </div>
          
          <!-- Creation Date -->
          <div class="min-w-0 flex items-center">
            <div class="flex items-center space-x-2">
              <div class="h-5 w-5 bg-gray-100 rounded flex items-center justify-center">
                <Clock class="h-3 w-3 text-gray-600" />
              </div>
              <span class="text-xs text-gray-500 font-medium">
                {formatDate(organization.created_at)}
              </span>
            </div>
          </div>
          
          <!-- Actions -->
          <div class="flex items-center justify-end space-x-2">
            <RoleGate roles={['OWNER', 'ADMIN', 'MANAGER']}>
              <div class="flex items-center space-x-1 opacity-0 group-hover:opacity-100 transition-opacity duration-200">
                <button
                  class="p-1.5 text-gray-400 hover:text-blue-600 hover:bg-blue-50 rounded-lg transition-all duration-200"
                  onclick={(e) => { e.stopPropagation(); handleEdit(e); }}
                  aria-label="Edit organization"
                >
                  <Edit2 class="h-3.5 w-3.5" />
                </button>
                <button
                  class="p-1.5 text-gray-400 hover:text-green-600 hover:bg-green-50 rounded-lg transition-all duration-200"
                  onclick={(e) => { e.stopPropagation(); handleToggleStatus(); }}
                  aria-label="Toggle organization status"
                >
                  <CheckCircle class="h-3.5 w-3.5" />
                </button>
                <button
                  class="p-1.5 text-gray-400 hover:text-red-600 hover:bg-red-50 rounded-lg transition-all duration-200"
                  onclick={(e) => { e.stopPropagation(); handleDelete(e); }}
                  aria-label="Delete organization"
                >
                  <Trash2 class="h-3.5 w-3.5" />
                </button>
              </div>
            </RoleGate>
            <Button size="sm" variant="gradient" onclick={(e) => { e.stopPropagation(); handleViewDetails(); }}>
              View
            </Button>
          </div>
        </div>
      </div>
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
