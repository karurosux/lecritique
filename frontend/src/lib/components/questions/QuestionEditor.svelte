<script lang="ts">
  import { createEventDispatcher } from 'svelte';
  import type { Question } from '$lib/stores/questions';
  import { Button, Input, Card, Select, Modal } from '$lib/components/ui';
  import {
    Plus,
    Trash2,
    X,
    Star,
    BarChart3,
    CheckSquare,
    Circle,
    MessageSquare,
    ToggleLeft,
  } from 'lucide-svelte';

  export let question: Question;

  const dispatch = createEventDispatcher();

  // Local state for editing
  let localQuestion = { ...question };
  let options = [...(question.options || [])];
  let newOption = '';

  // Question type definitions
  const questionTypes = [
    { value: 'rating', label: 'Star Rating (1-5 stars)' },
    { value: 'scale', label: 'Scale (1-10 with labels)' },
    { value: 'multi_choice', label: 'Multiple Choice' },
    { value: 'single_choice', label: 'Single Choice' },
    { value: 'text', label: 'Text Input' },
    { value: 'yes_no', label: 'Yes/No' },
  ];

  function addOption() {
    if (newOption.trim()) {
      options = [...options, newOption.trim()];
      newOption = '';
      updateQuestion();
    }
  }

  function removeOption(index: number) {
    options = options.filter((_, i) => i !== index);
    updateQuestion();
  }

  function updateQuestion() {
    localQuestion = {
      ...localQuestion,
      options: ['multi_choice', 'single_choice'].includes(localQuestion.type)
        ? options
        : undefined,
    };
  }

  function handleTypeChange(event: Event | { target: { value: string } }) {
    const newType = 'target' in event ? (event.target as any).value : event;
    localQuestion.type = newType as Question['type'];

    // Reset type-specific fields
    localQuestion.options = undefined;
    localQuestion.min_value = undefined;
    localQuestion.max_value = undefined;
    localQuestion.min_label = undefined;
    localQuestion.max_label = undefined;

    // Set defaults for specific types
    if (newType === 'scale') {
      localQuestion.min_value = 1;
      localQuestion.max_value = 10;
      localQuestion.min_label = 'Poor';
      localQuestion.max_label = 'Excellent';
    } else if (['multi_choice', 'single_choice'].includes(newType)) {
      options = options.length > 0 ? options : ['Option 1', 'Option 2'];
      localQuestion.options = [...options];
    }

    updateQuestion();
  }

  function save() {
    if (!localQuestion.text.trim()) {
      return;
    }

    // Final validation and cleanup
    const finalQuestion: Question = {
      ...localQuestion,
      text: localQuestion.text.trim(),
    };

    // Ensure options are set for choice types
    if (['multi_choice', 'single_choice'].includes(finalQuestion.type)) {
      finalQuestion.options = options.filter(opt => opt.trim());
    }

    dispatch('save', finalQuestion);
  }

  function cancel() {
    dispatch('cancel');
  }

  // Reactive updates
  $: if (localQuestion.type) {
    updateQuestion();
  }
</script>

