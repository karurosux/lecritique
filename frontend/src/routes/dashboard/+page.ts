import { requireActiveSubscription } from '$lib/subscription/route-guard';
import { requireAuth } from '$lib/utils/auth-guard';
import { browser } from '$app/environment';

export async function load() {
	if (browser) {
		requireAuth();
		requireActiveSubscription();
	}
	
	return {};
}