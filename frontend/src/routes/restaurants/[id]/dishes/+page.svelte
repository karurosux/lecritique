<script lang="ts">
	import type { PageData } from './$types';
	import { Button, Card, Input, Select } from '$lib/components/ui';
	import { Plus } from 'lucide-svelte';
	import DishCard from '$lib/components/dishes/DishCard.svelte';
	import AddDishModal from '$lib/components/dishes/AddDishModal.svelte';
	import { getApiClient } from '$lib/api';
	import { toast } from 'svelte-sonner';
	import { onMount } from 'svelte';
	import { page } from '$app/stores';
	import { goto } from '$app/navigation';
	import { QuestionApi } from '$lib/api/question';

	let { data }: { data: PageData } = $props();

	let showAddDishModal = $state(false);
	let editingDish = $state<any>(null);
	let searchQuery = $state('');
	let categoryFilter = $state('all');
	let availabilityFilter = $state('all');
	let sortBy = $state('name');
	let dishes = $state(data.dishes || []);
	let loading = $state(false);
	let dishesWithQuestions = $state<string[]>([]);

	// Get restaurant from layout/page data
	let restaurant = $derived(data.restaurant);
	let restaurantId = $derived($page.params.id);

	// Fetch dishes when restaurant becomes available
	onMount(async () => {
		if (!dishes.length && restaurant) {
			await fetchDishes();
		}
		if (restaurant) {
			await fetchDishesWithQuestions();
		}
	});

	// Watch for restaurant changes and fetch data
	$effect(async () => {
		if (restaurant && !dishes.length) {
			await fetchDishes();
		}
	});

	async function fetchDishes() {
		try {
			loading = true;
			const api = getApiClient();
			const response = await api.api.v1RestaurantsDishesList(restaurantId);
			
			if (response.data.success && response.data.data) {
				dishes = response.data.data.map((dish: any) => ({
					id: dish.id || '',
					name: dish.name || '',
					description: dish.description || '',
					price: dish.price || 0,
					category: dish.category || 'Uncategorized',
					is_available: dish.is_available !== false,
					allergens: dish.allergens || [],
					preparation_time: dish.preparation_time || 0,
					created_at: dish.created_at || '',
					updated_at: dish.updated_at || ''
				}));
			}
		} catch (error) {
			console.error('Error loading dishes:', error);
		} finally {
			loading = false;
		}
	}
	
	// Enhance dishes with questions information
	let dishesWithQuestionnaires = $derived(
		dishes.map(dish => {
			const hasQuestions = dishesWithQuestions.includes(dish.id);
			return {
				...dish,
				has_questionnaire: hasQuestions
			};
		})
	);

	// Get unique categories
	let categories = $derived(
		dishesWithQuestionnaires.reduce((cats: string[], dish: any) => {
			if (dish.category && !cats.includes(dish.category)) {
				cats.push(dish.category);
			}
			return cats;
		}, []).sort()
	);

	// Filter and sort dishes
	let filteredDishes = $derived(
		dishesWithQuestionnaires
			.filter((dish: any) => {
				const matchesSearch = dish.name.toLowerCase().includes(searchQuery.toLowerCase()) ||
					dish.description?.toLowerCase().includes(searchQuery.toLowerCase());
				const matchesCategory = categoryFilter === 'all' || dish.category === categoryFilter;
				const matchesAvailability = availabilityFilter === 'all' ||
					(availabilityFilter === 'available' && dish.is_available) ||
					(availabilityFilter === 'unavailable' && !dish.is_available);
				return matchesSearch && matchesCategory && matchesAvailability;
			})
			.sort((a: any, b: any) => {
				switch (sortBy) {
					case 'price':
						return a.price - b.price;
					case 'category':
						return a.category.localeCompare(b.category);
					case 'created_at':
						return new Date(b.created_at).getTime() - new Date(a.created_at).getTime();
					default:
						return a.name.localeCompare(b.name);
				}
			})
	);

	async function handleAddDish() {
		editingDish = null;
		showAddDishModal = true;
	}

	async function handleEditDish(dish: any) {
		editingDish = dish;
		showAddDishModal = true;
	}

	async function handleDeleteDish(dish: any) {
		if (!confirm(`Are you sure you want to delete "${dish.name}"?`)) {
			return;
		}

		try {
			const api = getApiClient();
			await api.api.v1DishesDelete(dish.id);
			toast.success('Dish deleted successfully');
			await fetchDishes();
		} catch (error) {
			toast.error('Failed to delete dish');
			console.error(error);
		}
	}

	async function handleToggleAvailability(dish: any) {
		try {
			const api = getApiClient();
			await api.api.v1DishesUpdate(dish.id, {
				...dish,
				is_available: !dish.is_available
			});
			toast.success(`${dish.name} ${dish.is_available ? 'disabled' : 'enabled'}`);
			await fetchDishes();
		} catch (error) {
			toast.error('Failed to update dish availability');
			console.error(error);
		}
	}

	async function handleManageQuestionnaire(dish: any) {
		// Navigate to questionnaire page
		goto(`/restaurants/${restaurantId}/questionnaire/${dish.id}`);
	}
	

	async function fetchDishesWithQuestions() {
		try {
			const api = getApiClient();
			const response = await api.api.v1RestaurantsQuestionsDishesWithQuestionsList(restaurantId);
			dishesWithQuestions = response.data.data || [];
			console.log('Dishes with questions:', dishesWithQuestions);
		} catch (error) {
			console.error('Failed to fetch dishes with questions:', error);
		}
	}
