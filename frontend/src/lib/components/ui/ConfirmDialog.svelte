<script lang="ts">
  import { Button, Modal } from '$lib/components/ui';
  import { AlertTriangle } from 'lucide-svelte';

  let {
    isOpen = false,
    title = 'Confirm Action',
    message = 'Are you sure you want to proceed?',
    confirmText = 'Confirm',
    cancelText = 'Cancel',
    variant = 'danger',
    clickOrigin = null,
    onConfirm = () => {},
    onCancel = () => {},
  }: {
    isOpen?: boolean;
    title?: string;
    message?: string;
    confirmText?: string;
    cancelText?: string;
    variant?: 'danger' | 'warning' | 'info';
    clickOrigin?: { x: number; y: number } | null;
    onConfirm?: () => void;
    onCancel?: () => void;
  } = $props();

  function handleConfirm() {
    onConfirm();
  }

  function handleCancel() {
    onCancel();
  }
</script>

<Modal {isOpen} {title} {clickOrigin} size="sm" onclose={handleCancel}>
  <div class="space-y-4">
    <div class="flex items-start gap-3">
      <div
        class="h-10 w-10 rounded-full flex items-center justify-center flex-shrink-0
        {variant === 'danger'
          ? 'bg-red-100'
          : variant === 'warning'
            ? 'bg-yellow-100'
            : 'bg-blue-100'}">
        <AlertTriangle
          class="h-5 w-5 
          {variant === 'danger'
            ? 'text-red-600'
            : variant === 'warning'
              ? 'text-yellow-600'
              : 'text-blue-600'}" />
      </div>
      <div class="flex-1">
        <p class="text-gray-600 text-sm leading-relaxed">
          {message}
        </p>
      </div>
    </div>

    <div
      class="mt-6 pt-6 border-t border-gray-200 flex items-center justify-end space-x-3">
      <Button onclick={handleCancel} variant="outline">
        {cancelText}
      </Button>
      <Button
        onclick={handleConfirm}
        variant={variant === 'danger' ? 'destructive' : 'gradient'}>
        {confirmText}
      </Button>
    </div>
  </div>
</Modal>
