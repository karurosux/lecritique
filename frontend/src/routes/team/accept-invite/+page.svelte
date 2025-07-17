<script lang="ts">
	import { onMount } from 'svelte';
	import { page } from '$app/stores';
	import { goto } from '$app/navigation';
	import { auth } from '$lib/stores/auth';
	import { Api } from '$lib/api/api';
	import { Button } from '$lib/components/ui';
	import { Loader2, UserPlus, CheckCircle, XCircle } from 'lucide-svelte';

	let loading = $state(true);
	let error = $state('');
	let success = $state(false);
	let invitationStatus = $state<'pending' | 'accepted' | null>(null);

	onMount(async () => {
		const token = $page.url.searchParams.get('token');
		
		if (!token) {
			error = 'Invalid invitation link';
			loading = false;
			return;
		}

		try {
			console.log('Accepting invitation with token:', token);
			console.log('Is authenticated:', $auth.isAuthenticated);
			
			// Get API instance - will include auth token if user is logged in
			const api = $auth.isAuthenticated ? auth.getApi() : new Api({
				baseURL: 'http://localhost:8080'
			});
			
			console.log('Calling API endpoint...');
			const response = await api.api.v1TeamAcceptInviteCreate({ token });
			console.log('API response:', response.data);
			
			if (response.data.success && response.data.data) {
				const data = response.data.data as any;
				invitationStatus = data.status;
				
				if (data.status === 'accepted') {
					// Invitation was accepted (user was authenticated)
					success = true;
					setTimeout(() => {
						goto('/settings/team');
					}, 2000);
				} else if (data.status === 'needs_registration' || data.status === 'pending') {
					// User needs to login or register
					invitationStatus = 'pending'; // Normalize status for UI
					success = true;
					setTimeout(() => {
						goto('/login');
					}, 3000);
				} else {
					throw new Error('Unknown invitation status: ' + data.status);
				}
			} else {
				throw new Error('Failed to process invitation');
			}
		} catch (err: any) {
			error = err.response?.data?.error?.message || 'Failed to accept invitation. The link may be expired or invalid.';
		} finally {
			loading = false;
		}
	});
</script>

<div class="min-h-screen bg-gradient-to-br from-slate-50 to-slate-100 flex items-center justify-center px-4">
	<div class="max-w-md w-full">
		<div class="bg-white rounded-2xl shadow-xl p-8">
			{#if loading}
				<div class="text-center">
					<Loader2 class="h-12 w-12 animate-spin text-blue-600 mx-auto mb-4" />
					<h2 class="text-xl font-semibold text-gray-900 mb-2">Processing Invitation</h2>
					<p class="text-gray-600">Please wait while we accept your team invitation...</p>
				</div>
			{:else if success}
				<div class="text-center">
					<CheckCircle class="h-12 w-12 text-green-600 mx-auto mb-4" />
					<h2 class="text-xl font-semibold text-gray-900 mb-2">
						{invitationStatus === 'accepted' ? 'Invitation Accepted!' : 'Valid Invitation!'}
					</h2>
					<p class="text-gray-600 mb-6">
						{#if invitationStatus === 'accepted'}
							You've successfully joined the team. Redirecting to team settings...
						{:else}
							Please login or register to join the team. Redirecting...
						{/if}
					</p>
					<Button
						variant="gradient"
						class="w-full"
						onclick={() => goto(invitationStatus === 'accepted' ? '/settings/team' : '/login')}
					>
						{invitationStatus === 'accepted' ? 'Go to Team Settings' : 'Go to Login'}
					</Button>
				</div>
			{:else if error}
				<div class="text-center">
					<XCircle class="h-12 w-12 text-red-600 mx-auto mb-4" />
					<h2 class="text-xl font-semibold text-gray-900 mb-2">Invitation Error</h2>
					<p class="text-gray-600 mb-6">{error}</p>
					<Button
						variant="outline"
						class="w-full"
						onclick={() => goto('/login')}
					>
						Go to Login
					</Button>
				</div>
			{/if}
		</div>
		
		<div class="mt-6 text-center">
			<p class="text-sm text-gray-600">
				Need help? <a href="/contact" class="text-blue-600 hover:text-blue-700 font-medium">Contact support</a>
			</p>
		</div>
	</div>
</div>