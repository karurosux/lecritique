<script lang="ts">
  import { onMount } from 'svelte';
  import { page } from '$app/stores';
  import { goto } from '$app/navigation';
  import { questionnaireStore, questionnaires, questionnaireLoading, questionnaireError } from '$lib/stores/questionnaires';
  import { Button, Card, Input } from '$lib/components/ui';
  import { Plus, FileText, Edit, Trash2, Eye, Sparkles } from 'lucide-svelte';
  import QuestionnaireBuilder from '$lib/components/questionnaires/QuestionnaireBuilder.svelte';
  import type { Questionnaire } from '$lib/stores/questionnaires';

  $: restaurantId = $page.params.id;
  $: dishId = $page.url.searchParams.get('dishId');
  $: dishName = $page.url.searchParams.get('dishName');
  
  // UI state
  let showBuilder = false;
  let editingQuestionnaire: Questionnaire | null = null;
  let confirmDeleteId = '';

  onMount(() => {
    if (restaurantId) {
      questionnaireStore.loadQuestionnaires(restaurantId);
    }
    
    // If we have a dishId, open the builder immediately
    if (dishId) {
      createNew();
    }
  });

  function createNew() {
    editingQuestionnaire = null;
    showBuilder = true;
  }

  function editQuestionnaire(questionnaire: Questionnaire) {
    editingQuestionnaire = questionnaire;
    showBuilder = true;
  }

  function viewQuestionnaire(questionnaire: Questionnaire) {
    goto(`/restaurants/${restaurantId}/questionnaires/${questionnaire.id}`);
  }

  async function deleteQuestionnaire(questionnaire: Questionnaire) {
    if (confirmDeleteId === questionnaire.id) {
      try {
        await questionnaireStore.deleteQuestionnaire(restaurantId, questionnaire.id!);
        confirmDeleteId = '';
      } catch (error) {
        console.error('Failed to delete questionnaire:', error);
      }
    } else {
      confirmDeleteId = questionnaire.id!;
      // Auto-clear confirmation after 3 seconds
      setTimeout(() => {
        if (confirmDeleteId === questionnaire.id) {
          confirmDeleteId = '';
        }
      }, 3000);
    }
  }

  function handleQuestionnaireSaved() {
    showBuilder = false;
    editingQuestionnaire = null;
    questionnaireStore.loadQuestionnaires(restaurantId);
    
    // Clear URL params if we came from a dish
    if (dishId) {
      goto(`/restaurants/${restaurantId}/questionnaires`, { replaceState: true });
    }
  }

  function handleBuilderCancelled() {
    showBuilder = false;
    editingQuestionnaire = null;
    
    // Clear URL params if we came from a dish
    if (dishId) {
      goto(`/restaurants/${restaurantId}/questionnaires`, { replaceState: true });
    }
  }

  function getQuestionnaireType(questionnaire: Questionnaire) {
    return questionnaire.dish_id ? 'Dish-specific' : 'General';
  }

  function formatDate(dateString: string) {
    return new Date(dateString).toLocaleDateString();
  }
</script>

<svelte:head>
  <title>Questionnaires - LeCritique</title>
</svelte:head>

