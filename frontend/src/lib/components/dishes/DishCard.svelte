<script lang="ts">
  import { Card, Button } from '$lib/components/ui';
  import { Edit2, Lightbulb, Eye, EyeOff, Trash2, ClipboardList, Clock, AlertTriangle } from 'lucide-svelte';

  interface Dish {
    id: string;
    name: string;
    description?: string;
    price: number;
    category: string;
    is_available: boolean;
    image_url?: string;
    allergens?: string[];
    preparation_time?: number;
    created_at: string;
    updated_at: string;
    has_questionnaire?: boolean;
    questionnaire_count?: number;
  }

  let {
    dish,
    onedit = (dish: Dish) => {},
    ontoggleavailability = (dish: Dish) => {},
    ondelete = (dish: Dish) => {},
    ongeneratequestionnaire = (dish: Dish) => {}
  }: {
    dish: Dish;
    onedit?: (dish: Dish) => void;
    ontoggleavailability?: (dish: Dish) => void;
    ondelete?: (dish: Dish) => void;
    ongeneratequestionnaire?: (dish: Dish) => void;
  } = $props();

  function formatPrice(price: number): string {
    return `$${price.toFixed(2)}`;
  }

  function handleEdit() {
    onedit(dish);
  }

  function handleToggleAvailability() {
    ontoggleavailability(dish);
  }

  function handleDelete() {
    ondelete(dish);
  }

  function handleGenerateQuestionnaire() {
    ongeneratequestionnaire(dish);
  }
</script>

<Card 
  variant="default" 
  hover 
  interactive 
  class="group transform transition-all duration-300 h-full"
>
  <div class="flex flex-col h-full">
    <div class="flex items-start justify-between mb-4">
    <div class="flex items-center space-x-3">
      <div class="w-12 h-12 bg-gradient-to-br from-green-500 to-emerald-600 rounded-xl flex items-center justify-center text-white font-bold text-lg shadow-lg shadow-green-500/25 group-hover:scale-110 transition-transform duration-200">
        {formatPrice(dish.price).replace('$', '')}
      </div>
      <div class="flex-1 min-w-0">
        <h3 class="font-bold text-lg text-gray-900 group-hover:text-green-600 transition-colors duration-200 truncate">
          {dish.name}
        </h3>
        <div class="flex items-center space-x-2 mt-1">
          <span class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium {dish.is_available ? 'bg-green-100 text-green-800' : 'bg-gray-100 text-gray-800'}">
            {dish.is_available ? 'Available' : 'Hidden'}
          </span>
          <span class="text-xs text-gray-500 bg-gray-100 px-2 py-1 rounded-md">
            {dish.category}
          </span>
          {#if dish.has_questionnaire}
            <span class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium bg-purple-100 text-purple-800">
              <ClipboardList class="h-3 w-3 mr-1" />
              Questionnaire
            </span>
          {/if}
        </div>
      </div>
    </div>
    <div class="flex items-center space-x-1 opacity-0 group-hover:opacity-100 transition-opacity duration-200">
      <button
        class="p-2 text-gray-400 hover:text-blue-600 hover:bg-blue-50 rounded-lg transition-all duration-200"
        onclick={(e) => { e.stopPropagation(); handleEdit(); }}
        title="Edit dish"
      >
        <Edit2 class="h-4 w-4" />
      </button>
      <button
        class="p-2 text-gray-400 hover:text-purple-600 hover:bg-purple-50 rounded-lg transition-all duration-200"
        onclick={(e) => { e.stopPropagation(); handleGenerateQuestionnaire(); }}
        title="{dish.has_questionnaire ? 'Manage questionnaire' : 'Create questionnaire'}"
      >
        {#if dish.has_questionnaire}
          <ClipboardList class="h-4 w-4" />
        {:else}
          <Lightbulb class="h-4 w-4" />
        {/if}
      </button>
      <button
        class="p-2 text-gray-400 hover:text-green-600 hover:bg-green-50 rounded-lg transition-all duration-200"
        onclick={(e) => { e.stopPropagation(); handleToggleAvailability(); }}
        title="Toggle availability"
      >
        {#if dish.is_available}
          <EyeOff class="h-4 w-4" />
        {:else}
          <Eye class="h-4 w-4" />
        {/if}
      </button>
      <button
        class="p-2 text-gray-400 hover:text-red-600 hover:bg-red-50 rounded-lg transition-all duration-200"
        onclick={(e) => { e.stopPropagation(); handleDelete(); }}
        title="Delete dish"
      >
        <Trash2 class="h-4 w-4" />
      </button>
    </div>
  </div>

  {#if dish.description}
    <p class="text-gray-600 text-sm mb-4 line-clamp-2 group-hover:text-gray-700 transition-colors duration-200">
      {dish.description}
    </p>
  {/if}

  <div class="grid grid-cols-1 gap-3 text-sm text-gray-600 min-h-[3rem]">
    {#if dish.preparation_time}
      <div class="flex items-center space-x-2">
        <Clock class="h-4 w-4 text-gray-400" />
        <span>{dish.preparation_time} minutes prep time</span>
      </div>
    {/if}
    {#if dish.allergens && dish.allergens.length > 0}
      <div class="flex items-center space-x-2">
        <AlertTriangle class="h-4 w-4 text-yellow-500" />
        <span class="truncate">
          {dish.allergens.slice(0, 3).join(', ')}
          {#if dish.allergens.length > 3}
            <span class="text-gray-400">+{dish.allergens.length - 3} more</span>
          {/if}
        </span>
      </div>
    {/if}
  </div>

  <!-- Spacer to push footer to bottom -->
  <div class="flex-grow"></div>

  <div class="flex items-center justify-between mt-4 pt-4 border-t border-gray-100">
    <span class="text-xs text-gray-500">
      Created {new Date(dish.created_at).toLocaleDateString()}
    </span>
  </div>
  </div>
</Card>

<style>
  .line-clamp-3 {
    display: -webkit-box;
    -webkit-line-clamp: 3;
    -webkit-box-orient: vertical;
    overflow: hidden;
  }
</style>