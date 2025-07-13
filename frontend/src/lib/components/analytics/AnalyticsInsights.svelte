<script lang="ts">
  import { Card } from '$lib/components/ui';
  import { Star, Heart, AlertTriangle } from 'lucide-svelte';

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
    
    const ratingDist = analyticsData.rating_distribution || {};
    const highRatings = (ratingDist['4'] || 0) + (ratingDist['5'] || 0);
    const totalFeedback = analyticsData.total_feedback || 0;
    const highRatingPercentage = getPercentage(highRatings, totalFeedback);
    const avgRating = analyticsData.average_rating || 0;
    
    // Calculate dish-specific insights
    const topDishes = analyticsData.top_rated_dishes || analyticsData.top_dishes || [];
    const bestDish = topDishes[0];
    const worstRatings = (ratingDist['1'] || 0) + (ratingDist['2'] || 0);
    const worstRatingPercentage = getPercentage(worstRatings, totalFeedback);
    
    const insights = [];
    
    // Best performing dish
    if (bestDish) {
      insights.push({
        id: 'best-dish',
        title: bestDish.dish_name || bestDish.name,
        subtitle: `Top rated: ${(bestDish.average_rating || 0).toFixed(1)} stars`,
        icon: Star,
        bgClass: 'bg-green-50',
        borderClass: 'border-green-200',
        iconClass: 'text-green-500',
        textClass: 'text-green-800',
        subtitleClass: 'text-green-600'
      });
    }
    
    // Customer satisfaction
    insights.push({
      id: 'satisfaction',
      title: `${highRatingPercentage.toFixed(0)}% Love It`,
      subtitle: 'Customers rating 4-5 stars',
      icon: Heart,
      bgClass: 'bg-pink-50',
      borderClass: 'border-pink-200',
      iconClass: 'text-pink-500',
      textClass: 'text-pink-800',
      subtitleClass: 'text-pink-600'
    });
    
    // Areas of concern
    if (worstRatingPercentage > 10) {
      insights.push({
        id: 'concern',
        title: `${worstRatingPercentage.toFixed(0)}% Need Work`,
        subtitle: 'Dishes rated 1-2 stars',
        icon: AlertTriangle,
        bgClass: 'bg-amber-50',
        borderClass: 'border-amber-200',
        iconClass: 'text-amber-500',
        textClass: 'text-amber-800',
        subtitleClass: 'text-amber-600'
      });
    }
    
    return insights;
  });
</script>

<div class="space-y-4">
  <div class="grid grid-cols-1 md:grid-cols-3 gap-4">
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
              <svelte:component this={insight.icon} class="h-5 w-5 {insight.iconClass}" />
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
</style>