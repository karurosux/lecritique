<script lang="ts">
  import ComparisonChart from './ComparisonChart.svelte';
  import {
    Activity,
    Star,
    BarChart2,
    CheckCircle,
    MessageCircle,
    List,
    TrendingUp,
  } from 'lucide-svelte';

  interface Props {
    data: any;
  }

  let { data }: Props = $props();

  let groupedComparisons = $derived(() => {
    if (!data?.comparisons) return {};

    const groups: Record<string, any[]> = {};

    data.comparisons.forEach((comparison: any) => {
      let questionType = 'general';

      if (comparison.metric_type.startsWith('question_')) {
        if (comparison.metadata) {
          let metadata = comparison.metadata;
          if (typeof metadata === 'string') {
            try {
              metadata = JSON.parse(metadata);
            } catch (e) {
            }
          }
          if (metadata?.question_type) {
            questionType = metadata.question_type;
          }
        }
      } else {
        if (comparison.metric_type.includes('rating')) questionType = 'rating';
        else if (comparison.metric_type.includes('scale'))
          questionType = 'scale';
        else if (comparison.metric_type.includes('yes_no'))
          questionType = 'yes_no';
        else if (comparison.metric_type.includes('text')) questionType = 'text';
        else if (comparison.metric_type.includes('choice'))
          questionType = 'choice';
        else if (comparison.metric_type.includes('survey_responses'))
          questionType = 'general';
      }

      if (!groups[questionType]) {
        groups[questionType] = [];
      }
      groups[questionType].push(comparison);
    });

    return groups;
  });

  function getGroupIcon(questionType: string) {
    switch (questionType) {
      case 'rating':
        return Star;
      case 'scale':
        return BarChart2;
      case 'yes_no':
        return CheckCircle;
      case 'text':
        return MessageCircle;
      case 'choice':
        return List;
      case 'general':
        return TrendingUp;
      default:
        return Activity;
    }
  }

  function getGroupTitle(questionType: string): string {
    switch (questionType) {
      case 'rating':
        return 'Rating Comparisons (1-5 Stars)';
      case 'scale':
        return 'Scale Comparisons (1-10)';
      case 'yes_no':
        return 'Yes/No Comparisons (%)';
      case 'text':
        return 'Text Sentiment Comparisons';
      case 'choice':
        return 'Choice Response Comparisons';
      case 'general':
        return 'General Metric Comparisons';
      default:
        return `${questionType} Comparisons`;
    }
  }

  function getGroupDescription(questionType: string): string {
    switch (questionType) {
      case 'rating':
        return 'Compare average star ratings between periods';
      case 'scale':
        return 'Compare average scale responses between periods';
      case 'yes_no':
        return 'Compare percentage of "Yes" responses between periods';
      case 'text':
        return 'Compare sentiment scores between periods';
      case 'choice':
        return 'Compare response counts for choice questions between periods';
      case 'general':
        return 'Compare general survey metrics between periods';
      default:
        return `Compare ${questionType} metrics between periods`;
    }
  }
</script>

<div class="smart-comparison-chart">
  {#if !data?.comparisons || data.comparisons.length === 0}
    <div class="text-center py-8 text-gray-500">
      <Activity class="w-12 h-12 mx-auto mb-4 opacity-50" />
      <p>No comparison data available</p>
    </div>
  {:else if Object.keys(groupedComparisons).length === 0}
    
    <ComparisonChart {data} />
  {:else}
    {#each Object.entries(groupedComparisons) as [questionType, comparisonGroup]}
      <div class="question-group mb-8 bg-white rounded-lg border">
        <div class="px-6 py-4 border-b border-gray-200">
          <div class="flex items-center gap-3">
            {#snippet icon()}
              {@const Icon = getGroupIcon(questionType)}
              <Icon class="w-5 h-5 text-gray-700" />
            {/snippet}
            {@render icon()}
            <div>
              <h3 class="text-lg font-semibold text-gray-900">
                {getGroupTitle(questionType)}
              </h3>
              <p class="text-sm text-gray-600 mt-1">
                {getGroupDescription(questionType)}
              </p>
            </div>
          </div>
        </div>

        <div class="p-6">
          <ComparisonChart
            data={{
              ...data,
              comparisons: comparisonGroup,
            }} />
        </div>
      </div>
    {/each}
  {/if}
</div>

<style>
  .question-group {
    box-shadow:
      0 1px 3px 0 rgba(0, 0, 0, 0.1),
      0 1px 2px 0 rgba(0, 0, 0, 0.06);
  }
</style>
