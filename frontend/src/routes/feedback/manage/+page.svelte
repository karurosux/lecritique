<script lang="ts">
  import { onMount } from 'svelte';
  import { page } from '$app/stores';
  import { Card, Button, Input, Select } from '$lib/components/ui';
  import { getApiClient, handleApiError } from '$lib/api/client';
  import { auth } from '$lib/stores/auth';
  import { goto } from '$app/navigation';
  import { RoleGate } from '$lib/components/auth';

  interface Feedback {
    id: string;
    customer_email?: string;
    rating: number;
    comment?: string;
    product_name?: string;
    organization_name?: string;
    location_name?: string;
    qr_code?: string;
    responses?: Record<string, any>;
    created_at: string;
  }

  interface FeedbackFilters {
    organization_id?: string;
    location_id?: string;
    rating_min?: number;
    rating_max?: number;
    date_from?: string;
    date_to?: string;
    search?: string;
  }

  let loading = $state(true);
  let error = $state('');
  let feedback = $state<Feedback[]>([]);
  let organizations = $state<any[]>([]);
  let products = $state<any[]>([]);
  let totalCount = $state(0);
  let currentPage = $state(1);
  let itemsPerPage = $state(20);
  
  // Filters
  let filters: FeedbackFilters = {};
  let searchQuery = $state('');
  let searchInput = $state(''); // Make it reactive for binding
  let selectedOrganization = $state('');
  let selectedProduct = $state('');
  let selectedRating = $state('');
  let dateFrom = $state('');
  let dateTo = $state('');

  // Track collapsed state for each feedback item
  let collapsedStates = $state<Record<string, boolean>>({});
  
  function toggleCollapse(feedbackId: string) {
    // If undefined, we want to set it to false (expanded) since default is true (collapsed)
    collapsedStates[feedbackId] = collapsedStates[feedbackId] === undefined ? false : !collapsedStates[feedbackId];
  }
  
  function isCollapsed(feedbackId: string): boolean {
    // Default to collapsed (true) if not set
    return collapsedStates[feedbackId] ?? true;
  }

  // Set default 15-day filter
  function setDefaultDateFilter() {
    const today = new Date();
    const fifteenDaysAgo = new Date(today);
    fifteenDaysAgo.setDate(today.getDate() - 15);
    
    dateFrom = fifteenDaysAgo.toISOString().split('T')[0];
    dateTo = today.toISOString().split('T')[0];
  }

  let isFirstLoad = $state(true);
  let searchTimeout: ReturnType<typeof setTimeout> | null = null;
  
  // Debounced search function
  function handleSearchInput() {
    if (!isFirstLoad) {
      // Clear existing timeout
      if (searchTimeout) {
        clearTimeout(searchTimeout);
      }
      
      // Set new timeout for debounced search
      searchTimeout = setTimeout(() => {
        searchQuery = searchInput; // Update reactive value after debounce
        loadFeedback();
      }, 500); // 500ms debounce
    }
  }

  // Handle other filter changes immediately (non-search filters)
  $effect(() => {
    if (!isFirstLoad && (selectedRating !== undefined || selectedOrganization !== undefined || selectedProduct !== undefined || dateFrom !== undefined || dateTo !== undefined)) {
      loadFeedback();
    }
  });

  // Reload products when organization changes
  $effect(() => {
    if (!isFirstLoad && selectedOrganization !== undefined) {
      selectedProduct = ''; // Clear product selection when organization changes
      loadProducts();
    }
  });

  let authState = $derived($auth);

  onMount(async () => {
    if (!authState.isAuthenticated) {
      goto('/login');
      return;
    }

    // Set default 15-day filter
    setDefaultDateFilter();
    
    const loadedOrganizations = await loadOrganizations();
    // Load products after organizations are loaded
    if (loadedOrganizations.length > 0) {
      await loadProducts();
    }
    await loadFeedback();
    
    // Enable filter changes after initial load
    isFirstLoad = false;
  });

  async function loadOrganizations() {
    try {
      const api = getApiClient();
      const response = await api.api.v1OrganizationsList();
      
      if (response.data.success && response.data.data) {
        organizations = response.data.data;
        return response.data.data; // Return the organizations
      }
      return [];
    } catch (err) {
      console.error('Error loading organizations:', err);
      return [];
    }
  }

  async function loadProducts() {
    try {
      const api = getApiClient();
      
      // Make sure we have organizations loaded first
      if (organizations.length === 0 && !selectedOrganization) {
        return;
      }
      
      // If a organization is selected, load products for that organization
      if (selectedOrganization) {
        const response = await api.api.v1OrganizationsProductsList(selectedOrganization);
        if (response.data.success && response.data.data) {
          products = response.data.data;
        }
      } else {
        // Load products from all organizations
        const productPromises = organizations.map(async (organization) => {
          try {
            const response = await api.api.v1OrganizationsProductsList(organization.id);
            return response.data.data || [];
          } catch (err) {
            console.error(`Error loading products for organization ${organization.id}:`, err);
            return [];
          }
        });
        
        const productArrays = await Promise.all(productPromises);
        const allProducts = productArrays.flat();
        
        // Remove duplicates by name
        const uniqueProducts = allProducts.reduce((acc: any[], product: any) => {
          if (!acc.find(d => d.name === product.name)) {
            acc.push(product);
          }
          return acc;
        }, []);
        
        products = uniqueProducts.sort((a, b) => a.name.localeCompare(b.name));
      }
    } catch (err) {
      console.error('Error loading products:', err);
    }
  }

  async function loadFeedback() {
    loading = true;
    error = '';

    try {
      const api = getApiClient();
      
      if (organizations.length === 0) {
        feedback = [];
        totalCount = 0;
        loading = false;
        return;
      }

      // Load feedback from all organizations if no specific organization selected
      let allFeedback: Feedback[] = [];
      
      if (selectedOrganization) {
        // Load feedback from specific organization with server-side filtering
        const query: any = {};
        
        if (searchQuery) query.search = searchQuery;
        if (selectedRating) {
          query.rating_min = parseInt(selectedRating);
          query.rating_max = parseInt(selectedRating);
        }
        if (selectedProduct) query.product_id = selectedProduct;
        if (dateFrom) query.date_from = dateFrom;
        if (dateTo) query.date_to = dateTo;
        
        const feedbackResponse = await api.api.v1OrganizationsFeedbackList(selectedOrganization, query);
        const feedbackData = feedbackResponse.data;
        const organizationFeedback = feedbackData?.data || [];
        
        allFeedback = organizationFeedback.map((fb: any) => ({
          id: fb.id,
          customer_email: fb.customer_email,
          rating: fb.overall_rating,
          comment: fb.comment,
          product_name: fb.product?.name || null,
          organization_name: organizations.find(r => r.id === selectedOrganization)?.name,
          location_name: fb.location_name,
          qr_code: fb.qr_code,
          responses: fb.responses,
          created_at: fb.created_at
        }));
      } else {
        // Load feedback from all organizations with server-side filtering
        const query: any = {};
        
        if (searchQuery) query.search = searchQuery;
        if (selectedRating) {
          query.rating_min = parseInt(selectedRating);
          query.rating_max = parseInt(selectedRating);
        }
        if (selectedProduct) query.product_id = selectedProduct;
        if (dateFrom) query.date_from = dateFrom;
        if (dateTo) query.date_to = dateTo;
        
        const feedbackPromises = organizations.map(async (organization) => {
          try {
            const feedbackResponse = await api.api.v1OrganizationsFeedbackList(organization.id, query);
            const feedbackData = feedbackResponse.data;
            const organizationFeedback = feedbackData?.data || [];
            
            return organizationFeedback.map((fb: any) => ({
              id: fb.id,
              customer_email: fb.customer_email,
              rating: fb.overall_rating,
              comment: fb.comment,
              product_name: fb.product?.name || null,
              organization_name: organization.name,
              location_name: fb.location_name,
              qr_code: fb.qr_code,
              responses: fb.responses,
              created_at: fb.created_at
            }));
          } catch (err) {
            console.error(`Error loading feedback for organization ${organization.id}:`, err);
            return [];
          }
        });
        
        const feedbackArrays = await Promise.all(feedbackPromises);
        allFeedback = feedbackArrays.flat();
      }

      // Server-side filtering is now handled by the backend
      // Sort by creation date (newest first) - though backend already sorts
      allFeedback.sort((a, b) => new Date(b.created_at).getTime() - new Date(a.created_at).getTime());

      feedback = allFeedback;
      totalCount = allFeedback.length;

    } catch (err) {
      error = handleApiError(err);
    } finally {
      loading = false;
    }
  }

  // Remove manual filter change handler since filters are now automatic
  // function handleFilterChange() {
  //   loadFeedback();
  // }

  function clearFilters() {
    searchQuery = '';
    searchInput = ''; // Clear the input field too
    selectedOrganization = '';
    selectedProduct = '';
    selectedRating = '';
    
    // Reset to default 15-day filter
    setDefaultDateFilter();
    
    // Clear any pending search timeout
    if (searchTimeout) {
      clearTimeout(searchTimeout);
    }
    
    loadFeedback();
  }

  function formatDate(dateString: string): string {
    try {
      const date = new Date(dateString);
      return date.toLocaleDateString() + ' ' + date.toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' });
    } catch {
      return dateString;
    }
  }

  function getQuestionText(response: any, index: number): string {
    // First priority: Backend-provided question text
    if (response.question_text && response.question_text.trim()) {
      return response.question_text;
    }
    
    // Second priority: Legacy question field (for backward compatibility)
    if (response.question && response.question.trim()) {
      return response.question;
    }
    
    // Fallback: Generate smart question text based on answer type and common patterns
    const answer = response.answer;
    
    if (typeof answer === 'boolean') {
      return `Question ${index + 1} (Yes/No)`;
    } else if (typeof answer === 'number' && answer >= 1 && answer <= 5) {
      return `Rating Question ${index + 1}`;
    } else if (typeof answer === 'number') {
      return `Numeric Question ${index + 1}`;
    } else if (Array.isArray(answer)) {
      return `Multiple Choice Question ${index + 1}`;
    } else if (typeof answer === 'string' && answer.length > 50) {
      return `Comment Question ${index + 1}`;
    } else {
      return `Question ${index + 1}`;
    }
  }

  function renderStars(rating: number): string {
    return '★'.repeat(rating) + '☆'.repeat(5 - rating);
  }

  function getRatingColor(rating: number): string {
    if (rating >= 4) return 'text-green-600';
    if (rating >= 3) return 'text-yellow-600';
    return 'text-red-600';
  }

  function getRatingBgColor(rating: number): string {
    if (rating >= 4) return 'bg-green-100 text-green-800';
    if (rating >= 3) return 'bg-yellow-100 text-yellow-800';
    return 'bg-red-100 text-red-800';
  }

  function exportToCSV() {
    // Get all unique questions across all feedback, prefixed with product name for clarity
    const allQuestions = new Set<string>();
    feedback.forEach(fb => {
      if (fb.responses) {
        fb.responses.forEach((response: any, index: number) => {
          const questionText = getQuestionText(response, index);
          const productName = fb.product_name || 'General';
          const prefixedQuestion = `[${productName}] ${questionText}`;
          allQuestions.add(prefixedQuestion);
        });
      }
    });
    
    const questionArray = Array.from(allQuestions).sort();
    const headers = ['ID', 'Date', 'Customer Email', 'Rating', 'Product', 'Organization', 'Location', 'Comment', ...questionArray];
    
    const csvContent = [
      headers.join(','),
      ...feedback.map(fb => {
        // Basic fields
        const basicFields = [
          fb.id,
          fb.created_at,
          fb.customer_email || 'Anonymous',
          fb.rating || '',
          fb.product_name || '',
          fb.organization_name || '',
          typeof fb.qr_code === 'object' ? (fb.qr_code?.location || fb.qr_code?.name || fb.qr_code?.identifier || '') : (fb.qr_code || ''),
          `"${fb.comment?.replace(/"/g, '""') || ''}"`
        ];
        
        // Response fields - create a map of question to answer
        const responseMap = new Map<string, string>();
        if (fb.responses) {
          fb.responses.forEach((response: any, index: number) => {
            const questionText = getQuestionText(response, index);
            const productName = fb.product_name || 'General';
            const prefixedQuestion = `[${productName}] ${questionText}`;
            let answerText = '';
            
            if (typeof response.answer === 'boolean') {
              answerText = response.answer ? 'Yes' : 'No';
            } else if (typeof response.answer === 'number') {
              answerText = response.answer.toString();
            } else if (Array.isArray(response.answer)) {
              answerText = response.answer.join('; ');
            } else {
              answerText = String(response.answer || '');
            }
            
            responseMap.set(prefixedQuestion, answerText);
          });
        }
        
        // Add response values in the same order as headers
        const responseFields = questionArray.map(question => {
          const answer = responseMap.get(question) || '';
          return `"${answer.replace(/"/g, '""')}"`;
        });
        
        return [...basicFields, ...responseFields].join(',');
      })
    ].join('\n');

    const blob = new Blob([csvContent], { type: 'text/csv;charset=utf-8;' });
    const link = document.createElement('a');
    const url = URL.createObjectURL(blob);
    link.setAttribute('href', url);
    link.setAttribute('download', `feedback-export-${new Date().toISOString().split('T')[0]}.csv`);
    link.style.visibility = 'hidden';
    document.body.appendChild(link);
    link.click();
    document.body.removeChild(link);
  }
