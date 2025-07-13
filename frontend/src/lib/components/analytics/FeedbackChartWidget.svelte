<script lang="ts">
  import { Card } from '$lib/components/ui';
  import { onMount } from 'svelte';
  import { Star, BarChart3, TrendingUp, Utensils, ToggleLeft, ToggleRight } from 'lucide-svelte';

  interface FeedbackRecord {
    id: string;
    overall_rating?: number;
    rating?: number;
    dish?: { name: string };
    dish_name?: string;
    responses?: Array<{
      question_id: string;
      question_text: string;
      answer: any;
    }>;
    comment?: string;
    created_at: string;
  }

  interface ChartData {
    questionId: string;
    questionText: string;
    type: 'rating' | 'choice' | 'text' | 'boolean';
    data: Array<{ label: string; value: number; color: string }>;
    combinationData?: Array<{ label: string; value: number; color: string }>;
    totalResponses: number;
    isMultiChoice: boolean;
  }

  let {
    feedbacks = [],
    title = "Feedback Analytics"
  }: {
    feedbacks: FeedbackRecord[];
    title?: string;
  } = $props();

  let dishesByName = $derived(groupFeedbacksByDish(feedbacks));
  let chartViewModes = $state(new Map<string, 'individual' | 'combinations'>());

  function groupFeedbacksByDish(feedbacks: FeedbackRecord[]) {
    if (!feedbacks || feedbacks.length === 0) return new Map();

    const dishGroups = new Map<string, {
      name: string;
      feedbacks: FeedbackRecord[];
      charts: ChartData[];
    }>();

    // Group feedbacks by dish
    feedbacks.forEach(feedback => {
      const dishName = feedback.dish_name || feedback.dish?.name || `Dish ${feedback.dish_id || 'Unknown'}`;
      
      if (!dishGroups.has(dishName)) {
        dishGroups.set(dishName, {
          name: dishName,
          feedbacks: [],
          charts: []
        });
      }
      
      dishGroups.get(dishName)!.feedbacks.push(feedback);
    });

    // Generate charts for each dish
    dishGroups.forEach(dishGroup => {
      dishGroup.charts = generateChartsForDish(dishGroup.feedbacks, dishGroup.name);
    });

    return dishGroups;
  }

  function isTextQuestion(answers: any[]): boolean {
    // If 80% or more of answers are detected as text, treat the whole question as text
    const textAnswers = answers.filter(answer => detectQuestionType(answer) === 'text');
    const textRatio = textAnswers.length / answers.length;
    
    if (textRatio >= 0.8) return true;
    
    // If answers are very diverse (many unique values), likely text
    const uniqueAnswers = new Set(answers.map(a => String(a).toLowerCase().trim()));
    const uniqueRatio = uniqueAnswers.size / answers.length;
    
    // If 90% of answers are unique, it's likely open text
    if (uniqueRatio >= 0.9 && answers.length > 3) return true;
    
    return false;
  }

  function generateChartsForDish(dishFeedbacks: FeedbackRecord[], dishName: string): ChartData[] {
    if (!dishFeedbacks || dishFeedbacks.length === 0) return [];

    // Collect all unique questions for this dish
    const questionMap = new Map<string, { text: string; answers: any[]; type: string }>();

    // Add overall rating as a special question
    const overallRatings = dishFeedbacks
      .map(f => f.overall_rating || f.rating)
      .filter(r => r && r > 0);
    
    if (overallRatings.length > 0) {
      questionMap.set('overall', {
        text: 'Overall Rating',
        answers: overallRatings,
        type: 'rating'
      });
    }

    // Process question responses for this dish
    dishFeedbacks.forEach(feedback => {
      if (feedback.responses) {
        feedback.responses.forEach(response => {
          const key = response.question_id || response.question_text;
          if (!questionMap.has(key)) {
            // Use actual question type if available, otherwise detect from answer
            const questionType = response.question_type || detectQuestionType(response.answer);
            questionMap.set(key, {
              text: response.question_text,
              answers: [],
              type: questionType
            });
          }
          questionMap.get(key)!.answers.push(response.answer);
        });
      }
    });

    // Generate chart data for each question (exclude text questions)
    return Array.from(questionMap.entries())
      .filter(([questionId, question]) => {
        // Filter out text questions and questions that look like open text
        return question.type !== 'text' && !isTextQuestion(question.answers);
      })
      .map(([questionId, question]) => {
        const { chartData, combinationData, isMultiChoice } = generateChartDataForQuestion(question.answers, question.type);
        return {
          questionId,
          questionText: question.text,
          type: question.type as 'rating' | 'choice',
          data: chartData,
          combinationData,
          totalResponses: question.answers.length,
          isMultiChoice
        };
      });
  }

  function detectQuestionType(answer: any): string {
    // Check for numeric ratings (1-5 scale or 1-10 scale)
    if (typeof answer === 'number') {
      if ((answer >= 1 && answer <= 5) || (answer >= 1 && answer <= 10)) {
        return 'rating';
      }
    }
    
    // Check for string representations of ratings
    if (typeof answer === 'string') {
      const trimmed = answer.trim();
      
      // Check for numeric strings that represent ratings
      const numValue = parseInt(trimmed);
      if (!isNaN(numValue) && ((numValue >= 1 && numValue <= 5) || (numValue >= 1 && numValue <= 10))) {
        return 'rating';
      }
      
      // Check for star ratings (★★★★☆, etc.)
      if (trimmed.includes('★') || trimmed.includes('⭐')) {
        return 'rating';
      }
      
      // Check for text responses (longer than 30 characters, contains sentences, or common text patterns)
      if (trimmed.length > 30 || 
          trimmed.includes('.') ||
          trimmed.includes('!') ||
          trimmed.includes('?') ||
          trimmed.split(' ').length > 5 ||  // More than 5 words
          trimmed.toLowerCase().includes('comment') ||
          trimmed.toLowerCase().includes('feedback') ||
          trimmed.toLowerCase().includes('suggestion') ||
          trimmed.toLowerCase().includes('recommend') ||
          trimmed.toLowerCase().includes('improve') ||
          trimmed.toLowerCase().includes('delicious') ||
          trimmed.toLowerCase().includes('terrible') ||
          trimmed.toLowerCase().includes('amazing') ||
          trimmed.toLowerCase().includes('love') ||
          trimmed.toLowerCase().includes('hate')) {
        return 'text';
      }
    }
    
    // Everything else (including yes/no, booleans, multiple choice) goes to 'choice'
    return 'choice';
  }

  function generateChartDataForQuestion(answers: any[], type: string): {
    chartData: Array<{ label: string; value: number; color: string }>;
    combinationData?: Array<{ label: string; value: number; color: string }>;
    isMultiChoice: boolean;
  } {
    if (type === 'rating') {
      return {
        chartData: generateRatingChart(answers),
        isMultiChoice: false
      };
    } else {
      const { individualData, combinationData, isMultiChoice } = generateChoiceChart(answers);
      return {
        chartData: individualData,
        combinationData,
        isMultiChoice
      };
    }
  }

  function generateRatingChart(ratings: number[]) {
    const counts = { '1': 0, '2': 0, '3': 0, '4': 0, '5': 0 };
    ratings.forEach(rating => {
      if (rating >= 1 && rating <= 5) {
        counts[rating.toString()]++;
      }
    });

    const colors = {
      '1': '#ef4444', // red
      '2': '#f97316', // orange  
      '3': '#eab308', // yellow
      '4': '#22c55e', // green
      '5': '#16a34a'  // dark green
    };

    return Object.entries(counts).map(([rating, count]) => ({
      label: `${rating} Star${rating === '1' ? '' : 's'}`,
      value: count,
      color: colors[rating as keyof typeof colors]
    }));
  }

  function generateChoiceChart(choices: any[]): {
    individualData: Array<{ label: string; value: number; color: string }>;
    combinationData: Array<{ label: string; value: number; color: string }>;
    isMultiChoice: boolean;
  } {
    const individualCounts = new Map<string, number>();
    const combinationCounts = new Map<string, number>();
    let hasMultipleSelections = false;
    
    choices.forEach(choice => {
      // Handle different types of multi-choice data
      let options: string[] = [];
      
      if (Array.isArray(choice)) {
        // If choice is already an array, use it directly
        options = choice.map(c => String(c).trim()).filter(c => c);
      } else if (typeof choice === 'string') {
        // Handle comma-separated values or single values
        if (choice.includes(',')) {
          options = choice.split(',').map(c => c.trim()).filter(c => c);
        } else if (choice.includes(';')) {
          options = choice.split(';').map(c => c.trim()).filter(c => c);
        } else if (choice.includes('|')) {
          options = choice.split('|').map(c => c.trim()).filter(c => c);
        } else {
          options = [choice.trim()].filter(c => c);
        }
      } else {
        // Convert other types to string
        const choiceStr = String(choice).trim();
        if (choiceStr) {
          options = [choiceStr];
        }
      }
      
      // Track if this is a multi-choice question
      if (options.length > 1) {
        hasMultipleSelections = true;
      }
      
      // Count each individual option
      options.forEach(option => {
        if (option) {
          individualCounts.set(option, (individualCounts.get(option) || 0) + 1);
        }
      });
      
      // Count combination patterns (only if multiple options selected)
      if (options.length > 0) {
        const combination = options.sort().join(' + ');
        combinationCounts.set(combination, (combinationCounts.get(combination) || 0) + 1);
      }
    });

    const colors = ['#3b82f6', '#8b5cf6', '#06b6d4', '#10b981', '#f59e0b', '#ef4444', '#ec4899', '#8b5a2b'];
    
    const individualData = Array.from(individualCounts.entries())
      .sort((a, b) => b[1] - a[1]) // Sort by count descending
      .map(([choice, count], index) => ({
        label: choice,
        value: count,
        color: colors[index % colors.length]
      }));
    
    const combinationData = Array.from(combinationCounts.entries())
      .sort((a, b) => b[1] - a[1]) // Sort by count descending
      .map(([combination, count], index) => ({
        label: combination,
        value: count,
        color: colors[index % colors.length]
      }));

    return {
      individualData,
      combinationData,
      isMultiChoice: hasMultipleSelections
    };
  }

  function generateTextChart(texts: string[]) {
    // For text responses, we'll create a simple sentiment analysis
    const sentiments = { positive: 0, neutral: 0, negative: 0 };
    
    const positiveWords = ['good', 'great', 'excellent', 'amazing', 'love', 'perfect', 'delicious', 'fresh'];
    const negativeWords = ['bad', 'terrible', 'awful', 'hate', 'cold', 'slow', 'expensive', 'bland'];

    texts.forEach(text => {
      if (!text || typeof text !== 'string') return;
      
      const lowerText = text.toLowerCase();
      const hasPositive = positiveWords.some(word => lowerText.includes(word));
      const hasNegative = negativeWords.some(word => lowerText.includes(word));
      
      if (hasPositive && !hasNegative) {
        sentiments.positive++;
      } else if (hasNegative && !hasPositive) {
        sentiments.negative++;
      } else {
        sentiments.neutral++;
      }
    });

    return [
      { label: 'Positive', value: sentiments.positive, color: '#22c55e' },
      { label: 'Neutral', value: sentiments.neutral, color: '#6b7280' },
      { label: 'Negative', value: sentiments.negative, color: '#ef4444' }
    ].filter(item => item.value > 0);
  }

  function getChartTypeIcon(type: string) {
    switch (type) {
      case 'rating': return Star;
      case 'choice': return BarChart3;
      default: return TrendingUp;
    }
  }

  function toggleChartView(questionId: string, currentMode: 'individual' | 'combinations') {
    const newMode = currentMode === 'individual' ? 'combinations' : 'individual';
    chartViewModes.set(questionId, newMode);
    chartViewModes = new Map(chartViewModes); // Trigger reactivity
  }

  function getChartViewMode(questionId: string): 'individual' | 'combinations' {
    return chartViewModes.get(questionId) || 'individual';
  }

  function getCurrentChartData(chart: ChartData, mode: 'individual' | 'combinations') {
    if (mode === 'combinations' && chart.combinationData) {
      return chart.combinationData;
    }
    return chart.data;
  }

  function getQuestionTypeLabel(chart: ChartData): string {
    // Map backend question types to display labels
    switch (chart.type) {
      case 'rating':
        return 'Rating';
      case 'scale':
        return 'Scale';
      case 'multi_choice':
        return 'Multi-Choice';
      case 'single_choice':
        return 'Choice';
      case 'yes_no':
        return 'Yes/No';
      case 'text':
        return 'Text';
      default:
        // Fallback for legacy data or unknown types
        if (chart.type === 'choice') {
          // Check if it's a yes/no question
          const hasYesNo = chart.data.some(item => 
            item.label.toLowerCase() === 'yes' || 
            item.label.toLowerCase() === 'no' ||
            item.label.toLowerCase() === 'true' ||
            item.label.toLowerCase() === 'false'
          );
          
          if (hasYesNo) {
            return 'Yes/No';
          }
          
          // Check if it's likely a scale question (numbers)
          const hasNumbers = chart.data.some(item => !isNaN(parseInt(item.label)));
          if (hasNumbers && chart.data.length <= 10) {
            return 'Scale';
          }
          
          // Check if it's multiple choice
          if (chart.isMultiChoice) {
            return 'Multi-Choice';
          }
          
          return 'Choice';
        }
        
        return chart.type || 'Unknown';
    }
  }
