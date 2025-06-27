<script lang="ts">
	import type { PageData } from './$types';
	import { Button, Card } from '$lib/components/ui';
	import { Plus, QrCode, Download, Trash2, Eye, EyeOff } from 'lucide-svelte';
	// Badge component not available, will use span with styling
	import CreateQRCodeModal from '$lib/components/qr-codes/CreateQRCodeModal.svelte';
	import QRCodeDisplay from '$lib/components/qr-codes/QRCodeDisplay.svelte';
	import { getApiClient } from '$lib/api';
	import { toast } from 'svelte-sonner';
	import { invalidateAll } from '$app/navigation';
	import { ModelsQRCodeType } from '$lib/api/api';

	let { data }: { data: PageData } = $props();

	let showCreateModal = $state(false);
	let selectedQRCode = $state<typeof data.qrCodes[0] | null>(null);
	let showQRCodeDisplay = $state(false);
	
	// Get restaurant from parent layout
	let restaurant = $derived(data.restaurant);

	const qrTypeColors: Record<string, string> = {
		[ModelsQRCodeType.QRCodeTypeTable]: 'bg-blue-100 text-blue-800',
		[ModelsQRCodeType.QRCodeTypeLocation]: 'bg-green-100 text-green-800',
		[ModelsQRCodeType.QRCodeTypeTakeaway]: 'bg-yellow-100 text-yellow-800',
		[ModelsQRCodeType.QRCodeTypeDelivery]: 'bg-purple-100 text-purple-800',
		[ModelsQRCodeType.QRCodeTypeGeneral]: 'bg-gray-100 text-gray-800'
	};

	async function handleDelete(qrCode: typeof data.qrCodes[0]) {
		if (!confirm(`Are you sure you want to delete QR code "${qrCode.label}"?`)) {
			return;
		}

		try {
			const api = getApiClient();
			await api.api.v1QrCodesDelete(qrCode.id);
			toast.success('QR code deleted successfully');
			await invalidateAll();
		} catch (error) {
			toast.error('Failed to delete QR code');
			console.error(error);
		}
	}

	async function handleToggleActive(qrCode: typeof data.qrCodes[0]) {
		try {
			// TODO: Add toggle active endpoint to API
			toast.info('Toggle active functionality coming soon');
		} catch (error) {
			toast.error('Failed to update QR code status');
			console.error(error);
		}
	}

	function showQRCode(qrCode: typeof data.qrCodes[0]) {
		selectedQRCode = qrCode;
		showQRCodeDisplay = true;
	}

	function formatDate(date: string | null) {
		if (!date) return 'Never';
		return new Date(date).toLocaleDateString();
	}
</script>

<svelte:head>
	<title>QR Codes - {data.restaurant.name} | LeCritique</title>
</svelte:head>

<div class="space-y-6">
	<!-- Header -->
	<div class="flex items-center justify-between">
		<div>
			<h1 class="text-2xl font-bold">QR Codes</h1>
			<p class="text-muted-foreground">
				Manage QR codes for {restaurant.name}
			</p>
		</div>
		<Button on:click={() => (showCreateModal = true)} class="gap-2">
			<Plus class="h-4 w-4" />
			Create QR Code
		</Button>
	</div>

	<!-- QR Codes Grid -->
	{#if data.qrCodes.length === 0}
		<Card class="p-12 text-center">
			<QrCode class="mx-auto h-12 w-12 text-muted-foreground mb-4" />
			<h3 class="text-lg font-semibold mb-2">No QR codes yet</h3>
			<p class="text-muted-foreground mb-4">
				Create your first QR code to start collecting feedback
			</p>
			<Button on:click={() => (showCreateModal = true)} variant="default">
				<Plus class="mr-2 h-4 w-4" />
				Create QR Code
			</Button>
		</Card>
	{:else}
		<div class="grid gap-4 md:grid-cols-2 lg:grid-cols-3">
			{#each data.qrCodes as qrCode}
				<Card class="p-6">
					<div class="flex items-start justify-between mb-4">
						<div class="flex items-center gap-3">
							<div class="p-2 bg-primary/10 rounded-lg">
								<QrCode class="h-6 w-6 text-primary" />
							</div>
							<div>
								<h3 class="font-semibold">{qrCode.label}</h3>
								<p class="text-sm text-muted-foreground">
									{qrCode.code}
								</p>
							</div>
						</div>
						<span class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium {qrTypeColors[qrCode.type] || 'bg-gray-100 text-gray-800'}">
							{qrCode.type}
						</span>
					</div>

					<div class="space-y-2 mb-4">
						<div class="flex justify-between text-sm">
							<span class="text-muted-foreground">Status</span>
							<span class="flex items-center gap-1">
								{#if qrCode.is_active}
									<span class="h-2 w-2 bg-green-500 rounded-full"></span>
									Active
								{:else}
									<span class="h-2 w-2 bg-gray-400 rounded-full"></span>
									Inactive
								{/if}
							</span>
						</div>
						<div class="flex justify-between text-sm">
							<span class="text-muted-foreground">Scans</span>
							<span>{qrCode.scans_count || 0}</span>
						</div>
						<div class="flex justify-between text-sm">
							<span class="text-muted-foreground">Last scan</span>
							<span>{formatDate(qrCode.last_scanned_at)}</span>
						</div>
						{#if qrCode.location_id}
							<div class="flex justify-between text-sm">
								<span class="text-muted-foreground">Location</span>
								<span>Location {qrCode.location_id}</span>
							</div>
						{/if}
					</div>

					<div class="flex gap-2">
						<Button
							variant="outline"
							size="sm"
							class="flex-1"
							on:click={() => showQRCode(qrCode)}
						>
							<Eye class="h-4 w-4" />
						</Button>
						<Button
							variant="outline"
							size="sm"
							class="flex-1"
							on:click={() => handleToggleActive(qrCode)}
						>
							{#if qrCode.is_active}
								<EyeOff class="h-4 w-4" />
							{:else}
								<Eye class="h-4 w-4" />
							{/if}
						</Button>
						<Button
							variant="outline"
							size="sm"
							class="flex-1"
							on:click={() => handleDelete(qrCode)}
						>
							<Trash2 class="h-4 w-4" />
						</Button>
					</div>
				</Card>
			{/each}
		</div>
	{/if}
</div>

<!-- Create QR Code Modal -->
{#if showCreateModal}
	<CreateQRCodeModal
		restaurantId={restaurant.id}
		onclose={() => (showCreateModal = false)}
		oncreated={() => {
			showCreateModal = false;
			invalidateAll();
		}}
	/>
{/if}

<!-- QR Code Display Modal -->
{#if showQRCodeDisplay && selectedQRCode}
	<QRCodeDisplay
		qrCode={selectedQRCode}
		restaurantName={restaurant.name}
		onclose={() => {
			showQRCodeDisplay = false;
			selectedQRCode = null;
		}}
	/>
{/if}