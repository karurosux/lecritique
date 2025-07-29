<script lang="ts">
  import { Button, Input, Modal } from '$lib/components/ui';
  import {
    AlertCircle,
    Mail,
    AlertTriangle,
    Edit,
    Calendar,
    Info,
  } from 'lucide-svelte';
  import type { User } from '$lib/stores/auth';
  import { getApiClient, handleApiError } from '$lib/api/client';

  interface Props {
    user: User | null;
    onSuccess?: (message: string) => void;
    onError?: (message: string) => void;
    onDeactivate?: () => void;
  }

  let { user, onSuccess, onError, onDeactivate }: Props = $props();

  let isLoading = $state(false);
  let showDeactivateConfirm = $state(false);
  let showEmailModal = $state(false);
  let deactivateConfirmText = $state('');
  let deactivationDate = $state<Date | null>(null);
  let emailForm = $state({
    newEmail: '',
    confirmEmail: '',
  });

  // Check if user has pending deactivation
  $effect(() => {
    if (user?.deactivation_requested_at) {
      const requestedAt = new Date(user.deactivation_requested_at);
      deactivationDate = new Date(
        requestedAt.getTime() + 15 * 24 * 60 * 60 * 1000
      ); // Add 15 days
    }
  });

  function openEmailModal() {
    emailForm = {
      newEmail: '',
      confirmEmail: '',
    };
    showEmailModal = true;
  }

  function closeEmailModal() {
    showEmailModal = false;
    emailForm = {
      newEmail: '',
      confirmEmail: '',
    };
  }

  async function handleEmailChange(event: Event) {
    event.preventDefault();

    // Validate emails match
    if (emailForm.newEmail !== emailForm.confirmEmail) {
      onError?.('Email addresses do not match');
      return;
    }

    // Validate email format
    const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;
    if (!emailRegex.test(emailForm.newEmail)) {
      onError?.('Please enter a valid email address');
      return;
    }

    // Check if email is different from current
    if (emailForm.newEmail === user?.email) {
      onError?.('New email must be different from current email');
      return;
    }

    isLoading = true;

    try {
      const api = getApiClient();
      const response = await api.api.v1AuthChangeEmailCreate({
        new_email: emailForm.newEmail,
      });

      // Check if we got a new token (dev mode)
      if (response.data.data?.token) {
        // Update auth state with new token
        const { auth } = await import('$lib/stores/auth');
        auth.updateToken(response.data.data.token);
        onSuccess?.('Email changed successfully!');
      } else {
        onSuccess?.(
          'Email change request sent. Please check your new email for verification.'
        );
      }

      closeEmailModal();
    } catch (error) {
      const errorMessage = handleApiError(error);
      onError?.(errorMessage);
    } finally {
      isLoading = false;
    }
  }

  function handleDeactivateClick() {
    showDeactivateConfirm = true;
  }

  async function confirmDeactivate() {
    isLoading = true;

    try {
      // TODO: Uncomment when API endpoints are available
      // const api = getApiClient();
      // const response = await api.api.v1AuthDeactivateCreate();

      // Temporary: Show success message and trigger logout
      const futureDate = new Date();
      futureDate.setDate(futureDate.getDate() + 15);
      deactivationDate = futureDate;

      onSuccess?.(
        `Your account will be deactivated on ${deactivationDate.toLocaleDateString()}. You can cancel this request by logging in before this date.`
      );

      // Show the deactivation info instead of logging out immediately
      showDeactivateConfirm = false;
      deactivateConfirmText = '';

      // Call the onDeactivate callback to handle logout
      onDeactivate?.();
    } catch (error) {
      const errorMessage = handleApiError(error);
      onError?.(errorMessage);
    } finally {
      isLoading = false;
    }
  }

  async function cancelDeactivation() {
    isLoading = true;

    try {
      // TODO: Uncomment when API endpoints are available
      // const api = getApiClient();
      // await api.api.v1AuthCancelDeactivationCreate();

      // Temporary: Just clear the deactivation date
      deactivationDate = null;
      onSuccess?.('Account deactivation cancelled successfully.');
    } catch (error) {
      const errorMessage = handleApiError(error);
      onError?.(errorMessage);
    } finally {
      isLoading = false;
    }
  }

  function cancelDeactivateDialog() {
    showDeactivateConfirm = false;
    deactivateConfirmText = '';
  }

  let isDeactivateEnabled = $derived(
    deactivateConfirmText.toLowerCase() === 'deactivate'
  );
