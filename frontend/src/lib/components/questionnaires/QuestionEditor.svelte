<script lang="ts">
  import { createEventDispatcher } from 'svelte';
  import type { Question } from '$lib/stores/questionnaires';
  import { Button, Input, Card, Select } from '$lib/components/ui';
  import { Plus, Trash2, X } from 'lucide-svelte';

  export let question: Question;

  const dispatch = createEventDispatcher();

  // Local state for editing
  let localQuestion = { ...question };
  let options = [...(question.options || [])];
  let newOption = '';

  // Question type definitions
  const questionTypes = [
    { value: 'rating', label: 'Star Rating (1-5 stars)', icon: '⭐' },
    { value: 'scale', label: 'Scale (1-10 with labels)', icon: '📊' },
    { value: 'multi_choice', label: 'Multiple Choice', icon: '☑️' },
    { value: 'single_choice', label: 'Single Choice', icon: '🔘' },
    { value: 'text', label: 'Text Input', icon: '💬' },
    { value: 'yes_no', label: 'Yes/No', icon: '✅' }
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
      options: ['multi_choice', 'single_choice'].includes(localQuestion.type) ? options : undefined
    };
  }

  function handleTypeChange(event: Event) {
    const target = event.target as HTMLSelectElement;
    const newType = target.value;
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
      text: localQuestion.text.trim()
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

<div class="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
  <Card class="w-full max-w-2xl max-h-[90vh] overflow-y-auto m-4">
    <div class="p-6">
      <div class="flex items-center justify-between mb-6">
        <h3 class="text-lg font-medium">Edit Question</h3>
        <Button onclick={cancel} variant="secondary">
          <X class="h-4 w-4" />
        </Button>
      </div>

      <div class="space-y-6">
        <!-- Question Text -->
        <div>
          <label for="question-text" class="block text-sm font-medium text-gray-700 mb-1">Question Text *</label>
          <textarea
            id="question-text"
            bind:value={localQuestion.text}
            placeholder="Enter your question here..."
            rows="2"
            class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
            required
          ></textarea>
        </div>

        <!-- Question Type -->
        <div>
          <label for="question-type" class="block text-sm font-medium text-gray-700 mb-1">Question Type *</label>
          <select
            id="question-type"
            bind:value={localQuestion.type}
            on:change={handleTypeChange}
            class="w-full px-3 py-2 border border-gray-300 rounded-md focus:outline-none focus:ring-2 focus:ring-blue-500"
          >
            {#each questionTypes as type}
              <option value={type.value}>
                {type.icon} {type.label}
              </option>
            {/each}
          </select>
        </div>

        <!-- Required Toggle -->
        <div class="flex items-center space-x-2">
          <input
            type="checkbox"
            id="required"
            bind:checked={localQuestion.is_required}
            class="rounded"
          />
          <label for="required" class="text-sm font-medium text-gray-700">Required question</label>
        </div>

        <!-- Type-specific configuration -->
        {#if localQuestion.type === 'scale'}
          <div class="space-y-4 p-4 border rounded-lg bg-gray-50">
            <h4 class="font-medium">Scale Configuration</h4>
            
            <div class="grid grid-cols-2 gap-4">
              <div>
                <label for="min-value" class="block text-sm font-medium text-gray-700 mb-1">Minimum Value</label>
                <Input
                  id="min-value"
                  type="number"
                  bind:value={localQuestion.min_value}
                  min="1"
                  max="10"
                />
              </div>
              <div>
                <label for="max-value" class="block text-sm font-medium text-gray-700 mb-1">Maximum Value</label>
                <Input
                  id="max-value"
                  type="number"
                  bind:value={localQuestion.max_value}
                  min="2"
                  max="10"
                />
              </div>
            </div>

            <div class="grid grid-cols-2 gap-4">
              <div>
                <label for="min-label" class="block text-sm font-medium text-gray-700 mb-1">Minimum Label</label>
                <Input
                  id="min-label"
                  bind:value={localQuestion.min_label}
                  placeholder="e.g., Poor, Terrible"
                />
              </div>
              <div>
                <label for="max-label" class="block text-sm font-medium text-gray-700 mb-1">Maximum Label</label>
                <Input
                  id="max-label"
                  bind:value={localQuestion.max_label}
                  placeholder="e.g., Excellent, Amazing"
                />
              </div>
            </div>
          </div>
        {/if}

        {#if ['multi_choice', 'single_choice'].includes(localQuestion.type)}
          <div class="space-y-4 p-4 border rounded-lg bg-gray-50">
            <h4 class="font-medium">
              {localQuestion.type === 'multi_choice' ? 'Multiple Choice' : 'Single Choice'} Options
            </h4>
            
            <!-- Existing options -->
            <div class="space-y-2">
              {#each options as option, index}
                <div class="flex items-center gap-2">
                  <Input
                    bind:value={options[index]}
                    on:input={updateQuestion}
                    placeholder="Option text"
                  />
                  <Button
                    onclick={() => removeOption(index)}
                    variant="secondary"
                    class="text-red-600 hover:text-red-700"
                  >
                    <Trash2 class="h-4 w-4" />
                  </Button>
                </div>
              {/each}
            </div>

            <!-- Add new option -->
            <div class="flex items-center gap-2">
              <Input
                bind:value={newOption}
                placeholder="Add new option..."
                on:keydown={(e) => e.key === 'Enter' && addOption()}
              />
              <Button onclick={addOption} disabled={!newOption.trim()} variant="secondary">
                <Plus class="h-4 w-4" />
              </Button>
            </div>

            {#if options.length === 0}
              <p class="text-sm text-gray-500">
                Add at least two options for this question type.
              </p>
            {/if}
          </div>
        {/if}

        <!-- Preview -->
        <div class="space-y-2 p-4 border rounded-lg bg-gray-50">
          <h4 class="font-medium">Preview</h4>
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
                  <span class="text-2xl text-yellow-400">⭐</span>
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
                {#each (options.length > 0 ? options : ['Option 1', 'Option 2']) as option}
                  <label class="flex items-center gap-2">
                    <input type="checkbox" class="rounded" />
                    <span class="text-sm">{option}</span>
                  </label>
                {/each}
              </div>
            {:else if localQuestion.type === 'single_choice'}
              <div class="space-y-1">
                {#each (options.length > 0 ? options : ['Option 1', 'Option 2']) as option}
                  <label class="flex items-center gap-2">
                    <input type="radio" name="preview" class="rounded-full" />
                    <span class="text-sm">{option}</span>
                  </label>
                {/each}
              </div>
            {:else if localQuestion.type === 'text'}
              <textarea class="w-full p-2 border rounded" rows="3" placeholder="Text input area..."></textarea>
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

        <!-- Actions -->
        <div class="flex justify-end gap-3 pt-4 border-t">
          <Button onclick={cancel} variant="secondary">
            Cancel
          </Button>
          <Button
            onclick={save}
            disabled={!localQuestion.text.trim() || 
                     (['multi_choice', 'single_choice'].includes(localQuestion.type) && options.filter(o => o.trim()).length < 2)}
          >
            Save Question
          </Button>
        </div>
      </div>
    </div>
  </Card>
</div>