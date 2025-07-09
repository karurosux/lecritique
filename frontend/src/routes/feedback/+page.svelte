<script lang="ts">
  import { page } from '$app/stores';
  import { onMount } from 'svelte';
  import { goto } from '$app/navigation';
  import { Card, Button, Input, Rating } from '$lib/components/ui';
  import { getApiClient, handleApiError } from '$lib/api/client';
  import { Loader2 } from 'lucide-svelte';
  import { questionnaireStore, currentQuestionnaire, questionnaireLoading, questionnaireError } from '$lib/stores/questionnaire';
  import type { Question } from '$lib/stores/questionnaire';

  interface FeedbackData {
    restaurant_id: string;
    dish_id?: string;
    qr_code?: string;
    rating?: number;
    responses: Record<string, any>;
    comment?: string;
  }

  // State variables
  let submitting = false;
  let error = '';
  let restaurantId = '';
  let locationId = '';
  let dishId = '';
  let qrCode = '';
  let restaurantName = '';
  let dishName = '';
  
  // Feedback form data
  let overallRating = 0;
  let responses: Record<string, any> = {};
  let comment = '';
  let customerEmail = '';
  
  // Get questions from store
  $: questions = $currentQuestionnaire?.questions || [];
  $: loading = $questionnaireLoading;
  $: questionnaireErrorMsg = $questionnaireError;

  // Get URL parameters
  $: {
    restaurantId = $page.url.searchParams.get('restaurant') || '';
    locationId = $page.url.searchParams.get('location') || '';
    dishId = $page.url.searchParams.get('dish') || '';
    qrCode = $page.url.searchParams.get('qr') || '';
  }

  onMount(async () => {
    if (!restaurantId) {
      error = 'Restaurant information is missing';
      return;
    }
    
    // Fetch questionnaire or dish questions
    if (dishId) {
      await questionnaireStore.fetchDishQuestions(restaurantId, dishId);
    } else {
      await questionnaireStore.fetchQuestionnaire(restaurantId, locationId);
    }
    
    // Initialize responses based on fetched questions
    if ($currentQuestionnaire) {
      $currentQuestionnaire.questions.forEach(question => {
        if (question.type === 'rating' || question.type === 'scale') {
          responses[question.id] = 0;
        } else if (question.type === 'multiple_choice') {
          responses[question.id] = [];
        }
      });
    }
    
    // Fetch restaurant info if needed
    try {
      const api = getApiClient();
      // For now, using placeholder data
      restaurantName = 'Demo Restaurant';
      dishName = dishId ? 'Sample Dish' : '';
    } catch (err) {
      console.error('Error fetching restaurant info:', err);
    }
  });

  function handleRatingClick(value: number) {
    overallRating = value;
  }

  function handleQuestionResponse(questionId: string, value: any) {
    responses[questionId] = value;
  }

  function validateForm(): boolean {
    // Check overall rating
    if (overallRating === 0) {
      error = 'Please provide an overall rating';
      return false;
    }

    // Check required questions
    for (const question of questions) {
      if (question.required) {
        const response = responses[question.id];
        
        if (question.type === 'multiple_choice') {
          if (!response || response.length === 0) {
            error = `Please answer: ${question.text}`;
            return false;
          }
        } else if (question.type === 'rating' || question.type === 'scale') {
          if (!response || response === 0) {
            error = `Please answer: ${question.text}`;
            return false;
          }
        } else if (!response && response !== false) {
          // For yes/no, false is a valid response
          error = `Please answer: ${question.text}`;
          return false;
        }
      }
    }

    return true;
  }

  async function handleSubmit() {
    error = '';
    
    if (!validateForm()) {
      return;
    }

    submitting = true;

    try {
      const api = getApiClient();
      
      const feedbackData: FeedbackData = {
        restaurant_id: restaurantId,
        rating: overallRating,
        responses,
        comment
      };

      if (dishId) {
        feedbackData.dish_id = dishId;
      }

      if (qrCode) {
        feedbackData.qr_code = qrCode;
      }

      await api.api.v1PublicFeedbackCreate(feedbackData as any);
      
      // Redirect to success page
      goto('/feedback/success');
    } catch (err) {
      error = handleApiError(err);
      submitting = false;
    }
  }

  // Helper function to render stars
  function renderStars(rating: number, max: number = 5, size: string = 'h-8 w-8') {
    return Array(max).fill(0).map((_, i) => ({
      filled: i < rating,
      index: i + 1
    }));
  }
