<script lang="ts">
  import TimeSeriesChart from './TimeSeriesChart.svelte';
  import { Card } from '$lib/components/ui';
  import { BarChart3, TrendingUp, Info, ChevronDown, ChevronUp } from 'lucide-svelte';

  // State for collapsible sections - all collapsed by default
  let expandedSections = $state({
    general: false,
    individual: false
  });

  interface Props {
    data: any;
  }

  let { data }: Props = $props();

  // Only general metrics and individual questions
  let generalMetrics = $derived(
    data?.series?.filter((s: any) => 
      s.metric_type === 'survey_responses' || 
      (!s.metric_type.includes('questions') && !s.metric_type.includes('question_'))
    ) || []
  );

  // For individual questions (question_xxx format)
  let individualQuestions = $derived(
    data?.series?.filter((s: any) => s.metric_type.startsWith('question_')) || []
  );

  function toggleSection(section: keyof typeof expandedSections) {
    expandedSections[section] = !expandedSections[section];
  }
</script>

<div class="separated-timeseries-charts space-y-6">
  <!-- Information Banner -->
  <Card variant="glass" class="border-blue-200/60 bg-gradient-to-r from-blue-50/80 to-indigo-50/60">
    <div class="flex items-start gap-4">
      <div class="h-10 w-10 bg-gradient-to-br from-blue-500 to-indigo-600 rounded-xl flex items-center justify-center shadow-lg shadow-blue-500/25 flex-shrink-0">
        <Info class="w-5 h-5 text-white" />
      </div>
      <div>
        <h4 class="font-semibold text-blue-900 mb-2">Individual Question Charts</h4>
        <p class="text-sm text-blue-700/80 leading-relaxed">
          Each question is displayed individually with its appropriate scale and formatting for precise analysis.
        </p>
      </div>
    </div>
  </Card>

  {#if !data?.series || data.series.length === 0}
    <Card variant="minimal" class="text-center py-12">
      <div class="text-gray-500">
        <BarChart3 class="w-16 h-16 mx-auto mb-4 opacity-40" />
        <h3 class="text-lg font-medium mb-2">No Data Available</h3>
        <p class="text-sm">No time series data is available for the selected filters.</p>
      </div>
    </Card>
  {:else}
    <!-- General Metrics -->
    {#if generalMetrics.length > 0}
      <Card variant="minimal" padding={false} class="group hover:shadow-lg transition-all duration-300 border">
        <button 
          class="w-full flex items-center gap-4 p-6 border-b border-gray-100/60 hover:bg-gray-50/50 transition-colors duration-200 text-left cursor-pointer"
          onclick={() => toggleSection('general')}
        >
          <div class="h-12 w-12 bg-gradient-to-br from-gray-500 to-gray-700 rounded-2xl flex items-center justify-center shadow-lg shadow-gray-500/25">
            <TrendingUp class="w-6 h-6 text-white" />
          </div>
          <div class="flex-1">
            <h3 class="text-xl font-bold bg-gradient-to-r from-gray-900 to-gray-700 bg-clip-text text-transparent">
              General Metrics
            </h3>
            <p class="text-gray-600 font-medium mt-1">Overall survey response counts and general metrics</p>
          </div>
          <div class="flex items-center gap-2 text-gray-500">
            <span class="text-sm font-medium">{generalMetrics.length} metric{generalMetrics.length !== 1 ? 's' : ''}</span>
            {#if expandedSections.general}
              <ChevronUp class="w-5 h-5 transition-transform duration-200" />
            {:else}
              <ChevronDown class="w-5 h-5 transition-transform duration-200" />
            {/if}
          </div>
        </button>
        {#if expandedSections.general}
          <div class="p-6 border-t border-gray-100/60">
            <TimeSeriesChart data={{ ...data, series: generalMetrics }} />
          </div>
        {/if}
      </Card>
    {/if}


    <!-- Individual Questions - Separated by question -->
    {#if individualQuestions.length > 0}
      <Card variant="minimal" padding={false} class="group hover:shadow-lg transition-all duration-300 border">
        <button 
          class="w-full flex items-center gap-4 p-6 border-b border-gray-100/60 hover:bg-gray-50/50 transition-colors duration-200 text-left cursor-pointer"
          onclick={() => toggleSection('individual')}
        >
          <div class="h-12 w-12 bg-gradient-to-br from-slate-500 to-slate-700 rounded-2xl flex items-center justify-center shadow-lg shadow-slate-500/25">
            <BarChart3 class="w-6 h-6 text-white" />
          </div>
          <div class="flex-1">
            <h3 class="text-xl font-bold bg-gradient-to-r from-gray-900 to-gray-700 bg-clip-text text-transparent">
              Individual Questions
            </h3>
            <p class="text-gray-600 font-medium mt-1">Metrics for specific survey questions</p>
          </div>
          <div class="flex items-center gap-2 text-gray-500">
            <span class="text-sm font-medium">{individualQuestions.length} question{individualQuestions.length !== 1 ? 's' : ''}</span>
            {#if expandedSections.individual}
              <ChevronUp class="w-5 h-5 transition-transform duration-200" />
            {:else}
              <ChevronDown class="w-5 h-5 transition-transform duration-200" />
            {/if}
          </div>
        </button>
        {#if expandedSections.individual}
          <div class="space-y-4 p-6 border-t border-gray-100/60">
            {#each individualQuestions as question, index}
              <div class="border border-gray-200 rounded-xl overflow-hidden">
                <div class="bg-gray-50 px-4 py-3 border-b border-gray-200">
                  <h4 class="font-semibold text-gray-900 text-sm">{question.metric_name}</h4>
                </div>
                <div class="p-4">
                  <TimeSeriesChart data={{ ...data, series: [question] }} />
                </div>
              </div>
            {/each}
          </div>
        {/if}
      </Card>
    {/if}
  {/if}
</div>

