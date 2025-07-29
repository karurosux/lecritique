<script lang="ts">
  import { onMount } from 'svelte';

  let {
    isOpen = $bindable(false),
    open = $bindable(false), // Backward compatibility
    title = '',
    showClose = true,
    size = 'md',
    onclose = () => {},
    clickOrigin = null,
    children,
  }: {
    isOpen?: boolean;
    open?: boolean;
    title?: string;
    showClose?: boolean;
    size?: 'sm' | 'md' | 'lg' | 'xl';
    onclose?: () => void;
    clickOrigin?: { x: number; y: number } | null;
    children?: any;
  } = $props();

  // Support both isOpen and open props for compatibility
  let modalOpen = $derived(isOpen || open);

  const sizeClasses = {
    sm: 'max-w-sm',
    md: 'max-w-md',
    lg: 'max-w-2xl',
    xl: 'max-w-4xl',
  };

  let modalElement = $state<HTMLDivElement | null>(null);
  let animationPhase = $state<'expanding' | 'fading' | 'showing' | 'idle'>(
    'idle'
  );

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

  // Move modal to document body when it opens and handle animation
  $effect(() => {
    if (modalOpen && modalElement) {
      document.body.appendChild(modalElement);

      // Start all animations simultaneously
      animationPhase = 'expanding';

      // Show modal content after a short delay
      setTimeout(() => {
        animationPhase = 'showing';
      }, 200);

      return () => {
        if (modalElement && modalElement.parentNode) {
          modalElement.parentNode.removeChild(modalElement);
        }
        animationPhase = 'idle';
      };
    }
  });
</script>

{#if modalOpen}
  <div
    bind:this={modalElement}
    class="fixed inset-0 z-[10000] flex items-center justify-center transition-all duration-200"
    class:circle-expanding={animationPhase === 'expanding'}
    class:backdrop-fading={animationPhase === 'fading'}
    class:backdrop-showing={animationPhase === 'showing'}
    style={clickOrigin
      ? `--click-x: ${clickOrigin.x}px; --click-y: ${clickOrigin.y}px;`
      : ''}
    onclick={handleBackdropClick}
    onkeydown={handleKeydown}
    role="dialog"
    aria-modal="true"
    tabindex="-1">
    <div
      class="relative bg-white rounded-lg shadow-xl {sizeClasses[
        size
      ]} w-full mx-4 my-8 max-h-[calc(100vh-4rem)] overflow-y-auto"
      class:animate-modal-enter={animationPhase === 'showing'}
      class:opacity-0={animationPhase !== 'showing'}>
      {#if title || showClose}
        <div
          class="flex items-center justify-between p-4 border-b border-gray-200">
          {#if title}
            <h3 class="text-lg font-semibold text-gray-900">{title}</h3>
          {/if}
          {#if showClose}
            <button
              type="button"
              class="text-gray-400 hover:text-gray-600 transition-colors cursor-pointer"
              onclick={closeModal}
              aria-label="Close modal">
              <svg
                class="w-6 h-6"
                fill="none"
                stroke="currentColor"
                viewBox="0 0 24 24">
                <path
                  stroke-linecap="round"
                  stroke-linejoin="round"
                  stroke-width="2"
                  d="M6 18L18 6M6 6l12 12" />
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

<style>
  :global(body:has(.fixed[role='dialog'])) {
    overflow: hidden;
  }

  @keyframes modal-enter {
    from {
      opacity: 0;
      transform: scale(0.95) translateY(10px);
    }
    to {
      opacity: 1;
      transform: scale(1) translateY(0);
    }
  }

  @keyframes circle-expand {
    from {
      clip-path: circle(0% at var(--click-x, 50%) var(--click-y, 50%));
      background: linear-gradient(
        135deg,
        rgb(243 244 246 / 0.3),
        rgb(209 213 219 / 0.3)
      );
      backdrop-filter: blur(4px);
      opacity: 0;
    }
    to {
      clip-path: circle(150% at var(--click-x, 50%) var(--click-y, 50%));
      background: linear-gradient(
        135deg,
        rgb(243 244 246 / 0.6),
        rgb(209 213 219 / 0.6)
      );
      backdrop-filter: blur(16px);
      opacity: 1;
    }
  }

  @keyframes backdrop-fade {
    from {
      background: linear-gradient(
        135deg,
        rgb(243 244 246 / 0.8),
        rgb(209 213 219 / 0.8)
      );
      backdrop-filter: blur(4px);
    }
    to {
      background: linear-gradient(
        135deg,
        rgb(243 244 246 / 0.6),
        rgb(209 213 219 / 0.6)
      );
      backdrop-filter: blur(16px);
    }
  }

  .circle-expanding {
    animation: circle-expand 0.4s ease-out forwards;
  }

  .backdrop-fading {
    animation: backdrop-fade 0.2s ease-out forwards;
  }

  .backdrop-showing {
    background: linear-gradient(
      135deg,
      rgb(243 244 246 / 0.6),
      rgb(209 213 219 / 0.6)
    );
    backdrop-filter: blur(16px);
  }

  :global(.animate-modal-enter) {
    animation: modal-enter 0.2s ease-out;
  }
</style>
