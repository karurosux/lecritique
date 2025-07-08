<script lang="ts">
  import { createEventDispatcher } from 'svelte';
  import { questionnaireStore, type Question, type GeneratedQuestion } from '$lib/stores/questionnaires';
  import { Button, Input, Card } from '$lib/components/ui';
  import { Plus, Trash2, Sparkles, X, Loader2 } from 'lucide-svelte';
  import QuestionEditor from './QuestionEditor.svelte';
  import { getApiClient } from '$lib/api';

  let { 
    restaurantId,
    dishId,
    initialData = null
  }: {
    restaurantId: string;
    dishId: string;
    initialData?: any;
  } = $props();

  const dispatch = createEventDispatcher();

  // Form state
  let selectedDishId = $state(dishId || initialData?.dish_id);
  let questions = $state<Question[]>(initialData?.questions || []);
  let generatedQuestions = $state<GeneratedQuestion[]>([]);
  
  // Auto-generate name from dish - we'll set this when saving

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

  // Debug: Log when questions change
  $effect(() => {
    console.log('Questions updated:', questions.length, questions);
  });

  // Debug: Log when initialData changes
  $effect(() => {
    console.log('InitialData:', initialData);
    if (initialData?.questions) {
      console.log('InitialData questions:', initialData.questions.length, initialData.questions);
    }
  });


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
      // Auto-generate name from existing data or use default
      const questionnaireData = {
        name: initialData?.name || 'Feedback Questionnaire',
        dish_id: selectedDishId
      };

      let questionnaire;
      if (initialData?.id) {
        // Update existing questionnaire
        questionnaire = await questionnaireStore.updateQuestionnaire(restaurantId, initialData.id, questionnaireData);
        
        // For existing questionnaires, we need to manage questions differently
        // This is a simplified approach - in production you'd want to handle adds/updates/deletes separately
        console.log('Updating existing questionnaire with questions:', questions);
        
        // For now, we'll just dispatch success - the questions management needs to be improved
        // TODO: Implement proper question update logic
      } else {
        // Create new questionnaire
        questionnaire = await questionnaireStore.createQuestionnaire(restaurantId, questionnaireData);
        
        // Add questions to the new questionnaire
        console.log('Adding questions to questionnaire:', questionnaire.id, 'Restaurant:', restaurantId);
        console.log('Questions to add:', questions);
        for (const question of questions) {
          console.log('Adding question:', question, 'to questionnaire:', questionnaire.id);
          await questionnaireStore.addQuestion(restaurantId, questionnaire.id!, question);
        }
      }

      dispatch('success', questionnaire);
    } catch (err) {
      error = err.message;
    } finally {
      loading = false;
    }
  }

  function cancel() {
    dispatch('cancel');
  }

  function getQuestionTypeIcon(type: string) {
    switch (type) {
      case 'rating': return 'Rating';
      case 'scale': return 'Scale';
      case 'multi_choice': return 'Multi';
      case 'single_choice': return 'Single';
      case 'text': return 'Text';
      case 'yes_no': return 'Y/N';
      default: return 'Other';
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

<div class="space-y-6">

  {#if error || storeError}
    <div class="bg-red-50 border border-red-200 rounded-md p-4">
      <div class="flex">
        <div class="flex-shrink-0">
          <X class="h-5 w-5 text-red-400" />
        </div>
        <div class="ml-3">
          <p class="text-sm text-red-800">{error || storeError}</p>
        </div>
      </div>
    </div>
  {/if}

  <!-- Questions -->
  <Card>
    <div class="p-6">
      <div class="flex items-center justify-between mb-4">
        <h3 class="text-lg font-medium text-gray-900">Questions ({questions.length})</h3>
      <div class="flex gap-2">
        <Button
          onclick={generateAIQuestions}
          disabled={generatingQuestions || !selectedDishId}
          variant="secondary"
          class="flex items-center gap-2"
          title={!selectedDishId ? 'Select a dish to generate AI questions' : ''}
        >
          {#if generatingQuestions}
            <Loader2 class="h-4 w-4 animate-spin" />
          {:else}
            <Sparkles class="h-4 w-4" />
          {/if}
          AI Questions
        </Button>
        <Button onclick={addQuestion} variant="gradient" class="flex items-center gap-2">
          <Plus class="h-4 w-4" />
          Add Question
        </Button>
      </div>
    </div>

      {#if questions.length === 0}
        <div class="text-center py-8 text-gray-500">
          <p>No questions added yet.</p>
          <p class="text-sm mt-1">Click "Add Question" or "AI Questions" to get started.</p>
        </div>
      {:else}
        <div class="space-y-3">
          {#each questions as question, index}
            <div class="flex items-center gap-3 p-4 bg-gray-50 rounded-lg border border-gray-200">
            <div class="flex-1">
              <div class="flex items-center gap-2 mb-1">
                <span class="px-2 py-1 bg-gray-100 text-xs font-medium rounded">
                  {getQuestionTypeLabel(question.type)}
                </span>
                {#if question.is_required}
                  <span class="px-2 py-1 bg-blue-100 text-blue-800 text-xs rounded">Required</span>
                {/if}
              </div>
              <p class="font-medium text-gray-900">{question.text}</p>
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
    </div>
  </Card>

  <!-- Actions -->
  <div class="flex items-center justify-end space-x-4 pt-6 border-t border-gray-200">
    <Button onclick={cancel} variant="outline">
      Cancel
    </Button>
    <Button 
      onclick={save} 
      disabled={loading || storeLoading}
      variant="gradient"
      class="min-w-32 shadow-lg"
    >
      {#if loading || storeLoading}
        <Loader2 class="h-4 w-4 mr-2 animate-spin" />
      {/if}
      {initialData?.id ? 'Update' : 'Create'} Questionnaire
    </Button>
  </div>
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
            <div class="p-4 bg-gray-50 rounded-lg border border-gray-200">
              <div class="flex items-center gap-2 mb-2">
                <span class="px-2 py-1 bg-gray-100 text-xs font-medium rounded">
                  {getQuestionTypeLabel(question.type)}
                </span>
              </div>
              <p class="font-medium text-gray-900">{question.text}</p>
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

        <div class="flex items-center justify-end space-x-4 pt-6 border-t border-gray-200">
          <Button
            onclick={() => { showGeneratedPreview = false; generatedQuestions = []; }}
            variant="outline"
          >
            Cancel
          </Button>
          <Button 
            onclick={addGeneratedQuestions} 
            variant="gradient"
            class="min-w-32 shadow-lg"
          >
            <Plus class="h-4 w-4 mr-2" />
            Add All Questions
          </Button>
        </div>
      </div>
    </Card>
  </div>
{/if}