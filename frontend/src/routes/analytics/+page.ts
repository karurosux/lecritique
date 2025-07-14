import { requireFeature } from '$lib/subscription/route-guard';
import { FEATURES } from '$lib/stores/subscription';

export async function load() {
	// Check if user has basic analytics feature
	requireFeature(FEATURES.BASIC_ANALYTICS);
	
	// Continue with normal page load
	return {};
}