</script>

<svelte:head>
  <title>Give Feedback - LeCritique</title>
  <meta name="description" content="Share your dining experience and help us improve" />
</svelte:head>

<div class="min-h-screen bg-gradient-to-br from-slate-50 via-blue-50/30 to-indigo-50/50 py-8 px-4">
  <div class="max-w-2xl mx-auto">
    {#if loading}
      <!-- Loading State -->
      <Card>
        <div class="text-center py-12">
          <Loader2 class="animate-spin h-8 w-8 text-blue-600 mx-auto mb-4" />
          <p class="text-gray-600">Loading feedback form...</p>
        </div>
      </Card>
    
    {:else if !restaurantId}
      <!-- Error State -->
      <Card>
        <div class="text-center py-12">
          <svg class="h-12 w-12 text-red-500 mx-auto mb-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-.833-1.964-.833-2.732 0L3.732 16.5c-.77.833.192 3 1.732 3z" />
          </svg>
          <h2 class="text-xl font-semibold text-gray-900 mb-2">Missing Information</h2>
          <p class="text-gray-600">{error}</p>
        </div>
      </Card>
    
    {:else}
      <!-- Feedback Form -->
      <form on:submit|preventDefault={handleSubmit} class="space-y-6">
        <!-- Header -->
        <Card variant="gradient">
          <div class="text-center space-y-4">
            <div class="h-16 w-16 bg-gradient-to-r from-blue-500 to-purple-600 rounded-2xl flex items-center justify-center mx-auto shadow-lg shadow-blue-500/25">
              <svg class="h-8 w-8 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 12h.01M12 12h.01M16 12h.01M21 12c0 4.418-4.03 8-9 8a9.863 9.863 0 01-4.255-.949L3 20l1.395-3.72C3.512 15.042 3 13.574 3 12c0-4.418 4.03-8 9-8s9 3.582 9 8z" />
              </svg>
            </div>
            <div>
              <h1 class="text-3xl font-bold bg-gradient-to-r from-gray-900 to-gray-700 bg-clip-text text-transparent mb-3">
                Share Your Experience
              </h1>
              <p class="text-gray-600 text-lg font-medium">
                {#if restaurantName}
                  at <span class="text-blue-600 font-semibold">{restaurantName}</span>
                {/if}
                {#if dishName}
                  <br />
                  <span class="inline-flex items-center px-3 py-1 rounded-full text-sm font-medium bg-blue-100 text-blue-800 mt-2">
                    🍽️ {dishName}
                  </span>
                {/if}
              </p>
            </div>
          </div>
        </Card>

        <!-- Overall Rating -->
        <Card variant="elevated" padding={false}>
          <div class="p-8 text-center space-y-6">
            <div class="space-y-2">
              <h2 class="text-2xl font-bold bg-gradient-to-r from-gray-900 to-gray-700 bg-clip-text text-transparent">
                Overall Rating
              </h2>
              <p class="text-gray-600">How was your overall experience?</p>
            </div>
            <div class="flex justify-center">
              <Rating bind:value={overallRating} size="lg" showLabel />
            </div>
            {#if overallRating > 0}
              <div class="transition-all duration-300 ease-out">
                <div class="inline-flex items-center px-4 py-2 rounded-full bg-gradient-to-r {overallRating >= 4 ? 'from-green-100 to-emerald-100 text-green-800' : overallRating >= 3 ? 'from-yellow-100 to-orange-100 text-yellow-800' : 'from-red-100 to-pink-100 text-red-800'}">
                  <span class="font-semibold">
                    {overallRating >= 4 ? '🎉 Excellent!' : overallRating >= 3 ? '👍 Good!' : '🤔 We can do better!'}
                  </span>
                </div>
              </div>
            {/if}
          </div>
        </Card>

        <!-- Dynamic Questions -->
        {#each questions as question}
          <Card variant="glass" hover>
            <div class="space-y-4">
              <div class="space-y-2">
                <label class="block text-lg font-semibold text-gray-900">
                  {question.text}
                  {#if question.required}
                    <span class="text-red-500 ml-1">*</span>
                  {/if}
                </label>
                <div class="h-1 w-16 bg-gradient-to-r from-blue-500 to-purple-600 rounded-full"></div>
              </div>

              {#if question.type === 'rating'}
                <Rating 
                  bind:value={responses[question.id]} 
                  max={question.max_value || 5}
                  size="md"
                />

              {:else if question.type === 'scale'}
                <div>
                  <input
                    type="range"
                    min={question.min_value || 1}
                    max={question.max_value || 10}
                    bind:value={responses[question.id]}
                    class="w-full h-2 bg-gray-200 rounded-lg appearance-none cursor-pointer"
                  />
                  <div class="flex justify-between text-xs text-gray-600 mt-1">
                    <span>{question.min_label || question.min_value || 1}</span>
                    <span class="font-medium text-blue-600">{responses[question.id] || 0}</span>
                    <span>{question.max_label || question.max_value || 10}</span>
                  </div>
                </div>

              {:else if question.type === 'single_choice' && question.options}
                <div class="space-y-2">
                  {#each question.options as option}
                    <label class="flex items-center cursor-pointer">
                      <input
                        type="radio"
                        name="question_{question.id}"
                        value={option.value || option.text}
                        checked={responses[question.id] === (option.value || option.text)}
                        on:change={() => handleQuestionResponse(question.id, option.value || option.text)}
                        class="h-4 w-4 text-blue-600 focus:ring-blue-500 cursor-pointer"
                      />
                      <span class="ml-2 text-gray-700">{option.text}</span>
                    </label>
                  {/each}
                </div>

              {:else if question.type === 'multiple_choice' && question.options}
                <div class="space-y-2">
                  {#each question.options as option}
                    <label class="flex items-center cursor-pointer">
                      <input
                        type="checkbox"
                        value={option.value || option.text}
                        checked={responses[question.id]?.includes(option.value || option.text)}
                        on:change={(e) => {
                          const target = e.target as HTMLInputElement;
                          const value = option.value || option.text;
                          if (target.checked) {
                            responses[question.id] = [...(responses[question.id] || []), value];
                          } else {
                            responses[question.id] = responses[question.id].filter((v: any) => v !== value);
                          }
                        }}
                        class="h-4 w-4 text-blue-600 rounded focus:ring-blue-500 cursor-pointer"
                      />
                      <span class="ml-2 text-gray-700">{option.text}</span>
                    </label>
                  {/each}
                </div>

              {:else if question.type === 'yes_no'}
                <div class="flex space-x-4">
                  <Button
                    type="button"
                    variant={responses[question.id] === true ? 'primary' : 'outline'}
                    on:click={() => handleQuestionResponse(question.id, true)}
                  >
                    Yes
                  </Button>
                  <Button
                    type="button"
                    variant={responses[question.id] === false ? 'primary' : 'outline'}
                    on:click={() => handleQuestionResponse(question.id, false)}
                  >
                    No
                  </Button>
                </div>

              {:else if question.type === 'text'}
                <textarea
                  class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
                  rows="3"
                  placeholder="Your answer..."
                  bind:value={responses[question.id]}
                ></textarea>
              {/if}
            </div>
          </Card>
        {/each}

        <!-- Additional Comments -->
        <Card>
          <div>
            <label for="comment" class="block text-sm font-medium text-gray-700 mb-2">
              Additional Comments
            </label>
            <textarea
              id="comment"
              rows="4"
              class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
              placeholder="Tell us more about your experience..."
              bind:value={comment}
            ></textarea>
          </div>
        </Card>

        <!-- Optional Email -->
        <Card>
          <div>
            <label for="email" class="block text-sm font-medium text-gray-700 mb-2">
              Email (Optional)
            </label>
            <p class="text-sm text-gray-600 mb-3">
              Leave your email if you'd like us to follow up on your feedback
            </p>
            <Input
              id="email"
              type="email"
              placeholder="your@email.com"
              bind:value={customerEmail}
            />
          </div>
        </Card>

        <!-- Error Display -->
        {#if error}
          <div class="bg-red-50 border border-red-200 rounded-md p-4">
            <p class="text-sm text-red-800">{error}</p>
          </div>
        {/if}

        <!-- Submit Button -->
        <Button
          type="submit"
          variant="gradient"
          size="xl"
          disabled={submitting}
          loading={submitting}
          class="w-full shadow-lg shadow-blue-500/25 hover:shadow-xl hover:shadow-blue-500/40"
        >
          {#if submitting}
            <Loader2 class="animate-spin -ml-1 mr-3 h-5 w-5 text-white inline" />
            Submitting...
          {:else}
            Submit Feedback
          {/if}
        </Button>
      </form>
    {/if}
  </div>
</div>