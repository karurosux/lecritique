<script lang="ts">
  import { goto } from '$app/navigation';
  import { auth } from '$lib/stores/auth';
  import { Button, Input, Card, Logo } from '$lib/components/ui';
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

<div class="min-h-screen bg-gradient-to-br from-blue-50 via-white to-purple-50 flex flex-col justify-center py-12 sm:px-6 lg:px-8">
  <!-- Background Pattern -->
  <div class="absolute inset-0 bg-grid-gray-100/50 bg-[size:20px_20px] opacity-30"></div>
  <div class="absolute top-0 left-0 w-full h-full bg-gradient-to-br from-blue-500/5 via-transparent to-purple-500/5"></div>
  
  <div class="relative z-10 sm:mx-auto sm:w-full sm:max-w-md">
    <div class="flex justify-center mb-8">
      <Logo size="xl" />
    </div>
    
    <div class="text-center space-y-3">
      <h2 class="text-4xl font-bold bg-gradient-to-r from-gray-900 via-blue-800 to-purple-800 bg-clip-text text-transparent">
        Welcome Back
      </h2>
      <p class="text-gray-600 text-lg">
        Sign in to your restaurant dashboard
      </p>
      <p class="text-sm text-gray-500">
        Don't have an account?
        <a href="/register" class="font-semibold text-blue-600 hover:text-blue-700 transition-colors duration-200 ml-1">
          Create one here
        </a>
      </p>
    </div>
  </div>

  <div class="relative z-10 mt-10 sm:mx-auto sm:w-full sm:max-w-md">
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
          variant="gradient"
          size="lg"
          disabled={isSubmitting || !email || !password}
          class="w-full shadow-lg hover:shadow-xl transition-all duration-300"
        >
          {#if isSubmitting}
            <svg class="animate-spin -ml-1 mr-3 h-5 w-5 text-white" fill="none" viewBox="0 0 24 24">
              <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
              <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
            </svg>
            Signing in...
          {:else}
            <svg class="w-5 h-5 mr-2" fill="none" stroke="currentColor" viewBox="0 0 24 24">
              <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11 16l-4-4m0 0l4-4m-4 4h14m-5 4v1a3 3 0 01-3 3H6a3 3 0 01-3-3V7a3 3 0 013-3h7a3 3 0 013 3v1"></path>
            </svg>
            Sign in to Account
          {/if}
        </Button>
      </form>
    </Card>
  </div>
</div>