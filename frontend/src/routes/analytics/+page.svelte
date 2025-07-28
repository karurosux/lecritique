<script lang="ts">
  import { getApiClient, handleApiError } from "$lib/api/client";
  import { auth } from "$lib/stores/auth";
  import { goto } from "$app/navigation";
  import {
    Card,
    Button,
    Select,
    Input,
    NoDataAvailable,
  } from "$lib/components/ui";
  import ChartDataWidgetGrouped from "$lib/components/analytics/ChartDataWidgetGrouped.svelte";
  import QRCodeAnalytics from "$lib/components/analytics/QRCodeAnalytics.svelte";
  import { FeatureGate, FEATURES } from "$lib/components/subscription";
  import {
    BarChart3,
    Activity,
    Package,
    Calendar,
    Search,
    Layers,
    Grid,
    RefreshCw,
    AlertTriangle,
    Building2,
  } from "lucide-svelte";

  const DaysFilters = {
    WEEK: "7",
    MONTH: "30",
    QUEARTER: "90",
    ALL_TIME: "all",
  } as const;

  type DaysFiltersType = (typeof DaysFilters)[keyof typeof DaysFilters];

  let loading = $state(true);
  let error = $state("");
  let organizations = $state<any[]>([]);
  let selectedOrganization = $state("");
  let chartData = $state<any>(null);
  let analyticsData = $state<any>(null);
  let hasInitialized = $state(false);

  let filters = $state<{ productId: string; days: DaysFiltersType }>({
    days: DaysFilters.WEEK, // 'all', '7', '30', '90'
    productId: "",
  });

  let searchQuery = $state("");
  let showOnlyWithData = $state(false);
  let viewMode = $state<"grouped" | "all">("grouped");

  let authState = $derived($auth);
  let availableProducts = $state<any[]>([]);

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

  async function loadOrganizations() {
    try {
      const api = getApiClient();
      const response = await api.api.v1OrganizationsList();

      if (response.data.success && response.data.data) {
        organizations = response.data.data;
        if (organizations.length > 0) {
          selectedOrganization = organizations[0].id;
          await loadProducts();
          loadAnalytics();
        }
      }
    } catch (err) {
      error = handleApiError(err);
    }
  }

  async function loadProducts() {
    if (!selectedOrganization) return;

    try {
      const api = getApiClient();
      const response =
        await api.api.v1OrganizationsProductsList(selectedOrganization);

      if (response.data.success && response.data.data) {
        availableProducts = response.data.data;
      }
    } catch (err) {
      console.error("Error loading products:", err);
    }
  }

  async function loadAnalytics() {
    if (!selectedOrganization) return;

    loading = true;
    error = "";

    try {
      const api = getApiClient();

      const chartParams: Record<string, any> = {};

      if (filters.days !== DaysFilters.ALL_TIME) {
        const today = new Date();
        const daysAgo = new Date(
          today.getTime() - parseInt(filters.days) * 24 * 60 * 60 * 1000,
        );
        chartParams.date_from = daysAgo.toISOString().split("T")[0];
      }

      if (filters.productId) {
        chartParams.product_id = filters.productId;
      }

      const [analyticsResponse, chartResponse, dashboardResponse] =
        await Promise.all([
          api.api.v1AnalyticsOrganizationsDetail(selectedOrganization),
          api.api.v1AnalyticsOrganizationsChartsList(
            selectedOrganization,
            chartParams,
          ),
          api.api.v1AnalyticsDashboardDetail(selectedOrganization),
        ]);

      if (analyticsResponse.data?.data) {
        analyticsData = analyticsResponse.data.data;
      }

      if (chartResponse.data?.data) {
        chartData = chartResponse.data.data;
      } else {
        chartData = null;
      }

      if (dashboardResponse.data?.data) {
        analyticsData = {
          ...analyticsData,
          qr_performance: dashboardResponse.data.data.qr_performance,
          total_qr_scans: dashboardResponse.data.data.total_qr_scans,
          total_active_codes: dashboardResponse.data.data.active_qr_codes,
          average_conversion_rate: dashboardResponse.data.data.completion_rate,
        };
      }
    } catch (err) {
      console.error("Error loading analytics:", err);
      error = handleApiError(err);
    } finally {
      loading = false;
    }
  }

  function handleOrganizationChange() {
    resetFilters();
    loadProducts();
    loadAnalytics();
  }

  function resetFilters() {
    filters.days = DaysFilters.WEEK;
    filters.productId = "";
  }

  function applyFilters() {
    loadAnalytics();
  }
</script>

<svelte:head>
  <title>Product Analytics - Kyooar</title>
</svelte:head>

