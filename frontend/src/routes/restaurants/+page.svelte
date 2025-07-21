<script lang="ts">
  import { onMount } from "svelte";
  import { getApiClient, handleApiError } from "$lib/api/client";
  import { auth } from "$lib/stores/auth";
  import { subscription } from "$lib/stores/subscription";
  import { goto } from "$app/navigation";
  import RestaurantHeader from "$lib/components/restaurants/RestaurantHeader.svelte";
  import SearchAndFilters from "$lib/components/restaurants/SearchAndFilters.svelte";
  import RestaurantList from "$lib/components/restaurants/RestaurantList.svelte";
  import CreateRestaurantModal from "$lib/components/restaurants/CreateRestaurantModal.svelte";
  import EditRestaurantModal from "$lib/components/restaurants/EditRestaurantModal.svelte";
  import { ConfirmDialog, Card, Button } from "$lib/components/ui";

  interface Restaurant {
    id: string;
    name: string;
    description?: string;
    address: string;
    phone?: string;
    email?: string;
    website?: string;
    cuisine_type?: string;
    status: "active" | "inactive";
    created_at: string;
    updated_at: string;
  }

  let loading = $state(true);
  let error = $state("");
  let restaurants = $state<Restaurant[]>([]);
  let searchQuery = $state("");
  let statusFilter = $state("all"); // 'all', 'active', 'inactive'
  let sortBy = $state("name"); // 'name', 'created_at', 'status'
  let viewMode = $state<"grid" | "list">("grid");
  let showCreateModal = $state(false);
  let showEditModal = $state(false);
  let editingRestaurant = $state<Restaurant | null>(null);
  let showDeleteConfirm = $state(false);
  let deletingRestaurant = $state<Restaurant | null>(null);
  let canCreateRestaurant = $state(false);
  let checkingPermissions = $state(true);
  let permissionReason = $state("");
  let clickOrigin = $state<{ x: number; y: number } | null>(null);

  let filteredRestaurants = $derived(
    restaurants
      .filter((restaurant) => {
        if (!restaurant) return false;
        const matchesSearch =
          restaurant.name?.toLowerCase().includes(searchQuery.toLowerCase()) ||
          restaurant.description
            ?.toLowerCase()
            .includes(searchQuery.toLowerCase()) ||
          restaurant.email?.toLowerCase().includes(searchQuery.toLowerCase());
        const matchesStatus =
          statusFilter === "all" || restaurant.status === statusFilter;
        return matchesSearch && matchesStatus;
      })
      .sort((a, b) => {
        if (!a || !b) return 0;
        switch (sortBy) {
          case "created_at":
            return (
              new Date(b.created_at || '').getTime() - new Date(a.created_at || '').getTime()
            );
          case "status":
            return (a.status || '').localeCompare(b.status || '');
          default:
            return (a.name || '').localeCompare(b.name || '');
        }
      })
  );

  let authState = $derived($auth);
  let subscriptionState = $derived($subscription);
  let teamMembers = $derived(subscriptionState.subscription?.team_members || []);
  let hasInitialized = $state(false);

  // Handle authentication and initial load
  $effect(() => {
    if (!authState.isAuthenticated) {
      goto("/login");
      return;
    }
    
    // Only load restaurants once when authenticated
    if (authState.isAuthenticated && !hasInitialized) {
      hasInitialized = true;
      loadRestaurants();
    }
  });

  // Update permissions when restaurants change
  $effect(() => {
    if (authState.isAuthenticated && restaurants.length >= 0) {
      checkCreatePermission();
    }
  });

  async function loadRestaurants() {
    loading = true;
    error = "";

    try {
      const api = getApiClient();

      // Use actual API client to get restaurants
      const response = await api.api.v1RestaurantsList();

      if (response.data.success && response.data.data) {
        restaurants = response.data.data.map((restaurant) => ({
          id: restaurant.id || "",
          name: restaurant.name || "",
          description: restaurant.description || "",
          address: "", // Note: address would come from locations array
          phone: restaurant.phone || "",
          email: restaurant.email || "",
          website: restaurant.website || "",
          cuisine_type: "", // Note: cuisine_type not in API model
          status: restaurant.is_active ? "active" : "inactive",
          created_at: restaurant.created_at || "",
          updated_at: restaurant.updated_at || "",
        }));
      } else {
        restaurants = [];
      }
    } catch (err) {
      error = handleApiError(err);
    } finally {
      loading = false;
    }
  }

  async function checkCreatePermission() {
    try {
      checkingPermissions = true;
      const api = getApiClient();
      
      const response = await api.api.v1UserCanCreateRestaurantList();
      
      if (response.data.success) {
        const data = response.data.data;
        canCreateRestaurant = data.can_create || false;
        
        if (!canCreateRestaurant) {
          permissionReason = data.reason || "Cannot create more restaurants";
        }
      } else {
        canCreateRestaurant = false;
        permissionReason = "Unable to verify permissions";
      }
      
    } catch (err) {
      console.error("Error checking create permission:", err);
      canCreateRestaurant = false;
      permissionReason = "Unable to verify permissions";
    } finally {
      checkingPermissions = false;
    }
  }

  function handleFiltersChanged(event: CustomEvent) {
    const { searchQuery: newSearchQuery, statusFilter: newStatusFilter, sortBy: newSortBy } = event.detail;
    searchQuery = newSearchQuery;
    statusFilter = newStatusFilter;
    sortBy = newSortBy;
  }

  function handleRestaurantClick(event: CustomEvent) {
    const restaurant = event.detail;
    goto(`/restaurants/${restaurant.id}`);
  }

  function handleRestaurantEdit(event: CustomEvent) {
    const { restaurant, event: mouseEvent } = event.detail;
    if (mouseEvent) {
      clickOrigin = { x: mouseEvent.clientX, y: mouseEvent.clientY };
    }
    editingRestaurant = restaurant;
    showEditModal = true;
  }

  function handleRestaurantToggleStatus(event: CustomEvent) {
    const restaurant = event.detail;
    toggleRestaurantStatus(restaurant);
  }

  function handleRestaurantDelete(event: CustomEvent) {
    const { restaurant, event: mouseEvent } = event.detail;
    if (mouseEvent) {
      clickOrigin = { x: mouseEvent.clientX, y: mouseEvent.clientY };
    }
    deletingRestaurant = restaurant;
    showDeleteConfirm = true;
  }

  function handleAddRestaurant() {
    showCreateModal = true;
  }

  function handleCloseModal() {
    showCreateModal = false;
  }

  function handleCloseEditModal() {
    showEditModal = false;
    editingRestaurant = null;
    clickOrigin = null;
  }

  function handleRestaurantUpdated() {
    showEditModal = false;
    editingRestaurant = null;
    clickOrigin = null;
    loadRestaurants(); // Refresh the restaurant list
  }

  function handleDeleteConfirm() {
    if (deletingRestaurant) {
      deleteRestaurant(deletingRestaurant);
    }
    showDeleteConfirm = false;
    deletingRestaurant = null;
    clickOrigin = null;
  }

  function handleDeleteCancel() {
    showDeleteConfirm = false;
    deletingRestaurant = null;
    clickOrigin = null;
  }

  function handleRestaurantCreated(event: CustomEvent) {
    // Reload the full list to ensure data consistency and get updated permissions
    loadRestaurants();
    // Refresh permissions since restaurant count changed
    checkCreatePermission();
  }

  function handleClearFilters() {
    searchQuery = '';
    statusFilter = 'all';
    sortBy = 'name';
  }

  async function toggleRestaurantStatus(restaurant: Restaurant) {
    try {
      const api = getApiClient();
      const newStatus = restaurant.status === "active" ? "inactive" : "active";

      // Update restaurant status via API
      await api.api.v1RestaurantsUpdate(restaurant.id, {
        is_active: newStatus === "active",
      });

      // Update local state
      restaurants = restaurants.map((r) =>
        r.id === restaurant.id
          ? { ...r, status: newStatus, updated_at: new Date().toISOString() }
          : r,
      );
    } catch (err) {
      error = handleApiError(err);
    }
  }

  async function deleteRestaurant(restaurant: Restaurant) {
    try {
      const api = getApiClient();

      // Delete restaurant via API
      await api.api.v1RestaurantsDelete(restaurant.id);

      // Update local state
      restaurants = restaurants.filter((r) => r.id !== restaurant.id);
    } catch (err) {
      error = handleApiError(err);
    }
  }
