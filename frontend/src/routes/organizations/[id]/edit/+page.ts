import type { PageLoad } from './$types';
import { requireAuth } from '$lib/utils/auth-guards';

export const load: PageLoad = async event => {
  // Only OWNER, ADMIN, and MANAGER can edit organizations
  await requireAuth(event, {
    roles: ['OWNER', 'ADMIN', 'MANAGER'],
    requireVerified: true,
  });

  return {};
};
