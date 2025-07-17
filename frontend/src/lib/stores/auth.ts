import { writable } from 'svelte/store';
import { browser } from '$app/environment';
import { Api, HttpClient, type HandlersAuthResponse, type HandlersLoginRequest, type HandlersRegisterRequest } from '$lib/api/api';

export interface User {
  id: string;
  email: string;
  company_name: string;
  phone?: string;
  email_verified: boolean;
  deactivation_requested_at?: string | null;
  account_id?: string; // The account they're accessing (important for team members)
}

export interface AuthState {
  user: User | null;
  token: string | null;
  isAuthenticated: boolean;
  isLoading: boolean;
  error: string | null;
}

// Check for stored auth data to set initial state properly
const getInitialState = (): AuthState => {
  const baseState = {
    user: null,
    token: null,
    isAuthenticated: false,
    isLoading: false,
    error: null
  };

  if (browser) {
    const storedToken = localStorage.getItem('auth_token');
    const storedUser = localStorage.getItem('auth_user');
    
    if (storedToken && storedUser) {
      try {
        const user = JSON.parse(storedUser);
        return {
          ...baseState,
          user,
          token: storedToken,
          isAuthenticated: true
        };
      } catch (error) {
        // Clear invalid stored data
        localStorage.removeItem('auth_token');
        localStorage.removeItem('auth_user');
      }
    }
  }

  return baseState;
};

const initialState = getInitialState();

function createAuthStore() {
  const { subscribe, set, update } = writable<AuthState>(initialState);

  // Initialize API client
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

  // Set initial API security data if we have a token
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
          const { token, account, subscription: subscriptionData } = response.data.data;
          
          if (token && account) {
            const user: User = {
              id: account.id,
              email: account.email,
              company_name: account.company_name,
              phone: account.phone,
              email_verified: account.email_verified,
              deactivation_requested_at: account.deactivation_requested_at,
              account_id: account.id // Store which account they're accessing
            };

            // Store in localStorage
            if (browser) {
              localStorage.setItem('auth_token', token);
              localStorage.setItem('auth_user', JSON.stringify(user));
            }

            // Set security data for API client
            api.setSecurityData(token);

            update(state => ({
              ...state,
              user,
              token,
              isAuthenticated: true,
              isLoading: false,
              error: null
            }));

            // Update subscription store if subscription data is provided
            if (subscriptionData) {
              const { subscription } = await import('./subscription');
              subscription.setSubscriptionData(subscriptionData);
              console.log('[Auth] Setting subscription data:', subscriptionData);
            } else {
              console.log('[Auth] No subscription data in login response');
            }

            return { success: true };
          }
        }
        
        throw new Error('Invalid response from server');
      } catch (error: any) {
        const errorCode = error.response?.data?.error?.code;
        const errorMessage = error.response?.data?.error?.message || error.message || 'Login failed';
        
        // Check if it's an email verification error
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
      // Clear localStorage
      if (browser) {
        localStorage.removeItem('auth_token');
        localStorage.removeItem('auth_user');
      }

      // Clear security data
      api.setSecurityData(null);

      // Reset subscription data
      const { subscription } = await import('./subscription');
      subscription.reset();

      // Reset store to clean initial state
      set({
        user: null,
        token: null,
        isAuthenticated: false,
        isLoading: false,
        error: null
      });
      
      // Return a promise to ensure async completion
      return Promise.resolve();
    },

    async refreshToken() {
      try {
        const response = await api.api.v1AuthRefreshCreate();
        
        if (response.data.success && response.data.data?.token) {
          const newToken = response.data.data.token;
          
          // Update stored token
          if (browser) {
            localStorage.setItem('auth_token', newToken);
          }

          // Update security data
          api.setSecurityData(newToken);

          update(state => ({
            ...state,
            token: newToken
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
      // Validate the token first
      try {
        const payload = JSON.parse(atob(newToken.split('.')[1]));
        
        // Update stored token
        if (browser) {
          localStorage.setItem('auth_token', newToken);
          
          // Update user info from token if email changed
          const storedUser = localStorage.getItem('auth_user');
          if (storedUser && payload.email) {
            const user = JSON.parse(storedUser);
            user.email = payload.email;
            localStorage.setItem('auth_user', JSON.stringify(user));
          }
        }

        // Update security data
        api.setSecurityData(newToken);

        update(state => ({
          ...state,
          token: newToken,
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
      
      // Update stored user
      if (browser) {
        localStorage.setItem('auth_user', JSON.stringify(updatedUser));
      }
    },

    // Expose API client for authenticated requests
    getApi() {
      return api;
    }
  };

  // Add response interceptor to handle authentication errors
  const originalRequest = api.request.bind(api);
  api.request = async (params: any) => {
    try {
      return await originalRequest(params);
    } catch (error: any) {
      // Check for authentication errors
      if (error.response?.status === 401 || error.response?.status === 403) {
        // Check if it's an email verification error - don't logout for this
        const errorCode = error.response?.data?.error?.code;
        if (errorCode === 'EMAIL_NOT_VERIFIED') {
          // Don't logout or redirect, let the login handler deal with it
          throw error;
        }
        
        // Token is invalid, logout user
        await authStore.logout();
        
        // Redirect to login if we're in the browser
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