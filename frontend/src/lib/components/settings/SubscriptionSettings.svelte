<script lang="ts">
  import { Button } from '$lib/components/ui';
  import { Check } from 'lucide-svelte';

  interface Props {
    currentPlan?: 'starter' | 'professional' | 'enterprise';
    onUpgrade?: (plan: string) => void;
    onDowngrade?: (plan: string) => void;
  }

  let { currentPlan = 'professional', onUpgrade, onDowngrade }: Props = $props();

  const plans = [
    {
      id: 'starter',
      name: 'Starter',
      price: 29,
      features: [
        '1 restaurant',
        '500 feedbacks/month',
        '10 QR codes',
        '2 team members'
      ]
    },
    {
      id: 'professional',
      name: 'Professional',
      price: 79,
      features: [
        '3 restaurants',
        '2,000 feedbacks/month',
        '50 QR codes per location',
        '5 team members',
        'Advanced analytics'
      ]
    },
    {
      id: 'enterprise',
      name: 'Enterprise',
      price: 199,
      features: [
        'Unlimited restaurants',
        'Unlimited feedbacks',
        'Unlimited QR codes',
        'Unlimited team members',
        'API access & priority support'
      ]
    }
  ];

  function handlePlanAction(planId: string) {
    const currentPlanIndex = plans.findIndex(p => p.id === currentPlan);
    const targetPlanIndex = plans.findIndex(p => p.id === planId);
    
    if (targetPlanIndex > currentPlanIndex) {
      onUpgrade?.(planId);
    } else if (targetPlanIndex < currentPlanIndex) {
      onDowngrade?.(planId);
    }
  }
</script>

<div>
  <div class="mb-8">
    <h2 class="text-2xl font-bold text-gray-900">Subscription Plan</h2>
    <p class="mt-1 text-sm text-gray-600">View your current plan and upgrade options</p>
  </div>
  
  <!-- Current Plan -->
  <div class="mb-8">
    <div class="bg-gradient-to-r from-blue-500 to-purple-600 rounded-xl p-6 text-white">
      <div class="flex items-start justify-between">
        <div>
          <h3 class="text-sm font-medium text-blue-100">Current Plan</h3>
          <p class="mt-1 text-2xl font-bold capitalize">{currentPlan}</p>
          <p class="mt-2 text-sm text-blue-100">
            ${plans.find(p => p.id === currentPlan)?.price}/month
          </p>
        </div>
        <div class="text-right">
          <p class="text-sm text-blue-100">Next billing date</p>
          <p class="text-lg font-semibold">January 15, 2025</p>
        </div>
      </div>
      
      <div class="mt-6 grid grid-cols-2 gap-4 border-t border-white/20 pt-4">
        <div>
          <p class="text-sm text-blue-100">Restaurants</p>
          <p class="text-lg font-semibold">2 of 3 used</p>
        </div>
        <div>
          <p class="text-sm text-blue-100">Monthly Feedback</p>
          <p class="text-lg font-semibold">1,234 of 2,000</p>
        </div>
      </div>
    </div>
  </div>

  <!-- Available Plans -->
  <div class="space-y-4">
    <h3 class="text-lg font-medium text-gray-900">Available Plans</h3>
    
    <div class="grid gap-4 md:grid-cols-3">
      {#each plans as plan}
        <div class="border rounded-xl p-6 transition-colors {plan.id === currentPlan ? 'border-2 border-blue-500' : 'border-gray-200 hover:border-gray-300'}">
          {#if plan.id === currentPlan}
            <div class="absolute -top-3 left-1/2 -translate-x-1/2">
              <span class="bg-blue-500 text-white text-xs font-semibold px-3 py-1 rounded-full">Current Plan</span>
            </div>
          {/if}
          
          <h4 class="text-lg font-semibold text-gray-900">{plan.name}</h4>
          <p class="mt-2 text-3xl font-bold text-gray-900">
            ${plan.price}<span class="text-lg font-normal text-gray-600">/month</span>
          </p>
          
          <ul class="mt-6 space-y-3 text-sm text-gray-600">
            {#each plan.features as feature}
              <li class="flex items-start">
                <Check class="h-5 w-5 text-green-500 mr-2 flex-shrink-0" />
                {feature}
              </li>
            {/each}
          </ul>
          
          {#if plan.id === currentPlan}
            <Button variant="primary" size="md" class="w-full mt-6" disabled>
              Current Plan
            </Button>
          {:else}
            <Button 
              variant={plans.findIndex(p => p.id === plan.id) > plans.findIndex(p => p.id === currentPlan) ? 'gradient' : 'outline'} 
              size="md" 
              class="w-full mt-6"
              onclick={() => handlePlanAction(plan.id)}
              disabled={plan.id === 'starter' && currentPlan !== 'starter'}
            >
              {plans.findIndex(p => p.id === plan.id) > plans.findIndex(p => p.id === currentPlan) ? 'Upgrade' : 'Downgrade'}
            </Button>
          {/if}
        </div>
      {/each}
    </div>
  </div>
</div>