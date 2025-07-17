import type { LayoutLoad } from './$types';
import { auth } from '$lib/stores/auth';
import { subscription } from '$lib/stores/subscription';
import { get } from 'svelte/store';
import { browser } from '$app/environment';

export const load: LayoutLoad = async ({ route }) => {
	// Get the current auth state
	const authState = get(auth);
	
	// Skip subscription loading for auth pages and public pages
	const isAuthPage = route?.id?.includes('login') || 
		route?.id?.includes('register') ||
		route?.id?.includes('forgot-password') ||
		route?.id?.includes('reset-password') ||
		route?.id?.includes('registration-success') ||
		route?.id?.includes('email-verification') ||
		route?.id?.includes('verify-email');
	const isPublicPage = route?.id?.includes('qr/');
	
	// If authenticated and not on auth/public pages, ensure subscription data is loaded
	if (browser && authState.isAuthenticated && !isAuthPage && !isPublicPage) {
		const subState = get(subscription);
		// Only fetch if we don't have subscription data (e.g., after page reload)
		if (!subState.subscription && !subState.isLoading) {
			await subscription.fetchSubscription();
		}
	}
	
	return {
		user: authState.user,
		auth: authState,
		subscription: get(subscription)
	};
};