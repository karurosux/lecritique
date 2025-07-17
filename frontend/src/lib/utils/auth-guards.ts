import { redirect } from '@sveltejs/kit';
import type { ModelsTeamMember } from '$lib/api/api';

export type Role = 'OWNER' | 'ADMIN' | 'MANAGER' | 'VIEWER';

interface GuardOptions {
	roles?: Role[];
	requireOwner?: boolean;
	requireVerified?: boolean;
	redirectTo?: string;
}

/**
 * Check if user has any of the specified roles
 */
export function hasRole(teamMembers: ModelsTeamMember[], userEmail: string, userId: string, allowedRoles: Role[]): boolean {
	const member = teamMembers.find((m: ModelsTeamMember) => {
		// For owner, check account_id match
		if (m.role && m.role.toString() === 'OWNER' && m.account_id === userId) {
			return true;
		}
		// For other roles, check email
		return m.member?.email === userEmail;
	});
	
	return member?.role ? allowedRoles.includes(member.role.toString() as Role) : false;
}

/**
 * Route guard for load functions
 * Usage in +page.ts or +layout.ts:
 * 
 * export async function load(event) {
 *   await requireAuth(event, { roles: ['OWNER', 'ADMIN'] });
 *   // ... rest of load function
 * }
 */
export async function requireAuth(event: any, options: GuardOptions = {}) {
	const { parent } = event;
	const parentData = await parent();
	
	// Check if user is authenticated
	if (!parentData.auth?.isAuthenticated) {
		throw redirect(303, options.redirectTo || '/login');
	}
	
	const user = parentData.auth.user;
	const subscription = parentData.subscription;
	
	// Check email verification if required
	if (options.requireVerified && !user?.email_verified) {
		throw redirect(303, '/email-verification');
	}
	
	// Check owner requirement
	if (options.requireOwner && user?.id !== user?.account_id) {
		throw redirect(303, '/dashboard');
	}
	
	// Check role requirements
	if (options.roles && options.roles.length > 0) {
		const teamMembers = subscription?.team_members || [];
		const hasRequiredRole = hasRole(teamMembers, user?.email || '', user?.id || '', options.roles);
		
		if (!hasRequiredRole) {
			throw redirect(303, '/dashboard');
		}
	}
	
	return {
		user,
		subscription,
		hasRole: (roles: Role[]) => hasRole(subscription?.team_members || [], user?.email || '', user?.id || '', roles)
	};
}

/**
 * Check if user can perform an action
 * Useful for API endpoints or server-side logic
 */
export function canPerformAction(
	teamMembers: ModelsTeamMember[], 
	userEmail: string, 
	userId: string, 
	action: string
): boolean {
	const rolePermissions: Record<Role, string[]> = {
		OWNER: ['*'], // Can do everything
		ADMIN: [
			'manage_team',
			'manage_restaurants',
			'view_analytics',
			'manage_feedback',
			'manage_qr_codes'
		],
		MANAGER: [
			'manage_restaurants',
			'view_analytics',
			'manage_feedback',
			'manage_qr_codes'
		],
		VIEWER: [
			'view_restaurants',
			'view_analytics',
			'view_feedback'
		]
	};
	
	const member = teamMembers.find((m: ModelsTeamMember) => {
		if (m.role && m.role.toString() === 'OWNER' && m.account_id === userId) {
			return true;
		}
		return m.member?.email === userEmail;
	});
	
	if (!member?.role) return false;
	
	const permissions = rolePermissions[member.role.toString() as Role];
	return permissions.includes('*') || permissions.includes(action);
}