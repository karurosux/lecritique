<script lang="ts">
  import { getApiClient, handleApiError } from '$lib/api/client';
  import { auth } from '$lib/stores/auth';
  import { goto } from '$app/navigation';
  import { Card, Button, Select, Input } from '$lib/components/ui';
  import TimeSeriesChart from '$lib/components/analytics/TimeSeriesChart.svelte';
  import ComparisonChart from '$lib/components/analytics/ComparisonChart.svelte';
  import { TrendingUp, Calendar, BarChart3, Activity, Filter, RefreshCw } from 'lucide-svelte';

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
    { value: 'survey_responses', label: 'Total Survey Responses' },
    { value: 'rating_questions', label: 'Rating Questions' },
    { value: 'scale_questions', label: 'Scale Questions' }, 
    { value: 'yes_no_questions', label: 'Yes/No Questions' },
    { value: 'text_questions', label: 'Text Sentiment' },
    { value: 'single_choice_questions', label: 'Single Choice Questions' },
    { value: 'multiple_choice_questions', label: 'Multiple Choice Questions' }
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

  async function loadOrganizations() {
    try {
      const api = getApiClient();
      const response = await api.api.v1OrganizationsList();
      
      if (response.data.success && response.data.data) {
        organizations = response.data.data;
        if (organizations.length > 0) {
          selectedOrganization = organizations[0].id;
          await loadProducts();
          loadTimeSeriesData();
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
      const allQuestions = [];
      
      // Load questions for each product
      for (const product of availableProducts) {
        const response = await api.api.v1OrganizationsProductsQuestionsList(selectedOrganization, product.id);
        if (response.data.success && response.data.data) {
          response.data.data.forEach(question => {
            allQuestions.push({
              ...question,
              productName: product.name
            });
          });
        }
      }
      
      availableQuestions = allQuestions;
      console.log('Loaded questions:', allQuestions);
      
      // Group questions by product for visual organization
      const productGroups = {};
      allQuestions.forEach(q => {
        if (!productGroups[q.productName]) {
          productGroups[q.productName] = {
            questionTypes: [],
            individualQuestions: []
          };
        }
        
        // Add question type (with deduplication)
        const questionType = {
          value: `${q.type}_questions`,
          label: `${q.type.replace('_', ' ').replace(/\b\w/g, l => l.toUpperCase())} Questions`,
          type: q.type
        };
        
        const typeExists = productGroups[q.productName].questionTypes.some(qt => qt.value === questionType.value);
        if (!typeExists) {
          productGroups[q.productName].questionTypes.push(questionType);
        }
        
        // Add individual question
        productGroups[q.productName].individualQuestions.push({
          value: `question_${q.id}`,
          label: `${q.text}`,
          type: q.type,
          questionId: q.id
        });
      });
      
      console.log('Product groups:', productGroups);
      
      // Create a flat list for the current implementation (include both types and individual questions)
      const allMetrics = [];
      
      // Add question types
      const allQuestionTypes = new Set();
      Object.values(productGroups).forEach(group => {
        group.questionTypes.forEach(metric => allQuestionTypes.add(metric.value));
      });
      
      allMetrics.push(...Array.from(allQuestionTypes).map(type => ({
        value: type,
        label: type.replace('_questions', '').replace('_', ' ').replace(/\b\w/g, l => l.toUpperCase()) + ' Questions'
      })));
      
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
          'survey_responses',
          'rating_questions',
          'text_questions'
        ];
        
        comparisonFilters.metricTypes = [
          'survey_responses', 
          'rating_questions',
          'text_questions'
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
      
      console.log('Time series response:', response.data);
      console.log('Series data:', response.data?.series);
      
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
      console.log('Collecting metrics for organization:', selectedOrganization);
      await api.api.v1AnalyticsOrganizationsCollectMetricsCreate(selectedOrganization);
      console.log('Metrics collected successfully');
      
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
      loadTimeSeriesData();
      loadComparisonData();
    }
  }

  function handleTimeSeriesFilterChange() {
    loadTimeSeriesData();
  }

  function handleComparisonFilterChange() {
    loadComparisonData();
  }
</script>

<div class="advanced-analytics min-h-screen bg-gray-50 p-6">
  <div class="max-w-7xl mx-auto">
    <!-- Header -->
    <div class="mb-8">
      <div class="flex items-center justify-between">
        <div>
          <h1 class="text-3xl font-bold text-gray-900 flex items-center gap-3">
            <TrendingUp class="w-8 h-8 text-blue-600" />
            Advanced Survey Analytics
          </h1>
          <p class="text-gray-600 mt-2">Time series analysis and comparisons of survey responses</p>
        </div>
        
        <div class="flex gap-3">
          <Button variant="secondary" onclick={() => goto('/analytics')}>
            <BarChart3 class="w-4 h-4 mr-2" />
            Standard Analytics
          </Button>
          
          <Button variant="secondary" onclick={collectMetrics} disabled={collectingMetrics}>
            <RefreshCw class={`w-4 h-4 mr-2 ${collectingMetrics ? 'animate-spin' : ''}`} />
            {collectingMetrics ? 'Collecting...' : 'Collect Metrics'}
          </Button>
        </div>
      </div>
    </div>

    <!-- Organization Selection -->
    <Card class="mb-8 p-6">
      <div class="flex items-center gap-4">
        <div class="flex-1">
          <label class="block text-sm font-medium text-gray-700 mb-2">Organization</label>
          <Select 
            bind:value={selectedOrganization}
            onchange={handleOrganizationChange}
            options={organizations.map(org => ({ value: org.id, label: org.name }))}
            placeholder="Select an organization"
          />
        </div>
      </div>
    </Card>

    <!-- Time Series Section -->
    <div class="mb-12">
      <h2 class="text-2xl font-bold text-gray-900 mb-6 flex items-center gap-2">
        <BarChart3 class="w-6 h-6 text-green-600" />
        Survey Response Trends
      </h2>
      
      <!-- Time Series Filters -->
      <Card class="mb-6 p-6">
        <h3 class="text-lg font-semibold mb-4 flex items-center gap-2">
          <Filter class="w-5 h-5" />
          Filters
        </h3>
        <div class="space-y-4">
          <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-2">Start Date</label>
              <Input 
                type="date" 
                bind:value={timeSeriesFilters.startDate}
                onchange={handleTimeSeriesFilterChange}
              />
            </div>
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-2">End Date</label>
              <Input 
                type="date" 
                bind:value={timeSeriesFilters.endDate}
                onchange={handleTimeSeriesFilterChange}
              />
            </div>
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-2">Granularity</label>
              <Select 
                bind:value={timeSeriesFilters.granularity}
                onchange={handleTimeSeriesFilterChange}
                options={granularityOptions}
              />
            </div>
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-2">Product</label>
              <Select 
                bind:value={timeSeriesFilters.productId}
                onchange={handleTimeSeriesFilterChange}
                options={[{ value: '', label: 'All Products' }, ...availableProducts.map(p => ({ value: p.id, label: p.name }))]}
              />
            </div>
          </div>
          
          <!-- Metric Types Selection -->
          <div>
            <label class="block text-sm font-medium text-gray-700 mb-3">Metric Types</label>
            
            <!-- Global Metrics -->
            <div class="mb-4">
              <h4 class="text-sm font-medium text-gray-600 mb-2">General Metrics</h4>
              <div class="grid grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-3">
                <label class="flex items-center space-x-2 cursor-pointer">
                  <input
                    type="checkbox"
                    class="rounded border-gray-300 text-blue-600 shadow-sm focus:border-blue-500 focus:ring-blue-500"
                    checked={timeSeriesFilters.metricTypes.includes('survey_responses')}
                    onchange={(e) => {
                      if (e.target.checked) {
                        timeSeriesFilters.metricTypes = [...timeSeriesFilters.metricTypes, 'survey_responses'];
                      } else {
                        timeSeriesFilters.metricTypes = timeSeriesFilters.metricTypes.filter(t => t !== 'survey_responses');
                      }
                      handleTimeSeriesFilterChange();
                    }}
                  />
                  <span class="text-sm text-gray-700">Total Survey Responses</span>
                </label>
              </div>
            </div>

            <!-- Product-Grouped Metrics -->
            {#if Object.keys(availableProductGroups).length > 0}
              <div>
                <h4 class="text-sm font-medium text-gray-600 mb-2">Metrics by Product</h4>
                {#each Object.entries(availableProductGroups) as [productName, productData]}
                  <details class="mb-3">
                    <summary class="cursor-pointer text-sm font-medium text-gray-700 hover:text-gray-900 py-2 px-3 bg-gray-50 rounded-lg">
                      📊 {productName}
                    </summary>
                    <div class="mt-2 pl-4">
                      <!-- Question Types for this product -->
                      {#if productData.questionTypes.length > 0}
                        <div class="mb-4">
                          <h5 class="text-xs font-medium text-gray-500 mb-2 uppercase tracking-wide">Question Types</h5>
                          <div class="grid grid-cols-2 md:grid-cols-3 gap-2">
                            {#each productData.questionTypes as questionType}
                              <label class="flex items-center space-x-2 cursor-pointer">
                                <input
                                  type="checkbox"
                                  class="rounded border-gray-300 text-blue-600 shadow-sm focus:border-blue-500 focus:ring-blue-500"
                                  checked={timeSeriesFilters.metricTypes.includes(questionType.value)}
                                  onchange={(e) => {
                                    if (e.target.checked) {
                                      timeSeriesFilters.metricTypes = [...timeSeriesFilters.metricTypes, questionType.value];
                                    } else {
                                      timeSeriesFilters.metricTypes = timeSeriesFilters.metricTypes.filter(t => t !== questionType.value);
                                    }
                                    handleTimeSeriesFilterChange();
                                  }}
                                />
                                <span class="text-sm text-gray-700">{questionType.label}</span>
                              </label>
                            {/each}
                          </div>
                        </div>
                      {/if}
                      
                      <!-- Individual questions for this product -->
                      {#if productData.individualQuestions.length > 0}
                        <div>
                          <h5 class="text-xs font-medium text-gray-500 mb-2 uppercase tracking-wide">Individual Questions</h5>
                          <div class="grid grid-cols-1 gap-2">
                            {#each productData.individualQuestions as question}
                              <label class="flex items-center space-x-2 cursor-pointer">
                                <input
                                  type="checkbox"
                                  class="rounded border-gray-300 text-blue-600 shadow-sm focus:border-blue-500 focus:ring-blue-500"
                                  checked={timeSeriesFilters.metricTypes.includes(question.value)}
                                  onchange={(e) => {
                                    if (e.target.checked) {
                                      timeSeriesFilters.metricTypes = [...timeSeriesFilters.metricTypes, question.value];
                                    } else {
                                      timeSeriesFilters.metricTypes = timeSeriesFilters.metricTypes.filter(t => t !== question.value);
                                    }
                                    handleTimeSeriesFilterChange();
                                  }}
                                />
                                <span class="text-sm text-gray-600">{question.label}</span>
                              </label>
                            {/each}
                          </div>
                        </div>
                      {/if}
                    </div>
                  </details>
                {/each}
              </div>
            {/if}
          </div>
        </div>
      </Card>
      
      <!-- Time Series Chart -->
      <Card class="p-6">
        {#if loading}
          <div class="flex justify-center py-12">
            <div class="animate-spin rounded-full h-8 w-8 border-b-2 border-blue-600"></div>
          </div>
        {:else if timeSeriesData}
          <TimeSeriesChart data={timeSeriesData} />
        {:else}
          <div class="text-center py-12 text-gray-500">
            <BarChart3 class="w-12 h-12 mx-auto mb-4 opacity-50" />
            <p>No time series data available</p>
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
      <Card class="mb-6 p-6">
        <h3 class="text-lg font-semibold mb-4 flex items-center gap-2">
          <Calendar class="w-5 h-5" />
          Compare Periods
        </h3>
        <div class="grid grid-cols-1 lg:grid-cols-2 gap-8">
          <div>
            <h4 class="font-medium text-gray-900 mb-3">Period 1</h4>
            <div class="grid grid-cols-2 gap-4">
              <div>
                <label class="block text-sm font-medium text-gray-700 mb-2">Start Date</label>
                <Input 
                  type="date" 
                  bind:value={comparisonFilters.period1Start}
                  onchange={handleComparisonFilterChange}
                />
              </div>
              <div>
                <label class="block text-sm font-medium text-gray-700 mb-2">End Date</label>
                <Input 
                  type="date" 
                  bind:value={comparisonFilters.period1End}
                  onchange={handleComparisonFilterChange}
                />
              </div>
            </div>
          </div>
          
          <div>
            <h4 class="font-medium text-gray-900 mb-3">Period 2</h4>
            <div class="grid grid-cols-2 gap-4">
              <div>
                <label class="block text-sm font-medium text-gray-700 mb-2">Start Date</label>
                <Input 
                  type="date" 
                  bind:value={comparisonFilters.period2Start}
                  onchange={handleComparisonFilterChange}
                />
              </div>
              <div>
                <label class="block text-sm font-medium text-gray-700 mb-2">End Date</label>
                <Input 
                  type="date" 
                  bind:value={comparisonFilters.period2End}
                  onchange={handleComparisonFilterChange}
                />
              </div>
            </div>
          </div>
        </div>
        
        <!-- Comparison Metric Types Selection -->
        <div class="mt-6">
          <label class="block text-sm font-medium text-gray-700 mb-3">Metric Types for Comparison</label>
          
          <!-- Global Metrics -->
          <div class="mb-4">
            <h4 class="text-sm font-medium text-gray-600 mb-2">General Metrics</h4>
            <div class="grid grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-3">
              <label class="flex items-center space-x-2 cursor-pointer">
                <input
                  type="checkbox"
                  class="rounded border-gray-300 text-blue-600 shadow-sm focus:border-blue-500 focus:ring-blue-500"
                  checked={comparisonFilters.metricTypes.includes('survey_responses')}
                  onchange={(e) => {
                    if (e.target.checked) {
                      comparisonFilters.metricTypes = [...comparisonFilters.metricTypes, 'survey_responses'];
                    } else {
                      comparisonFilters.metricTypes = comparisonFilters.metricTypes.filter(t => t !== 'survey_responses');
                    }
                    handleComparisonFilterChange();
                  }}
                />
                <span class="text-sm text-gray-700">Total Survey Responses</span>
              </label>
            </div>
          </div>

          <!-- Product-Grouped Metrics -->
          {#if Object.keys(availableProductGroups).length > 0}
            <div>
              <h4 class="text-sm font-medium text-gray-600 mb-2">Metrics by Product</h4>
              {#each Object.entries(availableProductGroups) as [productName, productData]}
                <details class="mb-3">
                  <summary class="cursor-pointer text-sm font-medium text-gray-700 hover:text-gray-900 py-2 px-3 bg-gray-50 rounded-lg">
                    📊 {productName}
                  </summary>
                  <div class="mt-2 pl-4">
                    <!-- Question Types for this product -->
                    {#if productData.questionTypes.length > 0}
                      <div class="mb-4">
                        <h5 class="text-xs font-medium text-gray-500 mb-2 uppercase tracking-wide">Question Types</h5>
                        <div class="grid grid-cols-2 md:grid-cols-3 gap-2">
                          {#each productData.questionTypes as questionType}
                            <label class="flex items-center space-x-2 cursor-pointer">
                              <input
                                type="checkbox"
                                class="rounded border-gray-300 text-blue-600 shadow-sm focus:border-blue-500 focus:ring-blue-500"
                                checked={comparisonFilters.metricTypes.includes(questionType.value)}
                                onchange={(e) => {
                                  if (e.target.checked) {
                                    comparisonFilters.metricTypes = [...comparisonFilters.metricTypes, questionType.value];
                                  } else {
                                    comparisonFilters.metricTypes = comparisonFilters.metricTypes.filter(t => t !== questionType.value);
                                  }
                                  handleComparisonFilterChange();
                                }}
                              />
                              <span class="text-sm text-gray-700">{questionType.label}</span>
                            </label>
                          {/each}
                        </div>
                      </div>
                    {/if}
                    
                    <!-- Individual questions for this product -->
                    {#if productData.individualQuestions.length > 0}
                      <div>
                        <h5 class="text-xs font-medium text-gray-500 mb-2 uppercase tracking-wide">Individual Questions</h5>
                        <div class="grid grid-cols-1 gap-2">
                          {#each productData.individualQuestions as question}
                            <label class="flex items-center space-x-2 cursor-pointer">
                              <input
                                type="checkbox"
                                class="rounded border-gray-300 text-blue-600 shadow-sm focus:border-blue-500 focus:ring-blue-500"
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
                              <span class="text-sm text-gray-600">{question.label}</span>
                            </label>
                          {/each}
                        </div>
                      </div>
                    {/if}
                  </div>
                </details>
              {/each}
            </div>
          {/if}
        </div>
      </Card>
      
      <!-- Comparison Chart -->
      <Card class="p-6">
        {#if comparisonData}
          <ComparisonChart data={comparisonData} />
        {:else}
          <div class="text-center py-12 text-gray-500">
            <Activity class="w-12 h-12 mx-auto mb-4 opacity-50" />
            <p>No comparison data available</p>
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
</div>

<style>
  .advanced-analytics :global(.chart-container) {
    height: 400px;
  }
</style>