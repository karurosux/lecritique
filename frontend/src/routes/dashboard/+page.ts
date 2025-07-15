import { requireActiveSubscription } from '$lib/subscription/route-guard';
import { browser } from '$app/environment';
import { get } from 'svelte/store';
import { subscription } from '$lib/stores/subscription';
import { auth } from '$lib/stores/auth';

export async function load() {
	// Server-side: just return empty (subscription check happens client-side)
	if (!browser) {
		return {};
	}
	
	// Client-side: ensure subscription is loaded before checking
	const authState = get(auth);
	
	if (authState.isAuthenticated) {
		const subState = get(subscription);
		
		// If subscription not loaded yet, fetch it first
		if (!subState.subscription && !subState.isLoading) {
			await subscription.fetchSubscription();
		}
		
		// Now check if user has active subscription
		requireActiveSubscription();
	}
	
	return {};
}