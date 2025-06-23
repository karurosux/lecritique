<script lang="ts">
  import { createEventDispatcher } from 'svelte';
  
  export let open = false;
  export let title = '';
  export let showClose = true;
  
  const dispatch = createEventDispatcher();
  
  function closeModal() {
    dispatch('close');
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

{#if open}
  <div 
    class="fixed inset-0 z-50 flex items-center justify-center bg-black bg-opacity-50"
    on:click={handleBackdropClick}
    on:keydown={handleKeydown}
    role="dialog"
    aria-modal="true"
    tabindex="-1"
  >
    <div class="relative bg-white rounded-lg shadow-xl max-w-md w-full mx-4 max-h-screen overflow-y-auto">
      {#if title || showClose}
        <div class="flex items-center justify-between p-4 border-b border-gray-200">
          {#if title}
            <h3 class="text-lg font-semibold text-gray-900">{title}</h3>
          {/if}
          {#if showClose}
            <button
              type="button"
              class="text-gray-400 hover:text-gray-600 transition-colors cursor-pointer"
              on:click={closeModal}
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
        <slot />
      </div>
    </div>
  </div>
{/if}