<script lang="ts">
  import { page } from '$app/stores';
  import { onMount } from 'svelte';
  import { goto } from '$app/navigation';
  import { Card, Button, Input, Rating } from '$lib/components/ui';
  import { getApiClient, handleApiError } from '$lib/api/client';

  // Question types based on the backend
  type QuestionType = 'rating' | 'scale' | 'single_choice' | 'multiple_choice' | 'yes_no' | 'text';

  interface Question {
    id: string;
    text: string;
    type: QuestionType;
    required: boolean;
    options?: string[];
    min_value?: number;
    max_value?: number;
  }

  interface FeedbackData {
    restaurant_id: string;
    dish_id?: string;
    qr_code?: string;
    rating?: number;
    responses: Record<string, any>;
    comment?: string;
  }

  // State variables
  let loading = true;
  let submitting = false;
  let error = '';
  let restaurantId = '';
  let dishId = '';
  let qrCode = '';
  let restaurantName = '';
  let dishName = '';
  
  // Feedback form data
  let overallRating = 0;
  let responses: Record<string, any> = {};
  let comment = '';
  let customerEmail = '';
  
  // For now, using sample questions until questionnaire API is implemented
  let questions: Question[] = [
    {
      id: 'q1',
      text: 'How would you rate the taste?',
      type: 'rating',
      required: true,
      min_value: 1,
      max_value: 5
    },
    {
      id: 'q2',
      text: 'How was the service speed?',
      type: 'rating',
      required: true,
      min_value: 1,
      max_value: 5
    },
    {
      id: 'q3',
      text: 'Was the food served at the right temperature?',
      type: 'single_choice',
      required: true,
      options: ['Too cold', 'Just right', 'Too hot']
    },
    {
      id: 'q4',
      text: 'Would you recommend this to a friend?',
      type: 'yes_no',
      required: true
    },
    {
      id: 'q5',
      text: 'Any additional comments?',
      type: 'text',
      required: false
    }
  ];

  // Get URL parameters
  $: {
    restaurantId = $page.url.searchParams.get('restaurant') || '';
    dishId = $page.url.searchParams.get('dish') || '';
    qrCode = $page.url.searchParams.get('qr') || '';
  }

  onMount(async () => {
    if (!restaurantId) {
      error = 'Restaurant information is missing';
      loading = false;
      return;
    }
    
    // Initialize rating responses to 0
    questions.forEach(question => {
      if (question.type === 'rating') {
        responses[question.id] = 0;
      }
    });
    
    // In a real implementation, we would load the questionnaire and restaurant/dish info
    // For now, using placeholder data
    restaurantName = 'Demo Restaurant';
    dishName = dishId ? 'Sample Dish' : '';
    loading = false;
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
      if (question.required && !responses[question.id]) {
        error = `Please answer: ${question.text}`;
        return false;
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

<div class="min-h-screen bg-gray-50 py-8 px-4">
  <div class="max-w-2xl mx-auto">
    {#if loading}
      <!-- Loading State -->
      <Card>
        <div class="text-center py-12">
          <svg class="animate-spin h-8 w-8 text-blue-600 mx-auto mb-4" fill="none" viewBox="0 0 24 24">
            <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
            <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
          </svg>
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
        <Card>
          <div class="text-center">
            <h1 class="text-2xl font-bold text-gray-900 mb-2">Share Your Experience</h1>
            <p class="text-gray-600">
              {#if restaurantName}
                at <span class="font-medium">{restaurantName}</span>
              {/if}
              {#if dishName}
                - <span class="font-medium">{dishName}</span>
              {/if}
            </p>
          </div>
        </Card>

        <!-- Overall Rating -->
        <Card>
          <div class="text-center">
            <h2 class="text-lg font-medium text-gray-900 mb-4">Overall Rating</h2>
            <div class="flex justify-center">
              <Rating bind:value={overallRating} size="lg" showLabel />
            </div>
          </div>
        </Card>

        <!-- Dynamic Questions -->
        {#each questions as question}
          <Card>
            <div>
              <label class="block text-sm font-medium text-gray-700 mb-3">
                {question.text}
                {#if question.required}
                  <span class="text-red-500">*</span>
                {/if}
              </label>

              {#if question.type === 'rating'}
                <Rating 
                  bind:value={responses[question.id]} 
                  max={question.max_value || 5}
                  size="md"
                />

              {:else if question.type === 'single_choice' && question.options}
                <div class="space-y-2">
                  {#each question.options as option}
                    <label class="flex items-center">
                      <input
                        type="radio"
                        name="question_{question.id}"
                        value={option}
                        checked={responses[question.id] === option}
                        on:change={() => handleQuestionResponse(question.id, option)}
                        class="h-4 w-4 text-blue-600 focus:ring-blue-500"
                      />
                      <span class="ml-2 text-gray-700">{option}</span>
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
          variant="primary"
          size="lg"
          disabled={submitting}
          class="w-full"
        >
          {#if submitting}
            <svg class="animate-spin -ml-1 mr-3 h-5 w-5 text-white inline" fill="none" viewBox="0 0 24 24">
              <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
              <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
            </svg>
            Submitting...
          {:else}
            Submit Feedback
          {/if}
        </Button>
      </form>
    {/if}
  </div>
</div>