<script lang="ts">
  import { Card } from '$lib/components/ui';

  interface OperationalMetrics {
    total_feedbacks: number;
    todays_feedback: number;
    active_qr_codes: number;
    total_qr_scans: number;
    completion_rate: number;
    average_response_time: number;
  }

  let {
    analyticsData = null,
    loading = false
  }: {
    analyticsData?: OperationalMetrics | null;
    loading?: boolean;
  } = $props();

  function formatTime(minutes: number): string {
    if (minutes < 1) return '< 1 min';
    if (minutes < 60) return `${Math.round(minutes)} min`;
    const hours = Math.floor(minutes / 60);
    const mins = Math.round(minutes % 60);
    return `${hours}h ${mins}m`;
  }

  const summaryCards = $derived(analyticsData ? [
    {
      id: 'total',
      title: 'Total Responses',
      value: analyticsData.total_feedbacks || 0,
      icon: 'M9 5H7a2 2 0 00-2 2v12a2 2 0 002 2h10a2 2 0 002-2V7a2 2 0 00-2-2h-2M9 5a2 2 0 002 2h2a2 2 0 002-2M9 5a2 2 0 012-2h2a2 2 0 012 2m-3 7h3m-3 4h3m-6-4h3m-3 4h3',
      colorClass: 'from-blue-500 to-blue-600',
      bgClass: 'bg-blue-50',
      iconColorClass: 'text-blue-600'
    },
    {
      id: 'today',
      title: 'Today\'s Activity',
      value: analyticsData.todays_feedback || 0,
      icon: 'M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z',
      colorClass: 'from-green-500 to-green-600',
      bgClass: 'bg-green-50',
      iconColorClass: 'text-green-600'
    },
    {
      id: 'qrcodes',
      title: 'Active QR Codes',
      value: analyticsData.active_qr_codes || 0,
      icon: 'M3 3h6v6H3V3zm12 0h6v6h-6V3zM3 15h6v6H3v-6zm12 0h6v6h-6v-6zm-6-6h6v6H9V9z',
      colorClass: 'from-purple-500 to-purple-600',
      bgClass: 'bg-purple-50',
      iconColorClass: 'text-purple-600'
    },
    {
      id: 'completion',
      title: 'Completion Rate',
      value: `${(analyticsData.completion_rate ?? 0).toFixed(1)}%`,
      subtitle: `${analyticsData.total_qr_scans || 0} total scans`,
      icon: 'M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z',
      colorClass: 'from-indigo-500 to-indigo-600',
      bgClass: 'bg-indigo-50',
      iconColorClass: 'text-indigo-600'
    },
    {
      id: 'response-time',
      title: 'Avg Response Time',
      value: formatTime(analyticsData.average_response_time || 0),
      subtitle: 'from scan to submit',
      icon: 'M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z',
      colorClass: 'from-orange-500 to-orange-600',
      bgClass: 'bg-orange-50',
      iconColorClass: 'text-orange-600'
    }
  ] : []);
</script>

<div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-5 gap-6 mb-8">
  {#if loading}
    {#each Array(5) as _}
      <Card variant="default" class="opacity-50">
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
    {#each summaryCards as card, index}
      <Card variant="default" hover interactive class="group transform transition-all duration-300 animate-fade-in-up" style="animation-delay: {index * 100}ms">
        <div class="relative">
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
      </Card>
    {/each}
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
</style>