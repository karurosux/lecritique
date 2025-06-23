<script lang="ts">
  let {
    value = $bindable(0),
    max = 5,
    size = 'md',
    readonly = false,
    showLabel = false
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
    lg: 'h-8 w-8'
  };

  const labels = ['Poor', 'Fair', 'Good', 'Great', 'Excellent'];

  function handleClick(rating: number) {
    if (!readonly) {
      value = rating;
    }
  }

  let stars = $derived(Array(max).fill(0).map((_, i) => ({
    filled: i < value,
    index: i + 1
  })));
</script>

<div class="inline-flex items-center space-x-1">
  {#each stars as star}
    <button
      type="button"
      class="focus:outline-none transition-transform {!readonly ? 'hover:scale-110 cursor-pointer' : 'cursor-default'}"
      onclick={() => handleClick(star.index)}
      disabled={readonly}
    >
      <svg
        class="{sizes[size]} {star.filled ? 'text-yellow-400 fill-current' : 'text-gray-300'}"
        fill="none"
        stroke="currentColor"
        viewBox="0 0 24 24"
      >
        <path
          stroke-linecap="round"
          stroke-linejoin="round"
          stroke-width="2"
          d="M11.049 2.927c.3-.921 1.603-.921 1.902 0l1.519 4.674a1 1 0 00.95.69h4.915c.969 0 1.371 1.24.588 1.81l-3.976 2.888a1 1 0 00-.363 1.118l1.518 4.674c.3.922-.755 1.688-1.538 1.118l-3.976-2.888a1 1 0 00-1.176 0l-3.976 2.888c-.783.57-1.838-.197-1.538-1.118l1.518-4.674a1 1 0 00-.363-1.118l-3.976-2.888c-.784-.57-.38-1.81.588-1.81h4.914a1 1 0 00.951-.69l1.519-4.674z"
        />
      </svg>
    </button>
  {/each}
  
  {#if showLabel && value > 0}
    <span class="ml-2 text-sm text-gray-600">
      {labels[value - 1] || ''}
    </span>
  {/if}
</div>