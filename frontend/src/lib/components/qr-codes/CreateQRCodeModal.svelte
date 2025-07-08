<script lang="ts">
	import { Modal, Button, Input, Select } from '$lib/components/ui';
	import { getApiClient } from '$lib/api';
	import { toast } from 'svelte-sonner';
	import { ModelsQRCodeType } from '$lib/api/api';
	import { Loader2 } from 'lucide-svelte';

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
	let location = $state('');

	async function handleSubmit() {
		if (!label.trim()) {
			toast.error('Please enter a label for the QR code');
			return;
		}

		loading = true;

		try {
			const payload = {
				label: label.trim(),
				restaurant_id: restaurantId,
				type: "general" as const, // Default to general type
				location: location.trim() || undefined
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
			<label for="location" class="block text-sm font-medium text-gray-700">Location (Optional)</label>
			<Input
				id="location"
				bind:value={location}
				placeholder="e.g., Table 5, Main Bar, Patio, Front Counter"
			/>
			<p class="text-sm text-gray-500">
				Describe where this QR code will be placed
			</p>
		</div>

		<div class="rounded-lg bg-gray-50 p-4">
			<h4 class="font-medium mb-2">How it works</h4>
			<ol class="space-y-1 text-sm text-gray-600">
				<li>1. Create the QR code with a descriptive label</li>
				<li>2. Download and print the QR code</li>
				<li>3. Place it at the designated location</li>
				<li>4. Customers scan to provide feedback</li>
			</ol>
		</div>

		<div class="mt-6 pt-6 border-t border-gray-200 flex justify-end space-x-3">
			<Button onclick={handleClose} variant="outline">
				Cancel
			</Button>
			<Button type="submit" disabled={loading} variant="gradient">
				{#if loading}
					<Loader2 class="w-4 h-4 mr-2 animate-spin" />
					Creating...
				{:else}
					Create QR Code
				{/if}
			</Button>
		</div>
	</form>
</Modal>