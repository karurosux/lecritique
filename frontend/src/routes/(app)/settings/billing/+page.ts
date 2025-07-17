import type { PageLoad } from './$types';
import { requireAuth } from '$lib/utils/auth-guards';

export const load: PageLoad = async (event) => {
	// Only owners can access billing settings
	await requireAuth(event, { 
		roles: ['OWNER'],
		requireOwner: true,
		redirectTo: '/settings'
	});
	
	// Load billing data here...
	return {
		// ... billing data
	};
};