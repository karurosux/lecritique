import { writable, derived } from 'svelte/store';
import { browser } from '$app/environment';
import { getApiClient } from '$lib/api/client';
import type { ModelsSubscription, ModelsSubscriptionPlan } from '$lib/api/api';


interface SubscriptionState {
  subscription: ModelsSubscription | null;
  plans: ModelsSubscriptionPlan[];
  usage: SubscriptionUsage | null;
  isLoading: boolean;
  error: string | null;
}

interface SubscriptionUsage {
  feedbacks_count: number;
  restaurants_count: number;
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
          case 'restaurant':
            const response = await api.api.v1UserCanCreateRestaurantList();
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

// Derived stores for easy access
export const currentPlan = derived(
  subscription,
  $subscription => $subscription.subscription?.plan || null
);

export const isSubscribed = derived(
  subscription,
  $subscription => {
    // Check if status matches the actual string value from backend
    // The backend sends "active" as a string, but TypeScript expects the enum
    const status = $subscription.subscription?.status;
    
    // Handle both enum constant and raw string value
    return status === 'active' || status === 'SubscriptionActive';
  }
);

export const planLimits = derived(
  subscription,
  $subscription => $subscription.subscription?.plan?.features || null
);


// Helper to check if a feature flag is enabled
export const hasFeature = derived(
  subscription,
  $subscription => (feature: string): boolean => {
    const plan = $subscription.subscription?.plan;
    if (!plan) return false;
    
    switch(feature) {
      case FEATURES.BASIC_ANALYTICS:
        return plan.has_basic_analytics || false;
      case FEATURES.ADVANCED_ANALYTICS:
        return plan.has_advanced_analytics || false;
      case FEATURES.FEEDBACK_EXPLORER:
        return plan.has_feedback_explorer || false;
      case FEATURES.CUSTOM_BRANDING:
        return plan.has_custom_branding || false;
      case FEATURES.PRIORITY_SUPPORT:
        return plan.has_priority_support || false;
      default:
        return false;
    }
  }
);

// Helper to get a limit value
export const getLimit = derived(
  subscription,
  $subscription => (limit: string): number => {
    const plan = $subscription.subscription?.plan;
    if (!plan) return 0;
    
    switch(limit) {
      case LIMITS.RESTAURANTS:
        return plan.max_restaurants || 0;
      case LIMITS.QR_CODES:
        return plan.max_qr_codes || 0;
      case LIMITS.FEEDBACKS_PER_MONTH:
        return plan.max_feedbacks_per_month || 0;
      case LIMITS.TEAM_MEMBERS:
        return plan.max_team_members || 0;
      default:
        return 0;
    }
  }
);

// Helper to check if limit is unlimited (-1)
export const isUnlimited = derived(
  subscription,
  $subscription => (limit: string): boolean => {
    const plan = $subscription.subscription?.plan;
    if (!plan) return false;
    
    switch(limit) {
      case LIMITS.RESTAURANTS:
        return plan.max_restaurants === -1;
      case LIMITS.QR_CODES:
        return plan.max_qr_codes === -1;
      case LIMITS.FEEDBACKS_PER_MONTH:
        return plan.max_feedbacks_per_month === -1;
      case LIMITS.TEAM_MEMBERS:
        return plan.max_team_members === -1;
      default:
        return false;
    }
  }
);

// Common feature flag constants
export const FEATURES = {
  BASIC_ANALYTICS: 'basic_analytics',
  ADVANCED_ANALYTICS: 'advanced_analytics',
  FEEDBACK_EXPLORER: 'feedback_explorer',
  CUSTOM_BRANDING: 'custom_branding',
  PRIORITY_SUPPORT: 'priority_support'
};

// Common limit constants
export const LIMITS = {
  RESTAURANTS: 'max_restaurants',
  QR_CODES: 'max_qr_codes',
  FEEDBACKS_PER_MONTH: 'max_feedbacks_per_month',
  TEAM_MEMBERS: 'max_team_members'
};
