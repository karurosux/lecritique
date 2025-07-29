<script lang="ts">
  import { onMount } from 'svelte';
  import { Button, Card } from '$lib/components/ui';
  import { Check, Loader2, CreditCard, ExternalLink } from 'lucide-svelte';
  import {
    subscription,
    currentPlan,
    planLimits,
  } from '$lib/stores/subscription';
  import type { ModelsSubscriptionPlan } from '$lib/api/api';
  import { _ } from 'svelte-i18n';

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
      await Promise.all([
        subscription.fetchSubscription(),
        subscription.fetchPlans(),
      ]);
    } catch (error) {
      onError?.($_('subscription.errors.loadFailed'));
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
      onError?.($_('subscription.errors.checkoutFailed'));
    } finally {
      isCreatingSession = false;
    }
  }

  // Get translated feature text
  function getFeatureText(plan: ModelsSubscriptionPlan): string[] {
    const features = plan.features;
    if (!features) return [];

    const items = [];

    // Use translation with pluralization
    items.push(
      $_('subscription.features.organizations', {
        values: { count: features.max_organizations },
      })
    );
    items.push(
      $_('subscription.features.feedbacks', {
        values: { count: features.max_feedbacks_per_month },
      })
    );
    items.push(
      $_('subscription.features.qrCodes', {
        values: { count: features.max_qr_codes_per_location },
      })
    );
    items.push(
      $_('subscription.features.teamMembers', {
        values: { count: features.max_team_members },
      })
    );

    // Additional features
    if (features.advanced_analytics) {
      items.push($_('subscription.features.advancedAnalytics'));
    }
    if (features.custom_branding) {
      items.push($_('subscription.features.customBranding'));
    }
    if (features.api_access) {
      items.push($_('subscription.features.apiAccess'));
    }
    if (features.priority_support) {
      items.push($_('subscription.features.prioritySupport'));
    }

    return items;
  }

  // Get translated plan name and description
  function getPlanName(plan: ModelsSubscriptionPlan): string {
    // Try to use translation key if it exists
    const translationKey = `subscription.plans.${plan.code}.name`;
    const translated = $_(translationKey);

    // If no translation found (key is returned), use database value
    return translated === translationKey ? plan.name : translated;
  }

  function getPlanDescription(plan: ModelsSubscriptionPlan): string {
    const translationKey = `subscription.plans.${plan.code}.description`;
    const translated = $_(translationKey);

    return translated === translationKey ? plan.description : translated;
  }

  let subscriptionData = $derived($subscription);
  let current = $derived($currentPlan);
  let limits = $derived($planLimits);
  let plans = $derived(subscriptionData.plans || []);
</script>

<div>
  <div class="mb-8">
    <h2 class="text-2xl font-bold text-gray-900">{$_('subscription.title')}</h2>
    <p class="mt-1 text-sm text-gray-600">{$_('subscription.subtitle')}</p>
  </div>

  {#if isLoading}
    <div class="flex items-center justify-center py-12">
      <Loader2 class="h-8 w-8 animate-spin text-gray-400" />
    </div>
  {:else if subscriptionData.subscription}
    <!-- Current Plan -->
    <div class="mb-8">
      <div
        class="bg-gradient-to-r from-blue-500 to-purple-600 rounded-xl p-6 text-white">
        <div class="flex items-start justify-between">
          <div>
            <h3 class="text-sm font-medium text-blue-100">
              {$_('subscription.currentPlan')}
            </h3>
            <p class="mt-1 text-2xl font-bold">
              {getPlanName(current) || $_('subscription.unknown')}
            </p>
            <!-- ... rest of the component ... -->
          </div>
        </div>
      </div>
    </div>
  {/if}
</div>
