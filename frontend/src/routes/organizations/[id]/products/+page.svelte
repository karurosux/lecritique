<script lang="ts">
  import type { PageData } from './$types';
  import {
    Button,
    Card,
    ConfirmDialog,
    NoDataAvailable,
    SearchInput,
    Select,
    FilterChip,
  } from '$lib/components/ui';
  import { Plus } from 'lucide-svelte';
  import ProductCard from '$lib/components/products/ProductCard.svelte';
  import ProductSearchAndFilters from '$lib/components/products/ProductSearchAndFilters.svelte';
  import AddProductModal from '$lib/components/products/AddProductModal.svelte';
  import { RoleGate } from '$lib/components/auth';
  import { getApiClient, getAuthToken } from '$lib/api';
  import { toast } from 'svelte-sonner';
  import { onMount } from 'svelte';
  import { page } from '$app/stores';
  import { goto } from '$app/navigation';
  import { QuestionApi } from '$lib/api/question';
  import { APP_CONFIG } from '$lib/constants/config';

  let { data }: { data: PageData } = $props();

  let showAddProductModal = $state(false);
  let editingProduct = $state<any>(null);
  let searchQuery = $state('');
  let categoryFilter = $state('all');
  let availabilityFilter = $state('all');
  let sortBy = $state('name');
  let products = $state(data.products || []);
  let loading = $state(false);
  let productsWithQuestions = $state<string[]>([]);
  let showDeleteConfirm = $state(false);
  let deletingProduct = $state<any>(null);
  let clickOrigin = $state<{ x: number; y: number } | null>(null);

  let organization = $derived(data.organization);
  let organizationId = $derived($page.params.id);

  onMount(async () => {
    if (organization) {
      await fetchProducts();
      await fetchProductsWithQuestions();
    }
  });

  async function fetchProducts() {
    try {
      loading = true;
      const api = getApiClient();
      const response =
        await api.api.v1OrganizationsProductsList(organizationId);

      if (response.data.success && response.data.data) {
        products = response.data.data.map((product: any) => ({
          id: product.id || '',
          name: product.name || '',
          description: product.description || '',
          price: product.price || 0,
          category: product.category || 'Uncategorized',
          is_available: product.is_available !== false,
          created_at: product.created_at || '',
          updated_at: product.updated_at || '',
        }));
      }
    } catch (error) {
      console.error('Error loading products:', error);
    } finally {
      loading = false;
    }
  }

  let productsWithQuestionnaires = $derived(
    products.map(product => {
      const hasQuestions = productsWithQuestions.includes(product.id);
      return {
        ...product,
        has_questionnaire: hasQuestions,
      };
    })
  );

  const categories = APP_CONFIG.productCategories;

  let filteredProducts = $derived(
    productsWithQuestionnaires
      .filter((product: any) => {
        const matchesSearch =
          product.name.toLowerCase().includes(searchQuery.toLowerCase()) ||
          product.description
            ?.toLowerCase()
            .includes(searchQuery.toLowerCase());
        const matchesCategory =
          categoryFilter === 'all' || product.category === categoryFilter;
        const matchesAvailability =
          availabilityFilter === 'all' ||
          (availabilityFilter === 'available' && product.is_available) ||
          (availabilityFilter === 'unavailable' && !product.is_available);
        return matchesSearch && matchesCategory && matchesAvailability;
      })
      .sort((a: any, b: any) => {
        switch (sortBy) {
          case 'price':
            return a.price - b.price;
          case 'category':
            return a.category.localeCompare(b.category);
          case 'created_at':
            return (
              new Date(b.created_at).getTime() -
              new Date(a.created_at).getTime()
            );
          default:
            return a.name.localeCompare(b.name);
        }
      })
  );

  async function handleAddProduct(event?: MouseEvent) {
    if (event) {
      clickOrigin = { x: event.clientX, y: event.clientY };
    }
    editingProduct = null;
    showAddProductModal = true;
  }

  async function handleEditProduct(product: any, event?: MouseEvent) {
    if (event) {
      clickOrigin = { x: event.clientX, y: event.clientY };
    }
    editingProduct = product;
    showAddProductModal = true;
  }

  async function handleDeleteProduct(product: any, event?: MouseEvent) {
    if (event) {
      clickOrigin = { x: event.clientX, y: event.clientY };
    }
    deletingProduct = product;
    showDeleteConfirm = true;
  }

  async function confirmDeleteProduct() {
    if (!deletingProduct) return;

    try {
      const api = getApiClient();
      await api.api.v1OrganizationsProductsDelete(organizationId, deletingProduct.id);
      toast.success('Product deleted successfully');
      await fetchProducts();
    } catch (error) {
      toast.error('Failed to delete product');
      console.error(error);
    } finally {
      showDeleteConfirm = false;
      deletingProduct = null;
      clickOrigin = null;
    }
  }

  function cancelDeleteProduct() {
    showDeleteConfirm = false;
    deletingProduct = null;
    clickOrigin = null;
  }

  async function handleToggleAvailability(product: any) {
    try {
      const api = getApiClient();
      await api.api.v1OrganizationsProductsUpdate(organizationId, product.id, {
        ...product,
        is_available: !product.is_available,
      });
      toast.success(
        `${product.name} ${product.is_available ? 'disabled' : 'enabled'}`
      );
      await fetchProducts();
    } catch (error) {
      toast.error('Failed to update product availability');
      console.error(error);
    }
  }

  async function handleManageQuestionnaire(product: any) {
    goto(`/organizations/${organizationId}/questionnaire/${product.id}`);
  }

  async function fetchProductsWithQuestions() {
    try {
      const api = getApiClient();
      const response =
        await api.api.v1OrganizationsQuestionsProductsWithQuestionsList(
          organizationId
        );
      productsWithQuestions = response.data.data || [];
    } catch (error) {
      console.error('Failed to fetch products with questions:', error);
    }
  }
