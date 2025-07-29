<script lang="ts">
  import { Button } from '$lib/components/ui';
  import { Building2, Grid3x3, List, Plus, Loader2 } from 'lucide-svelte';

  interface Organization {
    status: 'active' | 'inactive';
  }

  let {
    organizations = [],
    loading = false,
    viewMode = $bindable('grid'),
    canCreateOrganization = false,
    checkingPermissions = false,
    permissionReason = '',
    onaddorganization = () => {},
  }: {
    organizations?: Organization[];
    loading?: boolean;
    viewMode?: 'grid' | 'list';
    canCreateOrganization?: boolean;
    checkingPermissions?: boolean;
    permissionReason?: string;
    onaddorganization?: () => void;
  } = $props();

  let activeCount = $derived(
    organizations.filter(r => r?.status === 'active').length
  );
  let inactiveCount = $derived(
    organizations.filter(r => r?.status === 'inactive').length
  );
</script>

<div class="mb-8">
  <div
    class="flex flex-col lg:flex-row lg:justify-between lg:items-center gap-6">
    <div class="space-y-3">
      <div class="flex items-center space-x-3">
        <div
          class="h-12 w-12 bg-gradient-to-br from-blue-500 to-purple-600 rounded-2xl flex items-center justify-center shadow-lg shadow-blue-500/25">
          <Building2 class="h-6 w-6 text-white" />
        </div>
        <div>
          <h1
            class="text-3xl font-bold bg-gradient-to-r from-gray-900 to-gray-700 bg-clip-text text-transparent">
            Your Organizations
          </h1>
          <div class="flex items-center space-x-4 mt-1">
            <p class="text-gray-600 font-medium">
              Manage your organization locations and products
            </p>
            {#if !loading}
              <div class="flex items-center space-x-3 text-sm">
                <div class="flex items-center space-x-1">
                  <div class="w-2 h-2 bg-green-400 rounded-full"></div>
                  <span class="text-gray-600">{activeCount} Active</span>
                </div>
                <div class="flex items-center space-x-1">
                  <div class="w-2 h-2 bg-gray-400 rounded-full"></div>
                  <span class="text-gray-600">{inactiveCount} Inactive</span>
                </div>
              </div>
            {/if}
          </div>
        </div>
      </div>
    </div>

    <div class="flex items-center space-x-3">
      <!-- View Mode Toggle -->
      <div
        class="flex items-center bg-gradient-to-r from-gray-50 to-gray-100 rounded-xl border border-gray-200 p-1 shadow-lg">
        <button
          class="p-2 rounded-lg transition-all duration-200 cursor-pointer {viewMode ===
          'grid'
            ? 'bg-blue-100 text-blue-600'
            : 'text-gray-500 hover:text-gray-700'}"
          onclick={() => (viewMode = 'grid')}
          aria-label="Grid view">
          <Grid3x3 class="h-4 w-4" />
        </button>
        <button
          class="p-2 rounded-lg transition-all duration-200 cursor-pointer {viewMode ===
          'list'
            ? 'bg-blue-100 text-blue-600'
            : 'text-gray-500 hover:text-gray-700'}"
          onclick={() => (viewMode = 'list')}
          aria-label="List view">
          <List class="h-4 w-4" />
        </button>
      </div>

      {#if checkingPermissions}
        <!-- Loading spinner while checking permissions -->
        <Button
          variant="gradient"
          size="lg"
          disabled
          class="group relative overflow-hidden shadow-lg">
          <div class="flex items-center">
            <Loader2 class="h-5 w-5 mr-2 animate-spin" />
            <span>Loading...</span>
          </div>
        </Button>
      {:else if canCreateOrganization}
        <Button
          variant="gradient"
          size="lg"
          class="group relative overflow-hidden shadow-lg hover:shadow-xl transition-all duration-300"
          onclick={onaddorganization}>
          <div
            class="absolute inset-0 bg-gradient-to-r from-blue-600 to-purple-600 opacity-0 group-hover:opacity-100 transition-opacity duration-300">
          </div>
          <Plus
            class="h-5 w-5 mr-2 relative z-10 group-hover:scale-110 transition-transform duration-200" />
          <span class="relative z-10">Add Organization</span>
        </Button>
      {:else}
        <!-- Disabled button with tooltip when can't create -->
        <div class="relative group">
          <Button
            variant="outline"
            size="lg"
            disabled
            class="opacity-50 cursor-not-allowed">
            <Plus class="h-5 w-5 mr-2" />
            <span>Add Organization</span>
          </Button>
          <!-- Tooltip -->
          <div
            class="absolute top-full left-1/2 transform -translate-x-1/2 mt-2 px-3 py-2 bg-gray-800 text-white text-sm rounded-lg opacity-0 group-hover:opacity-100 transition-opacity duration-200 whitespace-nowrap z-10">
            {permissionReason || 'Cannot create more organizations'}
            <div
              class="absolute bottom-full left-1/2 transform -translate-x-1/2 border-4 border-transparent border-b-gray-800">
            </div>
          </div>
        </div>
      {/if}
    </div>
  </div>
</div>
