<script lang="ts">
  import { Card, Select, Button } from '$lib/components/ui';
  import { getContext } from 'svelte';
  import type { analyticsStore as AnalyticsStoreType } from '$lib/stores/analytics';
  import { ANALYTICS_CONTEXT_KEY } from '$lib/stores/analytics';
  import { CalendarIcon, ClockIcon, UsersIcon, BarChart3Icon } from 'lucide-svelte';
  
  // Get analytics store from context
  const analytics = getContext<typeof AnalyticsStoreType>(ANALYTICS_CONTEXT_KEY);
  const { filters, selection } = analytics;
  
  // Reactive values
  let currentFilters = $derived($filters);
  let currentSelection = $derived($selection);
  
  // Filter options
  const timeframeOptions = [
    { value: '24h', label: 'Last 24 Hours' },
    { value: '7d', label: 'Last 7 Days' },
    { value: '30d', label: 'Last 30 Days' },
    { value: '90d', label: 'Last 90 Days' }
  ];
  
  const segmentOptions = {
    timeOfDay: [
      { value: 'all', label: 'All Day' },
      { value: 'breakfast', label: 'Breakfast (6-11 AM)' },
      { value: 'lunch', label: 'Lunch (11 AM-3 PM)' },
      { value: 'dinner', label: 'Dinner (5-10 PM)' }
    ],
    dayType: [
      { value: 'all', label: 'All Days' },
      { value: 'weekday', label: 'Weekdays' },
      { value: 'weekend', label: 'Weekends' }
    ]
  };
  
  const comparisonModes = [
    { value: 'none', label: 'No Comparison' },
    { value: 'period', label: 'Period over Period' },
    { value: 'product', label: 'Compare Products' },
    { value: 'question', label: 'Compare Questions' }
  ];
  
  function handleTimeframeChange(value: string) {
    analytics.updateFilters({ timeframe: value as any });
  }
  
  function handleSegmentChange(type: string, value: string) {
    analytics.updateFilters({
      segments: {
        ...currentFilters.segments,
        [type]: value
      }
    });
  }
  
  function handleComparisonModeChange(value: string) {
    analytics.setComparisonMode(value as any);
  }
  
  function resetFilters() {
    analytics.updateFilters({
      timeframe: '7d',
      segments: {
        timeOfDay: 'all',
        dayType: 'all'
      }
    });
    analytics.setComparisonMode('none');
  }
  
  // Check if any filters are active
  const hasActiveFilters = $derived(() => {
    return currentFilters.timeframe !== '7d' ||
           currentFilters.segments.timeOfDay !== 'all' ||
           currentFilters.segments.dayType !== 'all' ||
           currentSelection.comparisonMode !== 'none';
  });
</script>

<Card variant="default" class="dynamic-filters">
  <div class="p-4">
    <div class="flex items-center justify-between mb-4">
      <h3 class="text-lg font-semibold text-gray-900">Analytics Filters</h3>
      {#if hasActiveFilters()}
        <Button 
          variant="ghost" 
          size="sm" 
          onclick={resetFilters}
        >
          Reset All
        </Button>
      {/if}
    </div>
    
    <div class="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">
      <!-- Timeframe -->
      <div>
        <label class="flex items-center gap-2 text-sm font-medium text-gray-700 mb-2">
          <CalendarIcon class="w-4 h-4" />
          Timeframe
        </label>
        <Select
          value={currentFilters.timeframe}
          onchange={(e) => handleTimeframeChange(e.target.value)}
          class="w-full"
        >
          {#each timeframeOptions as option}
            <option value={option.value}>{option.label}</option>
          {/each}
        </Select>
      </div>
      
      <!-- Time of Day -->
      <div>
        <label class="flex items-center gap-2 text-sm font-medium text-gray-700 mb-2">
          <ClockIcon class="w-4 h-4" />
          Time of Day
        </label>
        <Select
          value={currentFilters.segments.timeOfDay}
          onchange={(e) => handleSegmentChange('timeOfDay', e.target.value)}
          class="w-full"
        >
          {#each segmentOptions.timeOfDay as option}
            <option value={option.value}>{option.label}</option>
          {/each}
        </Select>
      </div>
      
      <!-- Day Type -->
      <div>
        <label class="flex items-center gap-2 text-sm font-medium text-gray-700 mb-2">
          <UsersIcon class="w-4 h-4" />
          Day Type
        </label>
        <Select
          value={currentFilters.segments.dayType}
          onchange={(e) => handleSegmentChange('dayType', e.target.value)}
          class="w-full"
        >
          {#each segmentOptions.dayType as option}
            <option value={option.value}>{option.label}</option>
          {/each}
        </Select>
      </div>
      
      <!-- Comparison Mode -->
      <div>
        <label class="flex items-center gap-2 text-sm font-medium text-gray-700 mb-2">
          <BarChart3Icon class="w-4 h-4" />
          Comparison
        </label>
        <Select
          value={currentSelection.comparisonMode}
          onchange={(e) => handleComparisonModeChange(e.target.value)}
          class="w-full"
        >
          {#each comparisonModes as mode}
            <option value={mode.value}>{mode.label}</option>
          {/each}
        </Select>
      </div>
    </div>
    
    <!-- Active Filters Summary -->
    {#if hasActiveFilters()}
      <div class="mt-4 p-3 bg-blue-50 rounded-lg">
        <div class="text-sm text-blue-900">
          <span class="font-medium">Active filters:</span>
          <span class="ml-2">
            {currentFilters.timeframe !== '7d' ? `${timeframeOptions.find(o => o.value === currentFilters.timeframe)?.label}` : ''}
            {currentFilters.segments.timeOfDay !== 'all' ? `, ${segmentOptions.timeOfDay.find(o => o.value === currentFilters.segments.timeOfDay)?.label}` : ''}
            {currentFilters.segments.dayType !== 'all' ? `, ${segmentOptions.dayType.find(o => o.value === currentFilters.segments.dayType)?.label}` : ''}
            {currentSelection.comparisonMode !== 'none' ? `, ${comparisonModes.find(m => m.value === currentSelection.comparisonMode)?.label}` : ''}
          </span>
        </div>
      </div>
    {/if}
  </div>
</Card>
