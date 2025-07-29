import { requireNoSubscription } from '$lib/subscription/route-guard';
import { browser } from '$app/environment';

export async function load() {
  if (!browser) {
    return {};
  }

  requireNoSubscription();

  return {};
}
