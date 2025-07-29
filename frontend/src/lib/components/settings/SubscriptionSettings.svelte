<script lang="ts">
  import { onMount } from 'svelte';
  import { Button, Card } from '$lib/components/ui';
  import { Loader2, CreditCard, ExternalLink } from 'lucide-svelte';
  import {
    subscription,
    currentPlan,
    planLimits,
  } from '$lib/stores/subscription';
  import { APP_CONFIG } from '$lib/constants/config';
  import { PlanSelector } from '$lib/components/subscription';
  import type { ModelsSubscriptionPlan } from '$lib/api/api';
  import { LIMITS } from '$lib/subscription/feature-registry';

  interface Props {
    onError?: (message: string) => void;
  }

  let { onError }: Props = $props();

  let isLoading = $state(false);
  let isCreatingSession = $state(false);

  onMount(() => {
    loadData();
  });

  async function loadData() {
    isLoading = true;
    try {
      // Fetch subscription, plans, and usage data
      await Promise.all([
        subscription.fetchSubscription(),
        subscription.fetchPlans(),
        subscription.fetchUsage(),
      ]);
    } catch (error) {
      onError?.('Failed to load subscription data');
    } finally {
      isLoading = false;
    }
  }

  async function handleUpgrade(plan: ModelsSubscriptionPlan) {
    isCreatingSession = true;
    try {
      const session = await subscription.createCheckoutSession(plan.id);
      if (session.checkout_url) {
        window.location.href = session.checkout_url;
      }
    } catch (error) {
      onError?.('Failed to create checkout session');
    } finally {
      isCreatingSession = false;
    }
  }

  async function handleManageBilling() {
    isCreatingSession = true;
    try {
      const session = await subscription.createPortalSession();
      if (session.portal_url) {
        window.open(session.portal_url, '_blank');
      }
    } catch (error) {
      onError?.('Failed to open billing portal');
    } finally {
      isCreatingSession = false;
    }
  }

  // Format price display
  function formatPrice(price: number) {
    return new Intl.NumberFormat('en-US', {
      style: 'currency',
      currency: 'USD',
      minimumFractionDigits: 0,
    }).format(price);
  }

  // Format date
  function formatDate(dateString: string) {
    return new Date(dateString).toLocaleDateString('en-US', {
      year: 'numeric',
      month: 'long',
      day: 'numeric',
    });
  }

  let subscriptionData = $derived($subscription);
  let current = $derived($currentPlan);
  let limits = $derived($planLimits);
  let plans = $derived(subscriptionData.plans || []);
  let usage = $derived(subscriptionData.usage);

  // Use subscription data from API instead of JWT for plan details
  let currentSubscription = $derived(subscriptionData.subscription);
  let actualCurrentPlan = $derived(
    currentSubscription && plans.length > 0
      ? plans.find(p => p.id === currentSubscription.plan_id)
      : null
  );

  // Check if current plan is custom (not in visible plans list)
  let isCustomPlan = $derived(
    currentSubscription &&
      plans.length > 0 &&
      !plans.find(p => p.id === currentSubscription.plan_id)
  );
</script>

<div>
  <div class="mb-8">
    <h2 class="text-2xl font-bold text-gray-900">Subscription Plan</h2>
    <p class="mt-1 text-sm text-gray-600">
      View your current plan and upgrade options
    </p>
  </div>

  {#if isLoading}
    <div class="flex items-center justify-center py-12">
      <Loader2 class="h-8 w-8 animate-spin text-gray-400" />
    </div>
  {:else if subscriptionData.subscription}
    <!-- Custom Plan Notice -->
    {#if isCustomPlan}
      <div class="mb-6 bg-amber-50 border border-amber-200 rounded-lg p-4">
        <div class="flex">
          <div class="flex-shrink-0">
            <svg
              class="h-5 w-5 text-amber-400"
              viewBox="0 0 20 20"
              fill="currentColor">
              <path
                fill-rule="evenodd"
                d="M10 18a8 8 0 100-16 8 8 0 000 16zm0-2a6 6 0 100-12 6 6 0 000 12zm0-10a1 1 0 011 1v4a1 1 0 11-2 0V7a1 1 0 011-1zm0 8a1 1 0 100-2 1 1 0 000 2z"
                clip-rule="evenodd" />
            </svg>
          </div>
          <div class="ml-3">
            <h3 class="text-sm font-medium text-amber-800">Custom Plan</h3>
            <p class="mt-1 text-sm text-amber-700">
              You are currently on a custom plan tailored specifically for your
              organization. To make any changes to your subscription, please
              contact our support team at
              <a
                href={`mailto:${APP_CONFIG.emails.support}`}
                class="font-medium underline">{APP_CONFIG.emails.support}</a
              >.
            </p>
          </div>
        </div>
      </div>
    {/if}

    <!-- Current Plan -->
    <div class="mb-8">
      <div
        class="bg-gradient-to-r from-blue-500 to-purple-600 rounded-xl p-6 text-white">
        <div class="flex items-start justify-between">
          <div>
            <h3 class="text-sm font-medium text-blue-100">Current Plan</h3>
            <div class="flex items-center gap-2 mt-1">
              <p class="text-2xl font-bold">
                {actualCurrentPlan?.name ||
                  currentSubscription?.plan_name ||
                  'Unknown'}
              </p>
              {#if isCustomPlan}
                <span
                  class="bg-amber-400 text-amber-900 text-xs font-semibold px-2 py-1 rounded-full">
                  CUSTOM
                </span>
              {/if}
            </div>
            <p class="mt-2 text-sm text-blue-100">
              {formatPrice(
                actualCurrentPlan?.price || 0
              )}/{actualCurrentPlan?.interval || 'month'}
            </p>
          </div>
          <div class="text-right">
            <p class="text-sm text-blue-100">Current period ends</p>
            <p class="text-lg font-semibold">
              {formatDate(subscriptionData.subscription.current_period_end)}
            </p>
          </div>
        </div>

        <div class="mt-6">
          <Button
            variant="glass"
            size="sm"
            onclick={handleManageBilling}
            disabled={isCreatingSession}
            class="text-white border-white/20 hover:bg-white/10">
            <CreditCard class="h-4 w-4 mr-2" />
            Manage Billing
            <ExternalLink class="h-3 w-3 ml-2" />
          </Button>
        </div>
      </div>
    </div>

    <!-- Available Plans -->
    {#if !isCustomPlan}
      <div class="space-y-4">
        <h3 class="text-lg font-medium text-gray-900">Available Plans</h3>

        <PlanSelector
          {plans}
          currentPlanId={currentSubscription?.plan_id}
          isLoading={isCreatingSession}
          onSelectPlan={handleUpgrade}
          actionLabel={plan => `Upgrade to ${plan.name}`} />
      </div>
    {/if}
  {:else}
    <!-- No Subscription -->
    <Card variant="glass" class="p-8 text-center">
      <h3 class="text-lg font-semibold text-gray-900 mb-2">
        No Active Subscription
      </h3>
      <p class="text-gray-600 mb-6">Choose a plan to get started with Kyooar</p>

      <div class="mt-8">
        <PlanSelector
          {plans}
          isLoading={isCreatingSession}
          onSelectPlan={handleUpgrade}
          actionLabel={() => 'Get Started'}
          showCurrentBadge={false} />
      </div>
    </Card>
  {/if}
</div>
