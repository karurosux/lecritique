<script lang="ts">
  import { auth } from '$lib/stores/auth';
  import type { Role } from '$lib/utils/auth-guards';

  interface Props {
    roles?: Role[];
    showLoading?: boolean;
    fallback?: string;
    children?: any;
  }

  let {
    roles = [],
    showLoading = false,
    fallback = '',
    children,
  }: Props = $props();

  let currentUserRole = $derived($auth.user?.role || null);

  let hasAccess = $derived.by(() => {
    return currentUserRole && roles.includes(currentUserRole as any);
  });

  let isLoading = $derived(showLoading && !$auth.user && $auth.isAuthenticated);
</script>

{#if isLoading}
  <div class="flex items-center justify-center p-4">
    <div
      class="animate-spin rounded-full h-6 w-6 border-b-2 border-primary-600">
    </div>
  </div>
{:else if hasAccess}
  {@render children?.()}
{:else if fallback}
  <div class="text-gray-500 text-sm">
    {fallback}
  </div>
{/if}
