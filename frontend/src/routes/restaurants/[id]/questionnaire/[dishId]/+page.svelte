<script lang="ts">
  import { page } from '$app/stores';
  import { goto } from '$app/navigation';
  import { Button, Card } from '$lib/components/ui';
  import { ArrowLeft, ClipboardList } from 'lucide-svelte';
  import QuestionBuilder from '$lib/components/questions/QuestionBuilder.svelte';
  import { getApiClient } from '$lib/api';
  import { onMount } from 'svelte';

  let restaurantId = $page.params.id;
  let dishId = $page.params.dishId;
  
  let dishName = $state('');
  let loading = $state(true);
  let error = $state('');

  onMount(async () => {
    await fetchDishName();
  });

  async function fetchDishName() {
    try {
      const api = getApiClient();
      const response = await api.api.v1RestaurantsDishesList(restaurantId);
      
      if (response.data.success && response.data.data) {
        const dish = response.data.data.find((d: any) => d.id === dishId);
        if (dish) {
          dishName = dish.name;
        } else {
          error = 'Dish not found';
        }
      }
    } catch (err) {
      console.error('Failed to load dish:', err);
      error = 'Failed to load dish information';
    } finally {
      loading = false;
    }
  }

  function goBack() {
    goto(`/restaurants/${restaurantId}/dishes`);
  }

</script>

<div class="space-y-6">
  <!-- Header -->
  <div class="flex items-center gap-4">
    <Button 
      onclick={goBack}
      variant="ghost" 
      size="sm"
      class="flex items-center gap-2"
    >
      <ArrowLeft class="h-4 w-4" />
      Back to Dishes
    </Button>
    <div class="flex items-center gap-3">
      <div class="h-10 w-10 bg-gradient-to-br from-purple-500 to-blue-600 rounded-xl flex items-center justify-center shadow-lg shadow-purple-500/25">
        <ClipboardList class="h-5 w-5 text-white" />
      </div>
      <div>
        <h1 class="text-xl font-bold text-gray-900">
          {loading ? 'Loading...' : dishName ? `${dishName} Questionnaire` : 'Create Questionnaire'}
        </h1>
        <p class="text-sm text-gray-600">Design targeted feedback questions for this dish</p>
      </div>
    </div>
  </div>

  <!-- Content -->
  <div class="space-y-6">
    {#if error}
      <Card class="p-6">
        <div class="bg-red-50 border border-red-200 rounded-lg p-4">
          <p class="text-red-800 font-medium">{error}</p>
        </div>
      </Card>
    {/if}

    {#if !loading && !error}
      <QuestionBuilder
        {restaurantId}
        {dishId}
      />
    {/if}
  </div>
</div>