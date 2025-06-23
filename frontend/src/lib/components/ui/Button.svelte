<script lang="ts">
  type ButtonVariant = 'primary' | 'secondary' | 'outline' | 'ghost' | 'destructive' | 'gradient' | 'glass';
  type ButtonSize = 'sm' | 'md' | 'lg' | 'xl';

  let {
    variant = 'primary',
    size = 'md',
    disabled = false,
    type = 'button',
    href = undefined,
    loading = false,
    onclick = undefined,
    children,
    ...restProps
  }: {
    variant?: ButtonVariant;
    size?: ButtonSize;
    disabled?: boolean;
    type?: 'button' | 'submit' | 'reset';
    href?: string;
    loading?: boolean;
    onclick?: ((event: MouseEvent) => void) | undefined;
    children?: any;
    class?: string;
    [key: string]: any;
  } = $props();
  
  let className = restProps.class || '';

  const variants = {
    primary: 'bg-gradient-to-r from-blue-600 to-blue-700 text-white shadow-lg shadow-blue-500/25 hover:shadow-xl hover:shadow-blue-500/40 hover:from-blue-700 hover:to-blue-800 focus-visible:ring-blue-500',
    secondary: 'bg-gradient-to-r from-gray-100 to-gray-200 text-gray-900 shadow-sm hover:from-gray-200 hover:to-gray-300 hover:shadow-md focus-visible:ring-gray-400',
    outline: 'border-2 border-gray-300 bg-white/50 backdrop-blur-sm text-gray-700 hover:border-gray-400 hover:bg-gray-50 hover:shadow-md focus-visible:ring-gray-400',
    ghost: 'text-gray-700 hover:bg-gray-100/80 hover:text-gray-900 focus-visible:ring-gray-400 backdrop-blur-sm',
    destructive: 'bg-gradient-to-r from-red-600 to-red-700 text-white shadow-lg shadow-red-500/25 hover:shadow-xl hover:shadow-red-500/40 hover:from-red-700 hover:to-red-800 focus-visible:ring-red-500',
    gradient: 'bg-gradient-to-r from-purple-600 via-blue-600 to-blue-700 text-white shadow-lg shadow-blue-500/30 hover:shadow-xl hover:shadow-blue-500/50 hover:from-purple-700 hover:via-blue-700 hover:to-blue-800 focus-visible:ring-blue-500',
    glass: 'bg-white/20 backdrop-blur-xl border border-white/30 text-gray-900 shadow-lg hover:bg-white/30 hover:shadow-xl focus-visible:ring-blue-500'
  };

  const sizes = {
    sm: 'px-4 py-2 text-sm font-medium',
    md: 'px-6 py-2.5 text-sm font-semibold',
    lg: 'px-8 py-3 text-base font-semibold',
    xl: 'px-10 py-4 text-lg font-semibold'
  };

  const baseClasses = 'cursor-pointer relative inline-flex items-center justify-center rounded-xl font-medium transition-all duration-300 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-offset-2 disabled:pointer-events-none disabled:opacity-50 disabled:cursor-not-allowed transform hover:scale-105 active:scale-95';
</script>

{#if href}
  <a
    {href}
    class="{baseClasses} {variants[variant]} {sizes[size]} {className}"
    class:opacity-50={disabled}
    class:pointer-events-none={disabled}
  >
    {#if loading}
      <svg class="animate-spin -ml-1 mr-3 h-4 w-4" fill="none" viewBox="0 0 24 24">
        <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
        <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
      </svg>
    {/if}
    {@render children?.()}
  </a>
{:else}
  <button
    {type}
    {disabled}
    class="{baseClasses} {variants[variant]} {sizes[size]} {className}"
    {onclick}
  >
    {#if loading}
      <svg class="animate-spin -ml-1 mr-3 h-4 w-4" fill="none" viewBox="0 0 24 24">
        <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
        <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 714 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
      </svg>
    {/if}
    {@render children?.()}
  </button>
{/if}