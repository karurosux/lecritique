import type { PageLoad } from './$types';
import { requireAuth } from '$lib/utils/auth-guards';

export const load: PageLoad = async (event) => {
	// Temporarily allow all authenticated users for testing questionnaires
	await requireAuth(event, { 
		requireVerified: true
	});
	
	return {};
};