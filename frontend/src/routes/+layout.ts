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
	
	// Subscription data is now loaded during login, so we don't need to fetch it here
	// The subscription data from login response is already available in the subscription store
	
	return {
		user: authState.user
	};
};