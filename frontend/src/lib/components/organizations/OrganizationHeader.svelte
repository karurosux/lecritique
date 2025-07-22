<script lang="ts">
  import { Button } from '$lib/components/ui';

  interface Organization {
    status: 'active' | 'inactive';
  }

  let {
    organizations = [],
    loading = false,
    viewMode = $bindable('grid'),
    canCreateOrganization = false,
    checkingPermissions = false,
    permissionReason = "",
    onaddorganization = () => {}
  }: {
    organizations?: Organization[];
    loading?: boolean;
    viewMode?: 'grid' | 'list';
    canCreateOrganization?: boolean;
    checkingPermissions?: boolean;
    permissionReason?: string;
    onaddorganization?: () => void;
  } = $props();

  let activeCount = $derived(organizations.filter(r => r?.status === 'active').length);
  let inactiveCount = $derived(organizations.filter(r => r?.status === 'inactive').length);
</script>

<div class="mb-8">
  <div class="flex flex-col lg:flex-row lg:justify-between lg:items-center gap-6">
    <div class="space-y-3">
      <div class="flex items-center space-x-3">
        <div class="h-12 w-12 bg-gradient-to-br from-blue-500 to-purple-600 rounded-2xl flex items-center justify-center shadow-lg shadow-blue-500/25">
          <svg class="h-6 w-6 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 21V5a2 2 0 00-2-2H7a2 2 0 00-2 2v16m14 0h2m-2 0h-4m-5 0H9m11 0a2 2 0 01-2 2H7a2 2 0 01-2-2m2-4h2.01M7 16h6M7 8h6v4H7V8z" />
          </svg>
        </div>
        <div>
          <h1 class="text-3xl font-bold bg-gradient-to-r from-gray-900 to-gray-700 bg-clip-text text-transparent">
            Your Organizations
          </h1>
          <div class="flex items-center space-x-4 mt-1">
            <p class="text-gray-600 font-medium">Manage your organization locations and products</p>
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
      <div class="flex items-center bg-white rounded-xl border border-gray-200 p-1 shadow-sm">
        <button
          class="p-2 rounded-lg transition-all duration-200 {viewMode === 'grid' ? 'bg-blue-100 text-blue-600' : 'text-gray-500 hover:text-gray-700'}"
          onclick={() => viewMode = 'grid'}
          aria-label="Grid view"
        >
          <svg class="h-4 w-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 6a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2H6a2 2 0 01-2-2V6zM14 6a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2h-2a2 2 0 01-2-2V6zM4 16a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2H6a2 2 0 01-2-2v-2zM14 16a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2h-2a2 2 0 01-2-2v-2z" />
          </svg>
        </button>
        <button
          class="p-2 rounded-lg transition-all duration-200 {viewMode === 'list' ? 'bg-blue-100 text-blue-600' : 'text-gray-500 hover:text-gray-700'}"
          onclick={() => viewMode = 'list'}
          aria-label="List view"
        >
          <svg class="h-4 w-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M4 6h16M4 10h16M4 14h16M4 18h16" />
          </svg>
        </button>
      </div>

      {#if checkingPermissions}
        <!-- Loading spinner while checking permissions -->
        <Button variant="gradient" size="lg" disabled class="group relative overflow-hidden shadow-lg">
          <div class="flex items-center">
            <svg class="animate-spin h-5 w-5 mr-2" fill="none" viewBox="0 0 24 24">
              <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
              <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
            </svg>
            <span>Loading...</span>
          </div>
        </Button>
      {:else if canCreateOrganization}
        <Button variant="gradient" size="lg" class="group relative overflow-hidden shadow-lg hover:shadow-xl transition-all duration-300" onclick={onaddorganization}>
          <div class="absolute inset-0 bg-gradient-to-r from-blue-600 to-purple-600 opacity-0 group-hover:opacity-100 transition-opacity duration-300"></div>
          <svg class="h-5 w-5 mr-2 relative z-10 group-hover:scale-110 transition-transform duration-200" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4" />
          </svg>
          <span class="relative z-10">Add Organization</span>
        </Button>
      {:else}
        <!-- Disabled button with tooltip when can't create -->
        <div class="relative group">
          <Button variant="outline" size="lg" disabled class="opacity-50 cursor-not-allowed">
            <svg class="h-5 w-5 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4" />
            </svg>
            <span>Add Organization</span>
          </Button>
          <!-- Tooltip -->
          <div class="absolute top-full left-1/2 transform -translate-x-1/2 mt-2 px-3 py-2 bg-gray-800 text-white text-sm rounded-lg opacity-0 group-hover:opacity-100 transition-opacity duration-200 whitespace-nowrap z-10">
            {permissionReason || "Cannot create more organizations"}
            <div class="absolute bottom-full left-1/2 transform -translate-x-1/2 border-4 border-transparent border-b-gray-800"></div>
          </div>
        </div>
      {/if}
    </div>
  </div>
</div>
