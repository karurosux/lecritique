<script lang="ts">
  import { auth } from '$lib/stores/auth';
  import { goto } from '$app/navigation';
  import { onMount } from 'svelte';
  import { Loader2 } from 'lucide-svelte';

  let authState = $derived($auth);
  let hasCheckedAuth = $state(false);

  $effect(() => {
    if (!hasCheckedAuth) {
      hasCheckedAuth = true;
      if (!authState.isAuthenticated) {
        goto('/login');
      }
    }
  });

  // Watch for auth state changes
  $effect(() => {
    if (
      !authState.isAuthenticated &&
      typeof window !== 'undefined' &&
      hasCheckedAuth
    ) {
      goto('/login');
    }
  });
</script>

{#if authState.isAuthenticated}
  <slot />
{:else}
  <div class="min-h-screen bg-gray-50 flex items-center justify-center">
    <div class="text-center">
      <Loader2 class="animate-spin h-8 w-8 text-blue-600 mx-auto mb-4" />
      <p class="text-gray-600">Redirecting...</p>
    </div>
  </div>
{/if}