</script>

<div class="feedback-chart-widget">
  {#if title}
    <div class="mb-8">
      <h2 class="text-2xl font-bold bg-gradient-to-r from-gray-900 to-gray-700 bg-clip-text text-transparent">{title}</h2>
      <p class="text-gray-600 mt-1 font-medium">Visual analysis of {feedbacks.length} feedback responses</p>
    </div>
  {/if}

  {#if dishesByName.size === 0}
    <Card variant="elevated">
      <div class="text-center py-16">
        <div class="h-16 w-16 bg-gray-100 rounded-2xl flex items-center justify-center mx-auto mb-4">
          <BarChart3 class="h-8 w-8 text-gray-400" />
        </div>
        <h3 class="text-lg font-semibold text-gray-900 mb-2">No Question Data</h3>
        <p class="text-gray-500">No feedback responses with questions available for analysis</p>
      </div>
    </Card>
  {:else}
    <div class="space-y-10">
      {#each Array.from(dishesByName.entries()) as [dishName, dishData]}
        {#if dishData.charts.length > 0}
          <!-- Dish Section Header -->
          <div class="dish-section">
            <div class="mb-6">
              <div class="flex items-center gap-3 mb-2">
                <div class="h-12 w-12 bg-gradient-to-br from-purple-500 to-pink-600 rounded-2xl flex items-center justify-center shadow-lg shadow-purple-500/25">
                  <Utensils class="h-6 w-6 text-white" />
                </div>
                <div>
                  <h3 class="text-xl font-bold bg-gradient-to-r from-purple-700 to-pink-600 bg-clip-text text-transparent">
                    {dishName}
                  </h3>
                  <p class="text-sm text-gray-600 font-medium">
                    {dishData.feedbacks.length} responses • {dishData.charts.length} questions analyzed
                  </p>
                </div>
              </div>
              <div class="w-full h-px bg-gradient-to-r from-purple-200 via-purple-300 to-transparent"></div>
            </div>

            <!-- Charts for this dish -->
            <div class="grid grid-cols-1 md:grid-cols-2 xl:grid-cols-3 gap-6">
              {#each dishData.charts as chart, index}
                {@const currentMode = getChartViewMode(chart.questionId)}
                {@const currentData = getCurrentChartData(chart, currentMode)}
                <Card variant="elevated" hover interactive class="group overflow-hidden">
                  <!-- Card Header -->
                  <div class="p-6 pb-4">
                    <div class="flex items-start gap-3 mb-3">
                      <div class="h-10 w-10 bg-gradient-to-br from-blue-500 to-purple-600 rounded-xl flex items-center justify-center shadow-lg shadow-blue-500/25 group-hover:scale-110 transition-transform duration-300">
                        <svelte:component this={getChartTypeIcon(chart.type)} class="h-5 w-5 text-white" />
                      </div>
                      <div class="flex-1 min-w-0">
                        <h4 class="font-semibold text-gray-900 text-base leading-tight mb-1">{chart.questionText}</h4>
                        <div class="flex items-center gap-2">
                          <span class="inline-flex items-center px-2 py-1 rounded-full text-xs font-medium bg-blue-50 text-blue-700">
                            {chart.totalResponses} responses
                          </span>
                          <span class="inline-flex items-center px-2 py-1 rounded-full text-xs font-medium bg-gray-100 text-gray-600">
                            {getQuestionTypeLabel(chart)}
                          </span>
                        </div>
                      </div>
                    </div>
                    
                    <!-- Toggle for multi-choice questions -->
                    {#if chart.isMultiChoice && chart.combinationData}
                      <div class="flex items-center justify-between mt-3 pt-3 border-t border-gray-100">
                        <span class="text-xs font-medium text-gray-600">
                          {currentMode === 'individual' ? 'Individual Options' : 'Combination Patterns'}
                        </span>
                        <button
                          onclick={() => toggleChartView(chart.questionId, currentMode)}
                          class="flex items-center gap-2 px-3 py-1.5 text-xs font-medium text-purple-700 bg-purple-50 hover:bg-purple-100 rounded-lg transition-colors duration-200"
                        >
                          <svelte:component this={currentMode === 'individual' ? ToggleLeft : ToggleRight} class="h-3 w-3" />
                          {currentMode === 'individual' ? 'Show Combinations' : 'Show Individual'}
                        </button>
                      </div>
                    {/if}
                  </div>

                  <!-- Chart Visualization -->
                  <div class="px-6 pb-4">
                    <div class="space-y-4">
                      {#each currentData as item, itemIndex}
                        {@const percentage = chart.totalResponses > 0 ? (item.value / chart.totalResponses) * 100 : 0}
                        <div class="group/item">
                          <div class="flex items-center justify-between text-sm mb-2">
                            <div class="flex items-center gap-2">
                              <div 
                                class="w-3 h-3 rounded-full shadow-sm" 
                                style="background-color: {item.color}"
                              ></div>
                              <span class="text-gray-700 font-medium text-xs leading-tight">{item.label}</span>
                            </div>
                            <div class="flex items-center gap-2">
                              <span class="text-gray-900 font-bold">{item.value}</span>
                              <span class="text-xs text-gray-500 font-medium">{percentage.toFixed(1)}%</span>
                            </div>
                          </div>
                          <div class="relative">
                            <div class="w-full bg-gray-100 rounded-full h-3 overflow-hidden shadow-inner">
                              <div 
                                class="h-full rounded-full transition-all duration-700 ease-out shadow-sm group-hover/item:shadow-md"
                                style="width: {percentage}%; background: linear-gradient(90deg, {item.color}, {item.color}dd)"
                              ></div>
                            </div>
                          </div>
                        </div>
                      {/each}
                    </div>
                  </div>

                  <!-- Summary Stats -->
                  {#if chart.type === 'rating'}
                    {@const ratings = chart.data.map((d, i) => ({ rating: i + 1, count: d.value })).filter(r => r.count > 0)}
                    {@const avgRating = ratings.length > 0 
                      ? ratings.reduce((sum, r) => sum + (r.rating * r.count), 0) / chart.totalResponses 
                      : 0}
                    <div class="bg-gradient-to-r from-yellow-50 to-orange-50 mx-6 mb-6 rounded-xl p-4 border border-yellow-200/50">
                      <div class="flex items-center justify-between">
                        <div>
                          <div class="text-2xl font-bold bg-gradient-to-r from-yellow-600 to-orange-600 bg-clip-text text-transparent">
                            {avgRating.toFixed(1)}/5.0
                          </div>
                          <div class="text-xs text-yellow-700 font-medium">Average Rating</div>
                        </div>
                        <div class="flex text-yellow-400">
                          {#each Array(5) as _, i}
                            <Star class="h-4 w-4 {i < Math.round(avgRating) ? 'fill-current' : 'text-gray-300'}" />
                          {/each}
                        </div>
                      </div>
                    </div>
                  {:else}
                    <!-- Most Popular Choice for non-rating questions -->
                    {@const topChoice = currentData.length > 0 ? currentData.reduce((max, item) => item.value > max.value ? item : max, currentData[0]) : null}
                    {#if topChoice}
                      <div class="bg-gradient-to-r from-blue-50 to-purple-50 mx-6 mb-6 rounded-xl p-4 border border-blue-200/50">
                        <div class="flex items-center justify-between">
                          <div>
                            <div class="text-sm font-medium text-blue-900 mb-1">
                              {currentMode === 'individual' ? 'Most Popular Option' : 'Most Common Combination'}
                            </div>
                            <div class="text-lg font-bold text-blue-800 leading-tight">{topChoice.label}</div>
                          </div>
                          <div class="text-right">
                            <div class="text-2xl font-bold bg-gradient-to-r from-blue-600 to-purple-600 bg-clip-text text-transparent">
                              {((topChoice.value / chart.totalResponses) * 100).toFixed(0)}%
                            </div>
                            <div class="text-xs text-blue-700 font-medium">{topChoice.value} selections</div>
                          </div>
                        </div>
                      </div>
                    {/if}
                  {/if}
                </Card>
              {/each}
            </div>
          </div>
        {/if}
      {/each}
    </div>
  {/if}
</div>

<style>
  .feedback-chart-widget {
    @apply w-full;
  }
</style>