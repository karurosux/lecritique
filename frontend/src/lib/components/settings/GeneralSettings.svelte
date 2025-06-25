<script lang="ts">
	import { Button, Input } from '$lib/components/ui';
	import { auth, type User } from '$lib/stores/auth';
	import { Loader2 } from 'lucide-svelte';

	interface Props {
		onSuccess?: (message: string) => void;
		onError?: (message: string) => void;
	}

	let { onSuccess, onError }: Props = $props();

	let loading = $state(false);
	
	// Form fields
	let companyName = $state('');
	let phone = $state('');
	
	// Subscribe to auth store
	let user = $derived($auth.user);
	
	// Initialize form values when user data is available
	$effect(() => {
		if (user) {
			companyName = user.company_name || '';
			phone = user.phone || '';
		}
	});
	
	// Check if form has changes
	let hasChanges = $derived(user && (
		companyName !== (user.company_name || '') ||
		phone !== (user.phone || '')
	));
	
	async function handleSubmit(event: Event) {
		event.preventDefault();
		if (!user || !hasChanges) return;
		
		loading = true;
		
		try {
			const api = auth.getApi();
			
			const response = await api.api.v1AuthProfileUpdate({
				company_name: companyName,
				phone: phone
			});
			
			if (response.data.success && response.data.data) {
				// Update local auth state with the updated account data
				auth.updateUser({
					...user,
					company_name: response.data.data.company_name || companyName,
					phone: response.data.data.phone || phone
				});
				
				onSuccess?.('Profile updated successfully');
			} else {
				throw new Error('Failed to update profile');
			}
		} catch (error: any) {
			console.error('Failed to update profile:', error);
			onError?.(error.message || 'Failed to update profile');
		} finally {
			loading = false;
		}
	}
	
	function handleReset() {
		if (user) {
			companyName = user.company_name || '';
			phone = user.phone || '';
		}
	}
</script>

<div>
	<div class="mb-8">
		<h2 class="text-2xl font-bold text-gray-900">General Information</h2>
		<p class="mt-1 text-sm text-gray-600">Update your personal and company information</p>
	</div>
	
	<div class="space-y-8">
		<form onsubmit={handleSubmit} class="space-y-6">
			<div class="grid grid-cols-1 md:grid-cols-2 gap-4">
				<!-- Company Information -->
				<div class="space-y-2">
					<label for="company-name" class="block text-sm font-medium text-gray-700">Company Name</label>
					<Input
						id="company-name"
						type="text"
						bind:value={companyName}
						placeholder="Enter company name"
						disabled={loading}
					/>
				</div>
				
				<div class="space-y-2">
					<label for="phone" class="block text-sm font-medium text-gray-700">Phone Number</label>
					<Input
						id="phone"
						type="tel"
						bind:value={phone}
						placeholder="Enter phone number"
						disabled={loading}
					/>
				</div>
			</div>
			
			<!-- Email (Read-only with note) -->
			<div class="space-y-2">
				<label for="email" class="block text-sm font-medium text-gray-700">Email Address</label>
				<Input
					id="email"
					type="email"
					value={user?.email || ''}
					disabled
					class="bg-gray-100"
				/>
				<p class="text-sm text-gray-500">
					To change your email address, use the Account Settings section
				</p>
			</div>
			
			<!-- Form Actions -->
			<div class="flex justify-end gap-2">
				<Button
					type="button"
					variant="outline"
					onclick={handleReset}
					disabled={loading || !hasChanges}
				>
					Reset
				</Button>
				<Button
					type="submit"
					disabled={loading || !hasChanges}
				>
					{#if loading}
						<Loader2 class="mr-2 h-4 w-4 animate-spin" />
						Saving...
					{:else}
						Save Changes
					{/if}
				</Button>
			</div>
		</form>
	</div>
</div>