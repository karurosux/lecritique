import { redirect } from '@sveltejs/kit';
import { get } from 'svelte/store';
import { subscription, hasFeature, getLimit, FEATURES, LIMITS } from '$lib/stores/subscription';

export type FeatureFlag = typeof FEATURES[keyof typeof FEATURES];
export type UsageLimit = typeof LIMITS[keyof typeof LIMITS];

export interface RouteProtectionConfig {
	requireSubscription?: boolean;
	requireFeature?: FeatureFlag;
	requireLimit?: {
		limit: UsageLimit;
		redirectOnExceeded?: boolean;
	};
	redirectTo?: string;
}

export function checkRouteAccess(config: RouteProtectionConfig): boolean {
	const sub = get(subscription);
	const checkFeature = get(hasFeature);
	const checkLimit = get(getLimit);
	
	// Check if subscription is required
	if (config.requireSubscription && sub.subscription?.status !== 'active') {
		return false;
	}
	
	// Check if specific feature is required
	if (config.requireFeature && !checkFeature(config.requireFeature)) {
		return false;
	}
	
	// Check if limit allows access
	if (config.requireLimit) {
		const limit = checkLimit(config.requireLimit.limit);
		
		// Map limit keys to usage keys
		const usageMap: Record<string, keyof SubscriptionUsage> = {
			[LIMITS.RESTAURANTS]: 'restaurants_count',
			[LIMITS.QR_CODES]: 'qr_codes_count',
			[LIMITS.FEEDBACKS_PER_MONTH]: 'feedbacks_count',
			[LIMITS.TEAM_MEMBERS]: 'team_members_count'
		};
		
		const usageKey = usageMap[config.requireLimit.limit];
		const currentUsage = sub.usage?.[usageKey] || 0;
		
		if (limit !== -1 && currentUsage >= limit) {
			return false;
		}
	}
	
	return true;
}

export function protectRoute(config: RouteProtectionConfig) {
	const hasAccess = checkRouteAccess(config);
	
	if (!hasAccess) {
		const redirectTo = config.redirectTo || '/subscription';
		throw redirect(303, redirectTo);
	}
}

// Convenience functions for common checks
export function requireFeature(feature: FeatureFlag, redirectTo?: string) {
	protectRoute({ requireFeature: feature, redirectTo });
}

export function requireSubscription(redirectTo?: string) {
	protectRoute({ requireSubscription: true, redirectTo });
}

export function requireLimit(limit: UsageLimit, redirectOnExceeded = true, redirectTo?: string) {
	protectRoute({ 
		requireLimit: { limit, redirectOnExceeded },
		redirectTo 
	});
}

// Type for subscription usage (should match backend)
interface SubscriptionUsage {
	feedbacks_count: number;
	restaurants_count: number;
	locations_count: number;
	qr_codes_count: number;
	team_members_count: number;
	period_start: string;
	period_end: string;
}