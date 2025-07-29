<script lang="ts">
  import { QrCode } from 'lucide-svelte';

  let {
    size = 'md',
    variant = 'default',
    showText = true,
    ...restProps
  }: {
    size?: 'sm' | 'md' | 'lg' | 'xl';
    variant?: 'default' | 'white' | 'dark';
    showText?: boolean;
    class?: string;
    [key: string]: any;
  } = $props();

  let className = restProps.class || '';

  const sizes = {
    sm: { icon: 'h-7 w-7', text: 'text-lg', lucide: 16 },
    md: { icon: 'h-9 w-9', text: 'text-xl', lucide: 20 },
    lg: { icon: 'h-11 w-11', text: 'text-2xl', lucide: 24 },
    xl: { icon: 'h-13 w-13', text: 'text-3xl', lucide: 28 },
  };

  const variants = {
    default: {
      containerBg: 'from-primary-500 to-primary-700',
      text: 'from-primary-600 to-primary-700',
    },
    white: {
      containerBg: 'from-primary-500 to-primary-700',
      text: 'from-primary-600 to-primary-700',
    },
    dark: {
      containerBg: 'from-primary-400 to-primary-600',
      text: 'from-primary-300 to-primary-500',
    },
  };

  let currentSize = $derived(sizes[size]);
  let currentVariant = $derived(variants[variant]);
</script>

<div class="flex items-center gap-2 {className}">
  <!-- Logo Icon -->
  <div class="relative">
    <!-- Main icon container with gradient -->
    <div
      class="relative {currentSize.icon} bg-gradient-to-r from-purple-600 via-blue-600 to-blue-700 rounded-lg flex items-center justify-center p-1 shadow-md">
      <!-- Inner background -->
      <div class="absolute inset-1 bg-white rounded"></div>

      <!-- QR Code icon -->
      <QrCode
        size={currentSize.lucide}
        strokeWidth={2.5}
        class="text-black relative z-10" />
    </div>
  </div>

  <!-- Logo Text -->
  {#if showText}
    <div
      class="font-black bg-gradient-to-r from-purple-600 via-blue-600 to-blue-700 bg-clip-text text-transparent {currentSize.text} tracking-tight select-none">
      Kyooar
    </div>
  {/if}
</div>
