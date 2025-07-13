import { writable, derived, get } from 'svelte/store';
import type { Writable, Readable } from 'svelte/store';

// Types for analytics state
export interface AnalyticsFilters {
  restaurantId: string;
  dishId: string | null;
  timeframe: '24h' | '7d' | '30d' | '90d' | 'custom';
  startDate?: Date;
  endDate?: Date;
  segments: {
    timeOfDay?: 'breakfast' | 'lunch' | 'dinner' | 'all';
    dayType?: 'weekday' | 'weekend' | 'all';
    customerType?: string;
  };
}

export interface AnalyticsSelection {
  highlightedQuestionId: string | null;
  highlightedMetric: string | null;
  comparisonMode: 'none' | 'period' | 'dish' | 'question';
  comparisonTarget?: string;
  hoveredDataPoint?: {
    type: string;
    id: string;
    value: any;
  };
}

export interface AnalyticsData {
  raw: any;
  aggregated: any;
  computed: any;
  loading: boolean;
  error: string | null;
}

export interface AnalyticsState {
  filters: AnalyticsFilters;
  selection: AnalyticsSelection;
  data: AnalyticsData;
}

// Create the main analytics store
function createAnalyticsStore() {
  // Individual writable stores for each part of the state
  const filters: Writable<AnalyticsFilters> = writable({
    restaurantId: '',
    dishId: null,
    timeframe: '7d',
    segments: {
      timeOfDay: 'all',
      dayType: 'all'
    }
  });

  const selection: Writable<AnalyticsSelection> = writable({
    highlightedQuestionId: null,
    highlightedMetric: null,
    comparisonMode: 'none',
    hoveredDataPoint: undefined
  });

  const data: Writable<AnalyticsData> = writable({
    raw: null,
    aggregated: null,
    computed: null,
    loading: false,
    error: null
  });

  // Derived store for filtered data based on current filters
  const filteredData: Readable<any> = derived(
    [data, filters],
    ([$data, $filters]) => {
      if (!$data.raw) return null;

      let filtered = $data.raw;

      // Apply time-based filtering
      if ($filters.timeframe !== 'all' && filtered.feedback) {
        const now = new Date();
        const cutoffDate = new Date();

        switch ($filters.timeframe) {
          case '24h':
            cutoffDate.setDate(now.getDate() - 1);
            break;
          case '7d':
            cutoffDate.setDate(now.getDate() - 7);
            break;
          case '30d':
            cutoffDate.setDate(now.getDate() - 30);
            break;
          case '90d':
            cutoffDate.setDate(now.getDate() - 90);
            break;
        }

        filtered = {
          ...filtered,
          feedback: filtered.feedback.filter((f: any) => 
            new Date(f.created_at) >= cutoffDate
          )
        };
      }

      // Apply segment filtering
      if ($filters.segments.timeOfDay !== 'all' && filtered.feedback) {
        filtered = {
          ...filtered,
          feedback: filtered.feedback.filter((f: any) => {
            const hour = new Date(f.created_at).getHours();
            switch ($filters.segments.timeOfDay) {
              case 'breakfast':
                return hour >= 6 && hour < 11;
              case 'lunch':
                return hour >= 11 && hour < 15;
              case 'dinner':
                return hour >= 17 && hour < 22;
              default:
                return true;
            }
          })
        };
      }

      if ($filters.segments.dayType !== 'all' && filtered.feedback) {
        filtered = {
          ...filtered,
          feedback: filtered.feedback.filter((f: any) => {
            const day = new Date(f.created_at).getDay();
            const isWeekend = day === 0 || day === 6;
            return $filters.segments.dayType === 'weekend' ? isWeekend : !isWeekend;
          })
        };
      }

      return filtered;
    }
  );

  // Derived store for comparison data
  const comparisonData: Readable<any> = derived(
    [data, selection, filters],
    ([$data, $selection, $filters]) => {
      if ($selection.comparisonMode === 'none' || !$data.raw) return null;

      switch ($selection.comparisonMode) {
        case 'period':
          // Compare current period with previous period
          return calculatePeriodComparison($data.raw, $filters.timeframe);
        
        case 'dish':
          // Compare selected dish with another dish
          if ($selection.comparisonTarget) {
            return {
              current: $data.raw,
              comparison: null // Would need to fetch comparison dish data
            };
          }
          break;
        
        case 'question':
          // Compare performance across questions
          return calculateQuestionComparison($data.raw);
      }

      return null;
    }
  );

  // Derived store for computed metrics
  const computedMetrics: Readable<any> = derived(
    [filteredData],
    ([$filteredData]) => {
      if (!$filteredData) return null;

      return {
        satisfactionIndex: calculateSatisfactionIndex($filteredData),
        improvementRate: calculateImprovementRate($filteredData),
        responseRate: calculateResponseRate($filteredData),
        sentimentScore: calculateSentimentScore($filteredData)
      };
    }
  );

  // Action functions
  function updateFilters(newFilters: Partial<AnalyticsFilters>) {
    filters.update(f => ({ ...f, ...newFilters }));
  }

  function updateSelection(newSelection: Partial<AnalyticsSelection>) {
    selection.update(s => ({ ...s, ...newSelection }));
  }

  function setHighlightedQuestion(questionId: string | null) {
    selection.update(s => ({ ...s, highlightedQuestionId: questionId }));
  }

  function setComparisonMode(mode: AnalyticsSelection['comparisonMode'], target?: string) {
    selection.update(s => ({ 
      ...s, 
      comparisonMode: mode,
      comparisonTarget: target 
    }));
  }

  function setHoveredDataPoint(dataPoint: AnalyticsSelection['hoveredDataPoint']) {
    selection.update(s => ({ ...s, hoveredDataPoint: dataPoint }));
  }

  function setData(newData: Partial<AnalyticsData>) {
    data.update(d => ({ ...d, ...newData }));
  }

  function reset() {
    filters.set({
      restaurantId: '',
      dishId: null,
      timeframe: '7d',
      segments: {
        timeOfDay: 'all',
        dayType: 'all'
      }
    });
    selection.set({
      highlightedQuestionId: null,
      highlightedMetric: null,
      comparisonMode: 'none',
      hoveredDataPoint: undefined
    });
    data.set({
      raw: null,
      aggregated: null,
      computed: null,
      loading: false,
      error: null
    });
  }

  return {
    // Stores
    filters,
    selection,
    data,
    filteredData,
    comparisonData,
    computedMetrics,
    
    // Actions
    updateFilters,
    updateSelection,
    setHighlightedQuestion,
    setComparisonMode,
    setHoveredDataPoint,
    setData,
    reset,
    
    // Utility
    subscribe: derived(
      [filters, selection, data, filteredData, comparisonData, computedMetrics],
      ([$filters, $selection, $data, $filteredData, $comparisonData, $computedMetrics]) => ({
        filters: $filters,
        selection: $selection,
        data: $data,
        filteredData: $filteredData,
        comparisonData: $comparisonData,
        computedMetrics: $computedMetrics
      })
    ).subscribe
  };
}

