<script lang="ts">
  import { Card } from '$lib/components/ui';
  import {
    SmileIcon,
    MehIcon,
    FrownIcon,
    TrendingUpIcon,
    MessageSquareIcon,
  } from 'lucide-svelte';
  import { onMount } from 'svelte';

  interface SentimentData {
    positive_rate: number;
    neutral_rate: number;
    negative_rate: number;
  }

  interface WordCloudItem {
    text: string;
    count: number;
    sentiment: 'positive' | 'negative' | 'neutral';
  }

  interface TextAnalysis {
    sentiment_summary: SentimentData;
    frequent_keywords?: Array<{ keyword: string; count: number }>;
    positive_mentions?: Array<{ text: string; count: number }>;
    negative_mentions?: Array<{ text: string; count: number }>;
  }

  let {
    textAnalysis = null,
    loading = false,
  }: {
    textAnalysis?: TextAnalysis | null;
    loading?: boolean;
  } = $props();

  let animationProgress = $state(0);
  let hoveredSentiment: string | null = null;

  const sentimentData = $derived(() => {
    if (!textAnalysis?.sentiment_summary) return null;

    const { positive_rate, neutral_rate, negative_rate } =
      textAnalysis.sentiment_summary;
    const total = positive_rate + neutral_rate + negative_rate;

    return [
      {
        type: 'positive',
        rate: positive_rate,
        percentage: total > 0 ? (positive_rate / total) * 100 : 0,
        color: '#10b981',
        bgColor: '#d1fae5',
        icon: SmileIcon,
        label: 'Positive',
      },
      {
        type: 'neutral',
        rate: neutral_rate,
        percentage: total > 0 ? (neutral_rate / total) * 100 : 0,
        color: '#6b7280',
        bgColor: '#f3f4f6',
        icon: MehIcon,
        label: 'Neutral',
      },
      {
        type: 'negative',
        rate: negative_rate,
        percentage: total > 0 ? (negative_rate / total) * 100 : 0,
        color: '#ef4444',
        bgColor: '#fee2e2',
        icon: FrownIcon,
        label: 'Negative',
      },
    ];
  });

  const wordCloud = $derived(() => {
    if (!textAnalysis) return [];

    const words: WordCloudItem[] = [];

    textAnalysis.positive_mentions?.forEach(mention => {
      words.push({
        text: mention.text,
        count: mention.count,
        sentiment: 'positive',
      });
    });

    textAnalysis.negative_mentions?.forEach(mention => {
      words.push({
        text: mention.text,
        count: mention.count,
        sentiment: 'negative',
      });
    });

    textAnalysis.frequent_keywords?.forEach(keyword => {
      if (
        !words.find(w => w.text.toLowerCase() === keyword.keyword.toLowerCase())
      ) {
        words.push({
          text: keyword.keyword,
          count: keyword.count,
          sentiment: 'neutral',
        });
      }
    });

    return words.sort((a, b) => b.count - a.count).slice(0, 20);
  });

  function getWordSize(count: number, maxCount: number): number {
    const minSize = 12;
    const maxSize = 32;
    return minSize + (count / maxCount) * (maxSize - minSize);
  }

  function getWordColor(sentiment: string): string {
    switch (sentiment) {
      case 'positive':
        return '#10b981';
      case 'negative':
        return '#ef4444';
      default:
        return '#6b7280';
    }
  }

  onMount(() => {
    const animate = () => {
      if (animationProgress < 1) {
        animationProgress = Math.min(animationProgress + 0.02, 1);
        requestAnimationFrame(animate);
      }
    };
    animate();
  });
</script>

