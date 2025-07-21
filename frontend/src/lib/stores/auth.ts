import { browser } from '$app/environment';
import { Api, type HandlersLoginRequest, type HandlersRegisterRequest } from '$lib/api/api';
import { APP_CONFIG } from '$lib/constants/config';
import { decodeJwt, getSubscriptionFeaturesFromToken, type JwtPayload } from '$lib/utils/jwt';
import { writable } from 'svelte/store';

type UserRole = 'OWNER' | 'ADMIN' | 'MANAGER' | 'VIEWER';

export interface User {
  id: string;
  email: string;
  name: string;
  phone?: string;
  email_verified: boolean;
  deactivation_requested_at?: string | null;
  account_id?: string; // The account they're accessing (important for team members)
  role?: UserRole; // User's role in the current account
}

export interface AuthState {
  user: User | null;
  token: string | null;
  subscriptionFeatures: JwtPayload['subscription_features'] | null;
  isAuthenticated: boolean;
  isLoading: boolean;
  error: string | null;
}

const getInitialState = (): AuthState => {
  const baseState = {
    user: null,
    token: null,
    subscriptionFeatures: null,
    isAuthenticated: false,
    isLoading: false,
    error: null
  };

  if (browser) {
    const storedToken = localStorage.getItem(APP_CONFIG.localStorageKeys.authToken);
    const storedUser = localStorage.getItem(APP_CONFIG.localStorageKeys.authUser);

    if (storedToken && storedUser) {
      try {
        const user = JSON.parse(storedUser);
        const subscriptionFeatures = getSubscriptionFeaturesFromToken(storedToken);
        
        return {
          ...baseState,
          user,
          token: storedToken,
          subscriptionFeatures,
          isAuthenticated: true
        };
      } catch (error) {
        localStorage.removeItem(APP_CONFIG.localStorageKeys.authToken);
        localStorage.removeItem(APP_CONFIG.localStorageKeys.authUser);
      }
    }
  }

  return baseState;
};

const initialState = getInitialState();

