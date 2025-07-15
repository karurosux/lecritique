import { requireActiveSubscription } from '$lib/subscription/route-guard';

export async function load() {
	// Dashboard requires an active subscription
	requireActiveSubscription();
	
	return {};
}