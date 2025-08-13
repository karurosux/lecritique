<script lang="ts">
  import { goto } from '$app/navigation';
  import { auth } from '$lib/stores/auth';
  import { Button, Input, Card, Logo } from '$lib/components/ui';
  import { onMount } from 'svelte';
  import { XCircle, Loader2, UserPlus } from 'lucide-svelte';

  let email = $state('');
  let password = $state('');
  let confirmPassword = $state('');
  let name = $state('');
  let isSubmitting = $state(false);
  let formErrors = $state<Record<string, string>>({});

  let authState = $derived($auth);

  $effect(() => {
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

    if (!name) {
      formErrors.name = 'Name is required';
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
      name: name,
    });

    if (result.success) {
      goto('/registration-success');
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
  <title>Register - Kyooar</title>
  <meta
    name="description"
    content="Create your Kyooar organization management account" />
</svelte:head>

<div class="min-h-full flex flex-col justify-center py-12 sm:px-6 lg:px-8">
  <!-- Background Pattern -->
  <div
    class="register-grid-pattern fixed inset-0 bg-grid-gray-100/50 bg-[size:20px_20px] opacity-30 -z-10">
  </div>

  <div class="relative z-10 sm:mx-auto sm:w-full sm:max-w-md">
    <div class="flex justify-center mb-8">
      <Logo size="xl" />
    </div>

    <div class="text-center space-y-3">
      <h2
        class="text-4xl font-bold bg-gradient-to-r from-gray-900 via-blue-800 to-purple-800 bg-clip-text text-transparent">
        Join Kyooar
      </h2>
      <p class="text-gray-600 text-lg">
        Create your organization management account
      </p>
      <p class="text-sm text-gray-500">
        Already have an account?
        <a
          href="/login"
          class="font-semibold text-blue-600 hover:text-blue-700 transition-colors duration-200 ml-1">
          Sign in here
        </a>
      </p>
    </div>
  </div>

  <div class="relative z-10 mt-10 sm:mx-auto sm:w-full sm:max-w-md">
    <Card>
      <form on:submit|preventDefault={handleSubmit} class="space-y-6">
        <Input
          id="name"
          type="text"
          label="Name"
          bind:value={name}
          required
          placeholder="Enter your name"
          disabled={isSubmitting}
          error={formErrors.name}
          on:keydown={handleKeyDown} />

        <Input
          id="email"
          type="email"
          label="Email address"
          bind:value={email}
          required
          placeholder="Enter your email"
          disabled={isSubmitting}
          error={formErrors.email}
          on:keydown={handleKeyDown} />

        <Input
          id="password"
          type="password"
          label="Password"
          bind:value={password}
          required
          placeholder="Create a password (min. 8 characters)"
          disabled={isSubmitting}
          error={formErrors.password}
          on:keydown={handleKeyDown} />

        <Input
          id="confirmPassword"
          type="password"
          label="Confirm Password"
          bind:value={confirmPassword}
          required
          placeholder="Confirm your password"
          disabled={isSubmitting}
          error={formErrors.confirmPassword}
          on:keydown={handleKeyDown} />

        {#if authState.error}
          <div class="bg-red-50 border border-red-200 rounded-md p-4">
            <div class="flex">
              <div class="flex-shrink-0">
                <XCircle class="h-5 w-5 text-red-400" />
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
            <a
              href="/terms"
              class="font-medium text-blue-600 hover:text-blue-500"
              >Terms of Service</a>
            and
            <a
              href="/privacy"
              class="font-medium text-blue-600 hover:text-blue-500"
              >Privacy Policy</a
            >.
          </p>
        </div>

        <Button
          type="submit"
          variant="gradient"
          size="lg"
          disabled={isSubmitting ||
            !email ||
            !password ||
            !confirmPassword ||
            !name}
          class="w-full shadow-lg hover:shadow-xl transition-all duration-300">
          {#if isSubmitting}
            <Loader2 class="animate-spin -ml-1 mr-3 h-5 w-5 text-white" />
            Creating account...
          {:else}
            <UserPlus class="w-5 h-5 mr-2" />
            Create Your Account
          {/if}
        </Button>
      </form>
    </Card>
  </div>
</div>
