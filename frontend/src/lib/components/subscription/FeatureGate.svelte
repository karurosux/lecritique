<script lang="ts">
  import { hasFeature, subscription } from "$lib/stores/subscription";
  import type { Snippet } from "svelte";

  interface Props {
    feature: string;
    children: Snippet;
    fallback?: Snippet;
    showUpgradePrompt?: boolean;
  }

  let {
    feature,
    children,
    fallback,
    showUpgradePrompt = false,
  }: Props = $props();

  let hasAccess = $derived($hasFeature(feature));
  let loading = $derived($subscription.isLoading);
</script>

{#if loading}
  <!-- Loading state while checking features -->
  <div class="animate-pulse">
    <div class="h-8 bg-gray-200 rounded"></div>
  </div>
{:else if hasAccess}
  {@render children()}
{:else if fallback}
  {@render fallback()}
{:else if showUpgradePrompt}
  <div class="p-4 border border-orange-200 bg-orange-50 rounded-lg">
    <p class="text-sm text-orange-800">
      This feature requires an upgraded plan.
      <a href="/subscription" class="font-medium underline">Upgrade now</a>
    </p>
  </div>
{/if}
