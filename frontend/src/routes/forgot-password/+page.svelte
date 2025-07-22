<script lang="ts">
	import { goto } from '$app/navigation';
	import { Api } from '$lib/api/api';
	import { Button, Card, Input, Logo } from '$lib/components/ui';
	import { Check, Mail } from 'lucide-svelte';

	const api = new Api({
		baseURL: import.meta.env.VITE_API_URL || 'http://localhost:8080'
	});

	let email = $state('');
	let loading = $state(false);
	let error = $state('');
	let success = $state(false);

	async function handleSubmit(e: Event) {
		e.preventDefault();
		loading = true;
		error = '';

		try {
			await api.api.v1AuthForgotPasswordCreate({ email });
			success = true;
		} catch (err: any) {
			error = err.response?.data?.message || 'Failed to send reset email. Please try again.';
		} finally {
			loading = false;
		}
	}
</script>

<svelte:head>
	<title>Forgot Password - Kyooar</title>
	<meta name="description" content="Reset your Kyooar account password" />
</svelte:head>

<div class="min-h-screen flex flex-col justify-center py-12 sm:px-6 lg:px-8">
	<div class="relative z-10 sm:mx-auto sm:w-full sm:max-w-md">
		<div class="flex justify-center mb-8">
			<Logo size="xl" />
		</div>
		
		<div class="text-center space-y-3">
			<h2 class="text-4xl font-bold bg-gradient-to-r from-gray-900 via-blue-800 to-purple-800 bg-clip-text text-transparent">
				Reset Password
			</h2>
			<p class="text-gray-600 text-lg">
				We'll send you a link to reset your password
			</p>
			<p class="text-sm text-gray-500">
				Remember your password?
				<a href="/login" class="font-semibold text-blue-600 hover:text-blue-700 transition-colors duration-200 ml-1">
					Sign in here
				</a>
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
						<h3 class="text-2xl font-bold text-gray-900">Check your email!</h3>
						<p class="text-gray-600 max-w-sm mx-auto">
							We've sent a password reset link to
						</p>
						<p class="font-semibold text-gray-900">{email}</p>
					</div>
					<Button 
						href="/login" 
						variant="gradient" 
						size="lg"
						class="w-full shadow-lg hover:shadow-xl transition-all duration-300"
					>
						<svg class="w-5 h-5 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
							<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M10 19l-7-7m0 0l7-7m-7 7h18"></path>
						</svg>
						Back to Sign In
					</Button>
				</div>
			</Card>
		{:else}
			<Card>
				<form on:submit|preventDefault={handleSubmit} class="space-y-6">
					<Input
						id="email"
						type="email"
						label="Email address"
						bind:value={email}
						required
						placeholder="Enter your email"
						disabled={loading}
						autocomplete="email"
					/>

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
						disabled={loading || !email}
						class="w-full shadow-lg hover:shadow-xl transition-all duration-300"
					>
						{#if loading}
							<svg class="animate-spin -ml-1 mr-3 h-5 w-5 text-white" fill="none" viewBox="0 0 24 24">
								<circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
								<path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
							</svg>
							Sending link...
						{:else}
							<svg class="w-5 h-5 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
								<path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M3 8l7.89 5.26a2 2 0 002.22 0L21 8M5 19h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v10a2 2 0 002 2z"></path>
							</svg>
							Send Reset Link
						{/if}
					</Button>
				</form>
			</Card>
		{/if}
	</div>
</div>