<script lang="ts">
  import { Modal, Button, Input, Card, Select } from '$lib/components/ui';
  import { APP_CONFIG } from '$lib/constants/config';

  interface Product {
    id: string;
    name: string;
    description?: string;
    price: number;
    category: string;
    is_available: boolean;
    image_url?: string;
    created_at: string;
    updated_at: string;
  }

  let { 
    isOpen = $bindable(false),
    editingProduct = null,
    loading = false,
    error = '',
    clickOrigin = null,
    onclose = () => {},
    onsave = (productData: any) => {}
  }: {
    isOpen?: boolean;
    editingProduct?: Product | null;
    loading?: boolean;
    error?: string;
    clickOrigin?: { x: number; y: number } | null;
    onclose?: () => void;
    onsave?: (productData: any) => void;
  } = $props();

  const productCategories = APP_CONFIG.productCategories;

  let formData = $state({
    name: '',
    description: '',
    price: 0,
    category: '',
    is_available: true
  });

  // Reset form when modal opens/closes or when editing product changes
  $effect(() => {
    if (isOpen) {
      if (editingProduct) {
        formData.name = editingProduct.name;
        formData.description = editingProduct.description || '';
        formData.price = editingProduct.price;
        formData.category = editingProduct.category;
        formData.is_available = editingProduct.is_available;
      } else {
        formData.name = '';
        formData.description = '';
        formData.price = 0;
        formData.category = '';
        formData.is_available = true;
      }
    }
  });

  function handleClose() {
    if (!loading) {
      onclose();
    }
  }

  function handleSubmit(event: Event) {
    event.preventDefault();
    
    if (!formData.name.trim() || !formData.category || formData.price <= 0) {
      return;
    }

    const productData = {
      name: formData.name.trim(),
      description: formData.description.trim() || undefined,
      price: formData.price,
      category: formData.category,
      is_available: formData.is_available
    };

    onsave(productData);
  }

  let isFormValid = $derived(
    formData.name.trim() !== '' && 
    formData.category !== '' && 
    formData.price > 0
  );
</script>

<Modal bind:isOpen title={editingProduct ? 'Edit Product' : 'Add New Product'} size="lg" {clickOrigin} onclose={handleClose}>
  <form onsubmit={handleSubmit} class="space-y-6">
    <!-- Error Message -->
    {#if error}
      <Card variant="minimal" class="border-red-200 bg-red-50 p-4">
        <div class="flex items-center space-x-3">
          <svg class="h-5 w-5 text-red-500 flex-shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-.833-1.964-.833-2.732 0L3.732 16.5c-.77.833.192 3 1.732 3z" />
          </svg>
          <p class="text-red-700 text-sm font-medium">{error}</p>
        </div>
      </Card>
    {/if}

    <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
      <!-- Product Name -->
      <div class="md:col-span-2">
        <Input
          id="product-name"
          type="text"
          label="Product Name"
          placeholder="Enter product name"
          bind:value={formData.name}
          disabled={loading}
          required
          variant="default"
        />
      </div>

      <!-- Description -->
      <div class="md:col-span-2">
        <Input
          id="product-description"
          type="text"
          label="Description"
          placeholder="Brief description of the product..."
          bind:value={formData.description}
          disabled={loading}
          variant="default"
        />
      </div>

      <!-- Price -->
      <div>
        <Input
          id="product-price"
          type="number"
          label="Price"
          placeholder="0.00"
          bind:value={formData.price}
          required
          disabled={loading}
          variant="default"
        />
      </div>

      <!-- Category -->
      <div>
        <label for="product-category" class="block text-sm font-semibold text-gray-700 mb-2">
          Category <span class="text-red-500">*</span>
        </label>
        <Select
          bind:value={formData.category}
          options={[
            { value: '', label: 'Select category' },
            ...productCategories.map(cat => ({ value: cat, label: cat }))
          ]}
        />
      </div>

      <!-- Availability -->
      <div>
        <label class="block text-sm font-semibold text-gray-700 mb-3">
          Availability
        </label>
        <label class="flex items-center space-x-3 cursor-pointer group">
          <input
            type="checkbox"
            class="h-5 w-5 text-blue-600 focus:ring-blue-500 border-2 border-gray-300 rounded cursor-pointer transition-all duration-200 disabled:opacity-50"
            bind:checked={formData.is_available}
            disabled={loading}
          />
          <span class="text-sm text-gray-700 group-hover:text-gray-900 transition-colors">Available for purchase</span>
        </label>
      </div>
    </div>


    <!-- Form Actions -->
    <div class="flex items-center justify-end space-x-4 pt-6 border-t border-gray-200">
      <Button 
        variant="outline" 
        type="button" 
        onclick={handleClose}
        disabled={loading}
      >
        Cancel
      </Button>
      <Button 
        type="submit"
        variant="gradient"
        disabled={loading || !isFormValid}
        class="min-w-32 shadow-lg"
      >
        {#if loading}
          <svg class="animate-spin h-4 w-4 mr-2" fill="none" viewBox="0 0 24 24">
            <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
            <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 714 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
          </svg>
          {editingProduct ? 'Saving...' : 'Creating...'}
        {:else}
          {editingProduct ? 'Save Changes' : 'Create Product'}
        {/if}
      </Button>
    </div>
  </form>
</Modal>
