<script lang="ts">
  import { Chart, registerables } from 'chart.js';
  import 'chartjs-adapter-date-fns';
  import { onMount, onDestroy } from 'svelte';
  import {
    TrendingUp,
    TrendingDown,
    Minus,
    AlertTriangle,
    Info,
    CheckCircle,
    BarChart3,
    Maximize2,
    Filter,
    Eye,
    EyeOff,
    Calendar,
    Activity,
    Star,
    BarChart2,
    Check,
    Circle,
    CheckSquare,
    MessageSquare,
  } from 'lucide-svelte';
  import { NoDataAvailable } from '$lib/components/ui';

  interface Props {
    data: any;
  }

  let { data }: Props = $props();
  let chartContainer: HTMLDivElement;
  let mounted = $state(false);
  let chartInstances: Chart[] = [];

  Chart.register(...registerables);

  // Interactive features
  let hoveredMetric = $state<string | null>(null);
  let showDetailedBreakdown = $state(false);
  let filteredMetricTypes = $state<string[]>([]);
  let sortBy = $state<'name' | 'change' | 'value'>('change');
  let sortOrder = $state<'asc' | 'desc'>('desc');
  let viewMode = $state<'cards' | 'chart' | 'table'>('cards');

  function formatValue(
    value: number,
    metricType: string,
    includeUnits: boolean = true,
    metricName?: string,
    isAverage: boolean = false,
    comparison?: any,
    periodData?: any
  ): string {
    if (metricType.includes('rate') || metricType.includes('completion')) {
      return value.toFixed(1) + '%';
    } else if (metricType.includes('time')) {
      return value.toFixed(1) + 'm';
    } else if (metricType.startsWith('question_')) {
      if (!includeUnits) return value.toFixed(2);

      // Try to get the actual question type from metadata
      let questionType = '';
      try {
        if (comparison?.metadata) {
          const metadata = JSON.parse(comparison.metadata);
          questionType = metadata.question_type || '';
        }
      } catch (e) {
        // Fallback - continue with name-based detection
      }

      // Use actual question type if available, otherwise fall back to name detection
      if (questionType === 'yes_no') {
        return value.toFixed(1) + '% yes';
      } else if (questionType === 'rating') {
        // Rating questions (1-5)
        if (isAverage) {
          return value.toFixed(1) + '/5';
        } else {
          // For totals, use the average from period data
          if (periodData?.average) {
            return periodData.average.toFixed(1) + '/5';
          }
          return value.toFixed(1) + '/5';
        }
      } else if (questionType === 'scale') {
        // Scale questions (1-10)
        let scaleValue = value;
        if (!isAverage && periodData?.average) {
          scaleValue = periodData.average;
        }

        // Try to get min/max labels from metadata
        let minLabel = '';
        let maxLabel = '';
        try {
          if (comparison?.metadata) {
            const metadata = JSON.parse(comparison.metadata);
            minLabel = metadata.min_label || '';
            maxLabel = metadata.max_label || '';
          }
        } catch (e) {}

        const baseValue = scaleValue.toFixed(1) + '/10';

        if (minLabel && maxLabel) {
          return `${baseValue}\n${minLabel} → ${maxLabel}`;
        }
        return baseValue;
      } else if (questionType === 'single_choice') {
        // Single choice questions - show most popular choice if available
        if (periodData?.most_popular_choice) {
          return `"${periodData.most_popular_choice.choice}" (${periodData.most_popular_choice.count})`;
        }
        // Fallback to response count
        if (includeUnits) {
          return `${value.toFixed(0)} responses`;
        }
        return value.toFixed(0);
      } else if (questionType === 'multi_choice') {
        // Multi choice questions - show top 3 choices if available
        if (periodData?.top_choices && periodData.top_choices.length > 0) {
          const topChoicesText = periodData.top_choices
            .slice(0, 3)
            .map(choice => `${choice.choice} (${choice.count})`)
            .join(', ');
          return topChoicesText;
        }
        // Fallback to response count
        if (includeUnits) {
          return `${value.toFixed(0)} responses`;
        }
        return value.toFixed(0);
      } else if (questionType === 'text') {
        // Text sentiment - show both score and meaningful description
        const getSentimentLabel = score => {
          if (score >= 0.5) return 'Very Positive';
          if (score >= 0.1) return 'Positive';
          if (score >= -0.1) return 'Neutral';
          if (score >= -0.5) return 'Negative';
          return 'Very Negative';
        };

        const sentimentLabel = getSentimentLabel(value);
        const scoreStr = value >= 0 ? `+${value.toFixed(2)}` : value.toFixed(2);
        return `${scoreStr} (${sentimentLabel})`;
      } else {
        // Fallback to name-based detection for compatibility
        const lowerName = metricName?.toLowerCase() || '';

        if (
          lowerName.includes('recommend') ||
          lowerName.includes('likelihood')
        ) {
          if (isAverage) {
            return value.toFixed(1) + '% likely';
          } else {
            return value.toFixed(0) + ' total %';
          }
        } else if (
          lowerName.includes('rate') ||
          lowerName.includes('rating') ||
          lowerName.includes('experience')
        ) {
          if (isAverage) {
            if (value <= 5) {
              return value.toFixed(1) + '/5';
            } else {
              return value.toFixed(1) + '/10';
            }
          } else {
            return value.toFixed(0) + ' total';
          }
        } else {
          return isAverage
            ? value.toFixed(2) + ' avg'
            : value.toFixed(0) + ' total';
        }
      }
    } else if (metricType === 'survey_responses') {
      return includeUnits
        ? Math.round(value).toLocaleString() + ' responses'
        : Math.round(value).toLocaleString();
    } else if (
      metricType.includes('sentiment') ||
      metricName?.toLowerCase().includes('sentiment')
    ) {
      // Handle sentiment metrics
      const getSentimentLabel = score => {
        if (score >= 0.5) return 'Very Positive';
        if (score >= 0.1) return 'Positive';
        if (score >= -0.1) return 'Neutral';
        if (score >= -0.5) return 'Negative';
        return 'Very Negative';
      };

      const sentimentLabel = getSentimentLabel(value);
      const scoreStr = value >= 0 ? `+${value.toFixed(2)}` : value.toFixed(2);
      return `${scoreStr} (${sentimentLabel})`;
    } else {
      return Math.round(value).toLocaleString();
    }
  }

  function formatChange(change: number, metricType: string): string {
    const absChange = Math.abs(change);
    if (metricType.includes('rate') || metricType.includes('completion')) {
      return absChange.toFixed(1) + '%';
    } else if (metricType.includes('time')) {
      return absChange.toFixed(1) + 'm';
    } else if (metricType.includes('rating')) {
      return absChange.toFixed(2);
    } else {
      return Math.round(absChange).toLocaleString();
    }
  }

  function getTrendIcon(trend: string) {
    switch (trend) {
      case 'improving':
        return TrendingUp;
      case 'declining':
        return TrendingDown;
      default:
        return Minus;
    }
  }

  function getTrendColor(trend: string, changePercent: number): string {
    if (Math.abs(changePercent) < 5) {
      return 'text-gray-600';
    }

    switch (trend) {
      case 'improving':
        return 'text-green-600';
      case 'declining':
        return 'text-red-600';
      default:
        return 'text-gray-600';
    }
  }

  function getBackgroundColor(trend: string, changePercent: number): string {
    if (Math.abs(changePercent) < 5) {
      return 'bg-gray-50 border-gray-200';
    }

    switch (trend) {
      case 'improving':
        return 'bg-green-50 border-green-200';
      case 'declining':
        return 'bg-red-50 border-red-200';
      default:
        return 'bg-gray-50 border-gray-200';
    }
  }

  function getInsightIcon(severity: string) {
    switch (severity) {
      case 'warning':
        return AlertTriangle;
      case 'info':
        return Info;
      default:
        return CheckCircle;
    }
  }

  function getInsightColor(severity: string): string {
    switch (severity) {
      case 'warning':
        return 'text-orange-600';
      case 'info':
        return 'text-blue-600';
      default:
        return 'text-green-600';
    }
  }

  function formatPeriodDate(dateString: string): string {
    return new Date(dateString).toLocaleDateString('en-US', {
      month: 'short',
      day: 'numeric',
      year: 'numeric',
    });
  }

  // Interactive functions

  function toggleMetricFilter(metricType: string) {
    if (filteredMetricTypes.includes(metricType)) {
      filteredMetricTypes = filteredMetricTypes.filter(t => t !== metricType);
    } else {
      filteredMetricTypes = [...filteredMetricTypes, metricType];
    }
  }

  function sortComparisons(comparisons: any[]) {
    return [...comparisons].sort((a, b) => {
      let aValue: number, bValue: number;

      switch (sortBy) {
        case 'name':
          aValue = a.metric_name.localeCompare(b.metric_name);
          bValue = 0;
          break;
        case 'change':
          aValue = Math.abs(a.change_percent);
          bValue = Math.abs(b.change_percent);
          break;
        case 'value':
          aValue = a.period2.value;
          bValue = b.period2.value;
          break;
        default:
          return 0;
      }

      const result = sortBy === 'name' ? aValue : aValue - bValue;
      return sortOrder === 'asc' ? result : -result;
    });
  }

  let allComparisons = $derived(data?.comparisons || []);
  let insights = $derived(data?.insights || []);

  // Filter and sort comparisons
  let comparisons = $derived(
    sortComparisons(
      filteredMetricTypes.length > 0
        ? allComparisons.filter((comp: any) =>
            filteredMetricTypes.some(type => comp.metric_type.includes(type))
          )
        : allComparisons
    )
  );

  // Get unique metric types for filtering
  let availableMetricTypes = $derived(
    Array.from(
      allComparisons.reduce((types: Set<string>, comp: any) => {
        if (comp.metric_type.startsWith('question_')) {
          types.add('Questions');
        } else if (comp.metric_type.includes('survey')) {
          types.add('Survey');
        } else {
          types.add('General');
        }
        return types;
      }, new Set<string>())
    )
  );

  // Group comparisons by question type for better visualization
  let groupedComparisons = $derived(() => {
    const groups = new Map();

    comparisons.forEach((comp: any) => {
      let questionType = 'other';
      try {
        if (comp.metadata) {
          const metadata = JSON.parse(comp.metadata);
          questionType = metadata.question_type || 'other';
        }
      } catch (e) {}

      if (!groups.has(questionType)) {
        groups.set(questionType, []);
      }
      groups.get(questionType).push(comp);
    });

    return groups;
  });

  function renderCharts() {
    if (!mounted || !chartContainer) return;

    // Destroy existing chart instances
    chartInstances.forEach(chart => {
      if (chart) {
        chart.destroy();
      }
    });
    chartInstances = [];

    // Clear previous charts
    chartContainer.innerHTML = '';

    if (comparisons.length === 0) {
      chartContainer.innerHTML =
        '<div class="text-center py-8 text-gray-500">No comparison data for chart</div>';
      return;
    }

    try {
      const containerWidth = chartContainer.offsetWidth || 800;
      const chartWidth = containerWidth - 48; // Use full width minus padding

      // Create charts for each question type
      groupedComparisons().forEach((questions, questionType) => {
        // Create section divider with title
        const sectionHeader = document.createElement('div');
        sectionHeader.className = 'flex items-center gap-3 mb-6 mt-8';

        // Add decorative line
        const leftLine = document.createElement('div');
        leftLine.className =
          'flex-1 h-px bg-gradient-to-r from-transparent to-gray-200';
        sectionHeader.appendChild(leftLine);

        // Add section title
        const titleContainer = document.createElement('div');
        titleContainer.className = 'flex items-center gap-2 px-4';
        titleContainer.innerHTML = `
          <div class="h-8 w-8 bg-gradient-to-br from-blue-500 to-indigo-600 rounded-lg flex items-center justify-center shadow-lg">
            <svg class="w-5 h-5 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v14a2 2 0 01-2 2h-2a2 2 0 01-2-2z" />
            </svg>
          </div>
          <h4 class="text-lg font-bold text-gray-900">${getQuestionTypeTitle(questionType)}</h4>
        `;
        sectionHeader.appendChild(titleContainer);

        // Add right decorative line
        const rightLine = document.createElement('div');
        rightLine.className =
          'flex-1 h-px bg-gradient-to-l from-transparent to-gray-200';
        sectionHeader.appendChild(rightLine);

        chartContainer.appendChild(sectionHeader);

        // Create grid container for charts
        const gridContainer = document.createElement('div');
        gridContainer.className = 'grid grid-cols-1 gap-6 mb-8';

        questions.forEach(comparison => {
          // Create chart container with subtle styling
          const chartCard = document.createElement('div');
          chartCard.className = 'group relative';

          // Add subtle background
          const bgDiv = document.createElement('div');
          bgDiv.className =
            'absolute inset-0 bg-gradient-to-br from-gray-50/50 to-white/50 rounded-xl opacity-0 group-hover:opacity-100 transition-opacity duration-300';
          chartCard.appendChild(bgDiv);

          // Main content container
          const contentDiv = document.createElement('div');
          contentDiv.className = 'relative p-6';

          // Add chart title
          const chartTitle = document.createElement('h5');
          chartTitle.className =
            'text-sm font-semibold text-gray-700 mb-4 flex items-center gap-2';
          chartTitle.innerHTML = `
            <div class="h-1.5 w-1.5 bg-blue-500 rounded-full"></div>
            <span class="truncate">${comparison.metric_name}</span>
          `;
          chartTitle.title = comparison.metric_name;
          contentDiv.appendChild(chartTitle);

          // Create plot container
          const plotContainer = document.createElement('div');
          plotContainer.className = 'flex justify-center';

          // Use full width for charts
          const fullChartWidth = chartWidth - 48;
          const chart = createChartForQuestionType(
            comparison,
            questionType,
            fullChartWidth
          );
          plotContainer.appendChild(chart);
          contentDiv.appendChild(plotContainer);

          chartCard.appendChild(contentDiv);
          gridContainer.appendChild(chartCard);
        });

        chartContainer.appendChild(gridContainer);
      });
    } catch (error) {
      console.error('Error rendering comparison charts:', error);
      chartContainer.innerHTML =
        '<div class="text-center py-8 text-red-500">Error rendering comparison charts</div>';
    }
  }

  function getQuestionTypeTitle(questionType: string): string {
    switch (questionType) {
      case 'single_choice':
        return 'Single Choice Questions';
      case 'multi_choice':
        return 'Multiple Choice Questions';
      case 'rating':
        return 'Rating Questions (1-5)';
      case 'scale':
        return 'Scale Questions (1-10)';
      case 'yes_no':
        return 'Yes/No Questions';
      case 'text':
        return 'Text/Sentiment Questions';
      default:
        return 'Other Metrics';
    }
  }

  function createChartForQuestionType(
    comparison: any,
    questionType: string,
    width: number
  ) {
    const canvas = document.createElement('canvas');
    canvas.style.width = width + 'px';
    canvas.style.height = '300px';
    canvas.style.maxWidth = '100%';

    try {
      if (questionType === 'single_choice') {
        return createDoughnutChart(
          canvas,
          comparison,
          'Single Choice Comparison'
        );
      } else if (questionType === 'multi_choice') {
        return createPolarChart(
          canvas,
          comparison,
          'Multiple Choice Comparison'
        );
      } else if (questionType === 'yes_no') {
        console.log('Creating yes/no chart for:', comparison);
        return createPieChart(canvas, comparison, 'Yes/No Comparison');
      } else if (questionType === 'rating' || questionType === 'scale') {
        return createBarChart(
          canvas,
          comparison,
          `${questionType.charAt(0).toUpperCase() + questionType.slice(1)} Comparison`
        );
      } else {
        return createBarChart(canvas, comparison, 'Value Comparison');
      }
    } catch (error) {
      console.error(
        `Error creating chart for question type ${questionType}:`,
        error
      );
      // Create error message on canvas
      const ctx = canvas.getContext('2d');
      if (ctx) {
        ctx.fillStyle = '#EF4444';
        ctx.font = '14px sans-serif';
        ctx.textAlign = 'center';
        ctx.fillText(
          'Error rendering chart',
          canvas.width / 2,
          canvas.height / 2
        );
      }
      return canvas;
    }
  }

  function createDoughnutChart(
    canvas: HTMLCanvasElement,
    comparison: any,
    title: string
  ) {
    const choiceData = getChoiceDistributionData(comparison);

    // Create subtitle with meaningful information
    const period1Info = formatValue(
      comparison.period1.value,
      comparison.metric_type,
      true,
      comparison.metric_name,
      false,
      comparison,
      comparison.period1
    );
    const period2Info = formatValue(
      comparison.period2.value,
      comparison.metric_type,
      true,
      comparison.metric_name,
      false,
      comparison,
      comparison.period2
    );

    const chart = new Chart(canvas, {
      type: 'doughnut',
      data: {
        labels: choiceData.labels,
        datasets: [
          {
            label: 'Period 1',
            data: choiceData.period1Data,
            backgroundColor: [
              '#3B82F680',
              '#EF444480',
              '#F59E0B80',
              '#10B98180',
              '#8B5CF680',
              '#F97316A0',
            ],
            borderColor: [
              '#3B82F6',
              '#EF4444',
              '#F59E0B',
              '#10B981',
              '#8B5CF6',
              '#F97316',
            ],
            borderWidth: 2,
          },
          {
            label: 'Period 2',
            data: choiceData.period2Data,
            backgroundColor: [
              '#3B82F6B0',
              '#EF4444B0',
              '#F59E0BB0',
              '#10B981B0',
              '#8B5CF6B0',
              '#F97316B0',
            ],
            borderColor: [
              '#3B82F6',
              '#EF4444',
              '#F59E0B',
              '#10B981',
              '#8B5CF6',
              '#F97316',
            ],
            borderWidth: 2,
          },
        ],
      },
      options: {
        responsive: true,
        maintainAspectRatio: false,
        plugins: {
          title: {
            display: true,
            text: formatChartTitle(title, period1Info, period2Info),
            font: { size: 12 },
          },
          legend: { position: 'bottom' },
          tooltip: {
            callbacks: {
              label: context =>
                `${context.dataset.label} - ${context.label}: ${context.parsed} selections`,
            },
          },
        },
      },
    });

    chartInstances.push(chart);
    return canvas;
  }

  function createPolarChart(
    canvas: HTMLCanvasElement,
    comparison: any,
    title: string
  ) {
    const choiceData = getChoiceDistributionData(comparison);

    // Create subtitle with meaningful information
    const period1Info = formatValue(
      comparison.period1.value,
      comparison.metric_type,
      true,
      comparison.metric_name,
      false,
      comparison,
      comparison.period1
    );
    const period2Info = formatValue(
      comparison.period2.value,
      comparison.metric_type,
      true,
      comparison.metric_name,
      false,
      comparison,
      comparison.period2
    );

    // For polar charts, use grouped bar instead for better comparison
    const chart = new Chart(canvas, {
      type: 'bar',
      data: {
        labels: choiceData.labels,
        datasets: [
          {
            label: 'Period 1',
            data: choiceData.period1Data,
            backgroundColor: '#3B82F680',
            borderColor: '#3B82F6',
            borderWidth: 2,
            maxBarThickness: 60,
            barPercentage: 0.8,
            categoryPercentage: 0.9,
          },
          {
            label: 'Period 2',
            data: choiceData.period2Data,
            backgroundColor: '#10B98180',
            borderColor: '#10B981',
            borderWidth: 2,
            maxBarThickness: 60,
            barPercentage: 0.8,
            categoryPercentage: 0.9,
          },
        ],
      },
      options: {
        responsive: true,
        maintainAspectRatio: false,
        plugins: {
          title: {
            display: true,
            text: formatChartTitle(title, period1Info, period2Info),
            font: { size: 12 },
          },
          legend: { position: 'bottom' },
          tooltip: {
            callbacks: {
              label: context =>
                `${context.dataset.label}: ${context.parsed.y} selections`,
            },
          },
        },
        scales: {
          y: {
            beginAtZero: true,
            grid: { display: true },
          },
          x: {
            grid: { display: false },
          },
        },
      },
    });

    chartInstances.push(chart);
    return canvas;
  }

  function createPieChart(
    canvas: HTMLCanvasElement,
    comparison: any,
    title: string
  ) {
    // For yes/no questions, handle different data formats
    console.log('Yes/No chart data:', comparison);
    let chartData;

    if (
      comparison.period1.choice_distribution ||
      comparison.period2.choice_distribution
    ) {
      // Use choice distribution if available
      console.log('Using choice distribution for yes/no');
      chartData = getChoiceDistributionData(comparison);
    } else {
      // Fallback to regular value data for yes/no
      console.log('Using fallback data for yes/no');
      chartData = {
        labels: ['Period 1', 'Period 2'],
        period1Data: [
          comparison.period1.average || comparison.period1.value || 0,
        ],
        period2Data: [
          comparison.period2.average || comparison.period2.value || 0,
        ],
      };
    }

    console.log('Final yes/no chart data:', chartData);

    // Create subtitle with meaningful information
    const period1Info = formatValue(
      comparison.period1.average || comparison.period1.value || 0,
      comparison.metric_type,
      true,
      comparison.metric_name,
      false,
      comparison,
      comparison.period1
    );
    const period2Info = formatValue(
      comparison.period2.average || comparison.period2.value || 0,
      comparison.metric_type,
      true,
      comparison.metric_name,
      false,
      comparison,
      comparison.period2
    );

    // For yes/no questions, use simple bar chart showing period comparison
    const chart = new Chart(canvas, {
      type: 'bar',
      data: {
        labels: ['Period 1', 'Period 2'],
        datasets: [
          {
            label: 'Value',
            data: [
              comparison.period1.average || comparison.period1.value || 0,
              comparison.period2.average || comparison.period2.value || 0,
            ],
            backgroundColor: ['#3B82F680', '#10B98180'],
            borderColor: ['#3B82F6', '#10B981'],
            borderWidth: 2,
            maxBarThickness: 80,
            barPercentage: 0.7,
            categoryPercentage: 0.8,
          },
        ],
      },
      options: {
        responsive: true,
        maintainAspectRatio: false,
        plugins: {
          title: {
            display: true,
            text: formatChartTitle(title, period1Info, period2Info),
            font: { size: 12 },
          },
          legend: { display: false },
          tooltip: {
            callbacks: {
              label: context => {
                const period = context.label;
                const value = context.parsed.y;
                const formattedValue = formatValue(
                  value,
                  comparison.metric_type,
                  true,
                  comparison.metric_name,
                  false,
                  comparison,
                  period === 'Period 1'
                    ? comparison.period1
                    : comparison.period2
                );
                return `${period}: ${formattedValue.split('\n')[0]}`;
              },
            },
          },
        },
        scales: {
          y: {
            beginAtZero: true,
            grid: { display: true },
          },
          x: {
            grid: { display: false },
          },
        },
      },
    });

    chartInstances.push(chart);
    return canvas;
  }

  function createBarChart(
    canvas: HTMLCanvasElement,
    comparison: any,
    title: string
  ) {
    const valueData = [
      comparison.period1.average || comparison.period1.value,
      comparison.period2.average || comparison.period2.value,
    ];

    // Create subtitle with meaningful information
    const period1Info = formatValue(
      valueData[0],
      comparison.metric_type,
      true,
      comparison.metric_name,
      true,
      comparison,
      comparison.period1
    );
    const period2Info = formatValue(
      valueData[1],
      comparison.metric_type,
      true,
      comparison.metric_name,
      true,
      comparison,
      comparison.period2
    );

    // Check if this is a sentiment metric
    const isSentiment =
      comparison.metric_type.includes('sentiment') ||
      comparison.metric_type.includes('text') ||
      (Math.abs(valueData[0]) <= 1 && Math.abs(valueData[1]) <= 1);

    // Check if this is a count/response metric
    const isCount =
      comparison.metric_type === 'survey_responses' ||
      comparison.metric_type.includes('count');

    const chart = new Chart(canvas, {
      type: 'bar',
      data: {
        labels: ['Period 1', 'Period 2'],
        datasets: [
          {
            label: 'Value',
            data: valueData,
            backgroundColor: isSentiment
              ? valueData.map(v => (v >= 0 ? '#10B98180' : '#EF444480'))
              : ['#3B82F680', '#10B98180'],
            borderColor: isSentiment
              ? valueData.map(v => (v >= 0 ? '#10B981' : '#EF4444'))
              : ['#3B82F6', '#10B981'],
            borderWidth: 2,
            maxBarThickness: 80,
            barPercentage: 0.7,
            categoryPercentage: 0.8,
          },
        ],
      },
      options: {
        responsive: true,
        maintainAspectRatio: false,
        plugins: {
          title: {
            display: true,
            text: formatChartTitle(title, period1Info, period2Info),
            font: { size: 12 },
          },
          legend: { display: false },
          tooltip: {
            callbacks: {
              label: context => {
                const period = context.label;
                const value = context.parsed.y;
                // For survey responses, ensure we show integers
                if (
                  comparison.metric_type === 'survey_responses' ||
                  comparison.metric_type.includes('count')
                ) {
                  return `${period}: ${Math.round(value).toLocaleString()} responses`;
                }
                const formattedValue = formatValue(
                  value,
                  comparison.metric_type,
                  true,
                  comparison.metric_name,
                  true,
                  comparison,
                  period === 'Period 1'
                    ? comparison.period1
                    : comparison.period2
                );
                return `${period}: ${formattedValue.split('\n')[0]}`;
              },
            },
          },
        },
        scales: {
          y: {
            beginAtZero: isCount || !isSentiment,
            grid: { display: true },
            ticks: {
              callback: function (value: any) {
                // Check if this is a survey response count
                if (
                  comparison.metric_type === 'survey_responses' ||
                  comparison.metric_type.includes('count') ||
                  comparison.metric_name?.toLowerCase().includes('responses')
                ) {
                  return Math.round(value).toLocaleString();
                }
                // Check if this is a sentiment metric (values between -1 and 1)
                if (
                  comparison.metric_type.includes('sentiment') ||
                  comparison.metric_type.includes('text') ||
                  (Math.abs(valueData[0]) <= 1 && Math.abs(valueData[1]) <= 1)
                ) {
                  const getSentimentLabel = (score: number) => {
                    if (score >= 0.5) return 'Very Positive';
                    if (score >= 0.1) return 'Positive';
                    if (score >= -0.1) return 'Neutral';
                    if (score >= -0.5) return 'Negative';
                    return 'Very Negative';
                  };
                  return `${value.toFixed(1)} (${getSentimentLabel(value)})`;
                }
                return value;
              },
              // Force integer steps for survey responses
              stepSize:
                comparison.metric_type === 'survey_responses' ||
                comparison.metric_type.includes('count')
                  ? 1
                  : undefined,
            },
          },
        },
      },
    });

    chartInstances.push(chart);
    return canvas;
  }

  function getChoiceDistributionData(comparison: any) {
    const allChoices = new Set();
    if (comparison.period1.choice_distribution) {
      Object.keys(comparison.period1.choice_distribution).forEach(choice =>
        allChoices.add(choice)
      );
    }
    if (comparison.period2.choice_distribution) {
      Object.keys(comparison.period2.choice_distribution).forEach(choice =>
        allChoices.add(choice)
      );
    }

    const labels = Array.from(allChoices);
    const period1Data = labels.map(
      choice => comparison.period1.choice_distribution?.[choice] || 0
    );
    const period2Data = labels.map(
      choice => comparison.period2.choice_distribution?.[choice] || 0
    );

    return { labels, period1Data, period2Data };
  }

  function formatChartTitle(
    title: string,
    period1Info: string,
    period2Info: string
  ): string[] {
    const titleLines = [title];

    // Split period info by newlines to handle scale labels properly
    const period1Lines = period1Info.split('\n');
    const period2Lines = period2Info.split('\n');

    // Add main period values
    titleLines.push(`Period 1: ${period1Lines[0]}`);
    titleLines.push(`Period 2: ${period2Lines[0]}`);

    // Add scale labels if they exist (for scale questions)
    if (period1Lines.length > 1 && period1Lines[1]) {
      titleLines.push(period1Lines[1]);
    }

    return titleLines;
  }

  // Re-render when data or view mode changes
  $effect(() => {
    if (mounted && viewMode === 'chart') {
      // Small delay to ensure DOM is ready
      setTimeout(() => {
        renderCharts();
      }, 100);
    }
  });

  // Re-render charts when comparison data changes
  $effect(() => {
    // Watch for changes in comparisons data
    comparisons.length;
    if (mounted && viewMode === 'chart' && comparisons.length > 0) {
      setTimeout(() => {
        renderCharts();
      }, 50);
    }
  });

  onMount(() => {
    mounted = true;
    if (viewMode === 'chart') {
      renderCharts();
    }
  });

  onDestroy(() => {
    // Clean up all chart instances
    chartInstances.forEach(chart => {
      if (chart) {
        chart.destroy();
      }
    });
    chartInstances = [];
  });
