<script lang="ts">
	import { goto } from '$app/navigation';
	import { page } from '$app/stores';
	import { Api } from '$lib/api/api';
	import { Button, Card, Input, Logo } from '$lib/components/ui';
	import { Check, Lock } from 'lucide-svelte';

	const api = new Api({
		baseURL: import.meta.env.VITE_API_URL || 'http://localhost:8080'
	});

	let newPassword = $state('');
	let confirmPassword = $state('');
	let loading = $state(false);
	let error = $state('');
	let success = $state(false);

	const token = $derived($page.url.searchParams.get('token') || '');

	$effect(() => {
		if (!token) {
			goto('/forgot-password');
		}
	});

	async function handleSubmit(e: Event) {
		e.preventDefault();
		
		if (newPassword !== confirmPassword) {
			error = 'Passwords do not match';
			return;
		}

		if (newPassword.length < 8) {
			error = 'Password must be at least 8 characters long';
			return;
		}

		loading = true;
		error = '';

		try {
			await api.api.v1AuthResetPasswordCreate({
				token,
				newPassword
			});
			success = true;
			setTimeout(() => {
				goto('/login');
			}, 3000);
		} catch (err: any) {
			error = err.response?.data?.message || 'Failed to reset password. The link may have expired.';
		} finally {
			loading = false;
		}
	}
</script>

<svelte:head>
	<title>Reset Password - LeCritique</title>
	<meta name="description" content="Create a new password for your LeCritique account" />
</svelte:head>

<div class="min-h-screen flex flex-col justify-center py-12 sm:px-6 lg:px-8">
	<div class="relative z-10 sm:mx-auto sm:w-full sm:max-w-md">
		<div class="flex justify-center mb-8">
			<Logo size="xl" />
		</div>
		
		<div class="text-center space-y-3">
			<h2 class="text-4xl font-bold bg-gradient-to-r from-gray-900 via-blue-800 to-purple-800 bg-clip-text text-transparent">
				Set New Password
			</h2>
			<p class="text-gray-600 text-lg">
				Choose a strong password for your account
			</p>
		</div>
	</div>

	<div class="relative z-10 mt-10 sm:mx-auto sm:w-full sm:max-w-md">

		{#if success}
			<Card>
				<div class="text-center space-y-6">
					<div class="mx-auto flex items-center justify-center h-16 w-16 rounded-full bg-gradient-to-r from-green-400 to-green-600 shadow-lg">
						<Check class="h-8 w-8 text-white" />
					</div>
					<div class="space-y-2">
						<h3 class="text-2xl font-bold text-gray-900">Password Reset!</h3>
						<p class="text-gray-600 max-w-sm mx-auto">
							Your password has been successfully updated. Redirecting to login...
						</p>
					</div>
					<Button 
						href="/login" 
						variant="gradient" 
						size="lg"
						class="w-full shadow-lg hover:shadow-xl transition-all duration-300"
					>
						<svg class="w-5 h-5 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
							<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11 16l-4-4m0 0l4-4m-4 4h14m-5 4v1a3 3 0 01-3 3H6a3 3 0 01-3-3V7a3 3 0 013-3h7a3 3 0 013 3v1"></path>
						</svg>
						Go to Sign In
					</Button>
				</div>
			</Card>
		{:else}
			<Card>
				<form on:submit|preventDefault={handleSubmit} class="space-y-6">
					<Input
						type="password"
						label="New password"
						id="new-password"
						bind:value={newPassword}
						required
						placeholder="Enter new password"
						autocomplete="new-password"
						disabled={loading}
						minlength={8}
					/>

					<Input
						type="password"
						label="Confirm new password"
						id="confirm-password"
						bind:value={confirmPassword}
						required
						placeholder="Confirm new password"
						autocomplete="new-password"
						disabled={loading}
						minlength={8}
					/>

					<p class="text-xs text-gray-500">
						Password must be at least 8 characters long
					</p>

					{#if error}
						<div class="bg-red-50 border border-red-200 rounded-md p-4">
							<div class="flex">
								<div class="flex-shrink-0">
									<svg class="h-5 w-5 text-red-400" fill="currentColor" viewBox="0 0 20 20">
										<path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zM8.707 7.293a1 1 0 00-1.414 1.414L8.586 10l-1.293 1.293a1 1 0 101.414 1.414L10 11.414l1.293 1.293a1 1 0 001.414-1.414L11.414 10l1.293-1.293a1 1 0 00-1.414-1.414L10 8.586 8.707 7.293z" clip-rule="evenodd" />
									</svg>
								</div>
								<div class="ml-3">
									<p class="text-sm text-red-800">{error}</p>
								</div>
							</div>
						</div>
					{/if}

					<Button
						type="submit"
						variant="gradient"
						size="lg"
						disabled={loading || !newPassword || !confirmPassword}
						class="w-full shadow-lg hover:shadow-xl transition-all duration-300"
					>
						{#if loading}
							<svg class="animate-spin -ml-1 mr-3 h-5 w-5 text-white" fill="none" viewBox="0 0 24 24">
								<circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
								<path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
							</svg>
							Resetting password...
						{:else}
							<svg class="w-5 h-5 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
								<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 7a2 2 0 012 2m4 0a6 6 0 01-7.743 5.743L11 17H9v2H7v2H4a1 1 0 01-1-1v-2.586a1 1 0 01.293-.707l5.964-5.964A6 6 0 1121 9z"></path>
							</svg>
							Reset Password
						{/if}
					</Button>
				</form>
			</Card>
		{/if}
	</div>
</div>
</script>