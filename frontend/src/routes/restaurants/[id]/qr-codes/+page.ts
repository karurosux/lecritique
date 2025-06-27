import type { PageLoad } from './$types';
import { getApiClient } from '$lib/api';
import { error } from '@sveltejs/kit';

export const load: PageLoad = async ({ params, parent }) => {
	const { restaurant } = await parent();

	try {
		// Fetch QR codes for this restaurant
		const api = getApiClient();
		const response = await api.api.v1RestaurantsQrCodesList(params.id);
		
		let qrCodes = [];
		if (response.data.success && response.data.data) {
			qrCodes = response.data.data;
		}
		
		return {
			qrCodes
		};
	} catch (err) {
		console.error('Error loading QR codes:', err);
		throw error(500, 'Failed to load QR codes');
	}
};