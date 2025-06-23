<script lang="ts">
  import { goto } from '$app/navigation';
  import { auth } from '$lib/stores/auth';
  import { Button, Input, Card } from '$lib/components/ui';
  import { onMount } from 'svelte';

  let email = $state('');
  let password = $state('');
  let isSubmitting = $state(false);

  let authState = $derived($auth);

  // Redirect if already authenticated
  $effect(() => {
    if (authState.isAuthenticated) {
      goto('/dashboard');
    }
  });

  async function handleSubmit() {
    if (isSubmitting) return;
    
    isSubmitting = true;
    auth.clearError();
    
    const result = await auth.login({ email, password });
    
    if (result.success) {
      goto('/dashboard');
    }
    
    isSubmitting = false;
  }

  function handleKeyDown(event: KeyboardEvent) {
    if (event.key === 'Enter') {
      handleSubmit();
    }
  }
</script>

<svelte:head>
  <title>Login - LeCritique</title>
  <meta name="description" content="Login to your LeCritique restaurant management account" />
</svelte:head>

<div class="min-h-screen bg-gray-50 flex flex-col justify-center py-12 sm:px-6 lg:px-8">
  <div class="sm:mx-auto sm:w-full sm:max-w-md">
    <div class="flex justify-center">
      <h1 class="text-3xl font-bold text-gray-900">LeCritique</h1>
    </div>
    <h2 class="mt-6 text-center text-3xl font-semibold text-gray-900">
      Sign in to your account
    </h2>
    <p class="mt-2 text-center text-sm text-gray-600">
      Or
      <a href="/register" class="font-medium text-blue-600 hover:text-blue-500">
        create a new account
      </a>
    </p>
  </div>

  <div class="mt-8 sm:mx-auto sm:w-full sm:max-w-md">
    <Card>
      <form on:submit|preventDefault={handleSubmit} class="space-y-6">
        <Input
          id="email"
          type="email"
          label="Email address"
          bind:value={email}
          required
          placeholder="Enter your email"
          disabled={isSubmitting}
          on:keydown={handleKeyDown}
        />

        <Input
          id="password"
          type="password"
          label="Password"
          bind:value={password}
          required
          placeholder="Enter your password"
          disabled={isSubmitting}
          on:keydown={handleKeyDown}
        />

        {#if authState.error}
          <div class="bg-red-50 border border-red-200 rounded-md p-4">
            <div class="flex">
              <div class="flex-shrink-0">
                <svg class="h-5 w-5 text-red-400" fill="currentColor" viewBox="0 0 20 20">
                  <path fill-rule="evenodd" d="M10 18a8 8 0 100-16 8 8 0 000 16zM8.707 7.293a1 1 0 00-1.414 1.414L8.586 10l-1.293 1.293a1 1 0 101.414 1.414L10 11.414l1.293 1.293a1 1 0 001.414-1.414L11.414 10l1.293-1.293a1 1 0 00-1.414-1.414L10 8.586 8.707 7.293z" clip-rule="evenodd" />
                </svg>
              </div>
              <div class="ml-3">
                <p class="text-sm text-red-800">{authState.error}</p>
              </div>
            </div>
          </div>
        {/if}

        <div class="flex items-center justify-between">
          <div class="text-sm">
            <a href="/forgot-password" class="font-medium text-blue-600 hover:text-blue-500">
              Forgot your password?
            </a>
          </div>
        </div>

        <Button
          type="submit"
          variant="primary"
          size="lg"
          disabled={isSubmitting || !email || !password}
          class="w-full"
        >
          {#if isSubmitting}
            <svg class="animate-spin -ml-1 mr-3 h-5 w-5 text-white" fill="none" viewBox="0 0 24 24">
              <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
              <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
            </svg>
            Signing in...
          {:else}
            Sign in
          {/if}
        </Button>
      </form>
    </Card>
  </div>
</div>