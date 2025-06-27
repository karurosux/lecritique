import type { PageLoad } from './$types';
import { getApiClient } from '$lib/api';
import { browser } from '$app/environment';

export const load: PageLoad = async ({ params, parent }) => {
	const { restaurant } = await parent();

	// On server, return empty QR codes and defer to client
	if (!browser || !restaurant) {
		return {
			qrCodes: []
		};
	}

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
		return {
			qrCodes: []
		};
	}
};