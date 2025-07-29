<script lang="ts">
  import { Card } from '$lib/components/ui';
  import {
    TrendingUp,
    Users,
    Star,
    BarChart3,
    Calendar,
    Clock,
    MessageSquare,
    ThumbsUp,
  } from 'lucide-svelte';

  interface AnalyticsData {
    total_feedback?: number;
    average_rating?: number;
    feedback_today?: number;
    feedback_this_week?: number;
    feedback_this_month?: number;
    completion_rate?: number;
    positive_sentiment_percentage?: number;
    total_feedbacks?: number;
    todays_feedback?: number;
    active_qr_codes?: number;
    total_qr_scans?: number;
    average_response_time?: number;
  }

  interface QuickStats {
    total: number;
    avgRating: number;
    thisWeek: number;
    hasData: boolean;
  }

  let {
    analyticsData = null,
    quickStats = null,
    loading = false,
  }: {
    analyticsData?: AnalyticsData | null;
    quickStats?: QuickStats | null;
    loading?: boolean;
  } = $props();

  // Compute display stats
  let displayStats = $derived(() => {
    if (analyticsData) {
      return {
        totalFeedback:
          analyticsData.total_feedback || analyticsData.total_feedbacks || 0,
        averageRating: analyticsData.average_rating || 0,
        feedbackToday:
          analyticsData.feedback_today || analyticsData.todays_feedback || 0,
        feedbackThisWeek: analyticsData.feedback_this_week || 0,
        completionRate: analyticsData.completion_rate,
        positiveSentiment: analyticsData.positive_sentiment_percentage,
      };
    } else if (quickStats?.hasData) {
      return {
        totalFeedback: quickStats.total,
        averageRating: quickStats.avgRating,
        feedbackToday: 0,
        feedbackThisWeek: quickStats.thisWeek,
        completionRate: null,
        positiveSentiment: null,
      };
    }
    return null;
  });

  const cardConfigs = [
    {
      id: 'total',
      title: 'Total Responses',
      getValue: (stats: any) => stats.totalFeedback,
      format: (v: number) => v.toString(),
      subtitle: 'All time feedback',
      gradientFrom: 'from-blue-600',
      gradientTo: 'to-purple-600',
      bgGradientFrom: 'from-blue-500',
      bgGradientTo: 'to-purple-600',
      shadowColor: 'shadow-blue-500/25',
      icon: MessageSquare,
      trendIcon: TrendingUp,
      trendColor: 'text-blue-600',
    },
    {
      id: 'rating',
      title: 'Average Rating',
      getValue: (stats: any) => stats.averageRating,
      format: (v: number) => v.toFixed(1),
      subtitle: null,
      gradientFrom: 'from-yellow-600',
      gradientTo: 'to-orange-600',
      bgGradientFrom: 'from-yellow-500',
      bgGradientTo: 'to-orange-500',
      shadowColor: 'shadow-yellow-500/25',
      icon: Star,
      showStars: true,
    },
    {
      id: 'week',
      title: 'This Week',
      getValue: (stats: any) => stats.feedbackThisWeek,
      format: (v: number) => v.toString(),
      subtitle: 'Weekly responses',
      gradientFrom: 'from-green-600',
      gradientTo: 'to-emerald-600',
      bgGradientFrom: 'from-green-500',
      bgGradientTo: 'to-emerald-500',
      shadowColor: 'shadow-green-500/25',
      icon: Calendar,
      trendIcon: TrendingUp,
      trendColor: 'text-green-600',
    },
    {
      id: 'today',
      title: 'Today',
      getValue: (stats: any) => stats.feedbackToday,
      format: (v: number) => v.toString(),
      subtitle: "Today's responses",
      gradientFrom: 'from-purple-600',
      gradientTo: 'to-indigo-600',
      bgGradientFrom: 'from-purple-500',
      bgGradientTo: 'to-indigo-500',
      shadowColor: 'shadow-purple-500/25',
      icon: Clock,
      trendIcon: Clock,
      trendColor: 'text-purple-600',
    },
  ];

  // Additional metrics if available
  const additionalMetrics = $derived(() => {
    const stats = displayStats();
    if (!stats) return [];

    const metrics = [];

    if (stats.completionRate != null) {
      metrics.push({
        id: 'completion',
        title: 'Completion Rate',
        value: stats.completionRate,
        format: (v: number) => `${v.toFixed(1)}%`,
        subtitle: 'QR scan to response',
        gradientFrom: 'from-indigo-600',
        gradientTo: 'to-blue-600',
        bgGradientFrom: 'from-indigo-500',
        bgGradientTo: 'to-blue-500',
        shadowColor: 'shadow-indigo-500/25',
        icon: Users,
      });
    }

    if (stats.positiveSentiment != null) {
      metrics.push({
        id: 'sentiment',
        title: 'Positive Sentiment',
        value: stats.positiveSentiment,
        format: (v: number) => `${v.toFixed(1)}%`,
        subtitle: 'Happy customers',
        gradientFrom: 'from-pink-600',
        gradientTo: 'to-rose-600',
        bgGradientFrom: 'from-pink-500',
        bgGradientTo: 'to-rose-500',
        shadowColor: 'shadow-pink-500/25',
        icon: ThumbsUp,
      });
    }

    return metrics;
  });
