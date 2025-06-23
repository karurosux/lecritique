<script lang="ts">
  import { createEventDispatcher } from 'svelte';

  export let value = '';
  export let options: Array<{ value: string; label: string }> = [];
  export let size: 'sm' | 'md' | 'lg' = 'md';
  export let minWidth = 'min-w-32';
  let className = '';
  export { className as class };

  const dispatch = createEventDispatcher();

  $: sizeClasses = {
    sm: 'px-3 py-2 pr-8 text-sm',
    md: 'px-4 py-3 pr-10',
    lg: 'px-5 py-4 pr-12 text-lg'
  }[size];

  $: iconSizeClasses = {
    sm: 'h-3 w-3',
    md: 'h-4 w-4',
    lg: 'h-5 w-5'
  }[size];

  $: iconPositionClasses = {
    sm: 'pr-2',
    md: 'pr-3',
    lg: 'pr-4'
  }[size];

  function handleChange(event: Event) {
    const target = event.target as HTMLSelectElement;
    value = target.value;
    dispatch('change', { value });
  }
</script>

<div class="relative {className}">
  <select
    bind:value
    on:change={handleChange}
    class="appearance-none {sizeClasses} border border-gray-200 rounded-xl bg-white/50 backdrop-blur-sm focus:outline-none focus:ring-2 focus:ring-blue-500 cursor-pointer {minWidth}"
  >
    {#each options as option}
      <option value={option.value}>{option.label}</option>
    {/each}
  </select>
  <div class="absolute inset-y-0 right-0 flex items-center {iconPositionClasses} pointer-events-none">
    <svg class="{iconSizeClasses} text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 9l-7 7-7-7" />
    </svg>
  </div>
</div>