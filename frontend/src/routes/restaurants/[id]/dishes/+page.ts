import { requireActiveSubscription } from '$lib/subscription/route-guard';

export async function load() {
	// Require active subscription to manage dishes
	requireActiveSubscription();
	
	return {};
}