import { requireFeature } from '$lib/subscription/route-guard';
import { FEATURES } from '$lib/stores/subscription';

export async function load() {
  // Check if user has advanced analytics feature
  requireFeature(FEATURES.ADVANCED_ANALYTICS);

  // Continue with normal page load
  return {};
}
