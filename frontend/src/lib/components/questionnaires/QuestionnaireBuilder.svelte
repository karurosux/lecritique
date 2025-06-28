<script lang="ts">
  import { createEventDispatcher } from 'svelte';
  import { questionnaireStore, type Question, type GeneratedQuestion } from '$lib/stores/questionnaires';
  import { Button, Input, Card } from '$lib/components/ui';
  import { Plus, Trash2, Sparkles, X } from 'lucide-svelte';
  import QuestionEditor from './QuestionEditor.svelte';

  export let restaurantId: string;
  export let dishId: string | null = null;
  export let initialData: any = null;

  const dispatch = createEventDispatcher();

  // Form state
  let name = initialData?.name || '';
  let description = initialData?.description || '';
  let isDefault = initialData?.is_default || false;
  let questions: Question[] = initialData?.questions || [];
  let generatedQuestions: GeneratedQuestion[] = [];

  // UI state
  let loading = false;
  let error = '';
  let showQuestionEditor = false;
  let editingQuestion: Question | null = null;
  let editingIndex = -1;
  let showGeneratedPreview = false;
  let generatingQuestions = false;

  // Subscribe to store
  $: storeError = $questionnaireStore?.error;
  $: storeLoading = $questionnaireStore?.loading;

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

  function saveQuestion(question: Question) {
    if (editingIndex >= 0) {
      questions[editingIndex] = question;
    } else {
      questions = [...questions, { ...question, display_order: questions.length + 1 }];
    }
    showQuestionEditor = false;
    editingQuestion = null;
    editingIndex = -1;
  }

  function deleteQuestion(index: number) {
    questions = questions.filter((_, i) => i !== index);
    // Reorder remaining questions
    questions = questions.map((q, i) => ({ ...q, display_order: i + 1 }));
  }

  async function generateAIQuestions() {
    if (!dishId) {
      error = 'Dish ID is required for AI question generation';
      return;
    }

    generatingQuestions = true;
    error = '';

    try {
      generatedQuestions = await questionnaireStore.generateQuestions(dishId);
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
    if (!name.trim()) {
      error = 'Questionnaire name is required';
      return;
    }

    loading = true;
    error = '';

    try {
      const data = {
        name: name.trim(),
        description: description.trim(),
        dish_id: dishId,
        is_default: isDefault
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
        Design feedback questions for your {dishId ? 'dish' : 'restaurant'}
      </p>
    </div>
    
    {#if dishId}
      <Button
        onclick={generateAIQuestions}
        disabled={generatingQuestions}
        class="flex items-center gap-2"
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
    {/if}
  </div>

  {#if error || storeError}
    <div class="bg-red-50 border border-red-200 text-red-700 px-4 py-3 rounded">
      {error || storeError}
    </div>
  {/if}

  <!-- Basic Information -->
  <Card class="p-6">
    <h3 class="text-lg font-medium mb-4">Basic Information</h3>
    <div class="space-y-4">
      <div>
        <label for="name" class="block text-sm font-medium text-gray-700 mb-1">Name *</label>
        <Input
          id="name"
          bind:value={name}
          placeholder="e.g., Margherita Pizza Feedback"
          required
        />
      </div>

      <div>
        <label for="description" class="block text-sm font-medium text-gray-700 mb-1">Description</label>
        <textarea
          id="description"
          bind:value={description}
          placeholder="Brief description of this questionnaire..."
          rows="2"
          class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
        ></textarea>
      </div>

      <div class="flex items-center space-x-2">
        <input
          type="checkbox"
          id="default"
          bind:checked={isDefault}
          class="rounded"
        />
        <label for="default" class="text-sm font-medium text-gray-700">Set as default questionnaire</label>
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
                onclick={() => deleteQuestion(index)}
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
</div>

<!-- Question Editor Modal -->
{#if showQuestionEditor && editingQuestion}
  <QuestionEditor
    bind:question={editingQuestion}
    on:save={(e) => saveQuestion(e.detail)}
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