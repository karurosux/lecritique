<script lang="ts">
  import { Button } from '$lib/components/ui';
  import { goto } from '$app/navigation';

  interface Product {
    id: string;
    name: string;
    description?: string;
    price: number;
    category: string;
    is_available: boolean;
    image_url?: string;
    allergens?: string[];
    preparation_time?: number;
    created_at: string;
    updated_at: string;
  }

  let {
    organizationName = '',
    organizationId = '',
    products = [],
    loading = false,
    onaddproduct = () => {}
  }: {
    organizationName?: string;
    organizationId?: string;
    products?: Product[];
    loading?: boolean;
    onaddproduct?: () => void;
  } = $props();
</script>

<div class="mb-8">
  <div class="flex flex-col lg:flex-row lg:justify-between lg:items-center gap-6">
    <div class="space-y-3">
      <div class="flex items-center space-x-3">
        <Button 
          variant="ghost" 
          onclick={() => goto('/organizations')} 
          class="p-2 mr-2" 
          aria-label="Back to organizations"
        >
          <svg class="h-5 w-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 19l-7-7 7-7" />
          </svg>
        </Button>
        
        
        <div>
          <h1 class="text-3xl font-bold bg-gradient-to-r from-gray-900 to-gray-700 bg-clip-text text-transparent">
            {organizationName} Products
          </h1>
          <div class="flex items-center space-x-4 mt-1">
            <p class="text-gray-600 font-medium">Manage your products and catalog</p>
            {#if !loading}
              <div class="flex items-center space-x-3 text-sm">
                <div class="flex items-center space-x-1">
                  <div class="w-2 h-2 bg-green-400 rounded-full"></div>
                  <span class="text-gray-600">{products.filter(d => d.is_available).length} Available</span>
                </div>
                <div class="flex items-center space-x-1">
                  <div class="w-2 h-2 bg-gray-400 rounded-full"></div>
                  <span class="text-gray-600">{products.filter(d => !d.is_available).length} Hidden</span>
                </div>
              </div>
            {/if}
          </div>
        </div>
      </div>
    </div>
    
    <div class="flex items-center space-x-3">
      <!-- Add Product Button -->
      <Button 
        variant="gradient" 
        size="lg" 
        class="group relative overflow-hidden shadow-lg hover:shadow-xl transition-all duration-300" 
        onclick={onaddproduct}
      >
        <div class="absolute inset-0 bg-gradient-to-r from-blue-600 to-purple-600 opacity-0 group-hover:opacity-100 transition-opacity duration-300"></div>
        <svg class="h-5 w-5 mr-2 relative z-10 group-hover:scale-110 transition-transform duration-200" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v16m8-8H4" />
        </svg>
        <span class="relative z-10">Add Product</span>
      </Button>
    </div>
  </div>
</div>
