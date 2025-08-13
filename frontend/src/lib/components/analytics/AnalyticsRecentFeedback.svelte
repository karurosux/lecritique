<script lang="ts">
  import { Card } from '$lib/components/ui';
  import { Clock, MessageSquare } from 'lucide-svelte';

  interface Feedback {
    id?: string;
    rating?: number;
    overall_rating?: number;
    comment?: string;
    product_name?: string;
    created_at: string;
  }

  interface AnalyticsData {
    recent_feedback: Feedback[];
  }

  let {
    analyticsData = null,
    feedbacks = [],
    loading = false,
  }: {
    analyticsData?: AnalyticsData | null;
    feedbacks?: Feedback[];
    loading?: boolean;
  } = $props();

  let displayFeedbacks = $derived(() => {
    if (analyticsData?.recent_feedback) {
      return analyticsData.recent_feedback;
    }
    return feedbacks;
  });

  function renderStars(rating: number): string {
    return '★'.repeat(Math.round(rating)) + '☆'.repeat(5 - Math.round(rating));
  }

  function getRatingColor(rating: number): string {
    if (rating >= 4) return 'text-green-600';
    if (rating >= 3) return 'text-yellow-600';
    return 'text-red-600';
  }

  function getBorderColor(rating: number): string {
    if (rating >= 4) return 'border-green-500';
    if (rating >= 3) return 'border-yellow-500';
    return 'border-red-500';
  }

  function formatDate(dateString: string): string {
    try {
      const date = new Date(dateString);
      const now = new Date();
      const diffTime = Math.abs(now.getTime() - date.getTime());
      const diffDays = Math.floor(diffTime / (1000 * 60 * 60 * 24));

      if (diffDays === 0) {
        const diffHours = Math.floor(diffTime / (1000 * 60 * 60));
        if (diffHours === 0) {
          const diffMinutes = Math.floor(diffTime / (1000 * 60));
          return `${diffMinutes} minutes ago`;
        }
        return `${diffHours} hours ago`;
      } else if (diffDays === 1) {
        return 'Yesterday';
      } else if (diffDays < 7) {
        return `${diffDays} days ago`;
      } else {
        return date.toLocaleDateString();
      }
    } catch {
      return dateString;
    }
  }
</script>

<div class="grid grid-cols-1 lg:grid-cols-2 gap-8 mb-8">
  
  <Card
    variant="glass"
    class="lg:col-span-2 hover:shadow-lg transition-all duration-300">
    <div class="mb-6">
      <h3
        class="text-lg font-semibold bg-gradient-to-r from-gray-900 to-gray-700 bg-clip-text text-transparent">
        Recent Feedback Activity
      </h3>
      <p class="text-sm text-gray-600 mt-1">
        Latest customer feedback and comments
      </p>
    </div>

    <div class="space-y-4 max-h-96 overflow-y-auto pr-2 custom-scrollbar">
      {#if loading}
        {#each Array(5) as _}
          <div class="animate-pulse">
            <div class="flex items-start gap-4 p-4 border-l-4 border-gray-200">
              <div class="flex-1 space-y-2">
                <div class="flex items-center gap-2">
                  <div class="h-4 bg-gray-200 rounded w-24"></div>
                  <div class="h-4 bg-gray-200 rounded w-32"></div>
                </div>
                <div class="h-3 bg-gray-200 rounded w-3/4"></div>
                <div class="h-3 bg-gray-200 rounded w-20"></div>
              </div>
            </div>
          </div>
        {/each}
      {:else if displayFeedbacks() && displayFeedbacks().length > 0}
        {#each displayFeedbacks().slice(0, 10) as feedback}
          {@const rating = feedback.rating || feedback.overall_rating || 0}
          <div
            class="group border-l-4 {getBorderColor(
              rating
            )} pl-4 py-3 hover:bg-gray-50 rounded-r-lg transition-all duration-200">
            <div class="flex items-start justify-between">
              <div class="flex-1">
                <div class="flex items-center gap-3 mb-1">
                  <span class="text-sm font-medium {getRatingColor(rating)}">
                    {renderStars(rating)}
                  </span>
                  <span
                    class="text-xs text-gray-500 bg-gray-100 px-2 py-0.5 rounded-full">
                    {rating}/5
                  </span>
                  {#if feedback.product_name}
                    <span class="text-xs text-gray-400">•</span>
                    <span
                      class="text-xs font-medium text-gray-700 bg-purple-50 px-2 py-0.5 rounded-full">
                      {feedback.product_name}
                    </span>
                  {/if}
                </div>

                {#if feedback.comment}
                  <p
                    class="text-gray-700 text-sm mb-2 line-clamp-2 group-hover:line-clamp-none transition-all duration-200">
                    "{feedback.comment}"
                  </p>
                {/if}

                <p class="text-xs text-gray-500">
                  <Clock class="inline-block h-3 w-3 mr-1" />
                  {formatDate(feedback.created_at)}
                </p>
              </div>
            </div>
          </div>
        {/each}
      {:else}
        <div class="text-center py-12">
          <div
            class="w-16 h-16 bg-gray-100 rounded-full flex items-center justify-center mx-auto mb-4">
            <MessageSquare class="h-8 w-8 text-gray-400" />
          </div>
          <p class="text-gray-500">No recent feedback available</p>
        </div>
      {/if}
    </div>
  </Card>
</div>

<style>
  .custom-scrollbar::-webkit-scrollbar {
    width: 6px;
  }

  .custom-scrollbar::-webkit-scrollbar-track {
    background: #f1f1f1;
    border-radius: 3px;
  }

  .custom-scrollbar::-webkit-scrollbar-thumb {
    background: #888;
    border-radius: 3px;
  }

  .custom-scrollbar::-webkit-scrollbar-thumb:hover {
    background: #555;
  }

  .line-clamp-2 {
    display: -webkit-box;
    -webkit-line-clamp: 2;
    -webkit-box-orient: vertical;
    overflow: hidden;
  }

  .line-clamp-none {
    display: block;
    -webkit-line-clamp: unset;
    -webkit-box-orient: unset;
    overflow: visible;
  }
</style>
