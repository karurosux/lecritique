<script lang="ts">
  import { createEventDispatcher, onMount } from 'svelte';
  import { questionStore, type Question } from '$lib/stores/questions';
  import { Button, Input, Card, ConfirmDialog } from '$lib/components/ui';
  import { Plus, Trash2, Edit2, Loader2, Sparkles } from 'lucide-svelte';
  import QuestionEditor from './QuestionEditor.svelte';
  import { QuestionApi } from '$lib/api/question';
  import { toast } from 'svelte-sonner';

  let { 
    restaurantId,
    dishId
  }: {
    restaurantId: string;
    dishId: string;
  } = $props();

  const dispatch = createEventDispatcher();

  // State
  let questions = $state<Question[]>([]);
  let loading = $state(false);
  let error = $state('');
  let showQuestionEditor = $state(false);
  let editingQuestion = $state<Question | null>(null);
  let editingIndex = $state(-1);
  let aiGenerating = $state(false);
  let showDeleteConfirm = $state(false);
  let questionToDelete = $state<Question | null>(null);

  // Subscribe to store
  let storeError = $derived($questionStore?.error);
  let storeLoading = $derived($questionStore?.loading);

  onMount(async () => {
    await loadQuestions();
  });

  async function loadQuestions() {
    try {
      questions = await questionStore.loadQuestions(restaurantId, dishId);
    } catch (err) {
      error = err.message;
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

  async function saveQuestion(questionData: Question) {
    try {
      if (editingIndex === -1) {
        // Create new question
        const newQuestion = await questionStore.createQuestion(restaurantId, dishId, {
          text: questionData.text!,
          type: questionData.type!,
          is_required: questionData.is_required || false,
          options: questionData.options || [],
          min_value: questionData.min_value,
          max_value: questionData.max_value,
          min_label: questionData.min_label || '',
          max_label: questionData.max_label || ''
        });
        questions = [...questions, newQuestion].sort((a, b) => (a.display_order || 0) - (b.display_order || 0));
      } else {
        // Update existing question
        const updatedQuestion = await questionStore.updateQuestion(restaurantId, dishId, questionData.id!, {
          text: questionData.text!,
          type: questionData.type!,
          is_required: questionData.is_required || false,
          options: questionData.options || [],
          min_value: questionData.min_value,
          max_value: questionData.max_value,
          min_label: questionData.min_label || '',
          max_label: questionData.max_label || ''
        });
        questions[editingIndex] = updatedQuestion;
      }
      
      showQuestionEditor = false;
      editingQuestion = null;
      editingIndex = -1;
    } catch (err) {
      error = err.message;
    }
  }

  function handleDeleteQuestion(question: Question) {
    questionToDelete = question;
    showDeleteConfirm = true;
  }

  async function confirmDeleteQuestion() {
    if (!questionToDelete) return;
    
    try {
      await questionStore.deleteQuestion(restaurantId, dishId, questionToDelete.id!);
      questions = questions.filter(q => q.id !== questionToDelete.id);
      toast.success('Question deleted successfully');
    } catch (err) {
      error = err.message;
      toast.error('Failed to delete question');
    } finally {
      questionToDelete = null;
      showDeleteConfirm = false;
    }
  }

  function cancelEdit() {
    showQuestionEditor = false;
    editingQuestion = null;
    editingIndex = -1;
  }

  async function generateAIQuestions() {
    try {
      aiGenerating = true;
      error = '';
      
      const generatedQuestions = await QuestionApi.generateQuestions(dishId);
      
      if (generatedQuestions && generatedQuestions.length > 0) {
        // Add generated questions to existing list
        for (const genQuestion of generatedQuestions) {
          const newQuestion = await questionStore.createQuestion(restaurantId, dishId, {
            text: genQuestion.text,
            type: genQuestion.type,
            is_required: genQuestion.is_required || false,
            options: genQuestion.options || [],
            min_value: genQuestion.min_value,
            max_value: genQuestion.max_value,
            min_label: genQuestion.min_label || '',
            max_label: genQuestion.max_label || ''
          });
          questions = [...questions, newQuestion];
        }
        
        // Sort by display order
        questions = questions.sort((a, b) => (a.display_order || 0) - (b.display_order || 0));
        
        toast.success(`Added ${generatedQuestions.length} AI-generated questions!`);
      } else {
        toast.error('No questions were generated. Please try again.');
      }
    } catch (err) {
      console.error('AI generation error:', err);
      error = err.message || 'Failed to generate AI questions';
      toast.error('Failed to generate AI questions. Please try again.');
    } finally {
      aiGenerating = false;
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

<div class="question-builder space-y-6">
  <!-- Error Display -->
  {#if error || storeError}
    <div class="bg-red-50 border border-red-200 rounded-md p-4">
      <p class="text-sm text-red-800">{error || storeError}</p>
    </div>
  {/if}

  <!-- Questions -->
  <Card>
    <div class="p-6 relative">
      {#if aiGenerating}
        <div class="absolute inset-0 bg-white bg-opacity-75 z-10 flex items-center justify-center rounded-lg">
          <div class="text-center">
            <Loader2 class="h-8 w-8 animate-spin mx-auto mb-4 text-blue-600" />
            <p class="text-lg font-medium text-gray-900">Generating AI Questions...</p>
            <p class="text-sm text-gray-500 mt-1">Please wait while we create relevant questions for your dish</p>
          </div>
        </div>
      {/if}
      <div class="flex items-center justify-between mb-4">
        <h3 class="text-lg font-medium text-gray-900">Questions ({questions.length})</h3>
        <div class="flex items-center gap-2">
          <div class="relative group">
            <Button 
              onclick={generateAIQuestions} 
              disabled={aiGenerating || questions.length > 0}
              variant="outline" 
              class="flex items-center gap-2"
            >
              {#if aiGenerating}
                <Loader2 class="h-4 w-4 animate-spin" />
              {:else}
                <Sparkles class="h-4 w-4" />
              {/if}
              {aiGenerating ? 'Generating...' : 'Generate AI Questions'}
            </Button>
            {#if questions.length > 0 && !aiGenerating}
              <div class="absolute bottom-full left-1/2 transform -translate-x-1/2 mb-2 px-3 py-2 bg-gray-900 text-white text-sm rounded-lg opacity-0 group-hover:opacity-100 transition-opacity duration-200 pointer-events-none whitespace-nowrap z-10">
                AI generation is available when starting from scratch. Delete existing questions to use this feature.
                <div class="absolute top-full left-1/2 transform -translate-x-1/2 border-4 border-transparent border-t-gray-900"></div>
              </div>
            {/if}
          </div>
          <Button 
            onclick={addQuestion} 
            disabled={aiGenerating}
            variant="gradient" 
            class="flex items-center gap-2"
          >
            <Plus class="h-4 w-4" />
            Add Question
          </Button>
        </div>
      </div>

      {#if questions.length === 0}
        <div class="text-center py-8 text-gray-500">
          <p>No questions added yet.</p>
          <p class="text-sm mt-1">Click "Add Question" to get started.</p>
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
              <div class="flex items-center gap-2">
                <Button 
                  onclick={() => editQuestion(question, index)}
                  disabled={aiGenerating}
                  variant="ghost" 
                  size="sm"
                  class="p-2"
                >
                  <Edit2 class="h-4 w-4" />
                </Button>
                <Button 
                  onclick={() => handleDeleteQuestion(question)}
                  disabled={aiGenerating}
                  variant="ghost" 
                  size="sm"
                  class="p-2 text-red-600 hover:text-red-700"
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

</div>

<!-- Question Editor Modal -->
{#if showQuestionEditor && editingQuestion}
  <QuestionEditor
    bind:question={editingQuestion}
    on:save={async (e) => await saveQuestion(e.detail)}
    on:cancel={cancelEdit}
  />
{/if}

<!-- Delete Confirmation Dialog -->
<ConfirmDialog
  bind:isOpen={showDeleteConfirm}
  title="Delete Question?"
  message={questionToDelete ? `Are you sure you want to delete "${questionToDelete.text}"? This action cannot be undone.` : 'Are you sure you want to delete this question?'}
  confirmText="Delete"
  cancelText="Cancel"
  variant="danger"
  onConfirm={confirmDeleteQuestion}
  onCancel={() => {
    showDeleteConfirm = false;
    questionToDelete = null;
  }}
/>

