<script lang="ts">
  export let size: 'sm' | 'md' | 'lg' | 'xl' = 'md';
  export let variant: 'default' | 'white' | 'dark' = 'default';
  export let showText = true;
  export let className = '';

  const sizes = {
    sm: { icon: 'h-6 w-6', text: 'text-sm' },
    md: { icon: 'h-8 w-8', text: 'text-lg' },
    lg: { icon: 'h-10 w-10', text: 'text-xl' },
    xl: { icon: 'h-12 w-12', text: 'text-2xl' }
  };

  const variants = {
    default: {
      icon: 'from-blue-600 to-purple-600',
      text: 'from-gray-900 to-gray-700',
      iconColor: 'text-white'
    },
    white: {
      icon: 'from-white/90 to-white/70',
      text: 'text-white',
      iconColor: 'text-blue-600'
    },
    dark: {
      icon: 'from-gray-800 to-gray-900',
      text: 'text-gray-900',
      iconColor: 'text-white'
    }
  };

  $: currentSize = sizes[size];
  $: currentVariant = variants[variant];
</script>

<div class="flex items-center space-x-3 {className}">
  <!-- Logo Icon -->
  <div class="relative">
    <div class="bg-gradient-to-br {currentVariant.icon} rounded-xl {currentSize.icon} flex items-center justify-center shadow-lg transition-transform duration-200 hover:scale-105">
      <!-- Clean creative LC monogram -->
      <svg class="{currentVariant.iconColor} {size === 'sm' ? 'h-3 w-3' : size === 'md' ? 'h-4 w-4' : size === 'lg' ? 'h-5 w-5' : 'h-6 w-6'}" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5">
        <!-- Letter L with subtle fork styling -->
        <path d="M5 4v14h8" stroke-linecap="round" stroke-linejoin="round"/>
        <!-- Simple fork prongs -->
        <path d="M4 4v5M5 4v5M6 4v5" stroke-linecap="round" stroke-linejoin="round"/>
        
        <!-- Letter C with elegant thickness variation -->
        <path d="M20 8c-2.5 0-4 2-4 4s1.5 4 4 4" stroke-linecap="round" stroke-linejoin="round" stroke-width="3"/>
      </svg>
    </div>
    
    <!-- Subtle glow effect -->
    <div class="absolute inset-0 bg-gradient-to-br {currentVariant.icon} rounded-xl opacity-20 blur-sm -z-10"></div>
  </div>

  <!-- Logo Text -->
  {#if showText}
    <div class="font-bold bg-gradient-to-r {currentVariant.text} bg-clip-text text-transparent {currentSize.text} tracking-tight">
      LeCritique
    </div>
  {/if}
</div>