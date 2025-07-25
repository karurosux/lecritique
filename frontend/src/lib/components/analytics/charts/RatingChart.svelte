<script lang="ts">
  import { onMount, onDestroy } from 'svelte';
  import { Chart, registerables } from 'chart.js';
  import 'chartjs-adapter-date-fns';
  import { Star, TrendingUp, TrendingDown, Minus } from 'lucide-svelte';

  interface Props {
    data: any;
    type: 'timeseries' | 'comparison';
  }

  let { data, type }: Props = $props();
  
  let chartCanvas: HTMLCanvasElement;
  let chart: Chart | null = null;
  let mounted = $state(false);

  Chart.register(...registerables);

  let series = $derived(data?.series || []);

  function formatRatingValue(value: number): string {
    return `${value.toFixed(2)}/5 ⭐`;
  }

  function getTrendIcon(trendDirection: string) {
    switch (trendDirection) {
      case 'improving': return TrendingUp;
      case 'declining': return TrendingDown;
      default: return Minus;
    }
  }

  function getTrendColor(trendDirection: string): string {
    switch (trendDirection) {
      case 'improving': return 'text-green-600';
      case 'declining': return 'text-red-600';
      default: return 'text-gray-600';
    }
  }

  function createChart() {
    if (!mounted || !chartCanvas || !data?.series) return;
    
    if (chart) {
      chart.destroy();
    }

    const datasets = data.series.map((seriesData: any, index: number) => {
      const points = (seriesData.points || []).map((point: any) => ({
        x: new Date(point.timestamp),
        y: point.value
      }));

      return {
        label: seriesData.metric_name,
        data: points,
        borderColor: `hsl(${45 + index * 60}, 70%, 50%)`, // Orange/yellow tones for ratings
        backgroundColor: `hsl(${45 + index * 60}, 70%, 50%)20`,
        borderWidth: 3,
        pointRadius: 6,
        pointHoverRadius: 8,
        fill: false,
        tension: 0.4
      };
    });

    const config = {
      type: 'line' as const,
      data: { datasets },
      options: {
        responsive: true,
        maintainAspectRatio: false,
        plugins: {
          title: {
            display: true,
            text: 'Rating Trends (1-5 Stars)',
            font: { size: 16, weight: 'bold' as const },
            color: '#1f2937'
          },
          legend: {
            display: true,
            position: 'top' as const
          },
          tooltip: {
            callbacks: {
              title: (context: any) => new Date(context[0].parsed.x).toLocaleDateString(),
              label: (context: any) => `${context.dataset.label}: ${formatRatingValue(context.parsed.y)}`
            }
          }
        },
        scales: {
          x: {
            type: 'time' as const,
            time: {
              unit: 'day' as const,
              displayFormats: { day: 'MMM dd' }
            },
            title: { display: true, text: 'Time' }
          },
          y: {
            min: 1,
            max: 5,
            ticks: {
              stepSize: 0.5,
              callback: function(value: any) {
                return `${value}/5`;
              }
            },
            title: { display: true, text: 'Average Rating' }
          }
        }
      }
    };

    chart = new Chart(chartCanvas, config);
  }

  $effect(() => {
    if (mounted && data) {
      createChart();
    }
  });

  onMount(() => {
    mounted = true;
    createChart();
  });

  onDestroy(() => {
    if (chart) {
      chart.destroy();
    }
  });
</script>

<div class="rating-chart bg-white rounded-lg border p-6">
  <div class="flex items-center gap-2 mb-4">
    <Star class="w-5 h-5 text-yellow-500" />
    <h3 class="text-lg font-semibold text-gray-900">Rating Analysis</h3>
  </div>

  {#if series.length === 0}
    <div class="text-center py-8 text-gray-500">
      <Star class="w-12 h-12 mx-auto mb-4 opacity-50" />
      <p>No rating data available</p>
    </div>
  {:else}
    <!-- Chart -->
    <div class="chart-container mb-6">
      <canvas bind:this={chartCanvas} class="w-full h-80"></canvas>
    </div>
    
    <!-- Statistics -->
    <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
      {#each series as seriesData, index}
        {#if seriesData.statistics}
          <div class="bg-yellow-50 border border-yellow-200 rounded-lg p-4">
            <div class="flex items-center justify-between mb-2">
              <h4 class="font-medium text-gray-900">{seriesData.metric_name}</h4>
              <div class="flex items-center gap-1">
                {#if seriesData.statistics.trend_direction}
                  {#snippet trendIcon()}
                    {@const TrendIcon = getTrendIcon(seriesData.statistics.trend_direction)}
                    <TrendIcon class="w-4 h-4 {getTrendColor(seriesData.statistics.trend_direction)}" />
                  {/snippet}
                  {@render trendIcon()}
                {/if}
              </div>
            </div>
            
            <div class="space-y-2 text-sm">
              <div class="flex justify-between">
                <span class="text-gray-600">Average:</span>
                <span class="font-semibold">{formatRatingValue(seriesData.statistics.average)}</span>
              </div>
              <div class="flex justify-between">
                <span class="text-gray-600">Highest:</span>
                <span class="font-semibold">{formatRatingValue(seriesData.statistics.max)}</span>
              </div>
              <div class="flex justify-between">
                <span class="text-gray-600">Lowest:</span>
                <span class="font-semibold">{formatRatingValue(seriesData.statistics.min)}</span>
              </div>
              {#if seriesData.statistics.trend_strength > 0}
                <div class="text-xs text-gray-500 pt-1 border-t">
                  Trend strength: {(seriesData.statistics.trend_strength * 100).toFixed(1)}%
                </div>
              {/if}
            </div>
          </div>
        {/if}
      {/each}
    </div>
  {/if}
</div>

<style>
  .chart-container {
    position: relative;
    height: 320px;
    width: 100%;
  }
</style>