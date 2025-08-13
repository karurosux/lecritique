<script lang="ts">
  import type { PageData } from './$types';
  import { Button, Card, QRCode } from '$lib/components/ui';
  import { Plus, QrCode, Download, Trash2, Eye, EyeOff } from 'lucide-svelte';
  import CreateQRCodeModal from '$lib/components/qr-codes/CreateQRCodeModal.svelte';
  import QRCodeDisplay from '$lib/components/qr-codes/QRCodeDisplay.svelte';
  import { RoleGate } from '$lib/components/auth';
  import { getApiClient } from '$lib/api';
  import { toast } from 'svelte-sonner';
  import { invalidateAll } from '$app/navigation';
  import { QrcodemodelQRCodeType, ContentType } from '$lib/api/api';
  import { onMount } from 'svelte';
  import { page } from '$app/stores';

  let { data }: { data: PageData } = $props();

  let showCreateModal = $state(false);
  let selectedQRCode = $state<(typeof data.qrCodes)[0] | null>(null);
  let showQRCodeDisplay = $state(false);
  let qrCodes = $state(data.qrCodes);
  let loading = $state(false);
  let clickOrigin = $state<{ x: number; y: number } | null>(null);

  let organization = $derived(data.organization);
  let organizationId = $derived($page.params.id);

  onMount(async () => {
    if (organization) {
      await fetchQRCodes();
    }
  });

  async function fetchQRCodes() {
    try {
      loading = true;
      const api = getApiClient();
      const response = await api.api.v1OrganizationsQrCodesList(organizationId);

      if (response.data.success && response.data.data) {
        qrCodes = response.data.data;
      }
    } catch (error) {
      console.error('Error loading QR codes:', error);
    } finally {
      loading = false;
    }
  }

  const qrTypeColors: Record<string, string> = {
    [QrcodemodelQRCodeType.QRCodeTypeTable]: 'bg-blue-100 text-blue-800',
    [QrcodemodelQRCodeType.QRCodeTypeLocation]: 'bg-green-100 text-green-800',
    [QrcodemodelQRCodeType.QRCodeTypeTakeaway]: 'bg-yellow-100 text-yellow-800',
    [QrcodemodelQRCodeType.QRCodeTypeDelivery]: 'bg-purple-100 text-purple-800',
    [QrcodemodelQRCodeType.QRCodeTypeGeneral]: 'bg-gray-100 text-gray-800',
  };

  async function handleDelete(qrCode: (typeof qrCodes)[0], event?: MouseEvent) {
    if (event) {
      clickOrigin = { x: event.clientX, y: event.clientY };
    }

    if (
      !confirm(`Are you sure you want to delete QR code "${qrCode.label}"?`)
    ) {
      clickOrigin = null;
      return;
    }

    try {
      const api = getApiClient();
      await api.api.v1QrCodesDelete(qrCode.id);
      toast.success('QR code deleted successfully');
      await fetchQRCodes();
    } catch (error) {
      toast.error('Failed to delete QR code');
      console.error(error);
    } finally {
      clickOrigin = null;
    }
  }

  async function handleToggleActive(qrCode: (typeof qrCodes)[0]) {
    try {
      const api = getApiClient();
      await api.request({
        path: `/api/v1/qr-codes/${qrCode.id}`,
        method: 'PATCH',
        secure: true,
        type: ContentType.Json,
        body: {
          is_active: !qrCode.is_active,
        },
        format: 'json',
      });

      qrCode.is_active = !qrCode.is_active;
      qrCodes = [...qrCodes];

      toast.success(
        `QR code ${qrCode.is_active ? 'activated' : 'deactivated'}`
      );
    } catch (error) {
      toast.error('Failed to update QR code status');
      console.error(error);
    }
  }

  function showQRCode(qrCode: (typeof qrCodes)[0], event?: MouseEvent) {
    if (event) {
      clickOrigin = { x: event.clientX, y: event.clientY };
    }
    selectedQRCode = qrCode;
    showQRCodeDisplay = true;
  }

  function formatDate(date: string | null) {
    if (!date) return 'Never';
    return new Date(date).toLocaleDateString();
  }

  function getFeedbackUrl(qrCode: (typeof qrCodes)[0]) {
    const baseUrl = 'https://kyooar.com';
    return `${baseUrl}/feedback/${qrCode.code}`;
  }
</script>

<svelte:head>
  <title>QR Codes - {organization?.name || 'Organization'} | Kyooar</title>
</svelte:head>

