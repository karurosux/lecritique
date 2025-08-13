<script lang="ts">
	import { onMount } from 'svelte';
	import { page } from '$app/stores';
	import { Card, Button, Input, Select } from '$lib/components/ui';
	import {
		MessageCircle,
		Download,
		Calendar,
		User,
		QrCode,
		Hash,
		MoreHorizontal,
		AlertTriangle,
		X,
		CheckCircle,
		Search,
		HelpCircle,
		Star,
		Book,
		Building2,
		ChevronDown,
		MessageSquare
	} from 'lucide-svelte';
	import { getApiClient, handleApiError } from '$lib/api/client';
	import { auth } from '$lib/stores/auth';
	import { goto } from '$app/navigation';
	import { RoleGate } from '$lib/components/auth';

	interface Feedback {
		id: string;
		customer_email?: string;
		rating: number;
		comment?: string;
		product_name?: string;
		organization_name?: string;
		location_name?: string;
		qr_code?: string;
		responses?: any[];
		created_at: string;
	}

	interface FeedbackFilters {
		organization_id?: string;
		location_id?: string;
		rating_min?: number;
		rating_max?: number;
		date_from?: string;
		date_to?: string;
		search?: string;
	}

	let loading = $state(true);
	let error = $state('');
	let feedback = $state<Feedback[]>([]);
	let organizations = $state<any[]>([]);
	let products = $state<any[]>([]);
	let totalCount = $state(0);
	let currentPage = $state(1);
	let itemsPerPage = $state(20);

	let filters: FeedbackFilters = {};
	let searchQuery = $state('');
	let searchInput = $state('');
	let selectedOrganization = $state('');
	let selectedProduct = $state('');
	let selectedRating = $state('');
	let dateFrom = $state('');
	let dateTo = $state('');

	let collapsedStates = $state<Record<string, boolean>>({});

	function toggleCollapse(feedbackId: string) {
		collapsedStates[feedbackId] =
			collapsedStates[feedbackId] === undefined ? false : !collapsedStates[feedbackId];
	}

	function isCollapsed(feedbackId: string): boolean {
		return collapsedStates[feedbackId] ?? true;
	}

	function setDefaultDateFilter() {
		const today = new Date();
		const fifteenDaysAgo = new Date(today);
		fifteenDaysAgo.setDate(today.getDate() - 15);

		dateFrom = fifteenDaysAgo.toISOString().split('T')[0];
		dateTo = today.toISOString().split('T')[0];
	}

	let isFirstLoad = $state(true);
	let searchTimeout: ReturnType<typeof setTimeout> | null = null;

	function handleSearchInput() {
		if (!isFirstLoad) {
			if (searchTimeout) {
				clearTimeout(searchTimeout);
			}

			searchTimeout = setTimeout(() => {
				searchQuery = searchInput;
				loadFeedback();
			}, 500);
		}
	}

	$effect(() => {
		if (
			!isFirstLoad &&
			(selectedRating !== undefined ||
				selectedOrganization !== undefined ||
				selectedProduct !== undefined ||
				dateFrom !== undefined ||
				dateTo !== undefined)
		) {
			loadFeedback();
		}
	});

	$effect(() => {
		if (!isFirstLoad && selectedOrganization !== undefined) {
			selectedProduct = '';
			loadProducts();
		}
	});

	let authState = $derived($auth);

	onMount(async () => {
		if (!authState.isAuthenticated) {
			goto('/login');
			return;
		}

		setDefaultDateFilter();

		const loadedOrganizations = await loadOrganizations();
		if (loadedOrganizations.length > 0) {
			await loadProducts();
		}
		await loadFeedback();

		isFirstLoad = false;
	});

	async function loadOrganizations() {
		try {
			const api = getApiClient();
			const response = await api.api.v1OrganizationsList();

			if (response.data.success && response.data.data) {
				organizations = response.data.data;
				return response.data.data;
			}
			return [];
		} catch (err) {
			console.error('Error loading organizations:', err);
			return [];
		}
	}

	async function loadProducts() {
		try {
			const api = getApiClient();

			if (organizations.length === 0 && !selectedOrganization) {
				return;
			}

			if (selectedOrganization) {
				const response = await api.api.v1OrganizationsProductsList(selectedOrganization);
				if (response.data.success && response.data.data) {
					products = response.data.data;
				}
			} else {
				const productPromises = organizations.map(async (organization) => {
					try {
						const response = await api.api.v1OrganizationsProductsList(organization.id);
						return response.data.data || [];
					} catch (err) {
						console.error(`Error loading products for organization ${organization.id}:`, err);
						return [];
					}
				});

				const productArrays = await Promise.all(productPromises);
				const allProducts = productArrays.flat();

				const uniqueProducts = allProducts.reduce((acc: any[], product: any) => {
					if (!acc.find((d) => d.name === product.name)) {
						acc.push(product);
					}
					return acc;
				}, []);

				products = uniqueProducts.sort((a, b) => a.name.localeCompare(b.name));
			}
		} catch (err) {
			console.error('Error loading products:', err);
		}
	}

	async function loadFeedback() {
		loading = true;
		error = '';

		try {
			const api = getApiClient();

			if (organizations.length === 0) {
				feedback = [];
				totalCount = 0;
				loading = false;
				return;
			}

			let allFeedback: Feedback[] = [];

			if (selectedOrganization) {
				const query: any = {};

				if (searchQuery) query.search = searchQuery;
				if (selectedRating) {
					query.rating_min = parseInt(selectedRating);
					query.rating_max = parseInt(selectedRating);
				}
				if (selectedProduct) query.product_id = selectedProduct;
				if (dateFrom) query.date_from = dateFrom;
				if (dateTo) query.date_to = dateTo;

				const feedbackResponse = await api.api.v1OrganizationsFeedbackList(
					selectedOrganization,
					query
				);
				const feedbackData = feedbackResponse.data;
				const organizationFeedback = feedbackData?.data || [];

				allFeedback = organizationFeedback.map((fb: any) => ({
					id: fb.id,
					customer_email: fb.customer_email,
					rating: fb.overall_rating,
					comment: fb.comment,
					product_name: fb.product?.name || null,
					organization_name: organizations.find((r) => r.id === selectedOrganization)?.name,
					location_name: fb.location_name,
					qr_code: fb.qr_code,
					responses: fb.responses,
					created_at: fb.created_at
				}));
			} else {
				const query: any = {};

				if (searchQuery) query.search = searchQuery;
				if (selectedRating) {
					query.rating_min = parseInt(selectedRating);
					query.rating_max = parseInt(selectedRating);
				}
				if (selectedProduct) query.product_id = selectedProduct;
				if (dateFrom) query.date_from = dateFrom;
				if (dateTo) query.date_to = dateTo;

				const feedbackPromises = organizations.map(async (organization) => {
					try {
						const feedbackResponse = await api.api.v1OrganizationsFeedbackList(
							organization.id,
							query
						);
						const feedbackData = feedbackResponse.data;
						const organizationFeedback = feedbackData?.data || [];

						return organizationFeedback.map((fb: any) => ({
							id: fb.id,
							customer_email: fb.customer_email,
							rating: fb.overall_rating,
							comment: fb.comment,
							product_name: fb.product?.name || null,
							organization_name: organization.name,
							location_name: fb.location_name,
							qr_code: fb.qr_code,
							responses: fb.responses,
							created_at: fb.created_at
						}));
					} catch (err) {
						console.error(`Error loading feedback for organization ${organization.id}:`, err);
						return [];
					}
				});

				const feedbackArrays = await Promise.all(feedbackPromises);
				allFeedback = feedbackArrays.flat();
			}

			allFeedback.sort(
				(a, b) => new Date(b.created_at).getTime() - new Date(a.created_at).getTime()
			);

			feedback = allFeedback;
			totalCount = allFeedback.length;
		} catch (err) {
			error = handleApiError(err);
		} finally {
			loading = false;
		}
	}

	function clearFilters() {
		searchQuery = '';
		searchInput = '';
		selectedOrganization = '';
		selectedProduct = '';
		selectedRating = '';

		setDefaultDateFilter();

		if (searchTimeout) {
			clearTimeout(searchTimeout);
		}

		loadFeedback();
	}

	function formatDate(dateString: string): string {
		try {
			const date = new Date(dateString);
			return (
				date.toLocaleDateString() +
				' ' +
				date.toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' })
			);
		} catch {
			return dateString;
		}
	}

	function getQuestionText(response: any, index: number): string {
		if (response.question_text && response.question_text.trim()) {
			return response.question_text;
		}

		if (response.question && response.question.trim()) {
			return response.question;
		}

		const answer = response.answer;

		if (typeof answer === 'boolean') {
			return `Question ${index + 1} (Yes/No)`;
		} else if (typeof answer === 'number' && answer >= 1 && answer <= 5) {
			return `Rating Question ${index + 1}`;
		} else if (typeof answer === 'number') {
			return `Numeric Question ${index + 1}`;
		} else if (Array.isArray(answer)) {
			return `Multiple Choice Question ${index + 1}`;
		} else if (typeof answer === 'string' && answer.length > 50) {
			return `Comment Question ${index + 1}`;
		} else {
			return `Question ${index + 1}`;
		}
	}

	function renderStars(rating: number): string {
		return '★'.repeat(rating) + '☆'.repeat(5 - rating);
	}

	function getRatingColor(rating: number): string {
		if (rating >= 4) return 'text-green-600';
		if (rating >= 3) return 'text-yellow-600';
		return 'text-red-600';
	}

	function getRatingBgColor(rating: number): string {
		if (rating >= 4) return 'bg-green-100 text-green-800';
		if (rating >= 3) return 'bg-yellow-100 text-yellow-800';
		return 'bg-red-100 text-red-800';
	}

	function exportToCSV() {
		const allQuestions = new Set<string>();
		feedback.forEach((fb) => {
			if (fb.responses) {
				fb.responses.forEach((response: any, index: number) => {
					const questionText = getQuestionText(response, index);
					const productName = fb.product_name || 'General';
					const prefixedQuestion = `[${productName}] ${questionText}`;
					allQuestions.add(prefixedQuestion);
				});
			}
		});

		const questionArray = Array.from(allQuestions).sort();
		const headers = [
			'ID',
			'Date',
			'Customer Email',
			'Rating',
			'Product',
			'Organization',
			'Location',
			'Comment',
			...questionArray
		];

		const csvContent = [
			headers.join(','),
			...feedback.map((fb) => {
				const basicFields = [
					fb.id,
					fb.created_at,
					fb.customer_email || 'Anonymous',
					fb.rating || '',
					fb.product_name || '',
					fb.organization_name || '',
					fb.qr_code || '',
					`"${fb.comment?.replace(/"/g, '""') || ''}"`
				];

				const responseMap = new Map<string, string>();
				if (fb.responses) {
					fb.responses.forEach((response: any, index: number) => {
						const questionText = getQuestionText(response, index);
						const productName = fb.product_name || 'General';
						const prefixedQuestion = `[${productName}] ${questionText}`;
						let answerText = '';

						if (typeof response.answer === 'boolean') {
							answerText = response.answer ? 'Yes' : 'No';
						} else if (typeof response.answer === 'number') {
							answerText = response.answer.toString();
						} else if (Array.isArray(response.answer)) {
							answerText = response.answer.join('; ');
						} else {
							answerText = String(response.answer || '');
						}

						responseMap.set(prefixedQuestion, answerText);
					});
				}

				const responseFields = questionArray.map((question) => {
					const answer = responseMap.get(question) || '';
					return `"${answer.replace(/"/g, '""')}"`;
				});

				return [...basicFields, ...responseFields].join(',');
			})
		].join('\n');

		const blob = new Blob([csvContent], { type: 'text/csv;charset=utf-8;' });
		const link = document.createElement('a');
		const url = URL.createObjectURL(blob);
		link.setAttribute('href', url);
		link.setAttribute('download', `feedback-export-${new Date().toISOString().split('T')[0]}.csv`);
		link.style.visibility = 'hidden';
		document.body.appendChild(link);
		link.click();
		document.body.removeChild(link);
	}
