import { auth } from '$lib/stores/auth';
import { get } from 'svelte/store';
import { Api } from '$lib/api/api';
import { APP_CONFIG } from '$lib/constants/config';

// Export the API client from auth store for easy access
export function getApiClient() {
  return auth.getApi();
}

// Export a public API client for unauthenticated requests
export function getPublicApiClient() {
  return new Api({
    baseURL: APP_CONFIG.API_URL,
  });
}

// Export a public API client for server-side load functions
export function getServerSideApiClient(fetch?: typeof globalThis.fetch) {
  return new Api({
    baseURL: APP_CONFIG.API_URL,
    customFetch: fetch,
  });
}

// Helper function to handle API errors consistently
export function handleApiError(error: any): string {
  if (error.response?.data?.error?.message) {
    return error.response.data.error.message;
  }

  if (error.response?.data?.message) {
    return error.response.data.message;
  }

  if (error.message) {
    return error.message;
  }

  return 'An unexpected error occurred';
}

// Helper function to check if user is authenticated
export function isAuthenticated(): boolean {
  return get(auth).isAuthenticated;
}

// Helper function to get current user
export function getCurrentUser() {
  return get(auth).user;
}

// Helper function to get auth token
export function getAuthToken(): string | null {
  return get(auth).token;
}
