<script lang="ts">
  import { BarChart3 } from 'lucide-svelte';
  import { Card } from '$lib/components/ui';
  import type { ComponentType } from 'svelte';

  interface Props {
    title?: string;
    description?: string;
    icon?: ComponentType;
    class?: string;
    variant?: 'card' | 'inline';
  }

  let {
    title = 'No Data Available',
    description = "There's nothing to show here right now",
    icon = BarChart3,
    class: className = '',
    variant = 'card',
  }: Props = $props();
</script>

{#if variant === 'card'}
  <Card variant="glass" class="p-12 text-center {className}">
    <div
      class="w-24 h-24 mx-auto bg-gradient-to-br from-blue-500 to-purple-600 rounded-full flex items-center justify-center mb-6 shadow-lg shadow-blue-500/25">
      <svelte:component this={icon} class="h-12 w-12 text-white" />
    </div>
    <h3
      class="text-xl font-bold mb-2 bg-gradient-to-r from-gray-900 to-gray-700 bg-clip-text text-transparent">
      {title}
    </h3>
    <p class="text-gray-600 mb-6 max-w-md mx-auto">
      {description}
    </p>
  </Card>
{:else}
  <div class="p-12 text-center {className}">
    <div
      class="w-24 h-24 mx-auto bg-gradient-to-br from-blue-500 to-purple-600 rounded-full flex items-center justify-center mb-6 shadow-lg shadow-blue-500/25">
      <svelte:component this={icon} class="h-12 w-12 text-white" />
    </div>
    <h3
      class="text-xl font-bold mb-2 bg-gradient-to-r from-gray-900 to-gray-700 bg-clip-text text-transparent">
      {title}
    </h3>
    <p class="text-gray-600 mb-6 max-w-md mx-auto">
      {description}
    </p>
  </div>
{/if}