</script>

<svelte:head>
  <title>Feedback Management - Kyooar</title>
  <meta name="description" content="Manage and analyze customer feedback" />
</svelte:head>

<div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
    <!-- Page Header -->
    <div class="mb-8">
      <div class="flex flex-col lg:flex-row lg:justify-between lg:items-center gap-6">
        <div class="space-y-3">
          <div class="flex items-center space-x-3">
            <div class="h-12 w-12 bg-gradient-to-br from-blue-500 to-purple-600 rounded-2xl flex items-center justify-center shadow-lg shadow-blue-500/25">
              <svg class="h-6 w-6 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 12h.01M12 12h.01M16 12h.01M21 12c0 4.418-4.03 8-9 8a9.863 9.863 0 01-4.255-.949L3 20l1.395-3.72C3.512 15.042 3 13.574 3 12c0-4.418 4.03-8 9-8s9 3.582 9 8z" />
              </svg>
            </div>
            <div>
              <h1 class="text-3xl font-bold bg-gradient-to-r from-gray-900 to-gray-700 bg-clip-text text-transparent">
                Feedback Management
              </h1>
              <div class="flex items-center space-x-4 mt-1">
                <p class="text-gray-600 font-medium">Review and analyze customer feedback from all your organizations</p>
                {#if !loading && feedback.length > 0}
                  <div class="flex items-center space-x-3 text-sm">
                    <div class="flex items-center space-x-1">
                      <div class="w-2 h-2 bg-blue-400 rounded-full"></div>
                      <span class="text-gray-600">{totalCount} Total</span>
                    </div>
                    <div class="flex items-center space-x-1">
                      <div class="w-2 h-2 bg-purple-400 rounded-full"></div>
                      <span class="text-gray-600">{feedback.filter(f => f.rating >= 4).length} Positive</span>
                    </div>
                  </div>
                {/if}
              </div>
            </div>
          </div>
        </div>
        
        <div class="flex items-center space-x-3">
          <!-- Export CSV Button -->
          <RoleGate roles={['OWNER', 'ADMIN', 'MANAGER']}>
            <Button 
              variant="gradient" 
              size="lg" 
              class="group relative overflow-hidden shadow-lg hover:shadow-xl transition-all duration-300" 
              onclick={exportToCSV}
              disabled={feedback.length === 0}
            >
              <div class="absolute inset-0 bg-gradient-to-r from-blue-600 to-purple-600 opacity-0 group-hover:opacity-100 transition-opacity duration-300"></div>
              <svg class="h-5 w-5 mr-2 relative z-10 group-hover:scale-110 transition-transform duration-200" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 10v6m0 0l-3-3m3 3l3-3m2 8H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
              </svg>
              <span class="relative z-10">Export CSV</span>
            </Button>
          </RoleGate>
        </div>
      </div>
    </div>
    <!-- Filters -->
    <Card variant="default" hover interactive class="mb-6 group transform transition-all duration-300 animate-fade-in-up">
      <div class="space-y-4">
        <!-- First Row: Search and Organization -->
        <div class="grid grid-cols-1 lg:grid-cols-3 gap-4">
          <!-- Search -->
          <div class="space-y-1">
            <label class="text-xs font-medium text-gray-500 uppercase tracking-wider">Search</label>
            <Input
              type="text"
              placeholder="Search comments, emails, products..."
              bind:value={searchInput}
              oninput={handleSearchInput}
              class="w-full"
            />
          </div>

          <!-- Organization Filter -->
          <div class="space-y-1">
            <label class="text-xs font-medium text-gray-500 uppercase tracking-wider">Organization</label>
            <Select
              bind:value={selectedOrganization}
              options={[
                { value: '', label: 'All Organizations' },
                ...organizations.map(r => ({ value: r.id, label: r.name }))
              ]}
              minWidth="min-w-full"
            />
          </div>

          <!-- Product Filter -->
          <div class="space-y-1">
            <label class="text-xs font-medium text-gray-500 uppercase tracking-wider">Product</label>
            <Select
              bind:value={selectedProduct}
              options={[
                { value: '', label: 'All Products' },
                ...products.map(d => ({ value: d.id, label: d.name }))
              ]}
              minWidth="min-w-full"
              disabled={products.length === 0}
            />
          </div>
        </div>

        <!-- Second Row: Rating and Date Range -->
        <div class="grid grid-cols-1 lg:grid-cols-2 gap-4">
          <!-- Rating Filter -->
          <div class="space-y-1">
            <label class="text-xs font-medium text-gray-500 uppercase tracking-wider">Rating</label>
            <Select
              bind:value={selectedRating}
              options={[
                { value: '', label: 'All Ratings' },
                { value: '5', label: '5★ Excellent' },
                { value: '4', label: '4★ Good' },
                { value: '3', label: '3★ Average' },
                { value: '2', label: '2★ Poor' },
                { value: '1', label: '1★ Very Poor' }
              ]}
              minWidth="min-w-full"
            />
          </div>

          <!-- Date Range -->
          <div class="space-y-1">
            <label class="text-xs font-medium text-gray-500 uppercase tracking-wider">Date Range</label>
            <div class="flex items-center gap-3">
              <div class="flex-1">
                <Input
                  type="date"
                  bind:value={dateFrom}
                  placeholder="From"
                  class="w-full"
                />
              </div>
              <span class="text-sm text-gray-500 font-medium">to</span>
              <div class="flex-1">
                <Input
                  type="date"
                  bind:value={dateTo}
                  placeholder="To"
                  class="w-full"
                />
              </div>
            </div>
          </div>
        </div>

        <!-- Filter Summary and Actions -->
        <div class="flex flex-col sm:flex-row sm:items-center sm:justify-between gap-3 pt-4 border-t border-gray-100">
          <div class="flex items-center gap-4 text-sm">
            <div class="flex items-center gap-2">
              <svg class="h-4 w-4 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5H7a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2" />
              </svg>
              <span class="text-gray-600 font-medium">
                {feedback.length} {feedback.length === 1 ? 'entry' : 'entries'}
              </span>
            </div>
            {#if totalCount > feedback.length}
              <div class="h-4 w-px bg-gray-200"></div>
              <span class="text-gray-500">
                {totalCount} total
              </span>
            {/if}
          </div>
          <Button variant="outline" size="sm" onclick={clearFilters}>
            <svg class="h-4 w-4 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
            </svg>
            Clear Filters
          </Button>
        </div>
      </div>
    </Card>

    {#if loading}
      <!-- Loading State -->
      <div class="space-y-4">
        {#each Array(5) as _}
          <Card variant="default" class="opacity-50">
            <div class="animate-pulse">
              <div class="flex justify-between items-start mb-4">
                <div class="space-y-2">
                  <div class="h-4 bg-gray-200 rounded w-32"></div>
                  <div class="h-4 bg-gray-200 rounded w-48"></div>
                </div>
                <div class="h-6 bg-gray-200 rounded w-20"></div>
              </div>
              <div class="h-16 bg-gray-200 rounded"></div>
            </div>
          </Card>
        {/each}
      </div>

    {:else if error}
      <!-- Error State -->
      <Card variant="default" hover interactive class="group">
        <div class="text-center py-12">
          <svg class="h-12 w-12 text-red-500 mx-auto mb-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-.833-1.964-.833-2.732 0L3.732 16.5c-.77.833.192 3 1.732 3z" />
          </svg>
          <h3 class="text-lg font-medium text-gray-900 mb-2">Failed to load feedback</h3>
          <p class="text-gray-600 mb-4">{error}</p>
          <Button onclick={loadFeedback}>Try Again</Button>
        </div>
      </Card>

    {:else if feedback.length === 0}
      <!-- Empty State -->
      <Card variant="default" hover interactive class="group">
        <div class="text-center py-12">
          <svg class="h-12 w-12 text-gray-400 mx-auto mb-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 12h.01M12 12h.01M16 12h.01M21 12c0 4.418-4.03 8-9 8a9.863 9.863 0 01-4.255-.949L3 20l1.395-3.72C3.512 15.042 3 13.574 3 12c0-4.418 4.03-8 9-8s9 3.582 9 8z" />
          </svg>
          <h3 class="text-lg font-medium text-gray-900 mb-2">No feedback available</h3>
          <p class="text-gray-600">There is no feedback to display at this time.</p>
        </div>
      </Card>

    {:else}
      <!-- Feedback List -->
      <div class="space-y-6">
        {#each feedback as fb, index}
          <div class="group relative animate-fade-in-up" style="animation-delay: {index * 50}ms">
            <!-- Modern card with gradient border on hover -->
            <div class="absolute -inset-0.5 bg-gradient-to-r from-blue-500 to-purple-600 rounded-2xl opacity-0 group-hover:opacity-20 blur transition duration-500"></div>
            <Card variant="default" class="relative overflow-hidden border-0 shadow-xl hover:shadow-2xl transition-all duration-500">
              <!-- Header Section -->
              <div class="relative">
                <!-- Background gradient accent -->
                <div class="absolute top-0 right-0 w-32 h-32 bg-gradient-to-br from-blue-500/5 to-purple-600/5 rounded-full blur-3xl"></div>
                
                <div class="relative flex flex-col sm:flex-row sm:items-start sm:justify-between gap-4 mb-6">
                  <!-- Left side info -->
                  <div class="flex-1 space-y-3">
                    <!-- Rating and badges row -->
                    <div class="flex flex-wrap items-center gap-3">
                      <!-- Modern rating badge -->
                      <div class="relative">
                        <div class="absolute inset-0 bg-gradient-to-r {fb.rating >= 4 ? 'from-green-400 to-emerald-500' : fb.rating >= 3 ? 'from-yellow-400 to-orange-500' : 'from-red-400 to-pink-500'} rounded-xl blur opacity-25"></div>
                        <div class="relative flex items-center gap-2 px-4 py-2 bg-gradient-to-r {fb.rating >= 4 ? 'from-green-50 to-emerald-50 border-green-200' : fb.rating >= 3 ? 'from-yellow-50 to-orange-50 border-yellow-200' : 'from-red-50 to-pink-50 border-red-200'} border rounded-xl">
                          <div class="flex text-lg">
                            {#each Array(5) as _, i}
                              {#if i < fb.rating}
                                <svg class="w-5 h-5 {fb.rating >= 4 ? 'text-green-500' : fb.rating >= 3 ? 'text-yellow-500' : 'text-red-500'}" fill="currentColor" viewBox="0 0 20 20">
                                  <path d="M9.049 2.927c.3-.921 1.603-.921 1.902 0l1.07 3.292a1 1 0 00.95.69h3.462c.969 0 1.371 1.24.588 1.81l-2.8 2.034a1 1 0 00-.364 1.118l1.07 3.292c.3.921-.755 1.688-1.54 1.118l-2.8-2.034a1 1 0 00-1.175 0l-2.8 2.034c-.784.57-1.838-.197-1.539-1.118l1.07-3.292a1 1 0 00-.364-1.118L2.98 8.72c-.783-.57-.38-1.81.588-1.81h3.461a1 1 0 00.951-.69l1.07-3.292z" />
                                </svg>
                              {:else}
                                <svg class="w-5 h-5 text-gray-300" fill="none" stroke="currentColor" stroke-width="1.5" viewBox="0 0 20 20">
                                  <path d="M9.049 2.927c.3-.921 1.603-.921 1.902 0l1.07 3.292a1 1 0 00.95.69h3.462c.969 0 1.371 1.24.588 1.81l-2.8 2.034a1 1 0 00-.364 1.118l1.07 3.292c.3.921-.755 1.688-1.54 1.118l-2.8-2.034a1 1 0 00-1.175 0l-2.8 2.034c-.784.57-1.838-.197-1.539-1.118l1.07-3.292a1 1 0 00-.364-1.118L2.98 8.72c-.783-.57-.38-1.81.588-1.81h3.461a1 1 0 00.951-.69l1.07-3.292z" />
                                </svg>
                              {/if}
                            {/each}
                          </div>
                          <span class="font-bold {fb.rating >= 4 ? 'text-green-700' : fb.rating >= 3 ? 'text-yellow-700' : 'text-red-700'}">{fb.rating}.0</span>
                        </div>
                      </div>
                      
                      {#if fb.product_name}
                        <div class="flex items-center gap-2 px-3 py-1.5 bg-purple-50 border border-purple-200 rounded-lg">
                          <svg class="w-4 h-4 text-purple-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 6.253v13m0-13C10.832 5.477 9.246 5 7.5 5S4.168 5.477 3 6.253v13C4.168 18.477 5.754 18 7.5 18s3.332.477 4.5 1.253m0-13C13.168 5.477 14.754 5 16.5 5c1.747 0 3.332.477 4.5 1.253v13C19.832 18.477 18.247 18 16.5 18c-1.746 0-3.332.477-4.5 1.253" />
                          </svg>
                          <span class="text-sm font-medium text-purple-700">{fb.product_name}</span>
                        </div>
                      {/if}
                      
                      {#if fb.organization_name}
                        <div class="flex items-center gap-2 px-3 py-1.5 bg-blue-50 border border-blue-200 rounded-lg">
                          <svg class="w-4 h-4 text-blue-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 21V5a2 2 0 00-2-2H7a2 2 0 00-2 2v16m14 0h2m-2 0h-5m-9 0H3m2 0h5M9 7h1m-1 4h1m4-4h1m-1 4h1m-5 10v-5a1 1 0 011-1h2a1 1 0 011 1v5m-4 0h4" />
                          </svg>
                          <span class="text-sm font-medium text-blue-700">{fb.organization_name}</span>
                        </div>
                      {/if}
                    </div>
                    
                    <!-- Meta information with modern icons -->
                    <div class="flex flex-wrap items-center gap-4 text-sm text-gray-600">
                      <div class="flex items-center gap-1.5">
                        <svg class="w-4 h-4 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z" />
                        </svg>
                        <span>{formatDate(fb.created_at)}</span>
                      </div>
                      
                      <div class="flex items-center gap-1.5">
                        <svg class="w-4 h-4 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M16 7a4 4 0 11-8 0 4 4 0 018 0zM12 14a7 7 0 00-7 7h14a7 7 0 00-7-7z" />
                        </svg>
                        <span>{fb.customer_email || 'Anonymous Guest'}</span>
                      </div>
                      
                      {#if fb.qr_code}
                        <div class="flex items-center gap-1.5">
                          <svg class="w-4 h-4 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v1m6 11h2m-6 0h-2v4m0-11v3m0 0h.01M12 12h4.01M16 20h4M4 12h4m12 0h.01M5 8h2a1 1 0 001-1V5a1 1 0 00-1-1H5a1 1 0 00-1 1v2a1 1 0 001 1zm12 0h2a1 1 0 001-1V5a1 1 0 00-1-1h-2a1 1 0 00-1 1v2a1 1 0 001 1zM5 20h2a1 1 0 001-1v-2a1 1 0 00-1-1H5a1 1 0 00-1 1v2a1 1 0 001 1z" />
                          </svg>
                          <span>{typeof fb.qr_code === 'object' ? (fb.qr_code?.location || fb.qr_code?.name || fb.qr_code?.identifier || 'QR Code') : fb.qr_code}</span>
                        </div>
                      {/if}
                    </div>
                  </div>
                  
                  <!-- Right side actions -->
                  <div class="flex items-center gap-2">
                    <button class="p-2 text-gray-400 hover:text-gray-600 hover:bg-gray-100 rounded-lg transition-colors">
                      <svg class="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 12h.01M12 12h.01M19 12h.01M6 12a1 1 0 11-2 0 1 1 0 012 0zm7 0a1 1 0 11-2 0 1 1 0 012 0zm7 0a1 1 0 11-2 0 1 1 0 012 0z" />
                      </svg>
                    </button>
                  </div>
                </div>
                
                <!-- Comment Section -->
                {#if fb.comment}
                  <div class="relative mb-6">
                    <div class="absolute -left-1 top-0 bottom-0 w-1 bg-gradient-to-b from-blue-500 to-purple-600 rounded-full"></div>
                    <div class="pl-6">
                      <div class="flex items-start gap-3">
                        <svg class="w-5 h-5 text-gray-400 mt-0.5 flex-shrink-0" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 12h.01M12 12h.01M16 12h.01M21 12c0 4.418-4.03 8-9 8a9.863 9.863 0 01-4.255-.949L3 20l1.395-3.72C3.512 15.042 3 13.574 3 12c0-4.418 4.03-8 9-8s9 3.582 9 8z" />
                        </svg>
                        <blockquote class="flex-1 text-gray-700 leading-relaxed">
                          <p class="text-base italic">"{fb.comment}"</p>
                        </blockquote>
                      </div>
                    </div>
                  </div>
                {/if}
                
                <!-- Questions & Responses Section -->
                {#if fb.responses && fb.responses.length > 0}
                  <div class="space-y-4">
                    <!-- Collapsible Header -->
                    <button
                      onclick={() => toggleCollapse(fb.id)}
                      class="w-full flex items-center justify-between p-3 -mx-3 rounded-lg border border-transparent hover:border-gray-200 hover:bg-gray-50/50 transition-all duration-200 cursor-pointer group/header focus:outline-none focus:ring-2 focus:ring-purple-500/20"
                    >
                      <div class="flex items-center gap-3">
                        <div class="p-2 bg-gradient-to-br from-purple-500 to-blue-600 rounded-lg group-hover/header:shadow-lg transition-shadow duration-200">
                          <svg class="w-5 h-5 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8.228 9c.549-1.165 2.03-2 3.772-2 2.21 0 4 1.343 4 3 0 1.4-1.278 2.575-3.006 2.907-.542.104-.994.54-.994 1.093m0 3h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
                          </svg>
                        </div>
                        <div class="text-left">
                          <h4 class="text-lg font-semibold text-gray-900 group-hover/header:text-purple-700 transition-colors duration-200">Customer Responses</h4>
                          <p class="text-sm text-gray-500 group-hover/header:text-gray-600 transition-colors duration-200">{fb.responses.length} question{fb.responses.length > 1 ? 's' : ''} answered</p>
                        </div>
                      </div>
                      
                      <!-- Toggle Indicator -->
                      <div class="flex items-center gap-2 px-3 py-2 text-sm font-medium text-purple-600 group-hover/header:text-purple-700 group-hover/header:bg-purple-50 rounded-lg transition-all duration-200 border border-purple-200/50 group-hover/header:border-purple-300">
                        <span class="hidden sm:block">{isCollapsed(fb.id) ? 'Show' : 'Hide'} Details</span>
                        <svg class="w-5 h-5 transform transition-transform duration-200 {isCollapsed(fb.id) ? 'rotate-0' : 'rotate-180'}" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 9l-7 7-7-7" />
                        </svg>
                      </div>
                    </button>
                    
                    <!-- Collapsed Preview -->
                    {#if isCollapsed(fb.id)}
                      <div class="ml-8 text-sm text-gray-400 italic">
                        Click above to view {fb.responses.length} detailed response{fb.responses.length > 1 ? 's' : ''}
                      </div>
                    {/if}
                    
                    <!-- Collapsible Content -->
                    {#if !isCollapsed(fb.id)}
                      <div class="space-y-4" style="animation: slideDown 0.3s ease-out;">
                        {#each fb.responses as response, index}
                          <div class="group/question relative">
                            <!-- Question Number Badge -->
                            <div class="absolute -left-3 top-0 z-10">
                              <div class="w-6 h-6 bg-gradient-to-r from-purple-500 to-blue-600 rounded-full flex items-center justify-center">
                                <span class="text-xs font-bold text-white">{index + 1}</span>
                              </div>
                            </div>
                            
                            <!-- Question Card -->
                            <div class="ml-6 bg-gradient-to-r from-gray-50 to-white border border-gray-200 rounded-xl p-5 group-hover/question:shadow-md transition-all duration-300">
                              <!-- Question Text -->
                              <div class="mb-3">
                                <p class="text-sm font-medium text-gray-900 leading-relaxed">
                                  {getQuestionText(response, index)}
                                </p>
                              </div>
                              
                              <!-- Answer Section -->
                              <div class="space-y-2">
                                <div class="flex items-start gap-3">
                                  <!-- Answer Type Icon -->
                                  <div class="mt-0.5">
                                    {#if typeof response.answer === 'boolean'}
                                      <svg class="w-4 h-4 {response.answer ? 'text-green-500' : 'text-red-500'}" fill="currentColor" viewBox="0 0 20 20">
                                        {#if response.answer}
                                          <path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zm3.707-9.293a1 1 0 00-1.414-1.414L9 10.586 7.707 9.293a1 1 0 00-1.414 1.414l2 2a1 1 0 001.414 0l4-4z" clip-rule="evenodd" />
                                        {:else}
                                          <path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zM8.707 7.293a1 1 0 00-1.414 1.414L8.586 10l-1.293 1.293a1 1 0 101.414 1.414L10 11.414l1.293 1.293a1 1 0 001.414-1.414L11.414 10l1.293-1.293a1 1 0 00-1.414-1.414L10 8.586 8.707 7.293z" clip-rule="evenodd" />
                                        {/if}
                                      </svg>
                                    {:else if typeof response.answer === 'number'}
                                      <svg class="w-4 h-4 text-blue-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M7 20l4-16m2 16l4-16M6 9h14M4 15h14" />
                                      </svg>
                                    {:else}
                                      <svg class="w-4 h-4 text-purple-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                        <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 12h.01M12 12h.01M16 12h.01M21 12c0 4.418-4.03 8-9 8a9.863 9.863 0 01-4.255-.949L3 20l1.395-3.72C3.512 15.042 3 13.574 3 12c0-4.418 4.03-8 9-8s9 3.582 9 8z" />
                                      </svg>
                                    {/if}
                                  </div>
                                  
                                  <!-- Answer Display -->
                                  <div class="flex-1">
                                    {#if typeof response.answer === 'boolean'}
                                      <div class="inline-flex items-center gap-2 px-3 py-1.5 rounded-lg {response.answer ? 'bg-green-100 text-green-800 border border-green-200' : 'bg-red-100 text-red-800 border border-red-200'}">
                                        <span class="font-semibold">{response.answer ? 'Yes' : 'No'}</span>
                                      </div>
                                    {:else if typeof response.answer === 'number'}
                                      {#if response.answer >= 1 && response.answer <= 5}
                                        <!-- Rating display -->
                                        <div class="flex items-center gap-2">
                                          <div class="flex">
                                            {#each Array(5) as _, i}
                                              <svg class="w-4 h-4 {i < response.answer ? 'text-yellow-400' : 'text-gray-300'}" fill="currentColor" viewBox="0 0 20 20">
                                                <path d="M9.049 2.927c.3-.921 1.603-.921 1.902 0l1.07 3.292a1 1 0 00.95.69h3.462c.969 0 1.371 1.24.588 1.81l-2.8 2.034a1 1 0 00-.364 1.118l1.07 3.292c.3.921-.755 1.688-1.54 1.118l-2.8-2.034a1 1 0 00-1.175 0l-2.8 2.034c-.784.57-1.838-.197-1.539-1.118l1.07-3.292a1 1 0 00-.364-1.118L2.98 8.72c-.783-.57-.38-1.81.588-1.81h3.461a1 1 0 00.951-.69l1.07-3.292z" />
                                              </svg>
                                            {/each}
                                          </div>
                                          <span class="font-semibold text-gray-900">{response.answer}/5</span>
                                        </div>
                                      {:else}
                                        <!-- Numeric value -->
                                        <div class="inline-flex items-center gap-2 px-3 py-1.5 bg-blue-100 text-blue-800 border border-blue-200 rounded-lg">
                                          <span class="font-semibold">{response.answer}</span>
                                        </div>
                                      {/if}
                                    {:else if Array.isArray(response.answer)}
                                      <!-- Multiple choice -->
                                      <div class="flex flex-wrap gap-2">
                                        {#each response.answer as item}
                                          <span class="inline-flex items-center px-3 py-1.5 bg-purple-100 text-purple-800 border border-purple-200 rounded-lg text-sm font-medium">
                                            {item}
                                          </span>
                                        {/each}
                                      </div>
                                    {:else}
                                      <!-- Text response -->
                                      <div class="bg-gray-50 border border-gray-200 rounded-lg p-3">
                                        <p class="text-gray-900 text-sm leading-relaxed">{response.answer}</p>
                                      </div>
                                    {/if}
                                  </div>
                                </div>
                              </div>
                            </div>
                          </div>
                        {/each}
                      </div>
                    {/if}
                  </div>
                {/if}
              </div>
            </Card>
          </div>
        {/each}
      </div>
    {/if}
</div>

<style>
  @keyframes fade-in-up {
    from {
      opacity: 0;
      transform: translateY(10px);
    }
    to {
      opacity: 1;
      transform: translateY(0);
    }
  }

  .animate-fade-in-up {
    animation: fade-in-up 0.6s ease-out forwards;
    opacity: 0;
  }

  @keyframes slideDown {
    from {
      opacity: 0;
      max-height: 0;
      transform: translateY(-10px);
    }
    to {
      opacity: 1;
      max-height: 1000px;
      transform: translateY(0);
    }
  }
</style>
