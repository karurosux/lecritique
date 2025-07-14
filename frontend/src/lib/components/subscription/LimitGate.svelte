<script lang="ts">
  import {
    getLimit,
    isUnlimited,
    planFeatures,
  } from "$lib/stores/subscription";
  import type { Snippet } from "svelte";

  interface Props {
    limit: string;
    currentCount: number;
    children: Snippet;
    fallback?: Snippet;
    showLimitMessage?: boolean;
  }

  let {
    limit,
    currentCount,
    children,
    fallback,
    showLimitMessage = true,
  }: Props = $props();

  let limitValue = $derived($getLimit(limit));
  let unlimited = $derived($isUnlimited(limit));
  let canAdd = $derived(unlimited || currentCount < limitValue);
  let loading = $derived($planFeatures === null);
</script>

{#if loading}
  <!-- Loading state while checking limits -->
  <div class="animate-pulse">
    <div class="h-8 bg-gray-200 rounded"></div>
  </div>
{:else if canAdd}
  {@render children()}
{:else if fallback}
  {@render fallback()}
{:else if showLimitMessage}
  <div class="p-4 border border-yellow-200 bg-yellow-50 rounded-lg">
    <p class="text-sm text-yellow-800">
      You've reached your plan limit of {limitValue}
      {limit.replace(/_/g, " ")}.
      <a href="/subscription" class="font-medium underline"
        >Upgrade to add more</a
      >
    </p>
  </div>
{/if}
