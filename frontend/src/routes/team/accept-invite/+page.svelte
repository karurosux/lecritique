<script lang="ts">
	import { onMount } from 'svelte';
	import { page } from '$app/stores';
	import { goto } from '$app/navigation';
	import { Api } from '$lib/api/api';
	import { Button } from '$lib/components/ui';
	import { Loader2, UserPlus, CheckCircle, XCircle } from 'lucide-svelte';

	let loading = $state(true);
	let error = $state('');
	let success = $state(false);

	onMount(async () => {
		const token = $page.url.searchParams.get('token');
		
		if (!token) {
			error = 'Invalid invitation link';
			loading = false;
			return;
		}

		try {
			const api = new Api({
				baseURL: 'http://localhost:8080'
			});
			const response = await api.api.v1TeamAcceptInviteCreate({ token });
			
			if (response.data.success) {
				success = true;
				setTimeout(() => {
					goto('/login');
				}, 3000);
			} else {
				throw new Error('Failed to accept invitation');
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
					<h2 class="text-xl font-semibold text-gray-900 mb-2">Invitation Accepted!</h2>
					<p class="text-gray-600 mb-6">You've successfully joined the team. Redirecting to login...</p>
					<Button
						variant="gradient"
						class="w-full"
						onclick={() => goto('/login')}
					>
						Go to Login
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