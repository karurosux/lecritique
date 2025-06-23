<script lang="ts">
  let {
    variant = 'default',
    padding = true,
    hover = false,
    interactive = false,
    children,
    ...restProps
  }: {
    variant?: 'default' | 'elevated' | 'glass' | 'gradient' | 'minimal';
    padding?: boolean;
    hover?: boolean;
    interactive?: boolean;
    children?: any;
    class?: string;
    [key: string]: any;
  } = $props();
  
  let className = restProps.class || '';
  
  const baseClasses = 'relative overflow-hidden transition-all duration-300 ease-out';
  
  const variantClasses = {
    default: 'bg-white rounded-xl border border-gray-300 shadow-lg shadow-gray-900/10 backdrop-blur-sm',
    elevated: 'bg-white rounded-2xl shadow-xl shadow-gray-900/15 border border-gray-200',
    glass: 'bg-white/95 backdrop-blur-xl rounded-2xl border border-gray-200/60 shadow-xl shadow-gray-900/20',
    gradient: 'bg-gradient-to-br from-white to-gray-50 rounded-2xl border border-gray-300/70 shadow-xl shadow-blue-500/15',
    minimal: 'bg-gray-50 rounded-xl border border-gray-300 shadow-md shadow-gray-900/8'
  };
  
  let paddingClasses = $derived(padding ? 'p-6 lg:p-8' : '');
  
  let hoverClasses = $derived(hover ? 'hover:shadow-2xl hover:shadow-gray-900/20 hover:-translate-y-1 hover:scale-[1.02] hover:border-gray-400' : '');
  
  let interactiveClasses = $derived(interactive ? 'cursor-pointer group' : '');
</script>

<div class="{baseClasses} {variantClasses[variant]} {paddingClasses} {hoverClasses} {interactiveClasses} {className}">
  {#if variant === 'gradient'}
    <div class="absolute inset-0 bg-gradient-to-br from-blue-500/5 via-purple-500/5 to-pink-500/5 rounded-2xl"></div>
  {/if}
  
  <div class="relative z-10">
    {@render children?.()}
  </div>
  
  {#if interactive}
    <div class="absolute inset-0 bg-gradient-to-r from-blue-500/0 via-purple-500/0 to-pink-500/0 group-hover:from-blue-500/5 group-hover:via-purple-500/5 group-hover:to-pink-500/5 transition-all duration-500 rounded-xl"></div>
  {/if}
</div>