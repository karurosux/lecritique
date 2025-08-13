<script lang="ts">
	import { Card, NoDataAvailable } from '$lib/components/ui';
	import {
		Star,
		BarChart3,
		ToggleLeft,
		ToggleRight,
		MessageSquare,
		ThumbsUp,
		ThumbsDown,
		Meh
	} from 'lucide-svelte';

	interface BackendChartData {
		question_id: string;
		question_text: string;
		question_type: string;
		chart_type: string;
		data: {
			scale?: number;
			distribution?: Record<string, number>;
			average?: number;
			total?: number;
			percentages?: Record<string, number>;
			options?: Record<string, number>;
			is_multi_choice?: boolean;
			combinations?: Array<{
				options: string[];
				count: number;
				percentage: number;
			}>;
			positive?: number;
			neutral?: number;
			negative?: number;
			samples?: string[];
			keywords?: string[];
		};
	}

	interface OrganizationChartData {
		organization_id: string;
		charts: BackendChartData[];
		summary: {
			total_responses: number;
			date_range: {
				start: string;
				end: string;
			};
			filters_applied: Record<string, any>;
		};
	}

	let {
		chartData = null,
		title = 'Analytics Dashboard',
		cardVariant = 'minimal'
	}: {
		chartData: OrganizationChartData | null;
		title?: string;
		cardVariant?: 'elevated' | 'default' | 'minimal';
	} = $props();

	let chartViewModes = $state(new Map<string, 'individual' | 'combinations'>());

	function getQuestionTypeLabel(type: string): string {
		switch (type) {
			case 'rating':
				return 'Rating';
			case 'scale':
				return 'Scale';
			case 'single_choice':
				return 'Single Choice';
			case 'multi_choice':
				return 'Multiple Choice';
			case 'yes_no':
				return 'Yes/No';
			case 'text':
				return 'Text Response';
			default:
				return 'Response';
		}
	}

	function getRatingStars(rating: number, scale: number = 5): string {
		const normalizedRating = scale === 10 ? rating / 2 : rating;
		const fullStars = Math.floor(normalizedRating);
		const hasHalfStar = normalizedRating % 1 >= 0.5;

		let stars = '';
		for (let i = 0; i < fullStars; i++) {
			stars += '★';
		}
		if (hasHalfStar) {
			stars += '☆';
		}
		while (stars.length < 5) {
			stars += '☆';
		}
		return stars;
	}

	function getColorForRating(rating: number, scale: number = 5): string {
		const normalizedRating = scale === 10 ? rating / 2 : rating;
		if (normalizedRating >= 4) return 'text-emerald-600';
		if (normalizedRating >= 3) return 'text-yellow-600';
		return 'text-red-600';
	}

	function getColorForChoice(index: number): string {
		const colors = [
			'bg-gradient-to-r from-blue-500 to-blue-600',
			'bg-gradient-to-r from-emerald-500 to-emerald-600',
			'bg-gradient-to-r from-purple-500 to-purple-600',
			'bg-gradient-to-r from-orange-500 to-orange-600',
			'bg-gradient-to-r from-pink-500 to-pink-600',
			'bg-gradient-to-r from-indigo-500 to-indigo-600',
			'bg-gradient-to-r from-teal-500 to-teal-600',
			'bg-gradient-to-r from-amber-500 to-amber-600'
		];
		return colors[index % colors.length];
	}

	function getColorForChoiceLight(index: number): string {
		const colors = [
			'bg-blue-50 border-blue-200',
			'bg-emerald-50 border-emerald-200',
			'bg-purple-50 border-purple-200',
			'bg-orange-50 border-orange-200',
			'bg-pink-50 border-pink-200',
			'bg-indigo-50 border-indigo-200',
			'bg-teal-50 border-teal-200',
			'bg-amber-50 border-amber-200'
		];
		return colors[index % colors.length];
	}

	function toggleViewMode(questionId: string) {
		const current = chartViewModes.get(questionId) || 'individual';
		chartViewModes.set(questionId, current === 'individual' ? 'combinations' : 'individual');
		chartViewModes = new Map(chartViewModes);
	}
