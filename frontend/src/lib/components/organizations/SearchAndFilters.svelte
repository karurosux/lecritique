<script lang="ts">
  import { Card, SearchInput, Select, FilterChip } from '$lib/components/ui';

  let {
    searchQuery = $bindable(''),
    statusFilter = $bindable('all'),
    sortBy = $bindable('name'),
    totalOrganizations = 0,
    filteredCount = 0,
    onfilterschanged = (filters: any) => {},
  }: {
    searchQuery?: string;
    statusFilter?: string;
    sortBy?: string;
    totalOrganizations?: number;
    filteredCount?: number;
    onfilterschanged?: (filters: any) => void;
  } = $props();

  function clearAllFilters() {
    searchQuery = '';
    statusFilter = 'all';
    sortBy = 'name';
    onfilterschanged({ searchQuery, statusFilter, sortBy });
  }

  function removeSearchFilter() {
    searchQuery = '';
    onfilterschanged({ searchQuery, statusFilter, sortBy });
  }

  function removeStatusFilter() {
    statusFilter = 'all';
    onfilterschanged({ searchQuery, statusFilter, sortBy });
  }

  function removeSortFilter() {
    sortBy = 'name';
    onfilterschanged({ searchQuery, statusFilter, sortBy });
  }

  function handleSearchInput(event: CustomEvent) {
    searchQuery = event.detail.value;
    onfilterschanged({ searchQuery, statusFilter, sortBy });
  }

  function handleStatusChange(event: CustomEvent) {
    statusFilter = event.detail.value;
    onfilterschanged({ searchQuery, statusFilter, sortBy });
  }

  function handleSortChange(event: CustomEvent) {
    sortBy = event.detail.value;
    onfilterschanged({ searchQuery, statusFilter, sortBy });
  }

  const statusOptions = [
    { value: 'all', label: 'All Status' },
    { value: 'active', label: 'Active Only' },
    { value: 'inactive', label: 'Inactive Only' },
  ];

  const sortOptions = [
    { value: 'name', label: 'Sort by Name' },
    { value: 'created_at', label: 'Sort by Date' },
    { value: 'status', label: 'Sort by Status' },
  ];

  let hasActiveFilters = $derived(
    searchQuery || statusFilter !== 'all' || sortBy !== 'name'
  );
</script>

<Card variant="gradient" class="mb-4">
  <div class="flex flex-col lg:flex-row gap-4">
    <!-- Search Input -->
    <SearchInput
      bind:value={searchQuery}
      placeholder="Search organizations by name, description, or email..."
      on:input={handleSearchInput} />

    <!-- Filters -->
    <div class="flex items-center space-x-3">
      <!-- Status Filter -->
      <Select
        bind:value={statusFilter}
        options={statusOptions}
        on:change={handleStatusChange} />

      <!-- Sort Options -->
      <Select
        bind:value={sortBy}
        options={sortOptions}
        on:change={handleSortChange} />
    </div>
  </div>

  <!-- Active Filters Display -->
  {#if hasActiveFilters}
    <div
      class="flex items-center justify-between mt-4 pt-4 border-t border-gray-200">
      <div class="flex items-center space-x-2">
        <span class="text-sm text-gray-600">Active filters:</span>
        {#if searchQuery}
          <FilterChip
            label="Search: "{searchQuery}""
            value="search"
            variant="blue"
            on:remove={removeSearchFilter} />
        {/if}
        {#if statusFilter !== 'all'}
          <FilterChip
            label="Status: {statusFilter}"
            value="status"
            variant="green"
            on:remove={removeStatusFilter} />
        {/if}
        {#if sortBy !== 'name'}
          <FilterChip
            label="Sort: {sortBy === 'created_at' ? 'Date' : 'Status'}"
            value="sort"
            variant="purple"
            on:remove={removeSortFilter} />
        {/if}
      </div>
      <div class="text-sm text-gray-500">
        {filteredCount} of {totalOrganizations} organizations
      </div>
    </div>
  {/if}
</Card>