{#if showBuilder}
  <QuestionnaireBuilder
    {restaurantId}
    dishId={editingQuestionnaire?.dish_id || dishId}
    initialData={editingQuestionnaire}
    on:saved={handleQuestionnaireSaved}
    on:cancelled={handleBuilderCancelled}
  />
{:else}
  <div class="space-y-6">
    <!-- Header -->
    <div class="flex items-center justify-between">
      <div>
        <h1 class="text-3xl font-bold">Questionnaires</h1>
        <p class="text-muted-foreground">
          Manage feedback questionnaires for your restaurant and dishes
        </p>
      </div>
      
      <div class="flex gap-2">
        <Button onclick={() => goto(`/restaurants/${restaurantId}/dishes`)} variant="secondary">
          <Sparkles class="h-4 w-4 mr-2" />
          Generate for Dishes
        </Button>
        <Button onclick={createNew}>
          <Plus class="h-4 w-4 mr-2" />
          New Questionnaire
        </Button>
      </div>
    </div>

    {#if $questionnaireError}
      <div class="bg-red-50 border border-red-200 text-red-700 px-4 py-3 rounded">
        {$questionnaireError}
      </div>
    {/if}

    <!-- Questionnaires List -->
    {#if $questionnaireLoading}
      <div class="grid gap-4 md:grid-cols-2 lg:grid-cols-3">
        {#each Array(6) as _}
          <Card class="p-6">
            <div class="animate-pulse">
              <div class="h-6 bg-gray-200 rounded w-3/4 mb-2"></div>
              <div class="h-4 bg-gray-200 rounded w-full mb-4"></div>
              <div class="space-y-2">
                <div class="h-4 bg-gray-200 rounded w-1/2"></div>
                <div class="h-4 bg-gray-200 rounded w-2/3"></div>
              </div>
            </div>
          </Card>
        {/each}
      </div>
    {:else if $questionnaires.length === 0}
      <Card class="p-6">
        <div class="text-center py-12">
          <FileText class="h-12 w-12 mx-auto text-gray-400 mb-4" />
          <h3 class="text-lg font-medium mb-2">No questionnaires yet</h3>
          <p class="text-gray-600 mb-4">
            Create your first questionnaire to start collecting detailed feedback from customers.
          </p>
          <div class="flex justify-center gap-2">
            <Button onclick={() => goto(`/restaurants/${restaurantId}/dishes`)} variant="secondary">
              <Sparkles class="h-4 w-4 mr-2" />
              Generate AI Questionnaires
            </Button>
            <Button onclick={createNew}>
              <Plus class="h-4 w-4 mr-2" />
              Create Manually
            </Button>
          </div>
        </div>
      </Card>
    {:else}
      <div class="grid gap-4 md:grid-cols-2 lg:grid-cols-3">
        {#each $questionnaires as questionnaire}
          <Card class="hover:shadow-md transition-shadow p-6">
            <div class="space-y-4">
              <div class="flex items-start justify-between">
                <div class="flex-1">
                  <h3 class="text-lg font-medium">{questionnaire.name}</h3>
                  {#if questionnaire.description}
                    <p class="text-sm text-gray-600 mt-1">
                      {questionnaire.description}
                    </p>
                  {/if}
                </div>
                
                <div class="flex gap-1">
                  {#if questionnaire.is_default}
                    <span class="px-2 py-1 bg-blue-100 text-blue-800 text-xs rounded">Default</span>
                  {/if}
                  {#if !questionnaire.is_active}
                    <span class="px-2 py-1 bg-gray-100 text-gray-800 text-xs rounded">Inactive</span>
                  {/if}
                </div>
              </div>
              
              <div class="flex items-center justify-between text-sm">
                <span class="text-gray-600">Type:</span>
                <span class="px-2 py-1 bg-gray-100 text-xs rounded">
                  {getQuestionnaireType(questionnaire)}
                </span>
              </div>
              
              <div class="flex items-center justify-between text-sm">
                <span class="text-gray-600">Questions:</span>
                <span class="font-medium">
                  {questionnaire.questions?.length || 0}
                </span>
              </div>
              
              {#if questionnaire.created_at}
                <div class="flex items-center justify-between text-sm">
                  <span class="text-gray-600">Created:</span>
                  <span>{formatDate(questionnaire.created_at)}</span>
                </div>
              {/if}

              <!-- Actions -->
              <div class="flex gap-2 pt-2 border-t">
                <Button
                  onclick={() => viewQuestionnaire(questionnaire)}
                  variant="secondary"
                  class="flex-1 text-sm"
                >
                  <Eye class="h-4 w-4 mr-1" />
                  View
                </Button>
                
                <Button
                  onclick={() => editQuestionnaire(questionnaire)}
                  variant="secondary"
                  class="flex-1 text-sm"
                >
                  <Edit class="h-4 w-4 mr-1" />
                  Edit
                </Button>
                
                <Button
                  onclick={() => deleteQuestionnaire(questionnaire)}
                  variant="secondary"
                  class="text-red-600 hover:text-red-700"
                  disabled={questionnaire.is_default}
                >
                  <Trash2 class="h-4 w-4" />
                  {#if confirmDeleteId === questionnaire.id}
                    Confirm?
                  {/if}
                </Button>
              </div>
              
              {#if questionnaire.is_default}
                <p class="text-xs text-gray-500">
                  Default questionnaires cannot be deleted
                </p>
              {/if}
            </div>
          </Card>
        {/each}
      </div>
    {/if}
  </div>
{/if}