<script lang="ts">
  import { createEventDispatcher } from 'svelte';

  export let value = '';
  export let placeholder = 'Search...';
  export let size: 'sm' | 'md' | 'lg' = 'md';
  let className = '';
  export { className as class };

  const dispatch = createEventDispatcher();

  $: sizeClasses = {
    sm: 'pl-8 pr-3 py-2 text-sm',
    md: 'pl-10 pr-4 py-3',
    lg: 'pl-12 pr-5 py-4 text-lg'
  }[size];

  $: iconSizeClasses = {
    sm: 'h-4 w-4',
    md: 'h-5 w-5', 
    lg: 'h-6 w-6'
  }[size];

  $: iconPositionClasses = {
    sm: 'pl-2',
    md: 'pl-3',
    lg: 'pl-4'
  }[size];

  function handleInput(event: Event) {
    const target = event.target as HTMLInputElement;
    value = target.value;
    dispatch('input', { value });
  }
</script>

<div class="relative flex-1 {className}">
  <div class="absolute inset-y-0 left-0 {iconPositionClasses} flex items-center pointer-events-none">
    <svg class="{iconSizeClasses} text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
      <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M21 21l-6-6m2-5a7 7 0 11-14 0 7 7 0 0114 0z" />
    </svg>
  </div>
  <input
    type="text"
    bind:value
    on:input={handleInput}
    {placeholder}
    class="w-full {sizeClasses} border border-gray-200 rounded-xl bg-white/50 backdrop-blur-sm focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500 transition-all duration-200"
  />
</div>