import { requireFeature } from '$lib/subscription/route-guard';
import { FEATURES } from '$lib/stores/subscription';

export async function load() {
  // Check if user has feedback explorer feature
  requireFeature(FEATURES.FEEDBACK_EXPLORER);

  // Continue with normal page load
  return {};
}