</script>

<div>
  <div class="mb-8">
    <h2 class="text-2xl font-bold text-gray-900">Account Settings</h2>
    <p class="mt-1 text-sm text-gray-600">
      Manage your email address and account status
    </p>
  </div>

  <div class="space-y-8">
    <!-- Email Display Section -->
    <div>
      <div class="flex items-center gap-2 mb-4">
        <Mail class="h-5 w-5 text-gray-400" />
        <h3 class="text-lg font-medium text-gray-900">Email Address</h3>
      </div>

      <div class="bg-gray-50 rounded-lg p-4">
        <div class="flex items-center justify-between">
          <div class="flex-1">
            <p class="text-sm text-gray-600 mb-1">Current email</p>
            <p class="font-medium text-gray-900">
              {user?.email || 'No email set'}
            </p>
            {#if user && !user.email_verified}
              <div class="mt-2 flex items-center gap-2">
                <AlertCircle class="h-4 w-4 text-yellow-500" />
                <p class="text-sm text-yellow-700">Email not verified</p>
              </div>
            {/if}
          </div>
          <Button
            variant="outline"
            size="sm"
            onclick={openEmailModal}
            class="flex items-center gap-2">
            <Edit class="h-4 w-4" />
            Change Email
          </Button>
        </div>
      </div>
    </div>

    <!-- Account Deactivation Section -->
    <div class="border-t pt-8">
      <div class="flex items-center gap-2 mb-4">
        <AlertTriangle class="h-5 w-5 text-red-500" />
        <h3 class="text-lg font-medium text-gray-900">Danger Zone</h3>
      </div>

      <div class="bg-red-50 border-2 border-red-200 rounded-lg p-6">
        <h4 class="text-base font-semibold text-red-900 mb-2">
          Deactivate Account
        </h4>

        {#if deactivationDate}
          <!-- Show pending deactivation status -->
          <div
            class="bg-yellow-50 border border-yellow-200 rounded-lg p-4 mb-4">
            <div class="flex items-start gap-3">
              <Info class="h-5 w-5 text-yellow-600 mt-0.5" />
              <div class="flex-1">
                <p class="text-sm font-medium text-yellow-900 mb-1">
                  Account deactivation scheduled
                </p>
                <p class="text-sm text-yellow-700 mb-2">
                  Your account will be permanently deactivated on {deactivationDate.toLocaleDateString()}.
                </p>
                <p class="text-sm text-yellow-700">
                  To keep your account active, simply log in before this date or
                  click cancel below.
                </p>
              </div>
            </div>
          </div>

          <Button
            variant="outline"
            class="border-yellow-300 text-yellow-700 hover:bg-yellow-50"
            onclick={cancelDeactivation}
            disabled={isLoading}>
            {#if isLoading}
              <svg
                class="animate-spin -ml-1 mr-3 h-5 w-5"
                fill="none"
                viewBox="0 0 24 24">
                <circle
                  class="opacity-25"
                  cx="12"
                  cy="12"
                  r="10"
                  stroke="currentColor"
                  stroke-width="4"></circle>
                <path
                  class="opacity-75"
                  fill="currentColor"
                  d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"
                ></path>
              </svg>
              Cancelling...
            {:else}
              Cancel Deactivation
            {/if}
          </Button>
        {:else}
          <p class="text-sm text-red-700 mb-4">
            Request account deactivation with a 15-day grace period. You can
            cancel anytime by logging in.
          </p>

          {#if !showDeactivateConfirm}
            <Button
              variant="outline"
              class="border-red-300 text-red-600 hover:bg-red-50"
              onclick={handleDeactivateClick}>
              Deactivate Account
            </Button>
          {:else}
            <div class="space-y-4">
              <div class="bg-white border border-red-300 rounded-lg p-4">
                <p class="text-sm font-medium text-red-900 mb-2">
                  Are you absolutely sure?
                </p>
                <p class="text-sm text-red-700 mb-4">
                  This will schedule your account for deactivation after a
                  15-day grace period. You can cancel this request anytime by
                  logging in during the grace period.
                </p>

                <div class="space-y-2">
                  <label
                    for="deactivate-confirm"
                    class="block text-sm font-medium text-red-900">
                    Type <span class="font-mono bg-red-100 px-1 rounded"
                      >deactivate</span> to confirm:
                  </label>
                  <Input
                    id="deactivate-confirm"
                    type="text"
                    bind:value={deactivateConfirmText}
                    placeholder="Type 'deactivate' to confirm"
                    class="border-red-300 focus:border-red-500 focus:ring-red-500" />
                </div>
              </div>

              <div class="flex gap-3">
                <Button
                  variant="primary"
                  class="bg-red-600 hover:bg-red-700 disabled:bg-red-300"
                  onclick={confirmDeactivate}
                  disabled={isLoading || !isDeactivateEnabled}>
                  {#if isLoading}
                    <svg
                      class="animate-spin -ml-1 mr-3 h-5 w-5 text-white"
                      fill="none"
                      viewBox="0 0 24 24">
                      <circle
                        class="opacity-25"
                        cx="12"
                        cy="12"
                        r="10"
                        stroke="currentColor"
                        stroke-width="4"></circle>
                      <path
                        class="opacity-75"
                        fill="currentColor"
                        d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"
                      ></path>
                    </svg>
                    Deactivating...
                  {:else}
                    Deactivate My Account
                  {/if}
                </Button>
                <Button
                  variant="outline"
                  onclick={cancelDeactivateDialog}
                  disabled={isLoading}>
                  Cancel
                </Button>
              </div>
            </div>
          {/if}
        {/if}
      </div>
    </div>
  </div>
</div>

<!-- Email Change Modal -->
<Modal bind:open={showEmailModal} onClose={closeEmailModal}>
  <div class="w-full max-w-md">
    <div class="mb-6">
      <h3 class="text-xl font-semibold text-gray-900">Change Email Address</h3>
      <p class="mt-1 text-sm text-gray-600">
        Enter your new email address below
      </p>
    </div>

    <form onsubmit={handleEmailChange} class="space-y-4">
      <Input
        label="New Email Address"
        type="email"
        bind:value={emailForm.newEmail}
        placeholder="Enter new email address"
        required />

      <Input
        label="Confirm New Email"
        type="email"
        bind:value={emailForm.confirmEmail}
        placeholder="Confirm new email address"
        required />

      <div class="bg-blue-50 border border-blue-200 rounded-lg p-4">
        <p class="text-sm text-blue-800">
          After changing your email, you'll need to verify your new email
          address. A verification link will be sent to both your old and new
          email addresses.
        </p>
      </div>

      <div class="flex gap-3 justify-end pt-4">
        <Button
          type="button"
          variant="outline"
          onclick={closeEmailModal}
          disabled={isLoading}>
          Cancel
        </Button>
        <Button
          type="submit"
          variant="gradient"
          disabled={isLoading ||
            !emailForm.newEmail ||
            !emailForm.confirmEmail}>
          {#if isLoading}
            <svg
              class="animate-spin -ml-1 mr-3 h-5 w-5 text-white"
              fill="none"
              viewBox="0 0 24 24">
              <circle
                class="opacity-25"
                cx="12"
                cy="12"
                r="10"
                stroke="currentColor"
                stroke-width="4"></circle>
              <path
                class="opacity-75"
                fill="currentColor"
                d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"
              ></path>
            </svg>
            Updating...
          {:else}
            Update Email
          {/if}
        </Button>
      </div>
    </form>
  </div>
</Modal>
