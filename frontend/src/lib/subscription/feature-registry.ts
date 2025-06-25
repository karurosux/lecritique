// Feature types and definitions matching the backend
export type FeatureType = 'limit' | 'flag' | 'custom';

export interface FeatureDefinition {
  key: string;
  type: FeatureType;
  displayName: string;
  description?: string;
  unit?: string;
  unlimitedText?: string;
  format?: string;
  icon?: string;
  category?: string;
  sortOrder: number;
}

// Feature keys matching backend constants
export const LIMITS = {
  RESTAURANTS: 'max_restaurants',
  LOCATIONS_PER_RESTAURANT: 'max_locations_per_restaurant',
  QR_CODES_PER_LOCATION: 'max_qr_codes_per_location',
  FEEDBACKS_PER_MONTH: 'max_feedbacks_per_month',
  TEAM_MEMBERS: 'max_team_members',
  STORAGE_GB: 'max_storage_gb',
  API_CALLS_PER_HOUR: 'max_api_calls_per_hour'
} as const;

export const FLAGS = {
  ADVANCED_ANALYTICS: 'advanced_analytics',
  CUSTOM_BRANDING: 'custom_branding',
  API_ACCESS: 'api_access',
  PRIORITY_SUPPORT: 'priority_support',
  WHITE_LABEL: 'white_label',
  CUSTOM_DOMAIN: 'custom_domain'
} as const;

// Feature registry with all feature definitions
export const featureRegistry: Record<string, FeatureDefinition> = {
  // Limits
  [LIMITS.RESTAURANTS]: {
    key: LIMITS.RESTAURANTS,
    type: 'limit',
    displayName: 'Restaurants',
    description: 'Maximum number of restaurants',
    unit: 'restaurants',
    unlimitedText: 'Unlimited restaurants',
    format: '{value} restaurant(s)',
    icon: 'store',
    category: 'core',
    sortOrder: 1
  },
  [LIMITS.FEEDBACKS_PER_MONTH]: {
    key: LIMITS.FEEDBACKS_PER_MONTH,
    type: 'limit',
    displayName: 'Monthly Feedbacks',
    description: 'Maximum feedbacks per month',
    unit: 'feedbacks/month',
    unlimitedText: 'Unlimited feedbacks',
    format: '{value} feedbacks/month',
    icon: 'message-square',
    category: 'core',
    sortOrder: 2
  },
  [LIMITS.QR_CODES_PER_LOCATION]: {
    key: LIMITS.QR_CODES_PER_LOCATION,
    type: 'limit',
    displayName: 'QR Codes',
    description: 'QR codes per location',
    unit: 'QR codes/location',
    unlimitedText: 'Unlimited QR codes',
    format: '{value} QR codes per location',
    icon: 'qr-code',
    category: 'core',
    sortOrder: 3
  },
  [LIMITS.TEAM_MEMBERS]: {
    key: LIMITS.TEAM_MEMBERS,
    type: 'limit',
    displayName: 'Team Members',
    description: 'Maximum team members',
    unit: 'members',
    unlimitedText: 'Unlimited team members',
    format: '{value} team member(s)',
    icon: 'users',
    category: 'collaboration',
    sortOrder: 4
  },
  [LIMITS.STORAGE_GB]: {
    key: LIMITS.STORAGE_GB,
    type: 'limit',
    displayName: 'Storage',
    description: 'Storage space for media files',
    unit: 'GB',
    unlimitedText: 'Unlimited storage',
    format: '{value} GB storage',
    icon: 'hard-drive',
    category: 'resources',
    sortOrder: 5
  },
  [LIMITS.API_CALLS_PER_HOUR]: {
    key: LIMITS.API_CALLS_PER_HOUR,
    type: 'limit',
    displayName: 'API Rate Limit',
    description: 'API calls per hour',
    unit: 'calls/hour',
    unlimitedText: 'Unlimited API calls',
    format: '{value} API calls/hour',
    icon: 'activity',
    category: 'developer',
    sortOrder: 10
  },

  // Flags
  [FLAGS.ADVANCED_ANALYTICS]: {
    key: FLAGS.ADVANCED_ANALYTICS,
    type: 'flag',
    displayName: 'Advanced Analytics',
    description: 'Detailed insights and reporting',
    icon: 'bar-chart',
    category: 'analytics',
    sortOrder: 20
  },
  [FLAGS.CUSTOM_BRANDING]: {
    key: FLAGS.CUSTOM_BRANDING,
    type: 'flag',
    displayName: 'Custom Branding',
    description: 'Customize with your brand',
    icon: 'palette',
    category: 'customization',
    sortOrder: 21
  },
  [FLAGS.API_ACCESS]: {
    key: FLAGS.API_ACCESS,
    type: 'flag',
    displayName: 'API Access',
    description: 'Programmatic access via API',
    icon: 'code',
    category: 'developer',
    sortOrder: 22
  },
  [FLAGS.PRIORITY_SUPPORT]: {
    key: FLAGS.PRIORITY_SUPPORT,
    type: 'flag',
    displayName: 'Priority Support',
    description: '24/7 priority customer support',
    icon: 'headphones',
    category: 'support',
    sortOrder: 23
  },
  [FLAGS.WHITE_LABEL]: {
    key: FLAGS.WHITE_LABEL,
    type: 'flag',
    displayName: 'White Label',
    description: 'Remove LeCritique branding',
    icon: 'eye-off',
    category: 'customization',
    sortOrder: 24
  },
  [FLAGS.CUSTOM_DOMAIN]: {
    key: FLAGS.CUSTOM_DOMAIN,
    type: 'flag',
    displayName: 'Custom Domain',
    description: 'Use your own domain',
    icon: 'globe',
    category: 'customization',
    sortOrder: 25
  }
};

// Helper functions
export function getFeatureDefinition(key: string): FeatureDefinition | undefined {
  return featureRegistry[key];
}

export function getFeaturesByCategory(category: string): FeatureDefinition[] {
  return Object.values(featureRegistry)
    .filter(def => def.category === category)
    .sort((a, b) => a.sortOrder - b.sortOrder);
}

export function formatFeatureValue(key: string, value: number | boolean): string {
  const def = featureRegistry[key];
  if (!def) return '';

  if (def.type === 'limit') {
    const limitValue = value as number;
    if (limitValue === -1) {
      return def.unlimitedText || 'Unlimited';
    }
    
    let result = def.format || '{value} {unit}';
    result = result.replace('{value}', limitValue.toLocaleString());
    result = result.replace('{unit}', def.unit || '');
    return result;
  }

  if (def.type === 'flag' && value === true) {
    return def.displayName;
  }

  return '';
}

// Get all features from a plan in display order
export function getPlanFeatures(plan: any): string[] {
  const features: string[] = [];
  
  // Process limits
  if (plan.features?.limits) {
    Object.entries(plan.features.limits)
      .map(([key, value]) => ({
        key,
        value,
        def: featureRegistry[key]
      }))
      .filter(item => item.def)
      .sort((a, b) => a.def.sortOrder - b.def.sortOrder)
      .forEach(item => {
        const formatted = formatFeatureValue(item.key, item.value as number);
        if (formatted) features.push(formatted);
      });
  }
  
  // Process flags
  if (plan.features?.flags) {
    Object.entries(plan.features.flags)
      .map(([key, value]) => ({
        key,
        value,
        def: featureRegistry[key]
      }))
      .filter(item => item.def && item.value === true)
      .sort((a, b) => a.def.sortOrder - b.def.sortOrder)
      .forEach(item => {
        features.push(item.def.displayName);
      });
  }
  
  return features;
}