</script>

<div class="comparison-chart">
  <div class="mb-8">
    {#if comparisons.length === 0}
      <NoDataAvailable
        title="No Comparison Data"
        description="Select questions above to compare metrics between periods"
        icon={Activity}
        variant="inline" />
    {:else}
      <!-- Period Summary Header -->
      {#if data?.request}
        <div
          class="bg-gradient-to-br from-indigo-50 via-purple-50 to-pink-50 rounded-2xl p-6 mb-8 relative overflow-hidden">
          <div
            class="absolute inset-0 bg-gradient-to-br from-indigo-500/5 via-purple-500/5 to-pink-500/5">
          </div>
          <div class="relative">
            <div class="flex items-center gap-3 mb-4">
              <div
                class="h-10 w-10 bg-gradient-to-br from-indigo-500 to-purple-600 rounded-xl flex items-center justify-center shadow-lg shadow-indigo-500/25">
                <Calendar class="w-5 h-5 text-white" />
              </div>
              <h3 class="text-xl font-bold text-gray-900">
                Period Comparison Analysis
              </h3>
            </div>

            <div class="grid grid-cols-1 md:grid-cols-2 gap-6">
              <div
                class="bg-white/70 backdrop-blur-sm rounded-xl p-4 border border-white/50 shadow-sm">
                <div class="flex items-center gap-2 mb-2">
                  <div
                    class="h-6 w-6 bg-gradient-to-br from-blue-500 to-blue-600 rounded-md flex items-center justify-center">
                    <span class="text-xs font-bold text-white">1</span>
                  </div>
                  <h4 class="font-semibold text-gray-900">First Period</h4>
                </div>
                <p class="text-sm text-gray-700">
                  {formatPeriodDate(data.request.period1_start)} → {formatPeriodDate(
                    data.request.period1_end
                  )}
                </p>
              </div>

              <div
                class="bg-white/70 backdrop-blur-sm rounded-xl p-4 border border-white/50 shadow-sm">
                <div class="flex items-center gap-2 mb-2">
                  <div
                    class="h-6 w-6 bg-gradient-to-br from-purple-500 to-purple-600 rounded-md flex items-center justify-center">
                    <span class="text-xs font-bold text-white">2</span>
                  </div>
                  <h4 class="font-semibold text-gray-900">Second Period</h4>
                </div>
                <p class="text-sm text-gray-700">
                  {formatPeriodDate(data.request.period2_start)} → {formatPeriodDate(
                    data.request.period2_end
                  )}
                </p>
              </div>
            </div>
          </div>
        </div>
      {/if}

      <!-- Interactive Controls -->
      <div
        class="bg-white/80 backdrop-blur-sm border border-gray-200/50 rounded-2xl p-5 mb-8 shadow-lg">
        <div class="flex flex-wrap items-center justify-between gap-4">
          <!-- View Mode Controls -->
          <div class="flex items-center gap-3">
            <div class="flex items-center gap-2">
              <Maximize2 class="w-4 h-4 text-gray-500" />
              <span class="text-sm font-semibold text-gray-700">View Mode</span>
            </div>
            <div
              class="flex bg-gradient-to-r from-gray-100 to-gray-50 rounded-xl p-1 shadow-inner">
              <button
                class="px-4 py-2 text-sm font-medium rounded-lg transition-all duration-200 {viewMode ===
                'cards'
                  ? 'bg-white text-gray-900 shadow-md transform scale-105'
                  : 'text-gray-600 hover:text-gray-900 hover:bg-white/50'}"
                onclick={() => (viewMode = 'cards')}>
                <div class="flex items-center gap-2">
                  <svg
                    class="w-4 h-4"
                    fill="none"
                    stroke="currentColor"
                    viewBox="0 0 24 24">
                    <path
                      stroke-linecap="round"
                      stroke-linejoin="round"
                      stroke-width="2"
                      d="M4 6a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2H6a2 2 0 01-2-2V6zM14 6a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2h-2a2 2 0 01-2-2V6zM4 16a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2H6a2 2 0 01-2-2v-2zM14 16a2 2 0 012-2h2a2 2 0 012 2v2a2 2 0 01-2 2h-2a2 2 0 01-2-2v-2z" />
                  </svg>
                  Cards
                </div>
              </button>
              <button
                class="px-4 py-2 text-sm font-medium rounded-lg transition-all duration-200 {viewMode ===
                'chart'
                  ? 'bg-white text-gray-900 shadow-md transform scale-105'
                  : 'text-gray-600 hover:text-gray-900 hover:bg-white/50'}"
                onclick={() => (viewMode = 'chart')}>
                <div class="flex items-center gap-2">
                  <BarChart3 class="w-4 h-4" />
                  Chart
                </div>
              </button>
              <button
                class="px-4 py-2 text-sm font-medium rounded-lg transition-all duration-200 {viewMode ===
                'table'
                  ? 'bg-white text-gray-900 shadow-md transform scale-105'
                  : 'text-gray-600 hover:text-gray-900 hover:bg-white/50'}"
                onclick={() => (viewMode = 'table')}>
                <div class="flex items-center gap-2">
                  <svg
                    class="w-4 h-4"
                    fill="none"
                    stroke="currentColor"
                    viewBox="0 0 24 24">
                    <path
                      stroke-linecap="round"
                      stroke-linejoin="round"
                      stroke-width="2"
                      d="M3 10h18M3 14h18m-9-4v8m-7 0h14a2 2 0 002-2V8a2 2 0 00-2-2H5a2 2 0 00-2 2v8a2 2 0 002 2z" />
                  </svg>
                  Table
                </div>
              </button>
            </div>
          </div>

          <!-- Sort Controls -->
          <div class="flex items-center gap-2">
            <span class="text-sm font-medium text-gray-700">Sort by:</span>
            <select
              bind:value={sortBy}
              class="text-sm border border-gray-300 rounded-md px-2 py-1 bg-white">
              <option value="change">Change %</option>
              <option value="value">Current Value</option>
              <option value="name">Name</option>
            </select>
            <button
              class="p-1 text-gray-600 hover:text-blue-600 rounded transition-colors"
              onclick={() => (sortOrder = sortOrder === 'asc' ? 'desc' : 'asc')}
              title="Toggle sort order">
              <svg
                class="w-4 h-4 {sortOrder === 'desc'
                  ? 'rotate-180'
                  : ''} transition-transform"
                fill="none"
                stroke="currentColor"
                viewBox="0 0 24 24">
                <path
                  stroke-linecap="round"
                  stroke-linejoin="round"
                  stroke-width="2"
                  d="M7 16V4m0 0L3 8m4-4l4 4m6 0v12m0 0l4-4m-4 4l-4-4"></path>
              </svg>
            </button>
          </div>
        </div>
      </div>

      <!-- Dynamic View Based on Mode -->
      {#if viewMode === 'cards'}
        <!-- Enhanced Interactive Comparison Cards -->
        <div class="grid grid-cols-1 lg:grid-cols-2 gap-6 mb-8">
          {#each comparisons as comparison}
            {@const isChoiceQuestion = (() => {
              try {
                if (comparison.metadata) {
                  const metadata = JSON.parse(comparison.metadata);
                  return (
                    metadata.question_type === 'single_choice' ||
                    metadata.question_type === 'multi_choice'
                  );
                }
              } catch (e) {}
              return false;
            })()}
            {@const isPositiveTrend =
              comparison.trend === 'improving' ||
              (comparison.metric_type.includes('rating') &&
                comparison.change_percent > 0) ||
              (comparison.metric_type.includes('recommend') &&
                comparison.change_percent > 0)}
            {@const trendColorClass = isChoiceQuestion
              ? 'from-blue-50 to-indigo-50 border-blue-200'
              : Math.abs(comparison.change_percent) < 5
                ? 'from-gray-50 to-gray-100 border-gray-200'
                : isPositiveTrend
                  ? 'from-emerald-50 to-green-50 border-emerald-200'
                  : 'from-red-50 to-rose-50 border-red-200'}

            <div
              class="group relative bg-gradient-to-br {trendColorClass} rounded-2xl border-2 transition-all duration-300 hover:shadow-xl hover:scale-[1.02] overflow-hidden"
              onmouseenter={() => (hoveredMetric = comparison.metric_name)}
              onmouseleave={() => (hoveredMetric = null)}>
              <!-- Background Pattern -->
              <div class="absolute inset-0 opacity-5">
                <svg class="w-full h-full" xmlns="http://www.w3.org/2000/svg">
                  <pattern
                    id="grid-{comparison.metric_type}"
                    x="0"
                    y="0"
                    width="40"
                    height="40"
                    patternUnits="userSpaceOnUse">
                    <path
                      d="M 40 0 L 0 0 0 40"
                      fill="none"
                      stroke="currentColor"
                      stroke-width="1" />
                  </pattern>
                  <rect
                    width="100%"
                    height="100%"
                    fill="url(#grid-{comparison.metric_type})" />
                </svg>
              </div>

              <div class="relative p-6">
                <!-- Header -->
                <div class="flex items-start justify-between mb-6">
                  <div class="flex-1 pr-4 min-w-0">
                    <h4
                      class="font-bold text-gray-900 text-lg mb-2 leading-tight group-hover:text-indigo-700 transition-colors truncate">
                      {comparison.metric_name}
                    </h4>
                    <div class="flex items-center gap-2">
                      <span
                        class="inline-flex items-center gap-1 px-3 py-1 rounded-full text-xs font-medium bg-white/80 text-gray-700 shadow-sm whitespace-nowrap">
                        {#snippet questionTypeIcon()}
                          {@const questionType = (() => {
                            if (
                              comparison.metric_type.startsWith('question_')
                            ) {
                              try {
                                if (comparison.metadata) {
                                  const metadata = JSON.parse(
                                    comparison.metadata
                                  );
                                  return metadata.question_type || 'question';
                                }
                              } catch (e) {}
                            }
                            return 'metric';
                          })()}

                          {#if questionType === 'rating'}
                            <Star class="w-3 h-3" />
                          {:else if questionType === 'scale'}
                            <BarChart2 class="w-3 h-3" />
                          {:else if questionType === 'yes_no'}
                            <Check class="w-3 h-3" />
                          {:else if questionType === 'single_choice'}
                            <Circle class="w-3 h-3" />
                          {:else if questionType === 'multi_choice'}
                            <CheckSquare class="w-3 h-3" />
                          {:else if questionType === 'text'}
                            <MessageSquare class="w-3 h-3" />
                          {:else}
                            <BarChart3 class="w-3 h-3" />
                          {/if}
                        {/snippet}
                        {@render questionTypeIcon()}
                        {(() => {
                          if (comparison.metric_type.startsWith('question_')) {
                            let questionType = 'Question';
                            try {
                              if (comparison.metadata) {
                                const metadata = JSON.parse(
                                  comparison.metadata
                                );
                                questionType =
                                  metadata.question_type || 'Question';
                              }
                            } catch (e) {}

                            switch (questionType) {
                              case 'rating':
                                return 'Rating (1-5)';
                              case 'scale':
                                return 'Scale (1-10)';
                              case 'yes_no':
                                return 'Yes/No';
                              case 'single_choice':
                                return 'Single Choice';
                              case 'multi_choice':
                                return 'Multiple Choice';
                              case 'text':
                                return 'Text/Sentiment';
                              default:
                                return 'Question';
                            }
                          }
                          return 'Metric';
                        })()}
                      </span>
                    </div>
                  </div>

                  <!-- Trend Indicator -->
                  {#if (() => {
                    try {
                      if (comparison.metadata) {
                        const metadata = JSON.parse(comparison.metadata);
                        return metadata.question_type !== 'single_choice' && metadata.question_type !== 'multi_choice';
                      }
                    } catch (e) {}
                    return true;
                  })()}
                    <div class="flex flex-col items-end gap-1">
                      <div
                        class="flex items-center gap-2 bg-white/90 rounded-xl px-3 py-2 shadow-md">
                        {#snippet trendIcon()}
                          {@const TrendIcon = getTrendIcon(comparison.trend)}
                          <TrendIcon
                            class="w-5 h-5 {Math.abs(
                              comparison.change_percent
                            ) < 5
                              ? 'text-gray-500'
                              : isPositiveTrend
                                ? 'text-emerald-600'
                                : 'text-red-600'}" />
                        {/snippet}
                        {@render trendIcon()}
                        <span
                          class="font-bold text-lg {Math.abs(
                            comparison.change_percent
                          ) < 5
                            ? 'text-gray-700'
                            : isPositiveTrend
                              ? 'text-emerald-700'
                              : 'text-red-700'}">
                          {comparison.change_percent > 0
                            ? '+'
                            : ''}{comparison.change_percent.toFixed(1)}%
                        </span>
                      </div>
                      <span
                        class="text-xs font-medium {Math.abs(
                          comparison.change_percent
                        ) < 5
                          ? 'text-gray-500'
                          : isPositiveTrend
                            ? 'text-emerald-600'
                            : 'text-red-600'}">
                        {Math.abs(comparison.change_percent) < 5
                          ? 'Stable'
                          : isPositiveTrend
                            ? 'Improving'
                            : 'Declining'}
                      </span>
                    </div>
                  {/if}
                </div>

                <!-- Choice Question Note -->
                {#if (() => {
                  try {
                    if (comparison.metadata) {
                      const metadata = JSON.parse(comparison.metadata);
                      return metadata.question_type === 'single_choice' || metadata.question_type === 'multi_choice';
                    }
                  } catch (e) {}
                  return false;
                })()}
                  <div
                    class="bg-blue-50/80 border border-blue-200/50 rounded-lg p-3 mb-4">
                    <div class="flex items-start gap-2">
                      <Info
                        class="w-4 h-4 text-blue-600 mt-0.5 flex-shrink-0" />
                      <div class="text-xs text-blue-800">
                        <span class="font-medium">Note:</span>
                        {(() => {
                          try {
                            if (comparison.metadata) {
                              const metadata = JSON.parse(comparison.metadata);
                              if (metadata.question_type === 'single_choice') {
                                return 'For single choice questions, this comparison shows the most popular choice selected in each period.';
                              } else if (
                                metadata.question_type === 'multi_choice'
                              ) {
                                return 'For multiple choice questions, this comparison shows the top 3 most selected options in each period.';
                              }
                            }
                          } catch (e) {}
                          return 'For choice questions, this comparison shows the most popular selections in each period.';
                        })()}
                      </div>
                    </div>
                  </div>
                {/if}

                <!-- Period Data -->
                <div class="grid grid-cols-2 gap-4 mb-4">
                  <!-- Period 1 -->
                  <div class="bg-white/70 rounded-xl p-4 backdrop-blur-sm">
                    <div class="flex items-center gap-2 mb-3">
                      <div
                        class="h-5 w-5 bg-blue-500 rounded flex items-center justify-center">
                        <span class="text-[10px] font-bold text-white">1</span>
                      </div>
                      <h5 class="text-sm font-semibold text-gray-900">
                        Period 1
                      </h5>
                    </div>
                    <div class="space-y-2">
                      <div class="flex justify-between items-center">
                        <span class="text-xs text-gray-600"
                          >{(() => {
                            try {
                              if (comparison.metadata) {
                                const metadata = JSON.parse(
                                  comparison.metadata
                                );
                                if (metadata.question_type === 'single_choice')
                                  return 'Top Choice';
                                if (metadata.question_type === 'multi_choice')
                                  return 'Top 3 Choices';
                                return 'Value';
                              }
                            } catch (e) {}
                            return 'Value';
                          })()}</span>
                        <div class="flex items-center gap-1 text-right ml-2">
                          <span
                            class="font-bold text-gray-900 text-xs whitespace-pre-line">
                            {formatValue(
                              comparison.period1.value,
                              comparison.metric_type,
                              true,
                              comparison.metric_name,
                              false,
                              comparison,
                              comparison.period1
                            )}
                          </span>
                          {#if (() => {
                            try {
                              if (comparison.metadata) {
                                const metadata = JSON.parse(comparison.metadata);
                                return metadata.question_type === 'rating';
                              }
                            } catch (e) {}
                            return false;
                          })()}
                            <Star class="w-3 h-3 text-yellow-500" />
                          {/if}
                        </div>
                      </div>
                      {#if (() => {
                        try {
                          if (comparison.metadata) {
                            const metadata = JSON.parse(comparison.metadata);
                            return metadata.question_type !== 'single_choice' && metadata.question_type !== 'multi_choice';
                          }
                        } catch (e) {}
                        return true;
                      })()}
                        <div class="flex justify-between items-center">
                          <span class="text-xs text-gray-600">Average</span>
                          <div class="flex items-center gap-1 text-right ml-2">
                            <span
                              class="font-semibold text-gray-800 text-xs whitespace-pre-line">
                              {formatValue(
                                comparison.period1.average,
                                comparison.metric_type,
                                true,
                                comparison.metric_name,
                                true,
                                comparison,
                                comparison.period1
                              )}
                            </span>
                            {#if (() => {
                              try {
                                if (comparison.metadata) {
                                  const metadata = JSON.parse(comparison.metadata);
                                  return metadata.question_type === 'rating';
                                }
                              } catch (e) {}
                              return false;
                            })()}
                              <Star class="w-3 h-3 text-yellow-500" />
                            {/if}
                          </div>
                        </div>
                      {/if}
                      <div class="flex justify-between items-center">
                        <span class="text-xs text-gray-600">Responses</span>
                        <span class="font-medium text-gray-700">
                          {comparison.period1.count}
                        </span>
                      </div>
                    </div>
                  </div>

                  <!-- Period 2 -->
                  <div class="bg-white/70 rounded-xl p-4 backdrop-blur-sm">
                    <div class="flex items-center gap-2 mb-3">
                      <div
                        class="h-5 w-5 bg-purple-500 rounded flex items-center justify-center">
                        <span class="text-[10px] font-bold text-white">2</span>
                      </div>
                      <h5 class="text-sm font-semibold text-gray-900">
                        Period 2
                      </h5>
                    </div>
                    <div class="space-y-2">
                      <div class="flex justify-between items-center">
                        <span class="text-xs text-gray-600"
                          >{(() => {
                            try {
                              if (comparison.metadata) {
                                const metadata = JSON.parse(
                                  comparison.metadata
                                );
                                if (metadata.question_type === 'single_choice')
                                  return 'Top Choice';
                                if (metadata.question_type === 'multi_choice')
                                  return 'Top 3 Choices';
                                return 'Value';
                              }
                            } catch (e) {}
                            return 'Value';
                          })()}</span>
                        <div class="flex items-center gap-1 text-right ml-2">
                          <span
                            class="font-bold text-gray-900 text-xs whitespace-pre-line">
                            {formatValue(
                              comparison.period2.value,
                              comparison.metric_type,
                              true,
                              comparison.metric_name,
                              false,
                              comparison,
                              comparison.period2
                            )}
                          </span>
                          {#if (() => {
                            try {
                              if (comparison.metadata) {
                                const metadata = JSON.parse(comparison.metadata);
                                return metadata.question_type === 'rating';
                              }
                            } catch (e) {}
                            return false;
                          })()}
                            <Star class="w-3 h-3 text-yellow-500" />
                          {/if}
                        </div>
                      </div>
                      {#if (() => {
                        try {
                          if (comparison.metadata) {
                            const metadata = JSON.parse(comparison.metadata);
                            return metadata.question_type !== 'single_choice' && metadata.question_type !== 'multi_choice';
                          }
                        } catch (e) {}
                        return true;
                      })()}
                        <div class="flex justify-between items-center">
                          <span class="text-xs text-gray-600">Average</span>
                          <div class="flex items-center gap-1 text-right ml-2">
                            <span
                              class="font-semibold text-gray-800 text-xs whitespace-pre-line">
                              {formatValue(
                                comparison.period2.average,
                                comparison.metric_type,
                                true,
                                comparison.metric_name,
                                true,
                                comparison,
                                comparison.period2
                              )}
                            </span>
                            {#if (() => {
                              try {
                                if (comparison.metadata) {
                                  const metadata = JSON.parse(comparison.metadata);
                                  return metadata.question_type === 'rating';
                                }
                              } catch (e) {}
                              return false;
                            })()}
                              <Star class="w-3 h-3 text-yellow-500" />
                            {/if}
                          </div>
                        </div>
                      {/if}
                      <div class="flex justify-between items-center">
                        <span class="text-xs text-gray-600">Responses</span>
                        <span class="font-medium text-gray-700">
                          {comparison.period2.count}
                        </span>
                      </div>
                    </div>
                  </div>
                </div>

                <!-- Change Summary -->
                {#if (() => {
                  try {
                    if (comparison.metadata) {
                      const metadata = JSON.parse(comparison.metadata);
                      return metadata.question_type !== 'single_choice' && metadata.question_type !== 'multi_choice';
                    }
                  } catch (e) {}
                  return true;
                })()}
                  <div
                    class="bg-gradient-to-r from-white/50 to-white/30 rounded-xl p-4 backdrop-blur-sm border border-white/50">
                    <div class="flex items-center justify-between">
                      <div class="flex items-center gap-2">
                        <svg
                          class="w-4 h-4 text-gray-600"
                          fill="none"
                          stroke="currentColor"
                          viewBox="0 0 24 24">
                          <path
                            stroke-linecap="round"
                            stroke-linejoin="round"
                            stroke-width="2"
                            d="M7 16V4m0 0L3 8m4-4l4 4m6 0v12m0 0l4-4m-4 4l-4-4" />
                        </svg>
                        <span class="text-sm font-medium text-gray-700"
                          >Net Change</span>
                      </div>
                      <span
                        class="font-bold text-lg {Math.abs(
                          comparison.change_percent
                        ) < 5
                          ? 'text-gray-700'
                          : isPositiveTrend
                            ? 'text-emerald-700'
                            : 'text-red-700'}">
                        {comparison.change > 0 ? '+' : ''}{formatChange(
                          comparison.change,
                          comparison.metric_type
                        )}
                      </span>
                    </div>
                  </div>
                {/if}
              </div>
            </div>
          {/each}
        </div>
      {:else if viewMode === 'table'}
        <!-- Table View -->
        <div
          class="bg-white border border-gray-200 rounded-xl overflow-hidden shadow-sm mb-8">
          <div class="overflow-x-auto">
            <table class="w-full">
              <thead class="bg-gray-50 border-b border-gray-200">
                <tr>
                  <th
                    class="px-4 py-3 text-left text-xs font-semibold text-gray-600 uppercase tracking-wider"
                    >Metric</th>
                  <th
                    class="px-4 py-3 text-left text-xs font-semibold text-gray-600 uppercase tracking-wider"
                    >Type</th>
                  <th
                    class="px-4 py-3 text-right text-xs font-semibold text-gray-600 uppercase tracking-wider"
                    >{(() => {
                      // Check if we have any choice questions to adjust header
                      const hasChoiceQuestions = comparisons.some(comp => {
                        try {
                          if (comp.metadata) {
                            const metadata = JSON.parse(comp.metadata);
                            return (
                              metadata.question_type === 'single_choice' ||
                              metadata.question_type === 'multi_choice'
                            );
                          }
                        } catch (e) {}
                        return false;
                      });
                      return hasChoiceQuestions ? 'Period 1' : 'Period 1 Value';
                    })()}</th>
                  <th
                    class="px-4 py-3 text-right text-xs font-semibold text-gray-600 uppercase tracking-wider"
                    >{(() => {
                      // Check if we have any choice questions to adjust header
                      const hasChoiceQuestions = comparisons.some(comp => {
                        try {
                          if (comp.metadata) {
                            const metadata = JSON.parse(comp.metadata);
                            return (
                              metadata.question_type === 'single_choice' ||
                              metadata.question_type === 'multi_choice'
                            );
                          }
                        } catch (e) {}
                        return false;
                      });
                      return hasChoiceQuestions ? 'Period 2' : 'Period 2 Value';
                    })()}</th>
                  <th
                    class="px-4 py-3 text-right text-xs font-semibold text-gray-600 uppercase tracking-wider whitespace-nowrap"
                    >Change %</th>
                  <th
                    class="px-4 py-3 text-center text-xs font-semibold text-gray-600 uppercase tracking-wider"
                    >Trend</th>
                </tr>
              </thead>
              <tbody class="divide-y divide-gray-200">
                {#each comparisons as comparison}
                  <tr class="hover:bg-gray-50 transition-colors">
                    <td class="px-4 py-3">
                      <div class="font-medium text-gray-900">
                        {comparison.metric_name}
                      </div>
                    </td>
                    <td class="px-4 py-3">
                      <span
                        class="inline-block px-2 py-1 text-xs bg-gray-100 text-gray-700 rounded-full capitalize whitespace-nowrap truncate max-w-24">
                        {(() => {
                          if (comparison.metric_type.startsWith('question_')) {
                            // Try to get question type from metadata if available
                            let questionType = 'Question';
                            try {
                              if (comparison.metadata) {
                                const metadata = JSON.parse(
                                  comparison.metadata
                                );
                                questionType =
                                  metadata.question_type || 'Question';
                              }
                            } catch (e) {
                              // Fallback to generic question type
                            }

                            // Map question types to display names
                            switch (questionType) {
                              case 'rating':
                                return 'Rating (1-5)';
                              case 'scale':
                                return 'Scale (1-10)';
                              case 'yes_no':
                                return 'Yes/No';
                              case 'single_choice':
                                return 'Single Choice';
                              case 'multi_choice':
                                return 'Multiple Choice';
                              case 'text':
                                return 'Text/Sentiment';
                              default:
                                return 'Question';
                            }
                          }
                          return comparison.metric_type.replace('_', ' ');
                        })()}
                      </span>
                    </td>
                    <td class="px-4 py-3 text-right font-medium">
                      {formatValue(
                        comparison.period1.value,
                        comparison.metric_type,
                        true,
                        comparison.metric_name,
                        false,
                        comparison,
                        comparison.period1
                      )}
                    </td>
                    <td class="px-4 py-3 text-right font-medium">
                      {formatValue(
                        comparison.period2.value,
                        comparison.metric_type,
                        true,
                        comparison.metric_name,
                        false,
                        comparison,
                        comparison.period2
                      )}
                    </td>
                    <td class="px-4 py-3 text-right">
                      {#if (() => {
                        try {
                          if (comparison.metadata) {
                            const metadata = JSON.parse(comparison.metadata);
                            return metadata.question_type !== 'single_choice' && metadata.question_type !== 'multi_choice';
                          }
                        } catch (e) {}
                        return true;
                      })()}
                        <span
                          class="font-semibold {getTrendColor(
                            comparison.trend,
                            comparison.change_percent
                          )}">
                          {comparison.change_percent > 0
                            ? '+'
                            : ''}{comparison.change_percent.toFixed(1)}%
                        </span>
                      {:else}
                        <span class="text-gray-400 text-sm">-</span>
                      {/if}
                    </td>
                    <td class="px-4 py-3 text-center">
                      {#if (() => {
                        try {
                          if (comparison.metadata) {
                            const metadata = JSON.parse(comparison.metadata);
                            return metadata.question_type !== 'single_choice' && metadata.question_type !== 'multi_choice';
                          }
                        } catch (e) {}
                        return true;
                      })()}
                        {#snippet trendIcon()}
                          {@const TrendIcon = getTrendIcon(comparison.trend)}
                          <TrendIcon
                            class="w-4 h-4 mx-auto {getTrendColor(
                              comparison.trend,
                              comparison.change_percent
                            )}" />
                        {/snippet}
                        {@render trendIcon()}
                      {:else}
                        <span class="text-gray-400 text-sm">-</span>
                      {/if}
                    </td>
                  </tr>
                {/each}
              </tbody>
            </table>
          </div>
        </div>
      {:else if viewMode === 'chart'}
        <!-- Enhanced Chart View -->
        <div style="height: auto; min-height: 0;">
          <div class="flex items-center justify-between mb-6">
            <h4 class="text-xl font-bold text-gray-900 flex items-center gap-3">
              <div
                class="h-10 w-10 bg-gradient-to-br from-blue-500 to-indigo-600 rounded-xl flex items-center justify-center shadow-lg shadow-blue-500/25">
                <BarChart3 class="w-5 h-5 text-white" />
              </div>
              Visual Comparison
            </h4>
            <div
              class="flex items-center gap-2 text-sm text-gray-500 bg-gray-50 px-3 py-1.5 rounded-lg">
              <Activity class="w-4 h-4" />
              <span>{comparisons.length} metrics</span>
            </div>
          </div>

          <!-- Charts Container -->
          <div
            bind:this={chartContainer}
            class="chart-container w-full"
            style="height: auto; overflow: visible;">
          </div>
        </div>
      {/if}

      <!-- Insights -->
      {#if insights.length > 0}
        <div class="mt-8">
          <div class="flex items-center gap-3 mb-6">
            <div
              class="h-10 w-10 bg-gradient-to-br from-amber-500 to-orange-600 rounded-xl flex items-center justify-center shadow-lg shadow-amber-500/25">
              <Info class="w-5 h-5 text-white" />
            </div>
            <h4 class="text-xl font-bold text-gray-900">Key Insights</h4>
          </div>

          <div class="space-y-4">
            {#each insights as insight}
              {@const bgClass =
                insight.severity === 'warning'
                  ? 'from-orange-50 to-amber-50 border-orange-300'
                  : insight.severity === 'info'
                    ? 'from-blue-50 to-indigo-50 border-blue-300'
                    : 'from-green-50 to-emerald-50 border-green-300'}
              {@const iconBgClass =
                insight.severity === 'warning'
                  ? 'from-orange-400 to-amber-500'
                  : insight.severity === 'info'
                    ? 'from-blue-400 to-indigo-500'
                    : 'from-green-400 to-emerald-500'}

              <div
                class="relative bg-gradient-to-r {bgClass} rounded-xl border-2 p-5 shadow-sm hover:shadow-md transition-shadow duration-200">
                <div class="flex items-start gap-4">
                  <div class="flex-shrink-0">
                    <div
                      class="h-10 w-10 bg-gradient-to-br {iconBgClass} rounded-lg flex items-center justify-center shadow-md">
                      {#snippet insightIcon()}
                        {@const InsightIcon = getInsightIcon(insight.severity)}
                        <InsightIcon class="w-5 h-5 text-white" />
                      {/snippet}
                      {@render insightIcon()}
                    </div>
                  </div>

                  <div class="flex-1">
                    <div
                      class="font-semibold text-gray-900 text-base leading-tight mb-2">
                      {insight.message}
                    </div>

                    {#if insight.recommendation}
                      <div class="mt-3 p-3 bg-white/70 rounded-lg">
                        <div class="flex items-start gap-2">
                          <svg
                            class="w-4 h-4 text-gray-600 mt-0.5 flex-shrink-0"
                            fill="none"
                            stroke="currentColor"
                            viewBox="0 0 24 24">
                            <path
                              stroke-linecap="round"
                              stroke-linejoin="round"
                              stroke-width="2"
                              d="M9 5l7 7-7 7" />
                          </svg>
                          <div class="text-sm text-gray-700">
                            <span class="font-medium">Action:</span>
                            {insight.recommendation}
                          </div>
                        </div>
                      </div>
                    {/if}

                    <div class="mt-3 flex items-center gap-3 text-xs">
                      <span
                        class="inline-flex items-center gap-1 px-2.5 py-1 rounded-full bg-white/70 text-gray-700 font-medium">
                        {#if insight.metric_type.startsWith('question_')}
                          <MessageSquare class="w-3 h-3" />
                          Question
                        {:else}
                          <BarChart3 class="w-3 h-3" />
                          Metric
                        {/if}
                      </span>
                      <span
                        class="inline-flex items-center px-2.5 py-1 rounded-full bg-white/70 font-bold {Math.abs(
                          insight.change
                        ) > 20
                          ? 'text-red-700'
                          : 'text-amber-700'}">
                        {Math.abs(insight.change).toFixed(1)}% change
                      </span>
                    </div>
                  </div>
                </div>
              </div>
            {/each}
          </div>
        </div>
      {/if}
    {/if}
  </div>
</div>

<style>
  .chart-container {
    min-height: 200px;
  }

  .chart-container :global(svg) {
    max-width: 100%;
    height: auto;
  }

  .chart-container :global(.plot-title) {
    font-size: 16px;
    font-weight: 600;
    fill: #1f2937;
  }
</style>
