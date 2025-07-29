<script lang="ts">
  import TimeSeriesChart from './TimeSeriesChart.svelte';
  import ChoiceDistributionChart from './ChoiceDistributionChart.svelte';
  import { Card, NoDataAvailable } from '$lib/components/ui';
  import {
    BarChart3,
    TrendingUp,
    Info,
    ChevronDown,
    ChevronUp,
  } from 'lucide-svelte';

  // State for collapsible sections - all collapsed by default
  let expandedSections = $state({
    general: false,
    individual: false,
  });

  interface Props {
    data: any;
  }

  let { data }: Props = $props();

  // Only general metrics and individual questions
  let generalMetrics = $derived(
    data?.series?.filter(
      (s: any) =>
        s.metric_type === 'survey_responses' ||
        (!s.metric_type.includes('questions') &&
          !s.metric_type.includes('question_'))
    ) || []
  );

  // For individual questions (question_xxx format) - include ALL question types
  let individualQuestions = $derived(
    data?.series?.filter((s: any) => s.metric_type.startsWith('question_')) ||
      []
  );

  function toggleSection(section: keyof typeof expandedSections) {
    expandedSections[section] = !expandedSections[section];
  }
</script>

<div class="separated-timeseries-charts space-y-6">
  {#if !data?.series || data.series.length === 0}
    <NoDataAvailable
      title="No Data Available"
      description="No time series data is available for the selected filters"
      icon={BarChart3}
      variant="inline" />
  {:else}
    <!-- General Metrics -->
    {#if generalMetrics.length > 0}
      <Card
        variant="minimal"
        padding={false}
        class="group hover:shadow-lg transition-all duration-300 border">
        <button
          class="w-full flex items-center gap-4 p-6 border-b border-gray-100/60 hover:bg-gray-50/50 transition-colors duration-200 text-left cursor-pointer"
          onclick={() => toggleSection('general')}>
          <div
            class="h-12 w-12 bg-gradient-to-br from-gray-500 to-gray-700 rounded-2xl flex items-center justify-center shadow-lg shadow-gray-500/25">
            <TrendingUp class="w-6 h-6 text-white" />
          </div>
          <div class="flex-1">
            <h3
              class="text-xl font-bold bg-gradient-to-r from-gray-900 to-gray-700 bg-clip-text text-transparent">
              General Metrics
            </h3>
            <p class="text-gray-600 font-medium mt-1">
              Overall survey response counts and general metrics
            </p>
          </div>
          <div class="flex items-center gap-2 text-gray-500">
            <span class="text-sm font-medium"
              >{generalMetrics.length} metric{generalMetrics.length !== 1
                ? 's'
                : ''}</span>
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
      <Card
        variant="minimal"
        padding={false}
        class="group hover:shadow-lg transition-all duration-300 border">
        <button
          class="w-full flex items-center gap-4 p-6 border-b border-gray-100/60 hover:bg-gray-50/50 transition-colors duration-200 text-left cursor-pointer"
          onclick={() => toggleSection('individual')}>
          <div
            class="h-12 w-12 bg-gradient-to-br from-slate-500 to-slate-700 rounded-2xl flex items-center justify-center shadow-lg shadow-slate-500/25">
            <BarChart3 class="w-6 h-6 text-white" />
          </div>
          <div class="flex-1">
            <h3
              class="text-xl font-bold bg-gradient-to-r from-gray-900 to-gray-700 bg-clip-text text-transparent">
              Individual Questions
            </h3>
            <p class="text-gray-600 font-medium mt-1">
              Metrics for specific survey questions
            </p>
          </div>
          <div class="flex items-center gap-2 text-gray-500">
            <span class="text-sm font-medium"
              >{individualQuestions.length} question{individualQuestions.length !==
              1
                ? 's'
                : ''}</span>
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
                  <div class="flex items-center justify-between">
                    <h4 class="font-semibold text-gray-900 text-sm">
                      {#if question.metric_name.includes(' - ')}
                        <span class="text-blue-600"
                          >{question.metric_name.split(' - ')[0]}</span>
                        <span class="text-gray-500 mx-1">•</span>
                        <span>{question.metric_name.split(' - ')[1]}</span>
                      {:else}
                        {question.metric_name}
                      {/if}
                    </h4>
                    <span
                      class="text-xs px-2 py-1 rounded-full font-medium capitalize {question
                        .metadata?.has_choice_series === true
                        ? 'bg-blue-200 text-blue-700'
                        : 'bg-gray-200 text-gray-700'}">
                      {(() => {
                        if (!question.metadata) {
                          return question.metric_type.startsWith('question_')
                            ? 'Question'
                            : 'Metric';
                        }
                        const metadata = question.metadata;
                        if (metadata?.has_choice_series === true) {
                          const questionType =
                            metadata?.question_type || 'single_choice';
                          const typeLabel =
                            questionType === 'multi_choice'
                              ? 'Multiple Choice'
                              : 'Single Choice';
                          return `${typeLabel} (${question.choice_series?.length || 0} options)`;
                        }
                        return (
                          metadata?.question_type?.replace('_', ' ') ||
                          'Unknown'
                        );
                      })()}
                    </span>
                  </div>
                </div>
                <div class="p-4">
                  {#if question.metadata?.has_choice_series === true}
                    <ChoiceDistributionChart data={question} />
                  {:else}
                    <TimeSeriesChart data={{ ...data, series: [question] }} />
                  {/if}
                </div>
              </div>
            {/each}
          </div>
        {/if}
      </Card>
    {/if}
  {/if}
</div>
