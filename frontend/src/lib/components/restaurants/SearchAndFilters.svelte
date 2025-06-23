<script lang="ts">
  import { Card, SearchInput, Select, FilterChip } from '$lib/components/ui';
  import { createEventDispatcher } from 'svelte';

  export let searchQuery = '';
  export let statusFilter = 'all';
  export let sortBy = 'name';
  export let totalRestaurants = 0;
  export let filteredCount = 0;

  const dispatch = createEventDispatcher();

  function clearAllFilters() {
    searchQuery = '';
    statusFilter = 'all';
    sortBy = 'name';
    dispatch('filtersChanged', { searchQuery, statusFilter, sortBy });
  }

  function removeSearchFilter() {
    searchQuery = '';
    dispatch('filtersChanged', { searchQuery, statusFilter, sortBy });
  }

  function removeStatusFilter() {
    statusFilter = 'all';
    dispatch('filtersChanged', { searchQuery, statusFilter, sortBy });
  }

  function removeSortFilter() {
    sortBy = 'name';
    dispatch('filtersChanged', { searchQuery, statusFilter, sortBy });
  }

  function handleSearchInput(event: CustomEvent) {
    searchQuery = event.detail.value;
    dispatch('filtersChanged', { searchQuery, statusFilter, sortBy });
  }

  function handleStatusChange(event: CustomEvent) {
    statusFilter = event.detail.value;
    dispatch('filtersChanged', { searchQuery, statusFilter, sortBy });
  }

  function handleSortChange(event: CustomEvent) {
    sortBy = event.detail.value;
    dispatch('filtersChanged', { searchQuery, statusFilter, sortBy });
  }

  $: statusOptions = [
    { value: 'all', label: 'All Status' },
    { value: 'active', label: 'Active Only' },
    { value: 'inactive', label: 'Inactive Only' }
  ];

  $: sortOptions = [
    { value: 'name', label: 'Sort by Name' },
    { value: 'created_at', label: 'Sort by Date' },
    { value: 'status', label: 'Sort by Status' }
  ];

  $: hasActiveFilters = searchQuery || statusFilter !== 'all' || sortBy !== 'name';
</script>

<Card variant="glass" class="mb-4">
  <div class="flex flex-col lg:flex-row gap-4">
    <!-- Search Input -->
    <SearchInput
      bind:value={searchQuery}
      placeholder="Search restaurants by name, description, or email..."
      on:input={handleSearchInput}
    />

    <!-- Filters -->
    <div class="flex items-center space-x-3">
      <!-- Status Filter -->
      <Select
        bind:value={statusFilter}
        options={statusOptions}
        on:change={handleStatusChange}
      />

      <!-- Sort Options -->
      <Select
        bind:value={sortBy}
        options={sortOptions}
        on:change={handleSortChange}
      />
    </div>
  </div>

  <!-- Active Filters Display -->
  {#if hasActiveFilters}
    <div class="flex items-center justify-between mt-4 pt-4 border-t border-gray-100">
      <div class="flex items-center space-x-2">
        <span class="text-sm text-gray-600">Active filters:</span>
        {#if searchQuery}
          <FilterChip
            label='Search: "{searchQuery}"'
            value="search"
            variant="blue"
            on:remove={removeSearchFilter}
          />
        {/if}
        {#if statusFilter !== 'all'}
          <FilterChip
            label="Status: {statusFilter}"
            value="status"
            variant="green"
            on:remove={removeStatusFilter}
          />
        {/if}
        {#if sortBy !== 'name'}
          <FilterChip
            label="Sort: {sortBy === 'created_at' ? 'Date' : 'Status'}"
            value="sort"
            variant="purple"
            on:remove={removeSortFilter}
          />
        {/if}
      </div>
      <div class="text-sm text-gray-500">
        {filteredCount} of {totalRestaurants} restaurants
      </div>
    </div>
  {/if}
</Card>