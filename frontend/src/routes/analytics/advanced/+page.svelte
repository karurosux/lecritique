<svelte:head>
  <title>Advanced Analytics - Kyooar</title>
</svelte:head>

<script lang="ts">
  import { getApiClient, handleApiError } from '$lib/api/client';
  import { auth } from '$lib/stores/auth';
  import { goto } from '$app/navigation';
  import { Card, Button, Select, Input } from '$lib/components/ui';
  import SeparatedTimeSeriesCharts from '$lib/components/analytics/SeparatedTimeSeriesCharts.svelte';
  import ComparisonChart from '$lib/components/analytics/ComparisonChart.svelte';
  import { TrendingUp, Calendar, BarChart3, Activity, Filter, RefreshCw, Building2 } from 'lucide-svelte';

  let loading = $state(true);
  let collectingMetrics = $state(false);
  let error = $state('');
  let organizations = $state<any[]>([]);
  let selectedOrganization = $state('');
  let timeSeriesData = $state<any>(null);
  let comparisonData = $state<any>(null);
  let hasInitialized = $state(false);
  
  // Time series filters
  let timeSeriesFilters = $state({
    startDate: new Date(Date.now() - 30 * 24 * 60 * 60 * 1000).toISOString().split('T')[0],
    endDate: new Date().toISOString().split('T')[0],
    granularity: 'daily',
    metricTypes: [], // Start empty, will be populated with questions
    productId: ''
  });
  
  // Comparison filters
  let comparisonFilters = $state({
    period1Start: new Date(Date.now() - 60 * 24 * 60 * 60 * 1000).toISOString().split('T')[0],
    period1End: new Date(Date.now() - 30 * 24 * 60 * 60 * 1000).toISOString().split('T')[0],
    period2Start: new Date(Date.now() - 30 * 24 * 60 * 60 * 1000).toISOString().split('T')[0],
    period2End: new Date().toISOString().split('T')[0],
    metricTypes: []
  });
  
  let availableProducts = $state<any[]>([]);
  let authState = $derived($auth);

  // Start with basic metrics, we'll add questions dynamically
  let metricTypeOptions = $state([
    { value: 'survey_responses', label: 'Total Survey Responses' }
  ]);
  
  // Store available questions from the API  
  let availableQuestions = $state<any[]>([]);
  let availableProductGroups = $state<any>({});

  const granularityOptions = [
    { value: 'hourly', label: 'Hourly' },
    { value: 'daily', label: 'Daily' },
    { value: 'weekly', label: 'Weekly' },
    { value: 'monthly', label: 'Monthly' }
  ];

  $effect(() => {
    if (!authState.isAuthenticated) {
      goto('/login');
      return;
    }
    
    if (authState.isAuthenticated && !hasInitialized) {
      hasInitialized = true;
      loadOrganizations();
    }
  });

  // Reactive effect to reload data when time series filters change
  $effect(() => {
    // Track specific filter changes
    timeSeriesFilters.startDate;
    timeSeriesFilters.endDate;
    timeSeriesFilters.granularity;
    timeSeriesFilters.metricTypes;
    timeSeriesFilters.productId;
    
    if (hasInitialized && selectedOrganization) {
      loadTimeSeriesData();
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
          loadComparisonData();
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
      const response = await api.api.v1OrganizationsProductsList(selectedOrganization);
      
      if (response.data.success && response.data.data) {
        availableProducts = response.data.data;
        // Load questions for the first product if any
        if (availableProducts.length > 0) {
          await loadQuestionsForProducts();
        }
      }
    } catch (err) {
      console.error('Error loading products:', err);
    }
  }
  
  async function loadQuestionsForProducts() {
    if (!selectedOrganization || availableProducts.length === 0) return;
    
    try {
      const api = getApiClient();
      const productIds = availableProducts.map(p => p.id);
      
      // Load questions for all products in a single batch request
      const response = await api.api.v1OrganizationsQuestionsBatchCreate(selectedOrganization, {
        product_ids: productIds
      });
      
      const allQuestions = [];
      if (response.data.success && response.data.data) {
        // Create a map of product ID to product name for quick lookup
        const productMap = {};
        availableProducts.forEach(product => {
          productMap[product.id] = product.name;
        });
        
        response.data.data.forEach(question => {
          allQuestions.push({
            ...question,
            productName: productMap[question.product_id] || 'Unknown Product'
          });
        });
      }
      
      availableQuestions = allQuestions;
      
      // Group questions by product for visual organization
      const productGroups = {};
      allQuestions.forEach(q => {
        if (!productGroups[q.productName]) {
          productGroups[q.productName] = {
            individualQuestions: []
          };
        }
        
        // Add individual question
        productGroups[q.productName].individualQuestions.push({
          value: `question_${q.id}`,
          label: `${q.text}`,
          type: q.type,
          questionId: q.id
        });
      });
      
      
      // Create a flat list for the current implementation (only individual questions)
      const allMetrics = [];
      
      // Add individual questions
      Object.values(productGroups).forEach(group => {
        allMetrics.push(...group.individualQuestions);
      });
      
      metricTypeOptions = [
        { value: 'survey_responses', label: 'Total Survey Responses' },
        ...allMetrics
      ];
      
      // Store product groups for the UI
      availableProductGroups = productGroups;
      
      if (timeSeriesFilters.metricTypes.length === 0) {
        timeSeriesFilters.metricTypes = [
          'survey_responses'
        ];
        
        comparisonFilters.metricTypes = [
          'survey_responses'
        ];
      }
    } catch (err) {
      console.error('Error loading questions:', err);
    }
  }

  async function loadTimeSeriesData() {
    if (!selectedOrganization) return;

    // Validate metric types
    if (!timeSeriesFilters.metricTypes || timeSeriesFilters.metricTypes.length === 0) {
      error = 'Please select at least one metric type';
      return;
    }

    loading = true;
    error = '';

    try {
      const api = getApiClient();
      const params = {
        metric_types: timeSeriesFilters.metricTypes,
        start_date: new Date(timeSeriesFilters.startDate).toISOString(),
        end_date: new Date(timeSeriesFilters.endDate).toISOString(),
        granularity: timeSeriesFilters.granularity,
        ...(timeSeriesFilters.productId && { product_id: timeSeriesFilters.productId })
      };


      const response = await api.api.v1AnalyticsOrganizationsTimeSeriesList(
        selectedOrganization,
        params
      );
      
      
      if (response.data) {
        timeSeriesData = response.data;
      }
    } catch (err) {
      error = handleApiError(err);
    } finally {
      loading = false;
    }
  }

  async function loadComparisonData() {
    if (!selectedOrganization) return;

    // Validate metric types
    if (!comparisonFilters.metricTypes || comparisonFilters.metricTypes.length === 0) {
      console.warn('No metric types selected for comparison');
      return;
    }

    try {
      const api = getApiClient();
      const requestBody = {
        organization_id: selectedOrganization,
        metric_types: comparisonFilters.metricTypes,
        period1_start: new Date(comparisonFilters.period1Start).toISOString(),
        period1_end: new Date(comparisonFilters.period1End).toISOString(),
        period2_start: new Date(comparisonFilters.period2Start).toISOString(),
        period2_end: new Date(comparisonFilters.period2End).toISOString()
      };

      const response = await api.api.v1AnalyticsOrganizationsCompareCreate(
        selectedOrganization,
        requestBody
      );
      
      if (response.data) {
        comparisonData = response.data;
      }
    } catch (err) {
      console.error('Error loading comparison data:', err);
    }
  }

  async function collectMetrics() {
    if (!selectedOrganization) return;

    collectingMetrics = true;
    error = '';

    try {
      const api = getApiClient();
      await api.api.v1AnalyticsOrganizationsCollectMetricsCreate(selectedOrganization);
      
      // Refresh data after collecting metrics
      await loadTimeSeriesData();
      await loadComparisonData();
    } catch (err) {
      console.error('Error collecting metrics:', err);
      error = handleApiError(err);
    } finally {
      collectingMetrics = false;
    }
  }

  function handleOrganizationChange() {
    if (selectedOrganization) {
      loadProducts();
      loadComparisonData();
    }
  }


  function handleComparisonFilterChange() {
    loadComparisonData();
  }
</script>

<div class="advanced-analytics max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
    <!-- Header -->
    <div class="mb-8 flex flex-col lg:flex-row lg:justify-between lg:items-center gap-6">
      <div class="space-y-3">
        <div class="flex items-center space-x-3">
          <div class="h-12 w-12 bg-gradient-to-br from-purple-500 to-pink-600 rounded-2xl flex items-center justify-center shadow-lg shadow-purple-500/25">
            <TrendingUp class="h-6 w-6 text-white" />
          </div>
          <div>
            <h1 class="text-3xl font-bold bg-gradient-to-r from-gray-900 to-gray-700 bg-clip-text text-transparent">
              Advanced Survey Analytics
            </h1>
            <p class="text-gray-600 font-medium">Time series analysis and comparisons of survey responses</p>
          </div>
        </div>
      </div>
      
      <div class="flex gap-3">
        <Button variant="secondary" size="lg" onclick={() => goto('/analytics')}>
          <BarChart3 class="w-4 h-4 mr-2" />
          Standard Analytics
        </Button>
        
        <Button variant="gradient" size="lg" onclick={collectMetrics} disabled={collectingMetrics}>
          <RefreshCw class={`w-4 h-4 mr-2 ${collectingMetrics ? 'animate-spin' : ''}`} />
          {collectingMetrics ? 'Collecting...' : 'Collect Metrics'}
        </Button>
      </div>
    </div>

    <!-- Organization Selection -->
    <Card variant="elevated" class="mb-8 overflow-hidden">
      <div class="bg-gradient-to-br from-purple-50 via-pink-50 to-indigo-50 p-6 rounded-2xl m-1">
        <div class="flex items-center justify-between gap-6">
          <div class="flex items-center gap-4">
            <div class="h-12 w-12 bg-gradient-to-br from-purple-500 to-pink-600 rounded-2xl flex items-center justify-center shadow-lg shadow-purple-500/25">
              <Building2 class="h-6 w-6 text-white" />
            </div>
            <div>
              <h3 class="text-lg font-semibold text-gray-900">Organization Context</h3>
              <p class="text-sm text-gray-600">Select the organization to analyze</p>
            </div>
          </div>
          
          <div class="flex-1 max-w-md">
            <Select 
              bind:value={selectedOrganization}
              onchange={handleOrganizationChange}
              options={organizations.map(org => ({ value: org.id, label: org.name }))}
              placeholder="Select an organization"
            />
          </div>
        </div>
        
        {#if selectedOrganization && organizations.find(org => org.id === selectedOrganization)}
          <div class="mt-4 pt-4 border-t border-purple-100/50 flex items-center gap-6">
            <span class="flex items-center gap-2 text-sm">
              <div class="h-8 w-8 bg-purple-100 rounded-lg flex items-center justify-center">
                <Activity class="w-4 h-4 text-purple-600" />
              </div>
              <span class="text-gray-600">Products: <span class="font-semibold text-gray-900">{availableProducts.length}</span></span>
            </span>
            <span class="flex items-center gap-2 text-sm">
              <div class="h-8 w-8 bg-pink-100 rounded-lg flex items-center justify-center">
                <BarChart3 class="w-4 h-4 text-pink-600" />
              </div>
              <span class="text-gray-600">Questions: <span class="font-semibold text-gray-900">{availableQuestions.length}</span></span>
            </span>
          </div>
        {/if}
      </div>
    </Card>

    <!-- Time Series Section -->
    <div class="mb-12">
      <h2 class="text-2xl font-bold text-gray-900 mb-6 flex items-center gap-2">
        <BarChart3 class="w-6 h-6 text-green-600" />
        Survey Response Trends
      </h2>
      
      <!-- Time Series Filters -->
      <Card variant="elevated" class="mb-6">
        <div class="p-6 border-b border-gray-100">
          <h3 class="text-lg font-semibold text-gray-900 flex items-center gap-2">
            <Filter class="w-5 h-5 text-blue-600" />
            Chart Configuration
          </h3>
          <p class="text-sm text-gray-600 mt-1">Configure date range, granularity, and select questions to analyze</p>
        </div>
        
        <div class="p-6 space-y-6">
          <!-- Essential Controls Row -->
          <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
            <div class="space-y-2">
              <label class="text-sm font-medium text-gray-700">Time Range</label>
              <div class="grid grid-cols-2 gap-2">
                <Input 
                  type="date" 
                  bind:value={timeSeriesFilters.startDate}
                  class="text-xs"
                />
                <Input 
                  type="date" 
                  bind:value={timeSeriesFilters.endDate}
                  class="text-xs"
                />
              </div>
            </div>
            
            <div class="space-y-2">
              <label class="text-sm font-medium text-gray-700">Granularity</label>
              <Select 
                bind:value={timeSeriesFilters.granularity}
                options={granularityOptions}
              />
            </div>
          </div>

          <!-- Question Selection -->
          {#if Object.keys(availableProductGroups).length > 0}
            <div class="border-t border-gray-100 pt-6">
              <div class="flex items-center justify-between mb-4">
                <h4 class="font-medium text-gray-900 flex items-center gap-2">
                  <BarChart3 class="w-4 h-4 text-blue-600" />
                  Select Questions to Analyze
                </h4>
                <div class="text-xs text-gray-500">
                  {timeSeriesFilters.metricTypes.filter(t => t.startsWith('question_')).length} questions selected
                </div>
              </div>
              
              <div class="space-y-3">
                {#each Object.entries(availableProductGroups) as [productName, productData]}
                  <Card variant="minimal" padding={false} class="overflow-hidden">
                    <details class="group">
                      <summary class="cursor-pointer flex items-center justify-between p-4 hover:bg-gray-50 transition-colors duration-200 list-none [&::-webkit-details-marker]:hidden">
                        <div class="flex items-center gap-3">
                          <div class="h-6 w-6 bg-gradient-to-br from-blue-500 to-indigo-600 rounded-md flex items-center justify-center">
                            <Activity class="w-3 h-3 text-white" />
                          </div>
                          <div>
                            <h5 class="font-medium text-gray-900">{productName}</h5>
                            <p class="text-xs text-gray-500">{productData.individualQuestions.length} question{productData.individualQuestions.length !== 1 ? 's' : ''}</p>
                          </div>
                        </div>
                        <div class="flex items-center gap-2">
                          <span class="text-xs bg-blue-100 text-blue-700 px-2 py-1 rounded-full font-medium">
                            {productData.individualQuestions.filter(q => timeSeriesFilters.metricTypes.includes(q.value)).length} selected
                          </span>
                          <svg class="w-4 h-4 text-gray-400 group-open:rotate-180 transition-transform duration-200" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 9l-7 7-7-7"></path>
                          </svg>
                        </div>
                      </summary>
                      
                      <div class="border-t border-gray-100 bg-gray-50/50 p-4">
                        <div class="space-y-2 max-h-48 overflow-y-auto">
                          {#each productData.individualQuestions as question}
                            <label class="flex items-center gap-2 cursor-pointer hover:bg-white p-2 rounded-md transition-colors">
                              <input
                                type="checkbox"
                                class="rounded border-gray-300 text-blue-600 shadow-sm focus:border-blue-500 focus:ring-blue-500 focus:ring-offset-0"
                                checked={timeSeriesFilters.metricTypes.includes(question.value)}
                                onchange={(e) => {
                                  if (e.target.checked) {
                                    timeSeriesFilters.metricTypes = [...timeSeriesFilters.metricTypes, question.value];
                                  } else {
                                    timeSeriesFilters.metricTypes = timeSeriesFilters.metricTypes.filter(t => t !== question.value);
                                  }
                                }}
                              />
                              <div class="flex-1 min-w-0">
                                <span class="text-sm text-gray-900 block truncate">{question.label}</span>
                              </div>
                            </label>
                          {/each}
                        </div>
                      </div>
                    </details>
                  </Card>
                {/each}
              </div>
            </div>
          {/if}
        </div>
      </Card>
      
      <!-- Time Series Chart -->
      <Card class="p-6">
        {#if loading}
          <div class="flex justify-center py-12">
            <div class="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600"></div>
          </div>
        {:else if timeSeriesData}
          <SeparatedTimeSeriesCharts data={timeSeriesData} />
        {:else}
          <div class="text-center py-12">
            <div class="bg-gradient-to-br from-purple-50 to-pink-50 rounded-2xl p-8 max-w-md mx-auto">
              <BarChart3 class="w-12 h-12 mx-auto mb-4 text-purple-400" />
              <h3 class="text-lg font-semibold text-gray-900 mb-2">No Analytics Data Yet</h3>
              <p class="text-gray-600 mb-6">Collect metrics to start analyzing your survey responses over time.</p>
              <Button variant="gradient" size="lg" onclick={collectMetrics} disabled={collectingMetrics}>
                <RefreshCw class={`w-4 h-4 mr-2 ${collectingMetrics ? 'animate-spin' : ''}`} />
                {collectingMetrics ? 'Collecting Metrics...' : 'Collect Metrics Now'}
              </Button>
            </div>
          </div>
        {/if}
      </Card>
    </div>

    <!-- Comparison Section -->
    <div>
      <h2 class="text-2xl font-bold text-gray-900 mb-6 flex items-center gap-2">
        <Activity class="w-6 h-6 text-purple-600" />
        Period Comparison
      </h2>
      
      <!-- Comparison Filters -->
      <Card variant="elevated" class="mb-6">
        <div class="p-6 border-b border-gray-100">
          <h3 class="text-lg font-semibold text-gray-900 flex items-center gap-2">
            <Calendar class="w-5 h-5 text-purple-600" />
            Period Comparison Setup
          </h3>
          <p class="text-sm text-gray-600 mt-1">Compare metrics between two time periods</p>
        </div>
        
        <div class="p-6 space-y-6">
          <!-- Period Configuration -->
          <div class="grid grid-cols-1 lg:grid-cols-2 gap-8">
            <Card variant="minimal" class="p-4">
              <div class="flex items-center gap-2 mb-3">
                <div class="h-6 w-6 bg-gradient-to-br from-blue-500 to-blue-600 rounded-md flex items-center justify-center">
                  <span class="text-xs font-bold text-white">1</span>
                </div>
                <h4 class="font-medium text-gray-900">First Period</h4>
              </div>
              <div class="grid grid-cols-2 gap-3">
                <div>
                  <label class="text-xs font-medium text-gray-600 mb-1 block">Start Date</label>
                  <Input 
                    type="date" 
                    bind:value={comparisonFilters.period1Start}
                    onchange={handleComparisonFilterChange}
                    class="text-sm"
                  />
                </div>
                <div>
                  <label class="text-xs font-medium text-gray-600 mb-1 block">End Date</label>
                  <Input 
                    type="date" 
                    bind:value={comparisonFilters.period1End}
                    onchange={handleComparisonFilterChange}
                    class="text-sm"
                  />
                </div>
              </div>
            </Card>
            
            <Card variant="minimal" class="p-4">
              <div class="flex items-center gap-2 mb-3">
                <div class="h-6 w-6 bg-gradient-to-br from-purple-500 to-purple-600 rounded-md flex items-center justify-center">
                  <span class="text-xs font-bold text-white">2</span>
                </div>
                <h4 class="font-medium text-gray-900">Second Period</h4>
              </div>
              <div class="grid grid-cols-2 gap-3">
                <div>
                  <label class="text-xs font-medium text-gray-600 mb-1 block">Start Date</label>
                  <Input 
                    type="date" 
                    bind:value={comparisonFilters.period2Start}
                    onchange={handleComparisonFilterChange}
                    class="text-sm"
                  />
                </div>
                <div>
                  <label class="text-xs font-medium text-gray-600 mb-1 block">End Date</label>
                  <Input 
                    type="date" 
                    bind:value={comparisonFilters.period2End}
                    onchange={handleComparisonFilterChange}
                    class="text-sm"
                  />
                </div>
              </div>
            </Card>
          </div>

          <!-- Question Selection -->
          {#if Object.keys(availableProductGroups).length > 0}
            <div class="border-t border-gray-100 pt-6">
              <div class="flex items-center justify-between mb-4">
                <h4 class="font-medium text-gray-900 flex items-center gap-2">
                  <BarChart3 class="w-4 h-4 text-purple-600" />
                  Select Questions to Compare
                </h4>
                <div class="text-xs text-gray-500">
                  {comparisonFilters.metricTypes.filter(t => t.startsWith('question_')).length} questions selected
                </div>
              </div>
              
              <div class="space-y-3">
                {#each Object.entries(availableProductGroups) as [productName, productData]}
                  <Card variant="minimal" padding={false} class="overflow-hidden">
                    <details class="group">
                      <summary class="cursor-pointer flex items-center justify-between p-4 hover:bg-gray-50 transition-colors duration-200 list-none [&::-webkit-details-marker]:hidden">
                        <div class="flex items-center gap-3">
                          <div class="h-6 w-6 bg-gradient-to-br from-purple-500 to-pink-600 rounded-md flex items-center justify-center">
                            <Activity class="w-3 h-3 text-white" />
                          </div>
                          <div>
                            <h5 class="font-medium text-gray-900">{productName}</h5>
                            <p class="text-xs text-gray-500">{productData.individualQuestions.length} question{productData.individualQuestions.length !== 1 ? 's' : ''}</p>
                          </div>
                        </div>
                        <div class="flex items-center gap-2">
                          <span class="text-xs bg-purple-100 text-purple-700 px-2 py-1 rounded-full font-medium">
                            {productData.individualQuestions.filter(q => comparisonFilters.metricTypes.includes(q.value)).length} selected
                          </span>
                          <svg class="w-4 h-4 text-gray-400 group-open:rotate-180 transition-transform duration-200" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 9l-7 7-7-7"></path>
                          </svg>
                        </div>
                      </summary>
                      
                      <div class="border-t border-gray-100 bg-gray-50/50 p-4">
                        <div class="space-y-2 max-h-48 overflow-y-auto">
                          {#each productData.individualQuestions as question}
                            <label class="flex items-center gap-2 cursor-pointer hover:bg-white p-2 rounded-md transition-colors">
                              <input
                                type="checkbox"
                                class="rounded border-gray-300 text-purple-600 shadow-sm focus:border-purple-500 focus:ring-purple-500 focus:ring-offset-0"
                                checked={comparisonFilters.metricTypes.includes(question.value)}
                                onchange={(e) => {
                                  if (e.target.checked) {
                                    comparisonFilters.metricTypes = [...comparisonFilters.metricTypes, question.value];
                                  } else {
                                    comparisonFilters.metricTypes = comparisonFilters.metricTypes.filter(t => t !== question.value);
                                  }
                                  handleComparisonFilterChange();
                                }}
                              />
                              <div class="flex-1 min-w-0">
                                <span class="text-sm text-gray-900 block truncate">{question.label}</span>
                              </div>
                            </label>
                          {/each}
                        </div>
                      </div>
                    </details>
                  </Card>
                {/each}
              </div>
            </div>
          {/if}
        </div>
      </Card>
      
      <!-- Comparison Chart -->
      <Card class="p-6">
        {#if comparisonData}
          <ComparisonChart data={comparisonData} />
        {:else}
          <div class="text-center py-12">
            <div class="bg-gradient-to-br from-purple-50 to-pink-50 rounded-2xl p-8 max-w-md mx-auto">
              <Activity class="w-12 h-12 mx-auto mb-4 text-purple-400" />
              <h3 class="text-lg font-semibold text-gray-900 mb-2">No Comparison Data Yet</h3>
              <p class="text-gray-600 mb-6">Collect metrics to compare performance between different time periods.</p>
              <Button variant="gradient" size="lg" onclick={collectMetrics} disabled={collectingMetrics}>
                <RefreshCw class={`w-4 h-4 mr-2 ${collectingMetrics ? 'animate-spin' : ''}`} />
                {collectingMetrics ? 'Collecting Metrics...' : 'Collect Metrics Now'}
              </Button>
            </div>
          </div>
        {/if}
      </Card>
    </div>

    <!-- Error Display -->
    {#if error}
      <div class="fixed bottom-4 right-4 bg-red-500 text-white px-4 py-2 rounded-lg shadow-lg">
        {error}
      </div>
    {/if}
</div>

<style>
  .advanced-analytics :global(.chart-container) {
    height: 400px;
  }
</style>