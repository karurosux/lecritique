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
		Grid3x3,
		Table,
		ArrowUpDown,
		ChevronRight
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

			let questionType = '';
			try {
				if (comparison?.metadata) {
					const metadata = JSON.parse(comparison.metadata);
					questionType = metadata.question_type || '';
				}
			} catch (e) {}

			if (questionType === 'yes_no') {
				return value.toFixed(1) + '% yes';
			} else if (questionType === 'rating') {
				if (isAverage) {
					return value.toFixed(1) + '/5';
				} else {
					if (periodData?.average) {
						return periodData.average.toFixed(1) + '/5';
					}
					return value.toFixed(1) + '/5';
				}
			} else if (questionType === 'scale') {
				let scaleValue = value;
				if (!isAverage && periodData?.average) {
					scaleValue = periodData.average;
				}

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
				if (periodData?.most_popular_choice) {
					return `"${periodData.most_popular_choice.choice}" (${periodData.most_popular_choice.count})`;
				}
				if (includeUnits) {
					return `${value.toFixed(0)} responses`;
				}
				return value.toFixed(0);
			} else if (questionType === 'multi_choice') {
				if (periodData?.top_choices && periodData.top_choices.length > 0) {
					const topChoicesText = periodData.top_choices
						.slice(0, 3)
						.map((choice: { choice: any; count: any }) => `${choice.choice} (${choice.count})`)
						.join(', ');
					return topChoicesText;
				}
				if (includeUnits) {
					return `${value.toFixed(0)} responses`;
				}
				return value.toFixed(0);
			} else if (questionType === 'text') {
				const getSentimentLabel = (score: number) => {
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
				const lowerName = metricName?.toLowerCase() || '';

				if (lowerName.includes('recommend') || lowerName.includes('likelihood')) {
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
					return isAverage ? value.toFixed(2) + ' avg' : value.toFixed(0) + ' total';
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
			const getSentimentLabel = (score: number) => {
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
			year: 'numeric'
		});
	}

	function toggleMetricFilter(metricType: string) {
		if (filteredMetricTypes.includes(metricType)) {
			filteredMetricTypes = filteredMetricTypes.filter((t) => t !== metricType);
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

	let comparisons = $derived(
		sortComparisons(
			filteredMetricTypes.length > 0
				? allComparisons.filter((comp: any) =>
						filteredMetricTypes.some((type) => comp.metric_type.includes(type))
					)
				: allComparisons
		)
	);

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

		chartInstances.forEach((chart) => {
			if (chart) {
				chart.destroy();
			}
		});
		chartInstances = [];

		chartContainer.innerHTML = '';

		if (comparisons.length === 0) {
			chartContainer.innerHTML =
				'<div class="text-center py-8 text-gray-500">No comparison data for chart</div>';
			return;
		}

		try {
			const containerWidth = chartContainer.offsetWidth || 800;
			const chartWidth = containerWidth - 48;

			groupedComparisons().forEach((questions, questionType) => {
				const sectionHeader = document.createElement('div');
				sectionHeader.className = 'flex items-center gap-3 mb-6 mt-8';

				const leftLine = document.createElement('div');
				leftLine.className = 'flex-1 h-px bg-gradient-to-r from-transparent to-gray-200';
				sectionHeader.appendChild(leftLine);

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

				const rightLine = document.createElement('div');
				rightLine.className = 'flex-1 h-px bg-gradient-to-l from-transparent to-gray-200';
				sectionHeader.appendChild(rightLine);

				chartContainer.appendChild(sectionHeader);

				const gridContainer = document.createElement('div');
				gridContainer.className = 'grid grid-cols-1 gap-6 mb-8';

				questions.forEach((comparison: { metric_name: string }) => {
					const chartCard = document.createElement('div');
					chartCard.className = 'group relative';

					const bgDiv = document.createElement('div');
					bgDiv.className =
						'absolute inset-0 bg-gradient-to-br from-gray-50/50 to-white/50 rounded-xl opacity-0 group-hover:opacity-100 transition-opacity duration-300';
					chartCard.appendChild(bgDiv);

					const contentDiv = document.createElement('div');
					contentDiv.className = 'relative p-6';

					const chartTitle = document.createElement('h5');
					chartTitle.className = 'text-sm font-semibold text-gray-700 mb-4 flex items-center gap-2';
					chartTitle.innerHTML = `
            <div class="h-1.5 w-1.5 bg-blue-500 rounded-full"></div>
            <span class="truncate">${comparison.metric_name}</span>
          `;
					chartTitle.title = comparison.metric_name;
					contentDiv.appendChild(chartTitle);

					const plotContainer = document.createElement('div');
					plotContainer.className = 'flex justify-center';

					const fullChartWidth = chartWidth - 48;
					const chart = createChartForQuestionType(comparison, questionType, fullChartWidth);
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

	function createChartForQuestionType(comparison: any, questionType: string, width: number) {
		const canvas = document.createElement('canvas');
		canvas.style.width = width + 'px';
		canvas.style.height = '300px';
		canvas.style.maxWidth = '100%';

		try {
			if (questionType === 'single_choice') {
				return createDoughnutChart(canvas, comparison, 'Single Choice Comparison');
			} else if (questionType === 'multi_choice') {
				return createPolarChart(canvas, comparison, 'Multiple Choice Comparison');
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
			console.error(`Error creating chart for question type ${questionType}:`, error);
			const ctx = canvas.getContext('2d');
			if (ctx) {
				ctx.fillStyle = '#EF4444';
				ctx.font = '14px sans-serif';
				ctx.textAlign = 'center';
				ctx.fillText('Error rendering chart', canvas.width / 2, canvas.height / 2);
			}
			return canvas;
		}
	}

	function createDoughnutChart(canvas: HTMLCanvasElement, comparison: any, title: string) {
		const choiceData = getChoiceDistributionData(comparison);

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
							'#F97316A0'
						],
						borderColor: ['#3B82F6', '#EF4444', '#F59E0B', '#10B981', '#8B5CF6', '#F97316'],
						borderWidth: 2
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
							'#F97316B0'
						],
						borderColor: ['#3B82F6', '#EF4444', '#F59E0B', '#10B981', '#8B5CF6', '#F97316'],
						borderWidth: 2
					}
				]
			},
			options: {
				responsive: true,
				maintainAspectRatio: false,
				plugins: {
					title: {
						display: true,
						text: formatChartTitle(title, period1Info, period2Info),
						font: { size: 12 }
					},
					legend: { position: 'bottom' },
					tooltip: {
						callbacks: {
							label: (context) =>
								`${context.dataset.label} - ${context.label}: ${context.parsed} selections`
						}
					}
				}
			}
		});

		chartInstances.push(chart);
		return canvas;
	}

	function createPolarChart(canvas: HTMLCanvasElement, comparison: any, title: string) {
		const choiceData = getChoiceDistributionData(comparison);

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
						categoryPercentage: 0.9
					},
					{
						label: 'Period 2',
						data: choiceData.period2Data,
						backgroundColor: '#10B98180',
						borderColor: '#10B981',
						borderWidth: 2,
						maxBarThickness: 60,
						barPercentage: 0.8,
						categoryPercentage: 0.9
					}
				]
			},
			options: {
				responsive: true,
				maintainAspectRatio: false,
				plugins: {
					title: {
						display: true,
						text: formatChartTitle(title, period1Info, period2Info),
						font: { size: 12 }
					},
					legend: { position: 'bottom' },
					tooltip: {
						callbacks: {
							label: (context) => `${context.dataset.label}: ${context.parsed.y} selections`
						}
					}
				},
				scales: {
					y: {
						beginAtZero: true,
						grid: { display: true }
					},
					x: {
						grid: { display: false }
					}
				}
			}
		});

		chartInstances.push(chart);
		return canvas;
	}

	function createPieChart(canvas: HTMLCanvasElement, comparison: any, title: string) {
		console.log('Yes/No chart data:', comparison);
		let chartData;

		if (comparison.period1.choice_distribution || comparison.period2.choice_distribution) {
			console.log('Using choice distribution for yes/no');
			chartData = getChoiceDistributionData(comparison);
		} else {
			console.log('Using fallback data for yes/no');
			chartData = {
				labels: ['Period 1', 'Period 2'],
				period1Data: [comparison.period1.average || comparison.period1.value || 0],
				period2Data: [comparison.period2.average || comparison.period2.value || 0]
			};
		}

		console.log('Final yes/no chart data:', chartData);

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

		const chart = new Chart(canvas, {
			type: 'bar',
			data: {
				labels: ['Period 1', 'Period 2'],
				datasets: [
					{
						label: 'Value',
						data: [
							comparison.period1.average || comparison.period1.value || 0,
							comparison.period2.average || comparison.period2.value || 0
						],
						backgroundColor: ['#3B82F680', '#10B98180'],
						borderColor: ['#3B82F6', '#10B981'],
						borderWidth: 2,
						maxBarThickness: 80,
						barPercentage: 0.7,
						categoryPercentage: 0.8
					}
				]
			},
			options: {
				responsive: true,
				maintainAspectRatio: false,
				plugins: {
					title: {
						display: true,
						text: formatChartTitle(title, period1Info, period2Info),
						font: { size: 12 }
					},
					legend: { display: false },
					tooltip: {
						callbacks: {
							label: (context) => {
								const period = context.label;
								const value = context.parsed.y;
								const formattedValue = formatValue(
									value,
									comparison.metric_type,
									true,
									comparison.metric_name,
									false,
									comparison,
									period === 'Period 1' ? comparison.period1 : comparison.period2
								);
								return `${period}: ${formattedValue.split('\n')[0]}`;
							}
						}
					}
				},
				scales: {
					y: {
						beginAtZero: true,
						grid: { display: true }
					},
					x: {
						grid: { display: false }
					}
				}
			}
		});

		chartInstances.push(chart);
		return canvas;
	}

	function createBarChart(canvas: HTMLCanvasElement, comparison: any, title: string) {
		const valueData = [
			comparison.period1.average || comparison.period1.value,
			comparison.period2.average || comparison.period2.value
		];

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

		const isSentiment =
			comparison.metric_type.includes('sentiment') ||
			comparison.metric_type.includes('text') ||
			(Math.abs(valueData[0]) <= 1 && Math.abs(valueData[1]) <= 1);

		const isCount =
			comparison.metric_type === 'survey_responses' || comparison.metric_type.includes('count');

		const chart = new Chart(canvas, {
			type: 'bar',
			data: {
				labels: ['Period 1', 'Period 2'],
				datasets: [
					{
						label: 'Value',
						data: valueData,
						backgroundColor: isSentiment
							? valueData.map((v) => (v >= 0 ? '#10B98180' : '#EF444480'))
							: ['#3B82F680', '#10B98180'],
						borderColor: isSentiment
							? valueData.map((v) => (v >= 0 ? '#10B981' : '#EF4444'))
							: ['#3B82F6', '#10B981'],
						borderWidth: 2,
						maxBarThickness: 80,
						barPercentage: 0.7,
						categoryPercentage: 0.8
					}
				]
			},
			options: {
				responsive: true,
				maintainAspectRatio: false,
				plugins: {
					title: {
						display: true,
						text: formatChartTitle(title, period1Info, period2Info),
						font: { size: 12 }
					},
					legend: { display: false },
					tooltip: {
						callbacks: {
							label: (context) => {
								const period = context.label;
								const value = context.parsed.y;
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
									period === 'Period 1' ? comparison.period1 : comparison.period2
								);
								return `${period}: ${formattedValue.split('\n')[0]}`;
							}
						}
					}
				},
				scales: {
					y: {
						beginAtZero: isCount || !isSentiment,
						grid: { display: true },
						ticks: {
							callback: function (value: any) {
								if (
									comparison.metric_type === 'survey_responses' ||
									comparison.metric_type.includes('count') ||
									comparison.metric_name?.toLowerCase().includes('responses')
								) {
									return Math.round(value).toLocaleString();
								}
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
							stepSize:
								comparison.metric_type === 'survey_responses' ||
								comparison.metric_type.includes('count')
									? 1
									: undefined
						}
					}
				}
			}
		});

		chartInstances.push(chart);
		return canvas;
	}

	function getChoiceDistributionData(comparison: any) {
		const allChoices = new Set();
		if (comparison.period1.choice_distribution) {
			Object.keys(comparison.period1.choice_distribution).forEach((choice) =>
				allChoices.add(choice)
			);
		}
		if (comparison.period2.choice_distribution) {
			Object.keys(comparison.period2.choice_distribution).forEach((choice) =>
				allChoices.add(choice)
			);
		}

		const labels = Array.from(allChoices) as string[];
		const period1Data = labels.map(
			(choice) => comparison.period1.choice_distribution?.[choice] || 0
		);
		const period2Data = labels.map(
			(choice) => comparison.period2.choice_distribution?.[choice] || 0
		);

		return { labels, period1Data, period2Data };
	}

	function formatChartTitle(title: string, period1Info: string, period2Info: string): string[] {
		const titleLines = [title];

		const period1Lines = period1Info.split('\n');
		const period2Lines = period2Info.split('\n');

		titleLines.push(`Period 1: ${period1Lines[0]}`);
		titleLines.push(`Period 2: ${period2Lines[0]}`);

		if (period1Lines.length > 1 && period1Lines[1]) {
			titleLines.push(period1Lines[1]);
		}

		return titleLines;
	}

	$effect(() => {
		if (mounted && viewMode === 'chart') {
			setTimeout(() => {
				renderCharts();
			}, 100);
		}
	});

	$effect(() => {
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
		chartInstances.forEach((chart) => {
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
				variant="inline"
			/>
		{:else}
			{#if data?.request}
				<div
					class="relative mb-8 overflow-hidden rounded-2xl bg-gradient-to-br from-indigo-50 via-purple-50 to-pink-50 p-6"
				>
					<div
						class="absolute inset-0 bg-gradient-to-br from-indigo-500/5 via-purple-500/5 to-pink-500/5"
					></div>
					<div class="relative">
						<div class="mb-4 flex items-center gap-3">
							<div
								class="flex h-10 w-10 items-center justify-center rounded-xl bg-gradient-to-br from-indigo-500 to-purple-600 shadow-lg shadow-indigo-500/25"
							>
								<Calendar class="h-5 w-5 text-white" />
							</div>
							<h3 class="text-xl font-bold text-gray-900">Period Comparison Analysis</h3>
						</div>

						<div class="grid grid-cols-1 gap-6 md:grid-cols-2">
							<div
								class="rounded-xl border border-white/50 bg-white/70 p-4 shadow-sm backdrop-blur-sm"
							>
								<div class="mb-2 flex items-center gap-2">
									<div
										class="flex h-6 w-6 items-center justify-center rounded-md bg-gradient-to-br from-blue-500 to-blue-600"
									>
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
								class="rounded-xl border border-white/50 bg-white/70 p-4 shadow-sm backdrop-blur-sm"
							>
								<div class="mb-2 flex items-center gap-2">
									<div
										class="flex h-6 w-6 items-center justify-center rounded-md bg-gradient-to-br from-purple-500 to-purple-600"
									>
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

			<div
				class="mb-8 rounded-2xl border border-gray-200/50 bg-white/80 p-5 shadow-lg backdrop-blur-sm"
			>
				<div class="flex flex-wrap items-center justify-between gap-4">
					<div class="flex items-center gap-3">
						<div class="flex items-center gap-2">
							<Maximize2 class="h-4 w-4 text-gray-500" />
							<span class="text-sm font-semibold text-gray-700">View Mode</span>
						</div>
						<div class="flex rounded-xl bg-gradient-to-r from-gray-100 to-gray-50 p-1 shadow-inner">
							<button
								class="rounded-lg px-4 py-2 text-sm font-medium transition-all duration-200 {viewMode ===
								'cards'
									? 'scale-105 transform bg-white text-gray-900 shadow-md'
									: 'text-gray-600 hover:bg-white/50 hover:text-gray-900'}"
								onclick={() => (viewMode = 'cards')}
							>
								<div class="flex items-center gap-2">
									<Grid3x3 class="h-4 w-4" />
									Cards
								</div>
							</button>
							<button
								class="rounded-lg px-4 py-2 text-sm font-medium transition-all duration-200 {viewMode ===
								'chart'
									? 'scale-105 transform bg-white text-gray-900 shadow-md'
									: 'text-gray-600 hover:bg-white/50 hover:text-gray-900'}"
								onclick={() => (viewMode = 'chart')}
							>
								<div class="flex items-center gap-2">
									<BarChart3 class="h-4 w-4" />
									Chart
								</div>
							</button>
							<button
								class="rounded-lg px-4 py-2 text-sm font-medium transition-all duration-200 {viewMode ===
								'table'
									? 'scale-105 transform bg-white text-gray-900 shadow-md'
									: 'text-gray-600 hover:bg-white/50 hover:text-gray-900'}"
								onclick={() => (viewMode = 'table')}
							>
								<div class="flex items-center gap-2">
									<Table class="h-4 w-4" />
									Table
								</div>
							</button>
						</div>
					</div>

					<div class="flex items-center gap-2">
						<span class="text-sm font-medium text-gray-700">Sort by:</span>
						<select
							bind:value={sortBy}
							class="rounded-md border border-gray-300 bg-white px-2 py-1 text-sm"
						>
							<option value="change">Change %</option>
							<option value="value">Current Value</option>
							<option value="name">Name</option>
						</select>
						<button
							class="rounded p-1 text-gray-600 transition-colors hover:text-blue-600"
							onclick={() => (sortOrder = sortOrder === 'asc' ? 'desc' : 'asc')}
							title="Toggle sort order"
						>
							<ArrowUpDown
								class="h-4 w-4 {sortOrder === 'desc' ? 'rotate-180' : ''} transition-transform"
							/>
						</button>
					</div>
				</div>
			</div>

			{#if viewMode === 'cards'}
				<div class="mb-8 grid grid-cols-1 gap-6 lg:grid-cols-2">
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
							(comparison.metric_type.includes('rating') && comparison.change_percent > 0) ||
							(comparison.metric_type.includes('recommend') && comparison.change_percent > 0)}
						{@const trendColorClass = isChoiceQuestion
							? 'from-blue-50 to-indigo-50 border-blue-200'
							: Math.abs(comparison.change_percent) < 5
								? 'from-gray-50 to-gray-100 border-gray-200'
								: isPositiveTrend
									? 'from-emerald-50 to-green-50 border-emerald-200'
									: 'from-red-50 to-rose-50 border-red-200'}

						<div
							class="group relative bg-gradient-to-br {trendColorClass} overflow-hidden rounded-2xl border-2 transition-all duration-300 hover:scale-[1.02] hover:shadow-xl"
							onmouseenter={() => (hoveredMetric = comparison.metric_name)}
							onmouseleave={() => (hoveredMetric = null)}
						>
							<div class="absolute inset-0 opacity-5">
								<svg class="h-full w-full" xmlns="http:">
									<pattern
										id="grid-{comparison.metric_type}"
										x="0"
										y="0"
										width="40"
										height="40"
										patternUnits="userSpaceOnUse"
									>
										<path
											d="M 40 0 L 0 0 0 40"
											fill="none"
											stroke="currentColor"
											stroke-width="1"
										/>
									</pattern>
									<rect width="100%" height="100%" fill="url(#grid-{comparison.metric_type})" />
								</svg>
							</div>

							<div class="relative p-6">
								<div class="mb-6 flex items-start justify-between">
									<div class="min-w-0 flex-1 pr-4">
										<h4
											class="mb-2 truncate text-lg leading-tight font-bold text-gray-900 transition-colors group-hover:text-indigo-700"
										>
											{comparison.metric_name}
										</h4>
										<div class="flex items-center gap-2">
											<span
												class="inline-flex items-center gap-1 rounded-full bg-white/80 px-3 py-1 text-xs font-medium whitespace-nowrap text-gray-700 shadow-sm"
											>
												{#snippet questionTypeIcon()}
													{@const questionType = (() => {
														if (comparison.metric_type.startsWith('question_')) {
															try {
																if (comparison.metadata) {
																	const metadata = JSON.parse(comparison.metadata);
																	return metadata.question_type || 'question';
																}
															} catch (e) {}
														}
														return 'metric';
													})()}

													{#if questionType === 'rating'}
														<Star class="h-3 w-3" />
													{:else if questionType === 'scale'}
														<BarChart2 class="h-3 w-3" />
													{:else if questionType === 'yes_no'}
														<Check class="h-3 w-3" />
													{:else if questionType === 'single_choice'}
														<Circle class="h-3 w-3" />
													{:else if questionType === 'multi_choice'}
														<CheckSquare class="h-3 w-3" />
													{:else if questionType === 'text'}
														<MessageSquare class="h-3 w-3" />
													{:else}
														<BarChart3 class="h-3 w-3" />
													{/if}
												{/snippet}
												{@render questionTypeIcon()}
												{(() => {
													if (comparison.metric_type.startsWith('question_')) {
														let questionType = 'Question';
														try {
															if (comparison.metadata) {
																const metadata = JSON.parse(comparison.metadata);
																questionType = metadata.question_type || 'Question';
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
												class="flex items-center gap-2 rounded-xl bg-white/90 px-3 py-2 shadow-md"
											>
												{#snippet trendIcon()}
													{@const TrendIcon = getTrendIcon(comparison.trend)}
													<TrendIcon
														class="h-5 w-5 {Math.abs(comparison.change_percent) < 5
															? 'text-gray-500'
															: isPositiveTrend
																? 'text-emerald-600'
																: 'text-red-600'}"
													/>
												{/snippet}
												{@render trendIcon()}
												<span
													class="text-lg font-bold {Math.abs(comparison.change_percent) < 5
														? 'text-gray-700'
														: isPositiveTrend
															? 'text-emerald-700'
															: 'text-red-700'}"
												>
													{comparison.change_percent > 0
														? '+'
														: ''}{comparison.change_percent.toFixed(1)}%
												</span>
											</div>
											<span
												class="text-xs font-medium {Math.abs(comparison.change_percent) < 5
													? 'text-gray-500'
													: isPositiveTrend
														? 'text-emerald-600'
														: 'text-red-600'}"
											>
												{Math.abs(comparison.change_percent) < 5
													? 'Stable'
													: isPositiveTrend
														? 'Improving'
														: 'Declining'}
											</span>
										</div>
									{/if}
								</div>

								{#if (() => {
									try {
										if (comparison.metadata) {
											const metadata = JSON.parse(comparison.metadata);
											return metadata.question_type === 'single_choice' || metadata.question_type === 'multi_choice';
										}
									} catch (e) {}
									return false;
								})()}
									<div class="mb-4 rounded-lg border border-blue-200/50 bg-blue-50/80 p-3">
										<div class="flex items-start gap-2">
											<Info class="mt-0.5 h-4 w-4 flex-shrink-0 text-blue-600" />
											<div class="text-xs text-blue-800">
												<span class="font-medium">Note:</span>
												{(() => {
													try {
														if (comparison.metadata) {
															const metadata = JSON.parse(comparison.metadata);
															if (metadata.question_type === 'single_choice') {
																return 'For single choice questions, this comparison shows the most popular choice selected in each period.';
															} else if (metadata.question_type === 'multi_choice') {
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

								<div class="mb-4 grid grid-cols-2 gap-4">
									<div class="rounded-xl bg-white/70 p-4 backdrop-blur-sm">
										<div class="mb-3 flex items-center gap-2">
											<div class="flex h-5 w-5 items-center justify-center rounded bg-blue-500">
												<span class="text-[10px] font-bold text-white">1</span>
											</div>
											<h5 class="text-sm font-semibold text-gray-900">Period 1</h5>
										</div>
										<div class="space-y-2">
											<div class="flex items-center justify-between">
												<span class="text-xs text-gray-600"
													>{(() => {
														try {
															if (comparison.metadata) {
																const metadata = JSON.parse(comparison.metadata);
																if (metadata.question_type === 'single_choice') return 'Top Choice';
																if (metadata.question_type === 'multi_choice')
																	return 'Top 3 Choices';
																return 'Value';
															}
														} catch (e) {}
														return 'Value';
													})()}</span
												>
												<div class="ml-2 flex items-center gap-1 text-right">
													<span class="text-xs font-bold whitespace-pre-line text-gray-900">
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
														<Star class="h-3 w-3 text-yellow-500" />
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
												<div class="flex items-center justify-between">
													<span class="text-xs text-gray-600">Average</span>
													<div class="ml-2 flex items-center gap-1 text-right">
														<span class="text-xs font-semibold whitespace-pre-line text-gray-800">
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
															<Star class="h-3 w-3 text-yellow-500" />
														{/if}
													</div>
												</div>
											{/if}
											<div class="flex items-center justify-between">
												<span class="text-xs text-gray-600">Responses</span>
												<span class="font-medium text-gray-700">
													{comparison.period1.count}
												</span>
											</div>
										</div>
									</div>

									<div class="rounded-xl bg-white/70 p-4 backdrop-blur-sm">
										<div class="mb-3 flex items-center gap-2">
											<div class="flex h-5 w-5 items-center justify-center rounded bg-purple-500">
												<span class="text-[10px] font-bold text-white">2</span>
											</div>
											<h5 class="text-sm font-semibold text-gray-900">Period 2</h5>
										</div>
										<div class="space-y-2">
											<div class="flex items-center justify-between">
												<span class="text-xs text-gray-600"
													>{(() => {
														try {
															if (comparison.metadata) {
																const metadata = JSON.parse(comparison.metadata);
																if (metadata.question_type === 'single_choice') return 'Top Choice';
																if (metadata.question_type === 'multi_choice')
																	return 'Top 3 Choices';
																return 'Value';
															}
														} catch (e) {}
														return 'Value';
													})()}</span
												>
												<div class="ml-2 flex items-center gap-1 text-right">
													<span class="text-xs font-bold whitespace-pre-line text-gray-900">
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
														<Star class="h-3 w-3 text-yellow-500" />
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
												<div class="flex items-center justify-between">
													<span class="text-xs text-gray-600">Average</span>
													<div class="ml-2 flex items-center gap-1 text-right">
														<span class="text-xs font-semibold whitespace-pre-line text-gray-800">
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
															<Star class="h-3 w-3 text-yellow-500" />
														{/if}
													</div>
												</div>
											{/if}
											<div class="flex items-center justify-between">
												<span class="text-xs text-gray-600">Responses</span>
												<span class="font-medium text-gray-700">
													{comparison.period2.count}
												</span>
											</div>
										</div>
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
									<div
										class="rounded-xl border border-white/50 bg-gradient-to-r from-white/50 to-white/30 p-4 backdrop-blur-sm"
									>
										<div class="flex items-center justify-between">
											<div class="flex items-center gap-2">
												<ArrowUpDown class="h-4 w-4 text-gray-600" />
												<span class="text-sm font-medium text-gray-700">Net Change</span>
											</div>
											<span
												class="text-lg font-bold {Math.abs(comparison.change_percent) < 5
													? 'text-gray-700'
													: isPositiveTrend
														? 'text-emerald-700'
														: 'text-red-700'}"
											>
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
				<div class="mb-8 overflow-hidden rounded-xl border border-gray-200 bg-white shadow-sm">
					<div class="overflow-x-auto">
						<table class="w-full">
							<thead class="border-b border-gray-200 bg-gray-50">
								<tr>
									<th
										class="px-4 py-3 text-left text-xs font-semibold tracking-wider text-gray-600 uppercase"
										>Metric</th
									>
									<th
										class="px-4 py-3 text-left text-xs font-semibold tracking-wider text-gray-600 uppercase"
										>Type</th
									>
									<th
										class="px-4 py-3 text-right text-xs font-semibold tracking-wider text-gray-600 uppercase"
										>{(() => {
											const hasChoiceQuestions = comparisons.some((comp) => {
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
										})()}</th
									>
									<th
										class="px-4 py-3 text-right text-xs font-semibold tracking-wider text-gray-600 uppercase"
										>{(() => {
											const hasChoiceQuestions = comparisons.some((comp) => {
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
										})()}</th
									>
									<th
										class="px-4 py-3 text-right text-xs font-semibold tracking-wider whitespace-nowrap text-gray-600 uppercase"
										>Change %</th
									>
									<th
										class="px-4 py-3 text-center text-xs font-semibold tracking-wider text-gray-600 uppercase"
										>Trend</th
									>
								</tr>
							</thead>
							<tbody class="divide-y divide-gray-200">
								{#each comparisons as comparison}
									<tr class="transition-colors hover:bg-gray-50">
										<td class="px-4 py-3">
											<div class="font-medium text-gray-900">
												{comparison.metric_name}
											</div>
										</td>
										<td class="px-4 py-3">
											<span
												class="inline-block max-w-24 truncate rounded-full bg-gray-100 px-2 py-1 text-xs whitespace-nowrap text-gray-700 capitalize"
											>
												{(() => {
													if (comparison.metric_type.startsWith('question_')) {
														let questionType = 'Question';
														try {
															if (comparison.metadata) {
																const metadata = JSON.parse(comparison.metadata);
																questionType = metadata.question_type || 'Question';
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
													)}"
												>
													{comparison.change_percent > 0
														? '+'
														: ''}{comparison.change_percent.toFixed(1)}%
												</span>
											{:else}
												<span class="text-sm text-gray-400">-</span>
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
														class="mx-auto h-4 w-4 {getTrendColor(
															comparison.trend,
															comparison.change_percent
														)}"
													/>
												{/snippet}
												{@render trendIcon()}
											{:else}
												<span class="text-sm text-gray-400">-</span>
											{/if}
										</td>
									</tr>
								{/each}
							</tbody>
						</table>
					</div>
				</div>
			{:else if viewMode === 'chart'}
				<div style="height: auto; min-height: 0;">
					<div class="mb-6 flex items-center justify-between">
						<h4 class="flex items-center gap-3 text-xl font-bold text-gray-900">
							<div
								class="flex h-10 w-10 items-center justify-center rounded-xl bg-gradient-to-br from-blue-500 to-indigo-600 shadow-lg shadow-blue-500/25"
							>
								<BarChart3 class="h-5 w-5 text-white" />
							</div>
							Visual Comparison
						</h4>
						<div
							class="flex items-center gap-2 rounded-lg bg-gray-50 px-3 py-1.5 text-sm text-gray-500"
						>
							<Activity class="h-4 w-4" />
							<span>{comparisons.length} metrics</span>
						</div>
					</div>

					<div
						bind:this={chartContainer}
						class="chart-container w-full"
						style="height: auto; overflow: visible;"
					></div>
				</div>
			{/if}

			{#if insights.length > 0}
				<div class="mt-8">
					<div class="mb-6 flex items-center gap-3">
						<div
							class="flex h-10 w-10 items-center justify-center rounded-xl bg-gradient-to-br from-amber-500 to-orange-600 shadow-lg shadow-amber-500/25"
						>
							<Info class="h-5 w-5 text-white" />
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
								class="relative bg-gradient-to-r {bgClass} rounded-xl border-2 p-5 shadow-sm transition-shadow duration-200 hover:shadow-md"
							>
								<div class="flex items-start gap-4">
									<div class="flex-shrink-0">
										<div
											class="h-10 w-10 bg-gradient-to-br {iconBgClass} flex items-center justify-center rounded-lg shadow-md"
										>
											{#snippet insightIcon()}
												{@const InsightIcon = getInsightIcon(insight.severity)}
												<InsightIcon class="h-5 w-5 text-white" />
											{/snippet}
											{@render insightIcon()}
										</div>
									</div>

									<div class="flex-1">
										<div class="mb-2 text-base leading-tight font-semibold text-gray-900">
											{insight.message}
										</div>

										{#if insight.recommendation}
											<div class="mt-3 rounded-lg bg-white/70 p-3">
												<div class="flex items-start gap-2">
													<ChevronRight class="mt-0.5 h-4 w-4 flex-shrink-0 text-gray-600" />
													<div class="text-sm text-gray-700">
														<span class="font-medium">Action:</span>
														{insight.recommendation}
													</div>
												</div>
											</div>
										{/if}

										<div class="mt-3 flex items-center gap-3 text-xs">
											<span
												class="inline-flex items-center gap-1 rounded-full bg-white/70 px-2.5 py-1 font-medium text-gray-700"
											>
												{#if insight.metric_type.startsWith('question_')}
													<MessageSquare class="h-3 w-3" />
													Question
												{:else}
													<BarChart3 class="h-3 w-3" />
													Metric
												{/if}
											</span>
											<span
												class="inline-flex items-center rounded-full bg-white/70 px-2.5 py-1 font-bold {Math.abs(
													insight.change
												) > 20
													? 'text-red-700'
													: 'text-amber-700'}"
											>
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
