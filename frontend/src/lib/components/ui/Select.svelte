<script lang="ts">
  import { ChevronDown } from 'lucide-svelte';
  let {
    value = $bindable(''),
    options = [],
    size = 'md',
    minWidth = 'min-w-32',
    placeholder = '',
    onchange = (value: string) => {},
    ...restProps
  }: {
    value?: string;
    options?: Array<{ value: string; label: string }>;
    size?: 'sm' | 'md' | 'lg';
    minWidth?: string;
    placeholder?: string;
    onchange?: (value: string) => void;
    class?: string;
    [key: string]: any;
  } = $props();

  let className = restProps.class || '';

  const sizeClasses = {
    sm: 'px-3 py-2 pr-8 text-sm',
    md: 'px-4 py-3 pr-10',
    lg: 'px-5 py-4 pr-12 text-lg',
  };

  const iconSizeClasses = {
    sm: 'h-3 w-3',
    md: 'h-4 w-4',
    lg: 'h-5 w-5',
  };

  const iconPositionClasses = {
    sm: 'pr-2',
    md: 'pr-3',
    lg: 'pr-4',
  };

  function handleChange(event: Event) {
    const target = event.target as HTMLSelectElement;
    value = target.value;
    onchange(value);
  }
</script>

<div class="relative {className}">
  <select
    bind:value
    onchange={handleChange}
    class="appearance-none w-full {sizeClasses[
      size
    ]} border border-gray-200 rounded-xl bg-white/50 backdrop-blur-sm hover:shadow-md focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-transparent cursor-pointer transition-all duration-200 {!value
      ? 'text-gray-400'
      : ''}"
    {...restProps}>
    {#if placeholder && !value}
      <option value="" disabled selected>{placeholder}</option>
    {/if}
    {#each options as option}
      <option value={option.value}>{option.label}</option>
    {/each}
  </select>
  <div
    class="absolute inset-y-0 right-0 flex items-center {iconPositionClasses[
      size
    ]} pointer-events-none">
    <ChevronDown class="{iconSizeClasses[size]} text-gray-400" />
  </div>
</div>
