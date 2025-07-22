<script lang="ts">
  import { Card } from "$lib/components/ui";
  import { Check, Loader2 } from "lucide-svelte";
  import type { ModelsSubscriptionPlan } from "$lib/api/api";

  interface Props {
    plans: ModelsSubscriptionPlan[];
    currentPlanId?: string;
    isLoading?: boolean;
    onSelectPlan: (plan: ModelsSubscriptionPlan) => void | Promise<void>;
    actionLabel?: (plan: ModelsSubscriptionPlan) => string;
    showCurrentBadge?: boolean;
    columns?: 1 | 2 | 3;
  }

  let {
    plans,
    currentPlanId = undefined,
    isLoading = false,
    onSelectPlan,
    actionLabel = (plan) => "Select Plan",
    showCurrentBadge = true,
    columns = 3,
  }: Props = $props();

  // Format price display
  function formatPrice(price: number) {
    return new Intl.NumberFormat("en-US", {
      style: "currency",
      currency: "USD",
      minimumFractionDigits: 0,
    }).format(price);
  }

  // Get plan features
  function getPlanFeatures(plan: ModelsSubscriptionPlan): string[] {
    const features = [];

    // Add limits
    if (plan.max_organizations === -1) {
      features.push("Unlimited organizations");
    } else {
      features.push(
        `Up to ${plan.max_organizations} organization${plan.max_organizations > 1 ? "s" : ""}`,
      );
    }

    if (plan.max_qr_codes === -1) {
      features.push("Unlimited QR codes");
    } else {
      features.push(`Up to ${plan.max_qr_codes} QR codes`);
    }

    if (plan.max_feedbacks_per_month === -1) {
      features.push("Unlimited monthly feedbacks");
    } else {
      features.push(`${plan.max_feedbacks_per_month} feedbacks per month`);
    }

    if (plan.max_team_members === -1) {
      features.push("Unlimited team members");
    } else {
      features.push(`Up to ${plan.max_team_members} team members`);
    }

    // Add feature flags
    if (plan.has_basic_analytics) features.push("Basic analytics dashboard");
    if (plan.has_advanced_analytics)
      features.push("Advanced analytics & insights");
    if (plan.has_feedback_explorer) features.push("Feedback explorer");
    if (plan.has_custom_branding) features.push("Custom branding");
    if (plan.has_priority_support) features.push("Priority support");

    return features;
  }

  const gridClasses = {
    1: "grid-cols-1",
    2: "grid-cols-1 md:grid-cols-2",
    3: "grid-cols-1 md:grid-cols-3",
  };
</script>

<div class="grid gap-4 {gridClasses[columns]}">
  {#each plans as plan}
    {@const isCurrent = plan.id === currentPlanId}
    <div class="relative">
      {#if isCurrent && showCurrentBadge}
        <div class="absolute -top-3 left-1/2 -translate-x-1/2 z-10">
          <span
            class="bg-blue-500 text-white text-xs font-semibold px-3 py-1 rounded-full whitespace-nowrap"
          >
            Current Plan
          </span>
        </div>
      {/if}

      <Card
        variant={isCurrent ? "glass" : "default"}
        class="relative p-6 overflow-visible h-full flex flex-col {isCurrent
          ? 'ring-2 ring-blue-500 ring-offset-2'
          : ''}"
      >
        <div class="flex-1">
          <h4 class="text-lg font-semibold text-gray-900">{plan.name}</h4>
          {#if plan.description}
            <p class="mt-1 text-sm text-gray-600">{plan.description}</p>
          {/if}
          <p class="mt-3 text-3xl font-bold text-gray-900">
            {formatPrice(plan.price || 0)}
            <span class="text-lg font-normal text-gray-600"
              >/{plan.interval}</span
            >
          </p>

          <ul class="mt-6 space-y-3 text-sm text-gray-600">
            {#each getPlanFeatures(plan) as feature}
              <li class="flex items-start">
                <Check class="h-5 w-5 text-green-500 mr-2 flex-shrink-0" />
                <span>{feature}</span>
              </li>
            {/each}
          </ul>
        </div>

        <div class="mt-6">
          {#if isCurrent}
            <button
              class="w-full px-4 py-2 bg-gray-100 text-gray-500 font-medium rounded-lg cursor-not-allowed"
              disabled
            >
              Current Plan
            </button>
          {:else}
            <button
              class="w-full px-4 py-2 bg-gradient-to-r from-blue-600 to-purple-600 text-white font-medium rounded-lg hover:from-blue-700 hover:to-purple-700 transition-all duration-200 disabled:opacity-50 disabled:cursor-not-allowed"
              onclick={() => onSelectPlan(plan)}
              disabled={isLoading}
            >
              {#if isLoading}
                <span class="inline-flex items-center">
                  <Loader2 class="h-4 w-4 mr-2 animate-spin" />
                  Processing...
                </span>
              {:else}
                {actionLabel(plan)}
              {/if}
            </button>
          {/if}
        </div>
      </Card>
    </div>
  {/each}
</div>

