<script lang="ts">
  import { page } from '$app/stores';
  import { goto } from '$app/navigation';
  import { Button } from '$lib/components/ui';
  import { ArrowLeft, ClipboardList } from 'lucide-svelte';
  import QuestionnaireBuilder from '$lib/components/questionnaires/QuestionnaireBuilder.svelte';
  import { getApiClient } from '$lib/api';
  import { onMount } from 'svelte';

  let organizationId = $page.params.id;
  let productId = $page.params.productId;

  let productName = $state('');
  let loading = $state(true);
  let error = $state('');

  onMount(async () => {
    await fetchProductName();
  });

  async function fetchProductName() {
    try {
      const api = getApiClient();
      const response =
        await api.api.v1OrganizationsProductsList(organizationId);

      if (response.data.success && response.data.data) {
        const product = response.data.data.find((d: any) => d.id === productId);
        if (product) {
          productName = product.name;
        } else {
          error = 'Product not found';
        }
      }
    } catch (err) {
      console.error('Failed to load product:', err);
      error = 'Failed to load product information';
    } finally {
      loading = false;
    }
  }

  function goBack() {
    goto(`/organizations/${organizationId}/products`);
  }
</script>

<div class="space-y-6">
  <!-- Header -->
  <div class="flex items-center gap-4">
    <Button
      onclick={goBack}
      variant="ghost"
      size="sm"
      class="flex items-center gap-2">
      <ArrowLeft class="h-4 w-4" />
      Back to Products
    </Button>
    <div class="flex items-center gap-3">
      <div
        class="h-10 w-10 bg-gradient-to-br from-purple-500 to-blue-600 rounded-xl flex items-center justify-center shadow-lg shadow-purple-500/25">
        <ClipboardList class="h-5 w-5 text-white" />
      </div>
      <div>
        <h1 class="text-xl font-bold text-gray-900">
          {loading
            ? 'Loading...'
            : productName
              ? `${productName} Questionnaire`
              : 'Create Questionnaire'}
        </h1>
        <p class="text-sm text-gray-600">
          Design targeted feedback questions for this product
        </p>
      </div>
    </div>
  </div>

  <!-- Content -->
  <div class="space-y-6">
    {#if !loading}
      <QuestionnaireBuilder {organizationId} {productId} />
    {/if}
  </div>
</div>
