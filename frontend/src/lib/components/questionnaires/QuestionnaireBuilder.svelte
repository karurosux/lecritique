<script lang="ts">
  import { onMount } from 'svelte';
  import { Button, Card, ConfirmDialog, NoDataAvailable } from '$lib/components/ui';
  import { Plus, Trash2, Sparkles, Loader2, Edit2, GripVertical, AlertTriangle, ClipboardList } from 'lucide-svelte';
  import QuestionEditor from './QuestionEditor.svelte';
  import { QuestionApi } from '$lib/api/question';
  import type { Question, GeneratedQuestion } from '$lib/api/questionnaire';
  import { toast } from 'svelte-sonner';

  let { 
    organizationId,
    productId
  }: {
    organizationId: string;
    productId: string;
  } = $props();

  // State
  let questions = $state<Question[]>([]);
  let loading = $state(false);
  let error = $state('');
  let showQuestionEditor = $state(false);
  let editingQuestion = $state<Question | null>(null);
  let editingIndex = $state(-1);
  let generatingQuestions = $state(false);
  let showDeleteConfirm = $state(false);
  let questionToDelete = $state<Question | null>(null);
  let draggingIndex = $state<number | null>(null);
  let dragOverIndex = $state<number | null>(null);
  let reordering = $state(false);

  onMount(async () => {
    await loadQuestionnaire();
  });

  async function loadQuestionnaire() {
    try {
      loading = true;
      error = '';
      
      // Load existing questions for this product directly
      const productQuestions = await QuestionApi.getQuestionsByProduct(organizationId, productId);
      questions = productQuestions || [];
      
    } catch (err: any) {
      error = err.message || 'Failed to load questions';
      console.error('Error loading questions:', err);
    } finally {
      loading = false;
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
      loading = true;
      error = '';

      if (editingIndex === -1) {
        // Create new question directly for the product
        const newQuestion = await QuestionApi.createQuestion(organizationId, productId, {
          text: questionData.text!,
          type: questionData.type!,
          is_required: questionData.is_required || false,
          options: questionData.options || [],
          min_value: questionData.min_value,
          max_value: questionData.max_value,
          min_label: questionData.min_label || '',
          max_label: questionData.max_label || ''
        });
        questions = [...questions, newQuestion];
        toast.success('Question added successfully');
      } else {
        // Update existing question
        const updatedQuestion = await QuestionApi.updateQuestion(
          organizationId,
          productId,
          questions[editingIndex].id!,
          {
            text: questionData.text!,
            type: questionData.type!,
            is_required: questionData.is_required || false,
            options: questionData.options || [],
            min_value: questionData.min_value,
            max_value: questionData.max_value,
            min_label: questionData.min_label || '',
            max_label: questionData.max_label || ''
          }
        );
        questions[editingIndex] = updatedQuestion;
        toast.success('Question updated successfully');
      }
      
      showQuestionEditor = false;
      editingQuestion = null;
      editingIndex = -1;
    } catch (err: any) {
      error = err.message || 'Failed to save question';
      toast.error('Failed to save question');
      console.error('Error saving question:', err);
    } finally {
      loading = false;
    }
  }

  function handleDeleteQuestion(question: Question) {
    questionToDelete = question;
    showDeleteConfirm = true;
  }

  async function confirmDeleteQuestion() {
    if (!questionToDelete) return;
    
    try {
      loading = true;
      await QuestionApi.deleteQuestion(organizationId, productId, questionToDelete.id!);
      questions = questions.filter(q => q.id !== questionToDelete.id);
      toast.success('Question deleted successfully');
    } catch (err: any) {
      error = err.message || 'Failed to delete question';
      toast.error('Failed to delete question');
      console.error('Error deleting question:', err);
    } finally {
      loading = false;
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
      generatingQuestions = true;
      error = '';
      
      const generatedQuestions = await QuestionApi.generateQuestions(organizationId, productId);
      
      if (generatedQuestions && generatedQuestions.length > 0) {
        // Add generated questions directly to the product
        for (const genQuestion of generatedQuestions) {
          const newQuestion = await QuestionApi.createQuestion(organizationId, productId, {
            text: genQuestion.text!,
            type: genQuestion.type!,
            is_required: genQuestion.is_required || false,
            options: genQuestion.options || [],
            min_value: genQuestion.min_value,
            max_value: genQuestion.max_value,
            min_label: genQuestion.min_label || '',
            max_label: genQuestion.max_label || ''
          });
          questions = [...questions, newQuestion];
        }
        
        toast.success(`Added ${generatedQuestions.length} AI-generated questions!`);
      } else {
        toast.error('No questions were generated. Please try again.');
      }
    } catch (err: any) {
      console.error('AI generation error:', err);
      error = err.message || 'Failed to generate AI questions';
      toast.error('Failed to generate AI questions. Please try again.');
    } finally {
      generatingQuestions = false;
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

  // Drag and drop handlers
  function handleDragStart(e: DragEvent, index: number) {
    draggingIndex = index;
    if (e.dataTransfer) {
      e.dataTransfer.effectAllowed = 'move';
    }
  }

  function handleDragOver(e: DragEvent, index: number) {
    e.preventDefault();
    if (e.dataTransfer) {
      e.dataTransfer.dropEffect = 'move';
    }
    dragOverIndex = index;
  }

  function handleDragLeave() {
    dragOverIndex = null;
  }

  async function handleDrop(e: DragEvent, dropIndex: number) {
    e.preventDefault();
    
    if (draggingIndex === null || draggingIndex === dropIndex) {
      draggingIndex = null;
      dragOverIndex = null;
      return;
    }

    // Reorder the questions array
    const newQuestions = [...questions];
    const [draggedQuestion] = newQuestions.splice(draggingIndex, 1);
    newQuestions.splice(dropIndex, 0, draggedQuestion);
    
    questions = newQuestions;
    
    // Reset drag state
    draggingIndex = null;
    dragOverIndex = null;
    
    // Save the new order to the backend
    await saveQuestionOrder();
  }

  async function saveQuestionOrder() {
    try {
      reordering = true;
      const questionIds = questions.map(q => q.id!);
      await QuestionApi.reorderQuestions(organizationId, productId, questionIds);
      toast.success('Question order updated');
    } catch (err: any) {
      error = err.message || 'Failed to update question order';
      toast.error('Failed to update question order');
      // Reload questions to restore original order
      await loadQuestionnaire();
    } finally {
      reordering = false;
    }
  }
</script>

<style>
  :global(.animate-modal-enter) {
    animation: modal-enter 0.2s ease-out forwards;
  }
  
  @keyframes modal-enter {
    from {
      opacity: 0;
      transform: scale(0.95) translateY(-10px);
    }
    to {
      opacity: 1;
      transform: scale(1) translateY(0);
    }
  }
</style>

<div class="question-builder space-y-6">
  <!-- Error Display -->
  {#if error}
    <NoDataAvailable
      title="Error Loading Questions"
      description={error}
      icon={AlertTriangle}
    />
  {/if}

  <!-- Loading overlay for AI generation -->
  {#if generatingQuestions}
    <div 
      class="fixed inset-0 z-[50000] flex items-center justify-center transition-all duration-200"
      style="background: linear-gradient(135deg, rgb(243 244 246 / 0.6), rgb(209 213 219 / 0.6)); backdrop-filter: blur(16px);"
      role="dialog"
      aria-modal="true"
      tabindex="-1"
    >
      <div class="relative bg-white rounded-lg shadow-xl max-w-md w-full mx-4 p-8 text-center animate-modal-enter">
        <!-- Animated AI Icon -->
        <div class="relative mb-6">
          <div class="w-20 h-20 mx-auto bg-gradient-to-br from-purple-500 to-blue-600 rounded-full flex items-center justify-center shadow-lg shadow-purple-500/25">
            <Sparkles class="h-10 w-10 text-white animate-pulse" />
          </div>
          <div class="absolute inset-0 w-20 h-20 mx-auto">
            <div class="w-full h-full border-4 border-purple-200 rounded-full animate-spin border-t-purple-500"></div>
          </div>
        </div>
        
        <!-- Title with gradient -->
        <h3 class="text-xl font-bold bg-gradient-to-r from-purple-600 to-blue-600 bg-clip-text text-transparent mb-3">
          Generating AI Questions...
        </h3>
        
        <!-- Description -->
        <p class="text-gray-600 mb-4 leading-relaxed">
          Our AI is analyzing your product and creating relevant feedback questions for your customers
        </p>
        
        <!-- Progress dots -->
        <div class="flex justify-center space-x-2">
          <div class="w-2 h-2 bg-purple-500 rounded-full animate-bounce" style="animation-delay: 0ms"></div>
          <div class="w-2 h-2 bg-purple-500 rounded-full animate-bounce" style="animation-delay: 150ms"></div>
          <div class="w-2 h-2 bg-purple-500 rounded-full animate-bounce" style="animation-delay: 300ms"></div>
        </div>
        
        <!-- Estimated time -->
        <p class="text-xs text-gray-500 mt-4">
          This usually takes 5-10 seconds
        </p>
      </div>
    </div>
  {/if}

  <!-- Questions -->
  <Card>
    <div class="p-6">
      <div class="flex items-center justify-between mb-4">
        <h3 class="text-lg font-medium text-gray-900">Questions ({questions.length})</h3>
        <div class="flex items-center gap-2">
          <div class="relative group">
            <Button 
              onclick={generateAIQuestions} 
              disabled={loading || generatingQuestions || questions.length > 0}
              variant="outline" 
              class="flex items-center gap-2"
            >
              {#if generatingQuestions}
                <Loader2 class="h-4 w-4 animate-spin" />
              {:else}
                <Sparkles class="h-4 w-4" />
              {/if}
              {generatingQuestions ? 'Generating...' : 'Generate AI Questions'}
            </Button>
            {#if questions.length > 0 && !generatingQuestions}
              <div class="absolute bottom-full left-1/2 transform -translate-x-1/2 mb-2 px-3 py-2 bg-gray-900 text-white text-sm rounded-lg opacity-0 group-hover:opacity-100 transition-opacity duration-200 pointer-events-none whitespace-nowrap z-10">
                AI generation is available when starting from scratch. Delete existing questions to use this feature.
                <div class="absolute top-full left-1/2 transform -translate-x-1/2 border-4 border-transparent border-t-gray-900"></div>
              </div>
            {/if}
          </div>
          <Button 
            onclick={addQuestion} 
            disabled={loading || generatingQuestions || reordering}
            variant="gradient" 
            class="flex items-center gap-2"
          >
            <Plus class="h-4 w-4" />
            Add Question
          </Button>
        </div>
      </div>

      {#if loading && !generatingQuestions}
        <NoDataAvailable
          title="Loading Questions..."
          description="Please wait while we load your questionnaire"
          icon={ClipboardList}
          variant="inline"
        />
      {:else if questions.length === 0}
        <NoDataAvailable
          title="No Questions Added Yet"
          description="Click 'Add Question' to get started or use AI to generate questions automatically"
          icon={ClipboardList}
          variant="inline"
        />
      {:else}
        <div class="space-y-3">
          {#each questions as question, index}
            <div 
              class="flex items-center gap-3 p-4 bg-gray-50 rounded-lg border border-gray-200 transition-all duration-200 {draggingIndex === index ? 'opacity-50' : ''} {dragOverIndex === index ? 'border-blue-400 border-2' : ''}"
              draggable="true"
              ondragstart={(e) => handleDragStart(e, index)}
              ondragover={(e) => handleDragOver(e, index)}
              ondragleave={handleDragLeave}
              ondrop={(e) => handleDrop(e, index)}
            >
              <div class="cursor-move" class:opacity-50={loading || generatingQuestions || reordering}>
                <GripVertical class="h-5 w-5 text-gray-400" />
              </div>
              <div class="flex-1">
                <div class="flex items-center gap-2 mb-1">
                  <span class="px-2 py-1 bg-gray-100 text-xs font-medium rounded">
                    {getQuestionTypeLabel(question.type!)}
                  </span>
                  {#if question.is_required}
                    <span class="px-2 py-1 bg-blue-100 text-blue-800 text-xs rounded">Required</span>
                  {/if}
                </div>
                <p class="font-medium text-gray-900">{question.text}</p>
                {#if question.options && question.options.length > 0}
                  <p class="text-sm text-gray-500">
                    Options: {Array.isArray(question.options) ? question.options.join(', ') : 'Invalid options'}
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
                  disabled={loading || generatingQuestions || reordering}
                  variant="ghost" 
                  size="sm"
                  class="p-2"
                >
                  <Edit2 class="h-4 w-4" />
                </Button>
                <Button 
                  onclick={() => handleDeleteQuestion(question)}
                  disabled={loading || generatingQuestions || reordering}
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
    bind:isOpen={showQuestionEditor}
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
