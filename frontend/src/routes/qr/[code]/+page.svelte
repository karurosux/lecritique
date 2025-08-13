<script lang="ts">
  import { page } from '$app/stores';
  import { goto } from '$app/navigation';
  import { Card, Button, Input, Rating } from '$lib/components/ui';
  import { getApiClient, getPublicApiClient, handleApiError } from '$lib/api/client';
  import { Api } from '$lib/api/api';
  import {
    Loader2,
    MapPin,
    Building2,
    CheckCircle,
    UtensilsCrossed,
    ChevronRight,
    Lightbulb,
    Menu,
    MessageSquare,
    ArrowLeft,
    X,
    AlertTriangle,
  } from 'lucide-svelte';
  import Logo from '$lib/components/ui/Logo.svelte';

  interface QRValidationData {
    valid: boolean;
    organization?: {
      id: string;
      name: string;
      description?: string;
      logo?: string;
    };
    location?: {
      id: string;
      name: string;
      address?: string;
    };
    qr_code?: {
      id: string;
      code: string;
      label: string;
      type: string;
    };
  }

  interface Product {
    id: string;
    name: string;
    description?: string;
    price: number;
    category: string;
    image_url?: string;
    question_count?: number;
  }

  interface QuestionOption {
    id: string;
    text: string;
    value?: string;
  }

  interface Question {
    id: string;
    text: string;
    type:
      | 'rating'
      | 'scale'
      | 'single_choice'
      | 'multiple_choice'
      | 'yes_no'
      | 'text';
    is_required: boolean;
    options?: QuestionOption[];
    min_value?: number;
    max_value?: number;
    min_label?: string;
    max_label?: string;
    display_order: number;
  }

  let loading = $state(true);
  let error = $state('');
  let qrData = $state<QRValidationData | null>(null);
  let productsWithQuestions = $state<Product[]>([]);
  let loadingProducts = $state(false);

  let selectedProduct = $state<Product | null>(null);
  let questions = $state<Question[]>([]);
  let loadingQuestions = $state(false);
  let submitting = $state(false);

  let overallRating = $state(0);
  let responses = $state<Record<string, any>>({});
  let comment = $state('');
  let customerEmail = $state('');

  const code = $derived($page.params.code);
  const pageTitle = $derived(
    qrData?.organization?.name
      ? `${qrData.organization.name} - Kyooar`
      : 'QR Code Validation - Kyooar'
  );

  $effect(() => {
    validateQRCode();
  });

  async function loadProductsWithQuestions() {
    if (!qrData?.organization?.id) {
      return;
    }

    try {
      loadingProducts = true;
      const publicApi = getPublicApiClient();
      const response =
        await publicApi.api.v1PublicOrganizationQuestionsProductsWithQuestionsList(
          qrData.organization.id
        );

      if (response.data.success && response.data.data) {
        productsWithQuestions = response.data.data;
      }
    } catch (error) {
    } finally {
      loadingProducts = false;
    }
  }

  async function validateQRCode() {
    if (!code) {
      error = 'No QR code provided';
      loading = false;
      return;
    }

    try {
      loading = true;
      error = '';

      const api = getApiClient();
      const response = await api.api.v1PublicQrDetail(code);

      if (response.data && response.data.success && response.data.data) {
        const qrCodeData = response.data.data;

        qrData = {
          valid: qrCodeData.is_active === true,
          organization: qrCodeData.organization
            ? {
                id: qrCodeData.organization.id,
                name: qrCodeData.organization.name,
                description: qrCodeData.organization.description,
                logo: qrCodeData.organization.logo,
              }
            : undefined,
          location: qrCodeData.location
            ? {
                id: qrCodeData.organization?.id || '',
                name: qrCodeData.location,
                address: qrCodeData.organization?.address,
              }
            : undefined,
          qr_code: {
            id: qrCodeData.id,
            code: qrCodeData.code,
            label: qrCodeData.label,
            type: qrCodeData.type,
          },
        };

        if (qrData && !qrData.valid) {
          error = 'This QR code is invalid or has expired';
        } else if (qrData && qrData.valid) {
          await loadProductsWithQuestions();
        }
      } else {
        error = 'Invalid QR code';
      }
    } catch (err) {
      error = handleApiError(err);
    } finally {
      loading = false;
    }
  }

  function handleViewMenu() {
    if (qrData?.organization?.id) {
      goto(`/organization/${qrData.organization.id}/menu?qr=${code}`);
    }
  }

  function handleGiveFeedback() {
    if (qrData?.organization?.id) {
      let url = `/feedback?organization=${qrData.organization.id}&qr=${code}`;
      if (qrData?.location?.id) {
        url += `&location=${qrData.location.id}`;
      }
      goto(url);
    }
  }

  async function handleProductFeedback(product: Product) {
    selectedProduct = product;
    await loadProductQuestions(product.id);
  }

  async function loadProductQuestions(productId: string) {
    if (!qrData?.organization?.id) {
      error = 'Organization information is missing';
      return;
    }

    try {
      loadingQuestions = true;
      error = '';

      const publicApi = getPublicApiClient();

      const response =
        await publicApi.api.v1PublicOrganizationProductsQuestionsList(
          qrData.organization.id,
          productId
        );

      if (response.data.success && response.data.data) {
        questions = response.data.data.questions || [];

        questions.forEach(question => {
          if (question.type === 'rating' || question.type === 'scale') {
            responses[question.id] = question.min_value || 0;
          } else if (question.type === 'multi_choice') {
            responses[question.id] = [];
          }
        });
      }
    } catch (err) {
      error = handleApiError(err);
    } finally {
      loadingQuestions = false;
    }
  }

  function handleBackToProducts() {
    selectedProduct = null;
    questions = [];
    responses = {};
    overallRating = 0;
    comment = '';
    customerEmail = '';
    error = '';
  }

  function handleQuestionResponse(questionId: string, value: any) {
    responses[questionId] = value;
  }

  function validateForm(): boolean {
    for (const question of questions) {
      if (question.is_required) {
        const response = responses[question.id];

        if (question.type === 'multiple_choice') {
          if (!response || response.length === 0) {
            error = `Please answer: ${question.text}`;
            return false;
          }
        } else if (question.type === 'rating' || question.type === 'scale') {
          if (!response || response === 0) {
            error = `Please answer: ${question.text}`;
            return false;
          }
        } else if (!response && response !== false) {
          error = `Please answer: ${question.text}`;
          return false;
        }
      }
    }

    return true;
  }

  async function handleSubmit(event: SubmitEvent) {
    event.preventDefault();

    if (!event.currentTarget?.checkValidity()) {
      return;
    }

    submitting = true;

    try {
      const api = getApiClient();

      const responseArray = Object.entries(responses).map(
        ([questionId, answer]) => {
          const question = questions.find(q => q.id === questionId);
          return {
            question_id: questionId,
            question_text: question?.text || '',
            question_type: question?.type || 'choice',
            answer: answer,
          };
        }
      );

      const feedbackData = {
        organization_id: qrData?.organization?.id,
        product_id: selectedProduct?.id,
        qr_code_id: qrData?.qr_code?.id,
        overall_rating: overallRating || 0,
        responses: responseArray,
        customer_email: customerEmail,
      };

      if (comment && comment.trim()) {
        feedbackData.responses.push({
          question_id: crypto.randomUUID(),
          question_text: 'Additional Comments',
          question_type: 'text',
          answer: comment,
        });
      }

      await api.api.v1PublicFeedbackCreate(feedbackData as any);

      goto('/feedback/success');
    } catch (err) {
      error = handleApiError(err);
      submitting = false;
    }
  }

  function formatPrice(price: number): string {
    return `$${price.toFixed(2)}`;
  }