<div class="analytics-page max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
  <!-- Analytics Header -->
  <div class="mb-8">
    <div
      class="flex flex-col lg:flex-row lg:justify-between lg:items-center gap-6"
    >
      <div class="space-y-3">
        <div class="flex items-center space-x-3">
          <div
            class="h-12 w-12 bg-gradient-to-br from-purple-500 to-pink-600 rounded-2xl flex items-center justify-center shadow-lg shadow-purple-500/25"
          >
            <Package class="h-6 w-6 text-white" />
          </div>
          <div>
            <h1
              class="text-3xl font-bold bg-gradient-to-r from-gray-900 to-gray-700 bg-clip-text text-transparent"
            >
              Product Performance Analytics
            </h1>
            <p class="text-gray-600 font-medium">
              Detailed insights on how your products are performing
            </p>
          </div>
        </div>
      </div>

      <div class="flex items-center space-x-3">
        <FeatureGate feature={FEATURES.ADVANCED_ANALYTICS}>
          <Button
            variant="secondary"
            size="lg"
            onclick={() => goto("/analytics/advanced")}
          >
            <BarChart3 class="h-5 w-5 mr-2" />
            Advanced Analytics
          </Button>
        </FeatureGate>
      </div>
    </div>
  </div>

  <!-- Analytics Controls -->
  <div class="analytics-controls mb-6">
    <Card variant="elevated" padding={false}>
      <div class="divide-y divide-gray-200">
        <!-- Primary Controls Row -->
        <div class="p-4">
          <div
            class="flex flex-col lg:flex-row items-start lg:items-center gap-4"
          >
            <!-- Left Side: Organization & Time -->
            <div class="flex flex-wrap items-center gap-3 flex-1">
              {#if organizations.length > 0}
                <div class="flex items-center gap-2">
                  <Activity class="h-4 w-4 text-gray-500" />
                  <Select
                    bind:value={selectedOrganization}
                    options={organizations.map((r) => ({
                      value: r.id,
                      label: r.name,
                    }))}
                    onchange={handleOrganizationChange}
                    minWidth="min-w-48"
                  />
                  {#if chartData?.summary?.total_responses}
                    <span class="text-sm text-gray-500 hidden sm:inline">
                      ({chartData.summary.total_responses} responses)
                    </span>
                  {/if}
                </div>
              {/if}

              <div class="flex items-center gap-2">
                <Calendar class="h-4 w-4 text-gray-500" />
                <div class="flex items-center gap-1 bg-gray-100 rounded-lg p-1">
                  <button
                    class="px-3 py-1.5 text-sm font-medium rounded-md transition-all {filters.days ===
                    DaysFilters.WEEK
                      ? 'bg-white text-gray-900 shadow-sm'
                      : 'text-gray-600 hover:text-gray-900'}"
                    onclick={() => {
                      filters.days = "7";
                      applyFilters();
                    }}
                  >
                    Week
                  </button>
                  <button
                    class="px-3 py-1.5 text-sm font-medium rounded-md transition-all {filters.days ===
                    DaysFilters.MONTH
                      ? 'bg-white text-gray-900 shadow-sm'
                      : 'text-gray-600 hover:text-gray-900'}"
                    onclick={() => {
                      filters.days = "30";
                      applyFilters();
                    }}
                  >
                    Month
                  </button>
                  <button
                    class="px-3 py-1.5 text-sm font-medium rounded-md transition-all {filters.days ===
                    DaysFilters.QUEARTER
                      ? 'bg-white text-gray-900 shadow-sm'
                      : 'text-gray-600 hover:text-gray-900'}"
                    onclick={() => {
                      filters.days = "90";
                      applyFilters();
                    }}
                  >
                    Quarter
                  </button>
                  <button
                    class="px-3 py-1.5 text-sm font-medium rounded-md transition-all {filters.days ===
                    DaysFilters.ALL_TIME
                      ? 'bg-white text-gray-900 shadow-sm'
                      : 'text-gray-600 hover:text-gray-900'}"
                    onclick={() => {
                      filters.days = "all";
                      applyFilters();
                    }}
                  >
                    All Time
                  </button>
                </div>
              </div>
            </div>

            <!-- Right Side: Refresh -->
            <!-- svelte-ignore a11y_consider_explicit_label -->
            <button
              onclick={loadAnalytics}
              disabled={loading}
              class="p-2 text-gray-600 hover:text-gray-900 hover:bg-gray-100 rounded-lg transition-all disabled:opacity-50 group"
              title="Refresh data"
            >
              <RefreshCw
                class="h-5 w-5 {loading
                  ? 'animate-spin'
                  : 'group-hover:rotate-180 transition-transform duration-300'}"
              />
            </button>
          </div>
        </div>

        <!-- Secondary Controls Row -->
        <div class="p-4 bg-gray-50">
          <div
            class="flex flex-col lg:flex-row items-start lg:items-center gap-4"
          >
            <!-- Left Side: Search & Product Filter -->
            <div class="flex flex-wrap items-center gap-3 flex-1">
              <div class="relative">
                <Search
                  class="absolute left-3 top-1/2 transform -translate-y-1/2 h-4 w-4 text-gray-400"
                />
                <Input
                  type="text"
                  bind:value={searchQuery}
                  placeholder="Search products or questions..."
                  class="pl-9 w-64"
                />
              </div>

              <div class="flex items-center gap-2">
                <Package class="h-4 w-4 text-gray-500" />
                <Select
                  bind:value={filters.productId}
                  options={[
                    { value: "", label: "All Products" },
                    ...availableProducts.map((d) => ({
                      value: d.id,
                      label: d.name,
                    })),
                  ]}
                  onchange={applyFilters}
                  minWidth="min-w-40"
                />
              </div>

              <label class="flex items-center gap-2 cursor-pointer">
                <input
                  type="checkbox"
                  bind:checked={showOnlyWithData}
                  class="w-4 h-4 text-purple-600 bg-gray-100 border-gray-300 rounded focus:ring-purple-500"
                />
                <span class="text-sm font-medium text-gray-700"
                  >With data only</span
                >
              </label>
            </div>

            <!-- Right Side: View Mode -->
            <div
              class="flex items-center gap-1 bg-white border border-gray-300 rounded-lg p-1"
            >
              <button
                class="px-3 py-1.5 text-sm font-medium rounded-md transition-all {viewMode ===
                'grouped'
                  ? 'bg-gray-900 text-white'
                  : 'text-gray-600 hover:text-gray-900 hover:bg-gray-100'}"
                onclick={() => (viewMode = "grouped")}
              >
                <Layers class="h-4 w-4 inline mr-1.5" />
                Grouped
              </button>
              <button
                class="px-3 py-1.5 text-sm font-medium rounded-md transition-all {viewMode ===
                'all'
                  ? 'bg-gray-900 text-white'
                  : 'text-gray-600 hover:text-gray-900 hover:bg-gray-100'}"
                onclick={() => (viewMode = "all")}
              >
                <Grid class="h-4 w-4 inline mr-1.5" />
                All
              </button>
            </div>
          </div>
        </div>
      </div>
    </Card>
  </div>

  {#if loading}
    <!-- Loading State -->
    <div class="space-y-8">
      <!-- Summary Cards Loading -->
      <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
        {#each Array(4) as _}
          <Card variant="gradient">
            <div class="animate-pulse">
              <div class="flex items-center justify-between mb-4">
                <div class="h-4 bg-gray-200 rounded w-1/2"></div>
                <div class="h-10 w-10 bg-gray-200 rounded-xl"></div>
              </div>
              <div class="h-8 bg-gray-200 rounded w-3/4 mb-2"></div>
              <div class="h-4 bg-gray-200 rounded w-1/3"></div>
            </div>
          </Card>
        {/each}
      </div>

      <!-- Charts Loading -->
      <div class="grid grid-cols-1 lg:grid-cols-2 gap-6">
        {#each Array(2) as _}
          <Card variant="elevated">
            <div class="animate-pulse">
              <div class="h-6 bg-gray-200 rounded w-1/3 mb-4"></div>
              <div class="h-64 bg-gray-100 rounded"></div>
            </div>
          </Card>
        {/each}
      </div>
    </div>
  {:else if error}
    <!-- Error State -->
    <NoDataAvailable
      title="Failed to load analytics"
      description={error}
      icon={AlertTriangle}
    />
  {:else if organizations.length === 0}
    <!-- No Organizations State -->
    <NoDataAvailable
      title="No Organizations Yet"
      description="Create your first organization to start collecting feedback and viewing analytics."
      icon={Building2}
    />
  {:else}
    <!-- Analytics Content -->
    <div class="space-y-8">
      <!-- Chart Analytics with Grouping -->
      <ChartDataWidgetGrouped
        {chartData}
        title=""
        groupByProduct={true}
        initiallyExpanded={false}
        bind:searchQuery
        bind:showOnlyWithData
        bind:viewMode
        hideControls={true}
      />

      {#if chartData?.feedbacks?.length > 0}
        <QRCodeAnalytics
          {analyticsData}
          feedbacks={chartData.feedbacks}
          {loading}
        />
      {/if}
    </div>
  {/if}
</div>
