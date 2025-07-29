<script lang="ts">
  import { Card } from '$lib/components/ui';
  import { BarChart3Icon, ArrowUpIcon, ArrowDownIcon } from 'lucide-svelte';
  import { getContext } from 'svelte';
  import type { analyticsStore as AnalyticsStoreType } from '$lib/stores/analytics';
  import { ANALYTICS_CONTEXT_KEY } from '$lib/stores/analytics';

  interface QuestionMetric {
    question_id?: string;
    question_text: string;
    question_type: string;
    response_count: number;
    average_score?: number;
    positive_rate?: number;
    negative_rate?: number;
  }

  interface ProductInsights {
    question_metrics: QuestionMetric[];
    best_aspects?: string[];
    areas_needing_attention?: string[];
  }

  let {
    productInsights = null,
    loading = false,
  }: {
    productInsights?: ProductInsights | null;
    loading?: boolean;
  } = $props();

  let hoveredBar: number | null = null;
  let animationProgress = $state(0);

  // Get analytics store from context
  const analytics = getContext<typeof AnalyticsStoreType>(
    ANALYTICS_CONTEXT_KEY
  );
  const { selection } = analytics;

  // Track highlighted question
  let highlightedQuestionId = $derived($selection.highlightedQuestionId);

  const chartData = $derived(() => {
    if (!productInsights?.question_metrics) return null;

    const scoredQuestions = productInsights.question_metrics
      .filter(q => q.average_score !== undefined && q.average_score !== null)
      .sort((a, b) => (b.average_score || 0) - (a.average_score || 0));

    if (scoredQuestions.length === 0) return null;

    const maxScore = 5;
    const avgScore =
      scoredQuestions.reduce((sum, q) => sum + (q.average_score || 0), 0) /
      scoredQuestions.length;

    return {
      questions: scoredQuestions,
      maxScore,
      avgScore,
      bestPerformers: scoredQuestions.filter(
        q => (q.average_score || 0) >= 4.5
      ),
      needsImprovement: scoredQuestions.filter(q => (q.average_score || 0) < 3),
    };
  });

  function getBarColor(score: number, avgScore: number): string {
    if (score >= 4.5) return '#10b981';
    if (score >= avgScore) return '#3b82f6';
    if (score >= 3) return '#f59e0b';
    return '#ef4444';
  }

  function truncateText(text: string, maxLength: number = 40): string {
    if (text.length <= maxLength) return text;
    return text.substring(0, maxLength - 3) + '...';
  }

  $effect(() => {
    // Animate on mount or data change
    animationProgress = 0;
    const animate = () => {
      if (animationProgress < 1) {
        animationProgress = Math.min(animationProgress + 0.02, 1);
        requestAnimationFrame(animate);
      }
    };
    animate();
  });
</script>

