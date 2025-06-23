<script lang="ts">
  import { Card, Button } from '$lib/components/ui';

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
  }

  let {
    dish,
    onedit = (dish: Dish) => {},
    ontoggleavailability = (dish: Dish) => {},
    ondelete = (dish: Dish) => {}
  }: {
    dish: Dish;
    onedit?: (dish: Dish) => void;
    ontoggleavailability?: (dish: Dish) => void;
    ondelete?: (dish: Dish) => void;
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
</script>

<Card 
  variant="default" 
  hover 
  interactive 
  class="group transform transition-all duration-300"
>
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
        </div>
      </div>
    </div>
    <div class="flex items-center space-x-1 opacity-0 group-hover:opacity-100 transition-opacity duration-200">
      <button
        class="p-2 text-gray-400 hover:text-blue-600 hover:bg-blue-50 rounded-lg transition-all duration-200"
        onclick={(e) => { e.stopPropagation(); handleEdit(); }}
        title="Edit dish"
      >
        <svg class="h-4 w-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M11 5H6a2 2 0 00-2 2v11a2 2 0 002 2h11a2 2 0 002-2v-5m-1.414-9.414a2 2 0 112.828 2.828L11.828 15H9v-2.828l8.586-8.586z" />
        </svg>
      </button>
      <button
        class="p-2 text-gray-400 hover:text-green-600 hover:bg-green-50 rounded-lg transition-all duration-200"
        onclick={(e) => { e.stopPropagation(); handleToggleAvailability(); }}
        title="Toggle availability"
      >
        {#if dish.is_available}
          <svg class="h-4 w-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M13.875 18.825A10.05 10.05 0 0112 19c-4.478 0-8.268-2.943-9.543-7a9.97 9.97 0 011.563-3.029m5.858.908a3 3 0 114.243 4.243M9.878 9.878l4.242 4.242M9.878 9.878L8.464 8.464M9.878 9.878L8.464 8.464m5.656 5.657L22 22m-5.7-5.7L17.72 17.72" />
          </svg>
        {:else}
          <svg class="h-4 w-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 12a3 3 0 11-6 0 3 3 0 016 0z M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z" />
          </svg>
        {/if}
      </button>
      <button
        class="p-2 text-gray-400 hover:text-red-600 hover:bg-red-50 rounded-lg transition-all duration-200"
        onclick={(e) => { e.stopPropagation(); handleDelete(); }}
        title="Delete dish"
      >
        <svg class="h-4 w-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
        </svg>
      </button>
    </div>
  </div>

  {#if dish.description}
    <p class="text-gray-600 text-sm mb-4 line-clamp-2 group-hover:text-gray-700 transition-colors duration-200">
      {dish.description}
    </p>
  {/if}

  <div class="grid grid-cols-1 gap-3 text-sm text-gray-600">
    {#if dish.preparation_time}
      <div class="flex items-center space-x-2">
        <svg class="h-4 w-4 text-gray-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 8v4l3 3m6-3a9 9 0 11-18 0 9 9 0 0118 0z" />
        </svg>
        <span>{dish.preparation_time} minutes prep time</span>
      </div>
    {/if}
    {#if dish.allergens && dish.allergens.length > 0}
      <div class="flex items-center space-x-2">
        <svg class="h-4 w-4 text-yellow-500" fill="none" stroke="currentColor" viewBox="0 0 24 24">
          <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-3L13.732 4c-.77-.833-1.964-.833-2.732 0L3.732 16.5c-.77.833.192 3 1.732 3z" />
        </svg>
        <span class="truncate">
          {dish.allergens.slice(0, 3).join(', ')}
          {#if dish.allergens.length > 3}
            <span class="text-gray-400">+{dish.allergens.length - 3} more</span>
          {/if}
        </span>
      </div>
    {/if}
  </div>

  <div class="flex items-center justify-between mt-4 pt-4 border-t border-gray-100">
    <span class="text-xs text-gray-500">
      Created {new Date(dish.created_at).toLocaleDateString()}
    </span>
    <div class="opacity-0 group-hover:opacity-100 transition-opacity duration-200">
      <Button size="sm" variant="outline" onclick={(e) => { e.stopPropagation(); handleEdit(); }}>
        Edit Dish
      </Button>
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