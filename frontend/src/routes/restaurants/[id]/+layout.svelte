<script lang="ts">
	import type { LayoutData } from './$types';
	import { page } from '$app/stores';
	import { ArrowLeft, Store, UtensilsCrossed, QrCode } from 'lucide-svelte';
	import { Button } from '$lib/components/ui';
	import { goto } from '$app/navigation';
	import { getApiClient, isAuthenticated } from '$lib/api';
	import { browser } from '$app/environment';
	import { onMount } from 'svelte';

	let { data }: { data: LayoutData } = $props();

	let currentPath = $derived($page.url.pathname);
	let restaurantId = $derived(data.restaurantId);
	let restaurant = $state(data.restaurant);
	let loading = $state(!restaurant && browser);

	const navItems = [
		{ href: `/restaurants/${restaurantId}/dishes`, label: 'Menu', icon: UtensilsCrossed },
		{ href: `/restaurants/${restaurantId}/qr-codes`, label: 'QR Codes', icon: QrCode }
	];

	let activeItem = $derived(navItems.find(item => currentPath === item.href) || navItems[0]);

	// Fetch restaurant data on client if not available
	onMount(async () => {
		if (!restaurant && isAuthenticated()) {
			try {
				loading = true;
				const api = getApiClient();
				const response = await api.api.v1RestaurantsDetail(restaurantId);
				
				if (response.data.success && response.data.data) {
					restaurant = response.data.data;
				}
			} catch (error) {
				console.error('Error loading restaurant:', error);
				goto('/restaurants');
			} finally {
				loading = false;
			}
		} else if (!restaurant) {
			loading = false;
		}
	});
</script>

<div class="min-h-screen bg-gradient-to-br from-blue-50/50 via-purple-50/30 to-pink-50/50 relative">
	<!-- Animated background gradients -->
	<div class="absolute inset-0 overflow-hidden">
		<div class="absolute -top-40 -right-40 w-80 h-80 bg-purple-300 rounded-full mix-blend-multiply filter blur-3xl opacity-20 animate-blob"></div>
		<div class="absolute -bottom-40 -left-40 w-80 h-80 bg-blue-300 rounded-full mix-blend-multiply filter blur-3xl opacity-20 animate-blob animation-delay-2000"></div>
		<div class="absolute top-1/2 left-1/2 transform -translate-x-1/2 -translate-y-1/2 w-80 h-80 bg-pink-300 rounded-full mix-blend-multiply filter blur-3xl opacity-20 animate-blob animation-delay-4000"></div>
	</div>
	<!-- Clean floating header -->
	<div class="relative pt-4 pb-2">
		<div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
			<!-- Restaurant Header -->
			<div class="flex items-center justify-between mb-3">
				<div class="flex items-center gap-4">
					<Button
						variant="ghost"
						size="icon"
						onclick={() => goto('/restaurants')}
						class="bg-white/90 backdrop-blur-sm shadow-lg hover:bg-white transition-all"
					>
						<ArrowLeft class="h-4 w-4" />
					</Button>
					<div class="flex items-center gap-3">
						<div class="h-10 w-10 bg-gradient-to-br from-blue-500 to-purple-600 rounded-xl flex items-center justify-center shadow-lg shadow-blue-500/25">
							<Store class="h-5 w-5 text-white" />
						</div>
						<h1 class="text-2xl font-bold bg-gradient-to-r from-gray-900 to-gray-700 bg-clip-text text-transparent">
							{#if loading}
								Loading...
							{:else if restaurant}
								{restaurant.name}
							{:else}
								Restaurant
							{/if}
						</h1>
					</div>
				</div>
			</div>

			<!-- Floating Navigation Tabs -->
			<div class="flex justify-center mb-4">
				<div class="bg-white backdrop-blur-sm rounded-xl shadow-lg shadow-gray-900/10 border border-gray-200/50 p-1">
					<div class="flex space-x-1">
						{#each navItems as item}
							<a
								href={item.href}
								class="
									relative inline-flex items-center gap-2 px-4 py-2 rounded-lg text-sm font-medium transition-all duration-300 group
									{currentPath === item.href
										? 'bg-gradient-to-r from-blue-500 to-purple-600 text-white shadow-md shadow-blue-500/25'
										: 'text-gray-600 hover:text-gray-800 hover:bg-gray-50'}
								"
							>
								{#if currentPath === item.href}
									<div class="absolute inset-0 bg-gradient-to-r from-blue-600 to-purple-700 rounded-lg opacity-0 group-hover:opacity-100 transition-opacity duration-300"></div>
								{/if}
								<svelte:component this={item.icon} class="h-4 w-4 relative z-10" />
								<span class="relative z-10">{item.label}</span>
							</a>
						{/each}
					</div>
				</div>
			</div>
		</div>
	</div>

	<!-- Content -->
	<div class="relative">
		<div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
			<slot />
		</div>
	</div>
</div>