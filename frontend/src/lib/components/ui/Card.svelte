<script lang="ts">
  export let variant: 'default' | 'elevated' | 'glass' | 'gradient' | 'minimal' = 'default';
  export let padding = true;
  export let hover = false;
  export let interactive = false;
  
  $: baseClasses = 'relative overflow-hidden transition-all duration-300 ease-out';
  
  $: variantClasses = {
    default: 'bg-white rounded-xl border border-gray-200/60 shadow-sm backdrop-blur-sm',
    elevated: 'bg-white rounded-2xl shadow-lg shadow-gray-900/5 border border-gray-100',
    glass: 'bg-white/80 backdrop-blur-xl rounded-2xl border border-white/20 shadow-xl shadow-gray-900/10',
    gradient: 'bg-gradient-to-br from-white to-gray-50/50 rounded-2xl border border-gray-200/50 shadow-lg shadow-blue-500/5',
    minimal: 'bg-gray-50/50 rounded-xl border border-gray-100'
  }[variant];
  
  $: paddingClasses = padding ? 'p-6 lg:p-8' : '';
  
  $: hoverClasses = hover ? 'hover:shadow-xl hover:shadow-gray-900/10 hover:-translate-y-1 hover:scale-[1.02]' : '';
  
  $: interactiveClasses = interactive ? 'cursor-pointer group' : '';
</script>

<div class="{baseClasses} {variantClasses} {paddingClasses} {hoverClasses} {interactiveClasses}">
  {#if variant === 'gradient'}
    <div class="absolute inset-0 bg-gradient-to-br from-blue-500/5 via-purple-500/5 to-pink-500/5 rounded-2xl"></div>
  {/if}
  
  <div class="relative z-10">
    <slot />
  </div>
  
  {#if interactive}
    <div class="absolute inset-0 bg-gradient-to-r from-blue-500/0 via-purple-500/0 to-pink-500/0 group-hover:from-blue-500/5 group-hover:via-purple-500/5 group-hover:to-pink-500/5 transition-all duration-500 rounded-xl"></div>
  {/if}
</div>