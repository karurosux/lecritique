<script lang="ts">
  let {
    isOpen = $bindable(false),
    open = $bindable(false), // Backward compatibility
    title = '',
    showClose = true,
    size = 'md',
    onclose = () => {},
    children
  }: {
    isOpen?: boolean;
    open?: boolean;
    title?: string;
    showClose?: boolean;
    size?: 'sm' | 'md' | 'lg' | 'xl';
    onclose?: () => void;
    children?: any;
  } = $props();
  
  // Support both isOpen and open props for compatibility
  let modalOpen = $derived(isOpen || open);
  
  const sizeClasses = {
    sm: 'max-w-sm',
    md: 'max-w-md',
    lg: 'max-w-2xl',
    xl: 'max-w-4xl'
  };
  
  function closeModal() {
    isOpen = false;
    open = false;
    onclose();
  }
  
  function handleBackdropClick(event: MouseEvent) {
    if (event.target === event.currentTarget) {
      closeModal();
    }
  }
  
  function handleKeydown(event: KeyboardEvent) {
    if (event.key === 'Escape') {
      closeModal();
    }
  }
</script>

{#if modalOpen}
  <div 
    class="fixed inset-0 z-50 flex items-center justify-center bg-black bg-opacity-50"
    onclick={handleBackdropClick}
    onkeydown={handleKeydown}
    role="dialog"
    aria-modal="true"
    tabindex="-1"
  >
    <div class="relative bg-white rounded-lg shadow-xl {sizeClasses[size]} w-full mx-4 max-h-screen overflow-y-auto">
      {#if title || showClose}
        <div class="flex items-center justify-between p-4 border-b border-gray-200">
          {#if title}
            <h3 class="text-lg font-semibold text-gray-900">{title}</h3>
          {/if}
          {#if showClose}
            <button
              type="button"
              class="text-gray-400 hover:text-gray-600 transition-colors cursor-pointer"
              onclick={closeModal}
              aria-label="Close modal"
            >
              <svg class="w-6 h-6" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M6 18L18 6M6 6l12 12" />
              </svg>
            </button>
          {/if}
        </div>
      {/if}
      
      <div class="p-4">
        {@render children?.()}
      </div>
    </div>
  </div>
{/if}