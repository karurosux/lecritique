<script lang="ts">
  import { Card, Button } from '$lib/components/ui';
  import { RoleGate } from '$lib/components/auth';
  import {
    Edit2,
    Lightbulb,
    Eye,
    EyeOff,
    Trash2,
    ClipboardList,
    Clock,
    AlertTriangle,
  } from 'lucide-svelte';

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
    index = 0,
    onedit = (product: Product, event?: MouseEvent) => {},
    ontoggleavailability = (product: Product) => {},
    ondelete = (product: Product, event?: MouseEvent) => {},
    ongeneratequestionnaire = (product: Product) => {},
  }: {
    product: Product;
    index?: number;
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

  function formatDate(dateString: string): string {
    return new Date(dateString).toLocaleDateString();
  }

  function getStatusColor(isAvailable: boolean): string {
    return isAvailable
      ? 'bg-green-100 text-green-800'
      : 'bg-gray-100 text-gray-800';
  }

  function getInitials(name: string): string {
    return name
      .split(' ')
      .map(word => word[0])
      .join('')
      .toUpperCase()
      .slice(0, 2);
  }
</script>

<Card
  variant="gradient"
  hover
  class="group relative transform transition-all duration-300 animate-fade-in-up !pb-3"
  style="animation-delay: {index * 100}ms">
  <!-- Header Section -->
  <div class="flex items-center space-x-4 mb-4">
    <div
      class="h-16 w-16 bg-gradient-to-br from-green-500 to-emerald-600 rounded-2xl flex items-center justify-center shadow-lg shadow-green-500/25 flex-shrink-0">
      <span class="text-white font-bold text-lg">
        {formatPrice(product.price)}
      </span>
    </div>
    <div class="flex-1 min-w-0">
      <div class="flex items-center space-x-2 mb-1">
        <p class="text-sm font-semibold text-gray-600 uppercase tracking-wide">
          Product
        </p>
        {#if product.has_questionnaire}
          <div
            class="w-5 h-5 bg-purple-500 rounded-full flex items-center justify-center shadow-sm flex-shrink-0"
            title="Has questionnaire">
            <ClipboardList class="h-3 w-3 text-white" />
          </div>
        {/if}
      </div>
      <p
        class="text-xl font-bold bg-gradient-to-r from-green-600 to-emerald-600 bg-clip-text text-transparent truncate max-w-[200px] mb-1"
        title={product.name}>
        {product.name}
      </p>
      <div class="flex items-center space-x-1">
        <span
          class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium {getStatusColor(
            product.is_available
          )}">
          {product.is_available ? 'Available' : 'Hidden'}
        </span>
        <span class="text-xs text-gray-500 bg-gray-100 px-2 py-1 rounded-md">
          {product.category}
        </span>
      </div>
    </div>
  </div>

  <!-- Description Section -->
  <div class="mb-4 min-h-[2.5rem]">
    {#if product.description}
      <p class="text-gray-600 text-sm line-clamp-2 leading-relaxed">
        {product.description}
      </p>
    {/if}
  </div>

  <!-- Product Information Section -->
  {#if product.preparation_time || (product.allergens && product.allergens.length > 0)}
    <div class="space-y-2 mb-4">
      {#if product.preparation_time}
        <div class="flex items-center space-x-3">
          <div
            class="h-8 w-8 bg-blue-100 rounded-lg flex items-center justify-center">
            <Clock class="h-4 w-4 text-blue-600" />
          </div>
          <span class="text-sm text-gray-700 font-medium"
            >{product.preparation_time} minutes prep time</span>
        </div>
      {/if}
      {#if product.allergens && product.allergens.length > 0}
        <div class="flex items-center space-x-3">
          <div
            class="h-8 w-8 bg-yellow-100 rounded-lg flex items-center justify-center">
            <AlertTriangle class="h-4 w-4 text-yellow-600" />
          </div>
          <span class="text-sm text-gray-700 font-medium">
            {product.allergens.slice(0, 3).join(', ')}
            {#if product.allergens.length > 3}
              <span class="text-gray-400">
                +{product.allergens.length - 3} more</span>
            {/if}
          </span>
        </div>
      {/if}
    </div>
  {/if}

  <!-- Footer Section -->
  <div class="flex items-center justify-between pt-4 border-t border-gray-200">
    <span class="text-xs text-gray-500">
      Created {formatDate(product.created_at)}
    </span>
    <!-- Action Buttons -->
    <RoleGate roles={['OWNER', 'ADMIN', 'MANAGER']}>
      <div
        class="flex items-center space-x-1 opacity-0 group-hover:opacity-100 transition-opacity duration-200">
        <button
          class="p-1.5 text-gray-400 hover:text-blue-600 hover:bg-blue-50 rounded-lg transition-all duration-200"
          onclick={e => {
            e.stopPropagation();
            handleEdit(e);
          }}
          aria-label="Edit product">
          <Edit2 class="h-3.5 w-3.5" />
        </button>
        <button
          class="p-1.5 text-gray-400 hover:text-green-600 hover:bg-green-50 rounded-lg transition-all duration-200"
          onclick={e => {
            e.stopPropagation();
            handleToggleAvailability();
          }}
          aria-label="Toggle product availability">
          {#if product.is_available}
            <EyeOff class="h-3.5 w-3.5" />
          {:else}
            <Eye class="h-3.5 w-3.5" />
          {/if}
        </button>
        <button
          class="p-1.5 text-gray-400 hover:text-purple-600 hover:bg-purple-50 rounded-lg transition-all duration-200"
          onclick={e => {
            e.stopPropagation();
            handleGenerateQuestionnaire();
          }}
          aria-label={product.has_questionnaire
            ? 'Manage questions'
            : 'Create questions'}>
          {#if product.has_questionnaire}
            <ClipboardList class="h-3.5 w-3.5" />
          {:else}
            <Lightbulb class="h-3.5 w-3.5" />
          {/if}
        </button>
        <button
          class="p-1.5 text-gray-400 hover:text-red-600 hover:bg-red-50 rounded-lg transition-all duration-200"
          onclick={e => {
            e.stopPropagation();
            handleDelete(e);
          }}
          aria-label="Delete product">
          <Trash2 class="h-3.5 w-3.5" />
        </button>
      </div>
    </RoleGate>
  </div>
</Card>

<style>
  @keyframes fade-in-up {
    from {
      opacity: 0;
      transform: translateY(10px);
    }
    to {
      opacity: 1;
      transform: translateY(0);
    }
  }

  .animate-fade-in-up {
    animation: fade-in-up 0.6s ease-out forwards;
    opacity: 0;
  }

  .line-clamp-1 {
    display: -webkit-box;
    -webkit-line-clamp: 1;
    -webkit-box-orient: vertical;
    overflow: hidden;
  }

  .line-clamp-2 {
    display: -webkit-box;
    -webkit-line-clamp: 2;
    -webkit-box-orient: vertical;
    overflow: hidden;
  }
</style>