</script>

<svelte:head>
	<title>Feedback Management - Kyooar</title>
	<meta name="description" content="Manage and analyze customer feedback" />
</svelte:head>

<div class="mx-auto max-w-7xl px-4 py-8 sm:px-6 lg:px-8">
	<!-- Page Header -->
	<div class="mb-8">
		<div class="flex flex-col gap-6 lg:flex-row lg:items-center lg:justify-between">
			<div class="space-y-3">
				<div class="flex items-center space-x-3">
					<div
						class="flex h-12 w-12 items-center justify-center rounded-2xl bg-gradient-to-br from-blue-500 to-purple-600 shadow-lg shadow-blue-500/25"
					>
						<MessageCircle class="h-6 w-6 text-white" />
					</div>
					<div>
						<h1
							class="bg-gradient-to-r from-gray-900 to-gray-700 bg-clip-text text-3xl font-bold text-transparent"
						>
							Feedback Management
						</h1>
						<div class="mt-1 flex items-center space-x-4">
							<p class="font-medium text-gray-600">
								Review and analyze customer feedback from all your organizations
							</p>
							{#if !loading && feedback.length > 0}
								<div class="flex items-center space-x-3 text-sm">
									<div class="flex items-center space-x-1">
										<div class="h-2 w-2 rounded-full bg-blue-400"></div>
										<span class="text-gray-600">{totalCount} Total</span>
									</div>
									<div class="flex items-center space-x-1">
										<div class="h-2 w-2 rounded-full bg-purple-400"></div>
										<span class="text-gray-600"
											>{feedback.filter((f) => f.rating >= 4).length} Positive</span
										>
									</div>
								</div>
							{/if}
						</div>
					</div>
				</div>
			</div>

			<div class="flex items-center space-x-3">
				<!-- Export CSV Button -->
				<RoleGate roles={['OWNER', 'ADMIN', 'MANAGER']}>
					<Button
						variant="gradient"
						size="lg"
						class="group relative overflow-hidden shadow-lg transition-all duration-300 hover:shadow-xl"
						onclick={exportToCSV}
						disabled={feedback.length === 0}
					>
						<div
							class="absolute inset-0 bg-gradient-to-r from-blue-600 to-purple-600 opacity-0 transition-opacity duration-300 group-hover:opacity-100"
						></div>
						<Download
							class="relative z-10 mr-2 h-5 w-5 transition-transform duration-200 group-hover:scale-110"
						/>
						<span class="relative z-10">Export CSV</span>
					</Button>
				</RoleGate>
			</div>
		</div>
	</div>
	<!-- Filters -->
	<Card
		variant="default"
		hover
		interactive
		class="group animate-fade-in-up mb-6 transform transition-all duration-300"
	>
		<div class="space-y-4">
			<!-- First Row: Search and Organization -->
			<div class="grid grid-cols-1 gap-4 lg:grid-cols-3">
				<!-- Search -->
				<div class="space-y-1">
					<label class="text-xs font-medium tracking-wider text-gray-500 uppercase">Search</label>
					<Input
						type="text"
						placeholder="Search comments, emails, products..."
						bind:value={searchInput}
						oninput={handleSearchInput}
						class="w-full"
					/>
				</div>

				<!-- Organization Filter -->
				<div class="space-y-1">
					<label class="text-xs font-medium tracking-wider text-gray-500 uppercase"
						>Organization</label
					>
					<Select
						bind:value={selectedOrganization}
						options={[
							{ value: '', label: 'All Organizations' },
							...organizations.map((r) => ({ value: r.id, label: r.name }))
						]}
						minWidth="min-w-full"
					/>
				</div>

				<!-- Product Filter -->
				<div class="space-y-1">
					<label class="text-xs font-medium tracking-wider text-gray-500 uppercase">Product</label>
					<Select
						bind:value={selectedProduct}
						options={[
							{ value: '', label: 'All Products' },
							...products.map((d) => ({ value: d.id, label: d.name }))
						]}
						minWidth="min-w-full"
						disabled={products.length === 0}
					/>
				</div>
			</div>

			<!-- Second Row: Rating and Date Range -->
			<div class="grid grid-cols-1 gap-4 lg:grid-cols-2">
				<!-- Rating Filter -->
				<div class="space-y-1">
					<label class="text-xs font-medium tracking-wider text-gray-500 uppercase">Rating</label>
					<Select
						bind:value={selectedRating}
						options={[
							{ value: '', label: 'All Ratings' },
							{ value: '5', label: '5★ Excellent' },
							{ value: '4', label: '4★ Good' },
							{ value: '3', label: '3★ Average' },
							{ value: '2', label: '2★ Poor' },
							{ value: '1', label: '1★ Very Poor' }
						]}
						minWidth="min-w-full"
					/>
				</div>

				<!-- Date Range -->
				<div class="space-y-1">
					<label class="text-xs font-medium tracking-wider text-gray-500 uppercase"
						>Date Range</label
					>
					<div class="flex items-center gap-3">
						<div class="flex-1">
							<Input type="date" bind:value={dateFrom} placeholder="From" class="w-full" />
						</div>
						<span class="text-sm font-medium text-gray-500">to</span>
						<div class="flex-1">
							<Input type="date" bind:value={dateTo} placeholder="To" class="w-full" />
						</div>
					</div>
				</div>
			</div>

			<!-- Filter Summary and Actions -->
			<div
				class="flex flex-col gap-3 border-t border-gray-100 pt-4 sm:flex-row sm:items-center sm:justify-between"
			>
				<div class="flex items-center gap-4 text-sm">
					<div class="flex items-center gap-2">
						<Search class="h-4 w-4 text-gray-400" />
						<span class="font-medium text-gray-600">
							{feedback.length}
							{feedback.length === 1 ? 'entry' : 'entries'}
						</span>
					</div>
					{#if totalCount > feedback.length}
						<div class="h-4 w-px bg-gray-200"></div>
						<span class="text-gray-500">
							{totalCount} total
						</span>
					{/if}
				</div>
				<Button variant="outline" size="sm" onclick={clearFilters}>
					<X class="mr-2 h-4 w-4" />
					Clear Filters
				</Button>
			</div>
		</div>
	</Card>

	{#if loading}
		<!-- Loading State -->
		<div class="space-y-4">
			{#each Array(5) as _}
				<Card variant="default" class="opacity-50">
					<div class="animate-pulse">
						<div class="mb-4 flex items-start justify-between">
							<div class="space-y-2">
								<div class="h-4 w-32 rounded bg-gray-200"></div>
								<div class="h-4 w-48 rounded bg-gray-200"></div>
							</div>
							<div class="h-6 w-20 rounded bg-gray-200"></div>
						</div>
						<div class="h-16 rounded bg-gray-200"></div>
					</div>
				</Card>
			{/each}
		</div>
	{:else if error}
		<!-- Error State -->
		<Card variant="default" hover interactive class="group">
			<div class="py-12 text-center">
				<AlertTriangle class="mx-auto mb-4 h-12 w-12 text-red-500" />
				<h3 class="mb-2 text-lg font-medium text-gray-900">Failed to load feedback</h3>
				<p class="mb-4 text-gray-600">{error}</p>
				<Button onclick={loadFeedback}>Try Again</Button>
			</div>
		</Card>
	{:else if feedback.length === 0}
		<!-- Empty State -->
		<Card variant="default" hover interactive class="group">
			<div class="py-12 text-center">
				<MessageCircle class="mx-auto mb-4 h-12 w-12 text-gray-400" />
				<h3 class="mb-2 text-lg font-medium text-gray-900">No feedback available</h3>
				<p class="text-gray-600">There is no feedback to display at this time.</p>
			</div>
		</Card>
	{:else}
		<!-- Feedback List -->
		<div class="space-y-6">
			{#each feedback as fb, index}
				<div class="group animate-fade-in-up relative" style="animation-delay: {index * 50}ms">
					<!-- Modern card with gradient border on hover -->
					<div
						class="absolute -inset-0.5 rounded-2xl bg-gradient-to-r from-blue-500 to-purple-600 opacity-0 blur transition duration-500 group-hover:opacity-20"
					></div>
					<Card
						variant="default"
						class="relative overflow-hidden border-0 shadow-xl transition-all duration-500 hover:shadow-2xl"
					>
						<!-- Header Section -->
						<div class="relative">
							<!-- Background gradient accent -->
							<div
								class="absolute top-0 right-0 h-32 w-32 rounded-full bg-gradient-to-br from-blue-500/5 to-purple-600/5 blur-3xl"
							></div>

							<div
								class="relative mb-6 flex flex-col gap-4 sm:flex-row sm:items-start sm:justify-between"
							>
								<!-- Left side info -->
								<div class="flex-1 space-y-3">
									<!-- Rating and badges row -->
									<div class="flex flex-wrap items-center gap-3">
										<!-- Modern rating badge -->
										<div class="relative">
											<div
												class="absolute inset-0 bg-gradient-to-r {fb.rating >= 4
													? 'from-green-400 to-emerald-500'
													: fb.rating >= 3
														? 'from-yellow-400 to-orange-500'
														: 'from-red-400 to-pink-500'} rounded-xl opacity-25 blur"
											></div>
											<div
												class="relative flex items-center gap-2 bg-gradient-to-r px-4 py-2 {fb.rating >=
												4
													? 'border-green-200 from-green-50 to-emerald-50'
													: fb.rating >= 3
														? 'border-yellow-200 from-yellow-50 to-orange-50'
														: 'border-red-200 from-red-50 to-pink-50'} rounded-xl border"
											>
												<div class="flex text-lg">
													{#each Array(5) as _, i}
														{#if i < fb.rating}
															<Star
																class="h-5 w-5 {fb.rating >= 4
																	? 'text-green-500'
																	: fb.rating >= 3
																		? 'text-yellow-500'
																		: 'text-red-500'}"
																fill="currentColor"
															/>
														{:else}
															<Star class="h-5 w-5 text-gray-300" />
														{/if}
													{/each}
												</div>
												<span
													class="font-bold {fb.rating >= 4
														? 'text-green-700'
														: fb.rating >= 3
															? 'text-yellow-700'
															: 'text-red-700'}">{fb.rating}.0</span
												>
											</div>
										</div>

										{#if fb.product_name}
											<div
												class="flex items-center gap-2 rounded-lg border border-purple-200 bg-purple-50 px-3 py-1.5"
											>
												<Book class="h-4 w-4 text-purple-600" />
												<span class="text-sm font-medium text-purple-700">{fb.product_name}</span>
											</div>
										{/if}

										{#if fb.organization_name}
											<div
												class="flex items-center gap-2 rounded-lg border border-blue-200 bg-blue-50 px-3 py-1.5"
											>
												<Building2 class="h-4 w-4 text-blue-600" />
												<span class="text-sm font-medium text-blue-700">{fb.organization_name}</span
												>
											</div>
										{/if}
									</div>

									<!-- Meta information with modern icons -->
									<div class="flex flex-wrap items-center gap-4 text-sm text-gray-600">
										<div class="flex items-center gap-1.5">
											<Calendar class="h-4 w-4 text-gray-400" />
											<span>{formatDate(fb.created_at)}</span>
										</div>

										<div class="flex items-center gap-1.5">
											<User class="h-4 w-4 text-gray-400" />
											<span>{fb.customer_email || 'Anonymous Guest'}</span>
										</div>

										{#if fb.qr_code}
											<div class="flex items-center gap-1.5">
												<QrCode class="h-4 w-4 text-gray-400" />
												<span>{fb.qr_code}</span>
											</div>
										{/if}
									</div>
								</div>

								<!-- Right side actions -->
								<div class="flex items-center gap-2">
									<button
										class="rounded-lg p-2 text-gray-400 transition-colors hover:bg-gray-100 hover:text-gray-600"
									>
										<MoreHorizontal class="h-5 w-5" />
									</button>
								</div>
							</div>

							<!-- Comment Section -->
							{#if fb.comment}
								<div class="relative mb-6">
									<div
										class="absolute top-0 bottom-0 -left-1 w-1 rounded-full bg-gradient-to-b from-blue-500 to-purple-600"
									></div>
									<div class="pl-6">
										<div class="flex items-start gap-3">
											<MessageCircle class="mt-0.5 h-5 w-5 flex-shrink-0 text-gray-400" />
											<blockquote class="flex-1 leading-relaxed text-gray-700">
												<p class="text-base italic">"{fb.comment}"</p>
											</blockquote>
										</div>
									</div>
								</div>
							{/if}

							<!-- Questions & Responses Section -->
							{#if fb.responses && fb.responses.length > 0}
								<div class="space-y-4">
									<!-- Collapsible Header -->
									<button
										onclick={() => toggleCollapse(fb.id)}
										class="group/header -mx-3 flex w-full cursor-pointer items-center justify-between rounded-lg border border-transparent p-3 transition-all duration-200 hover:border-gray-200 hover:bg-gray-50/50 focus:ring-2 focus:ring-purple-500/20 focus:outline-none"
									>
										<div class="flex items-center gap-3">
											<div
												class="rounded-lg bg-gradient-to-br from-purple-500 to-blue-600 p-2 transition-shadow duration-200 group-hover/header:shadow-lg"
											>
												<HelpCircle class="h-5 w-5 text-white" />
											</div>
											<div class="text-left">
												<h4
													class="text-lg font-semibold text-gray-900 transition-colors duration-200 group-hover/header:text-purple-700"
												>
													Customer Responses
												</h4>
												<p
													class="text-sm text-gray-500 transition-colors duration-200 group-hover/header:text-gray-600"
												>
													{fb.responses.length} question{fb.responses.length > 1 ? 's' : ''} answered
												</p>
											</div>
										</div>

										<!-- Toggle Indicator -->
										<div
											class="flex items-center gap-2 rounded-lg border border-purple-200/50 px-3 py-2 text-sm font-medium text-purple-600 transition-all duration-200 group-hover/header:border-purple-300 group-hover/header:bg-purple-50 group-hover/header:text-purple-700"
										>
											<span class="hidden sm:block"
												>{isCollapsed(fb.id) ? 'Show' : 'Hide'} Details</span
											>
											<ChevronDown
												class="h-5 w-5 transform transition-transform duration-200 {isCollapsed(
													fb.id
												)
													? 'rotate-0'
													: 'rotate-180'}"
											/>
										</div>
									</button>

									<!-- Collapsed Preview -->
									{#if isCollapsed(fb.id)}
										<div class="ml-8 text-sm text-gray-400 italic">
											Click above to view {fb.responses.length} detailed response{fb.responses
												.length > 1
												? 's'
												: ''}
										</div>
									{/if}

									<!-- Collapsible Content -->
									{#if !isCollapsed(fb.id)}
										<div class="space-y-4" style="animation: slideDown 0.3s ease-out;">
											{#each fb.responses as response, index}
												<div class="group/question relative">
													<!-- Question Number Badge -->
													<div class="absolute top-0 -left-3 z-10">
														<div
															class="flex h-6 w-6 items-center justify-center rounded-full bg-gradient-to-r from-purple-500 to-blue-600"
														>
															<span class="text-xs font-bold text-white">{index + 1}</span>
														</div>
													</div>

													<!-- Question Card -->
													<div
														class="ml-6 rounded-xl border border-gray-200 bg-gradient-to-r from-gray-50 to-white p-5 transition-all duration-300 group-hover/question:shadow-md"
													>
														<!-- Question Text -->
														<div class="mb-3">
															<p class="text-sm leading-relaxed font-medium text-gray-900">
																{getQuestionText(response, index)}
															</p>
														</div>

														<!-- Answer Section -->
														<div class="space-y-2">
															<div class="flex items-start gap-3">
																<!-- Answer Type Icon -->
																<div class="mt-0.5">
																	{#if typeof response.answer === 'boolean'}
																		{#if response.answer}
																			<CheckCircle class="h-4 w-4 text-green-500" />
																		{:else}
																			<X class="h-4 w-4 text-red-500" />
																		{/if}
																	{:else if typeof response.answer === 'number'}
																		<Hash class="h-4 w-4 text-blue-500" />
																	{:else}
																		<MessageSquare class="h-4 w-4 text-purple-500" />
																	{/if}
																</div>

																<!-- Answer Display -->
																<div class="flex-1">
																	{#if typeof response.answer === 'boolean'}
																		<div
																			class="inline-flex items-center gap-2 rounded-lg px-3 py-1.5 {response.answer
																				? 'border border-green-200 bg-green-100 text-green-800'
																				: 'border border-red-200 bg-red-100 text-red-800'}"
																		>
																			<span class="font-semibold"
																				>{response.answer ? 'Yes' : 'No'}</span
																			>
																		</div>
																	{:else if typeof response.answer === 'number'}
																		{#if response.answer >= 1 && response.answer <= 5}
																			<!-- Rating display -->
																			<div class="flex items-center gap-2">
																				<div class="flex">
																					{#each Array(5) as _, i}
																						<Star
																							class="h-4 w-4 {i < response.answer
																								? 'text-yellow-400'
																								: 'text-gray-300'}"
																							fill="currentColor"
																						/>
																					{/each}
																				</div>
																				<span class="font-semibold text-gray-900"
																					>{response.answer}/5</span
																				>
																			</div>
																		{:else}
																			<!-- Numeric value -->
																			<div
																				class="inline-flex items-center gap-2 rounded-lg border border-blue-200 bg-blue-100 px-3 py-1.5 text-blue-800"
																			>
																				<span class="font-semibold">{response.answer}</span>
																			</div>
																		{/if}
																	{:else if Array.isArray(response.answer)}
																		<!-- Multiple choice -->
																		<div class="flex flex-wrap gap-2">
																			{#each response.answer as item}
																				<span
																					class="inline-flex items-center rounded-lg border border-purple-200 bg-purple-100 px-3 py-1.5 text-sm font-medium text-purple-800"
																				>
																					{item}
																				</span>
																			{/each}
																		</div>
																	{:else}
																		<!-- Text response -->
																		<div class="rounded-lg border border-gray-200 bg-gray-50 p-3">
																			<p class="text-sm leading-relaxed text-gray-900">
																				{response.answer}
																			</p>
																		</div>
																	{/if}
																</div>
															</div>
														</div>
													</div>
												</div>
											{/each}
										</div>
									{/if}
								</div>
							{/if}
						</div>
					</Card>
				</div>
			{/each}
		</div>
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

	@keyframes slideDown {
		from {
			opacity: 0;
			max-height: 0;
			transform: translateY(-10px);
		}
		to {
			opacity: 1;
			max-height: 1000px;
			transform: translateY(0);
		}
	}
</style>
