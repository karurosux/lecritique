<script lang="ts">
  import { Card } from '$lib/components/ui';

  interface AnalyticsData {
    total_feedback: number;
    rating_distribution: Record<string, number>;
    top_dishes: Array<{
      id: string;
      name: string;
      average_rating: number;
      feedback_count: number;
    }>;
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

  function renderStars(rating: number): string {
    return '★'.repeat(Math.round(rating)) + '☆'.repeat(5 - Math.round(rating));
  }

  function getRatingColor(rating: number): string {
    if (rating >= 4) return 'text-green-600';
    if (rating >= 3) return 'text-yellow-600';
    return 'text-red-600';
  }

  function getRatingBarColor(rating: number): string {
    if (rating >= 4) return 'from-green-400 to-green-600';
    if (rating >= 3) return 'from-yellow-400 to-yellow-600';
    return 'from-red-400 to-red-600';
  }
</script>

<div class="grid grid-cols-1 lg:grid-cols-2 gap-8 mb-8">
  {#if loading}
    {#each Array(2) as _}
      <Card variant="glass">
        <div class="animate-pulse">
          <div class="h-5 bg-gray-200 rounded w-1/3 mb-2"></div>
          <div class="h-3 bg-gray-200 rounded w-1/2 mb-6"></div>
          <div class="space-y-4">
            {#each Array(5) as _}
              <div class="flex items-center gap-3">
                <div class="h-4 bg-gray-200 rounded w-12"></div>
                <div class="flex-1 h-2 bg-gray-200 rounded"></div>
                <div class="h-4 bg-gray-200 rounded w-16"></div>
              </div>
            {/each}
          </div>
        </div>
      </Card>
    {/each}
  {:else if analyticsData}
    <!-- Rating Distribution -->
    <Card variant="glass" class="hover:shadow-lg transition-all duration-300">
      <div class="mb-6">
        <h3 class="text-lg font-semibold bg-gradient-to-r from-gray-900 to-gray-700 bg-clip-text text-transparent">
          Rating Distribution
        </h3>
        <p class="text-sm text-gray-600 mt-1">Breakdown of customer ratings</p>
      </div>
      
      <div class="space-y-4">
        {#each [5, 4, 3, 2, 1] as rating}
          {@const count = analyticsData.rating_distribution[rating.toString()] || 0}
          {@const percentage = getPercentage(count, analyticsData.total_feedback)}
          <div class="group">
            <div class="flex items-center gap-3">
              <div class="w-14 text-sm font-medium text-gray-600 group-hover:text-gray-900 transition-colors">
                {rating} ★
              </div>
              <div class="flex-1">
                <div class="bg-gray-100 rounded-full h-2.5 overflow-hidden">
                  <div
                    class="h-full bg-gradient-to-r {getRatingBarColor(rating)} transition-all duration-500 ease-out rounded-full"
                    style="width: {percentage}%"
                  ></div>
                </div>
              </div>
              <div class="w-20 text-sm text-gray-600 text-right group-hover:text-gray-900 transition-colors">
                <span class="font-medium">{count}</span>
                <span class="text-gray-400 ml-1">({percentage.toFixed(1)}%)</span>
              </div>
            </div>
          </div>
        {/each}
      </div>
    </Card>

    <!-- Top Performing Dishes -->
    <Card variant="glass" class="hover:shadow-lg transition-all duration-300">
      <div class="mb-6">
        <h3 class="text-lg font-semibold bg-gradient-to-r from-gray-900 to-gray-700 bg-clip-text text-transparent">
          Top Performing Dishes
        </h3>
        <p class="text-sm text-gray-600 mt-1">Highest rated dishes with feedback</p>
      </div>
      
      <div class="space-y-4">
        {#each analyticsData.top_dishes.slice(0, 5) as dish, index}
          <div class="group flex items-center justify-between p-3 rounded-xl hover:bg-gray-50 transition-all duration-200">
            <div class="flex items-center space-x-3">
              <div class="flex items-center justify-center w-10 h-10 bg-gradient-to-br from-purple-100 to-pink-100 rounded-xl text-sm font-semibold text-purple-700 group-hover:scale-110 transition-transform duration-200">
                {index + 1}
              </div>
              <div>
                <div class="font-medium text-gray-900 group-hover:text-purple-700 transition-colors">
                  {dish.name}
                </div>
                <div class="text-sm text-gray-500">
                  {dish.feedback_count} {dish.feedback_count === 1 ? 'review' : 'reviews'}
                </div>
              </div>
            </div>
            <div class="text-right">
              <div class="font-semibold {getRatingColor(dish.average_rating)} text-lg">
                {dish.average_rating.toFixed(1)}
              </div>
              <div class="text-xs {getRatingColor(dish.average_rating)}">
                {renderStars(dish.average_rating)}
              </div>
            </div>
          </div>
        {:else}
          <div class="text-center py-12">
            <div class="w-16 h-16 bg-gray-100 rounded-full flex items-center justify-center mx-auto mb-4">
              <svg class="h-8 w-8 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 5H7a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2" />
              </svg>
            </div>
            <p class="text-gray-500">No dish data available</p>
          </div>
        {/each}
      </div>
    </Card>
  {/if}
</div>