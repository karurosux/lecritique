import { writable } from 'svelte/store';
import { browser } from '$app/environment';
import { Api, HttpClient, type HandlersAuthResponse, type HandlersLoginRequest, type HandlersRegisterRequest } from '$lib/api/api';

export interface User {
  id: string;
  email: string;
  company_name: string;
  email_verified: boolean;
}

export interface AuthState {
  user: User | null;
  token: string | null;
  isAuthenticated: boolean;
  isLoading: boolean;
  error: string | null;
}

const initialState: AuthState = {
  user: null,
  token: null,
  isAuthenticated: false,
  isLoading: false,
  error: null
};

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

  // Load auth state from localStorage on initialization
  if (browser) {
    const storedToken = localStorage.getItem('auth_token');
    const storedUser = localStorage.getItem('auth_user');
    
    if (storedToken && storedUser) {
      try {
        const user = JSON.parse(storedUser);
        api.setSecurityData(storedToken);
        set({
          user,
          token: storedToken,
          isAuthenticated: true,
          isLoading: false,
          error: null
        });
      } catch (error) {
        // Clear invalid stored data
        localStorage.removeItem('auth_token');
        localStorage.removeItem('auth_user');
      }
    }
  }

  return {
    subscribe,
    
    async login(credentials: HandlersLoginRequest) {
      update(state => ({ ...state, isLoading: true, error: null }));
      
      try {
        const response = await api.api.v1AuthLoginCreate(credentials);
        
        if (response.data.success && response.data.data) {
          const { token, account } = response.data.data;
          
          if (token && account) {
            const user: User = {
              id: account.id,
              email: account.email,
              company_name: account.company_name,
              email_verified: account.email_verified
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

            return { success: true };
          }
        }
        
        throw new Error('Invalid response from server');
      } catch (error: any) {
        const errorMessage = error.response?.data?.error?.message || error.message || 'Login failed';
        
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

      // Reset store
      set(initialState);
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

    // Expose API client for authenticated requests
    getApi() {
      return api;
    }
  };
}

export const auth = createAuthStore();