{#if !organization}
  <div class="space-y-6">
    <div class="text-center">
      <p class="text-gray-600">Loading organization...</p>
    </div>
  </div>
{:else}
  <div class="space-y-6">
    <!-- Header -->
    <div class="flex items-center justify-between mb-8">
      <div class="space-y-3">
        <div class="flex items-center space-x-3">
          <div
            class="h-12 w-12 bg-gradient-to-br from-blue-500 to-purple-600 rounded-2xl flex items-center justify-center shadow-lg shadow-blue-500/25">
            <QrCode class="h-6 w-6 text-white" />
          </div>
          <div>
            <h1
              class="text-3xl font-bold bg-gradient-to-r from-gray-900 to-gray-700 bg-clip-text text-transparent">
              QR Codes
            </h1>
            <p class="text-gray-600 font-medium">
              Manage QR codes for {organization.name}
            </p>
          </div>
        </div>
      </div>
      <RoleGate roles={['OWNER', 'ADMIN', 'MANAGER']}>
        <Button
          onclick={e => {
            clickOrigin = { x: e.clientX, y: e.clientY };
            showCreateModal = true;
          }}
          variant="gradient"
          size="lg"
          class="gap-2">
          <Plus class="h-4 w-4" />
          Create QR Code
        </Button>
      </RoleGate>
    </div>

    <!-- Loading State -->
    {#if loading}
      <div class="text-center">
        <p class="text-gray-600">Loading QR codes...</p>
      </div>
      <!-- QR Codes Grid -->
    {:else if qrCodes.length === 0}
      <Card variant="glass" class="p-12 text-center">
        <div
          class="w-24 h-24 mx-auto bg-gradient-to-br from-blue-500 to-purple-600 rounded-full flex items-center justify-center mb-6 shadow-lg shadow-blue-500/25">
          <QrCode class="h-12 w-12 text-white" />
        </div>
        <h3
          class="text-xl font-bold mb-2 bg-gradient-to-r from-gray-900 to-gray-700 bg-clip-text text-transparent">
          No QR codes yet
        </h3>
        <p class="text-gray-600 mb-6 max-w-md mx-auto">
          Create your first QR code to start collecting feedback from customers
        </p>
        <RoleGate roles={['OWNER', 'ADMIN', 'MANAGER']}>
          <Button
            onclick={e => {
              clickOrigin = { x: e.clientX, y: e.clientY };
              showCreateModal = true;
            }}
            variant="gradient"
            size="lg">
            <Plus class="mr-2 h-4 w-4" />
            Create First QR Code
          </Button>
        </RoleGate>
      </Card>
    {:else}
      <div class="grid gap-6 md:grid-cols-2 lg:grid-cols-3">
        {#each qrCodes as qrCode}
          <Card
            variant="glass"
            class="p-6 hover:shadow-xl transition-all duration-300">
            <!-- Header with QR Info -->
            <div class="flex items-start justify-between mb-4">
              <div>
                <h3 class="text-lg font-bold text-gray-900">{qrCode.label}</h3>
                <p class="text-sm text-gray-600 font-mono mt-1">
                  {qrCode.code}
                </p>
              </div>
              {#if qrCode.location}
                <span
                  class="inline-flex items-center px-3 py-1 rounded-full text-xs font-semibold bg-blue-100 text-blue-800">
                  {qrCode.location}
                </span>
              {/if}
            </div>

            <!-- Stats -->
            <div class="space-y-2 mb-4">
              <div class="flex justify-between items-center">
                <span class="text-sm text-gray-600">Status</span>
                <span class="flex items-center gap-2">
                  {#if qrCode.is_active}
                    <span class="h-2 w-2 bg-green-500 rounded-full"></span>
                    <span class="text-sm font-semibold text-green-700"
                      >Active</span>
                  {:else}
                    <span class="h-2 w-2 bg-gray-400 rounded-full"></span>
                    <span class="text-sm font-semibold text-gray-700"
                      >Inactive</span>
                  {/if}
                </span>
              </div>
              <div class="flex justify-between items-center">
                <span class="text-sm text-gray-600">Total Scans</span>
                <span class="text-sm font-bold text-gray-900"
                  >{qrCode.scans_count || 0}</span>
              </div>
              <div class="flex justify-between items-center">
                <span class="text-sm text-gray-600">Last Scan</span>
                <span class="text-sm font-semibold text-gray-900"
                  >{formatDate(qrCode.last_scanned_at)}</span>
              </div>
            </div>

            <!-- Actions -->
            <div class="flex gap-2">
              <Button
                variant="gradient"
                size="sm"
                class="flex-1"
                onclick={e => showQRCode(qrCode, e)}>
                <Download class="h-4 w-4 mr-1" />
                Download
              </Button>
              <RoleGate roles={['OWNER', 'ADMIN', 'MANAGER']}>
                <Button
                  variant="outline"
                  size="sm"
                  onclick={() => handleToggleActive(qrCode)}>
                  {#if qrCode.is_active}
                    <EyeOff class="h-4 w-4" />
                  {:else}
                    <Eye class="h-4 w-4" />
                  {/if}
                </Button>
              </RoleGate>
              <RoleGate roles={['OWNER', 'ADMIN', 'MANAGER']}>
                <Button
                  variant="outline"
                  size="sm"
                  onclick={e => handleDelete(qrCode, e)}>
                  <Trash2 class="h-4 w-4" />
                </Button>
              </RoleGate>
            </div>
          </Card>
        {/each}
      </div>
    {/if}
  </div>
{/if}

<!-- Create QR Code Modal -->
{#if showCreateModal && organization}
  <CreateQRCodeModal
    organizationId={organization.id}
    {clickOrigin}
    onclose={() => {
      showCreateModal = false;
      clickOrigin = null;
    }}
    oncreated={() => {
      showCreateModal = false;
      clickOrigin = null;
      fetchQRCodes();
    }} />
{/if}

<!-- QR Code Display Modal -->
{#if showQRCodeDisplay && selectedQRCode && organization}
  <QRCodeDisplay
    qrCode={selectedQRCode}
    organizationName={organization.name}
    {clickOrigin}
    onclose={() => {
      showQRCodeDisplay = false;
      selectedQRCode = null;
      clickOrigin = null;
    }} />
{/if}
