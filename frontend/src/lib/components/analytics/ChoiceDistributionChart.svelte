<script lang="ts">
  import { onMount, onDestroy } from 'svelte';
  import { Chart, registerables } from 'chart.js';
  import 'chartjs-adapter-date-fns';
  import {
    List,
    TrendingUp,
    TrendingDown,
    Minus,
    ZoomIn,
    ZoomOut,
    RotateCcw,
    Eye,
    EyeOff,
  } from 'lucide-svelte';

  interface Props {
    data: any;
  }

  let { data }: Props = $props();

  let chartCanvas: HTMLCanvasElement;
  let chart: Chart | null = null;
  let mounted = $state(false);

  let showDataPoints = $state(true);
  let lineSmooth = $state(true);
  let gridLines = $state(true);
  let showTooltips = $state(true);
  let lineThickness = $state(3);
  let pointSize = $state(6);
  let enableZoom = $state(true);
  let enablePan = $state(true);
  let showZoomControls = $state(true);

  let zoomPlugin: any = null;

  Chart.register(...registerables);

  const choiceColors = [
    '#3B82F6',
    '#10B981',
    '#F59E0B',
    '#EF4444',
    '#8B5CF6',
    '#F97316',
    '#06B6D4',
    '#84CC16',
    '#EC4899',
    '#6366F1',
    '#14B8A6',
    '#F43F5E',
  ];

  let series = $derived(data?.choice_series || data?.series || []);

  function formatChoiceValue(value: number): string {
    return `${Math.round(value).toLocaleString()} selections`;
  }

  function getTrendIcon(trendDirection: string) {
    switch (trendDirection) {
      case 'improving':
        return TrendingUp;
      case 'declining':
        return TrendingDown;
      default:
        return Minus;
    }
  }

  function getTrendColor(trendDirection: string): string {
    switch (trendDirection) {
      case 'improving':
        return 'text-green-600';
      case 'declining':
        return 'text-red-600';
      default:
        return 'text-gray-600';
    }
  }

  function resetZoom() {
    if (chart) {
      chart.resetZoom();
    }
  }

  function zoomIn() {
    if (chart) {
      chart.zoom(1.1);
    }
  }

  function zoomOut() {
    if (chart) {
      chart.zoom(0.9);
    }
  }

  function toggleSeries(seriesLabel: string) {
    const datasetIndex = chart?.data.datasets.findIndex(
      d => d.label === seriesLabel
    );
    if (chart && datasetIndex !== undefined && datasetIndex >= 0) {
      const meta = chart.getDatasetMeta(datasetIndex);
      meta.hidden = !meta.hidden;
      chart.update();
    }
  }

  function getChoiceLabel(seriesData: any): string {
    if (seriesData.metadata) {
      const metadata = seriesData.metadata;
      if (metadata?.choice_option) {
        return metadata.choice_option;
      }
    }

    const parts = seriesData.metric_name?.split(': ');
    if (parts && parts.length > 1) {
      return parts[parts.length - 1];
    }

    return seriesData.metric_name || 'Unknown Option';
  }

  function createChart() {
    if (!mounted || !chartCanvas || (!data?.series && !data?.choice_series))
      return;

    if (chart) {
      chart.destroy();
    }

    const dataSource = data.choice_series || data.series;

    const datasets = dataSource.map((seriesData: any, index: number) => {
      const points = (seriesData.points || []).map((point: any) => ({
        x: new Date(point.timestamp),
        y: point.value || point.count,
      }));

      const color = choiceColors[index % choiceColors.length];
      const choiceLabel = seriesData.choice || getChoiceLabel(seriesData);

      return {
        label: choiceLabel,
        data: points,
        borderColor: color,
        backgroundColor: color + '20',
        borderWidth: lineThickness,
        pointRadius: showDataPoints ? pointSize : 0,
        pointHoverRadius: showDataPoints ? pointSize + 2 : 4,
        fill: false,
        tension: lineSmooth ? 0.4 : 0,
      };
    });

    const config = {
      type: 'line' as const,
      data: { datasets },
      options: {
        responsive: true,
        maintainAspectRatio: false,
        animation: {
          duration: 1000,
          easing: 'easeInOutQuart',
        },
        plugins: {
          title: {
            display: true,
            text: 'Choice Distribution Over Time',
            font: { size: 16, weight: 'bold' as const },
            color: '#1f2937',
          },
          legend: {
            display: false,
          },
          tooltip: {
            enabled: showTooltips,
            mode: 'nearest' as const,
            intersect: false,
            backgroundColor: 'rgba(0, 0, 0, 0.9)',
            titleColor: 'white',
            bodyColor: 'white',
            borderColor: 'rgba(255, 255, 255, 0.2)',
            borderWidth: 1,
            cornerRadius: 8,
            padding: 12,
            displayColors: true,
            callbacks: {
              title: (context: any) => {
                const date = new Date(context[0].parsed.x);
                return date.toLocaleDateString('en-US', {
                  weekday: 'short',
                  year: 'numeric',
                  month: 'short',
                  day: 'numeric',
                });
              },
              label: (context: any) =>
                `${context.dataset.label}: ${formatChoiceValue(context.parsed.y)}`,
            },
          },
          zoom:
            enableZoom && zoomPlugin
              ? {
                  limits: {
                    x: { min: 'original', max: 'original' },
                    y: { min: 'original', max: 'original' },
                  },
                  pan: {
                    enabled: enablePan,
                    mode: 'xy' as const,
                  },
                  zoom: {
                    wheel: {
                      enabled: true,
                    },
                    pinch: {
                      enabled: true,
                    },
                    mode: 'xy' as const,
                  },
                }
              : {},
        },
        scales: {
          x: {
            type: 'time' as const,
            time: {
              unit: 'day' as const,
              displayFormats: {
                day: 'MMM dd',
              },
            },
            title: {
              display: true,
              text: 'Date',
              font: { size: 12, weight: 'bold' as const },
              color: '#374151',
            },
            grid: {
              display: gridLines,
              color: '#e5e7eb',
            },
          },
          y: {
            beginAtZero: true,
            ticks: {
              callback: function (value: any) {
                return Math.round(value).toLocaleString();
              },
            },
            title: {
              display: true,
              text: 'Number of Selections',
              font: { size: 12, weight: 'bold' as const },
              color: '#374151',
            },
            grid: {
              display: gridLines,
              color: '#e5e7eb',
            },
          },
        },
        interaction: {
          mode: 'nearest' as const,
          axis: 'x' as const,
          intersect: false,
        },
      },
    };

    chart = new Chart(chartCanvas, config);
  }

  $effect(() => {
    if (mounted && data) {
      createChart();
    }
  });

  onMount(async () => {
    if (typeof window !== 'undefined') {
      try {
        const zoomModule = await import('chartjs-plugin-zoom');
        zoomPlugin = zoomModule.default;
        Chart.register(zoomPlugin);
      } catch (e) {
        console.warn('Chart zoom plugin not available');
      }
    }

    mounted = true;
    createChart();
  });

  onDestroy(() => {
    if (chart) {
      chart.destroy();
    }
  });
