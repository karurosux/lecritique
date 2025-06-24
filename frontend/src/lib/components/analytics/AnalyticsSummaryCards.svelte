<script lang="ts">
  import { Card } from '$lib/components/ui';

  interface AnalyticsData {
    total_feedback: number;
    average_rating: number;
    feedback_today: number;
    feedback_this_week: number;
    feedback_this_month: number;
  }

  let {
    analyticsData = null,
    loading = false
  }: {
    analyticsData?: AnalyticsData | null;
    loading?: boolean;
  } = $props();

  function renderStars(rating: number): string {
    return '★'.repeat(Math.round(rating)) + '☆'.repeat(5 - Math.round(rating));
  }

  const summaryCards = $derived(analyticsData ? [
    {
      id: 'total',
      title: 'Total Feedback',
      value: analyticsData.total_feedback,
      icon: 'M9 5H7a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2m-3 7h3m-3 4h3m-6-4h3m-3 4h3',
      colorClass: 'from-blue-500 to-blue-600',
      bgClass: 'bg-blue-50',
      iconColorClass: 'text-blue-600'
    },
    {
      id: 'rating',
      title: 'Average Rating',
      value: analyticsData.average_rating.toFixed(1),
      subtitle: renderStars(analyticsData.average_rating),
      icon: 'M11.049 2.927c.3-.921 1.603-.921 1.902 0l1.519 4.674a1 1 0 00.95.69h4.915c.969 0 1.371 1.24.588 1.81l-3.976 2.888a1 1 0 00-.363 1.118l1.518 4.674c.3.922-.755 1.688-1.538 1.118l-3.976-2.888a1 1 0 00-1.176 0l-3.976 2.888c-.783.57-1.838-.197-1.538-1.118l1.518-4.674a1 1 0 00-.363-1.118l-3.976-2.888c-.784-.57-.38-1.81.588-1.81h4.914a1 1 0 00.951-.69l1.519-4.674z',
      colorClass: 'from-yellow-500 to-yellow-600',
      bgClass: 'bg-yellow-50',
      iconColorClass: 'text-yellow-600'
    },
    {
      id: 'today',
      title: 'Today',
      value: analyticsData.feedback_today,
      icon: 'M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z',
      colorClass: 'from-green-500 to-green-600',
      bgClass: 'bg-green-50',
      iconColorClass: 'text-green-600'
    },
    {
      id: 'week',
      title: 'This Week',
      value: analyticsData.feedback_this_week,
      icon: 'M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z',
      colorClass: 'from-purple-500 to-purple-600',
      bgClass: 'bg-purple-50',
      iconColorClass: 'text-purple-600'
    },
    {
      id: 'month',
      title: 'This Month',
      value: analyticsData.feedback_this_month,
      icon: 'M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z',
      colorClass: 'from-indigo-500 to-indigo-600',
      bgClass: 'bg-indigo-50',
      iconColorClass: 'text-indigo-600'
    }
  ] : []);
</script>

<div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-5 gap-6 mb-8">
  {#if loading}
    {#each Array(5) as _}
      <Card variant="glass">
        <div class="animate-pulse">
          <div class="flex items-center justify-between mb-4">
            <div class="h-10 w-10 bg-gray-200 rounded-xl"></div>
            <div class="h-4 bg-gray-200 rounded w-20"></div>
          </div>
          <div class="h-8 bg-gray-200 rounded w-24 mb-2"></div>
          <div class="h-3 bg-gray-200 rounded w-16"></div>
        </div>
      </Card>
    {/each}
  {:else if analyticsData}
    {#each summaryCards as card}
      <Card variant="glass" class="group hover:shadow-lg transition-all duration-300">
        <div class="relative overflow-hidden">
          <div class="absolute inset-0 bg-gradient-to-br {card.colorClass} opacity-0 group-hover:opacity-5 transition-opacity duration-300"></div>
          <div class="relative z-10">
            <div class="flex items-center justify-between mb-4">
              <div class="{card.bgClass} p-2.5 rounded-xl group-hover:scale-110 transition-transform duration-300">
                <svg class="h-5 w-5 {card.iconColorClass}" fill="currentColor" viewBox="0 0 24 24">
                  <path d={card.icon} />
                </svg>
              </div>
              <span class="text-xs font-medium text-gray-500 uppercase tracking-wider">{card.title}</span>
            </div>
            <div class="space-y-1">
              <div class="text-3xl font-bold bg-gradient-to-r {card.colorClass} bg-clip-text text-transparent">
                {card.value}
              </div>
              {#if card.subtitle}
                <div class="text-sm {card.iconColorClass}">{card.subtitle}</div>
              {/if}
            </div>
          </div>
        </div>
      </Card>
    {/each}
  {/if}
</div>