</script>

<svelte:head>
	<title>Menu - {restaurant?.name || 'Restaurant'} | LeCritique</title>
</svelte:head>

{#if !restaurant}
	<div class="space-y-6">
		<div class="text-center">
			<p class="text-gray-600">Loading restaurant...</p>
		</div>
	</div>
{:else}
	<div class="space-y-6">
		<!-- Loading State -->
		{#if loading}
			<div class="text-center">
				<p class="text-gray-600">Loading dishes...</p>
			</div>
		{:else}
	<!-- Search and Filters -->
	<Card variant="glass">
		<div class="flex flex-wrap gap-4 items-end">
			<div class="flex-1 min-w-64">
				<Input
					type="text"
					placeholder="Search dishes..."
					bind:value={searchQuery}
					variant="default"
				/>
			</div>
			
			<Select
				bind:value={categoryFilter}
				options={[
					{ value: 'all', label: 'All Categories' },
					...categories.map(cat => ({ value: cat, label: cat }))
				]}
			/>

			<Select
				bind:value={availabilityFilter}
				options={[
					{ value: 'all', label: 'All Status' },
					{ value: 'available', label: 'Available' },
					{ value: 'unavailable', label: 'Unavailable' }
				]}
			/>

			<Select
				bind:value={sortBy}
				options={[
					{ value: 'name', label: 'Sort by Name' },
					{ value: 'price', label: 'Sort by Price' },
					{ value: 'category', label: 'Sort by Category' },
					{ value: 'created_at', label: 'Sort by Date' }
				]}
			/>

			<Button onclick={handleAddDish} variant="gradient" size="lg" class="gap-2">
				<Plus class="h-4 w-4" />
				Add Dish
			</Button>
		</div>
	</Card>

	<!-- Dishes Grid -->
	{#if filteredDishes.length === 0}
		<div class="text-center py-12">
			<div class="w-24 h-24 mx-auto bg-gray-100 rounded-full flex items-center justify-center mb-4">
				<Plus class="h-8 w-8 text-gray-400" />
			</div>
			<h3 class="text-lg font-semibold mb-2">
				{dishesWithQuestionnaires.length === 0 ? 'No dishes yet' : 'No dishes match your filters'}
			</h3>
			<p class="text-gray-500 mb-4">
				{dishesWithQuestionnaires.length === 0 
					? 'Start building your menu by adding your first dish'
					: 'Try adjusting your search or filters'
				}
			</p>
			{#if dishesWithQuestionnaires.length === 0}
				<Button onclick={handleAddDish}>
					<Plus class="mr-2 h-4 w-4" />
					Add First Dish
				</Button>
			{/if}
		</div>
	{:else}
		<div class="grid gap-6 md:grid-cols-2 lg:grid-cols-3">
			{#each filteredDishes as dish (dish.id)}
				<DishCard
					{dish}
					onedit={() => handleEditDish(dish)}
					ondelete={() => handleDeleteDish(dish)}
					ontoggleavailability={() => handleToggleAvailability(dish)}
					ongeneratequestionnaire={() => handleManageQuestionnaire(dish)}
				/>
			{/each}
		</div>
	{/if}
	{/if}
</div>
{/if}

<!-- Add/Edit Dish Modal -->
{#if restaurant}
	<AddDishModal
		bind:isOpen={showAddDishModal}
		editingDish={editingDish}
		onclose={() => {
			showAddDishModal = false;
			editingDish = null;
		}}
		onsave={async (dishData) => {
			try {
				const api = getApiClient();
				if (editingDish) {
					// Update existing dish
					await api.api.v1DishesUpdate(editingDish.id, {
						...dishData,
						restaurant_id: restaurantId
					});
					toast.success('Dish updated successfully');
				} else {
					// Create new dish
					await api.api.v1DishesCreate({
						...dishData,
						restaurant_id: restaurantId
					});
					toast.success('Dish created successfully');
				}
				showAddDishModal = false;
				editingDish = null;
				await fetchDishes();
				await fetchDishesWithQuestions();
			} catch (error) {
				console.error('Error saving dish:', error);
				toast.error('Failed to save dish');
			}
		}}
	/>
{/if}