</script>

<div class="choice-distribution-chart">
  {#if series.length === 0}
    <div class="text-center py-8 text-gray-500">
      <List class="w-12 h-12 mx-auto mb-4 opacity-50" />
      <p>No choice distribution data available</p>
    </div>
  {:else}
    
    <div class="bg-white border border-gray-200 rounded-xl p-4 mb-6 shadow-sm">
      <div class="flex flex-wrap items-center justify-between gap-4">
        
        {#if showZoomControls}
          <div class="flex items-center gap-1">
            <span class="text-sm font-medium text-gray-700 mr-2">Zoom:</span>
            <button
              class="p-2 text-gray-600 hover:text-blue-600 hover:bg-blue-50 rounded-lg transition-all"
              onclick={zoomIn}
              title="Zoom In">
              <ZoomIn class="w-4 h-4" />
            </button>
            <button
              class="p-2 text-gray-600 hover:text-blue-600 hover:bg-blue-50 rounded-lg transition-all"
              onclick={zoomOut}
              title="Zoom Out">
              <ZoomOut class="w-4 h-4" />
            </button>
            <button
              class="p-2 text-gray-600 hover:text-blue-600 hover:bg-blue-50 rounded-lg transition-all"
              onclick={resetZoom}
              title="Reset Zoom">
              <RotateCcw class="w-4 h-4" />
            </button>
          </div>
        {/if}

        
        <div class="flex items-center gap-4 text-sm">
          <label class="flex items-center gap-2 cursor-pointer">
            <input
              type="checkbox"
              bind:checked={showDataPoints}
              class="rounded border-gray-300 text-blue-600 focus:ring-blue-500 focus:ring-offset-0" />
            <span class="text-gray-700">Show Points</span>
          </label>
          <label class="flex items-center gap-2 cursor-pointer">
            <input
              type="checkbox"
              bind:checked={lineSmooth}
              class="rounded border-gray-300 text-blue-600 focus:ring-blue-500 focus:ring-offset-0" />
            <span class="text-gray-700">Smooth Lines</span>
          </label>
        </div>
      </div>
    </div>

    
    {#if series.length > 0}
      <div
        class="bg-white border border-gray-200 rounded-xl p-4 mb-6 shadow-sm">
        <h4 class="text-sm font-semibold text-gray-900 mb-3">Choice Options</h4>
        <div class="flex flex-wrap gap-2">
          {#each series as seriesData, index}
            {@const isVisible = !chart?.getDatasetMeta(index)?.hidden}
            {@const choiceLabel =
              seriesData.choice || getChoiceLabel(seriesData)}
            <button
              class="flex items-center gap-2 px-3 py-2 rounded-lg border transition-all {isVisible
                ? 'border-gray-300 bg-white hover:bg-gray-50'
                : 'border-gray-200 bg-gray-100 text-gray-500'}"
              onclick={() => toggleSeries(choiceLabel)}>
              <div
                class="w-3 h-3 rounded-full border-2"
                style="background-color: {isVisible
                  ? choiceColors[index % choiceColors.length]
                  : 'transparent'}; border-color: {choiceColors[
                  index % choiceColors.length
                ]}">
              </div>
              <span class="text-sm font-medium">{choiceLabel}</span>
              {#if isVisible}
                <Eye class="w-3 h-3" />
              {:else}
                <EyeOff class="w-3 h-3" />
              {/if}
            </button>
          {/each}
        </div>
      </div>
    {/if}

    
    <div class="bg-white rounded-lg p-6 mb-6 shadow-sm">
      <div class="chart-container">
        <canvas bind:this={chartCanvas} class="w-full h-96"></canvas>
      </div>

      
      <div class="mt-4 pt-4 border-t border-gray-100">
        <div class="flex flex-wrap gap-4 text-xs text-gray-500">
          <span class="flex items-center gap-1">
            <ZoomIn class="w-3 h-3" />
            Mouse wheel to zoom
          </span>
        </div>
      </div>
    </div>

    
    <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4 mt-6">
      {#each series as seriesData, index}
        {#if seriesData.statistics}
          {@const choiceLabel = seriesData.choice || getChoiceLabel(seriesData)}
          <div class="bg-gray-50 rounded-lg p-4">
            <div class="flex items-center justify-between mb-2">
              <h4 class="font-medium text-gray-900">{choiceLabel}</h4>
              <div class="flex items-center gap-1">
                {#if seriesData.statistics.trend_direction}
                  {#snippet trendIcon()}
                    {@const TrendIcon = getTrendIcon(
                      seriesData.statistics.trend_direction
                    )}
                    <TrendIcon
                      class="w-4 h-4 {getTrendColor(
                        seriesData.statistics.trend_direction
                      )}" />
                  {/snippet}
                  {@render trendIcon()}
                {/if}
              </div>
            </div>

            <div class="grid grid-cols-2 gap-2 text-sm">
              <div>
                <span class="text-gray-500">Average:</span>
                <span class="font-semibold ml-1">
                  {formatChoiceValue(seriesData.statistics.average)}
                </span>
              </div>
              <div>
                <span class="text-gray-500">Total:</span>
                <span class="font-semibold ml-1">
                  {formatChoiceValue(seriesData.statistics.total)}
                </span>
              </div>
              <div>
                <span class="text-gray-500">Min:</span>
                <span class="font-semibold ml-1">
                  {formatChoiceValue(seriesData.statistics.min)}
                </span>
              </div>
              <div>
                <span class="text-gray-500">Max:</span>
                <span class="font-semibold ml-1">
                  {formatChoiceValue(seriesData.statistics.max)}
                </span>
              </div>
            </div>

            {#if seriesData.statistics.trend_strength > 0}
              <div class="mt-2 text-xs text-gray-500">
                Trend strength: {(
                  seriesData.statistics.trend_strength * 100
                ).toFixed(1)}%
              </div>
            {/if}
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
