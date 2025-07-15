import type { LayoutLoad } from './$types';
import { auth } from '$lib/stores/auth';
import { subscription } from '$lib/stores/subscription';
import { get } from 'svelte/store';
import { browser } from '$app/environment';

export const load: LayoutLoad = async ({ route }) => {
	// Get the current auth state
	const authState = get(auth);
	
	// Skip subscription loading for auth pages and public pages
	const isAuthPage = route?.id?.includes('login') || route?.id?.includes('register');
	const isPublicPage = route?.id?.includes('qr/');
	
	// Load subscription data if authenticated and not on auth/public pages
	if (browser && authState.isAuthenticated && !isAuthPage && !isPublicPage) {
		try {
			// Load subscription data before any page load functions run
			await Promise.all([
				subscription.fetchSubscription(),
				subscription.fetchPlanFeatures()
			]);
		} catch (error) {
			console.error('Failed to load subscription data:', error);
		}
	}
	
	return {
		user: authState.user
	};
};