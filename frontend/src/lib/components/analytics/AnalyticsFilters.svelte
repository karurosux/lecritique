<script lang="ts">
  import { Card, Select } from '$lib/components/ui';

  interface Organization {
    id: string;
    name: string;
  }

  let {
    organizations = [],
    selectedOrganization = $bindable(''),
    selectedTimeframe = $bindable('7d'),
    onorganizationchange = () => {},
    ontimeframechange = () => {},
  }: {
    organizations?: Organization[];
    selectedOrganization?: string;
    selectedTimeframe?: string;
    onorganizationchange?: () => void;
    ontimeframechange?: () => void;
  } = $props();

  let organizationCount = $derived(organizations.length);
  let selectedOrganizationName = $derived(
    organizations.find(r => r.id === selectedOrganization)?.name ||
      'Select Organization'
  );
</script>

<Card
  variant="default"
  hover
  interactive
  class="mb-6 group transform transition-all duration-300 animate-fade-in-up">
  <div
    class="flex flex-col sm:flex-row gap-4 items-start sm:items-center justify-between">
    <div class="flex flex-col sm:flex-row gap-4 items-start sm:items-center">
      <!-- Organization Selector -->
      <div class="space-y-1">
        <label
          class="text-xs font-medium text-gray-500 uppercase tracking-wider">
          Organization
        </label>
        <Select
          bind:value={selectedOrganization}
          options={organizations.map(r => ({ value: r.id, label: r.name }))}
          onchange={() => onorganizationchange()}
          minWidth="min-w-[200px]" />
      </div>

      <!-- Timeframe Selector -->
      <div class="space-y-1">
        <label
          class="text-xs font-medium text-gray-500 uppercase tracking-wider">
          Time Period
        </label>
        <Select
          bind:value={selectedTimeframe}
          options={[
            { value: '1d', label: 'Last 24 hours' },
            { value: '7d', label: 'Last 7 days' },
            { value: '30d', label: 'Last 30 days' },
            { value: '90d', label: 'Last 90 days' },
            { value: '1y', label: 'Last year' },
          ]}
          onchange={() => ontimeframechange()} />
      </div>
    </div>

    <div class="flex items-center space-x-4 text-sm">
      <div class="flex items-center space-x-2">
        <svg
          class="h-4 w-4 text-gray-400"
          fill="none"
          stroke="currentColor"
          viewBox="0 0 24 24">
          <path
            stroke-linecap="round"
            stroke-linejoin="round"
            stroke-width="2"
            d="M19 21V5a2 2 0 00-2-2H7a2 2 0 00-2 2v16m14 0h2m-2 0h-4m-5 0H9m11 0a2 2 0 01-2 2H7a2 2 0 01-2-2m2-4h2.01M7 16h6M7 8h6v4H7V8z" />
        </svg>
        <span class="text-gray-600 font-medium"
          >{organizationCount}
          {organizationCount === 1 ? 'Organization' : 'Organizations'}</span>
      </div>
      <div class="h-4 w-px bg-gray-200"></div>
      <div class="flex items-center space-x-2">
        <svg
          class="h-4 w-4 text-gray-400"
          fill="none"
          stroke="currentColor"
          viewBox="0 0 24 24">
          <path
            stroke-linecap="round"
            stroke-linejoin="round"
            stroke-width="2"
            d="M8 7V3m8 4V3m-9 8h10M5 21h14a2 2 0 002-2V7a2 2 0 00-2-2H5a2 2 0 00-2 2v12a2 2 0 002 2z" />
        </svg>
        <span class="text-gray-600 font-medium">
          {#if selectedTimeframe === '1d'}
            Last 24 hours
          {:else if selectedTimeframe === '7d'}
            Last 7 days
          {:else if selectedTimeframe === '30d'}
            Last 30 days
          {:else if selectedTimeframe === '90d'}
            Last 90 days
          {:else if selectedTimeframe === '1y'}
            Last year
          {/if}
        </span>
      </div>
    </div>
  </div>
</Card>

<style>
  @keyframes fade-in-up {
    from {
      opacity: 0;
      transform: translateY(10px);
    }
    to {
      opacity: 1;
      transform: translateY(0);
    }
  }

  .animate-fade-in-up {
    animation: fade-in-up 0.6s ease-out forwards;
    opacity: 0;
  }
</style>
