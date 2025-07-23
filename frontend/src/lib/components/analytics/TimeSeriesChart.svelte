<script lang="ts">
  import { onMount, onDestroy } from 'svelte';
  import { Chart, registerables } from 'chart.js';
  import 'chartjs-adapter-date-fns';
  import { TrendingUp, TrendingDown, Minus } from 'lucide-svelte';

  interface Props {
    data: any;
  }

  let { data }: Props = $props();
  
  let chartCanvas: HTMLCanvasElement;
  let chart: Chart | null = null;
  let mounted = $state(false);

  // Chart customization options
  let chartType = $state<'line' | 'bar'>('line');
  let showFill = $state(false);
  let showDataPoints = $state(true);
  let lineSmooth = $state(true);
  let gridLines = $state(true);
  let showLegend = $state(true);
  let showTitle = $state(true);
  let animateChart = $state(true);
  let lineThickness = $state(3);
  let pointSize = $state(6);
  let yAxisScale = $state<'linear' | 'logarithmic'>('linear');
  let showTooltips = $state(true);
  let chartTheme = $state<'default' | 'dark' | 'colorful'>('default');

  // Register Chart.js components
  Chart.register(...registerables);

  const colorThemes = {
    default: [
      '#3B82F6', // blue
      '#10B981', // green
      '#F59E0B', // amber
      '#EF4444', // red
      '#8B5CF6', // purple
      '#F97316', // orange
      '#06B6D4', // cyan
      '#84CC16'  // lime
    ],
    dark: [
      '#60A5FA', // lighter blue
      '#34D399', // lighter green
      '#FBBF24', // lighter amber
      '#F87171', // lighter red
      '#A78BFA', // lighter purple
      '#FB923C', // lighter orange
      '#22D3EE', // lighter cyan
      '#A3E635'  // lighter lime
    ],
    colorful: [
      '#FF6B6B', // bright red
      '#4ECDC4', // teal
      '#45B7D1', // bright blue
      '#96CEB4', // mint
      '#FFEAA7', // light yellow
      '#DDA0DD', // plum
      '#98D8C8', // mint green
      '#F7DC6F'  // light gold
    ]
  };

  let colors = $derived(colorThemes[chartTheme]);

  let series = $derived(data?.series || []);

  function formatValue(value: number, metricType: string, seriesData?: any): string {
    // Check if this is a question metric
    if (metricType.startsWith('question_')) {
      // Try to get question type from metadata
      let questionType = seriesData?.metadata?.question_type;
      
      // If metadata is a string, try to parse it as JSON
      if (typeof seriesData?.metadata === 'string') {
        try {
          const parsed = JSON.parse(seriesData.metadata);
          questionType = parsed.question_type;
        } catch (e) {
          console.warn('Failed to parse metadata:', seriesData.metadata);
        }
      }
      
      switch (questionType) {
        case 'rating':
          return value.toFixed(1) + '/5';
        case 'scale':
          return value.toFixed(1) + '/10';
        case 'yes_no':
          return value.toFixed(1) + '% Yes';
        case 'text':
          if (value >= -1 && value <= 1) {
            const sentiment = value > 0.1 ? 'Positive' : value < -0.1 ? 'Negative' : 'Neutral';
            return `${sentiment} (${value.toFixed(2)})`;
          }
          return Math.round(value).toLocaleString() + ' responses';
        case 'single_choice':
        case 'multiple_choice':
          return Math.round(value).toLocaleString() + ' responses';
        default:
          return value.toFixed(1);
      }
    }
    
    // Handle non-question metrics
    switch (metricType) {
      case 'survey_responses':
        return Math.round(value).toLocaleString() + ' responses';
      default:
        return Math.round(value).toLocaleString();
    }
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

  function createChart() {
    if (!mounted || !chartCanvas || !data?.series) return;
    
    // Destroy existing chart
    if (chart) {
      chart.destroy();
    }

    console.log('Creating Chart.js chart with data:', data);

    // Prepare datasets
    const datasets = data.series.map((seriesData: any, index: number) => {
      console.log(`Series ${index} (${seriesData.metric_name}):`, seriesData.points);
      
      const points = (seriesData.points || []).map((point: any) => {
        const timestamp = new Date(point.timestamp);
        console.log('Point:', { 
          originalTimestamp: point.timestamp, 
          parsedDate: timestamp, 
          value: point.value 
        });
        
        return {
          x: timestamp,
          y: point.value
        };
      });

      return {
        label: seriesData.metric_name,
        data: points,
        borderColor: colors[index % colors.length],
        backgroundColor: colors[index % colors.length] + (showFill ? '20' : '00'),
        borderWidth: lineThickness,
        pointRadius: showDataPoints ? pointSize : 0,
        pointHoverRadius: showDataPoints ? pointSize + 2 : 4,
        fill: showFill,
        tension: lineSmooth ? 0.4 : 0
      };
    });

    const config = {
      type: chartType,
      data: {
        datasets: datasets
      },
      options: {
        responsive: true,
        maintainAspectRatio: false,
        animation: animateChart ? {
          duration: 1000,
          easing: 'easeInOutQuart'
        } : false,
        plugins: {
          title: {
            display: showTitle,
            text: 'Survey Response Analysis Over Time',
            font: {
              size: 16,
              weight: 'bold' as const
            },
            color: '#1f2937',
            padding: 20
          },
          legend: {
            display: showLegend,
            position: 'top' as const,
            align: 'center' as const,
            labels: {
              usePointStyle: true,
              padding: 20,
              font: {
                size: 12
              }
            }
          },
          tooltip: {
            enabled: showTooltips,
            mode: 'point' as const,
            intersect: false,
            backgroundColor: 'rgba(0, 0, 0, 0.8)',
            titleColor: 'white',
            bodyColor: 'white',
            borderColor: 'rgba(255, 255, 255, 0.1)',
            borderWidth: 1,
            callbacks: {
              title: function(context: any) {
                return new Date(context[0].parsed.x).toLocaleDateString();
              },
              label: function(context: any) {
                const datasetLabel = context.dataset.label || '';
                const seriesData = data.series.find((s: any) => s.metric_name === datasetLabel);
                const formattedValue = seriesData 
                  ? formatValue(context.parsed.y, seriesData.metric_type, seriesData)
                  : context.parsed.y.toLocaleString();
                return `${datasetLabel}: ${formattedValue}`;
              }
            }
          }
        },
        scales: {
          x: {
            type: 'time' as const,
            time: {
              unit: 'day' as const,
              displayFormats: {
                day: 'MMM dd'
              }
            },
            title: {
              display: true,
              text: 'Time',
              font: {
                size: 12,
                weight: 'bold' as const
              },
              color: '#374151'
            },
            grid: {
              display: gridLines,
              color: '#e5e7eb',
              borderColor: '#d1d5db'
            }
          },
          y: {
            type: yAxisScale,
            beginAtZero: false,
            title: {
              display: true,
              text: 'Value',
              font: {
                size: 12,
                weight: 'bold' as const
              },
              color: '#374151'
            },
            grid: {
              display: gridLines,
              color: '#e5e7eb',
              borderColor: '#d1d5db'
            }
          }
        },
        interaction: {
          mode: 'nearest' as const,
          axis: 'x' as const,
          intersect: false
        }
      }
    };

    chart = new Chart(chartCanvas, config);
  }

  // Re-render when data changes
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

<div class="time-series-chart">
  <div class="mb-6">
    <h3 class="text-lg font-semibold text-gray-900 mb-4">Survey Response Trends</h3>
    
    {#if series.length === 0}
      <div class="text-center py-8 text-gray-500">
        <p>No data available for the selected time period</p>
      </div>
    {:else}
      <!-- Chart Controls -->
      <div class="bg-gray-50 rounded-lg p-4 mb-4">
        <h4 class="text-sm font-medium text-gray-700 mb-3">Chart Customization</h4>
        
        <!-- First Row: Basic Settings -->
        <div class="grid grid-cols-2 md:grid-cols-4 gap-4 mb-4">
          <!-- Chart Type -->
          <div>
            <label class="text-xs text-gray-600 mb-1 block">Chart Type</label>
            <select bind:value={chartType} class="w-full text-sm border border-gray-300 rounded px-2 py-1">
              <option value="line">Line Chart</option>
              <option value="bar">Bar Chart</option>
            </select>
          </div>

          <!-- Theme -->
          <div>
            <label class="text-xs text-gray-600 mb-1 block">Color Theme</label>
            <select bind:value={chartTheme} class="w-full text-sm border border-gray-300 rounded px-2 py-1">
              <option value="default">Default</option>
              <option value="dark">Dark</option>
              <option value="colorful">Colorful</option>
            </select>
          </div>

          <!-- Y-Axis Scale -->
          <div>
            <label class="text-xs text-gray-600 mb-1 block">Y-Axis Scale</label>
            <select bind:value={yAxisScale} class="w-full text-sm border border-gray-300 rounded px-2 py-1">
              <option value="linear">Linear</option>
              <option value="logarithmic">Logarithmic</option>
            </select>
          </div>

          <!-- Line Thickness -->
          <div>
            <label class="text-xs text-gray-600 mb-1 block">Line Thickness: {lineThickness}px</label>
            <input 
              type="range" 
              bind:value={lineThickness} 
              min="1" 
              max="8" 
              step="1"
              class="w-full text-sm"
            />
          </div>
        </div>

        <!-- Second Row: Visual Settings -->
        <div class="grid grid-cols-2 md:grid-cols-4 gap-4 mb-4">
          <!-- Point Size -->
          <div>
            <label class="text-xs text-gray-600 mb-1 block">Point Size: {pointSize}px</label>
            <input 
              type="range" 
              bind:value={pointSize} 
              min="2" 
              max="12" 
              step="1"
              class="w-full text-sm"
            />
          </div>

          <!-- Show Fill -->
          <div class="flex items-center pt-4">
            <input 
              type="checkbox" 
              bind:checked={showFill} 
              id="showFill"
              class="mr-2"
            />
            <label for="showFill" class="text-xs text-gray-600">Fill Area</label>
          </div>
          
          <!-- Show Data Points -->
          <div class="flex items-center pt-4">
            <input 
              type="checkbox" 
              bind:checked={showDataPoints} 
              id="showDataPoints"
              class="mr-2"
            />
            <label for="showDataPoints" class="text-xs text-gray-600">Data Points</label>
          </div>
          
          <!-- Smooth Lines -->
          <div class="flex items-center pt-4">
            <input 
              type="checkbox" 
              bind:checked={lineSmooth} 
              id="lineSmooth"
              class="mr-2"
            />
            <label for="lineSmooth" class="text-xs text-gray-600">Smooth Lines</label>
          </div>
        </div>

        <!-- Third Row: Display Settings -->
        <div class="grid grid-cols-2 md:grid-cols-5 gap-4">
          <!-- Grid Lines -->
          <div class="flex items-center">
            <input 
              type="checkbox" 
              bind:checked={gridLines} 
              id="gridLines"
              class="mr-2"
            />
            <label for="gridLines" class="text-xs text-gray-600">Grid Lines</label>
          </div>

          <!-- Show Legend -->
          <div class="flex items-center">
            <input 
              type="checkbox" 
              bind:checked={showLegend} 
              id="showLegend"
              class="mr-2"
            />
            <label for="showLegend" class="text-xs text-gray-600">Legend</label>
          </div>

          <!-- Show Title -->
          <div class="flex items-center">
            <input 
              type="checkbox" 
              bind:checked={showTitle} 
              id="showTitle"
              class="mr-2"
            />
            <label for="showTitle" class="text-xs text-gray-600">Title</label>
          </div>

          <!-- Show Tooltips -->
          <div class="flex items-center">
            <input 
              type="checkbox" 
              bind:checked={showTooltips} 
              id="showTooltips"
              class="mr-2"
            />
            <label for="showTooltips" class="text-xs text-gray-600">Tooltips</label>
          </div>

          <!-- Animate Chart -->
          <div class="flex items-center">
            <input 
              type="checkbox" 
              bind:checked={animateChart} 
              id="animateChart"
              class="mr-2"
            />
            <label for="animateChart" class="text-xs text-gray-600">Animations</label>
          </div>
        </div>
      </div>

      <!-- Chart Container -->
      <div class="bg-white rounded-lg p-6 mb-6">
        <div class="chart-container">
          <canvas bind:this={chartCanvas} class="w-full h-96"></canvas>
        </div>
      </div>
      
      <!-- Statistics Cards -->
      <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4 mt-6">
        {#each series as seriesData, index}
          {#if seriesData.statistics}
            <div class="bg-gray-50 rounded-lg p-4">
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
              
              <div class="grid grid-cols-2 gap-2 text-sm">
                <div>
                  <span class="text-gray-500">Average:</span>
                  <span class="font-semibold ml-1">
                    {formatValue(seriesData.statistics.average, seriesData.metric_type, seriesData)}
                  </span>
                </div>
                <div>
                  <span class="text-gray-500">Total:</span>
                  <span class="font-semibold ml-1">
                    {formatValue(seriesData.statistics.total, seriesData.metric_type, seriesData)}
                  </span>
                </div>
                <div>
                  <span class="text-gray-500">Min:</span>
                  <span class="font-semibold ml-1">
                    {formatValue(seriesData.statistics.min, seriesData.metric_type, seriesData)}
                  </span>
                </div>
                <div>
                  <span class="text-gray-500">Max:</span>
                  <span class="font-semibold ml-1">
                    {formatValue(seriesData.statistics.max, seriesData.metric_type, seriesData)}
                  </span>
                </div>
              </div>
              
              {#if seriesData.statistics.trend_strength > 0}
                <div class="mt-2 text-xs text-gray-500">
                  Trend strength: {(seriesData.statistics.trend_strength * 100).toFixed(1)}%
                </div>
              {/if}
            </div>
          {/if}
        {/each}
      </div>
    {/if}
  </div>
</div>

<style>
  .chart-container {
    position: relative;
    height: 400px;
    width: 100%;
  }
  
  .chart-container canvas {
    max-width: 100%;
    height: 100% !important;
  }
</style>