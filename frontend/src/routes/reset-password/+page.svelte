<script lang="ts">
  import { goto } from '$app/navigation';
  import { page } from '$app/stores';
  import { Api } from '$lib/api/api';
  import { Button, Card, Input, Logo } from '$lib/components/ui';
  import { Check, Lock, XCircle, Loader2, Key, LogIn } from 'lucide-svelte';

  const api = new Api({
    baseURL: import.meta.env.VITE_API_URL || 'http://localhost:8080',
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
        new_password: newPassword,
      });
      success = true;
      setTimeout(() => {
        goto('/login');
      }, 3000);
    } catch (err: any) {
      error =
        err.response?.data?.message ||
        'Failed to reset password. The link may have expired.';
    } finally {
      loading = false;
    }
  }
</script>

<svelte:head>
  <title>Reset Password - Kyooar</title>
  <meta
    name="description"
    content="Create a new password for your Kyooar account" />
</svelte:head>

<div class="min-h-screen flex flex-col justify-center py-12 sm:px-6 lg:px-8">
  <div class="relative z-10 sm:mx-auto sm:w-full sm:max-w-md">
    <div class="flex justify-center mb-8">
      <Logo size="xl" />
    </div>

    <div class="text-center space-y-3">
      <h2
        class="text-4xl font-bold bg-gradient-to-r from-gray-900 via-blue-800 to-purple-800 bg-clip-text text-transparent">
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
          <div
            class="mx-auto flex items-center justify-center h-16 w-16 rounded-full bg-gradient-to-r from-green-400 to-green-600 shadow-lg">
            <Check class="h-8 w-8 text-white" />
          </div>
          <div class="space-y-2">
            <h3 class="text-2xl font-bold text-gray-900">Password Reset!</h3>
            <p class="text-gray-600 max-w-sm mx-auto">
              Your password has been successfully updated. Redirecting to
              login...
            </p>
          </div>
          <Button
            href="/login"
            variant="gradient"
            size="lg"
            class="w-full shadow-lg hover:shadow-xl transition-all duration-300">
            <LogIn class="w-5 h-5 mr-2" />
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
            minlength={8} />

          <Input
            type="password"
            label="Confirm new password"
            id="confirm-password"
            bind:value={confirmPassword}
            required
            placeholder="Confirm new password"
            autocomplete="new-password"
            disabled={loading}
            minlength={8} />

          <p class="text-xs text-gray-500">
            Password must be at least 8 characters long
          </p>

          {#if error}
            <div class="bg-red-50 border border-red-200 rounded-md p-4">
              <div class="flex">
                <div class="flex-shrink-0">
                  <XCircle class="h-5 w-5 text-red-400" />
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
            class="w-full shadow-lg hover:shadow-xl transition-all duration-300">
            {#if loading}
              <Loader2 class="animate-spin -ml-1 mr-3 h-5 w-5 text-white" />
              Resetting password...
            {:else}
              <Key class="w-5 h-5 mr-2" />
              Reset Password
            {/if}
          </Button>
        </form>
      </Card>
    {/if}
  </div>
</div>