function createAuthStore() {
  const { subscribe, set, update } = writable<AuthState>(initialState);

  const api = new Api({
    baseURL: 'http://localhost:8080',
    securityWorker: (securityData) => {
      if (securityData) {
        return {
          headers: {
            Authorization: `Bearer ${securityData}`
          }
        };
      }
    }
  });

  if (initialState.token) {
    api.setSecurityData(initialState.token);
  }

  const authStore = {
    subscribe,

    async login(credentials: HandlersLoginRequest) {
      update(state => ({ ...state, isLoading: true, error: null }));

      try {
        const response = await api.api.v1AuthLoginCreate(credentials);

        if (response.data.success && response.data.data) {
          const { token } = response.data.data;

          if (token) {
            const payload = decodeJwt(token);
            if (!payload) {
              throw new Error('Invalid token received');
            }

            const subscriptionFeatures = getSubscriptionFeaturesFromToken(token);

            const user: User = {
              id: payload.member_id || payload.account_id,
              email: payload.email,
              name: payload.name || '', // Name comes from JWT token
              phone: '', // Phone is not in JWT, would need separate API call if needed
              email_verified: true, // If they can login, email is verified
              deactivation_requested_at: null,
              account_id: payload.account_id,
              role: payload.role as UserRole
            };

            if (browser) {
              localStorage.setItem(APP_CONFIG.localStorageKeys.authToken, token);
              localStorage.setItem(APP_CONFIG.localStorageKeys.authUser, JSON.stringify(user));
            }

            api.setSecurityData(token);

            update(state => ({
              ...state,
              user,
              token,
              subscriptionFeatures,
              isAuthenticated: true,
              isLoading: false,
              error: null
            }));

            // Note: We no longer need to update subscription store separately
            // as all subscription data is now in the JWT token

            return { success: true };
          }
        }

        throw new Error('Invalid response from server');
      } catch (error: any) {
        const errorCode = error.response?.data?.error?.code;
        const errorMessage = error.response?.data?.error?.message || error.message || 'Login failed';

        if (errorCode === 'EMAIL_NOT_VERIFIED') {
          update(state => ({
            ...state,
            isLoading: false,
            error: null // Don't show error since we're redirecting
          }));

          return { success: false, unverified: true, email: credentials.email };
        }

        update(state => ({
          ...state,
          isLoading: false,
          error: errorMessage
        }));

        return { success: false, error: errorMessage };
      }
    },

    async register(userData: HandlersRegisterRequest) {
      update(state => ({ ...state, isLoading: true, error: null }));

      try {
        const response = await api.api.v1AuthRegisterCreate(userData);

        if (response.data.success) {
          update(state => ({
            ...state,
            isLoading: false,
            error: null
          }));

          return { success: true };
        }

        throw new Error('Registration failed');
      } catch (error: any) {
        const errorMessage = error.response?.data?.error?.message || error.message || 'Registration failed';

        update(state => ({
          ...state,
          isLoading: false,
          error: errorMessage
        }));

        return { success: false, error: errorMessage };
      }
    },

    async logout() {
      if (browser) {
        localStorage.removeItem(APP_CONFIG.localStorageKeys.authToken);
        localStorage.removeItem(APP_CONFIG.localStorageKeys.authUser);
      }

      api.setSecurityData(null);

      const { subscription } = await import('./subscription');
      subscription.reset();

      set({
        user: null,
        token: null,
        subscriptionFeatures: null,
        isAuthenticated: false,
        isLoading: false,
        error: null
      });

      return Promise.resolve();
    },

    async refreshToken() {
      try {
        const response = await api.api.v1AuthRefreshCreate();

        if (response.data.success && response.data.data?.token) {
          const newToken = response.data.data.token;

          const subscriptionFeatures = getSubscriptionFeaturesFromToken(newToken);

          if (browser) {
            localStorage.setItem('auth_token', newToken);
          }

          api.setSecurityData(newToken);

          update(state => ({
            ...state,
            token: newToken,
            subscriptionFeatures
          }));

          return { success: true };
        }

        throw new Error('Token refresh failed');
      } catch (error) {
        // If refresh fails, logout user
        this.logout();
        return { success: false };
      }
    },

    clearError() {
      update(state => ({ ...state, error: null }));
    },

    updateToken(newToken: string) {
      try {
        const payload = decodeJwt(newToken);
        if (!payload) {
          throw new Error('Invalid token payload');
        }

        const subscriptionFeatures = getSubscriptionFeaturesFromToken(newToken);

        if (browser) {
          localStorage.setItem(APP_CONFIG.localStorageKeys.authToken, newToken);

          const storedUser = localStorage.getItem(APP_CONFIG.localStorageKeys.authUser);
          if (storedUser && payload.email) {
            const user = JSON.parse(storedUser);
            user.email = payload.email;
            localStorage.setItem(APP_CONFIG.localStorageKeys.authUser, JSON.stringify(user));
          }
        }

        api.setSecurityData(newToken);

        update(state => ({
          ...state,
          token: newToken,
          subscriptionFeatures,
          user: state.user ? {
            ...state.user,
            email: payload.email || state.user.email
          } : null
        }));
      } catch (error) {
        console.error('Failed to update token:', error);
      }
    },

    updateUser(updatedUser: User) {
      update(state => ({
        ...state,
        user: updatedUser
      }));

      if (browser) {
        localStorage.setItem(APP_CONFIG.localStorageKeys.authUser, JSON.stringify(updatedUser));
      }
    },

    getApi() {
      return api;
    }
  };

  const originalRequest = api.request.bind(api);
  api.request = async (params: any) => {
    try {
      return await originalRequest(params);
    } catch (error: any) {
      if (error.response?.status === 401 || error.response?.status === 403) {
        const errorCode = error.response?.data?.error?.code;
        if (errorCode === 'EMAIL_NOT_VERIFIED') {
          throw error;
        }

        await authStore.logout();

        if (browser && typeof window !== 'undefined') {
          window.location.href = '/login';
        }
      }
      throw error;
    }
  };

  return authStore;
}

export const auth = createAuthStore();
