<script lang="ts">
  import { page } from '$app/stores';
  import { goto } from '$app/navigation';
  import { Logo } from '$lib/components/ui';
  import { CheckCircle, XCircle, Loader2 } from 'lucide-svelte';
  import { onMount } from 'svelte';
  import { getApiClient } from '$lib/api/client';

  let token = $derived($page.url.searchParams.get('token') || '');
  let verifying = $state(true);
  let verified = $state(false);
  let error = $state('');

  onMount(async () => {
    if (!token) {
      error = 'No verification token provided';
      verifying = false;
      return;
    }

    try {
      const api = getApiClient();
      const response = await api.api.v1AuthVerifyEmailList({ token });

      if (response.data.success) {
        verified = true;

        setTimeout(() => {
          goto('/login');
        }, 3000);
      }
    } catch (err: any) {
      error = err.response?.data?.error?.message || 'Failed to verify email';
    } finally {
      verifying = false;
    }
  });
</script>

<svelte:head>
  <title>Verify Email - Kyooar</title>
  <meta name="description" content="Verify your email address for Kyooar" />
</svelte:head>

<div
  class="min-h-screen bg-gradient-to-b from-white to-gray-50/50 flex items-center justify-center px-4">
  <div class="verify-email-container max-w-md w-full">
    <div class="text-center">
      <!-- Logo -->
      <div class="flex justify-center mb-8">
        <Logo size="lg" />
      </div>

      {#if verifying}
        <!-- Verifying State -->
        <div class="mb-8">
          <div class="relative inline-flex">
            <div
              class="absolute inset-0 bg-blue-500/20 rounded-full blur-xl animate-pulse">
            </div>
            <div
              class="relative bg-gradient-to-br from-blue-400 to-blue-600 p-6 rounded-full">
              <Loader2 class="w-12 h-12 text-white animate-spin" />
            </div>
          </div>
        </div>

        <h1 class="text-3xl font-bold text-gray-900 mb-3">
          Verifying Your Email
        </h1>
        <p class="text-lg text-gray-600">
          Please wait while we verify your email address...
        </p>
      {:else if verified}
        <!-- Success State -->
        <div class="mb-8 success-icon-container">
          <div class="relative inline-flex">
            <div
              class="absolute inset-0 bg-green-500/20 rounded-full blur-xl animate-pulse">
            </div>
            <div
              class="relative bg-gradient-to-br from-green-400 to-green-600 p-6 rounded-full">
              <CheckCircle class="w-12 h-12 text-white" />
            </div>
          </div>
        </div>

        <h1 class="text-3xl font-bold text-gray-900 mb-3">Email Verified!</h1>
        <p class="text-lg text-gray-600 mb-6">
          Your email has been successfully verified.
        </p>
        <p class="text-sm text-gray-500">Redirecting to login...</p>
      {:else}
        <!-- Error State -->
        <div class="mb-8">
          <div class="relative inline-flex">
            <div
              class="absolute inset-0 bg-red-500/20 rounded-full blur-xl animate-pulse">
            </div>
            <div
              class="relative bg-gradient-to-br from-red-400 to-red-600 p-6 rounded-full">
              <XCircle class="w-12 h-12 text-white" />
            </div>
          </div>
        </div>

        <h1 class="text-3xl font-bold text-gray-900 mb-3">
          Verification Failed
        </h1>
        <p class="text-lg text-gray-600 mb-6">
          {error}
        </p>

        <div class="space-y-3">
          <a
            href="/login"
            class="inline-flex items-center justify-center w-full px-6 py-3 bg-gradient-to-r from-blue-600 to-purple-600 text-white rounded-lg font-medium hover:shadow-lg transition-all duration-300 hover:scale-[1.02]">
            Go to Login
          </a>

          <a
            href="/register"
            class="inline-flex items-center justify-center w-full px-6 py-3 bg-white text-gray-700 rounded-lg font-medium border border-gray-300 hover:bg-gray-50 transition-all duration-200">
            Create New Account
          </a>
        </div>
      {/if}
    </div>
  </div>
</div>

<style>
  .success-icon-container {
    animation: bounce-in 0.6s cubic-bezier(0.68, -0.55, 0.265, 1.55);
  }

  @keyframes bounce-in {
    0% {
      transform: scale(0);
      opacity: 0;
    }
    50% {
      transform: scale(1.1);
    }
    100% {
      transform: scale(1);
      opacity: 1;
    }
  }
</style>
