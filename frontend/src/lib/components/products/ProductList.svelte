<script lang="ts">
  import { Card, Button } from '$lib/components/ui';
  import { AlertTriangle, RefreshCw, FileText, Plus } from 'lucide-svelte';
  import ProductCard from './ProductCard.svelte';

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

  interface Category {
    id: string;
    name: string;
    products: Product[];
  }

  let {
    products = [],
    loading = false,
    error = '',
    onaddproduct = () => {},
    oneditproduct = (product: Product) => {},
    ontoggleavailability = (product: Product) => {},
    ondeleteproduct = (product: Product) => {},
    onretry = () => {},
  }: {
    products?: Product[];
    loading?: boolean;
    error?: string;
    onaddproduct?: () => void;
    oneditproduct?: (product: Product) => void;
    ontoggleavailability?: (product: Product) => void;
    ondeleteproduct?: (product: Product) => void;
    onretry?: () => void;
  } = $props();
</script>

<div>
  {#if loading}
    
    <div class="space-y-8">
      {#each Array(3) as _}
        <div class="space-y-6">
          
          <div class="flex items-center justify-between">
            <div class="h-7 bg-gray-200 rounded-lg w-48 animate-pulse"></div>
            <div class="h-5 bg-gray-200 rounded w-20 animate-pulse"></div>
          </div>

          
          <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
            {#each Array(3) as _}
              <Card class="p-6">
                <div class="animate-pulse space-y-4">
                  <div class="flex justify-between items-start">
                    <div class="space-y-2 flex-1">
                      <div class="h-6 bg-gray-200 rounded w-3/4"></div>
                      <div class="h-4 bg-gray-200 rounded w-full"></div>
                      <div class="h-4 bg-gray-200 rounded w-2/3"></div>
                    </div>
                    <div class="h-6 bg-gray-200 rounded-full w-20"></div>
                  </div>
                  <div class="h-8 bg-gray-200 rounded w-20"></div>
                  <div class="flex gap-2">
                    <div class="h-6 bg-gray-200 rounded w-16"></div>
                    <div class="h-6 bg-gray-200 rounded w-20"></div>
                  </div>
                  <div class="flex gap-2">
                    <div class="h-8 bg-gray-200 rounded flex-1"></div>
                    <div class="h-8 bg-gray-200 rounded flex-1"></div>
                    <div class="h-8 bg-gray-200 rounded w-12"></div>
                  </div>
                </div>
              </Card>
            {/each}
          </div>
        </div>
      {/each}
    </div>
  {:else if error}
    
    <Card class="p-12">
      <div class="text-center">
        <div
          class="mx-auto h-16 w-16 bg-red-100 rounded-full flex items-center justify-center mb-6">
          <AlertTriangle class="h-8 w-8 text-red-600" />
        </div>
        <h3 class="text-xl font-semibold text-gray-900 mb-2">
          Failed to load products
        </h3>
        <p class="text-gray-600 mb-6 max-w-md mx-auto">{error}</p>
        <Button onclick={onretry} class="mr-3">
          <RefreshCw class="h-4 w-4 mr-2" />
          Try Again
        </Button>
      </div>
    </Card>
  {:else if products.length === 0}
    
    <Card class="p-16">
      <div class="text-center">
        <div
          class="mx-auto h-20 w-20 bg-gray-100 rounded-full flex items-center justify-center mb-8">
          <FileText class="h-10 w-10 text-gray-400" />
        </div>
        <h3 class="text-2xl font-semibold text-gray-900 mb-3">
          No products found
        </h3>
        <p class="text-gray-600 mb-8 max-w-md mx-auto">
          Start building your product catalog by adding your first product.
          Create categories and organize your offerings to provide the best
          experience for your customers.
        </p>
        <Button onclick={onaddproduct} variant="gradient" class="shadow-lg">
          <Plus class="h-5 w-5 mr-2" />
          Add Your First Product
        </Button>
      </div>
    </Card>
  {:else}
    
    <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
      {#each products as product, productIndex}
        <div
          class="animate-fade-in-up"
          style="animation-delay: {productIndex * 50}ms">
          <ProductCard
            {product}
            onedit={oneditproduct}
            {ontoggleavailability}
            ondelete={ondeleteproduct} />
        </div>
      {/each}
    </div>
  {/if}
</div>

<style>
  @keyframes fade-in-up {
    from {
      opacity: 0;
      transform: translateY(20px);
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
</style>
