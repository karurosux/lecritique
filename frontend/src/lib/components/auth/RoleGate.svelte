<script lang="ts">
	import { auth } from '$lib/stores/auth';
	import type { ModelsTeamMember } from '$lib/api/api';
	
	interface Props {
		// Allowed roles - if any match, content is shown
		roles?: Array<'OWNER' | 'ADMIN' | 'MANAGER' | 'VIEWER'>;
		// Team members array to check against
		teamMembers?: ModelsTeamMember[];
		// If true, only shows content if user owns the account (not a team member)
		requireOwner?: boolean;
		// If true, shows loading state while checking
		showLoading?: boolean;
		// Custom fallback content
		fallback?: string;
		children?: any;
	}
	
	let { 
		roles = [], 
		teamMembers = [],
		requireOwner = false,
		showLoading = false,
		fallback = '',
		children 
	}: Props = $props();
	
	// Get current user's role from team members
	let currentUserRole = $derived.by(() => {
		const userEmail = $auth.user?.email;
		const userId = $auth.user?.id;
		
		if (!userEmail || !userId || teamMembers.length === 0) return null;
		
		// Find user's role in team members
		const member = teamMembers.find((m: ModelsTeamMember) => {
			// For owner, check account_id match
			if (m.role === 'OWNER' && m.account_id === userId) {
				return true;
			}
			// For other roles, check email
			return m.member?.email === userEmail;
		});
		
		return member?.role || null;
	});
	
	// Check if user is the account owner (not accessing via team membership)
	let isAccountOwner = $derived(
		$auth.user?.id === $auth.user?.account_id
	);
	
	// Determine if content should be shown
	let hasAccess = $derived.by(() => {
		// If requireOwner is true, user must be the account owner
		if (requireOwner && !isAccountOwner) {
			return false;
		}
		
		// If no roles specified, allow access
		if (roles.length === 0) {
			return true;
		}
		
		// Check if user has any of the allowed roles
		return currentUserRole && roles.includes(currentUserRole as any);
	});
	
	let isLoading = $derived(
		showLoading && teamMembers.length === 0 && $auth.isAuthenticated
	);
</script>

{#if isLoading}
	<div class="flex items-center justify-center p-4">
		<div class="animate-spin rounded-full h-6 w-6 border-b-2 border-primary-600"></div>
	</div>
{:else if hasAccess}
	{@render children?.()}
{:else if fallback}
	<div class="text-gray-500 text-sm">
		{fallback}
	</div>
{/if}