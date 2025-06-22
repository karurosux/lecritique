<script lang="ts">
  import { goto } from '$app/navigation';
  import { auth } from '$lib/stores/auth';
  import { Button, Input, Card } from '$lib/components/ui';
  import { onMount } from 'svelte';

  let email = '';
  let password = '';
  let confirmPassword = '';
  let companyName = '';
  let isSubmitting = false;
  let formErrors: Record<string, string> = {};

  $: authState = $auth;

  // Redirect if already authenticated
  onMount(() => {
    if (authState.isAuthenticated) {
      goto('/dashboard');
    }
  });

  function validateForm() {
    formErrors = {};
    
    if (!email) {
      formErrors.email = 'Email is required';
    } else if (!/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(email)) {
      formErrors.email = 'Please enter a valid email address';
    }
    
    if (!password) {
      formErrors.password = 'Password is required';
    } else if (password.length < 8) {
      formErrors.password = 'Password must be at least 8 characters long';
    }
    
    if (!confirmPassword) {
      formErrors.confirmPassword = 'Please confirm your password';
    } else if (password !== confirmPassword) {
      formErrors.confirmPassword = 'Passwords do not match';
    }
    
    if (!companyName) {
      formErrors.companyName = 'Company name is required';
    }
    
    return Object.keys(formErrors).length === 0;
  }

  async function handleSubmit() {
    if (isSubmitting) return;
    
    if (!validateForm()) return;
    
    isSubmitting = true;
    auth.clearError();
    
    const result = await auth.register({
      email,
      password,
      company_name: companyName
    });
    
    if (result.success) {
      // Registration successful, redirect to login with success message
      goto('/login?message=Registration successful. Please check your email for verification.');
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
  <title>Register - LeCritique</title>
  <meta name="description" content="Create your LeCritique restaurant management account" />
</svelte:head>

<div class="min-h-screen bg-gray-50 flex flex-col justify-center py-12 sm:px-6 lg:px-8">
  <div class="sm:mx-auto sm:w-full sm:max-w-md">
    <div class="flex justify-center">
      <h1 class="text-3xl font-bold text-gray-900">LeCritique</h1>
    </div>
    <h2 class="mt-6 text-center text-3xl font-semibold text-gray-900">
      Create your account
    </h2>
    <p class="mt-2 text-center text-sm text-gray-600">
      Or
      <a href="/login" class="font-medium text-blue-600 hover:text-blue-500">
        sign in to your existing account
      </a>
    </p>
  </div>

  <div class="mt-8 sm:mx-auto sm:w-full sm:max-w-md">
    <Card>
      <form on:submit|preventDefault={handleSubmit} class="space-y-6">
        <Input
          id="companyName"
          type="text"
          label="Company Name"
          bind:value={companyName}
          required
          placeholder="Enter your restaurant/company name"
          disabled={isSubmitting}
          error={formErrors.companyName}
          on:keydown={handleKeyDown}
        />

        <Input
          id="email"
          type="email"
          label="Email address"
          bind:value={email}
          required
          placeholder="Enter your email"
          disabled={isSubmitting}
          error={formErrors.email}
          on:keydown={handleKeyDown}
        />

        <Input
          id="password"
          type="password"
          label="Password"
          bind:value={password}
          required
          placeholder="Create a password (min. 8 characters)"
          disabled={isSubmitting}
          error={formErrors.password}
          on:keydown={handleKeyDown}
        />

        <Input
          id="confirmPassword"
          type="password"
          label="Confirm Password"
          bind:value={confirmPassword}
          required
          placeholder="Confirm your password"
          disabled={isSubmitting}
          error={formErrors.confirmPassword}
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

        <div class="text-sm text-gray-600">
          <p>
            By creating an account, you agree to our
            <a href="/terms" class="font-medium text-blue-600 hover:text-blue-500">Terms of Service</a>
            and
            <a href="/privacy" class="font-medium text-blue-600 hover:text-blue-500">Privacy Policy</a>.
          </p>
        </div>

        <Button
          type="submit"
          variant="primary"
          size="lg"
          disabled={isSubmitting || !email || !password || !confirmPassword || !companyName}
          class="w-full"
        >
          {#if isSubmitting}
            <svg class="animate-spin -ml-1 mr-3 h-5 w-5 text-white" fill="none" viewBox="0 0 24 24">
              <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
              <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
            </svg>
            Creating account...
          {:else}
            Create Account
          {/if}
        </Button>
      </form>
    </Card>
  </div>
</div>