</script>

{#if !chartData || !chartData.charts || chartData.charts.length === 0}
	<NoDataAvailable
		title="No Analytics Data Available"
		description="Start collecting customer feedback to unlock powerful insights and beautiful analytics visualizations."
		icon={BarChart3}
	/>
{:else}
	<div class="analytics-charts space-y-8">
		{#if title}
			<div class="mb-6">
				<div class="mb-3 flex items-center gap-3">
					<div
						class="flex h-8 w-8 items-center justify-center rounded-lg bg-gradient-to-br from-purple-500 to-pink-600"
					>
						<BarChart3 class="h-4 w-4 text-white" />
					</div>
					<div>
						<h2 class="text-xl font-semibold text-gray-900">{title}</h2>
						<p class="text-sm text-gray-600">
							{chartData.summary.total_responses.toLocaleString()} responses
							{#if chartData.summary.date_range.start}
								• {new Date(chartData.summary.date_range.start).toLocaleDateString()}
								to {new Date(chartData.summary.date_range.end).toLocaleDateString()}
							{/if}
						</p>
					</div>
				</div>
			</div>
		{/if}

		<div class="grid grid-cols-1 gap-8 lg:grid-cols-2">
			{#each chartData.charts as chart (chart.question_id)}
				<Card variant={cardVariant} class="group">
					<div class="chart-container p-6">
						<div class="mb-6 flex items-start justify-between">
							<div class="flex-1">
								<div class="mb-3 flex items-center gap-3">
									<div
										class="flex h-8 w-8 items-center justify-center rounded-lg bg-gradient-to-br from-blue-500 to-purple-600"
									>
										{#if chart.chart_type === 'rating'}
											<Star class="h-4 w-4 text-white" />
										{:else if chart.chart_type === 'scale'}
											<BarChart3 class="h-4 w-4 text-white" />
										{:else if chart.chart_type === 'text_sentiment'}
											<MessageSquare class="h-4 w-4 text-white" />
										{:else}
											<BarChart3 class="h-4 w-4 text-white" />
										{/if}
									</div>
									<div class="flex-1">
										<h3 class="mb-1 text-lg leading-tight font-semibold text-gray-900">
											{chart.question_text}
										</h3>
										<div class="flex items-center gap-2">
											<span
												class="rounded-md bg-gray-100 px-2 py-1 text-xs font-medium text-gray-700"
											>
												{getQuestionTypeLabel(chart.question_type)}
											</span>
											<span class="text-sm text-gray-500">
												{(chart.data.total || 0).toLocaleString()} responses
											</span>
										</div>
									</div>
								</div>
							</div>

							{#if chart.data.is_multi_choice && chart.data.combinations}
								<button
									onclick={() => toggleViewMode(chart.question_id)}
									class="flex items-center gap-2 rounded-lg border border-gray-200 bg-gray-50 px-3 py-1.5 text-sm font-medium text-gray-600 transition-all duration-200 hover:bg-gray-100 hover:text-gray-900"
								>
									<span class="text-xs">
										{chartViewModes.get(chart.question_id) === 'combinations'
											? 'Individual'
											: 'Combinations'}
									</span>
									{#if chartViewModes.get(chart.question_id) === 'combinations'}
										<ToggleRight class="h-3 w-3 text-emerald-600" />
									{:else}
										<ToggleLeft class="h-3 w-3 text-gray-400" />
									{/if}
								</button>
							{/if}
						</div>

						{#if chart.chart_type === 'rating'}
							<div class="space-y-4">
								{#if chart.data.average}
									<div
										class="relative mb-6 overflow-hidden rounded-2xl border border-amber-200 bg-gradient-to-br from-amber-50 via-orange-50 to-red-50 p-5 text-center"
									>
										<div class="absolute inset-0 opacity-10">
											<div class="absolute top-2 right-4 text-2xl">🍽️</div>
											<div class="absolute bottom-2 left-4 text-xl">✨</div>
											<div class="absolute top-4 left-8 text-lg">👨‍🍳</div>
										</div>
										<div class="relative z-10">
											<div class="mb-2 text-3xl font-black tracking-tight text-orange-700">
												{chart.data.average.toFixed(1)}
												<span class="ml-1 text-lg text-orange-600">/{chart.data.scale || 5}</span>
											</div>
											<div class="mb-3 text-2xl tracking-wider">
												{getRatingStars(chart.data.average, chart.data.scale)}
											</div>
											<div
												class="inline-flex items-center gap-2 rounded-full border border-orange-200 bg-white/80 px-3 py-1 backdrop-blur-sm"
											>
												<Star class="h-3 w-3 text-orange-600" />
												<span class="text-xs font-semibold text-orange-700">Average Rating</span>
											</div>
										</div>
									</div>
								{/if}

								{#if chart.data.distribution}
									<div class="space-y-3">
										{#each Object.entries(chart.data.distribution).sort(([a], [b]) => parseInt(b) - parseInt(a)) as [rating, count], index}
											{@const percentage = chart.data.percentages?.[rating] || 0}
											<div
												class="group rounded-xl border border-orange-100 p-3 transition-all duration-200 hover:border-orange-200 hover:bg-orange-50"
											>
												<div class="flex items-center gap-4">
													<div class="flex w-20 flex-shrink-0 items-center gap-3">
														<div
															class="flex h-8 w-8 items-center justify-center rounded-lg bg-gradient-to-br from-orange-100 to-orange-200 text-sm font-bold text-orange-800 transition-transform group-hover:scale-105"
														>
															{rating}
														</div>
														<Star
															class="h-4 w-4 {getColorForRating(
																parseInt(rating),
																chart.data.scale
															)} transition-transform group-hover:scale-105"
														/>
													</div>
													<div class="relative flex-1">
														<div
															class="h-6 overflow-hidden rounded-xl bg-gradient-to-r from-gray-200 to-gray-100 shadow-inner"
														>
															<div
																class="relative h-full rounded-xl bg-gradient-to-r from-orange-400 via-orange-500 to-orange-600 shadow-sm transition-all duration-700 ease-out"
																style="width: {percentage}%"
															>
																<div
																	class="absolute inset-0 rounded-xl bg-gradient-to-r from-transparent via-white/20 to-transparent"
																></div>
															</div>
														</div>

														{#if percentage > 20}
															<div
																class="pointer-events-none absolute inset-0 flex items-center justify-center text-xs font-bold text-white"
																style="text-shadow: 0 1px 2px rgba(0,0,0,0.5), 0 0 4px rgba(0,0,0,0.3);"
															>
																{percentage.toFixed(0)}%
															</div>
														{/if}
													</div>
													<div
														class="w-20 rounded-lg border border-orange-200 bg-white px-3 py-2 text-right text-sm font-bold text-gray-900 transition-colors group-hover:bg-orange-50"
													>
														{count.toLocaleString()}
														<span class="block text-xs font-normal text-orange-600">
															{percentage.toFixed(1)}%
														</span>
													</div>
												</div>
											</div>
										{/each}
									</div>
								{/if}
							</div>
						{:else if chart.chart_type === 'scale'}
							<div class="space-y-4">
								{#if chart.data.average}
									<div
										class="relative mb-6 overflow-hidden rounded-2xl border border-blue-200 bg-gradient-to-br from-blue-50 via-indigo-50 to-purple-50 p-5 text-center"
									>
										<div class="absolute inset-0 opacity-10">
											<div class="absolute top-2 right-4 text-2xl">📊</div>
											<div class="absolute bottom-2 left-4 text-xl">⚖️</div>
											<div class="absolute top-4 left-8 text-lg">📈</div>
										</div>
										<div class="relative z-10">
											<div class="mb-2 text-3xl font-black tracking-tight text-blue-700">
												{chart.data.average.toFixed(1)}
												<span class="ml-1 text-lg text-blue-600">/{chart.data.scale || 5}</span>
											</div>
											<div
												class="inline-flex items-center gap-2 rounded-full border border-blue-200 bg-white/80 px-3 py-1 backdrop-blur-sm"
											>
												<BarChart3 class="h-3 w-3 text-blue-600" />
												<span class="text-xs font-semibold text-blue-700">Average Score</span>
											</div>
										</div>
									</div>
								{/if}

								{#if chart.data.distribution}
									<div class="space-y-3">
										{#each Object.entries(chart.data.distribution).sort(([a], [b]) => parseInt(a) - parseInt(b)) as [value, count], index}
											{@const percentage = chart.data.percentages?.[value] || 0}
											<div
												class="group rounded-xl border border-blue-100 p-3 transition-all duration-200 hover:border-blue-200 hover:bg-blue-50"
											>
												<div class="flex items-center gap-4">
													<div class="flex w-20 flex-shrink-0 items-center gap-3">
														<div
															class="flex h-8 w-8 items-center justify-center rounded-lg bg-gradient-to-br from-blue-100 to-blue-200 text-sm font-bold text-blue-800 transition-transform group-hover:scale-105"
														>
															{value}
														</div>
														<BarChart3
															class="h-4 w-4 text-blue-600 transition-transform group-hover:scale-105"
														/>
													</div>
													<div class="relative flex-1">
														<div
															class="h-6 overflow-hidden rounded-xl bg-gradient-to-r from-gray-200 to-gray-100 shadow-inner"
														>
															<div
																class="relative h-full rounded-xl bg-gradient-to-r from-blue-400 via-blue-500 to-blue-600 shadow-sm transition-all duration-700 ease-out"
																style="width: {percentage}%"
															>
																<div
																	class="absolute inset-0 rounded-xl bg-gradient-to-r from-transparent via-white/20 to-transparent"
																></div>
															</div>
														</div>

														{#if percentage > 20}
															<div
																class="pointer-events-none absolute inset-0 flex items-center justify-center text-xs font-bold text-white"
																style="text-shadow: 0 1px 2px rgba(0,0,0,0.5), 0 0 4px rgba(0,0,0,0.3);"
															>
																{percentage.toFixed(0)}%
															</div>
														{/if}
													</div>
													<div
														class="w-20 rounded-lg border border-blue-200 bg-white px-3 py-2 text-right text-sm font-bold text-gray-900 transition-colors group-hover:bg-blue-50"
													>
														{count.toLocaleString()}
														<span class="block text-xs font-normal text-blue-600">
															{percentage.toFixed(1)}%
														</span>
													</div>
												</div>
											</div>
										{/each}
									</div>
								{/if}
							</div>
						{:else if chart.chart_type === 'choice'}
							{#if chart.data.is_multi_choice && chart.data.combinations && chartViewModes.get(chart.question_id) === 'combinations'}
								<div class="space-y-3">
									{#each chart.data.combinations.slice(0, 8) as combo, index}
										<div
											class="group rounded-xl border border-indigo-100 p-3 transition-all duration-200 hover:border-indigo-200 hover:bg-indigo-50"
										>
											<div class="flex items-center gap-4">
												<div
													class="h-6 w-6 rounded-lg {getColorForChoice(
														index
													)} flex flex-shrink-0 items-center justify-center shadow-sm transition-transform group-hover:scale-105"
												>
													<div class="h-2 w-2 rounded-full bg-white"></div>
												</div>
												<div class="min-w-0 flex-1">
													<span class="block truncate text-sm font-semibold text-gray-800">
														{combo.options.join(' + ')}
													</span>
													<span class="text-xs font-medium text-indigo-600">
														{combo.options.length} options combined
													</span>
												</div>
												<div class="flex items-center gap-4">
													<div class="relative w-28">
														<div
															class="h-5 overflow-hidden rounded-xl bg-gradient-to-r from-gray-200 to-gray-100 shadow-inner"
														>
															<div
																class="h-full {getColorForChoice(
																	index
																)} relative rounded-xl shadow-sm transition-all duration-700 ease-out"
																style="width: {combo.percentage}%"
															>
																<div
																	class="absolute inset-0 rounded-xl bg-gradient-to-r from-transparent via-white/30 to-transparent"
																></div>
															</div>
														</div>
														{#if combo.percentage > 15}
															<div
																class="pointer-events-none absolute inset-0 flex items-center justify-center text-xs font-bold text-white"
																style="text-shadow: 0 1px 2px rgba(0,0,0,0.5), 0 0 4px rgba(0,0,0,0.3);"
															>
																{combo.percentage.toFixed(0)}%
															</div>
														{/if}
													</div>
													<div
														class="w-20 rounded-lg border border-indigo-200 bg-white px-3 py-2 text-right text-sm font-bold text-gray-900 transition-colors group-hover:bg-indigo-50"
													>
														{combo.count.toLocaleString()}
														<span class="block text-xs font-normal text-indigo-600">
															{combo.percentage.toFixed(1)}%
														</span>
													</div>
												</div>
											</div>
										</div>
									{/each}
								</div>
							{:else}
								<div class="space-y-3">
									{#if chart.data.options}
										{#each Object.entries(chart.data.options).sort(([, a], [, b]) => b - a) as [option, count], index}
											{@const percentage = chart.data.percentages?.[option] || 0}
											<div
												class="group rounded-xl border border-purple-100 p-3 transition-all duration-200 hover:border-purple-200 hover:bg-purple-50"
											>
												<div class="flex items-center gap-4">
													<div
														class="h-6 w-6 rounded-lg {getColorForChoice(
															index
														)} flex flex-shrink-0 items-center justify-center shadow-sm transition-transform group-hover:scale-105"
													>
														<div class="h-2 w-2 rounded-full bg-white"></div>
													</div>
													<span class="flex-1 truncate text-sm font-semibold text-gray-800"
														>{option}</span
													>
													<div class="flex items-center gap-4">
														<div class="relative w-28">
															<div
																class="h-5 overflow-hidden rounded-xl bg-gradient-to-r from-gray-200 to-gray-100 shadow-inner"
															>
																<div
																	class="h-full {getColorForChoice(
																		index
																	)} relative rounded-xl shadow-sm transition-all duration-700 ease-out"
																	style="width: {percentage}%"
																>
																	<div
																		class="absolute inset-0 rounded-xl bg-gradient-to-r from-transparent via-white/30 to-transparent"
																	></div>
																</div>
															</div>
															{#if percentage > 15}
																<div
																	class="pointer-events-none absolute inset-0 flex items-center justify-center text-xs font-bold text-white"
																	style="text-shadow: 0 1px 2px rgba(0,0,0,0.5), 0 0 4px rgba(0,0,0,0.3);"
																>
																	{percentage.toFixed(0)}%
																</div>
															{/if}
														</div>
														<div
															class="w-20 rounded-lg border border-purple-200 bg-white px-3 py-2 text-right text-sm font-bold text-gray-900 transition-colors group-hover:bg-purple-50"
														>
															{count.toLocaleString()}
															<span class="block text-xs font-normal text-purple-600">
																{percentage.toFixed(1)}%
															</span>
														</div>
													</div>
												</div>
											</div>
										{/each}
									{/if}
								</div>
							{/if}
						{:else if chart.chart_type === 'text_sentiment'}
							<div class="space-y-4">
								<div class="mb-6 grid grid-cols-3 gap-3">
									<div
										class="group relative overflow-hidden rounded-xl border border-green-200 bg-gradient-to-br from-green-50 to-emerald-50 p-4 text-center transition-all duration-200 hover:shadow-lg"
									>
										<div class="absolute top-1 right-1 text-lg opacity-20">😋</div>
										<div
											class="mx-auto mb-3 flex h-10 w-10 items-center justify-center rounded-xl bg-gradient-to-br from-green-500 to-emerald-600 shadow-lg transition-transform group-hover:scale-110"
										>
											<ThumbsUp class="h-5 w-5 text-white" />
										</div>
										<div class="mb-1 text-2xl font-black text-green-700">
											{chart.data.positive || 0}
										</div>
										<div class="text-xs font-bold tracking-wide text-green-600 uppercase">
											Loved It
										</div>
									</div>
									<div
										class="group relative overflow-hidden rounded-xl border border-amber-200 bg-gradient-to-br from-amber-50 to-yellow-50 p-4 text-center transition-all duration-200 hover:shadow-lg"
									>
										<div class="absolute top-1 right-1 text-lg opacity-20">🤔</div>
										<div
											class="mx-auto mb-3 flex h-10 w-10 items-center justify-center rounded-xl bg-gradient-to-br from-amber-400 to-yellow-500 shadow-lg transition-transform group-hover:scale-110"
										>
											<Meh class="h-5 w-5 text-white" />
										</div>
										<div class="mb-1 text-2xl font-black text-amber-700">
											{chart.data.neutral || 0}
										</div>
										<div class="text-xs font-bold tracking-wide text-amber-600 uppercase">
											It's OK
										</div>
									</div>
									<div
										class="group relative overflow-hidden rounded-xl border border-red-200 bg-gradient-to-br from-red-50 to-rose-50 p-4 text-center transition-all duration-200 hover:shadow-lg"
									>
										<div class="absolute top-1 right-1 text-lg opacity-20">😞</div>
										<div
											class="mx-auto mb-3 flex h-10 w-10 items-center justify-center rounded-xl bg-gradient-to-br from-red-500 to-rose-600 shadow-lg transition-transform group-hover:scale-110"
										>
											<ThumbsDown class="h-5 w-5 text-white" />
										</div>
										<div class="mb-1 text-2xl font-black text-red-700">
											{chart.data.negative || 0}
										</div>
										<div class="text-xs font-bold tracking-wide text-red-600 uppercase">
											Not Great
										</div>
									</div>
								</div>

								{#if chart.data.samples && chart.data.samples.length > 0}
									<div class="rounded-lg border border-blue-100 bg-blue-50 p-4">
										<div class="mb-3 flex items-center gap-2">
											<MessageSquare class="h-4 w-4 text-blue-600" />
											<h4 class="text-sm font-semibold text-gray-900">Sample Responses</h4>
										</div>
										<div class="space-y-2">
											{#each chart.data.samples.slice(0, 3) as sample, index}
												<div class="rounded-lg border border-blue-100 bg-white p-3">
													<p class="text-sm text-gray-700 italic">
														"{sample}"
													</p>
												</div>
											{/each}
										</div>
									</div>
								{/if}

								{#if chart.data.keywords && chart.data.keywords.length > 0}
									<div class="rounded-lg border border-purple-100 bg-purple-50 p-4">
										<div class="mb-3 flex items-center gap-2">
											<span class="text-sm font-bold text-purple-600">#</span>
											<h4 class="text-sm font-semibold text-gray-900">Keywords</h4>
										</div>
										<div class="flex flex-wrap gap-2">
											{#each chart.data.keywords as keyword}
												<span
													class="rounded-full bg-purple-100 px-3 py-1 text-xs font-medium text-purple-800"
												>
													{keyword}
												</span>
											{/each}
										</div>
									</div>
								{/if}
							</div>
						{/if}
					</div>
				</Card>
			{/each}
		</div>
	</div>
{/if}
