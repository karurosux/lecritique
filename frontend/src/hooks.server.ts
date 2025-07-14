import type { Handle } from '@sveltejs/kit';
import { sequence } from '@sveltejs/kit/hooks';
import { FEATURES } from '$lib/stores/subscription';

// Define route protection configuration
const routeProtection = new Map([
	// Analytics routes
	['/analytics', { requireFeature: FEATURES.BASIC_ANALYTICS }],
	['/analytics/comparison', { requireFeature: FEATURES.ADVANCED_ANALYTICS }],
	['/analytics/grouped', { requireFeature: FEATURES.ADVANCED_ANALYTICS }],
	
	// Feedback routes
	['/feedback/manage', { requireFeature: FEATURES.FEEDBACK_EXPLORER }],
]);

const handleRouteProtection: Handle = async ({ event, resolve }) => {
	// Skip protection for non-protected routes
	const protection = routeProtection.get(event.url.pathname);
	if (!protection) {
		return resolve(event);
	}
	
	// Note: Server-side hooks can't access client stores directly
	// Route protection should be handled in +page.ts or +layout.ts files
	// This is just for documentation of protected routes
	
	return resolve(event);
};

export const handle = sequence(handleRouteProtection);