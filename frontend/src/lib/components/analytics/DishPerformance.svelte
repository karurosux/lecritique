<script lang="ts">
  import { Card } from '$lib/components/ui';
  import { TrendingUpIcon, TrendingDownIcon, StarIcon, Utensils } from 'lucide-svelte';
  
  interface DishData {
    id: string;
    name: string;
    average_rating: number;
    feedback_count: number;
    weekly_change?: number;
  }
  
  interface AnalyticsData {
    top_dishes?: DishData[];
    top_rated_dishes?: DishData[];
    bottom_dishes?: DishData[];
    average_rating?: number;
    total_feedback?: number;
  }
  
  interface Feedback {
    dish_name?: string;
    overall_rating?: number;
    rating?: number;
    created_at: string;
  }
  
  let {
    analyticsData = null,
    feedbacks = [],
    loading = false
  }: {
    analyticsData?: AnalyticsData | null;
    feedbacks?: Feedback[];
    loading?: boolean;
  } = $props();

  function getPerformanceColor(rating: number): string {
    if (rating >= 4.5) return 'text-green-600';
    if (rating >= 3.5) return 'text-blue-600';
    if (rating >= 2.5) return 'text-yellow-600';
    return 'text-red-600';
  }

  function getPerformanceLabel(rating: number): string {
    if (rating >= 4.5) return 'Excellent';
    if (rating >= 3.5) return 'Good';
    if (rating >= 2.5) return 'Fair';
    return 'Poor';
  }

  function getTrendIcon(change: number) {
    if (change > 0) return TrendingUpIcon;
    if (change < 0) return TrendingDownIcon;
    return null;
  }

  function getTrendColor(change: number): string {
    if (change > 0) return 'text-green-600';
    if (change < 0) return 'text-red-600';
    return 'text-gray-600';
  }

  // Process dishes from analytics data or feedbacks
  const processedDishes = $derived(() => {
    if (analyticsData?.top_dishes || analyticsData?.top_rated_dishes) {
      const dishes = analyticsData.top_dishes || analyticsData.top_rated_dishes || [];
      return {
        top: dishes.slice(0, 5),
        bottom: analyticsData.bottom_dishes?.slice(0, 3) || []
      };
    }
    
    // Process from feedbacks if no analytics data
    if (feedbacks.length > 0) {
      const dishStats = feedbacks.reduce((acc, f) => {
        if (f.dish_name) {
          if (!acc[f.dish_name]) {
            acc[f.dish_name] = { ratings: [], count: 0 };
          }
          const rating = f.overall_rating || f.rating || 0;
          if (rating > 0) {
            acc[f.dish_name].ratings.push(rating);
            acc[f.dish_name].count++;
          }
        }
        return acc;
      }, {} as Record<string, { ratings: number[]; count: number }>);
      
      const dishes = Object.entries(dishStats)
        .map(([name, data]) => ({
          id: name,
          name,
          average_rating: data.ratings.reduce((sum, r) => sum + r, 0) / data.ratings.length,
          feedback_count: data.count
        }))
        .filter(d => d.feedback_count > 0)
        .sort((a, b) => b.average_rating - a.average_rating);
      
      return {
        top: dishes.slice(0, 5),
        bottom: dishes.filter(d => d.average_rating < 3).slice(0, 3)
      };
    }
    
    return { top: [], bottom: [] };
  });
  
  const overallStats = $derived(() => {
    if (analyticsData) {
      return {
        avgRating: analyticsData.average_rating || 0,
        totalFeedback: analyticsData.total_feedback || 0
      };
    }
    
    if (feedbacks.length > 0) {
      const ratings = feedbacks
        .map(f => f.overall_rating || f.rating || 0)
        .filter(r => r > 0);
      
      return {
        avgRating: ratings.length > 0 ? ratings.reduce((sum, r) => sum + r, 0) / ratings.length : 0,
        totalFeedback: feedbacks.length
      };
    }
    
    return { avgRating: 0, totalFeedback: 0 };
  });

</script>

