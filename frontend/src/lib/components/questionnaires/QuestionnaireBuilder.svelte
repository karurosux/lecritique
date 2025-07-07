<script lang="ts">
  import { createEventDispatcher, onMount } from 'svelte';
  import { questionnaireStore, type Question, type GeneratedQuestion } from '$lib/stores/questionnaires';
  import { Button, Input, Card, Select } from '$lib/components/ui';
  import { Plus, Trash2, Sparkles, X, ChefHat } from 'lucide-svelte';
  import QuestionEditor from './QuestionEditor.svelte';
  import { getApiClient } from '$lib/api';
  import { goto } from '$app/navigation';

  let { 
    restaurantId,
    dishId = null,
    initialData = null
  }: {
    restaurantId: string;
    dishId?: string | null;
    initialData?: any;
  } = $props();

  const dispatch = createEventDispatcher();

  // Form state
  let selectedDishId = $state(dishId || initialData?.dish_id || null);
  let questions = $state<Question[]>(initialData?.questions || []);
  let generatedQuestions = $state<GeneratedQuestion[]>([]);
  
  // Dishes for dropdown
  let dishes = $state<any[]>([]);
  let loadingDishes = $state(false);
  
  // Auto-generate name from selected dish
  let selectedDish = $derived(dishes.find(d => d.value === selectedDishId));
  let name = $derived(
    initialData?.name || 
    (selectedDish ? `${selectedDish.label} Feedback` : '')
  );

  // UI state
  let loading = $state(false);
  let error = $state('');
  let showQuestionEditor = $state(false);
  let editingQuestion = $state<Question | null>(null);
  let editingIndex = $state(-1);
  let showGeneratedPreview = $state(false);
  let generatingQuestions = $state(false);

  // Subscribe to store
  let storeError = $derived($questionnaireStore?.error);
  let storeLoading = $derived($questionnaireStore?.loading);

  onMount(async () => {
    await fetchDishes();
  });

  async function fetchDishes() {
    try {
      loadingDishes = true;
      const api = getApiClient();
      const response = await api.api.v1RestaurantsDishesList(restaurantId);
      
      if (response.data.success && response.data.data) {
        dishes = response.data.data.map((dish: any) => ({
          value: dish.id,
          label: dish.name
        }));
        
        // If no dishes, show error
        if (dishes.length === 0) {
          error = 'No dishes found. Please add dishes to your menu first.';
        }
      }
    } catch (err) {
      console.error('Failed to load dishes:', err);
      error = 'Failed to load dishes. Please try again.';
    } finally {
      loadingDishes = false;
    }
  }

  function addQuestion() {
    editingQuestion = {
      text: '',
      type: 'rating',
      is_required: true,
      display_order: questions.length + 1,
      options: []
    };
    editingIndex = -1;
    showQuestionEditor = true;
  }

  function editQuestion(question: Question, index: number) {
    editingQuestion = { ...question };
    editingIndex = index;
    showQuestionEditor = true;
  }

  async function saveQuestion(question: Question) {
    // If we're editing an existing questionnaire, save questions directly to backend
    if (initialData?.id) {
      loading = true;
      try {
        if (editingIndex >= 0 && questions[editingIndex]?.id) {
          // Update existing question
          const updatedQuestion = await questionnaireStore.updateQuestion(
            restaurantId, 
            initialData.id, 
            questions[editingIndex].id!,
            question
          );
          questions[editingIndex] = updatedQuestion;
        } else {
          // Add new question
          const newQuestion = await questionnaireStore.addQuestion(
            restaurantId,
            initialData.id,
            { ...question, display_order: questions.length + 1 }
          );
          questions = [...questions, newQuestion];
        }
      } catch (err) {
        error = err.message;
      } finally {
        loading = false;
      }
    } else {
      // For new questionnaires, just update local state
      if (editingIndex >= 0) {
        questions[editingIndex] = question;
      } else {
        questions = [...questions, { ...question, display_order: questions.length + 1 }];
      }
    }
    
    showQuestionEditor = false;
    editingQuestion = null;
    editingIndex = -1;
  }

  async function deleteQuestion(index: number) {
    const question = questions[index];
    
    // If we're editing an existing questionnaire and the question has an ID, delete from backend
    if (initialData?.id && question?.id) {
      loading = true;
      try {
        await questionnaireStore.deleteQuestion(restaurantId, initialData.id, question.id);
      } catch (err) {
        error = err.message;
        return;
      } finally {
        loading = false;
      }
    }
    
    questions = questions.filter((_, i) => i !== index);
    // Reorder remaining questions
    questions = questions.map((q, i) => ({ ...q, display_order: i + 1 }));
  }

  async function generateAIQuestions() {
    if (!selectedDishId) {
      error = 'Please select a dish for AI question generation';
      return;
    }

    generatingQuestions = true;
    error = '';

    try {
      generatedQuestions = await questionnaireStore.generateQuestions(selectedDishId);
      showGeneratedPreview = true;
    } catch (err) {
      error = err.message;
    } finally {
      generatingQuestions = false;
    }
  }

  function addGeneratedQuestions() {
    const newQuestions: Question[] = generatedQuestions.map((gq, index) => ({
      text: gq.text,
      type: gq.type,
      is_required: true,
      display_order: questions.length + index + 1,
      options: gq.options || [],
      min_value: gq.min_value,
      max_value: gq.max_value,
      min_label: gq.min_label,
      max_label: gq.max_label
    }));

    questions = [...questions, ...newQuestions];
    showGeneratedPreview = false;
    generatedQuestions = [];
  }

  async function save() {
    if (!selectedDishId) {
      error = 'Please select a dish for this questionnaire';
      return;
    }

    loading = true;
    error = '';

    try {
      const data = {
        name: name.trim(),
        dish_id: selectedDishId
      };

      let questionnaire;
      if (initialData?.id) {
        // Update existing questionnaire
        questionnaire = await questionnaireStore.updateQuestionnaire(restaurantId, initialData.id, data);
      } else {
        // Create new questionnaire
        questionnaire = await questionnaireStore.createQuestionnaire(restaurantId, data);
        
        // Add questions to the new questionnaire
        for (const question of questions) {
          await questionnaireStore.addQuestion(restaurantId, questionnaire.id!, question);
        }
      }

      dispatch('saved', questionnaire);
    } catch (err) {
      error = err.message;
    } finally {
      loading = false;
    }
  }

  function cancel() {
    dispatch('cancelled');
  }

  function getQuestionTypeIcon(type: string) {
    switch (type) {
      case 'rating': return '⭐';
      case 'scale': return '📊';
      case 'multi_choice': return '☑️';
      case 'single_choice': return '🔘';
      case 'text': return '💬';
      case 'yes_no': return '✅';
      default: return '❓';
    }
  }

  function getQuestionTypeLabel(type: string) {
    switch (type) {
      case 'rating': return 'Star Rating';
      case 'scale': return 'Scale';
      case 'multi_choice': return 'Multiple Choice';
      case 'single_choice': return 'Single Choice';
      case 'text': return 'Text Input';
      case 'yes_no': return 'Yes/No';
      default: return type;
    }
  }