<Card variant="default" class="sentiment-visualization">
  <div class="mb-6">
    <h3
      class="text-lg font-semibold bg-gradient-to-r from-gray-900 to-gray-700 bg-clip-text text-transparent">
      Sentiment Analysis
    </h3>
    <p class="text-sm text-gray-600 mt-1">
      Customer sentiment from text responses
    </p>
  </div>

  {#if loading}
    <div class="space-y-6">
      <div class="animate-pulse">
        <div class="h-32 bg-gray-100 rounded"></div>
      </div>
      <div class="animate-pulse">
        <div class="h-48 bg-gray-100 rounded"></div>
      </div>
    </div>
  {:else if sentimentData}
    
    <div class="grid grid-cols-1 md:grid-cols-3 gap-4 mb-8">
      {#each sentimentData as sentiment}
        {@const Icon = sentiment.icon}
        <div
          class="relative overflow-hidden rounded-lg border transition-all duration-300 cursor-pointer"
          style="background-color: {hoveredSentiment === sentiment.type
            ? sentiment.bgColor
            : '#ffffff'}; border-color: {hoveredSentiment === sentiment.type
            ? sentiment.color
            : '#e5e7eb'}"
          onmouseenter={() => (hoveredSentiment = sentiment.type)}
          onmouseleave={() => (hoveredSentiment = null)}>
          <div class="p-4">
            <div class="flex items-center justify-between mb-3">
              <div class="flex items-center gap-2">
                <Icon class="w-5 h-5" style="color: {sentiment.color}" />
                <span class="font-medium text-gray-900">{sentiment.label}</span>
              </div>
              <span class="text-2xl font-bold" style="color: {sentiment.color}">
                {(sentiment.percentage * animationProgress).toFixed(0)}%
              </span>
            </div>

            
            <div class="w-full bg-gray-200 rounded-full h-2 overflow-hidden">
              <div
                class="h-full rounded-full transition-all duration-1000 ease-out"
                style="width: {sentiment.percentage *
                  animationProgress}%; background-color: {sentiment.color}">
              </div>
            </div>
          </div>

          
          <div
            class="absolute inset-0 opacity-10 pointer-events-none"
            style="background: radial-gradient(circle at {hoveredSentiment ===
            sentiment.type
              ? '50%'
              : '150%'} 50%, {sentiment.color} 0%, transparent 70%); transition: all 0.5s ease-out">
          </div>
        </div>
      {/each}
    </div>

    
    {#if wordCloud.length > 0}
      <div class="mt-8">
        <h4 class="font-medium text-gray-900 mb-4 flex items-center gap-2">
          <MessageSquareIcon class="w-4 h-4" />
          Key Themes & Mentions
        </h4>
        <div class="relative min-h-[200px] p-4 bg-gray-50 rounded-lg">
          <div class="flex flex-wrap gap-3 justify-center items-center">
            {#each wordCloud as word, index}
              {@const maxCount = Math.max(...wordCloud.map(w => w.count))}
              {@const delay = index * 50}
              <span
                class="inline-block transition-all duration-300 hover:scale-110 cursor-pointer"
                style="
                  font-size: {getWordSize(word.count, maxCount)}px;
                  color: {getWordColor(word.sentiment)};
                  opacity: {animationProgress};
                  animation: fadeInUp 0.5s ease-out {delay}ms;
                  font-weight: {word.count > maxCount * 0.7 ? '600' : '400'};
                "
                title="{word.count} mentions">
                {word.text}
              </span>
            {/each}
          </div>
        </div>
      </div>
    {/if}

    
    {#if sentimentData && sentimentData.length >= 3}
      {@const score =
        (sentimentData[0].percentage - sentimentData[2].percentage + 100) / 2}
      <div
        class="mt-6 p-4 bg-gradient-to-r from-blue-50 to-purple-50 rounded-lg">
        <div class="flex items-center justify-between">
          <div>
            <h4 class="font-medium text-gray-900">Overall Sentiment Score</h4>
            <p class="text-sm text-gray-600 mt-1">
              Based on all text responses
            </p>
          </div>
          <div class="text-right">
            <div
              class="text-3xl font-bold bg-gradient-to-r from-blue-600 to-purple-600 bg-clip-text text-transparent">
              {(score * animationProgress).toFixed(0)}
            </div>
            <div class="text-sm text-gray-600">out of 100</div>
          </div>
        </div>
      </div>
    {/if}
  {:else}
    <div class="text-center py-12">
      <MessageSquareIcon class="w-12 h-12 text-gray-400 mx-auto mb-4" />
      <div class="text-gray-500 text-sm">
        No sentiment data available. Collect more text feedback to see analysis.
      </div>
    </div>
  {/if}
</Card>

<style>
  @keyframes fadeInUp {
    from {
      opacity: 0;
      transform: translateY(10px);
    }
    to {
      opacity: 1;
      transform: translateY(0);
    }
  }
</style>
