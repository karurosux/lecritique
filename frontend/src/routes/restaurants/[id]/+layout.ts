import type { LayoutLoad } from './$types';
import { getApiClient } from '$lib/api';
import { error } from '@sveltejs/kit';

export const load: LayoutLoad = async ({ params, parent }) => {
	await parent();

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