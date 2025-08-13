<script lang="ts">
  import { QrCode } from 'lucide-svelte';

  let {
    size = 'md',
    showText = true,
    ...restProps
  }: {
    size?: 'sm' | 'md' | 'lg' | 'xl';
    showText?: boolean;
    class?: string;
    [key: string]: any;
  } = $props();

  let className = restProps.class || '';

  const iconSizes = {
    sm: 16,
    md: 20,
    lg: 24,
    xl: 28,
  };

  const textSizes = {
    sm: 'text-lg',
    md: 'text-xl',
    lg: 'text-2xl',
    xl: 'text-3xl',
  };

  const containerSizes = {
    sm: 'h-7 w-7',
    md: 'h-9 w-9',
    lg: 'h-11 w-11',
    xl: 'h-13 w-13',
  };

  let iconSize = $derived(iconSizes[size]);
  let textSize = $derived(textSizes[size]);
  let containerSize = $derived(containerSizes[size]);
</script>

<div class="flex items-center gap-2 {className}">
  <div class="relative">
    
    <div
      class="relative {containerSize} bg-gradient-to-r from-purple-600 via-blue-600 to-blue-700 rounded-lg flex items-center justify-center shadow-md">
      
      <QrCode size={iconSize} strokeWidth={2.5} class="text-white" />
    </div>
  </div>

  {#if showText}
    <span
      class="font-black {textSize} bg-gradient-to-r from-purple-600 via-blue-600 to-blue-700 bg-clip-text text-transparent tracking-tight select-none">
      Kyooar
    </span>
  {/if}
</div>
