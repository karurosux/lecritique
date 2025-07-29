import { redirect } from '@sveltejs/kit';
import { auth } from '$lib/stores/auth';
import { get } from 'svelte/store';

export function requireAuth() {
  const authState = get(auth);

  if (!authState.isAuthenticated) {
    throw redirect(302, '/login');
  }

  return authState;
}

export function requireGuest() {
  const authState = get(auth);

  if (authState.isAuthenticated) {
    throw redirect(302, '/dashboard');
  }

  return authState;
}