</script>

<div class="space-y-6 p-6">
  <!-- Header -->
  <div class="flex items-center justify-between">
    <div>
      <h2 class="text-2xl font-bold">
        {initialData?.id ? 'Edit' : 'Create'} Questionnaire
      </h2>
      <p class="text-gray-600">
        Design feedback questions for your dishes
      </p>
    </div>
    
    <Button
      onclick={generateAIQuestions}
      disabled={generatingQuestions || !selectedDishId}
      class="flex items-center gap-2"
      title={!selectedDishId ? 'Select a dish to generate AI questions' : ''}
    >
        {#if generatingQuestions}
          <svg class="animate-spin h-4 w-4" fill="none" viewBox="0 0 24 24">
            <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
            <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
          </svg>
        {:else}
          <Sparkles class="h-4 w-4" />
        {/if}
        Generate AI Questions
      </Button>
  </div>

  {#if error || storeError}
    <div class="bg-red-50 border border-red-200 text-red-700 px-4 py-3 rounded">
      {error || storeError}
    </div>
  {/if}

  {#if dishes.length === 0 && !loadingDishes}
    <Card class="p-6 bg-yellow-50 border-yellow-200">
      <div class="flex items-start gap-4">
        <ChefHat class="h-8 w-8 text-yellow-600 flex-shrink-0 mt-1" />
        <div class="flex-1">
          <h3 class="text-lg font-medium text-yellow-900 mb-1">No Dishes Found</h3>
          <p class="text-yellow-700 mb-3">
            You need to add dishes to your menu before creating questionnaires. 
            Each questionnaire is designed to collect feedback about a specific dish.
          </p>
          <Button 
            onclick={() => goto(`/restaurants/${restaurantId}/dishes`)}
            class="flex items-center gap-2"
          >
            <ChefHat class="h-4 w-4" />
            Go to Menu Management
          </Button>
        </div>
      </div>
    </Card>
  {:else}

  <!-- Dish Selection -->
  <Card class="p-6">
    <div class="space-y-4">
      <div>
        <label for="dish" class="block text-sm font-medium text-gray-700 mb-1">
          Select Dish *
        </label>
        <Select
          id="dish"
          bind:value={selectedDishId}
          options={dishes}
          disabled={loadingDishes}
          placeholder={loadingDishes ? 'Loading dishes...' : 'Select a dish'}
          required
        />
        {#if selectedDish}
          <p class="text-sm text-gray-600 mt-2">
            📋 This will create: <span class="font-medium">{name}</span>
          </p>
        {/if}
      </div>
    </div>
  </Card>

  <!-- Questions -->
  <Card class="p-6">
    <div class="flex items-center justify-between mb-4">
      <h3 class="text-lg font-medium">Questions ({questions.length})</h3>
      <Button onclick={addQuestion} class="flex items-center gap-2">
        <Plus class="h-4 w-4" />
        Add Question
      </Button>
    </div>

    {#if questions.length === 0}
      <div class="text-center py-8 text-gray-500">
        <p>No questions added yet.</p>
        <p class="text-sm">Click "Add Question" or "Generate AI Questions" to get started.</p>
      </div>
    {:else}
      <div class="space-y-3">
        {#each questions as question, index}
          <div class="flex items-center gap-3 p-3 border rounded-lg">
            <div class="flex-1">
              <div class="flex items-center gap-2 mb-1">
                <span class="text-sm">{getQuestionTypeIcon(question.type)}</span>
                <span class="px-2 py-1 bg-gray-100 text-xs rounded">
                  {getQuestionTypeLabel(question.type)}
                </span>
                {#if question.is_required}
                  <span class="px-2 py-1 bg-blue-100 text-blue-800 text-xs rounded">Required</span>
                {/if}
              </div>
              <p class="font-medium">{question.text}</p>
              {#if question.options && question.options.length > 0}
                <p class="text-sm text-gray-500">
                  Options: {question.options.join(', ')}
                </p>
              {/if}
              {#if question.min_value !== undefined && question.max_value !== undefined}
                <p class="text-sm text-gray-500">
                  Scale: {question.min_value} to {question.max_value}
                </p>
              {/if}
            </div>

            <div class="flex gap-1">
              <Button
                onclick={() => editQuestion(question, index)}
                variant="secondary"
                class="text-sm"
              >
                Edit
              </Button>
              <Button
                onclick={async () => await deleteQuestion(index)}
                variant="secondary"
                class="text-red-600 hover:text-red-700"
              >
                <Trash2 class="h-4 w-4" />
              </Button>
            </div>
          </div>
        {/each}
      </div>
    {/if}
  </Card>

  <!-- Actions -->
  <div class="flex justify-end gap-3">
    <Button onclick={cancel} variant="secondary">
      Cancel
    </Button>
    <Button onclick={save} disabled={loading || storeLoading}>
      {#if loading || storeLoading}
        <svg class="animate-spin h-4 w-4 mr-2" fill="none" viewBox="0 0 24 24">
          <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
          <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 014 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
        </svg>
      {/if}
      {initialData?.id ? 'Update' : 'Create'} Questionnaire
    </Button>
  </div>
  {/if}
</div>

<!-- Question Editor Modal -->
{#if showQuestionEditor && editingQuestion}
  <QuestionEditor
    bind:question={editingQuestion}
    on:save={async (e) => await saveQuestion(e.detail)}
    on:cancel={() => { showQuestionEditor = false; editingQuestion = null; }}
  />
{/if}

<!-- Generated Questions Preview -->
{#if showGeneratedPreview && generatedQuestions.length > 0}
  <div class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
    <Card class="w-full max-w-2xl max-h-[80vh] overflow-y-auto m-4">
      <div class="p-6">
        <div class="flex items-center justify-between mb-4">
          <h3 class="text-lg font-medium flex items-center gap-2">
            <Sparkles class="h-5 w-5" />
            AI Generated Questions
          </h3>
          <Button
            onclick={() => { showGeneratedPreview = false; generatedQuestions = []; }}
            variant="secondary"
          >
            <X class="h-4 w-4" />
          </Button>
        </div>

        <div class="space-y-3 mb-6">
          {#each generatedQuestions as question, index}
            <div class="p-3 border rounded-lg">
              <div class="flex items-center gap-2 mb-2">
                <span>{getQuestionTypeIcon(question.type)}</span>
                <span class="px-2 py-1 bg-gray-100 text-xs rounded">
                  {getQuestionTypeLabel(question.type)}
                </span>
              </div>
              <p class="font-medium">{question.text}</p>
              {#if question.options && question.options.length > 0}
                <p class="text-sm text-gray-500 mt-1">
                  Options: {question.options.join(', ')}
                </p>
              {/if}
              {#if question.min_value !== undefined && question.max_value !== undefined}
                <p class="text-sm text-gray-500 mt-1">
                  Scale: {question.min_value} to {question.max_value}
                  {#if question.min_label || question.max_label}
                    ({question.min_label || ''} - {question.max_label || ''})
                  {/if}
                </p>
              {/if}
            </div>
          {/each}
        </div>

        <div class="flex justify-end gap-3 pt-4 border-t">
          <Button
            onclick={() => { showGeneratedPreview = false; generatedQuestions = []; }}
            variant="secondary"
          >
            Cancel
          </Button>
          <Button onclick={addGeneratedQuestions} class="flex items-center gap-2">
            <Plus class="h-4 w-4" />
            Add All Questions
          </Button>
        </div>
      </div>
    </Card>
  </div>
{/if}