<Modal
  isOpen={true}
  title="{question.text ? 'Edit' : 'Add'} Question"
  size="lg"
  onclose={cancel}>
  <div class="space-y-6">
    <!-- Question Text -->
    <div>
      <label
        for="question-text"
        class="block text-sm font-medium text-gray-700 mb-2">
        Question Text <span class="text-red-500">*</span>
      </label>
      <textarea
        id="question-text"
        bind:value={localQuestion.text}
        placeholder="What would you like to ask your customers?"
        rows="3"
        class="w-full px-3 py-2 border border-gray-300 rounded-lg focus:outline-none focus:ring-2 focus:ring-blue-500 focus:border-blue-500 transition-colors"
        required></textarea>
    </div>

    <!-- Question Type -->
    <div>
      <label
        for="question-type"
        class="block text-sm font-medium text-gray-700 mb-2">
        Question Type <span class="text-red-500">*</span>
      </label>
      <Select
        id="question-type"
        value={localQuestion.type}
        options={questionTypes}
        onchange={value => {
          localQuestion.type = value;
          handleTypeChange({ target: { value } });
        }}
        placeholder="Select question type" />
    </div>

    <!-- Required Toggle -->
    <div class="flex items-center space-x-2">
      <input
        type="checkbox"
        id="required"
        bind:checked={localQuestion.is_required}
        class="rounded" />
      <label for="required" class="text-sm font-medium text-gray-700"
        >Required question</label>
    </div>

    <!-- Type-specific configuration -->
    {#if localQuestion.type === 'scale'}
      <Card variant="minimal">
        <div class="p-4 space-y-4">
          <h4 class="text-sm font-medium text-gray-900">Scale Configuration</h4>

          <div class="grid grid-cols-2 gap-4">
            <div>
              <label
                for="min-value"
                class="block text-sm font-medium text-gray-700 mb-1"
                >Minimum Value</label>
              <Input
                id="min-value"
                type="number"
                bind:value={localQuestion.min_value}
                min="1"
                max="10" />
            </div>
            <div>
              <label
                for="max-value"
                class="block text-sm font-medium text-gray-700 mb-1"
                >Maximum Value</label>
              <Input
                id="max-value"
                type="number"
                bind:value={localQuestion.max_value}
                min="2"
                max="10" />
            </div>
          </div>

          <div class="grid grid-cols-2 gap-4">
            <div>
              <label
                for="min-label"
                class="block text-sm font-medium text-gray-700 mb-1"
                >Minimum Label</label>
              <Input
                id="min-label"
                bind:value={localQuestion.min_label}
                placeholder="e.g., Poor, Terrible" />
            </div>
            <div>
              <label
                for="max-label"
                class="block text-sm font-medium text-gray-700 mb-1"
                >Maximum Label</label>
              <Input
                id="max-label"
                bind:value={localQuestion.max_label}
                placeholder="e.g., Excellent, Amazing" />
            </div>
          </div>
        </div>
      </Card>
    {/if}

    {#if ['multi_choice', 'single_choice'].includes(localQuestion.type)}
      <Card variant="minimal">
        <div class="p-4 space-y-4">
          <h4 class="text-sm font-medium text-gray-900">
            {localQuestion.type === 'multi_choice'
              ? 'Multiple Choice'
              : 'Single Choice'} Options
          </h4>

          <!-- Existing options -->
          <div class="space-y-2">
            {#each options as option, index}
              <div class="flex items-center gap-2">
                <Input
                  bind:value={options[index]}
                  on:input={updateQuestion}
                  placeholder="Option text" />
                <Button
                  onclick={() => removeOption(index)}
                  variant="outline"
                  class="p-2">
                  <Trash2 class="h-4 w-4 text-red-600" />
                </Button>
              </div>
            {/each}
          </div>

          <!-- Add new option -->
          <div class="flex items-center gap-2">
            <Input
              bind:value={newOption}
              placeholder="Add new option..."
              on:keydown={e => e.key === 'Enter' && addOption()} />
            <Button
              onclick={addOption}
              disabled={!newOption.trim()}
              variant="gradient"
              class="px-3">
              <Plus class="h-4 w-4" />
            </Button>
          </div>

          {#if options.length === 0}
            <p class="text-sm text-gray-500">
              Add at least two options for this question type.
            </p>
          {/if}
        </div>
      </Card>
    {/if}

    <!-- Preview -->
    <Card>
      <div class="p-4 space-y-2">
        <h4 class="text-sm font-medium text-gray-900">Preview</h4>
        <div class="space-y-2">
          <p class="font-medium">
            {localQuestion.text || 'Question text will appear here...'}
            {#if localQuestion.is_required}
              <span class="text-red-500">*</span>
            {/if}
          </p>

          {#if localQuestion.type === 'rating'}
            <div class="flex gap-1">
              {#each Array(5) as _, i}
                <Star class="h-6 w-6 text-yellow-400 fill-yellow-400" />
              {/each}
            </div>
          {:else if localQuestion.type === 'scale'}
            <div class="flex items-center gap-2">
              <span class="text-sm text-gray-500">
                {localQuestion.min_label || localQuestion.min_value || 1}
              </span>
              <div class="flex-1 h-2 bg-gray-200 rounded"></div>
              <span class="text-sm text-gray-500">
                {localQuestion.max_label || localQuestion.max_value || 10}
              </span>
            </div>
          {:else if localQuestion.type === 'multi_choice'}
            <div class="space-y-1">
              {#each options.length > 0 ? options : ['Option 1', 'Option 2'] as option}
                <label class="flex items-center gap-2">
                  <input type="checkbox" class="rounded" />
                  <span class="text-sm">{option}</span>
                </label>
              {/each}
            </div>
          {:else if localQuestion.type === 'single_choice'}
            <div class="space-y-1">
              {#each options.length > 0 ? options : ['Option 1', 'Option 2'] as option}
                <label class="flex items-center gap-2">
                  <input type="radio" name="preview" class="rounded-full" />
                  <span class="text-sm">{option}</span>
                </label>
              {/each}
            </div>
          {:else if localQuestion.type === 'text'}
            <textarea
              class="w-full p-2 border border-gray-300 rounded-lg"
              rows="3"
              placeholder="Customer's response will appear here..."></textarea>
          {:else if localQuestion.type === 'yes_no'}
            <div class="flex gap-4">
              <label class="flex items-center gap-2">
                <input type="radio" name="yesno" class="rounded-full" />
                <span class="text-sm">Yes</span>
              </label>
              <label class="flex items-center gap-2">
                <input type="radio" name="yesno" class="rounded-full" />
                <span class="text-sm">No</span>
              </label>
            </div>
          {/if}
        </div>
      </div>
    </Card>
  </div>

  <!-- Actions -->
  <div
    class="mt-6 pt-6 border-t border-gray-200 flex items-center justify-end space-x-3">
    <Button onclick={cancel} variant="outline">Cancel</Button>
    <Button
      onclick={save}
      disabled={!localQuestion.text.trim() ||
        (['multi_choice', 'single_choice'].includes(localQuestion.type) &&
          options.filter(o => o.trim()).length < 2)}
      variant="gradient"
      class="min-w-32">
      {question.text ? 'Save Changes' : 'Add Question'}
    </Button>
  </div>
</Modal>
