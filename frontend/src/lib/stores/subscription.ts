import { writable, derived } from 'svelte/store';
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
  $subscription => $subscription.subscription?.status === 'active'
);

export const planLimits = derived(
  subscription,
  $subscription => $subscription.subscription?.plan?.features || null
);