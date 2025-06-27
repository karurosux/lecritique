<script lang="ts">
	import type { PageData } from './$types';
	import { Button, Card, Input, Select } from '$lib/components/ui';
	import { Plus } from 'lucide-svelte';
	import DishCard from '$lib/components/dishes/DishCard.svelte';
	import AddDishModal from '$lib/components/dishes/AddDishModal.svelte';
	import { getApiClient } from '$lib/api';
	import { toast } from 'svelte-sonner';
	import { invalidateAll } from '$app/navigation';

	let { data }: { data: PageData } = $props();

	let showAddDishModal = $state(false);
	let editingDish = $state<any>(null);
	let searchQuery = $state('');
	let categoryFilter = $state('all');
	let availabilityFilter = $state('all');
	let sortBy = $state('name');

	// Get restaurant and dishes from layout/page data
	let restaurant = $derived(data.restaurant);
	let dishes = $derived(data.dishes || []);

	// Get unique categories
	let categories = $derived(
		dishes.reduce((cats: string[], dish: any) => {
			if (dish.category && !cats.includes(dish.category)) {
				cats.push(dish.category);
			}
			return cats;
		}, []).sort()
	);

	// Filter and sort dishes
	let filteredDishes = $derived(
		dishes
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
			await invalidateAll();
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
			await invalidateAll();
		} catch (error) {
			toast.error('Failed to update dish availability');
			console.error(error);
		}
	}
</script>

<svelte:head>
	<title>Menu - {restaurant.name} | LeCritique</title>
</svelte:head>

<div class="space-y-6">
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
				{dishes.length === 0 ? 'No dishes yet' : 'No dishes match your filters'}
			</h3>
			<p class="text-gray-500 mb-4">
				{dishes.length === 0 
					? 'Start building your menu by adding your first dish'
					: 'Try adjusting your search or filters'
				}
			</p>
			{#if dishes.length === 0}
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
					onEdit={() => handleEditDish(dish)}
					onDelete={() => handleDeleteDish(dish)}
					onToggleAvailability={() => handleToggleAvailability(dish)}
				/>
			{/each}
		</div>
	{/if}
</div>

<!-- Add/Edit Dish Modal -->
<AddDishModal
	bind:isOpen={showAddDishModal}
	editingDish={editingDish}
	onclose={() => {
		showAddDishModal = false;
		editingDish = null;
	}}
	onsave={() => {
		showAddDishModal = false;
		editingDish = null;
		invalidateAll();
	}}
/>