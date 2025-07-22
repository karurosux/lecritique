import { writable, derived } from 'svelte/store';
import { browser } from '$app/environment';
import { getApiClient } from '$lib/api/client';
import type { ModelsSubscription, ModelsSubscriptionPlan } from '$lib/api/api';
import { auth } from './auth';
import { hasFeatureFromToken, getLimitFromToken } from '$lib/utils/jwt';


interface SubscriptionState {
  subscription: ModelsSubscription | null;
  plans: ModelsSubscriptionPlan[];
  usage: SubscriptionUsage | null;
  isLoading: boolean;
  error: string | null;
}

interface SubscriptionUsage {
  feedbacks_count: number;
  organizations_count: number;
  locations_count: number;
  qr_codes_count: number;
  team_members_count: number;
  period_start: string;
  period_end: string;
}

function createSubscriptionStore() {
  const { subscribe, set, update } = writable<SubscriptionState>({
    subscription: null,
    plans: [],
    usage: null,
    isLoading: false,
    error: null
  });

  return {
    subscribe,

    async fetchSubscription() {
      update(state => ({ ...state, isLoading: true, error: null }));

      try {
        const api = getApiClient();
        const response = await api.api.v1UserSubscriptionList();

        if (response.data.success && response.data.data) {
          update(state => ({
            ...state,
            subscription: response.data.data,
            isLoading: false
          }));
        } else {
          throw new Error('Failed to fetch subscription');
        }
      } catch (error: any) {
        update(state => ({
          ...state,
          isLoading: false,
          error: error.message || 'Failed to fetch subscription'
        }));
      }
    },

    async fetchPlans() {
      update(state => ({ ...state, isLoading: true, error: null }));

      try {
        const api = getApiClient();
        const response = await api.api.v1PlansList();

        if (response.data.success && response.data.data) {
          update(state => ({
            ...state,
            plans: response.data.data,
            isLoading: false
          }));
        } else {
          throw new Error('Failed to fetch plans');
        }
      } catch (error: any) {
        update(state => ({
          ...state,
          isLoading: false,
          error: error.message || 'Failed to fetch plans'
        }));
      }
    },

    async fetchUsage() {
      update(state => ({ ...state, isLoading: true, error: null }));

      try {
        const api = getApiClient();
        const response = await api.api.v1UserSubscriptionUsageList();

        if (response.data.success && response.data.data) {
          update(state => ({
            ...state,
            usage: response.data.data,
            isLoading: false
          }));
        } else {
          throw new Error('Failed to fetch usage');
        }
      } catch (error: any) {
        update(state => ({
          ...state,
          isLoading: false,
          error: error.message || 'Failed to fetch usage'
        }));
      }
    },


    async createCheckoutSession(planId: string) {
      // TODO: Replace with actual API call when endpoints are available
      console.warn('Payment endpoints not yet available in API client');

      // Mock response for development
      return {
        session_id: 'mock_session_123',
        checkout_url: 'https://checkout.stripe.com/mock'
      };
    },

    async createPortalSession() {
      // TODO: Replace with actual API call when endpoints are available
      console.warn('Payment endpoints not yet available in API client');

      // Mock response for development
      return {
        portal_url: 'https://billing.stripe.com/mock'
      };
    },

    async checkPermission(resourceType: string) {
      try {
        const api = getApiClient();

        switch (resourceType) {
          case 'organization':
            const response = await api.api.v1UserCanCreateOrganizationList();
            if (response.data.success && response.data.data) {
              return response.data.data;
            }
            break;
          // Add other resource types as needed
        }

        throw new Error('Failed to check permission');
      } catch (error: any) {
        throw new Error(error.message || 'Failed to check permission');
      }
    },

    reset() {
      set({
        subscription: null,
        plans: [],
        usage: null,
        isLoading: false,
        error: null
      });
    },

    setSubscriptionData(subscriptionData: any) {
      update(state => ({
        ...state,
        subscription: subscriptionData,
        isLoading: false,
        error: null
      }));
    }
  };
}

export const subscription = createSubscriptionStore();

// Derived stores that use JWT subscription features instead of API calls
export const currentPlan = derived(
  auth,
  $auth => $auth.subscriptionFeatures || null
);

export const isSubscribed = derived(
  auth,
  $auth => {
    // User has subscription if JWT contains subscription features
    return $auth.subscriptionFeatures !== null;
  }
);

export const planLimits = derived(
  auth,
  $auth => $auth.subscriptionFeatures || null
);


// Helper to check if a feature flag is enabled
export const hasFeature = derived(
  auth,
  $auth => (feature: string): boolean => {
    if (!$auth.token) return false;
    return hasFeatureFromToken($auth.token, feature);
  }
);

// Helper to get a limit value
export const getLimit = derived(
  auth,
  $auth => (limit: string): number => {
    if (!$auth.token) return 0;
    return getLimitFromToken($auth.token, limit);
  }
);

// Helper to check if limit is unlimited (-1)
export const isUnlimited = derived(
  auth,
  $auth => (limit: string): boolean => {
    if (!$auth.token) return false;
    return getLimitFromToken($auth.token, limit) === -1;
  }
);

// Common feature flag constants
export const FEATURES = {
  BASIC_ANALYTICS: 'basic_analytics',
  ADVANCED_ANALYTICS: 'advanced_analytics',
  FEEDBACK_EXPLORER: 'feedback_explorer',
  CUSTOM_BRANDING: 'custom_branding',
  PRIORITY_SUPPORT: 'priority_support'
} as const;

// Common limit constants
export const LIMITS = {
  RESTAURANTS: 'max_organizations',
  QR_CODES: 'max_qr_codes',
  FEEDBACKS_PER_MONTH: 'max_feedbacks_per_month',
  TEAM_MEMBERS: 'max_team_members'
} as const;
