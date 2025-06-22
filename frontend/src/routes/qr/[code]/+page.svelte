<script lang="ts">
  import { page } from '$app/stores';
  import { onMount } from 'svelte';
  import { goto } from '$app/navigation';
  import { Card, Button } from '$lib/components/ui';
  import { getApiClient, handleApiError } from '$lib/api/client';

  interface QRValidationData {
    valid: boolean;
    restaurant?: {
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

  let loading = true;
  let error = '';
  let qrData: QRValidationData | null = null;
  
  $: code = $page.params.code;
  $: pageTitle = qrData?.restaurant?.name 
    ? `${qrData.restaurant.name} - LeCritique`
    : 'QR Code Validation - LeCritique';

  onMount(async () => {
    await validateQRCode();
  });

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
      
      if (response.data.success && response.data.data) {
        qrData = response.data.data as QRValidationData;
        
        if (!qrData.valid) {
          error = 'This QR code is invalid or has expired';
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
    if (qrData?.restaurant?.id) {
      goto(`/restaurant/${qrData.restaurant.id}/menu?qr=${code}`);
    }
  }

  function handleGiveFeedback() {
    if (qrData?.restaurant?.id) {
      goto(`/feedback?restaurant=${qrData.restaurant.id}&qr=${code}`);
    }
  }
</script>

<svelte:head>
  <title>{pageTitle}</title>
  <meta name="description" content="Restaurant feedback and menu access via QR code" />
  <meta name="robots" content="noindex, nofollow" />
</svelte:head>

<div class="min-h-screen bg-gray-50 py-8 px-4">
  <div class="max-w-2xl mx-auto">
    {#if loading}
      <!-- Loading State -->
      <Card>
        <div class="text-center py-12">
          <svg class="animate-spin h-8 w-8 text-blue-600 mx-auto mb-4" fill="none" viewBox="0 0 24 24">
            <circle class="opacity-25" cx="12" cy="12" r="10" stroke="currentColor" stroke-width="4"></circle>
            <path class="opacity-75" fill="currentColor" d="M4 12a8 8 0 018-8V0C5.373 0 0 5.373 0 12h4zm2 5.291A7.962 7.962 0 714 12H0c0 3.042 1.135 5.824 3 7.938l3-2.647z"></path>
          </svg>
          <p class="text-gray-600">Validating QR code...</p>
        </div>
      </Card>
    
    {:else if error}
      <!-- Error State -->
      <Card>
        <div class="text-center py-12">
          <svg class="h-12 w-12 text-red-500 mx-auto mb-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
            <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-2.5L13.732 4c-.77-.833-1.964-.833-2.732 0L3.732 16.5c-.77.833.192 2.5 1.732 2.5z" />
          </svg>
          <h2 class="text-xl font-semibold text-gray-900 mb-2">Invalid QR Code</h2>
          <p class="text-gray-600 mb-6">{error}</p>
          <div class="space-y-2">
            <p class="text-sm text-gray-500">Possible reasons:</p>
            <ul class="text-sm text-gray-500 space-y-1">
              <li>• The QR code has expired</li>
              <li>• The QR code is no longer active</li>
              <li>• The link may have been typed incorrectly</li>
            </ul>
          </div>
        </div>
      </Card>
    
    {:else if qrData && qrData.valid}
      <!-- Success State -->
      <div class="space-y-6">
        <!-- Restaurant Header -->
        <Card>
          <div class="text-center py-8">
            {#if qrData.restaurant?.logo}
              <img 
                src={qrData.restaurant.logo} 
                alt="{qrData.restaurant.name} logo"
                class="h-16 w-16 rounded-full mx-auto mb-4 object-cover"
              />
            {:else}
              <div class="h-16 w-16 rounded-full bg-blue-100 flex items-center justify-center mx-auto mb-4">
                <svg class="h-8 w-8 text-blue-600" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M19 21V5a2 2 0 00-2-2H7a2 2 0 00-2 2v16m14 0h2m-2 0h-4m-5 0H9m0 0H5m5 0v-4a1 1 0 011-1h2a1 1 0 011 1v4M7 7h3m3 0h3m-6 4h3m3 0h3" />
                </svg>
              </div>
            {/if}
            
            <h1 class="text-2xl font-bold text-gray-900 mb-2">
              Welcome to {qrData.restaurant?.name || 'Our Restaurant'}
            </h1>
            
            {#if qrData.restaurant?.description}
              <p class="text-gray-600 mb-4">{qrData.restaurant.description}</p>
            {/if}
            
            {#if qrData.location}
              <div class="flex items-center justify-center text-gray-500 text-sm">
                <svg class="h-4 w-4 mr-1" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M17.657 16.657L13.414 20.9a1.998 1.998 0 01-2.827 0l-4.244-4.243a8 8 0 1111.314 0z" />
                  <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M15 11a3 3 0 11-6 0 3 3 0 016 0z" />
                </svg>
                <span>{qrData.location.name}</span>
                {#if qrData.location.address}
                  <span class="ml-2">• {qrData.location.address}</span>
                {/if}
              </div>
            {/if}
            
            {#if qrData.qr_code?.label}
              <div class="mt-2 inline-flex items-center px-3 py-1 rounded-full text-xs font-medium bg-blue-100 text-blue-800">
                {qrData.qr_code.label}
              </div>
            {/if}
          </div>
        </Card>

        <!-- Action Buttons -->
        <div class="grid grid-cols-1 md:grid-cols-2 gap-4">
          <Card hover>
            <div class="text-center py-6">
              <svg class="h-12 w-12 text-green-600 mx-auto mb-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
              </svg>
              <h3 class="text-lg font-medium text-gray-900 mb-2">View Menu</h3>
              <p class="text-gray-600 text-sm mb-4">Browse our delicious dishes and current offerings</p>
              <Button variant="primary" on:click={handleViewMenu} class="w-full">
                View Menu
              </Button>
            </div>
          </Card>

          <Card hover>
            <div class="text-center py-6">
              <svg class="h-12 w-12 text-blue-600 mx-auto mb-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path stroke-linecap="round" stroke-linejoin="round" stroke-width="2" d="M8 12h.01M12 12h.01M16 12h.01M21 12c0 4.418-4.03 8-9 8a9.863 9.863 0 01-4.255-.949L3 20l1.395-3.72C3.512 15.042 3 13.574 3 12c0-4.418 4.03-8 9-8s9 3.582 9 8z" />
              </svg>
              <h3 class="text-lg font-medium text-gray-900 mb-2">Give Feedback</h3>
              <p class="text-gray-600 text-sm mb-4">Share your dining experience and help us improve</p>
              <Button variant="secondary" on:click={handleGiveFeedback} class="w-full">
                Give Feedback
              </Button>
            </div>
          </Card>
        </div>

        <!-- QR Code Info -->
        <Card>
          <div class="text-center py-4">
            <p class="text-xs text-gray-500">
              QR Code: {qrData.qr_code?.code || code}
              {#if qrData.qr_code?.type}
                • Type: {qrData.qr_code.type}
              {/if}
            </p>
          </div>
        </Card>
      </div>
    {/if}
  </div>
</div>