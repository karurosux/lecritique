import type { PageLoad } from './$types';
import { Api } from '$lib/api/api';
import { APP_CONFIG } from '$lib/constants/config';

export const load: PageLoad = async ({ fetch }) => {
  try {
    // Create an API instance without authentication for public endpoints
    const api = new Api({
      baseURL: APP_CONFIG.API_URL,
      customFetch: fetch,
    });

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
