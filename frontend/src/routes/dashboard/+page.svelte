<script lang="ts">
  import { onMount } from 'svelte';
  import { Card, Button, UserMenu, Logo } from '$lib/components/ui';
  import { getApiClient, handleApiError } from '$lib/api/client';
  import { auth } from '$lib/stores/auth';
  import { goto } from '$app/navigation';
  import { CheckCircle, Clock, Smartphone, QrCode } from 'lucide-svelte';

  interface DashboardStats {
    totalFeedback: number;
    averageRating: number;
    feedbackToday: number;
    activeQRCodes: number;
    topRatedProduct?: {
      name: string;
      rating: number;
    };
    recentFeedbackCount: number;
  }

  interface DashboardMetrics {
    todays_feedback: number;
    trend_vs_yesterday: number;
    overall_satisfaction: number;
    satisfaction_trend: string;
    recommendation_rate: number;
    positive_sentiment: number;
    top_issues: Array<{
      title: string;
      count: number;
      severity: string;
      action_link: string;
    }>;
    best_performers: Array<{
      product_id: string;
      product_name: string;
      score: number;
      feedback_count: number;
      trend: string;
    }>;
    needing_attention: Array<{
      product_id: string;
      product_name: string;
      score: number;
      feedback_count: number;
      trend: string;
    }>;
    recent_feedback: Array<{
      feedback_id: string;
      product_name: string;
      customer_name: string;
      score: number;
      sentiment: string;
      highlight: string;
      created_at: string;
    }>;
  }

  interface RecentFeedback {
    id: string;
    customer_email?: string;
    rating: number;
    comment?: string;
    product_name?: string;
    organization_name?: string;
    created_at: string;
  }

  let loading = $state(true);
  let error = $state('');
  let stats = $state<DashboardStats>({
    totalFeedback: 0,
    averageRating: 0,
    feedbackToday: 0,
    activeQRCodes: 0,
    recentFeedbackCount: 0
  });
  let recentFeedback = $state<RecentFeedback[]>([]);
  let dashboardMetrics = $state<DashboardMetrics | null>(null);
  let hasInitialized = $state(false);

  let authState = $derived($auth);

  $effect(() => {
    // Check if user is authenticated
    if (authState.isAuthenticated && !hasInitialized) {
      hasInitialized = true;
      loadDashboardData();
    } else if (!authState.isAuthenticated) {
      goto('/login');
    }
  });

  async function loadDashboardData() {
    loading = true;
    error = '';

    try {
      const api = getApiClient();
      
      // Get all organizations for the account
      const organizationsResponse = await api.api.v1OrganizationsList();
      
      if (organizationsResponse.data.success && organizationsResponse.data.data) {
        const organizations = organizationsResponse.data.data;
        
        // For now, using the first organization for demo purposes
        // In a real app, you'd either aggregate across all organizations or let user select
        if (organizations.length > 0) {
          const firstOrganization = organizations[0];
          
          // Get QR codes and analytics for the first organization
          const [qrCodesResponse, analyticsResponse] = await Promise.all([
            api.api.v1OrganizationsQrCodesList(firstOrganization.id!),
            api.api.v1AnalyticsOrganizationsDetail(firstOrganization.id!)
          ]);
          
          // Try to get new dashboard metrics (might not exist yet)
          try {
            const dashboardResponse = await api.api.v1AnalyticsDashboardDetail(firstOrganization.id!);
            if (dashboardResponse.data.success && dashboardResponse.data.data) {
              dashboardMetrics = dashboardResponse.data.data;
            }
          } catch (err) {
            // Dashboard metrics not available yet, using legacy analytics
          }
          
          // Calculate stats
          const activeQRCodes = qrCodesResponse.data.success && qrCodesResponse.data.data 
            ? qrCodesResponse.data.data.filter(qr => qr.is_active).length 
            : 0;
          
          // Parse analytics data (legacy)
          const analyticsData = analyticsResponse.data?.data || {};
          const totalFeedback = analyticsData?.total_feedback || 0;
          const averageRating = analyticsData?.average_rating || 0;
          const feedbackToday = analyticsData?.feedback_today || 0;
          const topProduct = analyticsData?.top_rated_products?.[0];
          const recentFeedbackData = analyticsData?.recent_feedback || [];
          
          // Both endpoints now return 1-5 scale for ratings
          let displayRating = dashboardMetrics?.overall_satisfaction || averageRating;
          
          // Clamp rating to valid 1-5 range and log if out of bounds
          if (displayRating > 5) {
            console.warn('Rating out of bounds (too high):', displayRating);
            displayRating = 5;
          } else if (displayRating < 0) {
            console.warn('Rating out of bounds (too low):', displayRating);
            displayRating = 0;
          }

          stats = {
            totalFeedback,
            averageRating: displayRating,
            feedbackToday: dashboardMetrics?.todays_feedback || feedbackToday,
            activeQRCodes,
            topRatedProduct: topProduct ? {
              name: topProduct.product_name,
              rating: topProduct.average_rating
            } : undefined,
            recentFeedbackCount: dashboardMetrics?.recent_feedback?.length || recentFeedbackData.length
          };
          
          // Map recent feedback - prefer new dashboard metrics
          const feedbackSource = dashboardMetrics?.recent_feedback || recentFeedbackData;
          recentFeedback = feedbackSource.slice(0, 5).map((fb: any) => ({
            id: fb.feedback_id || fb.id,
            customer_email: fb.customer_name || fb.customer_email,
            rating: fb.score || fb.rating,
            comment: fb.highlight || fb.comment,
            product_name: fb.product_name,
            organization_name: firstOrganization.name,
            created_at: fb.created_at
          }));
        } else {
          // No organizations yet
          stats = {
            totalFeedback: 0,
            averageRating: 0,
            feedbackToday: 0,
            activeQRCodes: 0,
            recentFeedbackCount: 0
          };
          recentFeedback = [];
        }
      } else {
        // Fallback to zero stats
        stats = {
          totalFeedback: 0,
          averageRating: 0,
          feedbackToday: 0,
          activeQRCodes: 0,
          recentFeedbackCount: 0
        };
        recentFeedback = [];
      }

    } catch (err) {
      error = handleApiError(err);
    } finally {
      loading = false;
    }
  }

  function handleLogout() {
    auth.logout();
    goto('/login');
  }

  function formatDate(dateString: string): string {
    try {
      const date = new Date(dateString);
      return date.toLocaleDateString() + ' ' + date.toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' });
    } catch {
      return dateString;
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
</script>

<svelte:head>
  <title>Dashboard - Kyooar</title>
  <meta name="description" content="Kyooar organization management dashboard" />
</svelte:head>

<div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-8">
    {#if loading}
      <!-- Loading State -->
      <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6 mb-8">
        {#each Array(4) as _}
          <Card>
            <div class="animate-pulse">
              <div class="h-4 bg-gray-200 rounded w-3/4 mb-2"></div>
              <div class="h-8 bg-gray-200 rounded w-1/2"></div>
            </div>
          </Card>
        {/each}
      </div>

    {:else if error}
      <!-- Error State -->
      <Card>
        <div class="text-center py-12">
          <svg class="h-12 w-12 text-red-500 mx-auto mb-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-.833-1.964-.833-2.732 0L3.732 16.5c-.77.833.192 3 1.732 3z" />
          </svg>
          <h3 class="text-lg font-medium text-gray-900 mb-2">Failed to load dashboard</h3>
          <p class="text-gray-600 mb-4">{error}</p>
          <Button on:click={loadDashboardData}>Try Again</Button>
        </div>
      </Card>

    {:else}
      <!-- Stats Cards -->
      <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6 mb-12">
        <!-- Total Feedback -->
        <Card variant="gradient" hover interactive>
          <div class="flex items-center justify-between">
            <div class="space-y-2">
              <p class="text-sm font-semibold text-gray-600 uppercase tracking-wide">Total Feedback</p>
              <p class="text-3xl font-bold bg-gradient-to-r from-blue-600 to-purple-600 bg-clip-text text-transparent">
                {stats.totalFeedback}
              </p>
              <div class="flex items-center space-x-1 text-green-600">
                <svg class="h-4 w-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 7h8m0 0v8m0-8l-8 8-4-4-6 6" />
                </svg>
                <span class="text-sm font-medium">All time</span>
              </div>
            </div>
            <div class="h-16 w-16 bg-gradient-to-br from-blue-500 to-purple-600 rounded-2xl flex items-center justify-center shadow-lg shadow-blue-500/25">
              <svg class="h-8 w-8 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 12h.01M12 12h.01M16 12h.01M21 12c0 4.418-4.03 8-9 8a9.863 9.863 0 01-4.255-.949L3 20l1.395-3.72C3.512 15.042 3 13.574 3 12c0-4.418 4.03-8 9-8s9 3.582 9 8z" />
              </svg>
            </div>
          </div>
        </Card>

        <!-- Average Rating -->
        <Card variant="gradient" hover interactive>
          <div class="flex items-center justify-between">
            <div class="space-y-2">
              <p class="text-sm font-semibold text-gray-600 uppercase tracking-wide">Average Rating</p>
              <p class="text-3xl font-bold bg-gradient-to-r from-yellow-600 to-orange-600 bg-clip-text text-transparent">
                {stats.averageRating.toFixed(1)}
              </p>
              <div class="flex items-center space-x-1">
                <div class="flex text-yellow-400">
                  {#each Array(5) as _, i}
                    <svg class="h-4 w-4 {i < Math.round(stats.averageRating) ? 'text-yellow-400' : 'text-gray-300'}" fill="currentColor" viewBox="0 0 24 24">
                      <path d="M11.049 2.927c.3-.921 1.603-.921 1.902 0l1.519 4.674a1 1 0 00.95.69h4.915c.969 0 1.371 1.24.588 1.81l-3.976 2.888a1 1 0 00-.363 1.118l1.518 4.674c.3.922-.755 1.688-1.538 1.118l-3.976-2.888a1 1 0 00-1.176 0l-3.976 2.888c-.783.57-1.838-.197-1.538-1.118l1.518-4.674a1 1 0 00-.363-1.118l-3.976-2.888c-.784-.57-.38-1.81.588-1.81h4.914a1 1 0 00.951-.69l1.519-4.674z" />
                    </svg>
                  {/each}
                </div>
              </div>
            </div>
            <div class="h-16 w-16 bg-gradient-to-br from-yellow-500 to-orange-500 rounded-2xl flex items-center justify-center shadow-lg shadow-yellow-500/25">
              <svg class="h-8 w-8 text-white" fill="currentColor" viewBox="0 0 24 24">
                <path d="M11.049 2.927c.3-.921 1.603-.921 1.902 0l1.519 4.674a1 1 0 00.95.69h4.915c.969 0 1.371 1.24.588 1.81l-3.976 2.888a1 1 0 00-.363 1.118l1.518 4.674c.3.922-.755 1.688-1.538 1.118l-3.976-2.888a1 1 0 00-1.176 0l-3.976 2.888c-.783.57-1.838-.197-1.538-1.118l1.518-4.674a1 1 0 00-.363-1.118l-3.976-2.888c-.784-.57-.38-1.81.588-1.81h4.914a1 1 0 00.951-.69l1.519-4.674z" />
              </svg>
            </div>
          </div>
        </Card>

        <!-- Today's Feedback -->
        <Card variant="gradient" hover interactive>
          <div class="flex items-center justify-between">
            <div class="space-y-2">
              <p class="text-sm font-semibold text-gray-600 uppercase tracking-wide">Today's Feedback</p>
              <p class="text-3xl font-bold bg-gradient-to-r from-green-600 to-emerald-600 bg-clip-text text-transparent">
                {stats.feedbackToday}
              </p>
              <div class="flex items-center space-x-1 text-green-600">
                <svg class="h-4 w-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 3v1m0 16v1m9-9h-1M4 12H3m15.364 6.364l-.707-.707M6.343 6.343l-.707-.707m12.728 0l-.707.707M6.343 17.657l-.707.707" />
                </svg>
                <span class="text-sm font-medium">Today</span>
              </div>
            </div>
            <div class="h-16 w-16 bg-gradient-to-br from-green-500 to-emerald-500 rounded-2xl flex items-center justify-center shadow-lg shadow-green-500/25">
              <svg class="h-8 w-8 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 7h8m0 0v8m0-8l-8 8-4-4-6 6" />
              </svg>
            </div>
          </div>
        </Card>

        <!-- Active QR Codes -->
        <Card variant="gradient" hover interactive>
          <div class="flex items-center justify-between">
            <div class="space-y-2">
              <p class="text-sm font-semibold text-gray-600 uppercase tracking-wide">Active QR Codes</p>
              <p class="text-3xl font-bold bg-gradient-to-r from-purple-600 to-indigo-600 bg-clip-text text-transparent">
                {stats.activeQRCodes}
              </p>
              <div class="flex items-center space-x-1 text-purple-600">
                <svg class="h-4 w-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
                </svg>
                <span class="text-sm font-medium">Active</span>
              </div>
            </div>
            <div class="h-16 w-16 bg-gradient-to-br from-purple-500 to-indigo-500 rounded-2xl flex items-center justify-center shadow-lg shadow-purple-500/25">
              <svg class="h-8 w-8 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 4v1m6 11h2m-6 0h-2v4m0-11v3m0 0h.01M12 12h4.01M16 20h4M4 12h4m12 0h.01M5 8h2a1 1 0 001-1V5a1 1 0 00-1-1H5a1 1 0 00-1 1v2a1 1 0 001 1zm12 0h2a1 1 0 001-1V5a1 1 0 00-1-1h-2a1 1 0 00-1 1v2a1 1 0 001 1zM5 20h2a1 1 0 001-1v-2a1 1 0 00-1-1H5a1 1 0 00-1 1v2a1 1 0 001 1z" />
              </svg>
            </div>
          </div>
        </Card>
      </div>

      <!-- Enhanced Analytics Section -->
      {#if dashboardMetrics}
        <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6 mb-8">
          <!-- Recommendation Rate -->
          <Card variant="gradient" hover interactive>
            <div class="space-y-3">
              <div class="flex items-center justify-between">
                <p class="text-sm font-semibold text-gray-600 uppercase tracking-wide">Completion Rate</p>
                <div class="h-8 w-8 bg-gradient-to-br from-green-500 to-emerald-600 rounded-lg flex items-center justify-center">
                  <CheckCircle class="h-4 w-4 text-white" />
                </div>
              </div>
              <p class="text-2xl font-bold text-green-600">{(dashboardMetrics.completion_rate ?? 0).toFixed(1)}%</p>
              <p class="text-sm text-gray-500">QR scans that become responses</p>
            </div>
          </Card>

          <!-- Positive Sentiment -->
          <Card variant="gradient" hover interactive>
            <div class="space-y-3">
              <div class="flex items-center justify-between">
                <p class="text-sm font-semibold text-gray-600 uppercase tracking-wide">Avg Response Time</p>
                <div class="h-8 w-8 bg-gradient-to-br from-blue-500 to-purple-600 rounded-lg flex items-center justify-center">
                  <Clock class="h-4 w-4 text-white" />
                </div>
              </div>
              <p class="text-2xl font-bold text-blue-600">{Math.round(dashboardMetrics.average_response_time ?? 0)} min</p>
              <p class="text-sm text-gray-500">From scan to submission</p>
            </div>
          </Card>

          <!-- Trend Indicator -->
          <Card variant="gradient" hover interactive>
            <div class="space-y-3">
              <div class="flex items-center justify-between">
                <p class="text-sm font-semibold text-gray-600 uppercase tracking-wide">Satisfaction Trend</p>
                <div class="h-8 w-8 bg-gradient-to-br from-yellow-500 to-orange-600 rounded-lg flex items-center justify-center">
                  <svg class="h-4 w-4 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    {#if dashboardMetrics.satisfaction_trend === 'up'}
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 7h8m0 0v8m0-8l-8 8-4-4-6 6" />
                    {:else if dashboardMetrics.satisfaction_trend === 'down'}
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13 17h8m0 0V9m0 8l-8-8-4 4-6-6" />
                    {:else}
                      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M5 12h14" />
                    {/if}
                  </svg>
                </div>
              </div>
              <p class="text-2xl font-bold {dashboardMetrics.satisfaction_trend === 'up' ? 'text-green-600' : dashboardMetrics.satisfaction_trend === 'down' ? 'text-red-600' : 'text-gray-600'}">
                {dashboardMetrics.satisfaction_trend === 'up' ? 'Improving' : dashboardMetrics.satisfaction_trend === 'down' ? 'Declining' : 'Stable'}
              </p>
              <p class="text-sm text-gray-500">vs previous period</p>
            </div>
          </Card>
        </div>

        <!-- Operational Insights -->
        {#if dashboardMetrics.device_breakdown && Object.keys(dashboardMetrics.device_breakdown).length > 0}
          <div class="grid grid-cols-1 lg:grid-cols-2 gap-6 mb-8">
            <!-- Device Analytics -->
            <Card variant="elevated" padding={false}>
              <div class="p-6">
                <h3 class="text-lg font-bold text-gray-900 mb-4">Device Analytics</h3>
                <div class="space-y-3">
                  {#each Object.entries(dashboardMetrics.device_breakdown || {}) as [device, count]}
                    <div class="flex items-center justify-between p-3 bg-blue-50 rounded-lg border border-blue-200">
                      <div class="flex items-center space-x-3">
                        <div class="h-8 w-8 bg-blue-500 rounded-full flex items-center justify-center">
                          <Smartphone class="h-4 w-4 text-white" />
                        </div>
                        <div>
                          <p class="font-semibold text-blue-900">{device}</p>
                          <p class="text-sm text-blue-600">{count} responses</p>
                        </div>
                      </div>
                      <div class="text-sm font-medium text-blue-600">
                        {((count / (dashboardMetrics.total_feedbacks || 1)) * 100).toFixed(1)}%
                      </div>
                    </div>
                  {/each}
                </div>
              </div>
            </Card>

            <!-- Peak Hours -->
            {#if dashboardMetrics.peak_hours && dashboardMetrics.peak_hours.length > 0}
              <Card variant="elevated" padding={false}>
                <div class="p-6">
                  <h3 class="text-lg font-bold text-gray-900 mb-4">Peak Activity Hours</h3>
                  <div class="space-y-3">
                    {#each dashboardMetrics.peak_hours as hour}
                      <div class="flex items-center justify-between p-3 bg-green-50 rounded-lg border border-green-200">
                        <div class="flex items-center space-x-3">
                          <div class="h-8 w-8 bg-green-500 rounded-full flex items-center justify-center">
                            <Clock class="h-4 w-4 text-white" />
                          </div>
                          <div>
                            <p class="font-semibold text-green-900">
                              {hour === 0 ? '12 AM' : hour <= 12 ? `${hour} ${hour === 12 ? 'PM' : 'AM'}` : `${hour - 12} PM`}
                            </p>
                            <p class="text-sm text-green-600">High activity period</p>
                          </div>
                        </div>
                      </div>
                    {/each}
                  </div>
                </div>
              </Card>
            {/if}
          </div>
        {/if}
      {/if}

      <!-- QR Code Performance Analytics -->
      {#if dashboardMetrics?.qr_performance}
        <div class="mb-8">
          <Card variant="elevated" padding={false}>
            <div class="p-6 lg:p-8">
              <div class="mb-6">
                <h3 class="text-xl font-bold bg-gradient-to-r from-gray-900 to-gray-700 bg-clip-text text-transparent">
                  QR Code Performance
                </h3>
                <p class="text-gray-600 font-medium">Usage and conversion analytics</p>
              </div>

              {#if dashboardMetrics.qr_performance.length > 0}
                <div class="space-y-4">
                  {#each dashboardMetrics.qr_performance as qr}
                    <div class="relative group">
                      <div class="absolute inset-0 bg-gradient-to-r from-blue-500/10 to-indigo-500/10 rounded-xl opacity-0 group-hover:opacity-100 transition-opacity duration-300"></div>
                      <div class="relative bg-white/50 backdrop-blur-sm rounded-xl border border-gray-100 p-5 hover:shadow-lg hover:shadow-gray-900/5 transition-all duration-300 group-hover:border-gray-200">
                        <div class="flex items-start justify-between mb-3">
                          <div class="flex items-center space-x-3">
                            <div class="h-10 w-10 bg-gradient-to-br from-blue-500 to-indigo-600 rounded-xl flex items-center justify-center shadow-lg">
                              <QrCode class="h-5 w-5 text-white" />
                            </div>
                            <div>
                              <p class="font-semibold text-gray-900">{qr.label || 'QR Code'}</p>
                              <p class="text-sm text-gray-500">{qr.organization_name}</p>
                              <p class="text-xs text-gray-400">{qr.location || 'No location set'}</p>
                            </div>
                          </div>
                          <div class="text-right">
                            <div class="flex flex-col items-end space-y-1">
                              <div class="flex items-center space-x-2">
                                <span class="text-sm text-gray-500">Scans:</span>
                                <span class="font-bold text-blue-600">{qr.scans_count}</span>
                              </div>
                              <div class="flex items-center space-x-2">
                                <span class="text-sm text-gray-500">Responses:</span>
                                <span class="font-bold text-green-600">{qr.feedback_count}</span>
                              </div>
                            </div>
                          </div>
                        </div>
                        <div class="flex items-center justify-between">
                          <span class="text-sm text-gray-600">Conversion Rate:</span>
                          <span class="font-bold text-indigo-600">{(qr.conversion_rate || 0).toFixed(1)}%</span>
                        </div>
                      </div>
                    </div>
                  {/each}
                </div>
              {:else}
                <div class="text-center py-8">
                  <div class="h-12 w-12 bg-gray-100 rounded-xl flex items-center justify-center mx-auto mb-3">
                    <QrCode class="h-6 w-6 text-gray-400" />
                  </div>
                  <p class="text-gray-500">No QR code data available</p>
                </div>
              {/if}
            </div>
          </Card>
        </div>
      {/if}
    {/if}
</div>
