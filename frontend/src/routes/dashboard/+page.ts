import { requireActiveSubscription } from '$lib/subscription/route-guard';
import { browser } from '$app/environment';
import { get } from 'svelte/store';
import { subscription } from '$lib/stores/subscription';

export async function load() {
	// Only check subscription in browser to ensure store is properly hydrated
	if (browser) {
		// Debug: log current subscription state
		const subState = get(subscription);
		console.log('Dashboard load - subscription state:', {
			subscription: subState.subscription,
			status: subState.subscription?.status,
			isLoading: subState.isLoading
		});
		
		// Dashboard requires an active subscription
		requireActiveSubscription();
	}
	
	return {};
}