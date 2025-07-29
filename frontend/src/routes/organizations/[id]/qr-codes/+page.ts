import type { PageLoad } from './$types';
import { getApiClient } from '$lib/api';
import { browser } from '$app/environment';
import { requireActiveSubscription } from '$lib/subscription/route-guard';
import { requireAuth } from '$lib/utils/auth-guard';

export const load: PageLoad = async ({ params, parent }) => {
  if (browser) {
    requireAuth();
    requireActiveSubscription();
  }

  const { organization } = await parent();

  // On server, return empty QR codes and defer to client
  if (!browser || !organization) {
    return {
      qrCodes: [],
    };
  }

  try {
    // Fetch QR codes for this organization
    const api = getApiClient();
    const response = await api.api.v1OrganizationsQrCodesList(params.id);

    let qrCodes = [];
    if (response.data.success && response.data.data) {
      qrCodes = response.data.data;
    }

    return {
      qrCodes,
    };
  } catch (err) {
    console.error('Error loading QR codes:', err);
    return {
      qrCodes: [],
    };
  }
};