</script>

<svelte:head>
  <title>Products - {organization?.name || 'Organization'} | Kyooar</title>
</svelte:head>

{#if !organization}
  <div class="space-y-6">
    <div class="text-center">
      <p class="text-gray-600">Loading organization...</p>
    </div>
  </div>
{:else}
  <div class="space-y-6">
    <!-- Loading State -->
    {#if loading}
      <div class="text-center">
        <p class="text-gray-600">Loading products...</p>
      </div>
    {:else}
      <!-- Header -->
      <div class="mb-6">
        <h1 class="text-2xl font-bold text-gray-900">Products</h1>
        <p class="text-gray-600">Manage your product catalog</p>
      </div>

      <!-- Search, Filters and Add Button -->
      <div class="mb-6">
        <Card variant="glass">
          <div class="flex flex-col lg:flex-row gap-4">
            <!-- Search Input -->
            <SearchInput
              bind:value={searchQuery}
              placeholder="Search products by name or description..."
              class="flex-1" />

            <!-- Filters -->
            <div class="flex items-center space-x-3">
              <!-- Category Filter -->
              <Select bind:value={categoryFilter} options={[
                { value: 'all', label: 'All Categories' },
                ...categories.map(cat => ({ value: cat.name, label: cat.name }))
              ]} />

              <!-- Availability Filter -->
              <Select bind:value={availabilityFilter} options={[
                { value: 'all', label: 'All Status' },
                { value: 'available', label: 'Available Only' },
                { value: 'unavailable', label: 'Hidden Only' }
              ]} />

              <!-- Sort Options -->
              <Select bind:value={sortBy} options={[
                { value: 'name', label: 'Sort by Name' },
                { value: 'price', label: 'Sort by Price' },
                { value: 'category', label: 'Sort by Category' },
                { value: 'created_at', label: 'Sort by Date Added' }
              ]} />

              <!-- Add Product Button -->
              <RoleGate roles={['OWNER', 'ADMIN', 'MANAGER']}>
                <Button 
                  onclick={handleAddProduct}
                  variant="gradient"
                  size="lg"
                  class="flex-shrink-0">
                  <Plus class="mr-2 h-5 w-5" />
                  Add Product
                </Button>
              </RoleGate>
            </div>
          </div>

          <!-- Active Filters Display -->
          {#if searchQuery !== '' || categoryFilter !== 'all' || availabilityFilter !== 'all' || sortBy !== 'name'}
            <div class="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4 mt-4 pt-4 border-t border-gray-100">
              <div class="flex items-center flex-wrap gap-2">
                <span class="text-sm text-gray-600">Active filters:</span>
                {#if searchQuery}
                  <FilterChip
                    label="Search: {searchQuery}"
                    value="search"
                    variant="blue"
                    onremove={() => searchQuery = ''} />
                {/if}
                {#if categoryFilter !== 'all'}
                  <FilterChip
                    label="Category: {categoryFilter}"
                    value="category"
                    variant="purple"
                    onremove={() => categoryFilter = 'all'} />
                {/if}
                {#if availabilityFilter !== 'all'}
                  <FilterChip
                    label="Status: {availabilityFilter === 'available' ? 'Available' : 'Hidden'}"
                    value="availability"
                    variant="green"
                    onremove={() => availabilityFilter = 'all'} />
                {/if}
                {#if sortBy !== 'name'}
                  <FilterChip
                    label="Sort: {sortBy === 'price' ? 'Price' : sortBy === 'category' ? 'Category' : 'Date Added'}"
                    value="sort"
                    variant="orange"
                    onremove={() => sortBy = 'name'} />
                {/if}
              </div>

              <Button
                variant="ghost"
                size="sm"
                onclick={() => {
                  searchQuery = '';
                  categoryFilter = 'all';
                  availabilityFilter = 'all';
                  sortBy = 'name';
                }}
                class="text-gray-500 hover:text-gray-700 flex-shrink-0">
                Clear all filters
              </Button>
            </div>
          {/if}
        </Card>
      </div>

      <!-- Products Grid -->
      {#if filteredProducts.length === 0}
        {#if productsWithQuestionnaires.length === 0}
          <!-- No products at all -->
          <div class="space-y-4">
            <NoDataAvailable
              title="No products yet"
              description="Start building your product catalog by adding your first product"
              icon={Plus} />
            <div class="text-center">
              <RoleGate roles={['OWNER', 'ADMIN', 'MANAGER']}>
                <Button onclick={handleAddProduct}>
                  <Plus class="mr-2 h-4 w-4" />
                  Add First Product
                </Button>
              </RoleGate>
            </div>
          </div>
        {:else}
          <!-- No products match filters -->
          <NoDataAvailable
            title="No products match your filters"
            description="Try adjusting your search or filters"
            icon={Plus} />
        {/if}
      {:else}
        <div class="grid gap-6 md:grid-cols-2 lg:grid-cols-3">
          {#each filteredProducts as product, index (product.id)}
            <ProductCard
              {product}
              {index}
              onedit={handleEditProduct}
              ondelete={handleDeleteProduct}
              ontoggleavailability={() => handleToggleAvailability(product)}
              ongeneratequestionnaire={() =>
                handleManageQuestionnaire(product)} />
          {/each}
        </div>
      {/if}
    {/if}
  </div>
{/if}

<!-- Add/Edit Product Modal -->
{#if organization}
  <AddProductModal
    bind:isOpen={showAddProductModal}
    {editingProduct}
    {clickOrigin}
    onclose={() => {
      showAddProductModal = false;
      editingProduct = null;
      clickOrigin = null;
    }}
    onsave={async productData => {
      try {
        const api = getApiClient();
        if (editingProduct) {
          await api.api.v1OrganizationsProductsUpdate(organizationId, editingProduct.id, {
            ...productData,
            organization_id: organizationId,
          });
          toast.success('Product updated successfully');
        } else {
          await api.api.v1OrganizationsProductsCreate(organizationId, {
            ...productData,
            organization_id: organizationId,
          });
          toast.success('Product created successfully');
        }
        showAddProductModal = false;
        editingProduct = null;
        await fetchProducts();
        await fetchProductsWithQuestions();
      } catch (error) {
        console.error('Error saving product:', error);
        toast.error('Failed to save product');
      }
    }} />
{/if}

<!-- Delete Confirmation Dialog -->
<ConfirmDialog
  isOpen={showDeleteConfirm}
  title="Delete Product"
  message={`Are you sure you want to delete "${deletingProduct?.name}"? This action cannot be undone.`}
  confirmText="Delete"
  cancelText="Cancel"
  variant="danger"
  {clickOrigin}
  onConfirm={confirmDeleteProduct}
  onCancel={cancelDeleteProduct} />
