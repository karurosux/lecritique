<script lang="ts">
  import { Card, SearchInput, Select, FilterChip, Button } from '$lib/components/ui';

  interface Category {
    id: string;
    name: string;
    products: any[];
  }

  let {
    searchQuery = $bindable(''),
    categoryFilter = $bindable('all'),
    availabilityFilter = $bindable('all'),
    sortBy = $bindable('name'),
    categories = [],
    totalProducts = 0,
    filteredCount = 0
  }: {
    searchQuery?: string;
    categoryFilter?: string;
    availabilityFilter?: string;
    sortBy?: string;
    categories?: Category[];
    totalProducts?: number;
    filteredCount?: number;
  } = $props();

  function clearAllFilters() {
    searchQuery = '';
    categoryFilter = 'all';
    availabilityFilter = 'all';
    sortBy = 'name';
  }

  function removeSearchFilter() {
    searchQuery = '';
  }

  function removeCategoryFilter() {
    categoryFilter = 'all';
  }

  function removeAvailabilityFilter() {
    availabilityFilter = 'all';
  }

  function removeSortFilter() {
    sortBy = 'name';
  }

  const categoryOptions = $derived([
    { value: 'all', label: 'All Categories' },
    ...categories.map(cat => ({ value: cat.name, label: cat.name }))
  ]);

  const availabilityOptions = [
    { value: 'all', label: 'All Status' },
    { value: 'available', label: 'Available Only' },
    { value: 'unavailable', label: 'Hidden Only' }
  ];

  const sortOptions = [
    { value: 'name', label: 'Sort by Name' },
    { value: 'price', label: 'Sort by Price' },
    { value: 'category', label: 'Sort by Category' },
    { value: 'created_at', label: 'Sort by Date Added' }
  ];

  let hasActiveFilters = $derived(
    searchQuery !== '' || 
    categoryFilter !== 'all' || 
    availabilityFilter !== 'all' || 
    sortBy !== 'name'
  );
</script>

<Card variant="glass" class="mb-4">
  <div class="flex flex-col lg:flex-row gap-4">
    <!-- Search Input -->
    <SearchInput
      bind:value={searchQuery}
      placeholder="Search products by name or description..."
    />

    <!-- Filters -->
    <div class="flex items-center space-x-3">
      <!-- Category Filter -->
      <Select
        bind:value={categoryFilter}
        options={categoryOptions}
      />

      <!-- Availability Filter -->
      <Select
        bind:value={availabilityFilter}
        options={availabilityOptions}
      />

      <!-- Sort Options -->
      <Select
        bind:value={sortBy}
        options={sortOptions}
      />
    </div>
  </div>

  <!-- Active Filters Display -->
  {#if hasActiveFilters}
    <div class="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-4 mt-4 pt-4 border-t border-gray-100">
      <div class="flex items-center flex-wrap gap-2">
        <span class="text-sm text-gray-600">Active filters:</span>
        {#if searchQuery}
          <FilterChip
            label='Search: "{searchQuery}"'
            value="search"
            variant="blue"
            onremove={removeSearchFilter}
          />
        {/if}
        {#if categoryFilter !== 'all'}
          <FilterChip
            label='Category: {categoryFilter}'
            value="category"
            variant="purple"
            onremove={removeCategoryFilter}
          />
        {/if}
        {#if availabilityFilter !== 'all'}
          <FilterChip
            label='Status: {availabilityFilter === "available" ? "Available" : "Hidden"}'
            value="availability"
            variant="green"
            onremove={removeAvailabilityFilter}
          />
        {/if}
        {#if sortBy !== 'name'}
          <FilterChip
            label='Sort: {sortBy === "price" ? "Price" : sortBy === "category" ? "Category" : "Date Added"}'
            value="sort"
            variant="orange"
            onremove={removeSortFilter}
          />
        {/if}
      </div>
      
      <Button 
        variant="ghost" 
        size="sm" 
        onclick={clearAllFilters}
        class="text-gray-500 hover:text-gray-700 flex-shrink-0"
      >
        Clear all filters
      </Button>
    </div>
  {/if}
</Card>
