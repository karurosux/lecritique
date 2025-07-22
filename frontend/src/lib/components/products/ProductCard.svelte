<script lang="ts">
  import { Card, Button } from '$lib/components/ui';
  import { RoleGate } from '$lib/components/auth';
  import { Edit2, Lightbulb, Eye, EyeOff, Trash2, ClipboardList, Clock, AlertTriangle, MoreVertical } from 'lucide-svelte';

  interface Product {
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
    product,
    onedit = (product: Product, event?: MouseEvent) => {},
    ontoggleavailability = (product: Product) => {},
    ondelete = (product: Product, event?: MouseEvent) => {},
    ongeneratequestionnaire = (product: Product) => {}
  }: {
    product: Product;
    onedit?: (product: Product, event?: MouseEvent) => void;
    ontoggleavailability?: (product: Product) => void;
    ondelete?: (product: Product, event?: MouseEvent) => void;
    ongeneratequestionnaire?: (product: Product) => void;
  } = $props();

  function formatPrice(price: number): string {
    return `$${price.toFixed(2)}`;
  }

  function handleEdit(event?: MouseEvent) {
    onedit(product, event);
  }

  function handleToggleAvailability() {
    ontoggleavailability(product);
  }

  function handleDelete(event?: MouseEvent) {
    ondelete(product, event);
  }

  function handleGenerateQuestionnaire() {
    ongeneratequestionnaire(product);
  }

  let showDropdown = $state(false);

  function toggleDropdown() {
    showDropdown = !showDropdown;
  }

  function closeDropdown() {
    showDropdown = false;
  }

  // Close dropdown when clicking outside
  function handleClickOutside(event: MouseEvent) {
    if (!event.target) return;
    const target = event.target as Element;
    if (!target.closest('.dropdown-container')) {
      showDropdown = false;
    }
  }
</script>

<svelte:window on:click={handleClickOutside} />

<Card 
  variant="default" 
  hover 
  interactive 
  class="group transform transition-all duration-300 h-full {product.has_questionnaire ? 'ring-2 ring-purple-200 border-purple-300' : ''}"
>
  <div class="flex flex-col h-full">
    <div class="flex items-start justify-between mb-4">
      <div class="flex items-center space-x-3">
        <div class="w-12 h-12 bg-gradient-to-br from-green-500 to-emerald-600 rounded-xl flex items-center justify-center text-white font-bold text-lg shadow-lg shadow-green-500/25 group-hover:scale-110 transition-transform duration-200">
          {formatPrice(product.price).replace('$', '')}
        </div>
        <div class="flex-1 min-w-0 mr-2">
          <div class="flex items-center space-x-2">
            <h3 class="font-bold text-lg text-gray-900 group-hover:text-green-600 transition-colors duration-200 truncate">
              {product.name}
            </h3>
            {#if product.has_questionnaire}
              <div class="w-5 h-5 bg-purple-500 rounded-full flex items-center justify-center shadow-sm flex-shrink-0" title="Has questionnaire">
                <ClipboardList class="h-3 w-3 text-white" />
              </div>
            {/if}
          </div>
          <div class="flex items-center space-x-2 mt-1">
            <span class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium {product.is_available ? 'bg-green-100 text-green-800' : 'bg-gray-100 text-gray-800'}">
              {product.is_available ? 'Available' : 'Hidden'}
            </span>
            <span class="text-xs text-gray-500 bg-gray-100 px-2 py-1 rounded-md">
              {product.category}
            </span>
          </div>
        </div>
      </div>
      <RoleGate roles={['OWNER', 'ADMIN', 'MANAGER']}>
        <div class="relative dropdown-container flex-shrink-0">
          <button
            type="button"
            class="p-2 text-gray-600 hover:text-gray-800 hover:bg-gray-100 hover:shadow-sm hover:border hover:border-gray-200 rounded-lg transition-all duration-200 cursor-pointer {showDropdown ? 'bg-gray-100 text-gray-800 shadow-sm border border-gray-200' : ''}"
            onclick={(e) => { e.stopPropagation(); toggleDropdown(); }}
            title="More actions"
          >
            <MoreVertical class="h-4 w-4" />
          </button>
          
          {#if showDropdown}
            <div class="absolute right-0 top-full mt-1 w-48 bg-white rounded-lg shadow-lg border border-gray-200 py-1 z-50">
              <button
                class="w-full px-4 py-2 text-left text-sm text-gray-700 hover:bg-gray-50 flex items-center gap-2 cursor-pointer"
                onclick={(e) => { e.stopPropagation(); handleEdit(e); closeDropdown(); }}
              >
                <Edit2 class="h-4 w-4 text-blue-500" />
                Edit product
              </button>
              <button
                class="w-full px-4 py-2 text-left text-sm text-gray-700 hover:bg-gray-50 flex items-center gap-2 cursor-pointer"
                onclick={(e) => { e.stopPropagation(); handleGenerateQuestionnaire(); closeDropdown(); }}
                title="{product.has_questionnaire ? 'Manage questions' : 'Create questions'}"
              >
                {#if product.has_questionnaire}
                  <ClipboardList class="h-4 w-4 text-purple-500" />
                  Questions
                {:else}
                  <Lightbulb class="h-4 w-4 text-purple-500" />
                  Create questions
                {/if}
              </button>
              <button
                class="w-full px-4 py-2 text-left text-sm text-gray-700 hover:bg-gray-50 flex items-center gap-2 cursor-pointer"
                onclick={(e) => { e.stopPropagation(); handleToggleAvailability(); closeDropdown(); }}
              >
                {#if product.is_available}
                  <EyeOff class="h-4 w-4 text-gray-500" />
                  Hide product
                {:else}
                  <Eye class="h-4 w-4 text-green-500" />
                  Show product
                {/if}
              </button>
              <hr class="my-1 border-gray-200" />
              <button
                class="w-full px-4 py-2 text-left text-sm text-red-700 hover:bg-red-50 flex items-center gap-2 cursor-pointer"
                onclick={(e) => { e.stopPropagation(); handleDelete(e); closeDropdown(); }}
              >
                <Trash2 class="h-4 w-4 text-red-500" />
                Delete product
              </button>
            </div>
          {/if}
        </div>
      </RoleGate>
    </div>

  {#if product.description}
    <p class="text-gray-600 text-sm mb-4 line-clamp-2 group-hover:text-gray-700 transition-colors duration-200">
      {product.description}
    </p>
  {/if}

  <div class="grid grid-cols-1 gap-3 text-sm text-gray-600 min-h-[3rem]">
    {#if product.preparation_time}
      <div class="flex items-center space-x-2">
        <Clock class="h-4 w-4 text-gray-400" />
        <span>{product.preparation_time} minutes prep time</span>
      </div>
    {/if}
    {#if product.allergens && product.allergens.length > 0}
      <div class="flex items-center space-x-2">
        <AlertTriangle class="h-4 w-4 text-yellow-500" />
        <span class="truncate">
          {product.allergens.slice(0, 3).join(', ')}
          {#if product.allergens.length > 3}
            <span class="text-gray-400">+{product.allergens.length - 3} more</span>
          {/if}
        </span>
      </div>
    {/if}
  </div>

  <!-- Spacer to push footer to bottom -->
  <div class="flex-grow"></div>

  <div class="flex items-center justify-between mt-4 pt-4 border-t border-gray-100">
    <span class="text-xs text-gray-500">
      Created {new Date(product.created_at).toLocaleDateString()}
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
