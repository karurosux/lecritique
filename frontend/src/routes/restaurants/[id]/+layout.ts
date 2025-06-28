import type { LayoutLoad } from './$types';
import { getApiClient, isAuthenticated, getAuthToken } from '$lib/api';
import { error, redirect } from '@sveltejs/kit';
import { browser } from '$app/environment';

export const load: LayoutLoad = async ({ params, parent }) => {
	await parent();

	// On server, skip authentication check and defer to client
	if (!browser) {
		return {
			restaurant: null,
			restaurantId: params.id
		};
	}

	// Check authentication on client side
	if (!isAuthenticated() || !getAuthToken()) {
		throw redirect(302, '/login');
	}

	try {
		const api = getApiClient();
		const response = await api.api.v1RestaurantsDetail(params.id);
		
		if (!response.data.success || !response.data.data) {
			throw error(404, 'Restaurant not found');
		}
		
		return {
			restaurant: response.data.data,
			restaurantId: params.id
		};
	} catch (err) {
		console.error('Error loading restaurant:', err);
		throw error(404, 'Restaurant not found');
	}
};