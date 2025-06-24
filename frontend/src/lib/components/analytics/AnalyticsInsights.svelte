<script lang="ts">
  import { Card } from '$lib/components/ui';

  interface AnalyticsData {
    total_feedback: number;
    average_rating: number;
    feedback_this_week: number;
    feedback_today: number;
    rating_distribution: Record<string, number>;
  }

  let {
    analyticsData = null,
    loading = false
  }: {
    analyticsData?: AnalyticsData | null;
    loading?: boolean;
  } = $props();

  function getPercentage(value: number, total: number): number {
    return total > 0 ? (value / total) * 100 : 0;
  }

  const insights = $derived(() => {
    if (!analyticsData) return [];
    
    const highRatings = (analyticsData.rating_distribution['4'] || 0) + (analyticsData.rating_distribution['5'] || 0);
    const highRatingPercentage = getPercentage(highRatings, analyticsData.total_feedback);
    const satisfactionLevel = analyticsData.average_rating >= 4 ? 'Excellent' : analyticsData.average_rating >= 3 ? 'Good' : 'Needs Improvement';
    const feedbackTrend = analyticsData.feedback_this_week > analyticsData.feedback_today * 7 / 2 ? 'Growing' : 'Stable';
    
    return [
      {
        id: 'high-ratings',
        title: `${highRatingPercentage.toFixed(1)}% High Ratings`,
        subtitle: '4-5 star reviews',
        icon: 'M11.049 2.927c.3-.921 1.603-.921 1.902 0l1.519 4.674a1 1 0 00.95.69h4.915c.969 0 1.371 1.24.588 1.81l-3.976 2.888a1 1 0 00-.363 1.118l1.518 4.674c.3.922-.755 1.688-1.538 1.118l-3.976-2.888a1 1 0 00-1.176 0l-3.976 2.888c-.783.57-1.838-.197-1.538-1.118l1.518-4.674a1 1 0 00-.363-1.118l-3.976-2.888c-.784-.57-.38-1.81.588-1.81h4.914a1 1 0 00.951-.69l1.519-4.674z',
        bgClass: 'bg-green-50',
        borderClass: 'border-green-200',
        iconClass: 'text-green-500',
        textClass: 'text-green-800',
        subtitleClass: 'text-green-600'
      },
      {
        id: 'satisfaction',
        title: satisfactionLevel,
        subtitle: 'Overall satisfaction',
        icon: 'M14.828 14.828a4 4 0 01-5.656 0M9 10h1.01M15 10h1.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z',
        bgClass: 'bg-blue-50',
        borderClass: 'border-blue-200',
        iconClass: 'text-blue-500',
        textClass: 'text-blue-800',
        subtitleClass: 'text-blue-600'
      },
      {
        id: 'trend',
        title: feedbackTrend,
        subtitle: 'Feedback trend',
        icon: 'M13 7h8m0 0v8m0-8l-8 8-4-4-6 6',
        bgClass: 'bg-purple-50',
        borderClass: 'border-purple-200',
        iconClass: 'text-purple-500',
        textClass: 'text-purple-800',
        subtitleClass: 'text-purple-600'
      }
    ];
  });
</script>

<Card variant="glass" class="hover:shadow-lg transition-all duration-300">
  <div class="mb-6">
    <h3 class="text-lg font-semibold bg-gradient-to-r from-gray-900 to-gray-700 bg-clip-text text-transparent">
      Key Insights
    </h3>
    <p class="text-sm text-gray-600 mt-1">Automated analysis and recommendations</p>
  </div>
  
  <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
    {#if loading}
      {#each Array(3) as _}
        <div class="animate-pulse">
          <div class="bg-gray-100 border border-gray-200 rounded-xl p-4">
            <div class="flex items-center gap-3">
              <div class="h-10 w-10 bg-gray-200 rounded-lg"></div>
              <div class="flex-1">
                <div class="h-4 bg-gray-200 rounded w-3/4 mb-1"></div>
                <div class="h-3 bg-gray-200 rounded w-1/2"></div>
              </div>
            </div>
          </div>
        </div>
      {/each}
    {:else if analyticsData}
      {#each insights() as insight}
        <div class="group {insight.bgClass} border {insight.borderClass} rounded-xl p-4 hover:shadow-md transition-all duration-300">
          <div class="flex items-center">
            <div class="p-2 bg-white rounded-lg shadow-sm group-hover:scale-110 transition-transform duration-200">
              <svg class="h-5 w-5 {insight.iconClass}" fill="currentColor" viewBox="0 0 24 24">
                <path d={insight.icon} />
              </svg>
            </div>
            <div class="ml-3">
              <div class="font-medium {insight.textClass}">{insight.title}</div>
              <div class="text-sm {insight.subtitleClass}">{insight.subtitle}</div>
            </div>
          </div>
        </div>
      {/each}
    {/if}
  </div>
</Card>