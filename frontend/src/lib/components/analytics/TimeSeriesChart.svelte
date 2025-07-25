<script lang="ts">
  import { onMount, onDestroy } from 'svelte';
  import { Chart, registerables } from 'chart.js';
  import 'chartjs-adapter-date-fns';
  import { TrendingUp, TrendingDown, Minus, ZoomIn, ZoomOut, Move, RotateCcw, Eye, EyeOff } from 'lucide-svelte';

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

  // Interactive features
  let enableZoom = $state(true);
  let enablePan = $state(true);
  let showZoomControls = $state(true);
  let selectedDataSeries = $state<string[]>([]);

  
  // Zoom plugin will be loaded dynamically
  let zoomPlugin: any = null;

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

  function cleanMetricName(metricName: string): string {
    // Remove UUID patterns (8-4-4-4-12 hex characters)
    let cleaned = metricName.replace(/[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}/gi, '');
    
    // Remove "product_" prefix
    cleaned = cleaned.replace(/product_/gi, '');
    
    // Remove double underscores or dashes
    cleaned = cleaned.replace(/[_-]{2,}/g, '_');
    
    // Remove leading/trailing underscores or dashes
    cleaned = cleaned.replace(/^[_-]+|[_-]+$/g, '');
    
    // Convert underscores to spaces and capitalize words
    cleaned = cleaned.replace(/_/g, ' ')
      .split(' ')
      .map(word => word.charAt(0).toUpperCase() + word.slice(1).toLowerCase())
      .join(' ');
    
    // Special cases for known metric types
    if (cleaned.toLowerCase().includes('rating questions')) return 'Rating Questions';
    if (cleaned.toLowerCase().includes('scale questions')) return 'Scale Questions';
    if (cleaned.toLowerCase().includes('yes no questions')) return 'Yes/No Questions';
    if (cleaned.toLowerCase().includes('text questions')) return 'Text Questions';
    if (cleaned.toLowerCase().includes('single choice questions')) return 'Single Choice Questions';
    if (cleaned.toLowerCase().includes('multiple choice questions')) return 'Multiple Choice Questions';
    if (cleaned.toLowerCase().includes('survey responses')) return 'Survey Responses';
    
    // If it's a question ID, try to extract just the question part
    if (cleaned.toLowerCase().startsWith('question ')) {
      return 'Question Response';
    }
    
    return cleaned || 'Metric';
  }

  function formatValue(value: number, metricType: string, seriesData?: any): string {
    // Check if this is a question metric
    if (metricType.startsWith('question_')) {
      // Try to get question type and labels from metadata
      let metadata = seriesData?.metadata;
      let questionType = metadata?.question_type;
      let minLabel = metadata?.min_label;
      let maxLabel = metadata?.max_label;
      let minValue = metadata?.min_value;
      let maxValue = metadata?.max_value;
      
      // If metadata is a string, try to parse it as JSON
      if (typeof seriesData?.metadata === 'string') {
        try {
          const parsed = JSON.parse(seriesData.metadata);
          questionType = parsed.question_type;
          minLabel = parsed.min_label;
          maxLabel = parsed.max_label;
          minValue = parsed.min_value;
          maxValue = parsed.max_value;
        } catch (e) {
          // ignore parsing errors
        }
      }
      
      switch (questionType) {
        case 'rating':
          const ratingMax = maxValue || 5;
          const ratingMaxLabel = maxLabel || ratingMax.toString();
          return `${value.toFixed(1)}/${ratingMax} (${ratingMaxLabel})`;
          
        case 'scale':
          return value.toFixed(1);
          
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

  // Interactive functions
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
    const datasetIndex = chart?.data.datasets.findIndex(d => d.label === seriesLabel);
    if (chart && datasetIndex !== undefined && datasetIndex >= 0) {
      const meta = chart.getDatasetMeta(datasetIndex);
      meta.hidden = !meta.hidden;
      chart.update();
    }
  }



  function getYAxisConfig(seriesData: any[]) {
    // Analyze the data to determine appropriate Y-axis configuration
    for (const series of seriesData) {
      const metricType = series.metric_type;
      
      // Check if this is a question metric
      if (metricType.startsWith('question_')) {
        let metadata = series.metadata;
        let questionType = metadata?.question_type;
        let minLabel = metadata?.min_label;
        let maxLabel = metadata?.max_label;
        let minValue = metadata?.min_value;
        let maxValue = metadata?.max_value;
        
        if (typeof series.metadata === 'string') {
          try {
            const parsed = JSON.parse(series.metadata);
            questionType = parsed.question_type;
            minLabel = parsed.min_label;
            maxLabel = parsed.max_label;
            minValue = parsed.min_value;
            maxValue = parsed.max_value;
          } catch (e) {
            // ignore parsing errors
          }
        }
        
        // If we have explicit question type, use it
        switch (questionType) {
          case 'rating':
            const ratingMin = minValue || 1;
            const ratingMax = maxValue || 5;
            const ratingMinLabel = minLabel || ratingMin.toString();
            const ratingMaxLabel = maxLabel || ratingMax.toString();
            
            return {
              min: ratingMin,
              max: ratingMax,
              stepSize: (ratingMax - ratingMin) / 8,
              label: `Rating (${ratingMinLabel} ← → ${ratingMaxLabel})`,
              formatter: (value: any) => `${ratingMinLabel} ← ${value.toFixed(1)} → ${ratingMaxLabel}`
            };
            
          case 'scale':
            const scaleMin = minValue || 1;
            const scaleMax = maxValue || 10;
            const scaleMinLabel = minLabel || scaleMin.toString();
            const scaleMaxLabel = maxLabel || scaleMax.toString();
            
            return {
              min: scaleMin,
              max: scaleMax,
              stepSize: 1,
              label: `Scale (${scaleMinLabel} ← → ${scaleMaxLabel})`,
              formatter: (value: any) => value.toFixed(1)
            };
            
          case 'yes_no':
            return {
              min: 0,
              max: 100,
              stepSize: 10,
              label: 'Yes Response Rate (%)',
              formatter: (value: any) => `${value.toFixed(1)}%`
            };
            
          case 'text':
            // Check if values are sentiment scores (-1 to 1) or counts
            const hasNegativeValues = series.points?.some((p: any) => p.value < 0);
            if (hasNegativeValues) {
              return {
                min: -1,
                max: 1,
                stepSize: 0.2,
                label: 'Sentiment Score (-1 to +1)',
                formatter: (value: any) => {
                  if (value > 0.1) return `+${value.toFixed(1)} (Positive)`;
                  if (value < -0.1) return `${value.toFixed(1)} (Negative)`;
                  return `${value.toFixed(1)} (Neutral)`;
                }
              };
            }
            break;
        }
        
        // If no metadata available, try to infer from question text and values
        const questionText = series.metric_name?.toLowerCase() || '';
        const sampleValue = series.points?.[0]?.value;
        
        // Infer question type based on patterns
        if (sampleValue !== undefined) {
          // If value is between 0-100 and looks like a percentage
          if (sampleValue >= 0 && sampleValue <= 100 && (sampleValue % 1 !== 0)) {
            // Check if question suggests yes/no
            if (questionText.includes('is ') || questionText.includes('do ') || 
                questionText.includes('can ') || questionText.includes('will ')) {
              return {
                min: 0,
                max: 100,
                stepSize: 10,
                label: 'Yes Response Rate (%)',
                formatter: (value: any) => `${value.toFixed(1)}% Yes`
              };
            }
          }
          
          // If value is between 1-5, likely a rating
          if (sampleValue >= 1 && sampleValue <= 5) {
            return {
              min: 1,
              max: 5,
              stepSize: 0.5,
              label: 'Average Rating (1-5)',
              formatter: (value: any) => `${value.toFixed(1)}/5`
            };
          }
          
          // If value is between 1-10, likely a scale
          if (sampleValue >= 1 && sampleValue <= 10) {
            // Try to infer scale meaning from question text
            let scaleLabel = 'Average Scale';
            let lowLabel = '1';
            let highLabel = '10';
            
            if (questionText.includes('spic') || questionText.includes('hot')) {
              scaleLabel = 'Spiciness Level';
              lowLabel = 'Mild';
              highLabel = 'Very Spicy';
            } else if (questionText.includes('satisf') || questionText.includes('happy')) {
              scaleLabel = 'Satisfaction Level'; 
              lowLabel = 'Poor';
              highLabel = 'Excellent';
            } else if (questionText.includes('quality') || questionText.includes('good')) {
              scaleLabel = 'Quality Rating';
              lowLabel = 'Poor';
              highLabel = 'Excellent';
            } else if (questionText.includes('recommend')) {
              scaleLabel = 'Recommendation Score';
              lowLabel = 'Never';
              highLabel = 'Definitely';
            } else if (questionText.includes('difficult') || questionText.includes('easy')) {
              scaleLabel = 'Difficulty Level';
              lowLabel = 'Very Easy';
              highLabel = 'Very Hard';
            }
            
            return {
              min: 1,
              max: 10,
              stepSize: 1,
              label: `${scaleLabel} (${lowLabel} ← → ${highLabel})`,
              formatter: (value: any) => value.toFixed(1)
            };
          }
        }
        
        // Default for individual questions
        return {
          beginAtZero: true,
          label: 'Response Value',
          formatter: (value: any) => value.toFixed(1)
        };
      }
      
      // Check for percentage-based metrics
      if (metricType.includes('rate') || metricType.includes('completion')) {
        return {
          min: 0,
          max: 100,
          stepSize: 10,
          label: 'Percentage (%)',
          formatter: (value: any) => `${value.toFixed(1)}%`
        };
      }
    }
    
    // Default configuration for count-based metrics
    return {
      beginAtZero: true,
      label: 'Value',
      formatter: (value: any) => Math.round(value).toLocaleString()
    };
  }

  function createChart() {
    if (!mounted || !chartCanvas || !data?.series) return;
    
    // Destroy existing chart
    if (chart) {
      chart.destroy();
    }

    // Get appropriate Y-axis configuration
    const yAxisConfig = getYAxisConfig(data.series);
    console.log('Y-axis config:', yAxisConfig);
    console.log('Series data:', data.series?.map(s => ({ 
      metric_type: s.metric_type, 
      metric_name: s.metric_name,
      metadata: s.metadata,
      parsed_metadata: typeof s.metadata === 'string' ? JSON.parse(s.metadata) : s.metadata
    })));

    // Prepare datasets
    const datasets = data.series.map((seriesData: any, index: number) => {
      const points = (seriesData.points || []).map((point: any) => {
        const timestamp = new Date(point.timestamp);
        
        // Validate the timestamp is a valid date
        if (isNaN(timestamp.getTime())) {
          console.warn('Invalid timestamp:', point.timestamp);
          return null;
        }
        
        return {
          x: timestamp, // Use Date object directly
          y: point.value
        };
      }).filter(Boolean); // Remove any null entries

      // Use different colors based on question type for better visual distinction
      let colorHue = 220 + index * 40; // default blue range
      const metricType = seriesData.metric_type;
      
      if (metricType.includes('rating')) colorHue = 45; // orange for ratings
      else if (metricType.includes('scale')) colorHue = 260; // purple for scales
      else if (metricType.includes('yes_no')) colorHue = 140; // green for yes/no
      else if (metricType.includes('text')) colorHue = 200; // blue for sentiment
      else if (metricType.includes('choice')) colorHue = 20; // red for choices

      const color = `hsl(${colorHue + index * 20}, 70%, 50%)`;

      return {
        label: cleanMetricName(seriesData.metric_name),
        data: points,
        borderColor: color,
        backgroundColor: color + (showFill ? '20' : '00'),
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
            display: false
          },
          legend: {
            display: false, // We'll create custom interactive legend
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
              title: function(context: any) {
                const date = new Date(context[0].parsed.x);
                return date.toLocaleDateString('en-US', { 
                  weekday: 'short',
                  year: 'numeric', 
                  month: 'short', 
                  day: 'numeric',
                  hour: 'numeric',
                  minute: '2-digit'
                });
              },
              label: function(context: any) {
                const datasetLabel = context.dataset.label || '';
                const seriesData = data.series.find((s: any) => cleanMetricName(s.metric_name) === datasetLabel);
                const formattedValue = seriesData 
                  ? formatValue(context.parsed.y, seriesData.metric_type, seriesData)
                  : context.parsed.y.toLocaleString();
                return `${datasetLabel}: ${formattedValue}`;
              },
              afterBody: function(context: any) {
                if (context.length > 0) {
                  const seriesData = data.series.find((s: any) => cleanMetricName(s.metric_name) === context[0].dataset.label);
                  if (seriesData?.metadata) {
                    try {
                      const metadata = typeof seriesData.metadata === 'string' ? JSON.parse(seriesData.metadata) : seriesData.metadata;
                      if (metadata.question_type) {
                        return [`Type: ${metadata.question_type.replace('_', ' ')}`];
                      }
                    } catch (e) {
                      // ignore
                    }
                  }
                }
                return [];
              }
            }
          },
          zoom: enableZoom && zoomPlugin ? {
            limits: {
              x: {min: 'original', max: 'original'},
              y: {min: 'original', max: 'original'}
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
                enabled: true
              },
              mode: 'xy' as const,
              onZoomComplete: function({chart}: any) {
                // Optional: emit zoom event
                const customEvent = new CustomEvent('chartZoom', {
                  detail: { chart }
                });
                chartCanvas?.dispatchEvent(customEvent);
              }
            }
          } : {}
        },
        scales: {
          x: {
            type: 'time' as const,
            time: {
              unit: 'day' as const,
              displayFormats: {
                day: 'MMM dd'
              },
              tooltipFormat: 'MMM dd, yyyy'
            },
            // Force a time range when we have limited data points
            min: datasets.length > 0 && datasets[0].data.length === 1 
              ? new Date(datasets[0].data[0].x.getTime() - 24 * 60 * 60 * 1000) // 1 day before
              : undefined,
            max: datasets.length > 0 && datasets[0].data.length === 1
              ? new Date(datasets[0].data[0].x.getTime() + 24 * 60 * 60 * 1000) // 1 day after  
              : undefined,
            title: {
              display: true,
              text: 'Date',
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
            },
            ticks: {
              display: true,
              maxTicksLimit: 8,
              color: '#6B7280',
              font: {
                size: 11
              },
              callback: function(value: any, index: number, values: any[]) {
                return new Date(value).toLocaleDateString('en-US', { month: 'short', day: 'numeric' });
              }
            }
          },
          y: {
            type: yAxisScale,
            ...(yAxisConfig.min !== undefined && { min: yAxisConfig.min }),
            ...(yAxisConfig.max !== undefined && { max: yAxisConfig.max }),
            ...(yAxisConfig.beginAtZero !== undefined && { beginAtZero: yAxisConfig.beginAtZero }),
            ticks: {
              ...(yAxisConfig.stepSize && { stepSize: yAxisConfig.stepSize }),
              callback: function(value: any, index: number, values: any[]) {
                return yAxisConfig.formatter ? yAxisConfig.formatter(value) : value;
              },
              color: '#6B7280',
              font: {
                size: 11
              }
            },
            title: {
              display: true,
              text: yAxisConfig.label || 'Value',
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

  onMount(async () => {
    // Dynamically import zoom plugin on client side only
    if (typeof window !== 'undefined') {
      const zoomModule = await import('chartjs-plugin-zoom');
      zoomPlugin = zoomModule.default;
      Chart.register(zoomPlugin);
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

<div class="time-series-chart">
  <div class="mb-6">
    
    {#if series.length === 0}
      <div class="text-center py-8 text-gray-500">
        <p>No data available for the selected time period</p>
      </div>
    {:else}
      <!-- Enhanced Chart Controls -->
      <div class="bg-white border border-gray-200 rounded-xl p-4 mb-6 shadow-sm">
        <div class="flex flex-wrap items-center justify-between gap-4">
          <!-- Chart Type Controls -->
          <div class="flex items-center gap-2">
            <span class="text-sm font-medium text-gray-700">Type:</span>
            <div class="flex bg-gray-100 rounded-lg p-1">
              <button
                class="px-3 py-1.5 text-sm font-medium rounded-md transition-all {chartType === 'line' ? 'bg-white text-gray-900 shadow-sm' : 'text-gray-600 hover:text-gray-900'}"
                onclick={() => chartType = 'line'}
              >
                Line
              </button>
              <button
                class="px-3 py-1.5 text-sm font-medium rounded-md transition-all {chartType === 'bar' ? 'bg-white text-gray-900 shadow-sm' : 'text-gray-600 hover:text-gray-900'}"
                onclick={() => chartType = 'bar'}
              >
                Bar
              </button>
            </div>
          </div>

          <!-- Zoom Controls -->
          {#if showZoomControls}
            <div class="flex items-center gap-1">
              <span class="text-sm font-medium text-gray-700 mr-2">Zoom:</span>
              <button
                class="p-2 text-gray-600 hover:text-blue-600 hover:bg-blue-50 rounded-lg transition-all"
                onclick={zoomIn}
                title="Zoom In"
              >
                <ZoomIn class="w-4 h-4" />
              </button>
              <button
                class="p-2 text-gray-600 hover:text-blue-600 hover:bg-blue-50 rounded-lg transition-all"
                onclick={zoomOut}
                title="Zoom Out"
              >
                <ZoomOut class="w-4 h-4" />
              </button>
              <button
                class="p-2 text-gray-600 hover:text-blue-600 hover:bg-blue-50 rounded-lg transition-all"
                onclick={resetZoom}
                title="Reset Zoom"
              >
                <RotateCcw class="w-4 h-4" />
              </button>
            </div>
          {/if}

        </div>
        
        <!-- Interactive Settings -->
        <div class="mt-4 pt-4 border-t border-gray-100">
          <div class="flex flex-wrap items-center gap-6 text-sm">
            <label class="flex items-center gap-2 cursor-pointer">
              <input 
                type="checkbox" 
                bind:checked={enableZoom}
                class="rounded border-gray-300 text-blue-600 focus:ring-blue-500 focus:ring-offset-0"
              />
              <span class="text-gray-700">Enable Zoom</span>
            </label>
            <label class="flex items-center gap-2 cursor-pointer">
              <input 
                type="checkbox" 
                bind:checked={enablePan}
                class="rounded border-gray-300 text-blue-600 focus:ring-blue-500 focus:ring-offset-0"
              />
              <span class="text-gray-700">Enable Pan</span>
            </label>
            <label class="flex items-center gap-2 cursor-pointer">
              <input 
                type="checkbox" 
                bind:checked={showTooltips}
                class="rounded border-gray-300 text-blue-600 focus:ring-blue-500 focus:ring-offset-0"
              />
              <span class="text-gray-700">Show Tooltips</span>
            </label>
            <label class="flex items-center gap-2 cursor-pointer">
              <input 
                type="checkbox" 
                bind:checked={animateChart}
                class="rounded border-gray-300 text-blue-600 focus:ring-blue-500 focus:ring-offset-0"
              />
              <span class="text-gray-700">Animations</span>
            </label>
          </div>
        </div>
      </div>

      <!-- Interactive Legend -->
      {#if series.length > 0}
        <div class="bg-white border border-gray-200 rounded-xl p-4 mb-6 shadow-sm">
          <h4 class="text-sm font-semibold text-gray-900 mb-3">Data Series</h4>
          <div class="flex flex-wrap gap-2">
            {#each series as seriesData, index}
              {@const isVisible = !chart?.getDatasetMeta(index)?.hidden}
              <button
                class="flex items-center gap-2 px-3 py-2 rounded-lg border transition-all {
                  isVisible 
                    ? 'border-gray-300 bg-white hover:bg-gray-50' 
                    : 'border-gray-200 bg-gray-100 text-gray-500'
                }"
                onclick={() => toggleSeries(cleanMetricName(seriesData.metric_name))}
              >
                <div 
                  class="w-3 h-3 rounded-full border-2"
                  style="background-color: {isVisible ? colors[index % colors.length] : 'transparent'}; border-color: {colors[index % colors.length]}"
                ></div>
                <span class="text-sm font-medium">{cleanMetricName(seriesData.metric_name)}</span>
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

      <!-- Chart Container -->
      <div class="bg-white rounded-lg p-6 mb-6 shadow-sm">
        <div class="chart-container">
          <canvas bind:this={chartCanvas} class="w-full h-96"></canvas>
        </div>
        
        <!-- Interactive hints -->
        <div class="mt-4 pt-4 border-t border-gray-100">
          <div class="flex flex-wrap gap-4 text-xs text-gray-500">
            <span class="flex items-center gap-1">
              <Move class="w-3 h-3" />
              Drag to pan
            </span>
            <span class="flex items-center gap-1">
              <ZoomIn class="w-3 h-3" />
              Mouse wheel to zoom
            </span>
          </div>
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