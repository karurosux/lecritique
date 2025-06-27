<script lang="ts">
	import { Modal, Button, Input, Select } from '$lib/components/ui';
	import { getApiClient } from '$lib/api';
	import { toast } from 'svelte-sonner';
	import { ModelsQRCodeType } from '$lib/api/api';

	let {
		restaurantId,
		onclose,
		oncreated
	}: {
		restaurantId: string;
		onclose: () => void;
		oncreated: () => void;
	} = $props();

	let loading = $state(false);
	let label = $state('');
	let type = $state(ModelsQRCodeType.QRCodeTypeTable);
	let locationId = $state('');

	const qrTypes = [
		{ value: ModelsQRCodeType.QRCodeTypeTable, label: 'Table', description: 'For dine-in customers at specific tables' },
		{ value: ModelsQRCodeType.QRCodeTypeLocation, label: 'Location', description: 'For specific areas in your restaurant' },
		{ value: ModelsQRCodeType.QRCodeTypeTakeaway, label: 'Takeaway', description: 'For takeout orders' },
		{ value: ModelsQRCodeType.QRCodeTypeDelivery, label: 'Delivery', description: 'For delivery orders' },
		{ value: ModelsQRCodeType.QRCodeTypeGeneral, label: 'General', description: 'For general use' }
	];

	async function handleSubmit() {
		if (!label.trim()) {
			toast.error('Please enter a label for the QR code');
			return;
		}

		loading = true;

		try {
			const payload = {
				label: label.trim(),
				type,
				location_id: locationId || undefined
			};

			const api = getApiClient();
			await api.api.v1RestaurantsQrCodesCreate(restaurantId, payload);
			toast.success('QR code created successfully');
			oncreated();
		} catch (error) {
			toast.error('Failed to create QR code');
			console.error(error);
		} finally {
			loading = false;
		}
	}

	function handleClose() {
		onclose();
	}
</script>

<Modal 
	isOpen={true} 
	title="Create QR Code"
	size="lg"
	onclose={handleClose}
>
	<form on:submit|preventDefault={handleSubmit} class="space-y-4">
		<div class="space-y-2">
			<label for="label" class="block text-sm font-medium text-gray-700">Label</label>
			<Input
				id="label"
				bind:value={label}
				placeholder="e.g., Table 1, Main Entrance, Takeaway Counter"
				required
			/>
			<p class="text-sm text-gray-500">
				A descriptive name to identify this QR code
			</p>
		</div>

		<div class="space-y-2">
			<label for="type" class="block text-sm font-medium text-gray-700">Type</label>
			<Select 
				bind:value={type}
				options={qrTypes.map(t => ({ value: t.value, label: t.label }))}
			/>
		</div>

		{#if type === ModelsQRCodeType.QRCodeTypeLocation}
			<div class="space-y-2">
				<label for="location" class="block text-sm font-medium text-gray-700">Location ID (Optional)</label>
				<Input
					id="location"
					bind:value={locationId}
					placeholder="Enter location ID"
				/>
				<p class="text-sm text-gray-500">
					Associate this QR code with a specific location
				</p>
			</div>
		{/if}

		<div class="rounded-lg bg-gray-50 p-4">
			<h4 class="font-medium mb-2">How it works</h4>
			<ol class="space-y-1 text-sm text-gray-600">
				<li>1. Create the QR code with a descriptive label</li>
				<li>2. Download and print the QR code</li>
				<li>3. Place it at the designated location</li>
				<li>4. Customers scan to provide feedback</li>
			</ol>
		</div>

		<div class="flex justify-end gap-2">
			<Button onclick={handleClose} variant="outline">
				Cancel
			</Button>
			<Button type="submit" disabled={loading}>
				{loading ? 'Creating...' : 'Create QR Code'}
			</Button>
		</div>
	</form>
</Modal>