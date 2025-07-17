import type { PageLoad } from './$types';
import { requireAuth } from '$lib/utils/auth-guards';

export const load: PageLoad = async (event) => {
	// All authenticated users can access settings
	await requireAuth(event, { 
		requireVerified: true
	});
	
	return {};
};