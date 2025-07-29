<script lang="ts">
  import { isSubscribed } from '$lib/stores/subscription';
  import { AlertCircle } from 'lucide-svelte';

  export let redirectTo = '/pricing';
  export let message = 'You need an active subscription to access this feature';
  export let showUpgradeButton = true;
  export let showAlert = false;
</script>

{#if $isSubscribed}
  <slot />
{:else if showAlert}
  <div class="flex flex-col items-center justify-center p-8">
    <div class="rounded-lg border border-red-200 bg-red-50 p-6 text-center">
      <AlertCircle class="mx-auto mb-4 h-12 w-12 text-red-500" />
      <h3 class="mb-2 text-lg font-semibold text-gray-900">
        Subscription Required
      </h3>
      <p class="mb-4 text-gray-600">{message}</p>
      {#if showUpgradeButton}
        <a
          href={redirectTo}
          class="inline-flex items-center justify-center rounded-md bg-primary px-4 py-2 text-sm font-medium text-primary-foreground hover:bg-primary/90">
          Upgrade Your Plan
        </a>
      {/if}
    </div>
    {#if $$slots.fallback}
      <slot name="fallback" />
    {/if}
  </div>
{/if}
