import type { PageLoad } from './$types';
import { getApiClient } from '$lib/api';
import { error } from '@sveltejs/kit';

export const load: PageLoad = async ({ params, parent }) => {
	const { restaurant } = await parent();

	try {
		// Fetch dishes for this restaurant
		const api = getApiClient();
		const response = await api.api.v1RestaurantsDishesList(params.id);
		
		let dishes = [];
		if (response.data.success && response.data.data) {
			dishes = response.data.data.map((dish: any) => ({
				id: dish.id || '',
				name: dish.name || '',
				description: dish.description || '',
				price: dish.price || 0,
				category: dish.category || 'Uncategorized',
				is_available: dish.is_available !== false,
				allergens: dish.allergens || [],
				preparation_time: dish.preparation_time || 0,
				created_at: dish.created_at || '',
				updated_at: dish.updated_at || ''
			}));
		}
		
		return {
			dishes
		};
	} catch (err) {
		console.error('Error loading dishes:', err);
		// Return empty array instead of throwing error to prevent crash
		return {
			dishes: []
		};
	}
};