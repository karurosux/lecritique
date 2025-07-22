/**
 * JWT decoding utilities
 * Centralizes JWT token parsing and provides type safety
 */

export interface JwtPayload {
  account_id: string;
  member_id: string;
  email: string;
  name: string;
  role: string;
  subscription_features?: {
    max_organizations: number;
    max_qr_codes: number;
    max_feedbacks_per_month: number;
    max_team_members: number;
    has_basic_analytics: boolean;
    has_advanced_analytics: boolean;
    has_feedback_explorer: boolean;
    has_custom_branding: boolean;
    has_priority_support: boolean;
  };
  exp: number;
  iat: number;
  iss: string;
  sub: string;
}

/**
 * Decode JWT token and extract payload
 */
export function decodeJwt(token: string): JwtPayload | null {
  try {
    // Split token and decode payload
    const parts = token.split('.');
    if (parts.length !== 3) {
      throw new Error('Invalid JWT format');
    }

    const payload = JSON.parse(atob(parts[1]));

    // Basic validation
    if (!payload.account_id || !payload.email) {
      throw new Error('Invalid JWT payload');
    }

    return payload as JwtPayload;
  } catch (error) {
    console.error('Failed to decode JWT:', error);
    return null;
  }
}

/**
 * Check if JWT token is expired
 */
export function isTokenExpired(token: string): boolean {
  const payload = decodeJwt(token);
  if (!payload) return true;

  return Date.now() >= payload.exp * 1000;
}

/**
 * Extract subscription features from JWT token
 */
export function getSubscriptionFeaturesFromToken(token: string): JwtPayload['subscription_features'] | null {
  const payload = decodeJwt(token);
  return payload?.subscription_features || null;
}

/**
 * Check if user has a specific feature based on JWT token
 */
export function hasFeatureFromToken(token: string, feature: string): boolean {
  const features = getSubscriptionFeaturesFromToken(token);
  if (!features) return false;

  switch (feature) {
    case 'basic_analytics':
      return features.has_basic_analytics;
    case 'advanced_analytics':
      return features.has_advanced_analytics;
    case 'feedback_explorer':
      return features.has_feedback_explorer;
    case 'custom_branding':
      return features.has_custom_branding;
    case 'priority_support':
      return features.has_priority_support;
    default:
      return false;
  }
}

/**
 * Get subscription limit from JWT token
 */
export function getLimitFromToken(token: string, limitType: string): number {
  const features = getSubscriptionFeaturesFromToken(token);
  if (!features) return 0;

  switch (limitType) {
    case 'max_organizations':
      return features.max_organizations;
    case 'max_qr_codes':
      return features.max_qr_codes;
    case 'max_feedbacks_per_month':
      return features.max_feedbacks_per_month;
    case 'max_team_members':
      return features.max_team_members;
    default:
      return 0;
  }
}