</script>

{#if loading}
  <!-- Loading State -->
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
{:else if displayStats()}
  <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-6">
    {#each cardConfigs as config}
      {@const value = config.getValue(displayStats())}
      <Card variant="gradient" hover interactive class="analytics-summary-card">
        <div class="flex items-center justify-between">
          <div class="space-y-2">
            <p
              class="text-sm font-semibold text-gray-600 uppercase tracking-wide">
              {config.title}
            </p>
            <p
              class="text-3xl font-bold bg-gradient-to-r {config.gradientFrom} {config.gradientTo} bg-clip-text text-transparent">
              {config.format(value)}
            </p>
            {#if config.showStars}
              <div class="flex text-yellow-400">
                {#each Array(5) as _, i}
                  <Star
                    class="h-4 w-4 {i < Math.round(value)
                      ? 'fill-current'
                      : 'text-gray-300'}" />
                {/each}
              </div>
            {:else if config.subtitle}
              <div
                class="flex items-center space-x-1 {config.trendColor ||
                  'text-gray-600'}">
                {#if config.trendIcon}
                  <svelte:component this={config.trendIcon} class="h-4 w-4" />
                {/if}
                <span class="text-sm font-medium">{config.subtitle}</span>
              </div>
            {/if}
          </div>
          <div
            class="h-16 w-16 bg-gradient-to-br {config.bgGradientFrom} {config.bgGradientTo} rounded-2xl flex items-center justify-center shadow-lg {config.shadowColor}">
            <svelte:component this={config.icon} class="h-8 w-8 text-white" />
          </div>
        </div>
      </Card>
    {/each}
  </div>

  <!-- Additional Metrics Row -->
  {#if additionalMetrics().length > 0}
    <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6 mt-6">
      {#each additionalMetrics() as metric}
        <Card
          variant="elevated"
          hover
          interactive
          class="analytics-metric-card">
          <div class="flex items-center justify-between">
            <div class="space-y-2">
              <p
                class="text-sm font-semibold text-gray-600 uppercase tracking-wide">
                {metric.title}
              </p>
              <p
                class="text-2xl font-bold bg-gradient-to-r {metric.gradientFrom} {metric.gradientTo} bg-clip-text text-transparent">
                {metric.format(metric.value)}
              </p>
              <p class="text-sm text-gray-500">{metric.subtitle}</p>
            </div>
            <div
              class="h-12 w-12 bg-gradient-to-br {metric.bgGradientFrom} {metric.bgGradientTo} rounded-xl flex items-center justify-center shadow-md {metric.shadowColor}">
              <svelte:component this={metric.icon} class="h-6 w-6 text-white" />
            </div>
          </div>
        </Card>
      {/each}
    </div>
  {/if}
{/if}
