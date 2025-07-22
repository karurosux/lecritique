<script lang="ts">
	import { subscription, isSubscribed } from '$lib/stores/subscription';
	import { PlanSelector } from '$lib/components/subscription';
	import { onMount } from 'svelte';
	import { goto } from '$app/navigation';
	import { APP_CONFIG } from '$lib/constants/config';
	import type { ModelsSubscriptionPlan } from '$lib/api/api';
	import { Check, HelpCircle, User, Frown } from 'lucide-svelte';

	let plans = $state([]);
	let isLoading = $state(true);
	let isProcessing = $state(false);

	onMount(async () => {
		try {
			// Check if user already has subscription from login data
			if ($isSubscribed) {
				goto('/settings');
				return;
			}

			// If not subscribed, fetch available plans
			await subscription.fetchPlans();
			const sub = $subscription;
			plans = sub.plans || [];
		} catch (error) {
			console.error('Failed to load plans:', error);
		} finally {
			isLoading = false;
		}
	});

	async function handleSelectPlan(plan: ModelsSubscriptionPlan) {
		// For pricing page, we redirect to registration with the plan
		goto(`/register?plan=${plan.code}`);
	}
</script>

<svelte:head>
	<title>Choose Your Plan - Kyooar</title>
	<meta name="description" content="Select the perfect Kyooar plan for your organization feedback needs." />
</svelte:head>

<div class="min-h-screen relative bg-gray-50">
	<!-- Light gradient background -->
	<div class="fixed inset-0 bg-gradient-to-br from-blue-50 via-white to-purple-50"></div>
	
	<!-- Animated gradient orbs with light colors -->
	<div class="fixed inset-0 overflow-hidden">
		<div class="absolute -top-[40%] -left-[20%] w-[60%] h-[60%] bg-gradient-to-br from-purple-400 to-pink-400 rounded-full mix-blend-multiply filter blur-3xl opacity-10 animate-pulse"></div>
		<div class="absolute -bottom-[40%] -right-[20%] w-[60%] h-[60%] bg-gradient-to-br from-blue-400 to-cyan-400 rounded-full mix-blend-multiply filter blur-3xl opacity-10 animate-pulse animation-delay-2000"></div>
		<div class="absolute top-[20%] right-[30%] w-[40%] h-[40%] bg-gradient-to-br from-indigo-400 to-purple-400 rounded-full mix-blend-multiply filter blur-3xl opacity-8 animate-pulse animation-delay-4000"></div>
	</div>

	<!-- Subtle grid pattern -->
	<div class="fixed inset-0 opacity-[0.03]">
		<div class="absolute inset-0" style="background-image: repeating-linear-gradient(0deg, transparent, transparent 59px, rgba(0,0,0,1) 59px, rgba(0,0,0,1) 60px), repeating-linear-gradient(90deg, transparent, transparent 59px, rgba(0,0,0,1) 59px, rgba(0,0,0,1) 60px);"></div>
	</div>

	<div class="relative z-10">
		<div class="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-20">
			<!-- Modern Header -->
			<div class="text-center mb-16">
				<h1 class="text-5xl sm:text-6xl font-bold text-gray-900 mb-4 tracking-tight">
					Choose Your
					<span class="block text-transparent bg-clip-text bg-gradient-to-r from-purple-600 to-blue-600">
						Perfect Plan
					</span>
				</h1>
				<p class="text-xl text-gray-600 max-w-2xl mx-auto">
					Start your journey with the plan that fits your organization best
				</p>
			</div>

			<!-- Plan Selection -->
			{#if isLoading}
				<div class="flex justify-center items-center py-24">
					<div class="animate-spin rounded-full h-12 w-12 border-4 border-purple-600 border-t-transparent"></div>
				</div>
			{:else if plans.length > 0}
				<div class="relative">
					<!-- Soft glow effect behind plans -->
					<div class="absolute inset-0 bg-gradient-to-r from-purple-300/20 via-blue-300/20 to-purple-300/20 blur-3xl"></div>
					<div class="relative">
						<PlanSelector
							{plans}
							isLoading={isProcessing}
							onSelectPlan={handleSelectPlan}
							actionLabel={() => 'Get Started'}
							showCurrentBadge={false}
						/>
					</div>
				</div>
				
				<!-- Trust Indicators -->
				<div class="mt-20 text-center">
					<div class="inline-flex items-center justify-center space-x-8 flex-wrap gap-y-4 bg-white/60 backdrop-blur-sm border border-gray-200 rounded-2xl px-8 py-6 shadow-lg">
						<div class="flex items-center text-sm text-gray-700">
							<div class="w-8 h-8 rounded-full bg-green-100 flex items-center justify-center mr-3">
								<Check class="w-4 h-4 text-green-600" />
							</div>
							Start collecting feedback instantly
						</div>
						<div class="flex items-center text-sm text-gray-700">
							<div class="w-8 h-8 rounded-full bg-green-100 flex items-center justify-center mr-3">
								<Check class="w-4 h-4 text-green-600" />
							</div>
							Cancel anytime
						</div>
						<div class="flex items-center text-sm text-gray-700">
							<div class="w-8 h-8 rounded-full bg-green-100 flex items-center justify-center mr-3">
								<Check class="w-4 h-4 text-green-600" />
							</div>
							No setup fees
						</div>
					</div>
				</div>
				
				<!-- Help Section -->
				<div class="mt-16 bg-white/80 backdrop-blur-md border border-gray-200 rounded-3xl p-8 max-w-2xl mx-auto shadow-xl">
					<div class="text-center">
						<div class="inline-flex items-center justify-center w-14 h-14 bg-gradient-to-br from-blue-500 to-purple-600 rounded-2xl mb-4 shadow-lg">
							<HelpCircle class="w-7 h-7 text-white" />
						</div>
						<h3 class="font-semibold text-gray-900 mb-2">Questions?</h3>
						<p class="text-sm text-gray-600">
							Get help choosing the right plan<br/>
							<a href={`mailto:${APP_CONFIG.emails.support}`} class="text-purple-600 hover:text-purple-700 font-medium transition-colors">{APP_CONFIG.emails.support}</a>
						</p>
					</div>
				</div>
			{:else}
				<div class="text-center py-24 bg-white/80 backdrop-blur-md border border-gray-200 rounded-3xl shadow-xl">
					<Frown class="w-16 h-16 text-gray-400 mx-auto mb-4" />
					<p class="text-gray-700 mb-2">No plans available at the moment.</p>
					<p class="text-sm text-gray-500">Please check back later or contact support.</p>
				</div>
			{/if}
		</div>
	</div>
</div>

<style>
	/* Styles moved to global app.css */
</style>

