import type { LayoutLoad } from './$types';
import { auth } from '$lib/stores/auth';
import { get } from 'svelte/store';

export const load: LayoutLoad = async () => {
	// Get the current auth state
	const authState = get(auth);
	
	return {
		user: authState.user
	};
};