<Card variant="default" class="question-comparison-chart">
  <div class="mb-6">
    <h3
      class="text-lg font-semibold bg-gradient-to-r from-gray-900 to-gray-700 bg-clip-text text-transparent flex items-center gap-2">
      <BarChart3Icon class="w-5 h-5" />
      Question Performance Comparison
    </h3>
    <p class="text-sm text-gray-600 mt-1">
      Side-by-side comparison of all question scores
    </p>
  </div>

  {#if loading}
    <div class="space-y-4">
      {#each Array(5) as _}
        <div class="animate-pulse">
          <div class="flex items-center gap-4">
            <div class="h-8 bg-gray-200 rounded flex-1"></div>
            <div class="h-8 bg-gray-200 rounded w-16"></div>
          </div>
        </div>
      {/each}
    </div>
  {:else if chartData}
    <!-- Summary Stats -->
    <div class="grid grid-cols-1 md:grid-cols-3 gap-4 mb-6">
      <div class="bg-blue-50 border border-blue-200 rounded-lg p-3">
        <div class="text-sm font-medium text-blue-900">Average Score</div>
        <div class="text-xl font-bold text-blue-600">
          {chartData.avgScore.toFixed(2)}/5.0
        </div>
      </div>

      <div class="bg-green-50 border border-green-200 rounded-lg p-3">
        <div class="text-sm font-medium text-green-900">Strong Points</div>
        <div class="text-xl font-bold text-green-600">
          {chartData.bestPerformers.length} questions
        </div>
      </div>

      <div class="bg-amber-50 border border-amber-200 rounded-lg p-3">
        <div class="text-sm font-medium text-amber-900">Need Attention</div>
        <div class="text-xl font-bold text-amber-600">
          {chartData.needsImprovement.length} questions
        </div>
      </div>
    </div>

    <!-- Horizontal Bar Chart -->
    <div class="space-y-3">
      {#each chartData.questions as question, index}
        {@const score = question.average_score || 0}
        {@const barWidth =
          (score / chartData.maxScore) * 100 * animationProgress}
        {@const isHovered = hoveredBar === index}

        <button
          class="w-full group text-left {highlightedQuestionId ===
          question.question_id
            ? 'scale-105 z-10 relative'
            : ''}"
          onmouseenter={() => (hoveredBar = index)}
          onmouseleave={() => (hoveredBar = null)}
          onclick={() =>
            analytics.setHighlightedQuestion(
              question.question_id === highlightedQuestionId
                ? null
                : question.question_id
            )}>
          <!-- Question Text -->
          <div class="flex items-center justify-between mb-1">
            <div class="flex items-center gap-2 flex-1">
              <span class="text-sm text-gray-700 font-medium">
                {truncateText(question.question_text)}
              </span>
              {#if score >= 4.5}
                <ArrowUpIcon class="w-3 h-3 text-green-600" />
              {:else if score < 3}
                <ArrowDownIcon class="w-3 h-3 text-red-600" />
              {/if}
            </div>
            <span
              class="text-sm font-bold ml-2"
              style="color: {getBarColor(score, chartData.avgScore)}">
              {score.toFixed(1)}
            </span>
          </div>

          <!-- Bar -->
          <div class="relative">
            <div class="w-full bg-gray-100 rounded-full h-6 overflow-hidden">
              <div
                class="h-full rounded-full transition-all duration-500 ease-out relative"
                style="width: {barWidth}%; background-color: {getBarColor(
                  score,
                  chartData.avgScore
                )}">
                <!-- Animated shine effect on hover -->
                {#if isHovered}
                  <div
                    class="absolute inset-0 bg-gradient-to-r from-transparent via-white to-transparent opacity-30 animate-shine">
                  </div>
                {/if}
              </div>
            </div>

            <!-- Average line -->
            <div
              class="absolute top-0 h-full w-0.5 bg-gray-600 opacity-50"
              style="left: {(chartData.avgScore / chartData.maxScore) * 100}%"
              title="Average: {chartData.avgScore.toFixed(2)}">
              {#if index === 0}
                <span
                  class="absolute -top-5 left-1/2 -translate-x-1/2 text-xs text-gray-600 whitespace-nowrap">
                  Avg: {chartData.avgScore.toFixed(1)}
                </span>
              {/if}
            </div>
          </div>

          <!-- Additional Info on Hover -->
          {#if isHovered}
            <div class="mt-1 text-xs text-gray-600 flex items-center gap-3">
              <span>{question.response_count} responses</span>
              {#if question.positive_rate !== undefined}
                <span class="text-green-600"
                  >{question.positive_rate.toFixed(0)}% positive</span>
              {/if}
              {#if question.negative_rate !== undefined}
                <span class="text-red-600"
                  >{question.negative_rate.toFixed(0)}% negative</span>
              {/if}
            </div>
          {/if}
        </button>
      {/each}
    </div>

    <!-- Insights -->
    {#if productInsights?.best_aspects?.length || productInsights?.areas_needing_attention?.length}
      <div class="mt-6 grid grid-cols-1 md:grid-cols-2 gap-4">
        {#if productInsights.best_aspects?.length}
          <div class="bg-green-50 border border-green-200 rounded-lg p-4">
            <h4 class="font-medium text-green-900 mb-2 flex items-center gap-2">
              <ArrowUpIcon class="w-4 h-4" />
              Best Performing Aspects
            </h4>
            <ul class="space-y-1 text-sm text-green-700">
              {#each productInsights.best_aspects.slice(0, 3) as aspect}
                <li class="flex items-start">
                  <span class="mr-2">•</span>
                  <span>{aspect}</span>
                </li>
              {/each}
            </ul>
          </div>
        {/if}

        {#if productInsights.areas_needing_attention?.length}
          <div class="bg-amber-50 border border-amber-200 rounded-lg p-4">
            <h4 class="font-medium text-amber-900 mb-2 flex items-center gap-2">
              <ArrowDownIcon class="w-4 h-4" />
              Areas Needing Attention
            </h4>
            <ul class="space-y-1 text-sm text-amber-700">
              {#each productInsights.areas_needing_attention.slice(0, 3) as area}
                <li class="flex items-start">
                  <span class="mr-2">•</span>
                  <span>{area}</span>
                </li>
              {/each}
            </ul>
          </div>
        {/if}
      </div>
    {/if}
  {:else}
    <div class="text-center py-12">
      <BarChart3Icon class="w-12 h-12 text-gray-400 mx-auto mb-4" />
      <div class="text-gray-500 text-sm">
        No scored questions available for comparison.
      </div>
    </div>
  {/if}
</Card>

<style>
  @keyframes shine {
    0% {
      transform: translateX(-100%);
    }
    100% {
      transform: translateX(100%);
    }
  }

  .animate-shine {
    animation: shine 1s ease-in-out;
  }
</style>
