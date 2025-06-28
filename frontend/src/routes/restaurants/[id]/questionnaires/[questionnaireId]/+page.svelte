<script lang="ts">
  import { onMount } from 'svelte';
  import { page } from '$app/stores';
  import { goto } from '$app/navigation';
  import { questionnaireStore, currentQuestionnaire, questionnaireLoading, questionnaireError } from '$lib/stores/questionnaires';
  import { Button } from '$lib/components/ui/button';
  import { Card, CardContent, CardHeader, CardTitle } from '$lib/components/ui/card';
  import { Badge } from '$lib/components/ui/badge';
  import { Alert, AlertDescription } from '$lib/components/ui/alert';
  import { Skeleton } from '$lib/components/ui/skeleton';
  import { ArrowLeft, Edit, Star, BarChart3, CheckSquare, Square, MessageSquare, CheckCircle, XCircle } from 'lucide-svelte';

  $: restaurantId = $page.params.id;
  $: questionnaireId = $page.params.questionnaireId;

  onMount(() => {
    if (restaurantId && questionnaireId) {
      questionnaireStore.loadQuestionnaire(restaurantId, questionnaireId);
    }
  });

  function goBack() {
    goto(`/restaurants/${restaurantId}/questionnaires`);
  }

  function editQuestionnaire() {
    goto(`/restaurants/${restaurantId}/questionnaires?edit=${questionnaireId}`);
  }

  function getQuestionIcon(type: string) {
    switch (type) {
      case 'rating': return Star;
      case 'scale': return BarChart3;
      case 'multi_choice': return CheckSquare;
      case 'single_choice': return Square;
      case 'text': return MessageSquare;
      case 'yes_no': return CheckCircle;
      default: return MessageSquare;
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

  function formatDate(dateString: string) {
    return new Date(dateString).toLocaleDateString('en-US', {
      year: 'numeric',
      month: 'long',
      day: 'numeric',
      hour: '2-digit',
      minute: '2-digit'
    });
  }
</script>

<svelte:head>
  <title>
    {$currentQuestionnaire ? `${$currentQuestionnaire.name} - Questionnaire` : 'Questionnaire'} - LeCritique
  </title>
</svelte:head>

<div class="space-y-6">
  <!-- Header -->
  <div class="flex items-center gap-4">
    <Button variant="ghost" size="sm" on:click={goBack}>
      <ArrowLeft class="h-4 w-4 mr-2" />
      Back to Questionnaires
    </Button>
  </div>

  {#if $questionnaireError}
    <Alert variant="destructive">
      <AlertDescription>{$questionnaireError}</AlertDescription>
    </Alert>
  {/if}

  {#if $questionnaireLoading}
    <div class="space-y-6">
      <Card>
        <CardHeader>
          <Skeleton class="h-8 w-1/2" />
          <Skeleton class="h-4 w-3/4" />
        </CardHeader>
        <CardContent>
          <div class="space-y-4">
            <Skeleton class="h-4 w-full" />
            <Skeleton class="h-4 w-2/3" />
          </div>
        </CardContent>
      </Card>
      
      {#each Array(3) as _}
        <Card>
          <CardHeader>
            <Skeleton class="h-6 w-1/3" />
          </CardHeader>
          <CardContent>
            <Skeleton class="h-20 w-full" />
          </CardContent>
        </Card>
      {/each}
    </div>
  {:else if $currentQuestionnaire}
    <div class="space-y-6">
      <!-- Questionnaire Info -->
      <Card>
        <CardHeader>
          <div class="flex items-start justify-between">
            <div class="flex-1">
              <CardTitle class="text-2xl">{$currentQuestionnaire.name}</CardTitle>
              {#if $currentQuestionnaire.description}
                <p class="text-muted-foreground mt-2">
                  {$currentQuestionnaire.description}
                </p>
              {/if}
            </div>
            
            <Button on:click={editQuestionnaire}>
              <Edit class="h-4 w-4 mr-2" />
              Edit
            </Button>
          </div>
        </CardHeader>
        
        <CardContent class="space-y-4">
          <div class="grid grid-cols-2 md:grid-cols-4 gap-4">
            <div class="space-y-1">
              <p class="text-sm text-muted-foreground">Type</p>
              <Badge variant="outline">
                {$currentQuestionnaire.dish_id ? 'Dish-specific' : 'General'}
              </Badge>
            </div>
            
            <div class="space-y-1">
              <p class="text-sm text-muted-foreground">Status</p>
              <div class="flex gap-1">
                {#if $currentQuestionnaire.is_default}
                  <Badge variant="default" class="text-xs">Default</Badge>
                {/if}
                <Badge variant={$currentQuestionnaire.is_active ? 'default' : 'secondary'} class="text-xs">
                  {$currentQuestionnaire.is_active ? 'Active' : 'Inactive'}
                </Badge>
              </div>
            </div>
            
            <div class="space-y-1">
              <p class="text-sm text-muted-foreground">Questions</p>
              <p class="font-medium">{$currentQuestionnaire.questions?.length || 0}</p>
            </div>
            
            <div class="space-y-1">
              <p class="text-sm text-muted-foreground">Created</p>
              <p class="text-sm">
                {$currentQuestionnaire.created_at ? formatDate($currentQuestionnaire.created_at) : 'Unknown'}
              </p>
            </div>
          </div>
        </CardContent>
      </Card>

      <!-- Questions -->
      {#if $currentQuestionnaire.questions && $currentQuestionnaire.questions.length > 0}
        <div class="space-y-4">
          <h2 class="text-xl font-semibold">Questions ({$currentQuestionnaire.questions.length})</h2>
          
          {#each $currentQuestionnaire.questions as question, index}
            <Card>
              <CardHeader>
                <div class="flex items-start gap-3">
                  <div class="flex items-center justify-center w-8 h-8 rounded-full bg-primary/10 text-primary font-medium text-sm">
                    {index + 1}
                  </div>
                  
                  <div class="flex-1">
                    <div class="flex items-center gap-2 mb-2">
                      <svelte:component this={getQuestionIcon(question.type)} class="h-4 w-4" />
                      <Badge variant="secondary" class="text-xs">
                        {getQuestionTypeLabel(question.type)}
                      </Badge>
                      {#if question.is_required}
                        <Badge variant="outline" class="text-xs">Required</Badge>
                      {/if}
                    </div>
                    
                    <h3 class="font-medium text-lg">
                      {question.text}
                      {#if question.is_required}
                        <span class="text-red-500">*</span>
                      {/if}
                    </h3>
                  </div>
                </div>
              </CardHeader>
              
              <CardContent>
                <!-- Question Preview -->
                <div class="space-y-3">
                  {#if question.type === 'rating'}
                    <div class="flex gap-1">
                      {#each Array(5) as _, i}
                        <Star class="h-6 w-6 text-yellow-400 fill-current" />
                      {/each}
                    </div>
                    <p class="text-sm text-muted-foreground">5-star rating scale</p>
                  
                  {:else if question.type === 'scale'}
                    <div class="space-y-2">
                      <div class="flex items-center gap-2">
                        <span class="text-sm text-muted-foreground">
                          {question.min_label || question.min_value || 1}
                        </span>
                        <div class="flex-1 h-2 bg-muted rounded">
                          <div class="h-full w-1/2 bg-primary rounded"></div>
                        </div>
                        <span class="text-sm text-muted-foreground">
                          {question.max_label || question.max_value || 10}
                        </span>
                      </div>
                      <p class="text-sm text-muted-foreground">
                        Scale from {question.min_value || 1} to {question.max_value || 10}
                      </p>
                    </div>
                  
                  {:else if question.type === 'multi_choice'}
                    <div class="space-y-2">
                      {#each question.options || [] as option}
                        <label class="flex items-center gap-2 cursor-pointer">
                          <CheckSquare class="h-4 w-4 text-muted-foreground" />
                          <span class="text-sm">{option}</span>
                        </label>
                      {/each}
                      <p class="text-sm text-muted-foreground">Multiple selections allowed</p>
                    </div>
                  
                  {:else if question.type === 'single_choice'}
                    <div class="space-y-2">
                      {#each question.options || [] as option}
                        <label class="flex items-center gap-2 cursor-pointer">
                          <Square class="h-4 w-4 text-muted-foreground" />
                          <span class="text-sm">{option}</span>
                        </label>
                      {/each}
                      <p class="text-sm text-muted-foreground">Single selection only</p>
                    </div>
                  
                  {:else if question.type === 'text'}
                    <div class="space-y-2">
                      <div class="w-full p-3 border rounded-md bg-muted/50">
                        <p class="text-sm text-muted-foreground italic">Text input area for customer feedback...</p>
                      </div>
                      <p class="text-sm text-muted-foreground">Open-ended text response</p>
                    </div>
                  
                  {:else if question.type === 'yes_no'}
                    <div class="space-y-2">
                      <div class="flex gap-4">
                        <label class="flex items-center gap-2 cursor-pointer">
                          <CheckCircle class="h-4 w-4 text-green-600" />
                          <span class="text-sm">Yes</span>
                        </label>
                        <label class="flex items-center gap-2 cursor-pointer">
                          <XCircle class="h-4 w-4 text-red-600" />
                          <span class="text-sm">No</span>
                        </label>
                      </div>
                      <p class="text-sm text-muted-foreground">Yes or No response</p>
                    </div>
                  {/if}
                </div>
              </CardContent>
            </Card>
          {/each}
        </div>
      {:else}
        <Card>
          <CardContent class="text-center py-12">
            <MessageSquare class="h-12 w-12 mx-auto text-muted-foreground mb-4" />
            <h3 class="text-lg font-medium mb-2">No questions</h3>
            <p class="text-muted-foreground mb-4">
              This questionnaire doesn't have any questions yet.
            </p>
            <Button on:click={editQuestionnaire}>
              <Edit class="h-4 w-4 mr-2" />
              Add Questions
            </Button>
          </CardContent>
        </Card>
      {/if}
    </div>
  {:else}
    <Card>
      <CardContent class="text-center py-12">
        <h3 class="text-lg font-medium mb-2">Questionnaire not found</h3>
        <p class="text-muted-foreground mb-4">
          The questionnaire you're looking for doesn't exist or has been deleted.
        </p>
        <Button on:click={goBack}>
          <ArrowLeft class="h-4 w-4 mr-2" />
          Back to Questionnaires
        </Button>
      </CardContent>
    </Card>
  {/if}
</div>