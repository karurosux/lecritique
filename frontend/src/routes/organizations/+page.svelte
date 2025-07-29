<script lang="ts">
  import { onMount } from "svelte";
  import { getApiClient, handleApiError } from "$lib/api/client";
  import { auth } from "$lib/stores/auth";
  import { subscription } from "$lib/stores/subscription";
  import { goto } from "$app/navigation";
  import OrganizationHeader from "$lib/components/organizations/OrganizationHeader.svelte";
  import SearchAndFilters from "$lib/components/organizations/SearchAndFilters.svelte";
  import OrganizationList from "$lib/components/organizations/OrganizationList.svelte";
  import CreateOrganizationModal from "$lib/components/organizations/CreateOrganizationModal.svelte";
  import EditOrganizationModal from "$lib/components/organizations/EditOrganizationModal.svelte";
  import { ConfirmDialog, Card, Button } from "$lib/components/ui";

  interface Organization {
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
  let organizations = $state<Organization[]>([]);
  let searchQuery = $state("");
  let statusFilter = $state("all"); // 'all', 'active', 'inactive'
  let sortBy = $state("name"); // 'name', 'created_at', 'status'
  let viewMode = $state<"grid" | "list">("grid");
  let showCreateModal = $state(false);
  let showEditModal = $state(false);
  let editingOrganization = $state<Organization | null>(null);
  let showDeleteConfirm = $state(false);
  let deletingOrganization = $state<Organization | null>(null);
  let canCreateOrganization = $state(false);
  let checkingPermissions = $state(true);
  let permissionReason = $state("");
  let clickOrigin = $state<{ x: number; y: number } | null>(null);

  let filteredOrganizations = $derived(
    organizations
      .filter((organization) => {
        if (!organization) return false;
        const matchesSearch =
          organization.name?.toLowerCase().includes(searchQuery.toLowerCase()) ||
          organization.description
            ?.toLowerCase()
            .includes(searchQuery.toLowerCase()) ||
          organization.email?.toLowerCase().includes(searchQuery.toLowerCase());
        const matchesStatus =
          statusFilter === "all" || organization.status === statusFilter;
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

  $effect(() => {
    if (!authState.isAuthenticated) {
      goto("/login");
      return;
    }
    
    if (authState.isAuthenticated && !hasInitialized) {
      hasInitialized = true;
      loadOrganizations();
    }
  });

  $effect(() => {
    if (authState.isAuthenticated && organizations.length >= 0) {
      checkCreatePermission();
    }
  });

  async function loadOrganizations() {
    loading = true;
    error = "";

    try {
      const api = getApiClient();

      // Use actual API client to get organizations
      const response = await api.api.v1OrganizationsList();

      if (response.data.success && response.data.data) {
        organizations = response.data.data.map((organization) => ({
          id: organization.id || "",
          name: organization.name || "",
          description: organization.description || "",
          address: organization.address || "",
          phone: organization.phone || "",
          email: organization.email || "",
          website: organization.website || "",
          cuisine_type: "", // Note: cuisine_type not in API model
          status: organization.is_active ? "active" : "inactive",
          created_at: organization.created_at || "",
          updated_at: organization.updated_at || "",
        }));
      } else {
        organizations = [];
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
      
      const response = await api.api.v1UserCanCreateOrganizationList();
      
      if (response.data.success) {
        const data = response.data.data;
        canCreateOrganization = data.can_create || false;
        
        if (!canCreateOrganization) {
          permissionReason = data.reason || "Cannot create more organizations";
        }
      } else {
        canCreateOrganization = false;
        permissionReason = "Unable to verify permissions";
      }
      
    } catch (err) {
      console.error("Error checking create permission:", err);
      canCreateOrganization = false;
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

  function handleOrganizationClick(event: CustomEvent) {
    const organization = event.detail;
    goto(`/organizations/${organization.id}`);
  }

  function handleOrganizationEdit(event: CustomEvent) {
    const { organization, event: mouseEvent } = event.detail;
    if (mouseEvent) {
      clickOrigin = { x: mouseEvent.clientX, y: mouseEvent.clientY };
    }
    editingOrganization = organization;
    showEditModal = true;
  }

  function handleOrganizationToggleStatus(event: CustomEvent) {
    const organization = event.detail;
    toggleOrganizationStatus(organization);
  }

  function handleOrganizationDelete(event: CustomEvent) {
    const { organization, event: mouseEvent } = event.detail;
    if (mouseEvent) {
      clickOrigin = { x: mouseEvent.clientX, y: mouseEvent.clientY };
    }
    deletingOrganization = organization;
    showDeleteConfirm = true;
  }

  function handleAddOrganization() {
    showCreateModal = true;
  }

  function handleCloseModal() {
    showCreateModal = false;
  }

  function handleCloseEditModal() {
    showEditModal = false;
    editingOrganization = null;
    clickOrigin = null;
  }

  function handleOrganizationUpdated() {
    showEditModal = false;
    editingOrganization = null;
    clickOrigin = null;
    loadOrganizations(); // Refresh the organization list
  }

  function handleDeleteConfirm() {
    if (deletingOrganization) {
      deleteOrganization(deletingOrganization);
    }
    showDeleteConfirm = false;
    deletingOrganization = null;
    clickOrigin = null;
  }

  function handleDeleteCancel() {
    showDeleteConfirm = false;
    deletingOrganization = null;
    clickOrigin = null;
  }

  function handleOrganizationCreated(event: CustomEvent) {
    // Reload the full list to ensure data consistency and get updated permissions
    loadOrganizations();
    // Refresh permissions since organization count changed
    checkCreatePermission();
  }

  function handleClearFilters() {
    searchQuery = '';
    statusFilter = 'all';
    sortBy = 'name';
  }

  async function toggleOrganizationStatus(organization: Organization) {
    try {
      const api = getApiClient();
      const newStatus = organization.status === "active" ? "inactive" : "active";

      // Update organization status via API
      await api.api.v1OrganizationsUpdate(organization.id, {
        is_active: newStatus === "active",
      });

      // Update local state
      organizations = organizations.map((r) =>
        r.id === organization.id
          ? { ...r, status: newStatus, updated_at: new Date().toISOString() }
          : r,
      );
    } catch (err) {
      error = handleApiError(err);
    }
  }

  async function deleteOrganization(organization: Organization) {
    try {
      const api = getApiClient();

      // Delete organization via API
      await api.api.v1OrganizationsDelete(organization.id);

      // Update local state
      organizations = organizations.filter((r) => r.id !== organization.id);
    } catch (err) {
      error = handleApiError(err);
    }
  }
</script>

<svelte:head>
  <title>Organizations - Kyooar</title>
  <meta name="description" content="Manage your organizations" />
</svelte:head>

<div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
  <!-- Organization Header -->
  <OrganizationHeader
    {organizations}
    {loading}
    {canCreateOrganization}
    {checkingPermissions}
    {permissionReason}
    bind:viewMode
    onaddorganization={handleAddOrganization}
  />

  <!-- Search and Filters -->
  <SearchAndFilters
    bind:searchQuery
    bind:statusFilter
    bind:sortBy
    totalOrganizations={organizations.length}
    filteredCount={filteredOrganizations.length}
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
          Failed to load organizations
        </h3>
        <p class="text-gray-600 mb-4">{error}</p>
        <Button on:click={loadOrganizations}>Try Again</Button>
      </div>
    </Card>
  {:else}
    <!-- Organization List -->
    <OrganizationList
      organizations={filteredOrganizations}
      {loading}
      {viewMode}
      on:organizationClick={handleOrganizationClick}
      on:organizationEdit={handleOrganizationEdit}
      on:organizationToggleStatus={handleOrganizationToggleStatus}
      on:organizationDelete={handleOrganizationDelete}
      on:addOrganization={handleAddOrganization}
      on:clearFilters={handleClearFilters}
    />
  {/if}
</div>

<!-- Create Organization Modal -->
<CreateOrganizationModal
  bind:isOpen={showCreateModal}
  onclose={handleCloseModal}
  onorganizationcreated={handleOrganizationCreated}
/>

<!-- Edit Organization Modal -->
{#if showEditModal && editingOrganization}
  <EditOrganizationModal
    organization={editingOrganization}
    {clickOrigin}
    onclose={handleCloseEditModal}
    onupdated={handleOrganizationUpdated}
  />
{/if}

<!-- Delete Confirmation Dialog -->
<ConfirmDialog
  isOpen={showDeleteConfirm}
  title="Delete Organization"
  message={`Are you sure you want to delete "${deletingOrganization?.name}"? This action cannot be undone.`}
  confirmText="Delete"
  cancelText="Cancel"
  variant="danger"
  {clickOrigin}
  onConfirm={handleDeleteConfirm}
  onCancel={handleDeleteCancel}
/>