<div class="space-y-6">
  <!-- Overall Performance Summary -->
  <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
    <Card variant="gradient" hover interactive>
      <div class="flex items-center justify-between">
        <div class="space-y-2">
          <p class="text-sm font-semibold text-gray-600 uppercase tracking-wide">Overall Rating</p>
          <p class="text-3xl font-bold bg-gradient-to-r from-yellow-600 to-orange-600 bg-clip-text text-transparent">
            {overallStats().avgRating.toFixed(1)}
          </p>
          <div class="flex text-yellow-400">
            {#each Array(5) as _, i}
              <StarIcon class="h-4 w-4 {i < Math.round(overallStats().avgRating) ? 'fill-current' : 'text-gray-300'}" />
            {/each}
          </div>
        </div>
        <div class="h-16 w-16 bg-gradient-to-br from-yellow-500 to-orange-500 rounded-2xl flex items-center justify-center shadow-lg shadow-yellow-500/25">
          <StarIcon class="h-8 w-8 text-white fill-current" />
        </div>
      </div>
    </Card>
    
    <Card variant="gradient" hover interactive>
      <div class="flex items-center justify-between">
        <div class="space-y-2">
          <p class="text-sm font-semibold text-gray-600 uppercase tracking-wide">Total Reviews</p>
          <p class="text-3xl font-bold bg-gradient-to-r from-blue-600 to-purple-600 bg-clip-text text-transparent">
            {overallStats().totalFeedback}
          </p>
          <p class="text-sm text-gray-600">Dish-specific feedback</p>
        </div>
        <div class="h-16 w-16 bg-gradient-to-br from-blue-500 to-purple-600 rounded-2xl flex items-center justify-center shadow-lg shadow-blue-500/25">
          <Utensils class="h-8 w-8 text-white" />
        </div>
      </div>
    </Card>
    
    <Card variant="gradient" hover interactive>
      <div class="flex items-center justify-between">
        <div class="space-y-2">
          <p class="text-sm font-semibold text-gray-600 uppercase tracking-wide">Top Rated</p>
          <p class="text-3xl font-bold bg-gradient-to-r from-green-600 to-emerald-600 bg-clip-text text-transparent">
            {processedDishes().top.length}
          </p>
          <p class="text-sm text-gray-600">High-performing dishes</p>
        </div>
        <div class="h-16 w-16 bg-gradient-to-br from-green-500 to-emerald-500 rounded-2xl flex items-center justify-center shadow-lg shadow-green-500/25">
          <TrendingUpIcon class="h-8 w-8 text-white" />
        </div>
      </div>
    </Card>
    
    <Card variant="gradient" hover interactive>
      <div class="flex items-center justify-between">
        <div class="space-y-2">
          <p class="text-sm font-semibold text-gray-600 uppercase tracking-wide">Need Attention</p>
          <p class="text-3xl font-bold bg-gradient-to-r from-red-600 to-pink-600 bg-clip-text text-transparent">
            {processedDishes().bottom.length}
          </p>
          <p class="text-sm text-gray-600">Below expectations</p>
        </div>
        <div class="h-16 w-16 bg-gradient-to-br from-red-500 to-pink-500 rounded-2xl flex items-center justify-center shadow-lg shadow-red-500/25">
          <TrendingDownIcon class="h-8 w-8 text-white" />
        </div>
      </div>
    </Card>
  </div>

  <!-- Dish Rankings -->
  <Card variant="elevated" padding={false}>
    <div class="p-6">
      <div class="mb-6">
        <h3 class="text-lg font-semibold bg-gradient-to-r from-gray-900 to-gray-700 bg-clip-text text-transparent">
          Dish Performance Rankings
        </h3>
        <p class="text-sm text-gray-600 mt-1">How your dishes are performing based on customer feedback</p>
      </div>

  {#if loading}
    <div class="space-y-6">
      <!-- Dishes loading -->
      <div class="space-y-4">
        {#each Array(5) as _}
          <div class="animate-pulse">
            <div class="bg-gray-100 border border-gray-200 rounded-lg p-4">
              <div class="flex items-center justify-between">
                <div class="h-4 bg-gray-200 rounded w-1/3"></div>
                <div class="h-4 bg-gray-200 rounded w-16"></div>
              </div>
            </div>
          </div>
        {/each}
      </div>
    </div>
  {:else if processedDishes().top.length > 0 || processedDishes().bottom.length > 0}
    <!-- Top Performing Dishes -->
    {#if processedDishes().top.length > 0}
    <div class="mb-6">
      <h4 class="font-medium text-gray-900 mb-3">Top Performing Dishes</h4>
      <div class="space-y-3">
        {#each processedDishes().top as dish, index}
          <div class="bg-white border border-gray-200 rounded-lg p-4 hover:shadow-md transition-shadow">
            <div class="flex items-center justify-between">
              <div class="flex-1">
                <div class="flex items-center gap-3">
                  <div class="flex items-center justify-center w-8 h-8 bg-blue-100 text-blue-600 rounded-full text-sm font-medium">
                    {index + 1}
                  </div>
                  <div>
                    <h5 class="font-medium text-gray-900">{dish.name}</h5>
                    <div class="flex items-center gap-2 mt-1">
                      <span class="text-xs text-gray-500">{dish.feedback_count} reviews</span>
                    </div>
                  </div>
                </div>
              </div>
              
              <div class="text-right">
                <div class="flex items-center gap-2">
                  <div class="flex items-center gap-1">
                    <StarIcon class="w-4 h-4 fill-yellow-400 text-yellow-400" />
                    <span class="font-medium {getPerformanceColor(dish.average_rating)}">
                      {dish.average_rating.toFixed(1)}
                    </span>
                  </div>
                  {#if dish.weekly_change !== undefined}
                    {@const TrendIcon = getTrendIcon(dish.weekly_change)}
                    {#if TrendIcon}
                      <div class="flex items-center gap-1">
                        <TrendIcon class="w-3 h-3 {getTrendColor(dish.weekly_change)}" />
                        <span class="text-xs {getTrendColor(dish.weekly_change)}">
                          {Math.abs(dish.weekly_change).toFixed(1)}%
                        </span>
                      </div>
                    {/if}
                  {/if}
                </div>
                <div class="text-xs text-gray-500 mt-1">
                  {getPerformanceLabel(dish.average_rating)}
                </div>
              </div>
            </div>
          </div>
        {/each}
      </div>
    </div>
    {/if}

    <!-- Bottom Performing Dishes -->
    {#if processedDishes().bottom.length > 0}
      <div>
        <h4 class="font-medium text-gray-900 mb-3">Areas for Improvement</h4>
        <div class="space-y-3">
          {#each processedDishes().bottom as dish, index}
            <div class="bg-red-50 border border-red-200 rounded-lg p-4">
              <div class="flex items-center justify-between">
                <div class="flex-1">
                  <div class="flex items-center gap-3">
                    <div class="flex items-center justify-center w-8 h-8 bg-red-100 text-red-600 rounded-full text-sm font-medium">
                      {index + 1}
                    </div>
                    <div>
                      <h5 class="font-medium text-red-900">{dish.name}</h5>
                      <div class="flex items-center gap-2 mt-1">
                        <span class="text-xs text-red-600">{dish.feedback_count} reviews</span>
                      </div>
                    </div>
                  </div>
                </div>
                
                <div class="text-right">
                  <div class="flex items-center gap-2">
                    <div class="flex items-center gap-1">
                      <StarIcon class="w-4 h-4 fill-red-400 text-red-400" />
                      <span class="font-medium text-red-600">
                        {dish.average_rating.toFixed(1)}
                      </span>
                    </div>
                  </div>
                  <div class="text-xs text-red-500 mt-1">
                    Needs attention
                  </div>
                </div>
              </div>
            </div>
          {/each}
        </div>
      </div>
    {/if}
  {:else}
    <Card variant="elevated">
      <div class="text-center py-8">
        <div class="text-gray-500 text-sm">
          No dish performance data available. Ensure you have collected feedback for your dishes.
        </div>
      </div>
    </Card>
  {/if}
    </div>
  </Card>
</div>