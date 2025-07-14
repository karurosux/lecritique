<script lang="ts">
	import { planFeatures } from '$lib/stores/subscription';

	let planName = $derived($planFeatures?.plan_name || 'Free');
	let planCode = $derived($planFeatures?.plan_code || 'free');
	let isActive = $derived($planFeatures?.is_active || false);

	let badgeClass = $derived(() => {
		if (!isActive) return 'bg-gray-100 text-gray-800';
		
		switch (planCode) {
			case 'premium':
			case 'enterprise':
				return 'bg-purple-100 text-purple-800';
			case 'professional':
				return 'bg-blue-100 text-blue-800';
			case 'starter':
				return 'bg-green-100 text-green-800';
			default:
				return 'bg-gray-100 text-gray-800';
		}
	});
</script>

<span class="inline-flex items-center px-2.5 py-0.5 rounded-full text-xs font-medium {badgeClass()}">
	{planName}
</span>