</script>

<svelte:head>
  <title>{pageTitle}</title>
  <meta
    name="description"
    content="Organization feedback and menu access via QR code" />
  <meta name="robots" content="noindex, nofollow" />
  <meta
    name="viewport"
    content="width=device-width, initial-scale=1.0, maximum-scale=1.0, user-scalable=no" />
  <meta name="theme-color" content="#3b82f6" />
</svelte:head>

<div
  class="min-h-screen bg-gradient-to-br from-blue-100 via-purple-50 to-pink-100 py-2 px-3 sm:py-4 sm:px-4">
  <div class="max-w-lg mx-auto qr-page">
    <!-- Kyooar Logo -->
    <div class="flex justify-center pt-2 pb-4">
      <Logo size="md" class="opacity-80" />
    </div>
    {#if loading}
      <!-- Loading State -->
      <Card>
        <div class="text-center py-8 sm:py-12">
          <Loader2 class="animate-spin h-8 w-8 text-blue-600 mx-auto mb-4" />
          <p class="text-gray-600">Validating QR code...</p>
        </div>
      </Card>
    {:else if error}
      <!-- Error State -->
      <Card>
        <div class="text-center py-8 sm:py-12 px-4">
          <AlertTriangle class="h-10 w-10 sm:h-12 sm:w-12 text-red-500 mx-auto mb-4" />
          <h2 class="text-lg sm:text-xl font-semibold text-gray-900 mb-2">
            Invalid QR Code
          </h2>
          <p class="text-gray-600 mb-4 sm:mb-6 text-sm sm:text-base">{error}</p>
          <div class="space-y-2">
            <p class="text-sm text-gray-500">Possible reasons:</p>
            <ul class="text-sm text-gray-500 space-y-1 text-left inline-block">
              <li>• The QR code has expired</li>
              <li>• The QR code is no longer active</li>
              <li>• The link may have been typed incorrectly</li>
            </ul>
          </div>
        </div>
      </Card>
    {:else if qrData && qrData.valid && !selectedProduct}
      <!-- Product Selection State -->
      <div class="space-y-2">
        <!-- Organization Header -->
        <div class="py-2 sm:py-3 px-2">
          <Card variant="glass" padding={false} class="p-3">
            {#if qrData.organization?.logo}
              <div
                class="h-10 w-10 sm:h-12 sm:w-12 rounded-2xl bg-gradient-to-br from-blue-100 to-purple-100 flex items-center justify-center shadow-sm flex-shrink-0 ring-2 ring-white/50">
                <img
                  src={qrData.organization.logo}
                  alt="{qrData.organization.name} logo"
                  class="h-9 w-9 sm:h-11 sm:w-11 rounded-2xl object-cover" />
              </div>
            {:else}
              <div
                class="h-10 w-10 sm:h-12 sm:w-12 rounded-2xl bg-gradient-to-br from-blue-500 to-purple-600 flex items-center justify-center shadow-sm flex-shrink-0 ring-2 ring-white/50">
                <Building2 class="h-5 w-5 sm:h-6 sm:w-6 text-white" />
              </div>
            {/if}

            <div class="flex-1 min-w-0 space-y-1">
              <div class="flex items-start justify-between">
                <h1
                  class="text-base sm:text-lg font-bold bg-gradient-to-r from-gray-900 to-gray-700 bg-clip-text text-transparent leading-tight">
                  {qrData.organization?.name || 'Our Organization'}
                </h1>
                {#if qrData.qr_code?.label}
                  <div
                    class="inline-flex items-center px-2 py-0.5 rounded-full text-xs font-medium bg-gradient-to-r from-blue-100 to-purple-100 text-blue-800 ml-2 flex-shrink-0 border border-blue-200/50">
                    {qrData.qr_code.label}
                  </div>
                {/if}
              </div>

              {#if qrData.organization?.description}
                <p class="text-gray-600 text-xs leading-relaxed line-clamp-2">
                  {qrData.organization.description}
                </p>
              {/if}

              {#if qrData.location}
                <div class="text-xs text-gray-500">
                  <div class="flex items-center">
                    <MapPin class="h-3 w-3 mr-1 flex-shrink-0" />
                    <span class="truncate">{qrData.location.name}</span>
                    {#if qrData.location.address}
                      <span class="ml-1 text-gray-400 hidden sm:inline"
                        >• {qrData.location.address}</span>
                    {/if}
                  </div>
                  {#if qrData.location.address}
                    <div class="text-gray-400 sm:hidden mt-0.5 truncate">
                      {qrData.location.address}
                    </div>
                  {/if}
                </div>
              {/if}
            </div>
          </Card>
        </div>

        <!-- Main Question: What did you eat? -->
        <div class="p-2 sm:p-3">
          <div class="text-center mb-3 sm:mb-4">
            <div
              class="inline-flex items-center justify-center w-10 h-10 sm:w-12 sm:h-12 rounded-2xl bg-gradient-to-br from-blue-500 via-purple-500 to-purple-600 text-white mb-3 sm:mb-4 shadow-lg shadow-blue-500/25 ring-2 ring-white/20">
              <CheckCircle class="h-5 w-5 sm:h-6 sm:w-6" />
            </div>
            <h2
              class="text-xl sm:text-2xl font-bold bg-gradient-to-r from-gray-900 via-gray-800 to-gray-700 bg-clip-text text-transparent mb-2 px-2">
              What would you like to give feedback on?
            </h2>
            <p class="text-gray-600 text-sm px-2">
              Select the item you'd like to review
            </p>
          </div>

          {#if loadingProducts}
            <div class="text-center py-8">
              <Loader2
                class="animate-spin h-8 w-8 text-blue-600 mx-auto mb-4" />
              <p class="text-gray-600">Loading products...</p>
            </div>
          {:else if productsWithQuestions.length > 0}
            <div class="space-y-3">
              {#each productsWithQuestions as product}
                <Card
                  variant="default"
                  hover
                  interactive
                  padding={false}
                  class="w-full p-4 sm:p-5 rounded-3xl focus:outline-none focus:ring-2 focus:ring-purple-500/30 active:scale-[0.98]"
                  onclick={() => handleProductFeedback(product)}>
                  <div class="flex items-start justify-between">
                    <div class="flex-1 min-w-0">
                      <div class="flex items-start gap-4">
                        <div
                          class="w-12 h-12 sm:w-14 sm:h-14 rounded-2xl bg-gradient-to-br from-purple-400 to-pink-400 flex items-center justify-center flex-shrink-0 shadow-lg group-hover:shadow-xl transition-shadow">
                          <UtensilsCrossed
                            class="h-6 w-6 sm:h-7 sm:w-7 text-white" />
                        </div>
                        <div class="min-w-0 flex-1">
                          <h3
                            class="font-bold text-gray-900 text-lg sm:text-xl mb-1 group-hover:text-purple-700 transition-colors">
                            {product.name}
                          </h3>
                          <div class="flex items-center gap-3 text-sm mb-2">
                            <span class="font-bold text-purple-600 text-base"
                              >{formatPrice(product.price)}</span>
                            <span
                              class="px-2 py-0.5 rounded-full bg-purple-100 text-purple-700 text-xs font-medium capitalize"
                              >{product.category}</span>
                          </div>
                          {#if product.description}
                            <p
                              class="text-sm text-gray-600 leading-relaxed line-clamp-2">
                              {product.description}
                            </p>
                          {/if}
                        </div>
                      </div>
                    </div>
                    <div
                      class="ml-4 flex-shrink-0 opacity-50 group-hover:opacity-100 transition-opacity">
                      <ChevronRight
                        class="h-5 w-5 text-gray-400 group-hover:text-purple-600 transition-colors" />
                    </div>
                  </div>
                </Card>
              {/each}
            </div>
          {:else}
            <div class="text-center py-8">
              <div
                class="w-16 h-16 rounded-full bg-gradient-to-br from-gray-100 to-gray-200 flex items-center justify-center mx-auto mb-4">
                <Lightbulb class="h-8 w-8 text-gray-400" />
              </div>
              <p class="text-gray-600 mb-2 font-medium">
                No products available for feedback
              </p>
              <p class="text-sm text-gray-500">
                Please check back later or contact the organization
              </p>
            </div>
          {/if}
        </div>
      </div>
    {:else if qrData && qrData.valid && selectedProduct}
      <!-- Questionnaire Form State -->
      <div class="space-y-4">
        <!-- Back Button and Header -->
        <div class="flex items-center gap-3 mb-4">
          <Button
            variant="ghost"
            size="sm"
            onclick={handleBackToProducts}
            class="flex items-center gap-2 hover:bg-white/50">
            <ArrowLeft class="h-4 w-4" />
            Back
          </Button>
        </div>

        <!-- Organization & Product Header -->
        <div class="bg-white rounded-3xl p-5 shadow-sm mb-6">
          <div class="flex items-center gap-4">
            {#if qrData.organization?.logo}
              <img
                src={qrData.organization.logo}
                alt="{qrData.organization.name} logo"
                class="h-14 w-14 rounded-2xl object-cover shadow-md" />
            {:else}
              <div
                class="h-14 w-14 rounded-2xl bg-gradient-to-br from-purple-500 to-pink-500 flex items-center justify-center shadow-md">
                <Building2 class="h-7 w-7 text-white" />
              </div>
            {/if}
            <div class="flex-1">
              <h2 class="font-bold text-gray-900 text-lg">
                {qrData.organization?.name}
              </h2>
              <div class="flex items-center gap-2 mt-1">
                <div class="h-1.5 w-1.5 rounded-full bg-purple-500"></div>
                <p class="text-purple-700 font-medium">
                  {selectedProduct.name}
                </p>
                <span class="text-purple-600 font-bold"
                  >{formatPrice(selectedProduct.price)}</span>
              </div>
            </div>
          </div>
        </div>

        {#if loadingQuestions}
          <!-- Loading Questions State -->
          <div class="bg-white rounded-3xl p-8 shadow-sm">
            <div class="text-center">
              <div
                class="inline-flex items-center justify-center w-16 h-16 rounded-full bg-gradient-to-br from-purple-100 to-pink-100 mb-4">
                <Loader2 class="animate-spin h-8 w-8 text-purple-600" />
              </div>
              <p class="text-gray-600 font-medium">Loading questions...</p>
            </div>
          </div>
        {:else}
          <!-- Feedback Form -->
          <form onsubmit={handleSubmit} class="space-y-4">
            <!-- Dynamic Questions -->
            {#each [...questions].sort((a, b) => a.display_order - b.display_order) as question, index}
              <div
                class="bg-white rounded-3xl p-6 shadow-sm hover:shadow-md transition-shadow duration-300 relative overflow-hidden group question-{question.type}">
                <div
                  class="absolute top-0 left-0 w-1 h-full bg-gradient-to-b from-purple-500 to-pink-500 opacity-0 group-hover:opacity-100 transition-opacity duration-300">
                </div>
                <div class="space-y-4">
                  <div class="flex items-start gap-3">
                    <div
                      class="w-8 h-8 rounded-full bg-gradient-to-br from-purple-100 to-pink-100 flex items-center justify-center flex-shrink-0 mt-0.5">
                      <span class="text-sm font-bold text-purple-700"
                        >{index + 1}</span>
                    </div>
                    <label
                      class="block text-lg font-semibold text-gray-900 flex-1">
                      {question.text}
                      {#if question.is_required}
                        <span class="text-pink-500 ml-1">*</span>
                      {/if}
                    </label>
                  </div>

                  {#if question.type === 'rating'}
                    <input
                      type="hidden"
                      name="rating_{question.id}"
                      bind:value={responses[question.id]}
                      required={question.is_required}
                      min={question.is_required ? 1 : 0}
                      max={question.max_value || 5} />
                    <Rating
                      bind:value={responses[question.id]}
                      max={question.max_value || 5}
                      size="md" />
                  {:else if question.type === 'scale'}
                    <div class="px-3">
                      <input
                        type="range"
                        name="scale_{question.id}"
                        min={question.min_value || 1}
                        max={question.max_value || 10}
                        bind:value={responses[question.id]}
                        required={question.is_required}
                        class="w-full h-3 bg-gradient-to-r from-purple-200 to-pink-200 rounded-lg appearance-none cursor-pointer slider"
                        style="background: linear-gradient(to right, rgb(168 85 247) 0%, rgb(168 85 247) {((responses[
                          question.id
                        ] -
                          (question.min_value || 1)) /
                          ((question.max_value || 10) -
                            (question.min_value || 1))) *
                          100}%, rgb(226 232 240) {((responses[question.id] -
                          (question.min_value || 1)) /
                          ((question.max_value || 10) -
                            (question.min_value || 1))) *
                          100}%, rgb(226 232 240) 100%)" />
                      <div class="flex justify-between text-sm mt-3">
                        <span class="text-gray-500 font-medium"
                          >{question.min_label ||
                            question.min_value ||
                            1}</span>
                        <span
                          class="font-bold text-purple-600 text-base bg-purple-50 px-3 py-1 rounded-full"
                          >{responses[question.id] || 0}</span>
                        <span class="text-gray-500 font-medium"
                          >{question.max_label ||
                            question.max_value ||
                            10}</span>
                      </div>
                    </div>
                  {:else if question.type === 'single_choice' && question.options}
                    <div class="space-y-3 px-3">
                      {#each question.options as option}
                        <label
                          class="flex items-center cursor-pointer p-4 rounded-2xl border-2 bg-white hover:bg-gray-50 transition-all duration-200
                                     {responses[question.id] === option
                            ? 'border-purple-500 bg-purple-50'
                            : 'border-gray-200 hover:border-purple-300'}">
                          <input
                            type="radio"
                            name="question_{question.id}"
                            value={option}
                            checked={responses[question.id] === option}
                            onchange={() =>
                              handleQuestionResponse(question.id, option)}
                            required={question.is_required}
                            class="h-5 w-5 text-purple-600 focus:ring-purple-500 mr-4" />
                          <span class="text-gray-800 font-medium flex-1"
                            >{option}</span>
                          {#if responses[question.id] === option}
                            <CheckCircle class="h-5 w-5 text-purple-600 ml-2" />
                          {/if}
                        </label>
                      {/each}
                    </div>
                  {:else if question.type === 'multi_choice' && question.options}
                    <div class="space-y-3 px-3">
                      {#each question.options as option}
                        {@const isSelected =
                          responses[question.id]?.includes(option)}
                        <label
                          class="flex items-center cursor-pointer p-4 rounded-2xl border-2 transition-all duration-200
                                     {isSelected
                            ? 'border-purple-500 bg-purple-50'
                            : 'border-gray-200 bg-white hover:border-purple-300 hover:bg-gray-50'}">
                          <input
                            type="checkbox"
                            name="multi_{question.id}"
                            value={option}
                            checked={isSelected}
                            onchange={e => {
                              const target = e.target as HTMLInputElement;
                              const value = option;
                              if (target.checked) {
                                responses[question.id] = [
                                  ...(responses[question.id] || []),
                                  value,
                                ];
                              } else {
                                responses[question.id] = responses[
                                  question.id
                                ].filter((v: any) => v !== value);
                              }
                            }}
                            required={question.is_required &&
                              (!responses[question.id] ||
                                responses[question.id].length === 0)}
                            class="h-5 w-5 text-purple-600 rounded focus:ring-purple-500 mr-4" />
                          <span class="text-gray-800 font-medium flex-1"
                            >{option}</span>
                          {#if isSelected}
                            <CheckCircle class="h-5 w-5 text-purple-600 ml-2" />
                          {/if}
                        </label>
                      {/each}
                    </div>
                  {:else if question.type === 'yes_no'}
                    <div class="flex gap-4 px-3">
                      <input
                        type="hidden"
                        name="yesno_{question.id}"
                        bind:value={responses[question.id]}
                        required={question.is_required} />
                      <button
                        type="button"
                        onclick={() =>
                          handleQuestionResponse(question.id, true)}
                        class="flex-1 py-4 px-6 rounded-2xl font-medium transition-all duration-200 {responses[
                          question.id
                        ] === true
                          ? 'bg-gradient-to-r from-green-500 to-emerald-500 text-white shadow-lg shadow-green-500/25'
                          : 'bg-gray-100 text-gray-700 hover:bg-gray-200'}">
                        <div class="flex items-center justify-center gap-2">
                          {#if responses[question.id] === true}
                            <CheckCircle class="h-5 w-5" />
                          {/if}
                          Yes
                        </div>
                      </button>
                      <button
                        type="button"
                        onclick={() =>
                          handleQuestionResponse(question.id, false)}
                        class="flex-1 py-4 px-6 rounded-2xl font-medium transition-all duration-200 {responses[
                          question.id
                        ] === false
                          ? 'bg-gradient-to-r from-red-500 to-pink-500 text-white shadow-lg shadow-red-500/25'
                          : 'bg-gray-100 text-gray-700 hover:bg-gray-200'}">
                        <div class="flex items-center justify-center gap-2">
                          {#if responses[question.id] === false}
                            <X class="h-5 w-5" />
                          {/if}
                          No
                        </div>
                      </button>
                    </div>
                  {:else if question.type === 'text'}
                    <div class="px-3">
                      <textarea
                        name="text_{question.id}"
                        class="w-full px-5 py-4 border-2 border-gray-200 rounded-2xl focus:outline-none focus:border-purple-500 focus:ring-4 focus:ring-purple-100 resize-none transition-all duration-200 font-medium text-gray-800 placeholder-gray-400"
                        rows="4"
                        placeholder="Share your thoughts..."
                        maxlength="500"
                        required={question.is_required}
                        bind:value={responses[question.id]}></textarea>
                      <div
                        class="flex justify-between items-center mt-2 text-xs text-gray-500">
                        <span>Maximum 500 characters</span>
                        <span
                          class="{(responses[question.id]?.length || 0) > 450
                            ? 'text-orange-500'
                            : ''} {(responses[question.id]?.length || 0) > 480
                            ? 'text-red-500'
                            : ''}">
                          {responses[question.id]?.length || 0}/500
                        </span>
                      </div>
                    </div>
                  {/if}
                </div>
              </div>
            {/each}

            <!-- Additional Comments -->
            <div
              class="bg-white rounded-3xl p-6 shadow-sm hover:shadow-md transition-shadow duration-300">
              <div class="flex items-center gap-3 mb-4">
                <div
                  class="w-8 h-8 rounded-full bg-gradient-to-br from-blue-100 to-purple-100 flex items-center justify-center">
                  <MessageSquare class="h-4 w-4 text-blue-600" />
                </div>
                <label
                  for="comment"
                  class="block text-lg font-semibold text-gray-900">
                  Additional Comments
                </label>
              </div>
              <textarea
                id="comment"
                rows="4"
                class="w-full px-5 py-4 border-2 border-gray-200 rounded-2xl focus:outline-none focus:border-purple-500 focus:ring-4 focus:ring-purple-100 resize-none transition-all duration-200 font-medium text-gray-800 placeholder-gray-400"
                placeholder="Tell us more about your experience..."
                maxlength="500"
                bind:value={comment}></textarea>
              <div
                class="flex justify-between items-center mt-2 text-xs text-gray-500">
                <span>Maximum 500 characters</span>
                <span
                  class="{(comment?.length || 0) > 450
                    ? 'text-orange-500'
                    : ''} {(comment?.length || 0) > 480 ? 'text-red-500' : ''}">
                  {comment?.length || 0}/500
                </span>
              </div>
            </div>

            <!-- Optional Email -->
            <div
              class="bg-white rounded-3xl p-6 shadow-sm hover:shadow-md transition-shadow duration-300">
              <div class="flex items-center gap-3 mb-2">
                <div
                  class="w-8 h-8 rounded-full bg-gradient-to-br from-green-100 to-emerald-100 flex items-center justify-center">
                  <Menu class="h-4 w-4 text-green-600" />
                </div>
                <label
                  for="email"
                  class="block text-lg font-semibold text-gray-900">
                  Email (Optional)
                </label>
              </div>
              <p class="text-sm text-gray-600 mb-4 ml-11">
                Leave your email to receive exclusive coupons and rewards for
                your feedback
              </p>
              <input
                id="email"
                type="email"
                placeholder="your@email.com"
                bind:value={customerEmail}
                class="w-full px-5 py-4 border-2 border-gray-200 rounded-2xl focus:outline-none focus:border-purple-500 focus:ring-4 focus:ring-purple-100 transition-all duration-200 font-medium text-gray-800 placeholder-gray-400" />
            </div>

            <!-- Error Display -->
            {#if error}
              <div
                class="bg-gradient-to-r from-red-50 to-pink-50 border-2 border-red-200 rounded-2xl p-5">
                <div class="flex items-center gap-3">
                  <div
                    class="w-8 h-8 rounded-full bg-red-100 flex items-center justify-center flex-shrink-0">
                    <X class="h-4 w-4 text-red-600" />
                  </div>
                  <p class="text-red-800 font-medium">{error}</p>
                </div>
              </div>
            {/if}

            <!-- Submit Button -->
            <button
              type="submit"
              disabled={submitting}
              class="w-full py-5 px-8 bg-gradient-to-r from-purple-600 via-purple-500 to-pink-500 text-white font-bold text-lg rounded-3xl shadow-xl shadow-purple-500/30 hover:shadow-2xl hover:shadow-purple-500/40 transform hover:scale-[1.02] active:scale-[0.98] transition-all duration-300 disabled:opacity-50 disabled:cursor-not-allowed disabled:transform-none">
              {#if submitting}
                <div class="flex items-center justify-center gap-3">
                  <Loader2 class="animate-spin h-6 w-6" />
                  <span>Submitting your feedback...</span>
                </div>
              {:else}
                <div class="flex items-center justify-center gap-3">
                  <CheckCircle class="h-6 w-6" />
                  <span>Submit Feedback</span>
                </div>
              {/if}
            </button>
          </form>
        {/if}
      </div>
    {/if}
  </div>
</div>
