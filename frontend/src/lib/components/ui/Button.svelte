<script lang="ts">
  type ButtonVariant = 'primary' | 'secondary' | 'outline' | 'ghost' | 'destructive';
  type ButtonSize = 'sm' | 'md' | 'lg';

  export let variant: ButtonVariant = 'primary';
  export let size: ButtonSize = 'md';
  export let disabled = false;
  export let type: 'button' | 'submit' | 'reset' = 'button';
  export let href: string | undefined = undefined;
  
  let className = '';
  export { className as class };

  const variants = {
    primary: 'bg-blue-600 text-white hover:bg-blue-700 focus-visible:ring-blue-600',
    secondary: 'bg-gray-100 text-gray-900 hover:bg-gray-200 focus-visible:ring-gray-500',
    outline: 'border border-gray-300 bg-transparent hover:bg-gray-50 focus-visible:ring-gray-500',
    ghost: 'hover:bg-gray-100 focus-visible:ring-gray-500',
    destructive: 'bg-red-600 text-white hover:bg-red-700 focus-visible:ring-red-600'
  };

  const sizes = {
    sm: 'px-3 py-1.5 text-sm',
    md: 'px-4 py-2 text-sm',
    lg: 'px-6 py-3 text-base'
  };

  const baseClasses = 'inline-flex items-center justify-center rounded-md font-medium transition-colors focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-offset-2 disabled:pointer-events-none disabled:opacity-50';
</script>

{#if href}
  <a
    {href}
    class="{baseClasses} {variants[variant]} {sizes[size]} {className}"
    class:opacity-50={disabled}
    class:pointer-events-none={disabled}
  >
    <slot />
  </a>
{:else}
  <button
    {type}
    {disabled}
    class="{baseClasses} {variants[variant]} {sizes[size]} {className}"
    on:click
  >
    <slot />
  </button>
{/if}