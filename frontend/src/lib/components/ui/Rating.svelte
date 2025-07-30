<script lang="ts">
  import { Star } from 'lucide-svelte';
  let {
    value = $bindable(0),
    max = 5,
    size = 'md',
    readonly = false,
    showLabel = false,
  }: {
    value?: number;
    max?: number;
    size?: 'sm' | 'md' | 'lg';
    readonly?: boolean;
    showLabel?: boolean;
  } = $props();

  const sizes = {
    sm: 'h-4 w-4',
    md: 'h-6 w-6',
    lg: 'h-8 w-8',
  };

  const labels = ['Poor', 'Fair', 'Good', 'Great', 'Excellent'];

  function handleClick(rating: number) {
    if (!readonly) {
      value = rating;
    }
  }

  let stars = $derived(
    Array(max)
      .fill(0)
      .map((_, i) => ({
        filled: i < value,
        index: i + 1,
      }))
  );
</script>

<div class="inline-flex items-center space-x-1">
  {#each stars as star}
    <button
      type="button"
      class="focus:outline-none transition-transform {!readonly
        ? 'hover:scale-110 cursor-pointer'
        : 'cursor-default'}"
      onclick={() => handleClick(star.index)}
      disabled={readonly}
      aria-label="Rate {star.index} out of {max} stars">
      <Star
        class="{sizes[size]} {star.filled
          ? 'text-yellow-400'
          : 'text-gray-300'}"
        fill={star.filled ? 'currentColor' : 'none'} />
    </button>
  {/each}

  {#if showLabel && value > 0}
    <span class="ml-2 text-sm text-gray-600">
      {labels[value - 1] || ''}
    </span>
  {/if}
</div>
