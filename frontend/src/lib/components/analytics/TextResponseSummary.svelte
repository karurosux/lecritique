<script lang="ts">
  import { Card } from '$lib/components/ui';
  import {
    MessageSquareIcon,
    ThumbsUpIcon,
    ThumbsDownIcon,
    HashIcon,
  } from 'lucide-svelte';

  interface TextAnalysis {
    frequent_keywords: Array<{ keyword: string; count: number }>;
    positive_mentions: Array<{ text: string; count: number }>;
    negative_mentions: Array<{ text: string; count: number }>;
    recent_feedback: Array<{
      id: string;
      text: string;
      rating?: number;
      product_name?: string;
      created_at: string;
    }>;
    sentiment_summary: {
      positive_rate: number;
      negative_rate: number;
      neutral_rate: number;
    };
  }

  let {
    textAnalysis = null,
    loading = false,
  }: {
    textAnalysis?: TextAnalysis | null;
    loading?: boolean;
  } = $props();

  function formatDate(dateString: string): string {
    const date = new Date(dateString);
    return date.toLocaleDateString('en-US', {
      month: 'short',
      day: 'numeric',
      hour: '2-digit',
      minute: '2-digit',
    });
  }

  function getSentimentColor(rate: number): string {
    if (!rate || rate < 0) return 'text-gray-600';
    if (rate >= 60) return 'text-green-600';
    if (rate >= 40) return 'text-yellow-600';
    return 'text-red-600';
  }

  function truncateText(text: string, maxLength: number = 100): string {
    if (text.length <= maxLength) return text;
    return text.substring(0, maxLength) + '...';
  }

  const topKeywords = $derived(() => {
    if (!textAnalysis?.frequent_keywords) return [];
    return textAnalysis.frequent_keywords.slice(0, 10);
  });

  const topPositive = $derived(() => {
    if (!textAnalysis?.positive_mentions) return [];
    return textAnalysis.positive_mentions.slice(0, 5);
  });

  const topNegative = $derived(() => {
    if (!textAnalysis?.negative_mentions) return [];
    return textAnalysis.negative_mentions.slice(0, 5);
  });

  const recentComments = $derived(() => {
    if (!textAnalysis?.recent_feedback) return [];
    return textAnalysis.recent_feedback.slice(0, 8);
  });
</script>