// Calculation helper functions
function calculatePeriodComparison(data: any, timeframe: string): any {
  const now = new Date();
  const currentPeriodStart = new Date();
  const previousPeriodStart = new Date();
  const previousPeriodEnd = new Date();

  switch (timeframe) {
    case '24h':
      currentPeriodStart.setDate(now.getDate() - 1);
      previousPeriodStart.setDate(now.getDate() - 2);
      previousPeriodEnd.setDate(now.getDate() - 1);
      break;
    case '7d':
      currentPeriodStart.setDate(now.getDate() - 7);
      previousPeriodStart.setDate(now.getDate() - 14);
      previousPeriodEnd.setDate(now.getDate() - 7);
      break;
    case '30d':
      currentPeriodStart.setDate(now.getDate() - 30);
      previousPeriodStart.setDate(now.getDate() - 60);
      previousPeriodEnd.setDate(now.getDate() - 30);
      break;
    default:
      return null;
  }

  // Filter data for each period
  const currentPeriodData = data.feedback?.filter((f: any) => 
    new Date(f.created_at) >= currentPeriodStart
  ) || [];

  const previousPeriodData = data.feedback?.filter((f: any) => 
    new Date(f.created_at) >= previousPeriodStart && 
    new Date(f.created_at) < previousPeriodEnd
  ) || [];

  return {
    current: {
      feedback: currentPeriodData,
      count: currentPeriodData.length,
      averageRating: calculateAverageRating(currentPeriodData)
    },
    previous: {
      feedback: previousPeriodData,
      count: previousPeriodData.length,
      averageRating: calculateAverageRating(previousPeriodData)
    },
    change: {
      count: currentPeriodData.length - previousPeriodData.length,
      percentage: previousPeriodData.length > 0 
        ? ((currentPeriodData.length - previousPeriodData.length) / previousPeriodData.length) * 100
        : 0
    }
  };
}

