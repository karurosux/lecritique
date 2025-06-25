<script lang="ts">
  import { onMount } from 'svelte';
  import { Button, Card } from '$lib/components/ui';
  import { Check, Loader2, CreditCard, ExternalLink } from 'lucide-svelte';
  import { subscription, currentPlan, planLimits } from '$lib/stores/subscription';
  import type { ModelsSubscriptionPlan } from '$lib/api/api';
  import { getPlanFeatures, LIMITS } from '$lib/subscription/feature-registry';

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
        subscription.fetchPlans()
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
      minimumFractionDigits: 0
    }).format(price);
  }

  // Format date
  function formatDate(dateString: string) {
    return new Date(dateString).toLocaleDateString('en-US', {
      year: 'numeric',
      month: 'long',
      day: 'numeric'
    });
  }

  // Get feature display text using the dynamic registry
  function getFeatureText(plan: ModelsSubscriptionPlan): string[] {
    return getPlanFeatures(plan);
  }

  let subscriptionData = $derived($subscription);
  let current = $derived($currentPlan);
  let limits = $derived($planLimits);
  let plans = $derived(subscriptionData.plans || []);
</script>

<div>
  <div class="mb-8">
    <h2 class="text-2xl font-bold text-gray-900">Subscription Plan</h2>
    <p class="mt-1 text-sm text-gray-600">View your current plan and upgrade options</p>
  </div>
  
  {#if isLoading}
    <div class="flex items-center justify-center py-12">
      <Loader2 class="h-8 w-8 animate-spin text-gray-400" />
    </div>
  {:else if subscriptionData.subscription}
    <!-- Current Plan -->
    <div class="mb-8">
      <div class="bg-gradient-to-r from-blue-500 to-purple-600 rounded-xl p-6 text-white">
        <div class="flex items-start justify-between">
          <div>
            <h3 class="text-sm font-medium text-blue-100">Current Plan</h3>
            <p class="mt-1 text-2xl font-bold">{current?.name || 'Unknown'}</p>
            <p class="mt-2 text-sm text-blue-100">
              {formatPrice(current?.price || 0)}/{current?.interval || 'month'}
            </p>
          </div>
          <div class="text-right">
            <p class="text-sm text-blue-100">Current period ends</p>
            <p class="text-lg font-semibold">
              {formatDate(subscriptionData.subscription.current_period_end)}
            </p>
          </div>
        </div>
        
        {#if limits}
          <div class="mt-6 grid grid-cols-2 gap-4 border-t border-white/20 pt-4">
            <div>
              <p class="text-sm text-blue-100">Restaurants</p>
              <p class="text-lg font-semibold">
                {limits.limits?.[LIMITS.RESTAURANTS] === -1 ? 'Unlimited' : `0 of ${limits.limits?.[LIMITS.RESTAURANTS] || 0} used`}
              </p>
            </div>
            <div>
              <p class="text-sm text-blue-100">Monthly Feedback</p>
              <p class="text-lg font-semibold">
                {limits.limits?.[LIMITS.FEEDBACKS_PER_MONTH] === -1 ? 'Unlimited' : `0 of ${(limits.limits?.[LIMITS.FEEDBACKS_PER_MONTH] || 0).toLocaleString()}`}
              </p>
            </div>
          </div>
        {/if}

        <div class="mt-6">
          <Button 
            variant="glass"
            size="sm"
            onclick={handleManageBilling}
            disabled={isCreatingSession}
            class="text-white border-white/20 hover:bg-white/10"
          >
            <CreditCard class="h-4 w-4 mr-2" />
            Manage Billing
            <ExternalLink class="h-3 w-3 ml-2" />
          </Button>
        </div>
      </div>
    </div>

    <!-- Available Plans -->
    <div class="space-y-4">
      <h3 class="text-lg font-medium text-gray-900">Available Plans</h3>
      
      <div class="grid gap-4 md:grid-cols-3">
        {#each plans as plan}
          {@const isCurrent = plan.id === current?.id}
          <div class="relative">
            {#if isCurrent}
              <div class="absolute -top-3 left-1/2 -translate-x-1/2 z-10">
                <span class="bg-blue-500 text-white text-xs font-semibold px-3 py-1 rounded-full whitespace-nowrap">
                  Current Plan
                </span>
              </div>
            {/if}
            <Card 
              variant={isCurrent ? 'glass' : 'default'} 
              class="relative p-6 overflow-visible {isCurrent ? 'ring-2 ring-blue-500 ring-offset-2' : ''}"
            >
            
            <h4 class="text-lg font-semibold text-gray-900">{plan.name}</h4>
            {#if plan.description}
              <p class="mt-1 text-sm text-gray-600">{plan.description}</p>
            {/if}
            <p class="mt-3 text-3xl font-bold text-gray-900">
              {formatPrice(plan.price)}
              <span class="text-lg font-normal text-gray-600">/{plan.interval}</span>
            </p>
            
            <ul class="mt-6 space-y-3 text-sm text-gray-600">
              {#each getFeatureText(plan) as feature}
                <li class="flex items-start">
                  <Check class="h-5 w-5 text-green-500 mr-2 flex-shrink-0" />
                  {feature}
                </li>
              {/each}
            </ul>
            
            {#if isCurrent}
              <Button variant="primary" size="md" class="w-full mt-6" disabled>
                Current Plan
              </Button>
            {:else}
              <Button 
                variant="gradient"
                size="md" 
                class="w-full mt-6"
                onclick={() => handleUpgrade(plan)}
                disabled={isCreatingSession}
              >
                {#if isCreatingSession}
                  <Loader2 class="h-4 w-4 mr-2 animate-spin" />
                  Processing...
                {:else}
                  Upgrade to {plan.name}
                {/if}
              </Button>
            {/if}
            </Card>
          </div>
        {/each}
      </div>
    </div>
  {:else}
    <!-- No Subscription -->
    <Card variant="glass" class="p-8 text-center">
      <h3 class="text-lg font-semibold text-gray-900 mb-2">No Active Subscription</h3>
      <p class="text-gray-600 mb-6">Choose a plan to get started with LeCritique</p>
      
      <div class="grid gap-4 md:grid-cols-3 mt-8">
        {#each plans as plan}
          <Card variant="default" class="p-6">
            <h4 class="text-lg font-semibold text-gray-900">{plan.name}</h4>
            {#if plan.description}
              <p class="mt-1 text-sm text-gray-600">{plan.description}</p>
            {/if}
            <p class="mt-3 text-3xl font-bold text-gray-900">
              {formatPrice(plan.price)}
              <span class="text-lg font-normal text-gray-600">/{plan.interval}</span>
            </p>
            
            <ul class="mt-6 space-y-3 text-sm text-gray-600">
              {#each getFeatureText(plan) as feature}
                <li class="flex items-start">
                  <Check class="h-5 w-5 text-green-500 mr-2 flex-shrink-0" />
                  {feature}
                </li>
              {/each}
            </ul>
            
            <Button 
              variant="gradient"
              size="md" 
              class="w-full mt-6"
              onclick={() => handleUpgrade(plan)}
              disabled={isCreatingSession}
            >
              {#if isCreatingSession}
                <Loader2 class="h-4 w-4 mr-2 animate-spin" />
                Processing...
              {:else}
                Get Started
              {/if}
            </Button>
          </Card>
        {/each}
      </div>
    </Card>
  {/if}
</div>