<Card variant="default" class="text-response-summary">
  <div class="mb-6">
    <h3
      class="text-lg font-semibold bg-gradient-to-r from-gray-900 to-gray-700 bg-clip-text text-transparent">
      Text Response Summary
    </h3>
    <p class="text-sm text-gray-600 mt-1">
      Key themes, sentiment analysis, and recent feedback
    </p>
  </div>

  {#if loading}
    <div class="space-y-6">
      
      <div class="grid grid-cols-1 md:grid-cols-3 gap-4">
        {#each Array(3) as _}
          <div class="animate-pulse">
            <div class="bg-gray-100 border border-gray-200 rounded-lg p-4">
              <div class="h-4 bg-gray-200 rounded w-1/2 mb-2"></div>
              <div class="h-6 bg-gray-200 rounded w-1/4"></div>
            </div>
          </div>
        {/each}
      </div>

      
      <div class="space-y-4">
        <div class="h-4 bg-gray-200 rounded w-32"></div>
        <div class="flex flex-wrap gap-2">
          {#each Array(8) as _}
            <div class="h-6 bg-gray-200 rounded w-16"></div>
          {/each}
        </div>
      </div>
    </div>
  {:else if textAnalysis}
    
    {#if textAnalysis.sentiment_summary}
      <div class="grid grid-cols-1 md:grid-cols-3 gap-4 mb-6">
        <div class="bg-green-50 border border-green-200 rounded-lg p-4">
          <div class="flex items-center">
            <ThumbsUpIcon class="w-5 h-5 text-green-600 mr-2" />
            <div>
              <div class="text-sm font-medium text-green-900">Positive</div>
              <div class="text-lg font-semibold text-green-600">
                {(textAnalysis.sentiment_summary.positive_rate || 0).toFixed(
                  1
                )}%
              </div>
            </div>
          </div>
        </div>

        <div class="bg-gray-50 border border-gray-200 rounded-lg p-4">
          <div class="flex items-center">
            <MessageSquareIcon class="w-5 h-5 text-gray-600 mr-2" />
            <div>
              <div class="text-sm font-medium text-gray-900">Neutral</div>
              <div class="text-lg font-semibold text-gray-600">
                {(textAnalysis.sentiment_summary.neutral_rate || 0).toFixed(1)}%
              </div>
            </div>
          </div>
        </div>

        <div class="bg-red-50 border border-red-200 rounded-lg p-4">
          <div class="flex items-center">
            <ThumbsDownIcon class="w-5 h-5 text-red-600 mr-2" />
            <div>
              <div class="text-sm font-medium text-red-900">Negative</div>
              <div class="text-lg font-semibold text-red-600">
                {(textAnalysis.sentiment_summary.negative_rate || 0).toFixed(
                  1
                )}%
              </div>
            </div>
          </div>
        </div>
      </div>
    {/if}

    <div class="grid grid-cols-1 lg:grid-cols-2 gap-6">
      
      <div>
        <h4 class="font-medium text-gray-900 mb-3 flex items-center gap-2">
          <HashIcon class="w-4 h-4" />
          Most Frequent Keywords
        </h4>
        {#if topKeywords.length > 0}
          <div class="flex flex-wrap gap-2">
            {#each topKeywords as keyword}
              <div
                class="inline-flex items-center gap-1 px-3 py-1 bg-blue-100 text-blue-800 rounded-full text-sm">
                <span class="font-medium">{keyword.keyword}</span>
                <span class="text-xs text-blue-600">({keyword.count})</span>
              </div>
            {/each}
          </div>
        {:else}
          <div class="text-gray-500 text-sm">No keywords available</div>
        {/if}
      </div>

      
      <div>
        <h4 class="font-medium text-gray-900 mb-3 flex items-center gap-2">
          <ThumbsUpIcon class="w-4 h-4 text-green-600" />
          Positive Mentions
        </h4>
        {#if topPositive.length > 0}
          <div class="space-y-2">
            {#each topPositive as mention}
              <div
                class="flex items-center justify-between bg-green-50 border border-green-200 rounded-lg px-3 py-2">
                <span class="text-sm text-green-800 truncate"
                  >{mention.text}</span>
                <span class="text-xs text-green-600 ml-2"
                  >({mention.count})</span>
              </div>
            {/each}
          </div>
        {:else}
          <div class="text-gray-500 text-sm">
            No positive mentions available
          </div>
        {/if}
      </div>
    </div>

    
    {#if topNegative.length > 0}
      <div class="mt-6">
        <h4 class="font-medium text-gray-900 mb-3 flex items-center gap-2">
          <ThumbsDownIcon class="w-4 h-4 text-red-600" />
          Areas for Improvement
        </h4>
        <div class="space-y-2">
          {#each topNegative as mention}
            <div
              class="flex items-center justify-between bg-red-50 border border-red-200 rounded-lg px-3 py-2">
              <span class="text-sm text-red-800 truncate">{mention.text}</span>
              <span class="text-xs text-red-600 ml-2">({mention.count})</span>
            </div>
          {/each}
        </div>
      </div>
    {/if}

    
    {#if recentComments.length > 0}
      <div class="mt-6">
        <h4 class="font-medium text-gray-900 mb-3">Recent Feedback</h4>
        <div class="space-y-3">
          {#each recentComments as feedback}
            <div class="bg-gray-50 border border-gray-200 rounded-lg p-3">
              <div class="flex items-start justify-between">
                <div class="flex-1">
                  <p class="text-sm text-gray-800 leading-relaxed">
                    {truncateText(feedback.text)}
                  </p>
                  <div class="flex items-center gap-2 mt-2">
                    {#if feedback.product_name}
                      <span
                        class="text-xs text-gray-500 px-2 py-1 bg-white rounded">
                        {feedback.product_name}
                      </span>
                    {/if}
                    <span class="text-xs text-gray-500">
                      {formatDate(feedback.created_at)}
                    </span>
                  </div>
                </div>
                {#if feedback.rating}
                  <div class="ml-3 flex items-center gap-1">
                    <div class="flex items-center gap-1">
                      {#each Array(5) as _, i}
                        <div
                          class="w-3 h-3 {i < feedback.rating
                            ? 'text-yellow-400'
                            : 'text-gray-300'}">
                          ⭐
                        </div>
                      {/each}
                    </div>
                  </div>
                {/if}
              </div>
            </div>
          {/each}
        </div>
      </div>
    {/if}
  {:else}
    <div class="text-center py-8">
      <MessageSquareIcon class="w-12 h-12 text-gray-400 mx-auto mb-4" />
      <div class="text-gray-500 text-sm">
        No text response data available. Collect more feedback to see analysis.
      </div>
    </div>
  {/if}
</Card>