function calculateQuestionComparison(data: any): any {
  if (!data.question_scores) return null;

  const questions = data.question_scores.map((q: any) => ({
    id: q.question_id,
    text: q.question_text,
    averageScore: q.average_score,
    responseCount: q.response_count,
    variance: calculateVariance(q.scores || [])
  }));

  return {
    questions,
    bestPerforming: questions.reduce((best: any, current: any) => 
      current.averageScore > (best?.averageScore || 0) ? current : best
    , null),
    worstPerforming: questions.reduce((worst: any, current: any) => 
      current.averageScore < (worst?.averageScore || 5) ? current : worst
    , null)
  };
}

function calculateSatisfactionIndex(data: any): number {
  if (!data.feedback || data.feedback.length === 0) return 0;
  
  const scores = data.feedback.map((f: any) => f.rating || 0);
  const average = scores.reduce((sum: number, score: number) => sum + score, 0) / scores.length;
  
  // Normalize to 0-100 scale
  return (average / 5) * 100;
}

function calculateImprovementRate(data: any): number {
  if (!data.trends?.daily_feedback || data.trends.daily_feedback.length < 2) return 0;
  
  const recentDays = data.trends.daily_feedback.slice(-7);
  if (recentDays.length < 2) return 0;
  
  const firstDay = recentDays[0].average_rating;
  const lastDay = recentDays[recentDays.length - 1].average_rating;
  
  return ((lastDay - firstDay) / firstDay) * 100;
}

function calculateResponseRate(data: any): number {
  // This would need actual visitor data to be accurate
  // For now, return a placeholder
  return data.feedback?.length || 0;
}

function calculateSentimentScore(data: any): number {
  if (!data.sentiment_summary) return 50;
  
  const { positive_rate = 0, negative_rate = 0 } = data.sentiment_summary;
  
  // Calculate sentiment score (0-100, where 100 is most positive)
  return ((positive_rate - negative_rate) + 100) / 2;
}

function calculateAverageRating(feedback: any[]): number {
  if (!feedback || feedback.length === 0) return 0;
  
  const sum = feedback.reduce((total, f) => total + (f.rating || 0), 0);
  return sum / feedback.length;
}

function calculateVariance(scores: number[]): number {
  if (scores.length === 0) return 0;
  
  const mean = scores.reduce((sum, score) => sum + score, 0) / scores.length;
  const squaredDifferences = scores.map(score => Math.pow(score - mean, 2));
  
  return squaredDifferences.reduce((sum, diff) => sum + diff, 0) / scores.length;
}

// Create and export the singleton store
export const analyticsStore = createAnalyticsStore();

// Export convenient helper hooks
export function useAnalytics() {
  return analyticsStore;
}

// Context key for providing the store through component tree
export const ANALYTICS_CONTEXT_KEY = Symbol('analytics');