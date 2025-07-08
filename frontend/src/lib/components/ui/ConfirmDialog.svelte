<script lang="ts">
  import { Button, Card } from '$lib/components/ui';
  import { AlertTriangle, X } from 'lucide-svelte';

  let {
    isOpen = false,
    title = 'Confirm Action',
    message = 'Are you sure you want to proceed?',
    confirmText = 'Confirm',
    cancelText = 'Cancel',
    variant = 'danger', // 'danger' | 'warning' | 'info'
    onConfirm = () => {},
    onCancel = () => {}
  }: {
    isOpen?: boolean;
    title?: string;
    message?: string;
    confirmText?: string;
    cancelText?: string;
    variant?: 'danger' | 'warning' | 'info';
    onConfirm?: () => void;
    onCancel?: () => void;
  } = $props();

  function handleConfirm() {
    onConfirm();
    isOpen = false;
  }

  function handleCancel() {
    onCancel();
    isOpen = false;
  }

  function handleBackdropClick(e: MouseEvent) {
    if (e.target === e.currentTarget) {
      handleCancel();
    }
  }

  function handleKeydown(e: KeyboardEvent) {
    if (e.key === 'Escape') {
      handleCancel();
    }
  }

  $effect(() => {
    if (isOpen) {
      document.body.style.overflow = 'hidden';
    } else {
      document.body.style.overflow = '';
    }

    return () => {
      document.body.style.overflow = '';
    };
  });
</script>

<svelte:window on:keydown={handleKeydown} />

{#if isOpen}
  <div 
    class="fixed inset-0 z-50 flex items-center justify-center bg-black bg-opacity-50 backdrop-blur-sm"
    onclick={handleBackdropClick}
  >
    <Card class="w-full max-w-md mx-4 relative" padding={false}>
      <div class="p-2">
        <!-- Header -->
        <div class="flex items-start justify-between mb-2">
          <div class="flex items-center gap-2">
            <div class="h-8 w-8 rounded-full flex items-center justify-center
              {variant === 'danger' ? 'bg-red-100' : 
               variant === 'warning' ? 'bg-yellow-100' : 
               'bg-blue-100'}">
              <AlertTriangle class="h-4 w-4 
                {variant === 'danger' ? 'text-red-600' : 
                 variant === 'warning' ? 'text-yellow-600' : 
                 'text-blue-600'}" />
            </div>
            <h3 class="text-lg font-semibold text-gray-900">{title}</h3>
          </div>
          <Button 
            onclick={handleCancel}
            variant="ghost" 
            size="sm"
            class="p-1"
          >
            <X class="h-4 w-4" />
          </Button>
        </div>

        <!-- Message -->
        <p class="text-gray-600 mb-3 ml-10 text-sm">
          {message}
        </p>

        <!-- Actions -->
        <div class="flex items-center justify-end gap-3">
          <Button 
            onclick={handleCancel}
            variant="outline"
          >
            {cancelText}
          </Button>
          <Button 
            onclick={handleConfirm}
            variant={variant === 'danger' ? 'destructive' : 'gradient'}
          >
            {confirmText}
          </Button>
        </div>
      </div>
    </Card>
  </div>
{/if}