</script>

<svelte:head>
  <title>Restaurants - LeCritique</title>
  <meta name="description" content="Manage your restaurants" />
</svelte:head>

<div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
  <!-- Restaurant Header -->
  <RestaurantHeader
    {restaurants}
    {loading}
    {canCreateRestaurant}
    {checkingPermissions}
    {permissionReason}
    bind:viewMode
    onaddrestaurant={handleAddRestaurant}
  />

  <!-- Search and Filters -->
  <SearchAndFilters
    bind:searchQuery
    bind:statusFilter
    bind:sortBy
    totalRestaurants={restaurants.length}
    filteredCount={filteredRestaurants.length}
    onfilterschanged={handleFiltersChanged}
  />

  {#if error}
    <!-- Error State -->
    <Card>
      <div class="text-center py-12">
        <svg
          class="h-12 w-12 text-red-500 mx-auto mb-4"
          fill="none"
          stroke="currentColor"
          viewBox="0 0 24 24"
        >
          <path
            stroke-linecap="round"
            stroke-linejoin="round"
            stroke-width="2"
            d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-.833-1.964-.833-2.732 0L3.732 16.5c-.77.833.192 3 1.732 3z"
          />
        </svg>
        <h3 class="text-lg font-medium text-gray-900 mb-2">
          Failed to load restaurants
        </h3>
        <p class="text-gray-600 mb-4">{error}</p>
        <Button on:click={loadRestaurants}>Try Again</Button>
      </div>
    </Card>
  {:else}
    <!-- Restaurant List -->
    <RestaurantList
      restaurants={filteredRestaurants}
      {loading}
      {viewMode}
      on:restaurantClick={handleRestaurantClick}
      on:restaurantEdit={handleRestaurantEdit}
      on:restaurantToggleStatus={handleRestaurantToggleStatus}
      on:restaurantDelete={handleRestaurantDelete}
      on:addRestaurant={handleAddRestaurant}
      on:clearFilters={handleClearFilters}
    />
  {/if}
</div>

<!-- Create Restaurant Modal -->
<CreateRestaurantModal
  bind:isOpen={showCreateModal}
  onclose={handleCloseModal}
  onrestaurantcreated={handleRestaurantCreated}
/>

<!-- Edit Restaurant Modal -->
{#if showEditModal && editingRestaurant}
  <EditRestaurantModal
    restaurant={editingRestaurant}
    {clickOrigin}
    onclose={handleCloseEditModal}
    onupdated={handleRestaurantUpdated}
  />
{/if}

<!-- Delete Confirmation Dialog -->
<ConfirmDialog
  isOpen={showDeleteConfirm}
  title="Delete Restaurant"
  message={`Are you sure you want to delete "${deletingRestaurant?.name}"? This action cannot be undone.`}
  confirmText="Delete"
  cancelText="Cancel"
  variant="danger"
  {clickOrigin}
  onConfirm={handleDeleteConfirm}
  onCancel={handleDeleteCancel}
/>

