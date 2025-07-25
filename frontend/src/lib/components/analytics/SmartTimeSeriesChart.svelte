<script lang="ts">
  import TimeSeriesChart from './TimeSeriesChart.svelte';
  import { BarChart3, Star, BarChart2, CheckCircle, MessageCircle, List, TrendingUp } from 'lucide-svelte';

  interface Props {
    data: any;
  }

  let { data }: Props = $props();

  // Group series by question type for separate visualization
  let groupedSeries = $derived(() => {
    if (!data?.series) {
      console.log('SmartTimeSeriesChart: No series data', data);
      return {};
    }
    
    console.log('SmartTimeSeriesChart: Processing series data', data.series);
    const groups: Record<string, any[]> = {};
    
    data.series.forEach((series: any) => {
      let questionType = 'general';
      
      console.log('Processing series:', series.metric_type, 'metadata:', series.metadata);
      
      // Determine question type
      if (series.metric_type.startsWith('question_')) {
        // Try to get question type from metadata
        if (series.metadata) {
          let metadata = series.metadata;
          if (typeof metadata === 'string') {
            try {
              metadata = JSON.parse(metadata);
            } catch (e) {
              // ignore parsing errors
            }
          }
          if (metadata?.question_type) {
            questionType = metadata.question_type;
          }
        }
      }
      
      // Also check the metric_type itself for question types
      if (questionType === 'general') {
        if (series.metric_type.includes('rating_questions')) questionType = 'rating';
        else if (series.metric_type.includes('scale_questions')) questionType = 'scale';
        else if (series.metric_type.includes('yes_no_questions')) questionType = 'yes_no';
        else if (series.metric_type.includes('text_questions')) questionType = 'text';
        else if (series.metric_type.includes('single_choice_questions')) questionType = 'choice';
        else if (series.metric_type.includes('multiple_choice_questions')) questionType = 'choice';
        else if (series.metric_type === 'survey_responses') questionType = 'general';
        
        // Check for any other patterns
        else if (series.metric_type.includes('rating')) questionType = 'rating';
        else if (series.metric_type.includes('scale')) questionType = 'scale';
        else if (series.metric_type.includes('yes_no')) questionType = 'yes_no';
        else if (series.metric_type.includes('text')) questionType = 'text';
        else if (series.metric_type.includes('choice')) questionType = 'choice';
      }
      
      console.log(`Series ${series.metric_type} assigned to group: ${questionType}`);
      
      if (!groups[questionType]) {
        groups[questionType] = [];
      }
      groups[questionType].push(series);
    });
    
    console.log('Final grouped series:', groups);
    return groups;
  });

  function getGroupIcon(questionType: string) {
    switch (questionType) {
      case 'rating': return Star;
      case 'scale': return BarChart2;
      case 'yes_no': return CheckCircle;
      case 'text': return MessageCircle;
      case 'choice': return List;
      case 'general': return TrendingUp;
      default: return BarChart3;
    }
  }

  function getGroupTitle(questionType: string): string {
    switch (questionType) {
      case 'rating': return 'Rating Questions (1-5 Stars)';
      case 'scale': return 'Scale Questions (1-10)';
      case 'yes_no': return 'Yes/No Questions (%)';
      case 'text': return 'Text Sentiment Analysis';
      case 'choice': return 'Choice Questions (Response Counts)';
      case 'general': return 'General Metrics';
      default: return `${questionType} Questions`;
    }
  }

  function getGroupDescription(questionType: string): string {
    switch (questionType) {
      case 'rating': return 'Average star ratings over time (1-5 scale)';
      case 'scale': return 'Average scale responses over time (1-10 scale)';
      case 'yes_no': return 'Percentage of "Yes" responses over time';
      case 'text': return 'Sentiment scores from text analysis (-1 to +1)';
      case 'choice': return 'Number of responses for choice-based questions';
      case 'general': return 'Overall survey response counts and general metrics';
      default: return `Metrics for ${questionType} type questions`;
    }
  }

  $effect(() => {
    console.log('SmartTimeSeriesChart - data:', data);
    console.log('SmartTimeSeriesChart - groupedSeries:', groupedSeries);
    console.log('SmartTimeSeriesChart - groupedSeries keys:', Object.keys(groupedSeries));
  });
</script>

<div class="smart-timeseries-chart">
  {#if !data?.series || data.series.length === 0}
    <div class="text-center py-8 text-gray-500">
      <BarChart3 class="w-12 h-12 mx-auto mb-4 opacity-50" />
      <p>No time series data available</p>
    </div>
  {:else if Object.keys(groupedSeries).length === 0}
    <!-- Fallback: show original chart if grouping produces no results -->
    <TimeSeriesChart {data} />
  {:else}
    {#each Object.entries(groupedSeries) as [questionType, seriesGroup]}
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
          <TimeSeriesChart 
            data={{
              ...data,
              series: seriesGroup
            }} 
          />
        </div>
      </div>
    {/each}
  {/if}
</div>

<style>
  .question-group {
    box-shadow: 0 1px 3px 0 rgba(0, 0, 0, 0.1), 0 1px 2px 0 rgba(0, 0, 0, 0.06);
  }
</style>