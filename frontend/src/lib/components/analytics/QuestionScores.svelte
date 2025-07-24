<script lang="ts">
  import { Card } from '$lib/components/ui';
  import { StarIcon } from 'lucide-svelte';
  import { getContext } from 'svelte';
  import type { analyticsStore as AnalyticsStoreType } from '$lib/stores/analytics';
  import { ANALYTICS_CONTEXT_KEY } from '$lib/stores/analytics';
  
  interface QuestionMetric {
    question_id?: string;
    question_text: string;
    question_type: string;
    response_count: number;
    average_score?: number;
    min_score?: number;
    max_score?: number;
    positive_rate?: number;
    negative_rate?: number;
    option_distribution?: Record<string, number>;
    text_responses?: string[];
  }
  
  interface ProductInsights {
    question_metrics: QuestionMetric[];
  }
  
  let {
    productInsights = null,
    loading = false
  }: {
    productInsights?: ProductInsights | null;
    loading?: boolean;
  } = $props();
  
  // Get analytics store from context
  const analytics = getContext<typeof AnalyticsStoreType>(ANALYTICS_CONTEXT_KEY);
  const { selection } = analytics;
  
  // Track highlighted question
  let highlightedQuestionId = $derived($selection.highlightedQuestionId);

  function getScoreTrend(score: number): { color: string; label: string } {
    if (score >= 4.5) return { color: 'text-green-600', label: 'Excellent' };
    if (score >= 3.5) return { color: 'text-blue-600', label: 'Good' };
    if (score >= 2.5) return { color: 'text-yellow-600', label: 'Fair' };
    return { color: 'text-red-600', label: 'Poor' };
  }

  function getResponseStats(metric: QuestionMetric): string {
    if (metric.question_type === 'numeric' && metric.average_score !== undefined && metric.average_score !== null) {
      return `${metric.average_score.toFixed(1)}/5.0`;
    }
    if (metric.question_type === 'choice' && metric.option_distribution) {
      const topOption = Object.entries(metric.option_distribution)
        .sort(([,a], [,b]) => b - a)[0];
      return topOption ? `${topOption[0]} (${topOption[1]})` : 'No responses';
    }
    if (metric.question_type === 'text' && metric.text_responses) {
      return `${metric.text_responses.length} responses`;
    }
    return `${metric.response_count || 0} responses`;
  }

  function getSentimentColor(positive: number, negative: number): string {
    if (positive > negative * 2) return 'bg-green-100 text-green-800';
    if (negative > positive * 2) return 'bg-red-100 text-red-800';
    return 'bg-gray-100 text-gray-800';
  }
</script>

<Card variant="default" class="question-scores">
  <div class="mb-6">
    <h3 class="text-lg font-semibold bg-gradient-to-r from-gray-900 to-gray-700 bg-clip-text text-transparent">
      Question Scores
    </h3>
    <p class="text-sm text-gray-600 mt-1">Average ratings and response counts per question</p>
  </div>

  {#if loading}
    <div class="space-y-4">
      {#each Array(5) as _}
        <div class="animate-pulse">
          <div class="bg-gray-100 border border-gray-200 rounded-lg p-4">
            <div class="flex items-center justify-between mb-2">
              <div class="h-4 bg-gray-200 rounded w-2/3"></div>
              <div class="h-4 bg-gray-200 rounded w-16"></div>
            </div>
            <div class="flex items-center gap-2 mb-2">
              <div class="h-3 bg-gray-200 rounded w-20"></div>
              <div class="h-3 bg-gray-200 rounded w-24"></div>
            </div>
            <div class="h-2 bg-gray-200 rounded w-full"></div>
          </div>
        </div>
      {/each}
    </div>
  {:else if productInsights?.question_metrics?.length}
    <div class="space-y-4">
      {#each productInsights.question_metrics as metric}
        <Card 
          variant="minimal"
          padding={false}
          hover={highlightedQuestionId !== metric.question_id}
          interactive
          class="w-full text-left p-4 cursor-pointer {highlightedQuestionId === metric.question_id ? 'border-blue-500 shadow-lg ring-2 ring-blue-200' : ''}"
          onclick={() => analytics.setHighlightedQuestion(metric.question_id === highlightedQuestionId ? null : metric.question_id)}
        >
          <div class="flex items-start justify-between mb-3">
            <div class="flex-1">
              <h4 class="font-medium text-gray-900 text-sm leading-relaxed">
                {metric.question_text}
              </h4>
              <div class="flex items-center gap-3 mt-2">
                <span class="text-xs text-gray-500 px-2 py-1 bg-gray-100 rounded">
                  {metric.question_type}
                </span>
                <span class="text-xs text-gray-600">
                  {metric.response_count} responses
                </span>
              </div>
            </div>
            <div class="text-right ml-4">
              <div class="text-sm font-medium text-gray-900">
                {getResponseStats(metric)}
              </div>
              {#if metric.average_score !== undefined}
                <div class="flex items-center gap-1 mt-1">
                  <StarIcon class="w-3 h-3 fill-yellow-400 text-yellow-400" />
                  <span class="text-xs {getScoreTrend(metric.average_score).color}">
                    {getScoreTrend(metric.average_score).label}
                  </span>
                </div>
              {/if}
            </div>
          </div>

          <!-- Progress bar for numeric questions -->
          {#if metric.question_type === 'numeric' && metric.average_score !== undefined && metric.average_score !== null}
            <div class="w-full bg-gray-200 rounded-full h-2 mb-2">
              <div 
                class="bg-gradient-to-r from-blue-500 to-purple-600 h-2 rounded-full transition-all duration-300"
                style="width: {(metric.average_score / 5) * 100}%"
              ></div>
            </div>
            <div class="flex justify-between text-xs text-gray-500">
              <span>Min: {metric.min_score?.toFixed(1) || 'N/A'}</span>
              <span>Max: {metric.max_score?.toFixed(1) || 'N/A'}</span>
            </div>
          {/if}

          <!-- Sentiment indicators -->
          {#if metric.positive_rate !== undefined && metric.negative_rate !== undefined}
            <div class="flex gap-2 mt-2">
              <span class="text-xs px-2 py-1 rounded-full {getSentimentColor(metric.positive_rate, metric.negative_rate)}">
                {metric.positive_rate.toFixed(0)}% positive, {metric.negative_rate.toFixed(0)}% negative
              </span>
            </div>
          {/if}

          <!-- Option distribution for choice questions -->
          {#if metric.question_type === 'choice' && metric.option_distribution}
            <div class="mt-3 space-y-1">
              {#each Object.entries(metric.option_distribution).sort(([,a], [,b]) => b - a).slice(0, 3) as [option, count]}
                <div class="flex items-center justify-between text-xs">
                  <span class="text-gray-600 truncate">{option}</span>
                  <span class="text-gray-500 font-medium">{count}</span>
                </div>
              {/each}
            </div>
          {/if}
        </Card>
      {/each}
    </div>
  {:else}
    <div class="text-center py-8">
      <div class="text-gray-500 text-sm">
        No question data available. Select a product to view question-level analytics.
      </div>
    </div>
  {/if}
</Card>
