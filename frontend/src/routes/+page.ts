import type { PageLoad } from './$types';
import { getServerSideApiClient } from '$lib/api/client';

export const load: PageLoad = async ({ fetch }) => {
  try {
    const api = getServerSideApiClient(fetch);

    const response = await api.api.v1PlansList();

    if (response.data.success && response.data.data) {
      return {
        plans: response.data.data,
      };
    }

    return {
      plans: [],
    };
  } catch (error) {
    console.error('Failed to fetch plans:', error);
    return {
      plans: [],
    };
  }
};
