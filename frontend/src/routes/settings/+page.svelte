<script lang="ts">
  import { auth } from '$lib/stores/auth';
  import { Card } from '$lib/components/ui';
  import { goto } from '$app/navigation';

  // Import icons from lucide-svelte
  import {
    User,
    Users,
    DollarSign,
    FileText,
    CreditCard,
    Settings,
    CheckCircle,
    XCircle,
  } from 'lucide-svelte';

  // Import settings components
  import {
    GeneralSettings,
    AccountSettings,
    TeamSettings,
    SubscriptionSettings,
    BillingHistory,
    PaymentMethods,
  } from '$lib/components/settings';
  import { ROLES } from '$lib/components/auth';
  import type { subscription } from '$lib/stores/subscription';
  import RoleGate from '$lib/components/auth/RoleGate.svelte';
  import type { Role } from '$lib/utils/auth-guards';

  let authState = $derived($auth);
  let user = $derived(authState.user);

  // Tab state
  let activeTab = $state('general');

  // Message states
  let successMessage = $state('');
  let errorMessage = $state('');

  function clearMessages() {
    successMessage = '';
    errorMessage = '';
  }

  function showSuccess(message: string) {
    clearMessages();
    successMessage = message;
    setTimeout(clearMessages, 5000);
  }

  function showError(message: string) {
    clearMessages();
    errorMessage = message;
    setTimeout(clearMessages, 5000);
  }

  async function handleDeactivation() {
    showSuccess(
      'Account deactivation scheduled. Your account will be deactivated in 15 days. You can cancel this anytime by logging in or from your account settings.'
    );

    // Log the user out after scheduling deactivation
    setTimeout(async () => {
      await auth.logout();
      // Small delay to ensure auth state is propagated
      await new Promise(resolve => setTimeout(resolve, 100));
      // Redirect to login page
      await goto('/login', { replaceState: true });
    }, 3000); // Give user time to read the message
  }

  const allowedRoles: Record<string, Role[]> = {
    general: [ROLES.manager, ROLES.owner, ROLES.admin, ROLES.viewer],
    account: [ROLES.manager, ROLES.owner, ROLES.admin, ROLES.viewer],
    team: [ROLES.owner, ROLES.admin],
    subscription: [ROLES.admin, ROLES.owner],
    billing: [ROLES.admin, ROLES.owner],
    payment: [ROLES.admin, ROLES.owner],
  } as const;

  const tabs = [
    { id: 'general', label: 'General', icon: Settings },
    { id: 'account', label: 'Account', icon: User },
    { id: 'team', label: 'Team', icon: Users },
    { id: 'subscription', label: 'Subscription', icon: DollarSign },
    { id: 'billing', label: 'Billing History', icon: FileText },
    { id: 'payment', label: 'Payment', icon: CreditCard },
  ];

  const getAllowedRolesById = (id: string): Role[] => allowedRoles[id];
</script>

<svelte:head>
  <title>Settings - Kyooar</title>
  <meta
    name="description"
    content="Manage your Kyooar account settings and preferences" />
</svelte:head>

<div class="min-h-screen p-4 sm:p-6 lg:p-8">
  <div class="mx-auto max-w-6xl">
    <!-- Header -->
    <div class="mb-8">
      <div class="flex items-center space-x-4">
        <div
          class="h-12 w-12 bg-gradient-to-br from-blue-500 to-purple-600 rounded-xl flex items-center justify-center shadow-lg">
          <Settings class="h-6 w-6 text-white" />
        </div>
        <div>
          <h1 class="text-3xl font-bold text-gray-900">Settings</h1>
          <p class="mt-1 text-gray-600">
            Manage your account, subscription, and billing
          </p>
        </div>
      </div>
    </div>

    <!-- Messages -->
    {#if successMessage}
      <div class="mb-6 rounded-lg bg-green-50 border border-green-200 p-4">
        <div class="flex">
          <CheckCircle class="h-5 w-5 text-green-400" />
          <p class="ml-3 text-sm text-green-800">{successMessage}</p>
        </div>
      </div>
    {/if}

    {#if errorMessage}
      <div class="mb-6 rounded-lg bg-red-50 border border-red-200 p-4">
        <div class="flex">
          <XCircle class="h-5 w-5 text-red-400" />
          <p class="ml-3 text-sm text-red-800">{errorMessage}</p>
        </div>
      </div>
    {/if}

    <div class="grid grid-cols-1 gap-6 lg:grid-cols-4">
      <!-- Tabs Navigation -->
      <nav class="lg:col-span-1">
        <Card variant="glass" class="p-4">
          <h3
            class="text-sm font-semibold text-gray-500 uppercase tracking-wider mb-4">
            Menu
          </h3>
          <ul class="space-y-2">
            {#each tabs as tab}
              <RoleGate roles={getAllowedRolesById(tab.id)}>
                <li>
                  <button
                    type="button"
                    class="w-full flex items-center space-x-3 px-4 py-3 text-sm font-medium rounded-xl transition-all duration-200 cursor-pointer {activeTab ===
                    tab.id
                      ? 'bg-gradient-to-r from-blue-500 to-purple-600 text-white shadow-lg'
                      : 'text-gray-700 hover:bg-gray-100'}"
                    onclick={() => (activeTab = tab.id)}>
                    <svelte:component this={tab.icon} class="h-5 w-5" />
                    <span>{tab.label}</span>
                  </button>
                </li>
              </RoleGate>
            {/each}
          </ul>
        </Card>
      </nav>

      <!-- Tab Content -->
      <div class="lg:col-span-3">
        <Card variant="glass" class="p-6">
          {#if activeTab === 'general'}
            <GeneralSettings onSuccess={showSuccess} onError={showError} />
          {:else if activeTab === 'account'}
            <AccountSettings
              {user}
              onSuccess={showSuccess}
              onError={showError}
              onDeactivate={handleDeactivation} />
          {:else if activeTab === 'team'}
            <TeamSettings onSuccess={showSuccess} onError={showError} />
          {:else if activeTab === 'subscription'}
            <SubscriptionSettings onError={showError} />
          {:else if activeTab === 'billing'}
            <BillingHistory
              onDownload={invoiceId =>
                showSuccess(`Downloading invoice ${invoiceId}...`)} />
          {:else if activeTab === 'payment'}
            <PaymentMethods
              onEditPayment={id =>
                showSuccess(`Editing payment method ${id}...`)}
              onAddPayment={type =>
                showSuccess(`Adding ${type} payment method...`)}
              onUpdateAddress={() =>
                showSuccess('Updating billing address...')} />
          {/if}
        </Card>
      </div>
    </div>
  </div>
</div>
